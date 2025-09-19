package services

import (
	"context"
	"net/http"
	"spark-park-cricket-backend/internal/models"

	"github.com/gorilla/sessions"
)

// SessionServiceInterface defines the interface for session management
type SessionServiceInterface interface {
	// CreateSession creates a new session for a user
	CreateSession(w http.ResponseWriter, r *http.Request, user *models.User) error

	// GetSession retrieves the current session
	GetSession(r *http.Request) (*models.User, error)

	// DestroySession destroys the current session
	DestroySession(w http.ResponseWriter, r *http.Request) error

	// IsAuthenticated checks if the user is authenticated
	IsAuthenticated(r *http.Request) bool

	// CleanupExpiredSessions removes expired sessions from the database
	CleanupExpiredSessions(ctx context.Context) error

	// GetStore returns the underlying session store (for OAuth state management)
	GetStore() *sessions.CookieStore
}
