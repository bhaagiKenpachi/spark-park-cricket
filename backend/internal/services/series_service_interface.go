package services

import (
	"context"
	"spark-park-cricket-backend/internal/models"
)

// SeriesServiceInterface defines the interface for series service operations
type SeriesServiceInterface interface {
	CreateSeries(ctx context.Context, req *models.CreateSeriesRequest) (*models.Series, error)
	GetSeries(ctx context.Context, id string) (*models.Series, error)
	ListSeries(ctx context.Context, filters *models.SeriesFilters) ([]*models.Series, error)
	UpdateSeries(ctx context.Context, id string, req *models.UpdateSeriesRequest) (*models.Series, error)
	DeleteSeries(ctx context.Context, id string) error
}
