package services

import (
	"context"
	"spark-park-cricket-backend/internal/models"
)

// MatchStateServiceInterface defines the interface for match state operations
type MatchStateServiceInterface interface {
	StartMatch(ctx context.Context, matchID string) (*models.Match, error)
	CompleteMatch(ctx context.Context, matchID string) (*models.Match, error)
	CancelMatch(ctx context.Context, matchID string) (*models.Match, error)
	GetMatchSummary(ctx context.Context, matchID string) (*MatchSummary, error)
}
