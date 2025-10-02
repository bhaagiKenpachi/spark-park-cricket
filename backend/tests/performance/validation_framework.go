package performance

import (
	"context"
	"fmt"
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

// createTestDataWithValidation creates test data with configurable validation level
func createTestDataWithValidation(dbClient *database.Client, userID string, level ValidationLevel, testType string) (string, string, error) {
	config := GetValidationConfig(level)

	// Apply mutex protection if configured
	var testDataMutex sync.Mutex
	if config.MutexProtection {
		testDataMutex.Lock()
		defer testDataMutex.Unlock()
		time.Sleep(config.RetryDelay)
	}

	ctx := context.Background()

	// Retry logic based on validation level
	for attempt := 1; attempt <= config.MaxRetries; attempt++ {
		if attempt > 1 {
			time.Sleep(time.Duration(attempt) * config.RetryDelay)
		}

		// Create test series with unique identifier
		uniqueID := fmt.Sprintf("%d-%d-%d-%s", time.Now().UnixNano(), time.Now().Unix(), attempt, testType)
		series := &models.Series{
			Name:      fmt.Sprintf("%s Test Series %s", testType, uniqueID),
			StartDate: time.Now(),
			EndDate:   time.Now().Add(24 * time.Hour),
			CreatedBy: userID,
		}

		err := dbClient.Repositories.Series.Create(ctx, series)
		if err != nil {
			if attempt == config.MaxRetries {
				return "", "", fmt.Errorf("failed to create test series after %d attempts: %v", config.MaxRetries, err)
			}
			continue
		}

		// Create test match
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
			if attempt == config.MaxRetries {
				return "", "", fmt.Errorf("failed to create test match after %d attempts: %v", config.MaxRetries, err)
			}
			continue
		}

		// Create first innings
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
			if attempt == config.MaxRetries {
				return "", "", fmt.Errorf("failed to create test innings after %d attempts: %v", config.MaxRetries, err)
			}
			continue
		}

		// Create first over
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
			if attempt == config.MaxRetries {
				return "", "", fmt.Errorf("failed to create test over after %d attempts: %v", config.MaxRetries, err)
			}
			continue
		}

		// Apply validation based on level
		if config.ValidateDataIntegrity {
			if err := validateTestDataIntegrityWithLevel(dbClient, series.ID, match.ID, innings.ID, over.ID, level); err != nil {
				if attempt == config.MaxRetries {
					return "", "", fmt.Errorf("failed to validate test data integrity after %d attempts: %v", config.MaxRetries, err)
				}
				continue
			}
		}

		// Always validate basic data creation for all levels
		if series.ID == "" || match.ID == "" || innings.ID == "" || over.ID == "" {
			if attempt == config.MaxRetries {
				return "", "", fmt.Errorf("failed to create complete test data hierarchy after %d attempts", config.MaxRetries)
			}
			continue
		}

		// Success - return the IDs
		return series.ID, match.ID, nil
	}

	return "", "", fmt.Errorf("unexpected error in createTestDataWithValidation")
}

// validateTestDataIntegrityWithLevel validates data based on validation level
func validateTestDataIntegrityWithLevel(dbClient *database.Client, seriesID, matchID, inningsID, overID string, level ValidationLevel) error {
	config := GetValidationConfig(level)

	// Basic validation - just check if IDs are not empty
	if !config.ValidateRelationships {
		if seriesID == "" || matchID == "" || inningsID == "" || overID == "" {
			return fmt.Errorf("one or more IDs are empty")
		}
		return nil
	}

	// Intermediate and Full validation - check relationships
	var series models.Series
	_, err := dbClient.Supabase.From("series").Select("*", "exact", false).Eq("id", seriesID).Single().ExecuteTo(&series)
	if err != nil {
		return fmt.Errorf("series validation failed: %v", err)
	}

	var match models.Match
	_, err = dbClient.Supabase.From("matches").Select("*", "exact", false).Eq("id", matchID).Single().ExecuteTo(&match)
	if err != nil {
		return fmt.Errorf("match validation failed: %v", err)
	}

	if config.ValidateRelationships && match.SeriesID != seriesID {
		return fmt.Errorf("match series reference invalid: expected %s, got %s", seriesID, match.SeriesID)
	}

	var innings models.Innings
	_, err = dbClient.Supabase.From("innings").Select("*", "exact", false).Eq("id", inningsID).Single().ExecuteTo(&innings)
	if err != nil {
		return fmt.Errorf("innings validation failed: %v", err)
	}

	if config.ValidateRelationships && innings.MatchID != matchID {
		return fmt.Errorf("innings match reference invalid: expected %s, got %s", matchID, innings.MatchID)
	}

	var over models.ScorecardOver
	_, err = dbClient.Supabase.From("overs").Select("*", "exact", false).Eq("id", overID).Single().ExecuteTo(&over)
	if err != nil {
		return fmt.Errorf("over validation failed: %v", err)
	}

	if config.ValidateRelationships && over.InningsID != inningsID {
		return fmt.Errorf("over innings reference invalid: expected %s, got %s", inningsID, over.InningsID)
	}

	// Full validation - additional data integrity checks
	if config.ValidateDataIntegrity {
		if series.ID == "" || match.ID == "" || innings.ID == "" || over.ID == "" {
			return fmt.Errorf("one or more entity IDs are empty")
		}

		// Check for reasonable data values
		if innings.InningsNumber < 1 || innings.InningsNumber > 4 {
			return fmt.Errorf("invalid innings number: %d", innings.InningsNumber)
		}

		if over.OverNumber < 1 {
			return fmt.Errorf("invalid over number: %d", over.OverNumber)
		}
	}

	return nil
}
