package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"spark-park-cricket-backend/internal/config"
	"spark-park-cricket-backend/internal/database"
	"spark-park-cricket-backend/internal/handlers"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/services"
	"spark-park-cricket-backend/pkg/testutils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Note: Helper functions are defined in match_completion_unit_test.go

// Helper function to create a test series for scorecard tests
func createTestSeriesForScorecard(t *testing.T, router http.Handler) string {
	seriesReq := map[string]interface{}{
		"name":        "Test Series " + time.Now().Format("2006-01-02 15:04:05"),
		"description": "Test series for scorecard integration tests",
		"start_date":  time.Now().AddDate(0, 0, 1).Format(time.RFC3339),
		"end_date":    time.Now().AddDate(0, 0, 7).Format(time.RFC3339),
	}

	body, err := json.Marshal(seriesReq)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/api/v1/series", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var response struct {
		Data models.Series `json:"data"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	return response.Data.ID
}

// Helper function to create a test match for scorecard tests
func createTestMatchForScorecard(t *testing.T, router http.Handler, seriesID string) string {
	matchReq := map[string]interface{}{
		"series_id":           seriesID,
		"match_number":        1,
		"date":                time.Now().AddDate(0, 0, 1).Format(time.RFC3339),
		"venue":               "Test Venue",
		"team_a_player_count": 11,
		"team_b_player_count": 11,
		"total_overs":         20,
		"toss_winner":         "A",
		"toss_type":           "H",
		"batting_team":        "A",
	}

	body, err := json.Marshal(matchReq)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/api/v1/matches", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var response struct {
		Data models.Match `json:"data"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	return response.Data.ID
}

// Helper function to update match status to live
func updateMatchToLive(t *testing.T, router http.Handler, matchID string) {
	updateReq := map[string]interface{}{
		"status": "live",
	}

	body, err := json.Marshal(updateReq)
	require.NoError(t, err)

	req := httptest.NewRequest("PUT", "/api/v1/matches/"+matchID, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestScorecardIntegration(t *testing.T) {
	// Load test configuration
	testConfig := config.LoadTestConfig()

	// Initialize database client
	dbClient, err := database.NewClient(testConfig.Config)
	require.NoError(t, err)
	defer dbClient.Close()

	// Clean up any existing test data
	testutils.CleanupScorecardTestData(t, dbClient)

	// Initialize services
	serviceContainer := services.NewContainer(dbClient.Repositories, testConfig.Config)
	scorecardHandler := handlers.NewScorecardHandler(serviceContainer.Scorecard)

	// Setup router
	router := testutils.SetupScorecardTestRouter(scorecardHandler, serviceContainer)

	// Create test series
	seriesID := createTestSeriesForScorecard(t, router)

	// Create test match
	matchID := createTestMatchForScorecard(t, router, seriesID)

	// Update match to live status
	updateMatchToLive(t, router, matchID)

	t.Run("StartScoring", func(t *testing.T) {
		req := map[string]interface{}{
			"match_id": matchID,
		}

		body, err := json.Marshal(req)
		require.NoError(t, err)

		reqHTTP := httptest.NewRequest("POST", "/api/v1/scorecard/start", bytes.NewBuffer(body))
		reqHTTP.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, reqHTTP)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		data, ok := response["data"].(map[string]interface{})
		require.True(t, ok, "Response should contain data field")
		assert.Equal(t, "Scoring started successfully", data["message"])
		assert.Equal(t, matchID, data["match_id"])
	})

	t.Run("StartScoring_AlreadyStarted", func(t *testing.T) {
		req := map[string]interface{}{
			"match_id": matchID,
		}

		body, err := json.Marshal(req)
		require.NoError(t, err)

		reqHTTP := httptest.NewRequest("POST", "/api/v1/scorecard/start", bytes.NewBuffer(body))
		reqHTTP.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, reqHTTP)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("AddBall", func(t *testing.T) {
		req := map[string]interface{}{
			"match_id":       matchID,
			"innings_number": 1,
			"ball_type":      "good",
			"run_type":       "4",
			"is_wicket":      false,
			"byes":           0,
		}

		body, err := json.Marshal(req)
		require.NoError(t, err)

		reqHTTP := httptest.NewRequest("POST", "/api/v1/scorecard/ball", bytes.NewBuffer(body))
		reqHTTP.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, reqHTTP)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		data, ok := response["data"].(map[string]interface{})
		require.True(t, ok, "Response should contain data field")
		assert.Equal(t, "Ball added successfully", data["message"])
		assert.Equal(t, matchID, data["match_id"])
		assert.Equal(t, 1, int(data["innings_number"].(float64)))
		assert.Equal(t, "good", data["ball_type"])
		assert.Equal(t, "4", data["run_type"])
		assert.Equal(t, 4, int(data["runs"].(float64)))
	})

	t.Run("AddBall_InvalidRequest", func(t *testing.T) {
		req := map[string]interface{}{
			"match_id":       matchID,
			"innings_number": 1,
			"ball_type":      "invalid_type",
			"run_type":       "4",
			"is_wicket":      false,
			"byes":           0,
		}

		body, err := json.Marshal(req)
		require.NoError(t, err)

		reqHTTP := httptest.NewRequest("POST", "/api/v1/scorecard/ball", bytes.NewBuffer(body))
		reqHTTP.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, reqHTTP)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("GetScorecard", func(t *testing.T) {
		reqHTTP := httptest.NewRequest("GET", "/api/v1/scorecard/"+matchID, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, reqHTTP)
		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Data models.ScorecardResponse `json:"data"`
		}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, matchID, response.Data.MatchID)
		assert.Len(t, response.Data.Innings, 1)
		assert.Equal(t, 1, response.Data.Innings[0].InningsNumber)
		assert.Equal(t, 4, response.Data.Innings[0].TotalRuns)
	})

	t.Run("GetScorecard_NotFound", func(t *testing.T) {
		reqHTTP := httptest.NewRequest("GET", "/api/v1/scorecard/nonexistent-match", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, reqHTTP)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("GetCurrentOver", func(t *testing.T) {
		reqHTTP := httptest.NewRequest("GET", "/api/v1/scorecard/"+matchID+"/current-over?innings=1", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, reqHTTP)
		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Data models.ScorecardOver `json:"data"`
		}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, 1, response.Data.OverNumber)
		assert.Equal(t, 4, response.Data.TotalRuns)
		assert.Equal(t, 1, response.Data.TotalBalls)
	})

	t.Run("GetCurrentOver_InvalidInnings", func(t *testing.T) {
		reqHTTP := httptest.NewRequest("GET", "/api/v1/scorecard/"+matchID+"/current-over?innings=3", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, reqHTTP)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("GetInnings", func(t *testing.T) {
		reqHTTP := httptest.NewRequest("GET", "/api/v1/scorecard/"+matchID+"/innings/1", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, reqHTTP)
		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Data models.InningsSummary `json:"data"`
		}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, 1, response.Data.InningsNumber)
		assert.Equal(t, 4, response.Data.TotalRuns)
		assert.Equal(t, 1, response.Data.TotalBalls)
	})

	t.Run("GetInnings_NotFound", func(t *testing.T) {
		reqHTTP := httptest.NewRequest("GET", "/api/v1/scorecard/"+matchID+"/innings/2", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, reqHTTP)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("GetOver", func(t *testing.T) {
		reqHTTP := httptest.NewRequest("GET", "/api/v1/scorecard/"+matchID+"/innings/1/over/1", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, reqHTTP)
		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Data models.OverSummary `json:"data"`
		}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, 1, response.Data.OverNumber)
		assert.Equal(t, 4, response.Data.TotalRuns)
		assert.Equal(t, 1, response.Data.TotalBalls)
	})

	t.Run("GetOver_NotFound", func(t *testing.T) {
		reqHTTP := httptest.NewRequest("GET", "/api/v1/scorecard/"+matchID+"/innings/1/over/2", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, reqHTTP)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	// Clean up test data
	testutils.CleanupScorecardTestData(t, dbClient)
}
