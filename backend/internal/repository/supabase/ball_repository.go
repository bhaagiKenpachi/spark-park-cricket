package supabase

import (
	"context"
	"fmt"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/repository/interfaces"

	"github.com/supabase-community/supabase-go"
)

type ballRepository struct {
	client *supabase.Client
}

// NewBallRepository creates a new ball repository
func NewBallRepository(client *supabase.Client) interfaces.BallRepository {
	return &ballRepository{client: client}
}

func (r *ballRepository) Create(ctx context.Context, ball *models.Ball) error {
	// Create a map to avoid UUID issues
	ballData := map[string]interface{}{
		"over_id":     ball.OverID,
		"ball_number": ball.BallNumber,
		"ball_type":   ball.BallType,
		"run_type":    ball.RunType,
		"is_wicket":   ball.IsWicket,
		"created_at":  ball.CreatedAt,
	}

	ballDataSlice := []map[string]interface{}{ballData}
	var result []models.Ball

	_, err := r.client.From("balls").Insert(ballDataSlice, false, "", "", "").ExecuteTo(&result)
	if err != nil {
		return err
	}

	if len(result) > 0 {
		*ball = result[0]
	}

	return nil
}

func (r *ballRepository) GetByID(ctx context.Context, id string) (*models.Ball, error) {
	var result []models.Ball
	_, err := r.client.From("balls").Select("*", "", false).Eq("id", id).ExecuteTo(&result)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("ball not found")
	}
	return &result[0], nil
}

func (r *ballRepository) GetByOverID(ctx context.Context, overID string) ([]*models.Ball, error) {
	var result []models.Ball
	_, err := r.client.From("balls").Select("*", "", false).Eq("over_id", overID).ExecuteTo(&result)
	if err != nil {
		return nil, err
	}

	// Convert to slice of pointers
	balls := make([]*models.Ball, len(result))
	for i := range result {
		balls[i] = &result[i]
	}
	return balls, nil
}

func (r *ballRepository) GetByMatchID(ctx context.Context, matchID string) ([]*models.Ball, error) {
	// This would require a join with overs table, but for simplicity, we'll return empty for now
	// In a real implementation, you'd need to join balls with overs to get balls by match_id
	return []*models.Ball{}, nil
}

func (r *ballRepository) Update(ctx context.Context, id string, ball *models.Ball) error {
	// Create a map to avoid UUID issues
	ballData := map[string]interface{}{
		"over_id":     ball.OverID,
		"ball_number": ball.BallNumber,
		"ball_type":   ball.BallType,
		"run_type":    ball.RunType,
		"is_wicket":   ball.IsWicket,
	}

	var result []models.Ball
	_, err := r.client.From("balls").Update(ballData, "", "").Eq("id", id).ExecuteTo(&result)
	if err != nil {
		return err
	}

	if len(result) > 0 {
		*ball = result[0]
	}

	return nil
}

func (r *ballRepository) Delete(ctx context.Context, id string) error {
	_, err := r.client.From("balls").Delete("", "").Eq("id", id).ExecuteTo(nil)
	return err
}

func (r *ballRepository) Count(ctx context.Context) (int64, error) {
	var result []models.Ball
	_, err := r.client.From("balls").Select("*", "", false).ExecuteTo(&result)
	if err != nil {
		return 0, err
	}
	return int64(len(result)), nil
}
