package main

import (
	"context"
	"fmt"
	"log"
	"spark-park-cricket-backend/internal/config"
	"spark-park-cricket-backend/internal/database"
	"spark-park-cricket-backend/internal/models"
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

	// Create sample series
	series := &models.Series{
		Name:      "IPL 2024",
		StartDate: time.Date(2024, 3, 22, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 5, 26, 23, 59, 59, 0, time.UTC),
	}

	err = dbClient.Repositories.Series.Create(ctx, series)
	if err != nil {
		log.Printf("Failed to create series: %v", err)
	} else {
		fmt.Printf("Created series: %s (ID: %s)\n", series.Name, series.ID)
	}

	// Create sample teams
	team1 := &models.Team{
		Name:         "Mumbai Indians",
		PlayersCount: 11,
	}

	team2 := &models.Team{
		Name:         "Chennai Super Kings",
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

	// Create sample match
	match := &models.Match{
		SeriesID:    series.ID,
		MatchNumber: 1,
		Date:        time.Date(2024, 3, 22, 19, 30, 0, 0, time.UTC),
		Status:      models.MatchStatusScheduled,
		Team1ID:     team1.ID,
		Team2ID:     team2.ID,
	}

	err = dbClient.Repositories.Match.Create(ctx, match)
	if err != nil {
		log.Printf("Failed to create match: %v", err)
	} else {
		fmt.Printf("Created match: %s vs %s (ID: %s)\n", team1.Name, team2.Name, match.ID)
	}

	fmt.Println("Database seeding completed successfully!")
}
