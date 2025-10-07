package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"spark-park-cricket-backend/internal/config"
	contextkeys "spark-park-cricket-backend/internal/context"
	"spark-park-cricket-backend/internal/database"
	"spark-park-cricket-backend/internal/handlers"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/services"
	"spark-park-cricket-backend/pkg/testutils"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestVoteRequirements tests all voting requirements
func TestVoteRequirements(t *testing.T) {
	// Load test configuration
	cfg := config.LoadTestConfig()

	// Initialize database client
	dbClient, err := database.NewTestClient(cfg)
	require.NoError(t, err)
	defer dbClient.Close()

	// Clean up any existing test data
	testutils.CleanupAllTestData(t, dbClient)

	// Initialize services
	serviceContainer := services.NewContainer(dbClient, cfg.Config)
	router := setupVoteRouter(serviceContainer)

	// Create two test users
	ctx := context.Background()
	user1 := &models.User{
		GoogleID:      fmt.Sprintf("test-vote-user1-%d", time.Now().UnixNano()),
		Email:         fmt.Sprintf("user1-%d@test.com", time.Now().UnixNano()),
		Name:          "User One",
		Picture:       "https://example.com/user1.jpg",
		EmailVerified: true,
	}
	err = dbClient.Repositories.User.CreateUser(ctx, user1)
	require.NoError(t, err)
	defer func() { _ = dbClient.Repositories.User.DeleteUser(ctx, user1.ID) }()

	user2 := &models.User{
		GoogleID:      fmt.Sprintf("test-vote-user2-%d", time.Now().UnixNano()),
		Email:         fmt.Sprintf("user2-%d@test.com", time.Now().UnixNano()),
		Name:          "User Two",
		Picture:       "https://example.com/user2.jpg",
		EmailVerified: true,
	}
	err = dbClient.Repositories.User.CreateUser(ctx, user2)
	require.NoError(t, err)
	defer func() { _ = dbClient.Repositories.User.DeleteUser(ctx, user2.ID) }()

	// Test 1: Only creator can update vote title/description
	t.Run("Only creator can update vote", func(t *testing.T) {
		// User1 creates a vote
		createReq := models.CreateVoteRequest{
			Title:       "Test Vote",
			Description: "This is a test vote for permissions",
			Type:        models.VoteTypeSingle,
			Options:     []string{"Option 1", "Option 2"},
		}

		voteID := createVote(t, router, user1.ID, &createReq)

		// User1 can update their own vote
		updateReq := models.UpdateVoteRequest{
			Title: stringPtr("Updated Title"),
		}
		updateVote(t, router, user1.ID, voteID, &updateReq, http.StatusOK)

		// User2 cannot update user1's vote
		updateVote(t, router, user2.ID, voteID, &updateReq, http.StatusInternalServerError)
	})

	// Test 2: Voter names are recorded in results
	t.Run("Voter names recorded in results", func(t *testing.T) {
		// Create a vote
		createReq := models.CreateVoteRequest{
			Title:       "Test Vote with Names",
			Description: "Testing voter name recording",
			Type:        models.VoteTypeSingle,
			Options:     []string{"Option A", "Option B"},
		}

		voteID := createVote(t, router, user1.ID, &createReq)

		// Get vote options
		vote := getVoteResults(t, router, voteID)
		require.NotEmpty(t, vote.Options)
		optionID := vote.Options[0].ID

		// User1 casts vote
		castVoteReq := models.VoteRequest{
			SelectedOptions: []string{optionID},
		}
		castVote(t, router, user1.ID, voteID, &castVoteReq, http.StatusOK)

		// Check results include voter name
		results := getVoteResults(t, router, voteID)
		require.NotNil(t, results.ResultsWithNames)
		require.Contains(t, results.ResultsWithNames, optionID)
		require.Len(t, results.ResultsWithNames[optionID], 1)
		assert.Equal(t, user1.ID, results.ResultsWithNames[optionID][0].UserID)
		assert.Equal(t, user1.Name, results.ResultsWithNames[optionID][0].UserName)
	})

	// Test 3: List votes is public (no auth required)
	t.Run("List votes is public", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/votes", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// Test 4: Voter can update their vote in active status
	t.Run("Voter can update vote in active status", func(t *testing.T) {
		// Create a vote
		createReq := models.CreateVoteRequest{
			Title:       "Test Vote Update",
			Description: "Testing vote update functionality",
			Type:        models.VoteTypeSingle,
			Options:     []string{"Option 1", "Option 2"},
		}

		voteID := createVote(t, router, user1.ID, &createReq)

		// Get vote options
		vote := getVoteResults(t, router, voteID)
		option1ID := vote.Options[0].ID
		option2ID := vote.Options[1].ID

		// User2 casts initial vote
		castVoteReq := models.VoteRequest{
			SelectedOptions: []string{option1ID},
		}
		castVote(t, router, user2.ID, voteID, &castVoteReq, http.StatusOK)

		// User2 updates their vote
		updateVoteReq := models.VoteRequest{
			SelectedOptions: []string{option2ID},
		}
		castVote(t, router, user2.ID, voteID, &updateVoteReq, http.StatusOK)

		// Verify the vote was updated
		results := getVoteResults(t, router, voteID)
		assert.Equal(t, 0, results.Results[option1ID])
		assert.Equal(t, 1, results.Results[option2ID])
	})

	// Test 5: Voting only allowed in active status
	t.Run("Voting only allowed in active status", func(t *testing.T) {
		// Create a vote
		createReq := models.CreateVoteRequest{
			Title:       "Test Closed Vote",
			Description: "Testing voting on closed votes",
			Type:        models.VoteTypeSingle,
			Options:     []string{"Option 1", "Option 2"},
		}

		voteID := createVote(t, router, user1.ID, &createReq)

		// Get vote options
		vote := getVoteResults(t, router, voteID)
		optionID := vote.Options[0].ID

		// Close the vote
		closeVote(t, router, user1.ID, voteID, http.StatusOK)

		// Try to vote on closed vote - should fail
		castVoteReq := models.VoteRequest{
			SelectedOptions: []string{optionID},
		}
		castVote(t, router, user2.ID, voteID, &castVoteReq, http.StatusBadRequest)
	})
}

// Helper functions

func setupVoteRouter(serviceContainer *services.Container) http.Handler {
	router := chi.NewRouter()
	voteHandler := handlers.NewVoteHandler(serviceContainer.VoteService)

	router.Get("/api/v1/votes", voteHandler.ListVotes)
	router.Get("/api/v1/votes/{id}/results", voteHandler.GetVoteWithResults)
	router.Post("/api/v1/votes", addUserContext(voteHandler.CreateVote))
	router.Put("/api/v1/votes/{id}", addUserContext(voteHandler.UpdateVote))
	router.Post("/api/v1/votes/{id}/vote", addUserContext(voteHandler.CastVote))
	router.Post("/api/v1/votes/{id}/close", addUserContext(voteHandler.CloseVote))

	return router
}

func addUserContext(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-User-ID")
		if userID != "" {
			ctx := context.WithValue(r.Context(), contextkeys.UserIDKey, userID)
			next(w, r.WithContext(ctx))
		} else {
			next(w, r)
		}
	}
}

func createVote(t *testing.T, router http.Handler, userID string, req *models.CreateVoteRequest) string {
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/api/v1/votes", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-User-ID", userID)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httpReq)

	require.Equal(t, http.StatusOK, w.Code, "Failed to create vote: %s", w.Body.String())

	var response struct {
		Data models.Vote `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	return response.Data.ID
}

func updateVote(t *testing.T, router http.Handler, userID, voteID string, req *models.UpdateVoteRequest, expectedStatus int) {
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("PUT", fmt.Sprintf("/api/v1/votes/%s", voteID), bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-User-ID", userID)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httpReq)

	assert.Equal(t, expectedStatus, w.Code, "Update response: %s", w.Body.String())
}

func castVote(t *testing.T, router http.Handler, userID, voteID string, req *models.VoteRequest, expectedStatus int) {
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", fmt.Sprintf("/api/v1/votes/%s/vote", voteID), bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-User-ID", userID)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httpReq)

	assert.Equal(t, expectedStatus, w.Code, "Cast vote response: %s", w.Body.String())
}

func getVoteResults(t *testing.T, router http.Handler, voteID string) *models.VoteWithResults {
	httpReq := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/votes/%s/results", voteID), nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httpReq)

	require.Equal(t, http.StatusOK, w.Code, "Failed to get vote results: %s", w.Body.String())

	var response struct {
		Data models.VoteWithResults `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	return &response.Data
}

func closeVote(t *testing.T, router http.Handler, userID, voteID string, expectedStatus int) {
	httpReq := httptest.NewRequest("POST", fmt.Sprintf("/api/v1/votes/%s/close", voteID), nil)
	httpReq.Header.Set("X-User-ID", userID)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, httpReq)

	assert.Equal(t, expectedStatus, w.Code, "Close vote response: %s", w.Body.String())
}

func stringPtr(s string) *string {
	return &s
}
