package models

import (
	"time"
)

// FallOfWickets represents a wicket fall event in a cricket match
type FallOfWickets struct {
	ID           string    `json:"id,omitempty" db:"id,omitempty"`
	MatchID      string    `json:"match_id" db:"match_id"`
	InningsID    string    `json:"innings_id" db:"innings_id"`
	OverID       string    `json:"over_id" db:"over_id"`
	BallID       string    `json:"ball_id" db:"ball_id"`
	WicketNumber int       `json:"wicket_number" db:"wicket_number"`
	Score        int       `json:"score" db:"score"`
	OverNumber   int       `json:"over_number" db:"over_number"`
	BallNumber   int       `json:"ball_number" db:"ball_number"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// CreateFallOfWicketsRequest represents the request to create a fall of wickets record
type CreateFallOfWicketsRequest struct {
	MatchID      string `json:"match_id" validate:"required"`
	InningsID    string `json:"innings_id" validate:"required"`
	OverID       string `json:"over_id" validate:"required"`
	BallID       string `json:"ball_id" validate:"required"`
	WicketNumber int    `json:"wicket_number" validate:"required,min=1,max=20"`
	Score        int    `json:"score" validate:"min=0"`
	OverNumber   int    `json:"over_number" validate:"required,min=1"`
	BallNumber   int    `json:"ball_number" validate:"required,min=1,max=20"`
}

// FallOfWicketsFilters represents filters for listing fall of wickets
type FallOfWicketsFilters struct {
	MatchID   *string `json:"match_id,omitempty"`
	InningsID *string `json:"innings_id,omitempty"`
	Limit     int     `json:"limit" validate:"min=1,max=100"`
	Offset    int     `json:"offset" validate:"min=0"`
}

// FallOfWicketsSummary represents a summary of fall of wickets for display
type FallOfWicketsSummary struct {
	MatchID      string       `json:"match_id"`
	InningsID    string       `json:"innings_id"`
	TotalWickets int          `json:"total_wickets"`
	Wickets      []WicketFall `json:"wickets"`
}

// WicketFall represents a single wicket fall entry
type WicketFall struct {
	WicketNumber int    `json:"wicket_number"`
	Score        int    `json:"score"`
	OverNumber   int    `json:"over_number"`
	BallNumber   int    `json:"ball_number"`
	OverPosition string `json:"over_position"` // e.g., "15.3" for over 15, ball 3
}

// UpdateFallOfWicketsRequest represents the request to update a fall of wickets record
type UpdateFallOfWicketsRequest struct {
	Score *int `json:"score,omitempty" validate:"omitempty,min=0"`
}
