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

type GroupHandler struct {
	groupService services.GroupServiceInterface
}

// NewGroupHandler creates a new group handler
func NewGroupHandler(groupService services.GroupServiceInterface) *GroupHandler {
	return &GroupHandler{
		groupService: groupService,
	}
}

// CreateGroup creates a new group
func (h *GroupHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	var req models.CreateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", err.Error())
		return
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", err.Error())
		return
	}

	// Get user ID from context
	userID, ok := r.Context().Value(contextkeys.UserIDKey).(string)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	group, err := h.groupService.CreateGroup(r.Context(), &req, userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create group", err.Error())
		return
	}

	utils.WriteCreated(w, group)
}

// GetGroup retrieves a group by ID
func (h *GroupHandler) GetGroup(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	if groupID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Group ID is required", nil)
		return
	}

	group, err := h.groupService.GetGroup(r.Context(), groupID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "NOT_FOUND", "Group not found", err.Error())
		return
	}

	utils.WriteSuccess(w, group)
}

// GetGroupWithMembers retrieves a group with its members
func (h *GroupHandler) GetGroupWithMembers(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	if groupID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Group ID is required", nil)
		return
	}

	group, err := h.groupService.GetGroupWithMembers(r.Context(), groupID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "NOT_FOUND", "Group not found", err.Error())
		return
	}

	utils.WriteSuccess(w, group)
}

// ListGroups retrieves a list of groups with pagination
func (h *GroupHandler) ListGroups(w http.ResponseWriter, r *http.Request) {
	limit, offset := h.parsePaginationParams(r)

	groups, err := h.groupService.ListGroups(r.Context(), limit, offset)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list groups", err.Error())
		return
	}

	utils.WriteSuccess(w, groups)
}

// ListGroupsByType retrieves groups filtered by type
func (h *GroupHandler) ListGroupsByType(w http.ResponseWriter, r *http.Request) {
	groupType := models.GroupType(chi.URLParam(r, "type"))
	if groupType == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Group type is required", nil)
		return
	}

	limit, offset := h.parsePaginationParams(r)

	groups, err := h.groupService.ListGroupsByType(r.Context(), groupType, limit, offset)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list groups by type", err.Error())
		return
	}

	utils.WriteSuccess(w, groups)
}

// ListGroupsByCreator retrieves groups created by a specific user
func (h *GroupHandler) ListGroupsByCreator(w http.ResponseWriter, r *http.Request) {
	creatorID := chi.URLParam(r, "creator_id")
	if creatorID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Creator ID is required", nil)
		return
	}

	limit, offset := h.parsePaginationParams(r)

	groups, err := h.groupService.ListGroupsByCreator(r.Context(), creatorID, limit, offset)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list groups by creator", err.Error())
		return
	}

	utils.WriteSuccess(w, groups)
}

// ListGroupsByUser retrieves groups that a user is a member of
func (h *GroupHandler) ListGroupsByUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "user_id")
	if userID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "User ID is required", nil)
		return
	}

	limit, offset := h.parsePaginationParams(r)

	groups, err := h.groupService.ListGroupsByUser(r.Context(), userID, limit, offset)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to list groups by user", err.Error())
		return
	}

	utils.WriteSuccess(w, groups)
}

// UpdateGroup updates an existing group
func (h *GroupHandler) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	if groupID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Group ID is required", nil)
		return
	}

	var req models.UpdateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", err.Error())
		return
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", err.Error())
		return
	}

	// Get user ID from context
	userID, ok := r.Context().Value(contextkeys.UserIDKey).(string)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	group, err := h.groupService.UpdateGroup(r.Context(), groupID, &req, userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update group", err.Error())
		return
	}

	utils.WriteSuccess(w, group)
}

// DeleteGroup deletes a group
func (h *GroupHandler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	if groupID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Group ID is required", nil)
		return
	}

	// Get user ID from context
	userID, ok := r.Context().Value(contextkeys.UserIDKey).(string)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	err := h.groupService.DeleteGroup(r.Context(), groupID, userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete group", err.Error())
		return
	}

	utils.WriteSuccess(w, map[string]string{"message": "Group deleted successfully"})
}

// AddGroupMember adds a user to a group
func (h *GroupHandler) AddGroupMember(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	if groupID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Group ID is required", nil)
		return
	}

	var req models.AddGroupMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", err.Error())
		return
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", err.Error())
		return
	}

	// Get user ID from context
	userID, ok := r.Context().Value(contextkeys.UserIDKey).(string)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	member, err := h.groupService.AddGroupMember(r.Context(), groupID, &req, userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to add group member", err.Error())
		return
	}

	utils.WriteCreated(w, member)
}

// RemoveGroupMember removes a user from a group
func (h *GroupHandler) RemoveGroupMember(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	if groupID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Group ID is required", nil)
		return
	}

	var req models.RemoveGroupMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", err.Error())
		return
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", err.Error())
		return
	}

	// Get user ID from context
	userID, ok := r.Context().Value(contextkeys.UserIDKey).(string)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	err := h.groupService.RemoveGroupMember(r.Context(), groupID, &req, userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to remove group member", err.Error())
		return
	}

	utils.WriteSuccess(w, map[string]string{"message": "Group member removed successfully"})
}

// GetGroupMembers retrieves all members of a group
func (h *GroupHandler) GetGroupMembers(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	if groupID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Group ID is required", nil)
		return
	}

	members, err := h.groupService.GetGroupMembers(r.Context(), groupID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get group members", err.Error())
		return
	}

	utils.WriteSuccess(w, members)
}

// UpdateGroupMemberRole updates a group member's role
func (h *GroupHandler) UpdateGroupMemberRole(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	if groupID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Group ID is required", nil)
		return
	}

	memberUserID := chi.URLParam(r, "user_id")
	if memberUserID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "User ID is required", nil)
		return
	}

	var req struct {
		Role string `json:"role" validate:"required,oneof=member admin moderator"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", err.Error())
		return
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", err.Error())
		return
	}

	// Get user ID from context
	userID, ok := r.Context().Value(contextkeys.UserIDKey).(string)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	member, err := h.groupService.UpdateGroupMemberRole(r.Context(), groupID, memberUserID, req.Role, userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update group member role", err.Error())
		return
	}

	utils.WriteSuccess(w, member)
}

// JoinGroup allows a user to join a group
func (h *GroupHandler) JoinGroup(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	if groupID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Group ID is required", nil)
		return
	}

	// Get user ID from context
	userID, ok := r.Context().Value(contextkeys.UserIDKey).(string)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	member, err := h.groupService.JoinGroup(r.Context(), groupID, userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to join group", err.Error())
		return
	}

	utils.WriteCreated(w, member)
}

// LeaveGroup allows a user to leave a group
func (h *GroupHandler) LeaveGroup(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	if groupID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Group ID is required", nil)
		return
	}

	// Get user ID from context
	userID, ok := r.Context().Value(contextkeys.UserIDKey).(string)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	err := h.groupService.LeaveGroup(r.Context(), groupID, userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to leave group", err.Error())
		return
	}

	utils.WriteSuccess(w, map[string]string{"message": "Left group successfully"})
}

// AssignGroupsToVote assigns groups to a vote
func (h *GroupHandler) AssignGroupsToVote(w http.ResponseWriter, r *http.Request) {
	voteID := chi.URLParam(r, "vote_id")
	if voteID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Vote ID is required", nil)
		return
	}

	var req models.AssignGroupsToVoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", err.Error())
		return
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", err.Error())
		return
	}

	// Get user ID from context
	userID, ok := r.Context().Value(contextkeys.UserIDKey).(string)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	err := h.groupService.AssignGroupsToVote(r.Context(), voteID, &req, userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to assign groups to vote", err.Error())
		return
	}

	utils.WriteSuccess(w, map[string]string{"message": "Groups assigned to vote successfully"})
}

// RemoveGroupsFromVote removes groups from a vote
func (h *GroupHandler) RemoveGroupsFromVote(w http.ResponseWriter, r *http.Request) {
	voteID := chi.URLParam(r, "vote_id")
	if voteID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Vote ID is required", nil)
		return
	}

	var req struct {
		GroupIDs []string `json:"group_ids" validate:"required,min=1,dive,uuid"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", err.Error())
		return
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", err.Error())
		return
	}

	// Get user ID from context
	userID, ok := r.Context().Value(contextkeys.UserIDKey).(string)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	err := h.groupService.RemoveGroupsFromVote(r.Context(), voteID, req.GroupIDs, userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to remove groups from vote", err.Error())
		return
	}

	utils.WriteSuccess(w, map[string]string{"message": "Groups removed from vote successfully"})
}

// GetVoteGroups retrieves all groups assigned to a vote
func (h *GroupHandler) GetVoteGroups(w http.ResponseWriter, r *http.Request) {
	voteID := chi.URLParam(r, "vote_id")
	if voteID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Vote ID is required", nil)
		return
	}

	groups, err := h.groupService.GetVoteGroups(r.Context(), voteID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get vote groups", err.Error())
		return
	}

	utils.WriteSuccess(w, groups)
}

// GetGroupVotes retrieves all votes assigned to a group
func (h *GroupHandler) GetGroupVotes(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	if groupID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Group ID is required", nil)
		return
	}

	votes, err := h.groupService.GetGroupVotes(r.Context(), groupID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get group votes", err.Error())
		return
	}

	utils.WriteSuccess(w, votes)
}

// GetGroupVoteResults retrieves voting results for a specific group
func (h *GroupHandler) GetGroupVoteResults(w http.ResponseWriter, r *http.Request) {
	voteID := chi.URLParam(r, "vote_id")
	if voteID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Vote ID is required", nil)
		return
	}

	groupID := chi.URLParam(r, "group_id")
	if groupID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Group ID is required", nil)
		return
	}

	results, err := h.groupService.GetGroupVoteResults(r.Context(), voteID, groupID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get group vote results", err.Error())
		return
	}

	utils.WriteSuccess(w, results)
}

// GetVoteResultsByGroups retrieves voting results broken down by groups
func (h *GroupHandler) GetVoteResultsByGroups(w http.ResponseWriter, r *http.Request) {
	voteID := chi.URLParam(r, "vote_id")
	if voteID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Vote ID is required", nil)
		return
	}

	results, err := h.groupService.GetVoteResultsByGroups(r.Context(), voteID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get vote results by groups", err.Error())
		return
	}

	utils.WriteSuccess(w, results)
}

// GetVoteWithGroupResults retrieves a vote with results broken down by groups
func (h *GroupHandler) GetVoteWithGroupResults(w http.ResponseWriter, r *http.Request) {
	voteID := chi.URLParam(r, "vote_id")
	if voteID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Vote ID is required", nil)
		return
	}

	results, err := h.groupService.GetVoteWithGroupResults(r.Context(), voteID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get vote with group results", err.Error())
		return
	}

	utils.WriteSuccess(w, results)
}

// SearchGroups searches for groups by name or description
func (h *GroupHandler) SearchGroups(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Search query is required", nil)
		return
	}

	limit, offset := h.parsePaginationParams(r)

	groups, err := h.groupService.SearchGroups(r.Context(), query, limit, offset)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to search groups", err.Error())
		return
	}

	utils.WriteSuccess(w, groups)
}

// GetGroupStats retrieves statistics for a group
func (h *GroupHandler) GetGroupStats(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	if groupID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Group ID is required", nil)
		return
	}

	stats, err := h.groupService.GetGroupStats(r.Context(), groupID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get group stats", err.Error())
		return
	}

	utils.WriteSuccess(w, stats)
}

// parsePaginationParams parses limit and offset from query parameters
func (h *GroupHandler) parsePaginationParams(r *http.Request) (int, int) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 20 // default
	offset := 0 // default

	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 100 {
			limit = parsedLimit
		}
	}

	if offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	return limit, offset
}
