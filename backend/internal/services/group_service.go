package services

import (
	"context"
	"fmt"
	"time"

	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/repository/interfaces"
	"spark-park-cricket-backend/internal/utils"

	"github.com/google/uuid"
)

type GroupService struct {
	groupRepo interfaces.GroupRepository
	userRepo  interfaces.UserRepository
	voteRepo  interfaces.VoteRepositoryInterface
}

// NewGroupService creates a new group service
func NewGroupService(groupRepo interfaces.GroupRepository, userRepo interfaces.UserRepository, voteRepo interfaces.VoteRepositoryInterface) GroupServiceInterface {
	return &GroupService{
		groupRepo: groupRepo,
		userRepo:  userRepo,
		voteRepo:  voteRepo,
	}
}

// CreateGroup creates a new group
func (s *GroupService) CreateGroup(ctx context.Context, req *models.CreateGroupRequest, creatorID string) (*models.Group, error) {
	// Validate creator exists
	_, err := s.userRepo.GetUserByID(ctx, creatorID)
	if err != nil {
		return nil, fmt.Errorf("creator not found: %w", err)
	}

	// Create group
	group := &models.Group{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		Status:      models.GroupStatusActive,
		CreatedBy:   creatorID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.groupRepo.CreateGroup(ctx, group); err != nil {
		utils.LogError(err, "Failed to create group", map[string]interface{}{
			"group_id": group.ID,
			"creator":  creatorID,
		})
		return nil, fmt.Errorf("failed to create group: %w", err)
	}

	// Add creator as admin member
	member := &models.GroupMember{
		ID:        uuid.New().String(),
		GroupID:   group.ID,
		UserID:    creatorID,
		Role:      "admin",
		JoinedAt:  time.Now(),
		CreatedAt: time.Now(),
	}

	if err := s.groupRepo.AddGroupMember(ctx, member); err != nil {
		utils.LogError(err, "Failed to add creator as group member", map[string]interface{}{
			"group_id": group.ID,
			"user_id":  creatorID,
		})
		// Don't fail the group creation, just log the error
	}

	utils.LogInfo("Group created successfully", map[string]interface{}{
		"group_id": group.ID,
		"name":     group.Name,
		"type":     group.Type,
		"creator":  creatorID,
	})

	return group, nil
}

// GetGroup retrieves a group by ID
func (s *GroupService) GetGroup(ctx context.Context, groupID string) (*models.GroupWithCreator, error) {
	group, err := s.groupRepo.GetGroupWithCreator(ctx, groupID)
	if err != nil {
		return nil, err
	}

	return group, nil
}

// GetGroupWithMembers retrieves a group with its members
func (s *GroupService) GetGroupWithMembers(ctx context.Context, groupID string) (*models.GroupWithMembersAndUsers, error) {
	group, err := s.groupRepo.GetGroupWithMembersAndUsers(ctx, groupID)
	if err != nil {
		return nil, err
	}

	return group, nil
}

// ListGroups retrieves a list of groups with pagination
func (s *GroupService) ListGroups(ctx context.Context, limit, offset int) ([]*models.GroupWithCreator, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	groups, err := s.groupRepo.ListGroups(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	return groups, nil
}

// ListGroupsByType retrieves groups filtered by type
func (s *GroupService) ListGroupsByType(ctx context.Context, groupType models.GroupType, limit, offset int) ([]*models.GroupWithCreator, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	groups, err := s.groupRepo.ListGroupsByType(ctx, groupType, limit, offset)
	if err != nil {
		return nil, err
	}

	return groups, nil
}

// ListGroupsByCreator retrieves groups created by a specific user
func (s *GroupService) ListGroupsByCreator(ctx context.Context, creatorID string, limit, offset int) ([]*models.GroupWithCreator, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	groups, err := s.groupRepo.ListGroupsByCreator(ctx, creatorID, limit, offset)
	if err != nil {
		return nil, err
	}

	return groups, nil
}

// ListGroupsByUser retrieves groups that a user is a member of
func (s *GroupService) ListGroupsByUser(ctx context.Context, userID string, limit, offset int) ([]*models.GroupWithCreator, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	groups, err := s.groupRepo.ListGroupsByUser(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	return groups, nil
}

// UpdateGroup updates an existing group
func (s *GroupService) UpdateGroup(ctx context.Context, groupID string, req *models.UpdateGroupRequest, userID string) (*models.Group, error) {
	// Check if user has admin access to the group
	hasAccess, err := s.ValidateGroupAdminAccess(ctx, groupID, userID)
	if err != nil {
		return nil, err
	}
	if !hasAccess {
		return nil, fmt.Errorf("insufficient permissions to update group")
	}

	// Get existing group
	group, err := s.groupRepo.GetGroup(ctx, groupID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != nil {
		group.Name = *req.Name
	}
	if req.Description != nil {
		group.Description = *req.Description
	}
	if req.Status != nil {
		group.Status = *req.Status
	}

	group.UpdatedAt = time.Now()

	if err := s.groupRepo.UpdateGroup(ctx, group); err != nil {
		utils.LogError(err, "Failed to update group", map[string]interface{}{
			"group_id": groupID,
			"user_id":  userID,
		})
		return nil, fmt.Errorf("failed to update group: %w", err)
	}

	utils.LogInfo("Group updated successfully", map[string]interface{}{
		"group_id": groupID,
		"user_id":  userID,
	})

	return &group.Group, nil
}

// DeleteGroup deletes a group
func (s *GroupService) DeleteGroup(ctx context.Context, groupID string, userID string) error {
	// Check if user has admin access to the group
	hasAccess, err := s.ValidateGroupAdminAccess(ctx, groupID, userID)
	if err != nil {
		return err
	}
	if !hasAccess {
		return fmt.Errorf("insufficient permissions to delete group")
	}

	if err := s.groupRepo.DeleteGroup(ctx, groupID); err != nil {
		utils.LogError(err, "Failed to delete group", map[string]interface{}{
			"group_id": groupID,
			"user_id":  userID,
		})
		return fmt.Errorf("failed to delete group: %w", err)
	}

	utils.LogInfo("Group deleted successfully", map[string]interface{}{
		"group_id": groupID,
		"user_id":  userID,
	})

	return nil
}

// AddGroupMember adds a user to a group
func (s *GroupService) AddGroupMember(ctx context.Context, groupID string, req *models.AddGroupMemberRequest, adminUserID string) (*models.GroupMember, error) {
	// Check if admin user has access to the group
	hasAccess, err := s.ValidateGroupAdminAccess(ctx, groupID, adminUserID)
	if err != nil {
		return nil, err
	}
	if !hasAccess {
		return nil, fmt.Errorf("insufficient permissions to add group members")
	}

	// Validate user exists
	_, err = s.userRepo.GetUserByID(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Check if user is already in the group
	isMember, err := s.groupRepo.IsUserInGroup(ctx, groupID, req.UserID)
	if err != nil {
		return nil, err
	}
	if isMember {
		return nil, fmt.Errorf("user is already a member of this group")
	}

	// Set default role if not provided
	role := req.Role
	if role == "" {
		role = "member"
	}

	member := &models.GroupMember{
		ID:        uuid.New().String(),
		GroupID:   groupID,
		UserID:    req.UserID,
		Role:      role,
		JoinedAt:  time.Now(),
		CreatedAt: time.Now(),
	}

	if err := s.groupRepo.AddGroupMember(ctx, member); err != nil {
		utils.LogError(err, "Failed to add group member", map[string]interface{}{
			"group_id": groupID,
			"user_id":  req.UserID,
		})
		return nil, fmt.Errorf("failed to add group member: %w", err)
	}

	utils.LogInfo("Group member added successfully", map[string]interface{}{
		"group_id": groupID,
		"user_id":  req.UserID,
		"role":     role,
	})

	return member, nil
}

// RemoveGroupMember removes a user from a group
func (s *GroupService) RemoveGroupMember(ctx context.Context, groupID string, req *models.RemoveGroupMemberRequest, adminUserID string) error {
	// Check if admin user has access to the group
	hasAccess, err := s.ValidateGroupAdminAccess(ctx, groupID, adminUserID)
	if err != nil {
		return err
	}
	if !hasAccess {
		return fmt.Errorf("insufficient permissions to remove group members")
	}

	// Check if user is in the group
	isMember, err := s.groupRepo.IsUserInGroup(ctx, groupID, req.UserID)
	if err != nil {
		return err
	}
	if !isMember {
		return fmt.Errorf("user is not a member of this group")
	}

	if err := s.groupRepo.RemoveGroupMember(ctx, groupID, req.UserID); err != nil {
		utils.LogError(err, "Failed to remove group member", map[string]interface{}{
			"group_id": groupID,
			"user_id":  req.UserID,
		})
		return fmt.Errorf("failed to remove group member: %w", err)
	}

	utils.LogInfo("Group member removed successfully", map[string]interface{}{
		"group_id": groupID,
		"user_id":  req.UserID,
	})

	return nil
}

// GetGroupMembers retrieves all members of a group
func (s *GroupService) GetGroupMembers(ctx context.Context, groupID string) ([]*models.GroupMemberWithUser, error) {
	members, err := s.groupRepo.GetGroupMembersWithUsers(ctx, groupID)
	if err != nil {
		return nil, err
	}

	return members, nil
}

// UpdateGroupMemberRole updates a group member's role
func (s *GroupService) UpdateGroupMemberRole(ctx context.Context, groupID, userID, role string, adminUserID string) (*models.GroupMember, error) {
	// Check if admin user has access to the group
	hasAccess, err := s.ValidateGroupAdminAccess(ctx, groupID, adminUserID)
	if err != nil {
		return nil, err
	}
	if !hasAccess {
		return nil, fmt.Errorf("insufficient permissions to update group member roles")
	}

	// Validate role
	if role != "member" && role != "admin" && role != "moderator" {
		return nil, fmt.Errorf("invalid role: %s", role)
	}

	if err := s.groupRepo.UpdateGroupMemberRole(ctx, groupID, userID, role); err != nil {
		utils.LogError(err, "Failed to update group member role", map[string]interface{}{
			"group_id": groupID,
			"user_id":  userID,
			"role":     role,
		})
		return nil, fmt.Errorf("failed to update group member role: %w", err)
	}

	// Get updated member
	member, err := s.groupRepo.GetGroupMember(ctx, groupID, userID)
	if err != nil {
		return nil, err
	}

	utils.LogInfo("Group member role updated successfully", map[string]interface{}{
		"group_id": groupID,
		"user_id":  userID,
		"role":     role,
	})

	return member, nil
}

// JoinGroup allows a user to join a group
func (s *GroupService) JoinGroup(ctx context.Context, groupID string, userID string) (*models.GroupMember, error) {
	// Check if user is already in the group
	isMember, err := s.groupRepo.IsUserInGroup(ctx, groupID, userID)
	if err != nil {
		return nil, err
	}
	if isMember {
		return nil, fmt.Errorf("user is already a member of this group")
	}

	member := &models.GroupMember{
		ID:        uuid.New().String(),
		GroupID:   groupID,
		UserID:    userID,
		Role:      "member",
		JoinedAt:  time.Now(),
		CreatedAt: time.Now(),
	}

	if err := s.groupRepo.AddGroupMember(ctx, member); err != nil {
		utils.LogError(err, "Failed to join group", map[string]interface{}{
			"group_id": groupID,
			"user_id":  userID,
		})
		return nil, fmt.Errorf("failed to join group: %w", err)
	}

	utils.LogInfo("User joined group successfully", map[string]interface{}{
		"group_id": groupID,
		"user_id":  userID,
	})

	return member, nil
}

// LeaveGroup allows a user to leave a group
func (s *GroupService) LeaveGroup(ctx context.Context, groupID string, userID string) error {
	// Check if user is in the group
	isMember, err := s.groupRepo.IsUserInGroup(ctx, groupID, userID)
	if err != nil {
		return err
	}
	if !isMember {
		return fmt.Errorf("user is not a member of this group")
	}

	if err := s.groupRepo.RemoveGroupMember(ctx, groupID, userID); err != nil {
		utils.LogError(err, "Failed to leave group", map[string]interface{}{
			"group_id": groupID,
			"user_id":  userID,
		})
		return fmt.Errorf("failed to leave group: %w", err)
	}

	utils.LogInfo("User left group successfully", map[string]interface{}{
		"group_id": groupID,
		"user_id":  userID,
	})

	return nil
}

// AssignGroupsToVote assigns groups to a vote
func (s *GroupService) AssignGroupsToVote(ctx context.Context, voteID string, req *models.AssignGroupsToVoteRequest, userID string) error {
	// Check if user has permission to modify the vote
	vote, err := s.voteRepo.GetVoteByID(ctx, voteID)
	if err != nil {
		return err
	}

	if vote.CreatedBy != userID {
		return fmt.Errorf("insufficient permissions to assign groups to vote")
	}

	// Validate all groups exist
	for _, groupID := range req.GroupIDs {
		_, err := s.groupRepo.GetGroup(ctx, groupID)
		if err != nil {
			return fmt.Errorf("group not found: %s", groupID)
		}
	}

	if err := s.groupRepo.AssignGroupsToVote(ctx, voteID, req.GroupIDs); err != nil {
		utils.LogError(err, "Failed to assign groups to vote", map[string]interface{}{
			"vote_id":   voteID,
			"group_ids": req.GroupIDs,
			"user_id":   userID,
		})
		return fmt.Errorf("failed to assign groups to vote: %w", err)
	}

	utils.LogInfo("Groups assigned to vote successfully", map[string]interface{}{
		"vote_id":   voteID,
		"group_ids": req.GroupIDs,
		"user_id":   userID,
	})

	return nil
}

// RemoveGroupsFromVote removes groups from a vote
func (s *GroupService) RemoveGroupsFromVote(ctx context.Context, voteID string, groupIDs []string, userID string) error {
	// Check if user has permission to modify the vote
	vote, err := s.voteRepo.GetVoteByID(ctx, voteID)
	if err != nil {
		return err
	}

	if vote.CreatedBy != userID {
		return fmt.Errorf("insufficient permissions to remove groups from vote")
	}

	if err := s.groupRepo.RemoveGroupsFromVote(ctx, voteID, groupIDs); err != nil {
		utils.LogError(err, "Failed to remove groups from vote", map[string]interface{}{
			"vote_id":   voteID,
			"group_ids": groupIDs,
			"user_id":   userID,
		})
		return fmt.Errorf("failed to remove groups from vote: %w", err)
	}

	utils.LogInfo("Groups removed from vote successfully", map[string]interface{}{
		"vote_id":   voteID,
		"group_ids": groupIDs,
		"user_id":   userID,
	})

	return nil
}

// GetVoteGroups retrieves all groups assigned to a vote
func (s *GroupService) GetVoteGroups(ctx context.Context, voteID string) ([]*models.Group, error) {
	groups, err := s.groupRepo.GetVoteGroups(ctx, voteID)
	if err != nil {
		return nil, err
	}

	return groups, nil
}

// GetGroupVotes retrieves all votes assigned to a group
func (s *GroupService) GetGroupVotes(ctx context.Context, groupID string) ([]*models.Vote, error) {
	votes, err := s.groupRepo.GetGroupVotes(ctx, groupID)
	if err != nil {
		return nil, err
	}

	return votes, nil
}

// GetGroupVoteResults retrieves voting results for a specific group
func (s *GroupService) GetGroupVoteResults(ctx context.Context, voteID, groupID string) (*models.GroupVoteResult, error) {
	results, err := s.groupRepo.GetGroupVoteResults(ctx, voteID, groupID)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// GetVoteResultsByGroups retrieves voting results broken down by groups
func (s *GroupService) GetVoteResultsByGroups(ctx context.Context, voteID string) ([]*models.GroupVoteResult, error) {
	results, err := s.groupRepo.GetVoteResultsByGroups(ctx, voteID)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// GetVoteWithGroupResults retrieves a vote with results broken down by groups
func (s *GroupService) GetVoteWithGroupResults(ctx context.Context, voteID string) (*models.VoteWithGroupResults, error) {
	results, err := s.groupRepo.GetVoteWithGroupResults(ctx, voteID)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// SearchGroups searches for groups by name or description
func (s *GroupService) SearchGroups(ctx context.Context, query string, limit, offset int) ([]*models.GroupWithCreator, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	groups, err := s.groupRepo.SearchGroups(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	return groups, nil
}

// GetGroupStats retrieves statistics for a group
func (s *GroupService) GetGroupStats(ctx context.Context, groupID string) (map[string]interface{}, error) {
	memberCount, err := s.groupRepo.GetGroupMemberCount(ctx, groupID)
	if err != nil {
		return nil, err
	}

	voteCount, err := s.groupRepo.GetGroupVoteCount(ctx, groupID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"member_count": memberCount,
		"vote_count":   voteCount,
	}, nil
}

// ValidateGroupAccess checks if a user has access to a group (is a member)
func (s *GroupService) ValidateGroupAccess(ctx context.Context, groupID, userID string) (bool, error) {
	isMember, err := s.groupRepo.IsUserInGroup(ctx, groupID, userID)
	if err != nil {
		return false, err
	}

	return isMember, nil
}

// ValidateGroupAdminAccess checks if a user has admin access to a group
func (s *GroupService) ValidateGroupAdminAccess(ctx context.Context, groupID, userID string) (bool, error) {
	member, err := s.groupRepo.GetGroupMember(ctx, groupID, userID)
	if err != nil {
		return false, err
	}

	return member.Role == "admin" || member.Role == "moderator", nil
}
