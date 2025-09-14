package supabase

import (
	"context"
	"fmt"
	"log"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/repository/interfaces"
	"time"

	"github.com/supabase-community/supabase-go"
)

type scorecardRepository struct {
	client *supabase.Client
}

// NewScorecardRepository creates a new scorecard repository
func NewScorecardRepository(client *supabase.Client) interfaces.ScorecardRepository {
	return &scorecardRepository{
		client: client,
	}
}

// CreateInnings creates a new innings
func (r *scorecardRepository) CreateInnings(ctx context.Context, innings *models.Innings) error {
	log.Printf("Creating innings for match %s, innings %d, batting team %s", innings.MatchID, innings.InningsNumber, innings.BattingTeam)

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
	_, err := r.client.From("innings").Insert(data, false, "", "", "").ExecuteTo(&result)
	if err != nil {
		log.Printf("Error creating innings: %v", err)
		return fmt.Errorf("failed to create innings: %w", err)
	}

	if len(result) > 0 {
		*innings = result[0]
	}

	log.Printf("Successfully created innings with ID: %s", innings.ID)
	return nil
}

// GetInningsByMatchID gets all innings for a match
func (r *scorecardRepository) GetInningsByMatchID(ctx context.Context, matchID string) ([]*models.Innings, error) {
	log.Printf("Getting innings for match %s", matchID)

	var innings []*models.Innings
	_, err := r.client.From("innings").
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

	var innings []*models.Innings
	_, err := r.client.From("innings").
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

	data := map[string]interface{}{
		"total_runs":    innings.TotalRuns,
		"total_wickets": innings.TotalWickets,
		"total_overs":   innings.TotalOvers,
		"total_balls":   innings.TotalBalls,
		"status":        innings.Status,
		"updated_at":    time.Now(),
	}

	var result []models.Innings
	_, err := r.client.From("innings").
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

	data := map[string]interface{}{
		"status":     string(models.InningsStatusCompleted),
		"updated_at": time.Now(),
	}

	var result []models.Innings
	_, err := r.client.From("innings").
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
	_, err := r.client.From("overs").Insert(data, false, "", "", "").ExecuteTo(&result)
	if err != nil {
		log.Printf("Error creating over: %v", err)
		return fmt.Errorf("failed to create over: %w", err)
	}

	if len(result) > 0 {
		*over = result[0]
	}

	log.Printf("Successfully created over with ID: %s", over.ID)
	return nil
}

// GetOverByInningsAndNumber gets a specific over
func (r *scorecardRepository) GetOverByInningsAndNumber(ctx context.Context, inningsID string, overNumber int) (*models.ScorecardOver, error) {
	log.Printf("Getting over %d for innings %s", overNumber, inningsID)

	var overs []*models.ScorecardOver
	_, err := r.client.From("overs").
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

	var overs []*models.ScorecardOver
	_, err := r.client.From("overs").
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

	var overs []*models.ScorecardOver
	_, err := r.client.From("overs").
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

	data := map[string]interface{}{
		"total_runs":    over.TotalRuns,
		"total_balls":   over.TotalBalls,
		"total_wickets": over.TotalWickets,
		"status":        over.Status,
		"updated_at":    time.Now(),
	}

	var result []models.ScorecardOver
	_, err := r.client.From("overs").
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

	data := map[string]interface{}{
		"status":     string(models.OverStatusCompleted),
		"updated_at": time.Now(),
	}

	var result []models.ScorecardOver
	_, err := r.client.From("overs").
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
	log.Printf("Creating ball %d for over %s, type: %s, run: %s, wicket: %v, wicket_type: %s",
		ball.BallNumber, ball.OverID, string(ball.BallType), string(ball.RunType), ball.IsWicket, ball.WicketType)

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
	_, err := r.client.From("balls").Insert(data, false, "", "", "").ExecuteTo(&result)
	if err != nil {
		log.Printf("Error creating ball: %v", err)
		return fmt.Errorf("failed to create ball: %w", err)
	}

	if len(result) > 0 {
		*ball = result[0]
	}

	log.Printf("Successfully created ball with ID: %s", ball.ID)
	return nil
}

// GetBallsByOver gets all balls for an over
func (r *scorecardRepository) GetBallsByOver(ctx context.Context, overID string) ([]*models.ScorecardBall, error) {
	log.Printf("Getting balls for over %s", overID)

	var balls []*models.ScorecardBall
	_, err := r.client.From("balls").
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

// GetLastBall gets the last ball of an over
func (r *scorecardRepository) GetLastBall(ctx context.Context, overID string) (*models.ScorecardBall, error) {
	log.Printf("Getting last ball for over %s", overID)

	var balls []*models.ScorecardBall
	_, err := r.client.From("balls").
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

// StartScoring starts scoring for a match
func (r *scorecardRepository) StartScoring(ctx context.Context, matchID string) error {
	log.Printf("Starting scoring for match %s", matchID)

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

	// Get match details
	var matches []*models.Match
	_, err := r.client.From("matches").
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
		_, err := r.client.From("overs").
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
		_, err := r.client.From("series").
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
