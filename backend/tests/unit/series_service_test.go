package unit

import (
	"context"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/services"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSeriesRepository is defined in match_completion_unit_test.go

func TestSeriesService_CreateSeries(t *testing.T) {
	tests := []struct {
		name        string
		request     *models.CreateSeriesRequest
		mockSetup   func(*MockSeriesRepository)
		expectError bool
		errorMsg    string
	}{
		{
			name: "successful series creation",
			request: &models.CreateSeriesRequest{
				Name:      "Test Series",
				StartDate: time.Now(),
				EndDate:   time.Now().AddDate(0, 0, 7),
			},
			mockSetup: func(mockRepo *MockSeriesRepository) {
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Series")).Return(nil)
			},
			expectError: false,
		},
		{
			name: "invalid date range",
			request: &models.CreateSeriesRequest{
				Name:      "Test Series",
				StartDate: time.Now().AddDate(0, 0, 7),
				EndDate:   time.Now(),
			},
			mockSetup:   func(mockRepo *MockSeriesRepository) {},
			expectError: true,
			errorMsg:    "end date must be after start date",
		},
		{
			name: "repository error",
			request: &models.CreateSeriesRequest{
				Name:      "Test Series",
				StartDate: time.Now(),
				EndDate:   time.Now().AddDate(0, 0, 7),
			},
			mockSetup: func(mockRepo *MockSeriesRepository) {
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Series")).Return(assert.AnError)
			},
			expectError: true,
			errorMsg:    "failed to create series",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockSeriesRepository)
			tt.mockSetup(mockRepo)

			service := services.NewSeriesService(mockRepo, new(MockMatchRepository))
			ctx := context.Background()

			result, err := service.CreateSeries(ctx, tt.request)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.request.Name, result.Name)
				assert.Equal(t, tt.request.StartDate, result.StartDate)
				assert.Equal(t, tt.request.EndDate, result.EndDate)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestSeriesService_GetSeries(t *testing.T) {
	tests := []struct {
		name        string
		seriesID    string
		mockSetup   func(*MockSeriesRepository)
		expectError bool
		errorMsg    string
	}{
		{
			name:     "successful series retrieval",
			seriesID: "test-series-id",
			mockSetup: func(mockRepo *MockSeriesRepository) {
				series := &models.Series{
					ID:        "test-series-id",
					Name:      "Test Series",
					StartDate: time.Now(),
					EndDate:   time.Now().AddDate(0, 0, 7),
				}
				mockRepo.On("GetByID", mock.Anything, "test-series-id").Return(series, nil)
			},
			expectError: false,
		},
		{
			name:        "empty series ID",
			seriesID:    "",
			mockSetup:   func(mockRepo *MockSeriesRepository) {},
			expectError: true,
			errorMsg:    "series ID is required",
		},
		{
			name:     "series not found",
			seriesID: "non-existent-id",
			mockSetup: func(mockRepo *MockSeriesRepository) {
				mockRepo.On("GetByID", mock.Anything, "non-existent-id").Return(nil, assert.AnError)
			},
			expectError: true,
			errorMsg:    "failed to get series",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockSeriesRepository)
			tt.mockSetup(mockRepo)

			service := services.NewSeriesService(mockRepo, new(MockMatchRepository))
			ctx := context.Background()

			result, err := service.GetSeries(ctx, tt.seriesID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.seriesID, result.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestSeriesService_ListSeries(t *testing.T) {
	tests := []struct {
		name        string
		filters     *models.SeriesFilters
		mockSetup   func(*MockSeriesRepository)
		expectError bool
		expectedLen int
	}{
		{
			name: "successful series listing",
			filters: &models.SeriesFilters{
				Limit:  10,
				Offset: 0,
			},
			mockSetup: func(mockRepo *MockSeriesRepository) {
				series := []*models.Series{
					{ID: "1", Name: "Series 1"},
					{ID: "2", Name: "Series 2"},
				}
				mockRepo.On("GetAll", mock.Anything, mock.AnythingOfType("*models.SeriesFilters")).Return(series, nil)
			},
			expectError: false,
			expectedLen: 2,
		},
		{
			name: "filters limit adjustment - too high",
			filters: &models.SeriesFilters{
				Limit:  150, // Over the max limit of 100
				Offset: 0,
			},
			mockSetup: func(mockRepo *MockSeriesRepository) {
				series := []*models.Series{
					{ID: "1", Name: "Series 1"},
				}
				// Expect the service to adjust the limit to 100
				mockRepo.On("GetAll", mock.Anything, mock.MatchedBy(func(filters *models.SeriesFilters) bool {
					return filters.Limit == 100
				})).Return(series, nil)
			},
			expectError: false,
			expectedLen: 1,
		},
		{
			name:    "default filters",
			filters: &models.SeriesFilters{},
			mockSetup: func(mockRepo *MockSeriesRepository) {
				series := []*models.Series{}
				mockRepo.On("GetAll", mock.Anything, mock.AnythingOfType("*models.SeriesFilters")).Return(series, nil)
			},
			expectError: false,
			expectedLen: 0,
		},
		{
			name: "repository error",
			filters: &models.SeriesFilters{
				Limit:  10,
				Offset: 0,
			},
			mockSetup: func(mockRepo *MockSeriesRepository) {
				mockRepo.On("GetAll", mock.Anything, mock.AnythingOfType("*models.SeriesFilters")).Return(nil, assert.AnError)
			},
			expectError: true,
			expectedLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockSeriesRepository)
			tt.mockSetup(mockRepo)

			service := services.NewSeriesService(mockRepo, new(MockMatchRepository))
			ctx := context.Background()

			result, err := service.ListSeries(ctx, tt.filters)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result, tt.expectedLen)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestSeriesService_UpdateSeries(t *testing.T) {
	tests := []struct {
		name        string
		seriesID    string
		request     *models.UpdateSeriesRequest
		mockSetup   func(*MockSeriesRepository)
		expectError bool
		errorMsg    string
	}{
		{
			name:     "successful series update - name only",
			seriesID: "test-series-id",
			request: &models.UpdateSeriesRequest{
				Name: stringPtr("Updated Series"),
			},
			mockSetup: func(mockRepo *MockSeriesRepository) {
				existingSeries := &models.Series{
					ID:        "test-series-id",
					Name:      "Original Series",
					StartDate: time.Now(),
					EndDate:   time.Now().AddDate(0, 0, 7),
				}
				mockRepo.On("GetByID", mock.Anything, "test-series-id").Return(existingSeries, nil)
				mockRepo.On("Update", mock.Anything, "test-series-id", mock.AnythingOfType("*models.Series")).Return(nil)
			},
			expectError: false,
		},
		{
			name:     "successful series update - dates only",
			seriesID: "test-series-id",
			request: &models.UpdateSeriesRequest{
				StartDate: timePtr(time.Now().AddDate(0, 0, 1)),
				EndDate:   timePtr(time.Now().AddDate(0, 0, 8)),
			},
			mockSetup: func(mockRepo *MockSeriesRepository) {
				existingSeries := &models.Series{
					ID:        "test-series-id",
					Name:      "Original Series",
					StartDate: time.Now(),
					EndDate:   time.Now().AddDate(0, 0, 7),
				}
				mockRepo.On("GetByID", mock.Anything, "test-series-id").Return(existingSeries, nil)
				mockRepo.On("Update", mock.Anything, "test-series-id", mock.AnythingOfType("*models.Series")).Return(nil)
			},
			expectError: false,
		},
		{
			name:     "invalid date range in update",
			seriesID: "test-series-id",
			request: &models.UpdateSeriesRequest{
				StartDate: timePtr(time.Now().AddDate(0, 0, 7)),
				EndDate:   timePtr(time.Now()), // End date before start date
			},
			mockSetup: func(mockRepo *MockSeriesRepository) {
				existingSeries := &models.Series{
					ID:        "test-series-id",
					Name:      "Original Series",
					StartDate: time.Now(),
					EndDate:   time.Now().AddDate(0, 0, 7),
				}
				mockRepo.On("GetByID", mock.Anything, "test-series-id").Return(existingSeries, nil)
			},
			expectError: true,
			errorMsg:    "end date must be after start date",
		},
		{
			name:        "empty series ID",
			seriesID:    "",
			request:     &models.UpdateSeriesRequest{},
			mockSetup:   func(mockRepo *MockSeriesRepository) {},
			expectError: true,
			errorMsg:    "series ID is required",
		},
		{
			name:     "series not found",
			seriesID: "non-existent-id",
			request:  &models.UpdateSeriesRequest{},
			mockSetup: func(mockRepo *MockSeriesRepository) {
				mockRepo.On("GetByID", mock.Anything, "non-existent-id").Return(nil, assert.AnError)
			},
			expectError: true,
			errorMsg:    "failed to get series",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockSeriesRepository)
			tt.mockSetup(mockRepo)

			service := services.NewSeriesService(mockRepo, new(MockMatchRepository))
			ctx := context.Background()

			result, err := service.UpdateSeries(ctx, tt.seriesID, tt.request)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestSeriesService_DeleteSeries(t *testing.T) {
	tests := []struct {
		name        string
		seriesID    string
		mockSetup   func(*MockSeriesRepository)
		expectError bool
		errorMsg    string
	}{
		{
			name:     "successful series deletion",
			seriesID: "test-series-id",
			mockSetup: func(mockRepo *MockSeriesRepository) {
				series := &models.Series{ID: "test-series-id", Name: "Test Series"}
				mockRepo.On("GetByID", mock.Anything, "test-series-id").Return(series, nil)
				mockRepo.On("Delete", mock.Anything, "test-series-id").Return(nil)
			},
			expectError: false,
		},
		{
			name:        "empty series ID",
			seriesID:    "",
			mockSetup:   func(mockRepo *MockSeriesRepository) {},
			expectError: true,
			errorMsg:    "series ID is required",
		},
		{
			name:     "series not found",
			seriesID: "non-existent-id",
			mockSetup: func(mockRepo *MockSeriesRepository) {
				mockRepo.On("GetByID", mock.Anything, "non-existent-id").Return(nil, assert.AnError)
			},
			expectError: true,
			errorMsg:    "series not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockSeriesRepository)
			mockMatchRepo := new(MockMatchRepository)
			tt.mockSetup(mockRepo)

			// Set up match repository expectations for DeleteSeries only for successful cases
			if !tt.expectError {
				mockMatchRepo.On("GetBySeriesID", mock.Anything, mock.Anything).Return([]*models.Match{}, nil)
			}

			service := services.NewSeriesService(mockRepo, mockMatchRepo)
			ctx := context.Background()

			err := service.DeleteSeries(ctx, tt.seriesID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
			mockMatchRepo.AssertExpectations(t)
		})
	}
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}

// Helper function to create time pointer
func timePtr(t time.Time) *time.Time {
	return &t
}
