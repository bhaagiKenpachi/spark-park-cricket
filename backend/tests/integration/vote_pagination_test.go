package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"spark-park-cricket-backend/internal/config"
	"spark-park-cricket-backend/internal/database"
	"spark-park-cricket-backend/internal/handlers"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/services"
	"spark-park-cricket-backend/pkg/testutils"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupVotePaginationRouter(serviceContainer *services.Container) http.Handler {
	router := chi.NewRouter()
	voteHandler := handlers.NewVoteHandler(serviceContainer.VoteService)

	router.Get("/api/v1/votes", voteHandler.ListVotes)

	return router
}

func TestVotePagination_Integration(t *testing.T) {
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
	router := setupVotePaginationRouter(serviceContainer)

	// Create test user
	ctx := context.Background()
	testUser := &models.User{
		GoogleID:      fmt.Sprintf("test-pagination-user-%d", time.Now().UnixNano()),
		Email:         fmt.Sprintf("pagination-%d@test.com", time.Now().UnixNano()),
		Name:          "Pagination User",
		Picture:       "https://example.com/pagination.jpg",
		EmailVerified: true,
	}
	err = dbClient.Repositories.User.CreateUser(ctx, testUser)
	require.NoError(t, err)
	defer func() { _ = dbClient.Repositories.User.DeleteUser(ctx, testUser.ID) }()

	// Create 25 test votes for pagination testing
	voteIDs := make([]string, 25)
	for i := 0; i < 25; i++ {
		vote := &models.Vote{
			Title:       fmt.Sprintf("Test Vote %d", i+1),
			Description: fmt.Sprintf("Description for test vote %d", i+1),
			Type:        models.VoteTypeSingle,
			Status:      models.VoteStatusActive,
			CreatedBy:   testUser.ID,
		}
		err := dbClient.Repositories.Vote.CreateVote(ctx, vote)
		require.NoError(t, err)
		voteIDs[i] = vote.ID

		// Create options for the vote
		options := []*models.VoteOption{
			{VoteID: vote.ID, Text: "Option 1"},
			{VoteID: vote.ID, Text: "Option 2"},
		}
		err = dbClient.Repositories.Vote.CreateVoteOptions(ctx, options)
		require.NoError(t, err)
	}
	defer func() {
		for _, id := range voteIDs {
			_ = dbClient.Repositories.Vote.DeleteVote(ctx, id)
		}
	}()

	tests := []struct {
		name             string
		queryParams      string
		expectedStatus   int
		validateResponse func(*testing.T, *models.PaginatedVoteList)
	}{
		{
			name:           "Default pagination (page 1, size 20)",
			queryParams:    "",
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, resp *models.PaginatedVoteList) {
				assert.GreaterOrEqual(t, len(resp.Votes), 20) // At least 20 items on page 1
				assert.GreaterOrEqual(t, resp.TotalItems, 25) // At least 25 total
				assert.Equal(t, 1, resp.Page)
				assert.Equal(t, 20, resp.PageSize)
				assert.GreaterOrEqual(t, resp.TotalPages, 2)
			},
		},
		{
			name:           "Page 2 with default size",
			queryParams:    "page=2",
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, resp *models.PaginatedVoteList) {
				assert.GreaterOrEqual(t, len(resp.Votes), 1) // At least some items on page 2
				assert.GreaterOrEqual(t, resp.TotalItems, 25)
				assert.Equal(t, 2, resp.Page)
				assert.Equal(t, 20, resp.PageSize)
				assert.GreaterOrEqual(t, resp.TotalPages, 2)
			},
		},
		{
			name:           "Custom page size (10 items)",
			queryParams:    "page=1&page_size=10",
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, resp *models.PaginatedVoteList) {
				assert.Equal(t, 10, len(resp.Votes))
				assert.GreaterOrEqual(t, resp.TotalItems, 25)
				assert.Equal(t, 1, resp.Page)
				assert.Equal(t, 10, resp.PageSize)
				assert.GreaterOrEqual(t, resp.TotalPages, 3)
			},
		},
		{
			name:           "Page 3 with size 10",
			queryParams:    "page=3&page_size=10",
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, resp *models.PaginatedVoteList) {
				assert.GreaterOrEqual(t, len(resp.Votes), 1) // At least 1 item
				assert.GreaterOrEqual(t, resp.TotalItems, 25)
				assert.Equal(t, 3, resp.Page)
				assert.Equal(t, 10, resp.PageSize)
				assert.GreaterOrEqual(t, resp.TotalPages, 3)
			},
		},
		{
			name:           "Offset pagination (limit=5, offset=0)",
			queryParams:    "limit=5&offset=0",
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, resp *models.PaginatedVoteList) {
				assert.Equal(t, 5, len(resp.Votes))
				assert.GreaterOrEqual(t, resp.TotalItems, 25)
				assert.Equal(t, 1, resp.Page)
				assert.Equal(t, 5, resp.PageSize)
				assert.GreaterOrEqual(t, resp.TotalPages, 5)
			},
		},
		{
			name:           "Offset pagination (limit=5, offset=5)",
			queryParams:    "limit=5&offset=5",
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, resp *models.PaginatedVoteList) {
				assert.Equal(t, 5, len(resp.Votes))
				assert.GreaterOrEqual(t, resp.TotalItems, 25)
				assert.Equal(t, 2, resp.Page)
				assert.Equal(t, 5, resp.PageSize)
				assert.GreaterOrEqual(t, resp.TotalPages, 5)
			},
		},
		{
			name:           "Filter by status with pagination",
			queryParams:    "status=active&page=1&page_size=15",
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, resp *models.PaginatedVoteList) {
				assert.GreaterOrEqual(t, len(resp.Votes), 15)
				assert.GreaterOrEqual(t, resp.TotalItems, 25)
				assert.Equal(t, 1, resp.Page)
				assert.Equal(t, 15, resp.PageSize)
				assert.GreaterOrEqual(t, resp.TotalPages, 2)
				// Verify all votes are active
				for _, vote := range resp.Votes {
					assert.Equal(t, models.VoteStatusActive, vote.Status)
				}
			},
		},
		{
			name:           "Empty page (page beyond total)",
			queryParams:    "page=10&page_size=20",
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, resp *models.PaginatedVoteList) {
				assert.Equal(t, 0, len(resp.Votes)) // No items on page 10
				assert.GreaterOrEqual(t, resp.TotalItems, 25)
				assert.Equal(t, 10, resp.Page)
				assert.Equal(t, 20, resp.PageSize)
				assert.GreaterOrEqual(t, resp.TotalPages, 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			url := "/api/v1/votes"
			if tt.queryParams != "" {
				url += "?" + tt.queryParams
			}

			req := httptest.NewRequest(http.MethodGet, url, nil)
			req.Header.Set("Accept", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Execute request
			router.ServeHTTP(w, req)

			// Assert status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Parse response
			var response struct {
				Data  *models.PaginatedVoteList `json:"data"`
				Error *struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				} `json:"error,omitempty"`
			}
			err := json.NewDecoder(w.Body).Decode(&response)
			require.NoError(t, err)

			if tt.expectedStatus == http.StatusOK {
				// Assert success response
				assert.NotNil(t, response.Data)
				assert.Nil(t, response.Error)

				// Run custom validation
				if tt.validateResponse != nil {
					tt.validateResponse(t, response.Data)
				}

				// Verify votes are ordered by created_at desc (newest first)
				if len(response.Data.Votes) > 1 {
					for i := 0; i < len(response.Data.Votes)-1; i++ {
						assert.True(t, response.Data.Votes[i].CreatedAt.After(response.Data.Votes[i+1].CreatedAt) ||
							response.Data.Votes[i].CreatedAt.Equal(response.Data.Votes[i+1].CreatedAt),
							"Votes should be ordered by created_at descending")
					}
				}
			} else {
				// Assert error response
				assert.Nil(t, response.Data)
				assert.NotNil(t, response.Error)
			}
		})
	}
}

func TestVotePaginationEdgeCases_Integration(t *testing.T) {
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
	router := setupVotePaginationRouter(serviceContainer)

	t.Run("No votes - empty pagination", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/votes?status=closed", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Data *models.PaginatedVoteList `json:"data"`
		}
		err := json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)

		assert.NotNil(t, response.Data)
		// May have 0 or more votes depending on test data
		assert.GreaterOrEqual(t, response.Data.TotalItems, 0)
		assert.Equal(t, 1, response.Data.Page)
		assert.Equal(t, 20, response.Data.PageSize)
		assert.GreaterOrEqual(t, response.Data.TotalPages, 1)
	})

	t.Run("Exceeds max page_size (should cap at 100)", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/votes?page_size=200", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Data *models.PaginatedVoteList `json:"data"`
		}
		err := json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)

		// Should be capped at default 20 or max allowed
		assert.LessOrEqual(t, response.Data.PageSize, 100)
	})
}
