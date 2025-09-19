package services

import (
	"context"
	"fmt"
	"log"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/repository/interfaces"
	"time"
)

// MatchService handles business logic for match operations
type MatchService struct {
	matchRepo  interfaces.MatchRepository
	seriesRepo interfaces.SeriesRepository
}

// NewMatchService creates a new match service
func NewMatchService(matchRepo interfaces.MatchRepository, seriesRepo interfaces.SeriesRepository) *MatchService {
	return &MatchService{
		matchRepo:  matchRepo,
		seriesRepo: seriesRepo,
	}
}

// CreateMatch creates a new match
func (s *MatchService) CreateMatch(ctx context.Context, req *models.CreateMatchRequest) (*models.Match, error) {
	fmt.Printf("DEBUG: MatchService.CreateMatch - Starting creation with request: %+v\n", req)

	// Get user ID from context
	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		fmt.Printf("DEBUG: MatchService.CreateMatch - Authentication failed: no user_id in context\n")
		return nil, fmt.Errorf("user authentication required")
	}
	fmt.Printf("DEBUG: MatchService.CreateMatch - User ID from context: %s\n", userID)
	log.Printf("Creating match for user ID: %s", userID)

	// Validate series exists
	fmt.Printf("DEBUG: MatchService.CreateMatch - Validating series exists with ID: %s\n", req.SeriesID)
	_, err := s.seriesRepo.GetByID(ctx, req.SeriesID)
	if err != nil {
		fmt.Printf("DEBUG: MatchService.CreateMatch - Series validation failed: %v\n", err)
		return nil, fmt.Errorf("series not found: %w", err)
	}
	fmt.Printf("DEBUG: MatchService.CreateMatch - Series validation successful\n")

	// Determine match number - use provided number or auto-increment
	var matchNumber int
	if req.MatchNumber != nil {
		matchNumber = *req.MatchNumber
		fmt.Printf("DEBUG: MatchService.CreateMatch - Using provided match number: %d\n", matchNumber)

		// Validate that the match number doesn't already exist for this series
		fmt.Printf("DEBUG: MatchService.CreateMatch - Checking match number uniqueness\n")
		exists, err := s.matchRepo.ExistsBySeriesAndMatchNumber(ctx, req.SeriesID, matchNumber)
		if err != nil {
			fmt.Printf("DEBUG: MatchService.CreateMatch - Failed to check match number uniqueness: %v\n", err)
			return nil, fmt.Errorf("failed to check match number uniqueness: %w", err)
		}
		if exists {
			fmt.Printf("DEBUG: MatchService.CreateMatch - Match number %d already exists for series %s\n", matchNumber, req.SeriesID)
			return nil, fmt.Errorf("match number %d already exists for series %s", matchNumber, req.SeriesID)
		}
		fmt.Printf("DEBUG: MatchService.CreateMatch - Match number %d is unique\n", matchNumber)
	} else {
		// Auto-increment match number for the series
		fmt.Printf("DEBUG: MatchService.CreateMatch - Getting next match number for series %s\n", req.SeriesID)
		matchNumber, err = s.matchRepo.GetNextMatchNumber(ctx, req.SeriesID)
		if err != nil {
			fmt.Printf("DEBUG: MatchService.CreateMatch - Failed to get next match number: %v\n", err)
			return nil, fmt.Errorf("failed to get next match number: %w", err)
		}
		fmt.Printf("DEBUG: MatchService.CreateMatch - Auto-assigned match number: %d\n", matchNumber)
	}

	// Create match model with toss winner as batting team by default
	match := &models.Match{
		SeriesID:         req.SeriesID,
		MatchNumber:      matchNumber,
		Date:             req.Date,
		Status:           models.MatchStatusLive, // Always live by default
		TeamAPlayerCount: req.TeamAPlayerCount,
		TeamBPlayerCount: req.TeamBPlayerCount,
		TotalOvers:       req.TotalOvers,
		TossWinner:       req.TossWinner,
		TossType:         req.TossType,
		BattingTeam:      req.TossWinner, // Winner of toss bats first
		CreatedBy:        userID,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	fmt.Printf("DEBUG: MatchService.CreateMatch - Created match model: %+v\n", match)

	// Save to repository
	fmt.Printf("DEBUG: MatchService.CreateMatch - Calling repository.Create\n")
	err = s.matchRepo.Create(ctx, match)
	if err != nil {
		fmt.Printf("DEBUG: MatchService.CreateMatch - Repository error: %v\n", err)
		return nil, fmt.Errorf("failed to create match: %w", err)
	}

	fmt.Printf("DEBUG: MatchService.CreateMatch - Successfully created match with ID: %s\n", match.ID)
	return match, nil
}

// GetMatch retrieves a match by ID
func (s *MatchService) GetMatch(ctx context.Context, id string) (*models.Match, error) {
	fmt.Printf("DEBUG: MatchService.GetMatch - Starting retrieval with ID: %s\n", id)

	if id == "" {
		fmt.Printf("DEBUG: MatchService.GetMatch - Validation failed: empty ID\n")
		return nil, fmt.Errorf("match ID is required")
	}

	fmt.Printf("DEBUG: MatchService.GetMatch - Calling repository.GetByID\n")
	match, err := s.matchRepo.GetByID(ctx, id)
	if err != nil {
		fmt.Printf("DEBUG: MatchService.GetMatch - Repository error: %v\n", err)
		return nil, fmt.Errorf("failed to get match: %w", err)
	}

	fmt.Printf("DEBUG: MatchService.GetMatch - Successfully retrieved match: %+v\n", match)
	return match, nil
}

// ListMatches retrieves all matches with optional filtering
func (s *MatchService) ListMatches(ctx context.Context, filters *models.MatchFilters) ([]*models.Match, error) {
	fmt.Printf("DEBUG: MatchService.ListMatches - Starting list with filters: %+v\n", filters)

	// Set default values
	if filters.Limit <= 0 {
		filters.Limit = 20
	}
	if filters.Limit > 100 {
		filters.Limit = 100
	}
	fmt.Printf("DEBUG: MatchService.ListMatches - Using filters: %+v\n", filters)

	fmt.Printf("DEBUG: MatchService.ListMatches - Calling repository.GetAll\n")
	matches, err := s.matchRepo.GetAll(ctx, filters)
	if err != nil {
		fmt.Printf("DEBUG: MatchService.ListMatches - Repository error: %v\n", err)
		return nil, fmt.Errorf("failed to list matches: %w", err)
	}

	fmt.Printf("DEBUG: MatchService.ListMatches - Successfully retrieved %d matches\n", len(matches))
	for i, m := range matches {
		fmt.Printf("DEBUG: MatchService.ListMatches - Match %d: ID=%s, SeriesID=%s, MatchNumber=%d, Status=%s\n", i+1, m.ID, m.SeriesID, m.MatchNumber, m.Status)
	}

	return matches, nil
}

// UpdateMatch updates an existing match
func (s *MatchService) UpdateMatch(ctx context.Context, id string, req *models.UpdateMatchRequest) (*models.Match, error) {
	fmt.Printf("DEBUG: MatchService.UpdateMatch - Starting update with ID: %s, request: %+v\n", id, req)

	if id == "" {
		fmt.Printf("DEBUG: MatchService.UpdateMatch - Validation failed: empty ID\n")
		return nil, fmt.Errorf("match ID is required")
	}

	// Get user ID from context
	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		fmt.Printf("DEBUG: MatchService.UpdateMatch - Authentication failed: no user_id in context\n")
		return nil, fmt.Errorf("user authentication required")
	}
	fmt.Printf("DEBUG: MatchService.UpdateMatch - User ID from context: %s\n", userID)

	// Get existing match
	fmt.Printf("DEBUG: MatchService.UpdateMatch - Getting existing match\n")
	match, err := s.matchRepo.GetByID(ctx, id)
	if err != nil {
		fmt.Printf("DEBUG: MatchService.UpdateMatch - Failed to get existing match: %v\n", err)
		return nil, fmt.Errorf("failed to get match: %w", err)
	}
	fmt.Printf("DEBUG: MatchService.UpdateMatch - Found existing match: %+v\n", match)

	// Check ownership
	if match.CreatedBy != userID {
		fmt.Printf("DEBUG: MatchService.UpdateMatch - Access denied: user %s cannot update match created by %s\n", userID, match.CreatedBy)
		return nil, fmt.Errorf("access denied: you can only update matches you created")
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
	if req.TeamAPlayerCount != nil {
		match.TeamAPlayerCount = *req.TeamAPlayerCount
	}
	if req.TeamBPlayerCount != nil {
		match.TeamBPlayerCount = *req.TeamBPlayerCount
	}
	if req.TotalOvers != nil {
		match.TotalOvers = *req.TotalOvers
	}
	if req.BattingTeam != nil {
		match.BattingTeam = *req.BattingTeam
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

	// Get user ID from context
	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		return fmt.Errorf("user authentication required")
	}

	// Check if match exists
	_, err := s.matchRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("match not found: %w", err)
	}

	// Check ownership
	if match.CreatedBy != userID {
		return fmt.Errorf("access denied: you can only delete matches you created")
	}

	// Cannot delete a live match
	if match.Status == models.MatchStatusLive {
		return fmt.Errorf("cannot delete a live match")
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
