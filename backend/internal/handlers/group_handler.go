package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/services"
	"spark-park-cricket-backend/internal/utils"

	"github.com/go-chi/chi/v5"
)

type GroupHandler struct {
	groupService services.GroupService
}

// NewGroupHandler creates a new group handler
func NewGroupHandler(groupService services.GroupService) *GroupHandler {
	return &GroupHandler{
		groupService: groupService,
	}
}

// CreateGroup creates a new group
func (h *GroupHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	var req models.CreateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Validation failed", err)
		return
	}

	// Get user ID from context
	userID, ok := utils.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	group, err := h.groupService.CreateGroup(r.Context(), &req, userID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create group", err)
		return
	}

	utils.WriteSuccessResponse(w, http.StatusCreated, "Group created successfully", group)
}

// GetGroup retrieves a group by ID
func (h *GroupHandler) GetGroup(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	if groupID == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Group ID is required", nil)
		return
	}

	group, err := h.groupService.GetGroup(r.Context(), groupID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusNotFound, "Group not found", err)
		return
	}

	utils.WriteSuccessResponse(w, http.StatusOK, "Group retrieved successfully", group)
}

// GetGroupWithMembers retrieves a group with its members
func (h *GroupHandler) GetGroupWithMembers(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	if groupID == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Group ID is required", nil)
		return
	}

	group, err := h.groupService.GetGroupWithMembers(r.Context(), groupID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusNotFound, "Group not found", err)
		return
	}

	utils.WriteSuccessResponse(w, http.StatusOK, "Group with members retrieved successfully", group)
}

// ListGroups retrieves a list of groups with pagination
func (h *GroupHandler) ListGroups(w http.ResponseWriter, r *http.Request) {
	limit, offset := h.parsePaginationParams(r)

	groups, err := h.groupService.ListGroups(r.Context(), limit, offset)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to list groups", err)
		return
	}

	utils.WriteSuccessResponse(w, http.StatusOK, "Groups retrieved successfully", groups)
}

// ListGroupsByType retrieves groups filtered by type
func (h *GroupHandler) ListGroupsByType(w http.ResponseWriter, r *http.Request) {
	groupType := models.GroupType(chi.URLParam(r, "type"))
	if groupType == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Group type is required", nil)
		return
	}

	limit, offset := h.parsePaginationParams(r)

	groups, err := h.groupService.ListGroupsByType(r.Context(), groupType, limit, offset)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to list groups by type", err)
		return
	}

	utils.WriteSuccessResponse(w, http.StatusOK, "Groups retrieved successfully", groups)
}

// ListGroupsByCreator retrieves groups created by a specific user
func (h *GroupHandler) ListGroupsByCreator(w http.ResponseWriter, r *http.Request) {
	creatorID := chi.URLParam(r, "creator_id")
	if creatorID == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Creator ID is required", nil)
		return
	}

	limit, offset := h.parsePaginationParams(r)

	groups, err := h.groupService.ListGroupsByCreator(r.Context(), creatorID, limit, offset)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to list groups by creator", err)
		return
	}

	utils.WriteSuccessResponse(w, http.StatusOK, "Groups retrieved successfully", groups)
}

// ListGroupsByUser retrieves groups that a user is a member of
func (h *GroupHandler) ListGroupsByUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "user_id")
	if userID == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "User ID is required", nil)
		return
	}

	limit, offset := h.parsePaginationParams(r)

	groups, err := h.groupService.ListGroupsByUser(r.Context(), userID, limit, offset)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to list groups by user", err)
		return
	}

	utils.WriteSuccessResponse(w, http.StatusOK, "Groups retrieved successfully", groups)
}

// UpdateGroup updates an existing group
func (h *GroupHandler) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	if groupID == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Group ID is required", nil)
		return
	}

	var req models.UpdateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Validation failed", err)
		return
	}

	// Get user ID from context
	userID, ok := utils.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	group, err := h.groupService.UpdateGroup(r.Context(), groupID, &req, userID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to update group", err)
		return
	}

	utils.WriteSuccessResponse(w, http.StatusOK, "Group updated successfully", group)
}

// DeleteGroup deletes a group
func (h *GroupHandler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	if groupID == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Group ID is required", nil)
		return
	}

	// Get user ID from context
	userID, ok := utils.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	err := h.groupService.DeleteGroup(r.Context(), groupID, userID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to delete group", err)
		return
	}

	utils.WriteSuccessResponse(w, http.StatusOK, "Group deleted successfully", nil)
}

// AddGroupMember adds a user to a group
func (h *GroupHandler) AddGroupMember(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	if groupID == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Group ID is required", nil)
		return
	}

	var req models.AddGroupMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Validation failed", err)
		return
	}

	// Get user ID from context
	userID, ok := utils.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	member, err := h.groupService.AddGroupMember(r.Context(), groupID, &req, userID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to add group member", err)
		return
	}

	utils.WriteSuccessResponse(w, http.StatusCreated, "Group member added successfully", member)
}

// RemoveGroupMember removes a user from a group
func (h *GroupHandler) RemoveGroupMember(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	if groupID == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Group ID is required", nil)
		return
	}

	var req models.RemoveGroupMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Validation failed", err)
		return
	}

	// Get user ID from context
	userID, ok := utils.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	err := h.groupService.RemoveGroupMember(r.Context(), groupID, &req, userID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to remove group member", err)
		return
	}

	utils.WriteSuccessResponse(w, http.StatusOK, "Group member removed successfully", nil)
}

// GetGroupMembers retrieves all members of a group
func (h *GroupHandler) GetGroupMembers(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	if groupID == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Group ID is required", nil)
		return
	}

	members, err := h.groupService.GetGroupMembers(r.Context(), groupID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get group members", err)
		return
	}

	utils.WriteSuccessResponse(w, http.StatusOK, "Group members retrieved successfully", members)
}

// UpdateGroupMemberRole updates a group member's role
func (h *GroupHandler) UpdateGroupMemberRole(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	if groupID == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Group ID is required", nil)
		return
	}

	memberUserID := chi.URLParam(r, "user_id")
	if memberUserID == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "User ID is required", nil)
		return
	}

	var req struct {
		Role string `json:"role" validate:"required,oneof=member admin moderator"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Validation failed", err)
		return
	}

	// Get user ID from context
	userID, ok := utils.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	member, err := h.groupService.UpdateGroupMemberRole(r.Context(), groupID, memberUserID, req.Role, userID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to update group member role", err)
		return
	}

	utils.WriteSuccessResponse(w, http.StatusOK, "Group member role updated successfully", member)
}

// JoinGroup allows a user to join a group
func (h *GroupHandler) JoinGroup(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	if groupID == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Group ID is required", nil)
		return
	}

	// Get user ID from context
	userID, ok := utils.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	member, err := h.groupService.JoinGroup(r.Context(), groupID, userID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to join group", err)
		return
	}

	utils.WriteSuccessResponse(w, http.StatusCreated, "Joined group successfully", member)
}

// LeaveGroup allows a user to leave a group
func (h *GroupHandler) LeaveGroup(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	if groupID == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Group ID is required", nil)
		return
	}

	// Get user ID from context
	userID, ok := utils.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	err := h.groupService.LeaveGroup(r.Context(), groupID, userID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to leave group", err)
		return
	}

	utils.WriteSuccessResponse(w, http.StatusOK, "Left group successfully", nil)
}

// AssignGroupsToVote assigns groups to a vote
func (h *GroupHandler) AssignGroupsToVote(w http.ResponseWriter, r *http.Request) {
	voteID := chi.URLParam(r, "vote_id")
	if voteID == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Vote ID is required", nil)
		return
	}

	var req models.AssignGroupsToVoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Validation failed", err)
		return
	}

	// Get user ID from context
	userID, ok := utils.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	err := h.groupService.AssignGroupsToVote(r.Context(), voteID, &req, userID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to assign groups to vote", err)
		return
	}

	utils.WriteSuccessResponse(w, http.StatusOK, "Groups assigned to vote successfully", nil)
}

// RemoveGroupsFromVote removes groups from a vote
func (h *GroupHandler) RemoveGroupsFromVote(w http.ResponseWriter, r *http.Request) {
	voteID := chi.URLParam(r, "vote_id")
	if voteID == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Vote ID is required", nil)
		return
	}

	var req struct {
		GroupIDs []string `json:"group_ids" validate:"required,min=1,dive,uuid"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Validation failed", err)
		return
	}

	// Get user ID from context
	userID, ok := utils.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	err := h.groupService.RemoveGroupsFromVote(r.Context(), voteID, req.GroupIDs, userID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to remove groups from vote", err)
		return
	}

	utils.WriteSuccessResponse(w, http.StatusOK, "Groups removed from vote successfully", nil)
}

// GetVoteGroups retrieves all groups assigned to a vote
func (h *GroupHandler) GetVoteGroups(w http.ResponseWriter, r *http.Request) {
	voteID := chi.URLParam(r, "vote_id")
	if voteID == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Vote ID is required", nil)
		return
	}

	groups, err := h.groupService.GetVoteGroups(r.Context(), voteID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get vote groups", err)
		return
	}

	utils.WriteSuccessResponse(w, http.StatusOK, "Vote groups retrieved successfully", groups)
}

// GetGroupVotes retrieves all votes assigned to a group
func (h *GroupHandler) GetGroupVotes(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	if groupID == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Group ID is required", nil)
		return
	}

	votes, err := h.groupService.GetGroupVotes(r.Context(), groupID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get group votes", err)
		return
	}

	utils.WriteSuccessResponse(w, http.StatusOK, "Group votes retrieved successfully", votes)
}

// GetGroupVoteResults retrieves voting results for a specific group
func (h *GroupHandler) GetGroupVoteResults(w http.ResponseWriter, r *http.Request) {
	voteID := chi.URLParam(r, "vote_id")
	if voteID == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Vote ID is required", nil)
		return
	}

	groupID := chi.URLParam(r, "group_id")
	if groupID == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Group ID is required", nil)
		return
	}

	results, err := h.groupService.GetGroupVoteResults(r.Context(), voteID, groupID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get group vote results", err)
		return
	}

	utils.WriteSuccessResponse(w, http.StatusOK, "Group vote results retrieved successfully", results)
}

// GetVoteResultsByGroups retrieves voting results broken down by groups
func (h *GroupHandler) GetVoteResultsByGroups(w http.ResponseWriter, r *http.Request) {
	voteID := chi.URLParam(r, "vote_id")
	if voteID == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Vote ID is required", nil)
		return
	}

	results, err := h.groupService.GetVoteResultsByGroups(r.Context(), voteID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get vote results by groups", err)
		return
	}

	utils.WriteSuccessResponse(w, http.StatusOK, "Vote results by groups retrieved successfully", results)
}

// GetVoteWithGroupResults retrieves a vote with results broken down by groups
func (h *GroupHandler) GetVoteWithGroupResults(w http.ResponseWriter, r *http.Request) {
	voteID := chi.URLParam(r, "vote_id")
	if voteID == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Vote ID is required", nil)
		return
	}

	results, err := h.groupService.GetVoteWithGroupResults(r.Context(), voteID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get vote with group results", err)
		return
	}

	utils.WriteSuccessResponse(w, http.StatusOK, "Vote with group results retrieved successfully", results)
}

// SearchGroups searches for groups by name or description
func (h *GroupHandler) SearchGroups(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Search query is required", nil)
		return
	}

	limit, offset := h.parsePaginationParams(r)

	groups, err := h.groupService.SearchGroups(r.Context(), query, limit, offset)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to search groups", err)
		return
	}

	utils.WriteSuccessResponse(w, http.StatusOK, "Groups search completed successfully", groups)
}

// GetGroupStats retrieves statistics for a group
func (h *GroupHandler) GetGroupStats(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "id")
	if groupID == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Group ID is required", nil)
		return
	}

	stats, err := h.groupService.GetGroupStats(r.Context(), groupID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get group stats", err)
		return
	}

	utils.WriteSuccessResponse(w, http.StatusOK, "Group stats retrieved successfully", stats)
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
