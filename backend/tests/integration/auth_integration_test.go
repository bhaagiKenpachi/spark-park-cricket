package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"spark-park-cricket-backend/internal/config"
	"spark-park-cricket-backend/internal/database"
	"spark-park-cricket-backend/internal/handlers"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/services"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthIntegration_UserFlow(t *testing.T) {
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

	// Create auth handler
	authHandler := handlers.NewAuthHandler(serviceContainer.AuthService, serviceContainer.SessionService)

	t.Run("Google Login Redirect", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/auth/google", nil)
		w := httptest.NewRecorder()

		authHandler.GoogleLogin(w, req)

		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
		assert.Contains(t, w.Header().Get("Location"), "accounts.google.com")
	})

	t.Run("Auth Status - Not Authenticated", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/auth/status", nil)
		w := httptest.NewRecorder()

		authHandler.AuthStatus(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		// Should return authenticated: false
	})

	t.Run("Get Current User - Not Authenticated", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/auth/me", nil)
		w := httptest.NewRecorder()

		authHandler.GetCurrentUser(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Logout - Not Authenticated", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/auth/logout", nil)
		w := httptest.NewRecorder()

		authHandler.Logout(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestAuthIntegration_SessionManagement(t *testing.T) {
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

	t.Run("Create User", func(t *testing.T) {
		// Create a test user
		user := &models.User{
			GoogleID:      "test-google-id-123",
			Email:         "test@example.com",
			Name:          "Test User",
			Picture:       "https://example.com/picture.jpg",
			EmailVerified: true,
		}

		err := dbClient.Repositories.User.CreateUser(context.Background(), user)
		require.NoError(t, err)
		assert.NotEmpty(t, user.ID)

		// Clean up
		defer dbClient.Repositories.User.DeleteUser(context.Background(), user.ID)
	})

	t.Run("Get User by Google ID", func(t *testing.T) {
		// Create a test user
		user := &models.User{
			GoogleID:      "test-google-id-456",
			Email:         "test2@example.com",
			Name:          "Test User 2",
			Picture:       "https://example.com/picture2.jpg",
			EmailVerified: true,
		}

		err := dbClient.Repositories.User.CreateUser(context.Background(), user)
		require.NoError(t, err)
		defer dbClient.Repositories.User.DeleteUser(context.Background(), user.ID)

		// Retrieve user by Google ID
		retrievedUser, err := dbClient.Repositories.User.GetUserByGoogleID(context.Background(), user.GoogleID)
		require.NoError(t, err)
		assert.Equal(t, user.GoogleID, retrievedUser.GoogleID)
		assert.Equal(t, user.Email, retrievedUser.Email)
		assert.Equal(t, user.Name, retrievedUser.Name)
	})

	t.Run("User Session Management", func(t *testing.T) {
		// Create a test user
		user := &models.User{
			GoogleID:      "test-google-id-789",
			Email:         "test3@example.com",
			Name:          "Test User 3",
			Picture:       "https://example.com/picture3.jpg",
			EmailVerified: true,
		}

		err := dbClient.Repositories.User.CreateUser(context.Background(), user)
		require.NoError(t, err)
		defer dbClient.Repositories.User.DeleteUser(context.Background(), user.ID)

		// Create a user session
		session := &models.UserSession{
			UserID:    user.ID,
			SessionID: "test-session-id-123",
			ExpiresAt: time.Now().Add(24 * time.Hour),
		}

		err = dbClient.Repositories.User.CreateUserSession(context.Background(), session)
		require.NoError(t, err)
		defer dbClient.Repositories.User.DeleteUserSession(context.Background(), session.SessionID)

		// Retrieve the session
		retrievedSession, err := dbClient.Repositories.User.GetUserSession(context.Background(), session.SessionID)
		require.NoError(t, err)
		assert.Equal(t, session.UserID, retrievedSession.UserID)
		assert.Equal(t, session.SessionID, retrievedSession.SessionID)

		// Update last login
		err = dbClient.Repositories.User.UpdateLastLogin(context.Background(), user.ID)
		require.NoError(t, err)

		// Verify user was updated
		updatedUser, err := dbClient.Repositories.User.GetUserByID(context.Background(), user.ID)
		require.NoError(t, err)
		assert.NotNil(t, updatedUser.LastLoginAt)
	})
}
