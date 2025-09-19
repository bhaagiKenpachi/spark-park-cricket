package interfaces

import (
	"context"
	"spark-park-cricket-backend/internal/models"
)

// ScorecardRepository defines the interface for scorecard data operations
type ScorecardRepository interface {
	// Innings operations
	CreateInnings(ctx context.Context, innings *models.Innings) error
	GetInningsByMatchID(ctx context.Context, matchID string) ([]*models.Innings, error)
	GetInningsByMatchAndNumber(ctx context.Context, matchID string, inningsNumber int) (*models.Innings, error)
	UpdateInnings(ctx context.Context, innings *models.Innings) error
	CompleteInnings(ctx context.Context, inningsID string) error

	// Over operations
	CreateOver(ctx context.Context, over *models.ScorecardOver) error
	GetOverByInningsAndNumber(ctx context.Context, inningsID string, overNumber int) (*models.ScorecardOver, error)
	GetCurrentOver(ctx context.Context, inningsID string) (*models.ScorecardOver, error)
	GetOversByInnings(ctx context.Context, inningsID string) ([]*models.ScorecardOver, error)
	UpdateOver(ctx context.Context, over *models.ScorecardOver) error
	CompleteOver(ctx context.Context, overID string) error

	// Ball operations
	CreateBall(ctx context.Context, ball *models.ScorecardBall) error
	GetBallsByOver(ctx context.Context, overID string) ([]*models.ScorecardBall, error)
	GetLastBall(ctx context.Context, overID string) (*models.ScorecardBall, error)
	DeleteBall(ctx context.Context, ballID string) error

	// Scorecard operations
	GetScorecard(ctx context.Context, matchID string) (*models.ScorecardResponse, error)
	StartScoring(ctx context.Context, matchID string) error
}
