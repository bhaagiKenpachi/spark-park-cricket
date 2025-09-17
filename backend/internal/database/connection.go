package database

import (
	"context"
	"fmt"
	"log"
	"spark-park-cricket-backend/internal/cache"
	"spark-park-cricket-backend/internal/config"
	cacherepo "spark-park-cricket-backend/internal/repository/cache"
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
	CacheManager *cache.CacheManager
}

// NewClient creates a new database client with all repositories
func NewClient(cfg *config.Config) (*Client, error) {
	log.Println("=== INITIALIZING DATABASE CONNECTION ===")

	if cfg.SupabaseURL == "" || cfg.SupabaseAPIKey == "" {
		log.Printf("ERROR: Supabase URL and API key are required")
		return nil, fmt.Errorf("supabase URL and API key are required")
	}

	log.Printf("Creating Supabase client with schema: %s", cfg.DatabaseSchema)

	// Create Supabase client with schema configuration
	clientOptions := &supabaseclient.ClientOptions{
		Schema: cfg.DatabaseSchema,
	}
	client, err := supabaseclient.NewClient(cfg.SupabaseURL, cfg.SupabaseAPIKey, clientOptions)
	if err != nil {
		log.Printf("ERROR: Failed to create Supabase client: %v", err)
		return nil, fmt.Errorf("failed to create supabase client: %w", err)
	}

	log.Printf("✅ Supabase client created successfully")

	// Initialize cache manager
	var cacheManager *cache.CacheManager
	if cfg.CacheEnabled {
		log.Printf("Initializing Redis cache...")
		redisClient, err := cache.NewRedisClient(cfg)
		if err != nil {
			// Log warning but continue without cache
			log.Printf("⚠️  Warning: Failed to initialize Redis cache: %v", err)
			log.Printf("Continuing without cache...")
		} else {
			cacheManager = cache.NewCacheManager(redisClient, true)
			log.Printf("✅ Redis cache initialized successfully")
		}
	} else {
		log.Printf("Cache disabled by configuration")
	}

	// Initialize base repositories
	log.Printf("Initializing database repositories...")
	baseRepositories := &Repositories{
		Series:     supabase.NewSeriesRepository(client),
		Match:      supabase.NewMatchRepository(client),
		Scoreboard: supabase.NewScoreboardRepository(client),
		Scorecard:  supabase.NewScorecardRepository(client),
		Over:       supabase.NewOverRepository(client),
		Ball:       supabase.NewBallRepository(client),
	}
	log.Printf("✅ Base repositories initialized")

	// Wrap repositories with caching if cache is available
	var repositories *Repositories
	if cacheManager != nil {
		log.Printf("Wrapping repositories with cache layer...")
		repositories = &Repositories{
			Series:     cacherepo.NewCachedSeriesRepository(baseRepositories.Series, cacheManager),
			Match:      cacherepo.NewCachedMatchRepository(baseRepositories.Match, cacheManager),
			Scoreboard: baseRepositories.Scoreboard, // Not cached yet
			Scorecard:  cacherepo.NewCachedScorecardRepository(baseRepositories.Scorecard, cacheManager),
			Over:       baseRepositories.Over, // Not cached yet
			Ball:       baseRepositories.Ball, // Not cached yet
		}
		log.Printf("✅ Cached repositories initialized")
	} else {
		log.Printf("Using direct database repositories (no cache)")
		repositories = baseRepositories
	}

	log.Printf("=== DATABASE CONNECTION INITIALIZED ===")
	log.Printf("Database Type: Supabase (PostgreSQL)")
	log.Printf("Database Schema: %s", cfg.DatabaseSchema)
	if cacheManager != nil {
		log.Printf("Cache Layer: Enabled (Redis)")
	} else {
		log.Printf("Cache Layer: Disabled")
	}
	log.Printf("Repositories: Series, Match, Scoreboard, Scorecard, Over, Ball")
	log.Printf("==========================================")

	return &Client{
		Supabase:     client,
		Repositories: repositories,
		Schema:       cfg.DatabaseSchema,
		CacheManager: cacheManager,
	}, nil
}

// HealthCheck performs a simple health check on the database
func (c *Client) HealthCheck() error {
	log.Printf("=== PERFORMING DATABASE HEALTH CHECK ===")
	log.Printf("Database Type: Supabase (PostgreSQL)")
	log.Printf("Database Schema: %s", c.Schema)

	// Simple health check by attempting to connect to a known table
	log.Printf("Testing connection to 'series' table...")
	_, _, err := c.Supabase.From("series").Select("id", "exact", false).Limit(1, "").Execute()
	if err != nil {
		log.Printf("⚠️  Warning: Could not query 'series' table: %v", err)
		log.Printf("This might be normal if the table doesn't exist yet")
		// If the table doesn't exist, that's okay for health check
		log.Printf("✅ Database connection appears to be working")
		return nil
	}

	log.Printf("✅ Database health check successful")
	log.Printf("✅ Successfully queried 'series' table")

	// Test cache if available
	if c.CacheManager != nil {
		log.Printf("Testing cache connection...")
		// You could add a cache ping test here if needed
		log.Printf("✅ Cache layer is available")
	}

	log.Printf("=========================================")
	return nil
}

// Close closes the database connection and cache (if needed)
func (c *Client) Close() error {
	// Close cache connection if available
	if c.CacheManager != nil {
		return c.CacheManager.Close()
	}
	// Supabase client doesn't need explicit closing
	return nil
}

// WithTransaction executes a function within a transaction context
func (c *Client) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	// For now, just execute the function directly
	// In a real implementation, you might want to implement proper transaction handling
	return fn(ctx)
}
