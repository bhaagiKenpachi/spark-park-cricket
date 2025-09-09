package interfaces

import (
	"context"
	"spark-park-cricket-backend/internal/models"
)

// BallRepository defines the interface for ball data operations
type BallRepository interface {
	Create(ctx context.Context, ball *models.Ball) error
	GetByID(ctx context.Context, id string) (*models.Ball, error)
	GetByOverID(ctx context.Context, overID string) ([]*models.Ball, error)
	GetByMatchID(ctx context.Context, matchID string) ([]*models.Ball, error)
	Update(ctx context.Context, id string, ball *models.Ball) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context) (int64, error)
}
