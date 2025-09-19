package services

import (
	"context"
	"spark-park-cricket-backend/internal/models"

	"golang.org/x/oauth2"
)

// AuthServiceInterface defines the interface for authentication operations
type AuthServiceInterface interface {
	// GetAuthURL returns the Google OAuth authorization URL
	GetAuthURL(state string) string

	// ExchangeCodeForToken exchanges authorization code for access token
	ExchangeCodeForToken(ctx context.Context, code string) (*oauth2.Token, error)

	// GetUserInfoFromGoogle fetches user information from Google
	GetUserInfoFromGoogle(ctx context.Context, token *oauth2.Token) (*models.GoogleUserInfo, error)

	// CreateOrUpdateUser creates a new user or updates existing one
	CreateOrUpdateUser(ctx context.Context, googleUserInfo *models.GoogleUserInfo) (*models.User, error)

	// AuthenticateUser performs the complete authentication flow
	AuthenticateUser(ctx context.Context, code string) (*models.User, error)
}
