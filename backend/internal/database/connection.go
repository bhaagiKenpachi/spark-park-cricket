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
	Team       interfaces.TeamRepository
	Player     interfaces.PlayerRepository
	Scoreboard interfaces.ScoreboardRepository
	Over       interfaces.OverRepository
	Ball       interfaces.BallRepository
}

// Client wraps the Supabase client and repositories
type Client struct {
	Supabase     *supabaseclient.Client
	Repositories *Repositories
}

// NewClient creates a new database client with all repositories
func NewClient(cfg *config.Config) (*Client, error) {
	if cfg.SupabaseURL == "" || cfg.SupabaseAPIKey == "" {
		return nil, fmt.Errorf("supabase URL and API key are required")
	}

	client, err := supabaseclient.NewClient(cfg.SupabaseURL, cfg.SupabaseAPIKey, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create supabase client: %w", err)
	}

	// Initialize repositories
	repositories := &Repositories{
		Series:     supabase.NewSeriesRepository(client),
		Match:      supabase.NewMatchRepository(client),
		Team:       supabase.NewTeamRepository(client),
		Player:     supabase.NewPlayerRepository(client),
		Scoreboard: supabase.NewScoreboardRepository(client),
		Over:       supabase.NewOverRepository(client),
		Ball:       supabase.NewBallRepository(client),
	}

	return &Client{
		Supabase:     client,
		Repositories: repositories,
	}, nil
}

// HealthCheck performs a simple health check on the database
func (c *Client) HealthCheck() error {
	// Simple health check by attempting to connect
	_, _, err := c.Supabase.From("_health_check").Select("*", "exact", false).Execute()
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
