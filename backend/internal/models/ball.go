package models

import (
	"time"
)

// BallType represents the type of ball in cricket
type BallType string

const (
	BallTypeGood     BallType = "good"
	BallTypeWide     BallType = "wide"
	BallTypeNoBall   BallType = "no_ball"
	BallTypeDeadBall BallType = "dead_ball"
)

// Ball represents a ball event in cricket
type Ball struct {
	ID         string    `json:"id" db:"id"`
	OverID     string    `json:"over_id" db:"over_id"`
	BallNumber int       `json:"ball_number" db:"ball_number"`
	BallType   BallType  `json:"ball_type" db:"ball_type"`
	Runs       int       `json:"runs" db:"runs"`
	IsWicket   bool      `json:"is_wicket" db:"is_wicket"`
	BatsmanID  string    `json:"batsman_id" db:"batsman_id"`
	BowlerID   string    `json:"bowler_id" db:"bowler_id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// BallEvent represents a ball event request
type BallEvent struct {
	BallType  BallType `json:"ball_type" validate:"required,oneof=good wide no_ball dead_ball"`
	Runs      int      `json:"runs" validate:"required,min=0,max=6"`
	IsWicket  bool     `json:"is_wicket"`
	BatsmanID string   `json:"batsman_id" validate:"required"`
	BowlerID  string   `json:"bowler_id" validate:"required"`
}
