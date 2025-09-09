package main

import (
	"fmt"
	"log"
	"net/http"
	"spark-park-cricket-backend/internal/config"
	"spark-park-cricket-backend/internal/database"
	"spark-park-cricket-backend/internal/handlers"
	"spark-park-cricket-backend/internal/utils"

	"github.com/joho/godotenv"
)

func main() {
	// Initialize logger
	utils.InitLogger()

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		utils.LogWarn("No .env file found, using system environment variables", nil)
	}

	// Load configuration
	cfg := config.Load()

	// Initialize Supabase client
	dbClient, err := database.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database client: %v", err)
	}

	// Test database connection
	if err := dbClient.HealthCheck(); err != nil {
		log.Printf("Database health check failed: %v", err)
	} else {
		log.Println("Database connection successful")
	}

	// Setup routes
	router := handlers.SetupRoutes(dbClient)

	fmt.Printf("Server starting on :%s\n", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}
