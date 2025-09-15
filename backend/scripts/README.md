# Spark Park Cricket - Scripts

This folder contains utility scripts for the Spark Park Cricket backend system.

## Available Scripts

### 🚀 **reset_and_migrate.go** (Main Script)
**Purpose**: Complete database reset and migration
**Usage**: `go run scripts/reset_and_migrate.go`

**What it does**:
1. Deletes all data from existing tables (series, matches, teams, players, etc.)
2. Provides SQL migration script for manual execution
3. Verifies table status after operations
4. Shows step-by-step instructions

**When to use**: 
- When you want to start fresh with empty tables
- After making schema changes
- For clean testing environment

### 📄 **reset_supabase.sql**
**Purpose**: SQL script for manual database reset
**Usage**: Copy and paste into Supabase SQL Editor

**What it does**:
- Drops all existing tables
- Recreates tables with latest schema
- Creates indexes and constraints
- Adds documentation comments

### 🧪 **setup_test_db.sql**
**Purpose**: Setup test database schema
**Usage**: Copy and paste into Supabase SQL Editor

**What it does**:
- Creates test database schema
- Sets up test tables and constraints
- Prepares isolated test environment

### 🧹 **cleanup_test_db.sql**
**Purpose**: Clean up test database
**Usage**: Copy and paste into Supabase SQL Editor

**What it does**:
- Removes test data and tables
- Cleans up test database
- Resets test environment

## Quick Start

### For Fresh Database Setup:
```bash
# 1. Reset and migrate database
go run scripts/reset_and_migrate.go

# 2. Follow the instructions to execute SQL in Supabase dashboard

# 3. Start the server
go run cmd/server/main.go

# 4. Test the APIs
curl http://localhost:8080/api/v1/series
curl http://localhost:8080/api/v1/matches
```

### For Testing:
```bash
# Setup test database
# Copy setup_test_db.sql content to Supabase SQL Editor and run

# Clean up test database after testing
# Copy cleanup_test_db.sql content to Supabase SQL Editor and run
```

## Environment Requirements

All scripts require the following environment variables in `.env`:
- `SUPABASE_URL`: Your Supabase project URL
- `SUPABASE_API_KEY`: Your Supabase API key

## Notes

- The main reset script (`reset_and_migrate.go`) handles both data deletion and provides migration instructions
- Manual SQL execution is required due to Supabase API limitations for DDL operations
- All scripts include comprehensive error handling and status reporting
- The simplified system uses Team A vs Team B instead of complex team management
- Migration scripts have been cleaned up after successful prod_v1 schema migration
