package performance

import (
	"bytes"
	"context"
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
	"spark-park-cricket-backend/internal/services"
	"spark-park-cricket-backend/pkg/testutils"
)

// BenchmarkValidationTest validates the performance improvements
type BenchmarkValidationTest struct {
	router           http.Handler
	matchID          string
	seriesID         string
	authCookie       string
	serviceContainer *services.Container
}

// SetupBenchmarkTest sets up the benchmark test environment
func SetupBenchmarkTest(t *testing.B) *BenchmarkValidationTest {
	// Load test configuration
	testCfg := config.LoadTestConfig()

	// Initialize test database
	dbClient, err := database.NewTestClient(testCfg)
	if err != nil {
		t.Fatalf("Failed to create test database client: %v", err)
	}

	// Setup test schema
	err = database.SetupTestSchema(testCfg)
	if err != nil {
		t.Fatalf("Failed to setup test schema: %v", err)
	}

	// Create service container
	serviceContainer := services.NewContainer(dbClient, testCfg.Config)

	// Create handlers
	seriesHandler := handlers.NewSeriesHandler(serviceContainer.Series)
	matchHandler := handlers.NewMatchHandler(serviceContainer.Match)
	scorecardHandler := handlers.NewScorecardHandler(serviceContainer.Scorecard)

	// Setup routes with proper authentication
	router := SetupTestRoutes(serviceContainer, seriesHandler, matchHandler, scorecardHandler, testCfg.Config)

	// Create authenticated test user
	mockT := &testing.T{}
	user, authCookie := testutils.CreateAuthenticatedTestUserWithSessionService(mockT, dbClient, serviceContainer.SessionService)

	// Create test data
	seriesID, matchID, err := createTestDataForBenchmark(dbClient, user.ID)
	if err != nil {
		t.Fatalf("Failed to create test data: %v", err)
	}

	// Start the match
	err = startMatchForBenchmark(router, matchID, authCookie)
	if err != nil {
		t.Fatalf("Failed to start match: %v", err)
	}

	return &BenchmarkValidationTest{
		router:           router,
		matchID:          matchID,
		seriesID:         seriesID,
		authCookie:       authCookie,
		serviceContainer: serviceContainer,
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
			req.Header.Set("Cookie", test.authCookie)

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
			req.Header.Set("Cookie", test.authCookie)

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
			req.Header.Set("Cookie", test.authCookie)

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
					req.Header.Set("Cookie", test.authCookie)

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
func createTestDataForBenchmark(dbClient *database.Client, userID string) (string, string, error) {
	ctx := context.Background()

	// Create test series
	series := &models.Series{
		Name:      fmt.Sprintf("Benchmark Test Series %d", time.Now().Unix()),
		StartDate: time.Now(),
		EndDate:   time.Now().Add(24 * time.Hour),
		CreatedBy: userID,
	}
	err := dbClient.Repositories.Series.Create(ctx, series)
	if err != nil {
		return "", "", fmt.Errorf("failed to create test series: %v", err)
	}

	// Create test match
	match := &models.Match{
		SeriesID:         series.ID,
		MatchNumber:      1,
		Date:             time.Now(),
		Status:           models.MatchStatusLive,
		TeamAPlayerCount: 11,
		TeamBPlayerCount: 11,
		TotalOvers:       20,
		TossWinner:       models.TeamTypeA,
		TossType:         models.TossTypeHeads,
		BattingTeam:      models.TeamTypeA,
		CreatedBy:        userID,
	}
	err = dbClient.Repositories.Match.Create(ctx, match)
	if err != nil {
		return "", "", fmt.Errorf("failed to create test match: %v", err)
	}

	return series.ID, match.ID, nil
}

// startMatchForBenchmark starts a match for benchmarking
func startMatchForBenchmark(router http.Handler, matchID string, authCookie string) error {
	// Start scoring for the match
	startReq := models.ScorecardRequest{
		MatchID: matchID,
	}

	reqBody, err := json.Marshal(startReq)
	if err != nil {
		return fmt.Errorf("failed to marshal start request: %v", err)
	}

	req := httptest.NewRequest("POST", "/api/v1/scorecard/start", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:     "user_session",
		Value:    authCookie,
		Path:     "/",
		HttpOnly: true,
	})
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		return fmt.Errorf("failed to start match: status %d, body: %s", w.Code, w.Body.String())
	}

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
	req.AddCookie(&http.Cookie{
		Name:     "user_session",
		Value:    test.authCookie,
		Path:     "/",
		HttpOnly: true,
	})

	start := time.Now()
	w := httptest.NewRecorder()
	test.router.ServeHTTP(w, req)
	duration := time.Since(start)

	// Validate performance targets
	if w.Code != http.StatusOK {
		t.Errorf("Request failed with status %d: %s", w.Code, w.Body.String())
	}

	if duration > 1500*time.Millisecond {
		t.Errorf("Add Ball API response time exceeds target: %v > 1500ms", duration)
	}

	// Test Get Scorecard API performance
	req = httptest.NewRequest("GET", fmt.Sprintf("/api/v1/scorecard/%s", test.matchID), nil)
	req.AddCookie(&http.Cookie{
		Name:     "user_session",
		Value:    test.authCookie,
		Path:     "/",
		HttpOnly: true,
	})

	start = time.Now()
	w = httptest.NewRecorder()
	test.router.ServeHTTP(w, req)
	duration = time.Since(start)

	if w.Code != http.StatusOK {
		t.Errorf("Get Scorecard request failed with status %d: %s", w.Code, w.Body.String())
	}

	if duration > 1000*time.Millisecond {
		t.Errorf("Get Scorecard API response time exceeds target: %v > 1000ms", duration)
	}

	t.Logf("✅ Performance targets validated successfully")
	t.Logf("   Add Ball API: %v (target: <1500ms)", duration)
	t.Logf("   Get Scorecard API: %v (target: <1000ms)", duration)
}
