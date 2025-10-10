package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScorecardInningsValidation_E2E(t *testing.T) {
	// Setup with authentication
	server, db, _, sessionCookie := testutils.SetupAuthenticatedE2ETestServerWithDB(t)
	defer server.Close()
	defer db.Close()

	// Clean up before test
	testutils.CleanupTestData(t, db)

	t.Run("API prevents adding ball to second innings before first innings", func(t *testing.T) {
		// Create a test series via API
		seriesReq := &models.CreateSeriesRequest{
			Name:      "Test E2E Innings Series",
			StartDate: time.Now(),
			EndDate:   time.Now().AddDate(0, 0, 7),
		}

		seriesBody, _ := json.Marshal(seriesReq)
		seriesHTTPReq, _ := http.NewRequest("POST", server.URL+"/api/v1/series", bytes.NewBuffer(seriesBody))
		seriesHTTPReq.Header.Set("Content-Type", "application/json")
		seriesHTTPReq.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    sessionCookie,
			Path:     "/",
			HttpOnly: true,
		})

		seriesResp, err := http.DefaultClient.Do(seriesHTTPReq)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, seriesResp.StatusCode)

		var seriesRespData map[string]interface{}
		err = json.NewDecoder(seriesResp.Body).Decode(&seriesRespData)
		require.NoError(t, err)

		// Extract series ID from response
		seriesData := seriesRespData["data"].(map[string]interface{})
		seriesID := seriesData["id"].(string)

		// Create a test match via API
		matchReq := &models.CreateMatchRequest{
			SeriesID:         seriesID,
			Date:             time.Now(),
			TeamAPlayerCount: 11,
			TeamBPlayerCount: 11,
			TotalOvers:       20,
			TossWinner:       models.TeamTypeA,
			TossType:         models.TossTypeHeads,
		}

		matchBody, _ := json.Marshal(matchReq)
		matchHTTPReq, _ := http.NewRequest("POST", server.URL+"/api/v1/matches", bytes.NewBuffer(matchBody))
		matchHTTPReq.Header.Set("Content-Type", "application/json")
		matchHTTPReq.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    sessionCookie,
			Path:     "/",
			HttpOnly: true,
		})

		matchResp, err := http.DefaultClient.Do(matchHTTPReq)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, matchResp.StatusCode)

		var matchRespData map[string]interface{}
		err = json.NewDecoder(matchResp.Body).Decode(&matchRespData)
		require.NoError(t, err)

		// Extract match ID from response
		matchData := matchRespData["data"].(map[string]interface{})
		matchID := matchData["id"].(string)

		// Start scoring for the match
		startReq := map[string]interface{}{
			"match_id": matchID,
		}
		startBody, _ := json.Marshal(startReq)
		startHTTPReq, _ := http.NewRequest("POST", server.URL+"/api/v1/scorecard/start", bytes.NewBuffer(startBody))
		startHTTPReq.Header.Set("Content-Type", "application/json")
		startHTTPReq.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    sessionCookie,
			Path:     "/",
			HttpOnly: true,
		})

		startResp, err := http.DefaultClient.Do(startHTTPReq)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, startResp.StatusCode)

		// Add some balls to first innings
		ballEvent := &models.BallEventRequest{
			MatchID:       matchID,
			InningsNumber: 1, // First innings
			BallType:      models.BallTypeGood,
			RunType:       models.RunTypeOne,
			IsWicket:      false,
		}

		ballBody, _ := json.Marshal(ballEvent)
		ballHTTPReq, _ := http.NewRequest("POST", server.URL+"/api/v1/scorecard/ball", bytes.NewBuffer(ballBody))
		ballHTTPReq.Header.Set("Content-Type", "application/json")
		ballHTTPReq.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    sessionCookie,
			Path:     "/",
			HttpOnly: true,
		})

		ballResp, err := http.DefaultClient.Do(ballHTTPReq)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, ballResp.StatusCode)

		// Try to add ball to second innings before first innings is complete
		ballEvent = &models.BallEventRequest{
			MatchID:       matchID,
			InningsNumber: 2, // Second innings
			BallType:      models.BallTypeGood,
			RunType:       models.RunTypeOne,
			IsWicket:      false,
		}

		ballBody, _ = json.Marshal(ballEvent)
		ballHTTPReq, _ = http.NewRequest("POST", server.URL+"/api/v1/scorecard/ball", bytes.NewBuffer(ballBody))
		ballHTTPReq.Header.Set("Content-Type", "application/json")
		ballHTTPReq.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    sessionCookie,
			Path:     "/",
			HttpOnly: true,
		})

		ballResp, err = http.DefaultClient.Do(ballHTTPReq)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, ballResp.StatusCode)

		var errorResp map[string]interface{}
		err = json.NewDecoder(ballResp.Body).Decode(&errorResp)
		require.NoError(t, err)
		assert.Contains(t, errorResp["error"].(map[string]interface{})["message"], "first innings is not complete, cannot start second innings")
	})

	t.Run("API prevents adding ball to wrong team", func(t *testing.T) {
		// Create a test series via API
		seriesReq := &models.CreateSeriesRequest{
			Name:      "Test E2E Team Validation Series",
			StartDate: time.Now(),
			EndDate:   time.Now().AddDate(0, 0, 7),
		}

		seriesBody, _ := json.Marshal(seriesReq)
		seriesHTTPReq, _ := http.NewRequest("POST", server.URL+"/api/v1/series", bytes.NewBuffer(seriesBody))
		seriesHTTPReq.Header.Set("Content-Type", "application/json")
		seriesHTTPReq.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    sessionCookie,
			Path:     "/",
			HttpOnly: true,
		})

		seriesResp, err := http.DefaultClient.Do(seriesHTTPReq)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, seriesResp.StatusCode)

		var seriesRespData map[string]interface{}
		err = json.NewDecoder(seriesResp.Body).Decode(&seriesRespData)
		require.NoError(t, err)

		// Extract series ID from response
		seriesData := seriesRespData["data"].(map[string]interface{})
		seriesID := seriesData["id"].(string)

		// Create a test match via API
		matchReq := &models.CreateMatchRequest{
			SeriesID:         seriesID,
			Date:             time.Now(),
			TeamAPlayerCount: 11,
			TeamBPlayerCount: 11,
			TotalOvers:       20,
			TossWinner:       models.TeamTypeA,
			TossType:         models.TossTypeHeads,
		}

		matchBody, _ := json.Marshal(matchReq)
		matchHTTPReq, _ := http.NewRequest("POST", server.URL+"/api/v1/matches", bytes.NewBuffer(matchBody))
		matchHTTPReq.Header.Set("Content-Type", "application/json")
		matchHTTPReq.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    sessionCookie,
			Path:     "/",
			HttpOnly: true,
		})

		matchResp, err := http.DefaultClient.Do(matchHTTPReq)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, matchResp.StatusCode)

		var matchRespData map[string]interface{}
		err = json.NewDecoder(matchResp.Body).Decode(&matchRespData)
		require.NoError(t, err)

		// Extract match ID from response
		matchData := matchRespData["data"].(map[string]interface{})
		matchID := matchData["id"].(string)

		// Start scoring for the match
		startReq := map[string]interface{}{
			"match_id": matchID,
		}
		startBody, _ := json.Marshal(startReq)
		startHTTPReq, _ := http.NewRequest("POST", server.URL+"/api/v1/scorecard/start", bytes.NewBuffer(startBody))
		startHTTPReq.Header.Set("Content-Type", "application/json")
		startHTTPReq.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    sessionCookie,
			Path:     "/",
			HttpOnly: true,
		})

		startResp, err := http.DefaultClient.Do(startHTTPReq)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, startResp.StatusCode)

		// Add some balls to first innings
		ballEvent := &models.BallEventRequest{
			MatchID:       matchID,
			InningsNumber: 1, // First innings
			BallType:      models.BallTypeGood,
			RunType:       models.RunTypeOne,
			IsWicket:      false,
		}

		ballBody, _ := json.Marshal(ballEvent)
		ballHTTPReq, _ := http.NewRequest("POST", server.URL+"/api/v1/scorecard/ball", bytes.NewBuffer(ballBody))
		ballHTTPReq.Header.Set("Content-Type", "application/json")
		ballHTTPReq.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    sessionCookie,
			Path:     "/",
			HttpOnly: true,
		})

		ballResp, err := http.DefaultClient.Do(ballHTTPReq)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, ballResp.StatusCode)

		// Try to add ball to second innings (should fail)
		ballEvent = &models.BallEventRequest{
			MatchID:       matchID,
			InningsNumber: 2, // Second innings
			BallType:      models.BallTypeGood,
			RunType:       models.RunTypeOne,
			IsWicket:      false,
		}

		ballBody, _ = json.Marshal(ballEvent)
		ballHTTPReq, _ = http.NewRequest("POST", server.URL+"/api/v1/scorecard/ball", bytes.NewBuffer(ballBody))
		ballHTTPReq.Header.Set("Content-Type", "application/json")
		ballHTTPReq.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    sessionCookie,
			Path:     "/",
			HttpOnly: true,
		})

		ballResp, err = http.DefaultClient.Do(ballHTTPReq)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, ballResp.StatusCode)

		var errorResp map[string]interface{}
		err = json.NewDecoder(ballResp.Body).Decode(&errorResp)
		require.NoError(t, err)
		assert.Contains(t, errorResp["error"].(map[string]interface{})["message"], "first innings is not complete, cannot start second innings")
	})

	t.Run("API prevents adding ball to first innings after second innings started", func(t *testing.T) {
		// Create a test series via API
		seriesReq := &models.CreateSeriesRequest{
			Name:      "Test E2E Innings Order Series",
			StartDate: time.Now(),
			EndDate:   time.Now().AddDate(0, 0, 7),
		}

		seriesBody, _ := json.Marshal(seriesReq)
		seriesHTTPReq, _ := http.NewRequest("POST", server.URL+"/api/v1/series", bytes.NewBuffer(seriesBody))
		seriesHTTPReq.Header.Set("Content-Type", "application/json")
		seriesHTTPReq.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    sessionCookie,
			Path:     "/",
			HttpOnly: true,
		})

		seriesResp, err := http.DefaultClient.Do(seriesHTTPReq)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, seriesResp.StatusCode)

		var seriesRespData map[string]interface{}
		err = json.NewDecoder(seriesResp.Body).Decode(&seriesRespData)
		require.NoError(t, err)

		// Extract series ID from response
		seriesData := seriesRespData["data"].(map[string]interface{})
		seriesID := seriesData["id"].(string)

		// Create a test match via API
		matchReq := &models.CreateMatchRequest{
			SeriesID:         seriesID,
			Date:             time.Now(),
			TeamAPlayerCount: 11,
			TeamBPlayerCount: 11,
			TotalOvers:       20,
			TossWinner:       models.TeamTypeA,
			TossType:         models.TossTypeHeads,
		}

		matchBody, _ := json.Marshal(matchReq)
		matchHTTPReq, _ := http.NewRequest("POST", server.URL+"/api/v1/matches", bytes.NewBuffer(matchBody))
		matchHTTPReq.Header.Set("Content-Type", "application/json")
		matchHTTPReq.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    sessionCookie,
			Path:     "/",
			HttpOnly: true,
		})

		matchResp, err := http.DefaultClient.Do(matchHTTPReq)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, matchResp.StatusCode)

		var matchRespData map[string]interface{}
		err = json.NewDecoder(matchResp.Body).Decode(&matchRespData)
		require.NoError(t, err)

		// Extract match ID from response
		matchData := matchRespData["data"].(map[string]interface{})
		matchID := matchData["id"].(string)

		// Start scoring for the match
		startReq := map[string]interface{}{
			"match_id": matchID,
		}
		startBody, _ := json.Marshal(startReq)
		startHTTPReq, _ := http.NewRequest("POST", server.URL+"/api/v1/scorecard/start", bytes.NewBuffer(startBody))
		startHTTPReq.Header.Set("Content-Type", "application/json")
		startHTTPReq.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    sessionCookie,
			Path:     "/",
			HttpOnly: true,
		})

		startResp, err := http.DefaultClient.Do(startHTTPReq)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, startResp.StatusCode)

		// Complete first innings by adding 120 balls (20 overs * 6 balls)
		for i := 1; i <= 20; i++ {
			for j := 1; j <= 6; j++ {
				ballEvent := &models.BallEventRequest{
					MatchID:       matchID,
					InningsNumber: 1, // First innings
					BallType:      models.BallTypeGood,
					RunType:       models.RunTypeOne,
					IsWicket:      false,
				}

				ballBody, _ := json.Marshal(ballEvent)
				ballHTTPReq, _ := http.NewRequest("POST", server.URL+"/api/v1/scorecard/ball", bytes.NewBuffer(ballBody))
				ballHTTPReq.Header.Set("Content-Type", "application/json")
				ballHTTPReq.AddCookie(&http.Cookie{
					Name:     "user_session",
					Value:    sessionCookie,
					Path:     "/",
					HttpOnly: true,
				})

				ballResp, err := http.DefaultClient.Do(ballHTTPReq)
				require.NoError(t, err)
				require.Equal(t, http.StatusOK, ballResp.StatusCode)
			}
		}

		// Start second innings
		ballEvent := &models.BallEventRequest{
			MatchID:       matchID,
			InningsNumber: 2, // Second innings
			BallType:      models.BallTypeGood,
			RunType:       models.RunTypeOne,
			IsWicket:      false,
		}

		ballBody, _ := json.Marshal(ballEvent)
		ballHTTPReq, _ := http.NewRequest("POST", server.URL+"/api/v1/scorecard/ball", bytes.NewBuffer(ballBody))
		ballHTTPReq.Header.Set("Content-Type", "application/json")
		ballHTTPReq.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    sessionCookie,
			Path:     "/",
			HttpOnly: true,
		})

		ballResp, err := http.DefaultClient.Do(ballHTTPReq)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, ballResp.StatusCode)

		// Try to add ball to first innings (should fail)
		ballEvent = &models.BallEventRequest{
			MatchID:       matchID,
			InningsNumber: 1, // First innings
			BallType:      models.BallTypeGood,
			RunType:       models.RunTypeOne,
			IsWicket:      false,
		}

		ballBody, _ = json.Marshal(ballEvent)
		ballHTTPReq, _ = http.NewRequest("POST", server.URL+"/api/v1/scorecard/ball", bytes.NewBuffer(ballBody))
		ballHTTPReq.Header.Set("Content-Type", "application/json")
		ballHTTPReq.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    sessionCookie,
			Path:     "/",
			HttpOnly: true,
		})

		ballResp, err = http.DefaultClient.Do(ballHTTPReq)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, ballResp.StatusCode)

		var errorResp map[string]interface{}
		err = json.NewDecoder(ballResp.Body).Decode(&errorResp)
		require.NoError(t, err)
		assert.Contains(t, errorResp["error"].(map[string]interface{})["message"], "innings is not in progress, cannot add ball")
	})

	t.Run("API allows adding ball to correct innings", func(t *testing.T) {
		// Create a test series via API
		seriesReq := &models.CreateSeriesRequest{
			Name:      "Test E2E Correct Innings Series",
			StartDate: time.Now(),
			EndDate:   time.Now().AddDate(0, 0, 7),
		}

		seriesBody, _ := json.Marshal(seriesReq)
		seriesHTTPReq, _ := http.NewRequest("POST", server.URL+"/api/v1/series", bytes.NewBuffer(seriesBody))
		seriesHTTPReq.Header.Set("Content-Type", "application/json")
		seriesHTTPReq.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    sessionCookie,
			Path:     "/",
			HttpOnly: true,
		})

		seriesResp, err := http.DefaultClient.Do(seriesHTTPReq)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, seriesResp.StatusCode)

		var seriesRespData map[string]interface{}
		err = json.NewDecoder(seriesResp.Body).Decode(&seriesRespData)
		require.NoError(t, err)

		// Extract series ID from response
		seriesData := seriesRespData["data"].(map[string]interface{})
		seriesID := seriesData["id"].(string)

		// Create a test match via API
		matchReq := &models.CreateMatchRequest{
			SeriesID:         seriesID,
			Date:             time.Now(),
			TeamAPlayerCount: 11,
			TeamBPlayerCount: 11,
			TotalOvers:       20,
			TossWinner:       models.TeamTypeA,
			TossType:         models.TossTypeHeads,
		}

		matchBody, _ := json.Marshal(matchReq)
		matchHTTPReq, _ := http.NewRequest("POST", server.URL+"/api/v1/matches", bytes.NewBuffer(matchBody))
		matchHTTPReq.Header.Set("Content-Type", "application/json")
		matchHTTPReq.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    sessionCookie,
			Path:     "/",
			HttpOnly: true,
		})

		matchResp, err := http.DefaultClient.Do(matchHTTPReq)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, matchResp.StatusCode)

		var matchRespData map[string]interface{}
		err = json.NewDecoder(matchResp.Body).Decode(&matchRespData)
		require.NoError(t, err)

		// Extract match ID from response
		matchData := matchRespData["data"].(map[string]interface{})
		matchID := matchData["id"].(string)

		// Start scoring for the match
		startReq := map[string]interface{}{
			"match_id": matchID,
		}
		startBody, _ := json.Marshal(startReq)
		startHTTPReq, _ := http.NewRequest("POST", server.URL+"/api/v1/scorecard/start", bytes.NewBuffer(startBody))
		startHTTPReq.Header.Set("Content-Type", "application/json")
		startHTTPReq.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    sessionCookie,
			Path:     "/",
			HttpOnly: true,
		})

		startResp, err := http.DefaultClient.Do(startHTTPReq)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, startResp.StatusCode)

		// Add ball to first innings (should succeed)
		ballEvent := &models.BallEventRequest{
			MatchID:       matchID,
			InningsNumber: 1, // First innings
			BallType:      models.BallTypeGood,
			RunType:       models.RunTypeFour,
			IsWicket:      false,
		}

		ballBody, _ := json.Marshal(ballEvent)
		ballHTTPReq, _ := http.NewRequest("POST", server.URL+"/api/v1/scorecard/ball", bytes.NewBuffer(ballBody))
		ballHTTPReq.Header.Set("Content-Type", "application/json")
		ballHTTPReq.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    sessionCookie,
			Path:     "/",
			HttpOnly: true,
		})

		ballResp, err := http.DefaultClient.Do(ballHTTPReq)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, ballResp.StatusCode)

		var successResp map[string]interface{}
		err = json.NewDecoder(ballResp.Body).Decode(&successResp)
		require.NoError(t, err)
		// The response is wrapped in a "data" field
		data := successResp["data"].(map[string]interface{})
		assert.Equal(t, "Ball added successfully", data["message"])
	})

	// Clean up after test
	testutils.CleanupTestData(t, db)
}
