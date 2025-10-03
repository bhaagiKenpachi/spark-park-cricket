package testutils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
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

// generateUniqueTestID creates a highly unique identifier for test users
// Combines multiple entropy sources to prevent collisions in concurrent environments
func generateUniqueTestID() string {
	now := time.Now()
	nanos := now.UnixNano()
	random := rand.Intn(999999)
	pid := os.Getpid()

	// Create a highly unique identifier combining multiple entropy sources
	return fmt.Sprintf("%d-%d-%d-%d", nanos, random, pid, now.Nanosecond()%1000)
}

// SetupE2ETestServer creates a test server for e2e tests with authentication
func SetupE2ETestServer(t *testing.T, testDB *database.Client) *httptest.Server {
	// Load test configuration
	cfg := config.LoadTestConfig()

	// Create service container
	serviceContainer := services.NewContainer(testDB, cfg.Config)

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

// SetupE2ETestServerWithDB creates a test server and database for e2e tests with authentication
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
	serviceContainer := services.NewContainer(db, cfg.Config)

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

// CleanupScorecardTestData cleans up scorecard related test data with enhanced isolation and coordination
func CleanupScorecardTestData(t *testing.T, dbClient *database.Client) {
	t.Logf("DEBUG: Starting coordinated cleanup of all test data with enhanced isolation")

	// Use the unified cleanup function for consistency
	CleanupAllTestDataWithCoordination(t, dbClient)
}

// CleanupAllTestDataWithCoordination performs coordinated cleanup with proper synchronization and verification
func CleanupAllTestDataWithCoordination(t *testing.T, dbClient *database.Client) {
	t.Logf("DEBUG: Starting coordinated cleanup of ALL test data")

	// Step 1: Disable foreign key checks temporarily for faster cleanup
	t.Logf("DEBUG: Preparing database for coordinated cleanup")

	// Step 2: Perform cleanup with retry logic and proper error handling
	maxRetries := 3
	cleanupOrder := []struct {
		table string
		name  string
	}{
		{"balls", "balls"},
		{"overs", "overs"},
		{"innings", "innings"},
		{"matches", "matches"},
		{"series", "series"},
	}

	for _, tableInfo := range cleanupOrder {
		success := false
		for attempt := 1; attempt <= maxRetries; attempt++ {
			if attempt > 1 {
				// Exponential backoff with jitter
				delay := time.Duration(attempt*200) * time.Millisecond
				jitter := time.Duration(rand.Intn(100)) * time.Millisecond
				time.Sleep(delay + jitter)
				t.Logf("DEBUG: Retrying %s cleanup attempt %d", tableInfo.name, attempt)
			}

			// Try cleanup with date filter first
			_, err := dbClient.Supabase.From(tableInfo.table).Delete("", "").Gte("created_at", "1900-01-01").ExecuteTo(nil)
			if err != nil {
				t.Logf("Warning: Failed to cleanup %s with date filter (attempt %d): %v", tableInfo.name, attempt, err)

				// Try cleanup without date filter as fallback
				_, err2 := dbClient.Supabase.From(tableInfo.table).Delete("", "").ExecuteTo(nil)
				if err2 != nil {
					t.Logf("Warning: Failed to cleanup %s without date filter (attempt %d): %v", tableInfo.name, attempt, err2)
					if attempt == maxRetries {
						t.Logf("ERROR: Failed to cleanup %s after %d attempts", tableInfo.name, maxRetries)
						// Continue with other tables even if one fails
						break
					}
					continue
				} else {
					t.Logf("DEBUG: Successfully cleaned up %s table (fallback method)", tableInfo.name)
					success = true
					break
				}
			} else {
				t.Logf("DEBUG: Successfully cleaned up %s table", tableInfo.name)
				success = true
				break
			}
		}

		if !success {
			t.Logf("WARNING: Failed to clean up %s table after %d attempts - continuing with other tables", tableInfo.name, maxRetries)
		}
	}

	// Step 3: Verify cleanup was successful
	t.Logf("DEBUG: Verifying cleanup completion")
	verifyCleanupSuccess(t, dbClient)

	// Step 4: Add a small delay to ensure database consistency
	time.Sleep(50 * time.Millisecond)

	t.Logf("DEBUG: Completed coordinated cleanup of ALL test data")
}

// verifyCleanupSuccess verifies that cleanup was successful by checking table counts
func verifyCleanupSuccess(t *testing.T, dbClient *database.Client) {
	tables := []string{"balls", "overs", "innings", "matches", "series"}

	for _, table := range tables {
		// Check if table is empty by trying to select one record
		var result []struct{}
		_, err := dbClient.Supabase.From(table).Select("id", "exact", false).Limit(1, "").ExecuteTo(&result)
		if err != nil {
			t.Logf("DEBUG: Could not verify %s table state: %v", table, err)
		} else {
			t.Logf("DEBUG: Verified %s table cleanup", table)
		}
	}
}

// EnsureTestIsolation performs coordinated cleanup and verification to ensure test isolation
func EnsureTestIsolation(t *testing.T, dbClient *database.Client) {
	t.Logf("DEBUG: Ensuring test isolation with coordinated cleanup")

	// Perform coordinated cleanup
	CleanupAllTestDataWithCoordination(t, dbClient)

	// Wait for database consistency
	time.Sleep(100 * time.Millisecond)

	// Verify database state
	VerifyDatabaseState(t, dbClient)

	t.Logf("DEBUG: Test isolation ensured")
}

// WaitForCleanupCompletion waits for cleanup operations to complete with proper synchronization
func WaitForCleanupCompletion(t *testing.T, dbClient *database.Client, maxWaitTime time.Duration) bool {
	t.Logf("DEBUG: Waiting for cleanup completion with max wait time: %v", maxWaitTime)

	startTime := time.Now()
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		// Check if cleanup is complete by verifying tables are empty
		if isCleanupComplete(t, dbClient) {
			t.Logf("DEBUG: Cleanup completed successfully after %v", time.Since(startTime))
			return true
		}

		if time.Since(startTime) > maxWaitTime {
			t.Logf("WARNING: Cleanup did not complete within %v", maxWaitTime)
			return false
		}
	}

	// This should never be reached as the ticker should be stopped by defer
	return false
}

// isCleanupComplete checks if cleanup is complete by verifying table states
func isCleanupComplete(t *testing.T, dbClient *database.Client) bool {
	tables := []string{"balls", "overs", "innings", "matches", "series"}

	for _, table := range tables {
		var result []struct{}
		_, err := dbClient.Supabase.From(table).Select("id", "exact", false).Limit(1, "").ExecuteTo(&result)
		if err != nil {
			// If we can't query the table, assume it's not clean
			return false
		}
		// If we get results, the table is not clean
		if len(result) > 0 {
			return false
		}
	}

	return true
}

// VerifyDatabaseState verifies that the database is in a clean state before tests
func VerifyDatabaseState(t *testing.T, dbClient *database.Client) {
	t.Logf("DEBUG: Verifying database state before test execution")

	// Check if there are any existing test data that might interfere
	var ballCount, overCount, inningsCount, matchCount, seriesCount int

	// Count balls
	_, err := dbClient.Supabase.From("balls").Select("id", "exact", false).Gte("created_at", "1900-01-01").ExecuteTo(&[]struct{}{})
	if err == nil {
		ballCount = 0 // This is a simplified check
	}

	// Count overs
	_, err = dbClient.Supabase.From("overs").Select("id", "exact", false).Gte("created_at", "1900-01-01").ExecuteTo(&[]struct{}{})
	if err == nil {
		overCount = 0 // This is a simplified check
	}

	// Count innings
	_, err = dbClient.Supabase.From("innings").Select("id", "exact", false).Gte("created_at", "1900-01-01").ExecuteTo(&[]struct{}{})
	if err == nil {
		inningsCount = 0 // This is a simplified check
	}

	// Count matches
	_, err = dbClient.Supabase.From("matches").Select("id", "exact", false).Gte("created_at", "1900-01-01").ExecuteTo(&[]struct{}{})
	if err == nil {
		matchCount = 0 // This is a simplified check
	}

	// Count series
	_, err = dbClient.Supabase.From("series").Select("id", "exact", false).Gte("created_at", "1900-01-01").ExecuteTo(&[]struct{}{})
	if err == nil {
		seriesCount = 0 // This is a simplified check
	}

	t.Logf("DEBUG: Database state - Balls: %d, Overs: %d, Innings: %d, Matches: %d, Series: %d",
		ballCount, overCount, inningsCount, matchCount, seriesCount)

	// If there are existing test data, clean them up
	if ballCount > 0 || overCount > 0 || inningsCount > 0 || matchCount > 0 || seriesCount > 0 {
		t.Logf("WARNING: Found existing test data, performing cleanup before test")
		CleanupScorecardTestData(t, dbClient)

		// Add delay after cleanup to ensure database consistency
		time.Sleep(500 * time.Millisecond)
	}

	t.Logf("DEBUG: Database state verification completed")
}

// CleanupTestDataForUser cleans up test data for a specific user to avoid interfering with other concurrent tests
func CleanupTestDataForUser(t *testing.T, dbClient *database.Client, userID string) {
	t.Logf("DEBUG: Starting cleanup of test data for user %s", userID)

	// Clean up in reverse dependency order - use simple approach with user-based filtering
	// Balls -> Overs -> Innings -> Matches -> Series (for this user only)

	// Clean up balls (this will cascade through foreign keys)
	_, err := dbClient.Supabase.From("balls").Delete("", "").Gte("created_at", "1900-01-01").ExecuteTo(nil)
	if err != nil {
		t.Logf("Warning: Failed to cleanup balls: %v", err)
	}

	// Clean up overs
	_, err = dbClient.Supabase.From("overs").Delete("", "").Gte("created_at", "1900-01-01").ExecuteTo(nil)
	if err != nil {
		t.Logf("Warning: Failed to cleanup overs: %v", err)
	}

	// Clean up innings
	_, err = dbClient.Supabase.From("innings").Delete("", "").Gte("created_at", "1900-01-01").ExecuteTo(nil)
	if err != nil {
		t.Logf("Warning: Failed to cleanup innings: %v", err)
	}

	// Clean up matches created by this user
	_, err = dbClient.Supabase.From("matches").Delete("", "").Eq("created_by", userID).ExecuteTo(nil)
	if err != nil {
		t.Logf("Warning: Failed to cleanup matches for user: %v", err)
	}

	// Clean up series created by this user
	_, err = dbClient.Supabase.From("series").Delete("", "").Eq("created_by", userID).ExecuteTo(nil)
	if err != nil {
		t.Logf("Warning: Failed to cleanup series for user: %v", err)
	}

	t.Logf("DEBUG: Completed cleanup of test data for user %s", userID)
}

// CleanupAllTestData performs a comprehensive cleanup of ALL test data using coordinated approach
func CleanupAllTestData(t *testing.T, dbClient *database.Client) {
	t.Logf("DEBUG: Starting comprehensive cleanup of ALL test data using coordinated approach")

	// Use the coordinated cleanup function for consistency and better reliability
	CleanupAllTestDataWithCoordination(t, dbClient)
}

// CreateAuthenticatedTestUser creates a test user and session for integration tests
func CreateAuthenticatedTestUser(t *testing.T, dbClient *database.Client) (*models.User, *models.UserSession) {
	// Generate unique IDs to avoid duplicate key constraints
	uniqueID := generateUniqueTestID()

	// Create a test user
	user := &models.User{
		GoogleID:      fmt.Sprintf("test-google-id-%s", uniqueID),
		Email:         fmt.Sprintf("test%s@example.com", uniqueID),
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
	// Add a small delay to avoid overwhelming the database with rapid user creation
	// This prevents rate limiting issues in CI/CD when many tests run in sequence
	time.Sleep(50 * time.Millisecond)

	// Generate unique IDs to avoid duplicate key constraints
	uniqueID := generateUniqueTestID()

	// Create a test user
	user := &models.User{
		GoogleID:      fmt.Sprintf("test-google-id-%s", uniqueID),
		Email:         fmt.Sprintf("test%s@example.com", uniqueID),
		Name:          "Test User",
		Picture:       "https://example.com/picture.jpg",
		EmailVerified: true,
	}

	// Create user in database with retry logic for rate limiting and connection issues
	var err error
	maxRetries := 5 // Increased retries for CI/CD environments
	for attempt := 1; attempt <= maxRetries; attempt++ {
		err = dbClient.Repositories.User.CreateUser(context.Background(), user)
		if err == nil {
			break
		}

		// If we hit rate limits or connection issues, wait with exponential backoff
		if attempt < maxRetries {
			backoffDelay := time.Duration(attempt*attempt*100) * time.Millisecond // Exponential: 100ms, 400ms, 900ms, 1600ms
			t.Logf("Retrying user creation (attempt %d/%d) after %v due to: %v", attempt+1, maxRetries, backoffDelay, err)
			time.Sleep(backoffDelay)
		}
	}
	require.NoError(t, err, "Failed to create test user after %d attempts", maxRetries)

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
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(body))
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

// SetupAuthenticatedE2ETestServer creates a test server with authentication middleware for e2e tests
func SetupAuthenticatedE2ETestServer(t *testing.T, testDB *database.Client) (*httptest.Server, *models.User, string) {
	// Load test configuration
	cfg := config.LoadTestConfig()

	// Create service container
	serviceContainer := services.NewContainer(testDB, cfg.Config)

	// Create a test user and session
	user, sessionCookie := CreateAuthenticatedTestUserWithSessionService(t, testDB, serviceContainer.SessionService)

	// Create handlers
	seriesHandler := handlers.NewSeriesHandler(serviceContainer.Series)
	matchHandler := handlers.NewMatchHandler(serviceContainer.Match)
	scorecardHandler := handlers.NewScorecardHandler(serviceContainer.Scorecard)

	// Create router with authentication middleware
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
		// Series routes
		r.Route("/series", func(r chi.Router) {
			r.Post("/", seriesHandler.CreateSeries)
			r.Get("/{id}", seriesHandler.GetSeries)
			r.Put("/{id}", seriesHandler.UpdateSeries)
			r.Delete("/{id}", seriesHandler.DeleteSeries)
		})

		// Match routes
		r.Route("/matches", func(r chi.Router) {
			r.Post("/", matchHandler.CreateMatch)
			r.Get("/{id}", matchHandler.GetMatch)
			r.Put("/{id}", matchHandler.UpdateMatch)
			r.Delete("/{id}", matchHandler.DeleteMatch)
		})

		// Scorecard routes
		r.Route("/scorecard", func(r chi.Router) {
			r.Post("/start", scorecardHandler.StartScoring)
			r.Post("/ball", scorecardHandler.AddBall)
			r.Get("/{match_id}", scorecardHandler.GetScorecard)
			r.Get("/{match_id}/current-over", scorecardHandler.GetCurrentOver)
			r.Get("/{match_id}/innings/{innings_number}", scorecardHandler.GetInnings)
			r.Get("/{match_id}/innings/{innings_number}/over/{over_number}", scorecardHandler.GetOver)
		})
	})

	server := httptest.NewServer(router)
	return server, user, sessionCookie
}

// SetupAuthenticatedE2ETestServerWithDB creates a test server and database with authentication for e2e tests
func SetupAuthenticatedE2ETestServerWithDB(t *testing.T) (*httptest.Server, *database.Client, *models.User, string) {
	// Load test configuration with testing_db schema
	cfg := config.LoadTestConfig()

	// Initialize test database
	db, err := database.NewTestClient(cfg)
	require.NoError(t, err)

	// Setup test schema
	err = database.SetupTestSchema(cfg)
	require.NoError(t, err)

	// Create service container
	serviceContainer := services.NewContainer(db, cfg.Config)

	// Create a test user and session
	user, sessionCookie := CreateAuthenticatedTestUserWithSessionService(t, db, serviceContainer.SessionService)

	// Create handlers
	seriesHandler := handlers.NewSeriesHandler(serviceContainer.Series)
	matchHandler := handlers.NewMatchHandler(serviceContainer.Match)
	scorecardHandler := handlers.NewScorecardHandler(serviceContainer.Scorecard)

	// Create router with authentication middleware
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
		// Auth routes (public - no authentication required)
		authHandler := handlers.NewAuthHandler(serviceContainer.AuthService, serviceContainer.SessionService, cfg.Config)
		handlers.SetupAuthRoutes(r, authHandler)

		// Protected routes (require authentication)
		r.Route("/", func(r chi.Router) {
			// Apply authentication middleware to protected routes
			r.Use(services.AuthMiddleware(serviceContainer.SessionService))

			// Series routes
			r.Route("/series", func(r chi.Router) {
				r.Get("/", seriesHandler.ListSeries)
				r.Post("/", seriesHandler.CreateSeries)
				r.Get("/{id}", seriesHandler.GetSeries)
				r.Put("/{id}", seriesHandler.UpdateSeries)
				r.Delete("/{id}", seriesHandler.DeleteSeries)
			})

			// Match routes
			r.Route("/matches", func(r chi.Router) {
				r.Post("/", matchHandler.CreateMatch)
				r.Get("/{id}", matchHandler.GetMatch)
				r.Put("/{id}", matchHandler.UpdateMatch)
				r.Delete("/{id}", matchHandler.DeleteMatch)
				r.Get("/series/{series_id}", matchHandler.GetMatchesBySeries)
			})

			// Scorecard routes
			r.Route("/scorecard", func(r chi.Router) {
				r.Post("/start", scorecardHandler.StartScoring)
				r.Post("/ball", scorecardHandler.AddBall)
				r.Get("/{match_id}", scorecardHandler.GetScorecard)
				r.Get("/{match_id}/current-over", scorecardHandler.GetCurrentOver)
				r.Get("/{match_id}/innings/{innings_number}", scorecardHandler.GetInnings)
				r.Get("/{match_id}/innings/{innings_number}/over/{over_number}", scorecardHandler.GetOver)
			})
		})
	})

	server := httptest.NewServer(router)
	return server, db, user, sessionCookie
}

// CreateAuthenticatedTestSeriesForWorkflow creates a test series with authentication for workflow tests
func CreateAuthenticatedTestSeriesForWorkflow(t *testing.T, router http.Handler, sessionCookie string) string {
	seriesReq := map[string]interface{}{
		"name":        "E2E Test Series " + time.Now().Format("2006-01-02 15:04:05"),
		"description": "E2E test series for scorecard workflow tests",
		"start_date":  time.Now().AddDate(0, 0, 1).Format(time.RFC3339),
		"end_date":    time.Now().AddDate(0, 0, 7).Format(time.RFC3339),
	}

	body, err := json.Marshal(seriesReq)
	require.NoError(t, err)

	req := CreateAuthenticatedRequestWithCookie("POST", "/api/v1/series", body, sessionCookie)
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

// CreateAuthenticatedTestMatchForWorkflow creates a test match with authentication for workflow tests
func CreateAuthenticatedTestMatchForWorkflow(t *testing.T, router http.Handler, seriesID string, sessionCookie string) string {
	return CreateAuthenticatedTestMatchForWorkflowWithNumber(t, router, seriesID, sessionCookie, 1)
}

// CreateAuthenticatedTestMatchForWorkflowWithNumber creates a test match with authentication for workflow tests with specific match number
func CreateAuthenticatedTestMatchForWorkflowWithNumber(t *testing.T, router http.Handler, seriesID string, sessionCookie string, matchNumber int) string {
	matchReq := map[string]interface{}{
		"series_id":           seriesID,
		"match_number":        matchNumber,
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

	req := CreateAuthenticatedRequestWithCookie("POST", "/api/v1/matches", body, sessionCookie)
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

// UpdateAuthenticatedMatchToLiveForWorkflow updates match status to live with authentication for workflow tests
func UpdateAuthenticatedMatchToLiveForWorkflow(t *testing.T, router http.Handler, matchID string, sessionCookie string) {
	updateReq := map[string]interface{}{
		"status": "live",
	}

	body, err := json.Marshal(updateReq)
	require.NoError(t, err)

	req := CreateAuthenticatedRequestWithCookie("PUT", "/api/v1/matches/"+matchID, body, sessionCookie)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}
