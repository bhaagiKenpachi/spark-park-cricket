package models

import (
	"time"
)

// Player represents a cricket player
type Player struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	TeamID    string    `json:"team_id" db:"team_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CreatePlayerRequest represents the request to create a new player
type CreatePlayerRequest struct {
	Name   string `json:"name" validate:"required,min=2,max=255"`
	TeamID string `json:"team_id" validate:"required"`
}

// UpdatePlayerRequest represents the request to update a player
type UpdatePlayerRequest struct {
	Name   *string `json:"name,omitempty" validate:"omitempty,min=2,max=255"`
	TeamID *string `json:"team_id,omitempty"`
}

// PlayerFilters represents filters for listing players
type PlayerFilters struct {
	TeamID *string `json:"team_id,omitempty"`
	Limit  int     `json:"limit" validate:"min=1,max=100"`
	Offset int     `json:"offset" validate:"min=0"`
}
