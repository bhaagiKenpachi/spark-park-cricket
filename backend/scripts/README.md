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

### 🌱 **seed.go**
**Purpose**: Seed database with sample data (legacy)
**Usage**: `go run scripts/seed.go`

**Note**: This script is from the old system and may not be compatible with the new simplified schema.

### 🧪 **test_scoring.go**
**Purpose**: Test cricket scoring functionality
**Usage**: `go run scripts/test_scoring.go`

**What it does**:
- Tests ball-by-ball scoring
- Validates run types and ball types
- Tests over completion logic

### 🔌 **test_websocket.go**
**Purpose**: Test WebSocket functionality
**Usage**: `go run scripts/test_websocket.go`

**What it does**:
- Tests real-time WebSocket connections
- Validates event broadcasting
- Tests match room functionality

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
# Test scoring functionality
go run scripts/test_scoring.go

# Test WebSocket functionality  
go run scripts/test_websocket.go
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
