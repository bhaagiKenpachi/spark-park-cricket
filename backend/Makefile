# Spark Park Cricket Backend - Makefile
# This Makefile provides commands for running tests and managing the project

.PHONY: help test test-unit test-integration test-e2e test-illegal test-all build run clean setup-test-db

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
	@echo "  test-all         Run all tests using Go test runner"
	@echo ""
	@echo "Build Commands:"
	@echo "  build            Build the application"
	@echo "  run              Run the server in background"
	@echo "  clean            Clean build artifacts"
	@echo ""
	@echo "Setup Commands:"
	@echo "  setup-test-db    Instructions for setting up test database"
	@echo ""
	@echo "Examples:"
	@echo "  make test-unit"
	@echo "  make test-all"
	@echo "  make build"

# Unit tests (no database required)
test-unit:
	@echo "🧪 Running Unit Tests..."
	@echo "------------------------"
	go test ./internal/tests/match_completion_unit_test.go -v

# Integration tests (requires database)
test-integration:
	@echo "🔗 Running Integration Tests..."
	@echo "-------------------------------"
	@echo "⚠️  Note: Integration tests require a running Supabase instance"
	go test ./internal/tests/match_completion_integration_test.go -v

# End-to-end tests (requires database)
test-e2e:
	@echo "🌐 Running End-to-End Tests..."
	@echo "------------------------------"
	@echo "⚠️  Note: E2E tests require a running Supabase instance"
	go test ./internal/tests/match_completion_e2e_test.go -v

# Illegal balls tests (requires database)
test-illegal:
	@echo "⚾ Running Illegal Balls Tests..."
	@echo "--------------------------------"
	go test ./internal/tests/illegal_balls_comprehensive_test.go -v

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
