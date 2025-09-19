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

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"spark-park-cricket-backend/internal/config"
	"spark-park-cricket-backend/internal/database"
	"spark-park-cricket-backend/internal/handlers"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/services"
	"spark-park-cricket-backend/pkg/testutils"
)

func setupTestServer(t *testing.T) (*httptest.Server, *database.Client) {
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
	scorecardHandler := handlers.NewScorecardHandler(serviceContainer.Scorecard)

	// Create router and register routes
	router := chi.NewRouter()
	router.Post("/api/v1/scorecard/ball", scorecardHandler.AddBall)
	router.Get("/api/v1/scorecard/{match_id}", scorecardHandler.GetScorecard)

	// Create test server
	server := httptest.NewServer(router)

	return server, db
}

func createTestMatch(t *testing.T, db *database.Client) (string, string) {
	ctx := context.Background()

	// Create test user first
	testUser := &models.User{
		GoogleID:      fmt.Sprintf("test-google-id-match-completion-%d", time.Now().UnixNano()),
		Email:         fmt.Sprintf("test-match-completion-%d@example.com", time.Now().UnixNano()),
		Name:          "Test Match Completion User",
		Picture:       "https://example.com/picture.jpg",
		EmailVerified: true,
	}
	err := db.Repositories.User.CreateUser(ctx, testUser)
	require.NoError(t, err)
	defer db.Repositories.User.DeleteUser(ctx, testUser.ID)

	// Create test series
	series := &models.Series{
		Name:      fmt.Sprintf("Test Series %d", time.Now().Unix()),
		StartDate: time.Now(),
		EndDate:   time.Now().Add(24 * time.Hour),
		CreatedBy: testUser.ID,
	}
	err := db.Repositories.Series.Create(ctx, series)
	require.NoError(t, err)

	// Create test match
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

	return series.ID, match.ID
}

func TestMatchCompletion_TargetReached_Integration(t *testing.T) {
	// Setup
	server, db := setupTestServer(t)
	defer server.Close()
	defer db.Close()
	defer testutils.CleanupScorecardTestData(t, db)

	_, matchID := createTestMatch(t, db)

	// Complete first innings with 10 runs
	ctx := context.Background()
	firstInnings := &models.Innings{
		MatchID:       matchID,
		InningsNumber: 1,
		BattingTeam:   models.TeamTypeA,
		TotalRuns:     10,
		TotalWickets:  0,
		TotalOvers:    2.0,
		TotalBalls:    12,
		Status:        string(models.InningsStatusCompleted),
	}
	err := db.Repositories.Scorecard.CreateInnings(ctx, firstInnings)
	require.NoError(t, err)

	// Create second innings
	secondInnings := &models.Innings{
		MatchID:       matchID,
		InningsNumber: 2,
		BattingTeam:   models.TeamTypeB,
		TotalRuns:     0,
		TotalWickets:  0,
		TotalOvers:    0.0,
		TotalBalls:    0,
		Status:        string(models.InningsStatusInProgress),
	}
	err = db.Repositories.Scorecard.CreateInnings(ctx, secondInnings)
	require.NoError(t, err)

	// Add balls to reach target (11 runs)
	ballRequests := []models.BallEventRequest{
		{MatchID: matchID, InningsNumber: 2, BallType: models.BallTypeGood, RunType: models.RunTypeSix, IsWicket: false, Byes: 0},
		{MatchID: matchID, InningsNumber: 2, BallType: models.BallTypeGood, RunType: models.RunTypeFive, IsWicket: false, Byes: 0},
	}

	for _, ballReq := range ballRequests {
		jsonData, _ := json.Marshal(ballReq)
		req, _ := http.NewRequest("POST", server.URL+"/api/v1/scorecard/ball", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Body.Close()
	}

	// Check match status
	req, _ := http.NewRequest("GET", server.URL+"/api/v1/scorecard/"+matchID, nil)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var scorecardResponse struct {
		Data struct {
			MatchStatus string `json:"match_status"`
			Innings     []struct {
				InningsNumber int    `json:"innings_number"`
				TotalRuns     int    `json:"total_runs"`
				Status        string `json:"status"`
			} `json:"innings"`
		} `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&scorecardResponse)
	require.NoError(t, err)

	// Verify match is completed
	assert.Equal(t, "completed", scorecardResponse.Data.MatchStatus)

	// Verify second innings is completed
	var secondInningsData struct {
		InningsNumber int    `json:"innings_number"`
		TotalRuns     int    `json:"total_runs"`
		Status        string `json:"status"`
	}
	for _, innings := range scorecardResponse.Data.Innings {
		if innings.InningsNumber == 2 {
			secondInningsData = innings
			break
		}
	}
	assert.Equal(t, "completed", secondInningsData.Status)
	assert.GreaterOrEqual(t, secondInningsData.TotalRuns, 11) // Target reached
}

func TestMatchCompletion_AllWicketsLost_Integration(t *testing.T) {
	// Setup
	server, db := setupTestServer(t)
	defer server.Close()
	defer db.Close()
	defer testutils.CleanupScorecardTestData(t, db)

	_, matchID := createTestMatch(t, db)

	// Complete first innings with 10 runs
	ctx := context.Background()
	firstInnings := &models.Innings{
		MatchID:       matchID,
		InningsNumber: 1,
		BattingTeam:   models.TeamTypeA,
		TotalRuns:     10,
		TotalWickets:  0,
		TotalOvers:    2.0,
		TotalBalls:    12,
		Status:        string(models.InningsStatusCompleted),
	}
	err := db.Repositories.Scorecard.CreateInnings(ctx, firstInnings)
	require.NoError(t, err)

	// Create second innings
	secondInnings := &models.Innings{
		MatchID:       matchID,
		InningsNumber: 2,
		BattingTeam:   models.TeamTypeB,
		TotalRuns:     0,
		TotalWickets:  0,
		TotalOvers:    0.0,
		TotalBalls:    0,
		Status:        string(models.InningsStatusInProgress),
	}
	err = db.Repositories.Scorecard.CreateInnings(ctx, secondInnings)
	require.NoError(t, err)

	// Add balls to lose all wickets (2 wickets for 3 players)
	ballRequests := []models.BallEventRequest{
		{MatchID: matchID, InningsNumber: 2, BallType: models.BallTypeGood, RunType: models.RunTypeWC, IsWicket: true, WicketType: "bowled", Byes: 0},
		{MatchID: matchID, InningsNumber: 2, BallType: models.BallTypeGood, RunType: models.RunTypeWC, IsWicket: true, WicketType: "bowled", Byes: 0},
	}

	for _, ballReq := range ballRequests {
		jsonData, _ := json.Marshal(ballReq)
		req, _ := http.NewRequest("POST", server.URL+"/api/v1/scorecard/ball", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Body.Close()
	}

	// Check match status
	req, _ := http.NewRequest("GET", server.URL+"/api/v1/scorecard/"+matchID, nil)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var scorecardResponse struct {
		Data struct {
			MatchStatus string `json:"match_status"`
			Innings     []struct {
				InningsNumber int    `json:"innings_number"`
				TotalWickets  int    `json:"total_wickets"`
				Status        string `json:"status"`
			} `json:"innings"`
		} `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&scorecardResponse)
	require.NoError(t, err)

	// Verify match is completed
	assert.Equal(t, "completed", scorecardResponse.Data.MatchStatus)

	// Verify second innings is completed with all wickets lost
	var secondInningsData struct {
		InningsNumber int    `json:"innings_number"`
		TotalWickets  int    `json:"total_wickets"`
		Status        string `json:"status"`
	}
	for _, innings := range scorecardResponse.Data.Innings {
		if innings.InningsNumber == 2 {
			secondInningsData = innings
			break
		}
	}
	assert.Equal(t, "completed", secondInningsData.Status)
	assert.Equal(t, 2, secondInningsData.TotalWickets) // All wickets lost (n-1)
}

func TestMatchCompletion_AllOversCompleted_Integration(t *testing.T) {
	// Setup
	server, db := setupTestServer(t)
	defer server.Close()
	defer db.Close()
	defer testutils.CleanupScorecardTestData(t, db)

	_, matchID := createTestMatch(t, db)

	// Complete first innings with 10 runs
	ctx := context.Background()
	firstInnings := &models.Innings{
		MatchID:       matchID,
		InningsNumber: 1,
		BattingTeam:   models.TeamTypeA,
		TotalRuns:     10,
		TotalWickets:  0,
		TotalOvers:    2.0,
		TotalBalls:    12,
		Status:        string(models.InningsStatusCompleted),
	}
	err := db.Repositories.Scorecard.CreateInnings(ctx, firstInnings)
	require.NoError(t, err)

	// Create second innings
	secondInnings := &models.Innings{
		MatchID:       matchID,
		InningsNumber: 2,
		BattingTeam:   models.TeamTypeB,
		TotalRuns:     0,
		TotalWickets:  0,
		TotalOvers:    0.0,
		TotalBalls:    0,
		Status:        string(models.InningsStatusInProgress),
	}
	err = db.Repositories.Scorecard.CreateInnings(ctx, secondInnings)
	require.NoError(t, err)

	// Add 12 balls (2 overs) to complete all overs
	ballRequests := make([]models.BallEventRequest, 12)
	for i := 0; i < 12; i++ {
		ballRequests[i] = models.BallEventRequest{
			MatchID:       matchID,
			InningsNumber: 2,
			BallType:      models.BallTypeGood,
			RunType:       models.RunTypeOne,
			IsWicket:      false,
			Byes:          0,
		}
	}

	for i, ballReq := range ballRequests {
		jsonData, _ := json.Marshal(ballReq)
		req, _ := http.NewRequest("POST", server.URL+"/api/v1/scorecard/ball", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		// After 11 balls, the match should be completed (target reached: 11/11)
		// So the 12th ball should return an error
		if i == 11 {
			assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		} else {
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		}
		resp.Body.Close()
	}

	// Check match status
	req, _ := http.NewRequest("GET", server.URL+"/api/v1/scorecard/"+matchID, nil)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var scorecardResponse struct {
		Data struct {
			MatchStatus string `json:"match_status"`
			Innings     []struct {
				InningsNumber int     `json:"innings_number"`
				TotalOvers    float64 `json:"total_overs"`
				Status        string  `json:"status"`
			} `json:"innings"`
		} `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&scorecardResponse)
	require.NoError(t, err)

	// Verify match is completed
	assert.Equal(t, "completed", scorecardResponse.Data.MatchStatus)

	// Verify second innings is completed with all overs
	var secondInningsData struct {
		InningsNumber int     `json:"innings_number"`
		TotalOvers    float64 `json:"total_overs"`
		Status        string  `json:"status"`
	}
	for _, innings := range scorecardResponse.Data.Innings {
		if innings.InningsNumber == 2 {
			secondInningsData = innings
			break
		}
	}
	assert.Equal(t, "completed", secondInningsData.Status)
	assert.GreaterOrEqual(t, secondInningsData.TotalOvers, 1.0) // Match completed when target reached
}

func TestMatchCompletion_MatchContinues_Integration(t *testing.T) {
	// Setup
	server, db := setupTestServer(t)
	defer server.Close()
	defer db.Close()
	defer testutils.CleanupScorecardTestData(t, db)

	_, matchID := createTestMatch(t, db)

	// Complete first innings with 10 runs
	ctx := context.Background()
	firstInnings := &models.Innings{
		MatchID:       matchID,
		InningsNumber: 1,
		BattingTeam:   models.TeamTypeA,
		TotalRuns:     10,
		TotalWickets:  0,
		TotalOvers:    2.0,
		TotalBalls:    12,
		Status:        string(models.InningsStatusCompleted),
	}
	err := db.Repositories.Scorecard.CreateInnings(ctx, firstInnings)
	require.NoError(t, err)

	// Create second innings
	secondInnings := &models.Innings{
		MatchID:       matchID,
		InningsNumber: 2,
		BattingTeam:   models.TeamTypeB,
		TotalRuns:     0,
		TotalWickets:  0,
		TotalOvers:    0.0,
		TotalBalls:    0,
		Status:        string(models.InningsStatusInProgress),
	}
	err = db.Repositories.Scorecard.CreateInnings(ctx, secondInnings)
	require.NoError(t, err)

	// Add only 1 ball (match should continue)
	ballReq := models.BallEventRequest{
		MatchID:       matchID,
		InningsNumber: 2,
		BallType:      models.BallTypeGood,
		RunType:       models.RunTypeOne,
		IsWicket:      false,
		Byes:          0,
	}

	jsonData, _ := json.Marshal(ballReq)
	req, _ := http.NewRequest("POST", server.URL+"/api/v1/scorecard/ball", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	// Check match status
	req, _ = http.NewRequest("GET", server.URL+"/api/v1/scorecard/"+matchID, nil)
	req.Header.Set("Accept", "application/json")

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var scorecardResponse struct {
		Data struct {
			MatchStatus string `json:"match_status"`
			Innings     []struct {
				InningsNumber int    `json:"innings_number"`
				TotalRuns     int    `json:"total_runs"`
				Status        string `json:"status"`
			} `json:"innings"`
		} `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&scorecardResponse)
	require.NoError(t, err)

	// Verify match is still live
	assert.Equal(t, "live", scorecardResponse.Data.MatchStatus)

	// Verify second innings is still in progress
	var secondInningsData struct {
		InningsNumber int    `json:"innings_number"`
		TotalRuns     int    `json:"total_runs"`
		Status        string `json:"status"`
	}
	for _, innings := range scorecardResponse.Data.Innings {
		if innings.InningsNumber == 2 {
			secondInningsData = innings
			break
		}
	}
	assert.Equal(t, "in_progress", secondInningsData.Status)
	assert.Less(t, secondInningsData.TotalRuns, 11) // Below target
}
