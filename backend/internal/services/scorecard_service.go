package services

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"spark-park-cricket-backend/internal/cache"
	contextkeys "spark-park-cricket-backend/internal/context"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/monitoring"
	"spark-park-cricket-backend/internal/repository/interfaces"
	"spark-park-cricket-backend/internal/utils"
	"strings"
	"sync"
	"time"
)

type ScorecardService struct {
	scorecardRepo interfaces.ScorecardRepository
	matchRepo     interfaces.MatchRepository
	metrics       *monitoring.Metrics
	cache         *cache.CacheManager
	// Distributed mutex for ball addition to prevent race conditions
	ballMutexes sync.Map // map[string]*sync.Mutex for match_innings keys
	// Distributed mutex for over creation to prevent race conditions
	overMutexes sync.Map // map[string]*sync.Mutex for match_innings keys
}

// NewScorecardService creates a new scorecard service
func NewScorecardService(scorecardRepo interfaces.ScorecardRepository, matchRepo interfaces.MatchRepository, metrics *monitoring.Metrics, cache *cache.CacheManager) *ScorecardService {
	return &ScorecardService{
		scorecardRepo: scorecardRepo,
		matchRepo:     matchRepo,
		metrics:       metrics,
		cache:         cache,
		ballMutexes:   sync.Map{},
		overMutexes:   sync.Map{},
	}
}

// getBallAdditionMutex returns a mutex for the specific match and innings combination
// This prevents race conditions when multiple requests try to add balls to the same innings
func (s *ScorecardService) getBallAdditionMutex(matchID string, inningsNumber int) *sync.Mutex {
	key := fmt.Sprintf("ball_addition_%s_%d", matchID, inningsNumber)

	mutex, _ := s.ballMutexes.LoadOrStore(key, &sync.Mutex{})
	return mutex.(*sync.Mutex)
}

// getOverCreationMutex returns a mutex for the specific match and innings combination
// This prevents race conditions when multiple requests try to create overs for the same innings
func (s *ScorecardService) getOverCreationMutex(matchID string, inningsNumber int) *sync.Mutex {
	key := fmt.Sprintf("over_creation_%s_%d", matchID, inningsNumber)

	mutex, _ := s.overMutexes.LoadOrStore(key, &sync.Mutex{})
	return mutex.(*sync.Mutex)
}

// isRetryableError checks if an error is retryable based on common retryable error patterns
func (s *ScorecardService) isRetryableError(err error) bool {
	errStr := err.Error()
	return strings.Contains(errStr, "duplicate key value violates unique constraint") ||
		strings.Contains(errStr, "violates foreign key constraint") ||
		strings.Contains(errStr, "connection reset") ||
		strings.Contains(errStr, "timeout") ||
		strings.Contains(errStr, "temporary failure") ||
		strings.Contains(errStr, "deadlock detected")
}

// retryWithExponentialBackoff executes an operation with exponential backoff and jitter
func (s *ScorecardService) retryWithExponentialBackoff(ctx context.Context, operation func() error) error {
	maxRetries := 5
	baseDelay := 10 * time.Millisecond

	for attempt := 1; attempt <= maxRetries; attempt++ {
		err := operation()
		if err == nil {
			return nil
		}

		// Only retry on specific errors
		if !s.isRetryableError(err) {
			return err
		}

		if attempt == maxRetries {
			return fmt.Errorf("failed after %d attempts: %w", maxRetries, err)
		}

		// Exponential backoff with jitter
		delay := baseDelay * time.Duration(1<<uint(attempt-1))
		jitter := time.Duration(rand.Intn(10)) * time.Millisecond

		log.Printf("Retryable error on attempt %d/%d, retrying in %v: %v", attempt, maxRetries, delay+jitter, err)
		time.Sleep(delay + jitter)
	}

	return fmt.Errorf("unexpected retry logic error")
}

// StartScoring starts scoring for a match
func (s *ScorecardService) StartScoring(ctx context.Context, matchID string) error {
	log.Printf("Starting scoring for match %s", matchID)

	// Get user ID from context
	userID, ok := ctx.Value(contextkeys.UserIDKey).(string)
	if !ok || userID == "" {
		return fmt.Errorf("user authentication required")
	}

	// Get match details
	match, err := s.matchRepo.GetByID(ctx, matchID)
	if err != nil {
		log.Printf("Error getting match: %v", err)
		return fmt.Errorf("match not found: %w", err)
	}

	// Check ownership - user must be the creator of the series
	if match.CreatedBy != userID {
		return fmt.Errorf("access denied: you can only start scoring for matches in series you created")
	}

	// Check if match is not_started or live
	if match.Status != models.MatchStatusNotStarted && match.Status != models.MatchStatusLive {
		return fmt.Errorf("match must be not started or live to begin scoring")
	}

	// If match is not_started, transition it to live
	if match.Status == models.MatchStatusNotStarted {
		now := time.Now()
		match.Status = models.MatchStatusLive
		match.StartTime = &now
		match.UpdatedAt = now

		err = s.matchRepo.Update(ctx, matchID, match)
		if err != nil {
			log.Printf("Error updating match status: %v", err)
			return fmt.Errorf("failed to start match: %w", err)
		}

		log.Printf("Match %s status changed from not_started to live at %s", matchID, now.Format(time.RFC3339))
	}

	// Check if scoring is already started
	innings, err := s.scorecardRepo.GetInningsByMatchID(ctx, matchID)
	if err == nil && len(innings) > 0 {
		return fmt.Errorf("scoring already started for this match")
	}

	// Create first innings with toss winner as batting team
	now := time.Now()
	firstInnings := &models.Innings{
		MatchID:       matchID,
		InningsNumber: 1,
		BattingTeam:   match.TossWinner, // Toss winner bats first
		TotalRuns:     0,
		TotalWickets:  0,
		TotalOvers:    0.0,
		TotalBalls:    0,
		Status:        string(models.InningsStatusInProgress),
		StartTime:     &now,
	}

	err = s.scorecardRepo.CreateInnings(ctx, firstInnings)
	if err != nil {
		log.Printf("Error creating first innings: %v", err)
		return fmt.Errorf("failed to start scoring: %w", err)
	}

	log.Printf("Successfully started scoring for match %s, first innings batting team: %s", matchID, match.TossWinner)
	return nil
}

// AddBall adds a ball to the scorecard with distributed mutex to prevent race conditions
func (s *ScorecardService) AddBall(ctx context.Context, req *models.BallEventRequest) error {
	// Get distributed mutex for this match and innings combination
	mutex := s.getBallAdditionMutex(req.MatchID, req.InningsNumber)
	mutex.Lock()
	defer mutex.Unlock()

	// Use optimized method for better performance
	return s.AddBallOptimized(ctx, req)
}

// AddBallOptimized adds a ball using optimized data fetching
func (s *ScorecardService) AddBallOptimized(ctx context.Context, req *models.BallEventRequest) error {
	// Get user ID from context
	userID, ok := ctx.Value(contextkeys.UserIDKey).(string)
	if !ok || userID == "" {
		return fmt.Errorf("user authentication required")
	}

	// Validate ball event
	if err := utils.ValidateBallEventRequest(req); err != nil {
		log.Printf("Invalid ball event: %v", err)
		return fmt.Errorf("invalid ball event: %w", err)
	}

	// Get all necessary data in a single optimized call
	data, err := s.scorecardRepo.GetMatchInningsOverData(ctx, req.MatchID, req.InningsNumber)
	if err != nil {
		log.Printf("Error getting optimized match data: %v", err)

		// Check if the error is due to non-existent innings
		if strings.Contains(err.Error(), "innings not found") {
			// Check if this is trying to add to second innings before first is complete
			if req.InningsNumber == 2 {
				// Check if first innings exists and is not complete
				firstInnings, err := s.scorecardRepo.GetInningsByMatchAndNumber(ctx, req.MatchID, 1)
				if err == nil && firstInnings.Status != string(models.InningsStatusCompleted) {
					return fmt.Errorf("first innings is not complete, cannot start second innings")
				}
			}
			return fmt.Errorf("innings not found")
		}

		return fmt.Errorf("failed to get match data: %w", err)
	}

	// Validate user access
	if data.CreatedBy != userID {
		return fmt.Errorf("access denied: you can only score balls for matches you created")
	}

	// Validate match status
	if data.MatchStatus != models.MatchStatusLive {
		return fmt.Errorf("match is not live, cannot add ball")
	}

	// Validate innings status
	if data.InningsStatus != models.InningsStatusInProgress {
		return fmt.Errorf("innings is not in progress, cannot add ball")
	}

	// Validate over status
	if data.OverStatus != models.OverStatusInProgress {
		return fmt.Errorf("over is not in progress, cannot add ball")
	}

	// Check if over is complete (6 legal balls)
	if data.LegalBallCount >= 6 {
		return fmt.Errorf("over is complete, cannot add more balls")
	}

	// CRITICAL PRE-VALIDATION: Check if adding this wicket would exceed maximum wickets
	// This check MUST happen BEFORE adding the ball to prevent balls being added after all wickets are down.
	// Issue fixed: Without this check, the system allowed 24 balls for 19 wickets (5 extra balls after last wicket).
	// DO NOT REMOVE: This prevents innings from continuing after all wickets have fallen.
	// The post-validation (after adding ball) only marks innings as complete but doesn't prevent the ball from being added.
	if req.IsWicket {
		maxWickets := data.TeamAPlayerCount - 1 // n-1 wickets for n players
		if data.InningsTotalWickets >= maxWickets {
			return fmt.Errorf("innings is complete, all wickets are down (%d/%d)", data.InningsTotalWickets, maxWickets)
		}
	}

	// Calculate next ball number - use fresh calculation to avoid stale data
	nextBallNumber, err := s.getNextBallNumber(ctx, data.OverID)
	if err != nil {
		log.Printf("Error getting next ball number: %v", err)
		return fmt.Errorf("failed to get next ball number: %w", err)
	}

	// Calculate runs from run type and byes
	runs := req.RunType.GetRunValue()
	byes := req.Byes
	if req.IsWicket && req.RunType == models.RunTypeWC {
		runs = 0 // Wicket doesn't count as runs
	}

	// Total runs = ball runs + byes
	totalRuns := runs + byes

	// Create ball
	ball := &models.ScorecardBall{
		OverID:     data.OverID,
		BallNumber: nextBallNumber,
		BallType:   req.BallType,
		RunType:    req.RunType,
		Runs:       runs,
		Byes:       byes,
		IsWicket:   req.IsWicket,
		WicketType: req.WicketType,
	}

	// Invalidate caches BEFORE database operations to prevent stale reads
	s.invalidateMatchCachesBeforeWrite(ctx, req.MatchID, data.InningsID, data.OverID)

	// Record ball addition metrics with retry logic for constraint violations
	start := time.Now()
	err = monitoring.WithDatabaseMonitoringContext(
		ctx, s.metrics, "INSERT", "balls", req.MatchID,
		func(ctx context.Context) error {
			// Use improved retry logic with exponential backoff and jitter
			return s.retryWithExponentialBackoff(ctx, func() error {
				err := s.scorecardRepo.CreateBall(ctx, ball)
				if err == nil {
					return nil
				}

				// If it's a constraint violation, try to get fresh data and retry
				if strings.Contains(err.Error(), "duplicate key value violates unique constraint") ||
					strings.Contains(err.Error(), "violates foreign key constraint") {
					log.Printf("Ball creation failed due to constraint violation, retrying with fresh data: %v", err)

					// Get fresh match data to ensure we have the latest over information
					freshData, freshErr := s.scorecardRepo.GetMatchInningsOverData(ctx, req.MatchID, req.InningsNumber)
					if freshErr != nil {
						log.Printf("Failed to get fresh match data: %v", freshErr)
						return fmt.Errorf("failed to get fresh match data: %w", freshErr)
					}

					// Update ball with fresh over ID and get fresh ball number
					ball.OverID = freshData.OverID
					freshBallNumber, freshErr := s.getNextBallNumber(ctx, ball.OverID)
					if freshErr != nil {
						log.Printf("Failed to get fresh ball number: %v", freshErr)
						return fmt.Errorf("failed to get fresh ball number: %w", freshErr)
					}
					ball.BallNumber = freshBallNumber

					// Return the original error to trigger retry
					return err
				}

				// If it's not a constraint violation, return the error immediately
				return err
			})
		},
	)
	duration := time.Since(start)

	// Record cricket-specific metrics
	s.metrics.RecordBallAddition(
		req.MatchID,
		fmt.Sprintf("innings_%d", req.InningsNumber),
		string(req.BallType),
		string(req.RunType),
		duration,
	)

	if err != nil {
		log.Printf("Error creating ball: %v", err)
		return fmt.Errorf("failed to add ball: %w", err)
	}

	// Update over statistics (optimized in-memory calculation)
	over := &models.ScorecardOver{
		ID:           data.OverID,
		InningsID:    data.InningsID,
		OverNumber:   data.OverNumber,
		TotalRuns:    data.OverTotalRuns + totalRuns,
		TotalBalls:   data.OverTotalBalls,
		TotalWickets: data.OverTotalWickets,
		Status:       string(data.OverStatus),
	}

	// Only count legal balls (good balls) for over completion
	if req.BallType == models.BallTypeGood {
		over.TotalBalls++
	}
	if req.IsWicket {
		over.TotalWickets++
	}

	// Check if over is complete (6 legal balls)
	// Note: Removed hardcoded wicket check - innings completion is handled separately
	if over.TotalBalls >= 6 {
		over.Status = string(models.OverStatusCompleted)
	}

	err = monitoring.WithDatabaseMonitoringContext(
		ctx, s.metrics, "UPDATE", "overs", req.MatchID,
		func(ctx context.Context) error {
			return s.scorecardRepo.UpdateOver(ctx, over)
		},
	)
	if err != nil {
		log.Printf("Error updating over: %v", err)
		return fmt.Errorf("failed to update over: %w", err)
	}

	// Update innings statistics (in-memory calculation)
	innings := &models.Innings{
		ID:            data.InningsID,
		MatchID:       req.MatchID,
		InningsNumber: data.InningsNumber,
		BattingTeam:   data.BattingTeam,
		TotalRuns:     data.InningsTotalRuns + totalRuns,
		TotalWickets:  data.InningsTotalWickets,
		TotalOvers:    data.InningsTotalOvers,
		TotalBalls:    data.InningsTotalBalls,
		Status:        string(data.InningsStatus),
	}

	// Only count legal balls for innings overs calculation
	if req.BallType == models.BallTypeGood {
		innings.TotalBalls++
	}
	if req.IsWicket {
		innings.TotalWickets++
	}

	// Calculate total overs properly (optimized in-memory calculation)
	completedOvers := data.CompletedOvers
	currentOverBalls := over.TotalBalls
	innings.TotalOvers = s.calculateOversInMemory(completedOvers, currentOverBalls)

	// Check if innings is complete (optimized in-memory calculation)
	maxWickets := data.TeamAPlayerCount - 1 // n-1 wickets for n players
	if req.InningsNumber == 1 {
		innings.Status = s.calculateInningsStatusInMemory(innings, maxWickets, data.TotalOvers)
		if innings.Status == string(models.InningsStatusCompleted) {
			log.Printf("First innings %d completed for match %s: wickets=%d/%d, overs=%.1f/%d",
				innings.InningsNumber, req.MatchID, innings.TotalWickets, maxWickets, innings.TotalOvers, data.TotalOvers)
		}
	}

	err = monitoring.WithDatabaseMonitoringContext(
		ctx, s.metrics, "UPDATE", "innings", req.MatchID,
		func(ctx context.Context) error {
			return s.scorecardRepo.UpdateInnings(ctx, innings)
		},
	)
	if err != nil {
		log.Printf("Error updating innings: %v", err)
		return fmt.Errorf("failed to update innings: %w", err)
	}

	// Handle match progression
	if req.InningsNumber == 1 {
		// First innings - check if completed and start second innings
		if innings.Status == string(models.InningsStatusCompleted) {
			err = s.startSecondInnings(ctx, req.MatchID, &models.Match{
				ID:               req.MatchID,
				Status:           models.MatchStatusLive,
				BattingTeam:      data.BattingTeam,
				TotalOvers:       data.TotalOvers,
				TeamAPlayerCount: data.TeamAPlayerCount,
			})
			if err != nil {
				log.Printf("Error starting second innings: %v", err)
				return fmt.Errorf("failed to start second innings: %w", err)
			}
			log.Printf("Second innings started for match %s", req.MatchID)
		}
	} else if req.InningsNumber == 2 {
		// Second innings - check for match completion after every ball
		shouldCompleteMatch, reason := s.ShouldCompleteMatch(ctx, req.MatchID, innings, &models.Match{
			ID:               req.MatchID,
			Status:           models.MatchStatusLive,
			BattingTeam:      data.BattingTeam,
			TotalOvers:       data.TotalOvers,
			TeamAPlayerCount: data.TeamAPlayerCount,
		})
		if shouldCompleteMatch {
			// Complete the innings first
			innings.Status = string(models.InningsStatusCompleted)
			// Set end time and calculate duration
			now := time.Now()
			innings.EndTime = &now
			if innings.StartTime != nil {
				innings.DurationSeconds = int(now.Sub(*innings.StartTime).Seconds())
			}
			err = s.scorecardRepo.UpdateInnings(ctx, innings)
			if err != nil {
				return fmt.Errorf("failed to update innings status: %w", err)
			}

			// Complete the match
			// Fetch the complete match data to ensure all fields are populated
			completeMatch, err := s.matchRepo.GetByID(ctx, req.MatchID)
			if err != nil {
				log.Printf("Error fetching match data for completion: %v", err)
				return fmt.Errorf("failed to fetch match data: %w", err)
			}

			now = time.Now()
			completeMatch.Status = models.MatchStatusCompleted
			completeMatch.EndTime = &now
			completeMatch.UpdatedAt = now
			err = s.matchRepo.Update(ctx, req.MatchID, completeMatch)
			if err != nil {
				return fmt.Errorf("failed to complete match: %w", err)
			}
			log.Printf("Match %s completed - %s", req.MatchID, reason)
		}
	}

	// Invalidate match-related caches after successful ball addition (async)
	go s.invalidateMatchCachesAsync(req.MatchID, data.InningsID, data.OverID)

	log.Printf("Successfully added ball: %s %d runs, byes: %d, total: %d, wicket: %v", req.RunType, runs, byes, totalRuns, req.IsWicket)
	return nil
}

// AddBallLegacy is the original AddBall method (kept for reference)
func (s *ScorecardService) AddBallLegacy(ctx context.Context, req *models.BallEventRequest) error {
	// Get user ID from context
	userID, ok := ctx.Value(contextkeys.UserIDKey).(string)
	if !ok || userID == "" {
		return fmt.Errorf("user authentication required")
	}

	// Validate ball event
	if err := utils.ValidateBallEventRequest(req); err != nil {
		log.Printf("Invalid ball event: %v", err)
		return fmt.Errorf("invalid ball event: %w", err)
	}

	// Get match details with monitoring
	var match *models.Match
	err := monitoring.WithDatabaseMonitoringContext(
		ctx, s.metrics, "SELECT", "matches", req.MatchID,
		func(ctx context.Context) error {
			var dbErr error
			match, dbErr = s.matchRepo.GetByID(ctx, req.MatchID)
			return dbErr
		},
	)
	if err != nil {
		log.Printf("Error getting match: %v", err)
		return fmt.Errorf("match not found: %w", err)
	}

	// Check ownership - user must be the creator of the match
	log.Printf("Ownership check - User ID: %s, Match CreatedBy: %s", userID, match.CreatedBy)
	if match.CreatedBy != userID {
		return fmt.Errorf("access denied: you can only score balls for matches you created")
	}

	// Check if match is live
	if match.Status != models.MatchStatusLive {
		return fmt.Errorf("match is not live, cannot add ball")
	}

	// Validate innings order
	log.Printf("DEBUG: Starting innings validation for match %s, innings %d, batting team %s, toss winner %s",
		req.MatchID, req.InningsNumber, match.BattingTeam, match.TossWinner)
	if err := s.ValidateInningsOrder(ctx, req.MatchID, match, req.InningsNumber); err != nil {
		log.Printf("DEBUG: Innings validation failed: %v", err)
		return fmt.Errorf("innings validation failed: %w", err)
	}
	log.Printf("DEBUG: Innings validation passed for match %s, innings %d", req.MatchID, req.InningsNumber)

	// Get innings or create if doesn't exist with monitoring
	var innings *models.Innings
	err = monitoring.WithDatabaseMonitoringContext(
		ctx, s.metrics, "SELECT", "innings", req.MatchID,
		func(ctx context.Context) error {
			var dbErr error
			innings, dbErr = s.scorecardRepo.GetInningsByMatchAndNumber(ctx, req.MatchID, req.InningsNumber)
			return dbErr
		},
	)
	if err != nil {
		log.Printf("Innings not found, creating new innings: %v", err)
		// Create innings if it doesn't exist
		now := time.Now()
		innings = &models.Innings{
			MatchID:       req.MatchID,
			InningsNumber: req.InningsNumber,
			BattingTeam:   match.BattingTeam,
			TotalRuns:     0,
			TotalWickets:  0,
			TotalOvers:    0.0,
			TotalBalls:    0,
			Status:        string(models.InningsStatusInProgress),
			StartTime:     &now,
		}
		err = monitoring.WithDatabaseMonitoringContext(
			ctx, s.metrics, "INSERT", "innings", req.MatchID,
			func(ctx context.Context) error {
				return s.scorecardRepo.CreateInnings(ctx, innings)
			},
		)
		if err != nil {
			log.Printf("Error creating innings: %v", err)
			return fmt.Errorf("failed to create innings: %w", err)
		}
		log.Printf("Created new innings %d for match %s", req.InningsNumber, req.MatchID)
	}

	// Check if innings is in progress
	if innings.Status != string(models.InningsStatusInProgress) {
		return fmt.Errorf("innings is not in progress, cannot add ball")
	}

	// Get current over or create new one with monitoring and mutex protection
	// Also get all overs for this innings to avoid duplicate query later
	var over *models.ScorecardOver
	var allOvers []*models.ScorecardOver

	// Use distributed mutex for over creation to prevent race conditions
	overMutex := s.getOverCreationMutex(req.MatchID, req.InningsNumber)
	overMutex.Lock()
	defer overMutex.Unlock()

	err = monitoring.WithDatabaseMonitoringContext(
		ctx, s.metrics, "SELECT", "overs", req.MatchID,
		func(ctx context.Context) error {
			var dbErr error
			over, allOvers, dbErr = s.getCurrentOverWithOvers(ctx, innings.ID)
			return dbErr
		},
	)
	if err != nil {
		log.Printf("Error getting current over: %v", err)
		return fmt.Errorf("failed to get current over: %w", err)
	}

	// Check if over is in progress
	if over.Status != string(models.OverStatusInProgress) {
		return fmt.Errorf("over is not in progress, cannot add ball")
	}

	// Get next ball number
	ballNumber, err := s.getNextBallNumber(ctx, over.ID)
	if err != nil {
		log.Printf("Error getting next ball number: %v", err)
		return fmt.Errorf("failed to get next ball number: %w", err)
	}

	// Calculate runs from run type and byes
	runs := req.RunType.GetRunValue()
	byes := req.Byes
	if req.IsWicket && req.RunType == models.RunTypeWC {
		runs = 0 // Wicket doesn't count as runs
	}

	// Total runs = ball runs + byes
	totalRuns := runs + byes

	// Create ball
	ball := &models.ScorecardBall{
		OverID:     over.ID,
		BallNumber: ballNumber,
		BallType:   req.BallType,
		RunType:    req.RunType,
		Runs:       runs,
		Byes:       byes,
		IsWicket:   req.IsWicket,
		WicketType: req.WicketType,
	}

	// Record ball addition metrics with retry logic for constraint violations
	start := time.Now()
	err = monitoring.WithDatabaseMonitoringContext(
		ctx, s.metrics, "INSERT", "balls", req.MatchID,
		func(ctx context.Context) error {
			// Retry logic for constraint violations
			maxRetries := 5
			for attempt := 1; attempt <= maxRetries; attempt++ {
				err := s.scorecardRepo.CreateBall(ctx, ball)
				if err == nil {
					return nil
				}

				// Check if it's a constraint violation
				if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
					log.Printf("Ball creation failed due to constraint violation (attempt %d/%d), retrying with new ball number", attempt, maxRetries)

					// Get fresh ball number with exponential backoff
					freshBallNumber, err := s.getNextBallNumber(ctx, over.ID)
					if err != nil {
						log.Printf("Failed to get fresh ball number on attempt %d: %v", attempt, err)
						if attempt == maxRetries {
							return fmt.Errorf("failed to get fresh ball number after %d attempts: %w", maxRetries, err)
						}
						// Continue to retry
						time.Sleep(time.Millisecond * time.Duration(attempt*10))
						continue
					}
					ball.BallNumber = freshBallNumber

					if attempt == maxRetries {
						return fmt.Errorf("failed to create ball after %d attempts due to constraint violations", maxRetries)
					}

					// Exponential backoff delay before retry
					time.Sleep(time.Millisecond * time.Duration(attempt*10))
					continue
				}

				// If it's not a constraint violation, return the error immediately
				return err
			}
			return fmt.Errorf("failed to create ball after %d attempts", maxRetries)
		},
	)
	duration := time.Since(start)

	// Record cricket-specific metrics
	s.metrics.RecordBallAddition(
		req.MatchID,
		fmt.Sprintf("innings_%d", req.InningsNumber),
		string(req.BallType),
		string(req.RunType),
		duration,
	)

	if err != nil {
		log.Printf("Error creating ball: %v", err)
		return fmt.Errorf("failed to add ball: %w", err)
	}

	// Update over statistics
	over.TotalRuns += totalRuns
	// Only count legal balls (good balls) for over completion
	if req.BallType == models.BallTypeGood {
		over.TotalBalls++
	}
	if req.IsWicket {
		over.TotalWickets++
	}

	// Check if over is complete (6 legal balls)
	// Note: Removed hardcoded wicket check - innings completion is handled separately
	if over.TotalBalls >= 6 {
		over.Status = string(models.OverStatusCompleted)
		// Set end time and calculate duration
		now := time.Now()
		over.EndTime = &now
		if over.StartTime != nil {
			over.DurationSeconds = int(now.Sub(*over.StartTime).Seconds())
		}
	}

	err = monitoring.WithDatabaseMonitoringContext(
		ctx, s.metrics, "UPDATE", "overs", req.MatchID,
		func(ctx context.Context) error {
			return s.scorecardRepo.UpdateOver(ctx, over)
		},
	)
	if err != nil {
		log.Printf("Error updating over: %v", err)
		return fmt.Errorf("failed to update over: %w", err)
	}

	// Update innings statistics
	innings.TotalRuns += totalRuns
	// Only count legal balls for innings overs calculation
	if req.BallType == models.BallTypeGood {
		innings.TotalBalls++
	}
	if req.IsWicket {
		innings.TotalWickets++
	}

	// Calculate total overs properly
	// Use cached overs from getCurrentOverWithOvers to avoid duplicate query
	overs := allOvers

	completedOvers := 0
	currentOverBalls := 0

	for _, over := range overs {
		if over.Status == string(models.OverStatusCompleted) {
			completedOvers++
		} else if over.Status == string(models.OverStatusInProgress) {
			currentOverBalls = over.TotalBalls
		}
	}

	// Total overs = completed overs + current over balls as decimal
	// In cricket scoring: 3 balls = 0.3 overs, 4 balls = 0.4 overs, etc.
	// We need to convert balls to the cricket scoring format
	var currentOverDecimal float64
	if currentOverBalls > 0 {
		// Convert balls to cricket scoring format (0.1, 0.2, 0.3, 0.4, 0.5, 1.0)
		if currentOverBalls == 6 {
			currentOverDecimal = 1.0
		} else {
			currentOverDecimal = float64(currentOverBalls) / 10.0
		}
	}
	innings.TotalOvers = float64(completedOvers) + currentOverDecimal

	// Check if innings is complete
	// For first innings: complete when all wickets are taken or all overs are completed
	// For second innings: completion is handled by shouldCompleteMatch method
	maxWickets := match.TeamAPlayerCount - 1 // n-1 wickets for n players
	if req.InningsNumber == 1 {
		if innings.TotalWickets >= maxWickets || innings.TotalOvers >= float64(match.TotalOvers) {
			innings.Status = string(models.InningsStatusCompleted)
			// Set end time and calculate duration
			now := time.Now()
			innings.EndTime = &now
			if innings.StartTime != nil {
				innings.DurationSeconds = int(now.Sub(*innings.StartTime).Seconds())
			}
			log.Printf("First innings %d completed for match %s: wickets=%d/%d, overs=%.1f/%d",
				innings.InningsNumber, match.ID, innings.TotalWickets, maxWickets, innings.TotalOvers, match.TotalOvers)
		}
	}
	// For second innings, we don't automatically complete here - let shouldCompleteMatch handle it

	err = monitoring.WithDatabaseMonitoringContext(
		ctx, s.metrics, "UPDATE", "innings", req.MatchID,
		func(ctx context.Context) error {
			return s.scorecardRepo.UpdateInnings(ctx, innings)
		},
	)
	if err != nil {
		log.Printf("Error updating innings: %v", err)
		return fmt.Errorf("failed to update innings: %w", err)
	}

	// Handle match progression
	if req.InningsNumber == 1 {
		// First innings - check if completed and start second innings
		if innings.Status == string(models.InningsStatusCompleted) {
			err = s.startSecondInnings(ctx, req.MatchID, match)
			if err != nil {
				log.Printf("Error starting second innings: %v", err)
				return fmt.Errorf("failed to start second innings: %w", err)
			}
			log.Printf("Second innings started for match %s", req.MatchID)
		}
	} else if req.InningsNumber == 2 {
		// Second innings - check for match completion after every ball
		shouldCompleteMatch, reason := s.ShouldCompleteMatch(ctx, req.MatchID, innings, match)
		if shouldCompleteMatch {
			// Complete the innings first
			innings.Status = string(models.InningsStatusCompleted)
			// Set end time and calculate duration
			now := time.Now()
			innings.EndTime = &now
			if innings.StartTime != nil {
				innings.DurationSeconds = int(now.Sub(*innings.StartTime).Seconds())
			}
			err = s.scorecardRepo.UpdateInnings(ctx, innings)
			if err != nil {
				return fmt.Errorf("failed to update innings status: %w", err)
			}

			// Complete the match
			matchEndTime := time.Now()
			match.Status = models.MatchStatusCompleted
			match.EndTime = &matchEndTime
			match.UpdatedAt = matchEndTime
			err = s.matchRepo.Update(ctx, req.MatchID, match)
			if err != nil {
				return fmt.Errorf("failed to complete match: %w", err)
			}
			log.Printf("Match %s completed at %s - %s", req.MatchID, matchEndTime.Format(time.RFC3339), reason)
		}
	}

	// Invalidate match-related caches after successful ball addition
	s.invalidateMatchCaches(ctx, req.MatchID, innings.ID, over.ID)

	log.Printf("Successfully added ball: %s %d runs, byes: %d, total: %d, wicket: %v", req.RunType, runs, byes, totalRuns, req.IsWicket)
	return nil
}

// GetTimeTracking retrieves time tracking data for a match
func (s *ScorecardService) GetTimeTracking(ctx context.Context, matchID string) (*models.TimeTrackingResponse, error) {
	log.Printf("Getting time tracking data for match %s", matchID)

	// Get all innings for the match
	innings, err := s.scorecardRepo.GetInningsByMatchID(ctx, matchID)
	if err != nil {
		log.Printf("Error getting innings: %v", err)
		return nil, fmt.Errorf("failed to get innings: %w", err)
	}

	if len(innings) == 0 {
		return &models.TimeTrackingResponse{
			MatchID:        matchID,
			Innings:        []models.InningsTimeTracking{},
			TotalMatchTime: 0,
		}, nil
	}

	var inningsTimeTracking []models.InningsTimeTracking
	var totalMatchTime int

	for _, innings := range innings {
		// Get overs for this innings
		overs, err := s.scorecardRepo.GetOversByInnings(ctx, innings.ID)
		if err != nil {
			log.Printf("Error getting overs for innings %s: %v", innings.ID, err)
			continue
		}

		var overTimeTracking []models.OverTimeTracking
		for _, over := range overs {
			overTimeTracking = append(overTimeTracking, models.OverTimeTracking{
				OverNumber:      over.OverNumber,
				StartTime:       over.StartTime,
				EndTime:         over.EndTime,
				DurationSeconds: over.DurationSeconds,
				Status:          over.Status,
				TotalRuns:       over.TotalRuns,
				TotalBalls:      over.TotalBalls,
				TotalWickets:    over.TotalWickets,
			})
		}

		inningsTimeTracking = append(inningsTimeTracking, models.InningsTimeTracking{
			InningsNumber:   innings.InningsNumber,
			BattingTeam:     innings.BattingTeam,
			StartTime:       innings.StartTime,
			EndTime:         innings.EndTime,
			DurationSeconds: innings.DurationSeconds,
			Status:          innings.Status,
			Overs:           overTimeTracking,
		})

		// Add to total match time
		totalMatchTime += innings.DurationSeconds
	}

	response := &models.TimeTrackingResponse{
		MatchID:        matchID,
		Innings:        inningsTimeTracking,
		TotalMatchTime: totalMatchTime,
	}

	log.Printf("Successfully retrieved time tracking data for match %s", matchID)
	return response, nil
}

// UndoBall removes the last ball from the current over and updates statistics
// Handles edge cases:
// 1. First ball of first over: prevents undo
// 2. First ball of any over: deletes the over and reverts to previous over
// 3. Last ball of innings: properly reverts innings status
func (s *ScorecardService) UndoBall(ctx context.Context, matchID string, inningsNumber int) error {
	log.Printf("Undoing last ball for match %s, innings %d", matchID, inningsNumber)

	// Get user ID from context
	userID, ok := ctx.Value(contextkeys.UserIDKey).(string)
	if !ok || userID == "" {
		return fmt.Errorf("user authentication required")
	}

	// Get match details
	match, err := s.matchRepo.GetByID(ctx, matchID)
	if err != nil {
		log.Printf("Error getting match: %v", err)
		return fmt.Errorf("match not found: %w", err)
	}

	// Check ownership - user must be the creator of the match
	if match.CreatedBy != userID {
		return fmt.Errorf("access denied: you can only undo balls for matches you created")
	}

	// Check if match is live
	if match.Status != models.MatchStatusLive {
		return fmt.Errorf("match is not live, cannot undo ball")
	}

	// Get innings
	innings, err := s.scorecardRepo.GetInningsByMatchAndNumber(ctx, matchID, inningsNumber)
	if err != nil {
		log.Printf("Error getting innings: %v", err)
		return fmt.Errorf("innings not found: %w", err)
	}

	// Check if innings is in progress
	if innings.Status != string(models.InningsStatusInProgress) {
		return fmt.Errorf("innings is not in progress, cannot undo ball")
	}

	// Get the last over (by over_number) - this works regardless of over status
	lastOver, err := s.scorecardRepo.GetLastOver(ctx, innings.ID)
	if err != nil {
		log.Printf("Error getting last over: %v", err)
		return fmt.Errorf("no over found to undo ball: %w", err)
	}

	// Get all balls for this over
	balls, err := s.scorecardRepo.GetBallsByOver(ctx, lastOver.ID)
	if err != nil {
		log.Printf("Error getting balls: %v", err)
		return fmt.Errorf("failed to get balls: %w", err)
	}

	// Edge case 1: Check if this is the first ball of the first over
	if lastOver.OverNumber == 1 && len(balls) == 0 {
		return fmt.Errorf("cannot undo: no balls have been bowled yet")
	}

	// Edge case 2: Check if this is the first ball of any over (but not first over)
	if len(balls) == 0 && lastOver.OverNumber > 1 {
		log.Printf("First ball of over %d - deleting the over and reverting to previous over", lastOver.OverNumber)

		// Delete this empty over
		err = s.scorecardRepo.DeleteOver(ctx, lastOver.ID)
		if err != nil {
			log.Printf("Error deleting over: %v", err)
			return fmt.Errorf("failed to delete over: %w", err)
		}

		// Get the previous over and mark it as in progress
		previousOver, err := s.scorecardRepo.GetOverByInningsAndNumber(ctx, innings.ID, lastOver.OverNumber-1)
		if err != nil {
			log.Printf("Error getting previous over: %v", err)
			return fmt.Errorf("failed to get previous over: %w", err)
		}

		// Mark previous over as in progress
		previousOver.Status = string(models.OverStatusInProgress)
		err = s.scorecardRepo.UpdateOver(ctx, previousOver)
		if err != nil {
			log.Printf("Error updating previous over: %v", err)
			return fmt.Errorf("failed to update previous over: %w", err)
		}

		// Recalculate innings total overs
		overs, err := s.scorecardRepo.GetOversByInnings(ctx, innings.ID)
		if err != nil {
			log.Printf("Error getting overs for innings: %v", err)
			return fmt.Errorf("failed to get overs: %w", err)
		}

		completedOvers := 0
		currentOverBalls := 0
		for _, over := range overs {
			if over.Status == string(models.OverStatusCompleted) {
				completedOvers++
			} else if over.Status == string(models.OverStatusInProgress) {
				currentOverBalls = over.TotalBalls
			}
		}

		var currentOverDecimal float64
		if currentOverBalls > 0 {
			if currentOverBalls == 6 {
				currentOverDecimal = 1.0
			} else {
				currentOverDecimal = float64(currentOverBalls) / 10.0
			}
		}
		innings.TotalOvers = float64(completedOvers) + currentOverDecimal

		err = s.scorecardRepo.UpdateInnings(ctx, innings)
		if err != nil {
			log.Printf("Error updating innings: %v", err)
			return fmt.Errorf("failed to update innings: %w", err)
		}

		// Invalidate caches including the deleted over and previous over
		s.invalidateScorecardCacheForMatch(matchID, innings.ID)
		s.invalidateOverCaches(innings.ID, lastOver.ID)
		if previousOver != nil {
			s.invalidateOverCaches(innings.ID, previousOver.ID)
		}

		log.Printf("Successfully deleted over %d and reverted to over %d", lastOver.OverNumber, previousOver.OverNumber)
		return nil
	}

	if len(balls) == 0 {
		return fmt.Errorf("no balls to undo in this over")
	}

	// Find the last ball (highest ball number)
	var lastBall *models.ScorecardBall
	maxBallNumber := 0
	for _, ball := range balls {
		if ball != nil && ball.BallNumber > maxBallNumber {
			maxBallNumber = ball.BallNumber
			lastBall = ball
		}
	}

	if lastBall == nil {
		return fmt.Errorf("no last ball found")
	}

	// Calculate runs to subtract
	runs := lastBall.Runs
	byes := lastBall.Byes
	totalRuns := runs + byes

	// Delete the ball
	err = s.scorecardRepo.DeleteBall(ctx, lastBall.ID)
	if err != nil {
		log.Printf("Error deleting ball: %v", err)
		return fmt.Errorf("failed to delete ball: %w", err)
	}

	// Update over statistics
	lastOver.TotalRuns -= totalRuns
	// Ensure total_runs doesn't go negative (database constraint)
	if lastOver.TotalRuns < 0 {
		lastOver.TotalRuns = 0
	}
	// Only count legal balls (good balls) for over completion
	if lastBall.BallType == models.BallTypeGood {
		lastOver.TotalBalls--
	}
	if lastBall.IsWicket {
		lastOver.TotalWickets--
	}

	// Check if over should be marked as in progress (if it was completed)
	if lastOver.Status == string(models.OverStatusCompleted) && lastOver.TotalBalls < 6 {
		lastOver.Status = string(models.OverStatusInProgress)
		log.Printf("Reverting over %d status from completed to in_progress", lastOver.OverNumber)
	}

	err = s.scorecardRepo.UpdateOver(ctx, lastOver)
	if err != nil {
		log.Printf("Error updating over: %v", err)
		return fmt.Errorf("failed to update over: %w", err)
	}

	// Update innings statistics
	innings.TotalRuns -= totalRuns
	// Ensure total_runs doesn't go negative (database constraint)
	if innings.TotalRuns < 0 {
		innings.TotalRuns = 0
	}
	// Only count legal balls for innings overs calculation
	if lastBall.BallType == models.BallTypeGood {
		innings.TotalBalls--
	}
	if lastBall.IsWicket {
		innings.TotalWickets--
	}

	// Recalculate total overs
	overs, err := s.scorecardRepo.GetOversByInnings(ctx, innings.ID)
	if err != nil {
		log.Printf("Error getting overs for innings: %v", err)
		return fmt.Errorf("failed to get overs: %w", err)
	}

	completedOvers := 0
	currentOverBalls := 0

	for _, over := range overs {
		if over.Status == string(models.OverStatusCompleted) {
			completedOvers++
		} else if over.Status == string(models.OverStatusInProgress) {
			currentOverBalls = over.TotalBalls
		}
	}

	// Total overs = completed overs + current over balls as decimal
	var currentOverDecimal float64
	if currentOverBalls > 0 {
		// Convert balls to cricket scoring format (0.1, 0.2, 0.3, 0.4, 0.5, 1.0)
		if currentOverBalls == 6 {
			currentOverDecimal = 1.0
		} else {
			currentOverDecimal = float64(currentOverBalls) / 10.0
		}
	}
	innings.TotalOvers = float64(completedOvers) + currentOverDecimal

	// Edge case 3: Check if innings should be marked as in progress (if it was completed)
	if innings.Status == string(models.InningsStatusCompleted) {
		maxWickets := match.TeamAPlayerCount - 1
		if innings.TotalWickets < maxWickets && innings.TotalOvers < float64(match.TotalOvers) {
			innings.Status = string(models.InningsStatusInProgress)
			log.Printf("Reverting innings %d status from completed to in_progress", innings.InningsNumber)
		}
	}

	err = s.scorecardRepo.UpdateInnings(ctx, innings)
	if err != nil {
		log.Printf("Error updating innings: %v", err)
		return fmt.Errorf("failed to update innings: %w", err)
	}

	// Handle match progression - if match was completed, revert it
	if match.Status == models.MatchStatusCompleted {
		match.Status = models.MatchStatusLive
		err = s.matchRepo.Update(ctx, matchID, match)
		if err != nil {
			log.Printf("Error reverting match status: %v", err)
			return fmt.Errorf("failed to revert match status: %w", err)
		}
		log.Printf("Reverted match %s status from completed to live", matchID)
	}

	// Invalidate caches including the specific over that was modified
	s.invalidateScorecardCacheForMatch(matchID, innings.ID)
	s.invalidateOverCaches(innings.ID, lastOver.ID)

	log.Printf("Successfully undone ball: %s %d runs, byes: %d, total: %d, wicket: %v", lastBall.RunType, runs, byes, totalRuns, lastBall.IsWicket)
	return nil
}

// getCurrentOverWithOvers gets the current over and all overs for the innings in one call
// This avoids duplicate GetOversByInnings calls
func (s *ScorecardService) getCurrentOverWithOvers(ctx context.Context, inningsID string) (*models.ScorecardOver, []*models.ScorecardOver, error) {
	// Try to get current over first
	over, err := s.scorecardRepo.GetCurrentOver(ctx, inningsID)
	if err == nil && over != nil {
		// If we found a current over, we still need to get all overs for calculations
		allOvers, err := s.scorecardRepo.GetOversByInnings(ctx, inningsID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get overs: %w", err)
		}
		return over, allOvers, nil
	}

	// Get all overs for this innings to determine next over number
	overs, err := s.scorecardRepo.GetOversByInnings(ctx, inningsID)
	if err != nil {
		log.Printf("Error getting overs: %v", err)
		return nil, nil, fmt.Errorf("failed to get overs: %w", err)
	}

	// Calculate next over number with defensive checks
	nextOverNumber := 1
	if len(overs) > 0 {
		// Find the highest over number and add 1
		maxOverNumber := 0
		for _, o := range overs {
			if o != nil && o.OverNumber > maxOverNumber {
				maxOverNumber = o.OverNumber
			}
		}
		nextOverNumber = maxOverNumber + 1
	}

	// Create new over with retry logic for constraint violations
	now := time.Now()
	newOver := &models.ScorecardOver{
		InningsID:    inningsID,
		OverNumber:   nextOverNumber,
		TotalRuns:    0,
		TotalBalls:   0,
		TotalWickets: 0,
		Status:       string(models.OverStatusInProgress),
		StartTime:    &now,
	}

	// Use improved retry logic with exponential backoff and jitter
	err = s.retryWithExponentialBackoff(ctx, func() error {
		err := s.scorecardRepo.CreateOver(ctx, newOver)
		if err == nil {
			return nil
		}

		// If it's a constraint violation, try to get fresh data and retry
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") ||
			strings.Contains(err.Error(), "violates foreign key constraint") {
			log.Printf("Over creation failed due to constraint violation, retrying with fresh data: %v", err)

			// Get fresh overs to recalculate over number
			freshOvers, freshErr := s.scorecardRepo.GetOversByInnings(ctx, inningsID)
			if freshErr != nil {
				log.Printf("Failed to get fresh overs: %v", freshErr)
				return fmt.Errorf("failed to get fresh overs: %w", freshErr)
			}

			// Recalculate next over number with fresh data and defensive checks
			nextOverNumber = 1
			if len(freshOvers) > 0 {
				maxOverNumber := 0
				for _, o := range freshOvers {
					if o != nil && o.OverNumber > maxOverNumber {
						maxOverNumber = o.OverNumber
					}
				}
				nextOverNumber = maxOverNumber + 1
			}
			newOver.OverNumber = nextOverNumber

			// Return the original error to trigger retry
			return err
		}

		// If it's not a constraint violation, return the error immediately
		return err
	})

	if err != nil {
		log.Printf("Error creating over: %v", err)
		return nil, nil, fmt.Errorf("failed to create over: %w", err)
	}

	log.Printf("Created new over %d for innings %s", nextOverNumber, inningsID)
	return newOver, overs, nil
}

// getNextBallNumber gets the next ball number for an over (optimized but preserves cricket logic)
func (s *ScorecardService) getNextBallNumber(ctx context.Context, overID string) (int, error) {
	// Get only necessary ball fields for cricket logic (optimized query)
	balls, err := s.scorecardRepo.GetBallsForNextNumber(ctx, overID)
	if err != nil {
		return 0, fmt.Errorf("failed to get balls: %w", err)
	}

	// Count legal balls (good balls only) and find max ball number
	legalBalls := 0
	maxBallNumber := 0

	for _, ball := range balls {
		if ball != nil && ball.BallNumber > maxBallNumber {
			maxBallNumber = ball.BallNumber
		}
		// Only count good balls as legal deliveries for over completion
		if ball != nil && ball.BallType == models.BallTypeGood {
			legalBalls++
		}
	}

	// An over is complete when it has 6 legal balls
	if legalBalls >= 6 {
		return 0, fmt.Errorf("over is complete, cannot add more balls")
	}

	// The next ball number is simply the next sequential number
	nextBallNumber := maxBallNumber + 1

	return nextBallNumber, nil
}

// invalidateMatchCachesBeforeWrite invalidates caches BEFORE write operations to prevent stale reads
func (s *ScorecardService) invalidateMatchCachesBeforeWrite(ctx context.Context, matchID, inningsID, overID string) {
	// Skip cache invalidation if cache is not available
	if s.cache == nil {
		return
	}

	log.Printf("Invalidating caches BEFORE write for match %s, innings %s, over %s", matchID, inningsID, overID)

	// Invalidate all related caches to prevent stale data reads during concurrent operations
	cacheKeys := []string{
		fmt.Sprintf("scorecard:%s", matchID),
		fmt.Sprintf("innings:match:%s", matchID),
		fmt.Sprintf("overs:innings:%s", inningsID),
		fmt.Sprintf("over:current:innings:%s", inningsID),
		fmt.Sprintf("match_innings_over:%s:%d", matchID, 1), // First innings
		fmt.Sprintf("match_innings_over:%s:%d", matchID, 2), // Second innings
		fmt.Sprintf("balls:over:%s", overID),
		fmt.Sprintf("last_ball:over:%s", overID),
	}

	for _, key := range cacheKeys {
		_ = s.cache.Invalidate(key)
	}

	log.Printf("Pre-write cache invalidation completed for match %s", matchID)
}

// invalidateMatchCaches invalidates all caches related to a match after ball addition
func (s *ScorecardService) invalidateMatchCaches(ctx context.Context, matchID, inningsID, overID string) {
	// Skip cache invalidation if cache is not available
	if s.cache == nil {
		return
	}

	// Invalidate scorecard cache for the match
	scorecardKey := fmt.Sprintf("scorecard:%s", matchID)
	_ = s.cache.Invalidate(scorecardKey)

	// Invalidate innings cache
	inningsKey := fmt.Sprintf("innings:match:%s", matchID)
	_ = s.cache.Invalidate(inningsKey)

	// Invalidate overs cache for this innings
	oversKey := fmt.Sprintf("overs:innings:%s", inningsID)
	_ = s.cache.Invalidate(oversKey)

	// Invalidate current over cache
	currentOverKey := fmt.Sprintf("over:current:innings:%s", inningsID)
	_ = s.cache.Invalidate(currentOverKey)

	// Invalidate balls cache for this over
	ballsKey := fmt.Sprintf("balls:over:%s", overID)
	_ = s.cache.Invalidate(ballsKey)

	// Invalidate ball count cache
	ballCountKey := fmt.Sprintf("ball_count:over:%s", overID)
	_ = s.cache.Invalidate(ballCountKey)

	// Invalidate balls for next number cache
	ballsNextNumberKey := fmt.Sprintf("balls_next_number:over:%s", overID)
	_ = s.cache.Invalidate(ballsNextNumberKey)

	// Invalidate last ball cache
	lastBallKey := fmt.Sprintf("ball:last:over:%s", overID)
	_ = s.cache.Invalidate(lastBallKey)

	log.Printf("Invalidated all caches for match %s, innings %s, over %s", matchID, inningsID, overID)
}

// invalidateMatchCachesAsync invalidates caches asynchronously
func (s *ScorecardService) invalidateMatchCachesAsync(matchID, inningsID, overID string) {
	ctx := context.Background()
	s.invalidateMatchCaches(ctx, matchID, inningsID, overID)
}

// invalidateScorecardCacheForMatch invalidates scorecard-related caches for a match
func (s *ScorecardService) invalidateScorecardCacheForMatch(matchID, inningsID string) {
	// Skip cache invalidation if cache is not available
	if s.cache == nil {
		return
	}

	log.Printf("Invalidating scorecard caches for match %s", matchID)

	// Invalidate scorecard cache
	scorecardKey := fmt.Sprintf("scorecard:%s", matchID)
	_ = s.cache.Invalidate(scorecardKey)

	// Invalidate innings cache
	inningsKey := fmt.Sprintf("innings:match:%s", matchID)
	_ = s.cache.Invalidate(inningsKey)

	// Invalidate overs cache for this innings
	oversKey := fmt.Sprintf("overs:innings:%s", inningsID)
	_ = s.cache.Invalidate(oversKey)

	// Invalidate current over cache
	currentOverKey := fmt.Sprintf("over:current:innings:%s", inningsID)
	_ = s.cache.Invalidate(currentOverKey)

	// Invalidate last over cache
	lastOverKey := fmt.Sprintf("over:last:innings:%s", inningsID)
	_ = s.cache.Invalidate(lastOverKey)

	// Invalidate match innings over data cache
	_ = s.cache.Invalidate(fmt.Sprintf("match_innings_over:%s:%d", matchID, 1))
	_ = s.cache.Invalidate(fmt.Sprintf("match_innings_over:%s:%d", matchID, 2))

	log.Printf("Scorecard cache invalidation completed for match %s", matchID)
}

// invalidateOverCaches invalidates all caches related to a specific over
func (s *ScorecardService) invalidateOverCaches(inningsID, overID string) {
	// Skip cache invalidation if cache is not available
	if s.cache == nil {
		return
	}

	log.Printf("Invalidating over-specific caches for over %s", overID)

	// Invalidate balls cache for this over
	ballsKey := fmt.Sprintf("balls:over:%s", overID)
	_ = s.cache.Invalidate(ballsKey)

	// Invalidate ball count cache
	ballCountKey := fmt.Sprintf("ball_count:over:%s", overID)
	_ = s.cache.Invalidate(ballCountKey)

	// Invalidate balls for next number cache
	ballsNextNumberKey := fmt.Sprintf("balls_next_number:over:%s", overID)
	_ = s.cache.Invalidate(ballsNextNumberKey)

	// Invalidate last ball cache
	lastBallKey := fmt.Sprintf("ball:last:over:%s", overID)
	_ = s.cache.Invalidate(lastBallKey)

	// Invalidate last over cache for innings (since over details changed)
	lastOverKey := fmt.Sprintf("over:last:innings:%s", inningsID)
	_ = s.cache.Invalidate(lastOverKey)

	log.Printf("Over-specific cache invalidation completed for over %s", overID)
}

// calculateOversInMemory performs optimized overs calculation in memory
func (s *ScorecardService) calculateOversInMemory(completedOvers int, currentOverBalls int) float64 {
	// Total overs = completed overs + current over balls as decimal
	var currentOverDecimal float64
	if currentOverBalls > 0 {
		// Convert balls to cricket scoring format (0.1, 0.2, 0.3, 0.4, 0.5, 1.0)
		if currentOverBalls == 6 {
			currentOverDecimal = 1.0
		} else {
			currentOverDecimal = float64(currentOverBalls) / 10.0
		}
	}
	return float64(completedOvers) + currentOverDecimal
}

// calculateInningsStatusInMemory performs optimized innings completion check in memory
func (s *ScorecardService) calculateInningsStatusInMemory(innings *models.Innings, maxWickets int, totalOvers int) string {
	// Check if innings is complete
	if innings.TotalWickets >= maxWickets || innings.TotalOvers >= float64(totalOvers) {
		return string(models.InningsStatusCompleted)
	}
	return string(models.InningsStatusInProgress)
}

// ShouldCompleteMatch determines if the match should be completed based on cricket rules
func (s *ScorecardService) ShouldCompleteMatch(ctx context.Context, matchID string, secondInnings *models.Innings, match *models.Match) (bool, string) {
	// Get first innings score
	firstInnings, err := s.scorecardRepo.GetInningsByMatchAndNumber(ctx, matchID, 1)
	if err != nil {
		return false, "error getting first innings"
	}

	target := firstInnings.TotalRuns + 1     // Target is first innings score + 1
	maxWickets := match.TeamAPlayerCount - 1 // n-1 wickets for n players

	// Check if target is reached
	if secondInnings.TotalRuns >= target {
		return true, fmt.Sprintf("target reached: %d/%d", secondInnings.TotalRuns, target)
	}

	// Check if all wickets are lost
	if secondInnings.TotalWickets >= maxWickets {
		return true, fmt.Sprintf("all wickets lost: %d/%d", secondInnings.TotalWickets, maxWickets)
	}

	// Check if all overs are completed
	if secondInnings.TotalOvers >= float64(match.TotalOvers) {
		return true, fmt.Sprintf("all overs completed: %.1f/%d", secondInnings.TotalOvers, match.TotalOvers)
	}

	return false, "match continues"
}

// startSecondInnings starts the second innings
func (s *ScorecardService) startSecondInnings(ctx context.Context, matchID string, match *models.Match) error {
	log.Printf("Starting second innings for match %s", matchID)

	// Try to fetch complete match data, but use fallback if database call fails
	var completeMatch *models.Match
	var err error

	completeMatch, err = s.matchRepo.GetByID(ctx, matchID)
	if err != nil {
		log.Printf("Warning: Failed to fetch complete match data from database: %v", err)
		log.Printf("Using fallback match data for second innings transition")

		// Use the provided match data as fallback
		completeMatch = match

		// Ensure we have the minimum required fields
		if completeMatch.TossWinner == "" {
			log.Printf("Error: TossWinner is missing from match data, cannot determine second innings batting team")
			return fmt.Errorf("toss winner information is missing from match data")
		}
	}

	// Determine batting team for second innings (non-toss-winning team)
	var battingTeam models.TeamType
	if completeMatch.TossWinner == models.TeamTypeA {
		battingTeam = models.TeamTypeB // Second innings should be played by non-toss winner
	} else {
		battingTeam = models.TeamTypeA // Second innings should be played by non-toss winner
	}

	// Create second innings
	now := time.Now()
	secondInnings := &models.Innings{
		MatchID:       matchID,
		InningsNumber: 2,
		BattingTeam:   battingTeam,
		TotalRuns:     0,
		TotalWickets:  0,
		TotalOvers:    0.0,
		TotalBalls:    0,
		Status:        string(models.InningsStatusInProgress),
		StartTime:     &now,
	}

	err = s.scorecardRepo.CreateInnings(ctx, secondInnings)
	if err != nil {
		log.Printf("Error creating second innings: %v", err)
		return fmt.Errorf("failed to start second innings: %w", err)
	}

	// Update match batting team
	completeMatch.BattingTeam = battingTeam
	err = s.matchRepo.Update(ctx, matchID, completeMatch)
	if err != nil {
		log.Printf("Warning: Failed to update match batting team in database: %v", err)
		log.Printf("Second innings created successfully, but match update failed - this may cause inconsistencies")
		// Don't fail the entire operation, just log the warning
		// The second innings is already created, which is the critical part
	} else {
		log.Printf("Successfully updated match batting team to %s", battingTeam)
	}

	log.Printf("Successfully started second innings for match %s, batting team: %s", matchID, battingTeam)
	return nil
}

// GetScorecard gets the complete scorecard for a match
func (s *ScorecardService) GetScorecard(ctx context.Context, matchID string) (*models.ScorecardResponse, error) {
	log.Printf("Getting scorecard for match %s", matchID)

	// Get scorecard from repository
	scorecard, err := s.scorecardRepo.GetScorecard(ctx, matchID)
	if err != nil {
		log.Printf("Error getting scorecard: %v", err)
		return nil, fmt.Errorf("failed to get scorecard: %w", err)
	}

	log.Printf("Successfully retrieved scorecard for match %s", matchID)
	return scorecard, nil
}

// GetCurrentOver gets the current over for a match
func (s *ScorecardService) GetCurrentOver(ctx context.Context, matchID string, inningsNumber int) (*models.ScorecardOver, error) {
	log.Printf("Getting current over for match %s, innings %d", matchID, inningsNumber)

	// Get innings
	innings, err := s.scorecardRepo.GetInningsByMatchAndNumber(ctx, matchID, inningsNumber)
	if err != nil {
		log.Printf("Error getting innings: %v", err)
		return nil, fmt.Errorf("innings not found: %w", err)
	}

	// Get current over
	over, err := s.scorecardRepo.GetCurrentOver(ctx, innings.ID)
	if err != nil {
		log.Printf("Error getting current over: %v", err)
		return nil, fmt.Errorf("no current over found: %w", err)
	}

	log.Printf("Found current over %d for match %s, innings %d", over.OverNumber, matchID, inningsNumber)
	return over, nil
}

// GetBallsByOver gets all balls for a specific over
func (s *ScorecardService) GetBallsByOver(ctx context.Context, overID string) ([]*models.ScorecardBall, error) {
	log.Printf("Getting balls for over %s", overID)

	balls, err := s.scorecardRepo.GetBallsByOver(ctx, overID)
	if err != nil {
		log.Printf("Error getting balls for over: %v", err)
		return nil, fmt.Errorf("failed to get balls for over: %w", err)
	}

	log.Printf("Found %d balls for over %s", len(balls), overID)
	return balls, nil
}

// GetInningsByMatchAndNumber gets innings by match ID and innings number
func (s *ScorecardService) GetInningsByMatchAndNumber(ctx context.Context, matchID string, inningsNumber int) (*models.Innings, error) {
	log.Printf("Getting innings for match %s, innings number %d", matchID, inningsNumber)

	innings, err := s.scorecardRepo.GetInningsByMatchAndNumber(ctx, matchID, inningsNumber)
	if err != nil {
		log.Printf("Error getting innings: %v", err)
		return nil, fmt.Errorf("failed to get innings: %w", err)
	}

	log.Printf("Found innings %d for match %s", inningsNumber, matchID)
	return innings, nil
}

// GetOversByInnings gets all overs for a specific innings
func (s *ScorecardService) GetOversByInnings(ctx context.Context, inningsID string) ([]*models.ScorecardOver, error) {
	log.Printf("Getting overs for innings %s", inningsID)

	overs, err := s.scorecardRepo.GetOversByInnings(ctx, inningsID)
	if err != nil {
		log.Printf("Error getting overs for innings: %v", err)
		return nil, fmt.Errorf("failed to get overs for innings: %w", err)
	}

	log.Printf("Found %d overs for innings %s", len(overs), inningsID)
	return overs, nil
}

// ValidateInningsOrder validates that balls can only be added to the correct innings
func (s *ScorecardService) ValidateInningsOrder(ctx context.Context, matchID string, match *models.Match, inningsNumber int) error {
	log.Printf("DEBUG: validateInningsOrder called - matchID: %s, inningsNumber: %d, battingTeam: %s, tossWinner: %s",
		matchID, inningsNumber, match.BattingTeam, match.TossWinner)

	// Get all innings for this match to determine current state
	innings, err := s.scorecardRepo.GetInningsByMatchID(ctx, matchID)
	if err != nil {
		log.Printf("DEBUG: No existing innings found for match %s, error: %v", matchID, err)
		// If no innings exist, this is the first ball of the match
		// The first innings should always be the toss-winning team
		if inningsNumber == 1 {
			// Check if the batting team matches the toss winner
			if match.BattingTeam != match.TossWinner {
				log.Printf("DEBUG: Validation failed - first innings must be played by toss winner %s, but batting team is %s",
					match.TossWinner, match.BattingTeam)
				return fmt.Errorf("first innings must be played by the toss-winning team (%s), but current batting team is %s",
					match.TossWinner, match.BattingTeam)
			}
			log.Printf("DEBUG: Validation passed - first innings with correct team")
			return nil
		} else {
			log.Printf("DEBUG: Validation failed - cannot start with innings %d, first innings must be played first", inningsNumber)
			return fmt.Errorf("cannot start with innings %d, first innings must be played first", inningsNumber)
		}
	}

	log.Printf("DEBUG: Found %d existing innings for match %s", len(innings), matchID)

	// Determine which innings exist
	firstInningsExists := false

	for _, inn := range innings {
		if inn.InningsNumber == 1 {
			firstInningsExists = true
		}
	}

	// Check if we're trying to add to first innings
	if inningsNumber == 1 {
		if !firstInningsExists {
			// First innings doesn't exist yet, check if batting team is correct
			if match.BattingTeam != match.TossWinner {
				return fmt.Errorf("first innings must be played by the toss-winning team (%s), but current batting team is %s",
					match.TossWinner, match.BattingTeam)
			}
			return nil
		} else {
			// First innings exists, check if it's complete
			firstInnings, err := s.scorecardRepo.GetInningsByMatchAndNumber(ctx, matchID, 1)
			if err != nil {
				return fmt.Errorf("failed to get first innings: %w", err)
			}

			// Check if first innings is complete (all wickets down or overs completed)
			maxWickets := match.TeamAPlayerCount - 1 // n-1 wickets for n players
			firstInningsComplete := firstInnings.TotalWickets >= maxWickets || firstInnings.TotalOvers >= float64(match.TotalOvers)

			if !firstInningsComplete {
				// First innings is not complete, only toss winner can bat
				if match.BattingTeam != match.TossWinner {
					return fmt.Errorf("first innings is not complete, only toss-winning team (%s) can bat, but current batting team is %s",
						match.TossWinner, match.BattingTeam)
				}
			} else {
				// First innings is complete, but we're trying to add to first innings again
				return fmt.Errorf("first innings is complete, cannot add more balls to first innings")
			}
			return nil
		}
	}

	// Check if we're trying to add to second innings
	if inningsNumber == 2 {
		if !firstInningsExists {
			return fmt.Errorf("cannot start second innings, first innings must be played first")
		}

		// Check if first innings is complete
		firstInnings, err := s.scorecardRepo.GetInningsByMatchAndNumber(ctx, matchID, 1)
		if err != nil {
			return fmt.Errorf("failed to get first innings: %w", err)
		}

		// First innings is complete if all wickets are down or overs are completed
		maxWickets := match.TeamAPlayerCount - 1 // n-1 wickets for n players
		firstInningsComplete := firstInnings.TotalWickets >= maxWickets || firstInnings.TotalOvers >= float64(match.TotalOvers)

		if !firstInningsComplete {
			return fmt.Errorf("first innings is not complete, cannot start second innings")
		}

		// Second innings should be played by the non-toss-winning team
		nonTossWinner := s.GetNonTossWinner(match.TossWinner)
		if match.BattingTeam != nonTossWinner {
			return fmt.Errorf("second innings must be played by the non-toss-winning team (%s), but current batting team is %s",
				nonTossWinner, match.BattingTeam)
		}

		return nil
	}

	return fmt.Errorf("invalid innings number: %d", inningsNumber)
}

// GetNonTossWinner returns the team that didn't win the toss
func (s *ScorecardService) GetNonTossWinner(tossWinner models.TeamType) models.TeamType {
	if tossWinner == models.TeamTypeA {
		return models.TeamTypeB
	}
	return models.TeamTypeA
}
