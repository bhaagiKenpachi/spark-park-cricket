package interfaces

import (
	"context"
	"spark-park-cricket-backend/internal/models"
)

// SeriesRepository defines the interface for series data operations
type SeriesRepository interface {
	Create(ctx context.Context, series *models.Series) error
	GetByID(ctx context.Context, id string) (*models.Series, error)
	GetAll(ctx context.Context, filters *models.SeriesFilters) ([]*models.Series, error)
	Update(ctx context.Context, id string, series *models.Series) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context) (int64, error)
}
