package mocks

import (
	"context"
	"spark-park-cricket-backend/internal/models"

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

// GetNonTossWinner returns the non-toss winner team type
func (m *MockScorecardService) GetNonTossWinner(tossWinner models.TeamType) models.TeamType {
	args := m.Called(tossWinner)
	return args.Get(0).(models.TeamType)
}
