package interfaces

import (
	"context"
	"spark-park-cricket-backend/internal/models"
)

// PlayerRepository defines the interface for player data operations
type PlayerRepository interface {
	Create(ctx context.Context, player *models.Player) error
	GetByID(ctx context.Context, id string) (*models.Player, error)
	GetAll(ctx context.Context, filters *models.PlayerFilters) ([]*models.Player, error)
	Update(ctx context.Context, id string, player *models.Player) error
	Delete(ctx context.Context, id string) error
	GetByTeamID(ctx context.Context, teamID string) ([]*models.Player, error)
	Count(ctx context.Context) (int64, error)
}
