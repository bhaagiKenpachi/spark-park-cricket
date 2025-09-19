package unit

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/services"
)

// Note: MockMatchRepository and MockSeriesRepository are defined in match_completion_unit_test.go

func TestMatchService_CreateMatch(t *testing.T) {
	tests := []struct {
		name          string
		request       *models.CreateMatchRequest
		mockSetup     func(*MockMatchRepository, *MockSeriesRepository)
		expectedError string
		expectedMatch *models.Match
	}{
		{
			name: "successful match creation with provided match number",
			request: &models.CreateMatchRequest{
				SeriesID:         "series-1",
				MatchNumber:      intPtr(1),
				Date:             time.Date(2025, 9, 14, 10, 0, 0, 0, time.UTC),
				TeamAPlayerCount: 11,
				TeamBPlayerCount: 11,
				TotalOvers:       20,
				TossWinner:       models.TeamTypeA,
				TossType:         models.TossTypeHeads,
			},
			mockSetup: func(matchRepo *MockMatchRepository, seriesRepo *MockSeriesRepository) {
				seriesRepo.On("GetByID", mock.Anything, "series-1").Return(&models.Series{ID: "series-1"}, nil)
				matchRepo.On("ExistsBySeriesAndMatchNumber", mock.Anything, "series-1", 1).Return(false, nil)
				matchRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Match")).Return(nil).Run(func(args mock.Arguments) {
					match := args.Get(1).(*models.Match)
					match.ID = "match-1"
				})
			},
			expectedMatch: &models.Match{
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
			},
		},
		{
			name: "successful match creation with auto-increment match number",
			request: &models.CreateMatchRequest{
				SeriesID:         "series-1",
				MatchNumber:      nil,
				Date:             time.Date(2025, 9, 14, 10, 0, 0, 0, time.UTC),
				TeamAPlayerCount: 11,
				TeamBPlayerCount: 11,
				TotalOvers:       20,
				TossWinner:       models.TeamTypeB,
				TossType:         models.TossTypeTails,
			},
			mockSetup: func(matchRepo *MockMatchRepository, seriesRepo *MockSeriesRepository) {
				seriesRepo.On("GetByID", mock.Anything, "series-1").Return(&models.Series{ID: "series-1"}, nil)
				matchRepo.On("GetNextMatchNumber", mock.Anything, "series-1").Return(2, nil)
				matchRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Match")).Return(nil).Run(func(args mock.Arguments) {
					match := args.Get(1).(*models.Match)
					match.ID = "match-2"
				})
			},
			expectedMatch: &models.Match{
				ID:               "match-2",
				SeriesID:         "series-1",
				MatchNumber:      2,
				Date:             time.Date(2025, 9, 14, 10, 0, 0, 0, time.UTC),
				Status:           models.MatchStatusLive,
				TeamAPlayerCount: 11,
				TeamBPlayerCount: 11,
				TotalOvers:       20,
				TossWinner:       models.TeamTypeB,
				TossType:         models.TossTypeTails,
				BattingTeam:      models.TeamTypeB,
			},
		},
		{
			name: "series not found",
			request: &models.CreateMatchRequest{
				SeriesID:         "nonexistent-series",
				MatchNumber:      intPtr(1),
				Date:             time.Date(2025, 9, 14, 10, 0, 0, 0, time.UTC),
				TeamAPlayerCount: 11,
				TeamBPlayerCount: 11,
				TotalOvers:       20,
				TossWinner:       models.TeamTypeA,
				TossType:         models.TossTypeHeads,
			},
			mockSetup: func(matchRepo *MockMatchRepository, seriesRepo *MockSeriesRepository) {
				seriesRepo.On("GetByID", mock.Anything, "nonexistent-series").Return(nil, errors.New("series not found"))
			},
			expectedError: "series not found",
		},
		{
			name: "duplicate match number",
			request: &models.CreateMatchRequest{
				SeriesID:         "series-1",
				MatchNumber:      intPtr(1),
				Date:             time.Date(2025, 9, 14, 10, 0, 0, 0, time.UTC),
				TeamAPlayerCount: 11,
				TeamBPlayerCount: 11,
				TotalOvers:       20,
				TossWinner:       models.TeamTypeA,
				TossType:         models.TossTypeHeads,
			},
			mockSetup: func(matchRepo *MockMatchRepository, seriesRepo *MockSeriesRepository) {
				seriesRepo.On("GetByID", mock.Anything, "series-1").Return(&models.Series{ID: "series-1"}, nil)
				matchRepo.On("ExistsBySeriesAndMatchNumber", mock.Anything, "series-1", 1).Return(true, nil)
			},
			expectedError: "match number 1 already exists for series series-1",
		},
		{
			name: "repository create error",
			request: &models.CreateMatchRequest{
				SeriesID:         "series-1",
				MatchNumber:      intPtr(1),
				Date:             time.Date(2025, 9, 14, 10, 0, 0, 0, time.UTC),
				TeamAPlayerCount: 11,
				TeamBPlayerCount: 11,
				TotalOvers:       20,
				TossWinner:       models.TeamTypeA,
				TossType:         models.TossTypeHeads,
			},
			mockSetup: func(matchRepo *MockMatchRepository, seriesRepo *MockSeriesRepository) {
				seriesRepo.On("GetByID", mock.Anything, "series-1").Return(&models.Series{ID: "series-1"}, nil)
				matchRepo.On("ExistsBySeriesAndMatchNumber", mock.Anything, "series-1", 1).Return(false, nil)
				matchRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Match")).Return(errors.New("database error"))
			},
			expectedError: "failed to create match",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMatchRepo := new(MockMatchRepository)
			mockSeriesRepo := new(MockSeriesRepository)
			tt.mockSetup(mockMatchRepo, mockSeriesRepo)

			// Create context with user_id for authentication
			ctx := context.WithValue(context.Background(), "user_id", "test-user-123")

			service := services.NewMatchService(mockMatchRepo, mockSeriesRepo)
			result, err := service.CreateMatch(ctx, tt.request)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedMatch.ID, result.ID)
				assert.Equal(t, tt.expectedMatch.SeriesID, result.SeriesID)
				assert.Equal(t, tt.expectedMatch.MatchNumber, result.MatchNumber)
				assert.Equal(t, tt.expectedMatch.Status, result.Status)
				assert.Equal(t, tt.expectedMatch.TossWinner, result.TossWinner)
				assert.Equal(t, tt.expectedMatch.BattingTeam, result.BattingTeam)
			}

			mockMatchRepo.AssertExpectations(t)
			mockSeriesRepo.AssertExpectations(t)
		})
	}
}

func TestMatchService_GetMatch(t *testing.T) {
	tests := []struct {
		name          string
		matchID       string
		mockSetup     func(*MockMatchRepository)
		expectedError string
		expectedMatch *models.Match
	}{
		{
			name:    "successful match retrieval",
			matchID: "match-1",
			mockSetup: func(mockRepo *MockMatchRepository) {
				mockRepo.On("GetByID", mock.Anything, "match-1").Return(&models.Match{
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
				}, nil)
			},
			expectedMatch: &models.Match{
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
		},
		{
			name:    "empty match ID",
			matchID: "",
			mockSetup: func(mockRepo *MockMatchRepository) {
				// No mock setup needed as validation happens before repository call
			},
			expectedError: "match ID is required",
		},
		{
			name:    "match not found",
			matchID: "nonexistent-match",
			mockSetup: func(mockRepo *MockMatchRepository) {
				mockRepo.On("GetByID", mock.Anything, "nonexistent-match").Return(nil, errors.New("match not found"))
			},
			expectedError: "failed to get match",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMatchRepo := new(MockMatchRepository)
			mockSeriesRepo := new(MockSeriesRepository)
			tt.mockSetup(mockMatchRepo)

			service := services.NewMatchService(mockMatchRepo, mockSeriesRepo)
			result, err := service.GetMatch(context.Background(), tt.matchID)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedMatch.ID, result.ID)
				assert.Equal(t, tt.expectedMatch.SeriesID, result.SeriesID)
				assert.Equal(t, tt.expectedMatch.Status, result.Status)
			}

			mockMatchRepo.AssertExpectations(t)
		})
	}
}

func TestMatchService_ListMatches(t *testing.T) {
	tests := []struct {
		name            string
		filters         *models.MatchFilters
		mockSetup       func(*MockMatchRepository)
		expectedError   string
		expectedMatches []*models.Match
	}{
		{
			name: "successful match listing with default filters",
			filters: &models.MatchFilters{
				Limit:  20,
				Offset: 0,
			},
			mockSetup: func(mockRepo *MockMatchRepository) {
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
					{
						ID:               "match-2",
						SeriesID:         "series-1",
						MatchNumber:      2,
						Status:           models.MatchStatusCompleted,
						TeamAPlayerCount: 11,
						TeamBPlayerCount: 11,
						TotalOvers:       20,
						TossWinner:       models.TeamTypeB,
						TossType:         models.TossTypeTails,
						BattingTeam:      models.TeamTypeB,
					},
				}
				mockRepo.On("GetAll", mock.Anything, mock.AnythingOfType("*models.MatchFilters")).Return(matches, nil)
			},
			expectedMatches: []*models.Match{
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
				{
					ID:               "match-2",
					SeriesID:         "series-1",
					MatchNumber:      2,
					Status:           models.MatchStatusCompleted,
					TeamAPlayerCount: 11,
					TeamBPlayerCount: 11,
					TotalOvers:       20,
					TossWinner:       models.TeamTypeB,
					TossType:         models.TossTypeTails,
					BattingTeam:      models.TeamTypeB,
				},
			},
		},
		{
			name: "filters limit adjustment - too high",
			filters: &models.MatchFilters{
				Limit:  150, // Should be capped to 100
				Offset: 0,
			},
			mockSetup: func(mockRepo *MockMatchRepository) {
				mockRepo.On("GetAll", mock.Anything, mock.MatchedBy(func(filters *models.MatchFilters) bool {
					return filters.Limit == 100 // Should be capped to 100
				})).Return([]*models.Match{}, nil)
			},
			expectedMatches: []*models.Match{},
		},
		{
			name: "filters limit adjustment - zero or negative",
			filters: &models.MatchFilters{
				Limit:  0, // Should be set to 20
				Offset: 0,
			},
			mockSetup: func(mockRepo *MockMatchRepository) {
				mockRepo.On("GetAll", mock.Anything, mock.MatchedBy(func(filters *models.MatchFilters) bool {
					return filters.Limit == 20 // Should be set to default 20
				})).Return([]*models.Match{}, nil)
			},
			expectedMatches: []*models.Match{},
		},
		{
			name: "repository error",
			filters: &models.MatchFilters{
				Limit:  20,
				Offset: 0,
			},
			mockSetup: func(mockRepo *MockMatchRepository) {
				mockRepo.On("GetAll", mock.Anything, mock.AnythingOfType("*models.MatchFilters")).Return(nil, errors.New("database error"))
			},
			expectedError: "failed to list matches",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMatchRepo := new(MockMatchRepository)
			mockSeriesRepo := new(MockSeriesRepository)
			tt.mockSetup(mockMatchRepo)

			service := services.NewMatchService(mockMatchRepo, mockSeriesRepo)
			result, err := service.ListMatches(context.Background(), tt.filters)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result, len(tt.expectedMatches))
				for i, expected := range tt.expectedMatches {
					assert.Equal(t, expected.ID, result[i].ID)
					assert.Equal(t, expected.Status, result[i].Status)
				}
			}

			mockMatchRepo.AssertExpectations(t)
		})
	}
}

func TestMatchService_UpdateMatch(t *testing.T) {
	tests := []struct {
		name          string
		matchID       string
		request       *models.UpdateMatchRequest
		mockSetup     func(*MockMatchRepository)
		expectedError string
		expectedMatch *models.Match
	}{
		{
			name:    "successful match update - status only",
			matchID: "match-1",
			request: &models.UpdateMatchRequest{
				Status: matchStatusPtr(models.MatchStatusCompleted),
			},
			mockSetup: func(mockRepo *MockMatchRepository) {
				existingMatch := &models.Match{
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
					CreatedBy:        "test-user-123",
				}
				mockRepo.On("GetByID", mock.Anything, "match-1").Return(existingMatch, nil)
				mockRepo.On("Update", mock.Anything, "match-1", mock.AnythingOfType("*models.Match")).Return(nil).Run(func(args mock.Arguments) {
					match := args.Get(2).(*models.Match)
					assert.Equal(t, models.MatchStatusCompleted, match.Status)
				})
			},
		},
		{
			name:    "successful match update - batting team only",
			matchID: "match-1",
			request: &models.UpdateMatchRequest{
				BattingTeam: teamTypePtr(models.TeamTypeB),
			},
			mockSetup: func(mockRepo *MockMatchRepository) {
				existingMatch := &models.Match{
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
					CreatedBy:        "test-user-123",
				}
				mockRepo.On("GetByID", mock.Anything, "match-1").Return(existingMatch, nil)
				mockRepo.On("Update", mock.Anything, "match-1", mock.AnythingOfType("*models.Match")).Return(nil).Run(func(args mock.Arguments) {
					match := args.Get(2).(*models.Match)
					assert.Equal(t, models.TeamTypeB, match.BattingTeam)
				})
			},
		},
		{
			name:    "empty match ID",
			matchID: "",
			request: &models.UpdateMatchRequest{
				Status: matchStatusPtr(models.MatchStatusCompleted),
			},
			mockSetup: func(mockRepo *MockMatchRepository) {
				// No mock setup needed as validation happens before repository call
			},
			expectedError: "match ID is required",
		},
		{
			name:    "match not found",
			matchID: "nonexistent-match",
			request: &models.UpdateMatchRequest{
				Status: matchStatusPtr(models.MatchStatusCompleted),
			},
			mockSetup: func(mockRepo *MockMatchRepository) {
				mockRepo.On("GetByID", mock.Anything, "nonexistent-match").Return(nil, errors.New("match not found"))
			},
			expectedError: "failed to get match",
		},
		{
			name:    "repository update error",
			matchID: "match-1",
			request: &models.UpdateMatchRequest{
				Status: matchStatusPtr(models.MatchStatusCompleted),
			},
			mockSetup: func(mockRepo *MockMatchRepository) {
				existingMatch := &models.Match{
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
					CreatedBy:        "test-user-123",
				}
				mockRepo.On("GetByID", mock.Anything, "match-1").Return(existingMatch, nil)
				mockRepo.On("Update", mock.Anything, "match-1", mock.AnythingOfType("*models.Match")).Return(errors.New("database error"))
			},
			expectedError: "failed to update match",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMatchRepo := new(MockMatchRepository)
			mockSeriesRepo := new(MockSeriesRepository)
			tt.mockSetup(mockMatchRepo)

			// Create context with user_id for authentication
			ctx := context.WithValue(context.Background(), "user_id", "test-user-123")

			service := services.NewMatchService(mockMatchRepo, mockSeriesRepo)
			result, err := service.UpdateMatch(ctx, tt.matchID, tt.request)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}

			mockMatchRepo.AssertExpectations(t)
		})
	}
}

func TestMatchService_DeleteMatch(t *testing.T) {
	tests := []struct {
		name          string
		matchID       string
		mockSetup     func(*MockMatchRepository)
		expectedError string
	}{
		{
			name:    "successful match deletion",
			matchID: "match-1",
			mockSetup: func(mockRepo *MockMatchRepository) {
				existingMatch := &models.Match{
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
					CreatedBy:        "test-user-123",
				}
				mockRepo.On("GetByID", mock.Anything, "match-1").Return(existingMatch, nil)
				mockRepo.On("Delete", mock.Anything, "match-1").Return(nil)
			},
		},
		{
			name:    "empty match ID",
			matchID: "",
			mockSetup: func(mockRepo *MockMatchRepository) {
				// No mock setup needed as validation happens before repository call
			},
			expectedError: "match ID is required",
		},
		{
			name:    "match not found",
			matchID: "nonexistent-match",
			mockSetup: func(mockRepo *MockMatchRepository) {
				mockRepo.On("GetByID", mock.Anything, "nonexistent-match").Return(nil, errors.New("match not found"))
			},
			expectedError: "match not found",
		},
		{
			name:    "cannot delete live match",
			matchID: "match-1",
			mockSetup: func(mockRepo *MockMatchRepository) {
				existingMatch := &models.Match{
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
					CreatedBy:        "test-user-123",
				}
				mockRepo.On("GetByID", mock.Anything, "match-1").Return(existingMatch, nil)
				// No Delete mock setup since the method should return early
			},
			expectedError: "cannot delete a live match",
		},
		{
			name:    "repository delete error",
			matchID: "match-1",
			mockSetup: func(mockRepo *MockMatchRepository) {
				existingMatch := &models.Match{
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
					CreatedBy:        "test-user-123",
				}
				mockRepo.On("GetByID", mock.Anything, "match-1").Return(existingMatch, nil)
				mockRepo.On("Delete", mock.Anything, "match-1").Return(errors.New("database error"))
			},
			expectedError: "failed to delete match",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMatchRepo := new(MockMatchRepository)
			mockSeriesRepo := new(MockSeriesRepository)
			tt.mockSetup(mockMatchRepo)

			// Create context with user_id for authentication
			ctx := context.WithValue(context.Background(), "user_id", "test-user-123")

			service := services.NewMatchService(mockMatchRepo, mockSeriesRepo)
			err := service.DeleteMatch(ctx, tt.matchID)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			mockMatchRepo.AssertExpectations(t)
		})
	}
}

func TestMatchService_GetMatchesBySeries(t *testing.T) {
	tests := []struct {
		name            string
		seriesID        string
		mockSetup       func(*MockMatchRepository)
		expectedError   string
		expectedMatches []*models.Match
	}{
		{
			name:     "successful retrieval of matches by series",
			seriesID: "series-1",
			mockSetup: func(mockRepo *MockMatchRepository) {
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
					{
						ID:               "match-2",
						SeriesID:         "series-1",
						MatchNumber:      2,
						Status:           models.MatchStatusCompleted,
						TeamAPlayerCount: 11,
						TeamBPlayerCount: 11,
						TotalOvers:       20,
						TossWinner:       models.TeamTypeB,
						TossType:         models.TossTypeTails,
						BattingTeam:      models.TeamTypeB,
					},
				}
				mockRepo.On("GetBySeriesID", mock.Anything, "series-1").Return(matches, nil)
			},
			expectedMatches: []*models.Match{
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
				{
					ID:               "match-2",
					SeriesID:         "series-1",
					MatchNumber:      2,
					Status:           models.MatchStatusCompleted,
					TeamAPlayerCount: 11,
					TeamBPlayerCount: 11,
					TotalOvers:       20,
					TossWinner:       models.TeamTypeB,
					TossType:         models.TossTypeTails,
					BattingTeam:      models.TeamTypeB,
				},
			},
		},
		{
			name:     "empty series ID",
			seriesID: "",
			mockSetup: func(mockRepo *MockMatchRepository) {
				// No mock setup needed as validation happens before repository call
			},
			expectedError: "series ID is required",
		},
		{
			name:     "repository error",
			seriesID: "series-1",
			mockSetup: func(mockRepo *MockMatchRepository) {
				mockRepo.On("GetBySeriesID", mock.Anything, "series-1").Return(nil, errors.New("database error"))
			},
			expectedError: "failed to get matches by series",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMatchRepo := new(MockMatchRepository)
			mockSeriesRepo := new(MockSeriesRepository)
			tt.mockSetup(mockMatchRepo)

			service := services.NewMatchService(mockMatchRepo, mockSeriesRepo)
			result, err := service.GetMatchesBySeries(context.Background(), tt.seriesID)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result, len(tt.expectedMatches))
				for i, expected := range tt.expectedMatches {
					assert.Equal(t, expected.ID, result[i].ID)
					assert.Equal(t, expected.SeriesID, result[i].SeriesID)
					assert.Equal(t, expected.Status, result[i].Status)
				}
			}

			mockMatchRepo.AssertExpectations(t)
		})
	}
}

// Helper functions
func intPtr(i int) *int {
	return &i
}

func matchStatusPtr(status models.MatchStatus) *models.MatchStatus {
	return &status
}

func teamTypePtr(teamType models.TeamType) *models.TeamType {
	return &teamType
}
