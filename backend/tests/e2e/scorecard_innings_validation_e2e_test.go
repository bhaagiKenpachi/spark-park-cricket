package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"spark-park-cricket-backend/internal/config"
	"spark-park-cricket-backend/internal/database"
	"spark-park-cricket-backend/internal/handlers"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/services"

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
	cleanupTestData(t, testDB)

	// Start test server
	server := setupE2ETestServer(t, testDB)
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

		assert.Equal(t, "Ball added successfully", successResp["message"])
	})

	// Clean up after test
	cleanupTestData(t, testDB)
}

func setupE2ETestServer(t *testing.T, testDB *database.Client) *httptest.Server {
	// Create service container
	serviceContainer := services.NewContainer(testDB.Repositories)

	// Create handlers
	seriesHandler := handlers.NewSeriesHandler(serviceContainer.Series)
	matchHandler := handlers.NewMatchHandler(serviceContainer.Match)
	scorecardHandler := handlers.NewScorecardHandler(serviceContainer.Scorecard)

	// Create router and register routes
	router := http.NewServeMux()

	// Series routes
	router.HandleFunc("/api/v1/series", seriesHandler.CreateSeries)

	// Match routes
	router.HandleFunc("/api/v1/matches", matchHandler.CreateMatch)
	router.HandleFunc("/api/v1/matches/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if len(path) > len("/api/v1/matches/") {
			matchID := path[len("/api/v1/matches/"):]
			switch r.Method {
			case http.MethodGet:
				matchHandler.GetMatch(w, r)
			case http.MethodPut:
				matchHandler.UpdateMatch(w, r)
			case http.MethodDelete:
				matchHandler.DeleteMatch(w, r)
			}
			// Store matchID in context for handlers to use
			_ = matchID
		}
	})

	// Scorecard routes
	router.HandleFunc("/api/v1/scorecard/ball", scorecardHandler.AddBall)

	return httptest.NewServer(router)
}

func cleanupTestData(t *testing.T, testDB *database.Client) {
	// Clean up matches
	_, err := testDB.Supabase.From("matches").Delete("", "").ExecuteTo(nil)
	if err != nil {
		t.Logf("Warning: Failed to cleanup matches: %v", err)
	}

	// Clean up series
	_, err = testDB.Supabase.From("series").Delete("", "").ExecuteTo(nil)
	if err != nil {
		t.Logf("Warning: Failed to cleanup series: %v", err)
	}
}
