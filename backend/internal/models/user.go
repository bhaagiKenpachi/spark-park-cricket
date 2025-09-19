package models

import (
	"time"
)

// User represents a system user
type User struct {
	ID            string    `json:"id" db:"id"`
	GoogleID      string    `json:"google_id" db:"google_id"`
	Email         string    `json:"email" db:"email"`
	Name          string    `json:"name" db:"name"`
	Picture       string    `json:"picture" db:"picture"`
	EmailVerified bool      `json:"email_verified" db:"email_verified"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
	LastLoginAt   time.Time `json:"last_login_at" db:"last_login_at"`
}

// UserSession represents a user session
type UserSession struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	SessionID string    `json:"session_id" db:"session_id"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// GoogleUserInfo represents the user info from Google OAuth
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	EmailVerified bool   `json:"verified_email"`
}

// CreateUserRequest represents the request to create a new user
type CreateUserRequest struct {
	GoogleID      string `json:"google_id" validate:"required"`
	Email         string `json:"email" validate:"required,email"`
	Name          string `json:"name" validate:"required,min=2,max=255"`
	Picture       string `json:"picture"`
	EmailVerified bool   `json:"email_verified"`
}

// UpdateUserRequest represents the request to update a user
type UpdateUserRequest struct {
	Name        *string    `json:"name,omitempty" validate:"omitempty,min=2,max=255"`
	Picture     *string    `json:"picture,omitempty"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
}

// UserFilters represents filters for listing users
type UserFilters struct {
	Email  *string `json:"email,omitempty"`
	Limit  int     `json:"limit" validate:"min=1,max=100"`
	Offset int     `json:"offset" validate:"min=0"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	User    *User  `json:"user"`
	Token   string `json:"token,omitempty"`
	Message string `json:"message"`
}
