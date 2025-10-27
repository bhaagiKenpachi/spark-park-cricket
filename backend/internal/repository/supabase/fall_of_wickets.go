package supabase

import (
	"context"
	"fmt"

	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/repository/interfaces"

	"github.com/supabase-community/postgrest-go"
	"github.com/supabase-community/supabase-go"
)

// fallOfWicketsRepository implements the FallOfWicketsRepository interface using Supabase
type fallOfWicketsRepository struct {
	client *supabase.Client
}

// NewFallOfWicketsRepository creates a new fall of wickets repository
func NewFallOfWicketsRepository(client *supabase.Client) interfaces.FallOfWicketsRepository {
	return &fallOfWicketsRepository{
		client: client,
	}
}

// Create creates a new fall of wickets record
func (r *fallOfWicketsRepository) Create(ctx context.Context, fallOfWickets *models.FallOfWickets) error {
	var result []models.FallOfWickets
	_, err := r.client.From("dev_v1.fall_of_wickets").Insert(fallOfWickets, false, "", "", "").ExecuteTo(&result)
	if err != nil {
		return err
	}

	if len(result) > 0 {
		*fallOfWickets = result[0]
	}

	return nil
}

// GetByID retrieves a fall of wickets record by ID
func (r *fallOfWicketsRepository) GetByID(ctx context.Context, id string) (*models.FallOfWickets, error) {
	var result []models.FallOfWickets
	_, err := r.client.From("dev_v1.fall_of_wickets").Select("*", "", false).Eq("id", id).ExecuteTo(&result)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("fall of wickets record not found")
	}

	return &result[0], nil
}

// List retrieves fall of wickets records with optional filters
func (r *fallOfWicketsRepository) List(ctx context.Context, filters *models.FallOfWicketsFilters) ([]*models.FallOfWickets, error) {
	query := r.client.From("dev_v1.fall_of_wickets").Select("*", "", false)

	if filters.MatchID != nil && *filters.MatchID != "" {
		query = query.Eq("match_id", *filters.MatchID)
	}
	if filters.InningsID != nil && *filters.InningsID != "" {
		query = query.Eq("innings_id", *filters.InningsID)
	}

	query = query.Order("wicket_number", &postgrest.OrderOpts{
		Ascending: true,
	})

	var result []models.FallOfWickets
	_, err := query.ExecuteTo(&result)
	if err != nil {
		return nil, err
	}

	// Convert to slice of pointers
	var records []*models.FallOfWickets
	for i := range result {
		records = append(records, &result[i])
	}

	return records, nil
}

// Update updates a fall of wickets record
func (r *fallOfWicketsRepository) Update(ctx context.Context, id string, req *models.UpdateFallOfWicketsRequest) (*models.FallOfWickets, error) {
	// First get the existing record
	existing, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update the fields
	if req.Score != nil {
		existing.Score = *req.Score
	}

	// Update in database
	_, err = r.client.From("dev_v1.fall_of_wickets").Update(existing, "", "").Eq("id", id).ExecuteTo(&[]models.FallOfWickets{})
	if err != nil {
		return nil, err
	}

	return existing, nil
}

// Delete deletes a fall of wickets record
func (r *fallOfWicketsRepository) Delete(ctx context.Context, id string) error {
	_, err := r.client.From("dev_v1.fall_of_wickets").Delete("", "").Eq("id", id).ExecuteTo(&[]models.FallOfWickets{})
	return err
}

// GetByMatchID retrieves all fall of wickets records for a match
func (r *fallOfWicketsRepository) GetByMatchID(ctx context.Context, matchID string) ([]*models.FallOfWickets, error) {
	var result []models.FallOfWickets
	_, err := r.client.From("dev_v1.fall_of_wickets").Select("*", "", false).Eq("match_id", matchID).Order("wicket_number", &postgrest.OrderOpts{
		Ascending: true,
	}).ExecuteTo(&result)
	if err != nil {
		return nil, err
	}

	// Convert to slice of pointers
	var records []*models.FallOfWickets
	for i := range result {
		records = append(records, &result[i])
	}

	return records, nil
}

// GetByInningsID retrieves all fall of wickets records for an innings
func (r *fallOfWicketsRepository) GetByInningsID(ctx context.Context, inningsID string) ([]*models.FallOfWickets, error) {
	var result []models.FallOfWickets
	_, err := r.client.From("dev_v1.fall_of_wickets").Select("*", "", false).Eq("innings_id", inningsID).Order("wicket_number", &postgrest.OrderOpts{
		Ascending: true,
	}).ExecuteTo(&result)
	if err != nil {
		return nil, err
	}

	// Convert to slice of pointers
	var records []*models.FallOfWickets
	for i := range result {
		records = append(records, &result[i])
	}

	return records, nil
}

// GetByBallID retrieves fall of wickets record for a specific ball
func (r *fallOfWicketsRepository) GetByBallID(ctx context.Context, ballID string) (*models.FallOfWickets, error) {
	var result []models.FallOfWickets
	_, err := r.client.From("dev_v1.fall_of_wickets").Select("*", "", false).Eq("ball_id", ballID).ExecuteTo(&result)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("fall of wickets record not found for ball %s", ballID)
	}

	return &result[0], nil
}

// GetWicketNumberForInnings gets the next wicket number for an innings
func (r *fallOfWicketsRepository) GetWicketNumberForInnings(ctx context.Context, inningsID string) (int, error) {
	var result []struct {
		Count int `json:"count"`
	}

	_, err := r.client.From("dev_v1.fall_of_wickets").Select("count", "", false).Eq("innings_id", inningsID).ExecuteTo(&result)
	if err != nil {
		return 0, err
	}

	if len(result) == 0 {
		return 1, nil // First wicket
	}

	return result[0].Count + 1, nil
}

// GetScoreAtWicket calculates the score at the time of wicket
func (r *fallOfWicketsRepository) GetScoreAtWicket(ctx context.Context, inningsID string, overNumber int, ballNumber int) (int, error) {
	// This is a simplified calculation - in a real implementation, you'd need to
	// calculate the actual score based on all balls up to this point
	// For now, we'll return a placeholder value
	return 0, nil
}

// GetSummary retrieves a summary of fall of wickets for a match/innings
func (r *fallOfWicketsRepository) GetSummary(ctx context.Context, matchID string, inningsID *string) (*models.FallOfWicketsSummary, error) {
	var records []*models.FallOfWickets
	var err error

	if inningsID != nil && *inningsID != "" {
		records, err = r.GetByInningsID(ctx, *inningsID)
	} else {
		records, err = r.GetByMatchID(ctx, matchID)
	}

	if err != nil {
		return nil, err
	}

	summary := &models.FallOfWicketsSummary{
		MatchID: matchID,
		InningsID: func() string {
			if inningsID != nil {
				return *inningsID
			}
			return ""
		}(),
		TotalWickets: len(records),
		Wickets:      make([]models.WicketFall, len(records)),
	}

	for i, record := range records {
		summary.Wickets[i] = models.WicketFall{
			WicketNumber: record.WicketNumber,
			Score:        record.Score,
			OverNumber:   record.OverNumber,
			BallNumber:   record.BallNumber,
			OverPosition: fmt.Sprintf("%d.%d", record.OverNumber, record.BallNumber),
		}
	}

	return summary, nil
}
