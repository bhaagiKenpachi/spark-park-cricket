package models

import (
	"time"
)

// MatchStatus represents the status of a cricket match
type MatchStatus string

const (
	MatchStatusLive      MatchStatus = "live"
	MatchStatusCompleted MatchStatus = "completed"
	MatchStatusCancelled MatchStatus = "cancelled"
)

// TossType represents the toss result
type TossType string

const (
	TossTypeHeads TossType = "H"
	TossTypeTails TossType = "T"
)

// TeamType represents the team type
type TeamType string

const (
	TeamTypeA TeamType = "A"
	TeamTypeB TeamType = "B"
)

// Match represents a cricket match
type Match struct {
	ID               string      `json:"id,omitempty" db:"id,omitempty"`
	SeriesID         string      `json:"series_id" db:"series_id"`
	MatchNumber      int         `json:"match_number" db:"match_number"`
	Date             time.Time   `json:"date" db:"date"`
	Status           MatchStatus `json:"status" db:"status"`
	TeamAPlayerCount int         `json:"team_a_player_count" db:"team_a_player_count"`
	TeamBPlayerCount int         `json:"team_b_player_count" db:"team_b_player_count"`
	TotalOvers       int         `json:"total_overs" db:"total_overs"`
	TossWinner       TeamType    `json:"toss_winner" db:"toss_winner"`
	TossType         TossType    `json:"toss_type" db:"toss_type"`
	BattingTeam      TeamType    `json:"batting_team" db:"batting_team"`
	CreatedAt        time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at" db:"updated_at"`
}

// CreateMatchRequest represents the request to create a new match
type CreateMatchRequest struct {
	SeriesID         string    `json:"series_id" validate:"required"`
	MatchNumber      *int      `json:"match_number,omitempty" validate:"omitempty,min=1"`
	Date             time.Time `json:"date" validate:"required"`
	TeamAPlayerCount int       `json:"team_a_player_count" validate:"required,min=1,max=20"`
	TeamBPlayerCount int       `json:"team_b_player_count" validate:"required,min=1,max=20"`
	TotalOvers       int       `json:"total_overs" validate:"required,min=1,max=20"`
	TossWinner       TeamType  `json:"toss_winner" validate:"required,oneof=A B"`
	TossType         TossType  `json:"toss_type" validate:"required,oneof=H T"`
}

// UpdateMatchRequest represents the request to update a match
type UpdateMatchRequest struct {
	MatchNumber      *int         `json:"match_number,omitempty" validate:"omitempty,min=1"`
	Date             *time.Time   `json:"date,omitempty"`
	Status           *MatchStatus `json:"status,omitempty" validate:"omitempty,oneof=live completed cancelled"`
	TeamAPlayerCount *int         `json:"team_a_player_count,omitempty" validate:"omitempty,min=1,max=20"`
	TeamBPlayerCount *int         `json:"team_b_player_count,omitempty" validate:"omitempty,min=1,max=20"`
	TotalOvers       *int         `json:"total_overs,omitempty" validate:"omitempty,min=1,max=20"`
	BattingTeam      *TeamType    `json:"batting_team,omitempty" validate:"omitempty,oneof=A B"`
}

// MatchFilters represents filters for listing matches
type MatchFilters struct {
	SeriesID *string      `json:"series_id,omitempty"`
	Status   *MatchStatus `json:"status,omitempty"`
	Limit    int          `json:"limit" validate:"min=1,max=100"`
	Offset   int          `json:"offset" validate:"min=0"`
}
