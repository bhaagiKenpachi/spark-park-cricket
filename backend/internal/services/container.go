package services

import (
	"spark-park-cricket-backend/internal/config"
	"spark-park-cricket-backend/internal/database"
	"spark-park-cricket-backend/internal/interfaces"
	"spark-park-cricket-backend/internal/monitoring"
)

// Container holds all service instances
type Container struct {
	DBClient   *database.Client
	Series     *SeriesService
	Match      *MatchService
	MatchState MatchStateServiceInterface
	Scorecard  interfaces.ScorecardServiceInterface
	// Authentication services
	AuthService    *AuthService
	SessionService *SessionService
	// User services
	UserService UserServiceInterface
	// Voting services
	VoteService     VoteServiceInterface
	VoteTeamService VoteTeamServiceInterface
	GroupService    GroupServiceInterface
	// Fall of wickets service
	FallOfWicketsService *FallOfWicketsService
	// Monitoring
	Metrics *monitoring.Metrics
}

// NewContainer creates a new service container with all services
func NewContainer(dbClient *database.Client, cfg *config.Config) *Container {
	// Initialize Prometheus metrics
	metrics := monitoring.NewMetrics()

	// Create base scorecard service
	scorecardService := NewScorecardService(dbClient.Repositories.Scorecard, dbClient.Repositories.Match, dbClient.Repositories.FallOfWickets, metrics, dbClient.CacheManager)

	// Create authentication services
	sessionService := NewSessionService(dbClient.Repositories.User, cfg)
	authService := NewAuthService(cfg, dbClient.Repositories.User, sessionService)

	// Create match state service
	matchStateService := NewMatchStateService(
		dbClient.Repositories.Match,
		dbClient.Repositories.Scoreboard,
		dbClient.Repositories.Over,
		dbClient.Repositories.Ball,
	)

	// Set metrics on cache manager if available
	if dbClient.CacheManager != nil {
		dbClient.CacheManager.SetMetrics(metrics)
	}

	// Create group service first
	groupService := NewGroupService(dbClient.Repositories.Group, dbClient.Repositories.User, dbClient.Repositories.Vote)

	// Create container
	container := &Container{
		DBClient:   dbClient,
		Series:     NewSeriesService(dbClient.Repositories.Series),
		Match:      NewMatchService(dbClient.Repositories.Match, dbClient.Repositories.Series),
		MatchState: matchStateService,
		Scorecard:  scorecardService,
		// Authentication services
		AuthService:    authService,
		SessionService: sessionService,
		// User services
		UserService: NewUserService(dbClient.Repositories.User),
		// Voting services
		VoteService:     NewVoteService(dbClient.Repositories.Vote, groupService),
		VoteTeamService: NewVoteTeamService(dbClient.Repositories.VoteTeam, dbClient.Repositories.Vote),
		GroupService:    groupService,
		// Fall of wickets service
		FallOfWicketsService: NewFallOfWicketsService(
			dbClient.Repositories.FallOfWickets,
			dbClient.Repositories.Scorecard,
			dbClient.Repositories.Match,
		),
		// Monitoring
		Metrics: metrics,
	}

	return container
}
