package cache

import (
	"context"
	"fmt"
	"spark-park-cricket-backend/internal/cache"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/repository/interfaces"
)

// CachedSeriesRepository wraps a series repository with caching
type CachedSeriesRepository struct {
	repo  interfaces.SeriesRepository
	cache *cache.CacheManager
}

// NewCachedSeriesRepository creates a new cached series repository
func NewCachedSeriesRepository(repo interfaces.SeriesRepository, cacheManager *cache.CacheManager) *CachedSeriesRepository {
	return &CachedSeriesRepository{
		repo:  repo,
		cache: cacheManager,
	}
}

// Create creates a new series and invalidates cache
func (r *CachedSeriesRepository) Create(ctx context.Context, series *models.Series) error {
	err := r.repo.Create(ctx, series)
	if err != nil {
		return err
	}

	// Invalidate series list cache - invalidate all possible cache keys
	keysToInvalidate := []string{
		"series:list:order:created_at_desc",
		"series:list:order:created_at_desc:limit:20",
		"series:list:order:created_at_desc:limit:50",
		"series:list:order:created_at_desc:limit:100",
		"series:list:order:created_at_desc:limit:20:offset:0",
		"series:list:order:created_at_desc:limit:20:offset:20",
		"series:list:order:created_at_desc:limit:20:offset:40",
		"series:list:order:created_at_desc:limit:50:offset:0",
		"series:list:order:created_at_desc:limit:100:offset:0",
	}

	for _, key := range keysToInvalidate {
		_ = r.cache.Invalidate(key)
	}

	// Cache the new series
	if series.ID != "" {
		key := r.cache.GetSeriesKey(series.ID)
		_ = r.cache.Set(key, series, cache.StaticDataTTL)
	}

	return nil
}

// GetByID retrieves a series by ID with caching
func (r *CachedSeriesRepository) GetByID(ctx context.Context, id string) (*models.Series, error) {
	key := r.cache.GetSeriesKey(id)

	var series models.Series
	err := r.cache.GetOrSet(key, &series, cache.StaticDataTTL, func() (interface{}, error) {
		return r.repo.GetByID(ctx, id)
	})

	if err != nil {
		return nil, err
	}

	return &series, nil
}

// GetAll retrieves all series with caching
func (r *CachedSeriesRepository) GetAll(ctx context.Context, filters *models.SeriesFilters) (*interfaces.PaginatedSeriesResult, error) {
	// Create cache key based on filters
	cacheKey := "series:list:order:created_at_desc"
	if filters != nil {
		// Add filter parameters to cache key
		if filters.Limit > 0 {
			cacheKey += fmt.Sprintf(":limit:%d", filters.Limit)
		}
		if filters.Offset > 0 {
			cacheKey += fmt.Sprintf(":offset:%d", filters.Offset)
		}
	}

	var result *interfaces.PaginatedSeriesResult
	err := r.cache.GetOrSet(cacheKey, &result, cache.MatchListTTL, func() (interface{}, error) {
		paginatedResult, err := r.repo.GetAll(ctx, filters)
		if err != nil {
			return nil, err
		}
		return paginatedResult, err
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// Update updates a series and invalidates cache
func (r *CachedSeriesRepository) Update(ctx context.Context, id string, series *models.Series) error {
	err := r.repo.Update(ctx, id, series)
	if err != nil {
		return err
	}

	// Invalidate series-specific cache
	key := r.cache.GetSeriesKey(id)
	_ = r.cache.Invalidate(key)

	// Invalidate all possible series list cache keys
	keysToInvalidate := []string{
		"series:list",
		"series:list:order:created_at_desc",
		"series:list:order:created_at_desc:limit:20",
		"series:list:order:created_at_desc:limit:50",
		"series:list:order:created_at_desc:limit:100",
		"series:list:order:created_at_desc:limit:20:offset:0",
		"series:list:order:created_at_desc:limit:20:offset:20",
		"series:list:order:created_at_desc:limit:20:offset:40",
		"series:list:order:created_at_desc:limit:50:offset:0",
		"series:list:order:created_at_desc:limit:100:offset:0",
	}

	for _, key := range keysToInvalidate {
		_ = r.cache.Invalidate(key)
	}

	// Update cache with new data
	_ = r.cache.Set(key, series, cache.StaticDataTTL)

	return nil
}

// Delete deletes a series and invalidates cache
func (r *CachedSeriesRepository) Delete(ctx context.Context, id string) error {
	err := r.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Invalidate series-specific cache
	key := r.cache.GetSeriesKey(id)
	_ = r.cache.Invalidate(key)

	// Invalidate all possible series list cache keys
	keysToInvalidate := []string{
		"series:list",
		"series:list:order:created_at_desc",
		"series:list:order:created_at_desc:limit:20",
		"series:list:order:created_at_desc:limit:50",
		"series:list:order:created_at_desc:limit:100",
		"series:list:order:created_at_desc:limit:20:offset:0",
		"series:list:order:created_at_desc:limit:20:offset:20",
		"series:list:order:created_at_desc:limit:20:offset:40",
		"series:list:order:created_at_desc:limit:50:offset:0",
		"series:list:order:created_at_desc:limit:100:offset:0",
	}

	for _, key := range keysToInvalidate {
		_ = r.cache.Invalidate(key)
	}

	return nil
}

// Count retrieves series count with caching
func (r *CachedSeriesRepository) Count(ctx context.Context) (int64, error) {
	cacheKey := "series:count"

	var count int64
	err := r.cache.GetOrSet(cacheKey, &count, cache.MatchListTTL, func() (interface{}, error) {
		return r.repo.Count(ctx)
	})

	if err != nil {
		return 0, err
	}

	return count, nil
}
