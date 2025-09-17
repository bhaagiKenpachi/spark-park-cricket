package services

import (
	"context"
	"log"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/repository/interfaces"
	"spark-park-cricket-backend/pkg/websocket"
)

// ScorecardServiceWithGraphQL wraps the scorecard service with GraphQL WebSocket integration
type ScorecardServiceWithGraphQL struct {
	*ScorecardService
	hub *websocket.Hub
}

// NewScorecardServiceWithGraphQL creates a new scorecard service with GraphQL integration
func NewScorecardServiceWithGraphQL(scorecardRepo interfaces.ScorecardRepository, matchRepo interfaces.MatchRepository, hub *websocket.Hub) *ScorecardServiceWithGraphQL {
	baseService := NewScorecardService(scorecardRepo, matchRepo)

	return &ScorecardServiceWithGraphQL{
		ScorecardService: baseService,
		hub:              hub,
	}
}

// AddBall adds a ball and broadcasts the update via WebSocket
func (s *ScorecardServiceWithGraphQL) AddBall(ctx context.Context, req *models.BallEventRequest) error {
	// Call the base service to add the ball
	err := s.ScorecardService.AddBall(ctx, req)
	if err != nil {
		return err
	}

	// Broadcast the update via WebSocket
	s.broadcastScorecardUpdate(req.MatchID)

	log.Printf("Ball added and broadcasted for match %s", req.MatchID)
	return nil
}

// UndoBall undoes a ball and broadcasts the update via WebSocket
func (s *ScorecardServiceWithGraphQL) UndoBall(ctx context.Context, matchID string, inningsNumber int) error {
	// Call the base service to undo the ball
	err := s.ScorecardService.UndoBall(ctx, matchID, inningsNumber)
	if err != nil {
		return err
	}

	// Broadcast the update via WebSocket
	s.broadcastScorecardUpdate(matchID)

	log.Printf("Ball undone and broadcasted for match %s", matchID)
	return nil
}

// StartScoring starts scoring and broadcasts the update via WebSocket
func (s *ScorecardServiceWithGraphQL) StartScoring(ctx context.Context, matchID string) error {
	// Call the base service to start scoring
	err := s.ScorecardService.StartScoring(ctx, matchID)
	if err != nil {
		return err
	}

	// Broadcast the update via WebSocket
	s.broadcastScorecardUpdate(matchID)

	log.Printf("Scoring started and broadcasted for match %s", matchID)
	return nil
}

// broadcastScorecardUpdate broadcasts a scorecard update to WebSocket clients
func (s *ScorecardServiceWithGraphQL) broadcastScorecardUpdate(matchID string) {
	// Get the current scorecard
	scorecard, err := s.ScorecardService.GetScorecard(context.Background(), matchID)
	if err != nil {
		log.Printf("Error getting scorecard for broadcast: %v", err)
		return
	}

	// Calculate current score for live updates
	currentScore := s.calculateCurrentScore(scorecard)

	// Get current over
	currentOver, err := s.ScorecardService.GetCurrentOver(context.Background(), matchID, scorecard.CurrentInnings)
	if err != nil {
		currentOver = nil
	}

	// Build the live scorecard response
	liveScorecard := map[string]interface{}{
		"match_id":        scorecard.MatchID,
		"match_number":    scorecard.MatchNumber,
		"series_name":     scorecard.SeriesName,
		"team_a":          scorecard.TeamA,
		"team_b":          scorecard.TeamB,
		"total_overs":     scorecard.TotalOvers,
		"toss_winner":     scorecard.TossWinner,
		"toss_type":       scorecard.TossType,
		"current_innings": scorecard.CurrentInnings,
		"match_status":    scorecard.MatchStatus,
		"current_score":   currentScore,
		"current_over":    currentOver,
		"innings":         scorecard.Innings,
	}

	// Create update message
	updateMessage := map[string]interface{}{
		"type": "scorecard_update",
		"data": liveScorecard,
	}

	// Broadcast to the match room
	s.hub.BroadcastToRoom(matchID, updateMessage)
	log.Printf("Broadcasted scorecard update for match %s", matchID)
}

// calculateCurrentScore calculates the current score from the scorecard
func (s *ScorecardServiceWithGraphQL) calculateCurrentScore(scorecard *models.ScorecardResponse) map[string]interface{} {
	var currentInnings *models.InningsSummary

	// Find the current innings
	for i := range scorecard.Innings {
		if scorecard.Innings[i].InningsNumber == scorecard.CurrentInnings {
			currentInnings = &scorecard.Innings[i]
			break
		}
	}

	if currentInnings == nil {
		return map[string]interface{}{
			"runs":     0,
			"wickets":  0,
			"overs":    0.0,
			"balls":    0,
			"run_rate": 0.0,
		}
	}

	// Calculate run rate
	runRate := 0.0
	if currentInnings.TotalOvers > 0 {
		runRate = float64(currentInnings.TotalRuns) / currentInnings.TotalOvers
	}

	return map[string]interface{}{
		"runs":     currentInnings.TotalRuns,
		"wickets":  currentInnings.TotalWickets,
		"overs":    currentInnings.TotalOvers,
		"balls":    currentInnings.TotalBalls,
		"run_rate": runRate,
	}
}
