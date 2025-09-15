package services

import (
	"context"
	"spark-park-cricket-backend/internal/models"
)

// MatchServiceInterface defines the interface for match business logic operations
type MatchServiceInterface interface {
	CreateMatch(ctx context.Context, req *models.CreateMatchRequest) (*models.Match, error)
	GetMatch(ctx context.Context, id string) (*models.Match, error)
	ListMatches(ctx context.Context, filters *models.MatchFilters) ([]*models.Match, error)
	UpdateMatch(ctx context.Context, id string, req *models.UpdateMatchRequest) (*models.Match, error)
	DeleteMatch(ctx context.Context, id string) error
	GetMatchesBySeries(ctx context.Context, seriesID string) ([]*models.Match, error)
}
