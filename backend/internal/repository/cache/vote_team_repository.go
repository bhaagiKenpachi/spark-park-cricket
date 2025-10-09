package cache

import (
	"context"
	"fmt"
	"spark-park-cricket-backend/internal/cache"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/repository/interfaces"
	"time"
)

// CachedVoteTeamRepository wraps VoteTeamRepositoryInterface with caching
type CachedVoteTeamRepository struct {
	repo  interfaces.VoteTeamRepositoryInterface
	cache *cache.CacheManager
}

// NewCachedVoteTeamRepository creates a new cached vote team repository
func NewCachedVoteTeamRepository(repo interfaces.VoteTeamRepositoryInterface, cacheManager *cache.CacheManager) *CachedVoteTeamRepository {
	return &CachedVoteTeamRepository{
		repo:  repo,
		cache: cacheManager,
	}
}

// Cache key patterns
const (
	voteTeamKeyPattern       = "vote_team:%s"
	voteTeamsListKeyPattern  = "vote:%s:teams"
	teamPlayersKeyPattern    = "team:%s:players"
	voteTeamsPrefixPattern   = "vote:%s:teams*"
	teamPlayersPrefixPattern = "team:*:players"
)

// CreateTeam creates a new team and invalidates related caches
func (r *CachedVoteTeamRepository) CreateTeam(ctx context.Context, team *models.VoteTeam) error {
	err := r.repo.CreateTeam(ctx, team)
	if err != nil {
		return err
	}

	// Invalidate vote teams list cache
	voteTeamsKey := fmt.Sprintf(voteTeamsListKeyPattern, team.VoteID)
	_ = r.cache.Invalidate(voteTeamsKey)

	return nil
}

// GetTeamByID retrieves a team by ID with caching
func (r *CachedVoteTeamRepository) GetTeamByID(ctx context.Context, id string) (*models.VoteTeam, error) {
	cacheKey := fmt.Sprintf(voteTeamKeyPattern, id)

	// Try cache first
	var team models.VoteTeam
	err := r.cache.Get(cacheKey, &team)
	if err == nil {
		return &team, nil
	}

	// Cache miss, get from database
	teamPtr, err := r.repo.GetTeamByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Cache the result
	_ = r.cache.Set(cacheKey, teamPtr, 5*time.Minute)

	return teamPtr, nil
}

// GetTeamsWithPlayersByVoteID retrieves all teams for a vote with caching
func (r *CachedVoteTeamRepository) GetTeamsWithPlayersByVoteID(ctx context.Context, voteID string) ([]*models.VoteTeamWithPlayers, error) {
	cacheKey := fmt.Sprintf(voteTeamsListKeyPattern, voteID)

	// Try cache first
	var teams []*models.VoteTeamWithPlayers
	err := r.cache.Get(cacheKey, &teams)
	if err == nil {
		return teams, nil
	}

	// Cache miss, get from database
	teams, err = r.repo.GetTeamsWithPlayersByVoteID(ctx, voteID)
	if err != nil {
		return nil, err
	}

	// Cache the result
	_ = r.cache.Set(cacheKey, teams, 5*time.Minute)

	return teams, nil
}

// GetTeamWithPlayers retrieves a team with players with caching
func (r *CachedVoteTeamRepository) GetTeamWithPlayers(ctx context.Context, teamID string) (*models.VoteTeamWithPlayers, error) {
	cacheKey := fmt.Sprintf(voteTeamKeyPattern, teamID)

	// Try cache first
	var team models.VoteTeamWithPlayers
	err := r.cache.Get(cacheKey, &team)
	if err == nil {
		return &team, nil
	}

	// Cache miss, get from database
	teamPtr, err := r.repo.GetTeamWithPlayers(ctx, teamID)
	if err != nil {
		return nil, err
	}

	// Cache the result
	_ = r.cache.Set(cacheKey, teamPtr, 5*time.Minute)

	return teamPtr, nil
}

// UpdateTeam updates a team and invalidates related caches
func (r *CachedVoteTeamRepository) UpdateTeam(ctx context.Context, id string, teamName *string, captainID *string) error {
	err := r.repo.UpdateTeam(ctx, id, teamName, captainID)
	if err != nil {
		return err
	}

	// Invalidate team cache
	teamKey := fmt.Sprintf(voteTeamKeyPattern, id)
	_ = r.cache.Invalidate(teamKey)

	// Get team to invalidate vote teams list
	team, err := r.repo.GetTeamByID(ctx, id)
	if err == nil {
		voteTeamsKey := fmt.Sprintf(voteTeamsListKeyPattern, team.VoteID)
		_ = r.cache.Invalidate(voteTeamsKey)
	}

	return nil
}

// GetTeamsByVoteID retrieves teams by vote ID (delegates to base repo, not cached separately)
func (r *CachedVoteTeamRepository) GetTeamsByVoteID(ctx context.Context, voteID string) ([]*models.VoteTeam, error) {
	return r.repo.GetTeamsByVoteID(ctx, voteID)
}

// DeleteTeam deletes a team and invalidates related caches
func (r *CachedVoteTeamRepository) DeleteTeam(ctx context.Context, id string) error {
	// Get team first to know which vote to invalidate
	team, err := r.repo.GetTeamByID(ctx, id)
	if err != nil {
		return err
	}

	err = r.repo.DeleteTeam(ctx, id)
	if err != nil {
		return err
	}

	// Invalidate caches
	teamKey := fmt.Sprintf(voteTeamKeyPattern, id)
	_ = r.cache.Invalidate(teamKey)

	voteTeamsKey := fmt.Sprintf(voteTeamsListKeyPattern, team.VoteID)
	_ = r.cache.Invalidate(voteTeamsKey)

	return nil
}

// GetTeamByVoteAndLetter retrieves a team by vote and letter (no cache)
func (r *CachedVoteTeamRepository) GetTeamByVoteAndLetter(ctx context.Context, voteID, teamLetter string) (*models.VoteTeam, error) {
	return r.repo.GetTeamByVoteAndLetter(ctx, voteID, teamLetter)
}

// AddPlayerToTeam adds a player and invalidates related caches
func (r *CachedVoteTeamRepository) AddPlayerToTeam(ctx context.Context, teamID, userID string) error {
	err := r.repo.AddPlayerToTeam(ctx, teamID, userID)
	if err != nil {
		return err
	}

	// Invalidate team cache
	teamKey := fmt.Sprintf(voteTeamKeyPattern, teamID)
	_ = r.cache.Invalidate(teamKey)

	// Invalidate team players cache
	playersKey := fmt.Sprintf(teamPlayersKeyPattern, teamID)
	_ = r.cache.Invalidate(playersKey)

	// Get team to invalidate vote teams list
	team, err := r.repo.GetTeamByID(ctx, teamID)
	if err == nil {
		voteTeamsKey := fmt.Sprintf(voteTeamsListKeyPattern, team.VoteID)
		_ = r.cache.Invalidate(voteTeamsKey)
	}

	return nil
}

// RemovePlayerFromTeam removes a player and invalidates related caches
func (r *CachedVoteTeamRepository) RemovePlayerFromTeam(ctx context.Context, teamID, userID string) error {
	err := r.repo.RemovePlayerFromTeam(ctx, teamID, userID)
	if err != nil {
		return err
	}

	// Invalidate team cache
	teamKey := fmt.Sprintf(voteTeamKeyPattern, teamID)
	_ = r.cache.Invalidate(teamKey)

	// Invalidate team players cache
	playersKey := fmt.Sprintf(teamPlayersKeyPattern, teamID)
	_ = r.cache.Invalidate(playersKey)

	// Get team to invalidate vote teams list
	team, err := r.repo.GetTeamByID(ctx, teamID)
	if err == nil {
		voteTeamsKey := fmt.Sprintf(voteTeamsListKeyPattern, team.VoteID)
		_ = r.cache.Invalidate(voteTeamsKey)
	}

	return nil
}

// GetPlayersByTeamID retrieves players with caching
func (r *CachedVoteTeamRepository) GetPlayersByTeamID(ctx context.Context, teamID string) ([]*models.User, error) {
	cacheKey := fmt.Sprintf(teamPlayersKeyPattern, teamID)

	// Try cache first
	var players []*models.User
	err := r.cache.Get(cacheKey, &players)
	if err == nil {
		return players, nil
	}

	// Cache miss, get from database
	players, err = r.repo.GetPlayersByTeamID(ctx, teamID)
	if err != nil {
		return nil, err
	}

	// Cache the result
	_ = r.cache.Set(cacheKey, players, 5*time.Minute)

	return players, nil
}

// HasUserVoted checks if user voted (no cache, always fresh)
func (r *CachedVoteTeamRepository) HasUserVoted(ctx context.Context, voteID, userID string) (bool, error) {
	return r.repo.HasUserVoted(ctx, voteID, userID)
}

// IsUserInTeam checks if user is in a team (no cache, always fresh)
func (r *CachedVoteTeamRepository) IsUserInTeam(ctx context.Context, teamID, userID string) (bool, error) {
	return r.repo.IsUserInTeam(ctx, teamID, userID)
}

// IsUserInAnyTeamForVote checks if user is in any team for a vote (no cache, always fresh)
func (r *CachedVoteTeamRepository) IsUserInAnyTeamForVote(ctx context.Context, voteID, userID string) (bool, error) {
	return r.repo.IsUserInAnyTeamForVote(ctx, voteID, userID)
}

// GetTeamPlayerCount returns player count (no cache, always fresh)
func (r *CachedVoteTeamRepository) GetTeamPlayerCount(ctx context.Context, teamID string) (int, error) {
	return r.repo.GetTeamPlayerCount(ctx, teamID)
}
