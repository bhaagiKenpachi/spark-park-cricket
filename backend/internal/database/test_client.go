package database

import (
	"fmt"
	"log"
	"spark-park-cricket-backend/internal/cache"
	"spark-park-cricket-backend/internal/config"
	cacherepo "spark-park-cricket-backend/internal/repository/cache"
	"spark-park-cricket-backend/internal/repository/supabase"

	supabaseclient "github.com/supabase-community/supabase-go"
)

// NewTestClient creates a new database client for testing with test schema
func NewTestClient(cfg *config.TestConfig) (*Client, error) {
	if cfg.SupabaseURL == "" || cfg.SupabaseAPIKey == "" {
		return nil, fmt.Errorf("supabase URL and API key are required")
	}

	// Create Supabase client with test schema
	clientOptions := &supabaseclient.ClientOptions{
		Schema:  cfg.TestSchema,
		Headers: cfg.GetSupabaseHeaders(),
	}
	client, err := supabaseclient.NewClient(cfg.SupabaseURL, cfg.SupabaseAPIKey, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create supabase client: %w", err)
	}

	// Initialize cache manager for tests (in-memory cache)
	cacheManager := cache.NewCacheManager(nil, false) // No Redis for tests, use in-memory
	log.Printf("Initialized in-memory cache for tests")

	// Initialize base repositories
	matchRepo := supabase.NewMatchRepository(client)
	baseRepositories := &Repositories{
		Series:     supabase.NewSeriesRepository(client),
		Match:      matchRepo,
		Scoreboard: supabase.NewScoreboardRepository(client),
		Scorecard:  supabase.NewScorecardRepository(client, "testing_db", matchRepo),
		Over:       supabase.NewOverRepository(client),
		Ball:       supabase.NewBallRepository(client),
		User:       supabase.NewUserRepository(client),
		Vote:       supabase.NewVoteRepository(client),
	}

	// Wrap repositories with caching (same as production)
	repositories := &Repositories{
		Series:     cacherepo.NewCachedSeriesRepository(baseRepositories.Series, cacheManager),
		Match:      cacherepo.NewCachedMatchRepository(baseRepositories.Match, cacheManager),
		Scoreboard: baseRepositories.Scoreboard, // Not cached yet
		Scorecard:  cacherepo.NewCachedScorecardRepository(baseRepositories.Scorecard, cacheManager),
		Over:       baseRepositories.Over, // Not cached yet
		Ball:       baseRepositories.Ball, // Not cached yet
		User:       baseRepositories.User, // Not cached yet
		Vote:       baseRepositories.Vote, // Not cached yet
	}

	return &Client{
		Supabase:     client,
		Repositories: repositories,
		Schema:       cfg.TestSchema,
		CacheManager: cacheManager,
	}, nil
}

// SetupTestSchema creates the test schema and tables if they don't exist
func SetupTestSchema(cfg *config.TestConfig) error {
	// This would typically involve running migrations on the test schema
	// For now, we'll assume the test schema is set up manually
	// In a real implementation, you might want to:
	// 1. Create the test schema
	// 2. Run migrations on the test schema
	// 3. Set up test data

	return nil
}

// CleanupTestSchema cleans up test data after tests
func CleanupTestSchema(cfg *config.TestConfig) error {
	// This would typically involve:
	// 1. Truncating test tables
	// 2. Resetting sequences
	// 3. Cleaning up any test data

	return nil
}
