package models

import "time"

// VoteTeam represents a team associated with a vote
type VoteTeam struct {
	ID         string    `json:"id" db:"id"`
	VoteID     string    `json:"vote_id" db:"vote_id"`
	TeamName   string    `json:"team_name" db:"team_name"`
	TeamLetter string    `json:"team_letter" db:"team_letter"` // 'A' or 'B'
	CaptainID  string    `json:"captain_id" db:"captain_id"`   // Must be a voter
	CreatedBy  string    `json:"created_by" db:"created_by"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// TeamPlayer represents a player assigned to a team
type TeamPlayer struct {
	ID        string    `json:"id" db:"id"`
	TeamID    string    `json:"team_id" db:"team_id"`
	UserID    string    `json:"user_id" db:"user_id"` // Must be a voter
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// VoteTeamWithPlayers represents a team with its players and captain info
type VoteTeamWithPlayers struct {
	VoteTeam
	Captain     *User    `json:"captain,omitempty"`
	Players     []*User  `json:"players,omitempty"`
	PlayerCount int      `json:"player_count"`
	CaptainName string   `json:"captain_name,omitempty"`
	PlayerNames []string `json:"player_names,omitempty"`
}

// CreateVoteTeamRequest represents a request to create a vote team
type CreateVoteTeamRequest struct {
	VoteID     string `json:"vote_id" binding:"required"`
	TeamName   string `json:"team_name" binding:"required,min=2,max=100"`
	TeamLetter string `json:"team_letter" binding:"required,oneof=A B"`
	CaptainID  string `json:"captain_id" binding:"required,uuid"`
}

// UpdateVoteTeamRequest represents a request to update vote team details
type UpdateVoteTeamRequest struct {
	TeamName  *string `json:"team_name,omitempty" binding:"omitempty,min=2,max=100"`
	CaptainID *string `json:"captain_id,omitempty" binding:"omitempty,uuid"`
}

// AddPlayerRequest represents a request to add a player to a team
type AddPlayerRequest struct {
	UserID string `json:"user_id" binding:"required,uuid"`
}

// TeamAssignmentRequest represents a request to assign multiple players
type TeamAssignmentRequest struct {
	UserIDs []string `json:"user_ids" binding:"required,min=1,max=20,dive,uuid"`
}
