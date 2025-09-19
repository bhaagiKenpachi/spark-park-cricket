package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"spark-park-cricket-backend/pkg/testutils"
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

func TestMatchIntegration(t *testing.T) {
	// Load test configuration
	testConfig := config.LoadTestConfig()

	// Initialize database client
	dbClient, err := database.NewClient(testConfig.Config)
	require.NoError(t, err)
	defer dbClient.Close()

	// Clean up any existing test data
	cleanupMatchTestData(t, dbClient)

	// Initialize services
	serviceContainer := services.NewContainer(dbClient.Repositories, testConfig.Config)
	matchHandler := handlers.NewMatchHandler(serviceContainer.Match)

	// Setup router
	router := setupMatchTestRouter(matchHandler, serviceContainer)

	t.Run("Complete Match CRUD Flow", func(t *testing.T) {
		testCompleteMatchCRUDFlow(t, router, dbClient)
	})

	t.Run("Match Pagination", func(t *testing.T) {
		// Clean up before pagination test to ensure isolation
		cleanupMatchTestData(t, dbClient)
		testMatchPagination(t, router, dbClient)
	})

	t.Run("Match Validation", func(t *testing.T) {
		// Clean up before validation test to ensure isolation
		cleanupMatchTestData(t, dbClient)
		testMatchValidation(t, router)
	})

	t.Run("Match Error Handling", func(t *testing.T) {
		// Clean up before error handling test to ensure isolation
		cleanupMatchTestData(t, dbClient)
		testMatchErrorHandling(t, router)
	})
}

func testCompleteMatchCRUDFlow(t *testing.T, router http.Handler, dbClient *database.Client) {
	// Initialize services to get session service
	serviceContainer := services.NewContainer(dbClient.Repositories, config.LoadTestConfig().Config)

	// Create authenticated test user with proper session cookie
	user, sessionCookie := testutils.CreateAuthenticatedTestUserWithSessionService(t, dbClient, serviceContainer.SessionService)
	_ = user // Use user if needed for assertions

	// First, create a series to associate with the match
	seriesID := createTestSeriesWithAuth(t, router, sessionCookie)

	// Create a match
	createReq := models.CreateMatchRequest{
		SeriesID:         seriesID,
		MatchNumber:      intPtr(1),
		Date:             time.Date(2025, 9, 14, 10, 0, 0, 0, time.UTC),
		TeamAPlayerCount: 11,
		TeamBPlayerCount: 11,
		TotalOvers:       20,
		TossWinner:       models.TeamTypeA,
		TossType:         models.TossTypeHeads,
	}

	createBody, err := json.Marshal(createReq)
	require.NoError(t, err)

	req := testutils.CreateAuthenticatedRequestWithCookie("POST", "/api/v1/matches", createBody, sessionCookie)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var response struct {
		Data models.Match `json:"data"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	createdMatch := response.Data
	assert.NotEmpty(t, createdMatch.ID)
	assert.Equal(t, createReq.SeriesID, createdMatch.SeriesID)
	assert.Equal(t, *createReq.MatchNumber, createdMatch.MatchNumber)
	assert.Equal(t, createReq.Date.UTC().Truncate(time.Second), createdMatch.Date.Truncate(time.Second))
	assert.Equal(t, models.MatchStatusLive, createdMatch.Status)
	assert.Equal(t, createReq.TossWinner, createdMatch.TossWinner)
	assert.Equal(t, createReq.TossWinner, createdMatch.BattingTeam) // Winner bats first

	// Get the created match
	req = testutils.CreateAuthenticatedRequestWithCookie("GET", fmt.Sprintf("/api/v1/matches/%s", createdMatch.ID), nil, sessionCookie)
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var getResponse struct {
		Data models.Match `json:"data"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &getResponse)
	require.NoError(t, err)
	retrievedMatch := getResponse.Data
	assert.Equal(t, createdMatch.ID, retrievedMatch.ID)
	assert.Equal(t, createdMatch.SeriesID, retrievedMatch.SeriesID)
	assert.Equal(t, createdMatch.Status, retrievedMatch.Status)

	// Update the match
	updateReq := models.UpdateMatchRequest{
		Status: matchStatusPtr(models.MatchStatusCompleted),
	}

	updateBody, err := json.Marshal(updateReq)
	require.NoError(t, err)

	req = testutils.CreateAuthenticatedRequestWithCookie("PUT", fmt.Sprintf("/api/v1/matches/%s", createdMatch.ID), updateBody, sessionCookie)
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var updateResponse struct {
		Data models.Match `json:"data"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &updateResponse)
	require.NoError(t, err)
	updatedMatch := updateResponse.Data
	assert.Equal(t, *updateReq.Status, updatedMatch.Status)

	// Verify the updated match directly by getting it
	req = testutils.CreateAuthenticatedRequestWithCookie("GET", fmt.Sprintf("/api/v1/matches/%s", updatedMatch.ID), nil, sessionCookie)
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var verifyResponse struct {
		Data models.Match `json:"data"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &verifyResponse)
	require.NoError(t, err)
	verifiedMatch := verifyResponse.Data
	assert.Equal(t, *updateReq.Status, verifiedMatch.Status)

	// Delete the match
	req = testutils.CreateAuthenticatedRequestWithCookie("DELETE", fmt.Sprintf("/api/v1/matches/%s", createdMatch.ID), nil, sessionCookie)
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify match is deleted
	req = httptest.NewRequest("GET", fmt.Sprintf("/api/v1/matches/%s", createdMatch.ID), nil)
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func testMatchPagination(t *testing.T, router http.Handler, dbClient *database.Client) {
	// Initialize services to get session service
	serviceContainer := services.NewContainer(dbClient.Repositories, config.LoadTestConfig().Config)

	// Create authenticated test user with proper session cookie
	user, sessionCookie := testutils.CreateAuthenticatedTestUserWithSessionService(t, dbClient, serviceContainer.SessionService)
	_ = user // Use user if needed for assertions

	// First, create a series to associate with matches
	seriesID := createTestSeriesWithAuth(t, router, sessionCookie)

	// Store created match IDs to verify they exist
	var createdMatchIDs []string

	// Create multiple matches for pagination testing
	for i := 1; i <= 5; i++ {
		createReq := models.CreateMatchRequest{
			SeriesID:         seriesID,
			MatchNumber:      intPtr(i),
			Date:             time.Date(2025, 9, 14+i, 10, 0, 0, 0, time.UTC),
			TeamAPlayerCount: 11,
			TeamBPlayerCount: 11,
			TotalOvers:       20,
			TossWinner:       models.TeamTypeA,
			TossType:         models.TossTypeHeads,
		}

		createBody, err := json.Marshal(createReq)
		require.NoError(t, err)

		req := testutils.CreateAuthenticatedRequestWithCookie("POST", "/api/v1/matches", createBody, sessionCookie)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		// Parse response to get match ID
		var createResponse struct {
			Data models.Match `json:"data"`
		}
		err = json.Unmarshal(w.Body.Bytes(), &createResponse)
		require.NoError(t, err)
		createdMatchIDs = append(createdMatchIDs, createResponse.Data.ID)
	}

	// Test pagination with limit
	req := testutils.CreateAuthenticatedRequestWithCookie("GET", "/api/v1/matches?limit=3", nil, sessionCookie)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var listResponse struct {
		Data []models.Match `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &listResponse)
	require.NoError(t, err)
	matchesList := listResponse.Data
	assert.GreaterOrEqual(t, len(matchesList), 3, "Should have at least 3 matches")

	// Test pagination with offset
	req = testutils.CreateAuthenticatedRequestWithCookie("GET", "/api/v1/matches?limit=2&offset=2", nil, sessionCookie)
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	err = json.Unmarshal(w.Body.Bytes(), &listResponse)
	require.NoError(t, err)
	matchesList = listResponse.Data
	assert.GreaterOrEqual(t, len(matchesList), 2, "Should have at least 2 matches")

	// Test invalid pagination parameters
	req = testutils.CreateAuthenticatedRequestWithCookie("GET", "/api/v1/matches?limit=invalid&offset=-1", nil, sessionCookie)
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code) // Should use default values

	err = json.Unmarshal(w.Body.Bytes(), &listResponse)
	require.NoError(t, err)
	matchesList = listResponse.Data

	// Verify that all 5 matches we created are present
	assert.Equal(t, 5, len(matchesList), "Should have exactly 5 matches with default pagination")

	// Verify that all created match IDs are present in the response
	responseMatchIDs := make(map[string]bool)
	for _, match := range matchesList {
		responseMatchIDs[match.ID] = true
	}

	for _, createdID := range createdMatchIDs {
		assert.True(t, responseMatchIDs[createdID], "Created match ID %s should be present in response", createdID)
	}
}

func testMatchValidation(t *testing.T, router http.Handler) {
	// Test invalid date range
	createReq := models.CreateMatchRequest{
		SeriesID:         "nonexistent-series",
		MatchNumber:      intPtr(1),
		Date:             time.Date(2025, 9, 14, 10, 0, 0, 0, time.UTC),
		TeamAPlayerCount: 11,
		TeamBPlayerCount: 11,
		TotalOvers:       20,
		TossWinner:       models.TeamTypeA,
		TossType:         models.TossTypeHeads,
	}

	createBody, err := json.Marshal(createReq)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/api/v1/matches", bytes.NewBuffer(createBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code) // Business logic error returns 500

	// Test missing required fields - this should return 500 for business logic errors (invalid UUID)
	invalidReq := map[string]interface{}{
		"series_id": "test-series", // Invalid UUID format
		"date":      "2025-09-14T10:00:00Z",
		// Missing team_a_player_count, team_b_player_count, total_overs, toss_winner, toss_type
	}

	invalidBody, err := json.Marshal(invalidReq)
	require.NoError(t, err)

	req = httptest.NewRequest("POST", "/api/v1/matches", bytes.NewBuffer(invalidBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)
	// Invalid UUID format causes business logic error, not validation error
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Test invalid JSON
	req = httptest.NewRequest("POST", "/api/v1/matches", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func testMatchErrorHandling(t *testing.T, router http.Handler) {
	// Test getting non-existent match
	req := httptest.NewRequest("GET", "/api/v1/matches/non-existent-id", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	// Test updating non-existent match
	updateReq := models.UpdateMatchRequest{
		Status: matchStatusPtr(models.MatchStatusCompleted),
	}

	updateBody, err := json.Marshal(updateReq)
	require.NoError(t, err)

	req = httptest.NewRequest("PUT", "/api/v1/matches/non-existent-id", bytes.NewBuffer(updateBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Test deleting non-existent match
	req = httptest.NewRequest("DELETE", "/api/v1/matches/non-existent-id", nil)
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Test empty match ID - use a route that would have an empty ID parameter
	req = httptest.NewRequest("GET", "/api/v1/matches/", nil)
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)
	// This actually hits the list endpoint, so it should return 200
	assert.Equal(t, http.StatusOK, w.Code)
}

// Helper function to create a test series
func createTestSeries(t *testing.T, router http.Handler) string {
	createReq := models.CreateSeriesRequest{
		Name:      "Test Series for Match",
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

// Helper function to create a test series with authentication
func createTestSeriesWithAuth(t *testing.T, router http.Handler, sessionCookie string) string {
	createReq := models.CreateSeriesRequest{
		Name:      "Test Series for Match",
		StartDate: time.Date(2025, 9, 14, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, 9, 21, 0, 0, 0, 0, time.UTC),
	}

	createBody, err := json.Marshal(createReq)
	require.NoError(t, err)

	req := testutils.CreateAuthenticatedRequestWithCookie("POST", "/api/v1/series", createBody, sessionCookie)
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

// Helper function to clean up test data
func cleanupMatchTestData(t *testing.T, dbClient *database.Client) {
	// Clean up matches table - delete all records
	_, err := dbClient.Supabase.From("matches").Delete("", "").Gte("created_at", "1970-01-01T00:00:00Z").ExecuteTo(nil)
	if err != nil {
		t.Logf("Warning: Failed to cleanup match test data: %v", err)
	}

	// Clean up series table as well
	_, err = dbClient.Supabase.From("series").Delete("", "").Gte("created_at", "1970-01-01T00:00:00Z").ExecuteTo(nil)
	if err != nil {
		t.Logf("Warning: Failed to cleanup series test data: %v", err)
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

// Helper function to setup test router for match tests
func setupMatchTestRouter(matchHandler *handlers.MatchHandler, serviceContainer *services.Container) http.Handler {
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
	router.Use(testutils.CORSMiddleware())
	router.Use(middleware.AuthMiddleware(serviceContainer.SessionService))

	// API routes
	router.Route("/api/v1", func(r chi.Router) {
		// Series routes (needed for creating matches)
		r.Route("/series", func(r chi.Router) {
			seriesHandler := handlers.NewSeriesHandler(serviceContainer.Series)
			r.Post("/", seriesHandler.CreateSeries)
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
