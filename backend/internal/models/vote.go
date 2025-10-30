package models

import (
	"time"
)

// VoteType represents the type of voting (single or multiple choice)
type VoteType string

const (
	VoteTypeSingle   VoteType = "single"   // Only one option can be selected
	VoteTypeMultiple VoteType = "multiple" // Multiple options can be selected
)

// VoteStatus represents the status of a vote
type VoteStatus string

const (
	VoteStatusActive    VoteStatus = "active"    // Vote is currently active
	VoteStatusClosed    VoteStatus = "closed"    // Vote is closed
	VoteStatusCancelled VoteStatus = "cancelled" // Vote is cancelled
)

// Vote represents a voting poll
type Vote struct {
	ID                   string     `json:"id" db:"id"`
	Title                string     `json:"title" db:"title"`
	Description          string     `json:"description" db:"description"`
	Type                 VoteType   `json:"type" db:"type"`
	Status               VoteStatus `json:"status" db:"status"`
	CreatedBy            string     `json:"created_by" db:"created_by"`
	TeamFormationEnabled bool       `json:"team_formation_enabled" db:"team_formation_enabled"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at" db:"updated_at"`
	ClosedAt             *time.Time `json:"closed_at,omitempty" db:"closed_at"`
}

// VoteWithCreator represents a vote with creator information
type VoteWithCreator struct {
	Vote
	CreatorName string `json:"creator_name"`
}

// VoteOption represents an option within a vote
type VoteOption struct {
	ID          string    `json:"id" db:"id"`
	VoteID      string    `json:"vote_id" db:"vote_id"`
	Text        string    `json:"text" db:"text"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// UserVote represents a user's vote selection
type UserVote struct {
	ID              string    `json:"id" db:"id"`
	VoteID          string    `json:"vote_id" db:"vote_id"`
	UserID          string    `json:"user_id" db:"user_id"`
	SelectedOptions []string  `json:"selected_options" db:"selected_options"` // Array of option IDs
	VotedAt         time.Time `json:"voted_at" db:"voted_at"`
}

// VoteWithOptions represents a vote with its options
type VoteWithOptions struct {
	Vote        *Vote         `json:"vote"`
	Options     []*VoteOption `json:"options"`
	CreatorName string        `json:"creator_name"`
}

// VoteWithResults represents a vote with results and user's vote
type VoteWithResults struct {
	Vote             *Vote                  `json:"vote"`
	Options          []*VoteOption          `json:"options"`
	Results          map[string]int         `json:"results"`            // option_id -> vote_count
	ResultsWithNames map[string][]VoterInfo `json:"results_with_names"` // option_id -> list of voters
	UserVote         *UserVote              `json:"user_vote,omitempty"`
	TotalVotes       int                    `json:"total_votes"`
	VotedUsers       []string               `json:"voted_users"` // List of user IDs who voted
	CreatorName      string                 `json:"creator_name"`
}

// VoterInfo represents voter information for a vote option
type VoterInfo struct {
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
	VotedAt  string `json:"voted_at"`
}

// CreateVoteRequest represents the request to create a new vote
type CreateVoteRequest struct {
	Title                string   `json:"title" validate:"required,min=3,max=255"`
	Description          string   `json:"description" validate:"required,min=10,max=1000"`
	Type                 VoteType `json:"type" validate:"required,oneof=single multiple"`
	Options              []string `json:"options" validate:"required,min=2,dive,required,min=1,max=255"`
	TeamFormationEnabled bool     `json:"team_formation_enabled"`
}

// UpdateVoteRequest represents the request to update a vote
type UpdateVoteRequest struct {
	Title                *string     `json:"title,omitempty" validate:"omitempty,min=3,max=255"`
	Description          *string     `json:"description,omitempty" validate:"omitempty,min=10,max=1000"`
	Status               *VoteStatus `json:"status,omitempty" validate:"omitempty,oneof=active closed cancelled"`
	TeamFormationEnabled *bool       `json:"team_formation_enabled,omitempty"`
}

// VoteRequest represents the request to cast a vote
type VoteRequest struct {
	SelectedOptions []string `json:"selected_options" validate:"required,min=1"`
}

// VoteFilters represents filters for listing votes
type VoteFilters struct {
	Status    *VoteStatus `json:"status,omitempty"`
	CreatedBy *string     `json:"created_by,omitempty"`
	Type      *VoteType   `json:"type,omitempty"`
	GroupID   *string     `json:"group_id,omitempty"`
	Limit     int         `json:"limit" validate:"min=1,max=100"`
	Offset    int         `json:"offset" validate:"min=0"`
}

// VoteResponse represents the response for vote operations
type VoteResponse struct {
	Message string `json:"message"`
	Vote    *Vote  `json:"vote,omitempty"`
}

// PaginatedVoteList represents a paginated list of votes
type PaginatedVoteList struct {
	Votes      []*VoteWithCreator `json:"votes"`
	TotalItems int                `json:"total_items"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
	TotalPages int                `json:"total_pages"`
}
