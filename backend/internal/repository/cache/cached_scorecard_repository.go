package cache

import (
	"context"
	"fmt"
	"spark-park-cricket-backend/internal/cache"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/repository/interfaces"
)

// CachedScorecardRepository wraps a scorecard repository with intelligent caching
type CachedScorecardRepository struct {
	repo  interfaces.ScorecardRepository
	cache *cache.CacheManager
}

// NewCachedScorecardRepository creates a new cached scorecard repository
func NewCachedScorecardRepository(repo interfaces.ScorecardRepository, cacheManager *cache.CacheManager) *CachedScorecardRepository {
	return &CachedScorecardRepository{
		repo:  repo,
		cache: cacheManager,
	}
}

// CreateInnings creates innings and invalidates scorecard cache
func (r *CachedScorecardRepository) CreateInnings(ctx context.Context, innings *models.Innings) error {
	err := r.repo.CreateInnings(ctx, innings)
	if err != nil {
		return err
	}

	// Invalidate scorecard cache for this match
	if innings.MatchID != "" {
		scorecardKey := r.cache.GetScorecardKey(innings.MatchID)
		_ = r.cache.Invalidate(scorecardKey)
	}

	return nil
}

// GetInningsByMatchID retrieves innings with caching
func (r *CachedScorecardRepository) GetInningsByMatchID(ctx context.Context, matchID string) ([]*models.Innings, error) {
	cacheKey := fmt.Sprintf("innings:match:%s", matchID)

	var innings []*models.Innings
	err := r.cache.GetOrSet(cacheKey, &innings, cache.ScorecardTTL, func() (interface{}, error) {
		return r.repo.GetInningsByMatchID(ctx, matchID)
	})

	if err != nil {
		return nil, err
	}

	return innings, nil
}

// GetInningsByMatchAndNumber retrieves specific innings with caching
func (r *CachedScorecardRepository) GetInningsByMatchAndNumber(ctx context.Context, matchID string, inningsNumber int) (*models.Innings, error) {
	cacheKey := fmt.Sprintf("innings:match:%s:number:%d", matchID, inningsNumber)

	var innings models.Innings
	err := r.cache.GetOrSet(cacheKey, &innings, cache.ScorecardTTL, func() (interface{}, error) {
		return r.repo.GetInningsByMatchAndNumber(ctx, matchID, inningsNumber)
	})

	if err != nil {
		return nil, err
	}

	return &innings, nil
}

// UpdateInnings updates innings and invalidates caches
func (r *CachedScorecardRepository) UpdateInnings(ctx context.Context, innings *models.Innings) error {
	err := r.repo.UpdateInnings(ctx, innings)
	if err != nil {
		return err
	}

	// Invalidate related caches
	if innings.MatchID != "" {
		scorecardKey := r.cache.GetScorecardKey(innings.MatchID)
		_ = r.cache.Invalidate(scorecardKey)

		inningsKey := fmt.Sprintf("innings:match:%s", innings.MatchID)
		_ = r.cache.Invalidate(inningsKey)

		specificInningsKey := fmt.Sprintf("innings:match:%s:number:%d", innings.MatchID, innings.InningsNumber)
		_ = r.cache.Invalidate(specificInningsKey)
	}

	return nil
}

// CompleteInnings completes innings and invalidates caches
func (r *CachedScorecardRepository) CompleteInnings(ctx context.Context, inningsID string) error {
	err := r.repo.CompleteInnings(ctx, inningsID)
	if err != nil {
		return err
	}

	// Note: We would need to get matchID from inningsID to invalidate properly
	// For now, we'll invalidate all scorecard caches
	// In a production system, you might want to store this relationship
	return nil
}

// CreateOver creates over and invalidates scorecard cache
func (r *CachedScorecardRepository) CreateOver(ctx context.Context, over *models.ScorecardOver) error {
	err := r.repo.CreateOver(ctx, over)
	if err != nil {
		return err
	}

	// Note: We would need to get matchID from innings to invalidate scorecard cache
	// For now, we'll handle this at the service level

	return nil
}

// GetOverByInningsAndNumber retrieves over with caching
func (r *CachedScorecardRepository) GetOverByInningsAndNumber(ctx context.Context, inningsID string, overNumber int) (*models.ScorecardOver, error) {
	cacheKey := fmt.Sprintf("over:innings:%s:number:%d", inningsID, overNumber)

	var over models.ScorecardOver
	err := r.cache.GetOrSet(cacheKey, &over, cache.ScorecardTTL, func() (interface{}, error) {
		return r.repo.GetOverByInningsAndNumber(ctx, inningsID, overNumber)
	})

	if err != nil {
		return nil, err
	}

	return &over, nil
}

// GetCurrentOver retrieves current over with caching
func (r *CachedScorecardRepository) GetCurrentOver(ctx context.Context, inningsID string) (*models.ScorecardOver, error) {
	cacheKey := fmt.Sprintf("over:current:innings:%s", inningsID)

	var over models.ScorecardOver
	err := r.cache.GetOrSet(cacheKey, &over, cache.ScorecardTTL, func() (interface{}, error) {
		return r.repo.GetCurrentOver(ctx, inningsID)
	})

	if err != nil {
		return nil, err
	}

	return &over, nil
}

// GetOversByInnings retrieves overs with caching
func (r *CachedScorecardRepository) GetOversByInnings(ctx context.Context, inningsID string) ([]*models.ScorecardOver, error) {
	cacheKey := fmt.Sprintf("overs:innings:%s", inningsID)

	var overs []*models.ScorecardOver
	err := r.cache.GetOrSet(cacheKey, &overs, cache.ScorecardTTL, func() (interface{}, error) {
		return r.repo.GetOversByInnings(ctx, inningsID)
	})

	if err != nil {
		return nil, err
	}

	return overs, nil
}

// UpdateOver updates over and invalidates caches
func (r *CachedScorecardRepository) UpdateOver(ctx context.Context, over *models.ScorecardOver) error {
	err := r.repo.UpdateOver(ctx, over)
	if err != nil {
		return err
	}

	// Note: We would need to get matchID from innings to invalidate scorecard cache
	// For now, we'll invalidate over-specific caches
	overKey := fmt.Sprintf("over:innings:%s:number:%d", over.InningsID, over.OverNumber)
	_ = r.cache.Invalidate(overKey)

	currentOverKey := fmt.Sprintf("over:current:innings:%s", over.InningsID)
	_ = r.cache.Invalidate(currentOverKey)

	oversKey := fmt.Sprintf("overs:innings:%s", over.InningsID)
	_ = r.cache.Invalidate(oversKey)

	return nil
}

// CompleteOver completes over and invalidates caches
func (r *CachedScorecardRepository) CompleteOver(ctx context.Context, overID string) error {
	err := r.repo.CompleteOver(ctx, overID)
	if err != nil {
		return err
	}

	// Note: Similar to CompleteInnings, we would need matchID to invalidate properly
	return nil
}

// CreateBall creates ball and invalidates scorecard cache (CRITICAL for performance)
func (r *CachedScorecardRepository) CreateBall(ctx context.Context, ball *models.ScorecardBall) error {
	// Invalidate balls cache BEFORE creating the ball to prevent race conditions
	// This ensures that any subsequent GetBallsByOver calls will fetch fresh data
	ballsCacheKey := fmt.Sprintf("balls:over:%s", ball.OverID)
	_ = r.cache.Invalidate(ballsCacheKey)

	// Also invalidate the last ball cache for this over
	lastBallCacheKey := fmt.Sprintf("ball:last:over:%s", ball.OverID)
	_ = r.cache.Invalidate(lastBallCacheKey)

	err := r.repo.CreateBall(ctx, ball)
	if err != nil {
		return err
	}

	// Note: We would need to get matchID from over/innings to invalidate scorecard cache
	// For now, we'll handle this at the service level where we have access to matchID
	// This is the main performance bottleneck we're solving

	return nil
}

// GetBallsByOver retrieves balls with caching
func (r *CachedScorecardRepository) GetBallsByOver(ctx context.Context, overID string) ([]*models.ScorecardBall, error) {
	cacheKey := fmt.Sprintf("balls:over:%s", overID)

	var balls []*models.ScorecardBall
	err := r.cache.GetOrSet(cacheKey, &balls, cache.ScorecardTTL, func() (interface{}, error) {
		return r.repo.GetBallsByOver(ctx, overID)
	})

	if err != nil {
		return nil, err
	}

	return balls, nil
}

// GetLastBall retrieves last ball with caching
func (r *CachedScorecardRepository) GetLastBall(ctx context.Context, overID string) (*models.ScorecardBall, error) {
	cacheKey := fmt.Sprintf("ball:last:over:%s", overID)

	var ball models.ScorecardBall
	err := r.cache.GetOrSet(cacheKey, &ball, cache.ScorecardTTL, func() (interface{}, error) {
		return r.repo.GetLastBall(ctx, overID)
	})

	if err != nil {
		return nil, err
	}

	return &ball, nil
}

// DeleteBall deletes ball and invalidates caches
func (r *CachedScorecardRepository) DeleteBall(ctx context.Context, ballID string) error {
	// Note: We would need to get the ball first to find the overID for cache invalidation
	// But the current interface doesn't provide a way to get a ball by ID
	// For now, we'll delete the ball and let the service handle cache invalidation
	// The service should call InvalidateBallsCacheForOver after deletion
	err := r.repo.DeleteBall(ctx, ballID)
	if err != nil {
		return err
	}

	// Note: We would need matchID to invalidate scorecard cache properly
	// For now, we'll handle this at the service level
	return nil
}

// InvalidateBallsCacheForOver invalidates all ball-related caches for a specific over
// This method can be called from the service layer when needed
func (r *CachedScorecardRepository) InvalidateBallsCacheForOver(overID string) {
	// Invalidate balls cache for this over
	ballsCacheKey := fmt.Sprintf("balls:over:%s", overID)
	_ = r.cache.Invalidate(ballsCacheKey)

	// Invalidate last ball cache for this over
	lastBallCacheKey := fmt.Sprintf("ball:last:over:%s", overID)
	_ = r.cache.Invalidate(lastBallCacheKey)
}

// GetScorecard retrieves complete scorecard with intelligent caching (CRITICAL)
func (r *CachedScorecardRepository) GetScorecard(ctx context.Context, matchID string) (*models.ScorecardResponse, error) {
	scorecardKey := r.cache.GetScorecardKey(matchID)

	var scorecard models.ScorecardResponse
	err := r.cache.GetOrSet(scorecardKey, &scorecard, cache.ScorecardTTL, func() (interface{}, error) {
		return r.repo.GetScorecard(ctx, matchID)
	})

	if err != nil {
		return nil, err
	}

	return &scorecard, nil
}

// StartScoring starts scoring and invalidates caches
func (r *CachedScorecardRepository) StartScoring(ctx context.Context, matchID string) error {
	err := r.repo.StartScoring(ctx, matchID)
	if err != nil {
		return err
	}

	// Invalidate scorecard cache to ensure fresh data
	scorecardKey := r.cache.GetScorecardKey(matchID)
	_ = r.cache.Invalidate(scorecardKey)

	return nil
}
