package interfaces

import (
	"context"
	"spark-park-cricket-backend/internal/models"
)

// ScorecardServiceInterface defines the interface for scorecard operations
type ScorecardServiceInterface interface {
	StartScoring(ctx context.Context, matchID string) error
	AddBall(ctx context.Context, req *models.BallEventRequest) error
	UndoBall(ctx context.Context, matchID string, inningsNumber int) error
	GetScorecard(ctx context.Context, matchID string) (*models.ScorecardResponse, error)
	GetCurrentOver(ctx context.Context, matchID string, inningsNumber int) (*models.ScorecardOver, error)
	GetBallsByOver(ctx context.Context, overID string) ([]*models.ScorecardBall, error)
	ShouldCompleteMatch(ctx context.Context, matchID string, secondInnings *models.Innings, match *models.Match) (bool, string)
	ValidateInningsOrder(ctx context.Context, matchID string, match *models.Match, inningsNumber int) error
	GetNonTossWinner(tossWinner models.TeamType) models.TeamType
}
