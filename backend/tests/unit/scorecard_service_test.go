package unit

import (
	"context"
	"errors"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Note: MockScorecardRepository is already defined in match_completion_unit_test.go

func TestScorecardService_StartScoring(t *testing.T) {
	tests := []struct {
		name               string
		matchID            string
		match              *models.Match
		getMatchError      error
		existingInnings    []*models.Innings
		getInningsError    error
		createInningsError error
		expectedError      string
	}{
		{
			name:    "successful start scoring",
			matchID: "match-1",
			match: &models.Match{
				ID:         "match-1",
				Status:     models.MatchStatusLive,
				TossWinner: models.TeamTypeA,
			},
			existingInnings: []*models.Innings{},
			expectedError:   "",
		},
		{
			name:          "match not found",
			matchID:       "nonexistent-match",
			getMatchError: errors.New("match not found"),
			expectedError: "match not found",
		},
		{
			name:    "match not live",
			matchID: "match-1",
			match: &models.Match{
				ID:         "match-1",
				Status:     models.MatchStatusCancelled,
				TossWinner: models.TeamTypeA,
			},
			expectedError: "match is not live, cannot start scoring",
		},
		{
			name:    "scoring already started",
			matchID: "match-1",
			match: &models.Match{
				ID:         "match-1",
				Status:     models.MatchStatusLive,
				TossWinner: models.TeamTypeA,
			},
			existingInnings: []*models.Innings{
				{ID: "innings-1", MatchID: "match-1", InningsNumber: 1},
			},
			expectedError: "scoring already started for this match",
		},
		{
			name:    "create innings error",
			matchID: "match-1",
			match: &models.Match{
				ID:         "match-1",
				Status:     models.MatchStatusLive,
				TossWinner: models.TeamTypeA,
			},
			existingInnings:    []*models.Innings{},
			createInningsError: errors.New("database error"),
			expectedError:      "failed to start scoring",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockScorecardRepo := new(MockScorecardRepository)
			mockMatchRepo := new(MockMatchRepository)

			// Setup expectations
			if tt.match != nil {
				mockMatchRepo.On("GetByID", mock.Anything, tt.matchID).Return(tt.match, tt.getMatchError)
			} else {
				mockMatchRepo.On("GetByID", mock.Anything, tt.matchID).Return(nil, tt.getMatchError)
			}

			if tt.getMatchError == nil && tt.match.Status == models.MatchStatusLive {
				mockScorecardRepo.On("GetInningsByMatchID", mock.Anything, tt.matchID).Return(tt.existingInnings, tt.getInningsError)

				if len(tt.existingInnings) == 0 {
					mockScorecardRepo.On("CreateInnings", mock.Anything, mock.AnythingOfType("*models.Innings")).Return(tt.createInningsError)
				}
			}

			// Create service
			service := services.NewScorecardService(mockScorecardRepo, mockMatchRepo)

			// Test
			err := service.StartScoring(context.Background(), tt.matchID)

			// Assertions
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			// Verify all expectations were met
			mockScorecardRepo.AssertExpectations(t)
			mockMatchRepo.AssertExpectations(t)
		})
	}
}

func TestScorecardService_AddBall(t *testing.T) {
	tests := []struct {
		name          string
		req           *models.BallEventRequest
		match         *models.Match
		getMatchError error
		expectedError string
	}{
		{
			name: "match not found",
			req: &models.BallEventRequest{
				MatchID:       "nonexistent-match",
				InningsNumber: 1,
				BallType:      models.BallTypeGood,
				RunType:       models.RunTypeOne,
			},
			getMatchError: errors.New("match not found"),
			expectedError: "match not found",
		},
		{
			name: "match not live",
			req: &models.BallEventRequest{
				MatchID:       "match-1",
				InningsNumber: 1,
				BallType:      models.BallTypeGood,
				RunType:       models.RunTypeOne,
			},
			match: &models.Match{
				ID:     "match-1",
				Status: models.MatchStatusCancelled,
			},
			expectedError: "match is not live, cannot add ball",
		},
		{
			name: "invalid innings order - second innings without first",
			req: &models.BallEventRequest{
				MatchID:       "match-1",
				InningsNumber: 2,
				BallType:      models.BallTypeGood,
				RunType:       models.RunTypeOne,
			},
			match: &models.Match{
				ID:               "match-1",
				Status:           models.MatchStatusLive,
				TossWinner:       models.TeamTypeA,
				BattingTeam:      models.TeamTypeA,
				TeamAPlayerCount: 11,
				TotalOvers:       20,
			},
			expectedError: "cannot start second innings, first innings must be played first",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockScorecardRepo := new(MockScorecardRepository)
			mockMatchRepo := new(MockMatchRepository)

			// Setup expectations
			if tt.match != nil {
				mockMatchRepo.On("GetByID", mock.Anything, tt.req.MatchID).Return(tt.match, tt.getMatchError)
			} else {
				mockMatchRepo.On("GetByID", mock.Anything, tt.req.MatchID).Return(nil, tt.getMatchError)
			}

			// For live matches, we need to mock the innings validation
			if tt.getMatchError == nil && tt.match.Status == models.MatchStatusLive {
				// Mock GetInningsByMatchID for innings validation
				mockScorecardRepo.On("GetInningsByMatchID", mock.Anything, tt.req.MatchID).Return([]*models.Innings{}, nil)
			}

			// Create service
			service := services.NewScorecardService(mockScorecardRepo, mockMatchRepo)

			// Test
			err := service.AddBall(context.Background(), tt.req)

			// Assertions
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			// Verify all expectations were met
			mockScorecardRepo.AssertExpectations(t)
			mockMatchRepo.AssertExpectations(t)
		})
	}
}

func TestScorecardService_GetScorecard(t *testing.T) {
	tests := []struct {
		name              string
		matchID           string
		scorecard         *models.ScorecardResponse
		getScorecardError error
		expectedError     string
	}{
		{
			name:    "successful get scorecard",
			matchID: "match-1",
			scorecard: &models.ScorecardResponse{
				MatchID: "match-1",
			},
			expectedError: "",
		},
		{
			name:              "scorecard not found",
			matchID:           "nonexistent-match",
			getScorecardError: errors.New("scorecard not found"),
			expectedError:     "failed to get scorecard",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockScorecardRepo := new(MockScorecardRepository)
			mockMatchRepo := new(MockMatchRepository)

			// Setup expectations
			mockScorecardRepo.On("GetScorecard", mock.Anything, tt.matchID).Return(tt.scorecard, tt.getScorecardError)

			// Create service
			service := services.NewScorecardService(mockScorecardRepo, mockMatchRepo)

			// Test
			result, err := service.GetScorecard(context.Background(), tt.matchID)

			// Assertions
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.scorecard, result)
			}

			// Verify all expectations were met
			mockScorecardRepo.AssertExpectations(t)
			mockMatchRepo.AssertExpectations(t)
		})
	}
}

func TestScorecardService_GetCurrentOver(t *testing.T) {
	tests := []struct {
		name                string
		matchID             string
		inningsNumber       int
		innings             *models.Innings
		getInningsError     error
		over                *models.ScorecardOver
		getCurrentOverError error
		expectedError       string
	}{
		{
			name:          "successful get current over",
			matchID:       "match-1",
			inningsNumber: 1,
			innings: &models.Innings{
				ID: "innings-1",
			},
			over: &models.ScorecardOver{
				ID:         "over-1",
				OverNumber: 1,
			},
			expectedError: "",
		},
		{
			name:            "innings not found",
			matchID:         "match-1",
			inningsNumber:   1,
			getInningsError: errors.New("innings not found"),
			expectedError:   "innings not found",
		},
		{
			name:          "no current over found",
			matchID:       "match-1",
			inningsNumber: 1,
			innings: &models.Innings{
				ID: "innings-1",
			},
			getCurrentOverError: errors.New("no current over"),
			expectedError:       "no current over found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockScorecardRepo := new(MockScorecardRepository)
			mockMatchRepo := new(MockMatchRepository)

			// Setup expectations
			mockScorecardRepo.On("GetInningsByMatchAndNumber", mock.Anything, tt.matchID, tt.inningsNumber).Return(tt.innings, tt.getInningsError)

			if tt.innings != nil {
				mockScorecardRepo.On("GetCurrentOver", mock.Anything, tt.innings.ID).Return(tt.over, tt.getCurrentOverError)
			}

			// Create service
			service := services.NewScorecardService(mockScorecardRepo, mockMatchRepo)

			// Test
			result, err := service.GetCurrentOver(context.Background(), tt.matchID, tt.inningsNumber)

			// Assertions
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.over, result)
			}

			// Verify all expectations were met
			mockScorecardRepo.AssertExpectations(t)
			mockMatchRepo.AssertExpectations(t)
		})
	}
}

func TestScorecardService_ShouldCompleteMatch(t *testing.T) {
	tests := []struct {
		name                 string
		matchID              string
		secondInnings        *models.Innings
		match                *models.Match
		firstInnings         *models.Innings
		getFirstInningsError error
		expectedComplete     bool
		expectedReason       string
	}{
		{
			name:    "target reached",
			matchID: "match-1",
			secondInnings: &models.Innings{
				TotalRuns: 150,
			},
			match: &models.Match{
				TeamAPlayerCount: 11,
			},
			firstInnings: &models.Innings{
				TotalRuns: 140,
			},
			expectedComplete: true,
			expectedReason:   "target reached: 150/141",
		},
		{
			name:    "all wickets lost",
			matchID: "match-1",
			secondInnings: &models.Innings{
				TotalRuns:    100,
				TotalWickets: 10,
			},
			match: &models.Match{
				TeamAPlayerCount: 11,
			},
			firstInnings: &models.Innings{
				TotalRuns: 140,
			},
			expectedComplete: true,
			expectedReason:   "all wickets lost: 10/10",
		},
		{
			name:    "all overs completed",
			matchID: "match-1",
			secondInnings: &models.Innings{
				TotalRuns:  100,
				TotalOvers: 20.0,
			},
			match: &models.Match{
				TeamAPlayerCount: 11,
				TotalOvers:       20,
			},
			firstInnings: &models.Innings{
				TotalRuns: 140,
			},
			expectedComplete: true,
			expectedReason:   "all overs completed: 20.0/20",
		},
		{
			name:    "match continues",
			matchID: "match-1",
			secondInnings: &models.Innings{
				TotalRuns:    100,
				TotalWickets: 5,
				TotalOvers:   15.0,
			},
			match: &models.Match{
				TeamAPlayerCount: 11,
				TotalOvers:       20,
			},
			firstInnings: &models.Innings{
				TotalRuns: 140,
			},
			expectedComplete: false,
			expectedReason:   "match continues",
		},
		{
			name:    "error getting first innings",
			matchID: "match-1",
			secondInnings: &models.Innings{
				TotalRuns: 100,
			},
			match: &models.Match{
				TeamAPlayerCount: 11,
			},
			getFirstInningsError: errors.New("database error"),
			expectedComplete:     false,
			expectedReason:       "error getting first innings",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockScorecardRepo := new(MockScorecardRepository)
			mockMatchRepo := new(MockMatchRepository)

			// Setup expectations
			mockScorecardRepo.On("GetInningsByMatchAndNumber", mock.Anything, tt.matchID, 1).Return(tt.firstInnings, tt.getFirstInningsError)

			// Create service
			service := services.NewScorecardService(mockScorecardRepo, mockMatchRepo)

			// Test
			complete, reason := service.ShouldCompleteMatch(context.Background(), tt.matchID, tt.secondInnings, tt.match)

			// Assertions
			assert.Equal(t, tt.expectedComplete, complete)
			assert.Equal(t, tt.expectedReason, reason)

			// Verify all expectations were met
			mockScorecardRepo.AssertExpectations(t)
			mockMatchRepo.AssertExpectations(t)
		})
	}
}

func TestScorecardService_GetNonTossWinner(t *testing.T) {
	service := services.NewScorecardService(nil, nil)

	tests := []struct {
		name       string
		tossWinner models.TeamType
		expected   models.TeamType
	}{
		{
			name:       "TeamA wins toss",
			tossWinner: models.TeamTypeA,
			expected:   models.TeamTypeB,
		},
		{
			name:       "TeamB wins toss",
			tossWinner: models.TeamTypeB,
			expected:   models.TeamTypeA,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.GetNonTossWinner(tt.tossWinner)
			assert.Equal(t, tt.expected, result)
		})
	}
}
