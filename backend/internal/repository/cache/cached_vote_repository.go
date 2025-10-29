package cache

import (
	"context"
	"fmt"
	"spark-park-cricket-backend/internal/cache"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/repository/interfaces"
)

// CachedVoteRepository wraps a vote repository with caching
type CachedVoteRepository struct {
	repo  interfaces.VoteRepositoryInterface
	cache *cache.CacheManager
}

// NewCachedVoteRepository creates a new cached vote repository
func NewCachedVoteRepository(repo interfaces.VoteRepositoryInterface, cacheManager *cache.CacheManager) interfaces.VoteRepositoryInterface {
	return &CachedVoteRepository{
		repo:  repo,
		cache: cacheManager,
	}
}

// CreateVote creates a new vote and invalidates list cache
func (r *CachedVoteRepository) CreateVote(ctx context.Context, vote *models.Vote) error {
	err := r.repo.CreateVote(ctx, vote)
	if err != nil {
		return err
	}

	// Invalidate vote list caches
	r.invalidateVoteListCaches()

	// Cache the new vote
	if vote.ID != "" {
		key := fmt.Sprintf("vote:%s", vote.ID)
		_ = r.cache.Set(key, vote, cache.StaticDataTTL)
	}

	return nil
}

// GetVoteByID retrieves a vote by ID with caching
func (r *CachedVoteRepository) GetVoteByID(ctx context.Context, id string) (*models.Vote, error) {
	key := fmt.Sprintf("vote:%s", id)

	var vote models.Vote
	err := r.cache.GetOrSet(key, &vote, cache.StaticDataTTL, func() (interface{}, error) {
		return r.repo.GetVoteByID(ctx, id)
	})

	if err != nil {
		return nil, err
	}

	return &vote, nil
}

// GetVoteWithOptions retrieves a vote with options with caching
func (r *CachedVoteRepository) GetVoteWithOptions(ctx context.Context, id string) (*models.VoteWithOptions, error) {
	key := fmt.Sprintf("vote:%s:with_options", id)

	var voteWithOptions models.VoteWithOptions
	err := r.cache.GetOrSet(key, &voteWithOptions, cache.StaticDataTTL, func() (interface{}, error) {
		return r.repo.GetVoteWithOptions(ctx, id)
	})

	if err != nil {
		return nil, err
	}

	return &voteWithOptions, nil
}

// GetVoteWithResults retrieves vote results with caching (shorter TTL for dynamic data)
func (r *CachedVoteRepository) GetVoteWithResults(ctx context.Context, id string, userID string) (*models.VoteWithResults, error) {
	// Use shorter TTL for results as they change frequently
	key := fmt.Sprintf("vote:%s:results:user:%s", id, userID)

	var voteWithResults models.VoteWithResults
	err := r.cache.GetOrSet(key, &voteWithResults, cache.LiveDataTTL, func() (interface{}, error) {
		return r.repo.GetVoteWithResults(ctx, id, userID)
	})

	if err != nil {
		return nil, err
	}

	return &voteWithResults, nil
}

// ListVotes retrieves votes with caching
func (r *CachedVoteRepository) ListVotes(ctx context.Context, filters *models.VoteFilters) (*models.PaginatedVoteList, error) {
	// Create cache key based on filters
	cacheKey := "votes:list"
	if filters != nil {
		if filters.Status != nil {
			cacheKey += fmt.Sprintf(":status:%s", *filters.Status)
		}
		if filters.Type != nil {
			cacheKey += fmt.Sprintf(":type:%s", *filters.Type)
		}
		if filters.CreatedBy != nil {
			cacheKey += fmt.Sprintf(":creator:%s", *filters.CreatedBy)
		}
		if filters.GroupID != nil && *filters.GroupID != "" {
			cacheKey += fmt.Sprintf(":group:%s", *filters.GroupID)
		}
		if filters.Limit > 0 {
			cacheKey += fmt.Sprintf(":limit:%d", filters.Limit)
		}
		if filters.Offset > 0 {
			cacheKey += fmt.Sprintf(":offset:%d", filters.Offset)
		}
	}

	var paginatedVotes models.PaginatedVoteList
	err := r.cache.GetOrSet(cacheKey, &paginatedVotes, cache.MatchListTTL, func() (interface{}, error) {
		return r.repo.ListVotes(ctx, filters)
	})

	if err != nil {
		return nil, err
	}

	return &paginatedVotes, nil
}

// CountVotes counts votes with caching
func (r *CachedVoteRepository) CountVotes(ctx context.Context, filters *models.VoteFilters) (int, error) {
	// Create cache key based on filters (same as ListVotes but with :count suffix)
	cacheKey := "votes:count"
	if filters != nil {
		if filters.Status != nil {
			cacheKey += fmt.Sprintf(":status:%s", *filters.Status)
		}
		if filters.Type != nil {
			cacheKey += fmt.Sprintf(":type:%s", *filters.Type)
		}
		if filters.CreatedBy != nil {
			cacheKey += fmt.Sprintf(":creator:%s", *filters.CreatedBy)
		}
		if filters.GroupID != nil && *filters.GroupID != "" {
			cacheKey += fmt.Sprintf(":group:%s", *filters.GroupID)
		}
	}

	var count int
	err := r.cache.GetOrSet(cacheKey, &count, cache.MatchListTTL, func() (interface{}, error) {
		return r.repo.CountVotes(ctx, filters)
	})

	if err != nil {
		return 0, err
	}

	return count, nil
}

// UpdateVote updates a vote and invalidates related caches
func (r *CachedVoteRepository) UpdateVote(ctx context.Context, vote *models.Vote) error {
	err := r.repo.UpdateVote(ctx, vote)
	if err != nil {
		return err
	}

	// Invalidate vote-specific caches
	if vote.ID != "" {
		r.invalidateVoteCaches(vote.ID)
	}

	// Invalidate list caches since vote title/description may have changed
	r.invalidateVoteListCaches()

	return nil
}

// DeleteVote deletes a vote and invalidates all related caches
func (r *CachedVoteRepository) DeleteVote(ctx context.Context, id string) error {
	err := r.repo.DeleteVote(ctx, id)
	if err != nil {
		return err
	}

	// Invalidate all vote-related caches
	r.invalidateVoteCaches(id)
	r.invalidateVoteListCaches()

	return nil
}

// Vote option operations
func (r *CachedVoteRepository) CreateVoteOptions(ctx context.Context, options []*models.VoteOption) error {
	err := r.repo.CreateVoteOptions(ctx, options)
	if err != nil {
		return err
	}

	// Invalidate vote options cache for each vote
	for _, option := range options {
		if option.VoteID != "" {
			r.invalidateVoteCaches(option.VoteID)
		}
	}
	return nil
}

func (r *CachedVoteRepository) GetVoteOptions(ctx context.Context, voteID string) ([]*models.VoteOption, error) {
	key := fmt.Sprintf("vote:%s:options", voteID)

	var options []*models.VoteOption
	err := r.cache.GetOrSet(key, &options, cache.StaticDataTTL, func() (interface{}, error) {
		return r.repo.GetVoteOptions(ctx, voteID)
	})

	if err != nil {
		return nil, err
	}

	return options, nil
}

func (r *CachedVoteRepository) UpdateVoteOption(ctx context.Context, option *models.VoteOption) error {
	err := r.repo.UpdateVoteOption(ctx, option)
	if err != nil {
		return err
	}

	// Invalidate vote options cache
	if option.VoteID != "" {
		r.invalidateVoteCaches(option.VoteID)
	}
	return nil
}

func (r *CachedVoteRepository) DeleteVoteOption(ctx context.Context, id string) error {
	// Note: We can't invalidate specific vote cache here without knowing voteID
	// The service layer should handle this or we need to fetch first
	return r.repo.DeleteVoteOption(ctx, id)
}

// User vote operations - invalidate results cache when votes are cast/updated
func (r *CachedVoteRepository) CreateUserVote(ctx context.Context, userVote *models.UserVote) error {
	err := r.repo.CreateUserVote(ctx, userVote)
	if err != nil {
		return err
	}

	// Invalidate results cache for this vote
	r.invalidateVoteResultsCaches(userVote.VoteID)
	return nil
}

func (r *CachedVoteRepository) UpdateUserVote(ctx context.Context, userVote *models.UserVote) error {
	err := r.repo.UpdateUserVote(ctx, userVote)
	if err != nil {
		return err
	}

	// Invalidate results cache for this vote
	r.invalidateVoteResultsCaches(userVote.VoteID)
	return nil
}

func (r *CachedVoteRepository) GetUserVote(ctx context.Context, voteID, userID string) (*models.UserVote, error) {
	key := fmt.Sprintf("vote:%s:user_vote:%s", voteID, userID)

	var userVote models.UserVote
	err := r.cache.GetOrSet(key, &userVote, cache.LiveDataTTL, func() (interface{}, error) {
		return r.repo.GetUserVote(ctx, voteID, userID)
	})

	if err != nil {
		return nil, err
	}

	return &userVote, nil
}

func (r *CachedVoteRepository) HasUserVoted(ctx context.Context, voteID, userID string) (bool, error) {
	// Don't cache boolean checks, just pass through
	return r.repo.HasUserVoted(ctx, voteID, userID)
}

func (r *CachedVoteRepository) GetVoteResults(ctx context.Context, voteID string) (map[string]int, error) {
	key := fmt.Sprintf("vote:%s:results", voteID)

	var results map[string]int
	err := r.cache.GetOrSet(key, &results, cache.LiveDataTTL, func() (interface{}, error) {
		return r.repo.GetVoteResults(ctx, voteID)
	})

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (r *CachedVoteRepository) GetVoteResultsWithNames(ctx context.Context, voteID string) (map[string][]models.VoterInfo, error) {
	key := fmt.Sprintf("vote:%s:results_with_names", voteID)

	var results map[string][]models.VoterInfo
	err := r.cache.GetOrSet(key, &results, cache.LiveDataTTL, func() (interface{}, error) {
		return r.repo.GetVoteResultsWithNames(ctx, voteID)
	})

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (r *CachedVoteRepository) GetVotedUsers(ctx context.Context, voteID string) ([]string, error) {
	key := fmt.Sprintf("vote:%s:voted_users", voteID)

	var users []string
	err := r.cache.GetOrSet(key, &users, cache.LiveDataTTL, func() (interface{}, error) {
		return r.repo.GetVotedUsers(ctx, voteID)
	})

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *CachedVoteRepository) GetTotalVoteCount(ctx context.Context, voteID string) (int, error) {
	key := fmt.Sprintf("vote:%s:total_count", voteID)

	var count int
	err := r.cache.GetOrSet(key, &count, cache.LiveDataTTL, func() (interface{}, error) {
		return r.repo.GetTotalVoteCount(ctx, voteID)
	})

	if err != nil {
		return 0, err
	}

	return count, nil
}

// Helper methods for cache invalidation
func (r *CachedVoteRepository) invalidateVoteCaches(voteID string) {
	keysToInvalidate := []string{
		fmt.Sprintf("vote:%s", voteID),
		fmt.Sprintf("vote:%s:with_options", voteID),
		fmt.Sprintf("vote:%s:options", voteID),
	}

	for _, key := range keysToInvalidate {
		_ = r.cache.Invalidate(key)
	}

	// Also invalidate results caches
	r.invalidateVoteResultsCaches(voteID)
}

func (r *CachedVoteRepository) invalidateVoteResultsCaches(voteID string) {
	// Invalidate exact keys
	exactKeys := []string{
		fmt.Sprintf("vote:%s:results", voteID),
		fmt.Sprintf("vote:%s:results_with_names", voteID),
		fmt.Sprintf("vote:%s:voted_users", voteID),
		fmt.Sprintf("vote:%s:total_count", voteID),
	}

	for _, key := range exactKeys {
		_ = r.cache.Invalidate(key)
	}

	// Pattern-based invalidation for user-specific caches
	// This invalidates ALL user-specific result caches for this vote
	resultPattern := fmt.Sprintf("vote:%s:results:user:*", voteID)
	_ = r.cache.InvalidatePattern(resultPattern)

	userVotePattern := fmt.Sprintf("vote:%s:user_vote:*", voteID)
	_ = r.cache.InvalidatePattern(userVotePattern)
}

func (r *CachedVoteRepository) invalidateVoteListCaches() {
	// Invalidate exact keys
	_ = r.cache.Invalidate("votes:list")
	_ = r.cache.Invalidate("votes:count")

	// Invalidate all pattern-based list and count caches
	_ = r.cache.InvalidatePattern("votes:list:*")
	_ = r.cache.InvalidatePattern("votes:count:*")
}
