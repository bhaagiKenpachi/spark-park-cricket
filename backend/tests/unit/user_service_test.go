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

func TestUserService_UpdateUserName(t *testing.T) {
	tests := []struct {
		name          string
		userID        string
		newName       string
		mockSetup     func(*MockUserRepository)
		expectedError string
	}{
		{
			name:    "Success - Valid name update",
			userID:  "user-123",
			newName: "John Doe",
			mockSetup: func(m *MockUserRepository) {
				// First GetUserByID call to verify user exists
				m.On("GetUserByID", mock.Anything, "user-123").Return(&models.User{
					ID:   "user-123",
					Name: "Old Name",
				}, nil).Once()

				// UpdateUser call
				m.On("UpdateUser", mock.Anything, "user-123", mock.MatchedBy(func(updates *models.UpdateUserRequest) bool {
					return updates.Name != nil && *updates.Name == "John Doe"
				})).Return(nil).Once()

				// Second GetUserByID call to return updated user
				m.On("GetUserByID", mock.Anything, "user-123").Return(&models.User{
					ID:   "user-123",
					Name: "John Doe",
				}, nil).Once()
			},
			expectedError: "",
		},
		{
			name:    "Error - Empty name",
			userID:  "user-123",
			newName: "",
			mockSetup: func(m *MockUserRepository) {
				// No mock setup needed as validation should fail before DB call
			},
			expectedError: "name cannot be empty",
		},
		{
			name:    "Error - Name too short",
			userID:  "user-123",
			newName: "A",
			mockSetup: func(m *MockUserRepository) {
				// No mock setup needed as validation should fail before DB call
			},
			expectedError: "name must be at least 2 characters long",
		},
		{
			name:    "Error - Name too long",
			userID:  "user-123",
			newName: string(make([]byte, 256)),
			mockSetup: func(m *MockUserRepository) {
				// No mock setup needed as validation should fail before DB call
			},
			expectedError: "name must be at most 255 characters long",
		},
		{
			name:    "Error - User not found",
			userID:  "non-existent",
			newName: "John Doe",
			mockSetup: func(m *MockUserRepository) {
				m.On("GetUserByID", mock.Anything, "non-existent").Return(nil, errors.New("user not found"))
			},
			expectedError: "user not found",
		},
		{
			name:    "Error - Update fails",
			userID:  "user-123",
			newName: "John Doe",
			mockSetup: func(m *MockUserRepository) {
				// First GetUserByID call succeeds
				m.On("GetUserByID", mock.Anything, "user-123").Return(&models.User{
					ID:   "user-123",
					Name: "Old Name",
				}, nil).Once()

				// UpdateUser call fails
				m.On("UpdateUser", mock.Anything, "user-123", mock.Anything).Return(errors.New("database error"))
			},
			expectedError: "failed to update user name",
		},
		{
			name:    "Success - Name with whitespace trimmed",
			userID:  "user-123",
			newName: "  John Doe  ",
			mockSetup: func(m *MockUserRepository) {
				// First GetUserByID call to verify user exists
				m.On("GetUserByID", mock.Anything, "user-123").Return(&models.User{
					ID:   "user-123",
					Name: "Old Name",
				}, nil).Once()

				// UpdateUser call with trimmed name
				m.On("UpdateUser", mock.Anything, "user-123", mock.MatchedBy(func(updates *models.UpdateUserRequest) bool {
					return updates.Name != nil && *updates.Name == "John Doe"
				})).Return(nil).Once()

				// Second GetUserByID call to return updated user
				m.On("GetUserByID", mock.Anything, "user-123").Return(&models.User{
					ID:   "user-123",
					Name: "John Doe",
				}, nil).Once()
			},
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := new(MockUserRepository)
			tt.mockSetup(mockRepo)

			service := services.NewUserService(mockRepo)

			// Execute
			user, err := service.UpdateUserName(context.Background(), tt.userID, tt.newName)

			// Assert
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.userID, user.ID)
				// Check name is trimmed
				expectedName := "John Doe"
				assert.Equal(t, expectedName, user.Name)
			}

			// Verify all expected calls were made
			mockRepo.AssertExpectations(t)
		})
	}
}
