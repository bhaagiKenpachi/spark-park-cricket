package interfaces

import (
	"context"
	"spark-park-cricket-backend/internal/models"
)

// VoteRepositoryInterface defines the interface for vote repository operations
type VoteRepositoryInterface interface {
	// Vote operations
	CreateVote(ctx context.Context, vote *models.Vote) error
	GetVoteByID(ctx context.Context, id string) (*models.Vote, error)
	GetVoteWithOptions(ctx context.Context, id string) (*models.VoteWithOptions, error)
	GetVoteWithResults(ctx context.Context, id string, userID string) (*models.VoteWithResults, error)
	UpdateVote(ctx context.Context, vote *models.Vote) error
	DeleteVote(ctx context.Context, id string) error
	ListVotes(ctx context.Context, filters *models.VoteFilters) ([]*models.Vote, error)

	// Vote option operations
	CreateVoteOptions(ctx context.Context, options []*models.VoteOption) error
	GetVoteOptions(ctx context.Context, voteID string) ([]*models.VoteOption, error)
	UpdateVoteOption(ctx context.Context, option *models.VoteOption) error
	DeleteVoteOption(ctx context.Context, id string) error

	// User vote operations
	CreateUserVote(ctx context.Context, userVote *models.UserVote) error
	UpdateUserVote(ctx context.Context, userVote *models.UserVote) error
	GetUserVote(ctx context.Context, voteID, userID string) (*models.UserVote, error)
	HasUserVoted(ctx context.Context, voteID, userID string) (bool, error)
	GetVoteResults(ctx context.Context, voteID string) (map[string]int, error)
	GetVoteResultsWithNames(ctx context.Context, voteID string) (map[string][]models.VoterInfo, error)
	GetVotedUsers(ctx context.Context, voteID string) ([]string, error)
	GetTotalVoteCount(ctx context.Context, voteID string) (int, error)
}
