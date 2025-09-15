package main

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// Migration represents a database migration
type Migration struct {
	Version string
	Name    string
	Path    string
}

// MigrationRunner handles database migrations
type MigrationRunner struct {
	db *sql.DB
}

// NewMigrationRunner creates a new migration runner
func NewMigrationRunner(db *sql.DB) *MigrationRunner {
	return &MigrationRunner{db: db}
}

// CreateMigrationsTable creates the migrations tracking table
func (mr *MigrationRunner) CreateMigrationsTable(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version VARCHAR(255) PRIMARY KEY,
		applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);`

	_, err := mr.db.ExecContext(ctx, query)
	return err
}

// GetAppliedMigrations returns list of applied migrations
func (mr *MigrationRunner) GetAppliedMigrations(ctx context.Context) (map[string]bool, error) {
	query := `SELECT version FROM schema_migrations ORDER BY version;`
	rows, err := mr.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}
	return applied, nil
}

// GetMigrationFiles returns sorted list of migration files
func (mr *MigrationRunner) GetMigrationFiles(migrationsDir string) ([]Migration, error) {
	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		return nil, err
	}

	var migrations []Migration
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			// Extract version from filename (e.g., "001_initial_schema.sql" -> "001")
			parts := strings.Split(file.Name(), "_")
			if len(parts) > 0 {
				version := parts[0]
				name := strings.TrimSuffix(file.Name(), ".sql")
				migrations = append(migrations, Migration{
					Version: version,
					Name:    name,
					Path:    filepath.Join(migrationsDir, file.Name()),
				})
			}
		}
	}

	// Sort migrations by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

// ReadMigrationFile reads the content of a migration file
func (mr *MigrationRunner) ReadMigrationFile(path string) (string, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// ApplyMigration applies a single migration
func (mr *MigrationRunner) ApplyMigration(ctx context.Context, migration Migration) error {
	fmt.Printf("Applying migration: %s (%s)\n", migration.Version, migration.Name)

	// Read migration file
	content, err := mr.ReadMigrationFile(migration.Path)
	if err != nil {
		return fmt.Errorf("failed to read migration file %s: %w", migration.Path, err)
	}

	// Start transaction
	tx, err := mr.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Execute migration SQL
	if _, err := tx.ExecContext(ctx, content); err != nil {
		return fmt.Errorf("failed to execute migration %s: %w", migration.Version, err)
	}

	// Record migration as applied
	recordQuery := `INSERT INTO schema_migrations (version) VALUES ($1) ON CONFLICT (version) DO NOTHING;`
	if _, err := tx.ExecContext(ctx, recordQuery, migration.Version); err != nil {
		return fmt.Errorf("failed to record migration %s: %w", migration.Version, err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit migration %s: %w", migration.Version, err)
	}

	fmt.Printf("✓ Migration %s applied successfully\n", migration.Version)
	return nil
}

// RunMigrations runs all pending migrations
func (mr *MigrationRunner) RunMigrations(ctx context.Context, migrationsDir string) error {
	fmt.Println("Starting database migrations...")

	// Create migrations table
	if err := mr.CreateMigrationsTable(ctx); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get applied migrations
	applied, err := mr.GetAppliedMigrations(ctx)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Get migration files
	migrations, err := mr.GetMigrationFiles(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to get migration files: %w", err)
	}

	if len(migrations) == 0 {
		fmt.Println("No migration files found")
		return nil
	}

	// Apply pending migrations
	pendingCount := 0
	for _, migration := range migrations {
		if !applied[migration.Version] {
			if err := mr.ApplyMigration(ctx, migration); err != nil {
				return fmt.Errorf("failed to apply migration %s: %w", migration.Version, err)
			}
			pendingCount++
		} else {
			fmt.Printf("⏭️  Migration %s already applied\n", migration.Version)
		}
	}

	if pendingCount == 0 {
		fmt.Println("✓ All migrations are up to date")
	} else {
		fmt.Printf("✓ Applied %d new migrations\n", pendingCount)
	}

	return nil
}

// GetConnectionString builds database connection string from environment variables
func GetConnectionString() (string, error) {
	// First try to use DATABASE_URL if it exists
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		return databaseURL, nil
	}

	// Fallback to building from individual components
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseAPIKey := os.Getenv("SUPABASE_API_KEY")

	if supabaseURL == "" {
		return "", fmt.Errorf("SUPABASE_URL environment variable is required")
	}

	if supabaseAPIKey == "" {
		return "", fmt.Errorf("SUPABASE_API_KEY environment variable is required")
	}

	// Extract project reference from Supabase URL
	// URL format: https://qehkpqubnnpbaejhcwvx.supabase.co
	parts := strings.Split(supabaseURL, "//")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid SUPABASE_URL format")
	}

	hostParts := strings.Split(parts[1], ".")
	if len(hostParts) < 2 {
		return "", fmt.Errorf("invalid SUPABASE_URL format")
	}

	projectRef := hostParts[0]

	// Build PostgreSQL connection string
	// Format: postgresql://postgres:[password]@db.[project-ref].supabase.co:5432/postgres
	connStr := fmt.Sprintf("postgresql://postgres:%s@db.%s.supabase.co:5432/postgres?sslmode=require",
		os.Getenv("DATABASE_PASSWORD"), projectRef)

	return connStr, nil
}

func main() {
	fmt.Println("🏏 Spark Park Cricket - Database Migration Runner")
	fmt.Println("================================================")

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Get database connection string
	connStr, err := GetConnectionString()
	if err != nil {
		log.Fatalf("Failed to get connection string: %v", err)
	}

	// Connect to database
	fmt.Println("Connecting to database...")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	fmt.Println("✓ Database connection established")

	// Create migration runner
	runner := NewMigrationRunner(db)

	// Get migrations directory
	migrationsDir := "internal/database/migrations"
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		log.Fatalf("Migrations directory not found: %s", migrationsDir)
	}

	// Run migrations
	if err := runner.RunMigrations(ctx, migrationsDir); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	fmt.Println("🎉 Database migrations completed successfully!")
}
