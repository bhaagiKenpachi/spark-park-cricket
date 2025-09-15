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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUndoBallIntegration(t *testing.T) {
	// Load test configuration
	cfg := config.LoadTestConfig()

	// Initialize test database
	testDB, err := database.NewTestClient(cfg)
	require.NoError(t, err)
	defer testDB.Close()

	// Setup test schema
	err = database.SetupTestSchema(cfg)
	require.NoError(t, err)
	defer testutils.CleanupScorecardTestData(t, testDB)

	// Setup service container and handlers
	serviceContainer := services.NewContainer(testDB.Repositories)
	scorecardHandler := handlers.NewScorecardHandler(serviceContainer.Scorecard)

	// Setup router using testutils
	router := testutils.SetupScorecardTestRouter(scorecardHandler, serviceContainer)

	t.Run("successful undo ball", func(t *testing.T) {
		// Create a test series first
		seriesID := testutils.CreateTestSeriesForWorkflow(t, router)

		// Create a test match
		matchID := testutils.CreateTestMatchForWorkflow(t, router, seriesID)

		// Update match to live status
		testutils.UpdateMatchToLiveForWorkflow(t, router, matchID)

		// Start scoring
		startReq := models.ScorecardRequest{
			MatchID: matchID,
		}
		startBody, _ := json.Marshal(startReq)
		startReqObj := httptest.NewRequest("POST", "/api/v1/scorecard/start", bytes.NewBuffer(startBody))
		startReqObj.Header.Set("Content-Type", "application/json")
		startW := httptest.NewRecorder()
		router.ServeHTTP(startW, startReqObj)
		assert.Equal(t, http.StatusOK, startW.Code)

		// Add first ball
		ballReq1 := models.BallEventRequest{
			MatchID:       matchID,
			InningsNumber: 1,
			BallType:      models.BallTypeGood,
			RunType:       models.RunTypeOne,
			IsWicket:      false,
		}
		ballBody1, _ := json.Marshal(ballReq1)
		ballReqObj1 := httptest.NewRequest("POST", "/api/v1/scorecard/ball", bytes.NewBuffer(ballBody1))
		ballReqObj1.Header.Set("Content-Type", "application/json")
		ballW1 := httptest.NewRecorder()
		router.ServeHTTP(ballW1, ballReqObj1)
		assert.Equal(t, http.StatusOK, ballW1.Code)

		// Add second ball
		ballReq2 := models.BallEventRequest{
			MatchID:       matchID,
			InningsNumber: 1,
			BallType:      models.BallTypeGood,
			RunType:       models.RunTypeTwo,
			IsWicket:      false,
		}
		ballBody2, _ := json.Marshal(ballReq2)
		ballReqObj2 := httptest.NewRequest("POST", "/api/v1/scorecard/ball", bytes.NewBuffer(ballBody2))
		ballReqObj2.Header.Set("Content-Type", "application/json")
		ballW2 := httptest.NewRecorder()
		router.ServeHTTP(ballW2, ballReqObj2)
		assert.Equal(t, http.StatusOK, ballW2.Code)

		// Get scorecard before undo
		getReq := httptest.NewRequest("GET", "/api/v1/scorecard/"+matchID, nil)
		getW := httptest.NewRecorder()
		router.ServeHTTP(getW, getReq)
		assert.Equal(t, http.StatusOK, getW.Code)

		var scorecardBefore models.ScorecardResponse
		err = json.Unmarshal(getW.Body.Bytes(), &scorecardBefore)
		require.NoError(t, err)
		assert.Equal(t, 3, scorecardBefore.Innings[0].TotalRuns)  // 1 + 2 = 3 runs
		assert.Equal(t, 2, scorecardBefore.Innings[0].TotalBalls) // 2 legal balls

		// Undo last ball
		undoReq := httptest.NewRequest("DELETE", "/api/v1/scorecard/"+matchID+"/ball?innings=1", nil)
		undoW := httptest.NewRecorder()
		router.ServeHTTP(undoW, undoReq)
		assert.Equal(t, http.StatusOK, undoW.Code)

		// Get scorecard after undo
		getReq2 := httptest.NewRequest("GET", "/api/v1/scorecard/"+matchID, nil)
		getW2 := httptest.NewRecorder()
		router.ServeHTTP(getW2, getReq2)
		assert.Equal(t, http.StatusOK, getW2.Code)

		var scorecardAfter models.ScorecardResponse
		err = json.Unmarshal(getW2.Body.Bytes(), &scorecardAfter)
		require.NoError(t, err)
		assert.Equal(t, 1, scorecardAfter.Innings[0].TotalRuns)           // Only 1 run left
		assert.Equal(t, 1, scorecardAfter.Innings[0].TotalBalls)          // Only 1 legal ball left
		assert.Equal(t, 1, len(scorecardAfter.Innings[0].Overs[0].Balls)) // Only 1 ball in over
	})

	t.Run("undo ball - no balls to undo", func(t *testing.T) {
		// Create a test series first
		seriesID := testutils.CreateTestSeriesForWorkflow(t, router)

		// Create a test match
		matchID := testutils.CreateTestMatchForWorkflow(t, router, seriesID)

		// Update match to live status
		testutils.UpdateMatchToLiveForWorkflow(t, router, matchID)

		// Start scoring
		startReq := models.ScorecardRequest{
			MatchID: matchID,
		}
		startBody, _ := json.Marshal(startReq)
		startReqObj := httptest.NewRequest("POST", "/api/v1/scorecard/start", bytes.NewBuffer(startBody))
		startReqObj.Header.Set("Content-Type", "application/json")
		startW := httptest.NewRecorder()
		router.ServeHTTP(startW, startReqObj)
		assert.Equal(t, http.StatusOK, startW.Code)

		// Try to undo ball when no balls exist
		undoReq := httptest.NewRequest("DELETE", "/api/v1/scorecard/"+matchID+"/ball?innings=1", nil)
		undoW := httptest.NewRecorder()
		router.ServeHTTP(undoW, undoReq)
		assert.Equal(t, http.StatusInternalServerError, undoW.Code)

		var errorResponse map[string]interface{}
		err = json.Unmarshal(undoW.Body.Bytes(), &errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse["message"], "no balls to undo")
	})

	t.Run("undo ball - match not found", func(t *testing.T) {
		// Try to undo ball for non-existent match
		undoReq := httptest.NewRequest("DELETE", "/api/v1/scorecard/nonexistent-match/ball?innings=1", nil)
		undoW := httptest.NewRecorder()
		router.ServeHTTP(undoW, undoReq)
		assert.Equal(t, http.StatusInternalServerError, undoW.Code)

		var errorResponse map[string]interface{}
		err := json.Unmarshal(undoW.Body.Bytes(), &errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse["message"], "match not found")
	})

	t.Run("undo ball - invalid innings number", func(t *testing.T) {
		// Create a test series first
		seriesID := testutils.CreateTestSeriesForWorkflow(t, router)

		// Create a test match
		matchID := testutils.CreateTestMatchForWorkflow(t, router, seriesID)

		// Try to undo ball with invalid innings number
		undoReq := httptest.NewRequest("DELETE", "/api/v1/scorecard/"+matchID+"/ball?innings=3", nil)
		undoW := httptest.NewRecorder()
		router.ServeHTTP(undoW, undoReq)
		assert.Equal(t, http.StatusBadRequest, undoW.Code)

		var errorResponse map[string]interface{}
		err := json.Unmarshal(undoW.Body.Bytes(), &errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse["message"], "innings must be 1 or 2")
	})
}
