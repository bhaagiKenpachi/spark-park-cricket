package handlers

import (
	"net/http"
	"spark-park-cricket-backend/internal/config"
	"spark-park-cricket-backend/internal/database"
	"spark-park-cricket-backend/internal/middleware"
	"spark-park-cricket-backend/internal/services"
	"spark-park-cricket-backend/internal/utils"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(dbClient *database.Client, cfg *config.Config) *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RecoveryMiddleware)
	r.Use(middleware.LoggerMiddleware)
	r.Use(middleware.RequestIDMiddleware)
	r.Use(chimiddleware.RealIP)
	r.Use(middleware.TimeoutMiddleware(60 * time.Second))
	r.Use(middleware.SecurityMiddleware)
	r.Use(middleware.ValidationMiddleware)
	r.Use(middleware.MetricsMiddleware)
	r.Use(middleware.RateLimitMiddleware(100)) // 100 requests per minute
	r.Use(corsMiddleware())

	// Initialize services
	serviceContainer := services.NewContainer(dbClient.Repositories, cfg)

	// Start WebSocket hub
	go serviceContainer.Hub.Run()

	// Initialize WebSocket handler
	wsHandler := NewWebSocketHandler(serviceContainer.Hub, serviceContainer)

	// Initialize health handler
	healthHandler := NewHealthHandler(dbClient)

	// Initialize auth handler
	authHandler := NewAuthHandler(serviceContainer.AuthService, serviceContainer.SessionService)

	// Health check routes
	r.Get("/", homeHandler)
	r.Get("/health", healthHandler.Health)
	r.Get("/health/database", healthHandler.DatabaseHealth)
	r.Get("/health/websocket", healthHandler.WebSocketHealth)
	r.Get("/health/system", healthHandler.SystemHealth)
	r.Get("/health/ready", healthHandler.Readiness)
	r.Get("/health/live", healthHandler.Liveness)
	r.Get("/health/metrics", healthHandler.Metrics)

	// Auth success page
	r.Get("/auth/success", authSuccessHandler)

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		// Authentication routes (public)
		SetupAuthRoutes(r, authHandler)

		// Series routes
		r.Route("/series", func(r chi.Router) {
			seriesHandler := NewSeriesHandler(serviceContainer.Series)
			// Public routes (view only)
			r.Get("/", seriesHandler.ListSeries)
			r.Get("/{id}", seriesHandler.GetSeries)

			// Protected routes (require authentication and ownership)
			r.With(middleware.AuthMiddleware(serviceContainer.SessionService)).Post("/", seriesHandler.CreateSeries)
			r.With(middleware.AuthMiddleware(serviceContainer.SessionService)).Put("/{id}", seriesHandler.UpdateSeries)
			r.With(middleware.AuthMiddleware(serviceContainer.SessionService)).Delete("/{id}", seriesHandler.DeleteSeries)
		})

		// Match routes
		r.Route("/matches", func(r chi.Router) {
			matchHandler := NewMatchHandler(serviceContainer.Match)
			// Public routes (view only)
			r.Get("/", matchHandler.ListMatches)
			r.Get("/{id}", matchHandler.GetMatch)
			r.Get("/series/{series_id}", matchHandler.GetMatchesBySeries)

			// Protected routes (require authentication and ownership)
			r.With(middleware.AuthMiddleware(serviceContainer.SessionService)).Post("/", matchHandler.CreateMatch)
			r.With(middleware.AuthMiddleware(serviceContainer.SessionService)).Put("/{id}", matchHandler.UpdateMatch)
			r.With(middleware.AuthMiddleware(serviceContainer.SessionService)).Delete("/{id}", matchHandler.DeleteMatch)
		})

		// Scorecard routes
		r.Route("/scorecard", func(r chi.Router) {
			scorecardHandler := NewScorecardHandler(serviceContainer.Scorecard)
			// Public routes (view only)
			r.Get("/{match_id}", scorecardHandler.GetScorecard)
			r.Get("/{match_id}/current-over", scorecardHandler.GetCurrentOver)
			r.Get("/{match_id}/innings/{innings_number}", scorecardHandler.GetInnings)
			r.Get("/{match_id}/innings/{innings_number}/over/{over_number}", scorecardHandler.GetOver)

			// Protected routes (require authentication and ownership)
			r.With(middleware.AuthMiddleware(serviceContainer.SessionService)).Post("/start", scorecardHandler.StartScoring)
			r.With(middleware.AuthMiddleware(serviceContainer.SessionService)).Post("/ball", scorecardHandler.AddBall)
			r.With(middleware.AuthMiddleware(serviceContainer.SessionService)).Delete("/{match_id}/ball", scorecardHandler.UndoBall)
		})

		// WebSocket routes
		r.Route("/ws", func(r chi.Router) {
			r.Get("/match/{match_id}", wsHandler.ServeWS)
			r.Get("/stats", wsHandler.GetConnectionStats)
			r.Get("/stats/{match_id}", wsHandler.GetRoomStats)
			r.Post("/test/{match_id}", wsHandler.TestBroadcast)
		})

		// GraphQL routes
		r.Route("/graphql", func(r chi.Router) {
			// Use GraphQL handler from the service
			graphqlHandler := serviceContainer.GraphQLWebSocket.GetGraphQLHandler()
			r.Post("/", graphqlHandler.ServeHTTP)
			r.Get("/playground", graphqlHandler.GetPlaygroundHandler().ServeHTTP)
		})
	})

	return r
}

// authSuccessHandler handles the authentication success page
func authSuccessHandler(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>Authentication Successful</title>
    <style>
        body { font-family: Arial, sans-serif; text-align: center; padding: 50px; }
        .success { color: #28a745; }
        .container { max-width: 600px; margin: 0 auto; }
    </style>
</head>
<body>
    <div class="container">
        <h1 class="success">✅ Authentication Successful!</h1>
        <p>You have been successfully authenticated with Google.</p>
        <p>You can now close this window and return to the application.</p>
        <script>
            // Auto-close window after 3 seconds
            setTimeout(function() {
                window.close();
            }, 3000);
        </script>
    </div>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

// corsMiddleware sets up CORS middleware
func corsMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the origin from the request
			origin := r.Header.Get("Origin")

			// Allow specific origins for credentials
			allowedOrigins := []string{
				"http://localhost:3000",
				"http://localhost:3001",
				"http://localhost:3002",
				"http://127.0.0.1:3000",
				"http://127.0.0.1:3001",
				"http://127.0.0.1:3002",
			}

			// Check if origin is allowed
			allowed := false
			for _, allowedOrigin := range allowedOrigins {
				if origin == allowedOrigin {
					allowed = true
					break
				}
			}

			if allowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			} else {
				// Fallback to wildcard for non-credential requests
				w.Header().Set("Access-Control-Allow-Origin", "*")
			}

			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Cache-Control, Pragma, Expires, Accept")
			w.Header().Set("Access-Control-Max-Age", "86400")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Handler functions (placeholder implementations)
func homeHandler(w http.ResponseWriter, r *http.Request) {
	utils.WriteSuccess(w, map[string]string{
		"message": "Welcome to Spark Park Cricket Backend!",
		"version": "1.0.0",
	})
}
