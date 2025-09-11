package services

import (
	"context"
	"fmt"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/repository/interfaces"
	"spark-park-cricket-backend/pkg/events"
	"time"
)

// RealtimeScoreboardService extends ScoreboardService with real-time WebSocket updates
type RealtimeScoreboardService struct {
	*ScoreboardService
	broadcaster *events.EventBroadcaster
}

// NewRealtimeScoreboardService creates a new real-time scoreboard service
func NewRealtimeScoreboardService(
	scoreboardRepo interfaces.ScoreboardRepository,
	overRepo interfaces.OverRepository,
	ballRepo interfaces.BallRepository,
	matchRepo interfaces.MatchRepository,
	teamRepo interfaces.TeamRepository,
	playerRepo interfaces.PlayerRepository,
	broadcaster *events.EventBroadcaster,
) *RealtimeScoreboardService {
	baseService := NewScoreboardService(
		scoreboardRepo, overRepo, ballRepo, matchRepo, teamRepo, playerRepo,
	)

	return &RealtimeScoreboardService{
		ScoreboardService: baseService,
		broadcaster:       broadcaster,
	}
}

// AddBall adds a ball event and broadcasts the update to WebSocket clients
func (s *RealtimeScoreboardService) AddBall(ctx context.Context, matchID string, ballEvent *models.BallEvent) (*models.LiveScoreboard, error) {
	// Use the base service to add the ball
	scoreboard, err := s.ScoreboardService.AddBall(ctx, matchID, ballEvent)
	if err != nil {
		return nil, err
	}

	// Broadcast the ball event to all connected clients
	s.broadcaster.BroadcastBallEvent(ctx, matchID, ballEvent, scoreboard)

	return scoreboard, nil
}

// UpdateScore updates the score and broadcasts the update to WebSocket clients
func (s *RealtimeScoreboardService) UpdateScore(ctx context.Context, matchID string, req *models.UpdateScoreRequest) (*models.LiveScoreboard, error) {
	// Use the base service to update the score
	scoreboard, err := s.ScoreboardService.UpdateScore(ctx, matchID, req)
	if err != nil {
		return nil, err
	}

	// Broadcast the score update to all connected clients
	s.broadcaster.BroadcastScoreUpdate(ctx, matchID, scoreboard)

	return scoreboard, nil
}

// UpdateWicket updates the wicket count and broadcasts the update to WebSocket clients
func (s *RealtimeScoreboardService) UpdateWicket(ctx context.Context, matchID string, req *models.UpdateWicketRequest) (*models.LiveScoreboard, error) {
	// Use the base service to update the wicket
	scoreboard, err := s.ScoreboardService.UpdateWicket(ctx, matchID, req)
	if err != nil {
		return nil, err
	}

	// Broadcast the wicket update to all connected clients
	s.broadcaster.BroadcastWicketUpdate(ctx, matchID, scoreboard)

	return scoreboard, nil
}

// StartMatch starts a match and broadcasts the start event
func (s *RealtimeScoreboardService) StartMatch(ctx context.Context, matchID string) (*models.Match, error) {
	// Get match state service
	matchStateService := NewMatchStateService(
		s.matchRepo, s.scoreboardRepo, s.overRepo, s.ballRepo,
	)

	// Start the match
	match, err := matchStateService.StartMatch(ctx, matchID)
	if err != nil {
		return nil, err
	}

	// Broadcast match start event
	s.broadcaster.BroadcastMatchStart(ctx, matchID, match)

	return match, nil
}

// CompleteMatch completes a match and broadcasts the end event
func (s *RealtimeScoreboardService) CompleteMatch(ctx context.Context, matchID string) (*models.Match, error) {
	// Get match state service
	matchStateService := NewMatchStateService(
		s.matchRepo, s.scoreboardRepo, s.overRepo, s.ballRepo,
	)

	// Complete the match
	match, err := matchStateService.CompleteMatch(ctx, matchID)
	if err != nil {
		return nil, err
	}

	// Get final scoreboard
	scoreboard, err := s.GetScoreboard(ctx, matchID)
	if err != nil {
		// If no scoreboard exists, create a default one
		scoreboard = &models.LiveScoreboard{
			MatchID: matchID,
			Score:   0,
			Wickets: 0,
			Overs:   0.0,
			Balls:   0,
		}
	}

	// Broadcast match end event
	s.broadcaster.BroadcastMatchEnd(ctx, matchID, match, scoreboard)

	return match, nil
}

// AddBallWithOverCompletion adds a ball and handles over completion broadcasting
func (s *RealtimeScoreboardService) AddBallWithOverCompletion(ctx context.Context, matchID string, ballEvent *models.BallEvent) (*models.LiveScoreboard, error) {
	// Get current scoreboard
	scoreboard, err := s.GetScoreboard(ctx, matchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get scoreboard: %w", err)
	}

	// Get current over
	over, err := s.getCurrentOver(ctx, matchID, string(scoreboard.BattingTeam))
	if err != nil {
		return nil, fmt.Errorf("failed to get current over: %w", err)
	}

	// Add the ball
	scoreboard, err = s.AddBall(ctx, matchID, ballEvent)
	if err != nil {
		return nil, err
	}

	// Check if over is complete and broadcast if so
	if s.isOverComplete(over) {
		// Get updated over data
		updatedOver, err := s.overRepo.GetByID(ctx, over.ID)
		if err == nil {
			s.broadcaster.BroadcastOverCompletion(ctx, matchID, updatedOver, scoreboard)
		}
	}

	return scoreboard, nil
}

// BroadcastCustomUpdate broadcasts a custom update to all clients watching the match
func (s *RealtimeScoreboardService) BroadcastCustomUpdate(ctx context.Context, matchID string, messageType string, data interface{}) {
	s.broadcaster.BroadcastCustomMessage(ctx, matchID, messageType, data)
}

// GetLiveStats returns live statistics for a match
func (s *RealtimeScoreboardService) GetLiveStats(ctx context.Context, matchID string) (*LiveMatchStats, error) {
	// Get match
	match, err := s.matchRepo.GetByID(ctx, matchID)
	if err != nil {
		return nil, fmt.Errorf("match not found: %w", err)
	}

	// Get scoreboard
	scoreboard, err := s.GetScoreboard(ctx, matchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get scoreboard: %w", err)
	}

	// Get overs
	overs, err := s.overRepo.GetByMatchID(ctx, matchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get overs: %w", err)
	}

	// Calculate statistics
	totalBalls := 0
	totalRuns := 0
	for _, over := range overs {
		totalBalls += over.TotalBalls
		totalRuns += over.TotalRuns
	}

	// Get recent balls (last 10)
	recentBalls := []*models.Ball{}
	if len(overs) > 0 {
		lastOver := overs[len(overs)-1]
		balls, err := s.ballRepo.GetByOverID(ctx, lastOver.ID)
		if err == nil && len(balls) > 0 {
			// Get last 10 balls
			start := 0
			if len(balls) > 10 {
				start = len(balls) - 10
			}
			recentBalls = balls[start:]
		}
	}

	return &LiveMatchStats{
		Match:       match,
		Scoreboard:  scoreboard,
		TotalOvers:  len(overs),
		TotalBalls:  totalBalls,
		TotalRuns:   totalRuns,
		RecentBalls: recentBalls,
		LastUpdated: time.Now(),
	}, nil
}

// LiveMatchStats represents live match statistics
type LiveMatchStats struct {
	Match       *models.Match          `json:"match"`
	Scoreboard  *models.LiveScoreboard `json:"scoreboard"`
	TotalOvers  int                    `json:"total_overs"`
	TotalBalls  int                    `json:"total_balls"`
	TotalRuns   int                    `json:"total_runs"`
	RecentBalls []*models.Ball         `json:"recent_balls"`
	LastUpdated time.Time              `json:"last_updated"`
}
