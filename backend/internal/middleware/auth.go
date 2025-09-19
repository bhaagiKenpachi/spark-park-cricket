package middleware

import (
	"context"
	"net/http"
	"spark-park-cricket-backend/internal/services"
	"spark-park-cricket-backend/internal/utils"
)

// AuthMiddleware provides authentication middleware
func AuthMiddleware(sessionSvc services.SessionServiceInterface) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if user is authenticated
			user, err := sessionSvc.GetSession(r)
			if err != nil {
				utils.LogWarn("Authentication failed", map[string]interface{}{
					"error":  err.Error(),
					"path":   r.URL.Path,
					"method": r.Method,
				})

				utils.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required", nil)
				return
			}

			// Add user to request context
			ctx := r.Context()
			ctx = context.WithValue(ctx, "user", user)
			ctx = context.WithValue(ctx, "user_id", user.ID)
			ctx = context.WithValue(ctx, "user_email", user.Email)

			// Continue with authenticated request
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// OptionalAuthMiddleware provides optional authentication middleware
// This allows both authenticated and unauthenticated requests
func OptionalAuthMiddleware(sessionSvc services.SessionServiceInterface) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Try to get user session, but don't fail if not authenticated
			user, err := sessionSvc.GetSession(r)
			if err == nil && user != nil {
				// User is authenticated, add to context
				ctx := r.Context()
				ctx = context.WithValue(ctx, "user", user)
				ctx = context.WithValue(ctx, "user_id", user.ID)
				ctx = context.WithValue(ctx, "user_email", user.Email)
				ctx = context.WithValue(ctx, "authenticated", true)
				r = r.WithContext(ctx)
			} else {
				// User is not authenticated, but continue anyway
				ctx := r.Context()
				ctx = context.WithValue(ctx, "authenticated", false)
				r = r.WithContext(ctx)
			}

			next.ServeHTTP(w, r)
		})
	}
}

// AdminMiddleware provides admin-only access middleware
func AdminMiddleware(sessionSvc services.SessionServiceInterface) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// First check authentication
			user, err := sessionSvc.GetSession(r)
			if err != nil {
				utils.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required", nil)
				return
			}

			// Check if user is admin (you can implement your own admin logic here)
			// For now, we'll check if the email contains "admin" or is a specific admin email
			isAdmin := false
			// Add your admin logic here - this is just an example
			if user.Email == "admin@sparkparkcricket.com" ||
				user.Email == "luffybhaagi@gmail.com" { // Replace with your admin email
				isAdmin = true
			}

			if !isAdmin {
				utils.LogWarn("Admin access denied", map[string]interface{}{
					"user_email": user.Email,
					"path":       r.URL.Path,
					"method":     r.Method,
				})

				utils.WriteError(w, http.StatusForbidden, "FORBIDDEN", "Admin access required", nil)
				return
			}

			// Add user to request context
			ctx := r.Context()
			ctx = context.WithValue(ctx, "user", user)
			ctx = context.WithValue(ctx, "user_id", user.ID)
			ctx = context.WithValue(ctx, "user_email", user.Email)
			ctx = context.WithValue(ctx, "is_admin", true)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
