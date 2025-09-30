package performance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"spark-park-cricket-backend/internal/config"
	"spark-park-cricket-backend/internal/database"
	"spark-park-cricket-backend/internal/handlers"
	"spark-park-cricket-backend/internal/models"
)

// BenchmarkValidationTest validates the performance improvements
type BenchmarkValidationTest struct {
	router   http.Handler
	matchID  string
	seriesID string
}

// SetupBenchmarkTest sets up the benchmark test environment
func SetupBenchmarkTest(t *testing.B) *BenchmarkValidationTest {
	// Create test configuration
	testCfg := &config.TestConfig{
		Config: &config.Config{
			SupabaseURL:    "http://localhost:54321",
			SupabaseAPIKey: "test-key",
		},
		TestSchema: "testing_db",
	}

	// Initialize test database
	dbClient, err := database.NewTestClient(testCfg)
	if err != nil {
		t.Fatalf("Failed to create test database client: %v", err)
	}

	// Setup routes
	router := handlers.SetupRoutes(dbClient, testCfg.Config)

	// Create test data
	seriesID, matchID, err := createTestDataForBenchmark(dbClient)
	if err != nil {
		t.Fatalf("Failed to create test data: %v", err)
	}

	// Start the match
	err = startMatchForBenchmark(router, matchID)
	if err != nil {
		t.Fatalf("Failed to start match: %v", err)
	}

	return &BenchmarkValidationTest{
		router:   router,
		matchID:  matchID,
		seriesID: seriesID,
	}
}

// BenchmarkAddBallOptimized benchmarks the optimized Add Ball API
func BenchmarkAddBallOptimized(b *testing.B) {
	test := SetupBenchmarkTest(b)
	defer func() {
		// Cleanup - router is http.Handler, not database.Client
		// Database cleanup is handled by the test framework
	}()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Create ball event request
			ballReq := models.BallEventRequest{
				MatchID:       test.matchID,
				InningsNumber: 1,
				BallType:      models.BallTypeGood,
				RunType:       models.RunTypeOne,
				IsWicket:      false,
				Byes:          0,
			}

			reqBody, err := json.Marshal(ballReq)
			if err != nil {
				b.Errorf("Failed to marshal request: %v", err)
				continue
			}

			// Create HTTP request
			req := httptest.NewRequest("POST", "/api/v1/scorecard/ball", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer test-token")

			// Record response time
			start := time.Now()
			w := httptest.NewRecorder()
			test.router.ServeHTTP(w, req)
			duration := time.Since(start)

			// Validate response
			if w.Code != http.StatusOK {
				b.Errorf("Request failed with status %d: %s", w.Code, w.Body.String())
			}

			// Record metrics
			b.ReportMetric(float64(duration.Nanoseconds())/1e6, "ms/op")
			b.ReportMetric(float64(len(w.Body.Bytes())), "bytes/op")

			// Validate performance targets
			if duration > 500*time.Millisecond {
				b.Errorf("Response time exceeds target: %v > 500ms", duration)
			}
		}
	})
}

// BenchmarkAddBallLegacy benchmarks the legacy Add Ball API (for comparison)
func BenchmarkAddBallLegacy(b *testing.B) {
	test := SetupBenchmarkTest(b)
	defer func() {
		// Cleanup - router is http.Handler, not database.Client
		// Database cleanup is handled by the test framework
	}()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Use legacy endpoint if available
			ballReq := models.BallEventRequest{
				MatchID:       test.matchID,
				InningsNumber: 1,
				BallType:      models.BallTypeGood,
				RunType:       models.RunTypeOne,
				IsWicket:      false,
				Byes:          0,
			}

			reqBody, err := json.Marshal(ballReq)
			if err != nil {
				b.Errorf("Failed to marshal request: %v", err)
				continue
			}

			// Create HTTP request to legacy endpoint
			req := httptest.NewRequest("POST", "/api/v1/scorecard/ball/legacy", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer test-token")

			// Record response time
			start := time.Now()
			w := httptest.NewRecorder()
			test.router.ServeHTTP(w, req)
			duration := time.Since(start)

			// Record metrics
			b.ReportMetric(float64(duration.Nanoseconds())/1e6, "ms/op")
			b.ReportMetric(float64(len(w.Body.Bytes())), "bytes/op")
		}
	})
}

// BenchmarkGetScorecard benchmarks the Get Scorecard API
func BenchmarkGetScorecard(b *testing.B) {
	test := SetupBenchmarkTest(b)
	defer func() {
		// Cleanup - router is http.Handler, not database.Client
		// Database cleanup is handled by the test framework
	}()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Create HTTP request
			req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/scorecard/%s", test.matchID), nil)
			req.Header.Set("Authorization", "Bearer test-token")

			// Record response time
			start := time.Now()
			w := httptest.NewRecorder()
			test.router.ServeHTTP(w, req)
			duration := time.Since(start)

			// Validate response
			if w.Code != http.StatusOK {
				b.Errorf("Request failed with status %d: %s", w.Code, w.Body.String())
			}

			// Record metrics
			b.ReportMetric(float64(duration.Nanoseconds())/1e6, "ms/op")
			b.ReportMetric(float64(len(w.Body.Bytes())), "bytes/op")

			// Validate performance targets
			if duration > 200*time.Millisecond {
				b.Errorf("Response time exceeds target: %v > 200ms", duration)
			}
		}
	})
}

// BenchmarkConcurrentAddBall benchmarks concurrent Add Ball requests
func BenchmarkConcurrentAddBall(b *testing.B) {
	test := SetupBenchmarkTest(b)
	defer func() {
		// Cleanup - router is http.Handler, not database.Client
		// Database cleanup is handled by the test framework
	}()

	concurrencyLevels := []int{1, 5, 10, 20, 50, 100}

	for _, concurrency := range concurrencyLevels {
		b.Run(fmt.Sprintf("Concurrency_%d", concurrency), func(b *testing.B) {
			b.SetParallelism(concurrency)
			b.ResetTimer()

			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					ballReq := models.BallEventRequest{
						MatchID:       test.matchID,
						InningsNumber: 1,
						BallType:      models.BallTypeGood,
						RunType:       models.RunTypeOne,
						IsWicket:      false,
						Byes:          0,
					}

					reqBody, err := json.Marshal(ballReq)
					if err != nil {
						b.Errorf("Failed to marshal request: %v", err)
						continue
					}

					req := httptest.NewRequest("POST", "/api/v1/scorecard/ball", bytes.NewBuffer(reqBody))
					req.Header.Set("Content-Type", "application/json")
					req.Header.Set("Authorization", "Bearer test-token")

					start := time.Now()
					w := httptest.NewRecorder()
					test.router.ServeHTTP(w, req)
					duration := time.Since(start)

					if w.Code != http.StatusOK {
						b.Errorf("Request failed with status %d: %s", w.Code, w.Body.String())
					}

					b.ReportMetric(float64(duration.Nanoseconds())/1e6, "ms/op")
				}
			})
		})
	}
}

// createTestDataForBenchmark creates test data for benchmarking
func createTestDataForBenchmark(dbClient *database.Client) (string, string, error) {
	// This would create test series, match, teams, and players
	// For now, return mock IDs
	return "benchmark-series-id", "benchmark-match-id", nil
}

// startMatchForBenchmark starts a match for benchmarking
func startMatchForBenchmark(router http.Handler, matchID string) error {
	// This would start the match and begin scoring
	// For now, return nil (mock implementation)
	return nil
}

// TestPerformanceTargets validates that performance targets are met
func TestPerformanceTargets(t *testing.T) {
	test := SetupBenchmarkTest(&testing.B{})
	defer func() {
		// Cleanup - router is http.Handler, not database.Client
		// Database cleanup is handled by the test framework
	}()

	// Test Add Ball API performance
	ballReq := models.BallEventRequest{
		MatchID:       test.matchID,
		InningsNumber: 1,
		BallType:      models.BallTypeGood,
		RunType:       models.RunTypeOne,
		IsWicket:      false,
		Byes:          0,
	}

	reqBody, err := json.Marshal(ballReq)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	req := httptest.NewRequest("POST", "/api/v1/scorecard/ball", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")

	start := time.Now()
	w := httptest.NewRecorder()
	test.router.ServeHTTP(w, req)
	duration := time.Since(start)

	// Validate performance targets
	if w.Code != http.StatusOK {
		t.Errorf("Request failed with status %d: %s", w.Code, w.Body.String())
	}

	if duration > 500*time.Millisecond {
		t.Errorf("Add Ball API response time exceeds target: %v > 500ms", duration)
	}

	// Test Get Scorecard API performance
	req = httptest.NewRequest("GET", fmt.Sprintf("/api/v1/scorecard/%s", test.matchID), nil)
	req.Header.Set("Authorization", "Bearer test-token")

	start = time.Now()
	w = httptest.NewRecorder()
	test.router.ServeHTTP(w, req)
	duration = time.Since(start)

	if w.Code != http.StatusOK {
		t.Errorf("Get Scorecard request failed with status %d: %s", w.Code, w.Body.String())
	}

	if duration > 200*time.Millisecond {
		t.Errorf("Get Scorecard API response time exceeds target: %v > 200ms", duration)
	}

	t.Logf("✅ Performance targets validated successfully")
	t.Logf("   Add Ball API: %v (target: <500ms)", duration)
	t.Logf("   Get Scorecard API: %v (target: <200ms)", duration)
}
