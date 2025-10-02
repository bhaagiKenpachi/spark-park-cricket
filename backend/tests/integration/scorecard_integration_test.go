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
func createAuthenticatedUserForScorecard(t *testing.T, db *database.Client, sessionService *services.SessionService) (*models.User, string) {
	// Create a unique test user
	user := &models.User{
		GoogleID:      "test-google-id-" + time.Now().Format("20060102150405"),
		Email:         "test-scorecard-" + time.Now().Format("20060102150405") + "@example.com",
		Name:          "Test Scorecard User",
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

// Helper function to create a test series for scorecard tests
func createTestSeriesForScorecard(t *testing.T, router http.Handler, sessionCookie string) string {
	seriesReq := map[string]interface{}{
		"name":        "Test Series " + time.Now().Format("2006-01-02 15:04:05"),
		"description": "Test series for scorecard integration tests",
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

// Helper function to create a test match for scorecard tests
func createTestMatchForScorecard(t *testing.T, router http.Handler, seriesID string, sessionCookie string) string {
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
	require.Equal(t, http.StatusCreated, w.Code)

	var response struct {
		Data models.Match `json:"data"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	return response.Data.ID
}

// Helper function to update match status to live
func updateMatchToLive(t *testing.T, router http.Handler, matchID string, sessionCookie string) {
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
	serviceContainer := services.NewContainer(dbClient, testConfig.Config)

	// Setup router with authentication
	router := handlers.SetupRoutes(dbClient, testConfig.Config)

	// Create authenticated user and session
	testUser, sessionCookie := createAuthenticatedUserForScorecard(t, dbClient, serviceContainer.SessionService)

	// Clean up test user at the end
	defer func() {
		_ = dbClient.Repositories.User.DeleteUser(context.Background(), testUser.ID)
	}()

	// Create test series
	seriesID := createTestSeriesForScorecard(t, router, sessionCookie)

	// Create test match
	matchID := createTestMatchForScorecard(t, router, seriesID, sessionCookie)

	// Update match to live status
	updateMatchToLive(t, router, matchID, sessionCookie)

	t.Run("StartScoring", func(t *testing.T) {
		req := map[string]interface{}{
			"match_id": matchID,
		}

		body, err := json.Marshal(req)
		require.NoError(t, err)

		reqHTTP := httptest.NewRequest("POST", "/api/v1/scorecard/start", bytes.NewBuffer(body))
		reqHTTP.Header.Set("Content-Type", "application/json")
		reqHTTP.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    sessionCookie,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		})
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
		reqHTTP.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    sessionCookie,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		})
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
		reqHTTP.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    sessionCookie,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		})
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
		reqHTTP.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    sessionCookie,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		})
		w := httptest.NewRecorder()

		router.ServeHTTP(w, reqHTTP)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("GetScorecard", func(t *testing.T) {
		reqHTTP := httptest.NewRequest("GET", "/api/v1/scorecard/"+matchID, nil)
		reqHTTP.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    sessionCookie,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		})
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
		require.NotNil(t, response.Data.Innings[0], "First innings should not be nil")
		assert.Equal(t, 1, response.Data.Innings[0].InningsNumber)
		assert.Equal(t, 4, response.Data.Innings[0].TotalRuns)
	})

	t.Run("GetScorecard_NotFound", func(t *testing.T) {
		reqHTTP := httptest.NewRequest("GET", "/api/v1/scorecard/nonexistent-match", nil)
		reqHTTP.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    sessionCookie,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		})
		w := httptest.NewRecorder()

		router.ServeHTTP(w, reqHTTP)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("GetCurrentOver", func(t *testing.T) {
		reqHTTP := httptest.NewRequest("GET", "/api/v1/scorecard/"+matchID+"/current-over?innings=1", nil)
		reqHTTP.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    sessionCookie,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		})
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
		reqHTTP.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    sessionCookie,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		})
		w := httptest.NewRecorder()

		router.ServeHTTP(w, reqHTTP)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("GetInnings", func(t *testing.T) {
		reqHTTP := httptest.NewRequest("GET", "/api/v1/scorecard/"+matchID+"/innings/1", nil)
		reqHTTP.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    sessionCookie,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		})
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
		reqHTTP.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    sessionCookie,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		})
		w := httptest.NewRecorder()

		router.ServeHTTP(w, reqHTTP)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("GetOver", func(t *testing.T) {
		reqHTTP := httptest.NewRequest("GET", "/api/v1/scorecard/"+matchID+"/innings/1/over/1", nil)
		reqHTTP.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    sessionCookie,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		})
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
		reqHTTP.AddCookie(&http.Cookie{
			Name:     "user_session",
			Value:    sessionCookie,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		})
		w := httptest.NewRecorder()

		router.ServeHTTP(w, reqHTTP)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	// Clean up test data
	testutils.CleanupScorecardTestData(t, dbClient)
}
