package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"spark-park-cricket-backend/internal/config"
	"spark-park-cricket-backend/internal/database"
	"spark-park-cricket-backend/internal/handlers"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/services"
	"spark-park-cricket-backend/pkg/testutils"
)

func TestIllegalBalls_Comprehensive_Scenario(t *testing.T) {
	// Setup
	cfg := config.LoadTestConfig()
	var server *httptest.Server

	db, err := database.NewTestClient(cfg)
	require.NoError(t, err)
	defer db.Close()

	// Setup test schema
	err = database.SetupTestSchema(cfg)
	require.NoError(t, err)

	// Use the standard router setup that includes authentication middleware
	router := handlers.SetupRoutes(db, cfg.Config)

	server = httptest.NewServer(router)
	defer server.Close()

	// Create test user first
	ctx := context.Background()
	testUser := &models.User{
		GoogleID:      fmt.Sprintf("test-google-id-illegal-%d", time.Now().UnixNano()),
		Email:         fmt.Sprintf("test-illegal-%d@example.com", time.Now().UnixNano()),
		Name:          "Test Illegal Balls User",
		Picture:       "https://example.com/picture.jpg",
		EmailVerified: true,
	}
	err = db.Repositories.User.CreateUser(ctx, testUser)
	require.NoError(t, err)
	defer func() { _ = db.Repositories.User.DeleteUser(ctx, testUser.ID) }()

	// Create user session for authentication
	serviceContainer := services.NewContainer(db, cfg.Config)
	sessionService := serviceContainer.SessionService

	mockReq := httptest.NewRequest("GET", "/", nil)
	mockWriter := httptest.NewRecorder()

	err = sessionService.CreateSession(mockWriter, mockReq, testUser)
	require.NoError(t, err)

	// Extract the session cookie
	cookies := mockWriter.Result().Cookies()
	var sessionCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "user_session" {
			sessionCookie = cookie
			break
		}
	}
	require.NotNil(t, sessionCookie, "Session cookie should be created")

	// Create test match
	series := &models.Series{
		Name:      fmt.Sprintf("Illegal Balls Test %d", time.Now().Unix()),
		StartDate: time.Now(),
		EndDate:   time.Now().Add(24 * time.Hour),
		CreatedBy: testUser.ID,
	}
	err = db.Repositories.Series.Create(ctx, series)
	require.NoError(t, err)

	match := &models.Match{
		SeriesID:         series.ID,
		MatchNumber:      1,
		Date:             time.Now(),
		Status:           models.MatchStatusLive,
		TeamAPlayerCount: 3,
		TeamBPlayerCount: 3,
		TotalOvers:       2,
		TossWinner:       models.TeamTypeA,
		TossType:         models.TossTypeHeads,
		BattingTeam:      models.TeamTypeA,
		CreatedBy:        testUser.ID,
	}
	err = db.Repositories.Match.Create(ctx, match)
	require.NoError(t, err)

	// Start scoring for the match
	startScoringReq := map[string]interface{}{
		"match_id": match.ID,
	}
	startScoringBody, _ := json.Marshal(startScoringReq)
	startScoringHTTPReq := httptest.NewRequest("POST", "/api/v1/scorecard/start", bytes.NewBuffer(startScoringBody))
	startScoringHTTPReq.Header.Set("Content-Type", "application/json")
	startScoringHTTPReq.AddCookie(sessionCookie)
	startScoringResp := httptest.NewRecorder()

	router.ServeHTTP(startScoringResp, startScoringHTTPReq)
	require.Equal(t, http.StatusOK, startScoringResp.Code, "Should start scoring successfully")

	// Test scenario: Over with illegal balls
	// Expected: 1 no_ball + 1 wide + 6 good balls = 8 total balls, but only 6 legal balls
	illegalBallsScenario := []models.BallEventRequest{
		// Ball 1: No ball with 5 byes
		{MatchID: match.ID, InningsNumber: 1, BallType: models.BallTypeNoBall, RunType: models.RunTypeNB, IsWicket: false, Byes: 5},
		// Ball 2: Wide ball
		{MatchID: match.ID, InningsNumber: 1, BallType: models.BallTypeWide, RunType: models.RunTypeWD, IsWicket: false, Byes: 0},
		// Ball 3: Good ball - 1 run
		{MatchID: match.ID, InningsNumber: 1, BallType: models.BallTypeGood, RunType: models.RunTypeOne, IsWicket: false, Byes: 0},
		// Ball 4: Good ball - 2 runs
		{MatchID: match.ID, InningsNumber: 1, BallType: models.BallTypeGood, RunType: models.RunTypeTwo, IsWicket: false, Byes: 0},
		// Ball 5: Good ball - 3 runs
		{MatchID: match.ID, InningsNumber: 1, BallType: models.BallTypeGood, RunType: models.RunTypeThree, IsWicket: false, Byes: 0},
		// Ball 6: Good ball - 4 runs
		{MatchID: match.ID, InningsNumber: 1, BallType: models.BallTypeGood, RunType: models.RunTypeFour, IsWicket: false, Byes: 0},
		// Ball 7: Good ball - 1 run (changed from 5 runs which is invalid in cricket)
		{MatchID: match.ID, InningsNumber: 1, BallType: models.BallTypeGood, RunType: models.RunTypeOne, IsWicket: false, Byes: 0},
		// Ball 8: Good ball - 6 runs
		{MatchID: match.ID, InningsNumber: 1, BallType: models.BallTypeGood, RunType: models.RunTypeSix, IsWicket: false, Byes: 0},
	}

	// Add all balls
	for i, ballReq := range illegalBallsScenario {
		ballJSON, _ := json.Marshal(ballReq)
		req, _ := http.NewRequest("POST", server.URL+"/api/v1/scorecard/ball", bytes.NewBuffer(ballJSON))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(sessionCookie) // Add authentication

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Ball %d should be added successfully", i+1)
		resp.Body.Close()
	}

	// Check scorecard
	req, _ := http.NewRequest("GET", server.URL+"/api/v1/scorecard/"+match.ID, nil)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var scorecardResponse struct {
		Data struct {
			Innings []struct {
				InningsNumber int     `json:"innings_number"`
				TotalRuns     int     `json:"total_runs"`
				TotalOvers    float64 `json:"total_overs"`
				TotalBalls    int     `json:"total_balls"`
				Status        string  `json:"status"`
				Extras        struct {
					Byes    int `json:"byes"`
					Wides   int `json:"wides"`
					NoBalls int `json:"no_balls"`
					Total   int `json:"total"`
				} `json:"extras"`
				Overs []struct {
					OverNumber int    `json:"over_number"`
					TotalRuns  int    `json:"total_runs"`
					TotalBalls int    `json:"total_balls"`
					Status     string `json:"status"`
					Balls      []struct {
						BallNumber int    `json:"ball_number"`
						BallType   string `json:"ball_type"`
						RunType    string `json:"run_type"`
						Runs       int    `json:"runs"`
						Byes       int    `json:"byes"`
						IsWicket   bool   `json:"is_wicket"`
					} `json:"balls"`
				} `json:"overs"`
			} `json:"innings"`
		} `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&scorecardResponse)
	require.NoError(t, err)
	resp.Body.Close()

	// Verify first innings data
	require.NotEmpty(t, scorecardResponse.Data.Innings, "Innings should not be empty")
	require.Len(t, scorecardResponse.Data.Innings, 1, "Should have exactly one innings")
	require.NotNil(t, scorecardResponse.Data.Innings[0], "First innings should not be nil")
	firstInnings := scorecardResponse.Data.Innings[0]

	// Verify total runs: 1 (no_ball) + 1 (wide) + 1+2+3+4+1+6 (good balls) + 5 (byes) = 24
	assert.Equal(t, 24, firstInnings.TotalRuns)

	// Verify total overs: 1.0 (1 completed over, no second over created yet)
	assert.Equal(t, 1.0, firstInnings.TotalOvers)

	// Verify total balls: 6 (only legal balls count towards over completion)
	assert.Equal(t, 6, firstInnings.TotalBalls)

	// Verify extras
	assert.Equal(t, 5, firstInnings.Extras.Byes)    // From no_ball with byes
	assert.Equal(t, 1, firstInnings.Extras.Wides)   // From wide ball
	assert.Equal(t, 1, firstInnings.Extras.NoBalls) // From no_ball
	assert.Equal(t, 7, firstInnings.Extras.Total)   // 5 + 1 + 1

	// Verify over data - first over should be completed with 6 legal balls
	require.NotEmpty(t, firstInnings.Overs, "First innings should have overs")
	require.Len(t, firstInnings.Overs, 1, "Should have exactly one over")
	require.NotNil(t, firstInnings.Overs[0], "First over should not be nil")
	firstOver := firstInnings.Overs[0]
	assert.Equal(t, 1, firstOver.OverNumber)
	assert.Equal(t, 24, firstOver.TotalRuns) // 1+1+1+2+3+4+1+6+5 = 24 runs from all balls in first over
	assert.Equal(t, 6, firstOver.TotalBalls) // All 6 legal balls in first over
	assert.Equal(t, "completed", firstOver.Status)

	// Verify ball details - first over should have 8 balls (2 illegal + 6 legal)
	assert.Len(t, firstOver.Balls, 8, "First over should have 8 balls (2 illegal + 6 legal)")

	// First over should contain all 8 balls (2 illegal + 6 legal)
	require.NotEmpty(t, firstOver.Balls, "First over should have balls")
	require.Len(t, firstOver.Balls, 8, "First over should have exactly 8 balls")

	// Defensive checks for ball access
	require.NotNil(t, firstOver.Balls[0], "Ball 1 should not be nil")
	require.NotNil(t, firstOver.Balls[1], "Ball 2 should not be nil")
	require.NotNil(t, firstOver.Balls[2], "Ball 3 should not be nil")
	require.NotNil(t, firstOver.Balls[3], "Ball 4 should not be nil")

	// Ball 1: No ball with byes
	assert.Equal(t, 1, firstOver.Balls[0].BallNumber)
	assert.Equal(t, "no_ball", firstOver.Balls[0].BallType)
	assert.Equal(t, "NB", firstOver.Balls[0].RunType)
	assert.Equal(t, 1, firstOver.Balls[0].Runs)
	assert.Equal(t, 5, firstOver.Balls[0].Byes)
	assert.False(t, firstOver.Balls[0].IsWicket)

	// Ball 2: Wide ball
	assert.Equal(t, 2, firstOver.Balls[1].BallNumber)
	assert.Equal(t, "wide", firstOver.Balls[1].BallType)
	assert.Equal(t, "WD", firstOver.Balls[1].RunType)
	assert.Equal(t, 1, firstOver.Balls[1].Runs)
	assert.Equal(t, 0, firstOver.Balls[1].Byes)
	assert.False(t, firstOver.Balls[1].IsWicket)

	// Ball 3: Good ball - 1 run
	assert.Equal(t, 3, firstOver.Balls[2].BallNumber)
	assert.Equal(t, "good", firstOver.Balls[2].BallType)
	assert.Equal(t, "1", firstOver.Balls[2].RunType)
	assert.Equal(t, 1, firstOver.Balls[2].Runs)
	assert.Equal(t, 0, firstOver.Balls[2].Byes)
	assert.False(t, firstOver.Balls[2].IsWicket)

	// Ball 4: Good ball - 2 runs
	assert.Equal(t, 4, firstOver.Balls[3].BallNumber)
	assert.Equal(t, "good", firstOver.Balls[3].BallType)
	assert.Equal(t, "2", firstOver.Balls[3].RunType)
	assert.Equal(t, 2, firstOver.Balls[3].Runs)

	// Ball 5: Good ball - 3 runs
	assert.Equal(t, 5, firstOver.Balls[4].BallNumber)
	assert.Equal(t, "good", firstOver.Balls[4].BallType)
	assert.Equal(t, "3", firstOver.Balls[4].RunType)
	assert.Equal(t, 3, firstOver.Balls[4].Runs)

	// Ball 6: Good ball - 4 runs
	assert.Equal(t, 6, firstOver.Balls[5].BallNumber)
	assert.Equal(t, "good", firstOver.Balls[5].BallType)
	assert.Equal(t, "4", firstOver.Balls[5].RunType)
	assert.Equal(t, 4, firstOver.Balls[5].Runs)

	// Ball 7: Good ball - 1 run
	assert.Equal(t, 7, firstOver.Balls[6].BallNumber)
	assert.Equal(t, "good", firstOver.Balls[6].BallType)
	assert.Equal(t, "1", firstOver.Balls[6].RunType)
	assert.Equal(t, 1, firstOver.Balls[6].Runs)
	assert.Equal(t, 0, firstOver.Balls[6].Byes)
	assert.False(t, firstOver.Balls[6].IsWicket)

	// Ball 8: Good ball - 6 runs
	assert.Equal(t, 8, firstOver.Balls[7].BallNumber)
	assert.Equal(t, "good", firstOver.Balls[7].BallType)
	assert.Equal(t, "6", firstOver.Balls[7].RunType)
	assert.Equal(t, 6, firstOver.Balls[7].Runs)
	assert.Equal(t, 0, firstOver.Balls[7].Byes)
	assert.False(t, firstOver.Balls[7].IsWicket)

	// Verify only one over exists (all balls in first over)
	assert.Len(t, firstInnings.Overs, 1, "Should have 1 over")

	// Test adding one more ball to start second over
	nextBallReq := models.BallEventRequest{
		MatchID:       match.ID,
		InningsNumber: 1,
		BallType:      models.BallTypeGood,
		RunType:       models.RunTypeOne,
		IsWicket:      false,
		Byes:          0,
	}

	ballJSON, _ := json.Marshal(nextBallReq)
	req, _ = http.NewRequest("POST", server.URL+"/api/v1/scorecard/ball", bytes.NewBuffer(ballJSON))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(sessionCookie) // Add authentication

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	// Check scorecard again to verify second over started
	req, _ = http.NewRequest("GET", server.URL+"/api/v1/scorecard/"+match.ID, nil)
	req.Header.Set("Accept", "application/json")

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	err = json.NewDecoder(resp.Body).Decode(&scorecardResponse)
	require.NoError(t, err)
	resp.Body.Close()

	// Verify second over started
	firstInnings = scorecardResponse.Data.Innings[0]
	assert.Equal(t, 1.1, firstInnings.TotalOvers) // 1 completed over + 1 ball in current over
	assert.Len(t, firstInnings.Overs, 2)          // Two overs now

	// Verify second over is in progress
	secondOverAfterBall := firstInnings.Overs[1]
	assert.Equal(t, 2, secondOverAfterBall.OverNumber)
	assert.Equal(t, "in_progress", secondOverAfterBall.Status)
	assert.Len(t, secondOverAfterBall.Balls, 1) // One ball in second over
}

func TestIllegalBalls_OverCompletion_Logic(t *testing.T) {
	// Setup
	cfg := config.LoadTestConfig()
	var server *httptest.Server

	db, err := database.NewTestClient(cfg)
	require.NoError(t, err)
	defer db.Close()

	// Setup test schema
	err = database.SetupTestSchema(cfg)
	require.NoError(t, err)

	// Clean up any existing test data BEFORE creating new test data
	testutils.CleanupScorecardTestData(t, db)

	// Use the standard router setup that includes authentication middleware
	router := handlers.SetupRoutes(db, cfg.Config)

	server = httptest.NewServer(router)
	defer server.Close()

	// Create test user first
	ctx := context.Background()
	testUser := &models.User{
		GoogleID:      fmt.Sprintf("test-google-id-illegal-over-%d", time.Now().UnixNano()),
		Email:         fmt.Sprintf("test-illegal-over-%d@example.com", time.Now().UnixNano()),
		Name:          "Test Illegal Balls Over User",
		Picture:       "https://example.com/picture.jpg",
		EmailVerified: true,
	}
	err = db.Repositories.User.CreateUser(ctx, testUser)
	require.NoError(t, err)
	defer func() { _ = db.Repositories.User.DeleteUser(ctx, testUser.ID) }()

	// Create user session for authentication
	serviceContainer := services.NewContainer(db, cfg.Config)
	sessionService := serviceContainer.SessionService

	mockReq := httptest.NewRequest("GET", "/", nil)
	mockWriter := httptest.NewRecorder()

	err = sessionService.CreateSession(mockWriter, mockReq, testUser)
	require.NoError(t, err)

	// Extract the session cookie
	cookies := mockWriter.Result().Cookies()
	var sessionCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "user_session" {
			sessionCookie = cookie
			break
		}
	}
	require.NotNil(t, sessionCookie, "Session cookie should be created")

	// Create test match
	series := &models.Series{
		Name:      fmt.Sprintf("Over Completion Test %d", time.Now().Unix()),
		StartDate: time.Now(),
		EndDate:   time.Now().Add(24 * time.Hour),
		CreatedBy: testUser.ID,
	}
	err = db.Repositories.Series.Create(ctx, series)
	require.NoError(t, err)

	match := &models.Match{
		SeriesID:         series.ID,
		MatchNumber:      1,
		Date:             time.Now(),
		Status:           models.MatchStatusLive,
		TeamAPlayerCount: 3,
		TeamBPlayerCount: 3,
		TotalOvers:       2,
		TossWinner:       models.TeamTypeA,
		TossType:         models.TossTypeHeads,
		BattingTeam:      models.TeamTypeA,
		CreatedBy:        testUser.ID,
	}
	err = db.Repositories.Match.Create(ctx, match)
	require.NoError(t, err)

	// Start scoring for the match
	startScoringReq := map[string]interface{}{
		"match_id": match.ID,
	}
	startScoringBody, _ := json.Marshal(startScoringReq)
	startScoringHTTPReq := httptest.NewRequest("POST", "/api/v1/scorecard/start", bytes.NewBuffer(startScoringBody))
	startScoringHTTPReq.Header.Set("Content-Type", "application/json")
	startScoringHTTPReq.AddCookie(sessionCookie)
	startScoringResp := httptest.NewRecorder()

	router.ServeHTTP(startScoringResp, startScoringHTTPReq)
	require.Equal(t, http.StatusOK, startScoringResp.Code, "Should start scoring successfully")

	// Test scenario: 5 legal balls + 3 illegal balls = 8 total balls, but over should not be complete
	overCompletionBalls := []models.BallEventRequest{
		// 5 legal balls
		{MatchID: match.ID, InningsNumber: 1, BallType: models.BallTypeGood, RunType: models.RunTypeOne, IsWicket: false, Byes: 0},
		{MatchID: match.ID, InningsNumber: 1, BallType: models.BallTypeGood, RunType: models.RunTypeOne, IsWicket: false, Byes: 0},
		{MatchID: match.ID, InningsNumber: 1, BallType: models.BallTypeGood, RunType: models.RunTypeOne, IsWicket: false, Byes: 0},
		{MatchID: match.ID, InningsNumber: 1, BallType: models.BallTypeGood, RunType: models.RunTypeOne, IsWicket: false, Byes: 0},
		{MatchID: match.ID, InningsNumber: 1, BallType: models.BallTypeGood, RunType: models.RunTypeOne, IsWicket: false, Byes: 0},
		// 3 illegal balls
		{MatchID: match.ID, InningsNumber: 1, BallType: models.BallTypeWide, RunType: models.RunTypeWD, IsWicket: false, Byes: 0},
		{MatchID: match.ID, InningsNumber: 1, BallType: models.BallTypeNoBall, RunType: models.RunTypeNB, IsWicket: false, Byes: 0},
		{MatchID: match.ID, InningsNumber: 1, BallType: models.BallTypeWide, RunType: models.RunTypeWD, IsWicket: false, Byes: 0},
	}

	// Add all balls
	for i, ballReq := range overCompletionBalls {
		ballJSON, _ := json.Marshal(ballReq)
		req, _ := http.NewRequest("POST", server.URL+"/api/v1/scorecard/ball", bytes.NewBuffer(ballJSON))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(sessionCookie) // Add authentication

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Ball %d should be added successfully", i+1)
		resp.Body.Close()
	}

	// Check scorecard
	req, _ := http.NewRequest("GET", server.URL+"/api/v1/scorecard/"+match.ID, nil)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var scorecardResponse struct {
		Data struct {
			MatchID string `json:"match_id"`
			Innings []struct {
				TotalOvers float64 `json:"total_overs"`
				Overs      []struct {
					OverNumber int    `json:"over_number"`
					TotalBalls int    `json:"total_balls"`
					Status     string `json:"status"`
				} `json:"overs"`
			} `json:"innings"`
		} `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&scorecardResponse)
	require.NoError(t, err)
	resp.Body.Close()

	// Verify over is not complete (only 5 legal balls)
	require.NotEmpty(t, scorecardResponse.Data.Innings, "Scorecard should have innings data")
	require.Len(t, scorecardResponse.Data.Innings, 1, "Should have exactly one innings")
	require.NotNil(t, scorecardResponse.Data.Innings[0], "First innings should not be nil")
	firstInnings := scorecardResponse.Data.Innings[0]

	require.NotEmpty(t, firstInnings.Overs, "First innings should have overs")
	require.Len(t, firstInnings.Overs, 1, "Should have exactly one over")
	require.NotNil(t, firstInnings.Overs[0], "First over should not be nil")
	firstOver := firstInnings.Overs[0]

	assert.Equal(t, 0.5, firstInnings.TotalOvers)    // 5 legal balls = 0.5 overs
	assert.Equal(t, 5, firstOver.TotalBalls)         // Only legal balls count
	assert.Equal(t, "in_progress", firstOver.Status) // Over not complete

	// Add one more legal ball to complete the over
	completingBallReq := models.BallEventRequest{
		MatchID:       match.ID,
		InningsNumber: 1,
		BallType:      models.BallTypeGood,
		RunType:       models.RunTypeOne,
		IsWicket:      false,
		Byes:          0,
	}

	ballJSON, _ := json.Marshal(completingBallReq)
	req, _ = http.NewRequest("POST", server.URL+"/api/v1/scorecard/ball", bytes.NewBuffer(ballJSON))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(sessionCookie) // Add authentication

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	// Check scorecard again
	req, _ = http.NewRequest("GET", server.URL+"/api/v1/scorecard/"+match.ID, nil)
	req.Header.Set("Accept", "application/json")

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	err = json.NewDecoder(resp.Body).Decode(&scorecardResponse)
	require.NoError(t, err)
	resp.Body.Close()

	// Verify over is now complete and second over started
	firstInnings = scorecardResponse.Data.Innings[0]
	assert.Equal(t, 1.0, firstInnings.TotalOvers) // 1 completed over, no second over created yet

	// Should have 1 over now (all balls in first over)
	assert.Len(t, firstInnings.Overs, 1, "Should have 1 over with all balls")

	firstOver = firstInnings.Overs[0]
	assert.Equal(t, 6, firstOver.TotalBalls)       // 6 legal balls (all in first over)
	assert.Equal(t, "completed", firstOver.Status) // Over complete
}
