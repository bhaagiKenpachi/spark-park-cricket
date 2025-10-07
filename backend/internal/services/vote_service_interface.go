package services

import (
	"context"
	"spark-park-cricket-backend/internal/models"
)

// VoteServiceInterface defines the interface for vote service operations
type VoteServiceInterface interface {
	// Vote operations
	CreateVote(ctx context.Context, req *models.CreateVoteRequest, userID string) (*models.Vote, error)
	GetVote(ctx context.Context, id string) (*models.VoteWithOptions, error)
	GetVoteWithResults(ctx context.Context, id string, userID string) (*models.VoteWithResults, error)
	UpdateVote(ctx context.Context, id string, req *models.UpdateVoteRequest, userID string) (*models.Vote, error)
	DeleteVote(ctx context.Context, id string, userID string) error
	ListVotes(ctx context.Context, filters *models.VoteFilters) ([]*models.Vote, error)

	// Voting operations
	CastVote(ctx context.Context, voteID string, req *models.VoteRequest, userID string) error
	GetUserVote(ctx context.Context, voteID string, userID string) (*models.UserVote, error)
	HasUserVoted(ctx context.Context, voteID string, userID string) (bool, error)

	// Vote management
	CloseVote(ctx context.Context, id string, userID string) error
	CancelVote(ctx context.Context, id string, userID string) error
}
