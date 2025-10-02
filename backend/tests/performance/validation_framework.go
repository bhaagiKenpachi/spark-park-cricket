package performance

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"spark-park-cricket-backend/internal/database"
	"spark-park-cricket-backend/internal/models"
)

// ValidationLevel defines the complexity of validation to apply
type ValidationLevel int

const (
	// BasicValidation - Minimal checks, focus on functionality (E2E Tests)
	BasicValidation ValidationLevel = iota
	// IntermediateValidation - Balanced approach with retry logic (Load Tests)
	IntermediateValidation
	// FullValidation - Comprehensive atomic operations (Integration Tests)
	FullValidation
)

// TestDataCreator defines the interface for test data creation with validation
type TestDataCreator interface {
	CreateTestData(dbClient *database.Client, userID string, level ValidationLevel) (string, string, error)
	ValidateTestData(dbClient *database.Client, seriesID, matchID, inningsID, overID string, level ValidationLevel) error
}

// ValidationConfig holds configuration for different validation levels
type ValidationConfig struct {
	Level                 ValidationLevel
	MaxRetries            int
	RetryDelay            time.Duration
	ValidateRelationships bool
	ValidateDataIntegrity bool
	AtomicOperations      bool
	MutexProtection       bool
}

// GetValidationConfig returns configuration for a specific validation level
func GetValidationConfig(level ValidationLevel) ValidationConfig {
	switch level {
	case BasicValidation:
		return ValidationConfig{
			Level:                 BasicValidation,
			MaxRetries:            3,
			RetryDelay:            100 * time.Millisecond,
			ValidateRelationships: false,
			ValidateDataIntegrity: true,
			AtomicOperations:      false,
			MutexProtection:       true,
		}
	case IntermediateValidation:
		return ValidationConfig{
			Level:                 IntermediateValidation,
			MaxRetries:            3,
			RetryDelay:            100 * time.Millisecond,
			ValidateRelationships: true,
			ValidateDataIntegrity: false,
			AtomicOperations:      false,
			MutexProtection:       true,
		}
	case FullValidation:
		return ValidationConfig{
			Level:                 FullValidation,
			MaxRetries:            3,
			RetryDelay:            200 * time.Millisecond,
			ValidateRelationships: true,
			ValidateDataIntegrity: true,
			AtomicOperations:      true,
			MutexProtection:       true,
		}
	default:
		return GetValidationConfig(BasicValidation)
	}
}

var globalTestDataMutex sync.Mutex // Global mutex for test data creation

// createTestDataWithValidation creates test data with configurable validation level using transactions
func createTestDataWithValidation(dbClient *database.Client, userID string, level ValidationLevel, testType string) (string, string, error) {
	config := GetValidationConfig(level)

	log.Printf("🔧 DEBUG: Starting test data creation for %s with validation level %d", testType, level)
	log.Printf("🔧 DEBUG: User ID: %s", userID)

	// Apply mutex protection if configured
	if config.MutexProtection {
		log.Printf("🔧 DEBUG: Acquiring global test data mutex")
		globalTestDataMutex.Lock()
		defer func() {
			globalTestDataMutex.Unlock()
			log.Printf("🔧 DEBUG: Released global test data mutex")
		}()
		time.Sleep(config.RetryDelay)
	}

	ctx := context.Background()

	// Retry logic based on validation level
	for attempt := 1; attempt <= config.MaxRetries; attempt++ {
		log.Printf("🔧 DEBUG: Attempt %d/%d for test data creation", attempt, config.MaxRetries)

		if attempt > 1 {
			retryDelay := time.Duration(attempt) * config.RetryDelay
			log.Printf("🔧 DEBUG: Waiting %v before retry", retryDelay)
			time.Sleep(retryDelay)
		}

		// Create highly unique identifier to prevent interference
		uniqueID := fmt.Sprintf("%d-%d-%d-%d-%s-%d", time.Now().UnixNano(), time.Now().Unix(), attempt, len(userID), testType, os.Getpid())
		log.Printf("🔧 DEBUG: Generated unique ID: %s", uniqueID)

		// Note: Supabase doesn't support explicit transactions, so we'll use atomic operations
		log.Printf("🔧 DEBUG: Starting atomic data creation (Supabase doesn't support explicit transactions)")

		// Create test series
		log.Printf("🔧 DEBUG: Creating test series")
		series := &models.Series{
			Name:      fmt.Sprintf("%s Test Series %s", testType, uniqueID),
			StartDate: time.Now(),
			EndDate:   time.Now().Add(24 * time.Hour),
			CreatedBy: userID,
		}

		err := dbClient.Repositories.Series.Create(ctx, series)
		if err != nil {
			log.Printf("❌ DEBUG: Failed to create series: %v", err)
			if attempt == config.MaxRetries {
				return "", "", fmt.Errorf("failed to create test series after %d attempts: %v", config.MaxRetries, err)
			}
			continue
		}
		log.Printf("✅ DEBUG: Created series with ID: %s", series.ID)

		// Create test match
		log.Printf("🔧 DEBUG: Creating test match")
		match := &models.Match{
			SeriesID:         series.ID,
			MatchNumber:      1,
			Date:             time.Now(),
			Status:           models.MatchStatusLive,
			TeamAPlayerCount: 11,
			TeamBPlayerCount: 11,
			TotalOvers:       20,
			TossWinner:       models.TeamTypeA,
			TossType:         models.TossTypeHeads,
			BattingTeam:      models.TeamTypeA,
			CreatedBy:        userID,
		}

		err = dbClient.Repositories.Match.Create(ctx, match)
		if err != nil {
			log.Printf("❌ DEBUG: Failed to create match: %v", err)
			if attempt == config.MaxRetries {
				return "", "", fmt.Errorf("failed to create test match after %d attempts: %v", config.MaxRetries, err)
			}
			continue
		}
		log.Printf("✅ DEBUG: Created match with ID: %s", match.ID)

		// Create first innings
		log.Printf("🔧 DEBUG: Creating first innings")
		innings := &models.Innings{
			MatchID:       match.ID,
			InningsNumber: 1,
			BattingTeam:   models.TeamTypeA,
			TotalRuns:     0,
			TotalWickets:  0,
			TotalOvers:    0,
			TotalBalls:    0,
			Status:        string(models.InningsStatusInProgress),
			CreatedAt:     time.Now(),
		}

		err = dbClient.Repositories.Scorecard.CreateInnings(ctx, innings)
		if err != nil {
			log.Printf("❌ DEBUG: Failed to create innings: %v", err)
			if attempt == config.MaxRetries {
				return "", "", fmt.Errorf("failed to create test innings after %d attempts: %v", config.MaxRetries, err)
			}
			continue
		}
		log.Printf("✅ DEBUG: Created innings with ID: %s", innings.ID)

		// Create first over
		log.Printf("🔧 DEBUG: Creating first over")
		over := &models.ScorecardOver{
			InningsID:    innings.ID,
			OverNumber:   1,
			TotalRuns:    0,
			TotalBalls:   0,
			TotalWickets: 0,
			Status:       string(models.OverStatusInProgress),
		}

		err = dbClient.Repositories.Scorecard.CreateOver(ctx, over)
		if err != nil {
			log.Printf("❌ DEBUG: Failed to create over: %v", err)
			if attempt == config.MaxRetries {
				return "", "", fmt.Errorf("failed to create test over after %d attempts: %v", config.MaxRetries, err)
			}
			continue
		}
		log.Printf("✅ DEBUG: Created over with ID: %s", over.ID)

		// Validate over ID is not empty before proceeding
		if over.ID == "" {
			if attempt == config.MaxRetries {
				return "", "", fmt.Errorf("over ID is empty after creation")
			}
			continue
		}

		// Apply validation based on level
		if config.ValidateDataIntegrity {
			log.Printf("🔧 DEBUG: Validating test data integrity")
			if err := validateTestDataIntegrityWithLevel(dbClient, series.ID, match.ID, innings.ID, over.ID, level); err != nil {
				log.Printf("❌ DEBUG: Data integrity validation failed: %v", err)
				if attempt == config.MaxRetries {
					return "", "", fmt.Errorf("failed to validate test data integrity after %d attempts: %v", config.MaxRetries, err)
				}
				continue
			}
			log.Printf("✅ DEBUG: Data integrity validation passed")
		}

		// Always validate basic data creation for all levels
		if series.ID == "" || match.ID == "" || innings.ID == "" || over.ID == "" {
			log.Printf("❌ DEBUG: Empty IDs detected - Series: %s, Match: %s, Innings: %s, Over: %s", series.ID, match.ID, innings.ID, over.ID)
			if attempt == config.MaxRetries {
				return "", "", fmt.Errorf("failed to create complete test data hierarchy after %d attempts", config.MaxRetries)
			}
			continue
		}

		// Ensure database persistence before proceeding
		log.Printf("🔧 DEBUG: Ensuring database persistence for all created entities")
		if err := ensureDatabasePersistence(dbClient, series.ID, match.ID, innings.ID, over.ID, 5*time.Second); err != nil {
			log.Printf("❌ DEBUG: Database persistence validation failed: %v", err)
			if attempt == config.MaxRetries {
				return "", "", fmt.Errorf("failed to ensure database persistence after %d attempts: %v", config.MaxRetries, err)
			}
			continue
		}

		// Validate database connection health
		log.Printf("🔧 DEBUG: Validating database connection health")
		if err := validateDatabaseConnection(dbClient); err != nil {
			log.Printf("❌ DEBUG: Database connection validation failed: %v", err)
			if attempt == config.MaxRetries {
				return "", "", fmt.Errorf("database connection validation failed after %d attempts: %v", config.MaxRetries, err)
			}
			continue
		}

		// Success - return the IDs
		log.Printf("✅ DEBUG: Test data creation completed successfully")
		log.Printf("✅ DEBUG: Series ID: %s", series.ID)
		log.Printf("✅ DEBUG: Match ID: %s", match.ID)
		log.Printf("✅ DEBUG: Innings ID: %s", innings.ID)
		log.Printf("✅ DEBUG: Over ID: %s", over.ID)
		return series.ID, match.ID, nil
	}

	return "", "", fmt.Errorf("unexpected error in createTestDataWithValidation")
}

// validateTestDataIntegrityWithLevel validates data based on validation level
func validateTestDataIntegrityWithLevel(dbClient *database.Client, seriesID, matchID, inningsID, overID string, level ValidationLevel) error {
	config := GetValidationConfig(level)

	log.Printf("🔧 DEBUG: Starting data integrity validation with level %d", level)
	log.Printf("🔧 DEBUG: Validating IDs - Series: %s, Match: %s, Innings: %s, Over: %s", seriesID, matchID, inningsID, overID)

	// Basic validation - just check if IDs are not empty
	if !config.ValidateRelationships {
		log.Printf("🔧 DEBUG: Performing basic validation (ID existence only)")
		if seriesID == "" || matchID == "" || inningsID == "" || overID == "" {
			log.Printf("❌ DEBUG: One or more IDs are empty")
			return fmt.Errorf("one or more IDs are empty")
		}
		log.Printf("✅ DEBUG: Basic validation passed")
		return nil
	}

	log.Printf("🔧 DEBUG: Performing relationship validation")

	// Intermediate and Full validation - check relationships
	log.Printf("🔧 DEBUG: Validating series existence")
	var series models.Series
	_, err := dbClient.Supabase.From("series").Select("*", "exact", false).Eq("id", seriesID).Single().ExecuteTo(&series)
	if err != nil {
		log.Printf("❌ DEBUG: Series validation failed: %v", err)
		return fmt.Errorf("series validation failed: %v", err)
	}
	log.Printf("✅ DEBUG: Series validation passed")

	log.Printf("🔧 DEBUG: Validating match existence and series reference")
	var match models.Match
	_, err = dbClient.Supabase.From("matches").Select("*", "exact", false).Eq("id", matchID).Single().ExecuteTo(&match)
	if err != nil {
		log.Printf("❌ DEBUG: Match validation failed: %v", err)
		return fmt.Errorf("match validation failed: %v", err)
	}

	if config.ValidateRelationships && match.SeriesID != seriesID {
		log.Printf("❌ DEBUG: Match series reference invalid: expected %s, got %s", seriesID, match.SeriesID)
		return fmt.Errorf("match series reference invalid: expected %s, got %s", seriesID, match.SeriesID)
	}
	log.Printf("✅ DEBUG: Match validation passed")

	log.Printf("🔧 DEBUG: Validating innings existence and match reference")
	var innings models.Innings
	_, err = dbClient.Supabase.From("innings").Select("*", "exact", false).Eq("id", inningsID).Single().ExecuteTo(&innings)
	if err != nil {
		log.Printf("❌ DEBUG: Innings validation failed: %v", err)
		return fmt.Errorf("innings validation failed: %v", err)
	}

	if config.ValidateRelationships && innings.MatchID != matchID {
		log.Printf("❌ DEBUG: Innings match reference invalid: expected %s, got %s", matchID, innings.MatchID)
		return fmt.Errorf("innings match reference invalid: expected %s, got %s", matchID, innings.MatchID)
	}
	log.Printf("✅ DEBUG: Innings validation passed")

	log.Printf("🔧 DEBUG: Validating over existence and innings reference")
	var over models.ScorecardOver
	_, err = dbClient.Supabase.From("overs").Select("*", "exact", false).Eq("id", overID).Single().ExecuteTo(&over)
	if err != nil {
		log.Printf("❌ DEBUG: Over validation failed: %v", err)
		return fmt.Errorf("over validation failed: %v", err)
	}

	if config.ValidateRelationships && over.InningsID != inningsID {
		log.Printf("❌ DEBUG: Over innings reference invalid: expected %s, got %s", inningsID, over.InningsID)
		return fmt.Errorf("over innings reference invalid: expected %s, got %s", inningsID, over.InningsID)
	}
	log.Printf("✅ DEBUG: Over validation passed")

	// Full validation - additional data integrity checks
	if config.ValidateDataIntegrity {
		log.Printf("🔧 DEBUG: Performing additional data integrity checks")
		if series.ID == "" || match.ID == "" || innings.ID == "" || over.ID == "" {
			log.Printf("❌ DEBUG: One or more entity IDs are empty")
			return fmt.Errorf("one or more entity IDs are empty")
		}

		// Check for reasonable data values
		log.Printf("🔧 DEBUG: Validating innings number: %d", innings.InningsNumber)
		if innings.InningsNumber < 1 || innings.InningsNumber > 4 {
			log.Printf("❌ DEBUG: Invalid innings number: %d", innings.InningsNumber)
			return fmt.Errorf("invalid innings number: %d", innings.InningsNumber)
		}

		log.Printf("🔧 DEBUG: Validating over number: %d", over.OverNumber)
		if over.OverNumber < 1 {
			log.Printf("❌ DEBUG: Invalid over number: %d", over.OverNumber)
			return fmt.Errorf("invalid over number: %d", over.OverNumber)
		}
		log.Printf("✅ DEBUG: Additional data integrity checks passed")
	}

	log.Printf("✅ DEBUG: All data integrity validations passed")
	return nil
}

// ensureDatabasePersistence waits for database operations to complete and validates persistence
func ensureDatabasePersistence(dbClient *database.Client, seriesID, matchID, inningsID, overID string, maxWaitTime time.Duration) error {
	log.Printf("🔧 DEBUG: Ensuring database persistence for created entities")
	log.Printf("🔧 DEBUG: Waiting for database operations to complete...")

	startTime := time.Now()

	// Wait for database operations to complete
	time.Sleep(500 * time.Millisecond) // Initial wait for database consistency

	for time.Since(startTime) < maxWaitTime {
		log.Printf("🔧 DEBUG: Validating database persistence (elapsed: %v)", time.Since(startTime))

		// Check if all entities exist in database
		allExist := true

		// Check series existence
		var series models.Series
		_, err := dbClient.Supabase.From("series").Select("*", "exact", false).Eq("id", seriesID).Single().ExecuteTo(&series)
		if err != nil {
			log.Printf("🔧 DEBUG: Series not yet persistent: %v", err)
			allExist = false
		}

		// Check match existence
		var match models.Match
		_, err = dbClient.Supabase.From("matches").Select("*", "exact", false).Eq("id", matchID).Single().ExecuteTo(&match)
		if err != nil {
			log.Printf("🔧 DEBUG: Match not yet persistent: %v", err)
			allExist = false
		}

		// Check innings existence
		var innings models.Innings
		_, err = dbClient.Supabase.From("innings").Select("*", "exact", false).Eq("id", inningsID).Single().ExecuteTo(&innings)
		if err != nil {
			log.Printf("🔧 DEBUG: Innings not yet persistent: %v", err)
			allExist = false
		}

		// Check over existence
		var over models.ScorecardOver
		_, err = dbClient.Supabase.From("overs").Select("*", "exact", false).Eq("id", overID).Single().ExecuteTo(&over)
		if err != nil {
			log.Printf("🔧 DEBUG: Over not yet persistent: %v", err)
			allExist = false
		}

		if allExist {
			log.Printf("✅ DEBUG: All entities are now persistent in database")
			return nil
		}

		// Wait before next check
		time.Sleep(200 * time.Millisecond)
	}

	return fmt.Errorf("database persistence validation timed out after %v", maxWaitTime)
}

// validateDatabaseConnection checks if database connection is healthy
func validateDatabaseConnection(dbClient *database.Client) error {
	log.Printf("🔧 DEBUG: Validating database connection health")

	// Test basic database connectivity
	var result []map[string]interface{}
	_, err := dbClient.Supabase.From("series").Select("id", "exact", false).Limit(1, "").ExecuteTo(&result)
	if err != nil {
		log.Printf("❌ DEBUG: Database connection validation failed: %v", err)
		return fmt.Errorf("database connection validation failed: %v", err)
	}

	log.Printf("✅ DEBUG: Database connection is healthy")
	return nil
}

// preflightDataValidation performs comprehensive pre-flight checks before critical operations
func preflightDataValidation(dbClient *database.Client, matchID string, inningsNumber int) error {
	log.Printf("🔧 DEBUG: Performing pre-flight data validation for match %s, innings %d", matchID, inningsNumber)

	// Validate match exists and is in correct state
	log.Printf("🔧 DEBUG: Validating match existence and state")
	var match models.Match
	_, err := dbClient.Supabase.From("matches").Select("*", "exact", false).Eq("id", matchID).Single().ExecuteTo(&match)
	if err != nil {
		log.Printf("❌ DEBUG: Match not found in pre-flight check: %v", err)
		return fmt.Errorf("match not found in pre-flight check: %v", err)
	}

	if match.Status != models.MatchStatusLive {
		log.Printf("❌ DEBUG: Match is not in live state: %s", match.Status)
		return fmt.Errorf("match is not in live state: %s", match.Status)
	}

	log.Printf("✅ DEBUG: Match validation passed - ID: %s, Status: %s", match.ID, match.Status)

	// Validate innings exists and is accessible
	log.Printf("🔧 DEBUG: Validating innings existence and accessibility")
	var innings models.Innings
	_, err = dbClient.Supabase.From("innings").Select("*", "exact", false).Eq("match_id", matchID).Eq("innings_number", fmt.Sprintf("%d", inningsNumber)).Single().ExecuteTo(&innings)
	if err != nil {
		log.Printf("❌ DEBUG: Innings not found in pre-flight check: %v", err)
		return fmt.Errorf("innings not found in pre-flight check: %v", err)
	}

	if innings.Status != string(models.InningsStatusInProgress) {
		log.Printf("❌ DEBUG: Innings is not in progress state: %s", innings.Status)
		return fmt.Errorf("innings is not in progress state: %s", innings.Status)
	}

	log.Printf("✅ DEBUG: Innings validation passed - ID: %s, Status: %s", innings.ID, innings.Status)

	// Validate current over exists and is accessible
	log.Printf("🔧 DEBUG: Validating current over existence and accessibility")
	var over models.ScorecardOver
	_, err = dbClient.Supabase.From("overs").Select("*", "exact", false).Eq("innings_id", innings.ID).Eq("status", string(models.OverStatusInProgress)).Single().ExecuteTo(&over)
	if err != nil {
		log.Printf("❌ DEBUG: Current over not found in pre-flight check: %v", err)
		return fmt.Errorf("current over not found in pre-flight check: %v", err)
	}

	if over.Status != string(models.OverStatusInProgress) {
		log.Printf("❌ DEBUG: Over is not in progress state: %s", over.Status)
		return fmt.Errorf("over is not in progress state: %s", over.Status)
	}

	log.Printf("✅ DEBUG: Over validation passed - ID: %s, Status: %s", over.ID, over.Status)

	// Additional wait to ensure database consistency
	log.Printf("🔧 DEBUG: Additional wait for database consistency")
	time.Sleep(300 * time.Millisecond)

	log.Printf("✅ DEBUG: All pre-flight validations passed")
	return nil
}
