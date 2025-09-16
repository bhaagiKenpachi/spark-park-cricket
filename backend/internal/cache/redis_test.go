package cache

import (
	"spark-park-cricket-backend/internal/config"
	"testing"
	"time"
)

func TestRedisClient(t *testing.T) {
	// Skip test if Redis is not available
	cfg := &config.Config{
		RedisURL:      "localhost:6379",
		RedisPassword: "",
		RedisDB:       0,
		CacheEnabled:  true,
	}

	client, err := NewRedisClient(cfg)
	if err != nil {
		t.Skipf("Redis not available, skipping test: %v", err)
	}
	defer client.Close()

	// Test basic operations
	key := "test:key"
	value := "test value"
	ttl := 1 * time.Minute

	// Test Set
	err = client.Set(key, value, ttl)
	if err != nil {
		t.Fatalf("Failed to set value: %v", err)
	}

	// Test Get
	var retrievedValue string
	err = client.Get(key, &retrievedValue)
	if err != nil {
		t.Fatalf("Failed to get value: %v", err)
	}

	if retrievedValue != value {
		t.Errorf("Expected %s, got %s", value, retrievedValue)
	}

	// Test Exists
	exists, err := client.Exists(key)
	if err != nil {
		t.Fatalf("Failed to check existence: %v", err)
	}

	if !exists {
		t.Error("Key should exist")
	}

	// Test Delete
	err = client.Delete(key)
	if err != nil {
		t.Fatalf("Failed to delete key: %v", err)
	}

	// Test that key no longer exists
	exists, err = client.Exists(key)
	if err != nil {
		t.Fatalf("Failed to check existence after delete: %v", err)
	}

	if exists {
		t.Error("Key should not exist after deletion")
	}
}

func TestCacheManager(t *testing.T) {
	// Skip test if Redis is not available
	cfg := &config.Config{
		RedisURL:      "localhost:6379",
		RedisPassword: "",
		RedisDB:       0,
		CacheEnabled:  true,
	}

	redisClient, err := NewRedisClient(cfg)
	if err != nil {
		t.Skipf("Redis not available, skipping test: %v", err)
	}
	defer redisClient.Close()

	cacheManager := NewCacheManager(redisClient, true)

	// Test GetOrSet
	key := "test:getorset"
	expectedValue := "cached value"
	ttl := 1 * time.Minute

	var retrievedValue string
	err = cacheManager.GetOrSet(key, &retrievedValue, ttl, func() (interface{}, error) {
		return expectedValue, nil
	})

	if err != nil {
		t.Fatalf("Failed to get or set value: %v", err)
	}

	if retrievedValue != expectedValue {
		t.Errorf("Expected %s, got %s", expectedValue, retrievedValue)
	}

	// Test that second call returns cached value
	var cachedValue string
	err = cacheManager.GetOrSet(key, &cachedValue, ttl, func() (interface{}, error) {
		return "should not be called", nil
	})

	if err != nil {
		t.Fatalf("Failed to get cached value: %v", err)
	}

	if cachedValue != expectedValue {
		t.Errorf("Expected cached value %s, got %s", expectedValue, cachedValue)
	}

	// Test key generation methods
	seriesKey := cacheManager.GetSeriesKey("test-series-id")
	if seriesKey != "series:test-series-id" {
		t.Errorf("Expected series:test-series-id, got %s", seriesKey)
	}

	matchKey := cacheManager.GetMatchKey("test-match-id")
	if matchKey != "match:test-match-id" {
		t.Errorf("Expected match:test-match-id, got %s", matchKey)
	}

	scorecardKey := cacheManager.GetScorecardKey("test-match-id")
	if scorecardKey != "scorecard:test-match-id" {
		t.Errorf("Expected scorecard:test-match-id, got %s", scorecardKey)
	}
}
