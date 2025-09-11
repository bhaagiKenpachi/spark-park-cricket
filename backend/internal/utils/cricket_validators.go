package utils

import (
	"fmt"
	"spark-park-cricket-backend/internal/models"
)

// CricketValidator provides cricket-specific validation functions
type CricketValidator struct{}

// NewCricketValidator creates a new cricket validator
func NewCricketValidator() *CricketValidator {
	return &CricketValidator{}
}

// ValidateBallEvent validates a ball event according to cricket rules
func (v *CricketValidator) ValidateBallEvent(ballEvent *models.BallEvent) error {
	// Validate ball type
	if !IsValidBallType(string(ballEvent.BallType)) {
		return fmt.Errorf("invalid ball type: %s", ballEvent.BallType)
	}

	// Validate run type based on ball type
	switch ballEvent.BallType {
	case models.BallTypeGood:
		if !ballEvent.RunType.IsValidRun() {
			return fmt.Errorf("invalid run type for good ball: %s", ballEvent.RunType)
		}
		if ballEvent.RunType == "5" {
			return fmt.Errorf("5 runs is not possible in cricket")
		}
	case models.BallTypeWide:
		if ballEvent.RunType != models.RunTypeWD && ballEvent.RunType.GetRunValue() < 1 {
			return fmt.Errorf("wide balls must have at least 1 run")
		}
	case models.BallTypeNoBall:
		if ballEvent.RunType != models.RunTypeNB && ballEvent.RunType.GetRunValue() < 1 {
			return fmt.Errorf("no balls must have at least 1 run")
		}
	case models.BallTypeDeadBall:
		if ballEvent.RunType.GetRunValue() != 0 {
			return fmt.Errorf("dead balls cannot have runs")
		}
	}

	// Validate wicket logic
	if ballEvent.IsWicket {
		if ballEvent.BallType == models.BallTypeWide || ballEvent.BallType == models.BallTypeNoBall {
			return fmt.Errorf("wickets cannot be taken on wide or no balls")
		}
		if ballEvent.RunType.GetRunValue() > 0 && ballEvent.IsWicket {
			return fmt.Errorf("wickets cannot be taken with runs on the same ball")
		}
	}

	return nil
}

// ValidateOverCompletion validates if an over is properly completed
func (v *CricketValidator) ValidateOverCompletion(over *models.Over) error {
	if over.TotalBalls < 0 || over.TotalBalls > 6 {
		return fmt.Errorf("over must have 0-6 balls")
	}

	if over.TotalRuns < 0 {
		return fmt.Errorf("over cannot have negative runs")
	}

	return nil
}

// ValidateScoreboard validates scoreboard data
func (v *CricketValidator) ValidateScoreboard(scoreboard *models.LiveScoreboard) error {
	if scoreboard.Score < 0 {
		return fmt.Errorf("score cannot be negative")
	}

	if scoreboard.Wickets < 0 || scoreboard.Wickets > 10 {
		return fmt.Errorf("wickets must be between 0 and 10")
	}

	if scoreboard.Overs < 0 {
		return fmt.Errorf("overs cannot be negative")
	}

	if scoreboard.Balls < 0 || scoreboard.Balls > 5 {
		return fmt.Errorf("balls must be between 0 and 5")
	}

	// Validate overs format (e.g., 15.3 means 15 overs and 3 balls)
	wholeOvers := int(scoreboard.Overs)
	fractionalPart := scoreboard.Overs - float64(wholeOvers)

	if fractionalPart > 0.5 {
		return fmt.Errorf("invalid overs format: fractional part cannot exceed 0.5")
	}

	return nil
}

// ValidateMatchStatus validates match status transitions
func (v *CricketValidator) ValidateMatchStatus(currentStatus, newStatus models.MatchStatus) error {
	validTransitions := map[models.MatchStatus][]models.MatchStatus{
		// Removed MatchStatusScheduled as matches are always live by default
		models.MatchStatusLive:      {models.MatchStatusCompleted, models.MatchStatusCancelled},
		models.MatchStatusCompleted: {}, // No transitions from completed
		models.MatchStatusCancelled: {}, // No transitions from cancelled
	}

	allowedTransitions, exists := validTransitions[currentStatus]
	if !exists {
		return fmt.Errorf("invalid current status: %s", currentStatus)
	}

	for _, allowed := range allowedTransitions {
		if allowed == newStatus {
			return nil
		}
	}

	return fmt.Errorf("invalid status transition from %s to %s", currentStatus, newStatus)
}

// ValidateTeamComposition validates team composition
func (v *CricketValidator) ValidateTeamComposition(team *models.Team, players []*models.Player) error {
	if team.PlayersCount < 1 || team.PlayersCount > 20 {
		return fmt.Errorf("team must have 1-20 players")
	}

	if len(players) > team.PlayersCount {
		return fmt.Errorf("team has more players than allowed")
	}

	// Check for duplicate player names in the same team
	playerNames := make(map[string]bool)
	for _, player := range players {
		if playerNames[player.Name] {
			return fmt.Errorf("duplicate player name: %s", player.Name)
		}
		playerNames[player.Name] = true
	}

	return nil
}

// ValidateMatchSetup validates match setup
func (v *CricketValidator) ValidateMatchSetup(match *models.Match) error {
	// Teams are now Team A and Team B, so no need to check for duplicates

	if match.MatchNumber < 1 {
		return fmt.Errorf("match number must be positive")
	}

	return nil
}

// CalculateRequiredRuns calculates runs required for a team to win
func (v *CricketValidator) CalculateRequiredRuns(targetScore, currentScore int) int {
	required := targetScore - currentScore
	if required < 0 {
		return 0
	}
	return required
}

// CalculateRequiredBalls calculates balls remaining in an innings
func (v *CricketValidator) CalculateRequiredBalls(totalOvers float64, currentOvers float64) int {
	totalBalls := int(totalOvers * 6)
	currentBalls := int(currentOvers*6) + int((currentOvers-float64(int(currentOvers)))*6)
	remaining := totalBalls - currentBalls
	if remaining < 0 {
		return 0
	}
	return remaining
}

// IsInningsComplete checks if an innings is complete
func (v *CricketValidator) IsInningsComplete(wickets int, overs float64, maxOvers float64) bool {
	// Innings complete if all wickets are down
	if wickets >= 10 {
		return true
	}

	// Innings complete if all overs are bowled
	if overs >= maxOvers {
		return true
	}

	return false
}

// ValidateBallEventRequest validates a ball event request for scorecard
func ValidateBallEventRequest(req *models.BallEventRequest) error {
	// Validate innings number
	if req.InningsNumber < 1 || req.InningsNumber > 2 {
		return fmt.Errorf("innings number must be 1 or 2")
	}

	// Validate ball type
	if !IsValidBallType(string(req.BallType)) {
		return fmt.Errorf("invalid ball type: %s", req.BallType)
	}

	// Validate run type
	if !req.RunType.IsValidRun() {
		return fmt.Errorf("invalid run type: %s", req.RunType)
	}

	// Validate byes (0-6)
	if req.Byes < 0 || req.Byes > 6 {
		return fmt.Errorf("byes must be between 0 and 6")
	}

	// Validate wicket logic
	if req.IsWicket {
		if req.BallType == models.BallTypeWide || req.BallType == models.BallTypeNoBall {
			return fmt.Errorf("wickets cannot be taken on wide or no balls")
		}
		if req.RunType == models.RunTypeWC {
			// WC (wicket) is valid
		} else if req.RunType.GetRunValue() > 0 {
			return fmt.Errorf("wickets cannot be taken with runs on the same ball")
		}
	}

	// Validate wicket type if wicket is taken
	if req.IsWicket && req.WicketType != "" {
		validWicketTypes := []string{"bowled", "caught", "lbw", "run_out", "stumped", "hit_wicket"}
		isValid := false
		for _, wt := range validWicketTypes {
			if req.WicketType == wt {
				isValid = true
				break
			}
		}
		if !isValid {
			return fmt.Errorf("invalid wicket type: %s", req.WicketType)
		}
	}

	return nil
}

// ValidateInnings validates innings data
func ValidateInnings(innings *models.Innings) error {
	if innings.InningsNumber < 1 || innings.InningsNumber > 2 {
		return fmt.Errorf("innings number must be 1 or 2")
	}

	if innings.TotalRuns < 0 {
		return fmt.Errorf("total runs cannot be negative")
	}

	if innings.TotalWickets < 0 || innings.TotalWickets > 10 {
		return fmt.Errorf("total wickets must be between 0 and 10")
	}

	if innings.TotalOvers < 0 {
		return fmt.Errorf("total overs cannot be negative")
	}

	if innings.TotalBalls < 0 {
		return fmt.Errorf("total balls cannot be negative")
	}

	if innings.Status != string(models.InningsStatusInProgress) && innings.Status != string(models.InningsStatusCompleted) {
		return fmt.Errorf("invalid innings status: %s", innings.Status)
	}

	return nil
}

// ValidateOver validates over data
func ValidateOver(over *models.ScorecardOver) error {
	if over.OverNumber < 1 {
		return fmt.Errorf("over number must be positive")
	}

	if over.TotalRuns < 0 {
		return fmt.Errorf("total runs cannot be negative")
	}

	if over.TotalBalls < 0 || over.TotalBalls > 6 {
		return fmt.Errorf("total balls must be between 0 and 6")
	}

	if over.TotalWickets < 0 {
		return fmt.Errorf("total wickets cannot be negative")
	}

	if over.Status != string(models.OverStatusInProgress) && over.Status != string(models.OverStatusCompleted) {
		return fmt.Errorf("invalid over status: %s", over.Status)
	}

	return nil
}

// ValidateBall validates ball data
func ValidateBall(ball *models.ScorecardBall) error {
	if ball.BallNumber < 1 || ball.BallNumber > 6 {
		return fmt.Errorf("ball number must be between 1 and 6")
	}

	if !IsValidBallType(string(ball.BallType)) {
		return fmt.Errorf("invalid ball type: %s", ball.BallType)
	}

	if !ball.RunType.IsValidRun() {
		return fmt.Errorf("invalid run type: %s", ball.RunType)
	}

	if ball.Runs < 0 {
		return fmt.Errorf("runs cannot be negative")
	}

	// Validate wicket type if wicket is taken
	if ball.IsWicket && ball.WicketType != "" {
		validWicketTypes := []string{"bowled", "caught", "lbw", "run_out", "stumped", "hit_wicket"}
		isValid := false
		for _, wt := range validWicketTypes {
			if ball.WicketType == wt {
				isValid = true
				break
			}
		}
		if !isValid {
			return fmt.Errorf("invalid wicket type: %s", ball.WicketType)
		}
	}

	return nil
}
