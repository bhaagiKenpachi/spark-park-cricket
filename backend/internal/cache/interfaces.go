package cache

import (
	"encoding/json"
	"fmt"
	"time"
)

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
}

// CacheManager handles cache operations with fallback to database
type CacheManager struct {
	cache   CacheInterface
	enabled bool
}

// NewCacheManager creates a new cache manager
func NewCacheManager(cache CacheInterface, enabled bool) *CacheManager {
	return &CacheManager{
		cache:   cache,
		enabled: enabled,
	}
}

// GetOrSet retrieves from cache or sets from database function
func (cm *CacheManager) GetOrSet(key string, dest interface{}, ttl time.Duration, dbFunc func() (interface{}, error)) error {
	if !cm.enabled {
		// Cache disabled, call database function directly
		value, err := dbFunc()
		if err != nil {
			return err
		}

		// Copy value to destination
		return copyValue(value, dest)
	}

	// Try to get from cache first
	err := cm.cache.Get(key, dest)
	if err == nil {
		// Cache hit
		return nil
	}

	// Cache miss, get from database
	value, err := dbFunc()
	if err != nil {
		return err
	}

	// Set in cache for next time
	cm.cache.Set(key, value, ttl)

	// Copy value to destination
	return copyValue(value, dest)
}

// Invalidate removes a key from cache
func (cm *CacheManager) Invalidate(key string) error {
	if !cm.enabled {
		return nil
	}
	return cm.cache.Delete(key)
}

// InvalidatePattern removes all keys matching a pattern
func (cm *CacheManager) InvalidatePattern(pattern string) error {
	if !cm.enabled {
		return nil
	}
	// Note: This would require Redis SCAN command implementation
	// For now, we'll handle this at the service level
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
