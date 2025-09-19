package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"spark-park-cricket-backend/internal/config"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/repository/interfaces"
	"spark-park-cricket-backend/internal/utils"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// AuthService handles authentication operations
type AuthService struct {
	config       *config.Config
	userRepo     interfaces.UserRepository
	sessionSvc   *SessionService
	oauth2Config *oauth2.Config
}

// NewAuthService creates a new authentication service
func NewAuthService(cfg *config.Config, userRepo interfaces.UserRepository, sessionSvc *SessionService) *AuthService {
	oauth2Config := &oauth2.Config{
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		RedirectURL:  cfg.GoogleRedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &AuthService{
		config:       cfg,
		userRepo:     userRepo,
		sessionSvc:   sessionSvc,
		oauth2Config: oauth2Config,
	}
}

// GetAuthURL returns the Google OAuth authorization URL
func (s *AuthService) GetAuthURL(state string) string {
	return s.oauth2Config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// ExchangeCodeForToken exchanges authorization code for access token
func (s *AuthService) ExchangeCodeForToken(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := s.oauth2Config.Exchange(ctx, code)
	if err != nil {
		utils.LogError(err, "Failed to exchange code for token", nil)
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	return token, nil
}

// GetUserInfoFromGoogle fetches user information from Google
func (s *AuthService) GetUserInfoFromGoogle(ctx context.Context, token *oauth2.Token) (*models.GoogleUserInfo, error) {
	client := s.oauth2Config.Client(ctx, token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info from Google: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Google API returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var userInfo models.GoogleUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user info: %w", err)
	}

	return &userInfo, nil
}

// CreateOrUpdateUser creates a new user or updates existing one
func (s *AuthService) CreateOrUpdateUser(ctx context.Context, googleUserInfo *models.GoogleUserInfo) (*models.User, error) {
	// Try to get existing user by Google ID
	user, err := s.userRepo.GetUserByGoogleID(ctx, googleUserInfo.ID)
	if err != nil {
		// User doesn't exist, create new one
		user = &models.User{
			GoogleID:      googleUserInfo.ID,
			Email:         googleUserInfo.Email,
			Name:          googleUserInfo.Name,
			Picture:       googleUserInfo.Picture,
			EmailVerified: googleUserInfo.EmailVerified,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			LastLoginAt:   time.Now(),
		}

		if err := s.userRepo.CreateUser(ctx, user); err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
	} else {
		// User exists, update last login and potentially other info
		updates := &models.UpdateUserRequest{
			LastLoginAt: &time.Time{},
		}
		*updates.LastLoginAt = time.Now()

		// Update picture if it has changed
		if user.Picture != googleUserInfo.Picture {
			updates.Picture = &googleUserInfo.Picture
		}

		// Update name if it has changed
		if user.Name != googleUserInfo.Name {
			updates.Name = &googleUserInfo.Name
		}

		if err := s.userRepo.UpdateUser(ctx, user.ID, updates); err != nil {
			return nil, fmt.Errorf("failed to update user: %w", err)
		}

		// Refresh user data
		user, err = s.userRepo.GetUserByID(ctx, user.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get updated user: %w", err)
		}
	}

	return user, nil
}

// AuthenticateUser performs the complete authentication flow
func (s *AuthService) AuthenticateUser(ctx context.Context, code string) (*models.User, error) {
	// Exchange code for token
	token, err := s.ExchangeCodeForToken(ctx, code)
	if err != nil {
		utils.LogError(err, "Failed to exchange code for token", nil)
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	// Get user info from Google
	googleUserInfo, err := s.GetUserInfoFromGoogle(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info from Google: %w", err)
	}

	// Create or update user
	user, err := s.CreateOrUpdateUser(ctx, googleUserInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to create or update user: %w", err)
	}

	return user, nil
}
