package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	contextkeys "spark-park-cricket-backend/internal/context"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/services"
	"spark-park-cricket-backend/internal/utils"

	"github.com/go-chi/chi/v5"
)

// VoteHandler handles voting-related HTTP requests
type VoteHandler struct {
	voteService services.VoteServiceInterface
}

// NewVoteHandler creates a new vote handler
func NewVoteHandler(voteService services.VoteServiceInterface) *VoteHandler {
	return &VoteHandler{
		voteService: voteService,
	}
}

// CreateVote creates a new vote
func (h *VoteHandler) CreateVote(w http.ResponseWriter, r *http.Request) {
	var req models.CreateVoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON", nil)
		return
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	// Get user ID from context
	userID := r.Context().Value(contextkeys.UserIDKey).(string)

	// Create vote
	vote, err := h.voteService.CreateVote(r.Context(), &req, userID)
	if err != nil {
		utils.LogError(err, "Failed to create vote", map[string]interface{}{
			"user_id": userID,
			"title":   req.Title,
		})
		utils.WriteError(w, http.StatusInternalServerError, "CREATE_ERROR", "Failed to create vote", nil)
		return
	}

	response := models.VoteResponse{
		Message: "Vote created successfully",
		Vote:    vote,
	}
	utils.WriteSuccess(w, response)
}

// GetVote retrieves a vote with its options
func (h *VoteHandler) GetVote(w http.ResponseWriter, r *http.Request) {
	voteID := chi.URLParam(r, "id")
	if voteID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Vote ID is required", nil)
		return
	}

	// Get vote
	voteWithOptions, err := h.voteService.GetVote(r.Context(), voteID)
	if err != nil {
		utils.LogError(err, "Failed to get vote", map[string]interface{}{
			"vote_id": voteID,
		})
		utils.WriteError(w, http.StatusNotFound, "NOT_FOUND", "Vote not found", nil)
		return
	}

	utils.WriteSuccess(w, voteWithOptions)
}

// GetVoteWithResults retrieves a vote with results and user's vote
func (h *VoteHandler) GetVoteWithResults(w http.ResponseWriter, r *http.Request) {
	voteID := chi.URLParam(r, "id")
	if voteID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Vote ID is required", nil)
		return
	}

	// Get user ID from context (optional for public results)
	var userID string
	if userIDVal := r.Context().Value(contextkeys.UserIDKey); userIDVal != nil {
		userID = userIDVal.(string)
	}

	// Get vote with results
	voteWithResults, err := h.voteService.GetVoteWithResults(r.Context(), voteID, userID)
	if err != nil {
		utils.LogError(err, "Failed to get vote with results", map[string]interface{}{
			"vote_id": voteID,
			"user_id": userID,
		})
		utils.WriteError(w, http.StatusNotFound, "NOT_FOUND", "Vote not found", nil)
		return
	}

	utils.WriteSuccess(w, voteWithResults)
}

// UpdateVote updates an existing vote
func (h *VoteHandler) UpdateVote(w http.ResponseWriter, r *http.Request) {
	voteID := chi.URLParam(r, "id")
	if voteID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Vote ID is required", nil)
		return
	}

	var req models.UpdateVoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON", nil)
		return
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	// Get user ID from context
	userID := r.Context().Value(contextkeys.UserIDKey).(string)

	// Update vote
	vote, err := h.voteService.UpdateVote(r.Context(), voteID, &req, userID)
	if err != nil {
		utils.LogError(err, "Failed to update vote", map[string]interface{}{
			"vote_id": voteID,
			"user_id": userID,
		})
		if err.Error() == "unauthorized: only vote creator can update vote" {
			utils.WriteError(w, http.StatusForbidden, "FORBIDDEN", "Only vote creator can update vote", nil)
			return
		}
		if err.Error() == "cannot update closed or cancelled vote" {
			utils.WriteError(w, http.StatusBadRequest, "INVALID_STATE", "Cannot update closed or cancelled vote", nil)
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "UPDATE_ERROR", "Failed to update vote", nil)
		return
	}

	response := models.VoteResponse{
		Message: "Vote updated successfully",
		Vote:    vote,
	}
	utils.WriteSuccess(w, response)
}

// DeleteVote deletes a vote
func (h *VoteHandler) DeleteVote(w http.ResponseWriter, r *http.Request) {
	voteID := chi.URLParam(r, "id")
	if voteID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Vote ID is required", nil)
		return
	}

	// Get user ID from context
	userID := r.Context().Value(contextkeys.UserIDKey).(string)

	// Delete vote
	err := h.voteService.DeleteVote(r.Context(), voteID, userID)
	if err != nil {
		utils.LogError(err, "Failed to delete vote", map[string]interface{}{
			"vote_id": voteID,
			"user_id": userID,
		})
		if err.Error() == "unauthorized: only vote creator can delete vote" {
			utils.WriteError(w, http.StatusForbidden, "FORBIDDEN", "Only vote creator can delete vote", nil)
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "DELETE_ERROR", "Failed to delete vote", nil)
		return
	}

	response := models.VoteResponse{
		Message: "Vote deleted successfully",
	}
	utils.WriteSuccess(w, response)
}

// ListVotes lists votes with filters
func (h *VoteHandler) ListVotes(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	filters := &models.VoteFilters{
		Limit:  20, // Default limit
		Offset: 0,  // Default offset
	}

	// Parse status filter
	if status := r.URL.Query().Get("status"); status != "" {
		voteStatus := models.VoteStatus(status)
		filters.Status = &voteStatus
	}

	// Parse type filter
	if voteType := r.URL.Query().Get("type"); voteType != "" {
		voteTypeEnum := models.VoteType(voteType)
		filters.Type = &voteTypeEnum
	}

	// Parse created_by filter
	if createdBy := r.URL.Query().Get("created_by"); createdBy != "" {
		filters.CreatedBy = &createdBy
	}

	// Parse limit
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 && limit <= 100 {
			filters.Limit = limit
		}
	}

	// Parse offset
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			filters.Offset = offset
		}
	}

	// List votes
	votes, err := h.voteService.ListVotes(r.Context(), filters)
	if err != nil {
		utils.LogError(err, "Failed to list votes", nil)
		utils.WriteError(w, http.StatusInternalServerError, "LIST_ERROR", "Failed to list votes", nil)
		return
	}

	utils.WriteSuccess(w, votes)
}

// CastVote allows a user to cast their vote
func (h *VoteHandler) CastVote(w http.ResponseWriter, r *http.Request) {
	voteID := chi.URLParam(r, "id")
	if voteID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Vote ID is required", nil)
		return
	}

	var req models.VoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON", nil)
		return
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	// Get user ID from context
	userID := r.Context().Value(contextkeys.UserIDKey).(string)

	// Cast vote
	err := h.voteService.CastVote(r.Context(), voteID, &req, userID)
	if err != nil {
		utils.LogError(err, "Failed to cast vote", map[string]interface{}{
			"vote_id": voteID,
			"user_id": userID,
		})
		if err.Error() == "cannot vote on closed or cancelled vote" {
			utils.WriteError(w, http.StatusBadRequest, "INVALID_STATE", "Cannot vote on closed or cancelled vote", nil)
			return
		}
		if err.Error() == "user has already voted on this poll" {
			utils.WriteError(w, http.StatusConflict, "ALREADY_VOTED", "User has already voted on this poll", nil)
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "VOTE_ERROR", "Failed to cast vote", nil)
		return
	}

	response := models.VoteResponse{
		Message: "Vote cast successfully",
	}
	utils.WriteSuccess(w, response)
}

// GetUserVote retrieves a user's vote for a specific vote
func (h *VoteHandler) GetUserVote(w http.ResponseWriter, r *http.Request) {
	voteID := chi.URLParam(r, "id")
	if voteID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Vote ID is required", nil)
		return
	}

	// Get user ID from context
	userID := r.Context().Value(contextkeys.UserIDKey).(string)

	// Get user vote
	userVote, err := h.voteService.GetUserVote(r.Context(), voteID, userID)
	if err != nil {
		utils.LogError(err, "Failed to get user vote", map[string]interface{}{
			"vote_id": voteID,
			"user_id": userID,
		})
		if err.Error() == "user vote not found" {
			utils.WriteError(w, http.StatusNotFound, "NOT_FOUND", "User vote not found", nil)
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "GET_ERROR", "Failed to get user vote", nil)
		return
	}

	utils.WriteSuccess(w, userVote)
}

// HasUserVoted checks if a user has voted for a specific vote
func (h *VoteHandler) HasUserVoted(w http.ResponseWriter, r *http.Request) {
	voteID := chi.URLParam(r, "id")
	if voteID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Vote ID is required", nil)
		return
	}

	// Get user ID from context
	userID := r.Context().Value(contextkeys.UserIDKey).(string)

	// Check if user has voted
	hasVoted, err := h.voteService.HasUserVoted(r.Context(), voteID, userID)
	if err != nil {
		utils.LogError(err, "Failed to check if user voted", map[string]interface{}{
			"vote_id": voteID,
			"user_id": userID,
		})
		utils.WriteError(w, http.StatusInternalServerError, "CHECK_ERROR", "Failed to check if user voted", nil)
		return
	}

	response := map[string]interface{}{
		"has_voted": hasVoted,
	}
	utils.WriteSuccess(w, response)
}

// CloseVote closes a vote
func (h *VoteHandler) CloseVote(w http.ResponseWriter, r *http.Request) {
	voteID := chi.URLParam(r, "id")
	if voteID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Vote ID is required", nil)
		return
	}

	// Get user ID from context
	userID := r.Context().Value(contextkeys.UserIDKey).(string)

	// Close vote
	err := h.voteService.CloseVote(r.Context(), voteID, userID)
	if err != nil {
		utils.LogError(err, "Failed to close vote", map[string]interface{}{
			"vote_id": voteID,
			"user_id": userID,
		})
		if err.Error() == "unauthorized: only vote creator can update vote" {
			utils.WriteError(w, http.StatusForbidden, "FORBIDDEN", "Only vote creator can close vote", nil)
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "CLOSE_ERROR", "Failed to close vote", nil)
		return
	}

	response := models.VoteResponse{
		Message: "Vote closed successfully",
	}
	utils.WriteSuccess(w, response)
}

// CancelVote cancels a vote
func (h *VoteHandler) CancelVote(w http.ResponseWriter, r *http.Request) {
	voteID := chi.URLParam(r, "id")
	if voteID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Vote ID is required", nil)
		return
	}

	// Get user ID from context
	userID := r.Context().Value(contextkeys.UserIDKey).(string)

	// Cancel vote
	err := h.voteService.CancelVote(r.Context(), voteID, userID)
	if err != nil {
		utils.LogError(err, "Failed to cancel vote", map[string]interface{}{
			"vote_id": voteID,
			"user_id": userID,
		})
		if err.Error() == "unauthorized: only vote creator can update vote" {
			utils.WriteError(w, http.StatusForbidden, "FORBIDDEN", "Only vote creator can cancel vote", nil)
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "CANCEL_ERROR", "Failed to cancel vote", nil)
		return
	}

	response := models.VoteResponse{
		Message: "Vote cancelled successfully",
	}
	utils.WriteSuccess(w, response)
}
