package cache

import (
	"encoding/json"
	"fmt"
	"time"
)

// MetricsRecorder interface for recording cache metrics
type MetricsRecorder interface {
	RecordCacheHit(cacheType, keyPattern string)
	RecordCacheMiss(cacheType, keyPattern string)
	RecordCacheOperation(operation, cacheType string, duration time.Duration)
	RecordCacheKeySize(keyPattern string, size int)
	RecordCacheError(operation, errorType string)
}

// CacheInterface defines the contract for caching operations
type CacheInterface interface {
	// Basic operations
	Set(key string, value interface{}, ttl time.Duration) error
	Get(key string, dest interface{}) error
	Delete(key string) error
	Exists(key string) (bool, error)
	SetNX(key string, value interface{}, ttl time.Duration) (bool, error)

	// Atomic operations
	Increment(key string) (int64, error)
	Expire(key string, ttl time.Duration) error

	// Connection management
	Close() error
	HealthCheck() error

	// Cricket-specific key generators
	GetSeriesKey(seriesID string) string
	GetMatchKey(matchID string) string
	GetScorecardKey(matchID string) string
	GetScorecardVersionKey(matchID string) string
	GetMatchesBySeriesKey(seriesID string) string

	// Pattern-based operations
	DeletePattern(pattern string) error
}

// CacheManager handles cache operations with fallback to database
type CacheManager struct {
	cache   CacheInterface
	enabled bool
	metrics MetricsRecorder
}

// NewCacheManager creates a new cache manager
func NewCacheManager(cache CacheInterface, enabled bool) *CacheManager {
	return &CacheManager{
		cache:   cache,
		enabled: enabled,
		metrics: nil, // Will be set later
	}
}

// SetMetrics sets the metrics recorder for the cache manager
func (cm *CacheManager) SetMetrics(metrics MetricsRecorder) {
	cm.metrics = metrics
	// Also set metrics on the underlying cache if it's a RedisClient
	if cm.cache != nil {
		if redisClient, ok := cm.cache.(*RedisClient); ok {
			redisClient.SetMetrics(metrics)
		}
	}
}

// GetOrSet retrieves from cache or sets from database function
func (cm *CacheManager) GetOrSet(key string, dest interface{}, ttl time.Duration, dbFunc func() (interface{}, error)) error {
	if !cm.enabled {
		// Cache disabled, call database function directly
		fmt.Printf("🔍 [CACHE] Cache disabled, calling database directly for key: %s\n", key)
		value, err := dbFunc()
		if err != nil {
			return err
		}

		// Copy value to destination
		return copyValue(value, dest)
	}

	// Try to get from cache first
	fmt.Printf("🔍 [CACHE] Attempting to get from cache for key: %s\n", key)
	err := cm.cache.Get(key, dest)
	if err == nil {
		// Cache hit
		fmt.Printf("✅ [CACHE] Cache HIT for key: %s\n", key)
		return nil
	}

	// Cache miss or error, get from database
	fmt.Printf("❌ [CACHE] Cache MISS for key: %s, error: %v\n", key, err)
	fmt.Printf("🔍 [CACHE] Calling database function for key: %s\n", key)
	value, err := dbFunc()
	if err != nil {
		return err
	}

	// Set in cache for next time (ignore cache errors)
	fmt.Printf("💾 [CACHE] Setting cache for key: %s with TTL: %v\n", key, ttl)
	if cacheErr := cm.cache.Set(key, value, ttl); cacheErr != nil {
		// Log cache error but don't fail the operation
		fmt.Printf("⚠️  [CACHE] Cache SET failed for key %s: %v (continuing without cache)\n", key, cacheErr)
	} else {
		fmt.Printf("✅ [CACHE] Cache SET successful for key: %s\n", key)
	}

	// Copy value to destination
	return copyValue(value, dest)
}

// Invalidate removes a key from cache
func (cm *CacheManager) Invalidate(key string) error {
	if !cm.enabled {
		fmt.Printf("🔍 [CACHE] Cache disabled, skipping invalidation for key: %s\n", key)
		return nil
	}

	fmt.Printf("🗑️  [CACHE] Invalidating cache key: %s\n", key)
	err := cm.cache.Delete(key)
	if err != nil {
		// Log cache error but don't fail the operation
		fmt.Printf("⚠️  [CACHE] Cache DELETE failed for key %s: %v (continuing without cache)\n", key, err)
		return nil
	}

	fmt.Printf("✅ [CACHE] Cache invalidation successful for key: %s\n", key)
	return nil
}

// InvalidatePattern removes all keys matching a pattern
func (cm *CacheManager) InvalidatePattern(pattern string) error {
	if !cm.enabled {
		fmt.Printf("🔍 [CACHE] Cache disabled, skipping pattern invalidation for: %s\n", pattern)
		return nil
	}

	fmt.Printf("🗑️  [CACHE] Invalidating cache pattern: %s\n", pattern)
	err := cm.cache.DeletePattern(pattern)
	if err != nil {
		fmt.Printf("⚠️  [CACHE] Pattern invalidation failed for %s: %v (continuing without cache)\n", pattern, err)
		return nil
	}

	fmt.Printf("✅ [CACHE] Pattern invalidation successful for: %s\n", pattern)
	return nil
}

// Set stores a value in cache
func (cm *CacheManager) Set(key string, value interface{}, ttl time.Duration) error {
	if !cm.enabled {
		return nil
	}
	return cm.cache.Set(key, value, ttl)
}

// Get retrieves a value from cache
func (cm *CacheManager) Get(key string, dest interface{}) error {
	if !cm.enabled {
		return fmt.Errorf("cache disabled")
	}
	return cm.cache.Get(key, dest)
}

// Exists checks if a key exists in cache
func (cm *CacheManager) Exists(key string) (bool, error) {
	if !cm.enabled {
		return false, nil
	}
	return cm.cache.Exists(key)
}

// IncrementVersion increments a version counter for cache invalidation
func (cm *CacheManager) IncrementVersion(key string) (int64, error) {
	if !cm.enabled {
		return 0, nil
	}
	return cm.cache.Increment(key)
}

// Close closes the cache connection
func (cm *CacheManager) Close() error {
	if cm.cache != nil {
		return cm.cache.Close()
	}
	return nil
}

// HealthCheck performs a health check on the cache
func (cm *CacheManager) HealthCheck() error {
	if !cm.enabled {
		return nil
	}
	return cm.cache.HealthCheck()
}

// InvalidateAllSeriesRelatedCaches invalidates all caches related to a series
func (cm *CacheManager) InvalidateAllSeriesRelatedCaches(seriesID string) {
	if !cm.enabled {
		return
	}

	fmt.Printf("🗑️  [CACHE] Invalidating all caches for series: %s\n", seriesID)

	// Invalidate series-specific caches
	seriesKey := cm.cache.GetSeriesKey(seriesID)
	_ = cm.Invalidate(seriesKey)

	// Invalidate matches by series cache
	matchesKey := cm.cache.GetMatchesBySeriesKey(seriesID)
	_ = cm.Invalidate(matchesKey)

	// Invalidate next match number cache
	nextNumberKey := fmt.Sprintf("match:next_number:series:%s", seriesID)
	_ = cm.Invalidate(nextNumberKey)

	// Invalidate all series list caches
	seriesListKeys := []string{
		"series:list",
		"series:list:order:created_at_desc",
		"series:list:order:created_at_desc:limit:20",
		"series:list:order:created_at_desc:limit:50",
		"series:list:order:created_at_desc:limit:100",
	}

	for _, key := range seriesListKeys {
		_ = cm.Invalidate(key)
	}

	// Invalidate all match-related caches that might contain matches from this series
	_ = cm.InvalidatePattern("match:list:*")
	_ = cm.InvalidatePattern("match:*")
	_ = cm.InvalidatePattern("scorecard:*")
	_ = cm.InvalidatePattern("scoreboard:*")

	// Invalidate match count
	_ = cm.Invalidate("match:count")

	fmt.Printf("✅ [CACHE] Completed comprehensive cache invalidation for series: %s\n", seriesID)
}

// InvalidateAllMatchRelatedCaches invalidates all caches related to a match
func (cm *CacheManager) InvalidateAllMatchRelatedCaches(matchID string, seriesID string) {
	if !cm.enabled {
		return
	}

	fmt.Printf("🗑️  [CACHE] Invalidating all caches for match: %s (series: %s)\n", matchID, seriesID)

	// Invalidate match-specific caches
	matchKey := cm.cache.GetMatchKey(matchID)
	_ = cm.Invalidate(matchKey)

	// Invalidate scorecard and scoreboard caches
	scorecardKey := cm.cache.GetScorecardKey(matchID)
	_ = cm.Invalidate(scorecardKey)
	scorecardVersionKey := cm.cache.GetScorecardVersionKey(matchID)
	_ = cm.Invalidate(scorecardVersionKey)

	// Invalidate match list caches
	_ = cm.Invalidate("match:list")
	_ = cm.Invalidate("match:count")

	// Invalidate common pagination cache keys
	paginationKeys := []string{
		"match:list:limit:20",
		"match:list:limit:10",
		"match:list:limit:5",
		"match:list:limit:3",
		"match:list:limit:2",
	}

	for _, key := range paginationKeys {
		_ = cm.Invalidate(key)
	}

	// Use pattern-based invalidation
	_ = cm.InvalidatePattern("match:list:*")
	_ = cm.InvalidatePattern("scorecard:*")
	_ = cm.InvalidatePattern("scoreboard:*")

	// If series ID is provided, invalidate series-related caches
	if seriesID != "" {
		seriesKey := cm.cache.GetMatchesBySeriesKey(seriesID)
		_ = cm.Invalidate(seriesKey)

		nextNumberKey := fmt.Sprintf("match:next_number:series:%s", seriesID)
		_ = cm.Invalidate(nextNumberKey)
	}

	fmt.Printf("✅ [CACHE] Completed comprehensive cache invalidation for match: %s\n", matchID)
}

// GetSeriesKey returns the cache key for a series
func (cm *CacheManager) GetSeriesKey(seriesID string) string {
	if cm.cache != nil {
		return cm.cache.GetSeriesKey(seriesID)
	}
	return fmt.Sprintf("series:%s", seriesID)
}

// GetMatchKey returns the cache key for a match
func (cm *CacheManager) GetMatchKey(matchID string) string {
	if cm.cache != nil {
		return cm.cache.GetMatchKey(matchID)
	}
	return fmt.Sprintf("match:%s", matchID)
}

// GetScorecardKey returns the cache key for a scorecard
func (cm *CacheManager) GetScorecardKey(matchID string) string {
	if cm.cache != nil {
		return cm.cache.GetScorecardKey(matchID)
	}
	return fmt.Sprintf("scorecard:%s", matchID)
}

// GetScorecardVersionKey returns the cache key for scorecard version
func (cm *CacheManager) GetScorecardVersionKey(matchID string) string {
	if cm.cache != nil {
		return cm.cache.GetScorecardVersionKey(matchID)
	}
	return fmt.Sprintf("scorecard:version:%s", matchID)
}

// GetMatchesBySeriesKey returns the cache key for matches by series
func (cm *CacheManager) GetMatchesBySeriesKey(seriesID string) string {
	if cm.cache != nil {
		return cm.cache.GetMatchesBySeriesKey(seriesID)
	}
	return fmt.Sprintf("matches:series:%s", seriesID)
}

// copyValue copies a value to destination interface
func copyValue(src, dest interface{}) error {
	// This is a simplified implementation
	// In a real scenario, you might want to use reflection or a more sophisticated approach
	srcBytes, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(srcBytes, dest)
}
