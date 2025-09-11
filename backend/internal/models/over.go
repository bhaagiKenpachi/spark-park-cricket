package models

import (
	"time"
)

// Over represents an over in a cricket match
type Over struct {
	ID          string    `json:"id,omitempty" db:"id,omitempty"`
	MatchID     string    `json:"match_id" db:"match_id"`
	OverNumber  int       `json:"over_number" db:"over_number"`
	BattingTeam TeamType  `json:"batting_team" db:"batting_team"`
	TotalRuns   int       `json:"total_runs" db:"total_runs"`
	TotalBalls  int       `json:"total_balls" db:"total_balls"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
