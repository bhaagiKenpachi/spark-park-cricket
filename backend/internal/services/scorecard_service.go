package services

import (
	"context"
	"fmt"
	"log"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/repository/interfaces"
	"spark-park-cricket-backend/internal/utils"
)

type ScorecardService struct {
	scorecardRepo interfaces.ScorecardRepository
	matchRepo     interfaces.MatchRepository
}

// NewScorecardService creates a new scorecard service
func NewScorecardService(scorecardRepo interfaces.ScorecardRepository, matchRepo interfaces.MatchRepository) *ScorecardService {
	return &ScorecardService{
		scorecardRepo: scorecardRepo,
		matchRepo:     matchRepo,
	}
}

// StartScoring starts scoring for a match
func (s *ScorecardService) StartScoring(ctx context.Context, matchID string) error {
	log.Printf("Starting scoring for match %s", matchID)

	// Get match details
	match, err := s.matchRepo.GetByID(ctx, matchID)
	if err != nil {
		log.Printf("Error getting match: %v", err)
		return fmt.Errorf("match not found: %w", err)
	}

	// Check if match is live
	if match.Status != models.MatchStatusLive {
		return fmt.Errorf("match is not live, cannot start scoring")
	}

	// Check if scoring is already started
	innings, err := s.scorecardRepo.GetInningsByMatchID(ctx, matchID)
	if err == nil && len(innings) > 0 {
		return fmt.Errorf("scoring already started for this match")
	}

	// Create first innings with toss winner as batting team
	firstInnings := &models.Innings{
		MatchID:       matchID,
		InningsNumber: 1,
		BattingTeam:   match.TossWinner, // Toss winner bats first
		TotalRuns:     0,
		TotalWickets:  0,
		TotalOvers:    0.0,
		TotalBalls:    0,
		Status:        string(models.InningsStatusInProgress),
	}

	err = s.scorecardRepo.CreateInnings(ctx, firstInnings)
	if err != nil {
		log.Printf("Error creating first innings: %v", err)
		return fmt.Errorf("failed to start scoring: %w", err)
	}

	log.Printf("Successfully started scoring for match %s, first innings batting team: %s", matchID, match.TossWinner)
	return nil
}

// AddBall adds a ball to the scorecard
func (s *ScorecardService) AddBall(ctx context.Context, req *models.BallEventRequest) error {

	// Validate ball event
	if err := utils.ValidateBallEventRequest(req); err != nil {
		log.Printf("Invalid ball event: %v", err)
		return fmt.Errorf("invalid ball event: %w", err)
	}

	// Get match details
	match, err := s.matchRepo.GetByID(ctx, req.MatchID)
	if err != nil {
		log.Printf("Error getting match: %v", err)
		return fmt.Errorf("match not found: %w", err)
	}

	// Check if match is live
	if match.Status != models.MatchStatusLive {
		return fmt.Errorf("match is not live, cannot add ball")
	}

	// Validate innings order
	if err := s.ValidateInningsOrder(ctx, req.MatchID, match, req.InningsNumber); err != nil {
		return fmt.Errorf("innings validation failed: %w", err)
	}

	// Get innings or create if doesn't exist
	innings, err := s.scorecardRepo.GetInningsByMatchAndNumber(ctx, req.MatchID, req.InningsNumber)
	if err != nil {
		log.Printf("Innings not found, creating new innings: %v", err)
		// Create innings if it doesn't exist
		innings = &models.Innings{
			MatchID:       req.MatchID,
			InningsNumber: req.InningsNumber,
			BattingTeam:   match.BattingTeam,
			TotalRuns:     0,
			TotalWickets:  0,
			TotalOvers:    0.0,
			TotalBalls:    0,
			Status:        string(models.InningsStatusInProgress),
		}
		err = s.scorecardRepo.CreateInnings(ctx, innings)
		if err != nil {
			log.Printf("Error creating innings: %v", err)
			return fmt.Errorf("failed to create innings: %w", err)
		}
		log.Printf("Created new innings %d for match %s", req.InningsNumber, req.MatchID)
	}

	// Check if innings is in progress
	if innings.Status != string(models.InningsStatusInProgress) {
		return fmt.Errorf("innings is not in progress, cannot add ball")
	}

	// Get current over or create new one
	over, err := s.getCurrentOver(ctx, innings.ID)
	if err != nil {
		log.Printf("Error getting current over: %v", err)
		return fmt.Errorf("failed to get current over: %w", err)
	}

	// Check if over is in progress
	if over.Status != string(models.OverStatusInProgress) {
		return fmt.Errorf("over is not in progress, cannot add ball")
	}


	// Get next ball number
	ballNumber, err := s.getNextBallNumber(ctx, over.ID)
	if err != nil {
		log.Printf("Error getting next ball number: %v", err)
		return fmt.Errorf("failed to get next ball number: %w", err)
	}

	// Calculate runs from run type and byes
	runs := req.RunType.GetRunValue()
	byes := req.Byes
	if req.IsWicket && req.RunType == models.RunTypeWC {
		runs = 0 // Wicket doesn't count as runs
	}

	// Total runs = ball runs + byes
	totalRuns := runs + byes

	// Create ball
	ball := &models.ScorecardBall{
		OverID:     over.ID,
		BallNumber: ballNumber,
		BallType:   req.BallType,
		RunType:    req.RunType,
		Runs:       runs,
		Byes:       byes,
		IsWicket:   req.IsWicket,
		WicketType: req.WicketType,
	}

	err = s.scorecardRepo.CreateBall(ctx, ball)
	if err != nil {
		log.Printf("Error creating ball: %v", err)
		return fmt.Errorf("failed to add ball: %w", err)
	}

	// Update over statistics
	over.TotalRuns += totalRuns
	// Only count legal balls (good balls) for over completion
	if req.BallType == models.BallTypeGood {
		over.TotalBalls++
	}
	if req.IsWicket {
		over.TotalWickets++
	}

	// Check if over is complete (6 legal balls or all wickets)
	if over.TotalBalls >= 6 || over.TotalWickets >= 10 {
		over.Status = string(models.OverStatusCompleted)
	}

	err = s.scorecardRepo.UpdateOver(ctx, over)
	if err != nil {
		log.Printf("Error updating over: %v", err)
		return fmt.Errorf("failed to update over: %w", err)
	}

	// If over is completed and we need to add more balls, create new over
	if over.Status == string(models.OverStatusCompleted) && req.BallType == models.BallTypeGood && over.TotalBalls > 6 {
		log.Printf("Over %d is complete with %d legal balls, creating new over", over.OverNumber, over.TotalBalls)
		
		// Get new over
		over, err = s.getCurrentOver(ctx, innings.ID)
		if err != nil {
			log.Printf("Error getting new current over: %v", err)
			return fmt.Errorf("failed to get new current over: %w", err)
		}
	}

	// Update innings statistics
	innings.TotalRuns += totalRuns
	// Only count legal balls for innings overs calculation
	if req.BallType == models.BallTypeGood {
		innings.TotalBalls++
	}
	if req.IsWicket {
		innings.TotalWickets++
	}

	// Calculate total overs properly
	// Get all overs for this innings to calculate completed overs + current over balls
	overs, err := s.scorecardRepo.GetOversByInnings(ctx, innings.ID)
	if err != nil {
		log.Printf("Error getting overs for innings: %v", err)
		return fmt.Errorf("failed to get overs: %w", err)
	}

	completedOvers := 0
	currentOverBalls := 0

	for _, over := range overs {
		if over.Status == string(models.OverStatusCompleted) {
			completedOvers++
		} else if over.Status == string(models.OverStatusInProgress) {
			currentOverBalls = over.TotalBalls
		}
	}

	// Total overs = completed overs + current over balls as decimal
	// In cricket scoring: 3 balls = 0.3 overs, 4 balls = 0.4 overs, etc.
	// We need to convert balls to the cricket scoring format
	var currentOverDecimal float64
	if currentOverBalls > 0 {
		// Convert balls to cricket scoring format (0.1, 0.2, 0.3, 0.4, 0.5, 1.0)
		if currentOverBalls == 6 {
			currentOverDecimal = 1.0
		} else {
			currentOverDecimal = float64(currentOverBalls) / 10.0
		}
	}
	innings.TotalOvers = float64(completedOvers) + currentOverDecimal

	// Check if innings is complete
	// For first innings: complete when all wickets are taken or all overs are completed
	// For second innings: completion is handled by shouldCompleteMatch method
	maxWickets := match.TeamAPlayerCount - 1 // n-1 wickets for n players
	if req.InningsNumber == 1 {
		if innings.TotalWickets >= maxWickets || innings.TotalOvers >= float64(match.TotalOvers) {
			innings.Status = string(models.InningsStatusCompleted)
			log.Printf("First innings %d completed for match %s: wickets=%d/%d, overs=%.1f/%d",
				innings.InningsNumber, match.ID, innings.TotalWickets, maxWickets, innings.TotalOvers, match.TotalOvers)
		}
	}
	// For second innings, we don't automatically complete here - let shouldCompleteMatch handle it

	err = s.scorecardRepo.UpdateInnings(ctx, innings)
	if err != nil {
		log.Printf("Error updating innings: %v", err)
		return fmt.Errorf("failed to update innings: %w", err)
	}

	// Handle match progression
	if req.InningsNumber == 1 {
		// First innings - check if completed and start second innings
		if innings.Status == string(models.InningsStatusCompleted) {
			err = s.startSecondInnings(ctx, req.MatchID, match)
			if err != nil {
				log.Printf("Error starting second innings: %v", err)
				return fmt.Errorf("failed to start second innings: %w", err)
			}
			log.Printf("Second innings started for match %s", req.MatchID)
		}
	} else if req.InningsNumber == 2 {
		// Second innings - check for match completion after every ball
		shouldCompleteMatch, reason := s.ShouldCompleteMatch(ctx, req.MatchID, innings, match)
		if shouldCompleteMatch {
			// Complete the innings first
			innings.Status = string(models.InningsStatusCompleted)
			err = s.scorecardRepo.UpdateInnings(ctx, innings)
			if err != nil {
				return fmt.Errorf("failed to update innings status: %w", err)
			}

			// Complete the match
			match.Status = models.MatchStatusCompleted
			err = s.matchRepo.Update(ctx, req.MatchID, match)
			if err != nil {
				return fmt.Errorf("failed to complete match: %w", err)
			}
			log.Printf("Match %s completed - %s", req.MatchID, reason)
		}
	}

	log.Printf("Successfully added ball: %s %d runs, byes: %d, total: %d, wicket: %v", req.RunType, runs, byes, totalRuns, req.IsWicket)
	return nil
}

// getCurrentOver gets the current in-progress over or creates a new one
func (s *ScorecardService) getCurrentOver(ctx context.Context, inningsID string) (*models.ScorecardOver, error) {
	// Try to get current over
	over, err := s.scorecardRepo.GetCurrentOver(ctx, inningsID)
	if err == nil && over != nil {
		return over, nil
	}

	// Get all overs for this innings to determine next over number
	overs, err := s.scorecardRepo.GetOversByInnings(ctx, inningsID)
	if err != nil {
		log.Printf("Error getting overs: %v", err)
		return nil, fmt.Errorf("failed to get overs: %w", err)
	}

	// Calculate next over number
	nextOverNumber := 1
	if len(overs) > 0 {
		// Find the highest over number and add 1
		maxOverNumber := 0
		for _, o := range overs {
			if o.OverNumber > maxOverNumber {
				maxOverNumber = o.OverNumber
			}
		}
		nextOverNumber = maxOverNumber + 1
	}

	// Create new over
	newOver := &models.ScorecardOver{
		InningsID:    inningsID,
		OverNumber:   nextOverNumber,
		TotalRuns:    0,
		TotalBalls:   0,
		TotalWickets: 0,
		Status:       string(models.OverStatusInProgress),
	}

	err = s.scorecardRepo.CreateOver(ctx, newOver)
	if err != nil {
		log.Printf("Error creating over: %v", err)
		return nil, fmt.Errorf("failed to create over: %w", err)
	}

	log.Printf("Created new over %d for innings %s", nextOverNumber, inningsID)
	return newOver, nil
}

// getNextBallNumber gets the next ball number for an over
func (s *ScorecardService) getNextBallNumber(ctx context.Context, overID string) (int, error) {
	// Get all balls for this over
	balls, err := s.scorecardRepo.GetBallsByOver(ctx, overID)
	if err != nil {
		return 0, fmt.Errorf("failed to get balls: %w", err)
	}

	// Count legal balls (good balls only) and find max ball number
	legalBalls := 0
	maxBallNumber := 0

	for _, ball := range balls {
		if ball.BallNumber > maxBallNumber {
			maxBallNumber = ball.BallNumber
		}
		// Only count good balls as legal deliveries for over completion
		if ball.BallType == models.BallTypeGood {
			legalBalls++
		}
	}

	// An over is complete when it has 6 legal balls
	if legalBalls >= 6 {
		return 0, fmt.Errorf("over is complete, cannot add more balls")
	}

	// The next ball number is simply the next sequential number
	nextBallNumber := maxBallNumber + 1

	return nextBallNumber, nil
}

// ShouldCompleteMatch determines if the match should be completed based on cricket rules
func (s *ScorecardService) ShouldCompleteMatch(ctx context.Context, matchID string, secondInnings *models.Innings, match *models.Match) (bool, string) {
	// Get first innings score
	firstInnings, err := s.scorecardRepo.GetInningsByMatchAndNumber(ctx, matchID, 1)
	if err != nil {
		return false, "error getting first innings"
	}

	target := firstInnings.TotalRuns + 1     // Target is first innings score + 1
	maxWickets := match.TeamAPlayerCount - 1 // n-1 wickets for n players

	// Check if target is reached
	if secondInnings.TotalRuns >= target {
		return true, fmt.Sprintf("target reached: %d/%d", secondInnings.TotalRuns, target)
	}

	// Check if all wickets are lost
	if secondInnings.TotalWickets >= maxWickets {
		return true, fmt.Sprintf("all wickets lost: %d/%d", secondInnings.TotalWickets, maxWickets)
	}

	// Check if all overs are completed
	if secondInnings.TotalOvers >= float64(match.TotalOvers) {
		return true, fmt.Sprintf("all overs completed: %.1f/%d", secondInnings.TotalOvers, match.TotalOvers)
	}

	return false, "match continues"
}

// startSecondInnings starts the second innings
func (s *ScorecardService) startSecondInnings(ctx context.Context, matchID string, match *models.Match) error {
	log.Printf("Starting second innings for match %s", matchID)

	// Determine batting team for second innings (opposite of first innings)
	var battingTeam models.TeamType
	if match.BattingTeam == models.TeamTypeA {
		battingTeam = models.TeamTypeB
	} else {
		battingTeam = models.TeamTypeA
	}

	// Create second innings
	secondInnings := &models.Innings{
		MatchID:       matchID,
		InningsNumber: 2,
		BattingTeam:   battingTeam,
		TotalRuns:     0,
		TotalWickets:  0,
		TotalOvers:    0.0,
		TotalBalls:    0,
		Status:        string(models.InningsStatusInProgress),
	}

	err := s.scorecardRepo.CreateInnings(ctx, secondInnings)
	if err != nil {
		log.Printf("Error creating second innings: %v", err)
		return fmt.Errorf("failed to start second innings: %w", err)
	}

	// Update match batting team
	match.BattingTeam = battingTeam
	err = s.matchRepo.Update(ctx, matchID, match)
	if err != nil {
		log.Printf("Error updating match batting team: %v", err)
		return fmt.Errorf("failed to update match batting team: %w", err)
	}

	log.Printf("Successfully started second innings for match %s, batting team: %s", matchID, battingTeam)
	return nil
}

// GetScorecard gets the complete scorecard for a match
func (s *ScorecardService) GetScorecard(ctx context.Context, matchID string) (*models.ScorecardResponse, error) {
	log.Printf("Getting scorecard for match %s", matchID)

	// Get scorecard from repository
	scorecard, err := s.scorecardRepo.GetScorecard(ctx, matchID)
	if err != nil {
		log.Printf("Error getting scorecard: %v", err)
		return nil, fmt.Errorf("failed to get scorecard: %w", err)
	}

	log.Printf("Successfully retrieved scorecard for match %s", matchID)
	return scorecard, nil
}

// GetCurrentOver gets the current over for a match
func (s *ScorecardService) GetCurrentOver(ctx context.Context, matchID string, inningsNumber int) (*models.ScorecardOver, error) {
	log.Printf("Getting current over for match %s, innings %d", matchID, inningsNumber)

	// Get innings
	innings, err := s.scorecardRepo.GetInningsByMatchAndNumber(ctx, matchID, inningsNumber)
	if err != nil {
		log.Printf("Error getting innings: %v", err)
		return nil, fmt.Errorf("innings not found: %w", err)
	}

	// Get current over (only existing ones, don't create new)
	over, err := s.scorecardRepo.GetCurrentOver(ctx, innings.ID)
	if err != nil {
		log.Printf("Error getting current over: %v", err)
		return nil, fmt.Errorf("no current over found: %w", err)
	}

	log.Printf("Found current over %d for match %s, innings %d", over.OverNumber, matchID, inningsNumber)
	return over, nil
}

// ValidateInningsOrder validates that balls can only be added to the correct innings
func (s *ScorecardService) ValidateInningsOrder(ctx context.Context, matchID string, match *models.Match, inningsNumber int) error {

	// Get all innings for this match to determine current state
	innings, err := s.scorecardRepo.GetInningsByMatchID(ctx, matchID)
	if err != nil {
		// If no innings exist, this is the first ball of the match
		// The first innings should always be the toss-winning team
		if inningsNumber == 1 {
			// Check if the batting team matches the toss winner
			if match.BattingTeam != match.TossWinner {
				return fmt.Errorf("first innings must be played by the toss-winning team (%s), but current batting team is %s",
					match.TossWinner, match.BattingTeam)
			}
			return nil
		} else {
			return fmt.Errorf("cannot start with innings %d, first innings must be played first", inningsNumber)
		}
	}

	// Determine which innings exist
	firstInningsExists := false

	for _, inn := range innings {
		if inn.InningsNumber == 1 {
			firstInningsExists = true
		}
	}

	// Check if we're trying to add to first innings
	if inningsNumber == 1 {
		if !firstInningsExists {
			// First innings doesn't exist yet, check if batting team is correct
			if match.BattingTeam != match.TossWinner {
				return fmt.Errorf("first innings must be played by the toss-winning team (%s), but current batting team is %s",
					match.TossWinner, match.BattingTeam)
			}
			return nil
		} else {
			// First innings exists, check if it's complete
			firstInnings, err := s.scorecardRepo.GetInningsByMatchAndNumber(ctx, matchID, 1)
			if err != nil {
				return fmt.Errorf("failed to get first innings: %w", err)
			}

			// Check if first innings is complete (all wickets down or overs completed)
			firstInningsComplete := firstInnings.TotalWickets >= 10 || firstInnings.TotalOvers >= float64(match.TotalOvers)

			if !firstInningsComplete {
				// First innings is not complete, only toss winner can bat
				if match.BattingTeam != match.TossWinner {
					return fmt.Errorf("first innings is not complete, only toss-winning team (%s) can bat, but current batting team is %s",
						match.TossWinner, match.BattingTeam)
				}
			} else {
				// First innings is complete, but we're trying to add to first innings again
				return fmt.Errorf("first innings is complete, cannot add more balls to first innings")
			}
			return nil
		}
	}

	// Check if we're trying to add to second innings
	if inningsNumber == 2 {
		if !firstInningsExists {
			return fmt.Errorf("cannot start second innings, first innings must be played first")
		}

		// Check if first innings is complete
		firstInnings, err := s.scorecardRepo.GetInningsByMatchAndNumber(ctx, matchID, 1)
		if err != nil {
			return fmt.Errorf("failed to get first innings: %w", err)
		}

		// First innings is complete if all wickets are down or overs are completed
		firstInningsComplete := firstInnings.TotalWickets >= 10 || firstInnings.TotalOvers >= float64(match.TotalOvers)

		if !firstInningsComplete {
			return fmt.Errorf("first innings is not complete, cannot start second innings")
		}

		// Second innings should be played by the non-toss-winning team
		nonTossWinner := s.GetNonTossWinner(match.TossWinner)
		if match.BattingTeam != nonTossWinner {
			return fmt.Errorf("second innings must be played by the non-toss-winning team (%s), but current batting team is %s",
				nonTossWinner, match.BattingTeam)
		}

		return nil
	}

	return fmt.Errorf("invalid innings number: %d", inningsNumber)
}

// GetNonTossWinner returns the team that didn't win the toss
func (s *ScorecardService) GetNonTossWinner(tossWinner models.TeamType) models.TeamType {
	if tossWinner == models.TeamTypeA {
		return models.TeamTypeB
	}
	return models.TeamTypeA
}
