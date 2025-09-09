package models

import (
	"time"
)

// Over represents an over in a cricket match
type Over struct {
	ID            string    `json:"id" db:"id"`
	MatchID       string    `json:"match_id" db:"match_id"`
	OverNumber    int       `json:"over_number" db:"over_number"`
	BattingTeamID string    `json:"batting_team_id" db:"batting_team_id"`
	TotalRuns     int       `json:"total_runs" db:"total_runs"`
	TotalBalls    int       `json:"total_balls" db:"total_balls"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}
