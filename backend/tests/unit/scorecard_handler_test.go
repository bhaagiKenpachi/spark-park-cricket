package unit

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"spark-park-cricket-backend/internal/handlers"
	"spark-park-cricket-backend/internal/models"
	"testing"

	"github.com/go-chi/chi/v5"
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

func (m *MockScorecardService) GetScorecard(ctx context.Context, matchID string) (*models.ScorecardResponse, error) {
	args := m.Called(ctx, matchID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ScorecardResponse), args.Error(1)
}

func (m *MockScorecardService) GetCurrentOver(ctx context.Context, matchID string, inningsNumber int) (*models.ScorecardOver, error) {
	args := m.Called(ctx, matchID, inningsNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ScorecardOver), args.Error(1)
}

func (m *MockScorecardService) ShouldCompleteMatch(ctx context.Context, matchID string, secondInnings *models.Innings, match *models.Match) (bool, string) {
	args := m.Called(ctx, matchID, secondInnings, match)
	return args.Bool(0), args.String(1)
}

func (m *MockScorecardService) ValidateInningsOrder(ctx context.Context, matchID string, match *models.Match, inningsNumber int) error {
	args := m.Called(ctx, matchID, match, inningsNumber)
	return args.Error(0)
}

func (m *MockScorecardService) GetNonTossWinner(tossWinner models.TeamType) models.TeamType {
	args := m.Called(tossWinner)
	return args.Get(0).(models.TeamType)
}

func TestScorecardHandler_StartScoring(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		serviceError   error
		expectedStatus int
	}{
		{
			name: "successful start scoring",
			requestBody: models.ScorecardRequest{
				MatchID: "match-1",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid request body",
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "missing match_id",
			requestBody: map[string]interface{}{
				// Missing match_id
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid match_id format",
			requestBody: map[string]interface{}{
				"match_id": "invalid-uuid",
			},
			expectedStatus: http.StatusOK, // The validation passes, service will handle the error
		},
		{
			name: "service error",
			requestBody: models.ScorecardRequest{
				MatchID: "match-1",
			},
			serviceError:   errors.New("match not found"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock service
			mockService := new(MockScorecardService)

			// Setup expectations
			if tt.serviceError != nil {
				mockService.On("StartScoring", mock.Anything, "match-1").Return(tt.serviceError)
			} else if tt.requestBody == "invalid json" {
				// Don't set expectations for invalid JSON
			} else if req, ok := tt.requestBody.(models.ScorecardRequest); ok {
				mockService.On("StartScoring", mock.Anything, req.MatchID).Return(nil)
			} else if reqMap, ok := tt.requestBody.(map[string]interface{}); ok {
				// For validation error cases, we still need to mock the service call
				if matchID, exists := reqMap["match_id"]; exists {
					mockService.On("StartScoring", mock.Anything, matchID).Return(nil)
				}
			}

			// Create handler
			handler := handlers.NewScorecardHandler(mockService)

			// Prepare request
			var body []byte
			if tt.requestBody == "invalid json" {
				body = []byte("invalid json")
			} else {
				var err error
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest("POST", "/api/v1/scorecard/start", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Test
			handler.StartScoring(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Verify expectations
			mockService.AssertExpectations(t)
		})
	}
}

func TestScorecardHandler_AddBall(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		serviceError   error
		expectedStatus int
	}{
		{
			name: "successful add ball",
			requestBody: models.BallEventRequest{
				MatchID:       "match-1",
				InningsNumber: 1,
				BallType:      models.BallTypeGood,
				RunType:       models.RunTypeOne,
				IsWicket:      false,
				Byes:          0,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid request body",
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing required fields",
			requestBody: map[string]interface{}{
				"match_id": "match-1",
				// Missing other required fields
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid ball type",
			requestBody: map[string]interface{}{
				"match_id":       "match-1",
				"innings_number": 1,
				"ball_type":      "invalid",
				"run_type":       "1",
				"is_wicket":      false,
				"byes":           0,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "service error",
			requestBody: models.BallEventRequest{
				MatchID:       "match-1",
				InningsNumber: 1,
				BallType:      models.BallTypeGood,
				RunType:       models.RunTypeOne,
			},
			serviceError:   errors.New("match not found"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock service
			mockService := new(MockScorecardService)

			// Setup expectations
			if tt.serviceError != nil {
				mockService.On("AddBall", mock.Anything, mock.AnythingOfType("*models.BallEventRequest")).Return(tt.serviceError)
			} else if tt.requestBody == "invalid json" {
				// Don't set expectations for invalid JSON
			} else if req, ok := tt.requestBody.(models.BallEventRequest); ok {
				mockService.On("AddBall", mock.Anything, &req).Return(nil)
			}

			// Create handler
			handler := handlers.NewScorecardHandler(mockService)

			// Prepare request
			var body []byte
			if tt.requestBody == "invalid json" {
				body = []byte("invalid json")
			} else {
				var err error
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest("POST", "/api/v1/scorecard/ball", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Test
			handler.AddBall(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Verify expectations
			mockService.AssertExpectations(t)
		})
	}
}

func TestScorecardHandler_GetScorecard(t *testing.T) {
	tests := []struct {
		name           string
		matchID        string
		scorecard      *models.ScorecardResponse
		serviceError   error
		expectedStatus int
	}{
		{
			name:    "successful get scorecard",
			matchID: "match-1",
			scorecard: &models.ScorecardResponse{
				MatchID: "match-1",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "missing match_id parameter",
			matchID:        "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "service error",
			matchID:        "nonexistent-match",
			serviceError:   errors.New("scorecard not found"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock service
			mockService := new(MockScorecardService)

			// Setup expectations
			if tt.matchID != "" {
				mockService.On("GetScorecard", mock.Anything, tt.matchID).Return(tt.scorecard, tt.serviceError)
			}

			// Create handler
			handler := handlers.NewScorecardHandler(mockService)

			if tt.matchID == "" {
				// For missing ID test, call handler directly without router
				req := httptest.NewRequest("GET", "/api/v1/scorecard/", nil)
				w := httptest.NewRecorder()
				handler.GetScorecard(w, req)
				assert.Equal(t, tt.expectedStatus, w.Code)
			} else {
				// Setup router to properly extract URL parameters
				router := chi.NewRouter()
				router.Route("/api/v1", func(r chi.Router) {
					r.Route("/scorecard", func(r chi.Router) {
						r.Get("/{match_id}", handler.GetScorecard)
					})
				})

				url := "/api/v1/scorecard/match-1"
				if tt.matchID == "nonexistent-match" {
					url = "/api/v1/scorecard/nonexistent-match"
				}

				req := httptest.NewRequest("GET", url, nil)
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)
				assert.Equal(t, tt.expectedStatus, w.Code)
			}

			// Verify expectations
			mockService.AssertExpectations(t)
		})
	}
}

func TestScorecardHandler_GetCurrentOver(t *testing.T) {
	tests := []struct {
		name           string
		matchID        string
		inningsQuery   string
		over           *models.ScorecardOver
		serviceError   error
		expectedStatus int
	}{
		{
			name:         "successful get current over",
			matchID:      "match-1",
			inningsQuery: "1",
			over: &models.ScorecardOver{
				OverNumber: 1,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:         "successful get current over with default innings",
			matchID:      "match-1",
			inningsQuery: "",
			over: &models.ScorecardOver{
				OverNumber: 1,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "missing match_id parameter",
			matchID:        "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid innings number",
			matchID:        "match-1",
			inningsQuery:   "3",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "service error",
			matchID:        "nonexistent-match",
			inningsQuery:   "1",
			serviceError:   errors.New("over not found"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock service
			mockService := new(MockScorecardService)

			// Setup expectations
			if tt.matchID != "" {
				inningsNumber := 1 // default
				if tt.inningsQuery == "3" {
					// Don't set expectations for invalid innings
				} else if tt.inningsQuery != "" {
					inningsNumber = 1 // for valid query
				}

				if tt.inningsQuery != "3" {
					mockService.On("GetCurrentOver", mock.Anything, tt.matchID, inningsNumber).Return(tt.over, tt.serviceError)
				}
			}

			// Create handler
			handler := handlers.NewScorecardHandler(mockService)

			if tt.matchID == "" {
				// For missing ID test, call handler directly without router
				req := httptest.NewRequest("GET", "/api/v1/scorecard/current-over", nil)
				w := httptest.NewRecorder()
				handler.GetCurrentOver(w, req)
				assert.Equal(t, tt.expectedStatus, w.Code)
			} else {
				// Setup router to properly extract URL parameters
				router := chi.NewRouter()
				router.Route("/api/v1", func(r chi.Router) {
					r.Route("/scorecard", func(r chi.Router) {
						r.Get("/{match_id}/current-over", handler.GetCurrentOver)
					})
				})

				url := "/api/v1/scorecard/match-1/current-over"
				if tt.matchID == "nonexistent-match" {
					url = "/api/v1/scorecard/nonexistent-match/current-over"
				}

				if tt.inningsQuery != "" {
					url += "?innings=" + tt.inningsQuery
				}

				req := httptest.NewRequest("GET", url, nil)
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)
				assert.Equal(t, tt.expectedStatus, w.Code)
			}

			// Verify expectations
			mockService.AssertExpectations(t)
		})
	}
}

func TestScorecardHandler_GetInnings(t *testing.T) {
	tests := []struct {
		name           string
		matchID        string
		inningsNumber  string
		scorecard      *models.ScorecardResponse
		serviceError   error
		expectedStatus int
	}{
		{
			name:          "successful get innings",
			matchID:       "match-1",
			inningsNumber: "1",
			scorecard: &models.ScorecardResponse{
				MatchID: "match-1",
				Innings: []models.InningsSummary{
					{InningsNumber: 1},
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "missing match_id parameter",
			matchID:        "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "missing innings_number parameter",
			matchID:        "match-1",
			inningsNumber:  "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid innings number",
			matchID:        "match-1",
			inningsNumber:  "3",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "service error",
			matchID:        "nonexistent-match",
			inningsNumber:  "1",
			serviceError:   errors.New("scorecard not found"),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:          "innings not found",
			matchID:       "match-1",
			inningsNumber: "2",
			scorecard: &models.ScorecardResponse{
				MatchID: "match-1",
				Innings: []models.InningsSummary{
					{InningsNumber: 1},
					// Missing innings 2
				},
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock service
			mockService := new(MockScorecardService)

			// Setup expectations
			if tt.matchID != "" && tt.inningsNumber != "" && tt.inningsNumber != "3" {
				mockService.On("GetScorecard", mock.Anything, tt.matchID).Return(tt.scorecard, tt.serviceError)
			}

			// Create handler
			handler := handlers.NewScorecardHandler(mockService)

			if tt.matchID == "" || tt.inningsNumber == "" {
				// For missing parameter tests, call handler directly without router
				req := httptest.NewRequest("GET", "/api/v1/scorecard/", nil)
				w := httptest.NewRecorder()
				handler.GetInnings(w, req)
				assert.Equal(t, tt.expectedStatus, w.Code)
			} else {
				// Setup router to properly extract URL parameters
				router := chi.NewRouter()
				router.Route("/api/v1", func(r chi.Router) {
					r.Route("/scorecard", func(r chi.Router) {
						r.Get("/{match_id}/innings/{innings_number}", handler.GetInnings)
					})
				})

				url := "/api/v1/scorecard/match-1/innings/1"
				if tt.matchID == "nonexistent-match" {
					url = "/api/v1/scorecard/nonexistent-match/innings/1"
				} else if tt.inningsNumber == "2" {
					url = "/api/v1/scorecard/match-1/innings/2"
				} else if tt.inningsNumber == "3" {
					url = "/api/v1/scorecard/match-1/innings/3"
				}

				req := httptest.NewRequest("GET", url, nil)
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)
				assert.Equal(t, tt.expectedStatus, w.Code)
			}

			// Verify expectations
			mockService.AssertExpectations(t)
		})
	}
}

func TestScorecardHandler_GetOver(t *testing.T) {
	tests := []struct {
		name           string
		matchID        string
		inningsNumber  string
		overNumber     string
		scorecard      *models.ScorecardResponse
		serviceError   error
		expectedStatus int
	}{
		{
			name:          "successful get over",
			matchID:       "match-1",
			inningsNumber: "1",
			overNumber:    "1",
			scorecard: &models.ScorecardResponse{
				MatchID: "match-1",
				Innings: []models.InningsSummary{
					{
						InningsNumber: 1,
						Overs: []models.OverSummary{
							{OverNumber: 1},
						},
					},
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "missing match_id parameter",
			matchID:        "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "missing innings_number parameter",
			matchID:        "match-1",
			inningsNumber:  "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "missing over_number parameter",
			matchID:        "match-1",
			inningsNumber:  "1",
			overNumber:     "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid innings number",
			matchID:        "match-1",
			inningsNumber:  "3",
			overNumber:     "1",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid over number",
			matchID:        "match-1",
			inningsNumber:  "1",
			overNumber:     "0",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "service error",
			matchID:        "nonexistent-match",
			inningsNumber:  "1",
			overNumber:     "1",
			serviceError:   errors.New("scorecard not found"),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:          "over not found",
			matchID:       "match-1",
			inningsNumber: "1",
			overNumber:    "2",
			scorecard: &models.ScorecardResponse{
				MatchID: "match-1",
				Innings: []models.InningsSummary{
					{
						InningsNumber: 1,
						Overs: []models.OverSummary{
							{OverNumber: 1},
							// Missing over 2
						},
					},
				},
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock service
			mockService := new(MockScorecardService)

			// Setup expectations
			if tt.matchID != "" && tt.inningsNumber != "" && tt.overNumber != "" && tt.inningsNumber != "3" && tt.overNumber != "0" {
				mockService.On("GetScorecard", mock.Anything, tt.matchID).Return(tt.scorecard, tt.serviceError)
			}

			// Create handler
			handler := handlers.NewScorecardHandler(mockService)

			if tt.matchID == "" || tt.inningsNumber == "" || tt.overNumber == "" {
				// For missing parameter tests, call handler directly without router
				req := httptest.NewRequest("GET", "/api/v1/scorecard/", nil)
				w := httptest.NewRecorder()
				handler.GetOver(w, req)
				assert.Equal(t, tt.expectedStatus, w.Code)
			} else {
				// Setup router to properly extract URL parameters
				router := chi.NewRouter()
				router.Route("/api/v1", func(r chi.Router) {
					r.Route("/scorecard", func(r chi.Router) {
						r.Get("/{match_id}/innings/{innings_number}/over/{over_number}", handler.GetOver)
					})
				})

				url := "/api/v1/scorecard/match-1/innings/1/over/1"
				if tt.matchID == "nonexistent-match" {
					url = "/api/v1/scorecard/nonexistent-match/innings/1/over/1"
				} else if tt.inningsNumber == "3" {
					url = "/api/v1/scorecard/match-1/innings/3/over/1"
				} else if tt.overNumber == "0" {
					url = "/api/v1/scorecard/match-1/innings/1/over/0"
				} else if tt.overNumber == "2" {
					url = "/api/v1/scorecard/match-1/innings/1/over/2"
				}

				req := httptest.NewRequest("GET", url, nil)
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)
				assert.Equal(t, tt.expectedStatus, w.Code)
			}

			// Verify expectations
			mockService.AssertExpectations(t)
		})
	}
}
