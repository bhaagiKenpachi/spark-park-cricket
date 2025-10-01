package supabase

import (
	"context"
	"fmt"
	"math"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/repository/interfaces"

	"github.com/supabase-community/postgrest-go"
	"github.com/supabase-community/supabase-go"
)

type seriesRepository struct {
	client *supabase.Client
}

// NewSeriesRepository creates a new series repository
func NewSeriesRepository(client *supabase.Client) interfaces.SeriesRepository {
	return &seriesRepository{
		client: client,
	}
}

func (r *seriesRepository) Create(ctx context.Context, series *models.Series) error {
	// Supabase returns an array even for single inserts, so we need to handle that
	var result []models.Series
	_, err := r.client.From("series").Insert(series, false, "", "", "").ExecuteTo(&result)
	if err != nil {
		return err
	}

	if len(result) > 0 {
		// Copy the result back to the original series
		*series = result[0]
	} else {
		// If no result returned, try to get the series by name to get the ID
		// This is a fallback for cases where Supabase doesn't return the created record
		allSeries, err := r.GetAll(ctx, &models.SeriesFilters{Limit: 1000, Offset: 0})
		if err != nil {
			return fmt.Errorf("failed to get created series: %w", err)
		}

		// Find the series by name (assuming unique names)
		for _, s := range allSeries.Series {
			if s.Name == series.Name {
				*series = *s
				break
			}
		}
	}

	return nil
}

func (r *seriesRepository) GetByID(ctx context.Context, id string) (*models.Series, error) {
	var result []models.Series
	_, err := r.client.From("series").Select("*", "", false).Eq("id", id).ExecuteTo(&result)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("series not found")
	}
	return &result[0], nil
}

func (r *seriesRepository) GetAll(ctx context.Context, filters *models.SeriesFilters) (*interfaces.PaginatedSeriesResult, error) {
	fmt.Printf("DEBUG: SupabaseSeriesRepository.GetAll - Starting with filters: %+v\n", filters)

	var result []models.Series
	query := r.client.From("series").Select("*", "", false).Order("created_at", &postgrest.OrderOpts{
		Ascending: false, // DESC order - newest first
	})

	// Always fetch all results and paginate in application layer
	// This ensures consistent ordering and proper pagination
	fmt.Printf("DEBUG: SupabaseSeriesRepository.GetAll - Executing query to database (fetching all for pagination)\n")
	_, err := query.ExecuteTo(&result)
	if err != nil {
		fmt.Printf("DEBUG: SupabaseSeriesRepository.GetAll - Database query error: %v\n", err)
		return nil, err
	}

	fmt.Printf("DEBUG: SupabaseSeriesRepository.GetAll - Database returned %d total series\n", len(result))

	// Apply pagination in application layer
	var paginatedResult []models.Series
	if filters != nil {
		start := filters.Offset
		end := start + filters.Limit

		if start >= len(result) {
			// Offset is beyond available data
			paginatedResult = []models.Series{}
		} else {
			if end > len(result) {
				end = len(result)
			}
			paginatedResult = result[start:end]
		}

		fmt.Printf("DEBUG: SupabaseSeriesRepository.GetAll - Applied pagination: offset=%d, limit=%d, result=%d\n",
			start, filters.Limit, len(paginatedResult))
	} else {
		paginatedResult = result
	}

	for i, s := range paginatedResult {
		fmt.Printf("DEBUG: SupabaseSeriesRepository.GetAll - DB Series %d: ID=%s, Name=%s, CreatedBy=%s, CreatedAt=%s\n",
			i+1, s.ID, s.Name, s.CreatedBy, s.CreatedAt.Format("2006-01-02 15:04:05"))
	}

	// Convert to slice of pointers
	series := make([]*models.Series, len(paginatedResult))
	for i := range paginatedResult {
		series[i] = &paginatedResult[i]
	}

	// Calculate pagination metadata
	totalItems := len(result) // Total items in database
	currentPage := 1
	if filters != nil && filters.Offset > 0 {
		currentPage = (filters.Offset / filters.Limit) + 1
	}
	pageSize := totalItems
	if filters != nil && filters.Limit > 0 {
		pageSize = filters.Limit
	}
	totalPages := 1
	if pageSize > 0 {
		totalPages = int(math.Ceil(float64(totalItems) / float64(pageSize)))
	}

	fmt.Printf("DEBUG: SupabaseSeriesRepository.GetAll - Returning %d series (page %d of %d, %d total items)\n",
		len(series), currentPage, totalPages, totalItems)

	return &interfaces.PaginatedSeriesResult{
		Series:     series,
		TotalItems: totalItems,
		Page:       currentPage,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (r *seriesRepository) Update(ctx context.Context, id string, series *models.Series) error {
	// Supabase returns an array even for single updates, so we need to handle that
	var result []models.Series
	_, err := r.client.From("series").Update(series, "", "").Eq("id", id).ExecuteTo(&result)
	if err != nil {
		return err
	}

	if len(result) > 0 {
		// Copy the result back to the original series
		*series = result[0]
	}

	return nil
}

func (r *seriesRepository) Delete(ctx context.Context, id string) error {
	_, err := r.client.From("series").Delete("", "").Eq("id", id).ExecuteTo(nil)
	return err
}

func (r *seriesRepository) Count(ctx context.Context) (int64, error) {
	var result []models.Series
	_, err := r.client.From("series").Select("*", "", false).ExecuteTo(&result)
	if err != nil {
		return 0, err
	}
	return int64(len(result)), nil
}
