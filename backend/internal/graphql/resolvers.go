package graphql

import (
	"fmt"
	"log"
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

// resolveInningsScore resolves the innings score query
func resolveInningsScore(p graphql.ResolveParams) (interface{}, error) {
	matchID, ok := p.Args["match_id"].(string)
	if !ok {
		return nil, fmt.Errorf("match_id is required")
	}

	inningsNumber, ok := p.Args["innings_number"].(int)
	if !ok {
		return nil, fmt.Errorf("innings_number is required")
	}

	// Get resolver context from the context
	resolverCtx, ok := p.Context.Value(resolverContextKey).(*ResolverContext)
	if !ok {
		return nil, fmt.Errorf("resolver context not found")
	}

	// Get the scorecard to find the specific innings
	scorecard, err := resolverCtx.ScorecardService.GetScorecard(p.Context, matchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get scorecard: %w", err)
	}

	// Find the specific innings
	for _, innings := range scorecard.Innings {
		if innings.InningsNumber == inningsNumber {
			return map[string]interface{}{
				"innings_number": innings.InningsNumber,
				"batting_team":   innings.BattingTeam,
				"total_runs":     innings.TotalRuns,
				"total_wickets":  innings.TotalWickets,
				"total_overs":    innings.TotalOvers,
				"total_balls":    innings.TotalBalls,
				"status":         innings.Status,
				"extras":         innings.Extras,
			}, nil
		}
	}

	return nil, fmt.Errorf("innings %d not found", inningsNumber)
}

// resolveLatestOver resolves the latest over query
func resolveLatestOver(p graphql.ResolveParams) (interface{}, error) {
	matchID, ok := p.Args["match_id"].(string)
	if !ok {
		return nil, fmt.Errorf("match_id is required")
	}

	inningsNumber, ok := p.Args["innings_number"].(int)
	if !ok {
		return nil, fmt.Errorf("innings_number is required")
	}

	// Get resolver context from the context
	resolverCtx, ok := p.Context.Value(resolverContextKey).(*ResolverContext)
	if !ok {
		return nil, fmt.Errorf("resolver context not found")
	}

	// Get the current over
	scorecardOver, err := resolverCtx.ScorecardService.GetCurrentOver(p.Context, matchID, inningsNumber)
	if err != nil {
		// If no current over exists (e.g., scoring just started, no balls added yet),
		// return a default over structure
		log.Printf("No current over found for match %s, innings %d: %v", matchID, inningsNumber, err)

		// Return a default over structure indicating no over exists yet
		overSummary := map[string]interface{}{
			"over_number":   1, // Default to over 1
			"total_runs":    0,
			"total_balls":   0,
			"total_wickets": 0,
			"status":        "not_started",          // Indicate over hasn't started yet
			"balls":         []models.BallSummary{}, // Empty balls array
		}

		return overSummary, nil
	}

	// Get balls for this over
	balls, err := resolverCtx.ScorecardService.GetBallsByOver(p.Context, scorecardOver.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get balls for over: %w", err)
	}

	// Convert ScorecardBall to BallSummary
	ballSummaries := make([]models.BallSummary, len(balls))
	for i, ball := range balls {
		ballSummaries[i] = models.BallSummary{
			BallNumber: ball.BallNumber,
			BallType:   ball.BallType,
			RunType:    ball.RunType,
			Runs:       ball.Runs,
			Byes:       ball.Byes,
			IsWicket:   ball.IsWicket,
			WicketType: ball.WicketType,
		}
	}

	// Convert ScorecardOver to OverSummary
	overSummary := map[string]interface{}{
		"over_number":   scorecardOver.OverNumber,
		"total_runs":    scorecardOver.TotalRuns,
		"total_balls":   scorecardOver.TotalBalls,
		"total_wickets": scorecardOver.TotalWickets,
		"status":        scorecardOver.Status,
		"balls":         ballSummaries,
	}

	return overSummary, nil
}

// resolveAllOvers resolves the all overs query
func resolveAllOvers(p graphql.ResolveParams) (interface{}, error) {
	matchID, ok := p.Args["match_id"].(string)
	if !ok {
		return nil, fmt.Errorf("match_id is required")
	}

	inningsNumber, ok := p.Args["innings_number"].(int)
	if !ok {
		return nil, fmt.Errorf("innings_number is required")
	}

	// Get resolver context from the context
	resolverCtx, ok := p.Context.Value(resolverContextKey).(*ResolverContext)
	if !ok {
		return nil, fmt.Errorf("resolver context not found")
	}

	// Get the scorecard to find the specific innings
	scorecard, err := resolverCtx.ScorecardService.GetScorecard(p.Context, matchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get scorecard: %w", err)
	}

	// Find the specific innings and return its overs
	for _, innings := range scorecard.Innings {
		if innings.InningsNumber == inningsNumber {
			// Convert overs to the expected format
			var overs []map[string]interface{}
			for _, over := range innings.Overs {
				overMap := map[string]interface{}{
					"over_number":   over.OverNumber,
					"total_runs":    over.TotalRuns,
					"total_balls":   over.TotalBalls,
					"total_wickets": over.TotalWickets,
					"status":        over.Status,
					"balls":         over.Balls,
				}
				overs = append(overs, overMap)
			}
			return overs, nil
		}
	}

	return nil, fmt.Errorf("innings %d not found", inningsNumber)
}

// resolveMatchDetails resolves the match details query
func resolveMatchDetails(p graphql.ResolveParams) (interface{}, error) {
	log.Printf("DEBUG: resolveMatchDetails called with args: %+v", p.Args)
	matchID, ok := p.Args["match_id"].(string)
	if !ok {
		log.Printf("DEBUG: match_id not found in args")
		return nil, fmt.Errorf("match_id is required")
	}
	log.Printf("DEBUG: match_id: %s", matchID)

	// Get resolver context from the context
	resolverCtx, ok := p.Context.Value(resolverContextKey).(*ResolverContext)
	if !ok {
		return nil, fmt.Errorf("resolver context not found")
	}

	// Get the scorecard to extract match details
	scorecard, err := resolverCtx.ScorecardService.GetScorecard(p.Context, matchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get scorecard: %w", err)
	}

	// Build match details response
	matchDetails := map[string]interface{}{
		"match_id":            scorecard.MatchID,
		"match_number":        scorecard.MatchNumber,
		"series_name":         scorecard.SeriesName,
		"team_a":              scorecard.TeamA,
		"team_b":              scorecard.TeamB,
		"total_overs":         scorecard.TotalOvers,
		"toss_winner":         scorecard.TossWinner,
		"toss_type":           scorecard.TossType,
		"current_innings":     scorecard.CurrentInnings,
		"match_status":        scorecard.MatchStatus,
		"batting_team":        nil, // This would need to be fetched from match service
		"team_a_player_count": nil, // This would need to be fetched from match service
		"team_b_player_count": nil, // This would need to be fetched from match service
	}

	return matchDetails, nil
}

// resolveMatchStatistics resolves the match statistics query
func resolveMatchStatistics(p graphql.ResolveParams) (interface{}, error) {
	matchID, ok := p.Args["match_id"].(string)
	if !ok {
		return nil, fmt.Errorf("match_id is required")
	}

	// Get resolver context from the context
	resolverCtx, ok := p.Context.Value(resolverContextKey).(*ResolverContext)
	if !ok {
		return nil, fmt.Errorf("resolver context not found")
	}

	// Get the scorecard to calculate statistics
	scorecard, err := resolverCtx.ScorecardService.GetScorecard(p.Context, matchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get scorecard: %w", err)
	}

	// Calculate match-level statistics
	totalRuns := 0
	totalWickets := 0
	totalOvers := 0.0
	totalBalls := 0
	totalExtras := models.ExtrasSummary{}

	for _, innings := range scorecard.Innings {
		totalRuns += innings.TotalRuns
		totalWickets += innings.TotalWickets
		totalOvers += innings.TotalOvers
		totalBalls += innings.TotalBalls

		if innings.Extras != nil {
			totalExtras.Byes += innings.Extras.Byes
			totalExtras.LegByes += innings.Extras.LegByes
			totalExtras.Wides += innings.Extras.Wides
			totalExtras.NoBalls += innings.Extras.NoBalls
			totalExtras.Total += innings.Extras.Total
		}
	}

	// Calculate run rate
	runRate := 0.0
	if totalOvers > 0 {
		runRate = float64(totalRuns) / totalOvers
	}

	// Build match statistics response
	matchStatistics := map[string]interface{}{
		"total_runs":    totalRuns,
		"total_wickets": totalWickets,
		"total_overs":   totalOvers,
		"total_balls":   totalBalls,
		"run_rate":      runRate,
		"extras":        totalExtras,
		"innings_count": len(scorecard.Innings),
	}

	return matchStatistics, nil
}

// resolveInningsDetails resolves the innings details query
func resolveInningsDetails(p graphql.ResolveParams) (interface{}, error) {
	matchID, ok := p.Args["match_id"].(string)
	if !ok {
		return nil, fmt.Errorf("match_id is required")
	}

	inningsNumber, ok := p.Args["innings_number"].(int)
	if !ok {
		return nil, fmt.Errorf("innings_number is required")
	}

	// Get resolver context from the context
	resolverCtx, ok := p.Context.Value(resolverContextKey).(*ResolverContext)
	if !ok {
		return nil, fmt.Errorf("resolver context not found")
	}

	// Get the scorecard to find the specific innings
	scorecard, err := resolverCtx.ScorecardService.GetScorecard(p.Context, matchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get scorecard: %w", err)
	}

	// Find the specific innings
	for _, innings := range scorecard.Innings {
		if innings.InningsNumber == inningsNumber {
			return map[string]interface{}{
				"innings_number": innings.InningsNumber,
				"batting_team":   innings.BattingTeam,
				"total_runs":     innings.TotalRuns,
				"total_wickets":  innings.TotalWickets,
				"total_overs":    innings.TotalOvers,
				"total_balls":    innings.TotalBalls,
				"status":         innings.Status,
				"extras":         innings.Extras,
				"overs":          innings.Overs,
			}, nil
		}
	}

	return nil, fmt.Errorf("innings %d not found", inningsNumber)
}

// resolveOverDetails resolves the over details query
func resolveOverDetails(p graphql.ResolveParams) (interface{}, error) {
	matchID, ok := p.Args["match_id"].(string)
	if !ok {
		return nil, fmt.Errorf("match_id is required")
	}

	inningsNumber, ok := p.Args["innings_number"].(int)
	if !ok {
		return nil, fmt.Errorf("innings_number is required")
	}

	overNumber, ok := p.Args["over_number"].(int)
	if !ok {
		return nil, fmt.Errorf("over_number is required")
	}

	// Get resolver context from the context
	resolverCtx, ok := p.Context.Value(resolverContextKey).(*ResolverContext)
	if !ok {
		return nil, fmt.Errorf("resolver context not found")
	}

	// Get the scorecard to find the specific over
	scorecard, err := resolverCtx.ScorecardService.GetScorecard(p.Context, matchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get scorecard: %w", err)
	}

	// Find the specific over
	for _, innings := range scorecard.Innings {
		if innings.InningsNumber == inningsNumber {
			for _, over := range innings.Overs {
				if over.OverNumber == overNumber {
					return map[string]interface{}{
						"over_number":   over.OverNumber,
						"total_runs":    over.TotalRuns,
						"total_balls":   over.TotalBalls,
						"total_wickets": over.TotalWickets,
						"status":        over.Status,
						"balls":         over.Balls,
					}, nil
				}
			}
			break
		}
	}

	return nil, fmt.Errorf("over %d not found in innings %d", overNumber, inningsNumber)
}

// resolveBallDetails resolves the ball details query
func resolveBallDetails(p graphql.ResolveParams) (interface{}, error) {
	matchID, ok := p.Args["match_id"].(string)
	if !ok {
		return nil, fmt.Errorf("match_id is required")
	}

	inningsNumber, ok := p.Args["innings_number"].(int)
	if !ok {
		return nil, fmt.Errorf("innings_number is required")
	}

	overNumber, ok := p.Args["over_number"].(int)
	if !ok {
		return nil, fmt.Errorf("over_number is required")
	}

	// Get resolver context from the context
	resolverCtx, ok := p.Context.Value(resolverContextKey).(*ResolverContext)
	if !ok {
		return nil, fmt.Errorf("resolver context not found")
	}

	// Get the scorecard to find the specific balls
	scorecard, err := resolverCtx.ScorecardService.GetScorecard(p.Context, matchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get scorecard: %w", err)
	}

	// Find the specific over and return its balls
	for _, innings := range scorecard.Innings {
		if innings.InningsNumber == inningsNumber {
			for _, over := range innings.Overs {
				if over.OverNumber == overNumber {
					// Convert balls to the expected format
					var balls []map[string]interface{}
					for _, ball := range over.Balls {
						ballMap := map[string]interface{}{
							"ball_number": ball.BallNumber,
							"ball_type":   ball.BallType,
							"run_type":    ball.RunType,
							"runs":        ball.Runs,
							"byes":        ball.Byes,
							"is_wicket":   ball.IsWicket,
							"wicket_type": ball.WicketType,
						}
						balls = append(balls, ballMap)
					}
					return balls, nil
				}
			}
			break
		}
	}

	return nil, fmt.Errorf("over %d not found in innings %d", overNumber, inningsNumber)
}

// resolveMatchTeams resolves the match teams query
func resolveMatchTeams(p graphql.ResolveParams) (interface{}, error) {
	matchID, ok := p.Args["match_id"].(string)
	if !ok {
		return nil, fmt.Errorf("match_id is required")
	}

	// Get resolver context from the context
	resolverCtx, ok := p.Context.Value(resolverContextKey).(*ResolverContext)
	if !ok {
		return nil, fmt.Errorf("resolver context not found")
	}

	// Get the scorecard to extract team information
	scorecard, err := resolverCtx.ScorecardService.GetScorecard(p.Context, matchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get scorecard: %w", err)
	}

	// For now, return basic team information from scorecard
	// In a real implementation, this would fetch full team details from team service
	teams := []map[string]interface{}{
		{
			"id":            nil, // Would need team service to get actual team IDs
			"name":          scorecard.TeamA,
			"players_count": nil, // Would need team service to get actual player count
			"created_at":    nil,
			"updated_at":    nil,
		},
		{
			"id":            nil, // Would need team service to get actual team IDs
			"name":          scorecard.TeamB,
			"players_count": nil, // Would need team service to get actual player count
			"created_at":    nil,
			"updated_at":    nil,
		},
	}

	return teams, nil
}

// resolveMatchPlayers resolves the match players query
func resolveMatchPlayers(p graphql.ResolveParams) (interface{}, error) {
	matchID, ok := p.Args["match_id"].(string)
	if !ok {
		return nil, fmt.Errorf("match_id is required")
	}

	// Get resolver context from the context
	resolverCtx, ok := p.Context.Value(resolverContextKey).(*ResolverContext)
	if !ok {
		return nil, fmt.Errorf("resolver context not found")
	}

	// Get the scorecard to extract team information
	_, err := resolverCtx.ScorecardService.GetScorecard(p.Context, matchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get scorecard: %w", err)
	}

	// For now, return empty players list
	// In a real implementation, this would fetch players from player service based on teams
	players := []map[string]interface{}{}

	return players, nil
}

// resolvePlayerStatistics resolves the player statistics query
func resolvePlayerStatistics(p graphql.ResolveParams) (interface{}, error) {
	matchID, ok := p.Args["match_id"].(string)
	if !ok {
		return nil, fmt.Errorf("match_id is required")
	}

	// Get resolver context from the context
	resolverCtx, ok := p.Context.Value(resolverContextKey).(*ResolverContext)
	if !ok {
		return nil, fmt.Errorf("resolver context not found")
	}

	// Get the scorecard to calculate player statistics
	_, err := resolverCtx.ScorecardService.GetScorecard(p.Context, matchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get scorecard: %w", err)
	}

	// For now, return empty player statistics
	// In a real implementation, this would calculate statistics from ball-by-ball data
	playerStats := []map[string]interface{}{}

	// This would require:
	// 1. Fetching all players for both teams
	// 2. Analyzing ball-by-ball data to calculate:
	//    - Runs scored by each batsman
	//    - Balls faced by each batsman
	//    - Wickets taken by each bowler
	//    - Overs bowled by each bowler
	//    - Runs conceded by each bowler
	//    - Strike rates and economy rates

	return playerStats, nil
}
