package config

import (
	"os"
)

type Config struct {
	SupabaseURL    string
	SupabaseAPIKey string
	Port           string
}

func Load() *Config {
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
