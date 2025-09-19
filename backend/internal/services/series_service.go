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
	fmt.Printf("DEBUG: SeriesService.CreateSeries - Starting creation with request: %+v\n", req)

	// Validate business rules
	if req.EndDate.Before(req.StartDate) {
		fmt.Printf("DEBUG: SeriesService.CreateSeries - Validation failed: end date before start date\n")
		return nil, fmt.Errorf("end date must be after start date")
	}

	// Get user ID from context
	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		fmt.Printf("DEBUG: SeriesService.CreateSeries - Authentication failed: no user_id in context\n")
		return nil, fmt.Errorf("user authentication required")
	}
	fmt.Printf("DEBUG: SeriesService.CreateSeries - User ID from context: %s\n", userID)

	// Create series model
	series := &models.Series{
		Name:      req.Name,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		CreatedBy: userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	fmt.Printf("DEBUG: SeriesService.CreateSeries - Created series model: %+v\n", series)

	// Save to repository
	fmt.Printf("DEBUG: SeriesService.CreateSeries - Calling repository.Create\n")
	err := s.seriesRepo.Create(ctx, series)
	if err != nil {
		fmt.Printf("DEBUG: SeriesService.CreateSeries - Repository error: %v\n", err)
		return nil, fmt.Errorf("failed to create series: %w", err)
	}

	fmt.Printf("DEBUG: SeriesService.CreateSeries - Successfully created series with ID: %s\n", series.ID)
	return series, nil
}

// GetSeries retrieves a series by ID
func (s *SeriesService) GetSeries(ctx context.Context, id string) (*models.Series, error) {
	fmt.Printf("DEBUG: SeriesService.GetSeries - Starting retrieval with ID: %s\n", id)

	if id == "" {
		fmt.Printf("DEBUG: SeriesService.GetSeries - Validation failed: empty ID\n")
		return nil, fmt.Errorf("series ID is required")
	}

	fmt.Printf("DEBUG: SeriesService.GetSeries - Calling repository.GetByID\n")
	series, err := s.seriesRepo.GetByID(ctx, id)
	if err != nil {
		fmt.Printf("DEBUG: SeriesService.GetSeries - Repository error: %v\n", err)
		return nil, fmt.Errorf("failed to get series: %w", err)
	}

	fmt.Printf("DEBUG: SeriesService.GetSeries - Successfully retrieved series: %+v\n", series)
	return series, nil
}

// ListSeries retrieves all series with optional filtering
func (s *SeriesService) ListSeries(ctx context.Context, filters *models.SeriesFilters) ([]*models.Series, error) {
	fmt.Printf("DEBUG: SeriesService.ListSeries - Starting list with filters: %+v\n", filters)

	// Set default values
	if filters.Limit <= 0 {
		filters.Limit = 20
	}
	if filters.Limit > 100 {
		filters.Limit = 100
	}
	fmt.Printf("DEBUG: SeriesService.ListSeries - Using filters: %+v\n", filters)

	fmt.Printf("DEBUG: SeriesService.ListSeries - Calling repository.GetAll\n")
	fmt.Printf("DEBUG: SeriesService.ListSeries - Repository type: %T\n", s.seriesRepo)
	series, err := s.seriesRepo.GetAll(ctx, filters)
	if err != nil {
		fmt.Printf("DEBUG: SeriesService.ListSeries - Repository error: %v\n", err)
		return nil, fmt.Errorf("failed to list series: %w", err)
	}

	fmt.Printf("DEBUG: SeriesService.ListSeries - Successfully retrieved %d series\n", len(series))
	for i, s := range series {
		fmt.Printf("DEBUG: SeriesService.ListSeries - Series %d: ID=%s, Name=%s, CreatedBy=%s, CreatedAt=%s\n",
			i+1, s.ID, s.Name, s.CreatedBy, s.CreatedAt.Format("2006-01-02 15:04:05"))
	}

	return series, nil
}

// UpdateSeries updates an existing series
func (s *SeriesService) UpdateSeries(ctx context.Context, id string, req *models.UpdateSeriesRequest) (*models.Series, error) {
	fmt.Printf("DEBUG: SeriesService.UpdateSeries - Starting update with ID: %s, request: %+v\n", id, req)

	if id == "" {
		fmt.Printf("DEBUG: SeriesService.UpdateSeries - Validation failed: empty ID\n")
		return nil, fmt.Errorf("series ID is required")
	}

	// Get user ID from context
	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		fmt.Printf("DEBUG: SeriesService.UpdateSeries - Authentication failed: no user_id in context\n")
		return nil, fmt.Errorf("user authentication required")
	}
	fmt.Printf("DEBUG: SeriesService.UpdateSeries - User ID from context: %s\n", userID)

	// Get existing series
	fmt.Printf("DEBUG: SeriesService.UpdateSeries - Getting existing series\n")
	series, err := s.seriesRepo.GetByID(ctx, id)
	if err != nil {
		fmt.Printf("DEBUG: SeriesService.UpdateSeries - Failed to get existing series: %v\n", err)
		return nil, fmt.Errorf("failed to get series: %w", err)
	}
	fmt.Printf("DEBUG: SeriesService.UpdateSeries - Found existing series: %+v\n", series)

	// Check ownership
	if series.CreatedBy != userID {
		fmt.Printf("DEBUG: SeriesService.UpdateSeries - Access denied: user %s cannot update series created by %s\n", userID, series.CreatedBy)
		return nil, fmt.Errorf("access denied: you can only update series you created")
	}

	// Update fields if provided
	if req.Name != nil {
		series.Name = *req.Name
		fmt.Printf("DEBUG: SeriesService.UpdateSeries - Updated name to: %s\n", *req.Name)
	}
	if req.StartDate != nil {
		series.StartDate = *req.StartDate
		fmt.Printf("DEBUG: SeriesService.UpdateSeries - Updated start date to: %s\n", req.StartDate.Format(time.RFC3339))
	}
	if req.EndDate != nil {
		series.EndDate = *req.EndDate
		fmt.Printf("DEBUG: SeriesService.UpdateSeries - Updated end date to: %s\n", req.EndDate.Format(time.RFC3339))
	}

	// Validate business rules
	if series.EndDate.Before(series.StartDate) {
		fmt.Printf("DEBUG: SeriesService.UpdateSeries - Validation failed: end date before start date\n")
		return nil, fmt.Errorf("end date must be after start date")
	}

	series.UpdatedAt = time.Now()
	fmt.Printf("DEBUG: SeriesService.UpdateSeries - Updated series model: %+v\n", series)

	// Save changes
	fmt.Printf("DEBUG: SeriesService.UpdateSeries - Calling repository.Update\n")
	err = s.seriesRepo.Update(ctx, id, series)
	if err != nil {
		fmt.Printf("DEBUG: SeriesService.UpdateSeries - Repository error: %v\n", err)
		return nil, fmt.Errorf("failed to update series: %w", err)
	}

	fmt.Printf("DEBUG: SeriesService.UpdateSeries - Successfully updated series: %+v\n", series)
	return series, nil
}

// DeleteSeries deletes a series
func (s *SeriesService) DeleteSeries(ctx context.Context, id string) error {
	fmt.Printf("DEBUG: SeriesService.DeleteSeries - Starting deletion with ID: %s\n", id)

	if id == "" {
		fmt.Printf("DEBUG: SeriesService.DeleteSeries - Validation failed: empty ID\n")
		return fmt.Errorf("series ID is required")
	}

	// Get user ID from context
	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		fmt.Printf("DEBUG: SeriesService.DeleteSeries - Authentication failed: no user_id in context\n")
		return fmt.Errorf("user authentication required")
	}
	fmt.Printf("DEBUG: SeriesService.DeleteSeries - User ID from context: %s\n", userID)

	// Check if series exists and get it
	fmt.Printf("DEBUG: SeriesService.DeleteSeries - Getting series to check ownership\n")
	series, err := s.seriesRepo.GetByID(ctx, id)
	if err != nil {
		fmt.Printf("DEBUG: SeriesService.DeleteSeries - Failed to get series: %v\n", err)
		return fmt.Errorf("series not found: %w", err)
	}
	fmt.Printf("DEBUG: SeriesService.DeleteSeries - Found series: %+v\n", series)

	// Check ownership
	if series.CreatedBy != userID {
		fmt.Printf("DEBUG: SeriesService.DeleteSeries - Access denied: user %s cannot delete series created by %s\n", userID, series.CreatedBy)
		return fmt.Errorf("access denied: you can only delete series you created")
	}

	// Delete series
	fmt.Printf("DEBUG: SeriesService.DeleteSeries - Calling repository.Delete\n")
	err = s.seriesRepo.Delete(ctx, id)
	if err != nil {
		fmt.Printf("DEBUG: SeriesService.DeleteSeries - Repository error: %v\n", err)
		return fmt.Errorf("failed to delete series: %w", err)
	}

	fmt.Printf("DEBUG: SeriesService.DeleteSeries - Successfully deleted series with ID: %s\n", id)
	return nil
}
