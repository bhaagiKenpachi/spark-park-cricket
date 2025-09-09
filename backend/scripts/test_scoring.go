package main

import (
	"context"
	"fmt"
	"log"
	"spark-park-cricket-backend/internal/config"
	"spark-park-cricket-backend/internal/database"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/services"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize database client
	dbClient, err := database.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database client: %v", err)
	}

	ctx := context.Background()

	// Create test data
	fmt.Println("Creating test data...")

	// Create series
	series := &models.Series{
		Name:      "Test Series 2024",
		StartDate: time.Now(),
		EndDate:   time.Now().AddDate(0, 0, 7),
	}
	err = dbClient.Repositories.Series.Create(ctx, series)
	if err != nil {
		log.Printf("Failed to create series: %v", err)
	} else {
		fmt.Printf("Created series: %s (ID: %s)\n", series.Name, series.ID)
	}

	// Create teams
	team1 := &models.Team{
		Name:         "Team A",
		PlayersCount: 11,
	}
	team2 := &models.Team{
		Name:         "Team B",
		PlayersCount: 11,
	}

	err = dbClient.Repositories.Team.Create(ctx, team1)
	if err != nil {
		log.Printf("Failed to create team1: %v", err)
	} else {
		fmt.Printf("Created team: %s (ID: %s)\n", team1.Name, team1.ID)
	}

	err = dbClient.Repositories.Team.Create(ctx, team2)
	if err != nil {
		log.Printf("Failed to create team2: %v", err)
	} else {
		fmt.Printf("Created team: %s (ID: %s)\n", team2.Name, team2.ID)
	}

	// Create match
	match := &models.Match{
		SeriesID:    series.ID,
		MatchNumber: 1,
		Date:        time.Now(),
		Status:      models.MatchStatusLive,
		Team1ID:     team1.ID,
		Team2ID:     team2.ID,
	}
	err = dbClient.Repositories.Match.Create(ctx, match)
	if err != nil {
		log.Printf("Failed to create match: %v", err)
	} else {
		fmt.Printf("Created match: %s vs %s (ID: %s)\n", team1.Name, team2.Name, match.ID)
	}

	// Create players
	players := []*models.Player{
		{Name: "Batsman 1", TeamID: team1.ID},
		{Name: "Batsman 2", TeamID: team1.ID},
		{Name: "Bowler 1", TeamID: team2.ID},
	}

	for _, player := range players {
		err = dbClient.Repositories.Player.Create(ctx, player)
		if err != nil {
			log.Printf("Failed to create player: %v", err)
		} else {
			fmt.Printf("Created player: %s (ID: %s)\n", player.Name, player.ID)
		}
	}

	// Test live scoring
	fmt.Println("\nTesting live scoring system...")

	// Initialize services
	serviceContainer := services.NewContainer(dbClient.Repositories)
	scoreboardService := serviceContainer.Scoreboard

	// Get initial scoreboard
	scoreboard, err := scoreboardService.GetScoreboard(ctx, match.ID)
	if err != nil {
		log.Printf("Failed to get scoreboard: %v", err)
	} else {
		fmt.Printf("Initial scoreboard: %d/%d in %.1f overs\n",
			scoreboard.Score, scoreboard.Wickets, scoreboard.Overs)
	}

	// Add some balls
	ballEvents := []models.BallEvent{
		{BallType: models.BallTypeGood, Runs: 1, IsWicket: false, BatsmanID: players[0].ID, BowlerID: players[2].ID},
		{BallType: models.BallTypeGood, Runs: 4, IsWicket: false, BatsmanID: players[1].ID, BowlerID: players[2].ID},
		{BallType: models.BallTypeGood, Runs: 0, IsWicket: false, BatsmanID: players[1].ID, BowlerID: players[2].ID},
		{BallType: models.BallTypeWide, Runs: 1, IsWicket: false, BatsmanID: players[1].ID, BowlerID: players[2].ID},
		{BallType: models.BallTypeGood, Runs: 2, IsWicket: false, BatsmanID: players[1].ID, BowlerID: players[2].ID},
		{BallType: models.BallTypeGood, Runs: 6, IsWicket: false, BatsmanID: players[1].ID, BowlerID: players[2].ID},
	}

	for i, ballEvent := range ballEvents {
		scoreboard, err = scoreboardService.AddBall(ctx, match.ID, &ballEvent)
		if err != nil {
			log.Printf("Failed to add ball %d: %v", i+1, err)
		} else {
			fmt.Printf("Ball %d: %s, %d runs, %d/%d in %.1f overs\n",
				i+1, ballEvent.BallType, ballEvent.Runs,
				scoreboard.Score, scoreboard.Wickets, scoreboard.Overs)
		}
	}

	fmt.Println("\nLive scoring test completed successfully!")
}
