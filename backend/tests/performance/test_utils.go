package performance

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"spark-park-cricket-backend/internal/config"
	"spark-park-cricket-backend/internal/handlers"
	"spark-park-cricket-backend/internal/services"
	"spark-park-cricket-backend/pkg/testutils"
)

// SetupTestRoutes sets up test routes with authentication
func SetupTestRoutes(serviceContainer *services.Container, seriesHandler *handlers.SeriesHandler, matchHandler *handlers.MatchHandler, scorecardHandler *handlers.ScorecardHandler, cfg *config.Config) http.Handler {
	router := chi.NewRouter()

	// Add middleware
	router.Use(chimiddleware.Recoverer)
	router.Use(chimiddleware.Logger)
	router.Use(chimiddleware.RequestID)
	router.Use(chimiddleware.RealIP)
	router.Use(chimiddleware.Timeout(60 * time.Second))
	router.Use(testutils.CORSMiddleware())

	// Initialize auth handler
	authHandler := handlers.NewAuthHandler(serviceContainer.AuthService, serviceContainer.SessionService, cfg)

	// Health check routes
	router.Get("/health/database", handlers.NewHealthHandler(serviceContainer.DBClient).DatabaseHealth)

	// API routes
	router.Route("/api/v1", func(r chi.Router) {
		// Auth routes (public - no authentication required)
		handlers.SetupAuthRoutes(r, authHandler)

		// Series routes
		r.Route("/series", func(r chi.Router) {
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
	})

	return router
}
