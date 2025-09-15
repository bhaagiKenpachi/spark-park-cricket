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
)

func TestIllegalBalls_Comprehensive_Scenario(t *testing.T) {
	// Setup
	cfg := config.LoadTestConfig()

	db, err := database.NewTestClient(cfg)
	require.NoError(t, err)
	defer db.Close()

	// Setup test schema
	err = database.SetupTestSchema(cfg)
	require.NoError(t, err)

	serviceContainer := services.NewContainer(db.Repositories)
	scorecardHandler := handlers.NewScorecardHandler(serviceContainer.Scorecard)

	router := http.NewServeMux()
	router.HandleFunc("/api/v1/scorecard/ball", scorecardHandler.AddBall)
	router.HandleFunc("/api/v1/scorecard/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if len(path) > len("/api/v1/scorecard/") {
			matchID := path[len("/api/v1/scorecard/"):]
			// Create a custom handler that extracts match_id from URL
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
		} else {
			http.NotFound(w, r)
		}
	})

	server := httptest.NewServer(router)
	defer server.Close()

	// Create test match
	ctx := context.Background()
	series := &models.Series{
		Name:      fmt.Sprintf("Illegal Balls Test %d", time.Now().Unix()),
		StartDate: time.Now(),
		EndDate:   time.Now().Add(24 * time.Hour),
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
	}
	err = db.Repositories.Match.Create(ctx, match)
	require.NoError(t, err)

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
		// Ball 7: Good ball - 5 runs
		{MatchID: match.ID, InningsNumber: 1, BallType: models.BallTypeGood, RunType: models.RunTypeFive, IsWicket: false, Byes: 0},
		// Ball 8: Good ball - 6 runs
		{MatchID: match.ID, InningsNumber: 1, BallType: models.BallTypeGood, RunType: models.RunTypeSix, IsWicket: false, Byes: 0},
	}

	// Add all balls
	for i, ballReq := range illegalBallsScenario {
		ballJSON, _ := json.Marshal(ballReq)
		req, _ := http.NewRequest("POST", server.URL+"/api/v1/scorecard/ball", bytes.NewBuffer(ballJSON))
		req.Header.Set("Content-Type", "application/json")

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
	firstInnings := scorecardResponse.Data.Innings[0]

	// Verify total runs: 1 (no_ball) + 1 (wide) + 1+2+3+4+5+6 (good balls) + 5 (byes) = 28
	assert.Equal(t, 28, firstInnings.TotalRuns)

	// Verify total overs: 1.1 (1 completed over + 1 ball in second over)
	assert.Equal(t, 1.1, firstInnings.TotalOvers)

	// Verify total balls: 6 (only legal balls count towards over completion)
	assert.Equal(t, 6, firstInnings.TotalBalls)

	// Verify extras
	assert.Equal(t, 5, firstInnings.Extras.Byes)    // From no_ball with byes
	assert.Equal(t, 1, firstInnings.Extras.Wides)   // From wide ball
	assert.Equal(t, 1, firstInnings.Extras.NoBalls) // From no_ball
	assert.Equal(t, 7, firstInnings.Extras.Total)   // 5 + 1 + 1

	// Verify over data - first over should be completed with 5 legal balls (6th legal ball goes to second over)
	firstOver := firstInnings.Overs[0]
	assert.Equal(t, 1, firstOver.OverNumber)
	assert.Equal(t, 22, firstOver.TotalRuns) // 6+1+1+2+3+4+5 = 22 runs from all balls in first over
	assert.Equal(t, 5, firstOver.TotalBalls) // Only legal balls count for over completion
	assert.Equal(t, "completed", firstOver.Status)

	// Verify ball details - first over should have 7 balls (2 illegal + 5 legal)
	assert.Len(t, firstOver.Balls, 7, "First over should have 7 balls (2 illegal + 5 legal)")

	// First over should contain all 7 balls (2 illegal + 5 legal)
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

	// Ball 7: Good ball - 5 runs
	assert.Equal(t, 7, firstOver.Balls[6].BallNumber)
	assert.Equal(t, "good", firstOver.Balls[6].BallType)
	assert.Equal(t, "5", firstOver.Balls[6].RunType)
	assert.Equal(t, 5, firstOver.Balls[6].Runs)

	// Verify second over exists and contains the 6th legal ball
	assert.Len(t, firstInnings.Overs, 2, "Should have 2 overs")
	secondOver := firstInnings.Overs[1]
	assert.Equal(t, 2, secondOver.OverNumber)
	assert.Equal(t, "in_progress", secondOver.Status)
	assert.Len(t, secondOver.Balls, 1, "Second over should have 1 legal ball")

	// Ball 1: Good ball - 6 runs (should be in second over, ball numbering resets per over)
	assert.Equal(t, 1, secondOver.Balls[0].BallNumber)
	assert.Equal(t, "good", secondOver.Balls[0].BallType)
	assert.Equal(t, "6", secondOver.Balls[0].RunType)
	assert.Equal(t, 6, secondOver.Balls[0].Runs)
	assert.Equal(t, 0, secondOver.Balls[0].Byes)
	assert.False(t, secondOver.Balls[0].IsWicket)

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
	assert.Equal(t, 1.2, firstInnings.TotalOvers) // 1 completed over + 2 balls in current over
	assert.Len(t, firstInnings.Overs, 2)          // Two overs now

	// Verify second over is in progress
	secondOverAfterBall := firstInnings.Overs[1]
	assert.Equal(t, 2, secondOverAfterBall.OverNumber)
	assert.Equal(t, "in_progress", secondOverAfterBall.Status)
	assert.Len(t, secondOverAfterBall.Balls, 2) // Two balls in second over
}

func TestIllegalBalls_OverCompletion_Logic(t *testing.T) {
	// Setup
	cfg := config.LoadTestConfig()

	db, err := database.NewTestClient(cfg)
	require.NoError(t, err)
	defer db.Close()

	// Setup test schema
	err = database.SetupTestSchema(cfg)
	require.NoError(t, err)

	serviceContainer := services.NewContainer(db.Repositories)
	scorecardHandler := handlers.NewScorecardHandler(serviceContainer.Scorecard)

	router := http.NewServeMux()
	router.HandleFunc("/api/v1/scorecard/ball", scorecardHandler.AddBall)
	router.HandleFunc("/api/v1/scorecard/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if len(path) > len("/api/v1/scorecard/") {
			matchID := path[len("/api/v1/scorecard/"):]
			// Create a custom handler that extracts match_id from URL
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
		} else {
			http.NotFound(w, r)
		}
	})

	server := httptest.NewServer(router)
	defer server.Close()

	// Create test match
	ctx := context.Background()
	series := &models.Series{
		Name:      fmt.Sprintf("Over Completion Test %d", time.Now().Unix()),
		StartDate: time.Now(),
		EndDate:   time.Now().Add(24 * time.Hour),
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
	}
	err = db.Repositories.Match.Create(ctx, match)
	require.NoError(t, err)

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
	firstInnings := scorecardResponse.Data.Innings[0]
	assert.Equal(t, 0.5, firstInnings.TotalOvers) // 5 legal balls = 0.5 overs

	firstOver := firstInnings.Overs[0]
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

	// Verify over is now complete
	firstInnings = scorecardResponse.Data.Innings[0]
	assert.Equal(t, 1.0, firstInnings.TotalOvers) // 6 legal balls = 1.0 overs

	firstOver = firstInnings.Overs[0]
	assert.Equal(t, 6, firstOver.TotalBalls)       // 6 legal balls
	assert.Equal(t, "completed", firstOver.Status) // Over complete
}
