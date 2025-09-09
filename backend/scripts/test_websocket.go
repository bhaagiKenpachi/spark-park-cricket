package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"spark-park-cricket-backend/internal/config"
	"spark-park-cricket-backend/internal/database"
	"spark-park-cricket-backend/internal/models"
	"spark-park-cricket-backend/internal/services"
	"time"

	"github.com/gorilla/websocket"
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
	fmt.Println("Creating test data for WebSocket testing...")

	// Create series
	series := &models.Series{
		Name:      "WebSocket Test Series 2024",
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
		Name:         "WebSocket Team A",
		PlayersCount: 11,
	}
	team2 := &models.Team{
		Name:         "WebSocket Team B",
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
		{Name: "WebSocket Batsman 1", TeamID: team1.ID},
		{Name: "WebSocket Batsman 2", TeamID: team1.ID},
		{Name: "WebSocket Bowler 1", TeamID: team2.ID},
	}

	for _, player := range players {
		err = dbClient.Repositories.Player.Create(ctx, player)
		if err != nil {
			log.Printf("Failed to create player: %v", err)
		} else {
			fmt.Printf("Created player: %s (ID: %s)\n", player.Name, player.ID)
		}
	}

	// Test WebSocket connection
	fmt.Println("\nTesting WebSocket connection...")

	// WebSocket URL
	wsURL := url.URL{
		Scheme: "ws",
		Host:   "localhost:8080",
		Path:   fmt.Sprintf("/api/v1/ws/match/%s", match.ID),
	}

	// Connect to WebSocket
	conn, _, err := websocket.DefaultDialer.Dial(wsURL.String(), nil)
	if err != nil {
		log.Printf("Failed to connect to WebSocket: %v", err)
		return
	}
	defer conn.Close()

	fmt.Printf("Connected to WebSocket for match %s\n", match.ID)

	// Start listening for messages in a goroutine
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Printf("WebSocket read error: %v", err)
				return
			}

			var wsMessage map[string]interface{}
			if err := json.Unmarshal(message, &wsMessage); err != nil {
				log.Printf("Failed to unmarshal WebSocket message: %v", err)
				continue
			}

			fmt.Printf("Received WebSocket message: %s\n", string(message))
		}
	}()

	// Test live scoring with WebSocket updates
	fmt.Println("\nTesting live scoring with WebSocket updates...")

	// Initialize services
	serviceContainer := services.NewContainer(dbClient.Repositories)
	scoreboardService := serviceContainer.Scoreboard

	// Add some balls and watch for WebSocket updates
	ballEvents := []models.BallEvent{
		{BallType: models.BallTypeGood, Runs: 1, IsWicket: false, BatsmanID: players[0].ID, BowlerID: players[2].ID},
		{BallType: models.BallTypeGood, Runs: 4, IsWicket: false, BatsmanID: players[1].ID, BowlerID: players[2].ID},
		{BallType: models.BallTypeWide, Runs: 1, IsWicket: false, BatsmanID: players[1].ID, BowlerID: players[2].ID},
		{BallType: models.BallTypeGood, Runs: 6, IsWicket: false, BatsmanID: players[1].ID, BowlerID: players[2].ID},
	}

	for i, ballEvent := range ballEvents {
		fmt.Printf("Adding ball %d: %s, %d runs\n", i+1, ballEvent.BallType, ballEvent.Runs)

		scoreboard, err := scoreboardService.AddBall(ctx, match.ID, &ballEvent)
		if err != nil {
			log.Printf("Failed to add ball %d: %v", i+1, err)
		} else {
			fmt.Printf("Ball %d added: %d/%d in %.1f overs\n",
				i+1, scoreboard.Score, scoreboard.Wickets, scoreboard.Overs)
		}

		// Wait a bit to see WebSocket messages
		time.Sleep(1 * time.Second)
	}

	// Test WebSocket stats endpoint
	fmt.Println("\nTesting WebSocket stats...")

	statsURL := fmt.Sprintf("http://localhost:8080/api/v1/ws/stats/%s", match.ID)
	resp, err := http.Get(statsURL)
	if err != nil {
		log.Printf("Failed to get WebSocket stats: %v", err)
	} else {
		defer resp.Body.Close()
		var stats map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
			log.Printf("Failed to decode stats response: %v", err)
		} else {
			fmt.Printf("WebSocket stats: %+v\n", stats)
		}
	}

	fmt.Println("\nWebSocket test completed successfully!")
}
