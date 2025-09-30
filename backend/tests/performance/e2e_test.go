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

// E2ETestSuite tests the complete end-to-end workflow
type E2ETestSuite struct {
	router           http.Handler
	dbClient         *database.Client
	cfg              *config.Config
	seriesID         string
	matchID          string
	authCookie       string
	serviceContainer *services.Container
}

// SetupE2ETest sets up the end-to-end test environment
func SetupE2ETest(t *testing.T) *E2ETestSuite {
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
	seriesID, matchID, err := createTestDataForE2E(dbClient, user.ID)
	if err != nil {
		t.Fatalf("Failed to create test data: %v", err)
	}

	return &E2ETestSuite{
		router:           router,
		dbClient:         dbClient,
		cfg:              testCfg.Config,
		seriesID:         seriesID,
		matchID:          matchID,
		authCookie:       authCookie,
		serviceContainer: serviceContainer,
	}
}

// TestCompleteBallAdditionWorkflow tests the complete ball addition workflow
func TestCompleteBallAdditionWorkflow(t *testing.T) {
	suite := SetupE2ETest(t)
	defer suite.dbClient.Close()

	// Step 1: Start the match
	t.Run("StartMatch", func(t *testing.T) {
		err := suite.startMatch()
		if err != nil {
			t.Fatalf("Failed to start match: %v", err)
		}
	})

	// Step 2: Begin scoring
	t.Run("BeginScoring", func(t *testing.T) {
		err := suite.beginScoring()
		if err != nil {
			t.Fatalf("Failed to begin scoring: %v", err)
		}
	})

	// Step 3: Add balls to complete an over
	t.Run("AddBallsToCompleteOver", func(t *testing.T) {
		err := suite.addBallsToCompleteOver()
		if err != nil {
			t.Fatalf("Failed to add balls to complete over: %v", err)
		}
	})

	// Step 4: Verify over completion
	t.Run("VerifyOverCompletion", func(t *testing.T) {
		err := suite.verifyOverCompletion()
		if err != nil {
			t.Fatalf("Failed to verify over completion: %v", err)
		}
	})

	// Step 5: Add balls to complete innings
	t.Run("AddBallsToCompleteInnings", func(t *testing.T) {
		err := suite.addBallsToCompleteInnings()
		if err != nil {
			t.Fatalf("Failed to add balls to complete innings: %v", err)
		}
	})

	// Step 6: Verify innings completion
	t.Run("VerifyInningsCompletion", func(t *testing.T) {
		err := suite.verifyInningsCompletion()
		if err != nil {
			t.Fatalf("Failed to verify innings completion: %v", err)
		}
	})

	// Step 7: Start second innings
	t.Run("StartSecondInnings", func(t *testing.T) {
		err := suite.startSecondInnings()
		if err != nil {
			t.Fatalf("Failed to start second innings: %v", err)
		}
	})

	// Step 8: Complete second innings
	t.Run("CompleteSecondInnings", func(t *testing.T) {
		err := suite.completeSecondInnings()
		if err != nil {
			t.Fatalf("Failed to complete second innings: %v", err)
		}
	})

	// Step 9: Verify match completion
	t.Run("VerifyMatchCompletion", func(t *testing.T) {
		err := suite.verifyMatchCompletion()
		if err != nil {
			t.Fatalf("Failed to verify match completion: %v", err)
		}
	})
}

// startMatch starts the match
func (suite *E2ETestSuite) startMatch() error {
	// This would start the match
	// For now, return nil (mock implementation)
	return nil
}

// beginScoring begins scoring for the match
func (suite *E2ETestSuite) beginScoring() error {
	req := httptest.NewRequest("POST", "/api/v1/scorecard/start", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", suite.authCookie)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		return fmt.Errorf("begin scoring failed with status %d: %s", w.Code, w.Body.String())
	}

	return nil
}

// addBallsToCompleteOver adds balls to complete an over
func (suite *E2ETestSuite) addBallsToCompleteOver() error {
	// Add 6 legal balls to complete an over
	for i := 0; i < 6; i++ {
		ballReq := models.BallEventRequest{
			MatchID:       suite.matchID,
			InningsNumber: 1,
			BallType:      models.BallTypeGood,
			RunType:       models.RunTypeOne,
			IsWicket:      false,
			Byes:          0,
		}

		err := suite.addBall(ballReq)
		if err != nil {
			return fmt.Errorf("failed to add ball %d: %w", i+1, err)
		}

		// Small delay between balls
		time.Sleep(10 * time.Millisecond)
	}

	return nil
}

// verifyOverCompletion verifies that the over is completed
func (suite *E2ETestSuite) verifyOverCompletion() error {
	// Get current over
	req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/scorecard/%s/current-over", suite.matchID), nil)
	req.Header.Set("Cookie", suite.authCookie)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		return fmt.Errorf("get current over failed with status %d: %s", w.Code, w.Body.String())
	}

	var over map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &over)
	if err != nil {
		return fmt.Errorf("failed to unmarshal over: %w", err)
	}

	// Check if over is completed
	if status, exists := over["status"]; !exists || status != "completed" {
		return fmt.Errorf("over is not completed, status: %v", status)
	}

	return nil
}

// addBallsToCompleteInnings adds balls to complete the innings
func (suite *E2ETestSuite) addBallsToCompleteInnings() error {
	// Add multiple overs to complete innings
	// For this test, we'll add 5 more overs (30 balls)
	for over := 1; over <= 5; over++ {
		for ball := 1; ball <= 6; ball++ {
			ballReq := models.BallEventRequest{
				MatchID:       suite.matchID,
				InningsNumber: 1,
				BallType:      models.BallTypeGood,
				RunType:       models.RunTypeOne,
				IsWicket:      false,
				Byes:          0,
			}

			err := suite.addBall(ballReq)
			if err != nil {
				return fmt.Errorf("failed to add ball in over %d, ball %d: %w", over, ball, err)
			}

			// Small delay between balls
			time.Sleep(10 * time.Millisecond)
		}
	}

	return nil
}

// verifyInningsCompletion verifies that the innings is completed
func (suite *E2ETestSuite) verifyInningsCompletion() error {
	// Get innings
	req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/scorecard/%s/innings/1", suite.matchID), nil)
	req.Header.Set("Cookie", suite.authCookie)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		return fmt.Errorf("get innings failed with status %d: %s", w.Code, w.Body.String())
	}

	var innings map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &innings)
	if err != nil {
		return fmt.Errorf("failed to unmarshal innings: %w", err)
	}

	// Check if innings is completed
	if status, exists := innings["status"]; !exists || status != "completed" {
		return fmt.Errorf("innings is not completed, status: %v", status)
	}

	return nil
}

// startSecondInnings starts the second innings
func (suite *E2ETestSuite) startSecondInnings() error {
	// This would start the second innings
	// For now, return nil (mock implementation)
	return nil
}

// completeSecondInnings completes the second innings
func (suite *E2ETestSuite) completeSecondInnings() error {
	// Add balls to complete second innings
	for over := 1; over <= 6; over++ {
		for ball := 1; ball <= 6; ball++ {
			ballReq := models.BallEventRequest{
				MatchID:       suite.matchID,
				InningsNumber: 2,
				BallType:      models.BallTypeGood,
				RunType:       models.RunTypeOne,
				IsWicket:      false,
				Byes:          0,
			}

			err := suite.addBall(ballReq)
			if err != nil {
				return fmt.Errorf("failed to add ball in second innings over %d, ball %d: %w", over, ball, err)
			}

			// Small delay between balls
			time.Sleep(10 * time.Millisecond)
		}
	}

	return nil
}

// verifyMatchCompletion verifies that the match is completed
func (suite *E2ETestSuite) verifyMatchCompletion() error {
	// Get match
	req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/matches/%s", suite.matchID), nil)
	req.Header.Set("Cookie", suite.authCookie)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		return fmt.Errorf("get match failed with status %d: %s", w.Code, w.Body.String())
	}

	var match map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &match)
	if err != nil {
		return fmt.Errorf("failed to unmarshal match: %w", err)
	}

	// Check if match is completed
	if status, exists := match["status"]; !exists || status != "completed" {
		return fmt.Errorf("match is not completed, status: %v", status)
	}

	return nil
}

// addBall adds a single ball
func (suite *E2ETestSuite) addBall(ballReq models.BallEventRequest) error {
	reqBody, err := json.Marshal(ballReq)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req := httptest.NewRequest("POST", "/api/v1/scorecard/ball", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", suite.authCookie)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		return fmt.Errorf("add ball failed with status %d: %s", w.Code, w.Body.String())
	}

	return nil
}

// TestPerformanceDuringE2EWorkflow tests performance during the complete workflow
func TestPerformanceDuringE2EWorkflow(t *testing.T) {
	suite := SetupE2ETest(t)
	defer suite.dbClient.Close()

	// Start the match and begin scoring
	err := suite.startMatch()
	if err != nil {
		t.Fatalf("Failed to start match: %v", err)
	}

	err = suite.beginScoring()
	if err != nil {
		t.Fatalf("Failed to begin scoring: %v", err)
	}

	// Measure performance during ball addition
	responseTimes := make([]time.Duration, 0, 100)

	// Add 100 balls and measure response time for each
	for i := 0; i < 100; i++ {
		ballReq := models.BallEventRequest{
			MatchID:       suite.matchID,
			InningsNumber: 1,
			BallType:      models.BallTypeGood,
			RunType:       models.RunTypeOne,
			IsWicket:      false,
			Byes:          0,
		}

		start := time.Now()
		err := suite.addBall(ballReq)
		duration := time.Since(start)

		if err != nil {
			t.Errorf("Failed to add ball %d: %v", i+1, err)
			continue
		}

		responseTimes = append(responseTimes, duration)

		// Validate performance for each ball
		if duration > 500*time.Millisecond {
			t.Errorf("Ball %d response time exceeds target: %v > 500ms", i+1, duration)
		}
	}

	// Calculate performance statistics
	if len(responseTimes) > 0 {
		var totalDuration time.Duration
		var maxDuration time.Duration
		var minDuration = responseTimes[0]

		for _, duration := range responseTimes {
			totalDuration += duration
			if duration > maxDuration {
				maxDuration = duration
			}
			if duration < minDuration {
				minDuration = duration
			}
		}

		avgDuration := totalDuration / time.Duration(len(responseTimes))

		t.Logf("✅ E2E Performance Results:")
		t.Logf("   Total Balls: %d", len(responseTimes))
		t.Logf("   Average Response Time: %v", avgDuration)
		t.Logf("   Min Response Time: %v", minDuration)
		t.Logf("   Max Response Time: %v", maxDuration)

		// Validate performance targets
		if avgDuration > 300*time.Millisecond {
			t.Errorf("Average response time exceeds target: %v > 300ms", avgDuration)
		}
		if maxDuration > 500*time.Millisecond {
			t.Errorf("Max response time exceeds target: %v > 500ms", maxDuration)
		}
	}
}

// createTestDataForE2E creates test data for end-to-end testing
func createTestDataForE2E(dbClient *database.Client, userID string) (string, string, error) {
	ctx := context.Background()

	// Create test series
	series := &models.Series{
		Name:      fmt.Sprintf("E2E Test Series %d", time.Now().Unix()),
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
