package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/services"
	"spark-park-cricket-backend/internal/utils"

	"github.com/go-chi/chi/v5"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	AuthService services.AuthServiceInterface
	SessionSvc  services.SessionServiceInterface
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(authService services.AuthServiceInterface, sessionSvc services.SessionServiceInterface) *AuthHandler {
	return &AuthHandler{
		AuthService: authService,
		SessionSvc:  sessionSvc,
	}
}

// GoogleLogin initiates Google OAuth login
func (h *AuthHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	// Generate a random state parameter for security
	state, err := h.generateState()
	if err != nil {
		utils.LogError(err, "Failed to generate state parameter", nil)
		utils.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to initiate login", nil)
		return
	}

	// Store state in session for validation
	session, err := h.SessionSvc.GetStore().Get(r, "oauth_state")
	if err != nil {
		utils.LogError(err, "Failed to get OAuth state session", nil)
		utils.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to initiate login", nil)
		return
	}

	session.Values["state"] = state
	session.Options.MaxAge = 600 // 10 minutes

	if err := session.Save(r, w); err != nil {
		utils.LogError(err, "Failed to save OAuth state session", nil)
		utils.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to initiate login", nil)
		return
	}

	// For debugging: also store state in a simple cookie as backup
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state_backup",
		Value:    state,
		Path:     "/",
		MaxAge:   600, // 10 minutes
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	// Get Google OAuth URL and redirect
	authURL := h.AuthService.GetAuthURL(state)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// GoogleCallback handles Google OAuth callback
func (h *AuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	errorParam := r.URL.Query().Get("error")

	// Check for OAuth errors
	if errorParam != "" {
		utils.LogWarn("OAuth error received", map[string]interface{}{
			"error": errorParam,
		})
		utils.WriteError(w, http.StatusBadRequest, "OAUTH_ERROR", "Authentication failed", nil)
		return
	}

	// Validate state parameter
	session, err := h.SessionSvc.GetStore().Get(r, "oauth_state")
	if err != nil {
		utils.LogError(err, "Failed to get OAuth state session", nil)
		utils.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Invalid session", nil)
		return
	}

	storedState, ok := session.Values["state"].(string)

	// Check backup cookie if session state is not available
	backupState := ""
	if !ok || storedState == "" {
		if cookie, err := r.Cookie("oauth_state_backup"); err == nil {
			backupState = cookie.Value
		}
	}

	// Use session state if available, otherwise use backup cookie
	validState := storedState
	if !ok || storedState == "" {
		validState = backupState
	}

	if validState == "" || validState != state {
		utils.LogWarn("Invalid state parameter", map[string]interface{}{
			"received_state": state,
			"stored_state":   storedState,
			"backup_state":   backupState,
		})
		utils.WriteError(w, http.StatusBadRequest, "INVALID_STATE", "Invalid state parameter", nil)
		return
	}

	// Clear the state session
	session.Values["state"] = nil
	session.Options.MaxAge = -1
	if err := session.Save(r, w); err != nil {
		utils.LogError(err, "Failed to clear OAuth state session", nil)
	}

	// Clear the backup cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state_backup",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	// Authenticate user
	user, err := h.AuthService.AuthenticateUser(r.Context(), code)
	if err != nil {
		utils.LogError(err, "Failed to authenticate user", nil)
		utils.WriteError(w, http.StatusInternalServerError, "AUTH_ERROR", "Authentication failed", nil)
		return
	}

	// Create user session
	if err := h.SessionSvc.CreateSession(w, r, user); err != nil {
		utils.LogError(err, "Failed to create user session", nil)
		utils.WriteError(w, http.StatusInternalServerError, "SESSION_ERROR", "Failed to create session", nil)
		return
	}

	utils.LogInfo("User session created successfully, redirecting to frontend", map[string]interface{}{
		"user_id": user.ID,
		"email":   user.Email,
		"cookies": w.Header().Get("Set-Cookie"),
	})

	// Redirect to frontend after successful authentication
	http.Redirect(w, r, "http://localhost:3000?auth=success", http.StatusTemporaryRedirect)
}

// Logout handles user logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Destroy session
	if err := h.SessionSvc.DestroySession(w, r); err != nil {
		utils.LogError(err, "Failed to destroy session", nil)
		utils.WriteError(w, http.StatusInternalServerError, "LOGOUT_ERROR", "Failed to logout", nil)
		return
	}

	utils.LogInfo("User logged out successfully", nil)

	// Return success response
	response := models.AuthResponse{
		Message: "Logged out successfully",
	}
	utils.WriteSuccess(w, response)
}

// GetCurrentUser returns the current authenticated user
func (h *AuthHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	user, err := h.SessionSvc.GetSession(r)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Not authenticated", nil)
		return
	}

	response := models.AuthResponse{
		User: user,
	}
	utils.WriteSuccess(w, response)
}

// AuthStatus returns the authentication status
func (h *AuthHandler) AuthStatus(w http.ResponseWriter, r *http.Request) {
	utils.LogInfo("=== AUTH STATUS CHECK ===", map[string]interface{}{
		"cookies": r.Header.Get("Cookie"),
		"method":  r.Method,
		"url":     r.URL.String(),
	})

	isAuthenticated := h.SessionSvc.IsAuthenticated(r)

	utils.LogInfo("Authentication check result", map[string]interface{}{
		"isAuthenticated": isAuthenticated,
	})

	response := map[string]interface{}{
		"authenticated": isAuthenticated,
	}

	if isAuthenticated {
		user, err := h.SessionSvc.GetSession(r)
		if err == nil {
			response["user"] = user
			utils.LogInfo("User session found", map[string]interface{}{
				"user_id": user.ID,
			})
		} else {
			utils.LogError(err, "Failed to get user session", nil)
		}
	}

	utils.WriteSuccess(w, response)
}

// generateState generates a cryptographically secure random state parameter
func (h *AuthHandler) generateState() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// SetupAuthRoutes sets up authentication routes
func SetupAuthRoutes(r chi.Router, authHandler *AuthHandler) {
	r.Route("/auth", func(r chi.Router) {
		r.Get("/google", authHandler.GoogleLogin)
		r.Get("/google/callback", authHandler.GoogleCallback)
		r.Post("/logout", authHandler.Logout)
		r.Get("/me", authHandler.GetCurrentUser)
		r.Get("/status", authHandler.AuthStatus)
	})
}
