package monitoring

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics holds all Prometheus metrics for the cricket application
type Metrics struct {
	// HTTP metrics
	HTTPRequestsTotal    *prometheus.CounterVec
	HTTPRequestDuration  *prometheus.HistogramVec
	HTTPRequestsInFlight *prometheus.GaugeVec

	// Cricket-specific metrics
	BallAdditionsTotal   *prometheus.CounterVec
	BallAdditionDuration *prometheus.HistogramVec
	MatchCreationsTotal  *prometheus.CounterVec
	SeriesCreationsTotal *prometheus.CounterVec

	// Database metrics
	DatabaseConnections   *prometheus.GaugeVec
	DatabaseQueryDuration *prometheus.HistogramVec
	DatabaseErrorsTotal   *prometheus.CounterVec

	// Ball API specific database metrics
	BallAPIDatabaseOperations *prometheus.CounterVec
	BallAPIDatabaseDuration   *prometheus.HistogramVec

	// Cache metrics
	CacheHitsTotal         *prometheus.CounterVec
	CacheMissesTotal       *prometheus.CounterVec
	CacheOperationsTotal   *prometheus.CounterVec
	CacheOperationDuration *prometheus.HistogramVec

	// WebSocket metrics
	WebSocketConnections   *prometheus.GaugeVec
	WebSocketMessagesTotal *prometheus.CounterVec

	// System metrics
	MemoryUsage     *prometheus.GaugeVec
	CPUUsage        *prometheus.GaugeVec
	GoroutinesCount *prometheus.GaugeVec
}

var (
	metricsInstance *Metrics
	metricsOnce     sync.Once
)

// NewMetrics creates a new Metrics instance with all Prometheus metrics (singleton)
func NewMetrics() *Metrics {
	metricsOnce.Do(func() {
		metricsInstance = &Metrics{
			// HTTP metrics
			HTTPRequestsTotal: promauto.NewCounterVec(
				prometheus.CounterOpts{
					Name: "http_requests_total",
					Help: "Total number of HTTP requests",
				},
				[]string{"method", "endpoint", "status"},
			),

			HTTPRequestDuration: promauto.NewHistogramVec(
				prometheus.HistogramOpts{
					Name:    "http_request_duration_seconds",
					Help:    "HTTP request duration in seconds",
					Buckets: prometheus.DefBuckets,
				},
				[]string{"method", "endpoint"},
			),

			HTTPRequestsInFlight: promauto.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: "http_requests_in_flight",
					Help: "Current number of HTTP requests being processed",
				},
				[]string{"method", "endpoint"},
			),

			// Cricket-specific metrics
			BallAdditionsTotal: promauto.NewCounterVec(
				prometheus.CounterOpts{
					Name: "cricket_ball_additions_total",
					Help: "Total number of ball additions",
				},
				[]string{"match_id", "innings", "ball_type", "run_type"},
			),

			BallAdditionDuration: promauto.NewHistogramVec(
				prometheus.HistogramOpts{
					Name:    "cricket_ball_addition_duration_seconds",
					Help:    "Time taken to add a ball",
					Buckets: []float64{0.1, 0.5, 1.0, 2.0, 5.0, 10.0, 30.0},
				},
				[]string{"match_id", "innings"},
			),

			MatchCreationsTotal: promauto.NewCounterVec(
				prometheus.CounterOpts{
					Name: "cricket_match_creations_total",
					Help: "Total number of match creations",
				},
				[]string{"series_id", "status"},
			),

			SeriesCreationsTotal: promauto.NewCounterVec(
				prometheus.CounterOpts{
					Name: "cricket_series_creations_total",
					Help: "Total number of series creations",
				},
				[]string{"status"},
			),

			// Database metrics
			DatabaseConnections: promauto.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: "database_connections_active",
					Help: "Number of active database connections",
				},
				[]string{"database", "status"},
			),

			DatabaseQueryDuration: promauto.NewHistogramVec(
				prometheus.HistogramOpts{
					Name:    "database_query_duration_seconds",
					Help:    "Database query duration in seconds",
					Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1.0, 5.0},
				},
				[]string{"operation", "table", "match_id"},
			),

			DatabaseErrorsTotal: promauto.NewCounterVec(
				prometheus.CounterOpts{
					Name: "database_errors_total",
					Help: "Total number of database errors",
				},
				[]string{"operation", "error_type", "table"},
			),

			// Ball API specific database metrics
			BallAPIDatabaseOperations: promauto.NewCounterVec(
				prometheus.CounterOpts{
					Name: "ball_api_database_operations_total",
					Help: "Total number of database operations during ball addition",
				},
				[]string{"operation", "table", "match_id", "status"},
			),

			BallAPIDatabaseDuration: promauto.NewHistogramVec(
				prometheus.HistogramOpts{
					Name:    "ball_api_database_duration_seconds",
					Help:    "Database operation duration during ball addition",
					Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1.0, 2.0, 5.0},
				},
				[]string{"operation", "table", "match_id"},
			),

			// Cache metrics
			CacheHitsTotal: promauto.NewCounterVec(
				prometheus.CounterOpts{
					Name: "cache_hits_total",
					Help: "Total number of cache hits",
				},
				[]string{"cache_type", "key_pattern"},
			),

			CacheMissesTotal: promauto.NewCounterVec(
				prometheus.CounterOpts{
					Name: "cache_misses_total",
					Help: "Total number of cache misses",
				},
				[]string{"cache_type", "key_pattern"},
			),

			CacheOperationsTotal: promauto.NewCounterVec(
				prometheus.CounterOpts{
					Name: "cache_operations_total",
					Help: "Total number of cache operations",
				},
				[]string{"operation", "cache_type"},
			),

			CacheOperationDuration: promauto.NewHistogramVec(
				prometheus.HistogramOpts{
					Name:    "cache_operation_duration_seconds",
					Help:    "Cache operation duration in seconds",
					Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1.0},
				},
				[]string{"operation", "cache_type"},
			),

			// WebSocket metrics
			WebSocketConnections: promauto.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: "websocket_connections_active",
					Help: "Number of active WebSocket connections",
				},
				[]string{"match_id"},
			),

			WebSocketMessagesTotal: promauto.NewCounterVec(
				prometheus.CounterOpts{
					Name: "websocket_messages_total",
					Help: "Total number of WebSocket messages sent",
				},
				[]string{"match_id", "message_type"},
			),

			// System metrics
			MemoryUsage: promauto.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: "cricket_memory_usage_bytes",
					Help: "Cricket application memory usage in bytes",
				},
				[]string{"service"},
			),

			CPUUsage: promauto.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: "cricket_cpu_seconds_total",
					Help: "Cricket application CPU time spent in seconds",
				},
				[]string{"service"},
			),

			GoroutinesCount: promauto.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: "cricket_goroutines_total",
					Help: "Number of goroutines in cricket application",
				},
				[]string{"service"},
			),
		}
	})
	return metricsInstance
}

// RecordHTTPRequest records HTTP request metrics
func (m *Metrics) RecordHTTPRequest(method, endpoint, status string, duration time.Duration) {
	m.HTTPRequestsTotal.WithLabelValues(method, endpoint, status).Inc()
	m.HTTPRequestDuration.WithLabelValues(method, endpoint).Observe(duration.Seconds())
}

// RecordBallAddition records ball addition metrics
func (m *Metrics) RecordBallAddition(matchID, innings, ballType, runType string, duration time.Duration) {
	m.BallAdditionsTotal.WithLabelValues(matchID, innings, ballType, runType).Inc()
	m.BallAdditionDuration.WithLabelValues(matchID, innings).Observe(duration.Seconds())
}

// RecordMatchCreation records match creation metrics
func (m *Metrics) RecordMatchCreation(seriesID, status string) {
	m.MatchCreationsTotal.WithLabelValues(seriesID, status).Inc()
}

// RecordSeriesCreation records series creation metrics
func (m *Metrics) RecordSeriesCreation(status string) {
	m.SeriesCreationsTotal.WithLabelValues(status).Inc()
}

// RecordDatabaseQuery records database query metrics
func (m *Metrics) RecordDatabaseQuery(operation, table string, duration time.Duration) {
	m.DatabaseQueryDuration.WithLabelValues(operation, table, "").Observe(duration.Seconds())
}

// RecordDatabaseError records database error metrics
func (m *Metrics) RecordDatabaseError(operation, errorType string) {
	m.DatabaseErrorsTotal.WithLabelValues(operation, errorType, "").Inc()
}

// RecordBallAPIDatabaseOperation records ball API specific database operations
func (m *Metrics) RecordBallAPIDatabaseOperation(operation, table, matchID, status string, duration time.Duration) {
	m.BallAPIDatabaseOperations.WithLabelValues(operation, table, matchID, status).Inc()
	m.BallAPIDatabaseDuration.WithLabelValues(operation, table, matchID).Observe(duration.Seconds())
}

// RecordBallAPIDatabaseError records ball API database errors
func (m *Metrics) RecordBallAPIDatabaseError(operation, table, matchID, errorType string) {
	m.DatabaseErrorsTotal.WithLabelValues(operation, errorType, table).Inc()
	m.BallAPIDatabaseOperations.WithLabelValues(operation, table, matchID, "error").Inc()
}

// RecordCacheHit records cache hit metrics
func (m *Metrics) RecordCacheHit(cacheType, keyPattern string) {
	m.CacheHitsTotal.WithLabelValues(cacheType, keyPattern).Inc()
}

// RecordCacheMiss records cache miss metrics
func (m *Metrics) RecordCacheMiss(cacheType, keyPattern string) {
	m.CacheMissesTotal.WithLabelValues(cacheType, keyPattern).Inc()
}

// RecordCacheOperation records cache operation metrics
func (m *Metrics) RecordCacheOperation(operation, cacheType string, duration time.Duration) {
	m.CacheOperationsTotal.WithLabelValues(operation, cacheType).Inc()
	m.CacheOperationDuration.WithLabelValues(operation, cacheType).Observe(duration.Seconds())
}

// RecordWebSocketConnection records WebSocket connection metrics
func (m *Metrics) RecordWebSocketConnection(matchID string, connected bool) {
	if connected {
		m.WebSocketConnections.WithLabelValues(matchID).Inc()
	} else {
		m.WebSocketConnections.WithLabelValues(matchID).Dec()
	}
}

// RecordWebSocketMessage records WebSocket message metrics
func (m *Metrics) RecordWebSocketMessage(matchID, messageType string) {
	m.WebSocketMessagesTotal.WithLabelValues(matchID, messageType).Inc()
}

// UpdateSystemMetrics updates system metrics
func (m *Metrics) UpdateSystemMetrics() {
	// This would typically be called periodically to update system metrics
	// Implementation depends on your specific needs
}
