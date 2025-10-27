package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/spark-park-cricket/backend/internal/config"
	"github.com/spark-park-cricket/backend/internal/database"
	"github.com/spark-park-cricket/backend/internal/handlers"
	"github.com/spark-park-cricket/backend/internal/models"
	"github.com/spark-park-cricket/backend/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFallOfWicketsIntegration(t *testing.T) {
	// Skip if not in integration test mode
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup test database
	cfg := &config.Config{
		SupabaseURL:    "http://localhost:54321",
		SupabaseAPIKey: "test-key",
		DatabaseSchema: "dev_v1",
		CacheEnabled:   false,
		AllowedOrigins: "*",
		PrometheusURL:  "http://localhost:9090",
		GrafanaURL:     "http://localhost:3000",
	}

	dbClient, err := database.NewClient(cfg)
	require.NoError(t, err)
	defer dbClient.Close()

	// Create service container
	serviceContainer := services.NewContainer(dbClient, cfg)

	// Create test match
	match, err := createTestMatch(serviceContainer)
	require.NoError(t, err)

	// Create test innings
	innings, err := createTestInnings(serviceContainer, match.ID)
	require.NoError(t, err)

	// Create test over
	over, err := createTestOver(serviceContainer, innings.ID)
	require.NoError(t, err)

	// Create test ball with wicket
	ball, err := createTestWicketBall(serviceContainer, over.ID)
	require.NoError(t, err)

	t.Run("Create Fall of Wickets Record", func(t *testing.T) {
		// Create fall of wickets record
		req := &models.CreateFallOfWicketsRequest{
			MatchID:       match.ID,
			InningsID:     innings.ID,
			OverID:        over.ID,
			BallID:        ball.ID,
			WicketNumber:  1,
			BatsmanName:   stringPtr("Test Batsman"),
			Runs:          25,
			BallsFaced:    30,
			WicketType:    "bowled",
			BowlerName:    stringPtr("Test Bowler"),
			OverNumber:    over.OverNumber,
			BallNumber:    ball.BallNumber,
			ScoreAtWicket: 25,
		}

		result, err := serviceContainer.FallOfWicketsService.CreateFallOfWickets(context.Background(), req)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, match.ID, result.MatchID)
		assert.Equal(t, innings.ID, result.InningsID)
		assert.Equal(t, over.ID, result.OverID)
		assert.Equal(t, ball.ID, result.BallID)
		assert.Equal(t, 1, result.WicketNumber)
		assert.Equal(t, "Test Batsman", *result.BatsmanName)
		assert.Equal(t, 25, result.Runs)
		assert.Equal(t, 30, result.BallsFaced)
		assert.Equal(t, "bowled", result.WicketType)
		assert.Equal(t, "Test Bowler", *result.BowlerName)
		assert.Equal(t, over.OverNumber, result.OverNumber)
		assert.Equal(t, ball.BallNumber, result.BallNumber)
		assert.Equal(t, 25, result.ScoreAtWicket)
	})

	t.Run("Get Fall of Wickets Summary", func(t *testing.T) {
		summary, err := serviceContainer.FallOfWicketsService.GetFallOfWicketsSummary(context.Background(), match.ID, nil)
		require.NoError(t, err)
		assert.NotNil(t, summary)
		assert.Equal(t, match.ID, summary.MatchID)
		assert.Equal(t, innings.ID, summary.InningsID)
		assert.Equal(t, 1, summary.TotalWickets)
		assert.Len(t, summary.Wickets, 1)

		wicket := summary.Wickets[0]
		assert.Equal(t, 1, wicket.WicketNumber)
		assert.Equal(t, "Test Batsman", *wicket.BatsmanName)
		assert.Equal(t, 25, wicket.Runs)
		assert.Equal(t, 30, wicket.BallsFaced)
		assert.Equal(t, "bowled", wicket.WicketType)
		assert.Equal(t, "Test Bowler", *wicket.BowlerName)
		assert.Equal(t, over.OverNumber, wicket.OverNumber)
		assert.Equal(t, ball.BallNumber, wicket.BallNumber)
		assert.Equal(t, 25, wicket.ScoreAtWicket)
	})

	t.Run("Get Fall of Wickets by Match ID", func(t *testing.T) {
		fallOfWickets, err := serviceContainer.FallOfWicketsService.GetFallOfWicketsByMatchID(context.Background(), match.ID)
		require.NoError(t, err)
		assert.Len(t, fallOfWickets, 1)

		fow := fallOfWickets[0]
		assert.Equal(t, match.ID, fow.MatchID)
		assert.Equal(t, innings.ID, fow.InningsID)
		assert.Equal(t, over.ID, fow.OverID)
		assert.Equal(t, ball.ID, fow.BallID)
		assert.Equal(t, 1, fow.WicketNumber)
	})

	t.Run("Get Fall of Wickets by Innings ID", func(t *testing.T) {
		fallOfWickets, err := serviceContainer.FallOfWicketsService.GetFallOfWicketsByInningsID(context.Background(), innings.ID)
		require.NoError(t, err)
		assert.Len(t, fallOfWickets, 1)

		fow := fallOfWickets[0]
		assert.Equal(t, innings.ID, fow.InningsID)
		assert.Equal(t, 1, fow.WicketNumber)
	})

	t.Run("Get Fall of Wickets by Ball ID", func(t *testing.T) {
		fallOfWickets, err := serviceContainer.FallOfWicketsService.GetFallOfWicketsByBallID(context.Background(), ball.ID)
		require.NoError(t, err)
		assert.NotNil(t, fallOfWickets)
		assert.Equal(t, ball.ID, fallOfWickets.BallID)
		assert.Equal(t, 1, fallOfWickets.WicketNumber)
	})

	t.Run("Update Fall of Wickets", func(t *testing.T) {
		// First get the fall of wickets record
		fallOfWickets, err := serviceContainer.FallOfWicketsService.GetFallOfWicketsByBallID(context.Background(), ball.ID)
		require.NoError(t, err)

		// Update the record
		updateReq := &models.UpdateFallOfWicketsRequest{
			BatsmanName: stringPtr("Updated Batsman"),
			Runs:        intPtr(30),
			BallsFaced:  intPtr(35),
		}

		updated, err := serviceContainer.FallOfWicketsService.UpdateFallOfWickets(context.Background(), fallOfWickets.ID, updateReq)
		require.NoError(t, err)
		assert.NotNil(t, updated)
		assert.Equal(t, "Updated Batsman", *updated.BatsmanName)
		assert.Equal(t, 30, updated.Runs)
		assert.Equal(t, 35, updated.BallsFaced)
	})

	t.Run("Create Fall of Wickets from Ball", func(t *testing.T) {
		// Create another ball with wicket
		ball2, err := createTestWicketBall(serviceContainer, over.ID)
		require.NoError(t, err)

		// Create fall of wickets record automatically from ball
		result, err := serviceContainer.FallOfWicketsService.CreateFallOfWicketsFromBall(
			context.Background(),
			ball2.ID,
			stringPtr("Auto Batsman"),
			stringPtr("Auto Bowler"),
			stringPtr("Auto Fielder"),
		)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ball2.ID, result.BallID)
		assert.Equal(t, "Auto Batsman", *result.BatsmanName)
		assert.Equal(t, "Auto Bowler", *result.BowlerName)
		assert.Equal(t, "Auto Fielder", *result.FielderName)
		assert.Equal(t, 2, result.WicketNumber) // Should be the second wicket
	})

	t.Run("HTTP API Tests", func(t *testing.T) {
		// Setup HTTP handler
		fallOfWicketsHandler := handlers.NewFallOfWicketsHandler(serviceContainer.FallOfWicketsService)

		t.Run("GET Fall of Wickets by Match ID", func(t *testing.T) {
			req, err := http.NewRequest("GET", fmt.Sprintf("/api/v1/matches/%s/fall-of-wickets", match.ID), nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			fallOfWicketsHandler.GetFallOfWicketsByMatchID(rr, req)

			assert.Equal(t, http.StatusOK, rr.Code)

			var response []models.FallOfWickets
			err = json.Unmarshal(rr.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.Len(t, response, 2) // Should have 2 wickets now
		})

		t.Run("GET Fall of Wickets Summary", func(t *testing.T) {
			req, err := http.NewRequest("GET", fmt.Sprintf("/api/v1/matches/%s/fall-of-wickets/summary", match.ID), nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			fallOfWicketsHandler.GetFallOfWicketsSummary(rr, req)

			assert.Equal(t, http.StatusOK, rr.Code)

			var response models.FallOfWicketsSummary
			err = json.Unmarshal(rr.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.Equal(t, match.ID, response.MatchID)
			assert.Equal(t, 2, response.TotalWickets)
			assert.Len(t, response.Wickets, 2)
		})

		t.Run("GET Fall of Wickets by Innings ID", func(t *testing.T) {
			req, err := http.NewRequest("GET", fmt.Sprintf("/api/v1/innings/%s/fall-of-wickets", innings.ID), nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			fallOfWicketsHandler.GetFallOfWicketsByInningsID(rr, req)

			assert.Equal(t, http.StatusOK, rr.Code)

			var response []models.FallOfWickets
			err = json.Unmarshal(rr.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.Len(t, response, 2)
		})
	})
}

// Helper functions for test setup
func createTestMatch(serviceContainer *services.Container) (*models.Match, error) {
	// Create a test series first
	seriesReq := &models.CreateSeriesRequest{
		Name:      "Test Series",
		StartDate: time.Now(),
		EndDate:   time.Now().Add(24 * time.Hour),
	}

	series, err := serviceContainer.Series.CreateSeries(context.Background(), seriesReq)
	if err != nil {
		return nil, err
	}

	// Create test match
	matchReq := &models.CreateMatchRequest{
		SeriesID:         series.ID,
		MatchNumber:      intPtr(1),
		Date:             time.Now(),
		TeamAPlayerCount: 11,
		TeamBPlayerCount: 11,
		TotalOvers:       20,
		TossWinner:       models.TeamTypeA,
		TossType:         models.TossTypeHeads,
	}

	return serviceContainer.Match.CreateMatch(context.Background(), matchReq)
}

func createTestInnings(serviceContainer *services.Container, matchID string) (*models.ScorecardInnings, error) {
	// Start scoring to create innings
	startReq := &models.StartScoringRequest{
		MatchID: matchID,
	}

	err := serviceContainer.Scorecard.StartScoring(context.Background(), startReq)
	if err != nil {
		return nil, err
	}

	// Get the innings
	scorecard, err := serviceContainer.Scorecard.GetScorecard(context.Background(), matchID)
	if err != nil {
		return nil, err
	}

	if len(scorecard.Innings) == 0 {
		return nil, fmt.Errorf("no innings found")
	}

	// Get the actual innings from database
	innings, err := serviceContainer.Scorecard.GetInningsByMatchAndNumber(context.Background(), matchID, 1)
	if err != nil {
		return nil, err
	}

	return innings, nil
}

func createTestOver(serviceContainer *services.Container, inningsID string) (*models.ScorecardOver, error) {
	// Get innings to find match ID
	innings, err := serviceContainer.Scorecard.GetInningsByID(context.Background(), inningsID)
	if err != nil {
		return nil, err
	}

	// Add a ball to create an over
	ballReq := &models.BallEventRequest{
		MatchID:       innings.MatchID,
		InningsNumber: 1,
		BallType:      models.BallTypeGood,
		RunType:       models.RunTypeOne,
		Runs:          1,
		Byes:          0,
		IsWicket:      false,
	}

	err = serviceContainer.Scorecard.AddBall(context.Background(), ballReq)
	if err != nil {
		return nil, err
	}

	// Get the over
	over, err := serviceContainer.Scorecard.GetCurrentOver(context.Background(), innings.MatchID, 1)
	if err != nil {
		return nil, err
	}

	return over, nil
}

func createTestWicketBall(serviceContainer *services.Container, overID string) (*models.ScorecardBall, error) {
	// Get over to find match ID and innings ID
	over, err := serviceContainer.Scorecard.GetOverByID(context.Background(), overID)
	if err != nil {
		return nil, err
	}

	// Get innings to find match ID
	innings, err := serviceContainer.Scorecard.GetInningsByID(context.Background(), over.InningsID)
	if err != nil {
		return nil, err
	}

	// Add a wicket ball
	ballReq := &models.BallEventRequest{
		MatchID:       innings.MatchID,
		InningsNumber: 1,
		BallType:      models.BallTypeGood,
		RunType:       models.RunTypeWC,
		Runs:          0,
		Byes:          0,
		IsWicket:      true,
		WicketType:    "bowled",
	}

	err = serviceContainer.Scorecard.AddBall(context.Background(), ballReq)
	if err != nil {
		return nil, err
	}

	// Get the ball that was just created
	balls, err := serviceContainer.Scorecard.GetBallsByOver(context.Background(), overID)
	if err != nil {
		return nil, err
	}

	if len(balls) == 0 {
		return nil, fmt.Errorf("no balls found")
	}

	// Find the wicket ball
	for _, ball := range balls {
		if ball.IsWicket {
			return ball, nil
		}
	}

	return nil, fmt.Errorf("no wicket ball found")
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}
