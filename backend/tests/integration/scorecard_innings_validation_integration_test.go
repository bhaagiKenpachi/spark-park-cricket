package tests

import (
	"context"
	"testing"
	"time"

	"spark-park-cricket-backend/internal/config"
	"spark-park-cricket-backend/internal/database"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScorecardInningsValidation_Integration(t *testing.T) {
	// Setup test database
	cfg := config.LoadTestConfig()
	testDB, err := database.NewTestClient(cfg)
	require.NoError(t, err)
	defer testDB.Close()

	// Clean up before test
	cleanupTestData(t, testDB)

	// Use repositories from test client
	seriesRepo := testDB.Repositories.Series
	matchRepo := testDB.Repositories.Match
	scorecardRepo := testDB.Repositories.Scorecard

	// Create services
	seriesService := services.NewSeriesService(seriesRepo)
	matchService := services.NewMatchService(matchRepo, seriesRepo)
	scorecardService := services.NewScorecardService(scorecardRepo, matchRepo)

	ctx := context.Background()

	t.Run("First ball must be played by toss-winning team", func(t *testing.T) {
		// Create a test series
		seriesReq := &models.CreateSeriesRequest{
			Name:      "Test Innings Validation Series",
			StartDate: time.Now(),
			EndDate:   time.Now().AddDate(0, 0, 7),
		}

		series, err := seriesService.CreateSeries(ctx, seriesReq)
		require.NoError(t, err)
		require.NotNil(t, series)

		// Create a match with Team A winning toss
		matchReq := &models.CreateMatchRequest{
			SeriesID:         series.ID,
			MatchNumber:      nil, // Auto-increment
			Date:             time.Now(),
			TeamAPlayerCount: 11,
			TeamBPlayerCount: 11,
			TotalOvers:       20,
			TossWinner:       models.TeamTypeA,
			TossType:         models.TossTypeHeads,
		}

		match, err := matchService.CreateMatch(ctx, matchReq)
		require.NoError(t, err)
		require.NotNil(t, match)
		assert.Equal(t, models.TeamTypeA, match.TossWinner)
		assert.Equal(t, models.TeamTypeA, match.BattingTeam) // Should be set to toss winner

		// Try to add a ball to first innings - this should work since Team A is the toss winner
		ballEvent := &models.BallEventRequest{
			MatchID:       match.ID,
			InningsNumber: 1,
			BallType:      models.BallTypeGood,
			RunType:       models.RunTypeOne,
			IsWicket:      false,
		}

		err = scorecardService.AddBall(ctx, ballEvent)
		require.NoError(t, err)
	})

	t.Run("Cannot add ball to second innings before first innings", func(t *testing.T) {
		// Create a test series
		seriesReq := &models.CreateSeriesRequest{
			Name:      "Test Innings Validation Series 2",
			StartDate: time.Now(),
			EndDate:   time.Now().AddDate(0, 0, 7),
		}

		series, err := seriesService.CreateSeries(ctx, seriesReq)
		require.NoError(t, err)
		require.NotNil(t, series)

		// Create a match with Team A winning toss
		matchReq := &models.CreateMatchRequest{
			SeriesID:         series.ID,
			MatchNumber:      nil, // Auto-increment
			Date:             time.Now(),
			TeamAPlayerCount: 11,
			TeamBPlayerCount: 11,
			TotalOvers:       20,
			TossWinner:       models.TeamTypeA,
			TossType:         models.TossTypeHeads,
		}

		match, err := matchService.CreateMatch(ctx, matchReq)
		require.NoError(t, err)
		require.NotNil(t, match)

		// Try to add a ball to second innings directly - this should fail
		ballEvent := &models.BallEventRequest{
			MatchID:       match.ID,
			InningsNumber: 2,
			BallType:      models.BallTypeGood,
			RunType:       models.RunTypeOne,
			IsWicket:      false,
		}

		err = scorecardService.AddBall(ctx, ballEvent)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot start second innings, first innings must be played first")
	})

	t.Run("Cannot add ball to wrong team in first innings", func(t *testing.T) {
		// Create a test series
		seriesReq := &models.CreateSeriesRequest{
			Name:      "Test Innings Validation Series 3",
			StartDate: time.Now(),
			EndDate:   time.Now().AddDate(0, 0, 7),
		}

		series, err := seriesService.CreateSeries(ctx, seriesReq)
		require.NoError(t, err)
		require.NotNil(t, series)

		// Create a match with Team A winning toss
		matchReq := &models.CreateMatchRequest{
			SeriesID:         series.ID,
			MatchNumber:      nil, // Auto-increment
			Date:             time.Now(),
			TeamAPlayerCount: 11,
			TeamBPlayerCount: 11,
			TotalOvers:       20,
			TossWinner:       models.TeamTypeA,
			TossType:         models.TossTypeHeads,
		}

		match, err := matchService.CreateMatch(ctx, matchReq)
		require.NoError(t, err)
		require.NotNil(t, match)

		// Manually change the batting team to Team B (non-toss winner)
		match.BattingTeam = models.TeamTypeB
		err = matchRepo.Update(ctx, match.ID, match)
		require.NoError(t, err)

		// Try to add a ball to first innings with Team B - this should fail
		ballEvent := &models.BallEventRequest{
			MatchID:       match.ID,
			InningsNumber: 1,
			BallType:      models.BallTypeGood,
			RunType:       models.RunTypeOne,
			IsWicket:      false,
		}

		err = scorecardService.AddBall(ctx, ballEvent)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "first innings must be played by the toss-winning team")
	})

	t.Run("Second innings can only start after first innings is complete", func(t *testing.T) {
		// Create a test series
		seriesReq := &models.CreateSeriesRequest{
			Name:      "Test Innings Validation Series 4",
			StartDate: time.Now(),
			EndDate:   time.Now().AddDate(0, 0, 7),
		}

		series, err := seriesService.CreateSeries(ctx, seriesReq)
		require.NoError(t, err)
		require.NotNil(t, series)

		// Create a match with Team A winning toss
		matchReq := &models.CreateMatchRequest{
			SeriesID:         series.ID,
			MatchNumber:      nil, // Auto-increment
			Date:             time.Now(),
			TeamAPlayerCount: 11,
			TeamBPlayerCount: 11,
			TotalOvers:       20,
			TossWinner:       models.TeamTypeA,
			TossType:         models.TossTypeHeads,
		}

		match, err := matchService.CreateMatch(ctx, matchReq)
		require.NoError(t, err)
		require.NotNil(t, match)

		// Add a few balls to the first innings (Team A)
		ballEvent := &models.BallEventRequest{
			MatchID:       match.ID,
			InningsNumber: 1,
			BallType:      models.BallTypeGood,
			RunType:       models.RunTypeOne,
			IsWicket:      false,
		}

		err = scorecardService.AddBall(ctx, ballEvent)
		require.NoError(t, err)

		// Try to add a ball to second innings before first innings is complete - this should fail
		ballEvent2 := &models.BallEventRequest{
			MatchID:       match.ID,
			InningsNumber: 2,
			BallType:      models.BallTypeGood,
			RunType:       models.RunTypeOne,
			IsWicket:      false,
		}

		err = scorecardService.AddBall(ctx, ballEvent2)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "first innings is not complete, cannot start second innings")
	})

	t.Run("Second innings must be played by non-toss-winning team", func(t *testing.T) {
		// Create a test series
		seriesReq := &models.CreateSeriesRequest{
			Name:      "Test Innings Validation Series 5",
			StartDate: time.Now(),
			EndDate:   time.Now().AddDate(0, 0, 7),
		}

		series, err := seriesService.CreateSeries(ctx, seriesReq)
		require.NoError(t, err)
		require.NotNil(t, series)

		// Create a match with Team A winning toss
		matchReq := &models.CreateMatchRequest{
			SeriesID:         series.ID,
			MatchNumber:      nil, // Auto-increment
			Date:             time.Now(),
			TeamAPlayerCount: 11,
			TeamBPlayerCount: 11,
			TotalOvers:       20,
			TossWinner:       models.TeamTypeA,
			TossType:         models.TossTypeHeads,
		}

		match, err := matchService.CreateMatch(ctx, matchReq)
		require.NoError(t, err)
		require.NotNil(t, match)

		// Complete the first innings by taking all wickets
		// Add balls until we have 10 wickets
		wicketBallEvent := &models.BallEventRequest{
			MatchID:       match.ID,
			InningsNumber: 1,
			BallType:      models.BallTypeGood,
			RunType:       models.RunTypeZero, // No runs with wicket
			IsWicket:      true,               // Wicket ball
			WicketType:    "bowled",           // Required wicket type
		}

		// Add 10 wicket balls to complete first innings
		for i := 0; i < 10; i++ {
			err = scorecardService.AddBall(ctx, wicketBallEvent)
			require.NoError(t, err)
		}

		// Now first innings should be complete
		// Try to change to Team B (non-toss winner) for second innings
		match.BattingTeam = models.TeamTypeB
		err = matchRepo.Update(ctx, match.ID, match)
		require.NoError(t, err)

		// Try to add a ball to Team B - this should work now
		ballEvent2 := &models.BallEventRequest{
			MatchID:       match.ID,
			InningsNumber: 2,
			BallType:      models.BallTypeGood,
			RunType:       models.RunTypeOne,
			IsWicket:      false,
		}

		err = scorecardService.AddBall(ctx, ballEvent2)
		require.NoError(t, err)

		// Try to change back to Team A (toss winner) - this should fail
		match.BattingTeam = models.TeamTypeA
		err = matchRepo.Update(ctx, match.ID, match)
		require.NoError(t, err)

		// Try to add a ball to Team A - this should fail
		ballEvent3 := &models.BallEventRequest{
			MatchID:       match.ID,
			InningsNumber: 2,
			BallType:      models.BallTypeGood,
			RunType:       models.RunTypeOne,
			IsWicket:      false,
		}

		err = scorecardService.AddBall(ctx, ballEvent3)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "second innings must be played by the non-toss-winning team")
	})

	// Clean up after test
	cleanupTestData(t, testDB)
}

func cleanupTestData(t *testing.T, testDB *database.Client) {
	// Clean up matches
	_, err := testDB.Supabase.From("matches").Delete("", "").ExecuteTo(nil)
	if err != nil {
		t.Logf("Warning: Failed to cleanup matches: %v", err)
	}

	// Clean up series
	_, err = testDB.Supabase.From("series").Delete("", "").ExecuteTo(nil)
	if err != nil {
		t.Logf("Warning: Failed to cleanup series: %v", err)
	}
}
