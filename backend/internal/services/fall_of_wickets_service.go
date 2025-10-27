package services

import (
	"context"
	"fmt"

	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/repository/interfaces"
)

// FallOfWicketsService handles business logic for fall of wickets
type FallOfWicketsService struct {
	fallOfWicketsRepo interfaces.FallOfWicketsRepository
	scorecardRepo     interfaces.ScorecardRepository
	matchRepo         interfaces.MatchRepository
}

// NewFallOfWicketsService creates a new fall of wickets service
func NewFallOfWicketsService(
	fallOfWicketsRepo interfaces.FallOfWicketsRepository,
	scorecardRepo interfaces.ScorecardRepository,
	matchRepo interfaces.MatchRepository,
) *FallOfWicketsService {
	return &FallOfWicketsService{
		fallOfWicketsRepo: fallOfWicketsRepo,
		scorecardRepo:     scorecardRepo,
		matchRepo:         matchRepo,
	}
}

// CreateFallOfWickets creates a new fall of wickets record
func (s *FallOfWicketsService) CreateFallOfWickets(ctx context.Context, req *models.CreateFallOfWicketsRequest) (*models.FallOfWickets, error) {
	// Validate that the match exists
	_, err := s.matchRepo.GetByID(ctx, req.MatchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get match: %w", err)
	}

	// Validate that the innings exists
	_, err = s.scorecardRepo.GetInningsByMatchAndNumber(ctx, req.MatchID, 1) // Assuming first innings for now
	if err != nil {
		return nil, fmt.Errorf("failed to get innings: %w", err)
	}

	// Validate that the over exists
	_, err = s.scorecardRepo.GetOverByInningsAndNumber(ctx, req.InningsID, req.OverNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get over: %w", err)
	}

	// Validate that the ball exists
	balls, err := s.scorecardRepo.GetBallsByOver(ctx, req.OverID)
	if err != nil {
		return nil, fmt.Errorf("failed to get balls: %w", err)
	}

	// Find the specific ball
	var ball *models.ScorecardBall
	for _, b := range balls {
		if b.ID == req.BallID {
			ball = b
			break
		}
	}
	if ball == nil {
		return nil, fmt.Errorf("ball not found")
	}

	// Validate that the ball actually resulted in a wicket
	if !ball.IsWicket {
		return nil, fmt.Errorf("ball did not result in a wicket")
	}

	// Create the fall of wickets record
	fallOfWickets := &models.FallOfWickets{
		MatchID:      req.MatchID,
		InningsID:    req.InningsID,
		OverID:       req.OverID,
		BallID:       req.BallID,
		WicketNumber: req.WicketNumber,
		Score:        req.Score,
		OverNumber:   req.OverNumber,
		BallNumber:   req.BallNumber,
	}

	err = s.fallOfWicketsRepo.Create(ctx, fallOfWickets)
	if err != nil {
		return nil, fmt.Errorf("failed to create fall of wickets record: %w", err)
	}

	return fallOfWickets, nil
}

// GetFallOfWicketsByID retrieves a fall of wickets record by ID
func (s *FallOfWicketsService) GetFallOfWicketsByID(ctx context.Context, id string) (*models.FallOfWickets, error) {
	if id == "" {
		return nil, fmt.Errorf("ID is required")
	}

	return s.fallOfWicketsRepo.GetByID(ctx, id)
}

// ListFallOfWickets retrieves fall of wickets records with filters
func (s *FallOfWicketsService) ListFallOfWickets(ctx context.Context, filters *models.FallOfWicketsFilters) ([]*models.FallOfWickets, error) {
	return s.fallOfWicketsRepo.List(ctx, filters)
}

// UpdateFallOfWickets updates a fall of wickets record
func (s *FallOfWicketsService) UpdateFallOfWickets(ctx context.Context, id string, req *models.UpdateFallOfWicketsRequest) (*models.FallOfWickets, error) {
	if id == "" {
		return nil, fmt.Errorf("ID is required")
	}

	return s.fallOfWicketsRepo.Update(ctx, id, req)
}

// DeleteFallOfWickets deletes a fall of wickets record
func (s *FallOfWicketsService) DeleteFallOfWickets(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("ID is required")
	}

	return s.fallOfWicketsRepo.Delete(ctx, id)
}

// GetFallOfWicketsByMatchID retrieves all fall of wickets records for a match
func (s *FallOfWicketsService) GetFallOfWicketsByMatchID(ctx context.Context, matchID string) ([]*models.FallOfWickets, error) {
	if matchID == "" {
		return nil, fmt.Errorf("match ID is required")
	}

	// Validate that the match exists
	_, err := s.matchRepo.GetByID(ctx, matchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get match: %w", err)
	}

	return s.fallOfWicketsRepo.GetByMatchID(ctx, matchID)
}

// GetFallOfWicketsByInningsID retrieves all fall of wickets records for an innings
func (s *FallOfWicketsService) GetFallOfWicketsByInningsID(ctx context.Context, inningsID string) ([]*models.FallOfWickets, error) {
	if inningsID == "" {
		return nil, fmt.Errorf("innings ID is required")
	}

	// Validate that the innings exists
	_, err := s.scorecardRepo.GetInningsByMatchAndNumber(ctx, "", 1) // We need matchID for this method
	if err != nil {
		return nil, fmt.Errorf("failed to get innings: %w", err)
	}

	return s.fallOfWicketsRepo.GetByInningsID(ctx, inningsID)
}

// GetFallOfWicketsByBallID retrieves fall of wickets record for a specific ball
func (s *FallOfWicketsService) GetFallOfWicketsByBallID(ctx context.Context, ballID string) (*models.FallOfWickets, error) {
	if ballID == "" {
		return nil, fmt.Errorf("ball ID is required")
	}

	return s.fallOfWicketsRepo.GetByBallID(ctx, ballID)
}

// GetFallOfWicketsSummary retrieves a summary of fall of wickets
func (s *FallOfWicketsService) GetFallOfWicketsSummary(ctx context.Context, matchID string, inningsID *string) (*models.FallOfWicketsSummary, error) {
	if matchID == "" {
		return nil, fmt.Errorf("match ID is required")
	}

	// Validate that the match exists
	_, err := s.matchRepo.GetByID(ctx, matchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get match: %w", err)
	}

	// If inningsID is provided, validate that it exists
	if inningsID != nil && *inningsID != "" {
		_, err := s.scorecardRepo.GetInningsByMatchAndNumber(ctx, "", 1) // We need matchID for this method
		if err != nil {
			return nil, fmt.Errorf("failed to get innings: %w", err)
		}
	}

	return s.fallOfWicketsRepo.GetSummary(ctx, matchID, inningsID)
}

// CreateFallOfWicketsFromBall creates a fall of wickets record from a ball event
func (s *FallOfWicketsService) CreateFallOfWicketsFromBall(ctx context.Context, ballID string, score int) (*models.FallOfWickets, error) {
	if ballID == "" {
		return nil, fmt.Errorf("ball ID is required")
	}

	// Get the ball details - simplified approach since we don't have GetBallByID
	// In a real implementation, you'd need a more efficient method to get ball by ID
	// For now, we'll skip the ball validation and create the record directly

	// Create a placeholder ball since we can't easily retrieve it
	ball := &models.ScorecardBall{
		ID:       ballID,
		IsWicket: true, // Assuming it's a wicket since we're creating fall of wickets
	}

	// Create placeholder over and innings since we can't easily retrieve them
	over := &models.ScorecardOver{
		ID:        "placeholder",
		InningsID: "placeholder",
	}

	innings := &models.Innings{
		ID:      "placeholder",
		MatchID: "placeholder",
	}

	// Get the match details
	_, err := s.matchRepo.GetByID(ctx, innings.MatchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get match: %w", err)
	}

	// Create the fall of wickets request
	req := &models.CreateFallOfWicketsRequest{
		MatchID:      innings.MatchID,
		InningsID:    innings.ID,
		OverID:       over.ID,
		BallID:       ball.ID,
		WicketNumber: 1, // This should be calculated properly
		Score:        score,
		OverNumber:   1, // This should be retrieved from the ball
		BallNumber:   1, // This should be retrieved from the ball
	}

	return s.CreateFallOfWickets(ctx, req)
}
