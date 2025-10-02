package tests

import (
	"bytes"
	"context"
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

// Helper function to create an authenticated user and session cookie
func createAuthenticatedUserForUndoBall(t *testing.T, db *database.Client, sessionService *services.SessionService) (*models.User, string) {
	// Create a unique test user
	user := &models.User{
		GoogleID:      "test-google-id-" + time.Now().Format("20060102150405"),
		Email:         "test-undo-ball-" + time.Now().Format("20060102150405") + "@example.com",
		Name:          "Test Undo Ball User",
		Picture:       "https://example.com/picture.jpg",
		EmailVerified: true,
	}

	// Create user in database
	err := db.Repositories.User.CreateUser(context.Background(), user)
	require.NoError(t, err, "Failed to create test user")

	// Create a proper session using the session service
	req := httptest.NewRequest("POST", "/auth/login", nil)
	w := httptest.NewRecorder()

	err = sessionService.CreateSession(w, req, user)
	require.NoError(t, err, "Failed to create test session")

	// Extract the session cookie from the response
	cookies := w.Result().Cookies()
	var sessionCookie string
	for _, cookie := range cookies {
		if cookie.Name == "user_session" {
			sessionCookie = cookie.Value
			break
		}
	}
	require.NotEmpty(t, sessionCookie, "Session cookie not found in response")

	return user, sessionCookie
}

// Helper function to create a test series for undo ball tests
func createTestSeriesForUndoBall(t *testing.T, router http.Handler, sessionCookie string) string {
	seriesReq := map[string]interface{}{
		"name":        "Test Series " + time.Now().Format("2006-01-02 15:04:05"),
		"description": "Test series for undo ball integration tests",
		"start_date":  time.Now().AddDate(0, 0, 1).Format("2006-01-02T15:04:05Z"),
		"end_date":    time.Now().AddDate(0, 0, 7).Format("2006-01-02T15:04:05Z"),
	}

	body, err := json.Marshal(seriesReq)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/api/v1/series", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:     "user_session",
		Value:    sessionCookie,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
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

// Helper function to create a test match for undo ball tests
func createTestMatchForUndoBall(t *testing.T, router http.Handler, seriesID string, sessionCookie string) string {
	t.Logf("DEBUG: createTestMatchForUndoBall - Creating match for seriesID: %s", seriesID)
	t.Logf("DEBUG: createTestMatchForUndoBall - Using session cookie: %s", sessionCookie)

	matchReq := map[string]interface{}{
		"series_id":           seriesID,
		"match_number":        1,
		"date":                time.Now().AddDate(0, 0, 1).Format("2006-01-02T15:04:05Z"),
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
	req.AddCookie(&http.Cookie{
		Name:     "user_session",
		Value:    sessionCookie,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	t.Logf("DEBUG: createTestMatchForUndoBall - Response status: %d", w.Code)
	t.Logf("DEBUG: createTestMatchForUndoBall - Response body: %s", w.Body.String())

	require.Equal(t, http.StatusCreated, w.Code)

	var response struct {
		Data models.Match `json:"data"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	t.Logf("DEBUG: createTestMatchForUndoBall - Created match with ID: %s", response.Data.ID)
	return response.Data.ID
}

// Helper function to update match status to live
func updateMatchToLiveForUndoBall(t *testing.T, router http.Handler, matchID string, sessionCookie string) {
	updateReq := map[string]interface{}{
		"status": "live",
	}

	body, err := json.Marshal(updateReq)
	require.NoError(t, err)

	req := httptest.NewRequest("PUT", "/api/v1/matches/"+matchID, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:     "user_session",
		Value:    sessionCookie,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

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
	serviceContainer := services.NewContainer(testDB, cfg.Config)

	// Setup router with authentication
	router := handlers.SetupRoutes(testDB, cfg.Config)

	// Create authenticated user and session
	testUser, sessionCookie := createAuthenticatedUserForUndoBall(t, testDB, serviceContainer.SessionService)

	// Clean up test user at the end
	defer func() {
		_ = testDB.Repositories.User.DeleteUser(context.Background(), testUser.ID)
	}()

	t.Run("successful undo ball", func(t *testing.T) {
		// Create a test series first
		seriesID := createTestSeriesForUndoBall(t, router, sessionCookie)

		// Create a test match
		matchID := createTestMatchForUndoBall(t, router, seriesID, sessionCookie)

		// Update match to live status
		updateMatchToLiveForUndoBall(t, router, matchID, sessionCookie)

		// Start scoring
		startReq := models.ScorecardRequest{
			MatchID: matchID,
		}
		startBody, _ := json.Marshal(startReq)
		startReqObj := httptest.NewRequest("POST", "/api/v1/scorecard/start", bytes.NewBuffer(startBody))
		startReqObj.Header.Set("Content-Type", "application/json")
		startReqObj.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    sessionCookie,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		})
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
		ballReqObj1.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    sessionCookie,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		})
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
		ballReqObj2.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    sessionCookie,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		})
		ballW2 := httptest.NewRecorder()
		router.ServeHTTP(ballW2, ballReqObj2)
		assert.Equal(t, http.StatusOK, ballW2.Code)

		// Get scorecard before undo
		getReq := httptest.NewRequest("GET", "/api/v1/scorecard/"+matchID, nil)
		getReq.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    sessionCookie,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		})
		getW := httptest.NewRecorder()
		router.ServeHTTP(getW, getReq)
		assert.Equal(t, http.StatusOK, getW.Code)

		var responseBefore struct {
			Data models.ScorecardResponse `json:"data"`
		}
		err = json.Unmarshal(getW.Body.Bytes(), &responseBefore)
		require.NoError(t, err)
		require.NotEmpty(t, responseBefore.Data.Innings, "Scorecard should have innings data")
		require.Len(t, responseBefore.Data.Innings, 1, "Should have exactly one innings")
		require.NotNil(t, responseBefore.Data.Innings[0], "First innings should not be nil")
		assert.Equal(t, 3, responseBefore.Data.Innings[0].TotalRuns)  // 1 + 2 = 3 runs
		assert.Equal(t, 2, responseBefore.Data.Innings[0].TotalBalls) // 2 legal balls

		// Undo last ball
		t.Logf("DEBUG: successful_undo_ball test - About to undo ball for matchID: %s", matchID)
		t.Logf("DEBUG: successful_undo_ball test - Session cookie: %s", sessionCookie)

		undoReq := httptest.NewRequest("DELETE", "/api/v1/scorecard/"+matchID+"/ball?innings=1", nil)
		undoReq.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    sessionCookie,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		})
		undoW := httptest.NewRecorder()
		router.ServeHTTP(undoW, undoReq)

		t.Logf("DEBUG: successful_undo_ball test - Undo response status: %d", undoW.Code)
		t.Logf("DEBUG: successful_undo_ball test - Undo response body: %s", undoW.Body.String())
		assert.Equal(t, http.StatusOK, undoW.Code)

		// Get scorecard after undo
		getReq2 := httptest.NewRequest("GET", "/api/v1/scorecard/"+matchID, nil)
		getReq2.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    sessionCookie,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		})
		getW2 := httptest.NewRecorder()
		router.ServeHTTP(getW2, getReq2)
		assert.Equal(t, http.StatusOK, getW2.Code)

		var responseAfter struct {
			Data models.ScorecardResponse `json:"data"`
		}
		err = json.Unmarshal(getW2.Body.Bytes(), &responseAfter)
		require.NoError(t, err)
		require.NotEmpty(t, responseAfter.Data.Innings, "Scorecard should have innings data after undo")
		require.Len(t, responseAfter.Data.Innings, 1, "Should have exactly one innings after undo")
		require.NotNil(t, responseAfter.Data.Innings[0], "First innings should not be nil after undo")
		require.NotEmpty(t, responseAfter.Data.Innings[0].Overs, "Innings should have overs data after undo")
		require.Len(t, responseAfter.Data.Innings[0].Overs, 1, "Should have exactly one over after undo")
		require.NotNil(t, responseAfter.Data.Innings[0].Overs[0], "First over should not be nil after undo")
		assert.Equal(t, 1, responseAfter.Data.Innings[0].TotalRuns)           // Only 1 run left
		assert.Equal(t, 1, responseAfter.Data.Innings[0].TotalBalls)          // Only 1 legal ball left
		assert.Equal(t, 1, len(responseAfter.Data.Innings[0].Overs[0].Balls)) // Only 1 ball in over
	})

	t.Run("undo ball - no balls to undo", func(t *testing.T) {
		// Create a test series first
		seriesID := createTestSeriesForUndoBall(t, router, sessionCookie)

		// Create a test match
		matchID := createTestMatchForUndoBall(t, router, seriesID, sessionCookie)

		// Update match to live status
		updateMatchToLiveForUndoBall(t, router, matchID, sessionCookie)

		// Start scoring
		startReq := models.ScorecardRequest{
			MatchID: matchID,
		}
		startBody, _ := json.Marshal(startReq)
		startReqObj := httptest.NewRequest("POST", "/api/v1/scorecard/start", bytes.NewBuffer(startBody))
		startReqObj.Header.Set("Content-Type", "application/json")
		startReqObj.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    sessionCookie,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		})
		startW := httptest.NewRecorder()
		router.ServeHTTP(startW, startReqObj)
		assert.Equal(t, http.StatusOK, startW.Code)

		// Try to undo ball when no balls exist
		undoReq := httptest.NewRequest("DELETE", "/api/v1/scorecard/"+matchID+"/ball?innings=1", nil)
		undoReq.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    sessionCookie,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		})
		undoW := httptest.NewRecorder()
		router.ServeHTTP(undoW, undoReq)
		assert.Equal(t, http.StatusInternalServerError, undoW.Code)

		var errorResponse struct {
			Error struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		err = json.Unmarshal(undoW.Body.Bytes(), &errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse.Error.Message, "no current over found")
	})

	t.Run("undo ball - match not found", func(t *testing.T) {
		// Try to undo ball for non-existent match
		undoReq := httptest.NewRequest("DELETE", "/api/v1/scorecard/nonexistent-match/ball?innings=1", nil)
		undoReq.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    sessionCookie,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		})
		undoW := httptest.NewRecorder()
		router.ServeHTTP(undoW, undoReq)
		assert.Equal(t, http.StatusInternalServerError, undoW.Code)

		var errorResponse struct {
			Error struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		err := json.Unmarshal(undoW.Body.Bytes(), &errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse.Error.Message, "match not found")
	})

	t.Run("undo ball - invalid innings number", func(t *testing.T) {
		// Create a test series first
		seriesID := createTestSeriesForUndoBall(t, router, sessionCookie)

		// Create a test match
		matchID := createTestMatchForUndoBall(t, router, seriesID, sessionCookie)

		// Try to undo ball with invalid innings number
		undoReq := httptest.NewRequest("DELETE", "/api/v1/scorecard/"+matchID+"/ball?innings=3", nil)
		undoReq.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    sessionCookie,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		})
		undoW := httptest.NewRecorder()
		router.ServeHTTP(undoW, undoReq)
		assert.Equal(t, http.StatusBadRequest, undoW.Code)

		var errorResponse struct {
			Error struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		err := json.Unmarshal(undoW.Body.Bytes(), &errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse.Error.Message, "innings must be 1 or 2")
	})
}
