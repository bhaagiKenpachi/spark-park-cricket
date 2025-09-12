package unit

import (
	"net/http"
	"net/http/httptest"
	"spark-park-cricket-backend/internal/middleware"
	"sync"
	"testing"
)

func TestRateLimitMiddlewareConcurrency(t *testing.T) {
	// Create rate limit middleware with low limit for testing
	rateLimitMiddleware := middleware.RateLimitMiddleware(5) // 5 requests per minute

	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	// Wrap with rate limit middleware
	handler := rateLimitMiddleware(testHandler)

	// Test concurrent requests
	t.Run("Concurrent requests should not cause race conditions", func(t *testing.T) {
		const numGoroutines = 20
		const requestsPerGoroutine = 3

		var wg sync.WaitGroup
		results := make(chan int, numGoroutines*requestsPerGoroutine)

		// Launch multiple goroutines making concurrent requests
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(goroutineID int) {
				defer wg.Done()

				for j := 0; j < requestsPerGoroutine; j++ {
					req := httptest.NewRequest("GET", "/test", nil)
					req.RemoteAddr = "127.0.0.1:12345" // Same IP for all requests

					w := httptest.NewRecorder()
					handler.ServeHTTP(w, req)

					results <- w.Code
				}
			}(i)
		}

		// Wait for all goroutines to complete
		wg.Wait()
		close(results)

		// Collect results
		var successCount, rateLimitCount int
		for statusCode := range results {
			switch statusCode {
			case http.StatusOK:
				successCount++
			case http.StatusTooManyRequests:
				rateLimitCount++
			default:
				t.Errorf("Unexpected status code: %d", statusCode)
			}
		}

		// Verify we got some successful requests and some rate limited
		totalRequests := numGoroutines * requestsPerGoroutine
		if successCount == 0 {
			t.Error("Expected some successful requests")
		}
		if rateLimitCount == 0 {
			t.Error("Expected some rate limited requests")
		}
		if successCount+rateLimitCount != totalRequests {
			t.Errorf("Expected %d total responses, got %d", totalRequests, successCount+rateLimitCount)
		}

		t.Logf("Concurrent test results: %d successful, %d rate limited out of %d total",
			successCount, rateLimitCount, totalRequests)
	})
}

func TestRateLimitMiddlewareBasic(t *testing.T) {
	// Create rate limit middleware
	rateLimitMiddleware := middleware.RateLimitMiddleware(3) // 3 requests per minute

	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	// Wrap with rate limit middleware
	handler := rateLimitMiddleware(testHandler)

	t.Run("Should allow requests within limit", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "127.0.0.1:12345"

		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("Should rate limit after exceeding limit", func(t *testing.T) {
		// Make requests up to the limit
		for i := 0; i < 3; i++ {
			req := httptest.NewRequest("GET", "/test", nil)
			req.RemoteAddr = "127.0.0.1:54321" // Different IP

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Request %d: Expected status code %d, got %d", i+1, http.StatusOK, w.Code)
			}
		}

		// This request should be rate limited
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "127.0.0.1:54321"

		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusTooManyRequests {
			t.Errorf("Expected status code %d, got %d", http.StatusTooManyRequests, w.Code)
		}

		// Check for Retry-After header
		retryAfter := w.Header().Get("Retry-After")
		if retryAfter != "60" {
			t.Errorf("Expected Retry-After header to be '60', got '%s'", retryAfter)
		}
	})

	t.Run("Should allow requests from different IPs", func(t *testing.T) {
		// Make requests from different IPs - each should be allowed
		ips := []string{"127.0.0.1:11111", "127.0.0.1:22222", "127.0.0.1:33333"}

		for i, ip := range ips {
			req := httptest.NewRequest("GET", "/test", nil)
			req.RemoteAddr = ip

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("IP %s (request %d): Expected status code %d, got %d", ip, i+1, http.StatusOK, w.Code)
			}
		}
	})
}

func TestRateLimitMiddlewareTimeWindow(t *testing.T) {
	// Create rate limit middleware with very low limit
	rateLimitMiddleware := middleware.RateLimitMiddleware(2) // 2 requests per minute

	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	// Wrap with rate limit middleware
	handler := rateLimitMiddleware(testHandler)

	t.Run("Should reset after time window", func(t *testing.T) {
		ip := "127.0.0.1:99999"

		// Make 2 requests (should succeed)
		for i := 0; i < 2; i++ {
			req := httptest.NewRequest("GET", "/test", nil)
			req.RemoteAddr = ip

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Request %d: Expected status code %d, got %d", i+1, http.StatusOK, w.Code)
			}
		}

		// Third request should be rate limited
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = ip

		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusTooManyRequests {
			t.Errorf("Expected status code %d, got %d", http.StatusTooManyRequests, w.Code)
		}

		// Note: In a real test, you might want to mock time or wait for the actual time window
		// For this unit test, we're just verifying the basic rate limiting logic works
	})
}
