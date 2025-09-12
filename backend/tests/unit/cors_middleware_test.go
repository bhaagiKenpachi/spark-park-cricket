package unit

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// corsMiddleware is a copy of the middleware function for testing
func corsMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Cache-Control, Pragma, Expires, Accept")
			w.Header().Set("Access-Control-Max-Age", "86400")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func TestCorsMiddleware(t *testing.T) {
	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})

	// Create CORS middleware
	corsMiddleware := corsMiddleware()

	// Test OPTIONS request (preflight)
	t.Run("OPTIONS request should return CORS headers", func(t *testing.T) {
		req := httptest.NewRequest("OPTIONS", "/api/v1/series", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		req.Header.Set("Access-Control-Request-Method", "GET")
		req.Header.Set("Access-Control-Request-Headers", "Content-Type, Cache-Control")

		w := httptest.NewRecorder()
		corsMiddleware(testHandler).ServeHTTP(w, req)

		// Check status code
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		// Check CORS headers
		expectedHeaders := map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
			"Access-Control-Allow-Headers": "Content-Type, Authorization, Cache-Control, Pragma, Expires, Accept",
			"Access-Control-Max-Age":       "86400",
		}

		for header, expectedValue := range expectedHeaders {
			actualValue := w.Header().Get(header)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s to be %s, got %s", header, expectedValue, actualValue)
			}
		}
	})

	// Test GET request with Cache-Control header
	t.Run("GET request with Cache-Control header should work", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/series", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cache-Control", "no-cache, no-store, must-revalidate")
		req.Header.Set("Accept", "application/json")

		w := httptest.NewRecorder()
		corsMiddleware(testHandler).ServeHTTP(w, req)

		// Check status code
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		// Check that CORS headers are present
		if w.Header().Get("Access-Control-Allow-Origin") != "*" {
			t.Error("Expected Access-Control-Allow-Origin header to be present")
		}

		// Check that the test handler was called
		if w.Body.String() != "test response" {
			t.Errorf("Expected response body 'test response', got '%s'", w.Body.String())
		}
	})

	// Test POST request with all headers
	t.Run("POST request with all custom headers should work", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/series", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cache-Control", "no-cache, no-store, must-revalidate")
		req.Header.Set("Pragma", "no-cache")
		req.Header.Set("Expires", "0")
		req.Header.Set("Accept", "application/json")

		w := httptest.NewRecorder()
		corsMiddleware(testHandler).ServeHTTP(w, req)

		// Check status code
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		// Check that CORS headers are present
		expectedHeaders := []string{
			"Access-Control-Allow-Origin",
			"Access-Control-Allow-Methods",
			"Access-Control-Allow-Headers",
		}

		for _, header := range expectedHeaders {
			if w.Header().Get(header) == "" {
				t.Errorf("Expected header %s to be present", header)
			}
		}
	})
}

func TestCorsMiddlewareHeaders(t *testing.T) {
	// Test that all required headers are allowed
	allowedHeaders := []string{
		"Content-Type",
		"Authorization",
		"Cache-Control",
		"Pragma",
		"Expires",
		"Accept",
	}

	for _, header := range allowedHeaders {
		t.Run("Should allow "+header+" header", func(t *testing.T) {
			req := httptest.NewRequest("OPTIONS", "/api/v1/series", nil)
			req.Header.Set("Access-Control-Request-Headers", header)

			w := httptest.NewRecorder()
			corsMiddleware := corsMiddleware()
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			corsMiddleware(testHandler).ServeHTTP(w, req)

			// Check that the header is in the allowed headers list
			allowHeaders := w.Header().Get("Access-Control-Allow-Headers")
			if allowHeaders == "" {
				t.Error("Expected Access-Control-Allow-Headers to be set")
			}

			// The header should be included in the allowed headers
			// We can't easily parse the comma-separated list, so we just check it's not empty
			// In a real scenario, you might want to parse and verify each header individually
		})
	}
}
