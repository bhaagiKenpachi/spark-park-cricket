package unit

import (
	"context"
	"net/http/httptest"
	"spark-park-cricket-backend/internal/config"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/services"
	"testing"
	"time"

	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByGoogleID(ctx context.Context, googleID string) (*models.User, error) {
	args := m.Called(ctx, googleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, id string, updates *models.UpdateUserRequest) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockUserRepository) ListUsers(ctx context.Context, filters *models.UserFilters) ([]*models.User, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *MockUserRepository) DeleteUser(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) CreateUserSession(ctx context.Context, session *models.UserSession) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserSession(ctx context.Context, sessionID string) (*models.UserSession, error) {
	args := m.Called(ctx, sessionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserSession), args.Error(1)
}

func (m *MockUserRepository) DeleteUserSession(ctx context.Context, sessionID string) error {
	args := m.Called(ctx, sessionID)
	return args.Error(0)
}

func (m *MockUserRepository) DeleteExpiredSessions(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateLastLogin(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func TestSessionService_CreateSession(t *testing.T) {
	// Setup
	cfg := &config.Config{
		SessionSecret: "test-secret-key-32-characters-long",
		SessionMaxAge: 3600,
	}

	mockUserRepo := new(MockUserRepository)
	sessionService := services.NewSessionService(mockUserRepo, cfg)

	// Test data
	testUser := &models.User{
		ID:    "user-123",
		Email: "test@example.com",
		Name:  "Test User",
	}

	// Mock expectations
	mockUserRepo.On("CreateUserSession", mock.Anything, mock.MatchedBy(func(session *models.UserSession) bool {
		return session.UserID == testUser.ID && session.SessionID != ""
	})).Return(nil)
	mockUserRepo.On("UpdateLastLogin", mock.Anything, testUser.ID).Return(nil)

	// Create request and response recorder
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Execute
	err := sessionService.CreateSession(w, req, testUser)

	// Assertions
	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)

	// Verify session cookie was set
	cookies := w.Result().Cookies()
	assert.Len(t, cookies, 1)
	assert.Equal(t, "user_session", cookies[0].Name)
}

func TestSessionService_GetSession_Success(t *testing.T) {
	// Setup
	cfg := &config.Config{
		SessionSecret: "test-secret-key-32-characters-long",
		SessionMaxAge: 3600,
	}

	mockUserRepo := new(MockUserRepository)
	sessionService := services.NewSessionService(mockUserRepo, cfg)

	// Test data
	testUser := &models.User{
		ID:    "user-123",
		Email: "test@example.com",
		Name:  "Test User",
	}

	testSession := &models.UserSession{
		UserID:    testUser.ID,
		SessionID: "session-123",
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	// Create a session first
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Mock expectations for session creation
	mockUserRepo.On("CreateUserSession", mock.Anything, mock.Anything).Return(nil)
	mockUserRepo.On("UpdateLastLogin", mock.Anything, testUser.ID).Return(nil)

	err := sessionService.CreateSession(w, req, testUser)
	assert.NoError(t, err)

	// Extract session ID from the created session
	session, _ := sessionService.GetStore().Get(req, "user_session")
	sessionID := session.Values["session_id"].(string)

	// Mock expectations for session retrieval
	mockUserRepo.On("GetUserSession", mock.Anything, sessionID).Return(testSession, nil)
	mockUserRepo.On("GetUserByID", mock.Anything, testUser.ID).Return(testUser, nil)

	// Create new request with the session cookie
	req2 := httptest.NewRequest("GET", "/test", nil)
	req2.Header.Set("Cookie", w.Header().Get("Set-Cookie"))

	// Execute
	user, err := sessionService.GetSession(req2)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, testUser.ID, user.ID)
	assert.Equal(t, testUser.Email, user.Email)
	assert.Equal(t, testUser.Name, user.Name)
	mockUserRepo.AssertExpectations(t)
}

func TestSessionService_GetSession_NoSession(t *testing.T) {
	// Setup
	cfg := &config.Config{
		SessionSecret: "test-secret-key-32-characters-long",
		SessionMaxAge: 3600,
	}

	mockUserRepo := new(MockUserRepository)
	sessionService := services.NewSessionService(mockUserRepo, cfg)

	// Create request without session
	req := httptest.NewRequest("GET", "/test", nil)

	// Execute
	user, err := sessionService.GetSession(req)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "user not authenticated")
}

func TestSessionService_GetSession_ExpiredSession(t *testing.T) {
	// Setup
	cfg := &config.Config{
		SessionSecret: "test-secret-key-32-characters-long",
		SessionMaxAge: 3600,
	}

	mockUserRepo := new(MockUserRepository)
	sessionService := services.NewSessionService(mockUserRepo, cfg)

	// Test data
	testUser := &models.User{
		ID:    "user-123",
		Email: "test@example.com",
		Name:  "Test User",
	}

	expiredSession := &models.UserSession{
		UserID:    testUser.ID,
		SessionID: "session-123",
		ExpiresAt: time.Now().Add(-1 * time.Hour), // Expired
	}

	// Create a session first
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Mock expectations for session creation
	mockUserRepo.On("CreateUserSession", mock.Anything, mock.Anything).Return(nil)
	mockUserRepo.On("UpdateLastLogin", mock.Anything, testUser.ID).Return(nil)

	err := sessionService.CreateSession(w, req, testUser)
	assert.NoError(t, err)

	// Extract session ID from the created session
	session, _ := sessionService.GetStore().Get(req, "user_session")
	sessionID := session.Values["session_id"].(string)

	// Mock expectations for expired session
	mockUserRepo.On("GetUserSession", mock.Anything, sessionID).Return(expiredSession, nil)
	mockUserRepo.On("DeleteUserSession", mock.Anything, sessionID).Return(nil)

	// Create new request with the session cookie
	req2 := httptest.NewRequest("GET", "/test", nil)
	req2.Header.Set("Cookie", w.Header().Get("Set-Cookie"))

	// Execute
	user, err := sessionService.GetSession(req2)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "session expired")
	mockUserRepo.AssertExpectations(t)
}

func TestSessionService_IsAuthenticated_Authenticated(t *testing.T) {
	// Setup
	cfg := &config.Config{
		SessionSecret: "test-secret-key-32-characters-long",
		SessionMaxAge: 3600,
	}

	mockUserRepo := new(MockUserRepository)
	sessionService := services.NewSessionService(mockUserRepo, cfg)

	// Test data
	testUser := &models.User{
		ID:    "user-123",
		Email: "test@example.com",
		Name:  "Test User",
	}

	testSession := &models.UserSession{
		UserID:    testUser.ID,
		SessionID: "session-123",
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	// Create a session first
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Mock expectations for session creation
	mockUserRepo.On("CreateUserSession", mock.Anything, mock.Anything).Return(nil)
	mockUserRepo.On("UpdateLastLogin", mock.Anything, testUser.ID).Return(nil)

	err := sessionService.CreateSession(w, req, testUser)
	assert.NoError(t, err)

	// Extract session ID from the created session
	session, _ := sessionService.GetStore().Get(req, "user_session")
	sessionID := session.Values["session_id"].(string)

	// Mock expectations for authentication check
	mockUserRepo.On("GetUserSession", mock.Anything, sessionID).Return(testSession, nil)
	mockUserRepo.On("GetUserByID", mock.Anything, testUser.ID).Return(testUser, nil)

	// Create new request with the session cookie
	req2 := httptest.NewRequest("GET", "/test", nil)
	req2.Header.Set("Cookie", w.Header().Get("Set-Cookie"))

	// Execute
	isAuth := sessionService.IsAuthenticated(req2)

	// Assertions
	assert.True(t, isAuth)
	mockUserRepo.AssertExpectations(t)
}

func TestSessionService_IsAuthenticated_NotAuthenticated(t *testing.T) {
	// Setup
	cfg := &config.Config{
		SessionSecret: "test-secret-key-32-characters-long",
		SessionMaxAge: 3600,
	}

	mockUserRepo := new(MockUserRepository)
	sessionService := services.NewSessionService(mockUserRepo, cfg)

	// Create request without session
	req := httptest.NewRequest("GET", "/test", nil)

	// Execute
	isAuth := sessionService.IsAuthenticated(req)

	// Assertions
	assert.False(t, isAuth)
}

func TestSessionService_DestroySession(t *testing.T) {
	// Setup
	cfg := &config.Config{
		SessionSecret: "test-secret-key-32-characters-long",
		SessionMaxAge: 3600,
	}

	mockUserRepo := new(MockUserRepository)
	sessionService := services.NewSessionService(mockUserRepo, cfg)

	// Test data
	testUser := &models.User{
		ID:    "user-123",
		Email: "test@example.com",
		Name:  "Test User",
	}

	// Create a session first
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Mock expectations for session creation
	mockUserRepo.On("CreateUserSession", mock.Anything, mock.Anything).Return(nil)
	mockUserRepo.On("UpdateLastLogin", mock.Anything, testUser.ID).Return(nil)

	err := sessionService.CreateSession(w, req, testUser)
	assert.NoError(t, err)

	// Extract session ID from the created session
	session, _ := sessionService.GetStore().Get(req, "user_session")
	sessionID := session.Values["session_id"].(string)

	// Mock expectations for session destruction
	mockUserRepo.On("DeleteUserSession", mock.Anything, sessionID).Return(nil)

	// Create new request with the session cookie
	req2 := httptest.NewRequest("GET", "/test", nil)
	req2.Header.Set("Cookie", w.Header().Get("Set-Cookie"))
	w2 := httptest.NewRecorder()

	// Execute
	err = sessionService.DestroySession(w2, req2)

	// Assertions
	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)

	// Verify session cookie was cleared
	cookies := w2.Result().Cookies()
	assert.Len(t, cookies, 1)
	assert.Equal(t, "user_session", cookies[0].Name)
	assert.Equal(t, -1, cookies[0].MaxAge)
}

func TestSessionService_CleanupExpiredSessions(t *testing.T) {
	// Setup
	cfg := &config.Config{
		SessionSecret: "test-secret-key-32-characters-long",
		SessionMaxAge: 3600,
	}

	mockUserRepo := new(MockUserRepository)
	sessionService := services.NewSessionService(mockUserRepo, cfg)

	// Mock expectations
	mockUserRepo.On("DeleteExpiredSessions", mock.Anything).Return(nil)

	// Execute
	err := sessionService.CleanupExpiredSessions(context.Background())

	// Assertions
	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
}

func TestSessionService_GetStore(t *testing.T) {
	// Setup
	cfg := &config.Config{
		SessionSecret: "test-secret-key-32-characters-long",
		SessionMaxAge: 3600,
	}

	mockUserRepo := new(MockUserRepository)
	sessionService := services.NewSessionService(mockUserRepo, cfg)

	// Execute
	store := sessionService.GetStore()

	// Assertions
	assert.NotNil(t, store)
	assert.IsType(t, &sessions.CookieStore{}, store)
}

// Test to verify that session values can be logged without JSON marshaling errors
func TestSessionService_LoggingSessionValues(t *testing.T) {
	// Setup
	cfg := &config.Config{
		SessionSecret: "test-secret-key-32-characters-long",
		SessionMaxAge: 3600,
	}

	mockUserRepo := new(MockUserRepository)
	sessionService := services.NewSessionService(mockUserRepo, cfg)

	// Test data
	testUser := &models.User{
		ID:    "user-123",
		Email: "test@example.com",
		Name:  "Test User",
	}

	// Create a session first
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Mock expectations for session creation
	mockUserRepo.On("CreateUserSession", mock.Anything, mock.Anything).Return(nil)
	mockUserRepo.On("UpdateLastLogin", mock.Anything, testUser.ID).Return(nil)

	// This should not cause JSON marshaling errors due to our sessionValuesToMap helper
	err := sessionService.CreateSession(w, req, testUser)
	assert.NoError(t, err)

	// Verify session was created successfully
	session, _ := sessionService.GetStore().Get(req, "user_session")
	assert.NotNil(t, session)
	assert.True(t, session.Values["authenticated"].(bool))

	mockUserRepo.AssertExpectations(t)
}
