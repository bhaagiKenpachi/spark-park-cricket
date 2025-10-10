package services

import (
	"context"
	"fmt"
	"log"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/repository/interfaces"
	"time"
)

// TeamService handles business logic for team operations
type TeamService struct {
	teamRepo   interfaces.TeamRepository
	playerRepo interfaces.PlayerRepository
}

// NewTeamService creates a new team service
func NewTeamService(teamRepo interfaces.TeamRepository, playerRepo interfaces.PlayerRepository) *TeamService {
	return &TeamService{
		teamRepo:   teamRepo,
		playerRepo: playerRepo,
	}
}

// CreateTeam creates a new team
func (s *TeamService) CreateTeam(ctx context.Context, req *models.CreateTeamRequest) (*models.Team, error) {
	log.Printf("DEBUG: CreateTeam called with request: %+v", req)

	// Validate business rules
	if req.PlayersCount < 1 || req.PlayersCount > 20 {
		log.Printf("DEBUG: Validation failed - players count out of range: %d", req.PlayersCount)
		return nil, fmt.Errorf("players count must be between 1 and 20")
	}

	// Create team model
	team := &models.Team{
		Name:         req.Name,
		PlayersCount: req.PlayersCount,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	log.Printf("DEBUG: Created team model: ID='%s', Name='%s', PlayersCount=%d", team.ID, team.Name, team.PlayersCount)

	// Save to repository
	log.Printf("DEBUG: Calling teamRepo.Create with team: %+v", team)
	err := s.teamRepo.Create(ctx, team)
	if err != nil {
		log.Printf("DEBUG: teamRepo.Create failed: %v", err)
		return nil, fmt.Errorf("failed to create team: %w", err)
	}

	log.Printf("DEBUG: Team created successfully: %+v", team)
	return team, nil
}

// GetTeam retrieves a team by ID
func (s *TeamService) GetTeam(ctx context.Context, id string) (*models.Team, error) {
	if id == "" {
		return nil, fmt.Errorf("team ID is required")
	}

	team, err := s.teamRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get team: %w", err)
	}

	return team, nil
}

// ListTeams retrieves all teams with optional filtering
func (s *TeamService) ListTeams(ctx context.Context, filters *models.TeamFilters) ([]*models.Team, error) {
	// Set default values
	if filters.Limit <= 0 {
		filters.Limit = 20
	}
	if filters.Limit > 100 {
		filters.Limit = 100
	}

	teams, err := s.teamRepo.GetAll(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list teams: %w", err)
	}

	return teams, nil
}

// UpdateTeam updates an existing team
func (s *TeamService) UpdateTeam(ctx context.Context, id string, req *models.UpdateTeamRequest) (*models.Team, error) {
	if id == "" {
		return nil, fmt.Errorf("team ID is required")
	}

	// Get existing team
	team, err := s.teamRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get team: %w", err)
	}

	// Update fields if provided
	if req.Name != nil {
		team.Name = *req.Name
	}
	if req.PlayersCount != nil {
		// Validate business rules
		if *req.PlayersCount < 1 || *req.PlayersCount > 20 {
			return nil, fmt.Errorf("players count must be between 1 and 20")
		}
		team.PlayersCount = *req.PlayersCount
	}

	team.UpdatedAt = time.Now()

	// Save changes
	err = s.teamRepo.Update(ctx, id, team)
	if err != nil {
		return nil, fmt.Errorf("failed to update team: %w", err)
	}

	return team, nil
}

// DeleteTeam deletes a team
func (s *TeamService) DeleteTeam(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("team ID is required")
	}

	// Check if team exists
	_, err := s.teamRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("team not found: %w", err)
	}

	// Delete team
	err = s.teamRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete team: %w", err)
	}

	return nil
}

// GetTeamPlayers retrieves all players for a specific team
func (s *TeamService) GetTeamPlayers(ctx context.Context, teamID string) ([]*models.Player, error) {
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
		return nil, fmt.Errorf("failed to get team players: %w", err)
	}

	return players, nil
}

// AddPlayerToTeam adds a player to a team
func (s *TeamService) AddPlayerToTeam(ctx context.Context, teamID string, req *models.CreatePlayerRequest) (*models.Player, error) {
	if teamID == "" {
		return nil, fmt.Errorf("team ID is required")
	}

	// Check if team exists
	_, err := s.teamRepo.GetByID(ctx, teamID)
	if err != nil {
		return nil, fmt.Errorf("team not found: %w", err)
	}

	// Create player model
	player := &models.Player{
		Name:      req.Name,
		TeamID:    teamID,
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
