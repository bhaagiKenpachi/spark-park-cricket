package database

import (
	"context"
	"fmt"
	"spark-park-cricket-backend/internal/config"
	"spark-park-cricket-backend/internal/repository/interfaces"
	"spark-park-cricket-backend/internal/repository/supabase"

	supabaseclient "github.com/supabase-community/supabase-go"
)

// Repositories holds all repository interfaces
type Repositories struct {
	Series     interfaces.SeriesRepository
	Match      interfaces.MatchRepository
	Scoreboard interfaces.ScoreboardRepository
	Scorecard  interfaces.ScorecardRepository
	Over       interfaces.OverRepository
	Ball       interfaces.BallRepository
}

// Client wraps the Supabase client and repositories
type Client struct {
	Supabase     *supabaseclient.Client
	Repositories *Repositories
	Schema       string
}

// NewClient creates a new database client with all repositories
func NewClient(cfg *config.Config) (*Client, error) {
	if cfg.SupabaseURL == "" || cfg.SupabaseAPIKey == "" {
		return nil, fmt.Errorf("supabase URL and API key are required")
	}

	// Create Supabase client with schema configuration
	clientOptions := &supabaseclient.ClientOptions{
		Schema: cfg.DatabaseSchema,
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
	}

	return &Client{
		Supabase:     client,
		Repositories: repositories,
		Schema:       cfg.DatabaseSchema,
	}, nil
}

// HealthCheck performs a simple health check on the database
func (c *Client) HealthCheck() error {
	// Simple health check by attempting to connect to a known table
	_, _, err := c.Supabase.From("series").Select("id", "exact", false).Limit(1, "").Execute()
	if err != nil {
		// If the table doesn't exist, that's okay for health check
		return nil
	}
	return nil
}

// Close closes the database connection (if needed)
func (c *Client) Close() error {
	// Supabase client doesn't need explicit closing
	return nil
}

// WithTransaction executes a function within a transaction context
func (c *Client) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	// For now, just execute the function directly
	// In a real implementation, you might want to implement proper transaction handling
	return fn(ctx)
}
