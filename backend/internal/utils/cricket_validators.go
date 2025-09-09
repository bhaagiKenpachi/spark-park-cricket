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

	// Validate runs based on ball type
	switch ballEvent.BallType {
	case models.BallTypeGood:
		if ballEvent.Runs < 0 || ballEvent.Runs > 6 {
			return fmt.Errorf("good balls can only have 0-6 runs")
		}
		if ballEvent.Runs == 5 {
			return fmt.Errorf("5 runs is not possible in cricket")
		}
	case models.BallTypeWide:
		if ballEvent.Runs < 1 {
			return fmt.Errorf("wide balls must have at least 1 run")
		}
		if ballEvent.Runs > 7 {
			return fmt.Errorf("wide balls cannot have more than 7 runs")
		}
	case models.BallTypeNoBall:
		if ballEvent.Runs < 1 {
			return fmt.Errorf("no balls must have at least 1 run")
		}
		if ballEvent.Runs > 7 {
			return fmt.Errorf("no balls cannot have more than 7 runs")
		}
	case models.BallTypeDeadBall:
		if ballEvent.Runs != 0 {
			return fmt.Errorf("dead balls cannot have runs")
		}
	}

	// Validate wicket logic
	if ballEvent.IsWicket {
		if ballEvent.BallType == models.BallTypeWide || ballEvent.BallType == models.BallTypeNoBall {
			return fmt.Errorf("wickets cannot be taken on wide or no balls")
		}
		if ballEvent.Runs > 0 && ballEvent.IsWicket {
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
		models.MatchStatusScheduled: {models.MatchStatusLive, models.MatchStatusCancelled},
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
	if match.Team1ID == match.Team2ID {
		return fmt.Errorf("team1 and team2 must be different")
	}

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
