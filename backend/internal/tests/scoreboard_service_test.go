package tests

import (
	"context"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockScoreboardRepository is a mock implementation of ScoreboardRepository
type MockScoreboardRepository struct {
	mock.Mock
}

func (m *MockScoreboardRepository) Create(ctx context.Context, scoreboard *models.LiveScoreboard) error {
	args := m.Called(ctx, scoreboard)
	return args.Error(0)
}

func (m *MockScoreboardRepository) GetByMatchID(ctx context.Context, matchID string) (*models.LiveScoreboard, error) {
	args := m.Called(ctx, matchID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LiveScoreboard), args.Error(1)
}

func (m *MockScoreboardRepository) Update(ctx context.Context, matchID string, scoreboard *models.LiveScoreboard) error {
	args := m.Called(ctx, matchID, scoreboard)
	return args.Error(0)
}

func (m *MockScoreboardRepository) Delete(ctx context.Context, matchID string) error {
	args := m.Called(ctx, matchID)
	return args.Error(0)
}

// MockOverRepository is a mock implementation of OverRepository
type MockOverRepository struct {
	mock.Mock
}

func (m *MockOverRepository) Create(ctx context.Context, over *models.Over) error {
	args := m.Called(ctx, over)
	return args.Error(0)
}

func (m *MockOverRepository) GetByID(ctx context.Context, id string) (*models.Over, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Over), args.Error(1)
}

func (m *MockOverRepository) GetByMatchID(ctx context.Context, matchID string) ([]*models.Over, error) {
	args := m.Called(ctx, matchID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Over), args.Error(1)
}

func (m *MockOverRepository) Update(ctx context.Context, id string, over *models.Over) error {
	args := m.Called(ctx, id, over)
	return args.Error(0)
}

func (m *MockOverRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockOverRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

// MockBallRepository is a mock implementation of BallRepository
type MockBallRepository struct {
	mock.Mock
}

func (m *MockBallRepository) Create(ctx context.Context, ball *models.Ball) error {
	args := m.Called(ctx, ball)
	return args.Error(0)
}

func (m *MockBallRepository) GetByID(ctx context.Context, id string) (*models.Ball, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Ball), args.Error(1)
}

func (m *MockBallRepository) GetByOverID(ctx context.Context, overID string) ([]*models.Ball, error) {
	args := m.Called(ctx, overID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Ball), args.Error(1)
}

func (m *MockBallRepository) Update(ctx context.Context, id string, ball *models.Ball) error {
	args := m.Called(ctx, id, ball)
	return args.Error(0)
}

func (m *MockBallRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockBallRepository) GetByMatchID(ctx context.Context, matchID string) ([]*models.Ball, error) {
	args := m.Called(ctx, matchID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Ball), args.Error(1)
}

func (m *MockBallRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

// MockMatchRepository is a mock implementation of MatchRepository
type MockMatchRepository struct {
	mock.Mock
}

func (m *MockMatchRepository) Create(ctx context.Context, match *models.Match) error {
	args := m.Called(ctx, match)
	return args.Error(0)
}

func (m *MockMatchRepository) GetByID(ctx context.Context, id string) (*models.Match, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Match), args.Error(1)
}

func (m *MockMatchRepository) GetAll(ctx context.Context, filters *models.MatchFilters) ([]*models.Match, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Match), args.Error(1)
}

func (m *MockMatchRepository) GetBySeriesID(ctx context.Context, seriesID string) ([]*models.Match, error) {
	args := m.Called(ctx, seriesID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Match), args.Error(1)
}

func (m *MockMatchRepository) Update(ctx context.Context, id string, match *models.Match) error {
	args := m.Called(ctx, id, match)
	return args.Error(0)
}

func (m *MockMatchRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockMatchRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

// MockTeamRepository is a mock implementation of TeamRepository
type MockTeamRepository struct {
	mock.Mock
}

func (m *MockTeamRepository) Create(ctx context.Context, team *models.Team) error {
	args := m.Called(ctx, team)
	return args.Error(0)
}

func (m *MockTeamRepository) GetByID(ctx context.Context, id string) (*models.Team, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Team), args.Error(1)
}

func (m *MockTeamRepository) GetAll(ctx context.Context, filters *models.TeamFilters) ([]*models.Team, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Team), args.Error(1)
}

func (m *MockTeamRepository) Update(ctx context.Context, id string, team *models.Team) error {
	args := m.Called(ctx, id, team)
	return args.Error(0)
}

func (m *MockTeamRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTeamRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

// MockPlayerRepository is a mock implementation of PlayerRepository
type MockPlayerRepository struct {
	mock.Mock
}

func (m *MockPlayerRepository) Create(ctx context.Context, player *models.Player) error {
	args := m.Called(ctx, player)
	return args.Error(0)
}

func (m *MockPlayerRepository) GetByID(ctx context.Context, id string) (*models.Player, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Player), args.Error(1)
}

func (m *MockPlayerRepository) GetAll(ctx context.Context, filters *models.PlayerFilters) ([]*models.Player, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Player), args.Error(1)
}

func (m *MockPlayerRepository) GetByTeamID(ctx context.Context, teamID string) ([]*models.Player, error) {
	args := m.Called(ctx, teamID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Player), args.Error(1)
}

func (m *MockPlayerRepository) Update(ctx context.Context, id string, player *models.Player) error {
	args := m.Called(ctx, id, player)
	return args.Error(0)
}

func (m *MockPlayerRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPlayerRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func TestScoreboardService_GetScoreboard(t *testing.T) {
	tests := []struct {
		name        string
		matchID     string
		mockSetup   func(*MockMatchRepository, *MockScoreboardRepository)
		expectError bool
		errorMsg    string
	}{
		{
			name:    "successful scoreboard retrieval",
			matchID: "test-match-id",
			mockSetup: func(mockMatchRepo *MockMatchRepository, mockScoreboardRepo *MockScoreboardRepository) {
				match := &models.Match{ID: "test-match-id", Status: models.MatchStatusLive}
				mockMatchRepo.On("GetByID", mock.Anything, "test-match-id").Return(match, nil)

				scoreboard := &models.LiveScoreboard{
					MatchID: "test-match-id",
					Score:   100,
					Wickets: 2,
					Overs:   15.3,
					Balls:   3,
				}
				mockScoreboardRepo.On("GetByMatchID", mock.Anything, "test-match-id").Return(scoreboard, nil)
			},
			expectError: false,
		},
		{
			name:        "empty match ID",
			matchID:     "",
			mockSetup:   func(mockMatchRepo *MockMatchRepository, mockScoreboardRepo *MockScoreboardRepository) {},
			expectError: true,
			errorMsg:    "match ID is required",
		},
		{
			name:    "match not found",
			matchID: "non-existent-id",
			mockSetup: func(mockMatchRepo *MockMatchRepository, mockScoreboardRepo *MockScoreboardRepository) {
				mockMatchRepo.On("GetByID", mock.Anything, "non-existent-id").Return(nil, assert.AnError)
			},
			expectError: true,
			errorMsg:    "match not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockScoreboardRepo := new(MockScoreboardRepository)
			mockOverRepo := new(MockOverRepository)
			mockBallRepo := new(MockBallRepository)
			mockMatchRepo := new(MockMatchRepository)
			mockTeamRepo := new(MockTeamRepository)
			mockPlayerRepo := new(MockPlayerRepository)

			tt.mockSetup(mockMatchRepo, mockScoreboardRepo)

			service := services.NewScoreboardService(
				mockScoreboardRepo, mockOverRepo, mockBallRepo,
				mockMatchRepo, mockTeamRepo, mockPlayerRepo,
			)
			ctx := context.Background()

			result, err := service.GetScoreboard(ctx, tt.matchID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.matchID, result.MatchID)
			}

			mockScoreboardRepo.AssertExpectations(t)
			mockMatchRepo.AssertExpectations(t)
		})
	}
}

func TestScoreboardService_AddBall(t *testing.T) {
	tests := []struct {
		name        string
		matchID     string
		ballEvent   *models.BallEvent
		mockSetup   func(*MockMatchRepository, *MockScoreboardRepository, *MockOverRepository, *MockBallRepository)
		expectError bool
		errorMsg    string
	}{
		{
			name:    "successful ball addition",
			matchID: "test-match-id",
			ballEvent: &models.BallEvent{
				BallType:  models.BallTypeGood,
				Runs:      4,
				IsWicket:  false,
				BatsmanID: "batsman-1",
				BowlerID:  "bowler-1",
			},
			mockSetup: func(mockMatchRepo *MockMatchRepository, mockScoreboardRepo *MockScoreboardRepository, mockOverRepo *MockOverRepository, mockBallRepo *MockBallRepository) {
				match := &models.Match{ID: "test-match-id", Status: models.MatchStatusLive}
				mockMatchRepo.On("GetByID", mock.Anything, "test-match-id").Return(match, nil)

				scoreboard := &models.LiveScoreboard{
					MatchID:       "test-match-id",
					BattingTeamID: "team-1",
					Score:         100,
					Wickets:       2,
					Overs:         15.3,
					Balls:         3,
				}
				mockScoreboardRepo.On("GetByMatchID", mock.Anything, "test-match-id").Return(scoreboard, nil)

				over := &models.Over{
					ID:            "over-1",
					MatchID:       "test-match-id",
					OverNumber:    16,
					BattingTeamID: "team-1",
					TotalRuns:     10,
					TotalBalls:    3,
				}
				mockOverRepo.On("GetByMatchID", mock.Anything, "test-match-id").Return([]*models.Over{over}, nil)

				mockBallRepo.On("GetByOverID", mock.Anything, "over-1").Return([]*models.Ball{}, nil)
				mockBallRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Ball")).Return(nil)
				mockScoreboardRepo.On("Update", mock.Anything, "test-match-id", mock.AnythingOfType("*models.LiveScoreboard")).Return(nil)
			},
			expectError: false,
		},
		{
			name:    "invalid ball event - negative runs",
			matchID: "test-match-id",
			ballEvent: &models.BallEvent{
				BallType:  models.BallTypeGood,
				Runs:      -1,
				IsWicket:  false,
				BatsmanID: "batsman-1",
				BowlerID:  "bowler-1",
			},
			mockSetup: func(mockMatchRepo *MockMatchRepository, mockScoreboardRepo *MockScoreboardRepository, mockOverRepo *MockOverRepository, mockBallRepo *MockBallRepository) {
				match := &models.Match{ID: "test-match-id", Status: models.MatchStatusLive}
				mockMatchRepo.On("GetByID", mock.Anything, "test-match-id").Return(match, nil)

				scoreboard := &models.LiveScoreboard{
					MatchID:       "test-match-id",
					BattingTeamID: "team-1",
					Score:         100,
					Wickets:       2,
					Overs:         15.3,
					Balls:         3,
				}
				mockScoreboardRepo.On("GetByMatchID", mock.Anything, "test-match-id").Return(scoreboard, nil)

				// Mock empty overs list to trigger new over creation
				mockOverRepo.On("GetByMatchID", mock.Anything, "test-match-id").Return([]*models.Over{}, nil)
				mockOverRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Over")).Return(nil)

				// Mock ball repository calls
				mockBallRepo.On("GetByOverID", mock.Anything, mock.AnythingOfType("string")).Return([]*models.Ball{}, nil)
				mockBallRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Ball")).Return(nil)

				// Mock scoreboard update
				mockScoreboardRepo.On("Update", mock.Anything, "test-match-id", mock.AnythingOfType("*models.LiveScoreboard")).Return(nil)
			},
			expectError: true,
			errorMsg:    "invalid ball event",
		},
		{
			name:    "invalid ball event - 5 runs",
			matchID: "test-match-id",
			ballEvent: &models.BallEvent{
				BallType:  models.BallTypeGood,
				Runs:      5,
				IsWicket:  false,
				BatsmanID: "batsman-1",
				BowlerID:  "bowler-1",
			},
			mockSetup: func(mockMatchRepo *MockMatchRepository, mockScoreboardRepo *MockScoreboardRepository, mockOverRepo *MockOverRepository, mockBallRepo *MockBallRepository) {
				match := &models.Match{ID: "test-match-id", Status: models.MatchStatusLive}
				mockMatchRepo.On("GetByID", mock.Anything, "test-match-id").Return(match, nil)

				scoreboard := &models.LiveScoreboard{
					MatchID:       "test-match-id",
					BattingTeamID: "team-1",
					Score:         100,
					Wickets:       2,
					Overs:         15.3,
					Balls:         3,
				}
				mockScoreboardRepo.On("GetByMatchID", mock.Anything, "test-match-id").Return(scoreboard, nil)

				// Mock empty overs list to trigger new over creation
				mockOverRepo.On("GetByMatchID", mock.Anything, "test-match-id").Return([]*models.Over{}, nil)
				mockOverRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Over")).Return(nil)

				// Mock ball repository calls
				mockBallRepo.On("GetByOverID", mock.Anything, mock.AnythingOfType("string")).Return([]*models.Ball{}, nil)
				mockBallRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Ball")).Return(nil)

				// Mock scoreboard update
				mockScoreboardRepo.On("Update", mock.Anything, "test-match-id", mock.AnythingOfType("*models.LiveScoreboard")).Return(nil)
			},
			expectError: true,
			errorMsg:    "invalid ball event",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockScoreboardRepo := new(MockScoreboardRepository)
			mockOverRepo := new(MockOverRepository)
			mockBallRepo := new(MockBallRepository)
			mockMatchRepo := new(MockMatchRepository)
			mockTeamRepo := new(MockTeamRepository)
			mockPlayerRepo := new(MockPlayerRepository)

			tt.mockSetup(mockMatchRepo, mockScoreboardRepo, mockOverRepo, mockBallRepo)

			service := services.NewScoreboardService(
				mockScoreboardRepo, mockOverRepo, mockBallRepo,
				mockMatchRepo, mockTeamRepo, mockPlayerRepo,
			)
			ctx := context.Background()

			result, err := service.AddBall(ctx, tt.matchID, tt.ballEvent)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}

			mockScoreboardRepo.AssertExpectations(t)
			mockOverRepo.AssertExpectations(t)
			mockBallRepo.AssertExpectations(t)
			mockMatchRepo.AssertExpectations(t)
		})
	}
}

func TestScoreboardService_UpdateScore(t *testing.T) {
	tests := []struct {
		name        string
		matchID     string
		request     *models.UpdateScoreRequest
		mockSetup   func(*MockMatchRepository, *MockScoreboardRepository)
		expectError bool
		errorMsg    string
	}{
		{
			name:    "successful score update",
			matchID: "test-match-id",
			request: &models.UpdateScoreRequest{
				Score: 150,
			},
			mockSetup: func(mockMatchRepo *MockMatchRepository, mockScoreboardRepo *MockScoreboardRepository) {
				match := &models.Match{ID: "test-match-id", Status: models.MatchStatusLive}
				mockMatchRepo.On("GetByID", mock.Anything, "test-match-id").Return(match, nil)

				scoreboard := &models.LiveScoreboard{
					MatchID: "test-match-id",
					Score:   100,
					Wickets: 2,
					Overs:   15.3,
					Balls:   3,
				}
				mockScoreboardRepo.On("GetByMatchID", mock.Anything, "test-match-id").Return(scoreboard, nil)
				mockScoreboardRepo.On("Update", mock.Anything, "test-match-id", mock.AnythingOfType("*models.LiveScoreboard")).Return(nil)
			},
			expectError: false,
		},
		{
			name:        "empty match ID",
			matchID:     "",
			request:     &models.UpdateScoreRequest{Score: 150},
			mockSetup:   func(mockMatchRepo *MockMatchRepository, mockScoreboardRepo *MockScoreboardRepository) {},
			expectError: true,
			errorMsg:    "match ID is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockScoreboardRepo := new(MockScoreboardRepository)
			mockOverRepo := new(MockOverRepository)
			mockBallRepo := new(MockBallRepository)
			mockMatchRepo := new(MockMatchRepository)
			mockTeamRepo := new(MockTeamRepository)
			mockPlayerRepo := new(MockPlayerRepository)

			tt.mockSetup(mockMatchRepo, mockScoreboardRepo)

			service := services.NewScoreboardService(
				mockScoreboardRepo, mockOverRepo, mockBallRepo,
				mockMatchRepo, mockTeamRepo, mockPlayerRepo,
			)
			ctx := context.Background()

			result, err := service.UpdateScore(ctx, tt.matchID, tt.request)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.request.Score, result.Score)
			}

			mockScoreboardRepo.AssertExpectations(t)
			mockMatchRepo.AssertExpectations(t)
		})
	}
}
