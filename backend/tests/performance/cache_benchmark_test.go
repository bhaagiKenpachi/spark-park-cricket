package performance

import (
	"fmt"
	"spark-park-cricket-backend/internal/cache"
	"spark-park-cricket-backend/internal/config"
	"testing"
	"time"
)

// BenchmarkCacheOperations benchmarks basic cache operations
func BenchmarkCacheOperations(b *testing.B) {
	cfg := &config.Config{
		RedisURL:      "localhost:6379",
		RedisPassword: "",
		RedisDB:       0,
		CacheEnabled:  true,
	}

	redisClient, err := cache.NewRedisClient(cfg)
	if err != nil {
		b.Skipf("Redis not available, skipping benchmark: %v", err)
	}
	defer redisClient.Close()

	cacheManager := cache.NewCacheManager(redisClient, true)

	b.Run("CacheSet", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("benchmark:set:%d", i)
			value := fmt.Sprintf("value:%d", i)
			err := cacheManager.Set(key, value, cache.StaticDataTTL)
			if err != nil {
				b.Fatalf("Failed to set cache value: %v", err)
			}
		}
	})

	b.Run("CacheGet", func(b *testing.B) {
		// Pre-populate cache
		for i := 0; i < 1000; i++ {
			key := fmt.Sprintf("benchmark:get:%d", i)
			value := fmt.Sprintf("value:%d", i)
			cacheManager.Set(key, value, cache.StaticDataTTL)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("benchmark:get:%d", i%1000)
			var value string
			err := cacheManager.Get(key, &value)
			if err != nil {
				b.Fatalf("Failed to get cache value: %v", err)
			}
		}
	})

	b.Run("CacheGetOrSet", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("benchmark:getorset:%d", i)
			var value string
			err := cacheManager.GetOrSet(key, &value, cache.StaticDataTTL, func() (interface{}, error) {
				return fmt.Sprintf("computed:value:%d", i), nil
			})
			if err != nil {
				b.Fatalf("Failed to get or set cache value: %v", err)
			}
		}
	})
}

// BenchmarkCacheVsDatabase simulates the performance difference between cache and database
func BenchmarkCacheVsDatabase(b *testing.B) {
	cfg := &config.Config{
		RedisURL:      "localhost:6379",
		RedisPassword: "",
		RedisDB:       0,
		CacheEnabled:  true,
	}

	redisClient, err := cache.NewRedisClient(cfg)
	if err != nil {
		b.Skipf("Redis not available, skipping benchmark: %v", err)
	}
	defer redisClient.Close()

	cacheManager := cache.NewCacheManager(redisClient, true)

	// Simulate database operation (slower)
	simulateDatabaseOperation := func() (interface{}, error) {
		time.Sleep(10 * time.Millisecond) // Simulate 10ms database query
		return "database_result", nil
	}

	b.Run("DatabaseOnly", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := simulateDatabaseOperation()
			if err != nil {
				b.Fatalf("Database operation failed: %v", err)
			}
		}
	})

	b.Run("CacheWithDatabaseFallback", func(b *testing.B) {
		// Pre-populate cache
		cacheManager.Set("benchmark:cache:test", "cached_result", cache.StaticDataTTL)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var result string
			err := cacheManager.GetOrSet("benchmark:cache:test", &result, cache.StaticDataTTL, simulateDatabaseOperation)
			if err != nil {
				b.Fatalf("Cache operation failed: %v", err)
			}
		}
	})
}

// BenchmarkScorecardCache simulates scorecard caching performance
func BenchmarkScorecardCache(b *testing.B) {
	cfg := &config.Config{
		RedisURL:      "localhost:6379",
		RedisPassword: "",
		RedisDB:       0,
		CacheEnabled:  true,
	}

	redisClient, err := cache.NewRedisClient(cfg)
	if err != nil {
		b.Skipf("Redis not available, skipping benchmark: %v", err)
	}
	defer redisClient.Close()

	cacheManager := cache.NewCacheManager(redisClient, true)

	// Simulate large scorecard data (similar to the 5,749 bytes from benchmark)
	scorecardData := map[string]interface{}{
		"match_id": "test-match-123",
		"innings": []map[string]interface{}{
			{
				"innings_number": 1,
				"total_runs":     150,
				"total_wickets":  3,
				"total_overs":    20.0,
				"overs": []map[string]interface{}{
					{
						"over_number": 1,
						"balls": []map[string]interface{}{
							{"ball_number": 1, "runs": 1},
							{"ball_number": 2, "runs": 0},
							{"ball_number": 3, "runs": 4},
							{"ball_number": 4, "runs": 2},
							{"ball_number": 5, "runs": 1},
							{"ball_number": 6, "runs": 0},
						},
					},
				},
			},
		},
	}

	b.Run("ScorecardCacheSet", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			key := cacheManager.GetScorecardKey(fmt.Sprintf("match-%d", i))
			err := cacheManager.Set(key, scorecardData, cache.ScorecardTTL)
			if err != nil {
				b.Fatalf("Failed to cache scorecard: %v", err)
			}
		}
	})

	b.Run("ScorecardCacheGet", func(b *testing.B) {
		// Pre-populate cache
		for i := 0; i < 100; i++ {
			key := cacheManager.GetScorecardKey(fmt.Sprintf("match-%d", i))
			cacheManager.Set(key, scorecardData, cache.ScorecardTTL)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			key := cacheManager.GetScorecardKey(fmt.Sprintf("match-%d", i%100))
			var result map[string]interface{}
			err := cacheManager.Get(key, &result)
			if err != nil {
				b.Fatalf("Failed to get cached scorecard: %v", err)
			}
		}
	})
}
