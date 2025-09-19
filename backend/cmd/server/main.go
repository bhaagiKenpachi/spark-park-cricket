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
	log.Println("=== SPARK PARK CRICKET BACKEND STARTING ===")

	// Initialize logger
	utils.InitLogger()
	log.Printf("✅ Logger initialized")

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		utils.LogWarn("No .env file found, using system environment variables", nil)
	} else {
		log.Printf("✅ Environment variables loaded from .env file")
	}

	// Load configuration
	log.Printf("Loading configuration...")
	cfg := config.Load()
	log.Printf("✅ Configuration loaded successfully")

	// Initialize Supabase client
	log.Printf("Initializing database connection...")
	dbClient, err := database.NewClient(cfg)
	if err != nil {
		log.Fatalf("❌ Failed to initialize database client: %v", err)
	}

	// Test database connection
	log.Printf("Testing database connection...")
	if err := dbClient.HealthCheck(); err != nil {
		log.Printf("❌ Database health check failed: %v", err)
		log.Printf("⚠️  Server will continue but database operations may fail")
	} else {
		log.Printf("✅ Database connection verified successfully")
	}

	// Setup routes
	log.Printf("Setting up API routes...")
	router := handlers.SetupRoutes(dbClient, cfg)
	log.Printf("✅ API routes configured")

	// Log startup information
	log.Println("=== SERVER STARTUP COMPLETE ===")
	log.Printf("🚀 Server starting on port: %s", cfg.Port)
	log.Printf("📊 Database: Supabase (PostgreSQL) - Schema: %s", cfg.DatabaseSchema)
	if dbClient.CacheManager != nil {
		log.Printf("⚡ Cache: Enabled (Redis)")
	} else {
		log.Printf("⚡ Cache: Disabled")
	}
	log.Printf("🌐 API Endpoints:")
	log.Printf("   - REST API: http://localhost:%s/api/v1/", cfg.Port)
	log.Printf("   - GraphQL: http://localhost:%s/api/v1/graphql", cfg.Port)
	log.Printf("   - GraphQL Playground: http://localhost:%s/api/v1/graphql/playground", cfg.Port)
	log.Printf("   - Health Check: http://localhost:%s/health", cfg.Port)
	log.Printf("   - WebSocket: ws://localhost:%s/api/v1/ws/match/{match_id}", cfg.Port)
	log.Println("===============================================")

	fmt.Printf("🚀 Spark Park Cricket Backend is running on :%s\n", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}
