package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"spark-park-cricket-backend/internal/config"
	"spark-park-cricket-backend/internal/database"
	"spark-park-cricket-backend/internal/handlers"
	"spark-park-cricket-backend/internal/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthE2E_CompleteAuthenticationFlow(t *testing.T) {
	// Skip if not in e2e test mode
	if testing.Short() {
		t.Skip("Skipping e2e test")
	}

	// Load test configuration
	cfg := config.Load()

	// Initialize database client
	dbClient, err := database.NewClient(cfg)
	require.NoError(t, err)
	defer dbClient.Close()

	// Setup routes
	router := handlers.SetupRoutes(dbClient, cfg)

	t.Run("Complete Authentication Flow", func(t *testing.T) {
		// Step 1: Check initial auth status (should be unauthenticated)
		req := httptest.NewRequest("GET", "/api/v1/auth/status", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var authStatusResponse map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &authStatusResponse)
		require.NoError(t, err)

		authenticated := authStatusResponse["data"].(map[string]interface{})["authenticated"].(bool)
		assert.False(t, authenticated)

		// Step 2: Try to access protected route (should fail)
		seriesData := map[string]interface{}{
			"name":       "Test E2E Series",
			"start_date": "2024-01-01",
			"end_date":   "2024-01-31",
		}

		jsonData, _ := json.Marshal(seriesData)
		req = httptest.NewRequest("POST", "/api/v1/series", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		// Step 3: Create a test user and session (simulating successful OAuth)
		testUser := &models.User{
			GoogleID:      "test-google-id-e2e",
			Email:         "test-e2e@example.com",
			Name:          "Test E2E User",
			Picture:       "https://example.com/picture.jpg",
			EmailVerified: true,
		}

		err = dbClient.Repositories.User.CreateUser(context.Background(), testUser)
		require.NoError(t, err)
		defer dbClient.Repositories.User.DeleteUser(context.Background(), testUser.ID)

		// Create user session
		session := &models.UserSession{
			UserID:    testUser.ID,
			SessionID: "test-session-e2e-123",
			ExpiresAt: time.Now().Add(24 * time.Hour),
		}

		err = dbClient.Repositories.User.CreateUserSession(context.Background(), session)
		require.NoError(t, err)
		defer dbClient.Repositories.User.DeleteUserSession(context.Background(), session.SessionID)

		// Step 4: Check auth status with valid session (should be authenticated)
		req = httptest.NewRequest("GET", "/api/v1/auth/status", nil)
		req.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    session.SessionID,
			Path:     "/",
			HttpOnly: true,
		})
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		err = json.Unmarshal(w.Body.Bytes(), &authStatusResponse)
		require.NoError(t, err)

		authenticated = authStatusResponse["data"].(map[string]interface{})["authenticated"].(bool)
		assert.True(t, authenticated)

		// Verify user data is present
		userData := authStatusResponse["data"].(map[string]interface{})["user"].(map[string]interface{})
		assert.Equal(t, testUser.ID, userData["id"])
		assert.Equal(t, testUser.Email, userData["email"])
		assert.Equal(t, testUser.Name, userData["name"])

		// Step 5: Access protected route with valid session (should succeed)
		req = httptest.NewRequest("POST", "/api/v1/series", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    session.SessionID,
			Path:     "/",
			HttpOnly: true,
		})
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		// Parse response to get series ID for cleanup
		var seriesResponse map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &seriesResponse)
		require.NoError(t, err)

		seriesID := seriesResponse["data"].(map[string]interface{})["id"].(string)
		defer dbClient.Repositories.Series.Delete(context.Background(), seriesID)

		// Step 6: Get current user info
		req = httptest.NewRequest("GET", "/api/v1/auth/me", nil)
		req.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    session.SessionID,
			Path:     "/",
			HttpOnly: true,
		})
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var userResponse map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &userResponse)
		require.NoError(t, err)

		userInfo := userResponse["data"].(map[string]interface{})["user"].(map[string]interface{})
		assert.Equal(t, testUser.ID, userInfo["id"])
		assert.Equal(t, testUser.Email, userInfo["email"])
		assert.Equal(t, testUser.Name, userInfo["name"])

		// Step 7: Logout
		req = httptest.NewRequest("POST", "/api/v1/auth/logout", nil)
		req.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    session.SessionID,
			Path:     "/",
			HttpOnly: true,
		})
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var logoutResponse map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &logoutResponse)
		require.NoError(t, err)

		message := logoutResponse["data"].(map[string]interface{})["message"].(string)
		assert.Equal(t, "Logged out successfully", message)

		// Step 8: Check auth status after logout (should be unauthenticated)
		req = httptest.NewRequest("GET", "/api/v1/auth/status", nil)
		req.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    session.SessionID,
			Path:     "/",
			HttpOnly: true,
		})
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		err = json.Unmarshal(w.Body.Bytes(), &authStatusResponse)
		require.NoError(t, err)

		authenticated = authStatusResponse["data"].(map[string]interface{})["authenticated"].(bool)
		assert.False(t, authenticated)

		// Step 9: Try to access protected route after logout (should fail)
		req = httptest.NewRequest("POST", "/api/v1/series", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    session.SessionID,
			Path:     "/",
			HttpOnly: true,
		})
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestAuthE2E_ProtectedRoutesAccess(t *testing.T) {
	// Skip if not in e2e test mode
	if testing.Short() {
		t.Skip("Skipping e2e test")
	}

	// Load test configuration
	cfg := config.Load()

	// Initialize database client
	dbClient, err := database.NewClient(cfg)
	require.NoError(t, err)
	defer dbClient.Close()

	// Setup routes
	router := handlers.SetupRoutes(dbClient, cfg)

	// Create test user
	testUser := &models.User{
		GoogleID:      "test-google-id-protected",
		Email:         "test-protected@example.com",
		Name:          "Test Protected User",
		Picture:       "https://example.com/picture.jpg",
		EmailVerified: true,
	}

	err = dbClient.Repositories.User.CreateUser(context.Background(), testUser)
	require.NoError(t, err)
	defer dbClient.Repositories.User.DeleteUser(context.Background(), testUser.ID)

	// Create user session
	session := &models.UserSession{
		UserID:    testUser.ID,
		SessionID: "test-session-protected-123",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	err = dbClient.Repositories.User.CreateUserSession(context.Background(), session)
	require.NoError(t, err)
	defer dbClient.Repositories.User.DeleteUserSession(context.Background(), session.SessionID)

	// Create test series for matches
	series := &models.Series{
		Name:      "Test Series for Protected Routes",
		StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
		CreatedBy: testUser.ID,
	}

	err = dbClient.Repositories.Series.Create(context.Background(), series)
	require.NoError(t, err)
	defer dbClient.Repositories.Series.Delete(context.Background(), series.ID)

	t.Run("Series CRUD Operations with Authentication", func(t *testing.T) {
		// Create series
		seriesData := map[string]interface{}{
			"name":       "E2E Test Series",
			"start_date": "2024-01-01",
			"end_date":   "2024-01-31",
		}

		jsonData, _ := json.Marshal(seriesData)
		req := httptest.NewRequest("POST", "/api/v1/series", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    session.SessionID,
			Path:     "/",
			HttpOnly: true,
		})
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var createResponse map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &createResponse)
		require.NoError(t, err)

		createdSeriesID := createResponse["data"].(map[string]interface{})["id"].(string)
		defer dbClient.Repositories.Series.Delete(context.Background(), createdSeriesID)

		// Update series
		updateData := map[string]interface{}{
			"name": "Updated E2E Test Series",
		}

		jsonUpdateData, _ := json.Marshal(updateData)
		req = httptest.NewRequest("PUT", fmt.Sprintf("/api/v1/series/%s", createdSeriesID), bytes.NewBuffer(jsonUpdateData))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    session.SessionID,
			Path:     "/",
			HttpOnly: true,
		})
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Delete series
		req = httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/series/%s", createdSeriesID), nil)
		req.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    session.SessionID,
			Path:     "/",
			HttpOnly: true,
		})
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Match CRUD Operations with Authentication", func(t *testing.T) {
		// Create match
		matchData := map[string]interface{}{
			"series_id":           series.ID,
			"date":                "2024-01-15",
			"team_a_player_count": 11,
			"team_b_player_count": 11,
			"total_overs":         20,
			"toss_winner":         "A",
			"toss_type":           "H",
		}

		jsonData, _ := json.Marshal(matchData)
		req := httptest.NewRequest("POST", "/api/v1/matches", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    session.SessionID,
			Path:     "/",
			HttpOnly: true,
		})
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var createResponse map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &createResponse)
		require.NoError(t, err)

		createdMatchID := createResponse["data"].(map[string]interface{})["id"].(string)
		defer dbClient.Repositories.Match.Delete(context.Background(), createdMatchID)

		// Update match
		updateData := map[string]interface{}{
			"team_a_player_count": 12,
		}

		jsonUpdateData, _ := json.Marshal(updateData)
		req = httptest.NewRequest("PUT", fmt.Sprintf("/api/v1/matches/%s", createdMatchID), bytes.NewBuffer(jsonUpdateData))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    session.SessionID,
			Path:     "/",
			HttpOnly: true,
		})
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Delete match
		req = httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/matches/%s", createdMatchID), nil)
		req.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    session.SessionID,
			Path:     "/",
			HttpOnly: true,
		})
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Public Routes Access Without Authentication", func(t *testing.T) {
		// List series (public)
		req := httptest.NewRequest("GET", "/api/v1/series", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Get specific series (public)
		req = httptest.NewRequest("GET", fmt.Sprintf("/api/v1/series/%s", series.ID), nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// List matches (public)
		req = httptest.NewRequest("GET", "/api/v1/matches", nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestAuthE2E_SessionExpiration(t *testing.T) {
	// Skip if not in e2e test mode
	if testing.Short() {
		t.Skip("Skipping e2e test")
	}

	// Load test configuration
	cfg := config.Load()

	// Initialize database client
	dbClient, err := database.NewClient(cfg)
	require.NoError(t, err)
	defer dbClient.Close()

	// Setup routes
	router := handlers.SetupRoutes(dbClient, cfg)

	// Create test user
	testUser := &models.User{
		GoogleID:      "test-google-id-expiration",
		Email:         "test-expiration@example.com",
		Name:          "Test Expiration User",
		Picture:       "https://example.com/picture.jpg",
		EmailVerified: true,
	}

	err = dbClient.Repositories.User.CreateUser(context.Background(), testUser)
	require.NoError(t, err)
	defer dbClient.Repositories.User.DeleteUser(context.Background(), testUser.ID)

	t.Run("Expired Session Handling", func(t *testing.T) {
		// Create expired session
		expiredSession := &models.UserSession{
			UserID:    testUser.ID,
			SessionID: "test-session-expired-e2e",
			ExpiresAt: time.Now().Add(-1 * time.Hour), // Expired 1 hour ago
		}

		err = dbClient.Repositories.User.CreateUserSession(context.Background(), expiredSession)
		require.NoError(t, err)
		defer dbClient.Repositories.User.DeleteUserSession(context.Background(), expiredSession.SessionID)

		// Try to access protected route with expired session
		seriesData := map[string]interface{}{
			"name":       "Test Expired Session Series",
			"start_date": "2024-01-01",
			"end_date":   "2024-01-31",
		}

		jsonData, _ := json.Marshal(seriesData)
		req := httptest.NewRequest("POST", "/api/v1/series", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
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

	t.Run("Valid Session After Expiration", func(t *testing.T) {
		// Create valid session
		validSession := &models.UserSession{
			UserID:    testUser.ID,
			SessionID: "test-session-valid-e2e",
			ExpiresAt: time.Now().Add(24 * time.Hour), // Valid for 24 hours
		}

		err = dbClient.Repositories.User.CreateUserSession(context.Background(), validSession)
		require.NoError(t, err)
		defer dbClient.Repositories.User.DeleteUserSession(context.Background(), validSession.SessionID)

		// Access protected route with valid session
		seriesData := map[string]interface{}{
			"name":       "Test Valid Session Series",
			"start_date": "2024-01-01",
			"end_date":   "2024-01-31",
		}

		jsonData, _ := json.Marshal(seriesData)
		req := httptest.NewRequest("POST", "/api/v1/series", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    validSession.SessionID,
			Path:     "/",
			HttpOnly: true,
		})
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		// Parse response to get series ID for cleanup
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		seriesID := response["data"].(map[string]interface{})["id"].(string)
		defer dbClient.Repositories.Series.Delete(context.Background(), seriesID)
	})
}
