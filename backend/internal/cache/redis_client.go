package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"spark-park-cricket-backend/internal/config"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient wraps the Redis client with cricket-specific caching methods
type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisClient creates a new Redis client
func NewRedisClient(cfg *config.Config) (*RedisClient, error) {
	if !cfg.CacheEnabled {
		return nil, fmt.Errorf("caching is disabled")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisURL,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	ctx := context.Background()

	// Test connection
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisClient{
		client: rdb,
		ctx:    ctx,
	}, nil
}

// Set stores a value in Redis with TTL
func (r *RedisClient) Set(key string, value interface{}, ttl time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return r.client.Set(r.ctx, key, jsonData, ttl).Err()
}

// Get retrieves a value from Redis
func (r *RedisClient) Get(key string, dest interface{}) error {
	val, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("key not found: %s", key)
		}
		return fmt.Errorf("failed to get value: %w", err)
	}

	return json.Unmarshal([]byte(val), dest)
}

// Delete removes a key from Redis
func (r *RedisClient) Delete(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

// Exists checks if a key exists in Redis
func (r *RedisClient) Exists(key string) (bool, error) {
	count, err := r.client.Exists(r.ctx, key).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// SetNX sets a key only if it doesn't exist (atomic operation)
func (r *RedisClient) SetNX(key string, value interface{}, ttl time.Duration) (bool, error) {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return false, fmt.Errorf("failed to marshal value: %w", err)
	}

	return r.client.SetNX(r.ctx, key, jsonData, ttl).Result()
}

// Increment atomically increments a counter
func (r *RedisClient) Increment(key string) (int64, error) {
	return r.client.Incr(r.ctx, key).Result()
}

// Expire sets expiration time for a key
func (r *RedisClient) Expire(key string, ttl time.Duration) error {
	return r.client.Expire(r.ctx, key, ttl).Err()
}

// Close closes the Redis connection
func (r *RedisClient) Close() error {
	return r.client.Close()
}

// HealthCheck performs a health check on Redis
func (r *RedisClient) HealthCheck() error {
	_, err := r.client.Ping(r.ctx).Result()
	return err
}

// Cache key generators for different data types
func (r *RedisClient) GetSeriesKey(seriesID string) string {
	return fmt.Sprintf("series:%s", seriesID)
}

func (r *RedisClient) GetMatchKey(matchID string) string {
	return fmt.Sprintf("match:%s", matchID)
}

func (r *RedisClient) GetScorecardKey(matchID string) string {
	return fmt.Sprintf("scorecard:%s", matchID)
}

func (r *RedisClient) GetScorecardVersionKey(matchID string) string {
	return fmt.Sprintf("scorecard:version:%s", matchID)
}

func (r *RedisClient) GetMatchesBySeriesKey(seriesID string) string {
	return fmt.Sprintf("matches:series:%s", seriesID)
}

// Cache TTL constants
const (
	// Static data (series, matches) - cache for 24 hours
	StaticDataTTL = 24 * time.Hour

	// Scorecard data - cache for 1 hour (refreshed on ball updates)
	ScorecardTTL = 1 * time.Hour

	// Match lists - cache for 30 minutes
	MatchListTTL = 30 * time.Minute

	// Version counters - cache for 1 hour
	VersionTTL = 1 * time.Hour
)
