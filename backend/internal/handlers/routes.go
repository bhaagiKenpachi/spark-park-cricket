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
			matchHandler := NewMatchHandler(serviceContainer.Match)
			r.Get("/", matchHandler.ListMatches)
			r.Post("/", matchHandler.CreateMatch)
			r.Get("/{id}", matchHandler.GetMatch)
			r.Put("/{id}", matchHandler.UpdateMatch)
			r.Delete("/{id}", matchHandler.DeleteMatch)
			r.Get("/series/{series_id}", matchHandler.GetMatchesBySeries)
		})

		// Scorecard routes
		r.Route("/scorecard", func(r chi.Router) {
			scorecardHandler := NewScorecardHandler(serviceContainer.Scorecard)
			r.Post("/start", scorecardHandler.StartScoring)
			r.Post("/ball", scorecardHandler.AddBall)
			r.Delete("/{match_id}/ball", scorecardHandler.UndoBall)
			r.Get("/{match_id}", scorecardHandler.GetScorecard)
			r.Get("/{match_id}/current-over", scorecardHandler.GetCurrentOver)
			r.Get("/{match_id}/innings/{innings_number}", scorecardHandler.GetInnings)
			r.Get("/{match_id}/innings/{innings_number}/over/{over_number}", scorecardHandler.GetOver)
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

// corsMiddleware sets up CORS middleware
func corsMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
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
