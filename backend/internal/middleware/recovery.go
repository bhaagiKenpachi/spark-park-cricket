package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"spark-park-cricket-backend/internal/utils"
	"sync"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

// RecoveryMiddleware provides panic recovery with structured logging
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic
				utils.LogError(fmt.Errorf("panic: %v", err), "Panic recovered", map[string]interface{}{
					"method":      r.Method,
					"path":        r.URL.Path,
					"remote_addr": r.RemoteAddr,
					"user_agent":  r.UserAgent(),
					"stack":       string(debug.Stack()),
				})

				// Write error response
				utils.WriteInternalError(w, "Internal server error")
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// TimeoutMiddleware provides request timeout handling
func TimeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
	return middleware.Timeout(timeout)
}

// RateLimitMiddleware provides basic rate limiting
func RateLimitMiddleware(requestsPerMinute int) func(http.Handler) http.Handler {
	// Simple in-memory rate limiter with mutex protection
	clients := make(map[string][]time.Time)
	var mutex sync.RWMutex

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientIP := r.RemoteAddr
			now := time.Now()

			// Lock for write operations
			mutex.Lock()

			// Clean old requests
			if requests, exists := clients[clientIP]; exists {
				var validRequests []time.Time
				for _, reqTime := range requests {
					if now.Sub(reqTime) < time.Minute {
						validRequests = append(validRequests, reqTime)
					}
				}
				clients[clientIP] = validRequests
			}

			// Check rate limit
			requestCount := len(clients[clientIP])
			if requestCount >= requestsPerMinute {
				mutex.Unlock() // Unlock before logging and responding

				utils.LogWarn("Rate limit exceeded", map[string]interface{}{
					"client_ip": clientIP,
					"requests":  requestCount,
					"limit":     requestsPerMinute,
				})

				w.Header().Set("Retry-After", "60")
				utils.WriteError(w, http.StatusTooManyRequests, "RATE_LIMIT_EXCEEDED", "Too many requests", nil)
				return
			}

			// Add current request
			clients[clientIP] = append(clients[clientIP], now)
			mutex.Unlock() // Unlock after all write operations

			next.ServeHTTP(w, r)
		})
	}
}

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware(next http.Handler) http.Handler {
	return middleware.RequestID(next)
}

// LoggerMiddleware provides structured request logging
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap the response writer to capture status code
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		// Process request
		next.ServeHTTP(ww, r)

		// Log request
		duration := time.Since(start)
		utils.LogHTTPRequest(
			r.Method,
			r.URL.Path,
			r.UserAgent(),
			ww.Status(),
			duration,
			map[string]interface{}{
				"request_id":    middleware.GetReqID(r.Context()),
				"remote_addr":   r.RemoteAddr,
				"bytes_written": ww.BytesWritten(),
			},
		)
	})
}

// SecurityMiddleware adds security headers
func SecurityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add security headers
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")

		next.ServeHTTP(w, r)
	})
}

// ErrorHandlerMiddleware provides centralized error handling
func ErrorHandlerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a custom response writer to capture errors
		ww := &errorResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(ww, r)

		// Handle any errors that occurred
		if ww.error != nil {
			utils.LogError(ww.error, "Request error", map[string]interface{}{
				"method":     r.Method,
				"path":       r.URL.Path,
				"status":     ww.statusCode,
				"request_id": middleware.GetReqID(r.Context()),
			})

			// Write error response if not already written
			if !ww.written {
				utils.WriteError(w, ww.statusCode, "REQUEST_ERROR", ww.error.Error(), nil)
			}
		}
	})
}

// errorResponseWriter wraps http.ResponseWriter to capture errors
type errorResponseWriter struct {
	http.ResponseWriter
	statusCode int
	error      error
	written    bool
}

func (w *errorResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.written = true
	w.ResponseWriter.WriteHeader(code)
}

func (w *errorResponseWriter) Write(b []byte) (int, error) {
	w.written = true
	return w.ResponseWriter.Write(b)
}

// ValidationMiddleware provides request validation
func ValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate request method
		if r.Method != http.MethodGet && r.Method != http.MethodPost &&
			r.Method != http.MethodPut && r.Method != http.MethodDelete &&
			r.Method != http.MethodPatch && r.Method != http.MethodOptions {
			utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", nil)
			return
		}

		// Validate content type for POST/PUT requests
		if (r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch) &&
			r.Header.Get("Content-Type") != "application/json" {
			utils.WriteError(w, http.StatusUnsupportedMediaType, "INVALID_CONTENT_TYPE", "Content-Type must be application/json", nil)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// MetricsMiddleware collects request metrics
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap response writer
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		// Process request
		next.ServeHTTP(ww, r)

		// Collect metrics
		duration := time.Since(start)

		// Log metrics (in a real implementation, you might send to a metrics system)
		utils.LogInfo("Request metrics", map[string]interface{}{
			"method":        r.Method,
			"path":          r.URL.Path,
			"status_code":   ww.Status(),
			"duration_ms":   duration.Milliseconds(),
			"bytes_written": ww.BytesWritten(),
		})
	})
}
