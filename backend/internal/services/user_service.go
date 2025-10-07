package services

import (
	"context"
	"fmt"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/repository/interfaces"
	"spark-park-cricket-backend/internal/utils"
	"strings"
)

// UserService handles user-related operations
type UserService struct {
	userRepo interfaces.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo interfaces.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// UpdateUserName updates the user's display name
func (s *UserService) UpdateUserName(ctx context.Context, userID string, name string) (*models.User, error) {
	// Validate input
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}
	if len(name) < 2 {
		return nil, fmt.Errorf("name must be at least 2 characters long")
	}
	if len(name) > 255 {
		return nil, fmt.Errorf("name must be at most 255 characters long")
	}

	// Get existing user to verify they exist
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		utils.LogError(err, "Failed to get user", map[string]interface{}{
			"user_id": userID,
		})
		return nil, fmt.Errorf("user not found")
	}

	// Update user name
	updates := &models.UpdateUserRequest{
		Name: &name,
	}

	if err := s.userRepo.UpdateUser(ctx, userID, updates); err != nil {
		utils.LogError(err, "Failed to update user name", map[string]interface{}{
			"user_id": userID,
			"name":    name,
		})
		return nil, fmt.Errorf("failed to update user name: %w", err)
	}

	// Get updated user
	user, err = s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		utils.LogError(err, "Failed to get updated user", map[string]interface{}{
			"user_id": userID,
		})
		return nil, fmt.Errorf("failed to get updated user: %w", err)
	}

	utils.LogInfo("User name updated successfully", map[string]interface{}{
		"user_id": userID,
		"name":    name,
	})

	return user, nil
}
