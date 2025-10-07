package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"spark-park-cricket-backend/internal/config"
	"spark-park-cricket-backend/internal/database"
	"spark-park-cricket-backend/internal/handlers"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVoteIntegration(t *testing.T) {
	// Skip if not in integration test mode
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Setup test database
	cfg := &config.Config{
		SupabaseURL:    "http://localhost:54321",
		SupabaseAPIKey: "test-key",
		DatabaseSchema: "testing_db",
	}

	dbClient, err := database.NewClient(cfg)
	require.NoError(t, err)
	defer dbClient.Close()

	// Setup services
	serviceContainer := services.NewContainer(dbClient, cfg)

	// Setup router
	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		voteHandler := handlers.NewVoteHandler(serviceContainer.VoteService)
		r.Route("/votes", func(r chi.Router) {
			r.Post("/", voteHandler.CreateVote)
			r.Get("/{id}", voteHandler.GetVote)
			r.Get("/{id}/results", voteHandler.GetVoteWithResults)
			r.With(services.AuthMiddleware(serviceContainer.SessionService)).Post("/", voteHandler.CreateVote)
			r.With(services.AuthMiddleware(serviceContainer.SessionService)).Post("/{id}/vote", voteHandler.CastVote)
		})
	})

	// Test data
	userID := "test-user-123"
	voteRequest := models.CreateVoteRequest{
		Title:       "Test Vote",
		Description: "This is a test vote for integration testing",
		Type:        models.VoteTypeSingle,
		Options:     []string{"Option A", "Option B", "Option C"},
	}

	t.Run("Create Vote", func(t *testing.T) {
		// Create a vote
		reqBody, _ := json.Marshal(voteRequest)
		req := httptest.NewRequest("POST", "/api/v1/votes", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		// Add user context (simulating authenticated user)
		ctx := context.WithValue(req.Context(), "user_id", userID)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Contains(t, response, "data")
		data := response["data"].(map[string]interface{})
		assert.Contains(t, data, "vote")
		vote := data["vote"].(map[string]interface{})
		assert.Equal(t, voteRequest.Title, vote["title"])
		assert.Equal(t, string(voteRequest.Type), vote["type"])
	})

	t.Run("Get Vote", func(t *testing.T) {
		// First create a vote
		reqBody, _ := json.Marshal(voteRequest)
		req := httptest.NewRequest("POST", "/api/v1/votes", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), "user_id", userID)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		var createResponse map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &createResponse)
		require.NoError(t, err)

		vote := createResponse["data"].(map[string]interface{})["vote"].(map[string]interface{})
		voteID := vote["id"].(string)

		// Now get the vote
		req = httptest.NewRequest("GET", fmt.Sprintf("/api/v1/votes/%s", voteID), nil)
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var getResponse map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &getResponse)
		require.NoError(t, err)

		assert.Contains(t, getResponse, "data")
		data := getResponse["data"].(map[string]interface{})
		assert.Contains(t, data, "vote")
		assert.Contains(t, data, "options")

		retrievedVote := data["vote"].(map[string]interface{})
		assert.Equal(t, voteRequest.Title, retrievedVote["title"])
		assert.Equal(t, string(voteRequest.Type), retrievedVote["type"])

		options := data["options"].([]interface{})
		assert.Len(t, options, 3)
	})

	t.Run("Get Vote Results", func(t *testing.T) {
		// First create a vote
		reqBody, _ := json.Marshal(voteRequest)
		req := httptest.NewRequest("POST", "/api/v1/votes", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), "user_id", userID)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		var createResponse map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &createResponse)
		require.NoError(t, err)

		vote := createResponse["data"].(map[string]interface{})["vote"].(map[string]interface{})
		voteID := vote["id"].(string)

		// Get vote results
		req = httptest.NewRequest("GET", fmt.Sprintf("/api/v1/votes/%s/results", voteID), nil)
		w = httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resultsResponse map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &resultsResponse)
		require.NoError(t, err)

		assert.Contains(t, resultsResponse, "data")
		data := resultsResponse["data"].(map[string]interface{})
		assert.Contains(t, data, "vote")
		assert.Contains(t, data, "options")
		assert.Contains(t, data, "results")
		assert.Contains(t, data, "total_votes")

		results := data["results"].(map[string]interface{})
		totalVotes := data["total_votes"].(float64)
		assert.Equal(t, float64(0), totalVotes) // No votes yet

		// Check that all options have 0 votes initially
		for _, option := range data["options"].([]interface{}) {
			optionID := option.(map[string]interface{})["id"].(string)
			assert.Equal(t, float64(0), results[optionID])
		}
	})
}

func TestVoteValidation(t *testing.T) {
	// Skip if not in integration test mode
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	cfg := &config.Config{
		SupabaseURL:    "http://localhost:54321",
		SupabaseAPIKey: "test-key",
		DatabaseSchema: "testing_db",
	}

	dbClient, err := database.NewClient(cfg)
	require.NoError(t, err)
	defer dbClient.Close()

	serviceContainer := services.NewContainer(dbClient, cfg)

	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		voteHandler := handlers.NewVoteHandler(serviceContainer.VoteService)
		r.Route("/votes", func(r chi.Router) {
			r.With(services.AuthMiddleware(serviceContainer.SessionService)).Post("/", voteHandler.CreateVote)
		})
	})

	tests := []struct {
		name          string
		request       models.CreateVoteRequest
		expectedError string
	}{
		{
			name: "Empty title",
			request: models.CreateVoteRequest{
				Title:       "",
				Description: "Valid description",
				Type:        models.VoteTypeSingle,
				Options:     []string{"Option 1", "Option 2"},
			},
			expectedError: "title is required",
		},
		{
			name: "Short title",
			request: models.CreateVoteRequest{
				Title:       "AB",
				Description: "Valid description",
				Type:        models.VoteTypeSingle,
				Options:     []string{"Option 1", "Option 2"},
			},
			expectedError: "title must be between 3 and 255 characters",
		},
		{
			name: "Empty description",
			request: models.CreateVoteRequest{
				Title:       "Valid Title",
				Description: "",
				Type:        models.VoteTypeSingle,
				Options:     []string{"Option 1", "Option 2"},
			},
			expectedError: "description is required",
		},
		{
			name: "Short description",
			request: models.CreateVoteRequest{
				Title:       "Valid Title",
				Description: "Short",
				Type:        models.VoteTypeSingle,
				Options:     []string{"Option 1", "Option 2"},
			},
			expectedError: "description must be between 10 and 1000 characters",
		},
		{
			name: "Invalid vote type",
			request: models.CreateVoteRequest{
				Title:       "Valid Title",
				Description: "Valid description",
				Type:        "invalid",
				Options:     []string{"Option 1", "Option 2"},
			},
			expectedError: "invalid vote type",
		},
		{
			name: "Too few options",
			request: models.CreateVoteRequest{
				Title:       "Valid Title",
				Description: "Valid description",
				Type:        models.VoteTypeSingle,
				Options:     []string{"Option 1"},
			},
			expectedError: "at least 2 options are required",
		},
		{
			name: "Empty option",
			request: models.CreateVoteRequest{
				Title:       "Valid Title",
				Description: "Valid description",
				Type:        models.VoteTypeSingle,
				Options:     []string{"Option 1", ""},
			},
			expectedError: "option 2 cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, _ := json.Marshal(tt.request)
			req := httptest.NewRequest("POST", "/api/v1/votes", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			ctx := context.WithValue(req.Context(), "user_id", "test-user")
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.Contains(t, response, "error")
			errorData := response["error"].(map[string]interface{})
			assert.Contains(t, errorData["message"], tt.expectedError)
		})
	}
}

func TestVoteService_ListVotes(t *testing.T) {
	// Skip if not in integration test mode
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	cfg := &config.Config{
		SupabaseURL:    "http://localhost:54321",
		SupabaseAPIKey: "test-key",
		DatabaseSchema: "testing_db",
	}

	dbClient, err := database.NewClient(cfg)
	require.NoError(t, err)
	defer dbClient.Close()

	serviceContainer := services.NewContainer(dbClient, cfg)

	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		voteHandler := handlers.NewVoteHandler(serviceContainer.VoteService)
		r.Route("/votes", func(r chi.Router) {
			r.Get("/", voteHandler.ListVotes)
			r.With(services.AuthMiddleware(serviceContainer.SessionService)).Post("/", voteHandler.CreateVote)
		})
	})

	// Create multiple votes
	userID := "test-user-123"
	votes := []models.CreateVoteRequest{
		{
			Title:       "Vote 1",
			Description: "First test vote",
			Type:        models.VoteTypeSingle,
			Options:     []string{"A", "B"},
		},
		{
			Title:       "Vote 2",
			Description: "Second test vote",
			Type:        models.VoteTypeMultiple,
			Options:     []string{"X", "Y", "Z"},
		},
	}

	for _, vote := range votes {
		reqBody, _ := json.Marshal(vote)
		req := httptest.NewRequest("POST", "/api/v1/votes", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), "user_id", userID)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	// Test listing votes
	req := httptest.NewRequest("GET", "/api/v1/votes", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "data")
	data := response["data"].([]interface{})
	assert.GreaterOrEqual(t, len(data), 2) // At least the 2 votes we created

	// Test filtering by type
	req = httptest.NewRequest("GET", "/api/v1/votes?type=single", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data = response["data"].([]interface{})
	for _, voteData := range data {
		vote := voteData.(map[string]interface{})
		assert.Equal(t, "single", vote["type"])
	}
}
