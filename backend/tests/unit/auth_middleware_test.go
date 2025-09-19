package unit

import (
	"net/http"
	"net/http/httptest"
	"spark-park-cricket-backend/internal/middleware"
	"spark-park-cricket-backend/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthMiddleware_AuthenticatedUser(t *testing.T) {
	// Create mock session service
	mockSessionService := new(MockSessionService)

	// Create test user
	testUser := &models.User{
		ID:    "user-123",
		Email: "test@example.com",
		Name:  "Test User",
	}

	// Setup mock expectations
	mockSessionService.On("GetSession", mock.AnythingOfType("*http.Request")).Return(testUser, nil)

	// Create middleware
	authMiddleware := middleware.AuthMiddleware(mockSessionService)

	// Create test handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if user is in context
		user := r.Context().Value("user").(*models.User)
		userID := r.Context().Value("user_id").(string)
		userEmail := r.Context().Value("user_email").(string)

		assert.Equal(t, testUser.ID, user.ID)
		assert.Equal(t, testUser.Email, user.Email)
		assert.Equal(t, testUser.ID, userID)
		assert.Equal(t, testUser.Email, userEmail)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	// Create request
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Apply middleware and serve
	middlewareHandler := authMiddleware(handler)
	middlewareHandler.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "success", w.Body.String())
	mockSessionService.AssertExpectations(t)
}

func TestAuthMiddleware_UnauthenticatedUser(t *testing.T) {
	// Create mock session service
	mockSessionService := new(MockSessionService)

	// Setup mock expectations - return error for unauthenticated user
	mockSessionService.On("GetSession", mock.AnythingOfType("*http.Request")).Return(nil, assert.AnError)

	// Create middleware
	authMiddleware := middleware.AuthMiddleware(mockSessionService)

	// Create test handler (should not be called)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called for unauthenticated user")
	})

	// Create request
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Apply middleware and serve
	middlewareHandler := authMiddleware(handler)
	middlewareHandler.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Check response body contains error
	responseBody := w.Body.String()
	assert.Contains(t, responseBody, "UNAUTHORIZED")
	assert.Contains(t, responseBody, "Authentication required")

	mockSessionService.AssertExpectations(t)
}

func TestOptionalAuthMiddleware_AuthenticatedUser(t *testing.T) {
	// Create mock session service
	mockSessionService := new(MockSessionService)

	// Create test user
	testUser := &models.User{
		ID:    "user-123",
		Email: "test@example.com",
		Name:  "Test User",
	}

	// Setup mock expectations
	mockSessionService.On("GetSession", mock.AnythingOfType("*http.Request")).Return(testUser, nil)

	// Create middleware
	optionalAuthMiddleware := middleware.OptionalAuthMiddleware(mockSessionService)

	// Create test handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if user is in context
		user := r.Context().Value("user").(*models.User)
		userID := r.Context().Value("user_id").(string)
		userEmail := r.Context().Value("user_email").(string)
		authenticated := r.Context().Value("authenticated").(bool)

		assert.Equal(t, testUser.ID, user.ID)
		assert.Equal(t, testUser.Email, user.Email)
		assert.Equal(t, testUser.ID, userID)
		assert.Equal(t, testUser.Email, userEmail)
		assert.True(t, authenticated)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	// Create request
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Apply middleware and serve
	middlewareHandler := optionalAuthMiddleware(handler)
	middlewareHandler.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "success", w.Body.String())
	mockSessionService.AssertExpectations(t)
}

func TestOptionalAuthMiddleware_UnauthenticatedUser(t *testing.T) {
	// Create mock session service
	mockSessionService := new(MockSessionService)

	// Setup mock expectations - return error for unauthenticated user
	mockSessionService.On("GetSession", mock.AnythingOfType("*http.Request")).Return(nil, assert.AnError)

	// Create middleware
	optionalAuthMiddleware := middleware.OptionalAuthMiddleware(mockSessionService)

	// Create test handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if authenticated is false in context
		authenticated := r.Context().Value("authenticated").(bool)
		assert.False(t, authenticated)

		// User should not be in context
		user := r.Context().Value("user")
		assert.Nil(t, user)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	// Create request
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Apply middleware and serve
	middlewareHandler := optionalAuthMiddleware(handler)
	middlewareHandler.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "success", w.Body.String())
	mockSessionService.AssertExpectations(t)
}

func TestAdminMiddleware_AdminUser(t *testing.T) {
	// Create mock session service
	mockSessionService := new(MockSessionService)

	// Create admin user
	adminUser := &models.User{
		ID:    "admin-123",
		Email: "admin@sparkparkcricket.com",
		Name:  "Admin User",
	}

	// Setup mock expectations
	mockSessionService.On("GetSession", mock.AnythingOfType("*http.Request")).Return(adminUser, nil)

	// Create middleware
	adminMiddleware := middleware.AdminMiddleware(mockSessionService)

	// Create test handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if user is in context
		user := r.Context().Value("user").(*models.User)
		userID := r.Context().Value("user_id").(string)
		userEmail := r.Context().Value("user_email").(string)
		isAdmin := r.Context().Value("is_admin").(bool)

		assert.Equal(t, adminUser.ID, user.ID)
		assert.Equal(t, adminUser.Email, user.Email)
		assert.Equal(t, adminUser.ID, userID)
		assert.Equal(t, adminUser.Email, userEmail)
		assert.True(t, isAdmin)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("admin success"))
	})

	// Create request
	req := httptest.NewRequest("GET", "/admin", nil)
	w := httptest.NewRecorder()

	// Apply middleware and serve
	middlewareHandler := adminMiddleware(handler)
	middlewareHandler.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "admin success", w.Body.String())
	mockSessionService.AssertExpectations(t)
}

func TestAdminMiddleware_NonAdminUser(t *testing.T) {
	// Create mock session service
	mockSessionService := new(MockSessionService)

	// Create non-admin user
	regularUser := &models.User{
		ID:    "user-123",
		Email: "user@example.com",
		Name:  "Regular User",
	}

	// Setup mock expectations
	mockSessionService.On("GetSession", mock.AnythingOfType("*http.Request")).Return(regularUser, nil)

	// Create middleware
	adminMiddleware := middleware.AdminMiddleware(mockSessionService)

	// Create test handler (should not be called)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called for non-admin user")
	})

	// Create request
	req := httptest.NewRequest("GET", "/admin", nil)
	w := httptest.NewRecorder()

	// Apply middleware and serve
	middlewareHandler := adminMiddleware(handler)
	middlewareHandler.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusForbidden, w.Code)

	// Check response body contains error
	responseBody := w.Body.String()
	assert.Contains(t, responseBody, "FORBIDDEN")
	assert.Contains(t, responseBody, "Admin access required")

	mockSessionService.AssertExpectations(t)
}

func TestAdminMiddleware_UnauthenticatedUser(t *testing.T) {
	// Create mock session service
	mockSessionService := new(MockSessionService)

	// Setup mock expectations - return error for unauthenticated user
	mockSessionService.On("GetSession", mock.AnythingOfType("*http.Request")).Return(nil, assert.AnError)

	// Create middleware
	adminMiddleware := middleware.AdminMiddleware(mockSessionService)

	// Create test handler (should not be called)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called for unauthenticated user")
	})

	// Create request
	req := httptest.NewRequest("GET", "/admin", nil)
	w := httptest.NewRecorder()

	// Apply middleware and serve
	middlewareHandler := adminMiddleware(handler)
	middlewareHandler.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Check response body contains error
	responseBody := w.Body.String()
	assert.Contains(t, responseBody, "UNAUTHORIZED")
	assert.Contains(t, responseBody, "Authentication required")

	mockSessionService.AssertExpectations(t)
}

func TestAdminMiddleware_AlternativeAdminEmail(t *testing.T) {
	// Create mock session service
	mockSessionService := new(MockSessionService)

	// Create admin user with alternative admin email
	adminUser := &models.User{
		ID:    "admin-456",
		Email: "luffybhaagi@gmail.com",
		Name:  "Alternative Admin",
	}

	// Setup mock expectations
	mockSessionService.On("GetSession", mock.AnythingOfType("*http.Request")).Return(adminUser, nil)

	// Create middleware
	adminMiddleware := middleware.AdminMiddleware(mockSessionService)

	// Create test handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if user is in context
		user := r.Context().Value("user").(*models.User)
		isAdmin := r.Context().Value("is_admin").(bool)

		assert.Equal(t, adminUser.ID, user.ID)
		assert.Equal(t, adminUser.Email, user.Email)
		assert.True(t, isAdmin)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("alternative admin success"))
	})

	// Create request
	req := httptest.NewRequest("GET", "/admin", nil)
	w := httptest.NewRecorder()

	// Apply middleware and serve
	middlewareHandler := adminMiddleware(handler)
	middlewareHandler.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "alternative admin success", w.Body.String())
	mockSessionService.AssertExpectations(t)
}
