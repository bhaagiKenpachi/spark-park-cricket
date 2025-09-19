package unit

import (
	"context"
	"net/http"
	"spark-park-cricket-backend/internal/models"

	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
)

// MockAuthService is a mock implementation of AuthServiceInterface
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) GetAuthURL(state string) string {
	args := m.Called(state)
	return args.String(0)
}

func (m *MockAuthService) ExchangeCodeForToken(ctx context.Context, code string) (*oauth2.Token, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*oauth2.Token), args.Error(1)
}

func (m *MockAuthService) GetUserInfoFromGoogle(ctx context.Context, token *oauth2.Token) (*models.GoogleUserInfo, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.GoogleUserInfo), args.Error(1)
}

func (m *MockAuthService) CreateOrUpdateUser(ctx context.Context, googleUserInfo *models.GoogleUserInfo) (*models.User, error) {
	args := m.Called(ctx, googleUserInfo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockAuthService) AuthenticateUser(ctx context.Context, code string) (*models.User, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

// MockSessionService is a mock implementation of SessionServiceInterface
type MockSessionService struct {
	mock.Mock
}

func (m *MockSessionService) GetSession(r *http.Request) (*models.User, error) {
	args := m.Called(r)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockSessionService) CreateSession(w http.ResponseWriter, r *http.Request, user *models.User) error {
	args := m.Called(w, r, user)
	return args.Error(0)
}

func (m *MockSessionService) DestroySession(w http.ResponseWriter, r *http.Request) error {
	args := m.Called(w, r)
	return args.Error(0)
}

func (m *MockSessionService) IsAuthenticated(r *http.Request) bool {
	args := m.Called(r)
	return args.Bool(0)
}

func (m *MockSessionService) CleanupExpiredSessions(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockSessionService) GetStore() *sessions.CookieStore {
	args := m.Called()
	return args.Get(0).(*sessions.CookieStore)
}
