package tests

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
	contextkeys "spark-park-cricket-backend/internal/context"
	"spark-park-cricket-backend/internal/database"
	"spark-park-cricket-backend/internal/handlers"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/services"
	"spark-park-cricket-backend/pkg/testutils"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Setup router with user routes
func setupUserRouter(serviceContainer *services.Container) http.Handler {
	router := chi.NewRouter()
	userHandler := handlers.NewUserHandler(serviceContainer.UserService)

	router.Put("/api/v1/users/me", addUserContextForUser(userHandler.UpdateUserName))

	return router
}

// Helper function to add user context to request
func addUserContextForUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-User-ID")
		if userID != "" {
			ctx := context.WithValue(r.Context(), contextkeys.UserIDKey, userID)
			next(w, r.WithContext(ctx))
		} else {
			next(w, r)
		}
	}
}

func TestUpdateUserName_Integration(t *testing.T) {
	// Load test configuration
	cfg := config.LoadTestConfig()

	// Initialize database client
	dbClient, err := database.NewTestClient(cfg)
	require.NoError(t, err)
	defer dbClient.Close()

	// Clean up any existing test data
	testutils.CleanupAllTestData(t, dbClient)

	// Initialize services
	serviceContainer := services.NewContainer(dbClient, cfg.Config)
	router := setupUserRouter(serviceContainer)

	// Create test user
	ctx := context.Background()
	testUser := &models.User{
		GoogleID:      fmt.Sprintf("test-user-%d", time.Now().UnixNano()),
		Email:         fmt.Sprintf("test-%d@example.com", time.Now().UnixNano()),
		Name:          "Test User",
		Picture:       "https://example.com/test.jpg",
		EmailVerified: true,
	}
	err = dbClient.Repositories.User.CreateUser(ctx, testUser)
	require.NoError(t, err)
	defer func() { _ = dbClient.Repositories.User.DeleteUser(ctx, testUser.ID) }()

	tests := []struct {
		name           string
		userID         string
		requestBody    map[string]interface{}
		expectedStatus int
		expectedError  string
		validate       func(*testing.T, *models.User)
	}{
		{
			name:   "Success - Update user name",
			userID: testUser.ID,
			requestBody: map[string]interface{}{
				"name": "Updated Name",
			},
			expectedStatus: http.StatusOK,
			validate: func(t *testing.T, user *models.User) {
				assert.Equal(t, "Updated Name", user.Name)
				assert.Equal(t, testUser.ID, user.ID)
				assert.Equal(t, testUser.Email, user.Email)
			},
		},
		{
			name:   "Error - Empty name",
			userID: testUser.ID,
			requestBody: map[string]interface{}{
				"name": "",
			},
			expectedStatus: http.StatusUnprocessableEntity, // 422 for JSON validation errors
			expectedError:  "VALIDATION_ERROR",
		},
		{
			name:   "Error - Name too short",
			userID: testUser.ID,
			requestBody: map[string]interface{}{
				"name": "A",
			},
			expectedStatus: http.StatusInternalServerError, // 500 because service validation returns error
			expectedError:  "UPDATE_ERROR",
		},
		{
			name:   "Error - Name too long",
			userID: testUser.ID,
			requestBody: map[string]interface{}{
				"name": string(make([]byte, 300)),
			},
			expectedStatus: http.StatusInternalServerError, // 500 because service validation returns error
			expectedError:  "UPDATE_ERROR",
		},
		{
			name:   "Error - No authentication",
			userID: "", // Empty user ID simulates no authentication
			requestBody: map[string]interface{}{
				"name": "New Name",
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "UNAUTHORIZED",
		},
		{
			name:   "Error - Invalid JSON",
			userID: testUser.ID,
			requestBody: map[string]interface{}{
				"invalid": "field",
			},
			expectedStatus: http.StatusUnprocessableEntity, // 422 for JSON validation errors
			expectedError:  "VALIDATION_ERROR",
		},
		{
			name:   "Success - Name with special characters",
			userID: testUser.ID,
			requestBody: map[string]interface{}{
				"name": "José María O'Connor",
			},
			expectedStatus: http.StatusOK,
			validate: func(t *testing.T, user *models.User) {
				assert.Equal(t, "José María O'Connor", user.Name)
			},
		},
		{
			name:   "Success - Name with whitespace (should be trimmed)",
			userID: testUser.ID,
			requestBody: map[string]interface{}{
				"name": "  Trimmed Name  ",
			},
			expectedStatus: http.StatusOK,
			validate: func(t *testing.T, user *models.User) {
				assert.Equal(t, "Trimmed Name", user.Name)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare request body
			body, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			// Create request
			req := httptest.NewRequest(http.MethodPut, "/api/v1/users/me", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			if tt.userID != "" {
				req.Header.Set("X-User-ID", tt.userID)
			}

			// Create response recorder
			w := httptest.NewRecorder()

			// Execute request
			router.ServeHTTP(w, req)

			// Assert status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Parse response
			var response struct {
				Data  *models.User `json:"data,omitempty"`
				Error *struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				} `json:"error,omitempty"`
			}
			err = json.NewDecoder(w.Body).Decode(&response)
			require.NoError(t, err)

			if tt.expectedStatus == http.StatusOK {
				// Assert success response
				assert.NotNil(t, response.Data)
				assert.Nil(t, response.Error)

				// Run custom validation if provided
				if tt.validate != nil {
					tt.validate(t, response.Data)
				}
			} else {
				// Assert error response
				assert.Nil(t, response.Data)
				assert.NotNil(t, response.Error)
				if tt.expectedError != "" {
					assert.Contains(t, response.Error.Code, tt.expectedError)
				}
			}
		})
	}
}

func TestUpdateUserName_Concurrent_Integration(t *testing.T) {
	// Load test configuration
	cfg := config.LoadTestConfig()

	// Initialize database client
	dbClient, err := database.NewTestClient(cfg)
	require.NoError(t, err)
	defer dbClient.Close()

	// Clean up any existing test data
	testutils.CleanupAllTestData(t, dbClient)

	// Initialize services
	serviceContainer := services.NewContainer(dbClient, cfg.Config)
	router := setupUserRouter(serviceContainer)

	// Create test user
	ctx := context.Background()
	testUser := &models.User{
		GoogleID:      fmt.Sprintf("test-concurrent-user-%d", time.Now().UnixNano()),
		Email:         fmt.Sprintf("concurrent-%d@example.com", time.Now().UnixNano()),
		Name:          "Concurrent User",
		Picture:       "https://example.com/concurrent.jpg",
		EmailVerified: true,
	}
	err = dbClient.Repositories.User.CreateUser(ctx, testUser)
	require.NoError(t, err)
	defer func() { _ = dbClient.Repositories.User.DeleteUser(ctx, testUser.ID) }()

	// Make concurrent requests
	done := make(chan bool, 5)
	for i := 0; i < 5; i++ {
		go func(index int) {
			requestBody := map[string]interface{}{
				"name": "Name " + string(rune('A'+index)),
			}
			body, _ := json.Marshal(requestBody)

			req := httptest.NewRequest(http.MethodPut, "/api/v1/users/me", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-User-ID", testUser.ID)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Each request should succeed
			assert.Equal(t, http.StatusOK, w.Code)
			done <- true
		}(i)
	}

	// Wait for all requests to complete
	for i := 0; i < 5; i++ {
		<-done
	}

	// Verify final state - name should be one of the values
	user, err := dbClient.Repositories.User.GetUserByID(ctx, testUser.ID)
	require.NoError(t, err)
	assert.NotEmpty(t, user.Name)
	t.Logf("Final user name after concurrent updates: %s", user.Name)
}
