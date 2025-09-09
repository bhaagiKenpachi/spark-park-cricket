package models

import (
	"time"
)

// MatchStatus represents the status of a cricket match
type MatchStatus string

const (
	MatchStatusScheduled MatchStatus = "scheduled"
	MatchStatusLive      MatchStatus = "live"
	MatchStatusCompleted MatchStatus = "completed"
	MatchStatusCancelled MatchStatus = "cancelled"
)

// Match represents a cricket match
type Match struct {
	ID          string      `json:"id" db:"id"`
	SeriesID    string      `json:"series_id" db:"series_id"`
	MatchNumber int         `json:"match_number" db:"match_number"`
	Date        time.Time   `json:"date" db:"date"`
	Status      MatchStatus `json:"status" db:"status"`
	Team1ID     string      `json:"team1_id" db:"team1_id"`
	Team2ID     string      `json:"team2_id" db:"team2_id"`
	CreatedAt   time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at" db:"updated_at"`
}

// CreateMatchRequest represents the request to create a new match
type CreateMatchRequest struct {
	SeriesID    string    `json:"series_id" validate:"required"`
	MatchNumber int       `json:"match_number" validate:"required,min=1"`
	Date        time.Time `json:"date" validate:"required"`
	Team1ID     string    `json:"team1_id" validate:"required"`
	Team2ID     string    `json:"team2_id" validate:"required,nefield=Team1ID"`
}

// UpdateMatchRequest represents the request to update a match
type UpdateMatchRequest struct {
	MatchNumber *int         `json:"match_number,omitempty" validate:"omitempty,min=1"`
	Date        *time.Time   `json:"date,omitempty"`
	Status      *MatchStatus `json:"status,omitempty" validate:"omitempty,oneof=scheduled live completed cancelled"`
	Team1ID     *string      `json:"team1_id,omitempty"`
	Team2ID     *string      `json:"team2_id,omitempty"`
}

// MatchFilters represents filters for listing matches
type MatchFilters struct {
	SeriesID *string      `json:"series_id,omitempty"`
	Status   *MatchStatus `json:"status,omitempty"`
	Limit    int          `json:"limit" validate:"min=1,max=100"`
	Offset   int          `json:"offset" validate:"min=0"`
}
