package unit

import (
	"context"
	"testing"

	"github.com/spark-park-cricket/backend/internal/models"
	"github.com/spark-park-cricket/backend/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockFallOfWicketsRepository is a mock implementation of FallOfWicketsRepository
type MockFallOfWicketsRepository struct {
	mock.Mock
}

func (m *MockFallOfWicketsRepository) Create(ctx context.Context, fallOfWickets *models.FallOfWickets) error {
	args := m.Called(ctx, fallOfWickets)
	return args.Error(0)
}

func (m *MockFallOfWicketsRepository) GetByID(ctx context.Context, id string) (*models.FallOfWickets, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.FallOfWickets), args.Error(1)
}

func (m *MockFallOfWicketsRepository) GetByMatchID(ctx context.Context, matchID string) ([]*models.FallOfWickets, error) {
	args := m.Called(ctx, matchID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.FallOfWickets), args.Error(1)
}

func (m *MockFallOfWicketsRepository) GetByInningsID(ctx context.Context, inningsID string) ([]*models.FallOfWickets, error) {
	args := m.Called(ctx, inningsID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.FallOfWickets), args.Error(1)
}

func (m *MockFallOfWicketsRepository) GetByBallID(ctx context.Context, ballID string) (*models.FallOfWickets, error) {
	args := m.Called(ctx, ballID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.FallOfWickets), args.Error(1)
}

func (m *MockFallOfWicketsRepository) List(ctx context.Context, filters *models.FallOfWicketsFilters) ([]*models.FallOfWickets, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.FallOfWickets), args.Error(1)
}

func (m *MockFallOfWicketsRepository) Update(ctx context.Context, id string, req *models.UpdateFallOfWicketsRequest) (*models.FallOfWickets, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.FallOfWickets), args.Error(1)
}

func (m *MockFallOfWicketsRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockFallOfWicketsRepository) GetSummary(ctx context.Context, matchID string, inningsID *string) (*models.FallOfWicketsSummary, error) {
	args := m.Called(ctx, matchID, inningsID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.FallOfWicketsSummary), args.Error(1)
}

func (m *MockFallOfWicketsRepository) GetWicketNumberForInnings(ctx context.Context, inningsID string) (int, error) {
	args := m.Called(ctx, inningsID)
	return args.Int(0), args.Error(1)
}

func (m *MockFallOfWicketsRepository) GetScoreAtWicket(ctx context.Context, inningsID string, overNumber int, ballNumber int) (int, error) {
	args := m.Called(ctx, inningsID, overNumber, ballNumber)
	return args.Int(0), args.Error(1)
}

// MockScorecardRepository is a mock implementation of ScorecardRepository
type MockScorecardRepository struct {
	mock.Mock
}

func (m *MockScorecardRepository) GetInningsByID(ctx context.Context, id string) (*models.ScorecardInnings, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ScorecardInnings), args.Error(1)
}

func (m *MockScorecardRepository) GetOverByID(ctx context.Context, id string) (*models.ScorecardOver, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ScorecardOver), args.Error(1)
}

func (m *MockScorecardRepository) GetBallByID(ctx context.Context, id string) (*models.ScorecardBall, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ScorecardBall), args.Error(1)
}

// MockMatchRepository is a mock implementation of MatchRepository
type MockMatchRepository struct {
	mock.Mock
}

func (m *MockMatchRepository) GetByID(ctx context.Context, id string) (*models.Match, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Match), args.Error(1)
}

// MockPlayerRepository is a mock implementation of PlayerRepository
type MockPlayerRepository struct {
	mock.Mock
}

func (m *MockPlayerRepository) GetByID(ctx context.Context, id string) (*models.Player, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Player), args.Error(1)
}

func TestFallOfWicketsService_CreateFallOfWickets(t *testing.T) {
	// Setup mocks
	mockFallOfWicketsRepo := &MockFallOfWicketsRepository{}
	mockScorecardRepo := &MockScorecardRepository{}
	mockMatchRepo := &MockMatchRepository{}
	mockPlayerRepo := &MockPlayerRepository{}

	// Create service
	service := services.NewFallOfWicketsService(
		mockFallOfWicketsRepo,
		mockScorecardRepo,
		mockMatchRepo,
		mockPlayerRepo,
	)

	ctx := context.Background()

	// Test data
	matchID := "test-match-id"
	inningsID := "test-innings-id"
	overID := "test-over-id"
	ballID := "test-ball-id"

	// Setup mock expectations
	mockMatchRepo.On("GetByID", ctx, matchID).Return(&models.Match{
		ID: matchID,
	}, nil)

	mockScorecardRepo.On("GetInningsByID", ctx, inningsID).Return(&models.ScorecardInnings{
		ID: inningsID,
	}, nil)

	mockScorecardRepo.On("GetOverByID", ctx, overID).Return(&models.ScorecardOver{
		ID: overID,
	}, nil)

	mockScorecardRepo.On("GetBallByID", ctx, ballID).Return(&models.ScorecardBall{
		ID:       ballID,
		IsWicket: true,
	}, nil)

	mockFallOfWicketsRepo.On("GetWicketNumberForInnings", ctx, inningsID).Return(1, nil)
	mockFallOfWicketsRepo.On("GetScoreAtWicket", ctx, inningsID, 1, 1).Return(50, nil)
	mockFallOfWicketsRepo.On("Create", ctx, mock.AnythingOfType("*models.FallOfWickets")).Return(nil)

	// Test request
	req := &models.CreateFallOfWicketsRequest{
		MatchID:       matchID,
		InningsID:     inningsID,
		OverID:        overID,
		BallID:        ballID,
		WicketNumber:  1,
		BatsmanName:   stringPtr("Test Batsman"),
		Runs:          25,
		BallsFaced:    30,
		WicketType:    "bowled",
		BowlerName:    stringPtr("Test Bowler"),
		OverNumber:    1,
		BallNumber:    1,
		ScoreAtWicket: 50,
	}

	// Execute test
	result, err := service.CreateFallOfWickets(ctx, req)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, matchID, result.MatchID)
	assert.Equal(t, inningsID, result.InningsID)
	assert.Equal(t, overID, result.OverID)
	assert.Equal(t, ballID, result.BallID)
	assert.Equal(t, 1, result.WicketNumber)
	assert.Equal(t, "Test Batsman", *result.BatsmanName)
	assert.Equal(t, 25, result.Runs)
	assert.Equal(t, 30, result.BallsFaced)
	assert.Equal(t, "bowled", result.WicketType)
	assert.Equal(t, "Test Bowler", *result.BowlerName)
	assert.Equal(t, 1, result.OverNumber)
	assert.Equal(t, 1, result.BallNumber)
	assert.Equal(t, 50, result.ScoreAtWicket)

	// Verify all expectations were met
	mockMatchRepo.AssertExpectations(t)
	mockScorecardRepo.AssertExpectations(t)
	mockFallOfWicketsRepo.AssertExpectations(t)
}

func TestFallOfWicketsService_GetFallOfWicketsSummary(t *testing.T) {
	// Setup mocks
	mockFallOfWicketsRepo := &MockFallOfWicketsRepository{}
	mockScorecardRepo := &MockScorecardRepository{}
	mockMatchRepo := &MockMatchRepository{}
	mockPlayerRepo := &MockPlayerRepository{}

	// Create service
	service := services.NewFallOfWicketsService(
		mockFallOfWicketsRepo,
		mockScorecardRepo,
		mockMatchRepo,
		mockPlayerRepo,
	)

	ctx := context.Background()
	matchID := "test-match-id"

	// Setup mock expectations
	mockMatchRepo.On("GetByID", ctx, matchID).Return(&models.Match{
		ID: matchID,
	}, nil)

	mockFallOfWicketsRepo.On("GetSummary", ctx, matchID, (*string)(nil)).Return(&models.FallOfWicketsSummary{
		MatchID:      matchID,
		InningsID:    "test-innings-id",
		TotalWickets: 2,
		Wickets: []models.WicketFall{
			{
				WicketNumber:  1,
				BatsmanName:   stringPtr("Batsman 1"),
				Runs:          25,
				BallsFaced:    30,
				WicketType:    "bowled",
				BowlerName:    stringPtr("Bowler 1"),
				OverNumber:    1,
				BallNumber:    3,
				ScoreAtWicket: 25,
				FallTime:      "10:30",
			},
			{
				WicketNumber:  2,
				BatsmanName:   stringPtr("Batsman 2"),
				Runs:          15,
				BallsFaced:    20,
				WicketType:    "caught",
				BowlerName:    stringPtr("Bowler 2"),
				FielderName:   stringPtr("Fielder 1"),
				OverNumber:    2,
				BallNumber:    1,
				ScoreAtWicket: 40,
				FallTime:      "10:35",
			},
		},
	}, nil)

	// Execute test
	result, err := service.GetFallOfWicketsSummary(ctx, matchID, nil)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, matchID, result.MatchID)
	assert.Equal(t, "test-innings-id", result.InningsID)
	assert.Equal(t, 2, result.TotalWickets)
	assert.Len(t, result.Wickets, 2)

	// Check first wicket
	assert.Equal(t, 1, result.Wickets[0].WicketNumber)
	assert.Equal(t, "Batsman 1", *result.Wickets[0].BatsmanName)
	assert.Equal(t, 25, result.Wickets[0].Runs)
	assert.Equal(t, 30, result.Wickets[0].BallsFaced)
	assert.Equal(t, "bowled", result.Wickets[0].WicketType)
	assert.Equal(t, "Bowler 1", *result.Wickets[0].BowlerName)
	assert.Equal(t, 1, result.Wickets[0].OverNumber)
	assert.Equal(t, 3, result.Wickets[0].BallNumber)
	assert.Equal(t, 25, result.Wickets[0].ScoreAtWicket)
	assert.Equal(t, "10:30", result.Wickets[0].FallTime)

	// Check second wicket
	assert.Equal(t, 2, result.Wickets[1].WicketNumber)
	assert.Equal(t, "Batsman 2", *result.Wickets[1].BatsmanName)
	assert.Equal(t, 15, result.Wickets[1].Runs)
	assert.Equal(t, 20, result.Wickets[1].BallsFaced)
	assert.Equal(t, "caught", result.Wickets[1].WicketType)
	assert.Equal(t, "Bowler 2", *result.Wickets[1].BowlerName)
	assert.Equal(t, "Fielder 1", *result.Wickets[1].FielderName)
	assert.Equal(t, 2, result.Wickets[1].OverNumber)
	assert.Equal(t, 1, result.Wickets[1].BallNumber)
	assert.Equal(t, 40, result.Wickets[1].ScoreAtWicket)
	assert.Equal(t, "10:35", result.Wickets[1].FallTime)

	// Verify all expectations were met
	mockMatchRepo.AssertExpectations(t)
	mockFallOfWicketsRepo.AssertExpectations(t)
}

func TestFallOfWicketsService_CreateFallOfWicketsFromBall(t *testing.T) {
	// Setup mocks
	mockFallOfWicketsRepo := &MockFallOfWicketsRepository{}
	mockScorecardRepo := &MockScorecardRepository{}
	mockMatchRepo := &MockMatchRepository{}
	mockPlayerRepo := &MockPlayerRepository{}

	// Create service
	service := services.NewFallOfWicketsService(
		mockFallOfWicketsRepo,
		mockScorecardRepo,
		mockMatchRepo,
		mockPlayerRepo,
	)

	ctx := context.Background()
	ballID := "test-ball-id"
	inningsID := "test-innings-id"
	overID := "test-over-id"
	matchID := "test-match-id"

	// Setup mock expectations
	mockScorecardRepo.On("GetBallByID", ctx, ballID).Return(&models.ScorecardBall{
		ID:         ballID,
		OverID:     overID,
		IsWicket:   true,
		WicketType: "bowled",
		BallNumber: 3,
	}, nil)

	mockScorecardRepo.On("GetOverByID", ctx, overID).Return(&models.ScorecardOver{
		ID:         overID,
		InningsID:  inningsID,
		OverNumber: 2,
	}, nil)

	mockScorecardRepo.On("GetInningsByID", ctx, inningsID).Return(&models.ScorecardInnings{
		ID:      inningsID,
		MatchID: matchID,
	}, nil)

	mockMatchRepo.On("GetByID", ctx, matchID).Return(&models.Match{
		ID: matchID,
	}, nil)

	mockFallOfWicketsRepo.On("GetWicketNumberForInnings", ctx, inningsID).Return(1, nil)
	mockFallOfWicketsRepo.On("GetScoreAtWicket", ctx, inningsID, 2, 3).Return(45, nil)
	mockFallOfWicketsRepo.On("Create", ctx, mock.AnythingOfType("*models.FallOfWickets")).Return(nil)

	// Execute test
	result, err := service.CreateFallOfWicketsFromBall(ctx, ballID, stringPtr("Test Batsman"), stringPtr("Test Bowler"), nil)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, matchID, result.MatchID)
	assert.Equal(t, inningsID, result.InningsID)
	assert.Equal(t, overID, result.OverID)
	assert.Equal(t, ballID, result.BallID)
	assert.Equal(t, 1, result.WicketNumber)
	assert.Equal(t, "Test Batsman", *result.BatsmanName)
	assert.Equal(t, "bowled", result.WicketType)
	assert.Equal(t, "Test Bowler", *result.BowlerName)
	assert.Equal(t, 2, result.OverNumber)
	assert.Equal(t, 3, result.BallNumber)
	assert.Equal(t, 45, result.ScoreAtWicket)

	// Verify all expectations were met
	mockScorecardRepo.AssertExpectations(t)
	mockMatchRepo.AssertExpectations(t)
	mockFallOfWicketsRepo.AssertExpectations(t)
}

func TestFallOfWicketsService_ValidationErrors(t *testing.T) {
	// Setup mocks
	mockFallOfWicketsRepo := &MockFallOfWicketsRepository{}
	mockScorecardRepo := &MockScorecardRepository{}
	mockMatchRepo := &MockMatchRepository{}
	mockPlayerRepo := &MockPlayerRepository{}

	// Create service
	service := services.NewFallOfWicketsService(
		mockFallOfWicketsRepo,
		mockScorecardRepo,
		mockMatchRepo,
		mockPlayerRepo,
	)

	ctx := context.Background()

	t.Run("Empty match ID", func(t *testing.T) {
		req := &models.CreateFallOfWicketsRequest{
			MatchID: "",
		}

		result, err := service.CreateFallOfWickets(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "match_id is required")
	})

	t.Run("Empty innings ID", func(t *testing.T) {
		req := &models.CreateFallOfWicketsRequest{
			MatchID:   "test-match-id",
			InningsID: "",
		}

		result, err := service.CreateFallOfWickets(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "innings_id is required")
	})

	t.Run("Empty ball ID", func(t *testing.T) {
		req := &models.CreateFallOfWicketsRequest{
			MatchID:   "test-match-id",
			InningsID: "test-innings-id",
			BallID:    "",
		}

		result, err := service.CreateFallOfWickets(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "ball_id is required")
	})

	t.Run("Invalid wicket type", func(t *testing.T) {
		req := &models.CreateFallOfWicketsRequest{
			MatchID:    "test-match-id",
			InningsID:  "test-innings-id",
			BallID:     "test-ball-id",
			WicketType: "invalid_type",
		}

		result, err := service.CreateFallOfWickets(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid wicket_type")
	})

	t.Run("Negative runs", func(t *testing.T) {
		req := &models.CreateFallOfWicketsRequest{
			MatchID:   "test-match-id",
			InningsID: "test-innings-id",
			BallID:    "test-ball-id",
			Runs:      -1,
		}

		result, err := service.CreateFallOfWickets(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "runs cannot be negative")
	})
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}
