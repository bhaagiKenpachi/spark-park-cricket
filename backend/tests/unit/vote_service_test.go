package unit

import (
	"context"
	"errors"
	"testing"

	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockVoteRepository is a mock implementation of VoteRepositoryInterface
type MockVoteRepository struct {
	mock.Mock
}

func (m *MockVoteRepository) CreateVote(ctx context.Context, vote *models.Vote) error {
	args := m.Called(ctx, vote)
	return args.Error(0)
}

func (m *MockVoteRepository) GetVoteByID(ctx context.Context, id string) (*models.Vote, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Vote), args.Error(1)
}

func (m *MockVoteRepository) GetVoteWithOptions(ctx context.Context, id string) (*models.VoteWithOptions, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.VoteWithOptions), args.Error(1)
}

func (m *MockVoteRepository) GetVoteWithResults(ctx context.Context, id string, userID string) (*models.VoteWithResults, error) {
	args := m.Called(ctx, id, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.VoteWithResults), args.Error(1)
}

func (m *MockVoteRepository) UpdateVote(ctx context.Context, vote *models.Vote) error {
	args := m.Called(ctx, vote)
	return args.Error(0)
}

func (m *MockVoteRepository) DeleteVote(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockVoteRepository) ListVotes(ctx context.Context, filters *models.VoteFilters) ([]*models.Vote, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Vote), args.Error(1)
}

func (m *MockVoteRepository) CreateVoteOptions(ctx context.Context, options []*models.VoteOption) error {
	args := m.Called(ctx, options)
	return args.Error(0)
}

func (m *MockVoteRepository) GetVoteOptions(ctx context.Context, voteID string) ([]*models.VoteOption, error) {
	args := m.Called(ctx, voteID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.VoteOption), args.Error(1)
}

func (m *MockVoteRepository) UpdateVoteOption(ctx context.Context, option *models.VoteOption) error {
	args := m.Called(ctx, option)
	return args.Error(0)
}

func (m *MockVoteRepository) DeleteVoteOption(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockVoteRepository) CreateUserVote(ctx context.Context, userVote *models.UserVote) error {
	args := m.Called(ctx, userVote)
	return args.Error(0)
}

func (m *MockVoteRepository) GetUserVote(ctx context.Context, voteID, userID string) (*models.UserVote, error) {
	args := m.Called(ctx, voteID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserVote), args.Error(1)
}

func (m *MockVoteRepository) HasUserVoted(ctx context.Context, voteID, userID string) (bool, error) {
	args := m.Called(ctx, voteID, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockVoteRepository) GetVoteResults(ctx context.Context, voteID string) (map[string]int, error) {
	args := m.Called(ctx, voteID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]int), args.Error(1)
}

func (m *MockVoteRepository) GetVotedUsers(ctx context.Context, voteID string) ([]string, error) {
	args := m.Called(ctx, voteID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockVoteRepository) GetTotalVoteCount(ctx context.Context, voteID string) (int, error) {
	args := m.Called(ctx, voteID)
	return args.Int(0), args.Error(1)
}

func TestVoteService_CreateVote(t *testing.T) {
	tests := []struct {
		name          string
		request       *models.CreateVoteRequest
		userID        string
		mockSetup     func(*MockVoteRepository)
		expectedError string
	}{
		{
			name: "successful vote creation",
			request: &models.CreateVoteRequest{
				Title:       "Test Vote",
				Description: "This is a test vote description",
				Type:        models.VoteTypeSingle,
				Options:     []string{"Option 1", "Option 2"},
			},
			userID: "user-123",
			mockSetup: func(mockRepo *MockVoteRepository) {
				mockRepo.On("CreateVote", mock.Anything, mock.AnythingOfType("*models.Vote")).Return(nil)
				mockRepo.On("CreateVoteOptions", mock.Anything, mock.AnythingOfType("[]*models.VoteOption")).Return(nil)
			},
		},
		{
			name: "validation error - empty title",
			request: &models.CreateVoteRequest{
				Title:       "",
				Description: "This is a test vote description",
				Type:        models.VoteTypeSingle,
				Options:     []string{"Option 1", "Option 2"},
			},
			userID:        "user-123",
			mockSetup:     func(mockRepo *MockVoteRepository) {},
			expectedError: "title is required",
		},
		{
			name: "validation error - too few options",
			request: &models.CreateVoteRequest{
				Title:       "Test Vote",
				Description: "This is a test vote description",
				Type:        models.VoteTypeSingle,
				Options:     []string{"Option 1"},
			},
			userID:        "user-123",
			mockSetup:     func(mockRepo *MockVoteRepository) {},
			expectedError: "at least 2 options are required",
		},
		{
			name: "repository error",
			request: &models.CreateVoteRequest{
				Title:       "Test Vote",
				Description: "This is a test vote description",
				Type:        models.VoteTypeSingle,
				Options:     []string{"Option 1", "Option 2"},
			},
			userID: "user-123",
			mockSetup: func(mockRepo *MockVoteRepository) {
				mockRepo.On("CreateVote", mock.Anything, mock.AnythingOfType("*models.Vote")).Return(errors.New("database error"))
			},
			expectedError: "failed to create vote: database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockVoteRepository)
			tt.mockSetup(mockRepo)

			service := services.NewVoteService(mockRepo)
			vote, err := service.CreateVote(context.Background(), tt.request, tt.userID)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, vote)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, vote)
				assert.Equal(t, tt.request.Title, vote.Title)
				assert.Equal(t, tt.request.Description, vote.Description)
				assert.Equal(t, tt.request.Type, vote.Type)
				assert.Equal(t, tt.userID, vote.CreatedBy)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestVoteService_CastVote(t *testing.T) {
	tests := []struct {
		name          string
		voteID        string
		request       *models.VoteRequest
		userID        string
		mockSetup     func(*MockVoteRepository)
		expectedError string
	}{
		{
			name:   "successful vote casting",
			voteID: "vote-123",
			request: &models.VoteRequest{
				SelectedOptions: []string{"option-1"},
			},
			userID: "user-123",
			mockSetup: func(mockRepo *MockVoteRepository) {
				vote := &models.Vote{
					ID:     "vote-123",
					Status: models.VoteStatusActive,
					Type:   models.VoteTypeSingle,
				}
				options := []*models.VoteOption{
					{ID: "option-1", VoteID: "vote-123", Text: "Option 1"},
					{ID: "option-2", VoteID: "vote-123", Text: "Option 2"},
				}
				mockRepo.On("GetVoteByID", mock.Anything, "vote-123").Return(vote, nil)
				mockRepo.On("HasUserVoted", mock.Anything, "vote-123", "user-123").Return(false, nil)
				mockRepo.On("GetVoteOptions", mock.Anything, "vote-123").Return(options, nil)
				mockRepo.On("CreateUserVote", mock.Anything, mock.AnythingOfType("*models.UserVote")).Return(nil)
			},
		},
		{
			name:   "user already voted",
			voteID: "vote-123",
			request: &models.VoteRequest{
				SelectedOptions: []string{"option-1"},
			},
			userID: "user-123",
			mockSetup: func(mockRepo *MockVoteRepository) {
				vote := &models.Vote{
					ID:     "vote-123",
					Status: models.VoteStatusActive,
					Type:   models.VoteTypeSingle,
				}
				mockRepo.On("GetVoteByID", mock.Anything, "vote-123").Return(vote, nil)
				mockRepo.On("HasUserVoted", mock.Anything, "vote-123", "user-123").Return(true, nil)
			},
			expectedError: "user has already voted on this poll",
		},
		{
			name:   "vote is closed",
			voteID: "vote-123",
			request: &models.VoteRequest{
				SelectedOptions: []string{"option-1"},
			},
			userID: "user-123",
			mockSetup: func(mockRepo *MockVoteRepository) {
				vote := &models.Vote{
					ID:     "vote-123",
					Status: models.VoteStatusClosed,
					Type:   models.VoteTypeSingle,
				}
				mockRepo.On("GetVoteByID", mock.Anything, "vote-123").Return(vote, nil)
			},
			expectedError: "cannot vote on closed or cancelled vote",
		},
		{
			name:   "invalid option selected",
			voteID: "vote-123",
			request: &models.VoteRequest{
				SelectedOptions: []string{"invalid-option"},
			},
			userID: "user-123",
			mockSetup: func(mockRepo *MockVoteRepository) {
				vote := &models.Vote{
					ID:     "vote-123",
					Status: models.VoteStatusActive,
					Type:   models.VoteTypeSingle,
				}
				options := []*models.VoteOption{
					{ID: "option-1", VoteID: "vote-123", Text: "Option 1"},
					{ID: "option-2", VoteID: "vote-123", Text: "Option 2"},
				}
				mockRepo.On("GetVoteByID", mock.Anything, "vote-123").Return(vote, nil)
				mockRepo.On("HasUserVoted", mock.Anything, "vote-123", "user-123").Return(false, nil)
				mockRepo.On("GetVoteOptions", mock.Anything, "vote-123").Return(options, nil)
			},
			expectedError: "invalid option selected: invalid-option",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockVoteRepository)
			tt.mockSetup(mockRepo)

			service := services.NewVoteService(mockRepo)
			err := service.CastVote(context.Background(), tt.voteID, tt.request, tt.userID)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestVoteService_GetVoteWithResults(t *testing.T) {
	tests := []struct {
		name      string
		voteID    string
		userID    string
		mockSetup func(*MockVoteRepository)
		expected  *models.VoteWithResults
	}{
		{
			name:   "successful retrieval",
			voteID: "vote-123",
			userID: "user-123",
			mockSetup: func(mockRepo *MockVoteRepository) {
				vote := &models.Vote{ID: "vote-123", Title: "Test Vote"}
				options := []*models.VoteOption{
					{ID: "option-1", VoteID: "vote-123", Text: "Option 1"},
				}
				results := map[string]int{"option-1": 5}
				votedUsers := []string{"user-1", "user-2"}
				userVote := &models.UserVote{UserID: "user-123", SelectedOptions: []string{"option-1"}}

				voteWithResults := &models.VoteWithResults{
					Vote:       vote,
					Options:    options,
					Results:    results,
					UserVote:   userVote,
					TotalVotes: 2,
					VotedUsers: votedUsers,
				}

				mockRepo.On("GetVoteWithResults", mock.Anything, "vote-123", "user-123").Return(voteWithResults, nil)
			},
			expected: &models.VoteWithResults{
				Vote:       &models.Vote{ID: "vote-123", Title: "Test Vote"},
				Options:    []*models.VoteOption{{ID: "option-1", VoteID: "vote-123", Text: "Option 1"}},
				Results:    map[string]int{"option-1": 5},
				UserVote:   &models.UserVote{UserID: "user-123", SelectedOptions: []string{"option-1"}},
				TotalVotes: 2,
				VotedUsers: []string{"user-1", "user-2"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockVoteRepository)
			tt.mockSetup(mockRepo)

			service := services.NewVoteService(mockRepo)
			result, err := service.GetVoteWithResults(context.Background(), tt.voteID, tt.userID)

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tt.expected.Vote.ID, result.Vote.ID)
			assert.Equal(t, tt.expected.TotalVotes, result.TotalVotes)
			assert.Equal(t, len(tt.expected.VotedUsers), len(result.VotedUsers))

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestVoteService_UpdateVote(t *testing.T) {
	tests := []struct {
		name          string
		voteID        string
		request       *models.UpdateVoteRequest
		userID        string
		mockSetup     func(*MockVoteRepository)
		expectedError string
	}{
		{
			name:   "successful update",
			voteID: "vote-123",
			request: &models.UpdateVoteRequest{
				Title: stringPtr("Updated Title"),
			},
			userID: "creator-123",
			mockSetup: func(mockRepo *MockVoteRepository) {
				vote := &models.Vote{
					ID:        "vote-123",
					CreatedBy: "creator-123",
					Status:    models.VoteStatusActive,
				}
				mockRepo.On("GetVoteByID", mock.Anything, "vote-123").Return(vote, nil)
				mockRepo.On("UpdateVote", mock.Anything, mock.AnythingOfType("*models.Vote")).Return(nil)
			},
		},
		{
			name:   "unauthorized update",
			voteID: "vote-123",
			request: &models.UpdateVoteRequest{
				Title: stringPtr("Updated Title"),
			},
			userID: "other-user-123",
			mockSetup: func(mockRepo *MockVoteRepository) {
				vote := &models.Vote{
					ID:        "vote-123",
					CreatedBy: "creator-123",
					Status:    models.VoteStatusActive,
				}
				mockRepo.On("GetVoteByID", mock.Anything, "vote-123").Return(vote, nil)
			},
			expectedError: "unauthorized: only vote creator can update vote",
		},
		{
			name:   "cannot update closed vote",
			voteID: "vote-123",
			request: &models.UpdateVoteRequest{
				Title: stringPtr("Updated Title"),
			},
			userID: "creator-123",
			mockSetup: func(mockRepo *MockVoteRepository) {
				vote := &models.Vote{
					ID:        "vote-123",
					CreatedBy: "creator-123",
					Status:    models.VoteStatusClosed,
				}
				mockRepo.On("GetVoteByID", mock.Anything, "vote-123").Return(vote, nil)
			},
			expectedError: "cannot update closed or cancelled vote",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockVoteRepository)
			tt.mockSetup(mockRepo)

			service := services.NewVoteService(mockRepo)
			vote, err := service.UpdateVote(context.Background(), tt.voteID, tt.request, tt.userID)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, vote)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, vote)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}
