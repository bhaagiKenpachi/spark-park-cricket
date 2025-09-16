package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	SupabaseURL            string
	SupabaseAPIKey         string
	SupabasePublishableKey string
	SupabaseSecretKey      string
	Port                   string
	DatabaseSchema         string
	RedisURL               string
	RedisPassword          string
	RedisDB                int
	RedisUseTLS            bool
	CacheEnabled           bool
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
		RedisURL:               getEnv("REDIS_URL", "localhost:6379"),
		RedisPassword:          getEnv("REDIS_PASSWORD", ""),
		RedisDB:                getEnvInt("REDIS_DB", 0),
		RedisUseTLS:            getEnvBool("REDIS_USE_TLS", false),
		CacheEnabled:           getEnvBool("CACHE_ENABLED", true),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
