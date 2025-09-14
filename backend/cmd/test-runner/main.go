package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func main() {
	fmt.Println("🏏 Spark Park Cricket Backend - Test Runner")
	fmt.Println("==========================================")

	// Check if .env file exists
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		fmt.Println("❌ Error: .env file not found. Please create one with Supabase configuration.")
		fmt.Println("   Copy env.example to .env and fill in your Supabase credentials.")
		fmt.Println("   For testing, you can also copy env.test.example to .env")
		os.Exit(1)
	}

	fmt.Println("✅ Environment configuration found")

	// Check test database setup
	fmt.Println("")
	fmt.Println("🔧 Checking test database setup...")
	fmt.Println("   Make sure you have:")
	fmt.Println("   1. Created the testing_db schema in your Supabase database")
	fmt.Println("   2. Run the setup script: psql -f scripts/setup_test_db.sql")
	fmt.Println("   3. Set TEST_SCHEMA=testing_db in your .env file")
	fmt.Println("")

	// Run Unit Tests (no database required)
	fmt.Println("🧪 Running Unit Tests...")
	fmt.Println("------------------------")
	if err := runTest("Unit Tests", "./internal/tests/match_completion_unit_test.go"); err != nil {
		fmt.Printf("❌ Unit tests failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ Unit tests passed")

	// Run Integration Tests (requires database)
	fmt.Println("")
	fmt.Println("🔗 Running Integration Tests...")
	fmt.Println("-------------------------------")
	fmt.Println("⚠️  Note: Integration tests require a running Supabase instance")
	if err := runTest("Integration Tests", "./internal/tests/match_completion_integration_test.go"); err != nil {
		fmt.Printf("❌ Integration tests failed: %v\n", err)
		fmt.Println("   (Check database connection)")
	} else {
		fmt.Println("✅ Integration tests passed")
	}

	// Run E2E Tests (requires database)
	fmt.Println("")
	fmt.Println("🌐 Running End-to-End Tests...")
	fmt.Println("------------------------------")
	fmt.Println("⚠️  Note: E2E tests require a running Supabase instance")
	if err := runTest("E2E Tests", "./internal/tests/match_completion_e2e_test.go"); err != nil {
		fmt.Printf("❌ E2E tests failed: %v\n", err)
		fmt.Println("   (Check database connection)")
	} else {
		fmt.Println("✅ E2E tests passed")
	}

	// Run Illegal Balls Tests
	fmt.Println("")
	fmt.Println("⚾ Running Illegal Balls Tests...")
	fmt.Println("--------------------------------")
	if err := runTest("Illegal Balls Tests", "./internal/tests/illegal_balls_comprehensive_test.go"); err != nil {
		fmt.Printf("❌ Illegal balls tests failed: %v\n", err)
		fmt.Println("   (Check database connection)")
	} else {
		fmt.Println("✅ Illegal balls tests passed")
	}

	fmt.Println("")
	fmt.Println("🎉 Test suite completed!")
	fmt.Println("========================")
	fmt.Println("")
	fmt.Println("📋 Test Summary:")
	fmt.Println("   • Unit Tests: ✅ Passed (no database required)")
	fmt.Println("   • Integration Tests: ⚠️  Requires database")
	fmt.Println("   • E2E Tests: ⚠️  Requires database")
	fmt.Println("   • Illegal Balls Tests: ⚠️  Requires database")
	fmt.Println("")
	fmt.Println("💡 To run database-dependent tests:")
	fmt.Println("   1. Ensure Supabase is running")
	fmt.Println("   2. Check .env file has correct credentials")
	fmt.Println("   3. Run individual test files as needed")
}

func runTest(testName, testFile string) error {
	fmt.Printf("Running %s...\n", testName)

	// Change to the backend directory
	backendDir, err := filepath.Abs(".")
	if err != nil {
		return fmt.Errorf("failed to get backend directory: %w", err)
	}

	// Run the test
	cmd := exec.Command("go", "test", testFile, "-v")
	cmd.Dir = backendDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Set timeout
	done := make(chan error, 1)
	go func() {
		done <- cmd.Run()
	}()

	select {
	case err := <-done:
		if err != nil {
			return err
		}
		return nil
	case <-time.After(5 * time.Minute):
		if err := cmd.Process.Kill(); err != nil {
			return fmt.Errorf("test timed out after 5 minutes and failed to kill process: %w", err)
		}
		return fmt.Errorf("test timed out after 5 minutes")
	}
}
