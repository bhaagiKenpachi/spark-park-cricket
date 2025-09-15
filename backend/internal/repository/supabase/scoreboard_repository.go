package supabase

import (
	"context"
	"fmt"
	"log"
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
	log.Printf("DEBUG: scoreboardRepository.GetByMatchID called with matchID: %s", matchID)

	var result []models.LiveScoreboard
	log.Printf("DEBUG: Calling Supabase From('live_scoreboard').Select().Eq('match_id', '%s')", matchID)

	_, err := r.client.From("live_scoreboard").Select("*", "", false).Eq("match_id", matchID).ExecuteTo(&result)
	if err != nil {
		log.Printf("DEBUG: Supabase query failed: %v", err)
		return nil, err
	}

	log.Printf("DEBUG: Supabase query successful, found %d scoreboards", len(result))

	if len(result) == 0 {
		log.Printf("DEBUG: No scoreboard found for matchID: %s", matchID)
		return nil, fmt.Errorf("scoreboard not found")
	}

	log.Printf("DEBUG: Returning scoreboard: %+v", result[0])
	return &result[0], nil
}

func (r *scoreboardRepository) Create(ctx context.Context, scoreboard *models.LiveScoreboard) error {
	log.Printf("DEBUG: scoreboardRepository.Create called with scoreboard: %+v", scoreboard)

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

	log.Printf("DEBUG: Created scoreboardData map: %+v", scoreboardData)

	scoreboardDataSlice := []map[string]interface{}{scoreboardData}
	log.Printf("DEBUG: Created scoreboardDataSlice: %+v", scoreboardDataSlice)

	var result []models.LiveScoreboard
	log.Printf("DEBUG: Calling Supabase Insert with scoreboardDataSlice")

	_, err := r.client.From("live_scoreboard").Insert(scoreboardDataSlice, false, "", "", "").ExecuteTo(&result)
	if err != nil {
		log.Printf("DEBUG: Supabase Insert failed: %v", err)
		return err
	}

	log.Printf("DEBUG: Supabase Insert successful, result: %+v", result)

	if len(result) > 0 {
		*scoreboard = result[0]
		log.Printf("DEBUG: Copied result back to scoreboard: %+v", scoreboard)
	}

	return nil
}

func (r *scoreboardRepository) Update(ctx context.Context, matchID string, scoreboard *models.LiveScoreboard) error {
	log.Printf("DEBUG: scoreboardRepository.Update called with matchID: %s, scoreboard: %+v", matchID, scoreboard)

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

	log.Printf("DEBUG: Created scoreboardData map for update: %+v", scoreboardData)

	var result []models.LiveScoreboard
	log.Printf("DEBUG: Calling Supabase Update with scoreboardData")

	_, err := r.client.From("live_scoreboard").Update(scoreboardData, "", "").Eq("match_id", matchID).ExecuteTo(&result)
	if err != nil {
		log.Printf("DEBUG: Supabase Update failed: %v", err)
		return err
	}

	log.Printf("DEBUG: Supabase Update successful, result: %+v", result)

	if len(result) > 0 {
		*scoreboard = result[0]
		log.Printf("DEBUG: Copied result back to scoreboard: %+v", scoreboard)
	}

	return nil
}

func (r *scoreboardRepository) Delete(ctx context.Context, matchID string) error {
	_, err := r.client.From("live_scoreboard").Delete("", "").Eq("match_id", matchID).ExecuteTo(nil)
	return err
}
