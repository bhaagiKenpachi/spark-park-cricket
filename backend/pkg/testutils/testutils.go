package testutils

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"spark-park-cricket-backend/internal/config"
	"spark-park-cricket-backend/internal/database"
	"spark-park-cricket-backend/internal/handlers"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/services"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/require"
)

// SetupE2ETestServer creates a test server for e2e tests
func SetupE2ETestServer(t *testing.T, testDB *database.Client) *httptest.Server {
	// Load test configuration
	cfg := config.LoadTestConfig()

	// Create service container
	serviceContainer := services.NewContainer(testDB.Repositories, cfg.Config)

	// Create handlers
	seriesHandler := handlers.NewSeriesHandler(serviceContainer.Series)
	matchHandler := handlers.NewMatchHandler(serviceContainer.Match)
	scorecardHandler := handlers.NewScorecardHandler(serviceContainer.Scorecard)

	// Create router and register routes
	router := http.NewServeMux()

	// Series routes
	router.HandleFunc("/api/v1/series", seriesHandler.CreateSeries)

	// Match routes
	router.HandleFunc("/api/v1/matches", matchHandler.CreateMatch)
	router.HandleFunc("/api/v1/matches/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if len(path) > len("/api/v1/matches/") {
			matchID := path[len("/api/v1/matches/"):]
			switch r.Method {
			case http.MethodGet:
				matchHandler.GetMatch(w, r)
			case http.MethodPut:
				matchHandler.UpdateMatch(w, r)
			case http.MethodDelete:
				matchHandler.DeleteMatch(w, r)
			}
			// Store matchID in context for handlers to use
			_ = matchID
		}
	})

	// Scorecard routes
	router.HandleFunc("/api/v1/scorecard/ball", scorecardHandler.AddBall)
	router.HandleFunc("/api/v1/scorecard/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if len(path) > len("/api/v1/scorecard/") {
			matchID := path[len("/api/v1/scorecard/"):]
			switch r.Method {
			case http.MethodGet:
				// Call service directly since we're using http.ServeMux
				scorecard, err := serviceContainer.Scorecard.GetScorecard(r.Context(), matchID)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(map[string]interface{}{
					"data": scorecard,
				}); err != nil {
					http.Error(w, "Failed to encode response", http.StatusInternalServerError)
					return
				}
			}
		}
	})

	return httptest.NewServer(router)
}

// SetupE2ETestServerWithDB creates a test server and database for e2e tests
func SetupE2ETestServerWithDB(t *testing.T) (*httptest.Server, *database.Client) {
	// Load test configuration with testing_db schema
	cfg := config.LoadTestConfig()

	// Initialize test database
	db, err := database.NewTestClient(cfg)
	require.NoError(t, err)

	// Setup test schema
	err = database.SetupTestSchema(cfg)
	require.NoError(t, err)

	// Create service container
	serviceContainer := services.NewContainer(db.Repositories, cfg.Config)

	// Create handlers
	seriesHandler := handlers.NewSeriesHandler(serviceContainer.Series)
	matchHandler := handlers.NewMatchHandler(serviceContainer.Match)
	scorecardHandler := handlers.NewScorecardHandler(serviceContainer.Scorecard)

	// Create router and register all routes
	router := http.NewServeMux()

	// Series routes
	router.HandleFunc("/api/v1/series", seriesHandler.CreateSeries)
	router.HandleFunc("/api/v1/series/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if len(path) > len("/api/v1/series/") {
			seriesID := path[len("/api/v1/series/"):]
			switch r.Method {
			case http.MethodGet:
				seriesHandler.GetSeries(w, r)
			case http.MethodPut:
				seriesHandler.UpdateSeries(w, r)
			case http.MethodDelete:
				seriesHandler.DeleteSeries(w, r)
			}
			_ = seriesID
		}
	})

	// Match routes
	router.HandleFunc("/api/v1/matches", matchHandler.CreateMatch)
	router.HandleFunc("/api/v1/matches/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if len(path) > len("/api/v1/matches/") {
			matchID := path[len("/api/v1/matches/"):]
			switch r.Method {
			case http.MethodGet:
				matchHandler.GetMatch(w, r)
			case http.MethodPut:
				matchHandler.UpdateMatch(w, r)
			case http.MethodDelete:
				matchHandler.DeleteMatch(w, r)
			}
			_ = matchID
		}
	})

	// Scorecard routes
	router.HandleFunc("/api/v1/scorecard/start", scorecardHandler.StartScoring)
	router.HandleFunc("/api/v1/scorecard/ball", scorecardHandler.AddBall)
	router.HandleFunc("/api/v1/scorecard/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if len(path) > len("/api/v1/scorecard/") {
			matchID := path[len("/api/v1/scorecard/"):]
			switch r.Method {
			case http.MethodGet:
				// Call service directly since we're using http.ServeMux
				scorecard, err := serviceContainer.Scorecard.GetScorecard(r.Context(), matchID)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(map[string]interface{}{
					"data": scorecard,
				}); err != nil {
					http.Error(w, "Failed to encode response", http.StatusInternalServerError)
					return
				}
			}
		}
	})

	server := httptest.NewServer(router)
	return server, db
}

// CORSMiddleware returns a CORS middleware function
func CORSMiddleware() func(http.Handler) http.Handler {
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

// SetupScorecardTestRouter creates a router for scorecard integration tests
func SetupScorecardTestRouter(scorecardHandler *handlers.ScorecardHandler, serviceContainer *services.Container) http.Handler {
	router := chi.NewRouter()

	// Add middleware
	router.Use(chimiddleware.Recoverer)
	router.Use(chimiddleware.Logger)
	router.Use(chimiddleware.RequestID)
	router.Use(chimiddleware.RealIP)
	router.Use(chimiddleware.Timeout(60 * time.Second))
	router.Use(CORSMiddleware())

	// API routes
	router.Route("/api/v1", func(r chi.Router) {
		// Series routes (needed for creating matches)
		r.Route("/series", func(r chi.Router) {
			seriesHandler := handlers.NewSeriesHandler(serviceContainer.Series)
			r.Post("/", seriesHandler.CreateSeries)
			r.Get("/{id}", seriesHandler.GetSeries)
			r.Put("/{id}", seriesHandler.UpdateSeries)
		})
		// Match routes (needed for creating matches)
		r.Route("/matches", func(r chi.Router) {
			matchHandler := handlers.NewMatchHandler(serviceContainer.Match)
			r.Post("/", matchHandler.CreateMatch)
			r.Get("/{id}", matchHandler.GetMatch)
			r.Put("/{id}", matchHandler.UpdateMatch)
		})
		// Scorecard routes
		r.Route("/scorecard", func(r chi.Router) {
			r.Post("/start", scorecardHandler.StartScoring)
			r.Post("/ball", scorecardHandler.AddBall)
			r.Delete("/{match_id}/ball", scorecardHandler.UndoBall)
			r.Get("/{match_id}", scorecardHandler.GetScorecard)
			r.Get("/{match_id}/current-over", scorecardHandler.GetCurrentOver)
			r.Get("/{match_id}/innings/{innings_number}", scorecardHandler.GetInnings)
			r.Get("/{match_id}/innings/{innings_number}/over/{over_number}", scorecardHandler.GetOver)
		})
	})

	return router
}

// StringPtr creates a pointer to a string value
func StringPtr(s string) *string {
	return &s
}

// TimePtr creates a pointer to a time.Time value
func TimePtr(t time.Time) *time.Time {
	return &t
}

// TeamTypePtr creates a pointer to a TeamType value
func TeamTypePtr(teamType models.TeamType) *models.TeamType {
	return &teamType
}

// MatchStatusPtr creates a pointer to a MatchStatus value
func MatchStatusPtr(status models.MatchStatus) *models.MatchStatus {
	return &status
}

// CleanupTestData cleans up test data from the database
func CleanupTestData(t *testing.T, testDB *database.Client) {
	// Clean up matches - use a condition that will match all records
	_, err := testDB.Supabase.From("matches").Delete("", "").Gte("created_at", "1900-01-01").ExecuteTo(nil)
	if err != nil {
		t.Logf("Warning: Failed to cleanup matches: %v", err)
	}

	// Clean up series - use a condition that will match all records
	_, err = testDB.Supabase.From("series").Delete("", "").Gte("created_at", "1900-01-01").ExecuteTo(nil)
	if err != nil {
		t.Logf("Warning: Failed to cleanup series: %v", err)
	}
}

// CreateTestSeriesForWorkflow creates a test series for workflow tests
func CreateTestSeriesForWorkflow(t *testing.T, router http.Handler) string {
	seriesReq := map[string]interface{}{
		"name":        "E2E Test Series " + time.Now().Format("2006-01-02 15:04:05"),
		"description": "E2E test series for scorecard workflow tests",
		"start_date":  time.Now().AddDate(0, 0, 1).Format(time.RFC3339),
		"end_date":    time.Now().AddDate(0, 0, 7).Format(time.RFC3339),
	}

	body, err := json.Marshal(seriesReq)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/api/v1/series", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var response struct {
		Data models.Series `json:"data"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	return response.Data.ID
}

// CreateTestMatchForWorkflow creates a test match for workflow tests
func CreateTestMatchForWorkflow(t *testing.T, router http.Handler, seriesID string) string {
	matchReq := map[string]interface{}{
		"series_id":           seriesID,
		"match_number":        1,
		"date":                time.Now().AddDate(0, 0, 1).Format(time.RFC3339),
		"venue":               "E2E Test Venue",
		"team_a_player_count": 11,
		"team_b_player_count": 11,
		"total_overs":         20,
		"toss_winner":         "A",
		"toss_type":           "H",
		"batting_team":        "A",
	}

	body, err := json.Marshal(matchReq)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/api/v1/matches", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var response struct {
		Data models.Match `json:"data"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	return response.Data.ID
}

// UpdateMatchToLiveForWorkflow updates match status to live for workflow tests
func UpdateMatchToLiveForWorkflow(t *testing.T, router http.Handler, matchID string) {
	updateReq := map[string]interface{}{
		"status": "live",
	}

	body, err := json.Marshal(updateReq)
	require.NoError(t, err)

	req := httptest.NewRequest("PUT", "/api/v1/matches/"+matchID, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

// CleanupScorecardTestData cleans up scorecard related test data
func CleanupScorecardTestData(t *testing.T, dbClient *database.Client) {
	t.Logf("DEBUG: Starting comprehensive cleanup of all test data")

	// Clean up scorecard related tables in reverse order of dependencies
	// Balls -> Overs -> Innings -> Matches -> Series

	// Clean up balls - delete ALL records using a condition that matches all
	_, err := dbClient.Supabase.From("balls").Delete("", "").Gte("created_at", "1900-01-01").ExecuteTo(nil)
	if err != nil {
		t.Logf("Warning: Failed to cleanup balls: %v", err)
	} else {
		t.Logf("DEBUG: Successfully cleaned up balls table")
	}

	// Clean up overs - delete ALL records using a condition that matches all
	_, err = dbClient.Supabase.From("overs").Delete("", "").Gte("created_at", "1900-01-01").ExecuteTo(nil)
	if err != nil {
		t.Logf("Warning: Failed to cleanup overs: %v", err)
	} else {
		t.Logf("DEBUG: Successfully cleaned up overs table")
	}

	// Clean up innings - delete ALL records using a condition that matches all
	_, err = dbClient.Supabase.From("innings").Delete("", "").Gte("created_at", "1900-01-01").ExecuteTo(nil)
	if err != nil {
		t.Logf("Warning: Failed to cleanup innings: %v", err)
	} else {
		t.Logf("DEBUG: Successfully cleaned up innings table")
	}

	// Clean up matches - delete ALL records using a condition that matches all
	_, err = dbClient.Supabase.From("matches").Delete("", "").Gte("created_at", "1900-01-01").ExecuteTo(nil)
	if err != nil {
		t.Logf("Warning: Failed to cleanup matches: %v", err)
	} else {
		t.Logf("DEBUG: Successfully cleaned up matches table")
	}

	// Clean up series - delete ALL records using a condition that matches all
	_, err = dbClient.Supabase.From("series").Delete("", "").Gte("created_at", "1900-01-01").ExecuteTo(nil)
	if err != nil {
		t.Logf("Warning: Failed to cleanup series: %v", err)
	} else {
		t.Logf("DEBUG: Successfully cleaned up series table")
	}

	t.Logf("DEBUG: Completed comprehensive cleanup of all test data")
}

// CreateAuthenticatedTestUser creates a test user and session for integration tests
func CreateAuthenticatedTestUser(t *testing.T, dbClient *database.Client) (*models.User, *models.UserSession) {
	// Create a test user
	user := &models.User{
		GoogleID:      "test-google-id-123",
		Email:         "test@example.com",
		Name:          "Test User",
		Picture:       "https://example.com/picture.jpg",
		EmailVerified: true,
	}

	// Create user in database
	err := dbClient.Repositories.User.CreateUser(context.Background(), user)
	require.NoError(t, err, "Failed to create test user")

	// Create session for the user
	session := &models.UserSession{
		UserID:    user.ID,
		SessionID: "test-session-123",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	err = dbClient.Repositories.User.CreateUserSession(context.Background(), session)
	require.NoError(t, err, "Failed to create test session")

	return user, session
}

// CreateAuthenticatedTestUserWithSessionService creates a test user and proper session using session service
func CreateAuthenticatedTestUserWithSessionService(t *testing.T, dbClient *database.Client, sessionService *services.SessionService) (*models.User, string) {
	// Create a test user
	user := &models.User{
		GoogleID:      "test-google-id-123",
		Email:         "test@example.com",
		Name:          "Test User",
		Picture:       "https://example.com/picture.jpg",
		EmailVerified: true,
	}

	// Create user in database
	err := dbClient.Repositories.User.CreateUser(context.Background(), user)
	require.NoError(t, err, "Failed to create test user")

	// Create a proper session using the session service
	req := httptest.NewRequest("POST", "/auth/login", nil)
	w := httptest.NewRecorder()

	err = sessionService.CreateSession(w, req, user)
	require.NoError(t, err, "Failed to create test session")

	// Extract the session cookie from the response
	cookies := w.Result().Cookies()
	var sessionCookie string
	for _, cookie := range cookies {
		if cookie.Name == "user_session" {
			sessionCookie = cookie.Value
			break
		}
	}
	require.NotEmpty(t, sessionCookie, "Session cookie not found in response")

	return user, sessionCookie
}

// CreateAuthenticatedRequestWithCookie creates an HTTP request with authentication cookie
func CreateAuthenticatedRequestWithCookie(method, url string, body []byte, sessionCookie string) *http.Request {
	req := httptest.NewRequest(method, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Add session cookie
	req.AddCookie(&http.Cookie{
		Name:     "user_session",
		Value:    sessionCookie,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production
		SameSite: http.SameSiteLaxMode,
	})

	return req
}

// CreateAuthenticatedRequest creates an HTTP request with authentication cookie
func CreateAuthenticatedRequest(method, url string, body []byte, session *models.UserSession) *http.Request {
	req := httptest.NewRequest(method, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Add session cookie
	req.AddCookie(&http.Cookie{
		Name:     "user_session",
		Value:    session.SessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production
		SameSite: http.SameSiteLaxMode,
	})

	return req
}
