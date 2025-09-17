package cache

import (
	"context"
	"fmt"
	"spark-park-cricket-backend/internal/cache"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/repository/interfaces"
)

// CachedMatchRepository wraps a match repository with caching
type CachedMatchRepository struct {
	repo  interfaces.MatchRepository
	cache *cache.CacheManager
}

// NewCachedMatchRepository creates a new cached match repository
func NewCachedMatchRepository(repo interfaces.MatchRepository, cacheManager *cache.CacheManager) *CachedMatchRepository {
	return &CachedMatchRepository{
		repo:  repo,
		cache: cacheManager,
	}
}

// Create creates a new match and invalidates cache
func (r *CachedMatchRepository) Create(ctx context.Context, match *models.Match) error {
	err := r.repo.Create(ctx, match)
	if err != nil {
		return err
	}

	// Invalidate caches
	_ = r.cache.Invalidate("match:list")
	_ = r.cache.Invalidate("match:count")

	// Invalidate common pagination cache keys
	_ = r.cache.Invalidate("match:list:limit:20")
	_ = r.cache.Invalidate("match:list:limit:10")
	_ = r.cache.Invalidate("match:list:limit:5")
	_ = r.cache.Invalidate("match:list:limit:3")
	_ = r.cache.Invalidate("match:list:limit:2")

	// Invalidate all possible match list cache keys with different pagination parameters
	// This is necessary because the cache keys include limit and offset parameters
	_ = r.cache.InvalidatePattern("match:list:*")

	if match.SeriesID != "" {
		seriesKey := r.cache.GetMatchesBySeriesKey(match.SeriesID)
		_ = r.cache.Invalidate(seriesKey)
	}

	// Cache the new match
	if match.ID != "" {
		key := r.cache.GetMatchKey(match.ID)
		_ = r.cache.Set(key, match, cache.StaticDataTTL)
	}

	return nil
}

// GetByID retrieves a match by ID with caching
func (r *CachedMatchRepository) GetByID(ctx context.Context, id string) (*models.Match, error) {
	key := r.cache.GetMatchKey(id)

	var match models.Match
	err := r.cache.GetOrSet(key, &match, cache.StaticDataTTL, func() (interface{}, error) {
		return r.repo.GetByID(ctx, id)
	})

	if err != nil {
		return nil, err
	}

	return &match, nil
}

// GetAll retrieves all matches with caching
func (r *CachedMatchRepository) GetAll(ctx context.Context, filters *models.MatchFilters) ([]*models.Match, error) {
	// Create cache key based on filters
	cacheKey := "match:list"
	if filters != nil {
		if filters.SeriesID != nil && *filters.SeriesID != "" {
			cacheKey += fmt.Sprintf(":series:%s", *filters.SeriesID)
		}
		if filters.Status != nil {
			cacheKey += fmt.Sprintf(":status:%s", *filters.Status)
		}
		if filters.Limit > 0 {
			cacheKey += fmt.Sprintf(":limit:%d", filters.Limit)
		}
		if filters.Offset > 0 {
			cacheKey += fmt.Sprintf(":offset:%d", filters.Offset)
		}
	}

	var matches []*models.Match
	err := r.cache.GetOrSet(cacheKey, &matches, cache.MatchListTTL, func() (interface{}, error) {
		return r.repo.GetAll(ctx, filters)
	})

	if err != nil {
		return nil, err
	}

	return matches, nil
}

// GetBySeriesID retrieves matches by series ID with caching
func (r *CachedMatchRepository) GetBySeriesID(ctx context.Context, seriesID string) ([]*models.Match, error) {
	key := r.cache.GetMatchesBySeriesKey(seriesID)

	var matches []*models.Match
	err := r.cache.GetOrSet(key, &matches, cache.MatchListTTL, func() (interface{}, error) {
		return r.repo.GetBySeriesID(ctx, seriesID)
	})

	if err != nil {
		return nil, err
	}

	return matches, nil
}

// Update updates a match and invalidates cache
func (r *CachedMatchRepository) Update(ctx context.Context, id string, match *models.Match) error {
	err := r.repo.Update(ctx, id, match)
	if err != nil {
		return err
	}

	// Invalidate caches
	key := r.cache.GetMatchKey(id)
	_ = r.cache.Invalidate(key)
	_ = r.cache.Invalidate("match:list")
	_ = r.cache.Invalidate("match:count")

	// Invalidate common pagination cache keys
	_ = r.cache.Invalidate("match:list:limit:20")
	_ = r.cache.Invalidate("match:list:limit:10")
	_ = r.cache.Invalidate("match:list:limit:5")
	_ = r.cache.Invalidate("match:list:limit:3")
	_ = r.cache.Invalidate("match:list:limit:2")

	_ = r.cache.InvalidatePattern("match:list:*")

	if match.SeriesID != "" {
		seriesKey := r.cache.GetMatchesBySeriesKey(match.SeriesID)
		_ = r.cache.Invalidate(seriesKey)
	}

	// Update cache with new data
	_ = r.cache.Set(key, match, cache.StaticDataTTL)

	return nil
}

// Delete deletes a match and invalidates cache
func (r *CachedMatchRepository) Delete(ctx context.Context, id string) error {
	// Get match first to know which series to invalidate
	match, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	err = r.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Invalidate caches
	key := r.cache.GetMatchKey(id)
	_ = r.cache.Invalidate(key)
	_ = r.cache.Invalidate("match:list")
	_ = r.cache.Invalidate("match:count")

	// Invalidate common pagination cache keys
	_ = r.cache.Invalidate("match:list:limit:20")
	_ = r.cache.Invalidate("match:list:limit:10")
	_ = r.cache.Invalidate("match:list:limit:5")
	_ = r.cache.Invalidate("match:list:limit:3")
	_ = r.cache.Invalidate("match:list:limit:2")

	_ = r.cache.InvalidatePattern("match:list:*")

	if match.SeriesID != "" {
		seriesKey := r.cache.GetMatchesBySeriesKey(match.SeriesID)
		_ = r.cache.Invalidate(seriesKey)
	}

	return nil
}

// Count retrieves match count with caching
func (r *CachedMatchRepository) Count(ctx context.Context) (int64, error) {
	cacheKey := "match:count"

	var count int64
	err := r.cache.GetOrSet(cacheKey, &count, cache.MatchListTTL, func() (interface{}, error) {
		return r.repo.Count(ctx)
	})

	if err != nil {
		return 0, err
	}

	return count, nil
}

// GetNextMatchNumber retrieves next match number with caching
func (r *CachedMatchRepository) GetNextMatchNumber(ctx context.Context, seriesID string) (int, error) {
	cacheKey := fmt.Sprintf("match:next_number:series:%s", seriesID)

	var nextNumber int
	err := r.cache.GetOrSet(cacheKey, &nextNumber, cache.MatchListTTL, func() (interface{}, error) {
		return r.repo.GetNextMatchNumber(ctx, seriesID)
	})

	if err != nil {
		return 0, err
	}

	return nextNumber, nil
}

// ExistsBySeriesAndMatchNumber checks if match exists with caching
func (r *CachedMatchRepository) ExistsBySeriesAndMatchNumber(ctx context.Context, seriesID string, matchNumber int) (bool, error) {
	cacheKey := fmt.Sprintf("match:exists:series:%s:number:%d", seriesID, matchNumber)

	var exists bool
	err := r.cache.GetOrSet(cacheKey, &exists, cache.MatchListTTL, func() (interface{}, error) {
		return r.repo.ExistsBySeriesAndMatchNumber(ctx, seriesID, matchNumber)
	})

	if err != nil {
		return false, err
	}

	return exists, nil
}
