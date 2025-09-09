package supabase

import (
	"context"
	"fmt"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/repository/interfaces"

	"github.com/supabase-community/supabase-go"
)

type scoreboardRepository struct {
	client *supabase.Client
}

type overRepository struct {
	client *supabase.Client
}

type ballRepository struct {
	client *supabase.Client
}

// NewScoreboardRepository creates a new scoreboard repository
func NewScoreboardRepository(client *supabase.Client) interfaces.ScoreboardRepository {
	return &scoreboardRepository{
		client: client,
	}
}

// NewOverRepository creates a new over repository
func NewOverRepository(client *supabase.Client) interfaces.OverRepository {
	return &overRepository{
		client: client,
	}
}

// NewBallRepository creates a new ball repository
func NewBallRepository(client *supabase.Client) interfaces.BallRepository {
	return &ballRepository{
		client: client,
	}
}

// ScoreboardRepository implementations
func (r *scoreboardRepository) GetByMatchID(ctx context.Context, matchID string) (*models.LiveScoreboard, error) {
	var result []models.LiveScoreboard
	_, err := r.client.From("live_scoreboards").Select("*", "", false).Eq("match_id", matchID).ExecuteTo(&result)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("scoreboard not found")
	}
	return &result[0], nil
}

func (r *scoreboardRepository) Create(ctx context.Context, scoreboard *models.LiveScoreboard) error {
	_, err := r.client.From("live_scoreboards").Insert(scoreboard, false, "", "", "").ExecuteTo(&scoreboard)
	return err
}

func (r *scoreboardRepository) Update(ctx context.Context, matchID string, scoreboard *models.LiveScoreboard) error {
	_, err := r.client.From("live_scoreboards").Update(scoreboard, "", "").Eq("match_id", matchID).ExecuteTo(&scoreboard)
	return err
}

func (r *scoreboardRepository) Delete(ctx context.Context, matchID string) error {
	_, err := r.client.From("live_scoreboards").Delete("", "").Eq("match_id", matchID).ExecuteTo(nil)
	return err
}

// OverRepository implementations
func (r *overRepository) Create(ctx context.Context, over *models.Over) error {
	_, err := r.client.From("overs").Insert(over, false, "", "", "").ExecuteTo(&over)
	return err
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

func (r *overRepository) Update(ctx context.Context, id string, over *models.Over) error {
	_, err := r.client.From("overs").Update(over, "", "").Eq("id", id).ExecuteTo(&over)
	return err
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

// BallRepository implementations
func (r *ballRepository) Create(ctx context.Context, ball *models.Ball) error {
	_, err := r.client.From("balls").Insert(ball, false, "", "", "").ExecuteTo(&ball)
	return err
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

func (r *ballRepository) GetByMatchID(ctx context.Context, matchID string) ([]*models.Ball, error) {
	var result []models.Ball
	_, err := r.client.From("balls").Select("*", "", false).Eq("match_id", matchID).ExecuteTo(&result)
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

func (r *ballRepository) Update(ctx context.Context, id string, ball *models.Ball) error {
	_, err := r.client.From("balls").Update(ball, "", "").Eq("id", id).ExecuteTo(&ball)
	return err
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
