package supabase

import (
	"context"
	"fmt"
	"log"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/repository/interfaces"
	"strings"
	"time"

	"github.com/supabase-community/supabase-go"
)

type scorecardRepository struct {
	client    *supabase.Client
	schema    string
	matchRepo interfaces.MatchRepository
}

// NewScorecardRepository creates a new scorecard repository
func NewScorecardRepository(client *supabase.Client, schema string, matchRepo interfaces.MatchRepository) interfaces.ScorecardRepository {
	return &scorecardRepository{
		client:    client,
		schema:    schema,
		matchRepo: matchRepo,
	}
}

// getTableName returns the table name (schema is handled by Supabase client configuration)
func (r *scorecardRepository) getTableName(table string) string {
	return table
}

// CreateInnings creates a new innings
func (r *scorecardRepository) CreateInnings(ctx context.Context, innings *models.Innings) error {
	log.Printf("Creating innings for match %s, innings %d, batting team %s", innings.MatchID, innings.InningsNumber, innings.BattingTeam)

	// Add timeout to prevent hanging queries (optimized for simple INSERT)
	_, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	data := map[string]interface{}{
		"match_id":       innings.MatchID,
		"innings_number": innings.InningsNumber,
		"batting_team":   string(innings.BattingTeam),
		"total_runs":     innings.TotalRuns,
		"total_wickets":  innings.TotalWickets,
		"total_overs":    innings.TotalOvers,
		"total_balls":    innings.TotalBalls,
		"status":         innings.Status,
		"created_at":     time.Now(),
		"updated_at":     time.Now(),
	}

	var result []models.Innings
	_, err := r.client.From(r.getTableName("innings")).Insert(data, false, "", "", "").ExecuteTo(&result)
	if err != nil {
		log.Printf("Error creating innings: %v", err)
		return fmt.Errorf("failed to create innings: %w", err)
	}

	if len(result) > 0 {
		*innings = result[0]
		log.Printf("Successfully created innings with ID: %s", innings.ID)
	} else {
		log.Printf("Warning: No result returned from innings creation, but no error occurred")
		// Try to get the innings by match_id and innings_number to get the ID
		createdInnings, err := r.GetInningsByMatchAndNumber(ctx, innings.MatchID, innings.InningsNumber)
		if err != nil {
			log.Printf("Error getting created innings: %v", err)
			return fmt.Errorf("failed to get created innings: %w", err)
		}
		*innings = *createdInnings
		log.Printf("Retrieved created innings with ID: %s", innings.ID)
	}

	log.Printf("Successfully created innings with ID: %s", innings.ID)
	return nil
}

// GetInningsByMatchID gets all innings for a match
func (r *scorecardRepository) GetInningsByMatchID(ctx context.Context, matchID string) ([]*models.Innings, error) {
	log.Printf("Getting innings for match %s", matchID)

	// Add timeout to prevent hanging queries
	_, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	var innings []*models.Innings
	_, err := r.client.From(r.getTableName("innings")).
		Select("*", "", false).
		Eq("match_id", matchID).
		ExecuteTo(&innings)

	if err != nil {
		log.Printf("Error getting innings: %v", err)
		return nil, fmt.Errorf("failed to get innings: %w", err)
	}

	log.Printf("Found %d innings for match %s", len(innings), matchID)
	return innings, nil
}

// GetInningsByMatchAndNumber gets a specific innings
func (r *scorecardRepository) GetInningsByMatchAndNumber(ctx context.Context, matchID string, inningsNumber int) (*models.Innings, error) {
	log.Printf("Getting innings %d for match %s", inningsNumber, matchID)

	// Add timeout to prevent hanging queries (optimized for simple SELECT)
	_, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var innings []*models.Innings
	_, err := r.client.From(r.getTableName("innings")).
		Select("*", "", false).
		Eq("match_id", matchID).
		Eq("innings_number", fmt.Sprintf("%d", inningsNumber)).
		ExecuteTo(&innings)

	if err != nil {
		log.Printf("Error getting innings: %v", err)
		return nil, fmt.Errorf("failed to get innings: %w", err)
	}

	if len(innings) == 0 {
		return nil, fmt.Errorf("innings not found")
	}

	log.Printf("Found innings %d for match %s", inningsNumber, matchID)
	return innings[0], nil
}

// UpdateInnings updates an innings
func (r *scorecardRepository) UpdateInnings(ctx context.Context, innings *models.Innings) error {
	log.Printf("Updating innings %s", innings.ID)

	// Add timeout to prevent hanging queries
	_, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	data := map[string]interface{}{
		"total_runs":    innings.TotalRuns,
		"total_wickets": innings.TotalWickets,
		"total_overs":   innings.TotalOvers,
		"total_balls":   innings.TotalBalls,
		"status":        innings.Status,
		"updated_at":    time.Now(),
	}

	var result []models.Innings
	_, err := r.client.From(r.getTableName("innings")).
		Update(data, "", "").
		Eq("id", innings.ID).
		ExecuteTo(&result)

	if err != nil {
		log.Printf("Error updating innings: %v", err)
		return fmt.Errorf("failed to update innings: %w", err)
	}

	log.Printf("Successfully updated innings %s", innings.ID)
	return nil
}

// CompleteInnings marks an innings as completed
func (r *scorecardRepository) CompleteInnings(ctx context.Context, inningsID string) error {
	log.Printf("Completing innings %s", inningsID)

	// Add timeout to prevent hanging queries
	_, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	data := map[string]interface{}{
		"status":     string(models.InningsStatusCompleted),
		"updated_at": time.Now(),
	}

	var result []models.Innings
	_, err := r.client.From(r.getTableName("innings")).
		Update(data, "", "").
		Eq("id", inningsID).
		ExecuteTo(&result)

	if err != nil {
		log.Printf("Error completing innings: %v", err)
		return fmt.Errorf("failed to complete innings: %w", err)
	}

	log.Printf("Successfully completed innings %s", inningsID)
	return nil
}

// CreateOver creates a new over
func (r *scorecardRepository) CreateOver(ctx context.Context, over *models.ScorecardOver) error {
	log.Printf("Creating over %d for innings %s", over.OverNumber, over.InningsID)

	// Add timeout to prevent hanging queries
	_, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	data := map[string]interface{}{
		"innings_id":    over.InningsID,
		"over_number":   over.OverNumber,
		"total_runs":    over.TotalRuns,
		"total_balls":   over.TotalBalls,
		"total_wickets": over.TotalWickets,
		"status":        over.Status,
		"created_at":    time.Now(),
		"updated_at":    time.Now(),
	}

	var result []models.ScorecardOver
	_, err := r.client.From(r.getTableName("overs")).Insert(data, false, "", "", "").ExecuteTo(&result)
	if err != nil {
		log.Printf("Error creating over: %v", err)
		return fmt.Errorf("failed to create over: %w", err)
	}

	if len(result) > 0 {
		*over = result[0]
		log.Printf("Successfully created over with ID: %s", over.ID)
	} else {
		log.Printf("Warning: No result returned from over creation, but no error occurred")
		// Try to get the over by innings_id and over_number to get the ID
		createdOver, err := r.GetOverByInningsAndNumber(ctx, over.InningsID, over.OverNumber)
		if err != nil {
			log.Printf("Error getting created over: %v", err)
			return fmt.Errorf("failed to get created over: %w", err)
		}
		*over = *createdOver
		log.Printf("Retrieved created over with ID: %s", over.ID)
	}

	log.Printf("Successfully created over with ID: %s", over.ID)
	return nil
}

// GetOverByInningsAndNumber gets a specific over
func (r *scorecardRepository) GetOverByInningsAndNumber(ctx context.Context, inningsID string, overNumber int) (*models.ScorecardOver, error) {
	log.Printf("Getting over %d for innings %s", overNumber, inningsID)

	// Add timeout to prevent hanging queries
	_, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	var overs []*models.ScorecardOver
	_, err := r.client.From(r.getTableName("overs")).
		Select("*", "", false).
		Eq("innings_id", inningsID).
		Eq("over_number", fmt.Sprintf("%d", overNumber)).
		ExecuteTo(&overs)

	if err != nil {
		log.Printf("Error getting over: %v", err)
		return nil, fmt.Errorf("failed to get over: %w", err)
	}

	if len(overs) == 0 {
		return nil, fmt.Errorf("over not found")
	}

	log.Printf("Found over %d for innings %s", overNumber, inningsID)
	return overs[0], nil
}

// GetCurrentOver gets the current in-progress over
func (r *scorecardRepository) GetCurrentOver(ctx context.Context, inningsID string) (*models.ScorecardOver, error) {
	log.Printf("Getting current over for innings %s", inningsID)

	// Add timeout to prevent hanging queries (optimized for simple SELECT with LIMIT)
	_, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var overs []*models.ScorecardOver
	_, err := r.client.From(r.getTableName("overs")).
		Select("*", "", false).
		Eq("innings_id", inningsID).
		Eq("status", string(models.OverStatusInProgress)).
		Limit(1, "").
		ExecuteTo(&overs)

	if err != nil {
		log.Printf("Error getting current over: %v", err)
		return nil, fmt.Errorf("failed to get current over: %w", err)
	}

	if len(overs) == 0 {
		return nil, fmt.Errorf("no current over found")
	}

	log.Printf("Found current over %d for innings %s", overs[0].OverNumber, inningsID)
	return overs[0], nil
}

// GetOversByInnings gets all overs for an innings
func (r *scorecardRepository) GetOversByInnings(ctx context.Context, inningsID string) ([]*models.ScorecardOver, error) {
	log.Printf("Getting all overs for innings %s", inningsID)

	// Add timeout to prevent hanging queries (longer timeout for complex queries)
	_, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	var overs []*models.ScorecardOver
	_, err := r.client.From(r.getTableName("overs")).
		Select("*", "", false).
		Eq("innings_id", inningsID).
		ExecuteTo(&overs)

	if err != nil {
		log.Printf("Error getting overs: %v", err)
		return nil, fmt.Errorf("failed to get overs: %w", err)
	}

	log.Printf("Found %d overs for innings %s", len(overs), inningsID)
	return overs, nil
}

// UpdateOver updates an over
func (r *scorecardRepository) UpdateOver(ctx context.Context, over *models.ScorecardOver) error {
	log.Printf("Updating over %s", over.ID)

	// Add timeout to prevent hanging queries
	_, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	data := map[string]interface{}{
		"total_runs":    over.TotalRuns,
		"total_balls":   over.TotalBalls,
		"total_wickets": over.TotalWickets,
		"status":        over.Status,
		"updated_at":    time.Now(),
	}

	var result []models.ScorecardOver
	_, err := r.client.From(r.getTableName("overs")).
		Update(data, "", "").
		Eq("id", over.ID).
		ExecuteTo(&result)

	if err != nil {
		log.Printf("Error updating over: %v", err)
		return fmt.Errorf("failed to update over: %w", err)
	}

	log.Printf("Successfully updated over %s", over.ID)
	return nil
}

// CompleteOver marks an over as completed
func (r *scorecardRepository) CompleteOver(ctx context.Context, overID string) error {
	log.Printf("Completing over %s", overID)

	// Add timeout to prevent hanging queries
	_, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	data := map[string]interface{}{
		"status":     string(models.OverStatusCompleted),
		"updated_at": time.Now(),
	}

	var result []models.ScorecardOver
	_, err := r.client.From(r.getTableName("overs")).
		Update(data, "", "").
		Eq("id", overID).
		ExecuteTo(&result)

	if err != nil {
		log.Printf("Error completing over: %v", err)
		return fmt.Errorf("failed to complete over: %w", err)
	}

	log.Printf("Successfully completed over %s", overID)
	return nil
}

// CreateBall creates a new ball
func (r *scorecardRepository) CreateBall(ctx context.Context, ball *models.ScorecardBall) error {
	log.Printf("Creating ball %d for over %s", ball.BallNumber, ball.OverID)

	// Note: Supabase client doesn't directly support context timeouts
	// Timeout is handled at the HTTP client level

	data := map[string]interface{}{
		"over_id":     ball.OverID,
		"ball_number": ball.BallNumber,
		"ball_type":   string(ball.BallType),
		"run_type":    string(ball.RunType),
		"runs":        ball.Runs,
		"byes":        ball.Byes,
		"is_wicket":   ball.IsWicket,
		"created_at":  time.Now(),
	}

	// Only include wicket_type if it's a wicket
	if ball.IsWicket && ball.WicketType != "" {
		data["wicket_type"] = ball.WicketType
	}

	var result []models.ScorecardBall
	_, err := r.client.From(r.getTableName("balls")).Insert(data, false, "", "", "").ExecuteTo(&result)
	if err != nil {
		log.Printf("Error creating ball: %v", err)
		return fmt.Errorf("failed to create ball: %w", err)
	}

	if len(result) > 0 {
		*ball = result[0]
		log.Printf("Successfully created ball with ID: %s", ball.ID)
	} else {
		log.Printf("Warning: No result returned from ball creation, but no error occurred")
		// Try to get the ball by over_id and ball_number to get the ID
		createdBall, err := r.GetBallByOverAndNumber(ctx, ball.OverID, ball.BallNumber)
		if err != nil {
			log.Printf("Error getting created ball: %v", err)
			return fmt.Errorf("failed to get created ball: %w", err)
		}
		*ball = *createdBall
		log.Printf("Retrieved created ball with ID: %s", ball.ID)
	}

	log.Printf("Successfully created ball with ID: %s", ball.ID)
	return nil
}

// GetBallsByOver gets all balls for an over
func (r *scorecardRepository) GetBallsByOver(ctx context.Context, overID string) ([]*models.ScorecardBall, error) {
	log.Printf("Getting balls for over %s", overID)

	// Add timeout to prevent hanging queries (longer timeout for complex queries)
	_, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	var balls []*models.ScorecardBall
	_, err := r.client.From(r.getTableName("balls")).
		Select("*", "", false).
		Eq("over_id", overID).
		ExecuteTo(&balls)

	if err != nil {
		log.Printf("Error getting balls: %v", err)
		return nil, fmt.Errorf("failed to get balls: %w", err)
	}

	log.Printf("Found %d balls for over %s", len(balls), overID)
	return balls, nil
}

// GetBallCountByOver gets the count of balls for an over (optimized for performance)
func (r *scorecardRepository) GetBallCountByOver(ctx context.Context, overID string) (int, error) {
	log.Printf("Getting ball count for over %s", overID)

	// Add timeout to prevent hanging queries
	_, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var balls []*models.ScorecardBall
	_, err := r.client.From(r.getTableName("balls")).
		Select("id", "", false). // Only select ID for counting
		Eq("over_id", overID).
		ExecuteTo(&balls)

	if err != nil {
		log.Printf("Error getting ball count: %v", err)
		return 0, fmt.Errorf("failed to get ball count: %w", err)
	}

	ballCount := len(balls)
	log.Printf("Found %d balls for over %s", ballCount, overID)
	return ballCount, nil
}

// GetBallsForNextNumber gets only the necessary fields for ball number calculation (optimized)
func (r *scorecardRepository) GetBallsForNextNumber(ctx context.Context, overID string) ([]*models.ScorecardBall, error) {
	log.Printf("Getting balls for next number calculation for over %s", overID)

	// Add timeout to prevent hanging queries (optimized for simple SELECT)
	_, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var balls []*models.ScorecardBall
	_, err := r.client.From(r.getTableName("balls")).
		Select("ball_number,ball_type", "", false). // Only select necessary fields
		Eq("over_id", overID).
		ExecuteTo(&balls)

	if err != nil {
		log.Printf("Error getting balls for next number: %v", err)
		return nil, fmt.Errorf("failed to get balls for next number: %w", err)
	}

	log.Printf("Found %d balls for next number calculation for over %s", len(balls), overID)
	return balls, nil
}

// GetBallByOverAndNumber gets a specific ball by over ID and ball number
func (r *scorecardRepository) GetBallByOverAndNumber(ctx context.Context, overID string, ballNumber int) (*models.ScorecardBall, error) {
	log.Printf("Getting ball %d for over %s", ballNumber, overID)

	// Add timeout to prevent hanging queries
	_, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	var balls []*models.ScorecardBall
	_, err := r.client.From(r.getTableName("balls")).
		Select("*", "", false).
		Eq("over_id", overID).
		Eq("ball_number", fmt.Sprintf("%d", ballNumber)).
		ExecuteTo(&balls)

	if err != nil {
		log.Printf("Error getting ball: %v", err)
		return nil, fmt.Errorf("failed to get ball: %w", err)
	}

	if len(balls) == 0 {
		return nil, fmt.Errorf("ball not found")
	}

	log.Printf("Found ball %d for over %s", ballNumber, overID)
	return balls[0], nil
}

// GetLastBall gets the last ball of an over
func (r *scorecardRepository) GetLastBall(ctx context.Context, overID string) (*models.ScorecardBall, error) {
	log.Printf("Getting last ball for over %s", overID)

	// Add timeout to prevent hanging queries
	_, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	var balls []*models.ScorecardBall
	_, err := r.client.From(r.getTableName("balls")).
		Select("*", "", false).
		Eq("over_id", overID).
		Limit(1, "").
		ExecuteTo(&balls)

	if err != nil {
		log.Printf("Error getting last ball: %v", err)
		return nil, fmt.Errorf("failed to get last ball: %w", err)
	}

	if len(balls) == 0 {
		return nil, fmt.Errorf("no balls found")
	}

	log.Printf("Found last ball %d for over %s", balls[0].BallNumber, overID)
	return balls[0], nil
}

// DeleteBall deletes a ball by ID
func (r *scorecardRepository) DeleteBall(ctx context.Context, ballID string) error {
	log.Printf("Deleting ball %s", ballID)

	// Add timeout to prevent hanging queries
	_, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	_, _, err := r.client.From(r.getTableName("balls")).
		Delete("", "").
		Eq("id", ballID).
		Execute()

	if err != nil {
		log.Printf("Error deleting ball: %v", err)
		return fmt.Errorf("failed to delete ball: %w", err)
	}

	log.Printf("Successfully deleted ball %s", ballID)
	return nil
}

// StartScoring starts scoring for a match
func (r *scorecardRepository) StartScoring(ctx context.Context, matchID string) error {
	log.Printf("Starting scoring for match %s", matchID)

	// Add timeout to prevent hanging queries (longer timeout for initialization)
	_, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	// Create first innings
	innings := &models.Innings{
		MatchID:       matchID,
		InningsNumber: 1,
		BattingTeam:   models.TeamTypeA, // Will be updated based on toss winner
		TotalRuns:     0,
		TotalWickets:  0,
		TotalOvers:    0.0,
		TotalBalls:    0,
		Status:        string(models.InningsStatusInProgress),
	}

	err := r.CreateInnings(ctx, innings)
	if err != nil {
		log.Printf("Error creating first innings: %v", err)
		return fmt.Errorf("failed to create first innings: %w", err)
	}

	log.Printf("Successfully started scoring for match %s", matchID)
	return nil
}

// GetScorecard gets the complete scorecard for a match
func (r *scorecardRepository) GetScorecard(ctx context.Context, matchID string) (*models.ScorecardResponse, error) {
	log.Printf("Getting scorecard for match %s", matchID)

	// Add timeout to prevent hanging queries (longer timeout for complex scorecard building)
	_, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	// Get match details
	var matches []*models.Match
	_, err := r.client.From(r.getTableName("matches")).
		Select("*, series(name)", "", false).
		Eq("id", matchID).
		ExecuteTo(&matches)

	if err != nil {
		log.Printf("Error getting match: %v", err)
		return nil, fmt.Errorf("failed to get match: %w", err)
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("match not found")
	}

	match := matches[0]

	// Get innings
	innings, err := r.GetInningsByMatchID(ctx, matchID)
	if err != nil {
		log.Printf("Error getting innings: %v", err)
		return nil, fmt.Errorf("failed to get innings: %w", err)
	}

	// Build innings summaries
	var inningsSummaries []models.InningsSummary
	for _, inn := range innings {
		// Get overs for this innings
		var overs []*models.ScorecardOver
		_, err := r.client.From(r.getTableName("overs")).
			Select("*", "", false).
			Eq("innings_id", inn.ID).
			ExecuteTo(&overs)

		if err != nil {
			log.Printf("Error getting overs: %v", err)
			continue
		}

		// Build over summaries and calculate extras
		var overSummaries []models.OverSummary
		extras := &models.ExtrasSummary{}

		for _, over := range overs {
			// Get balls for this over
			balls, err := r.GetBallsByOver(ctx, over.ID)
			if err != nil {
				log.Printf("Error getting balls: %v", err)
				continue
			}

			// Build ball summaries and calculate extras
			var ballSummaries []models.BallSummary
			for _, ball := range balls {
				ballSummaries = append(ballSummaries, models.BallSummary{
					BallNumber: ball.BallNumber,
					BallType:   ball.BallType,
					RunType:    ball.RunType,
					Runs:       ball.Runs,
					Byes:       ball.Byes,
					IsWicket:   ball.IsWicket,
					WicketType: ball.WicketType,
				})

				// Calculate extras
				switch ball.BallType {
				case models.BallTypeWide:
					extras.Wides += ball.Runs
					if ball.Byes > 0 {
						extras.Byes += ball.Byes
					}
				case models.BallTypeNoBall:
					extras.NoBalls += ball.Runs
					if ball.Byes > 0 {
						extras.Byes += ball.Byes
					}
				case models.BallTypeGood:
					if ball.RunType == models.RunTypeLB {
						extras.LegByes += ball.Runs
						if ball.Byes > 0 {
							extras.Byes += ball.Byes
						}
					} else if ball.Byes > 0 {
						extras.Byes += ball.Byes
					}
				}
			}

			overSummaries = append(overSummaries, models.OverSummary{
				OverNumber:   over.OverNumber,
				TotalRuns:    over.TotalRuns,
				TotalBalls:   over.TotalBalls,
				TotalWickets: over.TotalWickets,
				Status:       over.Status,
				Balls:        ballSummaries,
			})
		}

		// Calculate total extras
		extras.Total = extras.Byes + extras.LegByes + extras.Wides + extras.NoBalls

		inningsSummaries = append(inningsSummaries, models.InningsSummary{
			InningsNumber: inn.InningsNumber,
			BattingTeam:   inn.BattingTeam,
			TotalRuns:     inn.TotalRuns,
			TotalWickets:  inn.TotalWickets,
			TotalOvers:    inn.TotalOvers,
			TotalBalls:    inn.TotalBalls,
			Status:        inn.Status,
			Extras:        extras,
			Overs:         overSummaries,
		})
	}

	// Determine current innings
	currentInnings := 1
	if len(inningsSummaries) > 0 {
		for _, inn := range inningsSummaries {
			if inn.Status == string(models.InningsStatusInProgress) {
				currentInnings = inn.InningsNumber
				break
			}
		}
	}

	// Get series name
	seriesName := "Unknown Series"
	if match.SeriesID != "" {
		var series []*models.Series
		_, err := r.client.From(r.getTableName("series")).
			Select("name", "", false).
			Eq("id", match.SeriesID).
			ExecuteTo(&series)

		if err == nil && len(series) > 0 {
			seriesName = series[0].Name
		}
	}

	scorecard := &models.ScorecardResponse{
		MatchID:        matchID,
		MatchNumber:    match.MatchNumber,
		SeriesName:     seriesName,
		TeamA:          "Team A",
		TeamB:          "Team B",
		TotalOvers:     match.TotalOvers,
		TossWinner:     match.TossWinner,
		TossType:       match.TossType,
		CurrentInnings: currentInnings,
		Innings:        inningsSummaries,
		MatchStatus:    string(match.Status),
	}

	log.Printf("Successfully built scorecard for match %s", matchID)
	return scorecard, nil
}

// GetMatchInningsOverData gets all necessary data for add ball API in a single optimized query
func (r *scorecardRepository) GetMatchInningsOverData(ctx context.Context, matchID string, inningsNumber int) (*models.MatchInningsOverData, error) {
	log.Printf("Getting optimized match data for match %s, innings %d", matchID, inningsNumber)

	// Add timeout to prevent hanging queries (optimized for complex JOIN)
	_, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// For now, use fallback method since Supabase doesn't support complex JOINs easily
	// In production, this would use a stored procedure or raw SQL
	return r.getMatchInningsOverDataFallback(ctx, matchID, inningsNumber)
}

// getMatchInningsOverDataFallback provides fallback implementation using individual queries
func (r *scorecardRepository) getMatchInningsOverDataFallback(ctx context.Context, matchID string, inningsNumber int) (*models.MatchInningsOverData, error) {
	log.Printf("Using fallback method for match %s, innings %d", matchID, inningsNumber)

	// Get match data
	match, err := r.matchRepo.GetByID(ctx, matchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get match: %w", err)
	}

	// Get innings data
	innings, err := r.GetInningsByMatchAndNumber(ctx, matchID, inningsNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get innings: %w", err)
	}

	// Get current over data
	over, err := r.GetCurrentOver(ctx, innings.ID)
	if err != nil {
		// No current over exists, try to create the first over
		log.Printf("No current over found, creating first over for innings %s", innings.ID)

		// Get all overs for this innings to determine the next over number
		allOvers, err := r.GetOversByInnings(ctx, innings.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get overs for innings: %w", err)
		}

		overNumber := len(allOvers) + 1
		newOver := &models.ScorecardOver{
			InningsID:    innings.ID,
			OverNumber:   overNumber,
			TotalRuns:    0,
			TotalBalls:   0,
			TotalWickets: 0,
			Status:       string(models.OverStatusInProgress),
		}

		// Try to create the over in the database
		err = r.CreateOver(ctx, newOver)
		if err != nil {
			// If creation failed due to duplicate key (race condition), try to get the existing over
			if strings.Contains(err.Error(), "duplicate key") {
				log.Printf("Over creation failed due to race condition, attempting to get existing over")
				over, err = r.GetCurrentOver(ctx, innings.ID)
				if err != nil {
					return nil, fmt.Errorf("failed to get existing over after race condition: %w", err)
				}
				log.Printf("Successfully retrieved existing over %d for innings %s", over.OverNumber, innings.ID)
			} else {
				return nil, fmt.Errorf("failed to create first over: %w", err)
			}
		} else {
			log.Printf("Successfully created first over %d for innings %s", overNumber, innings.ID)
			over = newOver
		}
	}

	// Get ball count for current over
	ballCount, err := r.GetBallCountByOver(ctx, over.ID)
	if err != nil {
		ballCount = 0
	}

	// Get balls for next number calculation
	balls, err := r.GetBallsForNextNumber(ctx, over.ID)
	if err != nil {
		balls = []*models.ScorecardBall{}
	}

	// Calculate statistics
	legalBallCount := 0
	maxBallNumber := 0
	for _, ball := range balls {
		if ball.BallNumber > maxBallNumber {
			maxBallNumber = ball.BallNumber
		}
		if ball.BallType == models.BallTypeGood {
			legalBallCount++
		}
	}

	// Get overs for statistics
	overs, err := r.GetOversByInnings(ctx, innings.ID)
	if err != nil {
		overs = []*models.ScorecardOver{}
	}

	completedOvers := 0
	currentOverBalls := 0
	for _, o := range overs {
		if o.Status == string(models.OverStatusCompleted) {
			completedOvers++
		} else if o.Status == string(models.OverStatusInProgress) {
			currentOverBalls = o.TotalBalls
		}
	}

	return &models.MatchInningsOverData{
		// Match data
		MatchID:          match.ID,
		MatchStatus:      match.Status,
		CreatedBy:        match.CreatedBy,
		BattingTeam:      match.BattingTeam,
		TotalOvers:       match.TotalOvers,
		TeamAPlayerCount: match.TeamAPlayerCount,

		// Innings data
		InningsID:           innings.ID,
		InningsNumber:       innings.InningsNumber,
		InningsStatus:       models.InningsStatus(innings.Status),
		InningsTotalRuns:    innings.TotalRuns,
		InningsTotalWickets: innings.TotalWickets,
		InningsTotalOvers:   innings.TotalOvers,
		InningsTotalBalls:   innings.TotalBalls,

		// Over data
		OverID:           over.ID,
		OverNumber:       over.OverNumber,
		OverStatus:       models.OverStatus(over.Status),
		OverTotalRuns:    over.TotalRuns,
		OverTotalBalls:   over.TotalBalls,
		OverTotalWickets: over.TotalWickets,

		// Ball data
		BallCount:      ballCount,
		LegalBallCount: legalBallCount,
		MaxBallNumber:  maxBallNumber,

		// Overs data
		CompletedOvers:   completedOvers,
		CurrentOverBalls: currentOverBalls,

		// Timestamps
		CreatedAt: match.CreatedAt,
		UpdatedAt: match.UpdatedAt,
	}, nil
}
