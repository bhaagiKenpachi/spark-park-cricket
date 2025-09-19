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
	// Google OAuth Configuration
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	SessionSecret      string
	SessionMaxAge      int
}

func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	cfg := &Config{
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
		// Google OAuth Configuration
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8081/api/v1/auth/google/callback"),
		SessionSecret:      getEnv("SESSION_SECRET", "your-super-secret-session-key-change-this-in-production"),
		SessionMaxAge:      getEnvInt("SESSION_MAX_AGE", 86400), // 24 hours
	}

	// Log database configuration
	logDatabaseConfig(cfg)

	return cfg
}

// logDatabaseConfig logs the database configuration being used
func logDatabaseConfig(cfg *Config) {
	log.Println("=== DATABASE CONFIGURATION ===")

	// Log database type and URL
	if cfg.SupabaseURL != "" {
		// Extract host from URL for logging (hide sensitive info)
		log.Printf("Database Type: Supabase (PostgreSQL)")
		log.Printf("Database URL: %s", maskURL(cfg.SupabaseURL))
		log.Printf("Database Schema: %s", cfg.DatabaseSchema)

		if cfg.SupabaseAPIKey != "" {
			log.Printf("Database API Key: %s", maskAPIKey(cfg.SupabaseAPIKey))
		} else {
			log.Printf("Database API Key: NOT SET")
		}
	} else {
		log.Printf("Database Type: NOT CONFIGURED")
		log.Printf("Database URL: NOT SET")
	}

	// Log cache configuration
	if cfg.CacheEnabled {
		log.Printf("Cache: Enabled (Redis)")
		log.Printf("Cache URL: %s", cfg.RedisURL)
		log.Printf("Cache DB: %d", cfg.RedisDB)
		if cfg.RedisPassword != "" {
			log.Printf("Cache Password: %s", maskPassword(cfg.RedisPassword))
		} else {
			log.Printf("Cache Password: NOT SET")
		}
		log.Printf("Cache TLS: %t", cfg.RedisUseTLS)
	} else {
		log.Printf("Cache: Disabled")
	}

	log.Println("===============================")
}

// maskURL masks sensitive parts of the URL for logging
func maskURL(url string) string {
	if url == "" {
		return "NOT SET"
	}
	// For Supabase URLs, show the host but mask any sensitive parts
	// Example: https://xyz.supabase.co -> https://***.supabase.co
	if len(url) > 20 {
		return url[:10] + "***" + url[len(url)-15:]
	}
	return url
}

// maskAPIKey masks the API key for logging
func maskAPIKey(apiKey string) string {
	if apiKey == "" {
		return "NOT SET"
	}
	if len(apiKey) > 8 {
		return apiKey[:4] + "***" + apiKey[len(apiKey)-4:]
	}
	return "***"
}

// maskPassword masks the password for logging
func maskPassword(password string) string {
	if password == "" {
		return "NOT SET"
	}
	return "***"
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
