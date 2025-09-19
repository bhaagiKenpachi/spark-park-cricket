package models

import (
	"time"
)

// OAuthState represents an OAuth state parameter stored in the database
type OAuthState struct {
	ID        string     `json:"id" db:"id"`
	State     string     `json:"state" db:"state"`
	ExpiresAt time.Time  `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UsedAt    *time.Time `json:"used_at" db:"used_at"`
}

// CreateOAuthStateRequest represents the request to create an OAuth state
type CreateOAuthStateRequest struct {
	State     string    `json:"state"`
	ExpiresAt time.Time `json:"expires_at"`
}
