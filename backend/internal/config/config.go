package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	SupabaseURL    string
	SupabaseAPIKey string
	Port           string
}

func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	return &Config{
		SupabaseURL:    getEnv("SUPABASE_URL", ""),
		SupabaseAPIKey: getEnv("SUPABASE_API_KEY", ""),
		Port:           getEnv("PORT", "8081"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
