package graphql

import (
	"fmt"
	"spark-park-cricket-backend/internal/interfaces"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/pkg/websocket"

	"github.com/graphql-go/graphql"
)

// ResolverContext holds the services needed for GraphQL resolvers
type ResolverContext struct {
	ScorecardService interfaces.ScorecardServiceInterface
	Hub              *websocket.Hub
}

// resolveLiveScorecard resolves the live scorecard query
func resolveLiveScorecard(p graphql.ResolveParams) (interface{}, error) {
	matchID, ok := p.Args["match_id"].(string)
	if !ok {
		return nil, fmt.Errorf("match_id is required")
	}

	// Get resolver context from the context
	resolverCtx, ok := p.Context.Value(resolverContextKey).(*ResolverContext)
	if !ok {
		return nil, fmt.Errorf("resolver context not found")
	}

	// Get the full scorecard
	scorecard, err := resolverCtx.ScorecardService.GetScorecard(p.Context, matchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get scorecard: %w", err)
	}

	// Calculate current score for live updates
	currentScore := calculateCurrentScore(scorecard)

	// Get current over
	scorecardOver, err := resolverCtx.ScorecardService.GetCurrentOver(p.Context, matchID, scorecard.CurrentInnings)
	var currentOver *models.OverSummary
	if err != nil {
		// If no current over, return nil
		currentOver = nil
	} else {
		// Convert ScorecardOver to OverSummary
		currentOver = &models.OverSummary{
			OverNumber:   scorecardOver.OverNumber,
			TotalRuns:    scorecardOver.TotalRuns,
			TotalBalls:   scorecardOver.TotalBalls,
			TotalWickets: scorecardOver.TotalWickets,
			Status:       scorecardOver.Status,
			Balls:        []models.BallSummary{}, // Will be populated separately if needed
		}
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

	return liveScorecard, nil
}

// resolveScorecardSubscription resolves the scorecard subscription for real-time updates
func resolveScorecardSubscription(p graphql.ResolveParams) (interface{}, error) {
	matchID, ok := p.Args["match_id"].(string)
	if !ok {
		return nil, fmt.Errorf("match_id is required")
	}

	// Get resolver context from the context
	resolverCtx, ok := p.Context.Value(resolverContextKey).(*ResolverContext)
	if !ok {
		return nil, fmt.Errorf("resolver context not found")
	}

	// For now, return the current scorecard
	// In a real implementation, this would set up a subscription to WebSocket updates
	scorecard, err := resolverCtx.ScorecardService.GetScorecard(p.Context, matchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get scorecard: %w", err)
	}

	// Calculate current score for live updates
	currentScore := calculateCurrentScore(scorecard)

	// Get current over
	scorecardOver, err := resolverCtx.ScorecardService.GetCurrentOver(p.Context, matchID, scorecard.CurrentInnings)
	var currentOver *models.OverSummary
	if err != nil {
		currentOver = nil
	} else {
		// Convert ScorecardOver to OverSummary
		currentOver = &models.OverSummary{
			OverNumber:   scorecardOver.OverNumber,
			TotalRuns:    scorecardOver.TotalRuns,
			TotalBalls:   scorecardOver.TotalBalls,
			TotalWickets: scorecardOver.TotalWickets,
			Status:       scorecardOver.Status,
			Balls:        []models.BallSummary{}, // Will be populated separately if needed
		}
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

	return liveScorecard, nil
}

// calculateCurrentScore calculates the current score from the scorecard
func calculateCurrentScore(scorecard *models.ScorecardResponse) map[string]interface{} {
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
