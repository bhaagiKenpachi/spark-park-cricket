package models

import (
	"time"
)

// Team represents a cricket team
type Team struct {
	ID           string    `json:"id,omitempty" db:"id,omitempty"`
	Name         string    `json:"name" db:"name"`
	PlayersCount int       `json:"players_count" db:"players_count"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// CreateTeamRequest represents the request to create a new team
type CreateTeamRequest struct {
	Name         string `json:"name" validate:"required,min=3,max=255"`
	PlayersCount int    `json:"players_count" validate:"required,min=1,max=20"`
}

// CreateTeamData represents the data structure for team creation (without ID)
type CreateTeamData struct {
	Name         string    `json:"name" db:"name"`
	PlayersCount int       `json:"players_count" db:"players_count"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// UpdateTeamRequest represents the request to update a team
type UpdateTeamRequest struct {
	Name         *string `json:"name,omitempty" validate:"omitempty,min=3,max=255"`
	PlayersCount *int    `json:"players_count,omitempty" validate:"omitempty,min=1,max=20"`
}

// TeamFilters represents filters for listing teams
type TeamFilters struct {
	Limit  int `json:"limit" validate:"min=1,max=100"`
	Offset int `json:"offset" validate:"min=0"`
}
