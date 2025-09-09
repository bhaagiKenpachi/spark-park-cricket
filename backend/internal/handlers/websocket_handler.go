package handlers

import (
	"encoding/json"
	"net/http"
	"spark-park-cricket-backend/internal/services"
	"spark-park-cricket-backend/internal/utils"
	"spark-park-cricket-backend/pkg/websocket"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// WebSocketHandler handles WebSocket connections and real-time updates
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

// ServeWS handles WebSocket connections for live match updates
func (h *WebSocketHandler) ServeWS(w http.ResponseWriter, r *http.Request) {
	matchID := chi.URLParam(r, "match_id")
	if matchID == "" {
		utils.WriteValidationError(w, "Match ID is required", nil)
		return
	}

	// Validate match exists
	_, err := h.services.Match.GetMatch(r.Context(), matchID)
	if err != nil {
		utils.WriteNotFound(w, "Match")
		return
	}

	// Generate client ID
	clientID := uuid.New().String()

	// Serve WebSocket connection
	h.hub.ServeWS(w, r, matchID, clientID)
}

// GetConnectionStats returns WebSocket connection statistics
func (h *WebSocketHandler) GetConnectionStats(w http.ResponseWriter, r *http.Request) {
	stats := map[string]interface{}{
		"total_clients": h.hub.GetTotalClients(),
		"total_rooms":   h.hub.GetTotalRooms(),
		"timestamp":     time.Now().Unix(),
	}

	utils.WriteSuccess(w, stats)
}

// GetRoomStats returns statistics for a specific room (match)
func (h *WebSocketHandler) GetRoomStats(w http.ResponseWriter, r *http.Request) {
	matchID := chi.URLParam(r, "match_id")
	if matchID == "" {
		utils.WriteValidationError(w, "Match ID is required", nil)
		return
	}

	clientCount := h.hub.GetRoomClients(matchID)
	stats := map[string]interface{}{
		"match_id":     matchID,
		"client_count": clientCount,
		"timestamp":    time.Now().Unix(),
	}

	utils.WriteSuccess(w, stats)
}

// BroadcastScoreUpdate broadcasts a score update to all clients in a match room
func (h *WebSocketHandler) BroadcastScoreUpdate(matchID string, scoreboard interface{}) {
	message := websocket.Message{
		Type:   "score_update",
		RoomID: matchID,
		Data: map[string]interface{}{
			"scoreboard": scoreboard,
			"timestamp":  time.Now().Unix(),
		},
	}

	h.hub.BroadcastToRoom(matchID, message)
}

// BroadcastBallUpdate broadcasts a ball update to all clients in a match room
func (h *WebSocketHandler) BroadcastBallUpdate(matchID string, ballEvent interface{}, scoreboard interface{}) {
	message := websocket.Message{
		Type:   "ball_update",
		RoomID: matchID,
		Data: map[string]interface{}{
			"ball_event": ballEvent,
			"scoreboard": scoreboard,
			"timestamp":  time.Now().Unix(),
		},
	}

	h.hub.BroadcastToRoom(matchID, message)
}

// BroadcastOverUpdate broadcasts an over completion to all clients in a match room
func (h *WebSocketHandler) BroadcastOverUpdate(matchID string, over interface{}, scoreboard interface{}) {
	message := websocket.Message{
		Type:   "over_update",
		RoomID: matchID,
		Data: map[string]interface{}{
			"over":       over,
			"scoreboard": scoreboard,
			"timestamp":  time.Now().Unix(),
		},
	}

	h.hub.BroadcastToRoom(matchID, message)
}

// BroadcastMatchUpdate broadcasts a match status update to all clients in a match room
func (h *WebSocketHandler) BroadcastMatchUpdate(matchID string, match interface{}) {
	message := websocket.Message{
		Type:   "match_update",
		RoomID: matchID,
		Data: map[string]interface{}{
			"match":     match,
			"timestamp": time.Now().Unix(),
		},
	}

	h.hub.BroadcastToRoom(matchID, message)
}

// BroadcastWicketUpdate broadcasts a wicket update to all clients in a match room
func (h *WebSocketHandler) BroadcastWicketUpdate(matchID string, wicketData interface{}, scoreboard interface{}) {
	message := websocket.Message{
		Type:   "wicket_update",
		RoomID: matchID,
		Data: map[string]interface{}{
			"wicket":     wicketData,
			"scoreboard": scoreboard,
			"timestamp":  time.Now().Unix(),
		},
	}

	h.hub.BroadcastToRoom(matchID, message)
}

// TestBroadcast sends a test message to all clients in a match room
func (h *WebSocketHandler) TestBroadcast(w http.ResponseWriter, r *http.Request) {
	matchID := chi.URLParam(r, "match_id")
	if matchID == "" {
		utils.WriteValidationError(w, "Match ID is required", nil)
		return
	}

	var req struct {
		Message string `json:"message"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteValidationError(w, "Invalid request body", err.Error())
		return
	}

	message := websocket.Message{
		Type:   "test_message",
		RoomID: matchID,
		Data: map[string]interface{}{
			"message":   req.Message,
			"timestamp": time.Now().Unix(),
		},
	}

	h.hub.BroadcastToRoom(matchID, message)

	utils.WriteSuccess(w, map[string]interface{}{
		"message":  "Test message broadcasted successfully",
		"match_id": matchID,
		"clients":  h.hub.GetRoomClients(matchID),
	})
}
