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
	"spark-park-cricket-backend/internal/services"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

// Comprehensive database benchmark tests for all tables and operations
// Tests CRUD operations, queries, and performance across all entities

func BenchmarkSeriesOperations(b *testing.B) {
	cfg := config.LoadTestConfig()
	testDB, err := database.NewTestClient(cfg)
	if err != nil {
		b.Fatalf("Failed to create test database client: %v", err)
	}
	defer testDB.Close()

	err = database.SetupTestSchema(cfg)
	if err != nil {
		b.Fatalf("Failed to setup test schema: %v", err)
	}

	serviceContainer := services.NewContainer(testDB.Repositories)
	router := setupBenchmarkRouter(serviceContainer)

	b.Run("CreateSeries", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			seriesReq := models.CreateSeriesRequest{
				Name:      fmt.Sprintf("Benchmark Series %d", i),
				StartDate: time.Now(),
				EndDate:   time.Now().Add(24 * time.Hour),
			}

			reqBody, _ := json.Marshal(seriesReq)
			req := httptest.NewRequest("POST", "/api/v1/series", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			start := time.Now()
			router.ServeHTTP(w, req)
			duration := time.Since(start)

			if w.Code != http.StatusCreated {
				b.Errorf("Series creation failed: %d", w.Code)
			}

			b.ReportMetric(float64(duration.Nanoseconds())/1e6, "ms/op")
			b.ReportMetric(float64(len(w.Body.Bytes())), "bytes/op")
		}
	})

	b.Run("GetSeries", func(b *testing.B) {
		// Create test series first
		seriesID := createTestSeries(b, router)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			req := httptest.NewRequest("GET", "/api/v1/series/"+seriesID, nil)
			w := httptest.NewRecorder()

			start := time.Now()
			router.ServeHTTP(w, req)
			duration := time.Since(start)

			if w.Code != http.StatusOK {
				b.Errorf("Series retrieval failed: %d", w.Code)
			}

			b.ReportMetric(float64(duration.Nanoseconds())/1e6, "ms/op")
			b.ReportMetric(float64(len(w.Body.Bytes())), "bytes/op")
		}
	})
}

func BenchmarkMatchOperations(b *testing.B) {
	cfg := config.LoadTestConfig()
	testDB, err := database.NewTestClient(cfg)
	if err != nil {
		b.Fatalf("Failed to create test database client: %v", err)
	}
	defer testDB.Close()

	err = database.SetupTestSchema(cfg)
	if err != nil {
		b.Fatalf("Failed to setup test schema: %v", err)
	}

	serviceContainer := services.NewContainer(testDB.Repositories)
	router := setupBenchmarkRouter(serviceContainer)

	b.Run("CreateMatch", func(b *testing.B) {
		seriesID := createTestSeries(b, router)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			matchNumber := i + 1
			matchReq := models.CreateMatchRequest{
				SeriesID:         seriesID,
				MatchNumber:      &matchNumber,
				Date:             time.Now(),
				TeamAPlayerCount: 11,
				TeamBPlayerCount: 11,
				TotalOvers:       20,
				TossWinner:       models.TeamTypeA,
				TossType:         models.TossTypeHeads,
			}

			reqBody, _ := json.Marshal(matchReq)
			req := httptest.NewRequest("POST", "/api/v1/matches", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			start := time.Now()
			router.ServeHTTP(w, req)
			duration := time.Since(start)

			if w.Code != http.StatusCreated {
				b.Errorf("Match creation failed: %d", w.Code)
			}

			b.ReportMetric(float64(duration.Nanoseconds())/1e6, "ms/op")
			b.ReportMetric(float64(len(w.Body.Bytes())), "bytes/op")
		}
	})

	b.Run("UpdateMatch", func(b *testing.B) {
		seriesID := createTestSeries(b, router)
		matchID := createTestMatch(b, router, seriesID)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			status := models.MatchStatusLive
			updateReq := models.UpdateMatchRequest{
				Status: &status,
			}

			reqBody, _ := json.Marshal(updateReq)
			req := httptest.NewRequest("PUT", "/api/v1/matches/"+matchID, bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			start := time.Now()
			router.ServeHTTP(w, req)
			duration := time.Since(start)

			if w.Code != http.StatusOK {
				b.Errorf("Match update failed: %d", w.Code)
			}

			b.ReportMetric(float64(duration.Nanoseconds())/1e6, "ms/op")
			b.ReportMetric(float64(len(w.Body.Bytes())), "bytes/op")
		}
	})
}

func BenchmarkScorecardOperations(b *testing.B) {
	cfg := config.LoadTestConfig()
	testDB, err := database.NewTestClient(cfg)
	if err != nil {
		b.Fatalf("Failed to create test database client: %v", err)
	}
	defer testDB.Close()

	err = database.SetupTestSchema(cfg)
	if err != nil {
		b.Fatalf("Failed to setup test schema: %v", err)
	}

	serviceContainer := services.NewContainer(testDB.Repositories)
	router := setupBenchmarkRouter(serviceContainer)

	b.Run("StartScoring", func(b *testing.B) {
		seriesID := createTestSeries(b, router)
		matchID := createTestMatch(b, router, seriesID)
		updateMatchToLive(b, router, matchID)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			startReq := models.ScorecardRequest{
				MatchID: matchID,
			}

			reqBody, _ := json.Marshal(startReq)
			req := httptest.NewRequest("POST", "/api/v1/scorecard/start", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			start := time.Now()
			router.ServeHTTP(w, req)
			duration := time.Since(start)

			if w.Code != http.StatusOK {
				b.Errorf("Start scoring failed: %d", w.Code)
			}

			b.ReportMetric(float64(duration.Nanoseconds())/1e6, "ms/op")
			b.ReportMetric(float64(len(w.Body.Bytes())), "bytes/op")
		}
	})

	b.Run("AddBall", func(b *testing.B) {
		seriesID := createTestSeries(b, router)
		matchID := createTestMatch(b, router, seriesID)
		updateMatchToLive(b, router, matchID)
		startScoring(b, router, matchID)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			ballReq := models.BallEventRequest{
				MatchID:       matchID,
				InningsNumber: 1,
				BallType:      models.BallTypeGood,
				RunType:       models.RunTypeOne,
				IsWicket:      false,
			}

			reqBody, _ := json.Marshal(ballReq)
			req := httptest.NewRequest("POST", "/api/v1/scorecard/ball", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			start := time.Now()
			router.ServeHTTP(w, req)
			duration := time.Since(start)

			if w.Code != http.StatusOK {
				b.Errorf("Add ball failed: %d", w.Code)
			}

			b.ReportMetric(float64(duration.Nanoseconds())/1e6, "ms/op")
			b.ReportMetric(float64(len(w.Body.Bytes())), "bytes/op")
		}
	})

	b.Run("GetScorecard", func(b *testing.B) {
		seriesID := createTestSeries(b, router)
		matchID := createTestMatch(b, router, seriesID)
		updateMatchToLive(b, router, matchID)
		startScoring(b, router, matchID)

		// Add some balls to make scorecard more realistic
		for i := 0; i < 20; i++ {
			addTestBall(b, router, matchID)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			req := httptest.NewRequest("GET", "/api/v1/scorecard/"+matchID, nil)
			w := httptest.NewRecorder()

			start := time.Now()
			router.ServeHTTP(w, req)
			duration := time.Since(start)

			if w.Code != http.StatusOK {
				b.Errorf("Get scorecard failed: %d", w.Code)
			}

			b.ReportMetric(float64(duration.Nanoseconds())/1e6, "ms/op")
			b.ReportMetric(float64(len(w.Body.Bytes())), "bytes/op")
		}
	})

	b.Run("UndoBall", func(b *testing.B) {
		seriesID := createTestSeries(b, router)
		matchID := createTestMatch(b, router, seriesID)
		updateMatchToLive(b, router, matchID)
		startScoring(b, router, matchID)

		// Add a ball first
		addTestBall(b, router, matchID)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			req := httptest.NewRequest("DELETE", "/api/v1/scorecard/"+matchID+"/ball", nil)
			w := httptest.NewRecorder()

			start := time.Now()
			router.ServeHTTP(w, req)
			duration := time.Since(start)

			if w.Code != http.StatusOK {
				b.Errorf("Undo ball failed: %d", w.Code)
			}

			b.ReportMetric(float64(duration.Nanoseconds())/1e6, "ms/op")
			b.ReportMetric(float64(len(w.Body.Bytes())), "bytes/op")
		}
	})
}

func BenchmarkDatabaseQueries(b *testing.B) {
	cfg := config.LoadTestConfig()
	testDB, err := database.NewTestClient(cfg)
	if err != nil {
		b.Fatalf("Failed to create test database client: %v", err)
	}
	defer testDB.Close()

	err = database.SetupTestSchema(cfg)
	if err != nil {
		b.Fatalf("Failed to setup test schema: %v", err)
	}

	serviceContainer := services.NewContainer(testDB.Repositories)
	router := setupBenchmarkRouter(serviceContainer)

	// Create test data
	seriesID := createTestSeries(b, router)
	matchID := createTestMatch(b, router, seriesID)
	updateMatchToLive(b, router, matchID)
	startScoring(b, router, matchID)

	// Add multiple balls for realistic data
	for i := 0; i < 50; i++ {
		addTestBall(b, router, matchID)
	}

	b.Run("ComplexScorecardQuery", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			req := httptest.NewRequest("GET", "/api/v1/scorecard/"+matchID, nil)
			w := httptest.NewRecorder()

			start := time.Now()
			router.ServeHTTP(w, req)
			duration := time.Since(start)

			if w.Code != http.StatusOK {
				b.Errorf("Complex scorecard query failed: %d", w.Code)
			}

			b.ReportMetric(float64(duration.Nanoseconds())/1e6, "ms/op")
			b.ReportMetric(float64(len(w.Body.Bytes())), "bytes/op")
		}
	})

	b.Run("SeriesWithMatchesQuery", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			req := httptest.NewRequest("GET", "/api/v1/series/"+seriesID, nil)
			w := httptest.NewRecorder()

			start := time.Now()
			router.ServeHTTP(w, req)
			duration := time.Since(start)

			if w.Code != http.StatusOK {
				b.Errorf("Series with matches query failed: %d", w.Code)
			}

			b.ReportMetric(float64(duration.Nanoseconds())/1e6, "ms/op")
			b.ReportMetric(float64(len(w.Body.Bytes())), "bytes/op")
		}
	})
}

func BenchmarkConcurrentOperations(b *testing.B) {
	cfg := config.LoadTestConfig()
	testDB, err := database.NewTestClient(cfg)
	if err != nil {
		b.Fatalf("Failed to create test database client: %v", err)
	}
	defer testDB.Close()

	err = database.SetupTestSchema(cfg)
	if err != nil {
		b.Fatalf("Failed to setup test schema: %v", err)
	}

	serviceContainer := services.NewContainer(testDB.Repositories)
	router := setupBenchmarkRouter(serviceContainer)

	// Create test data
	seriesID := createTestSeries(b, router)
	matchID := createTestMatch(b, router, seriesID)
	updateMatchToLive(b, router, matchID)
	startScoring(b, router, matchID)

	b.Run("ConcurrentBallAdditions", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				ballReq := models.BallEventRequest{
					MatchID:       matchID,
					InningsNumber: 1,
					BallType:      models.BallTypeGood,
					RunType:       models.RunTypeOne,
					IsWicket:      false,
				}

				reqBody, _ := json.Marshal(ballReq)
				req := httptest.NewRequest("POST", "/api/v1/scorecard/ball", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()

				start := time.Now()
				router.ServeHTTP(w, req)
				duration := time.Since(start)

				if w.Code != http.StatusOK {
					b.Errorf("Concurrent ball addition failed: %d", w.Code)
				}

				b.ReportMetric(float64(duration.Nanoseconds())/1e6, "ms/op")
			}
		})
	})

	b.Run("ConcurrentScorecardReads", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				req := httptest.NewRequest("GET", "/api/v1/scorecard/"+matchID, nil)
				w := httptest.NewRecorder()

				start := time.Now()
				router.ServeHTTP(w, req)
				duration := time.Since(start)

				if w.Code != http.StatusOK {
					b.Errorf("Concurrent scorecard read failed: %d", w.Code)
				}

				b.ReportMetric(float64(duration.Nanoseconds())/1e6, "ms/op")
			}
		})
	})
}

// Helper functions for benchmark tests
func setupBenchmarkRouter(serviceContainer *services.Container) http.Handler {
	router := chi.NewRouter()

	// Add middleware
	router.Use(chimiddleware.Recoverer)
	router.Use(chimiddleware.Logger)
	router.Use(chimiddleware.RequestID)
	router.Use(chimiddleware.RealIP)
	router.Use(chimiddleware.Timeout(60 * time.Second))

	// API routes
	router.Route("/api/v1", func(r chi.Router) {
		// Series routes
		r.Route("/series", func(r chi.Router) {
			seriesHandler := handlers.NewSeriesHandler(serviceContainer.Series)
			r.Post("/", seriesHandler.CreateSeries)
			r.Get("/{id}", seriesHandler.GetSeries)
		})
		// Match routes
		r.Route("/matches", func(r chi.Router) {
			matchHandler := handlers.NewMatchHandler(serviceContainer.Match)
			r.Post("/", matchHandler.CreateMatch)
			r.Put("/{id}", matchHandler.UpdateMatch)
		})
		// Scorecard routes
		r.Route("/scorecard", func(r chi.Router) {
			scorecardHandler := handlers.NewScorecardHandler(serviceContainer.Scorecard)
			r.Post("/start", scorecardHandler.StartScoring)
			r.Post("/ball", scorecardHandler.AddBall)
			r.Delete("/{match_id}/ball", scorecardHandler.UndoBall)
			r.Get("/{match_id}", scorecardHandler.GetScorecard)
		})
	})

	return router
}

func createTestSeries(b *testing.B, router http.Handler) string {
	seriesReq := models.CreateSeriesRequest{
		Name:      "Benchmark Test Series",
		StartDate: time.Now(),
		EndDate:   time.Now().Add(24 * time.Hour),
	}

	reqBody, _ := json.Marshal(seriesReq)
	req := httptest.NewRequest("POST", "/api/v1/series", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		b.Fatalf("Failed to create series: %d", w.Code)
	}

	var response struct {
		Data models.Series `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &response)
	return response.Data.ID
}

func createTestMatch(b *testing.B, router http.Handler, seriesID string) string {
	matchNumber := 1
	matchReq := models.CreateMatchRequest{
		SeriesID:         seriesID,
		MatchNumber:      &matchNumber,
		Date:             time.Now(),
		TeamAPlayerCount: 11,
		TeamBPlayerCount: 11,
		TotalOvers:       20,
		TossWinner:       models.TeamTypeA,
		TossType:         models.TossTypeHeads,
	}

	reqBody, _ := json.Marshal(matchReq)
	req := httptest.NewRequest("POST", "/api/v1/matches", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		b.Fatalf("Failed to create match: %d", w.Code)
	}

	var response struct {
		Data models.Match `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &response)
	return response.Data.ID
}

func updateMatchToLive(b *testing.B, router http.Handler, matchID string) {
	status := models.MatchStatusLive
	updateReq := models.UpdateMatchRequest{
		Status: &status,
	}

	reqBody, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PUT", "/api/v1/matches/"+matchID, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		b.Fatalf("Failed to update match: %d", w.Code)
	}
}

func startScoring(b *testing.B, router http.Handler, matchID string) {
	startReq := models.ScorecardRequest{
		MatchID: matchID,
	}

	reqBody, _ := json.Marshal(startReq)
	req := httptest.NewRequest("POST", "/api/v1/scorecard/start", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		b.Fatalf("Failed to start scoring: %d", w.Code)
	}
}

func addTestBall(b *testing.B, router http.Handler, matchID string) {
	ballReq := models.BallEventRequest{
		MatchID:       matchID,
		InningsNumber: 1,
		BallType:      models.BallTypeGood,
		RunType:       models.RunTypeOne,
		IsWicket:      false,
	}

	reqBody, _ := json.Marshal(ballReq)
	req := httptest.NewRequest("POST", "/api/v1/scorecard/ball", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		b.Fatalf("Failed to add ball: %d", w.Code)
	}
}
