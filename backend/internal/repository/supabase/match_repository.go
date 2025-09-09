package supabase

import (
	"context"
	"fmt"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/repository/interfaces"

	"github.com/supabase-community/supabase-go"
)

type matchRepository struct {
	client *supabase.Client
}

// NewMatchRepository creates a new match repository
func NewMatchRepository(client *supabase.Client) interfaces.MatchRepository {
	return &matchRepository{
		client: client,
	}
}

func (r *matchRepository) Create(ctx context.Context, match *models.Match) error {
	_, err := r.client.From("matches").Insert(match, false, "", "", "").ExecuteTo(&match)
	return err
}

func (r *matchRepository) GetByID(ctx context.Context, id string) (*models.Match, error) {
	var result []models.Match
	_, err := r.client.From("matches").Select("*", "", false).Eq("id", id).ExecuteTo(&result)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("match not found")
	}
	return &result[0], nil
}

func (r *matchRepository) GetAll(ctx context.Context, filters *models.MatchFilters) ([]*models.Match, error) {
	var result []models.Match
	query := r.client.From("matches").Select("*", "", false)

	if filters != nil {
		if filters.Limit > 0 {
			query = query.Limit(filters.Limit, "")
		}
		// Note: Offset is not supported by this Supabase client version
		// Use Range method for pagination if needed
		if filters.SeriesID != nil && *filters.SeriesID != "" {
			query = query.Eq("series_id", *filters.SeriesID)
		}
		if filters.Status != nil {
			query = query.Eq("status", string(*filters.Status))
		}
	}

	_, err := query.ExecuteTo(&result)
	if err != nil {
		return nil, err
	}

	// Convert to slice of pointers
	matches := make([]*models.Match, len(result))
	for i := range result {
		matches[i] = &result[i]
	}
	return matches, nil
}

func (r *matchRepository) Update(ctx context.Context, id string, match *models.Match) error {
	_, err := r.client.From("matches").Update(match, "", "").Eq("id", id).ExecuteTo(&match)
	return err
}

func (r *matchRepository) Delete(ctx context.Context, id string) error {
	_, err := r.client.From("matches").Delete("", "").Eq("id", id).ExecuteTo(nil)
	return err
}

func (r *matchRepository) GetBySeriesID(ctx context.Context, seriesID string) ([]*models.Match, error) {
	var result []models.Match
	_, err := r.client.From("matches").Select("*", "", false).Eq("series_id", seriesID).ExecuteTo(&result)
	if err != nil {
		return nil, err
	}

	// Convert to slice of pointers
	matches := make([]*models.Match, len(result))
	for i := range result {
		matches[i] = &result[i]
	}
	return matches, nil
}

func (r *matchRepository) Count(ctx context.Context) (int64, error) {
	var result []models.Match
	_, err := r.client.From("matches").Select("*", "", false).ExecuteTo(&result)
	if err != nil {
		return 0, err
	}
	return int64(len(result)), nil
}
