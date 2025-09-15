package unit

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"spark-park-cricket-backend/internal/handlers"
	"spark-park-cricket-backend/internal/models"
)

// MockMatchService is a mock implementation of MatchServiceInterface
type MockMatchService struct {
	mock.Mock
}

func (m *MockMatchService) CreateMatch(ctx context.Context, req *models.CreateMatchRequest) (*models.Match, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Match), args.Error(1)
}

func (m *MockMatchService) GetMatch(ctx context.Context, id string) (*models.Match, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Match), args.Error(1)
}

func (m *MockMatchService) ListMatches(ctx context.Context, filters *models.MatchFilters) ([]*models.Match, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Match), args.Error(1)
}

func (m *MockMatchService) UpdateMatch(ctx context.Context, id string, req *models.UpdateMatchRequest) (*models.Match, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Match), args.Error(1)
}

func (m *MockMatchService) DeleteMatch(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockMatchService) GetMatchesBySeries(ctx context.Context, seriesID string) ([]*models.Match, error) {
	args := m.Called(ctx, seriesID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Match), args.Error(1)
}

func TestMatchHandler_ListMatches(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		mockSetup      func(*MockMatchService)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful list with default pagination",
			url:  "/api/v1/matches",
			mockSetup: func(mockService *MockMatchService) {
				matches := []*models.Match{
					{
						ID:               "match-1",
						SeriesID:         "series-1",
						MatchNumber:      1,
						Status:           models.MatchStatusLive,
						TeamAPlayerCount: 11,
						TeamBPlayerCount: 11,
						TotalOvers:       20,
						TossWinner:       models.TeamTypeA,
						TossType:         models.TossTypeHeads,
						BattingTeam:      models.TeamTypeA,
					},
				}
				mockService.On("ListMatches", mock.Anything, mock.MatchedBy(func(filters *models.MatchFilters) bool {
					return filters.Limit == 20 && filters.Offset == 0
				})).Return(matches, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "successful list with custom pagination",
			url:  "/api/v1/matches?limit=10&offset=5",
			mockSetup: func(mockService *MockMatchService) {
				matches := []*models.Match{}
				mockService.On("ListMatches", mock.Anything, mock.MatchedBy(func(filters *models.MatchFilters) bool {
					return filters.Limit == 10 && filters.Offset == 5
				})).Return(matches, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "successful list with series filter",
			url:  "/api/v1/matches?series_id=series-1",
			mockSetup: func(mockService *MockMatchService) {
				matches := []*models.Match{}
				mockService.On("ListMatches", mock.Anything, mock.MatchedBy(func(filters *models.MatchFilters) bool {
					return filters.SeriesID != nil && *filters.SeriesID == "series-1"
				})).Return(matches, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "successful list with status filter",
			url:  "/api/v1/matches?status=live",
			mockSetup: func(mockService *MockMatchService) {
				matches := []*models.Match{}
				mockService.On("ListMatches", mock.Anything, mock.MatchedBy(func(filters *models.MatchFilters) bool {
					return filters.Status != nil && *filters.Status == models.MatchStatusLive
				})).Return(matches, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid limit parameter",
			url:  "/api/v1/matches?limit=invalid",
			mockSetup: func(mockService *MockMatchService) {
				matches := []*models.Match{}
				mockService.On("ListMatches", mock.Anything, mock.MatchedBy(func(filters *models.MatchFilters) bool {
					return filters.Limit == 20 // Should use default
				})).Return(matches, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "negative limit parameter",
			url:  "/api/v1/matches?limit=-1",
			mockSetup: func(mockService *MockMatchService) {
				matches := []*models.Match{}
				mockService.On("ListMatches", mock.Anything, mock.MatchedBy(func(filters *models.MatchFilters) bool {
					return filters.Limit == 20 // Should use default
				})).Return(matches, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "service error",
			url:  "/api/v1/matches",
			mockSetup: func(mockService *MockMatchService) {
				mockService.On("ListMatches", mock.Anything, mock.AnythingOfType("*models.MatchFilters")).Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockMatchService)
			tt.mockSetup(mockService)

			handler := handlers.NewMatchHandler(mockService)
			req := httptest.NewRequest("GET", tt.url, nil)
			w := httptest.NewRecorder()

			handler.ListMatches(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestMatchHandler_CreateMatch(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(*MockMatchService)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful match creation",
			requestBody: models.CreateMatchRequest{
				SeriesID:         "series-1",
				MatchNumber:      intPtr(1),
				Date:             time.Date(2025, 9, 14, 10, 0, 0, 0, time.UTC),
				TeamAPlayerCount: 11,
				TeamBPlayerCount: 11,
				TotalOvers:       20,
				TossWinner:       models.TeamTypeA,
				TossType:         models.TossTypeHeads,
			},
			mockSetup: func(mockService *MockMatchService) {
				createdMatch := &models.Match{
					ID:               "match-1",
					SeriesID:         "series-1",
					MatchNumber:      1,
					Date:             time.Date(2025, 9, 14, 10, 0, 0, 0, time.UTC),
					Status:           models.MatchStatusLive,
					TeamAPlayerCount: 11,
					TeamBPlayerCount: 11,
					TotalOvers:       20,
					TossWinner:       models.TeamTypeA,
					TossType:         models.TossTypeHeads,
					BattingTeam:      models.TeamTypeA,
				}
				mockService.On("CreateMatch", mock.Anything, mock.AnythingOfType("*models.CreateMatchRequest")).Return(createdMatch, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:        "invalid JSON body",
			requestBody: "invalid json",
			mockSetup: func(mockService *MockMatchService) {
				// No mock setup needed as validation happens before service call
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "service error",
			requestBody: models.CreateMatchRequest{
				SeriesID:         "series-1",
				MatchNumber:      intPtr(1),
				Date:             time.Date(2025, 9, 14, 10, 0, 0, 0, time.UTC),
				TeamAPlayerCount: 11,
				TeamBPlayerCount: 11,
				TotalOvers:       20,
				TossWinner:       models.TeamTypeA,
				TossType:         models.TossTypeHeads,
			},
			mockSetup: func(mockService *MockMatchService) {
				mockService.On("CreateMatch", mock.Anything, mock.AnythingOfType("*models.CreateMatchRequest")).Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockMatchService)
			tt.mockSetup(mockService)

			var body []byte
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				var err error
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			handler := handlers.NewMatchHandler(mockService)
			req := httptest.NewRequest("POST", "/api/v1/matches", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.CreateMatch(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestMatchHandler_GetMatch(t *testing.T) {
	tests := []struct {
		name           string
		matchID        string
		mockSetup      func(*MockMatchService)
		expectedStatus int
		expectedError  string
	}{
		{
			name:    "successful match retrieval",
			matchID: "match-1",
			mockSetup: func(mockService *MockMatchService) {
				match := &models.Match{
					ID:               "match-1",
					SeriesID:         "series-1",
					MatchNumber:      1,
					Status:           models.MatchStatusLive,
					TeamAPlayerCount: 11,
					TeamBPlayerCount: 11,
					TotalOvers:       20,
					TossWinner:       models.TeamTypeA,
					TossType:         models.TossTypeHeads,
					BattingTeam:      models.TeamTypeA,
				}
				mockService.On("GetMatch", mock.Anything, "match-1").Return(match, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "empty match ID",
			matchID:        "",
			mockSetup:      func(mockService *MockMatchService) {},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:    "match not found",
			matchID: "nonexistent-match",
			mockSetup: func(mockService *MockMatchService) {
				mockService.On("GetMatch", mock.Anything, "nonexistent-match").Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockMatchService)
			tt.mockSetup(mockService)

			handler := handlers.NewMatchHandler(mockService)

			if tt.matchID == "" {
				// For empty ID test, call handler directly without router
				req := httptest.NewRequest("GET", "/api/v1/matches/", nil)
				w := httptest.NewRecorder()
				handler.GetMatch(w, req)
				assert.Equal(t, tt.expectedStatus, w.Code)
			} else {
				// Setup router to properly extract URL parameters
				router := chi.NewRouter()
				router.Route("/api/v1", func(r chi.Router) {
					r.Route("/matches", func(r chi.Router) {
						r.Get("/{id}", handler.GetMatch)
					})
				})

				url := "/api/v1/matches/match-1"
				if tt.matchID == "nonexistent-match" {
					url = "/api/v1/matches/nonexistent-match"
				}

				req := httptest.NewRequest("GET", url, nil)
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)
				assert.Equal(t, tt.expectedStatus, w.Code)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestMatchHandler_UpdateMatch(t *testing.T) {
	tests := []struct {
		name           string
		matchID        string
		requestBody    interface{}
		mockSetup      func(*MockMatchService)
		expectedStatus int
		expectedError  string
	}{
		{
			name:    "successful match update",
			matchID: "match-1",
			requestBody: models.UpdateMatchRequest{
				Status: matchStatusPtr(models.MatchStatusCompleted),
			},
			mockSetup: func(mockService *MockMatchService) {
				updatedMatch := &models.Match{
					ID:               "match-1",
					SeriesID:         "series-1",
					MatchNumber:      1,
					Status:           models.MatchStatusCompleted,
					TeamAPlayerCount: 11,
					TeamBPlayerCount: 11,
					TotalOvers:       20,
					TossWinner:       models.TeamTypeA,
					TossType:         models.TossTypeHeads,
					BattingTeam:      models.TeamTypeA,
				}
				mockService.On("UpdateMatch", mock.Anything, "match-1", mock.AnythingOfType("*models.UpdateMatchRequest")).Return(updatedMatch, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "empty match ID",
			matchID:        "",
			requestBody:    models.UpdateMatchRequest{},
			mockSetup:      func(mockService *MockMatchService) {},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:        "invalid JSON body",
			matchID:     "match-1",
			requestBody: "invalid json",
			mockSetup: func(mockService *MockMatchService) {
				// No mock setup needed as validation happens before service call
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:    "service error",
			matchID: "match-1",
			requestBody: models.UpdateMatchRequest{
				Status: matchStatusPtr(models.MatchStatusCompleted),
			},
			mockSetup: func(mockService *MockMatchService) {
				mockService.On("UpdateMatch", mock.Anything, "match-1", mock.AnythingOfType("*models.UpdateMatchRequest")).Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockMatchService)
			tt.mockSetup(mockService)

			var body []byte
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				var err error
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			handler := handlers.NewMatchHandler(mockService)

			if tt.matchID == "" {
				// For empty ID test, call handler directly without router
				req := httptest.NewRequest("PUT", "/api/v1/matches/", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				handler.UpdateMatch(w, req)
				assert.Equal(t, tt.expectedStatus, w.Code)
			} else {
				// Setup router to properly extract URL parameters
				router := chi.NewRouter()
				router.Route("/api/v1", func(r chi.Router) {
					r.Route("/matches", func(r chi.Router) {
						r.Put("/{id}", handler.UpdateMatch)
					})
				})

				url := "/api/v1/matches/match-1"

				req := httptest.NewRequest("PUT", url, bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)
				assert.Equal(t, tt.expectedStatus, w.Code)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestMatchHandler_DeleteMatch(t *testing.T) {
	tests := []struct {
		name           string
		matchID        string
		mockSetup      func(*MockMatchService)
		expectedStatus int
		expectedError  string
	}{
		{
			name:    "successful match deletion",
			matchID: "match-1",
			mockSetup: func(mockService *MockMatchService) {
				mockService.On("DeleteMatch", mock.Anything, "match-1").Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "empty match ID",
			matchID:        "",
			mockSetup:      func(mockService *MockMatchService) {},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:    "service error",
			matchID: "match-1",
			mockSetup: func(mockService *MockMatchService) {
				mockService.On("DeleteMatch", mock.Anything, "match-1").Return(assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockMatchService)
			tt.mockSetup(mockService)

			handler := handlers.NewMatchHandler(mockService)

			if tt.matchID == "" {
				// For empty ID test, call handler directly without router
				req := httptest.NewRequest("DELETE", "/api/v1/matches/", nil)
				w := httptest.NewRecorder()
				handler.DeleteMatch(w, req)
				assert.Equal(t, tt.expectedStatus, w.Code)
			} else {
				// Setup router to properly extract URL parameters
				router := chi.NewRouter()
				router.Route("/api/v1", func(r chi.Router) {
					r.Route("/matches", func(r chi.Router) {
						r.Delete("/{id}", handler.DeleteMatch)
					})
				})

				url := "/api/v1/matches/match-1"

				req := httptest.NewRequest("DELETE", url, nil)
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)
				assert.Equal(t, tt.expectedStatus, w.Code)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestMatchHandler_GetMatchesBySeries(t *testing.T) {
	tests := []struct {
		name           string
		seriesID       string
		mockSetup      func(*MockMatchService)
		expectedStatus int
		expectedError  string
	}{
		{
			name:     "successful retrieval of matches by series",
			seriesID: "series-1",
			mockSetup: func(mockService *MockMatchService) {
				matches := []*models.Match{
					{
						ID:               "match-1",
						SeriesID:         "series-1",
						MatchNumber:      1,
						Status:           models.MatchStatusLive,
						TeamAPlayerCount: 11,
						TeamBPlayerCount: 11,
						TotalOvers:       20,
						TossWinner:       models.TeamTypeA,
						TossType:         models.TossTypeHeads,
						BattingTeam:      models.TeamTypeA,
					},
				}
				mockService.On("GetMatchesBySeries", mock.Anything, "series-1").Return(matches, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "empty series ID",
			seriesID:       "",
			mockSetup:      func(mockService *MockMatchService) {},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:     "service error",
			seriesID: "series-1",
			mockSetup: func(mockService *MockMatchService) {
				mockService.On("GetMatchesBySeries", mock.Anything, "series-1").Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockMatchService)
			tt.mockSetup(mockService)

			handler := handlers.NewMatchHandler(mockService)

			if tt.seriesID == "" {
				// For empty series ID test, call handler directly without router
				req := httptest.NewRequest("GET", "/api/v1/matches/series/", nil)
				w := httptest.NewRecorder()
				handler.GetMatchesBySeries(w, req)
				assert.Equal(t, tt.expectedStatus, w.Code)
			} else {
				// Setup router to properly extract URL parameters
				router := chi.NewRouter()
				router.Route("/api/v1", func(r chi.Router) {
					r.Route("/matches", func(r chi.Router) {
						r.Get("/series/{series_id}", handler.GetMatchesBySeries)
					})
				})

				url := "/api/v1/matches/series/series-1"

				req := httptest.NewRequest("GET", url, nil)
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)
				assert.Equal(t, tt.expectedStatus, w.Code)
			}

			mockService.AssertExpectations(t)
		})
	}
}

// Note: Helper functions are defined in match_service_test.go
