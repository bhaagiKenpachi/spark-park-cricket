package models

import "time"

// GroupType represents the type of group
type GroupType string

const (
	GroupTypeCustom    GroupType = "custom"    // Custom group created by users
	GroupTypeTeam      GroupType = "team"      // Team-based group (Team A, Team B)
	GroupTypeSeries    GroupType = "series"    // Series-based group
	GroupTypeMatch     GroupType = "match"     // Match-based group
	GroupTypeLocation  GroupType = "location"  // Location-based group
	GroupTypeSkill     GroupType = "skill"     // Skill-based group (beginner, intermediate, advanced)
)

// GroupStatus represents the status of a group
type GroupStatus string

const (
	GroupStatusActive   GroupStatus = "active"   // Group is active and can be used for voting
	GroupStatusInactive GroupStatus = "inactive" // Group is inactive
	GroupStatusArchived GroupStatus = "archived" // Group is archived
)

// Group represents a voting group
type Group struct {
	ID          string      `json:"id" db:"id"`
	Name        string      `json:"name" db:"name"`
	Description string      `json:"description" db:"description"`
	Type        GroupType   `json:"type" db:"type"`
	Status      GroupStatus `json:"status" db:"status"`
	CreatedBy   string      `json:"created_by" db:"created_by"`
	CreatedAt   time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at" db:"updated_at"`
}

// GroupWithCreator represents a group with creator information
type GroupWithCreator struct {
	Group
	CreatorName string `json:"creator_name"`
}

// GroupMember represents a member of a group
type GroupMember struct {
	ID        string    `json:"id" db:"id"`
	GroupID   string    `json:"group_id" db:"group_id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Role      string    `json:"role" db:"role"` // "member", "admin", "moderator"
	JoinedAt  time.Time `json:"joined_at" db:"joined_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// GroupWithMembers represents a group with its members
type GroupWithMembers struct {
	Group
	Members     []*GroupMember `json:"members"`
	MemberCount int            `json:"member_count"`
	CreatorName string         `json:"creator_name"`
}

// GroupMemberWithUser represents a group member with user information
type GroupMemberWithUser struct {
	GroupMember
	UserName  string `json:"user_name"`
	UserEmail string `json:"user_email"`
}

// GroupWithMembersAndUsers represents a group with members and their user details
type GroupWithMembersAndUsers struct {
	Group
	Members     []*GroupMemberWithUser `json:"members"`
	MemberCount int                    `json:"member_count"`
	CreatorName string                 `json:"creator_name"`
}

// VoteGroup represents the association between a vote and groups
type VoteGroup struct {
	ID        string    `json:"id" db:"id"`
	VoteID    string    `json:"vote_id" db:"vote_id"`
	GroupID   string    `json:"group_id" db:"group_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// VoteWithGroups represents a vote with its associated groups
type VoteWithGroups struct {
	Vote
	Groups      []*Group `json:"groups"`
	GroupCount  int      `json:"group_count"`
	CreatorName string   `json:"creator_name"`
}

// GroupVoteResult represents voting results for a specific group
type GroupVoteResult struct {
	GroupID       string                 `json:"group_id"`
	GroupName     string                 `json:"group_name"`
	TotalVotes    int                    `json:"total_votes"`
	Results       map[string]int         `json:"results"`            // option_id -> vote_count
	ResultsWithNames map[string][]VoterInfo `json:"results_with_names"` // option_id -> list of voters
	VotedMembers  []string               `json:"voted_members"`      // List of user IDs who voted
}

// VoteWithGroupResults represents a vote with results broken down by groups
type VoteWithGroupResults struct {
	Vote
	Options        []*VoteOption     `json:"options"`
	GroupResults   []*GroupVoteResult `json:"group_results"`
	OverallResults map[string]int   `json:"overall_results"` // option_id -> total vote_count
	TotalVotes     int              `json:"total_votes"`
	CreatorName    string           `json:"creator_name"`
}

// CreateGroupRequest represents the request to create a new group
type CreateGroupRequest struct {
	Name        string    `json:"name" validate:"required,min=3,max=255"`
	Description string    `json:"description" validate:"required,min=10,max=1000"`
	Type        GroupType `json:"type" validate:"required,oneof=custom team series match location skill"`
}

// UpdateGroupRequest represents the request to update a group
type UpdateGroupRequest struct {
	Name        *string      `json:"name,omitempty" validate:"omitempty,min=3,max=255"`
	Description *string      `json:"description,omitempty" validate:"omitempty,min=10,max=1000"`
	Status      *GroupStatus `json:"status,omitempty" validate:"omitempty,oneof=active inactive archived"`
}

// AddGroupMemberRequest represents the request to add a member to a group
type AddGroupMemberRequest struct {
	UserID string `json:"user_id" validate:"required,uuid"`
	Role   string `json:"role" validate:"omitempty,oneof=member admin moderator"`
}

// RemoveGroupMemberRequest represents the request to remove a member from a group
type RemoveGroupMemberRequest struct {
	UserID string `json:"user_id" validate:"required,uuid"`
}

// AssignGroupsToVoteRequest represents the request to assign groups to a vote
type AssignGroupsToVoteRequest struct {
	GroupIDs []string `json:"group_ids" validate:"required,min=1,dive,uuid"`
}
