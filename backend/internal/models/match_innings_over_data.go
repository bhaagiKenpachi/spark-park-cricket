package models

import "time"

// MatchInningsOverData represents optimized data structure for add ball API
// This combines match, innings, over, and ball data in a single query result
type MatchInningsOverData struct {
	// Match data
	MatchID          string      `json:"match_id"`
	MatchStatus      MatchStatus `json:"match_status"`
	CreatedBy        string      `json:"created_by"`
	BattingTeam      TeamType    `json:"batting_team"`
	TotalOvers       int         `json:"total_overs"`
	TeamAPlayerCount int         `json:"team_a_player_count"`

	// Innings data
	InningsID           string        `json:"innings_id"`
	InningsNumber       int           `json:"innings_number"`
	InningsStatus       InningsStatus `json:"innings_status"`
	InningsTotalRuns    int           `json:"innings_total_runs"`
	InningsTotalWickets int           `json:"innings_total_wickets"`
	InningsTotalOvers   float64       `json:"innings_total_overs"`
	InningsTotalBalls   int           `json:"innings_total_balls"`
	InningsStartTime    *time.Time    `json:"innings_start_time"`
	InningsEndTime      *time.Time    `json:"innings_end_time"`
	InningsDuration     int           `json:"innings_duration_seconds"`

	// Over data (current over)
	OverID           string     `json:"over_id"`
	OverNumber       int        `json:"over_number"`
	OverStatus       OverStatus `json:"over_status"`
	OverTotalRuns    int        `json:"over_total_runs"`
	OverTotalBalls   int        `json:"over_total_balls"`
	OverTotalWickets int        `json:"over_total_wickets"`
	OverStartTime    *time.Time `json:"over_start_time"`
	OverEndTime      *time.Time `json:"over_end_time"`
	OverDuration     int        `json:"over_duration_seconds"`

	// Ball data (for next ball number calculation)
	BallCount      int `json:"ball_count"`
	LegalBallCount int `json:"legal_ball_count"`
	MaxBallNumber  int `json:"max_ball_number"`

	// Overs data (for overs calculation)
	CompletedOvers   int `json:"completed_overs"`
	CurrentOverBalls int `json:"current_over_balls"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
