package interfaces

import (
	"context"
	"spark-park-cricket-backend/internal/models"
)

// MatchRepository defines the interface for match data operations
type MatchRepository interface {
	Create(ctx context.Context, match *models.Match) error
	GetByID(ctx context.Context, id string) (*models.Match, error)
	GetAll(ctx context.Context, filters *models.MatchFilters) ([]*models.Match, error)
	Update(ctx context.Context, id string, match *models.Match) error
	Delete(ctx context.Context, id string) error
	GetBySeriesID(ctx context.Context, seriesID string) ([]*models.Match, error)
	Count(ctx context.Context) (int64, error)
}
