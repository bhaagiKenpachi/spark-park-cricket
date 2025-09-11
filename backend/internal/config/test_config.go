package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// TestConfig holds test-specific configuration
type TestConfig struct {
	*Config
	TestSchema string
}

// LoadTestConfig loads configuration for testing with test schema
func LoadTestConfig() *TestConfig {
	// Load .env file if it exists (try multiple possible locations)
	envPaths := []string{".env", "../.env", "../../.env"}
	for _, path := range envPaths {
		if err := godotenv.Load(path); err == nil {
			break
		}
	}

	config := Load()

	// Use testing_db schema for tests
	testSchema := "testing_db"
	if envSchema := os.Getenv("TEST_SCHEMA"); envSchema != "" {
		testSchema = envSchema
	}

	return &TestConfig{
		Config:     config,
		TestSchema: testSchema,
	}
}

// GetSupabaseURL returns the Supabase URL with test schema
func (tc *TestConfig) GetSupabaseURL() string {
	baseURL := tc.SupabaseURL
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}
	return baseURL + "rest/v1/"
}

// GetSupabaseHeaders returns headers with test schema
func (tc *TestConfig) GetSupabaseHeaders() map[string]string {
	headers := map[string]string{
		"apikey":        tc.SupabaseAPIKey,
		"Authorization": "Bearer " + tc.SupabaseAPIKey,
		"Content-Type":  "application/json",
		"Accept":        "application/json",
	}

	// Add schema header for test database
	if tc.TestSchema != "" {
		headers["Accept-Profile"] = tc.TestSchema
		headers["Content-Profile"] = tc.TestSchema
	}

	return headers
}
