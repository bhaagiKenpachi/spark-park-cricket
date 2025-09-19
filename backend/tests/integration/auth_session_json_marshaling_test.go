package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"spark-park-cricket-backend/internal/config"
	"spark-park-cricket-backend/internal/database"
	"spark-park-cricket-backend/internal/handlers"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAuthSessionJSONMarshaling tests that authentication status checks don't cause JSON marshaling errors
func TestAuthSessionJSONMarshaling(t *testing.T) {
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

	// Initialize services
	serviceContainer := services.NewContainer(dbClient.Repositories, cfg.Config)

	// Create test user
	testUser := &models.User{
		ID:            "test-user-123",
		Email:         "test@example.com",
		Name:          "Test User",
		EmailVerified: true,
	}

	// Create user in database
	err = dbClient.Repositories.User.CreateUser(context.Background(), testUser)
	require.NoError(t, err)

	// Create auth handler
	authHandler := handlers.NewAuthHandler(serviceContainer.AuthService, serviceContainer.SessionService)

	// Test 1: Auth status check without session (should not cause JSON marshaling error)
	t.Run("AuthStatus_NoSession", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/auth/status", nil)
		w := httptest.NewRecorder()

		// This should not cause JSON marshaling errors
		authHandler.AuthStatus(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, false, response["data"].(map[string]interface{})["authenticated"])
	})

	// Test 2: Create session and check auth status (should not cause JSON marshaling error)
	t.Run("AuthStatus_WithSession", func(t *testing.T) {
		// Create a session for the test user
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		err := serviceContainer.SessionService.CreateSession(w, req, testUser)
		require.NoError(t, err)

		// Create new request with session cookie
		req2 := httptest.NewRequest("GET", "/api/v1/auth/status", nil)
		req2.Header.Set("Cookie", w.Header().Get("Set-Cookie"))

		w2 := httptest.NewRecorder()

		// This should not cause JSON marshaling errors
		authHandler.AuthStatus(w2, req2)

		assert.Equal(t, http.StatusOK, w2.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w2.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, true, response["data"].(map[string]interface{})["authenticated"])

		// Verify user data is included
		userData := response["data"].(map[string]interface{})["user"].(map[string]interface{})
		assert.Equal(t, testUser.ID, userData["id"])
		assert.Equal(t, testUser.Email, userData["email"])
		assert.Equal(t, testUser.Name, userData["name"])
	})

	// Test 3: Multiple auth status checks (stress test for JSON marshaling)
	t.Run("AuthStatus_MultipleChecks", func(t *testing.T) {
		// Create a session for the test user
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		err := serviceContainer.SessionService.CreateSession(w, req, testUser)
		require.NoError(t, err)

		// Perform multiple auth status checks
		for i := 0; i < 5; i++ {
			req2 := httptest.NewRequest("GET", "/api/v1/auth/status", nil)
			req2.Header.Set("Cookie", w.Header().Get("Set-Cookie"))

			w2 := httptest.NewRecorder()

			// This should not cause JSON marshaling errors
			authHandler.AuthStatus(w2, req2)

			assert.Equal(t, http.StatusOK, w2.Code)

			var response map[string]interface{}
			err = json.Unmarshal(w2.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.Equal(t, true, response["data"].(map[string]interface{})["authenticated"])
		}
	})
}

// TestSessionValuesJSONMarshaling tests that session values can be safely logged
func TestSessionValuesJSONMarshaling(t *testing.T) {
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

	// Initialize services
	serviceContainer := services.NewContainer(dbClient.Repositories, cfg.Config)

	// Create test user
	testUser := &models.User{
		ID:            "test-user-456",
		Email:         "test2@example.com",
		Name:          "Test User 2",
		EmailVerified: true,
	}

	// Create user in database
	err = dbClient.Repositories.User.CreateUser(context.Background(), testUser)
	require.NoError(t, err)

	// Test that session creation with logging doesn't cause JSON marshaling errors
	t.Run("SessionCreation_WithLogging", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		// This should not cause JSON marshaling errors when logging session values
		err := serviceContainer.SessionService.CreateSession(w, req, testUser)
		require.NoError(t, err)

		// Verify session was created
		cookies := w.Result().Cookies()
		assert.Len(t, cookies, 1)
		assert.Equal(t, "user_session", cookies[0].Name)
	})

	// Test that session destruction with logging doesn't cause JSON marshaling errors
	t.Run("SessionDestruction_WithLogging", func(t *testing.T) {
		// Create a session first
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		err := serviceContainer.SessionService.CreateSession(w, req, testUser)
		require.NoError(t, err)

		// Create new request with session cookie
		req2 := httptest.NewRequest("GET", "/test", nil)
		req2.Header.Set("Cookie", w.Header().Get("Set-Cookie"))
		w2 := httptest.NewRecorder()

		// This should not cause JSON marshaling errors when logging
		err = serviceContainer.SessionService.DestroySession(w2, req2)
		require.NoError(t, err)

		// Verify session was destroyed
		cookies := w2.Result().Cookies()
		assert.Len(t, cookies, 1)
		assert.Equal(t, "user_session", cookies[0].Name)
		assert.Equal(t, -1, cookies[0].MaxAge)
	})
}
