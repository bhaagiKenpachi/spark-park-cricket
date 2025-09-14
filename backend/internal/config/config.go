package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	SupabaseURL            string
	SupabaseAPIKey         string
	SupabasePublishableKey string
	SupabaseSecretKey      string
	Port                   string
	DatabaseSchema         string
}

func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	return &Config{
		SupabaseURL:            getEnv("SUPABASE_URL", ""),
		SupabaseAPIKey:         getEnv("SUPABASE_API_KEY", ""),
		SupabasePublishableKey: getEnv("SUPABASE_PUBLISHABLE_KEY", ""),
		SupabaseSecretKey:      getEnv("SUPABASE_SECRET_KEY", ""),
		Port:                   getEnv("PORT", "8081"),
		DatabaseSchema:         getEnv("DATABASE_SCHEMA", "prod_v1"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
