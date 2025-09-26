package services

import (
	"net/http"
	"strconv"
	"time"

	"spark-park-cricket-backend/internal/monitoring"

	"github.com/go-chi/chi/v5/middleware"
)

// PrometheusMiddleware creates a middleware that records Prometheus metrics
func PrometheusMiddleware(metrics *monitoring.Metrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Create a response writer wrapper to capture status code
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			// Record request in flight
			metrics.HTTPRequestsInFlight.WithLabelValues(r.Method, r.URL.Path).Inc()
			defer func() {
				metrics.HTTPRequestsInFlight.WithLabelValues(r.Method, r.URL.Path).Dec()
			}()

			// Call the next handler
			next.ServeHTTP(ww, r)

			// Record metrics
			duration := time.Since(start)
			status := strconv.Itoa(ww.Status())

			metrics.RecordHTTPRequest(r.Method, r.URL.Path, status, duration)
		})
	}
}

// PrometheusHandler creates a handler for Prometheus metrics endpoint
func PrometheusHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This will be handled by the Prometheus client library
		// The actual metrics endpoint will be registered separately
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("# Prometheus metrics endpoint\n"))
	})
}
