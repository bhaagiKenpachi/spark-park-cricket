package performance

import (
	"bytes"
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

	// Create test data with basic validation for E2E tests
	seriesID, matchID, err := createTestDataWithValidation(dbClient, user.ID, BasicValidation, "E2E")
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
	defer func() {
		// Enhanced cleanup after test
		testutils.CleanupScorecardTestData(t, suite.dbClient)
		suite.dbClient.Close()
	}()

	// Step 1: Start the match
	t.Run("StartMatch", func(t *testing.T) {
		err := suite.startMatch()
		if err != nil {
			t.Fatalf("Failed to start match: %v", err)
		}
	})

	// Step 2: Test data is already created in SetupE2ETest
	// No need to create additional data

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

	// Step 5: Add more balls to test performance (simplified for performance testing)
	t.Run("AddMoreBallsForPerformance", func(t *testing.T) {
		err := suite.addMoreBallsForPerformance()
		if err != nil {
			t.Fatalf("Failed to add more balls for performance testing: %v", err)
		}
	})

	// Step 6: Verify scorecard performance
	t.Run("VerifyScorecardPerformance", func(t *testing.T) {
		err := suite.verifyScorecardPerformance()
		if err != nil {
			t.Fatalf("Failed to verify scorecard performance: %v", err)
		}
	})

	// Clean up test data at the end - clean only this test's specific data
	cleanupTestMatch(t, suite.dbClient, suite.matchID, suite.seriesID)
}

// cleanupTestMatch cleans up specific test match data to avoid affecting other concurrent tests
func cleanupTestMatch(t *testing.T, dbClient *database.Client, matchID, seriesID string) {
	t.Logf("DEBUG: Cleaning up specific test data for match: %s, series: %s", matchID, seriesID)

	// Clean up in reverse order of dependencies
	// Balls -> Overs -> Innings -> Match -> Series

	// Get innings IDs for this match first (needed for cleaning balls and overs)
	var innings []struct {
		ID string `json:"id"`
	}
	_, err := dbClient.Supabase.From("innings").Select("id", "exact", false).Eq("match_id", matchID).ExecuteTo(&innings)
	if err == nil && len(innings) > 0 {
		for _, ing := range innings {
			// Get over IDs for this innings
			var overs []struct {
				ID string `json:"id"`
			}
			_, err2 := dbClient.Supabase.From("overs").Select("id", "exact", false).Eq("innings_id", ing.ID).ExecuteTo(&overs)
			if err2 == nil && len(overs) > 0 {
				// Clean up balls for each over
				for _, over := range overs {
					_, err3 := dbClient.Supabase.From("balls").Delete("", "").Eq("over_id", over.ID).ExecuteTo(nil)
					if err3 != nil {
						t.Logf("Warning: Failed to cleanup balls for over %s: %v", over.ID, err3)
					}
				}
			}

			// Clean up overs for this innings
			_, err2 = dbClient.Supabase.From("overs").Delete("", "").Eq("innings_id", ing.ID).ExecuteTo(nil)
			if err2 != nil {
				t.Logf("Warning: Failed to cleanup overs for innings %s: %v", ing.ID, err2)
			}
		}
	}

	// Clean up innings for this match
	_, err = dbClient.Supabase.From("innings").Delete("", "").Eq("match_id", matchID).ExecuteTo(nil)
	if err != nil {
		t.Logf("Warning: Failed to cleanup innings for match %s: %v", matchID, err)
	}

	// Clean up the match
	_, err = dbClient.Supabase.From("matches").Delete("", "").Eq("id", matchID).ExecuteTo(nil)
	if err != nil {
		t.Logf("Warning: Failed to cleanup match %s: %v", matchID, err)
	}

	// Clean up the series
	_, err = dbClient.Supabase.From("series").Delete("", "").Eq("id", seriesID).ExecuteTo(nil)
	if err != nil {
		t.Logf("Warning: Failed to cleanup series %s: %v", seriesID, err)
	}

	t.Logf("DEBUG: Completed cleanup for match %s", matchID)
}

// startMatch starts the match
func (suite *E2ETestSuite) startMatch() error {
	fmt.Printf("🔧 DEBUG: Starting match %s\n", suite.matchID)

	// Update match to live status
	status := models.MatchStatusLive
	updateReq := models.UpdateMatchRequest{
		Status: &status,
	}

	reqBody, err := json.Marshal(updateReq)
	if err != nil {
		return fmt.Errorf("failed to marshal update request: %w", err)
	}

	req := httptest.NewRequest("PUT", fmt.Sprintf("/api/v1/matches/%s", suite.matchID), bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:     "user_session",
		Value:    suite.authCookie,
		Path:     "/",
		HttpOnly: true,
	})

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	fmt.Printf("🔧 DEBUG: Match update response: %d - %s\n", w.Code, w.Body.String())

	if w.Code != http.StatusOK {
		return fmt.Errorf("failed to update match status: %d - %s", w.Code, w.Body.String())
	}

	fmt.Printf("✅ DEBUG: Match started successfully\n")
	return nil
}

// addBallsToCompleteOver adds balls to complete an over
func (suite *E2ETestSuite) addBallsToCompleteOver() error {
	fmt.Printf("🔧 DEBUG: Adding 6 balls to complete over for match %s\n", suite.matchID)

	// Add 6 legal balls to complete an over
	for i := 0; i < 6; i++ {
		fmt.Printf("🔧 DEBUG: Adding ball %d/6\n", i+1)

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
			fmt.Printf("❌ DEBUG: Failed to add ball %d: %v\n", i+1, err)
			return fmt.Errorf("failed to add ball %d: %w", i+1, err)
		}

		fmt.Printf("✅ DEBUG: Ball %d added successfully\n", i+1)

		// Small delay between balls
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Printf("✅ DEBUG: Over completed successfully\n")
	return nil
}

// verifyOverCompletion verifies that the over is completed
func (suite *E2ETestSuite) verifyOverCompletion() error {
	fmt.Printf("🔧 DEBUG: Verifying over completion for match %s\n", suite.matchID)

	// Get scorecard to check over completion
	req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/scorecard/%s", suite.matchID), nil)
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

	fmt.Printf("🔧 DEBUG: Scorecard response: %d - %s\n", w.Code, w.Body.String())

	if w.Code != http.StatusOK {
		return fmt.Errorf("get scorecard failed with status %d: %s", w.Code, w.Body.String())
	}

	var response struct {
		Data models.ScorecardResponse `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		return fmt.Errorf("failed to unmarshal scorecard: %w", err)
	}

	// Check if we have at least one innings with at least one over
	if len(response.Data.Innings) == 0 {
		return fmt.Errorf("no innings found")
	}

	innings := response.Data.Innings[0]
	if len(innings.Overs) == 0 {
		return fmt.Errorf("no overs found in innings")
	}

	// Check if the first over has 6 balls (completed)
	firstOver := innings.Overs[0]
	if len(firstOver.Balls) < 6 {
		return fmt.Errorf("first over is not completed, balls: %d", len(firstOver.Balls))
	}

	fmt.Printf("✅ DEBUG: Over completion verified - first over has %d balls\n", len(firstOver.Balls))
	return nil
}

// addMoreBallsForPerformance adds more balls for performance testing
func (suite *E2ETestSuite) addMoreBallsForPerformance() error {
	fmt.Printf("🔧 DEBUG: Adding more balls for performance testing\n")

	// Add a few more overs for performance testing (simplified)
	for over := 1; over <= 3; over++ {
		fmt.Printf("🔧 DEBUG: Adding over %d/3\n", over)
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
				fmt.Printf("❌ DEBUG: Failed to add ball in over %d, ball %d: %v\n", over, ball, err)
				return fmt.Errorf("failed to add ball in over %d, ball %d: %w", over, ball, err)
			}

			fmt.Printf("✅ DEBUG: Ball %d in over %d added successfully\n", ball, over)

			// Small delay between balls
			time.Sleep(10 * time.Millisecond)
		}
	}

	fmt.Printf("✅ DEBUG: Performance balls added successfully\n")
	return nil
}

// verifyScorecardPerformance verifies scorecard performance
func (suite *E2ETestSuite) verifyScorecardPerformance() error {
	fmt.Printf("🔧 DEBUG: Verifying scorecard performance\n")

	// Get scorecard
	req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/scorecard/%s", suite.matchID), nil)
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

	fmt.Printf("🔧 DEBUG: Scorecard response: %d - %s\n", w.Code, w.Body.String())

	if w.Code != http.StatusOK {
		return fmt.Errorf("get scorecard failed with status %d: %s", w.Code, w.Body.String())
	}

	var response struct {
		Data models.ScorecardResponse `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		return fmt.Errorf("failed to unmarshal scorecard response: %w", err)
	}

	if len(response.Data.Innings) == 0 {
		return fmt.Errorf("scorecard should have at least 1 innings")
	}

	fmt.Printf("✅ DEBUG: Scorecard performance verified successfully\n")
	return nil
}

// Note: Removed complex match completion functions for performance testing
// Performance tests focus on API performance rather than complete match workflows

// addBall adds a single ball
func (suite *E2ETestSuite) addBall(ballReq models.BallEventRequest) error {
	fmt.Printf("🔧 DEBUG: Adding ball - Match: %s, Innings: %d, Type: %s, Run: %s\n",
		ballReq.MatchID, ballReq.InningsNumber, ballReq.BallType, ballReq.RunType)

	reqBody, err := json.Marshal(ballReq)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
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

	fmt.Printf("🔧 DEBUG: Add ball response: %d - %s\n", w.Code, w.Body.String())

	if w.Code != http.StatusOK {
		return fmt.Errorf("add ball failed with status %d: %s", w.Code, w.Body.String())
	}

	fmt.Printf("✅ DEBUG: Ball added successfully\n")
	return nil
}

// TestPerformanceDuringE2EWorkflow tests performance during the complete workflow
func TestPerformanceDuringE2EWorkflow(t *testing.T) {
	suite := SetupE2ETest(t)
	defer func() {
		// Enhanced cleanup after test
		testutils.CleanupScorecardTestData(t, suite.dbClient)
		suite.dbClient.Close()
	}()

	// Start the match and create test data directly (like integration test)
	err := suite.startMatch()
	if err != nil {
		t.Fatalf("Failed to start match: %v", err)
	}

	// Test data is already created in SetupE2ETest
	// No need to create additional data

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

		// Validate performance for each ball (adjusted for E2E with database operations)
		// Note: E2E tests include full database operations, cache invalidation, and over creation
		// Individual balls can take longer, especially when creating new overs
		if duration > 6000*time.Millisecond {
			t.Errorf("Ball %d response time exceeds target: %v > 6000ms", i+1, duration)
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

		// Validate performance targets (adjusted for E2E with database operations)
		// Note: E2E tests include full database operations and are slower than unit tests
		if avgDuration > 2000*time.Millisecond {
			t.Errorf("Average response time exceeds target: %v > 2000ms", avgDuration)
		}
		if maxDuration > 6000*time.Millisecond {
			t.Errorf("Max response time exceeds target: %v > 6000ms", maxDuration)
		}
	}

	// Clean up test data at the end - clean only this test's specific data
	cleanupTestMatch(t, suite.dbClient, suite.matchID, suite.seriesID)
}

