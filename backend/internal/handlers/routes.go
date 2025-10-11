package handlers

import (
	"context"
	"fmt"
	"net/http"
	"spark-park-cricket-backend/internal/config"
	"spark-park-cricket-backend/internal/database"
	"spark-park-cricket-backend/internal/graphql"
	"spark-park-cricket-backend/internal/monitoring"
	"spark-park-cricket-backend/internal/services"
	"spark-park-cricket-backend/internal/utils"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(dbClient *database.Client, cfg *config.Config) *chi.Mux {
	r := chi.NewRouter()

	// Initialize services
	serviceContainer := services.NewContainer(dbClient, cfg)

	// Middleware
	r.Use(services.RecoveryMiddleware)
	r.Use(services.LoggerMiddleware)
	r.Use(services.RequestIDMiddleware)
	r.Use(chimiddleware.RealIP)
	r.Use(services.TimeoutMiddleware(60 * time.Second))
	r.Use(services.SecurityMiddleware)
	r.Use(services.ValidationMiddleware)
	r.Use(services.MetricsMiddleware)
	r.Use(services.PrometheusMiddleware(serviceContainer.Metrics)) // Add Prometheus metrics middleware
	r.Use(services.RateLimitMiddleware(100))                       // 100 requests per minute
	r.Use(corsMiddleware(cfg))

	// Initialize health handler
	healthHandler := NewHealthHandler(dbClient)

	// Initialize auth handler
	authHandler := NewAuthHandler(serviceContainer.AuthService, serviceContainer.SessionService, cfg)

	// Health check routes
	r.Get("/", homeHandler)
	r.Get("/health", healthHandler.Health)
	r.Get("/health/database", healthHandler.DatabaseHealth)
	r.Get("/health/system", healthHandler.SystemHealth)
	r.Get("/health/ready", healthHandler.Readiness)
	r.Get("/health/live", healthHandler.Liveness)
	r.Get("/health/metrics", healthHandler.Metrics)

	// Prometheus metrics endpoint
	r.Get("/metrics", services.PrometheusHandler().ServeHTTP)

	// Monitoring endpoints
	r.Get("/monitoring/prometheus", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, cfg.PrometheusURL, http.StatusTemporaryRedirect)
	})
	r.Get("/monitoring/grafana", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, cfg.GrafanaURL, http.StatusTemporaryRedirect)
	})
	r.Get("/monitoring", func(w http.ResponseWriter, r *http.Request) {
		monitoringHandler(w, r, cfg)
	})

	// Test endpoint for database monitoring (no auth required)
	r.Get("/test-db-monitoring", func(w http.ResponseWriter, r *http.Request) {
		// Simulate a database operation with monitoring
		ctx := r.Context()
		metrics := serviceContainer.Metrics

		// Test the database monitoring
		err := monitoring.WithDatabaseMonitoringContext(
			ctx, metrics, "SELECT", "test_table", "test_match_id",
			func(ctx context.Context) error {
				// Simulate database operation
				time.Sleep(100 * time.Millisecond)
				return nil
			},
		)

		if err != nil {
			http.Error(w, "Database monitoring test failed", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("Database monitoring test completed successfully")); err != nil {
			// Log error if write fails
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
	})

	// Auth success page
	r.Get("/auth/success", authSuccessHandler)

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		// Authentication routes (public)
		SetupAuthRoutes(r, authHandler)

		// User routes
		r.Route("/users", func(r chi.Router) {
			userHandler := NewUserHandler(serviceContainer.UserService)
			// Protected routes (require authentication)
			r.With(services.AuthMiddleware(serviceContainer.SessionService)).Put("/me", userHandler.UpdateUserName)
		})

		// Series routes
		r.Route("/series", func(r chi.Router) {
			seriesHandler := NewSeriesHandler(serviceContainer.Series)
			// Public routes (view only)
			r.Get("/", seriesHandler.ListSeries)
			r.Get("/{id}", seriesHandler.GetSeries)

			// Protected routes (require authentication and ownership)
			r.With(services.AuthMiddleware(serviceContainer.SessionService)).Post("/", seriesHandler.CreateSeries)
			r.With(services.AuthMiddleware(serviceContainer.SessionService)).Put("/{id}", seriesHandler.UpdateSeries)
			r.With(services.AuthMiddleware(serviceContainer.SessionService)).Delete("/{id}", seriesHandler.DeleteSeries)
		})

		// Match routes
		r.Route("/matches", func(r chi.Router) {
			matchHandler := NewMatchHandler(serviceContainer.Match)
			// Public routes (view only)
			r.Get("/", matchHandler.ListMatches)
			r.Get("/{id}", matchHandler.GetMatch)
			r.Get("/series/{series_id}", matchHandler.GetMatchesBySeries)

			// Protected routes (require authentication and ownership)
			r.With(services.AuthMiddleware(serviceContainer.SessionService)).Post("/", matchHandler.CreateMatch)
			r.With(services.AuthMiddleware(serviceContainer.SessionService)).Put("/{id}", matchHandler.UpdateMatch)
			r.With(services.AuthMiddleware(serviceContainer.SessionService)).Delete("/{id}", matchHandler.DeleteMatch)
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
			r.With(services.AuthMiddleware(serviceContainer.SessionService)).Post("/start", scorecardHandler.StartScoring)
			r.With(services.AuthMiddleware(serviceContainer.SessionService)).Post("/ball", scorecardHandler.AddBall)
			r.With(services.AuthMiddleware(serviceContainer.SessionService)).Delete("/{match_id}/ball", scorecardHandler.UndoBall)
		})

		// Vote routes
		r.Route("/votes", func(r chi.Router) {
			voteHandler := NewVoteHandler(serviceContainer.VoteService)
			teamHandler := NewVoteTeamHandler(serviceContainer.VoteTeamService, serviceContainer.Metrics)

			// Public routes (view only)
			r.Get("/", voteHandler.ListVotes)
			r.Get("/{id}", voteHandler.GetVote)
			r.Get("/{id}/results", voteHandler.GetVoteWithResults)
			r.Get("/{vote_id}/teams", teamHandler.GetTeamsByVoteID)
			// Optional auth - works for both logged in and logged out users
			r.With(services.OptionalAuthMiddleware(serviceContainer.SessionService)).Get("/{id}/has-voted", voteHandler.HasUserVoted)

			// Protected routes (require authentication)
			r.With(services.AuthMiddleware(serviceContainer.SessionService)).Post("/", voteHandler.CreateVote)
			r.With(services.AuthMiddleware(serviceContainer.SessionService)).Put("/{id}", voteHandler.UpdateVote)
			r.With(services.AuthMiddleware(serviceContainer.SessionService)).Delete("/{id}", voteHandler.DeleteVote)
			r.With(services.AuthMiddleware(serviceContainer.SessionService)).Post("/{id}/vote", voteHandler.CastVote)
			r.With(services.AuthMiddleware(serviceContainer.SessionService)).Get("/{id}/my-vote", voteHandler.GetUserVote)
			r.With(services.AuthMiddleware(serviceContainer.SessionService)).Post("/{id}/close", voteHandler.CloseVote)
			r.With(services.AuthMiddleware(serviceContainer.SessionService)).Post("/{id}/cancel", voteHandler.CancelVote)

			// Team routes (protected - any logged-in user can manage teams)
			r.With(services.AuthMiddleware(serviceContainer.SessionService)).Post("/{vote_id}/teams", teamHandler.CreateTeam)
		})

		// Team management routes
		r.Route("/teams", func(r chi.Router) {
			teamHandler := NewVoteTeamHandler(serviceContainer.VoteTeamService, serviceContainer.Metrics)

			// Public routes
			r.Get("/{team_id}", teamHandler.GetTeamByID)
			r.Get("/{team_id}/players", teamHandler.GetTeamPlayers)

			// Protected routes (any logged-in user can manage teams)
			r.With(services.AuthMiddleware(serviceContainer.SessionService)).Put("/{team_id}", teamHandler.UpdateTeam)
			r.With(services.AuthMiddleware(serviceContainer.SessionService)).Delete("/{team_id}", teamHandler.DeleteTeam)
			r.With(services.AuthMiddleware(serviceContainer.SessionService)).Post("/{team_id}/players", teamHandler.AddPlayerToTeam)
			r.With(services.AuthMiddleware(serviceContainer.SessionService)).Post("/{team_id}/players/bulk", teamHandler.AddPlayersToTeam)
			r.With(services.AuthMiddleware(serviceContainer.SessionService)).Delete("/{team_id}/players/{player_id}", teamHandler.RemovePlayerFromTeam)
		})

		// GraphQL routes
		r.Route("/graphql", func(r chi.Router) {
			// Create GraphQL handler
			graphqlHandler := graphql.NewGraphQLHandler(serviceContainer.Scorecard)
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
	_, _ = w.Write([]byte(html))
}

// corsMiddleware sets up CORS middleware
func corsMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the origin from the request
			origin := r.Header.Get("Origin")

			// Parse allowed origins from config
			allowedOrigins := strings.Split(cfg.AllowedOrigins, ",")
			for i, allowedOrigin := range allowedOrigins {
				allowedOrigins[i] = strings.TrimSpace(allowedOrigin)
			}

			// Debug logging for CORS
			if origin != "" {
				fmt.Printf("DEBUG: CORS - Request from origin: %s\n", origin)
				fmt.Printf("DEBUG: CORS - Allowed origins: %v\n", allowedOrigins)
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
				if origin != "" {
					fmt.Printf("DEBUG: CORS - Origin %s is allowed\n", origin)
				}
			} else {
				// Fallback to wildcard for non-credential requests
				w.Header().Set("Access-Control-Allow-Origin", "*")
				if origin != "" {
					fmt.Printf("DEBUG: CORS - Origin %s not allowed, using wildcard\n", origin)
				}
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

// monitoringHandler provides a monitoring dashboard with links to Prometheus and Grafana
func monitoringHandler(w http.ResponseWriter, r *http.Request, cfg *config.Config) {
	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <title>Spark Park Cricket - Monitoring Dashboard</title>
    <style>
        body { 
            font-family: Arial, sans-serif; 
            margin: 0; 
            padding: 20px; 
            background-color: #f5f5f5; 
        }
        .container { 
            max-width: 1200px; 
            margin: 0 auto; 
            background: white; 
            padding: 30px; 
            border-radius: 10px; 
            box-shadow: 0 2px 10px rgba(0,0,0,0.1); 
        }
        .header { 
            text-align: center; 
            margin-bottom: 30px; 
            color: #333; 
        }
        .monitoring-grid { 
            display: grid; 
            grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); 
            gap: 20px; 
            margin-top: 30px; 
        }
        .monitoring-card { 
            background: #f8f9fa; 
            padding: 20px; 
            border-radius: 8px; 
            border-left: 4px solid #007bff; 
            text-align: center; 
        }
        .monitoring-card h3 { 
            margin-top: 0; 
            color: #007bff; 
        }
        .monitoring-card p { 
            color: #666; 
            margin-bottom: 20px; 
        }
        .btn { 
            display: inline-block; 
            padding: 10px 20px; 
            background: #007bff; 
            color: white; 
            text-decoration: none; 
            border-radius: 5px; 
            transition: background 0.3s; 
        }
        .btn:hover { 
            background: #0056b3; 
        }
        .status { 
            display: inline-block; 
            padding: 5px 10px; 
            border-radius: 15px; 
            font-size: 12px; 
            font-weight: bold; 
            margin-left: 10px; 
        }
        .status.active { 
            background: #d4edda; 
            color: #155724; 
        }
        .info-grid { 
            display: grid; 
            grid-template-columns: repeat(auto-fit, minmax(250px, 1fr)); 
            gap: 15px; 
            margin-top: 20px; 
        }
        .info-item { 
            background: #e9ecef; 
            padding: 15px; 
            border-radius: 5px; 
        }
        .info-item strong { 
            color: #495057; 
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🏏 Spark Park Cricket - Monitoring Dashboard</h1>
            <p>Real-time monitoring and observability for your cricket application</p>
        </div>

        <div class="monitoring-grid">
            <div class="monitoring-card">
                <h3>📊 Prometheus</h3>
                <p>Metrics collection and querying</p>
                <a href="%s" target="_blank" class="btn">Open Prometheus</a>
                <span class="status active">Active</span>
            </div>

            <div class="monitoring-card">
                <h3>📈 Grafana</h3>
                <p>Visualization and dashboards</p>
                <a href="%s" target="_blank" class="btn">Open Grafana</a>
                <span class="status active">Active</span>
            </div>
        </div>

        <div class="info-grid">
            <div class="info-item">
                <strong>Prometheus URL:</strong><br>
                <a href="%s" target="_blank">%s</a>
            </div>
            <div class="info-item">
                <strong>Grafana URL:</strong><br>
                <a href="%s" target="_blank">%s</a>
            </div>
            <div class="info-item">
                <strong>Grafana Login:</strong><br>
                Username: admin<br>
                Password: admin123
            </div>
            <div class="info-item">
                <strong>Backend Health:</strong><br>
                <a href="/health" target="_blank">Check Health Status</a>
            </div>
        </div>

        <div style="margin-top: 30px; text-align: center; color: #666;">
            <p>💡 <strong>Tip:</strong> Use Grafana to create custom dashboards for your cricket metrics!</p>
        </div>
    </div>
</body>
</html>`,
		cfg.PrometheusURL, cfg.GrafanaURL,
		cfg.PrometheusURL, cfg.PrometheusURL,
		cfg.GrafanaURL, cfg.GrafanaURL)

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(html))
}
