package supabase

import (
	"context"
	"fmt"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/repository/interfaces"
	"time"

	"github.com/google/uuid"
	"github.com/supabase-community/supabase-go"
)

type userRepository struct {
	client *supabase.Client
}

// NewUserRepository creates a new user repository
func NewUserRepository(client *supabase.Client) interfaces.UserRepository {
	return &userRepository{
		client: client,
	}
}

// CreateUser creates a new user
func (r *userRepository) CreateUser(ctx context.Context, user *models.User) error {
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	var result []models.User
	_, err := r.client.From("users").Insert(user, false, "", "", "").ExecuteTo(&result)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	if len(result) > 0 {
		*user = result[0]
	}

	return nil
}

// GetUserByID retrieves a user by ID
func (r *userRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	var users []models.User
	_, err := r.client.From("users").Select("*", "", false).Eq("id", id).ExecuteTo(&users)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return &users[0], nil
}

// GetUserByGoogleID retrieves a user by Google ID
func (r *userRepository) GetUserByGoogleID(ctx context.Context, googleID string) (*models.User, error) {
	var users []models.User
	_, err := r.client.From("users").Select("*", "", false).Eq("google_id", googleID).ExecuteTo(&users)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by Google ID: %w", err)
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return &users[0], nil
}

// GetUserByEmail retrieves a user by email
func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var users []models.User
	_, err := r.client.From("users").Select("*", "", false).Eq("email", email).ExecuteTo(&users)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return &users[0], nil
}

// UpdateUser updates an existing user
func (r *userRepository) UpdateUser(ctx context.Context, id string, updates *models.UpdateUserRequest) error {
	updateData := map[string]interface{}{
		"updated_at": time.Now(),
	}

	if updates.Name != nil {
		updateData["name"] = *updates.Name
	}
	if updates.Picture != nil {
		updateData["picture"] = *updates.Picture
	}
	if updates.LastLoginAt != nil {
		updateData["last_login_at"] = *updates.LastLoginAt
	}

	var result []models.User
	_, err := r.client.From("users").Update(updateData, "", "").Eq("id", id).ExecuteTo(&result)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// DeleteUser deletes a user
func (r *userRepository) DeleteUser(ctx context.Context, id string) error {
	var result []models.User
	_, err := r.client.From("users").Delete("", "").Eq("id", id).ExecuteTo(&result)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// ListUsers retrieves a list of users with optional filters
func (r *userRepository) ListUsers(ctx context.Context, filters *models.UserFilters) ([]*models.User, error) {
	query := r.client.From("users").Select("*", "", false)

	if filters.Email != nil {
		query = query.Eq("email", *filters.Email)
	}

	if filters.Limit > 0 {
		query = query.Limit(filters.Limit, "")
	}

	if filters.Offset > 0 {
		query = query.Range(filters.Offset, filters.Offset+filters.Limit-1, "")
	}

	var users []models.User
	_, err := query.ExecuteTo(&users)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	// Convert to pointers
	result := make([]*models.User, len(users))
	for i := range users {
		result[i] = &users[i]
	}

	return result, nil
}

// CreateUserSession creates a new user session
func (r *userRepository) CreateUserSession(ctx context.Context, session *models.UserSession) error {
	session.ID = uuid.New().String()
	session.CreatedAt = time.Now()
	session.UpdatedAt = time.Now()

	var result []models.UserSession
	_, err := r.client.From("user_sessions").Insert(session, false, "", "", "").ExecuteTo(&result)
	if err != nil {
		return fmt.Errorf("failed to create user session: %w", err)
	}

	return nil
}

// GetUserSession retrieves a user session by session ID
func (r *userRepository) GetUserSession(ctx context.Context, sessionID string) (*models.UserSession, error) {
	var sessions []models.UserSession
	_, err := r.client.From("user_sessions").Select("*", "", false).Eq("session_id", sessionID).ExecuteTo(&sessions)
	if err != nil {
		return nil, fmt.Errorf("failed to get user session: %w", err)
	}

	if len(sessions) == 0 {
		return nil, fmt.Errorf("session not found")
	}

	return &sessions[0], nil
}

// DeleteUserSession deletes a user session
func (r *userRepository) DeleteUserSession(ctx context.Context, sessionID string) error {
	var result []models.UserSession
	_, err := r.client.From("user_sessions").Delete("", "").Eq("session_id", sessionID).ExecuteTo(&result)
	if err != nil {
		return fmt.Errorf("failed to delete user session: %w", err)
	}

	return nil
}

// DeleteExpiredSessions deletes expired user sessions
func (r *userRepository) DeleteExpiredSessions(ctx context.Context) error {
	now := time.Now().Format(time.RFC3339)
	var result []models.UserSession
	_, err := r.client.From("user_sessions").Delete("", "").Lt("expires_at", now).ExecuteTo(&result)
	if err != nil {
		return fmt.Errorf("failed to delete expired sessions: %w", err)
	}

	return nil
}

// UpdateLastLogin updates the last login time for a user
func (r *userRepository) UpdateLastLogin(ctx context.Context, userID string) error {
	updateData := map[string]interface{}{
		"last_login_at": time.Now(),
		"updated_at":    time.Now(),
	}

	var result []models.User
	_, err := r.client.From("users").Update(updateData, "", "").Eq("id", userID).ExecuteTo(&result)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}

	return nil
}
