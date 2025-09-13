package services

import (
	"context"
	"fmt"
	"log"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/repository/interfaces"
	"time"
)

// ScoreboardService handles live scoring operations
type ScoreboardService struct {
	scoreboardRepo interfaces.ScoreboardRepository
	overRepo       interfaces.OverRepository
	ballRepo       interfaces.BallRepository
	matchRepo      interfaces.MatchRepository
	teamRepo       interfaces.TeamRepository
	playerRepo     interfaces.PlayerRepository
}

// NewScoreboardService creates a new scoreboard service
func NewScoreboardService(
	scoreboardRepo interfaces.ScoreboardRepository,
	overRepo interfaces.OverRepository,
	ballRepo interfaces.BallRepository,
	matchRepo interfaces.MatchRepository,
	teamRepo interfaces.TeamRepository,
	playerRepo interfaces.PlayerRepository,
) *ScoreboardService {
	return &ScoreboardService{
		scoreboardRepo: scoreboardRepo,
		overRepo:       overRepo,
		ballRepo:       ballRepo,
		matchRepo:      matchRepo,
		teamRepo:       teamRepo,
		playerRepo:     playerRepo,
	}
}

// GetScoreboard retrieves the live scoreboard for a match
func (s *ScoreboardService) GetScoreboard(ctx context.Context, matchID string) (*models.LiveScoreboard, error) {

	if matchID == "" {
		return nil, fmt.Errorf("match ID is required")
	}

	// Check if match exists
	match, err := s.matchRepo.GetByID(ctx, matchID)
	if err != nil {
		return nil, fmt.Errorf("match not found: %w", err)
	}

	// Get scoreboard
	scoreboard, err := s.scoreboardRepo.GetByMatchID(ctx, matchID)
	if err != nil {
		// If no scoreboard exists, create one
		if err.Error() == "scoreboard not found" {
			return s.initializeScoreboard(ctx, match)
		}
		return nil, fmt.Errorf("failed to get scoreboard: %w", err)
	}

	return scoreboard, nil
}

// AddBall adds a ball event to the match
func (s *ScoreboardService) AddBall(ctx context.Context, matchID string, ballEvent *models.BallEvent) (*models.LiveScoreboard, error) {
	if matchID == "" {
		return nil, fmt.Errorf("match ID is required")
	}

	// Validate ball event
	if err := s.validateBallEvent(ballEvent); err != nil {
		return nil, fmt.Errorf("invalid ball event: %w", err)
	}

	// Get match details
	match, err := s.matchRepo.GetByID(ctx, matchID)
	if err != nil {
		return nil, fmt.Errorf("match not found: %w", err)
	}

	// Validate innings order
	if err := s.ValidateInningsOrder(ctx, matchID, match); err != nil {
		return nil, fmt.Errorf("innings validation failed: %w", err)
	}

	// Get current scoreboard
	scoreboard, err := s.GetScoreboard(ctx, matchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get scoreboard: %w", err)
	}

	// Get current over or create new one
	over, err := s.getCurrentOver(ctx, matchID, string(scoreboard.BattingTeam))
	if err != nil {
		return nil, fmt.Errorf("failed to get current over: %w", err)
	}

	// Create ball record
	ball := &models.Ball{
		OverID:     over.ID,
		BallNumber: s.getNextBallNumber(ctx, over.ID),
		BallType:   ballEvent.BallType,
		RunType:    ballEvent.RunType,
		IsWicket:   ballEvent.IsWicket,
		CreatedAt:  time.Now(),
	}

	// Save ball
	err = s.ballRepo.Create(ctx, ball)
	if err != nil {
		return nil, fmt.Errorf("failed to create ball: %w", err)
	}

	// Update scoreboard
	err = s.updateScoreboardWithBall(ctx, scoreboard, ball)
	if err != nil {
		return nil, fmt.Errorf("failed to update scoreboard: %w", err)
	}

	// Check if over is complete
	if s.isOverComplete(over) {
		err = s.completeOver(ctx, over)
		if err != nil {
			return nil, fmt.Errorf("failed to complete over: %w", err)
		}
	}

	return scoreboard, nil
}

// UpdateScore updates the match score manually
func (s *ScoreboardService) UpdateScore(ctx context.Context, matchID string, req *models.UpdateScoreRequest) (*models.LiveScoreboard, error) {
	if matchID == "" {
		return nil, fmt.Errorf("match ID is required")
	}

	scoreboard, err := s.GetScoreboard(ctx, matchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get scoreboard: %w", err)
	}

	scoreboard.Score = req.Score
	scoreboard.UpdatedAt = time.Now()

	err = s.scoreboardRepo.Update(ctx, matchID, scoreboard)
	if err != nil {
		return nil, fmt.Errorf("failed to update score: %w", err)
	}

	return scoreboard, nil
}

// UpdateWicket updates the wicket count
func (s *ScoreboardService) UpdateWicket(ctx context.Context, matchID string, req *models.UpdateWicketRequest) (*models.LiveScoreboard, error) {
	if matchID == "" {
		return nil, fmt.Errorf("match ID is required")
	}

	if req.Wickets < 0 || req.Wickets > 10 {
		return nil, fmt.Errorf("wickets must be between 0 and 10")
	}

	scoreboard, err := s.GetScoreboard(ctx, matchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get scoreboard: %w", err)
	}

	scoreboard.Wickets = req.Wickets
	scoreboard.UpdatedAt = time.Now()

	err = s.scoreboardRepo.Update(ctx, matchID, scoreboard)
	if err != nil {
		return nil, fmt.Errorf("failed to update wickets: %w", err)
	}

	return scoreboard, nil
}

// initializeScoreboard creates a new scoreboard for a match
func (s *ScoreboardService) initializeScoreboard(ctx context.Context, match *models.Match) (*models.LiveScoreboard, error) {

	// Ensure the first innings always starts with the toss-winning team
	scoreboard := &models.LiveScoreboard{
		MatchID:     match.ID,
		BattingTeam: match.TossWinner, // Always start with toss-winning team for first innings
		Score:       0,
		Wickets:     0,
		Overs:       0.0,
		Balls:       0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}


	err := s.scoreboardRepo.Create(ctx, scoreboard)
	if err != nil {
		return nil, fmt.Errorf("failed to create scoreboard: %w", err)
	}

	return scoreboard, nil
}

// validateBallEvent validates a ball event
func (s *ScoreboardService) validateBallEvent(ballEvent *models.BallEvent) error {
	runValue := ballEvent.RunType.GetRunValue()
	if runValue < 0 || runValue > 6 {
		return fmt.Errorf("runs must be between 0 and 6")
	}

	if ballEvent.BallType == models.BallTypeGood && runValue > 6 {
		return fmt.Errorf("good balls cannot have more than 6 runs")
	}

	if ballEvent.BallType == models.BallTypeWide && runValue < 1 {
		return fmt.Errorf("wide balls must have at least 1 run")
	}

	if ballEvent.BallType == models.BallTypeNoBall && runValue < 1 {
		return fmt.Errorf("no balls must have at least 1 run")
	}

	return nil
}

// ValidateInningsOrder validates that balls can only be added to the correct team's innings
func (s *ScoreboardService) ValidateInningsOrder(ctx context.Context, matchID string, match *models.Match) error {
	// Get all overs for this match to determine current innings
	overs, err := s.overRepo.GetByMatchID(ctx, matchID)
	if err != nil {
		return fmt.Errorf("failed to get overs: %w", err)
	}

	// If no overs exist, this is the first ball of the match
	// The first innings should always be the toss-winning team
	if len(overs) == 0 {
		if match.BattingTeam != match.TossWinner {
			return fmt.Errorf("first innings must be played by the toss-winning team (%s), but current batting team is %s",
				match.TossWinner, match.BattingTeam)
		}
		return nil
	}

	// Determine which innings we're in based on overs
	firstInningsOvers := 0
	secondInningsOvers := 0

	for _, over := range overs {
		if over.BattingTeam == match.TossWinner {
			firstInningsOvers++
		} else {
			secondInningsOvers++
		}
	}

	// Check if we're in first innings
	if firstInningsOvers > 0 && secondInningsOvers == 0 {
		// Check if first innings is complete (all wickets down or overs completed)
		scoreboard, err := s.scoreboardRepo.GetByMatchID(ctx, matchID)
		if err != nil {
			return fmt.Errorf("failed to get scoreboard for innings validation: %w", err)
		}

		// First innings is complete if all wickets are down or overs are completed
		firstInningsComplete := scoreboard.Wickets >= 10 || scoreboard.Overs >= float64(match.TotalOvers)

		if !firstInningsComplete {
			// First innings is not complete, only toss winner can bat
			if match.BattingTeam != match.TossWinner {
				return fmt.Errorf("first innings is not complete, only toss-winning team (%s) can bat, but current batting team is %s",
					match.TossWinner, match.BattingTeam)
			}
		} else {
			// First innings is complete, second innings can start
			// Second innings should be the non-toss-winning team
			nonTossWinner := s.GetNonTossWinner(match.TossWinner)
			if match.BattingTeam != nonTossWinner {
				return fmt.Errorf("first innings is complete, second innings must be played by the non-toss-winning team (%s), but current batting team is %s",
					nonTossWinner, match.BattingTeam)
			}
		}
		return nil
	}

	// Check if we're in second innings
	if firstInningsOvers > 0 && secondInningsOvers > 0 {
		// Second innings is in progress - only non-toss winner can bat
		nonTossWinner := s.GetNonTossWinner(match.TossWinner)
		if match.BattingTeam != nonTossWinner {
			return fmt.Errorf("second innings is in progress, only non-toss-winning team (%s) can bat, but current batting team is %s",
				nonTossWinner, match.BattingTeam)
		}
		return nil
	}

	return nil
}

// GetNonTossWinner returns the team that didn't win the toss
func (s *ScoreboardService) GetNonTossWinner(tossWinner models.TeamType) models.TeamType {
	if tossWinner == models.TeamTypeA {
		return models.TeamTypeB
	}
	return models.TeamTypeA
}

// getCurrentOver gets the current over or creates a new one
func (s *ScoreboardService) getCurrentOver(ctx context.Context, matchID, battingTeamID string) (*models.Over, error) {
	// Get all overs for the match
	overs, err := s.overRepo.GetByMatchID(ctx, matchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get overs: %w", err)
	}

	// Find the current over (last incomplete over)
	for i := len(overs) - 1; i >= 0; i-- {
		if overs[i].TotalBalls < 6 {
			return overs[i], nil
		}
	}

	// Create new over
	overNumber := len(overs) + 1
	over := &models.Over{
		MatchID:     matchID,
		OverNumber:  overNumber,
		BattingTeam: models.TeamType(battingTeamID),
		TotalRuns:   0,
		TotalBalls:  0,
		CreatedAt:   time.Now(),
	}

	err = s.overRepo.Create(ctx, over)
	if err != nil {
		return nil, fmt.Errorf("failed to create over: %w", err)
	}

	return over, nil
}

// getNextBallNumber gets the next ball number for an over
func (s *ScoreboardService) getNextBallNumber(ctx context.Context, overID string) int {
	balls, err := s.ballRepo.GetByOverID(ctx, overID)
	if err != nil {
		return 1 // First ball
	}

	return len(balls) + 1
}

// updateScoreboardWithBall updates the scoreboard with a new ball
func (s *ScoreboardService) updateScoreboardWithBall(ctx context.Context, scoreboard *models.LiveScoreboard, ball *models.Ball) error {
	// Update score
	scoreboard.Score += ball.RunType.GetRunValue()

	// Update wickets
	if ball.IsWicket {
		scoreboard.Wickets++
	}

	// Update balls and overs
	if ball.BallType == models.BallTypeGood {
		scoreboard.Balls++
		if scoreboard.Balls == 6 {
			scoreboard.Overs += 1.0
			scoreboard.Balls = 0
		}
	}
	// Wide and no balls don't count as balls

	scoreboard.UpdatedAt = time.Now()

	return s.scoreboardRepo.Update(ctx, scoreboard.MatchID, scoreboard)
}

// isOverComplete checks if an over is complete
func (s *ScoreboardService) isOverComplete(over *models.Over) bool {
	return over.TotalBalls >= 6
}

// completeOver marks an over as complete
func (s *ScoreboardService) completeOver(ctx context.Context, over *models.Over) error {
	// Get all balls for this over
	balls, err := s.ballRepo.GetByOverID(ctx, over.ID)
	if err != nil {
		return fmt.Errorf("failed to get balls for over: %w", err)
	}

	// Calculate total runs and balls
	totalRuns := 0
	totalBalls := 0

	for _, ball := range balls {
		totalRuns += ball.RunType.GetRunValue()
		if ball.BallType == models.BallTypeGood {
			totalBalls++
		}
	}

	over.TotalRuns = totalRuns
	over.TotalBalls = totalBalls

	return s.overRepo.Update(ctx, over.ID, over)
}
