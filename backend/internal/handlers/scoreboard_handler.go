package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/services"
	"spark-park-cricket-backend/internal/utils"

	"github.com/go-chi/chi/v5"
)

// ScoreboardHandler handles scoreboard-related HTTP requests
type ScoreboardHandler struct {
	service *services.RealtimeScoreboardService
}

// NewScoreboardHandler creates a new scoreboard handler
func NewScoreboardHandler(service *services.RealtimeScoreboardService) *ScoreboardHandler {
	return &ScoreboardHandler{
		service: service,
	}
}

// GetScoreboard handles GET /api/v1/scoreboard/{match_id}
func (h *ScoreboardHandler) GetScoreboard(w http.ResponseWriter, r *http.Request) {

	matchID := chi.URLParam(r, "match_id")

	if matchID == "" {
		utils.WriteValidationError(w, "Match ID is required", nil)
		return
	}

	scoreboard, err := h.service.GetScoreboard(r.Context(), matchID)
	if err != nil {
		utils.WriteInternalError(w, err.Error())
		return
	}

	utils.WriteSuccess(w, scoreboard)
}

// AddBall handles POST /api/v1/scoreboard/{match_id}/ball
func (h *ScoreboardHandler) AddBall(w http.ResponseWriter, r *http.Request) {
	matchID := chi.URLParam(r, "match_id")
	if matchID == "" {
		utils.WriteValidationError(w, "Match ID is required", nil)
		return
	}

	var ballEvent models.BallEvent
	if err := json.NewDecoder(r.Body).Decode(&ballEvent); err != nil {
		utils.WriteValidationError(w, "Invalid request body", err.Error())
		return
	}

	scoreboard, err := h.service.AddBall(r.Context(), matchID, &ballEvent)
	if err != nil {
		utils.WriteInternalError(w, err.Error())
		return
	}

	utils.WriteCreated(w, map[string]interface{}{
		"message":    "Ball added successfully",
		"scoreboard": scoreboard,
		"ball_event": ballEvent,
	})
}

// UpdateScore handles PUT /api/v1/scoreboard/{match_id}/score
func (h *ScoreboardHandler) UpdateScore(w http.ResponseWriter, r *http.Request) {
	matchID := chi.URLParam(r, "match_id")
	if matchID == "" {
		utils.WriteValidationError(w, "Match ID is required", nil)
		return
	}

	var req models.UpdateScoreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteValidationError(w, "Invalid request body", err.Error())
		return
	}

	scoreboard, err := h.service.UpdateScore(r.Context(), matchID, &req)
	if err != nil {
		utils.WriteInternalError(w, err.Error())
		return
	}

	utils.WriteSuccess(w, map[string]interface{}{
		"message":    "Score updated successfully",
		"scoreboard": scoreboard,
	})
}

// UpdateWicket handles PUT /api/v1/scoreboard/{match_id}/wicket
func (h *ScoreboardHandler) UpdateWicket(w http.ResponseWriter, r *http.Request) {
	matchID := chi.URLParam(r, "match_id")
	if matchID == "" {
		utils.WriteValidationError(w, "Match ID is required", nil)
		return
	}

	var req models.UpdateWicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteValidationError(w, "Invalid request body", err.Error())
		return
	}

	scoreboard, err := h.service.UpdateWicket(r.Context(), matchID, &req)
	if err != nil {
		utils.WriteInternalError(w, err.Error())
		return
	}

	utils.WriteSuccess(w, map[string]interface{}{
		"message":    "Wicket updated successfully",
		"scoreboard": scoreboard,
	})
}
