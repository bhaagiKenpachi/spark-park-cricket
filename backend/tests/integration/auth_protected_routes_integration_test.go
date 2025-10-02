package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"spark-park-cricket-backend/internal/config"
	contextkeys "spark-park-cricket-backend/internal/context"
	"spark-park-cricket-backend/internal/database"
	"spark-park-cricket-backend/internal/handlers"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/services"
	"spark-park-cricket-backend/pkg/testutils"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthProtectedRoutesIntegration_SeriesRoutes(t *testing.T) {
	// Skip if not in integration test mode
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Load test configuration
	cfg := config.LoadTestConfig()

	// Initialize database client
	dbClient, err := database.NewTestClient(cfg)
	require.NoError(t, err)
	defer dbClient.Close()

	// Clean up any existing test data
	testutils.CleanupTestData(t, dbClient)

	// Initialize services
	serviceContainer := services.NewContainer(dbClient, cfg.Config)
	_ = serviceContainer // Use serviceContainer to avoid unused variable warning

	// Setup routes
	router := handlers.SetupRoutes(dbClient, cfg.Config)

	// Create test user
	testUser := &models.User{
		GoogleID:      fmt.Sprintf("test-google-id-series-%d", time.Now().UnixNano()),
		Email:         fmt.Sprintf("test-series-%d@example.com", time.Now().UnixNano()),
		Name:          "Test Series User",
		Picture:       "https://example.com/picture.jpg",
		EmailVerified: true,
	}

	err = dbClient.Repositories.User.CreateUser(context.Background(), testUser)
	require.NoError(t, err)
	defer func() { _ = dbClient.Repositories.User.DeleteUser(context.Background(), testUser.ID) }()

	// Create user session using session service
	sessionService := serviceContainer.SessionService

	// Create a mock HTTP request and response writer to create a proper session
	mockReq := httptest.NewRequest("GET", "/", nil)
	mockWriter := httptest.NewRecorder()

	err = sessionService.CreateSession(mockWriter, mockReq, testUser)
	require.NoError(t, err)

	// Extract the session cookie from the response
	cookies := mockWriter.Result().Cookies()
	var sessionCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "user_session" {
			sessionCookie = cookie
			break
		}
	}
	require.NotNil(t, sessionCookie, "Session cookie should be created")

	// Create a series for testing matches
	var seriesID string

	t.Run("Create Series - Unauthenticated", func(t *testing.T) {
		seriesData := map[string]interface{}{
			"name":       "Test Series",
			"start_date": "2024-01-01T00:00:00Z",
			"end_date":   "2024-01-31T00:00:00Z",
		}

		jsonData, _ := json.Marshal(seriesData)
		req := httptest.NewRequest("POST", "/api/v1/series", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		// Check response body contains error
		responseBody := w.Body.String()
		assert.Contains(t, responseBody, "UNAUTHORIZED")
		assert.Contains(t, responseBody, "Authentication required")
	})

	t.Run("Create Series - Authenticated", func(t *testing.T) {
		seriesData := map[string]interface{}{
			"name":       "Test Series Authenticated",
			"start_date": "2024-01-01T00:00:00Z",
			"end_date":   "2024-01-31T00:00:00Z",
		}

		jsonData, _ := json.Marshal(seriesData)
		req := httptest.NewRequest("POST", "/api/v1/series", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		// Add session cookie
		req.AddCookie(sessionCookie)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Should succeed with authentication
		assert.Equal(t, http.StatusCreated, w.Code)

		// Parse response to get series ID for cleanup and use in other tests
		if w.Body.Len() > 0 {
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			if err == nil && response["data"] != nil {
				if data, ok := response["data"].(map[string]interface{}); ok {
					if id, ok := data["id"].(string); ok {
						seriesID = id // Set the outer scope variable
						defer func() { _ = dbClient.Repositories.Series.Delete(context.Background(), seriesID) }()
					}
				}
			}
		}
	})

	t.Run("Update Series - Unauthenticated", func(t *testing.T) {
		// First create a series
		series := &models.Series{
			Name:      "Test Series for Update",
			StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
			CreatedBy: testUser.ID,
		}

		err := dbClient.Repositories.Series.Create(context.Background(), series)
		require.NoError(t, err)
		defer func() { _ = dbClient.Repositories.Series.Delete(context.Background(), series.ID) }()

		updateData := map[string]interface{}{
			"name": "Updated Series Name",
		}

		jsonData, _ := json.Marshal(updateData)
		req := httptest.NewRequest("PUT", fmt.Sprintf("/api/v1/series/%s", series.ID), bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		// Check response body contains error
		responseBody := w.Body.String()
		assert.Contains(t, responseBody, "UNAUTHORIZED")
		assert.Contains(t, responseBody, "Authentication required")
	})

	t.Run("Delete Series - Unauthenticated", func(t *testing.T) {
		// First create a series
		series := &models.Series{
			Name:      "Test Series for Delete",
			StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
			CreatedBy: testUser.ID,
		}

		err := dbClient.Repositories.Series.Create(context.Background(), series)
		require.NoError(t, err)
		defer func() { _ = dbClient.Repositories.Series.Delete(context.Background(), series.ID) }()

		req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/series/%s", series.ID), nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		// Check response body contains error
		responseBody := w.Body.String()
		assert.Contains(t, responseBody, "UNAUTHORIZED")
		assert.Contains(t, responseBody, "Authentication required")
	})

	t.Run("List Series - Public Access", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/series", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Should succeed without authentication
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Get Series - Public Access", func(t *testing.T) {
		// First create a series
		series := &models.Series{
			Name:      "Test Series for Get",
			StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
			CreatedBy: testUser.ID,
		}

		err := dbClient.Repositories.Series.Create(context.Background(), series)
		require.NoError(t, err)
		defer func() { _ = dbClient.Repositories.Series.Delete(context.Background(), series.ID) }()

		req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/series/%s", series.ID), nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Should succeed without authentication
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestAuthProtectedRoutesIntegration_MatchRoutes(t *testing.T) {
	// Skip if not in integration test mode
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Load test configuration
	cfg := config.LoadTestConfig()

	// Initialize database client
	dbClient, err := database.NewTestClient(cfg)
	require.NoError(t, err)
	defer dbClient.Close()

	// Clean up any existing test data
	testutils.CleanupTestData(t, dbClient)

	// Initialize services
	serviceContainer := services.NewContainer(dbClient, cfg.Config)
	_ = serviceContainer // Use serviceContainer to avoid unused variable warning

	// Setup routes
	router := handlers.SetupRoutes(dbClient, cfg.Config)

	// Create test user
	testUser := &models.User{
		GoogleID:      fmt.Sprintf("test-google-id-match-%d", time.Now().UnixNano()),
		Email:         fmt.Sprintf("test-match-%d@example.com", time.Now().UnixNano()),
		Name:          "Test Match User",
		Picture:       "https://example.com/picture.jpg",
		EmailVerified: true,
	}

	err = dbClient.Repositories.User.CreateUser(context.Background(), testUser)
	require.NoError(t, err)
	defer func() { _ = dbClient.Repositories.User.DeleteUser(context.Background(), testUser.ID) }()

	// Create user session using session service
	sessionService := serviceContainer.SessionService

	// Create a mock HTTP request and response writer to create a proper session
	mockReq := httptest.NewRequest("GET", "/", nil)
	mockWriter := httptest.NewRecorder()

	err = sessionService.CreateSession(mockWriter, mockReq, testUser)
	require.NoError(t, err)

	// Extract the session cookie from the response
	cookies := mockWriter.Result().Cookies()
	var sessionCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "user_session" {
			sessionCookie = cookie
			break
		}
	}
	require.NotNil(t, sessionCookie, "Session cookie should be created")

	// Create test series for matches
	series := &models.Series{
		Name:      "Test Series for Match",
		StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
		CreatedBy: testUser.ID,
	}

	err = dbClient.Repositories.Series.Create(context.Background(), series)
	require.NoError(t, err)
	defer func() { _ = dbClient.Repositories.Series.Delete(context.Background(), series.ID) }()

	t.Run("Create Match - Unauthenticated", func(t *testing.T) {
		matchData := map[string]interface{}{
			"series_id":           series.ID,
			"date":                "2024-01-15T00:00:00Z",
			"team_a_player_count": 11,
			"team_b_player_count": 11,
			"total_overs":         20,
			"toss_winner":         "A",
			"toss_type":           "H",
		}

		jsonData, _ := json.Marshal(matchData)
		req := httptest.NewRequest("POST", "/api/v1/matches", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		// Check response body contains error
		responseBody := w.Body.String()
		assert.Contains(t, responseBody, "UNAUTHORIZED")
		assert.Contains(t, responseBody, "Authentication required")
	})

	t.Run("Create Match - Authenticated", func(t *testing.T) {
		matchData := map[string]interface{}{
			"series_id":           series.ID,
			"date":                "2024-01-15T00:00:00Z",
			"team_a_player_count": 11,
			"team_b_player_count": 11,
			"total_overs":         20,
			"toss_winner":         "A",
			"toss_type":           "H",
		}

		jsonData, _ := json.Marshal(matchData)
		req := httptest.NewRequest("POST", "/api/v1/matches", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		// Add session cookie
		req.AddCookie(sessionCookie)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Should succeed with authentication
		assert.Equal(t, http.StatusCreated, w.Code)

		// Parse response to get match ID for cleanup
		if w.Body.Len() > 0 {
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			if err == nil && response["data"] != nil {
				if data, ok := response["data"].(map[string]interface{}); ok {
					if matchID, ok := data["id"].(string); ok {
						defer func() { _ = dbClient.Repositories.Match.Delete(context.Background(), matchID) }()
					}
				}
			}
		}
	})

	t.Run("List Matches - Public Access", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/matches", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Should succeed without authentication
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Get Match - Public Access", func(t *testing.T) {
		// First create a match
		match := &models.Match{
			SeriesID:         series.ID,
			MatchNumber:      1,
			Date:             time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			Status:           models.MatchStatusLive,
			TeamAPlayerCount: 11,
			TeamBPlayerCount: 11,
			TotalOvers:       20,
			TossWinner:       models.TeamTypeA,
			TossType:         models.TossTypeHeads,
			BattingTeam:      models.TeamTypeA,
			CreatedBy:        testUser.ID,
		}

		err := dbClient.Repositories.Match.Create(context.Background(), match)
		require.NoError(t, err)
		defer func() { _ = dbClient.Repositories.Match.Delete(context.Background(), match.ID) }()

		req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/matches/%s", match.ID), nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Should succeed without authentication
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestAuthProtectedRoutesIntegration_SessionValidation(t *testing.T) {
	// Skip if not in integration test mode
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Load test configuration
	cfg := config.LoadTestConfig()

	// Initialize database client
	dbClient, err := database.NewTestClient(cfg)
	require.NoError(t, err)
	defer dbClient.Close()

	// Setup routes
	router := handlers.SetupRoutes(dbClient, cfg.Config)

	t.Run("Invalid Session Cookie", func(t *testing.T) {
		seriesData := map[string]interface{}{
			"name":        "Test Series Invalid Session",
			"description": "Test Description Invalid Session",
			"start_date":  "2024-01-01",
			"end_date":    "2024-01-31",
		}

		jsonData, _ := json.Marshal(seriesData)
		req := httptest.NewRequest("POST", "/api/v1/series", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		// Add invalid session cookie
		req.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    "invalid-session-id",
			Path:     "/",
			HttpOnly: true,
		})

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		// Check response body contains error
		responseBody := w.Body.String()
		assert.Contains(t, responseBody, "UNAUTHORIZED")
		assert.Contains(t, responseBody, "Authentication required")
	})

	t.Run("Expired Session", func(t *testing.T) {
		// Create test user
		testUser := &models.User{
			GoogleID:      "test-google-id-expired",
			Email:         "test-expired@example.com",
			Name:          "Test Expired User",
			Picture:       "https://example.com/picture.jpg",
			EmailVerified: true,
		}

		err = dbClient.Repositories.User.CreateUser(context.Background(), testUser)
		require.NoError(t, err)
		defer func() { _ = dbClient.Repositories.User.DeleteUser(context.Background(), testUser.ID) }()

		// Create expired session
		expiredSession := &models.UserSession{
			UserID:    testUser.ID,
			SessionID: "test-session-expired-123",
			ExpiresAt: time.Now().Add(-1 * time.Hour), // Expired 1 hour ago
		}

		err = dbClient.Repositories.User.CreateUserSession(context.Background(), expiredSession)
		require.NoError(t, err)
		defer func() {
			_ = dbClient.Repositories.User.DeleteUserSession(context.Background(), expiredSession.SessionID)
		}()

		seriesData := map[string]interface{}{
			"name":        "Test Series Expired Session",
			"description": "Test Description Expired Session",
			"start_date":  "2024-01-01",
			"end_date":    "2024-01-31",
		}

		jsonData, _ := json.Marshal(seriesData)
		req := httptest.NewRequest("POST", "/api/v1/series", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		// Add expired session cookie
		req.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    expiredSession.SessionID,
			Path:     "/",
			HttpOnly: true,
		})

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		// Check response body contains error
		responseBody := w.Body.String()
		assert.Contains(t, responseBody, "UNAUTHORIZED")
		assert.Contains(t, responseBody, "Authentication required")
	})

	t.Run("No Session Cookie", func(t *testing.T) {
		seriesData := map[string]interface{}{
			"name":        "Test Series No Session",
			"description": "Test Description No Session",
			"start_date":  "2024-01-01",
			"end_date":    "2024-01-31",
		}

		jsonData, _ := json.Marshal(seriesData)
		req := httptest.NewRequest("POST", "/api/v1/series", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		// No session cookie added
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		// Check response body contains error
		responseBody := w.Body.String()
		assert.Contains(t, responseBody, "UNAUTHORIZED")
		assert.Contains(t, responseBody, "Authentication required")
	})
}

func TestAuthProtectedRoutesIntegration_ContextValues(t *testing.T) {
	// Skip if not in integration test mode
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Load test configuration
	cfg := config.LoadTestConfig()

	// Initialize database client
	dbClient, err := database.NewTestClient(cfg)
	require.NoError(t, err)
	defer dbClient.Close()

	// Clean up any existing test data
	testutils.CleanupTestData(t, dbClient)

	// Initialize services
	serviceContainer := services.NewContainer(dbClient, cfg.Config)
	_ = serviceContainer // Use serviceContainer to avoid unused variable warning

	// Setup routes
	router := handlers.SetupRoutes(dbClient, cfg.Config)
	_ = router // Use router to avoid unused variable warning

	// Create test user
	testUser := &models.User{
		GoogleID:      fmt.Sprintf("test-google-id-context-%d", time.Now().UnixNano()),
		Email:         fmt.Sprintf("test-context-%d@example.com", time.Now().UnixNano()),
		Name:          "Test Context User",
		Picture:       "https://example.com/picture.jpg",
		EmailVerified: true,
	}

	err = dbClient.Repositories.User.CreateUser(context.Background(), testUser)
	require.NoError(t, err)
	defer func() { _ = dbClient.Repositories.User.DeleteUser(context.Background(), testUser.ID) }()

	// Create user session using session service
	sessionService := serviceContainer.SessionService

	// Create a mock HTTP request and response writer to create a proper session
	mockReq := httptest.NewRequest("GET", "/", nil)
	mockWriter := httptest.NewRecorder()

	err = sessionService.CreateSession(mockWriter, mockReq, testUser)
	require.NoError(t, err)

	// Extract the session cookie from the response
	cookies := mockWriter.Result().Cookies()
	var sessionCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "user_session" {
			sessionCookie = cookie
			break
		}
	}
	require.NotNil(t, sessionCookie, "Session cookie should be created")

	t.Run("Context Values in Authenticated Request", func(t *testing.T) {
		// Create a test handler that checks context values
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if user context values are set using the correct context keys
			user := r.Context().Value(contextkeys.UserKey)
			userID := r.Context().Value(contextkeys.UserIDKey)
			userEmail := r.Context().Value(contextkeys.UserEmailKey)

			assert.NotNil(t, user)
			assert.NotNil(t, userID)
			assert.NotNil(t, userEmail)

			// Verify values match test user
			if user != nil {
				userObj := user.(*models.User)
				assert.Equal(t, testUser.ID, userObj.ID)
				assert.Equal(t, testUser.Email, userObj.Email)
			}

			if userID != nil {
				assert.Equal(t, testUser.ID, userID.(string))
			}

			if userEmail != nil {
				assert.Equal(t, testUser.Email, userEmail.(string))
			}

			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("context test passed"))
		})

		// Create a test router with auth middleware
		testRouter := chi.NewRouter()
		testRouter.Use(services.AuthMiddleware(serviceContainer.SessionService))
		testRouter.Post("/test-context", testHandler)

		req := httptest.NewRequest("POST", "/test-context", nil)

		// Add session cookie
		req.AddCookie(sessionCookie)

		w := httptest.NewRecorder()

		testRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "context test passed", w.Body.String())
	})
}
