package performance

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"spark-park-cricket-backend/internal/config"
	"spark-park-cricket-backend/internal/database"
	"spark-park-cricket-backend/internal/handlers"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/services"
	"spark-park-cricket-backend/pkg/testutils"
)

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

	// Create test data
	seriesID, matchID, err := createTestDataForIntegration(dbClient, user.ID)
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
	suite := SetupIntegrationTest(t)
	defer suite.dbClient.Close()

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
				req.Header.Set("Cookie", suite.authCookie)

				// Record response time
				start := time.Now()
				w := httptest.NewRecorder()
				suite.router.ServeHTTP(w, req)
				duration := time.Since(start)

				// Validate response
				if w.Code != http.StatusOK {
					t.Errorf("Request failed with status %d: %s", w.Code, w.Body.String())
				}

				// Validate performance
				if duration > 500*time.Millisecond {
					t.Errorf("Response time exceeds target: %v > 500ms", duration)
				}

				// Validate response structure
				var response map[string]interface{}
				err = json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}

				// Check required fields
				requiredFields := []string{"message", "match_id", "innings_number", "ball_type", "run_type", "runs", "byes", "is_wicket"}
				for _, field := range requiredFields {
					if _, exists := response[field]; !exists {
						t.Errorf("Missing required field: %s", field)
					}
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
	req.Header.Set("Cookie", suite.authCookie)

	start := time.Now()
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	duration := time.Since(start)

	// Validate response
	if w.Code != http.StatusOK {
		t.Errorf("Request failed with status %d: %s", w.Code, w.Body.String())
	}

	// Validate performance
	if duration > 200*time.Millisecond {
		t.Errorf("Response time exceeds target: %v > 200ms", duration)
	}

	// Validate response structure
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	// Check required fields
	requiredFields := []string{"match_id", "status", "innings", "teams"}
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
			httpReq.Header.Set("Cookie", suite.authCookie)

			start := time.Now()
			w1 := httptest.NewRecorder()
			suite.router.ServeHTTP(w1, httpReq)
			firstDuration := time.Since(start)

			// Make second request (cache hit)
			httpReq2 := httptest.NewRequest(req.method, req.endpoint, bytes.NewBuffer(req.body))
			httpReq2.Header.Set("Cookie", suite.authCookie)

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
		req.Header.Set("Cookie", suite.authCookie)

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Add ball request failed with status %d: %s", w.Code, w.Body.String())
		}

		// Verify ball was added by getting scorecard
		req2 := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/scorecard/%s", suite.matchID), nil)
		req2.Header.Set("Cookie", suite.authCookie)

		w2 := httptest.NewRecorder()
		suite.router.ServeHTTP(w2, req2)

		if w2.Code != http.StatusOK {
			t.Errorf("Get scorecard request failed with status %d: %s", w2.Code, w2.Body.String())
		}

		// Validate scorecard contains the added ball
		var scorecard map[string]interface{}
		err = json.Unmarshal(w2.Body.Bytes(), &scorecard)
		if err != nil {
			t.Errorf("Failed to unmarshal scorecard: %v", err)
		}

		// Check that innings data exists
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
		req.Header.Set("Cookie", suite.authCookie)

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
		req.Header.Set("Cookie", suite.authCookie)

		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		// Should return error status
		if w.Code == http.StatusOK {
			t.Errorf("Expected error for invalid ball type, got success")
		}
	})
}

// createTestDataForIntegration creates test data for integration testing
func createTestDataForIntegration(dbClient *database.Client, userID string) (string, string, error) {
	ctx := context.Background()

	// Create test series
	series := &models.Series{
		Name:      fmt.Sprintf("Integration Test Series %d", time.Now().Unix()),
		StartDate: time.Now(),
		EndDate:   time.Now().Add(24 * time.Hour),
		CreatedBy: userID,
	}
	err := dbClient.Repositories.Series.Create(ctx, series)
	if err != nil {
		return "", "", fmt.Errorf("failed to create test series: %v", err)
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
		return "", "", fmt.Errorf("failed to create test match: %v", err)
	}

	return series.ID, match.ID, nil
}
