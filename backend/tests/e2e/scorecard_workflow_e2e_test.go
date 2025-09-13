package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"spark-park-cricket-backend/internal/config"
	"spark-park-cricket-backend/internal/database"
	"spark-park-cricket-backend/internal/handlers"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/services"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function to setup test router for scorecard workflow tests
func setupScorecardWorkflowTestRouter(scorecardHandler *handlers.ScorecardHandler, serviceContainer *services.Container) http.Handler {
	router := chi.NewRouter()

	// Add middleware
	router.Use(chimiddleware.Recoverer)
	router.Use(chimiddleware.Logger)
	router.Use(chimiddleware.RequestID)
	router.Use(chimiddleware.RealIP)
	router.Use(chimiddleware.Timeout(60 * time.Second))
	router.Use(corsMiddleware())

	// API routes
	router.Route("/api/v1", func(r chi.Router) {
		// Series routes (needed for creating matches)
		r.Route("/series", func(r chi.Router) {
			seriesHandler := handlers.NewSeriesHandler(serviceContainer.Series)
			r.Get("/", seriesHandler.ListSeries)
			r.Post("/", seriesHandler.CreateSeries)
			r.Get("/{id}", seriesHandler.GetSeries)
			r.Put("/{id}", seriesHandler.UpdateSeries)
			r.Delete("/{id}", seriesHandler.DeleteSeries)
		})
		// Match routes (needed for creating matches)
		r.Route("/matches", func(r chi.Router) {
			matchHandler := handlers.NewMatchHandler(serviceContainer.Match)
			r.Get("/", matchHandler.ListMatches)
			r.Post("/", matchHandler.CreateMatch)
			r.Get("/{id}", matchHandler.GetMatch)
			r.Put("/{id}", matchHandler.UpdateMatch)
			r.Delete("/{id}", matchHandler.DeleteMatch)
			r.Get("/series/{series_id}", matchHandler.GetMatchesBySeries)
		})
		// Scorecard routes
		r.Route("/scorecard", func(r chi.Router) {
			r.Post("/start", scorecardHandler.StartScoring)
			r.Post("/ball", scorecardHandler.AddBall)
			r.Get("/{match_id}", scorecardHandler.GetScorecard)
			r.Get("/{match_id}/current-over", scorecardHandler.GetCurrentOver)
			r.Get("/{match_id}/innings/{innings_number}", scorecardHandler.GetInnings)
			r.Get("/{match_id}/innings/{innings_number}/over/{over_number}", scorecardHandler.GetOver)
		})
	})

	return router
}

func corsMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "*")
			next.ServeHTTP(w, r)
		})
	}
}

// Helper function to clean up test data
func cleanupScorecardWorkflowTestData(t *testing.T, dbClient *database.Client) {
	// Clean up scorecard related tables in reverse order of dependencies
	// Balls -> Overs -> Innings -> Matches -> Series

	// Clean up balls
	_, err := dbClient.Supabase.From("scorecard_balls").Delete("", "").Gte("id", "").ExecuteTo(nil)
	if err != nil {
		t.Logf("Warning: Failed to cleanup balls: %v", err)
	}

	// Clean up overs
	_, err = dbClient.Supabase.From("scorecard_overs").Delete("", "").Gte("id", "").ExecuteTo(nil)
	if err != nil {
		t.Logf("Warning: Failed to cleanup overs: %v", err)
	}

	// Clean up innings
	_, err = dbClient.Supabase.From("scorecard_innings").Delete("", "").Gte("id", "").ExecuteTo(nil)
	if err != nil {
		t.Logf("Warning: Failed to cleanup innings: %v", err)
	}

	// Clean up matches
	_, err = dbClient.Supabase.From("matches").Delete("", "").Gte("id", "").ExecuteTo(nil)
	if err != nil {
		t.Logf("Warning: Failed to cleanup matches: %v", err)
	}

	// Clean up series
	_, err = dbClient.Supabase.From("series").Delete("", "").Gte("id", "").ExecuteTo(nil)
	if err != nil {
		t.Logf("Warning: Failed to cleanup series: %v", err)
	}
}

// Helper function to create a test series for workflow tests
func createTestSeriesForWorkflow(t *testing.T, router http.Handler) string {
	seriesReq := map[string]interface{}{
		"name":        "E2E Test Series " + time.Now().Format("2006-01-02 15:04:05"),
		"description": "E2E test series for scorecard workflow tests",
		"start_date":  time.Now().AddDate(0, 0, 1).Format(time.RFC3339),
		"end_date":    time.Now().AddDate(0, 0, 7).Format(time.RFC3339),
	}

	body, err := json.Marshal(seriesReq)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/api/v1/series", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var response struct {
		Data models.Series `json:"data"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	return response.Data.ID
}

// Helper function to create a test match for workflow tests
func createTestMatchForWorkflow(t *testing.T, router http.Handler, seriesID string) string {
	matchReq := map[string]interface{}{
		"series_id":           seriesID,
		"match_number":        1,
		"date":                time.Now().AddDate(0, 0, 1).Format(time.RFC3339),
		"venue":               "E2E Test Venue",
		"team_a_player_count": 11,
		"team_b_player_count": 11,
		"total_overs":         20,
		"toss_winner":         "A",
		"toss_type":           "H",
		"batting_team":        "A",
	}

	body, err := json.Marshal(matchReq)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/api/v1/matches", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var response struct {
		Data models.Match `json:"data"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	return response.Data.ID
}

// Helper function to update match status to live
func updateMatchToLiveForWorkflow(t *testing.T, router http.Handler, matchID string) {
	updateReq := map[string]interface{}{
		"status": "live",
	}

	body, err := json.Marshal(updateReq)
	require.NoError(t, err)

	req := httptest.NewRequest("PUT", "/api/v1/matches/"+matchID, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

// Helper function to add a ball
func addBallToMatch(t *testing.T, router http.Handler, matchID string, inningsNumber int, ballType, runType string, isWicket bool, byes int, wicketType ...string) {
	req := map[string]interface{}{
		"match_id":       matchID,
		"innings_number": inningsNumber,
		"ball_type":      ballType,
		"run_type":       runType,
		"is_wicket":      isWicket,
		"byes":           byes,
	}

	// Add wicket_type if provided and isWicket is true
	if isWicket && len(wicketType) > 0 && wicketType[0] != "" {
		req["wicket_type"] = wicketType[0]
	}

	body, err := json.Marshal(req)
	require.NoError(t, err)

	reqHTTP := httptest.NewRequest("POST", "/api/v1/scorecard/ball", bytes.NewBuffer(body))
	reqHTTP.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, reqHTTP)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestCompleteScorecardWorkflow(t *testing.T) {
	// Load test configuration
	testConfig := config.LoadTestConfig()

	// Initialize database client
	dbClient, err := database.NewClient(testConfig.Config)
	require.NoError(t, err)
	defer dbClient.Close()

	// Clean up any existing test data
	cleanupScorecardWorkflowTestData(t, dbClient)

	// Initialize services
	serviceContainer := services.NewContainer(dbClient.Repositories)
	scorecardHandler := handlers.NewScorecardHandler(serviceContainer.Scorecard)

	// Setup router
	router := setupScorecardWorkflowTestRouter(scorecardHandler, serviceContainer)

	t.Run("CompleteMatchWorkflow", func(t *testing.T) {
		// Create test series
		seriesID := createTestSeriesForWorkflow(t, router)
		assert.NotEmpty(t, seriesID)

		// Create test match
		matchID := createTestMatchForWorkflow(t, router, seriesID)
		assert.NotEmpty(t, matchID)

		// Update match to live status
		updateMatchToLiveForWorkflow(t, router, matchID)

		// Start scoring
		startReq := map[string]interface{}{
			"match_id": matchID,
		}

		body, err := json.Marshal(startReq)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/api/v1/scorecard/start", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		// Simulate first innings with some balls
		// Over 1: 4, 1, 0, 6, 2, 1 = 14 runs
		addBallToMatch(t, router, matchID, 1, "good", "4", false, 0)
		addBallToMatch(t, router, matchID, 1, "good", "1", false, 0)
		addBallToMatch(t, router, matchID, 1, "good", "0", false, 0)
		addBallToMatch(t, router, matchID, 1, "good", "6", false, 0)
		addBallToMatch(t, router, matchID, 1, "good", "2", false, 0)
		addBallToMatch(t, router, matchID, 1, "good", "1", false, 0)

		// Over 2: 1, W, 0, 4, 1, 2 = 8 runs + 1 wicket
		addBallToMatch(t, router, matchID, 1, "good", "1", false, 0)
		addBallToMatch(t, router, matchID, 1, "good", "WC", true, 0, "bowled") // Wicket
		addBallToMatch(t, router, matchID, 1, "good", "0", false, 0)
		addBallToMatch(t, router, matchID, 1, "good", "4", false, 0)
		addBallToMatch(t, router, matchID, 1, "good", "1", false, 0)
		addBallToMatch(t, router, matchID, 1, "good", "2", false, 0)

		// Check that there's no current over since all 2 overs are complete
		req = httptest.NewRequest("GET", "/api/v1/scorecard/"+matchID+"/current-over?innings=1", nil)
		w = httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code) // No current over exists

		// Check scorecard after first 2 overs
		req = httptest.NewRequest("GET", "/api/v1/scorecard/"+matchID, nil)
		w = httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var scorecardResponse struct {
			Data models.ScorecardResponse `json:"data"`
		}
		err = json.Unmarshal(w.Body.Bytes(), &scorecardResponse)
		require.NoError(t, err)
		assert.Equal(t, matchID, scorecardResponse.Data.MatchID)
		assert.Len(t, scorecardResponse.Data.Innings, 1)
		assert.Equal(t, 1, scorecardResponse.Data.Innings[0].InningsNumber)
		assert.Equal(t, 22, scorecardResponse.Data.Innings[0].TotalRuns) // 14 + 8
		assert.Equal(t, 1, scorecardResponse.Data.Innings[0].TotalWickets)
		assert.Equal(t, 2.0, scorecardResponse.Data.Innings[0].TotalOvers)

		// Check specific innings
		req = httptest.NewRequest("GET", "/api/v1/scorecard/"+matchID+"/innings/1", nil)
		w = httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var inningsResponse struct {
			Data models.InningsSummary `json:"data"`
		}
		err = json.Unmarshal(w.Body.Bytes(), &inningsResponse)
		require.NoError(t, err)
		assert.Equal(t, 1, inningsResponse.Data.InningsNumber)
		assert.Equal(t, 22, inningsResponse.Data.TotalRuns)
		assert.Equal(t, 1, inningsResponse.Data.TotalWickets)
		assert.Equal(t, 2.0, inningsResponse.Data.TotalOvers)

		// Check specific over
		req = httptest.NewRequest("GET", "/api/v1/scorecard/"+matchID+"/innings/1/over/1", nil)
		w = httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var overResponse struct {
			Data models.OverSummary `json:"data"`
		}
		err = json.Unmarshal(w.Body.Bytes(), &overResponse)
		require.NoError(t, err)
		assert.Equal(t, 1, overResponse.Data.OverNumber)
		assert.Equal(t, 14, overResponse.Data.TotalRuns)
		assert.Equal(t, 6, overResponse.Data.TotalBalls)
		assert.Equal(t, 0, overResponse.Data.TotalWickets)
	})

	t.Run("MultipleOversWorkflow", func(t *testing.T) {
		// Create a new series and match for this test
		seriesID := createTestSeriesForWorkflow(t, router)
		matchID := createTestMatchForWorkflow(t, router, seriesID)
		updateMatchToLiveForWorkflow(t, router, matchID)

		// Start scoring
		startReq := map[string]interface{}{
			"match_id": matchID,
		}

		body, err := json.Marshal(startReq)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/api/v1/scorecard/start", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		// Simulate 3 complete overs
		// Over 1: 1, 2, 3, 4, 5, 6 = 21 runs
		addBallToMatch(t, router, matchID, 1, "good", "1", false, 0)
		addBallToMatch(t, router, matchID, 1, "good", "2", false, 0)
		addBallToMatch(t, router, matchID, 1, "good", "3", false, 0)
		addBallToMatch(t, router, matchID, 1, "good", "4", false, 0)
		addBallToMatch(t, router, matchID, 1, "good", "5", false, 0)
		addBallToMatch(t, router, matchID, 1, "good", "6", false, 0)

		// Over 2: 0, 1, 2, 3, 4, 5 = 15 runs
		addBallToMatch(t, router, matchID, 1, "good", "0", false, 0)
		addBallToMatch(t, router, matchID, 1, "good", "1", false, 0)
		addBallToMatch(t, router, matchID, 1, "good", "2", false, 0)
		addBallToMatch(t, router, matchID, 1, "good", "3", false, 0)
		addBallToMatch(t, router, matchID, 1, "good", "4", false, 0)
		addBallToMatch(t, router, matchID, 1, "good", "5", false, 0)

		// Over 3: 6, 0, 1, 2, 3, 4 = 16 runs
		addBallToMatch(t, router, matchID, 1, "good", "6", false, 0)
		addBallToMatch(t, router, matchID, 1, "good", "0", false, 0)
		addBallToMatch(t, router, matchID, 1, "good", "1", false, 0)
		addBallToMatch(t, router, matchID, 1, "good", "2", false, 0)
		addBallToMatch(t, router, matchID, 1, "good", "3", false, 0)
		addBallToMatch(t, router, matchID, 1, "good", "4", false, 0)

		// Check final scorecard
		req = httptest.NewRequest("GET", "/api/v1/scorecard/"+matchID, nil)
		w = httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var scorecardResponse struct {
			Data models.ScorecardResponse `json:"data"`
		}
		err = json.Unmarshal(w.Body.Bytes(), &scorecardResponse)
		require.NoError(t, err)
		assert.Equal(t, matchID, scorecardResponse.Data.MatchID)
		assert.Len(t, scorecardResponse.Data.Innings, 1)
		assert.Equal(t, 1, scorecardResponse.Data.Innings[0].InningsNumber)
		assert.Equal(t, 52, scorecardResponse.Data.Innings[0].TotalRuns) // 21 + 15 + 16
		assert.Equal(t, 0, scorecardResponse.Data.Innings[0].TotalWickets)
		assert.Equal(t, 3.0, scorecardResponse.Data.Innings[0].TotalOvers)

		// Check that there's no current over since all 3 overs are complete
		req = httptest.NewRequest("GET", "/api/v1/scorecard/"+matchID+"/current-over?innings=1", nil)
		w = httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code) // No current over exists
	})

	t.Run("WideAndNoBallWorkflow", func(t *testing.T) {
		// Create a new series and match for this test
		seriesID := createTestSeriesForWorkflow(t, router)
		matchID := createTestMatchForWorkflow(t, router, seriesID)
		updateMatchToLiveForWorkflow(t, router, matchID)

		// Start scoring
		startReq := map[string]interface{}{
			"match_id": matchID,
		}

		body, err := json.Marshal(startReq)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/api/v1/scorecard/start", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		// Simulate over with wides and no balls
		// Ball 1: Wide + 1 run = 2 runs total
		addBallToMatch(t, router, matchID, 1, "wide", "1", false, 0)

		// Ball 2: Good ball, 4 runs
		addBallToMatch(t, router, matchID, 1, "good", "4", false, 0)

		// Ball 3: No ball + 2 runs = 3 runs total
		addBallToMatch(t, router, matchID, 1, "no_ball", "2", false, 0)

		// Ball 4: Good ball, 1 run
		addBallToMatch(t, router, matchID, 1, "good", "1", false, 0)

		// Ball 5: Wide + 2 runs = 3 runs total
		addBallToMatch(t, router, matchID, 1, "wide", "2", false, 0)

		// Ball 6: Good ball, 6 runs
		addBallToMatch(t, router, matchID, 1, "good", "6", false, 0)

		// Check scorecard - should have 17 runs but only 3 legal balls
		req = httptest.NewRequest("GET", "/api/v1/scorecard/"+matchID, nil)
		w = httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var scorecardResponse struct {
			Data models.ScorecardResponse `json:"data"`
		}
		err = json.Unmarshal(w.Body.Bytes(), &scorecardResponse)
		require.NoError(t, err)
		assert.Equal(t, matchID, scorecardResponse.Data.MatchID)
		assert.Len(t, scorecardResponse.Data.Innings, 1)
		assert.Equal(t, 16, scorecardResponse.Data.Innings[0].TotalRuns)   // 1 + 4 + 2 + 1 + 2 + 6
		assert.Equal(t, 3, scorecardResponse.Data.Innings[0].TotalBalls)   // Only good balls count
		assert.Equal(t, 0.3, scorecardResponse.Data.Innings[0].TotalOvers) // 3 balls = 0.3 overs
	})

	// Clean up test data
	cleanupScorecardWorkflowTestData(t, dbClient)
}
