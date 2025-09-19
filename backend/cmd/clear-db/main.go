package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"spark-park-cricket-backend/internal/config"

	supabaseclient "github.com/supabase-community/supabase-go"
)

func main() {
	log.Println("=== SPARK PARK CRICKET - DATABASE CLEAR SCRIPT ===")
	log.Println("This script will clear ALL data from the database tables")
	log.Println("==================================================")

	// Load configuration
	cfg := config.Load()
	if cfg.SupabaseURL == "" || cfg.SupabaseAPIKey == "" {
		log.Fatal("ERROR: Supabase URL and API key are required. Please check your environment variables.")
	}

	// Create Supabase client
	clientOptions := &supabaseclient.ClientOptions{
		Schema: cfg.DatabaseSchema,
	}
	client, err := supabaseclient.NewClient(cfg.SupabaseURL, cfg.SupabaseAPIKey, clientOptions)
	if err != nil {
		log.Fatalf("ERROR: Failed to create Supabase client: %v", err)
	}

	log.Printf("✅ Connected to Supabase database")
	log.Printf("Database Schema: %s", cfg.DatabaseSchema)

	// Display schema information prominently
	log.Println("\n" + strings.Repeat("=", 60))
	log.Printf("🗄️  CLEARING DATA FROM SCHEMA: %s", strings.ToUpper(cfg.DatabaseSchema))
	log.Println(strings.Repeat("=", 60))

	// Show current table counts
	showTableCounts(client, cfg.DatabaseSchema)

	// Ask for confirmation
	if !confirmClear(cfg.DatabaseSchema) {
		log.Println("Operation cancelled by user")
		return
	}

	// Clear all tables
	if err := clearAllTables(client, cfg.DatabaseSchema); err != nil {
		log.Fatalf("ERROR: Failed to clear tables: %v", err)
	}

	// Show final table counts
	log.Println("\n=== CLEARING COMPLETED ===")
	log.Printf("🗄️  Schema '%s' has been cleared", cfg.DatabaseSchema)
	showTableCounts(client, cfg.DatabaseSchema)
	log.Println("✅ All tables have been cleared successfully!")
}

// showTableCounts displays the current count of records in each table
func showTableCounts(client *supabaseclient.Client, schema string) {
	log.Printf("\n=== TABLE COUNTS IN SCHEMA '%s' ===", schema)

	tables := []string{"balls", "overs", "innings", "live_scoreboard", "matches", "series", "schema_version"}

	for _, table := range tables {
		count, err := getTableCount(client, schema, table)
		if err != nil {
			log.Printf("❌ %s.%s: Error getting count - %v", schema, table, err)
		} else if count == -1 {
			log.Printf("⚠️  %s.%s: Table does not exist", schema, table)
		} else {
			log.Printf("📊 %s.%s: %d records", schema, table, count)
		}
	}
	log.Println("=============================")
}

// getTableCount returns the number of records in a table
func getTableCount(client *supabaseclient.Client, schema, table string) (int, error) {
	// Try to get a sample of records to estimate count
	// This is a workaround since Supabase client doesn't easily support COUNT queries
	var result []map[string]interface{}
	_, err := client.From(table).Select("id", "exact", false).Limit(1000, "").ExecuteTo(&result)
	if err != nil {
		// Check if the error is because the table doesn't exist
		if strings.Contains(err.Error(), "Could not find the table") {
			return -1, nil // Return -1 to indicate table doesn't exist
		}
		// Table might be empty or other error
		return 0, nil
	}

	// Count the returned records (limited to 1000)
	return len(result), nil
}

// confirmClear asks the user for confirmation before clearing the database
func confirmClear(schema string) bool {
	fmt.Printf("\n⚠️  WARNING: This will permanently delete ALL data from the following tables in schema '%s':\n", schema)
	fmt.Print("   - balls\n")
	fmt.Print("   - overs\n")
	fmt.Print("   - innings\n")
	fmt.Print("   - live_scoreboard\n")
	fmt.Print("   - matches\n")
	fmt.Print("   - series\n")
	fmt.Print("   - schema_version\n\n")

	fmt.Print("Are you sure you want to continue? Type 'YES' to confirm: ")

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Error reading input: %v", err)
		return false
	}

	response = strings.TrimSpace(response)
	return response == "YES"
}

// clearAllTables clears all tables in the correct order to respect foreign key constraints
func clearAllTables(client *supabaseclient.Client, schema string) error {
	log.Println("\n=== STARTING TABLE CLEARING ===")

	// Define tables in order (child tables first due to foreign key constraints)
	tables := []string{
		"balls",           // References overs
		"overs",           // References innings
		"innings",         // References matches
		"live_scoreboard", // References matches
		"matches",         // References series
		"series",          // No dependencies
		"schema_version",  // No dependencies
	}

	for _, table := range tables {
		log.Printf("🗑️  Clearing table: %s", table)

		if err := clearTable(client, schema, table); err != nil {
			return fmt.Errorf("failed to clear table %s: %w", table, err)
		}

		log.Printf("✅ Cleared table: %s", table)
		time.Sleep(100 * time.Millisecond) // Small delay between operations
	}

	log.Println("✅ All tables cleared successfully")
	return nil
}

// clearTable clears all records from a specific table
func clearTable(client *supabaseclient.Client, schema, table string) error {
	// Supabase requires a WHERE clause for DELETE operations
	// Try multiple approaches to ensure we delete all records

	// First try: created_at >= 1900-01-01
	_, err := client.From(table).Delete("", "").Gte("created_at", "1900-01-01").ExecuteTo(nil)
	if err != nil {
		// Check if the error is because the table doesn't exist
		if strings.Contains(err.Error(), "Could not find the table") {
			log.Printf("⚠️  Table %s does not exist, skipping...", table)
			return nil
		}
		log.Printf("⚠️  First delete attempt failed for %s: %v", table, err)
	}

	// Second try: id is not null (should match all records with UUIDs)
	_, err2 := client.From(table).Delete("", "").Not("id", "is", "null").ExecuteTo(nil)
	if err2 != nil {
		log.Printf("⚠️  Second delete attempt failed for %s: %v", table, err2)
	}

	// Third try: created_at is not null
	_, err3 := client.From(table).Delete("", "").Not("created_at", "is", "null").ExecuteTo(nil)
	if err3 != nil {
		log.Printf("⚠️  Third delete attempt failed for %s: %v", table, err3)
	}

	// If all attempts failed, return the first error
	if err != nil && err2 != nil && err3 != nil {
		return fmt.Errorf("failed to delete records from table %s: %w", table, err)
	}

	return nil
}
