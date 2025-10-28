package services

import (
	"context"
	"spark-park-cricket-backend/internal/models"
)

// GroupServiceInterface defines the interface for group business logic operations
type GroupServiceInterface interface {
	// Group CRUD operations
	CreateGroup(ctx context.Context, req *models.CreateGroupRequest, creatorID string) (*models.Group, error)
	GetGroup(ctx context.Context, groupID string) (*models.GroupWithCreator, error)
	GetGroupWithMembers(ctx context.Context, groupID string) (*models.GroupWithMembersAndUsers, error)
	ListGroups(ctx context.Context, limit, offset int) ([]*models.GroupWithCreator, error)
	ListGroupsByType(ctx context.Context, groupType models.GroupType, limit, offset int) ([]*models.GroupWithCreator, error)
	ListGroupsByCreator(ctx context.Context, creatorID string, limit, offset int) ([]*models.GroupWithCreator, error)
	ListGroupsByUser(ctx context.Context, userID string, limit, offset int) ([]*models.GroupWithCreator, error)
	UpdateGroup(ctx context.Context, groupID string, req *models.UpdateGroupRequest, userID string) (*models.Group, error)
	DeleteGroup(ctx context.Context, groupID string, userID string) error

	// Group member operations
	AddGroupMember(ctx context.Context, groupID string, req *models.AddGroupMemberRequest, adminUserID string) (*models.GroupMember, error)
	RemoveGroupMember(ctx context.Context, groupID string, req *models.RemoveGroupMemberRequest, adminUserID string) error
	GetGroupMembers(ctx context.Context, groupID string) ([]*models.GroupMemberWithUser, error)
	UpdateGroupMemberRole(ctx context.Context, groupID, userID, role string, adminUserID string) (*models.GroupMember, error)
	JoinGroup(ctx context.Context, groupID string, userID string) (*models.GroupMember, error)
	LeaveGroup(ctx context.Context, groupID string, userID string) error

	// Vote-Group association operations
	AssignGroupsToVote(ctx context.Context, voteID string, req *models.AssignGroupsToVoteRequest, userID string) error
	RemoveGroupsFromVote(ctx context.Context, voteID string, groupIDs []string, userID string) error
	GetVoteGroups(ctx context.Context, voteID string) ([]*models.Group, error)
	GetGroupVotes(ctx context.Context, groupID string) ([]*models.Vote, error)

	// Group voting results
	GetGroupVoteResults(ctx context.Context, voteID, groupID string) (*models.GroupVoteResult, error)
	GetVoteResultsByGroups(ctx context.Context, voteID string) ([]*models.GroupVoteResult, error)
	GetVoteWithGroupResults(ctx context.Context, voteID string) (*models.VoteWithGroupResults, error)

	// Utility operations
	SearchGroups(ctx context.Context, query string, limit, offset int) ([]*models.GroupWithCreator, error)
	GetGroupStats(ctx context.Context, groupID string) (map[string]interface{}, error)
	ValidateGroupAccess(ctx context.Context, groupID, userID string) (bool, error)
	ValidateGroupAdminAccess(ctx context.Context, groupID, userID string) (bool, error)
}
