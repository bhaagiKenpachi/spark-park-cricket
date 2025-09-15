package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"spark-park-cricket-backend/internal/config"
	"spark-park-cricket-backend/internal/database"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/pkg/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScorecardInningsValidation_E2E(t *testing.T) {
	// Setup test database
	cfg := config.LoadTestConfig()
	testDB, err := database.NewTestClient(cfg)
	require.NoError(t, err)
	defer testDB.Close()

	// Clean up before test
	testutils.CleanupTestData(t, testDB)

	// Start test server
	server := testutils.SetupE2ETestServer(t, testDB)
	defer server.Close()

	t.Run("API prevents adding ball to second innings before first innings", func(t *testing.T) {
		// Create a test series via API
		seriesReq := &models.CreateSeriesRequest{
			Name:      "Test E2E Innings Series",
			StartDate: time.Now(),
			EndDate:   time.Now().AddDate(0, 0, 7),
		}

		seriesBody, _ := json.Marshal(seriesReq)
		seriesResp, err := http.Post(server.URL+"/api/v1/series", "application/json", bytes.NewBuffer(seriesBody))
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, seriesResp.StatusCode)

		var seriesRespData map[string]interface{}
		err = json.NewDecoder(seriesResp.Body).Decode(&seriesRespData)
		require.NoError(t, err)
		seriesResp.Body.Close()

		seriesData := seriesRespData["data"].(map[string]interface{})
		seriesID := seriesData["id"].(string)

		// Create a match with Team A winning toss
		matchReq := &models.CreateMatchRequest{
			SeriesID:         seriesID,
			MatchNumber:      nil, // Auto-increment
			Date:             time.Now(),
			TeamAPlayerCount: 11,
			TeamBPlayerCount: 11,
			TotalOvers:       20,
			TossWinner:       models.TeamTypeA,
			TossType:         models.TossTypeHeads,
		}

		matchBody, _ := json.Marshal(matchReq)
		matchResp, err := http.Post(server.URL+"/api/v1/matches", "application/json", bytes.NewBuffer(matchBody))
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, matchResp.StatusCode)

		var matchRespData map[string]interface{}
		err = json.NewDecoder(matchResp.Body).Decode(&matchRespData)
		require.NoError(t, err)
		matchResp.Body.Close()

		matchData := matchRespData["data"].(map[string]interface{})
		matchID := matchData["id"].(string)

		// Try to add a ball to second innings directly - this should fail
		ballEvent := &models.BallEventRequest{
			MatchID:       matchID,
			InningsNumber: 2,
			BallType:      models.BallTypeGood,
			RunType:       models.RunTypeOne,
			IsWicket:      false,
		}

		ballBody, _ := json.Marshal(ballEvent)
		ballResp, err := http.Post(server.URL+"/api/v1/scorecard/ball", "application/json", bytes.NewBuffer(ballBody))
		require.NoError(t, err)
		require.Equal(t, http.StatusInternalServerError, ballResp.StatusCode)

		var errorResp map[string]interface{}
		err = json.NewDecoder(ballResp.Body).Decode(&errorResp)
		require.NoError(t, err)
		ballResp.Body.Close()

		assert.Contains(t, errorResp["error"].(map[string]interface{})["message"], "cannot start second innings, first innings must be played first")
	})

	t.Run("API prevents adding ball to wrong team in first innings", func(t *testing.T) {
		// Create a test series via API
		seriesReq := &models.CreateSeriesRequest{
			Name:      "Test E2E Innings Series 2",
			StartDate: time.Now(),
			EndDate:   time.Now().AddDate(0, 0, 7),
		}

		seriesBody, _ := json.Marshal(seriesReq)
		seriesResp, err := http.Post(server.URL+"/api/v1/series", "application/json", bytes.NewBuffer(seriesBody))
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, seriesResp.StatusCode)

		var seriesRespData map[string]interface{}
		err = json.NewDecoder(seriesResp.Body).Decode(&seriesRespData)
		require.NoError(t, err)
		seriesResp.Body.Close()

		seriesData := seriesRespData["data"].(map[string]interface{})
		seriesID := seriesData["id"].(string)

		// Create a match with Team A winning toss
		matchReq := &models.CreateMatchRequest{
			SeriesID:         seriesID,
			MatchNumber:      nil, // Auto-increment
			Date:             time.Now(),
			TeamAPlayerCount: 11,
			TeamBPlayerCount: 11,
			TotalOvers:       20,
			TossWinner:       models.TeamTypeA,
			TossType:         models.TossTypeHeads,
		}

		matchBody, _ := json.Marshal(matchReq)
		matchResp, err := http.Post(server.URL+"/api/v1/matches", "application/json", bytes.NewBuffer(matchBody))
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, matchResp.StatusCode)

		var matchRespData map[string]interface{}
		err = json.NewDecoder(matchResp.Body).Decode(&matchRespData)
		require.NoError(t, err)
		matchResp.Body.Close()

		matchData := matchRespData["data"].(map[string]interface{})
		matchID := matchData["id"].(string)

		// Manually update the match to change batting team to Team B (non-toss winner)
		updateReq := &models.UpdateMatchRequest{
			BattingTeam: &[]models.TeamType{models.TeamTypeB}[0],
		}

		updateBody, _ := json.Marshal(updateReq)
		updateResp, err := http.NewRequest("PUT", server.URL+"/api/v1/matches/"+matchID, bytes.NewBuffer(updateBody))
		require.NoError(t, err)
		updateResp.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		updateRespResult, err := client.Do(updateResp)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, updateRespResult.StatusCode)
		updateRespResult.Body.Close()

		// Try to add a ball to first innings with Team B - this should fail
		ballEvent := &models.BallEventRequest{
			MatchID:       matchID,
			InningsNumber: 1,
			BallType:      models.BallTypeGood,
			RunType:       models.RunTypeOne,
			IsWicket:      false,
		}

		ballBody, _ := json.Marshal(ballEvent)
		ballResp, err := http.Post(server.URL+"/api/v1/scorecard/ball", "application/json", bytes.NewBuffer(ballBody))
		require.NoError(t, err)
		require.Equal(t, http.StatusInternalServerError, ballResp.StatusCode)

		var errorResp map[string]interface{}
		err = json.NewDecoder(ballResp.Body).Decode(&errorResp)
		require.NoError(t, err)
		ballResp.Body.Close()

		assert.Contains(t, errorResp["error"].(map[string]interface{})["message"], "first innings must be played by the toss-winning team")
	})

	t.Run("API allows correct first innings ball", func(t *testing.T) {
		// Create a test series via API
		seriesReq := &models.CreateSeriesRequest{
			Name:      "Test E2E Innings Series 3",
			StartDate: time.Now(),
			EndDate:   time.Now().AddDate(0, 0, 7),
		}

		seriesBody, _ := json.Marshal(seriesReq)
		seriesResp, err := http.Post(server.URL+"/api/v1/series", "application/json", bytes.NewBuffer(seriesBody))
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, seriesResp.StatusCode)

		var seriesRespData map[string]interface{}
		err = json.NewDecoder(seriesResp.Body).Decode(&seriesRespData)
		require.NoError(t, err)
		seriesResp.Body.Close()

		seriesData := seriesRespData["data"].(map[string]interface{})
		seriesID := seriesData["id"].(string)

		// Create a match with Team A winning toss
		matchReq := &models.CreateMatchRequest{
			SeriesID:         seriesID,
			MatchNumber:      nil, // Auto-increment
			Date:             time.Now(),
			TeamAPlayerCount: 11,
			TeamBPlayerCount: 11,
			TotalOvers:       20,
			TossWinner:       models.TeamTypeA,
			TossType:         models.TossTypeHeads,
		}

		matchBody, _ := json.Marshal(matchReq)
		matchResp, err := http.Post(server.URL+"/api/v1/matches", "application/json", bytes.NewBuffer(matchBody))
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, matchResp.StatusCode)

		var matchRespData map[string]interface{}
		err = json.NewDecoder(matchResp.Body).Decode(&matchRespData)
		require.NoError(t, err)
		matchResp.Body.Close()

		matchData := matchRespData["data"].(map[string]interface{})
		matchID := matchData["id"].(string)

		// Try to add a ball to first innings with Team A (toss winner) - this should work
		ballEvent := &models.BallEventRequest{
			MatchID:       matchID,
			InningsNumber: 1,
			BallType:      models.BallTypeGood,
			RunType:       models.RunTypeOne,
			IsWicket:      false,
		}

		ballBody, _ := json.Marshal(ballEvent)
		ballResp, err := http.Post(server.URL+"/api/v1/scorecard/ball", "application/json", bytes.NewBuffer(ballBody))
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, ballResp.StatusCode)

		var successResp map[string]interface{}
		err = json.NewDecoder(ballResp.Body).Decode(&successResp)
		require.NoError(t, err)
		ballResp.Body.Close()

		// The response is wrapped in a "data" field
		data := successResp["data"].(map[string]interface{})
		assert.Equal(t, "Ball added successfully", data["message"])
	})

	// Clean up after test
	testutils.CleanupTestData(t, testDB)
}
