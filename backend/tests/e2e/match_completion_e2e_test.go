package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/pkg/testutils"
)

func TestCompleteMatchFlow_TargetReached_E2E(t *testing.T) {
	// Setup
	server, db := testutils.SetupE2ETestServerWithDB(t)
	defer server.Close()
	defer db.Close()

	// Step 1: Create Series
	seriesReq := models.CreateSeriesRequest{
		Name:      fmt.Sprintf("E2E Test Series %d", time.Now().Unix()),
		StartDate: time.Now(),
		EndDate:   time.Now().Add(24 * time.Hour),
	}
	seriesJSON, _ := json.Marshal(seriesReq)
	req, _ := http.NewRequest("POST", server.URL+"/api/v1/series", bytes.NewBuffer(seriesJSON))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var seriesResponse struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&seriesResponse)
	require.NoError(t, err)
	resp.Body.Close()

	seriesID := seriesResponse.Data.ID
	require.NotEmpty(t, seriesID)

	// Step 2: Create Match
	matchNumber := 1
	matchReq := models.CreateMatchRequest{
		SeriesID:         seriesID,
		MatchNumber:      &matchNumber,
		Date:             time.Now(),
		TeamAPlayerCount: 3,
		TeamBPlayerCount: 3,
		TotalOvers:       2,
		TossWinner:       models.TeamTypeA,
		TossType:         models.TossTypeHeads,
	}
	matchJSON, _ := json.Marshal(matchReq)
	req, _ = http.NewRequest("POST", server.URL+"/api/v1/matches", bytes.NewBuffer(matchJSON))
	req.Header.Set("Content-Type", "application/json")

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var matchResponse struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&matchResponse)
	require.NoError(t, err)
	resp.Body.Close()

	matchID := matchResponse.Data.ID
	require.NotEmpty(t, matchID)

	// Step 3: Complete First Innings (12 balls = 2 overs)
	firstInningsBalls := []models.BallEventRequest{
		{MatchID: matchID, InningsNumber: 1, BallType: models.BallTypeGood, RunType: models.RunTypeOne, IsWicket: false, Byes: 0},
		{MatchID: matchID, InningsNumber: 1, BallType: models.BallTypeGood, RunType: models.RunTypeTwo, IsWicket: false, Byes: 0},
		{MatchID: matchID, InningsNumber: 1, BallType: models.BallTypeGood, RunType: models.RunTypeThree, IsWicket: false, Byes: 0},
		{MatchID: matchID, InningsNumber: 1, BallType: models.BallTypeGood, RunType: models.RunTypeFour, IsWicket: false, Byes: 0},
		{MatchID: matchID, InningsNumber: 1, BallType: models.BallTypeGood, RunType: models.RunTypeFive, IsWicket: false, Byes: 0},
		{MatchID: matchID, InningsNumber: 1, BallType: models.BallTypeGood, RunType: models.RunTypeSix, IsWicket: false, Byes: 0},
		{MatchID: matchID, InningsNumber: 1, BallType: models.BallTypeGood, RunType: models.RunTypeOne, IsWicket: false, Byes: 0},
		{MatchID: matchID, InningsNumber: 1, BallType: models.BallTypeGood, RunType: models.RunTypeTwo, IsWicket: false, Byes: 0},
		{MatchID: matchID, InningsNumber: 1, BallType: models.BallTypeGood, RunType: models.RunTypeThree, IsWicket: false, Byes: 0},
		{MatchID: matchID, InningsNumber: 1, BallType: models.BallTypeGood, RunType: models.RunTypeFour, IsWicket: false, Byes: 0},
		{MatchID: matchID, InningsNumber: 1, BallType: models.BallTypeGood, RunType: models.RunTypeFive, IsWicket: false, Byes: 0},
		{MatchID: matchID, InningsNumber: 1, BallType: models.BallTypeGood, RunType: models.RunTypeSix, IsWicket: false, Byes: 0},
	}

	for _, ballReq := range firstInningsBalls {
		ballJSON, _ := json.Marshal(ballReq)
		req, _ := http.NewRequest("POST", server.URL+"/api/v1/scorecard/ball", bytes.NewBuffer(ballJSON))
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Body.Close()
	}

	// Step 4: Check First Innings Completed and Second Innings Started
	req, _ = http.NewRequest("GET", server.URL+"/api/v1/scorecard/"+matchID, nil)
	req.Header.Set("Accept", "application/json")

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var scorecardResponse struct {
		Data struct {
			MatchStatus    string `json:"match_status"`
			CurrentInnings int    `json:"current_innings"`
			Innings        []struct {
				InningsNumber int    `json:"innings_number"`
				BattingTeam   string `json:"batting_team"`
				TotalRuns     int    `json:"total_runs"`
				Status        string `json:"status"`
			} `json:"innings"`
		} `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&scorecardResponse)
	require.NoError(t, err)
	resp.Body.Close()

	// Verify first innings completed and second innings started
	assert.Equal(t, "live", scorecardResponse.Data.MatchStatus)
	assert.Equal(t, 2, scorecardResponse.Data.CurrentInnings)

	var firstInnings, secondInnings struct {
		InningsNumber int    `json:"innings_number"`
		BattingTeam   string `json:"batting_team"`
		TotalRuns     int    `json:"total_runs"`
		Status        string `json:"status"`
	}
	for _, innings := range scorecardResponse.Data.Innings {
		if innings.InningsNumber == 1 {
			firstInnings = innings
		} else if innings.InningsNumber == 2 {
			secondInnings = innings
		}
	}

	assert.Equal(t, "completed", firstInnings.Status)
	assert.Equal(t, "A", firstInnings.BattingTeam)
	assert.Equal(t, 42, firstInnings.TotalRuns) // 1+2+3+4+5+6+1+2+3+4+5+6 = 42

	assert.Equal(t, "in_progress", secondInnings.Status)
	assert.Equal(t, "B", secondInnings.BattingTeam)
	assert.Equal(t, 0, secondInnings.TotalRuns)

	// Step 5: Add balls to second innings to reach target (43 runs)
	secondInningsBalls := []models.BallEventRequest{
		{MatchID: matchID, InningsNumber: 2, BallType: models.BallTypeGood, RunType: models.RunTypeSix, IsWicket: false, Byes: 0},
		{MatchID: matchID, InningsNumber: 2, BallType: models.BallTypeGood, RunType: models.RunTypeSix, IsWicket: false, Byes: 0},
		{MatchID: matchID, InningsNumber: 2, BallType: models.BallTypeGood, RunType: models.RunTypeSix, IsWicket: false, Byes: 0},
		{MatchID: matchID, InningsNumber: 2, BallType: models.BallTypeGood, RunType: models.RunTypeSix, IsWicket: false, Byes: 0},
		{MatchID: matchID, InningsNumber: 2, BallType: models.BallTypeGood, RunType: models.RunTypeSix, IsWicket: false, Byes: 0},
		{MatchID: matchID, InningsNumber: 2, BallType: models.BallTypeGood, RunType: models.RunTypeSix, IsWicket: false, Byes: 0},
		{MatchID: matchID, InningsNumber: 2, BallType: models.BallTypeGood, RunType: models.RunTypeSix, IsWicket: false, Byes: 0},
		{MatchID: matchID, InningsNumber: 2, BallType: models.BallTypeGood, RunType: models.RunTypeOne, IsWicket: false, Byes: 0}, // 43rd run
	}

	for _, ballReq := range secondInningsBalls {
		ballJSON, _ := json.Marshal(ballReq)
		req, _ := http.NewRequest("POST", server.URL+"/api/v1/scorecard/ball", bytes.NewBuffer(ballJSON))
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Body.Close()
	}

	// Step 6: Verify Match Completed
	req, _ = http.NewRequest("GET", server.URL+"/api/v1/scorecard/"+matchID, nil)
	req.Header.Set("Accept", "application/json")

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	err = json.NewDecoder(resp.Body).Decode(&scorecardResponse)
	require.NoError(t, err)
	resp.Body.Close()

	// Verify match is completed
	assert.Equal(t, "completed", scorecardResponse.Data.MatchStatus)

	// Verify second innings is completed
	for _, innings := range scorecardResponse.Data.Innings {
		if innings.InningsNumber == 2 {
			assert.Equal(t, "completed", innings.Status)
			assert.GreaterOrEqual(t, innings.TotalRuns, 43) // Target reached
			break
		}
	}
}

func TestCompleteMatchFlow_AllWicketsLost_E2E(t *testing.T) {
	// Setup
	server, db := testutils.SetupE2ETestServerWithDB(t)
	defer server.Close()
	defer db.Close()

	// Step 1: Create Series
	seriesReq := models.CreateSeriesRequest{
		Name:      fmt.Sprintf("E2E Test Series Wickets %d", time.Now().Unix()),
		StartDate: time.Now(),
		EndDate:   time.Now().Add(24 * time.Hour),
	}
	seriesJSON, _ := json.Marshal(seriesReq)
	req, _ := http.NewRequest("POST", server.URL+"/api/v1/series", bytes.NewBuffer(seriesJSON))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var seriesResponse struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&seriesResponse)
	require.NoError(t, err)
	resp.Body.Close()

	seriesID := seriesResponse.Data.ID

	// Step 2: Create Match
	matchNumber := 1
	matchReq := models.CreateMatchRequest{
		SeriesID:         seriesID,
		MatchNumber:      &matchNumber,
		Date:             time.Now(),
		TeamAPlayerCount: 3,
		TeamBPlayerCount: 3,
		TotalOvers:       2,
		TossWinner:       models.TeamTypeA,
		TossType:         models.TossTypeHeads,
	}
	matchJSON, _ := json.Marshal(matchReq)
	req, _ = http.NewRequest("POST", server.URL+"/api/v1/matches", bytes.NewBuffer(matchJSON))
	req.Header.Set("Content-Type", "application/json")

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var matchResponse struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&matchResponse)
	require.NoError(t, err)
	resp.Body.Close()

	matchID := matchResponse.Data.ID

	// Step 3: Complete First Innings with 10 runs
	firstInningsBalls := []models.BallEventRequest{
		{MatchID: matchID, InningsNumber: 1, BallType: models.BallTypeGood, RunType: models.RunTypeOne, IsWicket: false, Byes: 0},
		{MatchID: matchID, InningsNumber: 1, BallType: models.BallTypeGood, RunType: models.RunTypeOne, IsWicket: false, Byes: 0},
		{MatchID: matchID, InningsNumber: 1, BallType: models.BallTypeGood, RunType: models.RunTypeOne, IsWicket: false, Byes: 0},
		{MatchID: matchID, InningsNumber: 1, BallType: models.BallTypeGood, RunType: models.RunTypeOne, IsWicket: false, Byes: 0},
		{MatchID: matchID, InningsNumber: 1, BallType: models.BallTypeGood, RunType: models.RunTypeOne, IsWicket: false, Byes: 0},
		{MatchID: matchID, InningsNumber: 1, BallType: models.BallTypeGood, RunType: models.RunTypeOne, IsWicket: false, Byes: 0},
		{MatchID: matchID, InningsNumber: 1, BallType: models.BallTypeGood, RunType: models.RunTypeOne, IsWicket: false, Byes: 0},
		{MatchID: matchID, InningsNumber: 1, BallType: models.BallTypeGood, RunType: models.RunTypeOne, IsWicket: false, Byes: 0},
		{MatchID: matchID, InningsNumber: 1, BallType: models.BallTypeGood, RunType: models.RunTypeOne, IsWicket: false, Byes: 0},
		{MatchID: matchID, InningsNumber: 1, BallType: models.BallTypeGood, RunType: models.RunTypeOne, IsWicket: false, Byes: 0},
		{MatchID: matchID, InningsNumber: 1, BallType: models.BallTypeGood, RunType: models.RunTypeOne, IsWicket: false, Byes: 0},
		{MatchID: matchID, InningsNumber: 1, BallType: models.BallTypeGood, RunType: models.RunTypeOne, IsWicket: false, Byes: 0},
	}

	for _, ballReq := range firstInningsBalls {
		ballJSON, _ := json.Marshal(ballReq)
		req, _ := http.NewRequest("POST", server.URL+"/api/v1/scorecard/ball", bytes.NewBuffer(ballJSON))
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Body.Close()
	}

	// Step 4: Add balls to second innings to lose all wickets (2 wickets for 3 players)
	secondInningsBalls := []models.BallEventRequest{
		{MatchID: matchID, InningsNumber: 2, BallType: models.BallTypeGood, RunType: models.RunTypeWC, IsWicket: true, WicketType: "bowled", Byes: 0},
		{MatchID: matchID, InningsNumber: 2, BallType: models.BallTypeGood, RunType: models.RunTypeWC, IsWicket: true, WicketType: "bowled", Byes: 0},
	}

	for _, ballReq := range secondInningsBalls {
		ballJSON, _ := json.Marshal(ballReq)
		req, _ := http.NewRequest("POST", server.URL+"/api/v1/scorecard/ball", bytes.NewBuffer(ballJSON))
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Body.Close()
	}

	// Step 5: Verify Match Completed
	req, _ = http.NewRequest("GET", server.URL+"/api/v1/scorecard/"+matchID, nil)
	req.Header.Set("Accept", "application/json")

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
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
	resp.Body.Close()

	// Verify match is completed
	assert.Equal(t, "completed", scorecardResponse.Data.MatchStatus)

	// Verify second innings is completed with all wickets lost
	for _, innings := range scorecardResponse.Data.Innings {
		if innings.InningsNumber == 2 {
			assert.Equal(t, "completed", innings.Status)
			assert.Equal(t, 2, innings.TotalWickets) // All wickets lost (n-1)
			break
		}
	}
}

func TestCompleteMatchFlow_AllOversCompleted_E2E(t *testing.T) {
	// Setup
	server, db := testutils.SetupE2ETestServerWithDB(t)
	defer server.Close()
	defer db.Close()

	// Step 1: Create Series
	seriesReq := models.CreateSeriesRequest{
		Name:      fmt.Sprintf("E2E Test Series Overs %d", time.Now().Unix()),
		StartDate: time.Now(),
		EndDate:   time.Now().Add(24 * time.Hour),
	}
	seriesJSON, _ := json.Marshal(seriesReq)
	req, _ := http.NewRequest("POST", server.URL+"/api/v1/series", bytes.NewBuffer(seriesJSON))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var seriesResponse struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&seriesResponse)
	require.NoError(t, err)
	resp.Body.Close()

	seriesID := seriesResponse.Data.ID

	// Step 2: Create Match
	matchNumber := 1
	matchReq := models.CreateMatchRequest{
		SeriesID:         seriesID,
		MatchNumber:      &matchNumber,
		Date:             time.Now(),
		TeamAPlayerCount: 3,
		TeamBPlayerCount: 3,
		TotalOvers:       2,
		TossWinner:       models.TeamTypeA,
		TossType:         models.TossTypeHeads,
	}
	matchJSON, _ := json.Marshal(matchReq)
	req, _ = http.NewRequest("POST", server.URL+"/api/v1/matches", bytes.NewBuffer(matchJSON))
	req.Header.Set("Content-Type", "application/json")

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var matchResponse struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&matchResponse)
	require.NoError(t, err)
	resp.Body.Close()

	matchID := matchResponse.Data.ID

	// Step 3: Complete First Innings with 5 runs (so target is 6, won't be reached with 12 balls of 1 run each)
	firstInningsBalls := []models.BallEventRequest{
		{MatchID: matchID, InningsNumber: 1, BallType: models.BallTypeGood, RunType: models.RunTypeOne, IsWicket: false, Byes: 0},
		{MatchID: matchID, InningsNumber: 1, BallType: models.BallTypeGood, RunType: models.RunTypeOne, IsWicket: false, Byes: 0},
		{MatchID: matchID, InningsNumber: 1, BallType: models.BallTypeGood, RunType: models.RunTypeOne, IsWicket: false, Byes: 0},
		{MatchID: matchID, InningsNumber: 1, BallType: models.BallTypeGood, RunType: models.RunTypeOne, IsWicket: false, Byes: 0},
		{MatchID: matchID, InningsNumber: 1, BallType: models.BallTypeGood, RunType: models.RunTypeOne, IsWicket: false, Byes: 0},
		{MatchID: matchID, InningsNumber: 1, BallType: models.BallTypeGood, RunType: models.RunTypeZero, IsWicket: false, Byes: 0},
		{MatchID: matchID, InningsNumber: 1, BallType: models.BallTypeGood, RunType: models.RunTypeZero, IsWicket: false, Byes: 0},
		{MatchID: matchID, InningsNumber: 1, BallType: models.BallTypeGood, RunType: models.RunTypeZero, IsWicket: false, Byes: 0},
		{MatchID: matchID, InningsNumber: 1, BallType: models.BallTypeGood, RunType: models.RunTypeZero, IsWicket: false, Byes: 0},
		{MatchID: matchID, InningsNumber: 1, BallType: models.BallTypeGood, RunType: models.RunTypeZero, IsWicket: false, Byes: 0},
		{MatchID: matchID, InningsNumber: 1, BallType: models.BallTypeGood, RunType: models.RunTypeZero, IsWicket: false, Byes: 0},
		{MatchID: matchID, InningsNumber: 1, BallType: models.BallTypeGood, RunType: models.RunTypeZero, IsWicket: false, Byes: 0},
	}

	for _, ballReq := range firstInningsBalls {
		ballJSON, _ := json.Marshal(ballReq)
		req, _ := http.NewRequest("POST", server.URL+"/api/v1/scorecard/ball", bytes.NewBuffer(ballJSON))
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Body.Close()
	}

	// Step 4: Add 12 balls to second innings to complete all overs (score 0 runs so target is never reached)
	secondInningsBalls := make([]models.BallEventRequest, 12)
	for i := 0; i < 12; i++ {
		secondInningsBalls[i] = models.BallEventRequest{
			MatchID:       matchID,
			InningsNumber: 2,
			BallType:      models.BallTypeGood,
			RunType:       models.RunTypeZero,
			IsWicket:      false,
			Byes:          0,
		}
	}

	for _, ballReq := range secondInningsBalls {
		ballJSON, _ := json.Marshal(ballReq)
		req, _ := http.NewRequest("POST", server.URL+"/api/v1/scorecard/ball", bytes.NewBuffer(ballJSON))
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Body.Close()
	}

	// Step 5: Verify Match Completed
	req, _ = http.NewRequest("GET", server.URL+"/api/v1/scorecard/"+matchID, nil)
	req.Header.Set("Accept", "application/json")

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
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
	resp.Body.Close()

	// Verify match is completed
	assert.Equal(t, "completed", scorecardResponse.Data.MatchStatus)

	// Verify second innings is completed with all overs
	for _, innings := range scorecardResponse.Data.Innings {
		if innings.InningsNumber == 2 {
			assert.Equal(t, "completed", innings.Status)
			assert.GreaterOrEqual(t, innings.TotalOvers, 2.0) // All overs completed
			break
		}
	}
}
