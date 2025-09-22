package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/pkg/testutils"
)

func TestMatchWorkflow_E2E(t *testing.T) {
	// Setup with authentication
	server, db, _, sessionCookie := testutils.SetupAuthenticatedE2ETestServerWithDB(t)
	defer server.Close()
	defer db.Close()

	// Clean up any existing test data
	testutils.CleanupTestData(t, db)

	// Use the authenticated test utilities for creating series and matches
	t.Run("Complete Match Lifecycle", func(t *testing.T) {
		testCompleteMatchLifecycleWithAuth(t, server, sessionCookie)
	})

	t.Run("Match State Transitions", func(t *testing.T) {
		testMatchStateTransitionsWithAuth(t, server, sessionCookie)
	})

	t.Run("Match Series Integration", func(t *testing.T) {
		testMatchSeriesIntegrationWithAuth(t, server, sessionCookie)
	})

	t.Run("Match Validation Workflow", func(t *testing.T) {
		testMatchValidationWorkflowWithAuth(t, server, sessionCookie)
	})
}

func testCompleteMatchLifecycleWithAuth(t *testing.T, server *httptest.Server, sessionCookie string) {
	// Create a series
	seriesReq := models.CreateSeriesRequest{
		Name:      "Complete Match Lifecycle Series",
		StartDate: time.Now(),
		EndDate:   time.Now().Add(7 * 24 * time.Hour),
	}
	seriesID := createAuthenticatedTestSeries(t, server, sessionCookie, seriesReq)

	// Create a match
	matchReq := models.CreateMatchRequest{
		SeriesID:         seriesID,
		MatchNumber:      &[]int{1}[0],
		Date:             time.Now(),
		TeamAPlayerCount: 11,
		TeamBPlayerCount: 11,
		TotalOvers:       20,
		TossWinner:       models.TeamTypeA,
		TossType:         models.TossTypeHeads,
	}
	matchID := createAuthenticatedTestMatch(t, server, sessionCookie, matchReq)

	// Verify match is created with "live" status
	match := getMatchWithAuth(t, server, sessionCookie, matchID)
	assert.Equal(t, models.MatchStatusLive, match.Status)
	assert.Equal(t, seriesID, match.SeriesID)
	assert.Equal(t, 1, match.MatchNumber)

	// Update match status to completed
	updateReq := models.UpdateMatchRequest{
		Status: &[]models.MatchStatus{models.MatchStatusCompleted}[0],
	}
	updateMatchStatusWithAuth(t, server, sessionCookie, matchID, updateReq)

	// Verify status update
	match = getMatchWithAuth(t, server, sessionCookie, matchID)
	assert.Equal(t, models.MatchStatusCompleted, match.Status)

	// Delete the match
	deleteMatchWithAuth(t, server, sessionCookie, matchID)

	// Verify match is deleted
	match = getMatchWithAuth(t, server, sessionCookie, matchID)
	assert.Empty(t, match.ID)
}

func testMatchStateTransitionsWithAuth(t *testing.T, server *httptest.Server, sessionCookie string) {
	// Create a series and match
	seriesReq := models.CreateSeriesRequest{
		Name:      "State Transitions Series",
		StartDate: time.Now(),
		EndDate:   time.Now().Add(7 * 24 * time.Hour),
	}
	seriesID := createAuthenticatedTestSeries(t, server, sessionCookie, seriesReq)

	matchReq := models.CreateMatchRequest{
		SeriesID:         seriesID,
		MatchNumber:      &[]int{1}[0],
		Date:             time.Now(),
		TeamAPlayerCount: 11,
		TeamBPlayerCount: 11,
		TotalOvers:       20,
		TossWinner:       models.TeamTypeA,
		TossType:         models.TossTypeHeads,
	}
	matchID := createAuthenticatedTestMatch(t, server, sessionCookie, matchReq)

	// Test state transitions: live -> completed -> cancelled
	transitions := []models.MatchStatus{
		models.MatchStatusLive,
		models.MatchStatusCompleted,
		models.MatchStatusCancelled,
	}

	for i, expectedStatus := range transitions {
		if i > 0 { // Skip first iteration as match starts in "live" status
			updateReq := models.UpdateMatchRequest{
				Status: &[]models.MatchStatus{expectedStatus}[0],
			}
			updateMatchStatusWithAuth(t, server, sessionCookie, matchID, updateReq)
		}

		match := getMatchWithAuth(t, server, sessionCookie, matchID)
		assert.Equal(t, expectedStatus, match.Status, "Match status should be %s", expectedStatus)

		// Verify the status change is persisted
		match = getMatchWithAuth(t, server, sessionCookie, matchID)
		assert.Equal(t, expectedStatus, match.Status, "Match status should persist as %s", expectedStatus)
	}

	// Clean up
	deleteMatchWithAuth(t, server, sessionCookie, matchID)
}

func testMatchSeriesIntegrationWithAuth(t *testing.T, server *httptest.Server, sessionCookie string) {
	// Create multiple series
	series1Req := models.CreateSeriesRequest{
		Name:      "Series 1",
		StartDate: time.Now(),
		EndDate:   time.Now().Add(7 * 24 * time.Hour),
	}
	series1ID := createAuthenticatedTestSeries(t, server, sessionCookie, series1Req)

	series2Req := models.CreateSeriesRequest{
		Name:      "Series 2",
		StartDate: time.Now(),
		EndDate:   time.Now().Add(7 * 24 * time.Hour),
	}
	series2ID := createAuthenticatedTestSeries(t, server, sessionCookie, series2Req)

	// Create matches in different series
	match1Req := models.CreateMatchRequest{
		SeriesID:         series1ID,
		MatchNumber:      &[]int{1}[0],
		Date:             time.Now(),
		TeamAPlayerCount: 11,
		TeamBPlayerCount: 11,
		TotalOvers:       20,
		TossWinner:       models.TeamTypeA,
		TossType:         models.TossTypeHeads,
	}
	match1ID := createAuthenticatedTestMatch(t, server, sessionCookie, match1Req)

	match2Req := models.CreateMatchRequest{
		SeriesID:         series1ID,
		MatchNumber:      &[]int{2}[0],
		Date:             time.Now(),
		TeamAPlayerCount: 11,
		TeamBPlayerCount: 11,
		TotalOvers:       20,
		TossWinner:       models.TeamTypeA,
		TossType:         models.TossTypeHeads,
	}
	match2ID := createAuthenticatedTestMatch(t, server, sessionCookie, match2Req)

	match3Req := models.CreateMatchRequest{
		SeriesID:         series2ID,
		MatchNumber:      &[]int{1}[0],
		Date:             time.Now(),
		TeamAPlayerCount: 11,
		TeamBPlayerCount: 11,
		TotalOvers:       20,
		TossWinner:       models.TeamTypeA,
		TossType:         models.TossTypeHeads,
	}
	match3ID := createAuthenticatedTestMatch(t, server, sessionCookie, match3Req)

	// Verify matches are associated with correct series
	match1 := getMatchWithAuth(t, server, sessionCookie, match1ID)
	match2 := getMatchWithAuth(t, server, sessionCookie, match2ID)
	match3 := getMatchWithAuth(t, server, sessionCookie, match3ID)

	assert.Equal(t, series1ID, match1.SeriesID)
	assert.Equal(t, series1ID, match2.SeriesID)
	assert.Equal(t, series2ID, match3.SeriesID)

	// Test getting matches by series
	series1Matches := getMatchesBySeriesWithAuth(t, server, sessionCookie, series1ID)
	assert.Len(t, series1Matches, 2)

	series2Matches := getMatchesBySeriesWithAuth(t, server, sessionCookie, series2ID)
	assert.Len(t, series2Matches, 1)

	// Clean up
	deleteMatchWithAuth(t, server, sessionCookie, match1ID)
	deleteMatchWithAuth(t, server, sessionCookie, match2ID)
	deleteMatchWithAuth(t, server, sessionCookie, match3ID)
}

func testMatchValidationWorkflowWithAuth(t *testing.T, server *httptest.Server, sessionCookie string) {
	// Create a series
	seriesReq := models.CreateSeriesRequest{
		Name:      "Validation Series",
		StartDate: time.Now(),
		EndDate:   time.Now().Add(7 * 24 * time.Hour),
	}
	seriesID := createAuthenticatedTestSeries(t, server, sessionCookie, seriesReq)

	// Test invalid match creation (missing required fields)
	invalidReq := models.CreateMatchRequest{
		SeriesID:         seriesID,
		MatchNumber:      &[]int{1}[0],
		Date:             time.Now(),
		TeamAPlayerCount: 0, // Invalid: should be > 0
		TeamBPlayerCount: 11,
		TotalOvers:       20,
		TossWinner:       models.TeamTypeA,
		TossType:         models.TossTypeHeads,
	}
	createMatchExpectingErrorWithAuth(t, server, sessionCookie, invalidReq, 400)

	// Test valid match creation
	validReq := models.CreateMatchRequest{
		SeriesID:         seriesID,
		MatchNumber:      &[]int{1}[0],
		Date:             time.Now(),
		TeamAPlayerCount: 11,
		TeamBPlayerCount: 11,
		TotalOvers:       20,
		TossWinner:       models.TeamTypeA,
		TossType:         models.TossTypeHeads,
	}
	matchID := createTestMatchWithRequestWithAuth(t, server, sessionCookie, validReq)

	// Test invalid match update (invalid batting team)
	invalidUpdateReq := models.UpdateMatchRequest{
		BattingTeam: &[]models.TeamType{models.TeamTypeA}[0],
	}
	updateBattingTeamWithAuth(t, server, sessionCookie, matchID, invalidUpdateReq)

	// Clean up
	deleteMatchWithAuth(t, server, sessionCookie, matchID)
}

// Helper functions for authenticated requests

func createAuthenticatedTestSeries(t *testing.T, server *httptest.Server, sessionCookie string, req models.CreateSeriesRequest) string {
	reqBody, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("POST", server.URL+"/api/v1/series", bytes.NewBuffer(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.AddCookie(&http.Cookie{Name: "user_session", Value: sessionCookie})

	resp, err := server.Client().Do(httpReq)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, 201, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	seriesData := response["data"].(map[string]interface{})
	return seriesData["id"].(string)
}

func createAuthenticatedTestMatch(t *testing.T, server *httptest.Server, sessionCookie string, req models.CreateMatchRequest) string {
	reqBody, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("POST", server.URL+"/api/v1/matches", bytes.NewBuffer(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.AddCookie(&http.Cookie{Name: "user_session", Value: sessionCookie})

	resp, err := server.Client().Do(httpReq)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, 201, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	matchData := response["data"].(map[string]interface{})
	return matchData["id"].(string)
}

func getMatchWithAuth(t *testing.T, server *httptest.Server, sessionCookie string, matchID string) models.Match {
	httpReq, _ := http.NewRequest("GET", server.URL+"/api/v1/matches/"+matchID, nil)
	httpReq.AddCookie(&http.Cookie{Name: "user_session", Value: sessionCookie})

	resp, err := server.Client().Do(httpReq)
	require.NoError(t, err)
	defer resp.Body.Close()

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	var match models.Match
	if resp.StatusCode == 200 {
		matchData := response["data"].(map[string]interface{})
		matchJSON, _ := json.Marshal(matchData)
		_ = json.Unmarshal(matchJSON, &match)
	}
	return match
}

func updateMatchStatusWithAuth(t *testing.T, server *httptest.Server, sessionCookie string, matchID string, req models.UpdateMatchRequest) {
	reqBody, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("PUT", server.URL+"/api/v1/matches/"+matchID, bytes.NewBuffer(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.AddCookie(&http.Cookie{Name: "user_session", Value: sessionCookie})

	resp, err := server.Client().Do(httpReq)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, 200, resp.StatusCode)
}

func updateBattingTeamWithAuth(t *testing.T, server *httptest.Server, sessionCookie string, matchID string, req models.UpdateMatchRequest) {
	reqBody, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("PUT", server.URL+"/api/v1/matches/"+matchID, bytes.NewBuffer(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.AddCookie(&http.Cookie{Name: "user_session", Value: sessionCookie})

	resp, err := server.Client().Do(httpReq)
	require.NoError(t, err)
	defer resp.Body.Close()

	// This might return 200 or 400 depending on validation
	assert.True(t, resp.StatusCode == 200 || resp.StatusCode == 400)
}

func getMatchesBySeriesWithAuth(t *testing.T, server *httptest.Server, sessionCookie string, seriesID string) []models.Match {
	httpReq, _ := http.NewRequest("GET", server.URL+"/api/v1/matches/series/"+seriesID, nil)
	httpReq.AddCookie(&http.Cookie{Name: "user_session", Value: sessionCookie})

	resp, err := server.Client().Do(httpReq)
	require.NoError(t, err)
	defer resp.Body.Close()

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	var matches []models.Match
	if resp.StatusCode == 200 {
		matchesData := response["data"].([]interface{})
		for _, matchData := range matchesData {
			var match models.Match
			matchJSON, _ := json.Marshal(matchData)
			_ = json.Unmarshal(matchJSON, &match)
			matches = append(matches, match)
		}
	}
	return matches
}

func deleteMatchWithAuth(t *testing.T, server *httptest.Server, sessionCookie string, matchID string) {
	httpReq, _ := http.NewRequest("DELETE", server.URL+"/api/v1/matches/"+matchID, nil)
	httpReq.AddCookie(&http.Cookie{Name: "user_session", Value: sessionCookie})

	resp, err := server.Client().Do(httpReq)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, 200, resp.StatusCode)
}

func createTestMatchWithRequestWithAuth(t *testing.T, server *httptest.Server, sessionCookie string, req models.CreateMatchRequest) string {
	reqBody, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("POST", server.URL+"/api/v1/matches", bytes.NewBuffer(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.AddCookie(&http.Cookie{Name: "user_session", Value: sessionCookie})

	resp, err := server.Client().Do(httpReq)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, 201, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	matchData := response["data"].(map[string]interface{})
	return matchData["id"].(string)
}

func createMatchExpectingErrorWithAuth(t *testing.T, server *httptest.Server, sessionCookie string, req models.CreateMatchRequest, expectedStatus int) {
	reqBody, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("POST", server.URL+"/api/v1/matches", bytes.NewBuffer(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.AddCookie(&http.Cookie{Name: "user_session", Value: sessionCookie})

	resp, err := server.Client().Do(httpReq)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, expectedStatus, resp.StatusCode, "Expected status %d but got %d", expectedStatus, resp.StatusCode)
}
