package models

import (
	"time"
)

// Ball represents a ball in a cricket match
type Ball struct {
	ID         string    `json:"id,omitempty" db:"id,omitempty"`
	OverID     string    `json:"over_id" db:"over_id"`
	BallNumber int       `json:"ball_number" db:"ball_number"`
	BallType   BallType  `json:"ball_type" db:"ball_type"`
	RunType    RunType   `json:"run_type" db:"run_type"`
	IsWicket   bool      `json:"is_wicket" db:"is_wicket"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}
