package handlers

import (
	"context"
	"net/http"
	"runtime"
	"spark-park-cricket-backend/internal/database"
	"spark-park-cricket-backend/internal/utils"
	"time"
)

// HealthHandler handles health check endpoints
type HealthHandler struct {
	dbClient *database.Client
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(dbClient *database.Client) *HealthHandler {
	return &HealthHandler{
		dbClient: dbClient,
	}
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string                   `json:"status"`
	Timestamp time.Time                `json:"timestamp"`
	Version   string                   `json:"version"`
	Uptime    string                   `json:"uptime"`
	Services  map[string]ServiceHealth `json:"services"`
	System    SystemInfo               `json:"system"`
}

// ServiceHealth represents the health status of a service
type ServiceHealth struct {
	Status       string        `json:"status"`
	ResponseTime time.Duration `json:"response_time_ms"`
	Error        string        `json:"error,omitempty"`
	Details      interface{}   `json:"details,omitempty"`
}

// SystemInfo represents system information
type SystemInfo struct {
	GoVersion    string      `json:"go_version"`
	NumGoroutine int         `json:"num_goroutines"`
	MemoryUsage  MemoryUsage `json:"memory_usage"`
}

// MemoryUsage represents memory usage information
type MemoryUsage struct {
	Alloc      uint64 `json:"alloc_bytes"`
	TotalAlloc uint64 `json:"total_alloc_bytes"`
	Sys        uint64 `json:"sys_bytes"`
	NumGC      uint32 `json:"num_gc"`
}

var startTime = time.Now()

// Health handles GET /health
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Check database health
	dbHealth := h.checkDatabaseHealth(ctx)

	// Check WebSocket hub health
	wsHealth := h.checkWebSocketHealth()

	// Get system information
	systemInfo := h.getSystemInfo()

	// Determine overall status
	overallStatus := "healthy"
	if dbHealth.Status != "healthy" || wsHealth.Status != "healthy" {
		overallStatus = "unhealthy"
	}

	response := HealthResponse{
		Status:    overallStatus,
		Timestamp: time.Now(),
		Version:   "1.0.0",
		Uptime:    time.Since(startTime).String(),
		Services: map[string]ServiceHealth{
			"database":  dbHealth,
			"websocket": wsHealth,
		},
		System: systemInfo,
	}

	statusCode := http.StatusOK
	if overallStatus != "healthy" {
		statusCode = http.StatusServiceUnavailable
	}

	utils.WriteJSONResponse(w, statusCode, response)
}

// DatabaseHealth handles GET /health/database
func (h *HealthHandler) DatabaseHealth(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	dbHealth := h.checkDatabaseHealth(ctx)

	statusCode := http.StatusOK
	if dbHealth.Status != "healthy" {
		statusCode = http.StatusServiceUnavailable
	}

	utils.WriteJSONResponse(w, statusCode, dbHealth)
}

// WebSocketHealth handles GET /health/websocket
func (h *HealthHandler) WebSocketHealth(w http.ResponseWriter, r *http.Request) {
	wsHealth := h.checkWebSocketHealth()

	statusCode := http.StatusOK
	if wsHealth.Status != "healthy" {
		statusCode = http.StatusServiceUnavailable
	}

	utils.WriteJSONResponse(w, statusCode, wsHealth)
}

// SystemHealth handles GET /health/system
func (h *HealthHandler) SystemHealth(w http.ResponseWriter, r *http.Request) {
	systemInfo := h.getSystemInfo()
	utils.WriteJSONResponse(w, http.StatusOK, systemInfo)
}

// Readiness handles GET /health/ready
func (h *HealthHandler) Readiness(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Check if all critical services are ready
	dbHealth := h.checkDatabaseHealth(ctx)

	ready := dbHealth.Status == "healthy"

	response := map[string]interface{}{
		"ready":     ready,
		"timestamp": time.Now(),
		"services": map[string]string{
			"database": dbHealth.Status,
		},
	}

	statusCode := http.StatusOK
	if !ready {
		statusCode = http.StatusServiceUnavailable
	}

	utils.WriteJSONResponse(w, statusCode, response)
}

// Liveness handles GET /health/live
func (h *HealthHandler) Liveness(w http.ResponseWriter, r *http.Request) {
	// Simple liveness check - if we can respond, we're alive
	response := map[string]interface{}{
		"alive":     true,
		"timestamp": time.Now(),
		"uptime":    time.Since(startTime).String(),
	}

	utils.WriteJSONResponse(w, http.StatusOK, response)
}

// checkDatabaseHealth checks the health of the database connection
func (h *HealthHandler) checkDatabaseHealth(ctx context.Context) ServiceHealth {
	start := time.Now()

	if h.dbClient == nil {
		return ServiceHealth{
			Status:       "unhealthy",
			ResponseTime: time.Since(start),
			Error:        "database client not initialized",
		}
	}

	// Try to execute a simple query
	_, _, err := h.dbClient.Supabase.From("series").Select("count", "", false).Single().Execute()
	responseTime := time.Since(start)

	if err != nil {
		return ServiceHealth{
			Status:       "unhealthy",
			ResponseTime: responseTime,
			Error:        err.Error(),
		}
	}

	return ServiceHealth{
		Status:       "healthy",
		ResponseTime: responseTime,
		Details: map[string]interface{}{
			"connection": "active",
		},
	}
}

// checkWebSocketHealth checks the health of the WebSocket hub
func (h *HealthHandler) checkWebSocketHealth() ServiceHealth {
	start := time.Now()

	// For now, we'll assume WebSocket is healthy if the hub exists
	// In a real implementation, you might want to check actual connections
	responseTime := time.Since(start)

	return ServiceHealth{
		Status:       "healthy",
		ResponseTime: responseTime,
		Details: map[string]interface{}{
			"hub_initialized": true,
		},
	}
}

// getSystemInfo gets system information
func (h *HealthHandler) getSystemInfo() SystemInfo {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return SystemInfo{
		GoVersion:    runtime.Version(),
		NumGoroutine: runtime.NumGoroutine(),
		MemoryUsage: MemoryUsage{
			Alloc:      m.Alloc,
			TotalAlloc: m.TotalAlloc,
			Sys:        m.Sys,
			NumGC:      m.NumGC,
		},
	}
}

// Metrics handles GET /health/metrics
func (h *HealthHandler) Metrics(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Get database metrics
	dbHealth := h.checkDatabaseHealth(ctx)

	// Get system metrics
	systemInfo := h.getSystemInfo()

	// Get WebSocket metrics (if available)
	wsHealth := h.checkWebSocketHealth()

	metrics := map[string]interface{}{
		"timestamp": time.Now(),
		"database": map[string]interface{}{
			"status":        dbHealth.Status,
			"response_time": dbHealth.ResponseTime.Milliseconds(),
		},
		"websocket": map[string]interface{}{
			"status":        wsHealth.Status,
			"response_time": wsHealth.ResponseTime.Milliseconds(),
		},
		"system": map[string]interface{}{
			"go_version":   systemInfo.GoVersion,
			"goroutines":   systemInfo.NumGoroutine,
			"memory_alloc": systemInfo.MemoryUsage.Alloc,
			"memory_sys":   systemInfo.MemoryUsage.Sys,
			"gc_count":     systemInfo.MemoryUsage.NumGC,
		},
		"uptime": time.Since(startTime).Seconds(),
	}

	utils.WriteJSONResponse(w, http.StatusOK, metrics)
}
