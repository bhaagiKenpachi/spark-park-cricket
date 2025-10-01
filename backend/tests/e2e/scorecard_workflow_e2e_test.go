package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/pkg/testutils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompleteScorecardWorkflow(t *testing.T) {
	t.Run("CompleteMatchWorkflow", func(t *testing.T) {
		// Setup with authentication
		server, db, user, sessionCookie := testutils.SetupAuthenticatedE2ETestServerWithDB(t)
		defer server.Close()
		defer db.Close()

		// Clean up before test
		t.Logf("DEBUG: Starting CompleteMatchWorkflow test")
		testutils.CleanupTestDataForUser(t, db, user.ID)

		// Create a test series via API
		seriesReq := &models.CreateSeriesRequest{
			Name:      "E2E Test Series " + time.Now().Format("2006-01-02 15:04:05"),
			StartDate: time.Now().AddDate(0, 0, 1),
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
			Date:             time.Now().AddDate(0, 0, 1),
			TeamAPlayerCount: 11,
			TeamBPlayerCount: 11,
			TotalOvers:       2, // Short match for testing
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

		// Log response details for debugging
		if startResp.StatusCode != http.StatusOK {
			t.Logf("DEBUG: Start scoring failed with status %d", startResp.StatusCode)
			body, _ := io.ReadAll(startResp.Body)
			t.Logf("DEBUG: Response body: %s", string(body))
		}

		require.Equal(t, http.StatusOK, startResp.StatusCode)

		// Add balls to complete the match
		for i := 1; i <= 12; i++ { // 2 overs * 6 balls = 12 balls
			ballEvent := &models.BallEventRequest{
				MatchID:       matchID,
				InningsNumber: 1,
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

		// Get final scorecard
		scorecardHTTPReq, _ := http.NewRequest("GET", server.URL+"/api/v1/scorecard/"+matchID, nil)
		scorecardHTTPReq.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    sessionCookie,
			Path:     "/",
			HttpOnly: true,
		})

		scorecardResp, err := http.DefaultClient.Do(scorecardHTTPReq)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, scorecardResp.StatusCode)

		var scorecardData map[string]interface{}
		err = json.NewDecoder(scorecardResp.Body).Decode(&scorecardData)
		require.NoError(t, err)

		// Verify scorecard structure
		assert.Contains(t, scorecardData, "data")
		scorecard := scorecardData["data"].(map[string]interface{})
		assert.Contains(t, scorecard, "match_id")
		assert.Contains(t, scorecard, "innings")

		// Clean up after test
		testutils.CleanupTestDataForUser(t, db, user.ID)
	})

	t.Run("MultipleOversWorkflow", func(t *testing.T) {
		// Setup with authentication
		server, db, user, sessionCookie := testutils.SetupAuthenticatedE2ETestServerWithDB(t)
		defer server.Close()
		defer db.Close()

		// Clean up before test
		testutils.CleanupTestDataForUser(t, db, user.ID)

		// Create a test series via API
		seriesReq := &models.CreateSeriesRequest{
			Name:      "E2E Test Series " + time.Now().Format("2006-01-02 15:04:05"),
			StartDate: time.Now().AddDate(0, 0, 1),
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
			Date:             time.Now().AddDate(0, 0, 1),
			TeamAPlayerCount: 11,
			TeamBPlayerCount: 11,
			TotalOvers:       3, // Multiple overs for testing
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

		// Log response details for debugging
		if startResp.StatusCode != http.StatusOK {
			t.Logf("DEBUG: Start scoring failed with status %d", startResp.StatusCode)
			body, _ := io.ReadAll(startResp.Body)
			t.Logf("DEBUG: Response body: %s", string(body))
		}

		require.Equal(t, http.StatusOK, startResp.StatusCode)

		// Add balls across multiple overs
		for over := 1; over <= 3; over++ {
			for ball := 1; ball <= 6; ball++ {
				ballEvent := &models.BallEventRequest{
					MatchID:       matchID,
					InningsNumber: 1,
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

		// Get final scorecard
		scorecardHTTPReq, _ := http.NewRequest("GET", server.URL+"/api/v1/scorecard/"+matchID, nil)
		scorecardHTTPReq.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    sessionCookie,
			Path:     "/",
			HttpOnly: true,
		})

		scorecardResp, err := http.DefaultClient.Do(scorecardHTTPReq)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, scorecardResp.StatusCode)

		var scorecardData map[string]interface{}
		err = json.NewDecoder(scorecardResp.Body).Decode(&scorecardData)
		require.NoError(t, err)

		// Verify scorecard structure
		assert.Contains(t, scorecardData, "data")
		scorecard := scorecardData["data"].(map[string]interface{})
		assert.Contains(t, scorecard, "match_id")
		assert.Contains(t, scorecard, "innings")

		// Clean up after test
		testutils.CleanupTestDataForUser(t, db, user.ID)
	})

	t.Run("WideAndNoBallWorkflow", func(t *testing.T) {
		// Setup with authentication
		server, db, user, sessionCookie := testutils.SetupAuthenticatedE2ETestServerWithDB(t)
		defer server.Close()
		defer db.Close()

		// Clean up before test
		testutils.CleanupTestDataForUser(t, db, user.ID)

		// Create a test series via API
		seriesReq := &models.CreateSeriesRequest{
			Name:      "E2E Test Series " + time.Now().Format("2006-01-02 15:04:05"),
			StartDate: time.Now().AddDate(0, 0, 1),
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
			Date:             time.Now().AddDate(0, 0, 1),
			TeamAPlayerCount: 11,
			TeamBPlayerCount: 11,
			TotalOvers:       1, // Single over for testing wides and no-balls
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

		// Log response details for debugging
		if startResp.StatusCode != http.StatusOK {
			t.Logf("DEBUG: Start scoring failed with status %d", startResp.StatusCode)
			body, _ := io.ReadAll(startResp.Body)
			t.Logf("DEBUG: Response body: %s", string(body))
		}

		require.Equal(t, http.StatusOK, startResp.StatusCode)

		// Add various types of balls
		ballTypes := []struct {
			ballType models.BallType
			runType  models.RunType
		}{
			{models.BallTypeGood, models.RunTypeOne},
			{models.BallTypeWide, models.RunTypeOne},
			{models.BallTypeGood, models.RunTypeTwo},
			{models.BallTypeNoBall, models.RunTypeFour},
			{models.BallTypeGood, models.RunTypeSix},
			{models.BallTypeGood, models.RunTypeZero},
		}

		for _, ball := range ballTypes {
			ballEvent := &models.BallEventRequest{
				MatchID:       matchID,
				InningsNumber: 1,
				BallType:      ball.ballType,
				RunType:       ball.runType,
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

		// Get final scorecard
		scorecardHTTPReq, _ := http.NewRequest("GET", server.URL+"/api/v1/scorecard/"+matchID, nil)
		scorecardHTTPReq.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    sessionCookie,
			Path:     "/",
			HttpOnly: true,
		})

		scorecardResp, err := http.DefaultClient.Do(scorecardHTTPReq)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, scorecardResp.StatusCode)

		var scorecardData map[string]interface{}
		err = json.NewDecoder(scorecardResp.Body).Decode(&scorecardData)
		require.NoError(t, err)

		// Verify scorecard structure
		assert.Contains(t, scorecardData, "data")
		scorecard := scorecardData["data"].(map[string]interface{})
		assert.Contains(t, scorecard, "match_id")
		assert.Contains(t, scorecard, "innings")

		// Clean up after test
		testutils.CleanupTestDataForUser(t, db, user.ID)
	})
}
