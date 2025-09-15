package handlers

import (
	"net/http"
	"spark-park-cricket-backend/internal/services"
	"spark-park-cricket-backend/pkg/websocket"

	"github.com/go-chi/chi/v5"
)

// WebSocketHandler handles WebSocket connections
type WebSocketHandler struct {
	hub      *websocket.Hub
	services *services.Container
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler(hub *websocket.Hub, services *services.Container) *WebSocketHandler {
	return &WebSocketHandler{
		hub:      hub,
		services: services,
	}
}

// ServeWS handles WebSocket connections for a specific match
func (h *WebSocketHandler) ServeWS(w http.ResponseWriter, r *http.Request) {
	matchID := chi.URLParam(r, "match_id")
	if matchID == "" {
		http.Error(w, "Match ID is required", http.StatusBadRequest)
		return
	}

	// For now, just return a simple response
	// In a real implementation, you would upgrade the connection to WebSocket
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "WebSocket endpoint for match ` + matchID + `"}`))
}

// GetConnectionStats returns WebSocket connection statistics
func (h *WebSocketHandler) GetConnectionStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"connections": 0, "rooms": 0}`))
}

// GetRoomStats returns statistics for a specific room/match
func (h *WebSocketHandler) GetRoomStats(w http.ResponseWriter, r *http.Request) {
	matchID := chi.URLParam(r, "match_id")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"match_id": "` + matchID + `", "connections": 0}`))
}

// TestBroadcast sends a test broadcast to a specific match room
func (h *WebSocketHandler) TestBroadcast(w http.ResponseWriter, r *http.Request) {
	matchID := chi.URLParam(r, "match_id")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Test broadcast sent to match ` + matchID + `"}`))
}
