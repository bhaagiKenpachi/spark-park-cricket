package services

import (
	"context"
	"spark-park-cricket-backend/internal/models"
)

// VoteTeamServiceInterface defines methods for vote team business logic
type VoteTeamServiceInterface interface {
	// Team operations
	CreateTeam(ctx context.Context, userID string, req *models.CreateVoteTeamRequest) (*models.VoteTeam, error)
	GetTeamByID(ctx context.Context, teamID string) (*models.VoteTeamWithPlayers, error)
	GetTeamsByVoteID(ctx context.Context, voteID string) ([]*models.VoteTeamWithPlayers, error)
	UpdateTeam(ctx context.Context, userID, teamID string, req *models.UpdateVoteTeamRequest) (*models.VoteTeam, error)
	DeleteTeam(ctx context.Context, userID, teamID string) error

	// Player operations
	AddPlayerToTeam(ctx context.Context, userID, teamID string, req *models.AddPlayerRequest) error
	AddPlayersToTeam(ctx context.Context, userID, teamID string, req *models.TeamAssignmentRequest) error
	RemovePlayerFromTeam(ctx context.Context, userID, teamID, playerID string) error
	GetTeamPlayers(ctx context.Context, teamID string) ([]*models.User, error)
}
