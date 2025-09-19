# Spark Park Cricket Backend - Makefile
# This Makefile provides commands for running tests and managing the project

.PHONY: help test test-unit test-integration test-e2e test-illegal test-series test-match test-scorecard test-all build run clean setup-test-db clear-db

# Default target
help:
	@echo "🏏 Spark Park Cricket Backend - Available Commands"
	@echo "=================================================="
	@echo ""
	@echo "Test Commands:"
	@echo "  test-unit        Run unit tests (no database required)"
	@echo "  test-integration Run integration tests (requires database)"
	@echo "  test-e2e         Run end-to-end tests (requires database)"
	@echo "  test-illegal     Run illegal balls tests (requires database)"
	@echo "  test-series      Run series-specific tests (unit + integration)"
	@echo "  test-match       Run match-specific tests (unit + integration + e2e)"
	@echo "  test-scorecard   Run scorecard-specific tests (unit + integration + e2e)"
	@echo "  test-all         Run all tests using Go test runner"
	@echo ""
	@echo "Build Commands:"
	@echo "  build            Build the application"
	@echo "  run              Run the server in background"
	@echo "  clean            Clean build artifacts"
	@echo ""
	@echo "Setup Commands:"
	@echo "  setup-test-db    Instructions for setting up test database"
	@echo "  clear-db         Clear all data from database tables (DANGEROUS!)"
	@echo ""
	@echo "Examples:"
	@echo "  make test-unit"
	@echo "  make test-all"
	@echo "  make build"

# Unit tests (no database required)
test-unit:
	@echo "🧪 Running Unit Tests..."
	@echo "------------------------"
	go test ./tests/unit/ -v

# Integration tests (requires database)
test-integration:
	@echo "🔗 Running Integration Tests..."
	@echo "-------------------------------"
	@echo "⚠️  Note: Integration tests require a running Supabase instance"
	@echo "🧹 Cleaning database before integration tests..."
	@echo "YES" | go run cmd/clear-db/main.go
	@echo "✅ Database cleaned, starting integration tests..."
	go test ./tests/integration/ -v

# End-to-end tests (requires database)
test-e2e:
	@echo "🌐 Running End-to-End Tests..."
	@echo "------------------------------"
	@echo "⚠️  Note: E2E tests require a running Supabase instance"
	@echo "🧹 Cleaning database before E2E tests..."
	@echo "YES" | go run cmd/clear-db/main.go
	@echo "✅ Database cleaned, starting E2E tests..."
	go test ./tests/e2e/ -v

# Illegal balls tests (requires database)
test-illegal:
	@echo "⚾ Running Illegal Balls Tests..."
	@echo "--------------------------------"
	@echo "🧹 Cleaning database before illegal balls tests..."
	@echo "YES" | go run cmd/clear-db/main.go
	@echo "✅ Database cleaned, starting illegal balls tests..."
	go test ./tests/integration/illegal_balls_comprehensive_test.go -v

# Series-specific tests
test-series:
	@echo "🏏 Running Series Tests..."
	@echo "========================="
	@echo "Unit tests (no database required):"
	go test ./tests/unit/match_completion_unit_test.go ./tests/unit/series_service_test.go ./tests/unit/series_handler_test.go -v
	@echo ""
	@echo "Integration tests (requires database):"
	@echo "⚠️  Note: Integration tests require a running Supabase instance"
	@echo "🧹 Cleaning database before series integration tests..."
	@echo "YES" | go run cmd/clear-db/main.go
	@echo "✅ Database cleaned, starting series integration tests..."
	go test ./tests/integration/series_integration_test.go -v

# Match-specific tests
test-match:
	@echo "🏏 Running Match Tests..."
	@echo "========================="
	@echo "Unit tests (no database required):"
	go test ./tests/unit/match_completion_unit_test.go ./tests/unit/match_service_test.go ./tests/unit/match_handler_test.go -v
	@echo ""
	@echo "Integration tests (requires database):"
	@echo "⚠️  Note: Integration tests require a running Supabase instance"
	@echo "🧹 Cleaning database before match integration tests..."
	@echo "YES" | go run cmd/clear-db/main.go
	@echo "✅ Database cleaned, starting match integration tests..."
	go test ./tests/integration/match_integration_test.go -v
	@echo ""
	@echo "E2E tests (requires database):"
	@echo "⚠️  Note: E2E tests require a running Supabase instance"
	@echo "🧹 Cleaning database before match E2E tests..."
	@echo "YES" | go run cmd/clear-db/main.go
	@echo "✅ Database cleaned, starting match E2E tests..."
	go test ./tests/e2e/match_workflow_e2e_test.go -v

# Scorecard-specific tests
test-scorecard:
	@echo "🏏 Running Scorecard Tests..."
	@echo "============================="
	@echo "Unit tests (no database required):"
	go test ./tests/unit/match_completion_unit_test.go ./tests/unit/scorecard_service_test.go ./tests/unit/scorecard_handler_test.go -v
	@echo ""
	@echo "Integration tests (requires database):"
	@echo "⚠️  Note: Integration tests require a running Supabase instance"
	@echo "🧹 Cleaning database before scorecard integration tests..."
	@echo "YES" | go run cmd/clear-db/main.go
	@echo "✅ Database cleaned, starting scorecard integration tests..."
	go test ./tests/integration/scorecard_integration_test.go -v
	@echo ""
	@echo "E2E tests (requires database):"
	@echo "⚠️  Note: E2E tests require a running Supabase instance"
	@echo "🧹 Cleaning database before scorecard E2E tests..."
	@echo "YES" | go run cmd/clear-db/main.go
	@echo "✅ Database cleaned, starting scorecard E2E tests..."
	go test ./tests/e2e/scorecard_workflow_e2e_test.go -v

# Run all tests using the Go test runner
test-all:
	@echo "🏏 Running All Tests..."
	@echo "======================"
	go run cmd/test-runner/main.go

# Build the application
build:
	@echo "🔨 Building application..."
	go build -o bin/server cmd/server/main.go
	@echo "✅ Build completed: bin/server"

# Run the server in background
run:
	@echo "🚀 Starting server in background..."
	go run cmd/server/main.go &
	@echo "✅ Server started in background"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -rf bin/
	go clean
	@echo "✅ Clean completed"

# Setup test database instructions
setup-test-db:
	@echo "🗄️  Test Database Setup Instructions"
	@echo "===================================="
	@echo ""
	@echo "1. Connect to your Supabase database:"
	@echo "   psql -h your-supabase-host -U postgres -d postgres"
	@echo ""
	@echo "2. Run the test schema setup script:"
	@echo "   \\i scripts/setup_test_db.sql"
	@echo ""
	@echo "   Or from command line:"
	@echo "   psql -h your-supabase-host -U postgres -d postgres -f scripts/setup_test_db.sql"
	@echo ""
	@echo "3. Verify the schema was created:"
	@echo "   \\dn"
	@echo "   \\dt testing_db.*"
	@echo ""
	@echo "4. Make sure your .env file has:"
	@echo "   TEST_SCHEMA=testing_db"
	@echo "   SUPABASE_URL=your_supabase_url"
	@echo "   SUPABASE_API_KEY=your_supabase_key"

# Clear all database tables (DANGEROUS!)
clear-db:
	@echo "🗑️  Database Clear Script"
	@echo "========================"
	@echo "⚠️  WARNING: This will permanently delete ALL data from the database!"
	@echo "⚠️  Make sure you have backups before proceeding!"
	@echo ""
	@echo "Running database clear script..."
	go run cmd/clear-db/main.go
