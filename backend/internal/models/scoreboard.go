package models

import (
	"time"
)

// LiveScoreboard represents the live scoreboard for a match
type LiveScoreboard struct {
	ID            string    `json:"id" db:"id"`
	MatchID       string    `json:"match_id" db:"match_id"`
	BattingTeamID string    `json:"batting_team_id" db:"batting_team_id"`
	Score         int       `json:"score" db:"score"`
	Wickets       int       `json:"wickets" db:"wickets"`
	Overs         float64   `json:"overs" db:"overs"`
	Balls         int       `json:"balls" db:"balls"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// UpdateScoreRequest represents the request to update match score
type UpdateScoreRequest struct {
	Score int `json:"score" validate:"required,min=0"`
}

// UpdateWicketRequest represents the request to update wickets
type UpdateWicketRequest struct {
	Wickets int `json:"wickets" validate:"required,min=0,max=10"`
}
