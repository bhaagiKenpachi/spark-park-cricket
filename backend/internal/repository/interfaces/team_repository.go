package interfaces

import (
	"context"
	"spark-park-cricket-backend/internal/models"
)

// TeamRepository defines the interface for team data operations
type TeamRepository interface {
	Create(ctx context.Context, team *models.Team) error
	GetByID(ctx context.Context, id string) (*models.Team, error)
	GetAll(ctx context.Context, filters *models.TeamFilters) ([]*models.Team, error)
	Update(ctx context.Context, id string, team *models.Team) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context) (int64, error)
}
