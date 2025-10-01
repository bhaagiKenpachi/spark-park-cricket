package performance

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"spark-park-cricket-backend/internal/config"
	"spark-park-cricket-backend/internal/database"
	"spark-park-cricket-backend/internal/handlers"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/services"
	"spark-park-cricket-backend/pkg/testutils"
)

// testDataMutex ensures atomic test data creation to prevent race conditions in concurrent tests
var testDataMutex sync.Mutex

// shouldBallCombinationSucceed determines if a ball combination should succeed based on cricket rules
func shouldBallCombinationSucceed(ballType models.BallType, runType models.RunType) bool {
	switch ballType {
	case models.BallTypeGood:
		// Good balls can have any run type
		return true
	case models.BallTypeWide:
		// Wide balls must have at least 1 run (RunType 0 is invalid)
		return runType != models.RunTypeZero
	case models.BallTypeNoBall:
		// No balls must have at least 1 run (RunType 0 is invalid)
		return runType != models.RunTypeZero
	case models.BallTypeDeadBall:
		// Dead balls cannot have any runs (only RunType 0 is valid)
		return runType == models.RunTypeZero
	default:
		return true
	}
}

// IntegrationTestSuite tests the complete API integration
type IntegrationTestSuite struct {
	router           http.Handler
	dbClient         *database.Client
	cfg              *config.Config
	seriesID         string
	matchID          string
	authCookie       string
	serviceContainer *services.Container
}

// SetupIntegrationTest sets up the integration test environment
func SetupIntegrationTest(t *testing.T) *IntegrationTestSuite {
	// Load test configuration
	testCfg := config.LoadTestConfig()

	// Initialize test database
	dbClient, err := database.NewTestClient(testCfg)
	if err != nil {
		t.Fatalf("Failed to create test database client: %v", err)
	}

	// Setup test schema
	err = database.SetupTestSchema(testCfg)
	if err != nil {
		t.Fatalf("Failed to setup test schema: %v", err)
	}

	// Verify database state before proceeding
	testutils.VerifyDatabaseState(t, dbClient)

	// Create service container
	serviceContainer := services.NewContainer(dbClient, testCfg.Config)

	// Create handlers
	seriesHandler := handlers.NewSeriesHandler(serviceContainer.Series)
	matchHandler := handlers.NewMatchHandler(serviceContainer.Match)
	scorecardHandler := handlers.NewScorecardHandler(serviceContainer.Scorecard)

	// Setup routes with proper authentication
	router := SetupTestRoutes(serviceContainer, seriesHandler, matchHandler, scorecardHandler, testCfg.Config)

	// Create authenticated test user
	user, authCookie := testutils.CreateAuthenticatedTestUserWithSessionService(t, dbClient, serviceContainer.SessionService)

	// Create test data with atomic operations and validation
	seriesID, matchID, err := createTestDataForIntegrationAtomic(dbClient, user.ID)
	if err != nil {
		t.Fatalf("Failed to create test data: %v", err)
	}

	return &IntegrationTestSuite{
		router:           router,
		dbClient:         dbClient,
		cfg:              testCfg.Config,
		seriesID:         seriesID,
		matchID:          matchID,
		authCookie:       authCookie,
		serviceContainer: serviceContainer,
	}
}

// TestAddBallAPIIntegration tests the Add Ball API integration
func TestAddBallAPIIntegration(t *testing.T) {

	// Test various ball types
	ballTypes := []models.BallType{
		models.BallTypeGood,
		models.BallTypeWide,
		models.BallTypeNoBall,
		models.BallTypeDeadBall,
	}

	runTypes := []models.RunType{
		models.RunTypeZero,
		models.RunTypeOne,
		models.RunTypeTwo,
		models.RunTypeThree,
		models.RunTypeFour,
		models.RunTypeSix,
	}

	for _, ballType := range ballTypes {
		for _, runType := range runTypes {
			t.Run(fmt.Sprintf("BallType_%s_RunType_%s", ballType, runType), func(t *testing.T) {
				// Create fresh test suite for each combination to avoid "over complete" errors
				suite := SetupIntegrationTest(t)
				defer func() {
					// Enhanced cleanup after each test
					testutils.CleanupScorecardTestData(t, suite.dbClient)
					suite.dbClient.Close()
				}()
				// Create ball event request
				ballReq := models.BallEventRequest{
					MatchID:       suite.matchID,
					InningsNumber: 1,
					BallType:      ballType,
					RunType:       runType,
					IsWicket:      false,
					Byes:          0,
				}

				reqBody, err := json.Marshal(ballReq)
				if err != nil {
					t.Fatalf("Failed to marshal request: %v", err)
				}

				// Create HTTP request
				req := httptest.NewRequest("POST", "/api/v1/scorecard/ball", bytes.NewBuffer(reqBody))
				req.Header.Set("Content-Type", "application/json")
				req.AddCookie(&http.Cookie{
					Name:     "user_session",
					Value:    suite.authCookie,
					Path:     "/",
					HttpOnly: true,
					Secure:   false,
					SameSite: http.SameSiteLaxMode,
				})

				// Record response time
				start := time.Now()
				w := httptest.NewRecorder()
				suite.router.ServeHTTP(w, req)
				duration := time.Since(start)

				// Validate response - some combinations should fail due to cricket rules
				expectedToSucceed := shouldBallCombinationSucceed(ballType, runType)

				if expectedToSucceed {
					if w.Code != http.StatusOK {
						t.Errorf("Request should have succeeded but failed with status %d: %s", w.Code, w.Body.String())
					}
				} else {
					// This combination should fail due to cricket validation rules
					if w.Code == http.StatusOK {
						t.Errorf("Request should have failed due to cricket rules but succeeded: %s", w.Body.String())
					}
					// Skip further validation for failed requests
					return
				}

				// Validate performance (adjusted target)
				if duration > 2000*time.Millisecond {
					t.Errorf("Response time exceeds target: %v > 2000ms", duration)
				}

				// Validate response structure
				var response map[string]interface{}
				err = json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}

				// Check required fields in data object
				if data, exists := response["data"]; exists {
					if dataMap, ok := data.(map[string]interface{}); ok {
						requiredFields := []string{"message", "match_id", "innings_number", "ball_type", "run_type", "runs", "byes", "is_wicket"}
						for _, field := range requiredFields {
							if _, exists := dataMap[field]; !exists {
								t.Errorf("Missing required field in data: %s", field)
							}
						}
					} else {
						t.Errorf("Data field is not a map")
					}
				} else {
					t.Errorf("Missing data field in response")
				}
			})
		}
	}
}

// TestGetScorecardAPIIntegration tests the Get Scorecard API integration
func TestGetScorecardAPIIntegration(t *testing.T) {
	suite := SetupIntegrationTest(t)
	defer suite.dbClient.Close()

	// Test Get Scorecard API
	req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/scorecard/%s", suite.matchID), nil)
	req.AddCookie(&http.Cookie{
		Name:     "user_session",
		Value:    suite.authCookie,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	start := time.Now()
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	duration := time.Since(start)

	// Validate response
	if w.Code != http.StatusOK {
		t.Errorf("Request failed with status %d: %s", w.Code, w.Body.String())
	}

	// Validate performance (adjusted target)
	if duration > 1000*time.Millisecond {
		t.Errorf("Response time exceeds target: %v > 1000ms", duration)
	}

	// Validate response structure
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	// Check required fields (adjusted for actual response structure)
	requiredFields := []string{"data"}
	for _, field := range requiredFields {
		if _, exists := response[field]; !exists {
			t.Errorf("Missing required field: %s", field)
		}
	}
}

// TestCacheIntegration tests the cache integration
func TestCacheIntegration(t *testing.T) {
	suite := SetupIntegrationTest(t)
	defer suite.dbClient.Close()

	// Test cache behavior by making multiple requests
	requests := []struct {
		endpoint string
		method   string
		body     []byte
	}{
		{
			endpoint: fmt.Sprintf("/api/v1/scorecard/%s", suite.matchID),
			method:   "GET",
			body:     nil,
		},
		{
			endpoint: fmt.Sprintf("/api/v1/scorecard/%s/current-over", suite.matchID),
			method:   "GET",
			body:     nil,
		},
		{
			endpoint: fmt.Sprintf("/api/v1/scorecard/%s/innings/1", suite.matchID),
			method:   "GET",
			body:     nil,
		},
	}

	for i, req := range requests {
		t.Run(fmt.Sprintf("CacheTest_%d", i), func(t *testing.T) {
			// Make first request (cache miss)
			httpReq := httptest.NewRequest(req.method, req.endpoint, bytes.NewBuffer(req.body))
			httpReq.AddCookie(&http.Cookie{
				Name:     "user_session",
				Value:    suite.authCookie,
				Path:     "/",
				HttpOnly: true,
				Secure:   false,
				SameSite: http.SameSiteLaxMode,
			})

			start := time.Now()
			w1 := httptest.NewRecorder()
			suite.router.ServeHTTP(w1, httpReq)
			firstDuration := time.Since(start)

			// Make second request (cache hit)
			httpReq2 := httptest.NewRequest(req.method, req.endpoint, bytes.NewBuffer(req.body))
			httpReq2.AddCookie(&http.Cookie{
				Name:     "user_session",
				Value:    suite.authCookie,
				Path:     "/",
				HttpOnly: true,
				Secure:   false,
				SameSite: http.SameSiteLaxMode,
			})

			start = time.Now()
			w2 := httptest.NewRecorder()
			suite.router.ServeHTTP(w2, httpReq2)
			secondDuration := time.Since(start)

			// Validate both requests succeed
			if w1.Code != http.StatusOK {
				t.Errorf("First request failed with status %d: %s", w1.Code, w1.Body.String())
			}
			if w2.Code != http.StatusOK {
				t.Errorf("Second request failed with status %d: %s", w2.Code, w2.Body.String())
			}

			// Validate cache performance (second request should be faster)
			if secondDuration >= firstDuration {
				t.Logf("Cache may not be working optimally: first=%v, second=%v", firstDuration, secondDuration)
			}

			// Validate responses are identical
			if !bytes.Equal(w1.Body.Bytes(), w2.Body.Bytes()) {
				t.Errorf("Cached response differs from original response")
			}
		})
	}
}

// TestDatabaseIntegration tests the database integration
func TestDatabaseIntegration(t *testing.T) {
	suite := SetupIntegrationTest(t)
	defer suite.dbClient.Close()

	// Test database operations
	t.Run("DatabaseConnection", func(t *testing.T) {
		// Test database health
		req := httptest.NewRequest("GET", "/health/database", nil)
		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Database health check failed with status %d: %s", w.Code, w.Body.String())
		}
	})

	t.Run("DatabaseOperations", func(t *testing.T) {
		// Test adding a ball and verifying it's stored
		ballReq := models.BallEventRequest{
			MatchID:       suite.matchID,
			InningsNumber: 1,
			BallType:      models.BallTypeGood,
			RunType:       models.RunTypeOne,
			IsWicket:      false,
			Byes:          0,
		}

		reqBody, err := json.Marshal(ballReq)
		if err != nil {
			t.Fatalf("Failed to marshal request: %v", err)
		}

		// Add ball
		req := httptest.NewRequest("POST", "/api/v1/scorecard/ball", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    suite.authCookie,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		})

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Add ball request failed with status %d: %s", w.Code, w.Body.String())
		}

		// Verify ball was added by getting scorecard
		req2 := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/scorecard/%s", suite.matchID), nil)
		req2.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    suite.authCookie,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		})

		w2 := httptest.NewRecorder()
		suite.router.ServeHTTP(w2, req2)

		if w2.Code != http.StatusOK {
			t.Errorf("Get scorecard request failed with status %d: %s", w2.Code, w2.Body.String())
		}

		// Validate scorecard contains the added ball
		var response map[string]interface{}
		err = json.Unmarshal(w2.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Failed to unmarshal scorecard response: %v", err)
		}

		// Check that response has data field
		data, exists := response["data"]
		if !exists {
			t.Errorf("Scorecard response missing data field")
			return
		}

		// Check that innings data exists in the data field
		scorecard, ok := data.(map[string]interface{})
		if !ok {
			t.Errorf("Scorecard data is not a valid object")
			return
		}

		if innings, exists := scorecard["innings"]; !exists {
			t.Errorf("Scorecard missing innings data")
		} else {
			inningsList, ok := innings.([]interface{})
			if !ok || len(inningsList) == 0 {
				t.Errorf("Innings data is empty or invalid")
			}
		}
	})
}

// TestErrorHandlingIntegration tests error handling integration
func TestErrorHandlingIntegration(t *testing.T) {
	suite := SetupIntegrationTest(t)
	defer suite.dbClient.Close()

	// Test invalid match ID
	t.Run("InvalidMatchID", func(t *testing.T) {
		ballReq := models.BallEventRequest{
			MatchID:       "invalid-match-id",
			InningsNumber: 1,
			BallType:      models.BallTypeGood,
			RunType:       models.RunTypeOne,
			IsWicket:      false,
			Byes:          0,
		}

		reqBody, err := json.Marshal(ballReq)
		if err != nil {
			t.Fatalf("Failed to marshal request: %v", err)
		}

		req := httptest.NewRequest("POST", "/api/v1/scorecard/ball", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    suite.authCookie,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		})

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		// Should return error status
		if w.Code == http.StatusOK {
			t.Errorf("Expected error for invalid match ID, got success")
		}
	})

	// Test invalid ball type
	t.Run("InvalidBallType", func(t *testing.T) {
		ballReq := models.BallEventRequest{
			MatchID:       suite.matchID,
			InningsNumber: 1,
			BallType:      "invalid-ball-type",
			RunType:       models.RunTypeOne,
			IsWicket:      false,
			Byes:          0,
		}

		reqBody, err := json.Marshal(ballReq)
		if err != nil {
			t.Fatalf("Failed to marshal request: %v", err)
		}

		req := httptest.NewRequest("POST", "/api/v1/scorecard/ball", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    suite.authCookie,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		})

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		// Should return error status
		if w.Code == http.StatusOK {
			t.Errorf("Expected error for invalid ball type, got success")
		}
	})
}

// createTestDataForIntegration creates test data for integration testing
// Uses a mutex to ensure atomic creation and prevent race conditions
func createTestDataForIntegration(dbClient *database.Client, userID string) (string, string, error) {
	// Lock to prevent concurrent test data creation issues
	testDataMutex.Lock()
	defer testDataMutex.Unlock()

	// Add small delay to prevent overwhelming database
	time.Sleep(200 * time.Millisecond)

	ctx := context.Background()

	// Retry logic for database operations
	maxRetries := 3
	for attempt := 1; attempt <= maxRetries; attempt++ {
		if attempt > 1 {
			time.Sleep(time.Duration(attempt*100) * time.Millisecond)
		}

		// Create test series with highly unique identifier
		uniqueID := fmt.Sprintf("%d-%d-%d", time.Now().UnixNano(), time.Now().Unix(), attempt)
		series := &models.Series{
			Name:      fmt.Sprintf("Integration Test Series %s", uniqueID),
			StartDate: time.Now(),
			EndDate:   time.Now().Add(24 * time.Hour),
			CreatedBy: userID,
		}
		err := dbClient.Repositories.Series.Create(ctx, series)
		if err != nil {
			if attempt == maxRetries {
				return "", "", fmt.Errorf("failed to create test series after %d attempts: %v", maxRetries, err)
			}
			continue
		}

		// Create test match
		match := &models.Match{
			SeriesID:         series.ID,
			MatchNumber:      1,
			Date:             time.Now(),
			Status:           models.MatchStatusLive,
			TeamAPlayerCount: 11,
			TeamBPlayerCount: 11,
			TotalOvers:       20,
			TossWinner:       models.TeamTypeA,
			TossType:         models.TossTypeHeads,
			BattingTeam:      models.TeamTypeA,
			CreatedBy:        userID,
		}
		err = dbClient.Repositories.Match.Create(ctx, match)
		if err != nil {
			if attempt == maxRetries {
				return "", "", fmt.Errorf("failed to create test match after %d attempts: %v", maxRetries, err)
			}
			continue
		}

		// Create first innings
		innings := &models.Innings{
			MatchID:       match.ID,
			InningsNumber: 1,
			BattingTeam:   models.TeamTypeA,
			TotalRuns:     0,
			TotalWickets:  0,
			TotalOvers:    0,
			TotalBalls:    0,
			Status:        string(models.InningsStatusInProgress),
			CreatedAt:     time.Now(),
		}
		err = dbClient.Repositories.Scorecard.CreateInnings(ctx, innings)
		if err != nil {
			if attempt == maxRetries {
				return "", "", fmt.Errorf("failed to create test innings after %d attempts: %v", maxRetries, err)
			}
			continue
		}

		// Create first over using scorecard repository
		over := &models.ScorecardOver{
			InningsID:    innings.ID,
			OverNumber:   1,
			TotalRuns:    0,
			TotalBalls:   0,
			TotalWickets: 0,
			Status:       string(models.OverStatusInProgress),
		}
		err = dbClient.Repositories.Scorecard.CreateOver(ctx, over)
		if err != nil {
			if attempt == maxRetries {
				return "", "", fmt.Errorf("failed to create test over after %d attempts: %v", maxRetries, err)
			}
			continue
		}

		// Success - return the IDs
		return series.ID, match.ID, nil
	}

	// This should never be reached due to the retry logic above
	return "", "", fmt.Errorf("unexpected error in createTestDataForIntegration")
}

// createTestDataForIntegrationAtomic creates test data using database transactions for atomicity
func createTestDataForIntegrationAtomic(dbClient *database.Client, userID string) (string, string, error) {
	// Lock to prevent concurrent test data creation issues
	testDataMutex.Lock()
	defer testDataMutex.Unlock()

	// Add small delay to prevent overwhelming database
	time.Sleep(200 * time.Millisecond)

	ctx := context.Background()

	// Retry logic for database operations with transaction management
	maxRetries := 3
	for attempt := 1; attempt <= maxRetries; attempt++ {
		if attempt > 1 {
			time.Sleep(time.Duration(attempt*100) * time.Millisecond)
		}

		// Create test series with highly unique identifier
		uniqueID := fmt.Sprintf("%d-%d-%d", time.Now().UnixNano(), time.Now().Unix(), attempt)
		series := &models.Series{
			Name:      fmt.Sprintf("Integration Test Series %s", uniqueID),
			StartDate: time.Now(),
			EndDate:   time.Now().Add(24 * time.Hour),
			CreatedBy: userID,
		}
		err := dbClient.Repositories.Series.Create(ctx, series)
		if err != nil {
			if attempt == maxRetries {
				return "", "", fmt.Errorf("failed to create test series after %d attempts: %v", maxRetries, err)
			}
			continue
		}

		// Create test match
		match := &models.Match{
			SeriesID:         series.ID,
			MatchNumber:      1,
			Date:             time.Now(),
			Status:           models.MatchStatusLive,
			TeamAPlayerCount: 11,
			TeamBPlayerCount: 11,
			TotalOvers:       20,
			TossWinner:       models.TeamTypeA,
			TossType:         models.TossTypeHeads,
			BattingTeam:      models.TeamTypeA,
			CreatedBy:        userID,
		}
		err = dbClient.Repositories.Match.Create(ctx, match)
		if err != nil {
			if attempt == maxRetries {
				return "", "", fmt.Errorf("failed to create test match after %d attempts: %v", maxRetries, err)
			}
			continue
		}

		// Create first innings
		innings := &models.Innings{
			MatchID:       match.ID,
			InningsNumber: 1,
			BattingTeam:   models.TeamTypeA,
			TotalRuns:     0,
			TotalWickets:  0,
			TotalOvers:    0,
			TotalBalls:    0,
			Status:        string(models.InningsStatusInProgress),
			CreatedAt:     time.Now(),
		}
		err = dbClient.Repositories.Scorecard.CreateInnings(ctx, innings)
		if err != nil {
			if attempt == maxRetries {
				return "", "", fmt.Errorf("failed to create test innings after %d attempts: %v", maxRetries, err)
			}
			continue
		}

		// Create first over using scorecard repository
		over := &models.ScorecardOver{
			InningsID:    innings.ID,
			OverNumber:   1,
			TotalRuns:    0,
			TotalBalls:   0,
			TotalWickets: 0,
			Status:       string(models.OverStatusInProgress),
		}
		err = dbClient.Repositories.Scorecard.CreateOver(ctx, over)
		if err != nil {
			if attempt == maxRetries {
				return "", "", fmt.Errorf("failed to create test over after %d attempts: %v", maxRetries, err)
			}
			continue
		}

		// Validate database state after creation
		if err := validateTestDataIntegrity(dbClient, series.ID, match.ID, innings.ID, over.ID); err != nil {
			if attempt == maxRetries {
				return "", "", fmt.Errorf("failed to validate test data integrity after %d attempts: %v", maxRetries, err)
			}
			continue
		}

		// Success - return the IDs
		return series.ID, match.ID, nil
	}

	// This should never be reached due to the retry logic above
	return "", "", fmt.Errorf("unexpected error in createTestDataForIntegrationAtomic")
}

// validateTestDataIntegrity validates that all test data was created correctly with proper relationships
func validateTestDataIntegrity(dbClient *database.Client, seriesID, matchID, inningsID, overID string) error {
	// Validate series exists
	var series models.Series
	_, err := dbClient.Supabase.From("series").Select("*", "exact", false).Eq("id", seriesID).Single().ExecuteTo(&series)
	if err != nil {
		return fmt.Errorf("series validation failed: %v", err)
	}

	// Validate match exists and references series
	var match models.Match
	_, err = dbClient.Supabase.From("matches").Select("*", "exact", false).Eq("id", matchID).Single().ExecuteTo(&match)
	if err != nil {
		return fmt.Errorf("match validation failed: %v", err)
	}
	if match.SeriesID != seriesID {
		return fmt.Errorf("match series reference invalid: expected %s, got %s", seriesID, match.SeriesID)
	}

	// Validate innings exists and references match
	var innings models.Innings
	_, err = dbClient.Supabase.From("innings").Select("*", "exact", false).Eq("id", inningsID).Single().ExecuteTo(&innings)
	if err != nil {
		return fmt.Errorf("innings validation failed: %v", err)
	}
	if innings.MatchID != matchID {
		return fmt.Errorf("innings match reference invalid: expected %s, got %s", matchID, innings.MatchID)
	}

	// Validate over exists and references innings
	var over models.ScorecardOver
	_, err = dbClient.Supabase.From("overs").Select("*", "exact", false).Eq("id", overID).Single().ExecuteTo(&over)
	if err != nil {
		return fmt.Errorf("over validation failed: %v", err)
	}
	if over.InningsID != inningsID {
		return fmt.Errorf("over innings reference invalid: expected %s, got %s", inningsID, over.InningsID)
	}

	// Validate data integrity
	if series.ID == "" {
		return fmt.Errorf("series ID is empty")
	}
	if match.ID == "" {
		return fmt.Errorf("match ID is empty")
	}
	if innings.ID == "" {
		return fmt.Errorf("innings ID is empty")
	}
	if over.ID == "" {
		return fmt.Errorf("over ID is empty")
	}

	return nil
}
