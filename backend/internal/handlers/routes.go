package handlers

import (
	"net/http"
	"spark-park-cricket-backend/internal/database"
	"spark-park-cricket-backend/internal/middleware"
	"spark-park-cricket-backend/internal/services"
	"spark-park-cricket-backend/internal/utils"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(dbClient *database.Client) *chi.Mux {
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
	serviceContainer := services.NewContainer(dbClient.Repositories)

	// Start WebSocket hub
	go serviceContainer.Hub.Run()

	// Initialize WebSocket handler
	wsHandler := NewWebSocketHandler(serviceContainer.Hub, serviceContainer)

	// Initialize health handler
	healthHandler := NewHealthHandler(dbClient)

	// Health check routes
	r.Get("/", homeHandler)
	r.Get("/health", healthHandler.Health)
	r.Get("/health/database", healthHandler.DatabaseHealth)
	r.Get("/health/websocket", healthHandler.WebSocketHealth)
	r.Get("/health/system", healthHandler.SystemHealth)
	r.Get("/health/ready", healthHandler.Readiness)
	r.Get("/health/live", healthHandler.Liveness)
	r.Get("/health/metrics", healthHandler.Metrics)

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		// Series routes
		r.Route("/series", func(r chi.Router) {
			seriesHandler := NewSeriesHandler(serviceContainer.Series)
			r.Get("/", seriesHandler.ListSeries)
			r.Post("/", seriesHandler.CreateSeries)
			r.Get("/{id}", seriesHandler.GetSeries)
			r.Put("/{id}", seriesHandler.UpdateSeries)
			r.Delete("/{id}", seriesHandler.DeleteSeries)
		})

		// Match routes
		r.Route("/matches", func(r chi.Router) {
			r.Get("/", listMatchesHandler)
			r.Post("/", createMatchHandler)
			r.Get("/{id}", getMatchHandler)
			r.Put("/{id}", updateMatchHandler)
			r.Delete("/{id}", deleteMatchHandler)
		})

		// Team routes
		r.Route("/teams", func(r chi.Router) {
			r.Get("/", listTeamsHandler)
			r.Post("/", createTeamHandler)
			r.Get("/{id}", getTeamHandler)
			r.Put("/{id}", updateTeamHandler)
			r.Get("/{id}/players", listTeamPlayersHandler)
			r.Post("/{id}/players", addTeamPlayerHandler)
		})

		// Player routes
		r.Route("/players", func(r chi.Router) {
			r.Get("/", listPlayersHandler)
			r.Post("/", createPlayerHandler)
			r.Get("/{id}", getPlayerHandler)
			r.Put("/{id}", updatePlayerHandler)
			r.Delete("/{id}", deletePlayerHandler)
		})

		// Scoreboard routes
		r.Route("/scoreboard", func(r chi.Router) {
			scoreboardHandler := NewScoreboardHandler(serviceContainer.Scoreboard)
			r.Get("/{match_id}", scoreboardHandler.GetScoreboard)
			r.Post("/{match_id}/ball", scoreboardHandler.AddBall)
			r.Put("/{match_id}/score", scoreboardHandler.UpdateScore)
			r.Put("/{match_id}/wicket", scoreboardHandler.UpdateWicket)
		})

		// WebSocket routes
		r.Route("/ws", func(r chi.Router) {
			r.Get("/match/{match_id}", wsHandler.ServeWS)
			r.Get("/stats", wsHandler.GetConnectionStats)
			r.Get("/stats/{match_id}", wsHandler.GetRoomStats)
			r.Post("/test/{match_id}", wsHandler.TestBroadcast)
		})
	})

	return r
}

// corsMiddleware sets up CORS middleware
func corsMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
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

func healthHandler(w http.ResponseWriter, r *http.Request) {
	utils.WriteSuccess(w, map[string]string{
		"status":  "OK",
		"service": "spark-park-cricket-backend",
	})
}

func dbHealthHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement database health check
	utils.WriteSuccess(w, map[string]string{
		"status":   "OK",
		"database": "connected",
	})
}

// Placeholder handlers for API endpoints
func listSeriesHandler(w http.ResponseWriter, r *http.Request) {
	utils.WriteSuccess(w, []interface{}{})
}

func createSeriesHandler(w http.ResponseWriter, r *http.Request) {
	utils.WriteCreated(w, map[string]string{"message": "Series created"})
}

func getSeriesHandler(w http.ResponseWriter, r *http.Request) {
	utils.WriteSuccess(w, map[string]string{"message": "Series details"})
}

func updateSeriesHandler(w http.ResponseWriter, r *http.Request) {
	utils.WriteSuccess(w, map[string]string{"message": "Series updated"})
}

func deleteSeriesHandler(w http.ResponseWriter, r *http.Request) {
	utils.WriteSuccess(w, map[string]string{"message": "Series deleted"})
}

func listMatchesHandler(w http.ResponseWriter, r *http.Request) {
	utils.WriteSuccess(w, []interface{}{})
}

func createMatchHandler(w http.ResponseWriter, r *http.Request) {
	utils.WriteCreated(w, map[string]string{"message": "Match created"})
}

func getMatchHandler(w http.ResponseWriter, r *http.Request) {
	utils.WriteSuccess(w, map[string]string{"message": "Match details"})
}

func updateMatchHandler(w http.ResponseWriter, r *http.Request) {
	utils.WriteSuccess(w, map[string]string{"message": "Match updated"})
}

func deleteMatchHandler(w http.ResponseWriter, r *http.Request) {
	utils.WriteSuccess(w, map[string]string{"message": "Match deleted"})
}

func listTeamsHandler(w http.ResponseWriter, r *http.Request) {
	utils.WriteSuccess(w, []interface{}{})
}

func createTeamHandler(w http.ResponseWriter, r *http.Request) {
	utils.WriteCreated(w, map[string]string{"message": "Team created"})
}

func getTeamHandler(w http.ResponseWriter, r *http.Request) {
	utils.WriteSuccess(w, map[string]string{"message": "Team details"})
}

func updateTeamHandler(w http.ResponseWriter, r *http.Request) {
	utils.WriteSuccess(w, map[string]string{"message": "Team updated"})
}

func listTeamPlayersHandler(w http.ResponseWriter, r *http.Request) {
	utils.WriteSuccess(w, []interface{}{})
}

func addTeamPlayerHandler(w http.ResponseWriter, r *http.Request) {
	utils.WriteCreated(w, map[string]string{"message": "Player added to team"})
}

func listPlayersHandler(w http.ResponseWriter, r *http.Request) {
	utils.WriteSuccess(w, []interface{}{})
}

func createPlayerHandler(w http.ResponseWriter, r *http.Request) {
	utils.WriteCreated(w, map[string]string{"message": "Player created"})
}

func getPlayerHandler(w http.ResponseWriter, r *http.Request) {
	utils.WriteSuccess(w, map[string]string{"message": "Player details"})
}

func updatePlayerHandler(w http.ResponseWriter, r *http.Request) {
	utils.WriteSuccess(w, map[string]string{"message": "Player updated"})
}

func deletePlayerHandler(w http.ResponseWriter, r *http.Request) {
	utils.WriteSuccess(w, map[string]string{"message": "Player deleted"})
}

func getScoreboardHandler(w http.ResponseWriter, r *http.Request) {
	utils.WriteSuccess(w, map[string]string{"message": "Scoreboard details"})
}

func addBallHandler(w http.ResponseWriter, r *http.Request) {
	utils.WriteCreated(w, map[string]string{"message": "Ball added"})
}

func updateScoreHandler(w http.ResponseWriter, r *http.Request) {
	utils.WriteSuccess(w, map[string]string{"message": "Score updated"})
}

func updateWicketHandler(w http.ResponseWriter, r *http.Request) {
	utils.WriteSuccess(w, map[string]string{"message": "Wicket updated"})
}
