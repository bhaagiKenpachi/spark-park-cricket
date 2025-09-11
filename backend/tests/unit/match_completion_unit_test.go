package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/services"
)

// MockScorecardRepository for testing
type MockScorecardRepository struct {
	mock.Mock
}

func (m *MockScorecardRepository) CreateInnings(ctx context.Context, innings *models.Innings) error {
	args := m.Called(ctx, innings)
	return args.Error(0)
}

func (m *MockScorecardRepository) GetInningsByMatchAndNumber(ctx context.Context, matchID string, inningsNumber int) (*models.Innings, error) {
	args := m.Called(ctx, matchID, inningsNumber)
	return args.Get(0).(*models.Innings), args.Error(1)
}

func (m *MockScorecardRepository) UpdateInnings(ctx context.Context, innings *models.Innings) error {
	args := m.Called(ctx, innings)
	return args.Error(0)
}

func (m *MockScorecardRepository) CreateOver(ctx context.Context, over *models.ScorecardOver) error {
	args := m.Called(ctx, over)
	return args.Error(0)
}

func (m *MockScorecardRepository) GetCurrentOver(ctx context.Context, inningsID string) (*models.ScorecardOver, error) {
	args := m.Called(ctx, inningsID)
	return args.Get(0).(*models.ScorecardOver), args.Error(1)
}

func (m *MockScorecardRepository) UpdateOver(ctx context.Context, over *models.ScorecardOver) error {
	args := m.Called(ctx, over)
	return args.Error(0)
}

func (m *MockScorecardRepository) CreateBall(ctx context.Context, ball *models.ScorecardBall) error {
	args := m.Called(ctx, ball)
	return args.Error(0)
}

func (m *MockScorecardRepository) GetBallsByOver(ctx context.Context, overID string) ([]*models.ScorecardBall, error) {
	args := m.Called(ctx, overID)
	return args.Get(0).([]*models.ScorecardBall), args.Error(1)
}

func (m *MockScorecardRepository) GetOversByInnings(ctx context.Context, inningsID string) ([]*models.ScorecardOver, error) {
	args := m.Called(ctx, inningsID)
	return args.Get(0).([]*models.ScorecardOver), args.Error(1)
}

func (m *MockScorecardRepository) GetScorecard(ctx context.Context, matchID string) (*models.ScorecardResponse, error) {
	args := m.Called(ctx, matchID)
	return args.Get(0).(*models.ScorecardResponse), args.Error(1)
}

func (m *MockScorecardRepository) CompleteInnings(ctx context.Context, inningsID string) error {
	args := m.Called(ctx, inningsID)
	return args.Error(0)
}

func (m *MockScorecardRepository) GetOverByInningsAndNumber(ctx context.Context, inningsID string, overNumber int) (*models.ScorecardOver, error) {
	args := m.Called(ctx, inningsID, overNumber)
	return args.Get(0).(*models.ScorecardOver), args.Error(1)
}

func (m *MockScorecardRepository) CompleteOver(ctx context.Context, overID string) error {
	args := m.Called(ctx, overID)
	return args.Error(0)
}

func (m *MockScorecardRepository) GetLastBall(ctx context.Context, overID string) (*models.ScorecardBall, error) {
	args := m.Called(ctx, overID)
	return args.Get(0).(*models.ScorecardBall), args.Error(1)
}

func (m *MockScorecardRepository) GetInningsByMatchID(ctx context.Context, matchID string) ([]*models.Innings, error) {
	args := m.Called(ctx, matchID)
	return args.Get(0).([]*models.Innings), args.Error(1)
}

func (m *MockScorecardRepository) StartScoring(ctx context.Context, matchID string) error {
	args := m.Called(ctx, matchID)
	return args.Error(0)
}

// MockMatchRepository for testing
type MockMatchRepository struct {
	mock.Mock
}

func (m *MockMatchRepository) Create(ctx context.Context, match *models.Match) error {
	args := m.Called(ctx, match)
	return args.Error(0)
}

func (m *MockMatchRepository) GetByID(ctx context.Context, id string) (*models.Match, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Match), args.Error(1)
}

func (m *MockMatchRepository) Update(ctx context.Context, id string, match *models.Match) error {
	args := m.Called(ctx, id, match)
	return args.Error(0)
}

func (m *MockMatchRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockMatchRepository) GetAll(ctx context.Context, filters *models.MatchFilters) ([]*models.Match, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).([]*models.Match), args.Error(1)
}

func (m *MockMatchRepository) GetBySeriesID(ctx context.Context, seriesID string) ([]*models.Match, error) {
	args := m.Called(ctx, seriesID)
	return args.Get(0).([]*models.Match), args.Error(1)
}

func (m *MockMatchRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockMatchRepository) GetNextMatchNumber(ctx context.Context, seriesID string) (int, error) {
	args := m.Called(ctx, seriesID)
	return args.Get(0).(int), args.Error(1)
}

func (m *MockMatchRepository) ExistsBySeriesAndMatchNumber(ctx context.Context, seriesID string, matchNumber int) (bool, error) {
	args := m.Called(ctx, seriesID, matchNumber)
	return args.Get(0).(bool), args.Error(1)
}

func TestShouldCompleteMatch_TargetReached(t *testing.T) {
	// Setup
	mockScorecardRepo := &MockScorecardRepository{}
	mockMatchRepo := &MockMatchRepository{}
	service := services.NewScorecardService(mockScorecardRepo, mockMatchRepo)

	ctx := context.Background()
	matchID := "test-match-id"

	// First innings with 10 runs
	firstInnings := &models.Innings{
		MatchID:       matchID,
		InningsNumber: 1,
		TotalRuns:     10,
		Status:        string(models.InningsStatusCompleted),
	}

	// Second innings with 11 runs (target reached: 11 >= 10+1)
	secondInnings := &models.Innings{
		MatchID:       matchID,
		InningsNumber: 2,
		TotalRuns:     11,
		TotalWickets:  0,
		TotalOvers:    1.0,
		Status:        string(models.InningsStatusInProgress),
	}

	// Match with 3 players per team
	match := &models.Match{
		ID:               matchID,
		TeamAPlayerCount: 3,
		TeamBPlayerCount: 3,
		TotalOvers:       2,
		Status:           models.MatchStatusLive,
	}

	// Mock expectations
	mockScorecardRepo.On("GetInningsByMatchAndNumber", ctx, matchID, 1).Return(firstInnings, nil)

	// Test
	shouldComplete, reason := service.ShouldCompleteMatch(ctx, matchID, secondInnings, match)

	// Assertions
	assert.True(t, shouldComplete)
	assert.Contains(t, reason, "target reached")
	assert.Contains(t, reason, "11/11")

	mockScorecardRepo.AssertExpectations(t)
}

func TestShouldCompleteMatch_AllWicketsLost(t *testing.T) {
	// Setup
	mockScorecardRepo := &MockScorecardRepository{}
	mockMatchRepo := &MockMatchRepository{}
	service := services.NewScorecardService(mockScorecardRepo, mockMatchRepo)

	ctx := context.Background()
	matchID := "test-match-id"

	// First innings with 10 runs
	firstInnings := &models.Innings{
		MatchID:       matchID,
		InningsNumber: 1,
		TotalRuns:     10,
		Status:        string(models.InningsStatusCompleted),
	}

	// Second innings with 2 wickets lost (n-1 = 3-1 = 2)
	secondInnings := &models.Innings{
		MatchID:       matchID,
		InningsNumber: 2,
		TotalRuns:     5, // Below target
		TotalWickets:  2, // All wickets lost (n-1)
		TotalOvers:    1.0,
		Status:        string(models.InningsStatusInProgress),
	}

	// Match with 3 players per team
	match := &models.Match{
		ID:               matchID,
		TeamAPlayerCount: 3,
		TeamBPlayerCount: 3,
		TotalOvers:       2,
		Status:           models.MatchStatusLive,
	}

	// Mock expectations
	mockScorecardRepo.On("GetInningsByMatchAndNumber", ctx, matchID, 1).Return(firstInnings, nil)

	// Test
	shouldComplete, reason := service.ShouldCompleteMatch(ctx, matchID, secondInnings, match)

	// Assertions
	assert.True(t, shouldComplete)
	assert.Contains(t, reason, "all wickets lost")
	assert.Contains(t, reason, "2/2")

	mockScorecardRepo.AssertExpectations(t)
}

func TestShouldCompleteMatch_AllOversCompleted(t *testing.T) {
	// Setup
	mockScorecardRepo := &MockScorecardRepository{}
	mockMatchRepo := &MockMatchRepository{}
	service := services.NewScorecardService(mockScorecardRepo, mockMatchRepo)

	ctx := context.Background()
	matchID := "test-match-id"

	// First innings with 10 runs
	firstInnings := &models.Innings{
		MatchID:       matchID,
		InningsNumber: 1,
		TotalRuns:     10,
		Status:        string(models.InningsStatusCompleted),
	}

	// Second innings with all overs completed
	secondInnings := &models.Innings{
		MatchID:       matchID,
		InningsNumber: 2,
		TotalRuns:     8,   // Below target
		TotalWickets:  1,   // Not all wickets lost
		TotalOvers:    2.0, // All overs completed
		Status:        string(models.InningsStatusInProgress),
	}

	// Match with 2 overs
	match := &models.Match{
		ID:               matchID,
		TeamAPlayerCount: 3,
		TeamBPlayerCount: 3,
		TotalOvers:       2,
		Status:           models.MatchStatusLive,
	}

	// Mock expectations
	mockScorecardRepo.On("GetInningsByMatchAndNumber", ctx, matchID, 1).Return(firstInnings, nil)

	// Test
	shouldComplete, reason := service.ShouldCompleteMatch(ctx, matchID, secondInnings, match)

	// Assertions
	assert.True(t, shouldComplete)
	assert.Contains(t, reason, "all overs completed")
	assert.Contains(t, reason, "2.0/2")

	mockScorecardRepo.AssertExpectations(t)
}

func TestShouldCompleteMatch_MatchContinues(t *testing.T) {
	// Setup
	mockScorecardRepo := &MockScorecardRepository{}
	mockMatchRepo := &MockMatchRepository{}
	service := services.NewScorecardService(mockScorecardRepo, mockMatchRepo)

	ctx := context.Background()
	matchID := "test-match-id"

	// First innings with 10 runs
	firstInnings := &models.Innings{
		MatchID:       matchID,
		InningsNumber: 1,
		TotalRuns:     10,
		Status:        string(models.InningsStatusCompleted),
	}

	// Second innings still in progress
	secondInnings := &models.Innings{
		MatchID:       matchID,
		InningsNumber: 2,
		TotalRuns:     5,   // Below target (11)
		TotalWickets:  1,   // Not all wickets lost (need 2)
		TotalOvers:    1.0, // Not all overs completed (need 2)
		Status:        string(models.InningsStatusInProgress),
	}

	// Match with 2 overs
	match := &models.Match{
		ID:               matchID,
		TeamAPlayerCount: 3,
		TeamBPlayerCount: 3,
		TotalOvers:       2,
		Status:           models.MatchStatusLive,
	}

	// Mock expectations
	mockScorecardRepo.On("GetInningsByMatchAndNumber", ctx, matchID, 1).Return(firstInnings, nil)

	// Test
	shouldComplete, reason := service.ShouldCompleteMatch(ctx, matchID, secondInnings, match)

	// Assertions
	assert.False(t, shouldComplete)
	assert.Equal(t, "match continues", reason)

	mockScorecardRepo.AssertExpectations(t)
}

func TestShouldCompleteMatch_ErrorGettingFirstInnings(t *testing.T) {
	// Setup
	mockScorecardRepo := &MockScorecardRepository{}
	mockMatchRepo := &MockMatchRepository{}
	service := services.NewScorecardService(mockScorecardRepo, mockMatchRepo)

	ctx := context.Background()
	matchID := "test-match-id"

	// Second innings
	secondInnings := &models.Innings{
		MatchID:       matchID,
		InningsNumber: 2,
		TotalRuns:     11,
		Status:        string(models.InningsStatusInProgress),
	}

	// Match
	match := &models.Match{
		ID:               matchID,
		TeamAPlayerCount: 3,
		TeamBPlayerCount: 3,
		TotalOvers:       2,
		Status:           models.MatchStatusLive,
	}

	// Mock expectations - return error
	mockScorecardRepo.On("GetInningsByMatchAndNumber", ctx, matchID, 1).Return((*models.Innings)(nil), assert.AnError)

	// Test
	shouldComplete, reason := service.ShouldCompleteMatch(ctx, matchID, secondInnings, match)

	// Assertions
	assert.False(t, shouldComplete)
	assert.Equal(t, "error getting first innings", reason)

	mockScorecardRepo.AssertExpectations(t)
}

func TestShouldCompleteMatch_EdgeCase_ExactTarget(t *testing.T) {
	// Setup
	mockScorecardRepo := &MockScorecardRepository{}
	mockMatchRepo := &MockMatchRepository{}
	service := services.NewScorecardService(mockScorecardRepo, mockMatchRepo)

	ctx := context.Background()
	matchID := "test-match-id"

	// First innings with 10 runs
	firstInnings := &models.Innings{
		MatchID:       matchID,
		InningsNumber: 1,
		TotalRuns:     10,
		Status:        string(models.InningsStatusCompleted),
	}

	// Second innings with exactly the target (11 runs)
	secondInnings := &models.Innings{
		MatchID:       matchID,
		InningsNumber: 2,
		TotalRuns:     11, // Exactly target (10+1)
		TotalWickets:  0,
		TotalOvers:    1.0,
		Status:        string(models.InningsStatusInProgress),
	}

	// Match with 3 players per team
	match := &models.Match{
		ID:               matchID,
		TeamAPlayerCount: 3,
		TeamBPlayerCount: 3,
		TotalOvers:       2,
		Status:           models.MatchStatusLive,
	}

	// Mock expectations
	mockScorecardRepo.On("GetInningsByMatchAndNumber", ctx, matchID, 1).Return(firstInnings, nil)

	// Test
	shouldComplete, reason := service.ShouldCompleteMatch(ctx, matchID, secondInnings, match)

	// Assertions
	assert.True(t, shouldComplete)
	assert.Contains(t, reason, "target reached")
	assert.Contains(t, reason, "11/11")

	mockScorecardRepo.AssertExpectations(t)
}

func TestShouldCompleteMatch_EdgeCase_ExactWickets(t *testing.T) {
	// Setup
	mockScorecardRepo := &MockScorecardRepository{}
	mockMatchRepo := &MockMatchRepository{}
	service := services.NewScorecardService(mockScorecardRepo, mockMatchRepo)

	ctx := context.Background()
	matchID := "test-match-id"

	// First innings with 10 runs
	firstInnings := &models.Innings{
		MatchID:       matchID,
		InningsNumber: 1,
		TotalRuns:     10,
		Status:        string(models.InningsStatusCompleted),
	}

	// Second innings with exactly n-1 wickets (2 wickets for 3 players)
	secondInnings := &models.Innings{
		MatchID:       matchID,
		InningsNumber: 2,
		TotalRuns:     5, // Below target
		TotalWickets:  2, // Exactly n-1 wickets (3-1=2)
		TotalOvers:    1.0,
		Status:        string(models.InningsStatusInProgress),
	}

	// Match with 3 players per team
	match := &models.Match{
		ID:               matchID,
		TeamAPlayerCount: 3,
		TeamBPlayerCount: 3,
		TotalOvers:       2,
		Status:           models.MatchStatusLive,
	}

	// Mock expectations
	mockScorecardRepo.On("GetInningsByMatchAndNumber", ctx, matchID, 1).Return(firstInnings, nil)

	// Test
	shouldComplete, reason := service.ShouldCompleteMatch(ctx, matchID, secondInnings, match)

	// Assertions
	assert.True(t, shouldComplete)
	assert.Contains(t, reason, "all wickets lost")
	assert.Contains(t, reason, "2/2")

	mockScorecardRepo.AssertExpectations(t)
}
