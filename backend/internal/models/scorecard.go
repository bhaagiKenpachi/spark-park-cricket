package models

import (
	"time"
)

// Innings represents a cricket innings
type Innings struct {
	ID            string    `json:"id" db:"id"`
	MatchID       string    `json:"match_id" db:"match_id"`
	InningsNumber int       `json:"innings_number" db:"innings_number"`
	BattingTeam   TeamType  `json:"batting_team" db:"batting_team"`
	TotalRuns     int       `json:"total_runs" db:"total_runs"`
	TotalWickets  int       `json:"total_wickets" db:"total_wickets"`
	TotalOvers    float64   `json:"total_overs" db:"total_overs"`
	TotalBalls    int       `json:"total_balls" db:"total_balls"`
	Status        string    `json:"status" db:"status"` // "in_progress", "completed"
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// ScorecardOver represents a cricket over in scorecard
type ScorecardOver struct {
	ID           string    `json:"id" db:"id"`
	InningsID    string    `json:"innings_id" db:"innings_id"`
	OverNumber   int       `json:"over_number" db:"over_number"`
	TotalRuns    int       `json:"total_runs" db:"total_runs"`
	TotalBalls   int       `json:"total_balls" db:"total_balls"`
	TotalWickets int       `json:"total_wickets" db:"total_wickets"`
	Status       string    `json:"status" db:"status"` // "in_progress", "completed"
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// ScorecardBall represents a cricket ball in scorecard
type ScorecardBall struct {
	ID         string    `json:"id" db:"id"`
	OverID     string    `json:"over_id" db:"over_id"`
	BallNumber int       `json:"ball_number" db:"ball_number"`
	BallType   BallType  `json:"ball_type" db:"ball_type"`
	RunType    RunType   `json:"run_type" db:"run_type"`
	Runs       int       `json:"runs" db:"runs"`
	Byes       int       `json:"byes" db:"byes"` // Additional runs from byes
	IsWicket   bool      `json:"is_wicket" db:"is_wicket"`
	WicketType string    `json:"wicket_type,omitempty" db:"wicket_type"` // "bowled", "caught", "lbw", "run_out", "stumped", "hit_wicket"
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// ScorecardRequest represents the request to start scoring
type ScorecardRequest struct {
	MatchID string `json:"match_id" validate:"required,uuid"`
}

// BallEventRequest represents a ball event
type BallEventRequest struct {
	MatchID       string   `json:"match_id" validate:"required,uuid"`
	InningsNumber int      `json:"innings_number" validate:"required,min=1,max=2"`
	BallType      BallType `json:"ball_type" validate:"required"`
	RunType       RunType  `json:"run_type" validate:"required"`
	IsWicket      bool     `json:"is_wicket"`
	WicketType    string   `json:"wicket_type,omitempty"`
	Byes          int      `json:"byes,omitempty"` // Additional runs from byes (0-6)
}

// ScorecardResponse represents the complete scorecard
type ScorecardResponse struct {
	MatchID        string           `json:"match_id"`
	MatchNumber    int              `json:"match_number"`
	SeriesName     string           `json:"series_name"`
	TeamA          string           `json:"team_a"`
	TeamB          string           `json:"team_b"`
	TotalOvers     int              `json:"total_overs"`
	TossWinner     TeamType         `json:"toss_winner"`
	TossType       TossType         `json:"toss_type"`
	CurrentInnings int              `json:"current_innings"`
	Innings        []InningsSummary `json:"innings"`
	MatchStatus    string           `json:"match_status"`
}

// ExtrasSummary represents extras in an innings
type ExtrasSummary struct {
	Byes    int `json:"byes"`     // Byes
	LegByes int `json:"leg_byes"` // Leg byes
	Wides   int `json:"wides"`    // Wides
	NoBalls int `json:"no_balls"` // No balls
	Total   int `json:"total"`    // Total extras
}

// InningsSummary represents a summary of an innings
type InningsSummary struct {
	InningsNumber int            `json:"innings_number"`
	BattingTeam   TeamType       `json:"batting_team"`
	TotalRuns     int            `json:"total_runs"`
	TotalWickets  int            `json:"total_wickets"`
	TotalOvers    float64        `json:"total_overs"`
	TotalBalls    int            `json:"total_balls"`
	Status        string         `json:"status"`
	Extras        *ExtrasSummary `json:"extras"`
	Overs         []OverSummary  `json:"overs"`
}

// OverSummary represents a summary of an over
type OverSummary struct {
	OverNumber   int           `json:"over_number"`
	TotalRuns    int           `json:"total_runs"`
	TotalBalls   int           `json:"total_balls"`
	TotalWickets int           `json:"total_wickets"`
	Status       string        `json:"status"`
	Balls        []BallSummary `json:"balls"`
}

// BallSummary represents a summary of a ball
type BallSummary struct {
	BallNumber int      `json:"ball_number"`
	BallType   BallType `json:"ball_type"`
	RunType    RunType  `json:"run_type"`
	Runs       int      `json:"runs"`
	Byes       int      `json:"byes"`
	IsWicket   bool     `json:"is_wicket"`
	WicketType string   `json:"wicket_type,omitempty"`
}

// WicketType represents different types of wickets
type WicketType string

const (
	WicketTypeBowled    WicketType = "bowled"
	WicketTypeCaught    WicketType = "caught"
	WicketTypeLBW       WicketType = "lbw"
	WicketTypeRunOut    WicketType = "run_out"
	WicketTypeStumped   WicketType = "stumped"
	WicketTypeHitWicket WicketType = "hit_wicket"
)

// InningsStatus represents the status of an innings
type InningsStatus string

const (
	InningsStatusInProgress InningsStatus = "in_progress"
	InningsStatusCompleted  InningsStatus = "completed"
)

// OverStatus represents the status of an over
type OverStatus string

const (
	OverStatusInProgress OverStatus = "in_progress"
	OverStatusCompleted  OverStatus = "completed"
)
