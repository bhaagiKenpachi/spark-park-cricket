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

// NewScoreboardRepository creates a new scoreboard repository
func NewScoreboardRepository(client *supabase.Client) interfaces.ScoreboardRepository {
	return &scoreboardRepository{
		client: client,
	}
}

// ScoreboardRepository implementations
func (r *scoreboardRepository) GetByMatchID(ctx context.Context, matchID string) (*models.LiveScoreboard, error) {

	var result []models.LiveScoreboard

	_, err := r.client.From("live_scoreboard").Select("*", "", false).Eq("match_id", matchID).ExecuteTo(&result)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("scoreboard not found")
	}

	return &result[0], nil
}

func (r *scoreboardRepository) Create(ctx context.Context, scoreboard *models.LiveScoreboard) error {

	// Create a map to avoid UUID issues
	scoreboardData := map[string]interface{}{
		"match_id":     scoreboard.MatchID,
		"batting_team": scoreboard.BattingTeam,
		"score":        scoreboard.Score,
		"wickets":      scoreboard.Wickets,
		"overs":        scoreboard.Overs,
		"balls":        scoreboard.Balls,
		"created_at":   scoreboard.CreatedAt,
		"updated_at":   scoreboard.UpdatedAt,
	}

	scoreboardDataSlice := []map[string]interface{}{scoreboardData}

	var result []models.LiveScoreboard

	_, err := r.client.From("live_scoreboard").Insert(scoreboardDataSlice, false, "", "", "").ExecuteTo(&result)
	if err != nil {
		return err
	}

	if len(result) > 0 {
		*scoreboard = result[0]
	}

	return nil
}

func (r *scoreboardRepository) Update(ctx context.Context, matchID string, scoreboard *models.LiveScoreboard) error {

	// Create a map to avoid UUID issues
	scoreboardData := map[string]interface{}{
		"match_id":     scoreboard.MatchID,
		"batting_team": scoreboard.BattingTeam,
		"score":        scoreboard.Score,
		"wickets":      scoreboard.Wickets,
		"overs":        scoreboard.Overs,
		"balls":        scoreboard.Balls,
		"updated_at":   scoreboard.UpdatedAt,
	}

	var result []models.LiveScoreboard

	_, err := r.client.From("live_scoreboard").Update(scoreboardData, "", "").Eq("match_id", matchID).ExecuteTo(&result)
	if err != nil {
		return err
	}

	if len(result) > 0 {
		*scoreboard = result[0]
	}

	return nil
}

func (r *scoreboardRepository) Delete(ctx context.Context, matchID string) error {
	_, err := r.client.From("live_scoreboard").Delete("", "").Eq("match_id", matchID).ExecuteTo(nil)
	return err
}
