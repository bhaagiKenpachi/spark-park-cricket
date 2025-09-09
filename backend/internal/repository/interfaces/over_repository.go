package interfaces

import (
	"context"
	"spark-park-cricket-backend/internal/models"
)

// OverRepository defines the interface for over data operations
type OverRepository interface {
	Create(ctx context.Context, over *models.Over) error
	GetByID(ctx context.Context, id string) (*models.Over, error)
	GetByMatchID(ctx context.Context, matchID string) ([]*models.Over, error)
	Update(ctx context.Context, id string, over *models.Over) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context) (int64, error)
}
