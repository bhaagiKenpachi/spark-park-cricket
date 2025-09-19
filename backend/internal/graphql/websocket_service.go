package graphql

import (
	"context"
	"log"
	"spark-park-cricket-backend/internal/interfaces"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/pkg/websocket"
)

// GraphQLWebSocketService handles GraphQL queries with WebSocket integration
type GraphQLWebSocketService struct {
	scorecardService interfaces.ScorecardServiceInterface
	hub              *websocket.Hub
	graphqlHandler   *GraphQLHandler
}

// NewGraphQLWebSocketService creates a new GraphQL WebSocket service
func NewGraphQLWebSocketService(scorecardService interfaces.ScorecardServiceInterface, hub *websocket.Hub) *GraphQLWebSocketService {
	graphqlHandler := NewGraphQLHandler(scorecardService, hub)

	return &GraphQLWebSocketService{
		scorecardService: scorecardService,
		hub:              hub,
		graphqlHandler:   graphqlHandler,
	}
}

// BroadcastScorecardUpdate broadcasts a scorecard update to WebSocket clients
func (s *GraphQLWebSocketService) BroadcastScorecardUpdate(matchID string) {
	// Get the current scorecard
	scorecard, err := s.scorecardService.GetScorecard(context.Background(), matchID)
	if err != nil {
		log.Printf("Error getting scorecard for broadcast: %v", err)
		return
	}

	// Calculate current score for live updates
	currentScore := s.calculateCurrentScore(scorecard)

	// Get current over
	currentOver, err := s.scorecardService.GetCurrentOver(context.Background(), matchID, scorecard.CurrentInnings)
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

// BroadcastBallUpdate broadcasts a ball update to WebSocket clients
func (s *GraphQLWebSocketService) BroadcastBallUpdate(matchID string, ballEvent *models.BallEventRequest) {
	// Get the current scorecard
	scorecard, err := s.scorecardService.GetScorecard(context.Background(), matchID)
	if err != nil {
		log.Printf("Error getting scorecard for ball broadcast: %v", err)
		return
	}

	// Calculate current score for live updates
	currentScore := s.calculateCurrentScore(scorecard)

	// Get current over
	currentOver, err := s.scorecardService.GetCurrentOver(context.Background(), matchID, scorecard.CurrentInnings)
	if err != nil {
		currentOver = nil
	}

	// Build the ball update response
	ballUpdate := map[string]interface{}{
		"match_id":      scorecard.MatchID,
		"current_score": currentScore,
		"current_over":  currentOver,
		"ball_event": map[string]interface{}{
			"ball_type":      ballEvent.BallType,
			"run_type":       ballEvent.RunType,
			"runs":           ballEvent.RunType.GetRunValue(),
			"byes":           ballEvent.Byes,
			"is_wicket":      ballEvent.IsWicket,
			"wicket_type":    ballEvent.WicketType,
			"innings_number": ballEvent.InningsNumber,
		},
	}

	// Create update message
	updateMessage := map[string]interface{}{
		"type": "ball_update",
		"data": ballUpdate,
	}

	// Broadcast to the match room
	s.hub.BroadcastToRoom(matchID, updateMessage)
	log.Printf("Broadcasted ball update for match %s", matchID)
}

// calculateCurrentScore calculates the current score from the scorecard
func (s *GraphQLWebSocketService) calculateCurrentScore(scorecard *models.ScorecardResponse) map[string]interface{} {
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

// GetGraphQLHandler returns the GraphQL handler
func (s *GraphQLWebSocketService) GetGraphQLHandler() *GraphQLHandler {
	return s.graphqlHandler
}
