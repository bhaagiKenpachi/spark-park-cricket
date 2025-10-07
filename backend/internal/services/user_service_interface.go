package services

import (
	"context"
	"spark-park-cricket-backend/internal/models"
)

// UserServiceInterface defines the interface for user operations
type UserServiceInterface interface {
	// UpdateUserName updates the user's display name
	UpdateUserName(ctx context.Context, userID string, name string) (*models.User, error)
}
