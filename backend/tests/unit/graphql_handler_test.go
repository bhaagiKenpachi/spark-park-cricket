package unit

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"spark-park-cricket-backend/internal/graphql"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/pkg/websocket"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockScorecardService is a mock implementation of ScorecardServiceInterface
type MockScorecardService struct {
	mock.Mock
}

func (m *MockScorecardService) StartScoring(ctx context.Context, matchID string) error {
	args := m.Called(ctx, matchID)
	return args.Error(0)
}

func (m *MockScorecardService) AddBall(ctx context.Context, req *models.BallEventRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockScorecardService) UndoBall(ctx context.Context, matchID string, inningsNumber int) error {
	args := m.Called(ctx, matchID, inningsNumber)
	return args.Error(0)
}

func (m *MockScorecardService) GetScorecard(ctx context.Context, matchID string) (*models.ScorecardResponse, error) {
	args := m.Called(ctx, matchID)
	return args.Get(0).(*models.ScorecardResponse), args.Error(1)
}

func (m *MockScorecardService) GetCurrentOver(ctx context.Context, matchID string, inningsNumber int) (*models.OverSummary, error) {
	args := m.Called(ctx, matchID, inningsNumber)
	return args.Get(0).(*models.OverSummary), args.Error(1)
}

func TestGraphQLHandler_ServeHTTP(t *testing.T) {
	// Create mock services
	mockScorecardService := &MockScorecardService{}
	mockHub := websocket.NewHub()

	// Create GraphQL handler
	handler := graphql.NewGraphQLHandler(mockScorecardService, mockHub)

	// Test data
	matchID := "test-match-id"
	scorecard := &models.ScorecardResponse{
		MatchID:        matchID,
		MatchNumber:    1,
		SeriesName:     "Test Series",
		TeamA:          "Team A",
		TeamB:          "Team B",
		TotalOvers:     20,
		TossWinner:     models.TeamTypeA,
		TossType:       models.TossTypeHeads,
		CurrentInnings: 1,
		MatchStatus:    "live",
		Innings: []models.InningsSummary{
			{
				InningsNumber: 1,
				BattingTeam:   models.TeamTypeA,
				TotalRuns:     150,
				TotalWickets:  3,
				TotalOvers:    25.3,
				TotalBalls:    153,
				Status:        "in_progress",
			},
		},
	}

	currentOver := &models.OverSummary{
		OverNumber:   26,
		TotalRuns:    8,
		TotalBalls:   3,
		TotalWickets: 0,
		Status:       "in_progress",
		Balls: []models.BallSummary{
			{
				BallNumber: 1,
				BallType:   models.BallTypeGood,
				RunType:    models.RunTypeFour,
				Runs:       4,
				Byes:       0,
				IsWicket:   false,
			},
		},
	}

	// Set up mock expectations
	mockScorecardService.On("GetScorecard", mock.Anything, matchID).Return(scorecard, nil)
	mockScorecardService.On("GetCurrentOver", mock.Anything, matchID, 1).Return(currentOver, nil)

	// Test query
	query := `
		query GetLiveScorecard($matchId: String!) {
			liveScorecard(match_id: $matchId) {
				match_id
				current_score {
					runs
					wickets
					overs
					run_rate
				}
			}
		}
	`

	variables := map[string]interface{}{
		"matchId": matchID,
	}

	requestBody := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}

	jsonBody, _ := json.Marshal(requestBody)

	// Create HTTP request
	req := httptest.NewRequest("POST", "/graphql", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Call handler
	handler.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	// Parse response
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify response structure
	assert.Contains(t, response, "data")
	data := response["data"].(map[string]interface{})
	assert.Contains(t, data, "liveScorecard")

	liveScorecard := data["liveScorecard"].(map[string]interface{})
	assert.Equal(t, matchID, liveScorecard["match_id"])

	currentScore := liveScorecard["current_score"].(map[string]interface{})
	assert.Equal(t, float64(150), currentScore["runs"])
	assert.Equal(t, float64(3), currentScore["wickets"])
	assert.Equal(t, 25.3, currentScore["overs"])
	assert.Equal(t, 5.93, currentScore["run_rate"])

	// Verify mock calls
	mockScorecardService.AssertExpectations(t)
}

func TestGraphQLHandler_InvalidQuery(t *testing.T) {
	// Create mock services
	mockScorecardService := &MockScorecardService{}
	mockHub := websocket.NewHub()

	// Create GraphQL handler
	handler := graphql.NewGraphQLHandler(mockScorecardService, mockHub)

	// Invalid query
	query := `
		query {
			invalidField
		}
	`

	requestBody := map[string]interface{}{
		"query": query,
	}

	jsonBody, _ := json.Marshal(requestBody)

	// Create HTTP request
	req := httptest.NewRequest("POST", "/graphql", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Call handler
	handler.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse response
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Should have errors
	assert.Contains(t, response, "errors")
	errors := response["errors"].([]interface{})
	assert.Greater(t, len(errors), 0)
}

func TestGraphQLHandler_OPTIONS(t *testing.T) {
	// Create mock services
	mockScorecardService := &MockScorecardService{}
	mockHub := websocket.NewHub()

	// Create GraphQL handler
	handler := graphql.NewGraphQLHandler(mockScorecardService, mockHub)

	// Create OPTIONS request
	req := httptest.NewRequest("OPTIONS", "/graphql", nil)

	// Create response recorder
	w := httptest.NewRecorder()

	// Call handler
	handler.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
}

func TestGraphQLHandler_InvalidMethod(t *testing.T) {
	// Create mock services
	mockScorecardService := &MockScorecardService{}
	mockHub := websocket.NewHub()

	// Create GraphQL handler
	handler := graphql.NewGraphQLHandler(mockScorecardService, mockHub)

	// Create GET request (should be rejected)
	req := httptest.NewRequest("GET", "/graphql", nil)

	// Create response recorder
	w := httptest.NewRecorder()

	// Call handler
	handler.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}
