package interfaces

import (
	"context"
	"spark-park-cricket-backend/internal/models"
)

// GroupRepository defines the interface for group data operations
type GroupRepository interface {
	// Group CRUD operations
	CreateGroup(ctx context.Context, group *models.Group) error
	GetGroup(ctx context.Context, groupID string) (*models.Group, error)
	GetGroupWithCreator(ctx context.Context, groupID string) (*models.GroupWithCreator, error)
	GetGroupWithMembers(ctx context.Context, groupID string) (*models.GroupWithMembers, error)
	GetGroupWithMembersAndUsers(ctx context.Context, groupID string) (*models.GroupWithMembersAndUsers, error)
	ListGroups(ctx context.Context, limit, offset int) ([]*models.GroupWithCreator, error)
	ListGroupsByType(ctx context.Context, groupType models.GroupType, limit, offset int) ([]*models.GroupWithCreator, error)
	ListGroupsByCreator(ctx context.Context, creatorID string, limit, offset int) ([]*models.GroupWithCreator, error)
	ListGroupsByUser(ctx context.Context, userID string, limit, offset int) ([]*models.GroupWithCreator, error)
	UpdateGroup(ctx context.Context, group *models.Group) error
	DeleteGroup(ctx context.Context, groupID string) error

	// Group member operations
	AddGroupMember(ctx context.Context, member *models.GroupMember) error
	RemoveGroupMember(ctx context.Context, groupID, userID string) error
	GetGroupMember(ctx context.Context, groupID, userID string) (*models.GroupMember, error)
	GetGroupMembers(ctx context.Context, groupID string) ([]*models.GroupMember, error)
	GetGroupMembersWithUsers(ctx context.Context, groupID string) ([]*models.GroupMemberWithUser, error)
	UpdateGroupMemberRole(ctx context.Context, groupID, userID, role string) error
	IsUserInGroup(ctx context.Context, groupID, userID string) (bool, error)
	GetUserGroups(ctx context.Context, userID string) ([]*models.Group, error)

	// Vote-Group association operations
	AssignGroupsToVote(ctx context.Context, voteID string, groupIDs []string) error
	RemoveGroupsFromVote(ctx context.Context, voteID string, groupIDs []string) error
	GetVoteGroups(ctx context.Context, voteID string) ([]*models.Group, error)
	GetGroupVotes(ctx context.Context, groupID string) ([]*models.Vote, error)
	IsGroupAssignedToVote(ctx context.Context, voteID, groupID string) (bool, error)

	// Group voting results
	GetGroupVoteResults(ctx context.Context, voteID, groupID string) (*models.GroupVoteResult, error)
	GetVoteResultsByGroups(ctx context.Context, voteID string) ([]*models.GroupVoteResult, error)
	GetVoteWithGroupResults(ctx context.Context, voteID string) (*models.VoteWithGroupResults, error)

	// Utility operations
	GetGroupMemberCount(ctx context.Context, groupID string) (int, error)
	GetGroupVoteCount(ctx context.Context, groupID string) (int, error)
	SearchGroups(ctx context.Context, query string, limit, offset int) ([]*models.GroupWithCreator, error)
}
