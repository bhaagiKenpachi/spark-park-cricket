package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// Cache represents an in-memory cache
type Cache struct {
	items map[string]*CacheItem
	mutex sync.RWMutex
}

// CacheItem represents a cached item
type CacheItem struct {
	Value     interface{}
	ExpiresAt time.Time
}

// NewCache creates a new cache instance
func NewCache() *Cache {
	cache := &Cache{
		items: make(map[string]*CacheItem),
	}

	// Start cleanup goroutine
	go cache.cleanup()

	return cache
}

// Set stores a value in the cache with an expiration time
func (c *Cache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.items[key] = &CacheItem{
		Value:     value,
		ExpiresAt: time.Now().Add(expiration),
	}

	return nil
}

// Get retrieves a value from the cache
func (c *Cache) Get(ctx context.Context, key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return nil, false
	}

	// Check if item has expired
	if time.Now().After(item.ExpiresAt) {
		return nil, false
	}

	return item.Value, true
}

// Delete removes a value from the cache
func (c *Cache) Delete(ctx context.Context, key string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.items, key)
	return nil
}

// Clear removes all items from the cache
func (c *Cache) Clear(ctx context.Context) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.items = make(map[string]*CacheItem)
	return nil
}

// GetOrSet retrieves a value from the cache, or sets it if it doesn't exist
func (c *Cache) GetOrSet(ctx context.Context, key string, setter func() (interface{}, error), expiration time.Duration) (interface{}, error) {
	// Try to get from cache first
	if value, exists := c.Get(ctx, key); exists {
		return value, nil
	}

	// Set the value using the setter function
	value, err := setter()
	if err != nil {
		return nil, err
	}

	// Store in cache
	err = c.Set(ctx, key, value, expiration)
	if err != nil {
		return nil, err
	}

	return value, nil
}

// cleanup removes expired items from the cache
func (c *Cache) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mutex.Lock()
		now := time.Now()
		for key, item := range c.items {
			if now.After(item.ExpiresAt) {
				delete(c.items, key)
			}
		}
		c.mutex.Unlock()
	}
}

// GetStats returns cache statistics
func (c *Cache) GetStats() map[string]interface{} {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	totalItems := len(c.items)
	expiredItems := 0
	now := time.Now()

	for _, item := range c.items {
		if now.After(item.ExpiresAt) {
			expiredItems++
		}
	}

	return map[string]interface{}{
		"total_items":   totalItems,
		"expired_items": expiredItems,
		"active_items":  totalItems - expiredItems,
	}
}

// CacheManager manages multiple caches
type CacheManager struct {
	caches map[string]*Cache
	mutex  sync.RWMutex
}

// NewCacheManager creates a new cache manager
func NewCacheManager() *CacheManager {
	return &CacheManager{
		caches: make(map[string]*Cache),
	}
}

// GetCache returns a cache by name, creating it if it doesn't exist
func (cm *CacheManager) GetCache(name string) *Cache {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if cache, exists := cm.caches[name]; exists {
		return cache
	}

	cache := NewCache()
	cm.caches[name] = cache
	return cache
}

// GetStats returns statistics for all caches
func (cm *CacheManager) GetStats() map[string]interface{} {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	stats := make(map[string]interface{})
	for name, cache := range cm.caches {
		stats[name] = cache.GetStats()
	}

	return stats
}

// CacheKeyBuilder helps build cache keys
type CacheKeyBuilder struct {
	parts []string
}

// NewCacheKeyBuilder creates a new cache key builder
func NewCacheKeyBuilder() *CacheKeyBuilder {
	return &CacheKeyBuilder{
		parts: make([]string, 0),
	}
}

// Add adds a part to the cache key
func (ckb *CacheKeyBuilder) Add(part string) *CacheKeyBuilder {
	ckb.parts = append(ckb.parts, part)
	return ckb
}

// AddInt adds an integer part to the cache key
func (ckb *CacheKeyBuilder) AddInt(part int) *CacheKeyBuilder {
	ckb.parts = append(ckb.parts, fmt.Sprintf("%d", part))
	return ckb
}

// AddJSON adds a JSON representation of an object to the cache key
func (ckb *CacheKeyBuilder) AddJSON(obj interface{}) *CacheKeyBuilder {
	jsonBytes, _ := json.Marshal(obj)
	ckb.parts = append(ckb.parts, string(jsonBytes))
	return ckb
}

// Build builds the final cache key
func (ckb *CacheKeyBuilder) Build() string {
	key := ""
	for i, part := range ckb.parts {
		if i > 0 {
			key += ":"
		}
		key += part
	}
	return key
}

// Cache constants
const (
	DefaultExpiration = 5 * time.Minute
	LongExpiration    = 1 * time.Hour
	ShortExpiration   = 1 * time.Minute
)

// Cache names
const (
	SeriesCache     = "series"
	MatchCache      = "match"
	TeamCache       = "team"
	PlayerCache     = "player"
	ScoreboardCache = "scoreboard"
)
