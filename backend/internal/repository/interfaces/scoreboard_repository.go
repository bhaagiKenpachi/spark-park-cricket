package interfaces

import (
	"context"
	"spark-park-cricket-backend/internal/models"
)

// ScoreboardRepository defines the interface for scoreboard data operations
type ScoreboardRepository interface {
	GetByMatchID(ctx context.Context, matchID string) (*models.LiveScoreboard, error)
	Create(ctx context.Context, scoreboard *models.LiveScoreboard) error
	Update(ctx context.Context, matchID string, scoreboard *models.LiveScoreboard) error
	Delete(ctx context.Context, matchID string) error
}
