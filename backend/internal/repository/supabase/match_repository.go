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
	return &matchRepository{client: client}
}

func (r *matchRepository) Create(ctx context.Context, match *models.Match) error {
	// Create a map to avoid UUID issues
	matchData := map[string]interface{}{
		"series_id":           match.SeriesID,
		"match_number":        match.MatchNumber,
		"date":                match.Date,
		"status":              match.Status,
		"team_a_player_count": match.TeamAPlayerCount,
		"team_b_player_count": match.TeamBPlayerCount,
		"total_overs":         match.TotalOvers,
		"toss_winner":         match.TossWinner,
		"toss_type":           match.TossType,
		"batting_team":        match.BattingTeam,
		"created_at":          match.CreatedAt,
		"updated_at":          match.UpdatedAt,
	}

	matchDataSlice := []map[string]interface{}{matchData}
	var result []models.Match

	_, err := r.client.From("matches").Insert(matchDataSlice, false, "", "", "").ExecuteTo(&result)
	if err != nil {
		return err
	}

	if len(result) > 0 {
		*match = result[0]
	}

	return nil
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

	if filters.SeriesID != nil {
		query = query.Eq("series_id", *filters.SeriesID)
	}
	if filters.Status != nil {
		query = query.Eq("status", string(*filters.Status))
	}

	query = query.Range(filters.Offset, filters.Offset+filters.Limit-1, "")

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

func (r *matchRepository) Update(ctx context.Context, id string, match *models.Match) error {
	// Create a map to avoid UUID issues
	matchData := map[string]interface{}{
		"series_id":           match.SeriesID,
		"match_number":        match.MatchNumber,
		"date":                match.Date,
		"status":              match.Status,
		"team_a_player_count": match.TeamAPlayerCount,
		"team_b_player_count": match.TeamBPlayerCount,
		"total_overs":         match.TotalOvers,
		"toss_winner":         match.TossWinner,
		"toss_type":           match.TossType,
		"batting_team":        match.BattingTeam,
		"updated_at":          match.UpdatedAt,
	}

	var result []models.Match
	_, err := r.client.From("matches").Update(matchData, "", "").Eq("id", id).ExecuteTo(&result)
	if err != nil {
		return err
	}

	if len(result) > 0 {
		*match = result[0]
	}

	return nil
}

func (r *matchRepository) Delete(ctx context.Context, id string) error {
	_, err := r.client.From("matches").Delete("", "").Eq("id", id).ExecuteTo(nil)
	return err
}

func (r *matchRepository) Count(ctx context.Context) (int64, error) {
	var result []models.Match
	_, err := r.client.From("matches").Select("*", "", false).ExecuteTo(&result)
	if err != nil {
		return 0, err
	}
	return int64(len(result)), nil
}

func (r *matchRepository) GetNextMatchNumber(ctx context.Context, seriesID string) (int, error) {
	// Get all matches for the series to find the highest match number
	var result []models.Match
	_, err := r.client.From("matches").
		Select("match_number", "", false).
		Eq("series_id", seriesID).
		ExecuteTo(&result)

	if err != nil {
		return 0, err
	}

	// If no matches exist for this series, start with match number 1
	if len(result) == 0 {
		return 1, nil
	}

	// Find the highest match number
	maxMatchNumber := 0
	for _, match := range result {
		if match.MatchNumber > maxMatchNumber {
			maxMatchNumber = match.MatchNumber
		}
	}

	// Return the next match number
	return maxMatchNumber + 1, nil
}

func (r *matchRepository) ExistsBySeriesAndMatchNumber(ctx context.Context, seriesID string, matchNumber int) (bool, error) {
	var result []models.Match
	_, err := r.client.From("matches").
		Select("id", "", false).
		Eq("series_id", seriesID).
		Eq("match_number", fmt.Sprintf("%d", matchNumber)).
		ExecuteTo(&result)

	if err != nil {
		return false, err
	}

	// Return true if any match exists with the given series ID and match number
	return len(result) > 0, nil
}
