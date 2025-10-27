package interfaces

import (
	"context"

	"spark-park-cricket-backend/internal/models"
)

// FallOfWicketsRepository defines the interface for fall of wickets data operations
type FallOfWicketsRepository interface {
	// Create creates a new fall of wickets record
	Create(ctx context.Context, fallOfWickets *models.FallOfWickets) error

	// GetByID retrieves a fall of wickets record by ID
	GetByID(ctx context.Context, id string) (*models.FallOfWickets, error)

	// GetByMatchID retrieves all fall of wickets for a specific match
	GetByMatchID(ctx context.Context, matchID string) ([]*models.FallOfWickets, error)

	// GetByInningsID retrieves all fall of wickets for a specific innings
	GetByInningsID(ctx context.Context, inningsID string) ([]*models.FallOfWickets, error)

	// GetByBallID retrieves fall of wickets record by ball ID
	GetByBallID(ctx context.Context, ballID string) (*models.FallOfWickets, error)

	// List retrieves fall of wickets records with filters
	List(ctx context.Context, filters *models.FallOfWicketsFilters) ([]*models.FallOfWickets, error)

	// Update updates an existing fall of wickets record
	Update(ctx context.Context, id string, req *models.UpdateFallOfWicketsRequest) (*models.FallOfWickets, error)

	// Delete deletes a fall of wickets record
	Delete(ctx context.Context, id string) error

	// GetSummary retrieves a summary of fall of wickets for a match/innings
	GetSummary(ctx context.Context, matchID string, inningsID *string) (*models.FallOfWicketsSummary, error)

	// GetWicketNumberForInnings gets the next wicket number for an innings
	GetWicketNumberForInnings(ctx context.Context, inningsID string) (int, error)

	// GetScoreAtWicket calculates the score when a wicket fell
	GetScoreAtWicket(ctx context.Context, inningsID string, overNumber int, ballNumber int) (int, error)
}
