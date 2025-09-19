package supabase

import (
	"context"
	"fmt"
	"log"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/repository/interfaces"

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

func (r *seriesRepository) GetAll(ctx context.Context, filters *models.SeriesFilters) ([]*models.Series, error) {
	log.Printf("=== SERIES REPOSITORY: GetAll ===")
	log.Printf("Filters: %+v", filters)

	var result []models.Series
	query := r.client.From("series").Select("*", "", false)

	if filters != nil {
		if filters.Limit > 0 {
			query = query.Limit(filters.Limit, "")
			log.Printf("Applied limit: %d", filters.Limit)
		}
		// Note: Offset is not supported by this Supabase client version
		// Use Range method for pagination if needed
	}

	log.Printf("Executing query to 'series' table...")
	_, err := query.ExecuteTo(&result)
	if err != nil {
		log.Printf("ERROR: Query failed: %v", err)
		return nil, err
	}

	log.Printf("Query successful, found %d series", len(result))
	for i, s := range result {
		log.Printf("Series %d: ID=%s, Name=%s, CreatedBy=%s", i+1, s.ID, s.Name, s.CreatedBy)
	}

	// Convert to slice of pointers
	series := make([]*models.Series, len(result))
	for i := range result {
		series[i] = &result[i]
	}

	log.Printf("Returning %d series", len(series))
	return series, nil
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
