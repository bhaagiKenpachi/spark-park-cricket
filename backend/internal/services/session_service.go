package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"spark-park-cricket-backend/internal/config"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/repository/interfaces"
	"spark-park-cricket-backend/internal/utils"
	"time"

	"github.com/gorilla/sessions"
)

// SessionService handles session management
type SessionService struct {
	Store    *sessions.CookieStore
	userRepo interfaces.UserRepository
	config   *config.Config
}

// NewSessionService creates a new session service
func NewSessionService(userRepo interfaces.UserRepository, cfg *config.Config) *SessionService {
	// Create secure cookie store
	store := sessions.NewCookieStore([]byte(cfg.SessionSecret))

	// Configure session options for OAuth compatibility
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   cfg.SessionMaxAge,
		HttpOnly: true,
		Secure:   false,                // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode, // Changed from None to Lax for localhost
		// Don't set Domain to avoid localhost issues
	}

	return &SessionService{
		Store:    store,
		userRepo: userRepo,
		config:   cfg,
	}
}

// CreateSession creates a new session for a user
func (s *SessionService) CreateSession(w http.ResponseWriter, r *http.Request, user *models.User) error {
	utils.LogInfo("Creating user session", map[string]interface{}{
		"user_id": user.ID,
		"email":   user.Email,
		"cookies": r.Header.Get("Cookie"),
	})

	session, err := s.Store.Get(r, "user_session")
	if err != nil {
		utils.LogError(err, "Failed to get user session", nil)
		return fmt.Errorf("failed to get session: %w", err)
	}

	utils.LogInfo("Session retrieved for creation", map[string]interface{}{
		"session_id": session.ID,
		"is_new":     session.IsNew,
		"name":       session.Name,
		"options":    session.Options,
	})

	// Generate a unique session ID
	sessionID, err := s.generateSessionID()
	if err != nil {
		return fmt.Errorf("failed to generate session ID: %w", err)
	}

	// Store user data in session
	session.Values["user_id"] = user.ID
	session.Values["session_id"] = sessionID
	session.Values["email"] = user.Email
	session.Values["name"] = user.Name
	session.Values["authenticated"] = true

	utils.LogInfo("Session values set", map[string]interface{}{
		"user_id":        user.ID,
		"session_id":     sessionID,
		"email":          user.Email,
		"name":           user.Name,
		"authenticated":  true,
		"session_values": session.Values,
	})

	// Save session to database
	userSession := &models.UserSession{
		UserID:    user.ID,
		SessionID: sessionID,
		ExpiresAt: time.Now().Add(time.Duration(s.config.SessionMaxAge) * time.Second),
	}

	if err := s.userRepo.CreateUserSession(r.Context(), userSession); err != nil {
		return fmt.Errorf("failed to create user session: %w", err)
	}

	// Update last login time
	if err := s.userRepo.UpdateLastLogin(r.Context(), user.ID); err != nil {
		// Log error but don't fail the session creation
		fmt.Printf("Warning: failed to update last login: %v\n", err)
	}

	// Save session to cookie
	if err := session.Save(r, w); err != nil {
		utils.LogError(err, "Failed to save session", nil)
		return fmt.Errorf("failed to save session: %w", err)
	}

	utils.LogInfo("Session saved successfully", map[string]interface{}{
		"session_id":  sessionID,
		"user_id":     user.ID,
		"cookies":     w.Header().Get("Set-Cookie"),
		"all_headers": w.Header(),
	})

	return nil
}

// GetSession retrieves the current session
func (s *SessionService) GetSession(r *http.Request) (*models.User, error) {
	utils.LogInfo("=== GET SESSION ===", map[string]interface{}{
		"cookies": r.Header.Get("Cookie"),
	})

	session, err := s.Store.Get(r, "user_session")
	if err != nil {
		utils.LogError(err, "Failed to get session", nil)
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	utils.LogInfo("Session retrieved", map[string]interface{}{
		"session_id": session.ID,
		"is_new":     session.IsNew,
		"values":     session.Values,
	})

	// Check if user is authenticated
	authenticated, ok := session.Values["authenticated"].(bool)
	if !ok || !authenticated {
		return nil, fmt.Errorf("user not authenticated")
	}

	// Get user ID from session
	userID, ok := session.Values["user_id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid session: missing user ID")
	}

	// Get session ID from session
	sessionID, ok := session.Values["session_id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid session: missing session ID")
	}

	// Verify session exists in database and is not expired
	userSession, err := s.userRepo.GetUserSession(r.Context(), sessionID)
	if err != nil {
		return nil, fmt.Errorf("invalid session: %w", err)
	}

	// Check if session is expired
	if time.Now().After(userSession.ExpiresAt) {
		// Clean up expired session
		s.userRepo.DeleteUserSession(r.Context(), sessionID)
		return nil, fmt.Errorf("session expired")
	}

	// Get user from database
	user, err := s.userRepo.GetUserByID(r.Context(), userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// DestroySession destroys the current session
func (s *SessionService) DestroySession(w http.ResponseWriter, r *http.Request) error {
	session, err := s.Store.Get(r, "user_session")
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}

	// Get session ID from session
	if sessionID, ok := session.Values["session_id"].(string); ok {
		// Delete session from database
		if err := s.userRepo.DeleteUserSession(r.Context(), sessionID); err != nil {
			// Log error but don't fail the logout
			fmt.Printf("Warning: failed to delete user session: %v\n", err)
		}
	}

	// Clear session values
	session.Values["user_id"] = nil
	session.Values["session_id"] = nil
	session.Values["email"] = nil
	session.Values["name"] = nil
	session.Values["authenticated"] = false

	// Set max age to -1 to delete the cookie
	session.Options.MaxAge = -1

	// Save session
	if err := session.Save(r, w); err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}

	return nil
}

// IsAuthenticated checks if the user is authenticated
func (s *SessionService) IsAuthenticated(r *http.Request) bool {
	utils.LogInfo("=== SESSION AUTHENTICATION CHECK ===", map[string]interface{}{
		"cookies": r.Header.Get("Cookie"),
	})

	_, err := s.GetSession(r)
	isAuth := err == nil

	utils.LogInfo("Authentication check result", map[string]interface{}{
		"isAuthenticated": isAuth,
		"error":           err,
	})

	return isAuth
}

// generateSessionID generates a cryptographically secure random session ID
func (s *SessionService) generateSessionID() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// CleanupExpiredSessions removes expired sessions from the database
func (s *SessionService) CleanupExpiredSessions(ctx context.Context) error {
	return s.userRepo.DeleteExpiredSessions(ctx)
}

// GetStore returns the underlying session store
func (s *SessionService) GetStore() *sessions.CookieStore {
	return s.Store
}
