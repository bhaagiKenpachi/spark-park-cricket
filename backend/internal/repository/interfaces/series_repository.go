package interfaces

import (
	"context"
	"spark-park-cricket-backend/internal/models"
)

// PaginatedSeriesResult represents a paginated result with metadata
type PaginatedSeriesResult struct {
	Series     []*models.Series `json:"series"`
	TotalItems int              `json:"total_items"`
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
	TotalPages int              `json:"total_pages"`
}

// SeriesRepository defines the interface for series data operations
type SeriesRepository interface {
	Create(ctx context.Context, series *models.Series) error
	GetByID(ctx context.Context, id string) (*models.Series, error)
	GetAll(ctx context.Context, filters *models.SeriesFilters) (*PaginatedSeriesResult, error)
	Update(ctx context.Context, id string, series *models.Series) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context) (int64, error)
}
