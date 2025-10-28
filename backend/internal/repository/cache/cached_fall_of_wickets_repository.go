package cache

import (
	"context"
	"fmt"
	"spark-park-cricket-backend/internal/cache"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/repository/interfaces"
)

// CachedFallOfWicketsRepository wraps a fall of wickets repository with caching
type CachedFallOfWicketsRepository struct {
	repo  interfaces.FallOfWicketsRepository
	cache *cache.CacheManager
}

// NewCachedFallOfWicketsRepository creates a new cached fall of wickets repository
func NewCachedFallOfWicketsRepository(repo interfaces.FallOfWicketsRepository, cacheManager *cache.CacheManager) interfaces.FallOfWicketsRepository {
	return &CachedFallOfWicketsRepository{
		repo:  repo,
		cache: cacheManager,
	}
}

// Create creates a new fall of wickets record and invalidates related caches
func (r *CachedFallOfWicketsRepository) Create(ctx context.Context, fallOfWickets *models.FallOfWickets) error {
	err := r.repo.Create(ctx, fallOfWickets)
	if err != nil {
		return err
	}

	// Invalidate related caches
	r.invalidateFallOfWicketsCaches(fallOfWickets.MatchID, fallOfWickets.InningsID)

	// Cache the new record
	if fallOfWickets.ID != "" {
		key := fmt.Sprintf("fall_of_wickets:%s", fallOfWickets.ID)
		_ = r.cache.Set(key, fallOfWickets, cache.FallOfWicketsTTL)
	}

	return nil
}

// GetByID retrieves a fall of wickets record by ID with caching
func (r *CachedFallOfWicketsRepository) GetByID(ctx context.Context, id string) (*models.FallOfWickets, error) {
	key := fmt.Sprintf("fall_of_wickets:%s", id)

	var fallOfWickets models.FallOfWickets
	err := r.cache.GetOrSet(key, &fallOfWickets, cache.FallOfWicketsTTL, func() (interface{}, error) {
		return r.repo.GetByID(ctx, id)
	})

	if err != nil {
		return nil, err
	}

	return &fallOfWickets, nil
}

// GetByMatchID retrieves all fall of wickets for a specific match with caching
func (r *CachedFallOfWicketsRepository) GetByMatchID(ctx context.Context, matchID string) ([]*models.FallOfWickets, error) {
	key := fmt.Sprintf("fall_of_wickets:match:%s", matchID)

	var fallOfWickets []*models.FallOfWickets
	err := r.cache.GetOrSet(key, &fallOfWickets, cache.FallOfWicketsTTL, func() (interface{}, error) {
		return r.repo.GetByMatchID(ctx, matchID)
	})

	if err != nil {
		return nil, err
	}

	return fallOfWickets, nil
}

// GetByInningsID retrieves all fall of wickets for a specific innings with caching
func (r *CachedFallOfWicketsRepository) GetByInningsID(ctx context.Context, inningsID string) ([]*models.FallOfWickets, error) {
	key := fmt.Sprintf("fall_of_wickets:innings:%s", inningsID)

	var fallOfWickets []*models.FallOfWickets
	err := r.cache.GetOrSet(key, &fallOfWickets, cache.FallOfWicketsTTL, func() (interface{}, error) {
		return r.repo.GetByInningsID(ctx, inningsID)
	})

	if err != nil {
		return nil, err
	}

	return fallOfWickets, nil
}

// GetByBallID retrieves fall of wickets record by ball ID with caching
func (r *CachedFallOfWicketsRepository) GetByBallID(ctx context.Context, ballID string) (*models.FallOfWickets, error) {
	key := fmt.Sprintf("fall_of_wickets:ball:%s", ballID)

	var fallOfWickets models.FallOfWickets
	err := r.cache.GetOrSet(key, &fallOfWickets, cache.FallOfWicketsTTL, func() (interface{}, error) {
		return r.repo.GetByBallID(ctx, ballID)
	})

	if err != nil {
		return nil, err
	}

	return &fallOfWickets, nil
}

// List retrieves fall of wickets records with filters (no caching for complex filters)
func (r *CachedFallOfWicketsRepository) List(ctx context.Context, filters *models.FallOfWicketsFilters) ([]*models.FallOfWickets, error) {
	// For complex filters, don't cache - call repository directly
	return r.repo.List(ctx, filters)
}

// Update updates an existing fall of wickets record and invalidates caches
func (r *CachedFallOfWicketsRepository) Update(ctx context.Context, id string, req *models.UpdateFallOfWicketsRequest) (*models.FallOfWickets, error) {
	result, err := r.repo.Update(ctx, id, req)
	if err != nil {
		return nil, err
	}

	// Invalidate related caches
	if result != nil {
		r.invalidateFallOfWicketsCaches(result.MatchID, result.InningsID)

		// Update cache with new data
		key := fmt.Sprintf("fall_of_wickets:%s", id)
		_ = r.cache.Set(key, result, cache.FallOfWicketsTTL)
	}

	return result, nil
}

// Delete deletes a fall of wickets record and invalidates caches
func (r *CachedFallOfWicketsRepository) Delete(ctx context.Context, id string) error {
	// Get the record first to know which caches to invalidate
	record, err := r.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	err = r.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Invalidate related caches
	if record != nil {
		r.invalidateFallOfWicketsCaches(record.MatchID, record.InningsID)

		// Remove from cache
		key := fmt.Sprintf("fall_of_wickets:%s", id)
		_ = r.cache.Invalidate(key)
	}

	return nil
}

// GetSummary retrieves a summary of fall of wickets for a match/innings with caching
func (r *CachedFallOfWicketsRepository) GetSummary(ctx context.Context, matchID string, inningsID *string) (*models.FallOfWicketsSummary, error) {
	var key string
	if inningsID != nil && *inningsID != "" {
		key = fmt.Sprintf("fall_of_wickets:summary:%s:%s", matchID, *inningsID)
	} else {
		key = fmt.Sprintf("fall_of_wickets:summary:%s", matchID)
	}

	var summary models.FallOfWicketsSummary
	err := r.cache.GetOrSet(key, &summary, cache.FallOfWicketsTTL, func() (interface{}, error) {
		return r.repo.GetSummary(ctx, matchID, inningsID)
	})

	if err != nil {
		return nil, err
	}

	return &summary, nil
}

// GetWicketNumberForInnings gets the next wicket number for an innings (no caching - always fresh)
func (r *CachedFallOfWicketsRepository) GetWicketNumberForInnings(ctx context.Context, inningsID string) (int, error) {
	// This should always be fresh - no caching
	return r.repo.GetWicketNumberForInnings(ctx, inningsID)
}

// invalidateFallOfWicketsCaches invalidates all fall of wickets related caches
func (r *CachedFallOfWicketsRepository) invalidateFallOfWicketsCaches(matchID, inningsID string) {
	if r.cache == nil {
		return
	}

	// Invalidate match-level caches
	matchKey := fmt.Sprintf("fall_of_wickets:match:%s", matchID)
	_ = r.cache.Invalidate(matchKey)

	summaryKey := fmt.Sprintf("fall_of_wickets:summary:%s", matchID)
	_ = r.cache.Invalidate(summaryKey)

	// Invalidate innings-level caches
	if inningsID != "" {
		inningsKey := fmt.Sprintf("fall_of_wickets:innings:%s", inningsID)
		_ = r.cache.Invalidate(inningsKey)

		summaryInningsKey := fmt.Sprintf("fall_of_wickets:summary:%s:%s", matchID, inningsID)
		_ = r.cache.Invalidate(summaryInningsKey)
	}

	// Invalidate pattern-based caches
	pattern := fmt.Sprintf("fall_of_wickets:*%s*", matchID)
	_ = r.cache.InvalidatePattern(pattern)

	if inningsID != "" {
		inningsPattern := fmt.Sprintf("fall_of_wickets:*%s*", inningsID)
		_ = r.cache.InvalidatePattern(inningsPattern)
	}
}
