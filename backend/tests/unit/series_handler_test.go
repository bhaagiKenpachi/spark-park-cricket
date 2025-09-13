package unit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"spark-park-cricket-backend/internal/handlers"
	"spark-park-cricket-backend/internal/models"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSeriesService is a mock implementation of SeriesService
type MockSeriesService struct {
	mock.Mock
}

func (m *MockSeriesService) CreateSeries(ctx context.Context, req *models.CreateSeriesRequest) (*models.Series, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Series), args.Error(1)
}

func (m *MockSeriesService) GetSeries(ctx context.Context, id string) (*models.Series, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Series), args.Error(1)
}

func (m *MockSeriesService) ListSeries(ctx context.Context, filters *models.SeriesFilters) ([]*models.Series, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Series), args.Error(1)
}

func (m *MockSeriesService) UpdateSeries(ctx context.Context, id string, req *models.UpdateSeriesRequest) (*models.Series, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Series), args.Error(1)
}

func (m *MockSeriesService) DeleteSeries(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestSeriesHandler_ListSeries(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		mockSetup      func(*MockSeriesService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:        "successful list with default pagination",
			queryParams: "",
			mockSetup: func(mockService *MockSeriesService) {
				series := []*models.Series{
					{ID: "1", Name: "Test Series 1"},
					{ID: "2", Name: "Test Series 2"},
				}
				mockService.On("ListSeries", mock.Anything, mock.MatchedBy(func(filters *models.SeriesFilters) bool {
					return filters.Limit == 20 && filters.Offset == 0
				})).Return(series, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:        "successful list with custom pagination",
			queryParams: "?limit=10&offset=5",
			mockSetup: func(mockService *MockSeriesService) {
				series := []*models.Series{
					{ID: "1", Name: "Test Series 1"},
				}
				mockService.On("ListSeries", mock.Anything, mock.MatchedBy(func(filters *models.SeriesFilters) bool {
					return filters.Limit == 10 && filters.Offset == 5
				})).Return(series, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:        "service error",
			queryParams: "",
			mockSetup: func(mockService *MockSeriesService) {
				mockService.On("ListSeries", mock.Anything, mock.AnythingOfType("*models.SeriesFilters")).Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
		},
		{
			name:        "invalid limit parameter",
			queryParams: "?limit=invalid",
			mockSetup: func(mockService *MockSeriesService) {
				series := []*models.Series{}
				mockService.On("ListSeries", mock.Anything, mock.MatchedBy(func(filters *models.SeriesFilters) bool {
					return filters.Limit == 20 && filters.Offset == 0 // Should use defaults
				})).Return(series, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:        "negative limit parameter",
			queryParams: "?limit=-5",
			mockSetup: func(mockService *MockSeriesService) {
				series := []*models.Series{}
				mockService.On("ListSeries", mock.Anything, mock.MatchedBy(func(filters *models.SeriesFilters) bool {
					return filters.Limit == 20 && filters.Offset == 0 // Should use defaults
				})).Return(series, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockSeriesService)
			tt.mockSetup(mockService)

			handler := handlers.NewSeriesHandler(mockService)
			req := httptest.NewRequest("GET", "/api/v1/series"+tt.queryParams, nil)
			w := httptest.NewRecorder()

			handler.ListSeries(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestSeriesHandler_CreateSeries(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(*MockSeriesService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "successful series creation",
			requestBody: models.CreateSeriesRequest{
				Name:      "Test Series",
				StartDate: time.Now(),
				EndDate:   time.Now().AddDate(0, 0, 7),
			},
			mockSetup: func(mockService *MockSeriesService) {
				series := &models.Series{
					ID:        "test-id",
					Name:      "Test Series",
					StartDate: time.Now(),
					EndDate:   time.Now().AddDate(0, 0, 7),
				}
				mockService.On("CreateSeries", mock.Anything, mock.AnythingOfType("*models.CreateSeriesRequest")).Return(series, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedError:  false,
		},
		{
			name:           "invalid JSON body",
			requestBody:    "invalid json",
			mockSetup:      func(mockService *MockSeriesService) {},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedError:  true,
		},
		{
			name: "service error",
			requestBody: models.CreateSeriesRequest{
				Name:      "Test Series",
				StartDate: time.Now(),
				EndDate:   time.Now().AddDate(0, 0, 7),
			},
			mockSetup: func(mockService *MockSeriesService) {
				mockService.On("CreateSeries", mock.Anything, mock.AnythingOfType("*models.CreateSeriesRequest")).Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockSeriesService)
			tt.mockSetup(mockService)

			handler := handlers.NewSeriesHandler(mockService)

			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest("POST", "/api/v1/series", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.CreateSeries(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestSeriesHandler_GetSeries(t *testing.T) {
	tests := []struct {
		name           string
		seriesID       string
		mockSetup      func(*MockSeriesService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:     "successful series retrieval",
			seriesID: "test-series-id",
			mockSetup: func(mockService *MockSeriesService) {
				series := &models.Series{
					ID:   "test-series-id",
					Name: "Test Series",
				}
				mockService.On("GetSeries", mock.Anything, "test-series-id").Return(series, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "empty series ID",
			seriesID:       "",
			mockSetup:      func(mockService *MockSeriesService) {},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedError:  true,
		},
		{
			name:     "series not found",
			seriesID: "non-existent-id",
			mockSetup: func(mockService *MockSeriesService) {
				mockService.On("GetSeries", mock.Anything, "non-existent-id").Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockSeriesService)
			tt.mockSetup(mockService)

			handler := handlers.NewSeriesHandler(mockService)

			req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/series/%s", tt.seriesID), nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.seriesID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			w := httptest.NewRecorder()

			handler.GetSeries(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestSeriesHandler_UpdateSeries(t *testing.T) {
	tests := []struct {
		name           string
		seriesID       string
		requestBody    interface{}
		mockSetup      func(*MockSeriesService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:     "successful series update",
			seriesID: "test-series-id",
			requestBody: models.UpdateSeriesRequest{
				Name: stringPtr("Updated Series"),
			},
			mockSetup: func(mockService *MockSeriesService) {
				series := &models.Series{
					ID:   "test-series-id",
					Name: "Updated Series",
				}
				mockService.On("UpdateSeries", mock.Anything, "test-series-id", mock.AnythingOfType("*models.UpdateSeriesRequest")).Return(series, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "empty series ID",
			seriesID:       "",
			requestBody:    models.UpdateSeriesRequest{},
			mockSetup:      func(mockService *MockSeriesService) {},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedError:  true,
		},
		{
			name:           "invalid JSON body",
			seriesID:       "test-series-id",
			requestBody:    "invalid json",
			mockSetup:      func(mockService *MockSeriesService) {},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedError:  true,
		},
		{
			name:     "service error",
			seriesID: "test-series-id",
			requestBody: models.UpdateSeriesRequest{
				Name: stringPtr("Updated Series"),
			},
			mockSetup: func(mockService *MockSeriesService) {
				mockService.On("UpdateSeries", mock.Anything, "test-series-id", mock.AnythingOfType("*models.UpdateSeriesRequest")).Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockSeriesService)
			tt.mockSetup(mockService)

			handler := handlers.NewSeriesHandler(mockService)

			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest("PUT", fmt.Sprintf("/api/v1/series/%s", tt.seriesID), bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.seriesID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			w := httptest.NewRecorder()

			handler.UpdateSeries(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestSeriesHandler_DeleteSeries(t *testing.T) {
	tests := []struct {
		name           string
		seriesID       string
		mockSetup      func(*MockSeriesService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:     "successful series deletion",
			seriesID: "test-series-id",
			mockSetup: func(mockService *MockSeriesService) {
				mockService.On("DeleteSeries", mock.Anything, "test-series-id").Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "empty series ID",
			seriesID:       "",
			mockSetup:      func(mockService *MockSeriesService) {},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedError:  true,
		},
		{
			name:     "service error",
			seriesID: "test-series-id",
			mockSetup: func(mockService *MockSeriesService) {
				mockService.On("DeleteSeries", mock.Anything, "test-series-id").Return(assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockSeriesService)
			tt.mockSetup(mockService)

			handler := handlers.NewSeriesHandler(mockService)

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/series/%s", tt.seriesID), nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.seriesID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			w := httptest.NewRecorder()

			handler.DeleteSeries(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}
