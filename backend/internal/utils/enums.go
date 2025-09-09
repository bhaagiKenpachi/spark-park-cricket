package utils

// MatchStatus represents the status of a cricket match
type MatchStatus string

const (
	MatchStatusScheduled MatchStatus = "scheduled"
	MatchStatusLive      MatchStatus = "live"
	MatchStatusCompleted MatchStatus = "completed"
	MatchStatusCancelled MatchStatus = "cancelled"
)

// BallType represents the type of ball in cricket
type BallType string

const (
	BallTypeGood     BallType = "good"
	BallTypeWide     BallType = "wide"
	BallTypeNoBall   BallType = "no_ball"
	BallTypeDeadBall BallType = "dead_ball"
)

// IsValidMatchStatus checks if the status is valid
func IsValidMatchStatus(status string) bool {
	switch MatchStatus(status) {
	case MatchStatusScheduled, MatchStatusLive, MatchStatusCompleted, MatchStatusCancelled:
		return true
	default:
		return false
	}
}

// IsValidBallType checks if the ball type is valid
func IsValidBallType(ballType string) bool {
	switch BallType(ballType) {
	case BallTypeGood, BallTypeWide, BallTypeNoBall, BallTypeDeadBall:
		return true
	default:
		return false
	}
}

// GetValidMatchStatuses returns all valid match statuses
func GetValidMatchStatuses() []string {
	return []string{
		string(MatchStatusScheduled),
		string(MatchStatusLive),
		string(MatchStatusCompleted),
		string(MatchStatusCancelled),
	}
}

// GetValidBallTypes returns all valid ball types
func GetValidBallTypes() []string {
	return []string{
		string(BallTypeGood),
		string(BallTypeWide),
		string(BallTypeNoBall),
		string(BallTypeDeadBall),
	}
}
