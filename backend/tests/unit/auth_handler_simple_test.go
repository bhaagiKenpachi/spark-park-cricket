package unit

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"spark-park-cricket-backend/internal/handlers"
	"spark-park-cricket-backend/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthHandler_GetCurrentUser_Authenticated(t *testing.T) {
	// Create mocks
	mockAuthService := new(MockAuthService)
	mockSessionService := new(MockSessionService)

	// Create test user
	testUser := &models.User{
		ID:    "user-123",
		Email: "test@example.com",
		Name:  "Test User",
	}

	// Setup mock expectations
	mockSessionService.On("GetSession", mock.AnythingOfType("*http.Request")).Return(testUser, nil)

	// Create auth handler
	authHandler := handlers.NewAuthHandler(mockAuthService, mockSessionService)

	// Create request
	req := httptest.NewRequest("GET", "/auth/me", nil)
	w := httptest.NewRecorder()

	// Call handler
	authHandler.GetCurrentUser(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	// Check response body contains user data
	responseBody := w.Body.String()
	assert.Contains(t, responseBody, testUser.ID)
	assert.Contains(t, responseBody, testUser.Email)
	assert.Contains(t, responseBody, testUser.Name)

	mockSessionService.AssertExpectations(t)
}

func TestAuthHandler_GetCurrentUser_Unauthenticated(t *testing.T) {
	// Create mocks
	mockAuthService := new(MockAuthService)
	mockSessionService := new(MockSessionService)

	// Setup mock expectations - return error for unauthenticated user
	mockSessionService.On("GetSession", mock.AnythingOfType("*http.Request")).Return(nil, errors.New("not authenticated"))

	// Create auth handler
	authHandler := handlers.NewAuthHandler(mockAuthService, mockSessionService)

	// Create request
	req := httptest.NewRequest("GET", "/auth/me", nil)
	w := httptest.NewRecorder()

	// Call handler
	authHandler.GetCurrentUser(w, req)

	// Assertions
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Check response body contains error
	responseBody := w.Body.String()
	assert.Contains(t, responseBody, "UNAUTHORIZED")
	assert.Contains(t, responseBody, "Not authenticated")

	mockSessionService.AssertExpectations(t)
}

func TestAuthHandler_Logout_Success(t *testing.T) {
	// Create mocks
	mockAuthService := new(MockAuthService)
	mockSessionService := new(MockSessionService)

	// Setup mock expectations
	mockSessionService.On("DestroySession", mock.AnythingOfType("*httptest.ResponseRecorder"), mock.AnythingOfType("*http.Request")).Return(nil)

	// Create auth handler
	authHandler := handlers.NewAuthHandler(mockAuthService, mockSessionService)

	// Create request
	req := httptest.NewRequest("POST", "/auth/logout", nil)
	w := httptest.NewRecorder()

	// Call handler
	authHandler.Logout(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	// Check response body contains success message
	responseBody := w.Body.String()
	assert.Contains(t, responseBody, "Logged out successfully")

	mockSessionService.AssertExpectations(t)
}

func TestAuthHandler_Logout_Error(t *testing.T) {
	// Create mocks
	mockAuthService := new(MockAuthService)
	mockSessionService := new(MockSessionService)

	// Setup mock expectations - return error
	mockSessionService.On("DestroySession", mock.AnythingOfType("*httptest.ResponseRecorder"), mock.AnythingOfType("*http.Request")).Return(errors.New("logout error"))

	// Create auth handler
	authHandler := handlers.NewAuthHandler(mockAuthService, mockSessionService)

	// Create request
	req := httptest.NewRequest("POST", "/auth/logout", nil)
	w := httptest.NewRecorder()

	// Call handler
	authHandler.Logout(w, req)

	// Assertions
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Check response body contains error
	responseBody := w.Body.String()
	assert.Contains(t, responseBody, "LOGOUT_ERROR")
	assert.Contains(t, responseBody, "Failed to logout")

	mockSessionService.AssertExpectations(t)
}

func TestAuthHandler_AuthStatus_Authenticated(t *testing.T) {
	// Create mocks
	mockAuthService := new(MockAuthService)
	mockSessionService := new(MockSessionService)

	// Create test user
	testUser := &models.User{
		ID:    "user-123",
		Email: "test@example.com",
		Name:  "Test User",
	}

	// Setup mock expectations
	mockSessionService.On("IsAuthenticated", mock.AnythingOfType("*http.Request")).Return(true)
	mockSessionService.On("GetSession", mock.AnythingOfType("*http.Request")).Return(testUser, nil)

	// Create auth handler
	authHandler := handlers.NewAuthHandler(mockAuthService, mockSessionService)

	// Create request
	req := httptest.NewRequest("GET", "/auth/status", nil)
	w := httptest.NewRecorder()

	// Call handler
	authHandler.AuthStatus(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	// Check response body contains authenticated status and user data
	responseBody := w.Body.String()
	assert.Contains(t, responseBody, "authenticated")
	assert.Contains(t, responseBody, "true")
	assert.Contains(t, responseBody, testUser.ID)
	assert.Contains(t, responseBody, testUser.Email)

	mockSessionService.AssertExpectations(t)
}

func TestAuthHandler_AuthStatus_Unauthenticated(t *testing.T) {
	// Create mocks
	mockAuthService := new(MockAuthService)
	mockSessionService := new(MockSessionService)

	// Setup mock expectations
	mockSessionService.On("IsAuthenticated", mock.AnythingOfType("*http.Request")).Return(false)

	// Create auth handler
	authHandler := handlers.NewAuthHandler(mockAuthService, mockSessionService)

	// Create request
	req := httptest.NewRequest("GET", "/auth/status", nil)
	w := httptest.NewRecorder()

	// Call handler
	authHandler.AuthStatus(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	// Check response body contains unauthenticated status
	responseBody := w.Body.String()
	assert.Contains(t, responseBody, "authenticated")
	assert.Contains(t, responseBody, "false")

	mockSessionService.AssertExpectations(t)
}
