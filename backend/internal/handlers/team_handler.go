package handlers

import (
	"encoding/json"
	"net/http"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/services"
	"spark-park-cricket-backend/internal/utils"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// TeamHandler handles HTTP requests for team operations
type TeamHandler struct {
	service *services.TeamService
}

// NewTeamHandler creates a new team handler
func NewTeamHandler(service *services.TeamService) *TeamHandler {
	return &TeamHandler{
		service: service,
	}
}

// ListTeams handles GET /api/v1/teams
func (h *TeamHandler) ListTeams(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	// Set default values
	limit := 20
	offset := 0

	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	if offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	// Create filters
	filters := &models.TeamFilters{
		Limit:  limit,
		Offset: offset,
	}

	// Get teams from service
	teams, err := h.service.ListTeams(r.Context(), filters)
	if err != nil {
		utils.WriteInternalError(w, err.Error())
		return
	}

	utils.WriteSuccess(w, teams)
}

// CreateTeam handles POST /api/v1/teams
func (h *TeamHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {

	var req models.CreateTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteValidationError(w, "Invalid request body", err.Error())
		return
	}

	// Create team
	team, err := h.service.CreateTeam(r.Context(), &req)
	if err != nil {
		utils.WriteInternalError(w, err.Error())
		return
	}

	utils.WriteCreated(w, team)
}

// GetTeam handles GET /api/v1/teams/{id}
func (h *TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.WriteValidationError(w, "Team ID is required", nil)
		return
	}

	team, err := h.service.GetTeam(r.Context(), id)
	if err != nil {
		if err.Error() == "team not found" {
			utils.WriteNotFound(w, "Team")
			return
		}
		utils.WriteInternalError(w, err.Error())
		return
	}

	utils.WriteSuccess(w, team)
}

// UpdateTeam handles PUT /api/v1/teams/{id}
func (h *TeamHandler) UpdateTeam(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.WriteValidationError(w, "Team ID is required", nil)
		return
	}

	var req models.UpdateTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteValidationError(w, "Invalid request body", err.Error())
		return
	}

	team, err := h.service.UpdateTeam(r.Context(), id, &req)
	if err != nil {
		if err.Error() == "team not found" {
			utils.WriteNotFound(w, "Team")
			return
		}
		utils.WriteInternalError(w, err.Error())
		return
	}

	utils.WriteSuccess(w, team)
}

// DeleteTeam handles DELETE /api/v1/teams/{id}
func (h *TeamHandler) DeleteTeam(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.WriteValidationError(w, "Team ID is required", nil)
		return
	}

	err := h.service.DeleteTeam(r.Context(), id)
	if err != nil {
		if err.Error() == "team not found" {
			utils.WriteNotFound(w, "Team")
			return
		}
		utils.WriteInternalError(w, err.Error())
		return
	}

	utils.WriteSuccess(w, map[string]string{"message": "Team deleted successfully"})
}

// ListTeamPlayers handles GET /api/v1/teams/{id}/players
func (h *TeamHandler) ListTeamPlayers(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.WriteValidationError(w, "Team ID is required", nil)
		return
	}

	players, err := h.service.GetTeamPlayers(r.Context(), id)
	if err != nil {
		if err.Error() == "team not found" {
			utils.WriteNotFound(w, "Team")
			return
		}
		utils.WriteInternalError(w, err.Error())
		return
	}

	utils.WriteSuccess(w, players)
}

// AddTeamPlayer handles POST /api/v1/teams/{id}/players
func (h *TeamHandler) AddTeamPlayer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.WriteValidationError(w, "Team ID is required", nil)
		return
	}

	var req models.CreatePlayerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteValidationError(w, "Invalid request body", err.Error())
		return
	}

	player, err := h.service.AddPlayerToTeam(r.Context(), id, &req)
	if err != nil {
		if err.Error() == "team not found" {
			utils.WriteNotFound(w, "Team")
			return
		}
		utils.WriteInternalError(w, err.Error())
		return
	}

	utils.WriteCreated(w, player)
}
