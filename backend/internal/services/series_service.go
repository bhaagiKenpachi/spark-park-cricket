package services

import (
	"context"
	"fmt"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/repository/interfaces"
	"time"
)

// SeriesService handles business logic for series operations
type SeriesService struct {
	seriesRepo interfaces.SeriesRepository
}

// NewSeriesService creates a new series service
func NewSeriesService(seriesRepo interfaces.SeriesRepository) *SeriesService {
	return &SeriesService{
		seriesRepo: seriesRepo,
	}
}

// CreateSeries creates a new series
func (s *SeriesService) CreateSeries(ctx context.Context, req *models.CreateSeriesRequest) (*models.Series, error) {
	// Validate business rules
	if req.EndDate.Before(req.StartDate) {
		return nil, fmt.Errorf("end date must be after start date")
	}

	// Create series model
	series := &models.Series{
		Name:      req.Name,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}

	// Save to repository
	err := s.seriesRepo.Create(ctx, series)
	if err != nil {
		return nil, fmt.Errorf("failed to create series: %w", err)
	}

	return series, nil
}

// GetSeries retrieves a series by ID
func (s *SeriesService) GetSeries(ctx context.Context, id string) (*models.Series, error) {
	if id == "" {
		return nil, fmt.Errorf("series ID is required")
	}

	series, err := s.seriesRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get series: %w", err)
	}

	return series, nil
}

// ListSeries retrieves all series with optional filtering
func (s *SeriesService) ListSeries(ctx context.Context, filters *models.SeriesFilters) ([]*models.Series, error) {
	// Set default values
	if filters.Limit <= 0 {
		filters.Limit = 20
	}
	if filters.Limit > 100 {
		filters.Limit = 100
	}

	series, err := s.seriesRepo.GetAll(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list series: %w", err)
	}

	return series, nil
}

// UpdateSeries updates an existing series
func (s *SeriesService) UpdateSeries(ctx context.Context, id string, req *models.UpdateSeriesRequest) (*models.Series, error) {
	if id == "" {
		return nil, fmt.Errorf("series ID is required")
	}

	// Get existing series
	series, err := s.seriesRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get series: %w", err)
	}

	// Update fields if provided
	if req.Name != nil {
		series.Name = *req.Name
	}
	if req.StartDate != nil {
		series.StartDate = *req.StartDate
	}
	if req.EndDate != nil {
		series.EndDate = *req.EndDate
	}

	// Validate business rules
	if series.EndDate.Before(series.StartDate) {
		return nil, fmt.Errorf("end date must be after start date")
	}

	series.UpdatedAt = time.Now()

	// Save changes
	err = s.seriesRepo.Update(ctx, id, series)
	if err != nil {
		return nil, fmt.Errorf("failed to update series: %w", err)
	}

	return series, nil
}

// DeleteSeries deletes a series
func (s *SeriesService) DeleteSeries(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("series ID is required")
	}

	// Check if series exists
	_, err := s.seriesRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("series not found: %w", err)
	}

	// Delete series
	err = s.seriesRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete series: %w", err)
	}

	return nil
}
