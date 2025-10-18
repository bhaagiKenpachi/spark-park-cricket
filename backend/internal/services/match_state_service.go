package services

import (
	"context"
	"fmt"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/repository/interfaces"
	"time"
)

// MatchStateService handles match state transitions and management
type MatchStateService struct {
	matchRepo      interfaces.MatchRepository
	scoreboardRepo interfaces.ScoreboardRepository
	overRepo       interfaces.OverRepository
	ballRepo       interfaces.BallRepository
}

// NewMatchStateService creates a new match state service
func NewMatchStateService(
	matchRepo interfaces.MatchRepository,
	scoreboardRepo interfaces.ScoreboardRepository,
	overRepo interfaces.OverRepository,
	ballRepo interfaces.BallRepository,
) *MatchStateService {
	return &MatchStateService{
		matchRepo:      matchRepo,
		scoreboardRepo: scoreboardRepo,
		overRepo:       overRepo,
		ballRepo:       ballRepo,
	}
}

// StartMatch starts a match (changes status from not_started to live)
func (s *MatchStateService) StartMatch(ctx context.Context, matchID string) (*models.Match, error) {
	if matchID == "" {
		return nil, fmt.Errorf("match ID is required")
	}

	match, err := s.matchRepo.GetByID(ctx, matchID)
	if err != nil {
		return nil, fmt.Errorf("match not found: %w", err)
	}

	// Check if match is in not_started status
	if match.Status != models.MatchStatusNotStarted {
		return nil, fmt.Errorf("match must be not started to begin")
	}

	// Change status to live
	match.Status = models.MatchStatusLive
	match.UpdatedAt = time.Now()

	err = s.matchRepo.Update(ctx, matchID, match)
	if err != nil {
		return nil, fmt.Errorf("failed to start match: %w", err)
	}

	return match, nil
}

// CompleteMatch completes a match (changes status from live to completed)
func (s *MatchStateService) CompleteMatch(ctx context.Context, matchID string) (*models.Match, error) {
	if matchID == "" {
		return nil, fmt.Errorf("match ID is required")
	}

	match, err := s.matchRepo.GetByID(ctx, matchID)
	if err != nil {
		return nil, fmt.Errorf("match not found: %w", err)
	}

	if match.Status != models.MatchStatusLive {
		return nil, fmt.Errorf("match must be live to complete")
	}

	now := time.Now()
	match.Status = models.MatchStatusCompleted
	match.EndTime = &now
	match.UpdatedAt = now

	err = s.matchRepo.Update(ctx, matchID, match)
	if err != nil {
		return nil, fmt.Errorf("failed to complete match: %w", err)
	}

	return match, nil
}

// CancelMatch cancels a match
func (s *MatchStateService) CancelMatch(ctx context.Context, matchID string) (*models.Match, error) {
	if matchID == "" {
		return nil, fmt.Errorf("match ID is required")
	}

	match, err := s.matchRepo.GetByID(ctx, matchID)
	if err != nil {
		return nil, fmt.Errorf("match not found: %w", err)
	}

	if match.Status == models.MatchStatusCompleted {
		return nil, fmt.Errorf("cannot cancel a completed match")
	}

	match.Status = models.MatchStatusCancelled
	match.UpdatedAt = time.Now()

	err = s.matchRepo.Update(ctx, matchID, match)
	if err != nil {
		return nil, fmt.Errorf("failed to cancel match: %w", err)
	}

	return match, nil
}

// GetMatchSummary gets a comprehensive match summary
func (s *MatchStateService) GetMatchSummary(ctx context.Context, matchID string) (*MatchSummary, error) {
	if matchID == "" {
		return nil, fmt.Errorf("match ID is required")
	}

	match, err := s.matchRepo.GetByID(ctx, matchID)
	if err != nil {
		return nil, fmt.Errorf("match not found: %w", err)
	}

	// Get scoreboard
	scoreboard, err := s.scoreboardRepo.GetByMatchID(ctx, matchID)
	if err != nil {
		// If no scoreboard exists, create a basic summary
		return &MatchSummary{
			Match:      match,
			Scoreboard: nil,
			Overs:      []*models.Over{},
			TotalBalls: 0,
		}, nil
	}

	// Get overs
	overs, err := s.overRepo.GetByMatchID(ctx, matchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get overs: %w", err)
	}

	// Calculate total balls
	totalBalls := 0
	for _, over := range overs {
		totalBalls += over.TotalBalls
	}

	return &MatchSummary{
		Match:      match,
		Scoreboard: scoreboard,
		Overs:      overs,
		TotalBalls: totalBalls,
	}, nil
}

// CheckMatchCompletion checks if a match should be completed based on cricket rules
func (s *MatchStateService) CheckMatchCompletion(ctx context.Context, matchID string) (bool, error) {
	if matchID == "" {
		return false, fmt.Errorf("match ID is required")
	}

	match, err := s.matchRepo.GetByID(ctx, matchID)
	if err != nil {
		return false, fmt.Errorf("match not found: %w", err)
	}

	if match.Status != models.MatchStatusLive {
		return false, nil
	}

	scoreboard, err := s.scoreboardRepo.GetByMatchID(ctx, matchID)
	if err != nil {
		return false, nil // No scoreboard means match not started
	}

	// Check if all wickets are down - use dynamic calculation based on team size
	maxWickets := match.TeamAPlayerCount - 1 // n-1 wickets for n players
	if scoreboard.Wickets >= maxWickets {
		return true, nil
	}

	// Check if target overs are completed (assuming 20 overs for T20)
	// This would need to be configurable based on match format
	if scoreboard.Overs >= 20.0 {
		return true, nil
	}

	return false, nil
}

// MatchSummary represents a comprehensive match summary
type MatchSummary struct {
	Match      *models.Match          `json:"match"`
	Scoreboard *models.LiveScoreboard `json:"scoreboard"`
	Overs      []*models.Over         `json:"overs"`
	TotalBalls int                    `json:"total_balls"`
}
