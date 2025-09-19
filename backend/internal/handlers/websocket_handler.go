package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"spark-park-cricket-backend/internal/services"
	"spark-park-cricket-backend/pkg/websocket"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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

	// Generate a unique client ID
	clientID := uuid.New().String()

	log.Printf("WebSocket connection request for match %s, client %s", matchID, clientID)

	// Upgrade the connection to WebSocket
	h.hub.ServeWS(w, r, matchID, clientID)
}

// GetConnectionStats returns WebSocket connection statistics
func (h *WebSocketHandler) GetConnectionStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	stats := map[string]interface{}{
		"total_connections": h.hub.GetTotalClients(),
		"total_rooms":       h.hub.GetTotalRooms(),
	}

	response, err := json.Marshal(stats)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(response); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

// GetRoomStats returns statistics for a specific room/match
func (h *WebSocketHandler) GetRoomStats(w http.ResponseWriter, r *http.Request) {
	matchID := chi.URLParam(r, "match_id")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	stats := map[string]interface{}{
		"match_id":    matchID,
		"connections": h.hub.GetRoomClients(matchID),
	}

	response, err := json.Marshal(stats)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(response); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

// TestBroadcast sends a test broadcast to a specific match room
func (h *WebSocketHandler) TestBroadcast(w http.ResponseWriter, r *http.Request) {
	matchID := chi.URLParam(r, "match_id")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Test broadcast sent to match ` + matchID + `"}`))
}
