package interfaces

import (
	"context"
	"spark-park-cricket-backend/internal/models"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	// CreateUser creates a new user
	CreateUser(ctx context.Context, user *models.User) error

	// GetUserByID retrieves a user by ID
	GetUserByID(ctx context.Context, id string) (*models.User, error)

	// GetUserByGoogleID retrieves a user by Google ID
	GetUserByGoogleID(ctx context.Context, googleID string) (*models.User, error)

	// GetUserByEmail retrieves a user by email
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)

	// UpdateUser updates an existing user
	UpdateUser(ctx context.Context, id string, updates *models.UpdateUserRequest) error

	// DeleteUser deletes a user
	DeleteUser(ctx context.Context, id string) error

	// ListUsers retrieves a list of users with optional filters
	ListUsers(ctx context.Context, filters *models.UserFilters) ([]*models.User, error)

	// CreateUserSession creates a new user session
	CreateUserSession(ctx context.Context, session *models.UserSession) error

	// GetUserSession retrieves a user session by session ID
	GetUserSession(ctx context.Context, sessionID string) (*models.UserSession, error)

	// DeleteUserSession deletes a user session
	DeleteUserSession(ctx context.Context, sessionID string) error

	// DeleteExpiredSessions deletes expired user sessions
	DeleteExpiredSessions(ctx context.Context) error

	// UpdateLastLogin updates the last login time for a user
	UpdateLastLogin(ctx context.Context, userID string) error
}
