package services

import (
	"context"
	"fmt"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/repository/interfaces"
	"time"
)

// PlayerService handles business logic for player operations
type PlayerService struct {
	playerRepo interfaces.PlayerRepository
	teamRepo   interfaces.TeamRepository
}

// NewPlayerService creates a new player service
func NewPlayerService(playerRepo interfaces.PlayerRepository, teamRepo interfaces.TeamRepository) *PlayerService {
	return &PlayerService{
		playerRepo: playerRepo,
		teamRepo:   teamRepo,
	}
}

// CreatePlayer creates a new player
func (s *PlayerService) CreatePlayer(ctx context.Context, req *models.CreatePlayerRequest) (*models.Player, error) {
	// Validate team exists
	_, err := s.teamRepo.GetByID(ctx, req.TeamID)
	if err != nil {
		return nil, fmt.Errorf("team not found: %w", err)
	}

	// Create player model
	player := &models.Player{
		Name:      req.Name,
		TeamID:    req.TeamID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save to repository
	err = s.playerRepo.Create(ctx, player)
	if err != nil {
		return nil, fmt.Errorf("failed to create player: %w", err)
	}

	return player, nil
}

// GetPlayer retrieves a player by ID
func (s *PlayerService) GetPlayer(ctx context.Context, id string) (*models.Player, error) {
	if id == "" {
		return nil, fmt.Errorf("player ID is required")
	}

	player, err := s.playerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get player: %w", err)
	}

	return player, nil
}

// ListPlayers retrieves all players with optional filtering
func (s *PlayerService) ListPlayers(ctx context.Context, filters *models.PlayerFilters) ([]*models.Player, error) {
	// Set default values
	if filters.Limit <= 0 {
		filters.Limit = 20
	}
	if filters.Limit > 100 {
		filters.Limit = 100
	}

	players, err := s.playerRepo.GetAll(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list players: %w", err)
	}

	return players, nil
}

// UpdatePlayer updates an existing player
func (s *PlayerService) UpdatePlayer(ctx context.Context, id string, req *models.UpdatePlayerRequest) (*models.Player, error) {
	if id == "" {
		return nil, fmt.Errorf("player ID is required")
	}

	// Get existing player
	player, err := s.playerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get player: %w", err)
	}

	// Update fields if provided
	if req.Name != nil {
		player.Name = *req.Name
	}
	if req.TeamID != nil {
		// Validate team exists
		_, err = s.teamRepo.GetByID(ctx, *req.TeamID)
		if err != nil {
			return nil, fmt.Errorf("team not found: %w", err)
		}
		player.TeamID = *req.TeamID
	}

	player.UpdatedAt = time.Now()

	// Save changes
	err = s.playerRepo.Update(ctx, id, player)
	if err != nil {
		return nil, fmt.Errorf("failed to update player: %w", err)
	}

	return player, nil
}

// DeletePlayer deletes a player
func (s *PlayerService) DeletePlayer(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("player ID is required")
	}

	// Check if player exists
	_, err := s.playerRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("player not found: %w", err)
	}

	// Delete player
	err = s.playerRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete player: %w", err)
	}

	return nil
}

// GetPlayersByTeam retrieves all players for a specific team
func (s *PlayerService) GetPlayersByTeam(ctx context.Context, teamID string) ([]*models.Player, error) {
	if teamID == "" {
		return nil, fmt.Errorf("team ID is required")
	}

	// Check if team exists
	_, err := s.teamRepo.GetByID(ctx, teamID)
	if err != nil {
		return nil, fmt.Errorf("team not found: %w", err)
	}

	// Get players for the team
	players, err := s.playerRepo.GetByTeamID(ctx, teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to get players by team: %w", err)
	}

	return players, nil
}
