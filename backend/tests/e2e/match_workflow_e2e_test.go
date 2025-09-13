package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"spark-park-cricket-backend/internal/config"
	"spark-park-cricket-backend/internal/database"
	"spark-park-cricket-backend/internal/handlers"
	"spark-park-cricket-backend/internal/middleware"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/services"
)

func TestMatchWorkflow_E2E(t *testing.T) {
	// Load test configuration
	testConfig := config.LoadTestConfig()

	// Initialize database client
	dbClient, err := database.NewClient(testConfig.Config)
	require.NoError(t, err)
	defer dbClient.Close()

	// Clean up any existing test data
	cleanupMatchWorkflowTestData(t, dbClient)

	// Initialize services
	serviceContainer := services.NewContainer(dbClient.Repositories)
	matchHandler := handlers.NewMatchHandler(serviceContainer.Match)

	// Setup router
	router := setupMatchWorkflowTestRouter(matchHandler, serviceContainer)

	t.Run("Complete Match Lifecycle", func(t *testing.T) {
		testCompleteMatchLifecycle(t, router, dbClient)
	})

	t.Run("Match State Transitions", func(t *testing.T) {
		testMatchStateTransitions(t, router, dbClient)
	})

	t.Run("Match Series Integration", func(t *testing.T) {
		testMatchSeriesIntegration(t, router, dbClient)
	})

	t.Run("Match Validation Workflow", func(t *testing.T) {
		testMatchValidationWorkflow(t, router, dbClient)
	})
}

func testCompleteMatchLifecycle(t *testing.T, router http.Handler, dbClient *database.Client) {
	// Step 1: Create a series
	seriesID := createTestSeriesForWorkflow(t, router, "Test Series for Match Lifecycle")

	// Step 2: Create a match in the series
	matchID := createTestMatchForWorkflow(t, router, seriesID, 1)

	// Step 3: Verify match is in "live" status
	match := getMatch(t, router, matchID)
	assert.Equal(t, models.MatchStatusLive, match.Status)
	assert.Equal(t, models.TeamTypeA, match.BattingTeam) // Toss winner bats first

	// Step 4: Update match status to completed
	updateMatchStatus(t, router, matchID, models.MatchStatusCompleted)

	// Step 5: Verify match status change
	updatedMatch := getMatch(t, router, matchID)
	assert.Equal(t, models.MatchStatusCompleted, updatedMatch.Status)

	// Step 6: Create multiple matches in the same series
	match2ID := createTestMatchForWorkflow(t, router, seriesID, 2)
	match3ID := createTestMatchForWorkflow(t, router, seriesID, 3)

	// Step 7: List all matches in the series
	matches := getMatchesBySeries(t, router, seriesID)
	assert.GreaterOrEqual(t, len(matches), 3, "Should have at least 3 matches in the series")

	// Step 8: Verify match numbers are unique
	matchNumbers := make(map[int]bool)
	for _, match := range matches {
		assert.False(t, matchNumbers[match.MatchNumber], "Match numbers should be unique")
		matchNumbers[match.MatchNumber] = true
	}

	// Step 9: Delete the matches
	deleteMatch(t, router, matchID)
	deleteMatch(t, router, match2ID)
	deleteMatch(t, router, match3ID)

	// Step 10: Verify matches are deleted
	assertMatchNotFound(t, router, matchID)
	assertMatchNotFound(t, router, match2ID)
	assertMatchNotFound(t, router, match3ID)
}

func testMatchStateTransitions(t *testing.T, router http.Handler, dbClient *database.Client) {
	// Create a series and match
	seriesID := createTestSeriesForWorkflow(t, router, "Test Series for State Transitions")
	matchID := createTestMatchForWorkflow(t, router, seriesID, 1)

	// Test state transitions: live -> completed -> cancelled
	transitions := []models.MatchStatus{
		models.MatchStatusLive,
		models.MatchStatusCompleted,
		models.MatchStatusCancelled,
	}

	for i, expectedStatus := range transitions {
		if i > 0 { // Skip first iteration as match starts in "live" status
			updateMatchStatus(t, router, matchID, expectedStatus)
		}

		match := getMatch(t, router, matchID)
		assert.Equal(t, expectedStatus, match.Status, "Match status should be %s", expectedStatus)

		// Verify the status change is persisted
		match = getMatch(t, router, matchID)
		assert.Equal(t, expectedStatus, match.Status, "Match status should persist as %s", expectedStatus)
	}

	// Clean up
	deleteMatch(t, router, matchID)
}

func testMatchSeriesIntegration(t *testing.T, router http.Handler, dbClient *database.Client) {
	// Create multiple series
	series1ID := createTestSeriesForWorkflow(t, router, "Series 1")
	series2ID := createTestSeriesForWorkflow(t, router, "Series 2")

	// Create matches in different series
	match1ID := createTestMatchForWorkflow(t, router, series1ID, 1)
	match2ID := createTestMatchForWorkflow(t, router, series1ID, 2)
	match3ID := createTestMatchForWorkflow(t, router, series2ID, 1)

	// Verify matches are associated with correct series
	match1 := getMatch(t, router, match1ID)
	match2 := getMatch(t, router, match2ID)
	match3 := getMatch(t, router, match3ID)

	assert.Equal(t, series1ID, match1.SeriesID)
	assert.Equal(t, series1ID, match2.SeriesID)
	assert.Equal(t, series2ID, match3.SeriesID)

	// Test getting matches by series
	series1Matches := getMatchesBySeries(t, router, series1ID)
	series2Matches := getMatchesBySeries(t, router, series2ID)

	// Series 1 should have 2 matches, Series 2 should have 1 match
	assert.GreaterOrEqual(t, len(series1Matches), 2, "Series 1 should have at least 2 matches")
	assert.GreaterOrEqual(t, len(series2Matches), 1, "Series 2 should have at least 1 match")

	// Verify all matches in series 1 belong to series 1
	for _, match := range series1Matches {
		assert.Equal(t, series1ID, match.SeriesID, "All matches in series 1 should belong to series 1")
	}

	// Verify all matches in series 2 belong to series 2
	for _, match := range series2Matches {
		assert.Equal(t, series2ID, match.SeriesID, "All matches in series 2 should belong to series 2")
	}

	// Test match number uniqueness within series
	series1MatchNumbers := make(map[int]bool)
	for _, match := range series1Matches {
		assert.False(t, series1MatchNumbers[match.MatchNumber], "Match numbers should be unique within series 1")
		series1MatchNumbers[match.MatchNumber] = true
	}

	// Clean up
	deleteMatch(t, router, match1ID)
	deleteMatch(t, router, match2ID)
	deleteMatch(t, router, match3ID)
}

func testMatchValidationWorkflow(t *testing.T, router http.Handler, dbClient *database.Client) {
	// Create a series for testing
	seriesID := createTestSeriesForWorkflow(t, router, "Test Series for Validation")

	// Test 1: Try to create match with non-existent series
	invalidMatchReq := models.CreateMatchRequest{
		SeriesID:         "non-existent-series",
		MatchNumber:      intPtr(1),
		Date:             time.Date(2025, 9, 14, 10, 0, 0, 0, time.UTC),
		TeamAPlayerCount: 11,
		TeamBPlayerCount: 11,
		TotalOvers:       20,
		TossWinner:       models.TeamTypeA,
		TossType:         models.TossTypeHeads,
	}

	createMatchExpectingError(t, router, invalidMatchReq, http.StatusInternalServerError)

	// Test 2: Create match with duplicate match number
	validMatchReq := models.CreateMatchRequest{
		SeriesID:         seriesID,
		MatchNumber:      intPtr(1),
		Date:             time.Date(2025, 9, 14, 10, 0, 0, 0, time.UTC),
		TeamAPlayerCount: 11,
		TeamBPlayerCount: 11,
		TotalOvers:       20,
		TossWinner:       models.TeamTypeA,
		TossType:         models.TossTypeHeads,
	}

	// Create first match successfully
	matchID := createTestMatchWithRequest(t, router, validMatchReq)

	// Try to create second match with same number - should fail
	duplicateMatchReq := validMatchReq
	duplicateMatchReq.Date = time.Date(2025, 9, 15, 10, 0, 0, 0, time.UTC)
	createMatchExpectingError(t, router, duplicateMatchReq, http.StatusInternalServerError)

	// Test 3: Create match without match number (auto-increment)
	autoMatchReq := models.CreateMatchRequest{
		SeriesID:         seriesID,
		MatchNumber:      nil, // Should auto-increment to 2
		Date:             time.Date(2025, 9, 16, 10, 0, 0, 0, time.UTC),
		TeamAPlayerCount: 11,
		TeamBPlayerCount: 11,
		TotalOvers:       20,
		TossWinner:       models.TeamTypeB,
		TossType:         models.TossTypeTails,
	}

	autoMatchID := createTestMatchWithRequest(t, router, autoMatchReq)
	autoMatch := getMatch(t, router, autoMatchID)
	assert.Equal(t, 2, autoMatch.MatchNumber, "Auto-incremented match number should be 2")

	// Test 4: Verify toss winner becomes batting team
	assert.Equal(t, models.TeamTypeB, autoMatch.TossWinner)
	assert.Equal(t, models.TeamTypeB, autoMatch.BattingTeam, "Toss winner should be batting team")

	// Test 5: Update batting team
	updateBattingTeam(t, router, autoMatchID, models.TeamTypeA)
	updatedMatch := getMatch(t, router, autoMatchID)
	assert.Equal(t, models.TeamTypeA, updatedMatch.BattingTeam, "Batting team should be updated")

	// Clean up
	deleteMatch(t, router, matchID)
	deleteMatch(t, router, autoMatchID)
}

// Helper functions for workflow testing

func createTestSeriesForWorkflow(t *testing.T, router http.Handler, name string) string {
	createReq := models.CreateSeriesRequest{
		Name:      name,
		StartDate: time.Date(2025, 9, 14, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, 9, 21, 0, 0, 0, 0, time.UTC),
	}

	createBody, err := json.Marshal(createReq)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/api/v1/series", bytes.NewBuffer(createBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var response struct {
		Data models.Series `json:"data"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	return response.Data.ID
}

func createTestMatchForWorkflow(t *testing.T, router http.Handler, seriesID string, matchNumber int) string {
	createReq := models.CreateMatchRequest{
		SeriesID:         seriesID,
		MatchNumber:      intPtr(matchNumber),
		Date:             time.Date(2025, 9, 14+matchNumber, 10, 0, 0, 0, time.UTC),
		TeamAPlayerCount: 11,
		TeamBPlayerCount: 11,
		TotalOvers:       20,
		TossWinner:       models.TeamTypeA,
		TossType:         models.TossTypeHeads,
	}

	return createTestMatchWithRequest(t, router, createReq)
}

func createTestMatchWithRequest(t *testing.T, router http.Handler, createReq models.CreateMatchRequest) string {
	createBody, err := json.Marshal(createReq)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/api/v1/matches", bytes.NewBuffer(createBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var response struct {
		Data models.Match `json:"data"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	return response.Data.ID
}

func createMatchExpectingError(t *testing.T, router http.Handler, createReq models.CreateMatchRequest, expectedStatus int) {
	createBody, err := json.Marshal(createReq)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/api/v1/matches", bytes.NewBuffer(createBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, expectedStatus, w.Code, "Expected status %d but got %d", expectedStatus, w.Code)
}

func getMatch(t *testing.T, router http.Handler, matchID string) models.Match {
	req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/matches/%s", matchID), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Data models.Match `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	return response.Data
}

func updateMatchStatus(t *testing.T, router http.Handler, matchID string, status models.MatchStatus) {
	updateReq := models.UpdateMatchRequest{
		Status: &status,
	}

	updateBody, err := json.Marshal(updateReq)
	require.NoError(t, err)

	req := httptest.NewRequest("PUT", fmt.Sprintf("/api/v1/matches/%s", matchID), bytes.NewBuffer(updateBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func updateBattingTeam(t *testing.T, router http.Handler, matchID string, battingTeam models.TeamType) {
	updateReq := models.UpdateMatchRequest{
		BattingTeam: &battingTeam,
	}

	updateBody, err := json.Marshal(updateReq)
	require.NoError(t, err)

	req := httptest.NewRequest("PUT", fmt.Sprintf("/api/v1/matches/%s", matchID), bytes.NewBuffer(updateBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func getMatchesBySeries(t *testing.T, router http.Handler, seriesID string) []models.Match {
	req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/matches/series/%s", seriesID), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Data []models.Match `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	return response.Data
}

func deleteMatch(t *testing.T, router http.Handler, matchID string) {
	// First, update the match status to completed so it can be deleted
	updateReq := map[string]interface{}{
		"status": "completed",
	}
	reqBody, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PUT", fmt.Sprintf("/api/v1/matches/%s", matchID), bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Now delete the match
	req = httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/matches/%s", matchID), nil)
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func assertMatchNotFound(t *testing.T, router http.Handler, matchID string) {
	req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/matches/%s", matchID), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// Helper function to clean up test data
func cleanupMatchWorkflowTestData(t *testing.T, dbClient *database.Client) {
	// Clean up matches table
	_, err := dbClient.Supabase.From("matches").Delete("", "").Gte("id", "").ExecuteTo(nil)
	if err != nil {
		_, err = dbClient.Supabase.From("matches").Delete("", "").Gte("created_at", "1970-01-01T00:00:00Z").ExecuteTo(nil)
		if err != nil {
			t.Logf("Warning: Failed to cleanup match workflow test data: %v", err)
		}
	}

	// Clean up series table
	_, err = dbClient.Supabase.From("series").Delete("", "").Gte("id", "").ExecuteTo(nil)
	if err != nil {
		_, err = dbClient.Supabase.From("series").Delete("", "").Gte("created_at", "1970-01-01T00:00:00Z").ExecuteTo(nil)
		if err != nil {
			t.Logf("Warning: Failed to cleanup series workflow test data: %v", err)
		}
	}
}

// Helper function to create int pointer
func intPtr(i int) *int {
	return &i
}

// Helper function to create match status pointer
func matchStatusPtr(status models.MatchStatus) *models.MatchStatus {
	return &status
}

// Helper function to create team type pointer
func teamTypePtr(teamType models.TeamType) *models.TeamType {
	return &teamType
}

// Helper function to setup test router for match workflow tests
func setupMatchWorkflowTestRouter(matchHandler *handlers.MatchHandler, serviceContainer *services.Container) http.Handler {
	router := chi.NewRouter()

	// Add middleware
	router.Use(middleware.RecoveryMiddleware)
	router.Use(middleware.LoggerMiddleware)
	router.Use(middleware.RequestIDMiddleware)
	router.Use(chimiddleware.RealIP)
	router.Use(middleware.TimeoutMiddleware(60 * time.Second))
	router.Use(middleware.SecurityMiddleware)
	router.Use(middleware.ValidationMiddleware)
	router.Use(middleware.MetricsMiddleware)
	router.Use(middleware.RateLimitMiddleware(100))
	router.Use(corsMiddleware())

	// API routes
	router.Route("/api/v1", func(r chi.Router) {
		// Series routes
		r.Route("/series", func(r chi.Router) {
			seriesHandler := handlers.NewSeriesHandler(serviceContainer.Series)
			r.Get("/", seriesHandler.ListSeries)
			r.Post("/", seriesHandler.CreateSeries)
			r.Get("/{id}", seriesHandler.GetSeries)
			r.Put("/{id}", seriesHandler.UpdateSeries)
			r.Delete("/{id}", seriesHandler.DeleteSeries)
		})

		// Match routes
		r.Route("/matches", func(r chi.Router) {
			r.Get("/", matchHandler.ListMatches)
			r.Post("/", matchHandler.CreateMatch)
			r.Get("/{id}", matchHandler.GetMatch)
			r.Put("/{id}", matchHandler.UpdateMatch)
			r.Delete("/{id}", matchHandler.DeleteMatch)
			r.Get("/series/{series_id}", matchHandler.GetMatchesBySeries)
		})
	})

	return router
}

// CORS middleware
func corsMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Cache-Control, Pragma, Expires, Accept")
			w.Header().Set("Access-Control-Max-Age", "86400")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
