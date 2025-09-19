package database

import (
	"fmt"
	"spark-park-cricket-backend/internal/config"
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

	// Initialize repositories
	repositories := &Repositories{
		Series:     supabase.NewSeriesRepository(client),
		Match:      supabase.NewMatchRepository(client),
		Scoreboard: supabase.NewScoreboardRepository(client),
		Scorecard:  supabase.NewScorecardRepository(client),
		Over:       supabase.NewOverRepository(client),
		Ball:       supabase.NewBallRepository(client),
		User:       supabase.NewUserRepository(client),
	}

	return &Client{
		Supabase:     client,
		Repositories: repositories,
		Schema:       cfg.TestSchema,
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
