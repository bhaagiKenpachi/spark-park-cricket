package services

import (
	"spark-park-cricket-backend/internal/config"
	"spark-park-cricket-backend/internal/database"
	"spark-park-cricket-backend/internal/interfaces"
	"spark-park-cricket-backend/internal/monitoring"
)

// Container holds all service instances
type Container struct {
	Series    *SeriesService
	Match     *MatchService
	Scorecard interfaces.ScorecardServiceInterface
	// Authentication services
	AuthService    *AuthService
	SessionService *SessionService
	// Monitoring
	Metrics *monitoring.Metrics
}

// NewContainer creates a new service container with all services
func NewContainer(dbClient *database.Client, cfg *config.Config) *Container {
	// Initialize Prometheus metrics
	metrics := monitoring.NewMetrics()

	// Create base scorecard service
	scorecardService := NewScorecardService(dbClient.Repositories.Scorecard, dbClient.Repositories.Match, metrics, dbClient.CacheManager)

	// Create authentication services
	sessionService := NewSessionService(dbClient.Repositories.User, cfg)
	authService := NewAuthService(cfg, dbClient.Repositories.User, sessionService)

	// Create container
	container := &Container{
		Series:    NewSeriesService(dbClient.Repositories.Series),
		Match:     NewMatchService(dbClient.Repositories.Match, dbClient.Repositories.Series),
		Scorecard: scorecardService,
		// Authentication services
		AuthService:    authService,
		SessionService: sessionService,
		// Monitoring
		Metrics: metrics,
	}

	return container
}
