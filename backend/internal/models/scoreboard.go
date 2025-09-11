package models

import (
	"time"
)

// RunType represents the type of run scored
type RunType string

const (
	RunTypeZero  RunType = "0" // Dot Ball - 0 runs
	RunTypeOne   RunType = "1"
	RunTypeTwo   RunType = "2"
	RunTypeThree RunType = "3"
	RunTypeFour  RunType = "4"
	RunTypeFive  RunType = "5"
	RunTypeSix   RunType = "6"
	RunTypeSeven RunType = "7"
	RunTypeEight RunType = "8"
	RunTypeNine  RunType = "9"
	RunTypeNB    RunType = "NB" // No Ball
	RunTypeWD    RunType = "WD" // Wide
	RunTypeLB    RunType = "LB" // Leg Byes
	RunTypeWC    RunType = "WC" // Wicket
)

// BallType represents the type of ball bowled
type BallType string

const (
	BallTypeGood     BallType = "good"
	BallTypeWide     BallType = "wide"
	BallTypeNoBall   BallType = "no_ball"
	BallTypeDeadBall BallType = "dead_ball"
)

// ExtrasType represents the type of extras
type ExtrasType string

const (
	ExtrasTypeByes    ExtrasType = "byes"     // Byes
	ExtrasTypeLegByes ExtrasType = "leg_byes" // Leg byes
	ExtrasTypeWides   ExtrasType = "wides"    // Wides
	ExtrasTypeNoBalls ExtrasType = "no_balls" // No balls
)

// LiveScoreboard represents the live scoreboard for a match
type LiveScoreboard struct {
	ID          string    `json:"id,omitempty" db:"id,omitempty"`
	MatchID     string    `json:"match_id" db:"match_id"`
	BattingTeam TeamType  `json:"batting_team" db:"batting_team"`
	Score       int       `json:"score" db:"score"`
	Wickets     int       `json:"wickets" db:"wickets"`
	Overs       float64   `json:"overs" db:"overs"`
	Balls       int       `json:"balls" db:"balls"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// BallEvent represents a ball event in the match
type BallEvent struct {
	BallType  BallType  `json:"ball_type" validate:"required,oneof=good wide no_ball dead_ball"`
	RunType   RunType   `json:"run_type" validate:"required,oneof=1 2 3 4 5 6 7 8 9 NB WD LB"`
	IsWicket  bool      `json:"is_wicket"`
	CreatedAt time.Time `json:"created_at"`
}

// UpdateScoreRequest represents the request to update the score
type UpdateScoreRequest struct {
	Score int `json:"score" validate:"required,min=0"`
}

// UpdateWicketRequest represents the request to update wickets
type UpdateWicketRequest struct {
	Wickets int `json:"wickets" validate:"required,min=0,max=10"`
}

// GetRunValue returns the numeric value of a run type
func (rt RunType) GetRunValue() int {
	switch rt {
	case RunTypeZero:
		return 0
	case RunTypeOne:
		return 1
	case RunTypeTwo:
		return 2
	case RunTypeThree:
		return 3
	case RunTypeFour:
		return 4
	case RunTypeFive:
		return 5
	case RunTypeSix:
		return 6
	case RunTypeSeven:
		return 7
	case RunTypeEight:
		return 8
	case RunTypeNine:
		return 9
	case RunTypeNB, RunTypeWD, RunTypeLB:
		return 1
	case RunTypeWC:
		return 0
	default:
		return 0
	}
}

// IsValidRun returns true if the run type is valid
func (rt RunType) IsValidRun() bool {
	switch rt {
	case RunTypeZero, RunTypeOne, RunTypeTwo, RunTypeThree, RunTypeFour, RunTypeFive,
		RunTypeSix, RunTypeSeven, RunTypeEight, RunTypeNine, RunTypeNB, RunTypeWD, RunTypeLB, RunTypeWC:
		return true
	default:
		return false
	}
}
