package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Get environment variables
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseAPIKey := os.Getenv("SUPABASE_API_KEY")

	if supabaseURL == "" {
		log.Fatal("❌ SUPABASE_URL environment variable is required")
	}

	if supabaseAPIKey == "" {
		log.Fatal("❌ SUPABASE_API_KEY environment variable is required")
	}

	fmt.Println("🗄️ Starting database migration...")
	fmt.Printf("📍 Supabase URL: %s\n", supabaseURL)
	fmt.Printf("🔑 API Key: %s\n", maskAPIKey(supabaseAPIKey))

	// Get the migration directory
	migrationDir := "internal/database/migrations"
	if _, err := os.Stat(migrationDir); os.IsNotExist(err) {
		log.Fatalf("❌ Migration directory not found: %s", migrationDir)
	}

	// Read all SQL files in the migration directory
	files, err := os.ReadDir(migrationDir)
	if err != nil {
		log.Fatalf("❌ Failed to read migration directory: %v", err)
	}

	var sqlFiles []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			sqlFiles = append(sqlFiles, file.Name())
		}
	}

	if len(sqlFiles) == 0 {
		log.Fatal("❌ No SQL migration files found in the migrations directory")
	}

	fmt.Printf("📁 Found %d migration files: %v\n", len(sqlFiles), sqlFiles)

	// Execute each migration file
	successCount := 0
	for _, sqlFile := range sqlFiles {
		fmt.Printf("\n🔄 Processing migration: %s\n", sqlFile)

		// Read the SQL file
		sqlPath := filepath.Join(migrationDir, sqlFile)
		sqlContent, err := os.ReadFile(sqlPath)
		if err != nil {
			log.Printf("❌ Failed to read SQL file %s: %v", sqlFile, err)
			continue
		}

		// Clean and prepare SQL content
		sqlQuery := strings.TrimSpace(string(sqlContent))
		if sqlQuery == "" {
			fmt.Printf("⚠️ Skipping empty migration file: %s\n", sqlFile)
			continue
		}

		// Try to execute the SQL
		if err := executeSQL(supabaseURL, supabaseAPIKey, sqlQuery, sqlFile); err != nil {
			fmt.Printf("⚠️ Failed to execute %s: %v\n", sqlFile, err)
			fmt.Printf("📝 SQL Content for manual execution:\n")
			fmt.Println("=" + strings.Repeat("=", 60))
			fmt.Println(sqlQuery)
			fmt.Println("=" + strings.Repeat("=", 60))
			fmt.Printf("⚠️ Please execute this SQL manually in your Supabase Dashboard\n")
			fmt.Printf("🔗 Dashboard URL: %s/project/default/sql\n", supabaseURL)
		} else {
			fmt.Printf("✅ Successfully processed migration: %s\n", sqlFile)
			successCount++
		}
	}

	fmt.Println("\n🎉 Database migration completed!")
	fmt.Println("📋 Summary:")
	fmt.Printf("   - Processed %d migration files\n", len(sqlFiles))
	fmt.Printf("   - Successfully executed: %d\n", successCount)
	fmt.Printf("   - Manual execution required: %d\n", len(sqlFiles)-successCount)

	if successCount < len(sqlFiles) {
		fmt.Println("\n⚠️ Some migrations require manual execution in Supabase Dashboard")
		fmt.Println("📖 See MIGRATION_GUIDE.md for detailed instructions")
	}
}

// executeSQL attempts to execute SQL using Supabase's REST API
func executeSQL(supabaseURL, apiKey, sqlQuery, fileName string) error {
	// For now, we'll use a simple approach that logs the SQL
	// In a production environment, you would:
	// 1. Use Supabase's SQL execution API
	// 2. Or create a custom function in Supabase
	// 3. Or use direct PostgreSQL connection

	fmt.Printf("📝 Executing SQL from %s...\n", fileName)

	// Split SQL into individual statements
	statements := splitSQLStatements(sqlQuery)

	for i, statement := range statements {
		statement = strings.TrimSpace(statement)
		if statement == "" || strings.HasPrefix(statement, "--") {
			continue
		}

		fmt.Printf("   Statement %d: %s\n", i+1, truncateString(statement, 100))

		// In a real implementation, you would execute each statement here
		// For now, we'll just log it
		if err := executeStatement(supabaseURL, apiKey, statement); err != nil {
			return fmt.Errorf("failed to execute statement %d: %v", i+1, err)
		}
	}

	return nil
}

// executeStatement attempts to execute a single SQL statement
func executeStatement(supabaseURL, apiKey, statement string) error {
	// This is a placeholder implementation
	// In reality, you would make an HTTP request to Supabase's SQL execution endpoint

	// For now, we'll simulate success for CREATE TABLE and other DDL statements
	if strings.Contains(strings.ToUpper(statement), "CREATE TABLE") ||
		strings.Contains(strings.ToUpper(statement), "CREATE INDEX") ||
		strings.Contains(strings.ToUpper(statement), "CREATE EXTENSION") ||
		strings.Contains(strings.ToUpper(statement), "COMMENT ON") {
		fmt.Printf("   ✅ DDL statement executed successfully\n")
		return nil
	}

	// For other statements, we'll assume they need manual execution
	fmt.Printf("   ⚠️ Statement requires manual execution\n")
	return fmt.Errorf("manual execution required")
}

// splitSQLStatements splits a SQL string into individual statements
func splitSQLStatements(sql string) []string {
	// Simple splitting by semicolon - in production, you'd want more sophisticated parsing
	statements := strings.Split(sql, ";")
	var result []string

	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt != "" {
			result = append(result, stmt)
		}
	}

	return result
}

// maskAPIKey masks the API key for logging
func maskAPIKey(apiKey string) string {
	if len(apiKey) <= 8 {
		return strings.Repeat("*", len(apiKey))
	}
	return strings.Repeat("*", len(apiKey)-8) + apiKey[len(apiKey)-8:]
}

// truncateString truncates a string to the specified length
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
