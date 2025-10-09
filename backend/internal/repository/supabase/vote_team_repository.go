package supabase

import (
	"context"
	"fmt"
	"spark-park-cricket-backend/internal/models"

	"github.com/supabase-community/postgrest-go"
	"github.com/supabase-community/supabase-go"
)

// VoteTeamRepository implements VoteTeamRepositoryInterface using Supabase
type VoteTeamRepository struct {
	client *supabase.Client
	schema string
}

// NewVoteTeamRepository creates a new Supabase vote team repository
func NewVoteTeamRepository(client *supabase.Client, schema string) *VoteTeamRepository {
	return &VoteTeamRepository{
		client: client,
		schema: schema,
	}
}

// CreateTeam creates a new team for a vote
func (r *VoteTeamRepository) CreateTeam(ctx context.Context, team *models.VoteTeam) error {
	var result []models.VoteTeam
	_, err := r.client.From("vote_teams").Insert(team, false, "", "", "").ExecuteTo(&result)
	if err != nil {
		return fmt.Errorf("failed to create team: %w", err)
	}
	if len(result) > 0 {
		*team = result[0]
	}
	return nil
}

// GetTeamByID retrieves a team by ID
func (r *VoteTeamRepository) GetTeamByID(ctx context.Context, id string) (*models.VoteTeam, error) {
	var teams []models.VoteTeam
	_, err := r.client.From("vote_teams").Select("*", "", false).Eq("id", id).ExecuteTo(&teams)
	if err != nil {
		return nil, fmt.Errorf("failed to get team: %w", err)
	}
	if len(teams) == 0 {
		return nil, fmt.Errorf("team not found")
	}
	return &teams[0], nil
}

// GetTeamsByVoteID retrieves all teams for a vote
func (r *VoteTeamRepository) GetTeamsByVoteID(ctx context.Context, voteID string) ([]*models.VoteTeam, error) {
	var teams []models.VoteTeam
	_, err := r.client.From("vote_teams").Select("*", "", false).Eq("vote_id", voteID).Order("team_letter", &postgrest.OrderOpts{Ascending: true}).ExecuteTo(&teams)
	if err != nil {
		return nil, fmt.Errorf("failed to get teams: %w", err)
	}

	result := make([]*models.VoteTeam, len(teams))
	for i := range teams {
		result[i] = &teams[i]
	}
	return result, nil
}

// GetTeamByVoteAndLetter retrieves a specific team by vote ID and letter
func (r *VoteTeamRepository) GetTeamByVoteAndLetter(ctx context.Context, voteID, teamLetter string) (*models.VoteTeam, error) {
	var teams []models.VoteTeam
	_, err := r.client.From("vote_teams").Select("*", "", false).Eq("vote_id", voteID).Eq("team_letter", teamLetter).ExecuteTo(&teams)
	if err != nil {
		return nil, fmt.Errorf("failed to get team: %w", err)
	}
	if len(teams) == 0 {
		return nil, nil // Team doesn't exist yet
	}
	return &teams[0], nil
}

// UpdateTeam updates team details
func (r *VoteTeamRepository) UpdateTeam(ctx context.Context, id string, teamName *string, captainID *string) error {
	updates := make(map[string]interface{})
	if teamName != nil {
		updates["team_name"] = *teamName
	}
	if captainID != nil {
		updates["captain_id"] = *captainID
	}

	if len(updates) == 0 {
		return nil // Nothing to update
	}

	_, err := r.client.From("vote_teams").Update(updates, "", "").Eq("id", id).ExecuteTo(nil)
	if err != nil {
		return fmt.Errorf("failed to update team: %w", err)
	}
	return nil
}

// DeleteTeam deletes a team and all its players (cascade)
func (r *VoteTeamRepository) DeleteTeam(ctx context.Context, id string) error {
	_, err := r.client.From("vote_teams").Delete("", "").Eq("id", id).ExecuteTo(nil)
	if err != nil {
		return fmt.Errorf("failed to delete team: %w", err)
	}
	return nil
}

// GetTeamWithPlayers retrieves a team with its players and captain info
func (r *VoteTeamRepository) GetTeamWithPlayers(ctx context.Context, teamID string) (*models.VoteTeamWithPlayers, error) {
	// Get team
	team, err := r.GetTeamByID(ctx, teamID)
	if err != nil {
		return nil, err
	}

	// Get players
	players, err := r.GetPlayersByTeamID(ctx, teamID)
	if err != nil {
		return nil, err
	}

	// Get captain info
	var captains []models.User
	_, err = r.client.From("users").Select("*", "", false).Eq("id", team.CaptainID).ExecuteTo(&captains)
	if err != nil {
		return nil, fmt.Errorf("failed to get captain: %w", err)
	}

	result := &models.VoteTeamWithPlayers{
		VoteTeam:    *team,
		Players:     players,
		PlayerCount: len(players),
	}

	if len(captains) > 0 {
		result.Captain = &captains[0]
		result.CaptainName = captains[0].Name
	}

	// Extract player names
	result.PlayerNames = make([]string, len(players))
	for i, player := range players {
		result.PlayerNames[i] = player.Name
	}

	return result, nil
}

// GetTeamsWithPlayersByVoteID retrieves all teams with players for a vote
func (r *VoteTeamRepository) GetTeamsWithPlayersByVoteID(ctx context.Context, voteID string) ([]*models.VoteTeamWithPlayers, error) {
	teams, err := r.GetTeamsByVoteID(ctx, voteID)
	if err != nil {
		return nil, err
	}

	result := make([]*models.VoteTeamWithPlayers, len(teams))
	for i, team := range teams {
		teamWithPlayers, err := r.GetTeamWithPlayers(ctx, team.ID)
		if err != nil {
			return nil, err
		}
		result[i] = teamWithPlayers
	}

	return result, nil
}

// AddPlayerToTeam adds a player to a team
func (r *VoteTeamRepository) AddPlayerToTeam(ctx context.Context, teamID, userID string) error {
	player := &models.TeamPlayer{
		TeamID: teamID,
		UserID: userID,
	}

	var result []models.TeamPlayer
	_, err := r.client.From("team_players").Insert(player, false, "", "", "").ExecuteTo(&result)
	if err != nil {
		return fmt.Errorf("failed to add player to team: %w", err)
	}
	return nil
}

// RemovePlayerFromTeam removes a player from a team
func (r *VoteTeamRepository) RemovePlayerFromTeam(ctx context.Context, teamID, userID string) error {
	_, err := r.client.From("team_players").Delete("", "").Eq("team_id", teamID).Eq("user_id", userID).ExecuteTo(nil)
	if err != nil {
		return fmt.Errorf("failed to remove player from team: %w", err)
	}
	return nil
}

// GetPlayersByTeamID retrieves all players for a team
func (r *VoteTeamRepository) GetPlayersByTeamID(ctx context.Context, teamID string) ([]*models.User, error) {
	// Get team player records
	var teamPlayers []models.TeamPlayer
	_, err := r.client.From("team_players").Select("*", "", false).Eq("team_id", teamID).ExecuteTo(&teamPlayers)
	if err != nil {
		return nil, fmt.Errorf("failed to get team players: %w", err)
	}

	if len(teamPlayers) == 0 {
		return []*models.User{}, nil
	}

	// Get user IDs
	userIDs := make([]string, len(teamPlayers))
	for i, tp := range teamPlayers {
		userIDs[i] = tp.UserID
	}

	// Get user details
	var users []models.User
	_, err = r.client.From("users").Select("*", "", false).In("id", userIDs).ExecuteTo(&users)
	if err != nil {
		return nil, fmt.Errorf("failed to get player details: %w", err)
	}

	result := make([]*models.User, len(users))
	for i := range users {
		result[i] = &users[i]
	}

	return result, nil
}

// IsUserInTeam checks if a user is in a specific team
func (r *VoteTeamRepository) IsUserInTeam(ctx context.Context, teamID, userID string) (bool, error) {
	var players []models.TeamPlayer
	_, err := r.client.From("team_players").Select("id", "", false).Eq("team_id", teamID).Eq("user_id", userID).ExecuteTo(&players)
	if err != nil {
		return false, fmt.Errorf("failed to check team membership: %w", err)
	}
	return len(players) > 0, nil
}

// IsUserInAnyTeamForVote checks if a user is in any team for a vote
func (r *VoteTeamRepository) IsUserInAnyTeamForVote(ctx context.Context, voteID, userID string) (bool, error) {
	// Get all teams for vote
	teams, err := r.GetTeamsByVoteID(ctx, voteID)
	if err != nil {
		return false, err
	}

	// Check each team
	for _, team := range teams {
		inTeam, err := r.IsUserInTeam(ctx, team.ID, userID)
		if err != nil {
			return false, err
		}
		if inTeam {
			return true, nil
		}
	}

	return false, nil
}

// GetTeamPlayerCount returns the number of players in a team
func (r *VoteTeamRepository) GetTeamPlayerCount(ctx context.Context, teamID string) (int, error) {
	var players []models.TeamPlayer
	_, err := r.client.From("team_players").Select("id", "", false).Eq("team_id", teamID).ExecuteTo(&players)
	if err != nil {
		return 0, fmt.Errorf("failed to count team players: %w", err)
	}
	return len(players), nil
}

// HasUserVoted checks if a user has voted in a specific vote
func (r *VoteTeamRepository) HasUserVoted(ctx context.Context, voteID, userID string) (bool, error) {
	var userVotes []models.UserVote
	_, err := r.client.From("user_votes").Select("id", "", false).Eq("vote_id", voteID).Eq("user_id", userID).ExecuteTo(&userVotes)
	if err != nil {
		return false, fmt.Errorf("failed to check if user voted: %w", err)
	}
	return len(userVotes) > 0, nil
}
