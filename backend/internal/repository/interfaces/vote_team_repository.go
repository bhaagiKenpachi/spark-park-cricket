package interfaces

import (
	"context"
	"spark-park-cricket-backend/internal/models"
)

// VoteTeamRepositoryInterface defines methods for vote team data operations
type VoteTeamRepositoryInterface interface {
	// Team CRUD operations
	CreateTeam(ctx context.Context, team *models.VoteTeam) error
	GetTeamByID(ctx context.Context, id string) (*models.VoteTeam, error)
	GetTeamsByVoteID(ctx context.Context, voteID string) ([]*models.VoteTeam, error)
	GetTeamByVoteAndLetter(ctx context.Context, voteID, teamLetter string) (*models.VoteTeam, error)
	UpdateTeam(ctx context.Context, id string, teamName *string, captainID *string) error
	DeleteTeam(ctx context.Context, id string) error

	// Team with players
	GetTeamWithPlayers(ctx context.Context, teamID string) (*models.VoteTeamWithPlayers, error)
	GetTeamsWithPlayersByVoteID(ctx context.Context, voteID string) ([]*models.VoteTeamWithPlayers, error)

	// Player operations
	AddPlayerToTeam(ctx context.Context, teamID, userID string) error
	RemovePlayerFromTeam(ctx context.Context, teamID, userID string) error
	GetPlayersByTeamID(ctx context.Context, teamID string) ([]*models.User, error)
	IsUserInTeam(ctx context.Context, teamID, userID string) (bool, error)
	IsUserInAnyTeamForVote(ctx context.Context, voteID, userID string) (bool, error)
	GetTeamPlayerCount(ctx context.Context, teamID string) (int, error)

	// Validation helpers
	HasUserVoted(ctx context.Context, voteID, userID string) (bool, error)
}
