package services

import (
	"context"
	"fmt"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/repository/interfaces"
	"time"
)

// MatchService handles business logic for match operations
type MatchService struct {
	matchRepo  interfaces.MatchRepository
	seriesRepo interfaces.SeriesRepository
	teamRepo   interfaces.TeamRepository
}

// NewMatchService creates a new match service
func NewMatchService(matchRepo interfaces.MatchRepository, seriesRepo interfaces.SeriesRepository, teamRepo interfaces.TeamRepository) *MatchService {
	return &MatchService{
		matchRepo:  matchRepo,
		seriesRepo: seriesRepo,
		teamRepo:   teamRepo,
	}
}

// CreateMatch creates a new match
func (s *MatchService) CreateMatch(ctx context.Context, req *models.CreateMatchRequest) (*models.Match, error) {
	// Validate series exists
	_, err := s.seriesRepo.GetByID(ctx, req.SeriesID)
	if err != nil {
		return nil, fmt.Errorf("series not found: %w", err)
	}

	// Validate teams exist
	_, err = s.teamRepo.GetByID(ctx, req.Team1ID)
	if err != nil {
		return nil, fmt.Errorf("team1 not found: %w", err)
	}

	_, err = s.teamRepo.GetByID(ctx, req.Team2ID)
	if err != nil {
		return nil, fmt.Errorf("team2 not found: %w", err)
	}

	// Validate teams are different
	if req.Team1ID == req.Team2ID {
		return nil, fmt.Errorf("team1 and team2 must be different")
	}

	// Create match model
	match := &models.Match{
		SeriesID:    req.SeriesID,
		MatchNumber: req.MatchNumber,
		Date:        req.Date,
		Status:      models.MatchStatusScheduled,
		Team1ID:     req.Team1ID,
		Team2ID:     req.Team2ID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Save to repository
	err = s.matchRepo.Create(ctx, match)
	if err != nil {
		return nil, fmt.Errorf("failed to create match: %w", err)
	}

	return match, nil
}

// GetMatch retrieves a match by ID
func (s *MatchService) GetMatch(ctx context.Context, id string) (*models.Match, error) {
	if id == "" {
		return nil, fmt.Errorf("match ID is required")
	}

	match, err := s.matchRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get match: %w", err)
	}

	return match, nil
}

// ListMatches retrieves all matches with optional filtering
func (s *MatchService) ListMatches(ctx context.Context, filters *models.MatchFilters) ([]*models.Match, error) {
	// Set default values
	if filters.Limit <= 0 {
		filters.Limit = 20
	}
	if filters.Limit > 100 {
		filters.Limit = 100
	}

	matches, err := s.matchRepo.GetAll(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list matches: %w", err)
	}

	return matches, nil
}

// UpdateMatch updates an existing match
func (s *MatchService) UpdateMatch(ctx context.Context, id string, req *models.UpdateMatchRequest) (*models.Match, error) {
	if id == "" {
		return nil, fmt.Errorf("match ID is required")
	}

	// Get existing match
	match, err := s.matchRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get match: %w", err)
	}

	// Update fields if provided
	if req.MatchNumber != nil {
		match.MatchNumber = *req.MatchNumber
	}
	if req.Date != nil {
		match.Date = *req.Date
	}
	if req.Status != nil {
		match.Status = *req.Status
	}
	if req.Team1ID != nil {
		match.Team1ID = *req.Team1ID
	}
	if req.Team2ID != nil {
		match.Team2ID = *req.Team2ID
	}

	// Validate business rules
	if match.Team1ID == match.Team2ID {
		return nil, fmt.Errorf("team1 and team2 must be different")
	}

	match.UpdatedAt = time.Now()

	// Save changes
	err = s.matchRepo.Update(ctx, id, match)
	if err != nil {
		return nil, fmt.Errorf("failed to update match: %w", err)
	}

	return match, nil
}

// DeleteMatch deletes a match
func (s *MatchService) DeleteMatch(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("match ID is required")
	}

	// Check if match exists
	_, err := s.matchRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("match not found: %w", err)
	}

	// Delete match
	err = s.matchRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete match: %w", err)
	}

	return nil
}

// GetMatchesBySeries retrieves all matches for a specific series
func (s *MatchService) GetMatchesBySeries(ctx context.Context, seriesID string) ([]*models.Match, error) {
	if seriesID == "" {
		return nil, fmt.Errorf("series ID is required")
	}

	matches, err := s.matchRepo.GetBySeriesID(ctx, seriesID)
	if err != nil {
		return nil, fmt.Errorf("failed to get matches by series: %w", err)
	}

	return matches, nil
}
