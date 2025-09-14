package supabase

import (
	"context"
	"fmt"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/repository/interfaces"

	"github.com/supabase-community/supabase-go"
)

type overRepository struct {
	client *supabase.Client
	schema string
}

// NewOverRepository creates a new over repository
func NewOverRepository(client *supabase.Client, schema string) interfaces.OverRepository {
	return &overRepository{
		client: client,
		schema: schema,
	}
}

func (r *overRepository) Create(ctx context.Context, over *models.Over) error {
	// Create a map to avoid UUID issues
	overData := map[string]interface{}{
		"match_id":     over.MatchID,
		"over_number":  over.OverNumber,
		"batting_team": over.BattingTeam,
		"total_runs":   over.TotalRuns,
		"total_balls":  over.TotalBalls,
		"created_at":   over.CreatedAt,
	}

	overDataSlice := []map[string]interface{}{overData}
	var result []models.Over

	_, err := r.client.From("overs").Insert(overDataSlice, false, "", "", "").ExecuteTo(&result)
	if err != nil {
		return err
	}

	if len(result) > 0 {
		*over = result[0]
	}

	return nil
}

func (r *overRepository) GetByID(ctx context.Context, id string) (*models.Over, error) {
	var result []models.Over
	_, err := r.client.From("overs").Select("*", "", false).Eq("id", id).ExecuteTo(&result)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("over not found")
	}
	return &result[0], nil
}

func (r *overRepository) GetByMatchID(ctx context.Context, matchID string) ([]*models.Over, error) {
	var result []models.Over
	_, err := r.client.From("overs").Select("*", "", false).Eq("match_id", matchID).ExecuteTo(&result)
	if err != nil {
		return nil, err
	}

	// Convert to slice of pointers
	overs := make([]*models.Over, len(result))
	for i := range result {
		overs[i] = &result[i]
	}
	return overs, nil
}

func (r *overRepository) Update(ctx context.Context, id string, over *models.Over) error {
	// Create a map to avoid UUID issues
	overData := map[string]interface{}{
		"match_id":     over.MatchID,
		"over_number":  over.OverNumber,
		"batting_team": over.BattingTeam,
		"total_runs":   over.TotalRuns,
		"total_balls":  over.TotalBalls,
	}

	var result []models.Over
	_, err := r.client.From("overs").Update(overData, "", "").Eq("id", id).ExecuteTo(&result)
	if err != nil {
		return err
	}

	if len(result) > 0 {
		*over = result[0]
	}

	return nil
}

func (r *overRepository) Delete(ctx context.Context, id string) error {
	_, err := r.client.From("overs").Delete("", "").Eq("id", id).ExecuteTo(nil)
	return err
}

func (r *overRepository) Count(ctx context.Context) (int64, error) {
	var result []models.Over
	_, err := r.client.From("overs").Select("*", "", false).ExecuteTo(&result)
	if err != nil {
		return 0, err
	}
	return int64(len(result)), nil
}
