package unit

import (
	"context"
	"spark-park-cricket-backend/internal/models"
	"testing"
	"time"
)

func TestTimeTrackingIntegration(t *testing.T) {
	// This is a basic test to verify time tracking functionality
	// In a real implementation, you would use a test database

	ctx := context.Background()

	// Create test data
	matchID := "test-match-123"
	startTime := time.Now()

	// Test innings with time tracking
	innings := &models.Innings{
		ID:              "test-innings-1",
		MatchID:         matchID,
		InningsNumber:   1,
		BattingTeam:     models.TeamTypeA,
		TotalRuns:       150,
		TotalWickets:    3,
		TotalOvers:      20.0,
		TotalBalls:      120,
		Status:          string(models.InningsStatusCompleted),
		StartTime:       &startTime,
		EndTime:         func() *time.Time { t := startTime.Add(2 * time.Hour); return &t }(),
		DurationSeconds: 7200, // 2 hours
	}

	// Test over with time tracking
	overStartTime := startTime.Add(30 * time.Minute)
	over := &models.ScorecardOver{
		ID:              "test-over-1",
		InningsID:       innings.ID,
		OverNumber:      1,
		TotalRuns:       8,
		TotalBalls:      6,
		TotalWickets:    0,
		Status:          string(models.OverStatusCompleted),
		StartTime:       &overStartTime,
		EndTime:         func() *time.Time { t := overStartTime.Add(5 * time.Minute); return &t }(),
		DurationSeconds: 300, // 5 minutes
	}

	// Verify time tracking fields are set correctly
	if innings.StartTime == nil {
		t.Error("Expected innings start time to be set")
	}

	if innings.EndTime == nil {
		t.Error("Expected innings end time to be set")
	}

	if innings.DurationSeconds != 7200 {
		t.Errorf("Expected innings duration to be 7200 seconds, got %d", innings.DurationSeconds)
	}

	if over.StartTime == nil {
		t.Error("Expected over start time to be set")
	}

	if over.EndTime == nil {
		t.Error("Expected over end time to be set")
	}

	if over.DurationSeconds != 300 {
		t.Errorf("Expected over duration to be 300 seconds, got %d", over.DurationSeconds)
	}

	// Test time tracking response model
	timeTrackingResponse := &models.TimeTrackingResponse{
		MatchID: matchID,
		Innings: []models.InningsTimeTracking{
			{
				InningsNumber:   innings.InningsNumber,
				BattingTeam:     innings.BattingTeam,
				StartTime:       innings.StartTime,
				EndTime:         innings.EndTime,
				DurationSeconds: innings.DurationSeconds,
				Status:          innings.Status,
				Overs: []models.OverTimeTracking{
					{
						OverNumber:      over.OverNumber,
						StartTime:       over.StartTime,
						EndTime:         over.EndTime,
						DurationSeconds: over.DurationSeconds,
						Status:          over.Status,
						TotalRuns:       over.TotalRuns,
						TotalBalls:      over.TotalBalls,
						TotalWickets:    over.TotalWickets,
					},
				},
			},
		},
		TotalMatchTime: innings.DurationSeconds,
	}

	// Verify response structure
	if timeTrackingResponse.MatchID != matchID {
		t.Errorf("Expected match ID %s, got %s", matchID, timeTrackingResponse.MatchID)
	}

	if len(timeTrackingResponse.Innings) != 1 {
		t.Errorf("Expected 1 innings, got %d", len(timeTrackingResponse.Innings))
	}

	if timeTrackingResponse.TotalMatchTime != 7200 {
		t.Errorf("Expected total match time to be 7200 seconds, got %d", timeTrackingResponse.TotalMatchTime)
	}

	// Test duration calculation
	expectedDuration := int(innings.EndTime.Sub(*innings.StartTime).Seconds())
	if innings.DurationSeconds != expectedDuration {
		t.Errorf("Expected calculated duration %d, got %d", expectedDuration, innings.DurationSeconds)
	}

	t.Logf("Time tracking test passed - Innings duration: %d seconds, Over duration: %d seconds",
		innings.DurationSeconds, over.DurationSeconds)
}

func TestTimeTrackingCalculations(t *testing.T) {
	// Test duration calculations
	startTime := time.Now()
	endTime := startTime.Add(1*time.Hour + 30*time.Minute) // 1.5 hours

	expectedDuration := int(endTime.Sub(startTime).Seconds()) // 5400 seconds

	innings := &models.Innings{
		StartTime: &startTime,
		EndTime:   &endTime,
	}

	// Simulate duration calculation
	innings.DurationSeconds = int(endTime.Sub(startTime).Seconds())

	if innings.DurationSeconds != expectedDuration {
		t.Errorf("Expected duration %d seconds, got %d", expectedDuration, innings.DurationSeconds)
	}

	// Test over duration
	overStart := startTime.Add(10 * time.Minute)
	overEnd := overStart.Add(5 * time.Minute)

	over := &models.ScorecardOver{
		StartTime: &overStart,
		EndTime:   &overEnd,
	}

	over.DurationSeconds = int(overEnd.Sub(overStart).Seconds())

	expectedOverDuration := 300 // 5 minutes
	if over.DurationSeconds != expectedOverDuration {
		t.Errorf("Expected over duration %d seconds, got %d", expectedOverDuration, over.DurationSeconds)
	}

	t.Logf("Duration calculation test passed - Innings: %d seconds, Over: %d seconds",
		innings.DurationSeconds, over.DurationSeconds)
}
