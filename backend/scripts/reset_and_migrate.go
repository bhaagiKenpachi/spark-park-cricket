package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type TableData struct {
	Data []map[string]interface{} `json:"data"`
}

func main() {
	fmt.Println("🚀 Spark Park Cricket - Complete Database Reset & Migration")
	fmt.Println("==========================================================")

	// Load environment variables from .env file
	envFile, err := os.Open(".env")
	if err != nil {
		log.Fatalf("Failed to open .env file: %v", err)
	}
	defer envFile.Close()

	scanner := bufio.NewScanner(envFile)
	envVars := make(map[string]string)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			envVars[parts[0]] = parts[1]
		}
	}

	supabaseURL := envVars["SUPABASE_URL"]
	supabaseAPIKey := envVars["SUPABASE_API_KEY"]

	if supabaseURL == "" || supabaseAPIKey == "" {
		log.Fatal("SUPABASE_URL and SUPABASE_API_KEY not found in .env file")
	}

	fmt.Println("✅ Environment variables loaded")
	fmt.Printf("📍 Supabase URL: %s\n", supabaseURL)

	// Create HTTP client
	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	// Step 1: Delete all data from existing tables
	fmt.Println("\n🗑️  Step 1: Deleting all existing data...")
	fmt.Println("=========================================")

	// All possible tables that might exist
	allTables := []string{
		"balls", "overs", "live_scoreboard", "matches", "series",
		"players", "teams", // Old tables that should be removed
	}

	for _, table := range allTables {
		fmt.Printf("   Deleting from table: %s\n", table)

		// Create DELETE request with WHERE clause
		deleteURL := fmt.Sprintf("%s/rest/v1/%s?id=not.is.null", supabaseURL, table)
		req, err := http.NewRequest("DELETE", deleteURL, nil)
		if err != nil {
			fmt.Printf("   ⚠️  Failed to create DELETE request: %v\n", err)
			continue
		}

		req.Header.Set("Authorization", "Bearer "+supabaseAPIKey)
		req.Header.Set("apikey", supabaseAPIKey)
		req.Header.Set("Prefer", "return=minimal")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("   ⚠️  Failed to delete from %s: %v\n", table, err)
			continue
		}

		_, _ = io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			fmt.Printf("   ✅ Successfully deleted from %s\n", table)
		} else if resp.StatusCode == 404 {
			fmt.Printf("   ℹ️  Table %s does not exist (OK)\n", table)
		} else {
			fmt.Printf("   ⚠️  Delete from %s returned status %d\n", table, resp.StatusCode)
		}
	}

	// Step 2: Drop and recreate tables using SQL
	fmt.Println("\n🔄 Step 2: Applying database migration...")
	fmt.Println("=======================================")

	// Read the migration script
	migrationPath := "scripts/reset_supabase.sql"
	migrationSQL, err := os.ReadFile(migrationPath)
	if err != nil {
		log.Fatalf("Failed to read migration file: %v", err)
	}

	fmt.Println("📄 Migration script loaded")
	fmt.Println("⚠️  IMPORTANT: Manual SQL execution required")
	fmt.Println()
	fmt.Println("Due to Supabase API limitations, you need to manually execute the migration.")
	fmt.Println("Here's what you need to do:")
	fmt.Println()
	fmt.Println("1. Go to your Supabase Dashboard:")
	fmt.Printf("   %s\n", supabaseURL)
	fmt.Println()
	fmt.Println("2. Navigate to SQL Editor")
	fmt.Println("3. Copy and paste the following SQL script:")
	fmt.Println()
	fmt.Println("```sql")
	fmt.Println(string(migrationSQL))
	fmt.Println("```")
	fmt.Println()
	fmt.Println("4. Click 'Run' to execute the script")
	fmt.Println()

	// Step 3: Verify the migration
	fmt.Println("🔍 Step 3: Verification instructions...")
	fmt.Println("======================================")
	fmt.Println()
	fmt.Println("After running the SQL script, verify the migration by:")
	fmt.Println()
	fmt.Println("1. Check that these tables exist and are empty:")
	fmt.Println("   - series")
	fmt.Println("   - matches")
	fmt.Println("   - live_scoreboard")
	fmt.Println("   - overs")
	fmt.Println("   - balls")
	fmt.Println()
	fmt.Println("2. Test the APIs:")
	fmt.Println("   curl http://localhost:8080/api/v1/series")
	fmt.Println("   curl http://localhost:8080/api/v1/matches")
	fmt.Println()
	fmt.Println("3. Both should return: {\"data\":[]}")
	fmt.Println()

	// Step 4: Show current table status
	fmt.Println("📊 Step 4: Current table status...")
	fmt.Println("=================================")

	for _, table := range []string{"series", "matches", "teams", "players", "live_scoreboard", "overs", "balls"} {
		fmt.Printf("🔍 Checking table: %s\n", table)

		url := fmt.Sprintf("%s/rest/v1/%s?select=*&limit=1", supabaseURL, table)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Printf("   ❌ Failed to create request: %v\n", err)
			continue
		}

		req.Header.Set("Authorization", "Bearer "+supabaseAPIKey)
		req.Header.Set("apikey", supabaseAPIKey)
		req.Header.Set("Prefer", "count=exact")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("   ❌ Failed to query table: %v\n", err)
			continue
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode == 200 {
			var tableData TableData
			err := json.Unmarshal(body, &tableData)
			if err != nil {
				// Handle direct array response
				var directData []map[string]interface{}
				err = json.Unmarshal(body, &directData)
				if err != nil {
					fmt.Printf("   ⚠️  Response format error: %v\n", err)
					continue
				}
				tableData.Data = directData
			}

			count := len(tableData.Data)
			if count == 0 {
				fmt.Printf("   ✅ Table %s exists and is EMPTY\n", table)
			} else {
				fmt.Printf("   📊 Table %s has %d records\n", table, count)
			}
		} else if resp.StatusCode == 404 {
			fmt.Printf("   ❌ Table %s does NOT exist\n", table)
		} else {
			fmt.Printf("   ⚠️  Table %s returned status %d\n", table, resp.StatusCode)
		}
	}

	fmt.Println("\n🎉 Database reset and migration process completed!")
	fmt.Println("💡 Next steps:")
	fmt.Println("   1. Execute the SQL script in Supabase dashboard")
	fmt.Println("   2. Start the server: go run cmd/server/main.go")
	fmt.Println("   3. Test the APIs to verify empty responses")
}
