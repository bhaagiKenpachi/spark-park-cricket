package cache

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"spark-park-cricket-backend/internal/config"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient wraps the Redis client with cricket-specific caching methods
type RedisClient struct {
	client  *redis.Client
	ctx     context.Context
	metrics MetricsRecorder
}

// NewRedisClient creates a new Redis client with fallback support
func NewRedisClient(cfg *config.Config) (*RedisClient, error) {
	if !cfg.CacheEnabled {
		log.Printf("Cache disabled by configuration")
		return nil, fmt.Errorf("caching is disabled")
	}

	log.Printf("Configuring Redis client...")
	log.Printf("Redis URL: %s", cfg.RedisURL)
	log.Printf("Redis DB: %d", cfg.RedisDB)
	log.Printf("Redis TLS: %t", cfg.RedisUseTLS)
	if cfg.RedisPassword != "" {
		log.Printf("Redis Password: ***")
	} else {
		log.Printf("Redis Password: NOT SET")
	}

	// Configure Redis client options with timeout and retry settings
	options := &redis.Options{
		Addr:         cfg.RedisURL,
		Password:     cfg.RedisPassword,
		DB:           cfg.RedisDB,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 1,
		MaxRetries:   3,
	}

	// Enable TLS only if configured
	if cfg.RedisUseTLS {
		log.Printf("Enabling TLS for Redis connection")
		options.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	log.Printf("Creating Redis client...")
	rdb := redis.NewClient(options)

	ctx := context.Background()

	// Test connection with timeout
	log.Printf("Testing Redis connection...")
	testCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(testCtx).Result()
	if err != nil {
		log.Printf("❌ Failed to connect to Redis: %v", err)
		log.Printf("⚠️  Continuing without cache - system will work with database-only mode")
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Printf("✅ Redis connection successful")

	return &RedisClient{
		client:  rdb,
		ctx:     ctx,
		metrics: nil, // Will be set later via SetMetrics
	}, nil
}

// SetMetrics sets the metrics recorder for the Redis client
func (r *RedisClient) SetMetrics(metrics MetricsRecorder) {
	r.metrics = metrics
}

// Set stores a value in Redis with TTL
func (r *RedisClient) Set(key string, value interface{}, ttl time.Duration) error {
	start := time.Now()
	jsonData, err := json.Marshal(value)
	if err != nil {
		if r.metrics != nil {
			r.metrics.RecordCacheError("set", "marshal_error")
		}
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	// Add timeout to prevent hanging
	ctx, cancel := context.WithTimeout(r.ctx, 3*time.Second)
	defer cancel()

	err = r.client.Set(ctx, key, jsonData, ttl).Err()
	if err != nil {
		log.Printf("⚠️  Cache SET failed for key %s: %v", key, err)
		if r.metrics != nil {
			r.metrics.RecordCacheError("set", "redis_error")
			r.metrics.RecordCacheOperation("set", "redis", time.Since(start))
		}
		return fmt.Errorf("failed to set cache key %s: %w", key, err)
	}

	if r.metrics != nil {
		r.metrics.RecordCacheOperation("set", "redis", time.Since(start))
		r.metrics.RecordCacheKeySize(extractKeyPattern(key), len(jsonData))
	}

	return nil
}

// Get retrieves a value from Redis
func (r *RedisClient) Get(key string, dest interface{}) error {
	start := time.Now()
	// Add timeout to prevent hanging
	ctx, cancel := context.WithTimeout(r.ctx, 3*time.Second)
	defer cancel()

	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			if r.metrics != nil {
				r.metrics.RecordCacheMiss("redis", extractKeyPattern(key))
				r.metrics.RecordCacheOperation("get", "redis", time.Since(start))
			}
			return fmt.Errorf("key not found: %s", key)
		}
		log.Printf("⚠️  Cache GET failed for key %s: %v", key, err)
		if r.metrics != nil {
			r.metrics.RecordCacheError("get", "redis_error")
			r.metrics.RecordCacheOperation("get", "redis", time.Since(start))
		}
		return fmt.Errorf("failed to get cache key %s: %w", key, err)
	}

	err = json.Unmarshal([]byte(val), dest)
	if err != nil {
		if r.metrics != nil {
			r.metrics.RecordCacheError("get", "unmarshal_error")
		}
		return err
	}

	if r.metrics != nil {
		r.metrics.RecordCacheHit("redis", extractKeyPattern(key))
		r.metrics.RecordCacheOperation("get", "redis", time.Since(start))
		r.metrics.RecordCacheKeySize(extractKeyPattern(key), len(val))
	}

	return nil
}

// Delete removes a key from Redis
func (r *RedisClient) Delete(key string) error {
	start := time.Now()
	result := r.client.Del(r.ctx, key)
	err := result.Err()
	if err != nil {
		if r.metrics != nil {
			r.metrics.RecordCacheError("delete", "redis_error")
			r.metrics.RecordCacheOperation("delete", "redis", time.Since(start))
		}
		return err
	}
	if r.metrics != nil {
		r.metrics.RecordCacheOperation("delete", "redis", time.Since(start))
	}
	return nil
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
	// Create a new context with timeout for health check
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.client.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("redis health check failed: %w", err)
	}

	return nil
}

// GetSeriesKey generates a cache key for series data
func (r *RedisClient) GetSeriesKey(seriesID string) string {
	return fmt.Sprintf("series:%s", seriesID)
}

// GetMatchKey generates a cache key for match data
func (r *RedisClient) GetMatchKey(matchID string) string {
	return fmt.Sprintf("match:%s", matchID)
}

// GetScorecardKey generates a cache key for scorecard data
func (r *RedisClient) GetScorecardKey(matchID string) string {
	return fmt.Sprintf("scorecard:%s", matchID)
}

// GetScorecardVersionKey generates a cache key for scorecard version data
func (r *RedisClient) GetScorecardVersionKey(matchID string) string {
	return fmt.Sprintf("scorecard:version:%s", matchID)
}

// GetMatchesBySeriesKey generates a cache key for matches by series data
func (r *RedisClient) GetMatchesBySeriesKey(seriesID string) string {
	return fmt.Sprintf("matches:series:%s", seriesID)
}

// DeletePattern removes all keys matching a pattern using SCAN
func (r *RedisClient) DeletePattern(pattern string) error {
	ctx, cancel := context.WithTimeout(r.ctx, 10*time.Second)
	defer cancel()

	iter := r.client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		if err := r.client.Del(ctx, key).Err(); err != nil {
			log.Printf("⚠️  Failed to delete cache key %s: %v", key, err)
		}
	}

	if err := iter.Err(); err != nil {
		return fmt.Errorf("failed to scan pattern %s: %w", pattern, err)
	}

	return nil
}

// Cache TTL constants
const (
	// Static data (series, matches, votes) - cache for 24 hours
	StaticDataTTL = 24 * time.Hour

	// Scorecard data - cache for 1 hour (refreshed on ball updates)
	ScorecardTTL = 1 * time.Hour

	// Match lists - cache for 30 minutes
	MatchListTTL = 30 * time.Minute

	// Live/dynamic data (vote results, user votes) - cache for 5 minutes
	LiveDataTTL = 5 * time.Minute

	// Version counters - cache for 1 hour
	VersionTTL = 1 * time.Hour
)

// PublishToStream publishes an event to a Redis stream
func (r *RedisClient) PublishToStream(streamKey string, data map[string]interface{}) (string, error) {
	ctx, cancel := context.WithTimeout(r.ctx, 3*time.Second)
	defer cancel()

	// Convert data to string map as required by Redis streams
	stringData := make(map[string]interface{})
	for k, v := range data {
		switch val := v.(type) {
		case string:
			stringData[k] = val
		default:
			// Marshal non-string values to JSON
			jsonBytes, err := json.Marshal(val)
			if err != nil {
				return "", fmt.Errorf("failed to marshal field %s: %w", k, err)
			}
			stringData[k] = string(jsonBytes)
		}
	}

	// Use XADD to add event to stream with automatic ID generation
	id, err := r.client.XAdd(ctx, &redis.XAddArgs{
		Stream: streamKey,
		MaxLen: 1000, // Keep last 1000 events per stream
		Approx: true, // Use approximate trimming for better performance
		Values: stringData,
	}).Result()

	if err != nil {
		log.Printf("⚠️  Failed to publish to stream %s: %v", streamKey, err)
		return "", fmt.Errorf("failed to publish to stream: %w", err)
	}

	return id, nil
}

// ReadFromStream reads events from a Redis stream
func (r *RedisClient) ReadFromStream(ctx context.Context, streamKey string, lastID string, count int64, block time.Duration) ([]redis.XMessage, error) {
	streams, err := r.client.XRead(ctx, &redis.XReadArgs{
		Streams: []string{streamKey, lastID},
		Count:   count,
		Block:   block,
	}).Result()

	if err != nil {
		if err == redis.Nil {
			// No new messages
			return []redis.XMessage{}, nil
		}
		return nil, fmt.Errorf("failed to read from stream: %w", err)
	}

	if len(streams) == 0 {
		return []redis.XMessage{}, nil
	}

	return streams[0].Messages, nil
}

// GetStreamKey generates a stream key for match ball events
func (r *RedisClient) GetStreamKey(matchID string) string {
	return fmt.Sprintf("stream:match:%s:balls", matchID)
}

// ReadStreamEvents reads all events from a Redis stream
func (r *RedisClient) ReadStreamEvents(streamKey string, start, end string) ([]redis.XMessage, error) {
	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()

	// Read all events from the stream using XRANGE
	streams, err := r.client.XRange(ctx, streamKey, start, end).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to read stream events: %w", err)
	}

	// Convert to XMessage format for consistency
	var messages []redis.XMessage
	for _, stream := range streams {
		messages = append(messages, redis.XMessage{
			ID:     stream.ID,
			Values: stream.Values,
		})
	}

	return messages, nil
}

// extractKeyPattern extracts the pattern from a cache key for metrics
func extractKeyPattern(key string) string {
	if len(key) == 0 {
		return "unknown"
	}

	// Common patterns
	patterns := map[string]string{
		"series:":    "series",
		"match:":     "match",
		"scorecard:": "scorecard",
		"vote:":      "vote",
		"vote_team:": "vote_team",
		"team:":      "team",
		"user:":      "user",
	}

	for prefix, pattern := range patterns {
		if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
			return pattern
		}
	}

	return "other"
}
