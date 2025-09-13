package services

import (
	"context"
	"spark-park-cricket-backend/internal/models"
)

// ScorecardServiceInterface defines the interface for scorecard business logic operations
type ScorecardServiceInterface interface {
	StartScoring(ctx context.Context, matchID string) error
	AddBall(ctx context.Context, req *models.BallEventRequest) error
	GetScorecard(ctx context.Context, matchID string) (*models.ScorecardResponse, error)
	GetCurrentOver(ctx context.Context, matchID string, inningsNumber int) (*models.ScorecardOver, error)
	ShouldCompleteMatch(ctx context.Context, matchID string, secondInnings *models.Innings, match *models.Match) (bool, string)
	ValidateInningsOrder(ctx context.Context, matchID string, match *models.Match, inningsNumber int) error
	GetNonTossWinner(tossWinner models.TeamType) models.TeamType
}
