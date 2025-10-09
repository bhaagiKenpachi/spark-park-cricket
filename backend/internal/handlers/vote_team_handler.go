package handlers

import (
	"encoding/json"
	"net/http"
	contextkeys "spark-park-cricket-backend/internal/context"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/services"
	"spark-park-cricket-backend/internal/utils"

	"github.com/go-chi/chi/v5"
)

// Helper functions matching the project's response pattern
func errorResponse(w http.ResponseWriter, statusCode int, code, message string) {
	utils.WriteErrorResponse(w, statusCode, code, message)
}

func successResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	utils.WriteJSONResponse(w, statusCode, map[string]interface{}{"data": data})
}

// VoteTeamHandler handles vote team-related HTTP requests
type VoteTeamHandler struct {
	teamService services.VoteTeamServiceInterface
}

// NewVoteTeamHandler creates a new vote team handler
func NewVoteTeamHandler(teamService services.VoteTeamServiceInterface) *VoteTeamHandler {
	return &VoteTeamHandler{
		teamService: teamService,
	}
}

// CreateTeam handles POST /api/v1/votes/:vote_id/teams
func (h *VoteTeamHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	voteID := chi.URLParam(r, "vote_id")

	// Get user ID from context using the correct context key
	userID, ok := ctx.Value(contextkeys.UserIDKey).(string)
	if !ok || userID == "" {
		errorResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	// Parse request
	var req models.CreateVoteTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	// Set vote ID from URL
	req.VoteID = voteID

	// Create team
	team, err := h.teamService.CreateTeam(ctx, userID, &req)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "CREATE_TEAM_ERROR", err.Error())
		return
	}

	successResponse(w, http.StatusCreated, team)
}

// GetTeamsByVoteID handles GET /api/v1/votes/:vote_id/teams
func (h *VoteTeamHandler) GetTeamsByVoteID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	voteID := chi.URLParam(r, "vote_id")

	teams, err := h.teamService.GetTeamsByVoteID(ctx, voteID)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "GET_TEAMS_ERROR", "Failed to get teams")
		return
	}

	successResponse(w, http.StatusOK, teams)
}

// GetTeamByID handles GET /api/v1/teams/:team_id
func (h *VoteTeamHandler) GetTeamByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	teamID := chi.URLParam(r, "team_id")

	team, err := h.teamService.GetTeamByID(ctx, teamID)
	if err != nil {
		errorResponse(w, http.StatusNotFound, "TEAM_NOT_FOUND", "Team not found")
		return
	}

	successResponse(w, http.StatusOK, team)
}

// UpdateTeam handles PUT /api/v1/teams/:team_id
func (h *VoteTeamHandler) UpdateTeam(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	teamID := chi.URLParam(r, "team_id")

	// Get user ID from context
	userID, ok := ctx.Value(contextkeys.UserIDKey).(string)
	if !ok || userID == "" {
		errorResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	// Parse request
	var req models.UpdateVoteTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	// Update team
	team, err := h.teamService.UpdateTeam(ctx, userID, teamID, &req)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "UPDATE_TEAM_ERROR", err.Error())
		return
	}

	successResponse(w, http.StatusOK, team)
}

// DeleteTeam handles DELETE /api/v1/teams/:team_id
func (h *VoteTeamHandler) DeleteTeam(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	teamID := chi.URLParam(r, "team_id")

	// Get user ID from context
	userID, ok := ctx.Value(contextkeys.UserIDKey).(string)
	if !ok || userID == "" {
		errorResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	// Delete team
	err := h.teamService.DeleteTeam(ctx, userID, teamID)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "DELETE_TEAM_ERROR", err.Error())
		return
	}

	successResponse(w, http.StatusOK, map[string]string{"message": "Team deleted successfully"})
}

// AddPlayerToTeam handles POST /api/v1/teams/:team_id/players
func (h *VoteTeamHandler) AddPlayerToTeam(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	teamID := chi.URLParam(r, "team_id")

	// Get user ID from context
	userID, ok := ctx.Value(contextkeys.UserIDKey).(string)
	if !ok || userID == "" {
		errorResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	// Parse request
	var req models.AddPlayerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	// Add player
	err := h.teamService.AddPlayerToTeam(ctx, userID, teamID, &req)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "ADD_PLAYER_ERROR", err.Error())
		return
	}

	successResponse(w, http.StatusOK, map[string]string{"message": "Player added successfully"})
}

// AddPlayersToTeam handles POST /api/v1/teams/:team_id/players/bulk
func (h *VoteTeamHandler) AddPlayersToTeam(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	teamID := chi.URLParam(r, "team_id")

	// Get user ID from context
	userID, ok := ctx.Value(contextkeys.UserIDKey).(string)
	if !ok || userID == "" {
		errorResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	// Parse request
	var req models.TeamAssignmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	// Add players
	err := h.teamService.AddPlayersToTeam(ctx, userID, teamID, &req)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "ADD_PLAYERS_ERROR", err.Error())
		return
	}

	successResponse(w, http.StatusOK, map[string]string{"message": "Players added successfully"})
}

// RemovePlayerFromTeam handles DELETE /api/v1/teams/:team_id/players/:player_id
func (h *VoteTeamHandler) RemovePlayerFromTeam(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	teamID := chi.URLParam(r, "team_id")
	playerID := chi.URLParam(r, "player_id")

	// Get user ID from context
	userID, ok := ctx.Value(contextkeys.UserIDKey).(string)
	if !ok || userID == "" {
		errorResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	// Remove player
	err := h.teamService.RemovePlayerFromTeam(ctx, userID, teamID, playerID)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "REMOVE_PLAYER_ERROR", err.Error())
		return
	}

	successResponse(w, http.StatusOK, map[string]string{"message": "Player removed successfully"})
}

// GetTeamPlayers handles GET /api/v1/teams/:team_id/players
func (h *VoteTeamHandler) GetTeamPlayers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	teamID := chi.URLParam(r, "team_id")

	players, err := h.teamService.GetTeamPlayers(ctx, teamID)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "GET_PLAYERS_ERROR", "Failed to get team players")
		return
	}

	successResponse(w, http.StatusOK, players)
}
