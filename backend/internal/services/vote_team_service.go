package services

import (
	"context"
	"fmt"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/repository/interfaces"
	"strings"
)

// VoteTeamService implements VoteTeamServiceInterface
type VoteTeamService struct {
	teamRepo interfaces.VoteTeamRepositoryInterface
	voteRepo interfaces.VoteRepositoryInterface
}

// NewVoteTeamService creates a new vote team service
func NewVoteTeamService(teamRepo interfaces.VoteTeamRepositoryInterface, voteRepo interfaces.VoteRepositoryInterface) *VoteTeamService {
	return &VoteTeamService{
		teamRepo: teamRepo,
		voteRepo: voteRepo,
	}
}

// CreateTeam creates a new team for a vote
func (s *VoteTeamService) CreateTeam(ctx context.Context, userID string, req *models.CreateVoteTeamRequest) (*models.VoteTeam, error) {
	// Validate team name
	req.TeamName = strings.TrimSpace(req.TeamName)
	if len(req.TeamName) < 2 {
		return nil, fmt.Errorf("team name must be at least 2 characters")
	}
	if len(req.TeamName) > 100 {
		return nil, fmt.Errorf("team name cannot exceed 100 characters")
	}

	// Validate team letter
	req.TeamLetter = strings.ToUpper(req.TeamLetter)
	if req.TeamLetter != "A" && req.TeamLetter != "B" {
		return nil, fmt.Errorf("team letter must be A or B")
	}

	// Verify vote exists
	_, err := s.voteRepo.GetVoteByID(ctx, req.VoteID)
	if err != nil {
		return nil, fmt.Errorf("vote not found: %w", err)
	}

	// Verify captain has voted
	hasVoted, err := s.teamRepo.HasUserVoted(ctx, req.VoteID, req.CaptainID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify captain voted: %w", err)
	}
	if !hasVoted {
		return nil, fmt.Errorf("captain must be one of the voters")
	}

	// Check if team with this letter already exists for this vote
	existingTeam, err := s.teamRepo.GetTeamByVoteAndLetter(ctx, req.VoteID, req.TeamLetter)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing team: %w", err)
	}
	if existingTeam != nil {
		return nil, fmt.Errorf("team %s already exists for this vote", req.TeamLetter)
	}

	// Create team
	team := &models.VoteTeam{
		VoteID:     req.VoteID,
		TeamName:   req.TeamName,
		TeamLetter: req.TeamLetter,
		CaptainID:  req.CaptainID,
		CreatedBy:  userID,
	}

	err = s.teamRepo.CreateTeam(ctx, team)
	if err != nil {
		return nil, fmt.Errorf("failed to create team: %w", err)
	}

	// Add captain as first player
	err = s.teamRepo.AddPlayerToTeam(ctx, team.ID, req.CaptainID)
	if err != nil {
		// Rollback team creation if adding captain fails
		_ = s.teamRepo.DeleteTeam(ctx, team.ID)
		return nil, fmt.Errorf("failed to add captain to team: %w", err)
	}

	return team, nil
}

// GetTeamByID retrieves a team with its players
func (s *VoteTeamService) GetTeamByID(ctx context.Context, teamID string) (*models.VoteTeamWithPlayers, error) {
	team, err := s.teamRepo.GetTeamWithPlayers(ctx, teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to get team: %w", err)
	}
	return team, nil
}

// GetTeamsByVoteID retrieves all teams for a vote
func (s *VoteTeamService) GetTeamsByVoteID(ctx context.Context, voteID string) ([]*models.VoteTeamWithPlayers, error) {
	teams, err := s.teamRepo.GetTeamsWithPlayersByVoteID(ctx, voteID)
	if err != nil {
		return nil, fmt.Errorf("failed to get teams: %w", err)
	}
	return teams, nil
}

// UpdateTeam updates team details
func (s *VoteTeamService) UpdateTeam(ctx context.Context, userID, teamID string, req *models.UpdateVoteTeamRequest) (*models.VoteTeam, error) {
	// Get existing team
	team, err := s.teamRepo.GetTeamByID(ctx, teamID)
	if err != nil {
		return nil, fmt.Errorf("team not found: %w", err)
	}

	// Validate team name if provided
	if req.TeamName != nil {
		*req.TeamName = strings.TrimSpace(*req.TeamName)
		if len(*req.TeamName) < 2 {
			return nil, fmt.Errorf("team name must be at least 2 characters")
		}
		if len(*req.TeamName) > 100 {
			return nil, fmt.Errorf("team name cannot exceed 100 characters")
		}
	}

	// Validate new captain if provided
	if req.CaptainID != nil {
		// Verify new captain has voted
		hasVoted, err := s.teamRepo.HasUserVoted(ctx, team.VoteID, *req.CaptainID)
		if err != nil {
			return nil, fmt.Errorf("failed to verify captain voted: %w", err)
		}
		if !hasVoted {
			return nil, fmt.Errorf("captain must be one of the voters")
		}

		// Verify new captain is in the team
		isInTeam, err := s.teamRepo.IsUserInTeam(ctx, teamID, *req.CaptainID)
		if err != nil {
			return nil, fmt.Errorf("failed to verify captain in team: %w", err)
		}
		if !isInTeam {
			// Add captain to team if not already there
			err = s.teamRepo.AddPlayerToTeam(ctx, teamID, *req.CaptainID)
			if err != nil {
				return nil, fmt.Errorf("failed to add new captain to team: %w", err)
			}
		}
	}

	// Update team
	err = s.teamRepo.UpdateTeam(ctx, teamID, req.TeamName, req.CaptainID)
	if err != nil {
		return nil, fmt.Errorf("failed to update team: %w", err)
	}

	// Get updated team
	updatedTeam, err := s.teamRepo.GetTeamByID(ctx, teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated team: %w", err)
	}

	return updatedTeam, nil
}

// DeleteTeam deletes a team and all its players
func (s *VoteTeamService) DeleteTeam(ctx context.Context, userID, teamID string) error {
	// Verify team exists
	_, err := s.teamRepo.GetTeamByID(ctx, teamID)
	if err != nil {
		return fmt.Errorf("team not found: %w", err)
	}

	// Delete team (cascade will delete players)
	err = s.teamRepo.DeleteTeam(ctx, teamID)
	if err != nil {
		return fmt.Errorf("failed to delete team: %w", err)
	}

	return nil
}

// AddPlayerToTeam adds a single player to a team
func (s *VoteTeamService) AddPlayerToTeam(ctx context.Context, userID, teamID string, req *models.AddPlayerRequest) error {
	// Get team
	team, err := s.teamRepo.GetTeamByID(ctx, teamID)
	if err != nil {
		return fmt.Errorf("team not found: %w", err)
	}

	// Verify user has voted
	hasVoted, err := s.teamRepo.HasUserVoted(ctx, team.VoteID, req.UserID)
	if err != nil {
		return fmt.Errorf("failed to verify user voted: %w", err)
	}
	if !hasVoted {
		return fmt.Errorf("user must have voted to be added to a team")
	}

	// Check if user is already in any team for this vote
	inAnyTeam, err := s.teamRepo.IsUserInAnyTeamForVote(ctx, team.VoteID, req.UserID)
	if err != nil {
		return fmt.Errorf("failed to check team membership: %w", err)
	}
	if inAnyTeam {
		return fmt.Errorf("user is already in a team for this vote")
	}

	// Check team player count
	count, err := s.teamRepo.GetTeamPlayerCount(ctx, teamID)
	if err != nil {
		return fmt.Errorf("failed to check team size: %w", err)
	}
	if count >= 20 {
		return fmt.Errorf("team is full (max 20 players)")
	}

	// Add player
	err = s.teamRepo.AddPlayerToTeam(ctx, teamID, req.UserID)
	if err != nil {
		return fmt.Errorf("failed to add player to team: %w", err)
	}

	return nil
}

// AddPlayersToTeam adds multiple players to a team
func (s *VoteTeamService) AddPlayersToTeam(ctx context.Context, userID, teamID string, req *models.TeamAssignmentRequest) error {
	// Get team
	team, err := s.teamRepo.GetTeamByID(ctx, teamID)
	if err != nil {
		return fmt.Errorf("team not found: %w", err)
	}

	// Check current team size
	currentCount, err := s.teamRepo.GetTeamPlayerCount(ctx, teamID)
	if err != nil {
		return fmt.Errorf("failed to check team size: %w", err)
	}

	// Check if adding all players would exceed limit
	if currentCount+len(req.UserIDs) > 20 {
		return fmt.Errorf("adding %d players would exceed team limit of 20 (current: %d)", len(req.UserIDs), currentCount)
	}

	// Validate and add each player
	for _, playerID := range req.UserIDs {
		// Verify user has voted
		hasVoted, err := s.teamRepo.HasUserVoted(ctx, team.VoteID, playerID)
		if err != nil {
			return fmt.Errorf("failed to verify user %s voted: %w", playerID, err)
		}
		if !hasVoted {
			return fmt.Errorf("user %s has not voted and cannot be added to team", playerID)
		}

		// Check if user is already in any team for this vote
		inAnyTeam, err := s.teamRepo.IsUserInAnyTeamForVote(ctx, team.VoteID, playerID)
		if err != nil {
			return fmt.Errorf("failed to check team membership for user %s: %w", playerID, err)
		}
		if inAnyTeam {
			return fmt.Errorf("user %s is already in a team for this vote", playerID)
		}

		// Add player
		err = s.teamRepo.AddPlayerToTeam(ctx, teamID, playerID)
		if err != nil {
			return fmt.Errorf("failed to add player %s to team: %w", playerID, err)
		}
	}

	return nil
}

// RemovePlayerFromTeam removes a player from a team
func (s *VoteTeamService) RemovePlayerFromTeam(ctx context.Context, userID, teamID, playerID string) error {
	// Get team
	team, err := s.teamRepo.GetTeamByID(ctx, teamID)
	if err != nil {
		return fmt.Errorf("team not found: %w", err)
	}

	// Prevent removing captain
	if team.CaptainID == playerID {
		return fmt.Errorf("cannot remove captain from team. Change captain first")
	}

	// Verify player is in team
	isInTeam, err := s.teamRepo.IsUserInTeam(ctx, teamID, playerID)
	if err != nil {
		return fmt.Errorf("failed to check team membership: %w", err)
	}
	if !isInTeam {
		return fmt.Errorf("user is not in this team")
	}

	// Remove player
	err = s.teamRepo.RemovePlayerFromTeam(ctx, teamID, playerID)
	if err != nil {
		return fmt.Errorf("failed to remove player from team: %w", err)
	}

	return nil
}

// GetTeamPlayers retrieves all players for a team
func (s *VoteTeamService) GetTeamPlayers(ctx context.Context, teamID string) ([]*models.User, error) {
	players, err := s.teamRepo.GetPlayersByTeamID(ctx, teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to get team players: %w", err)
	}
	return players, nil
}
