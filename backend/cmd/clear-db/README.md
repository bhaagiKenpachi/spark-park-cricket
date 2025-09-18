# Database Clear Script

This script clears all data from the Spark Park Cricket database tables in the `prod_v1` schema.

## ⚠️ WARNING

**THIS SCRIPT WILL PERMANENTLY DELETE ALL DATA FROM THE DATABASE!**

Make sure you have proper backups before running this script.

## Usage

### Using Makefile (Recommended)

```bash
make clear-db
```

### Direct Execution

```bash
go run cmd/clear-db/main.go
```

### Build and Run

```bash
go build -o bin/clear-db cmd/clear-db/main.go
./bin/clear-db
```

## What it does

1. **Displays target schema** - Clearly shows which database schema will be cleared
2. **Shows current table counts** - Displays the number of records in each table before clearing
3. **Asks for confirmation** - Requires you to type 'YES' to confirm the operation
4. **Clears tables in correct order** - Respects foreign key constraints by clearing child tables first:
   - `balls` (references overs)
   - `overs` (references innings)
   - `innings` (references matches)
   - `live_scoreboard` (references matches)
   - `matches` (references series)
   - `series` (no dependencies)
   - `schema_version` (no dependencies)
5. **Shows final counts** - Displays table counts after clearing (should all be 0)

## Tables Cleared

- `balls` - Ball-by-ball events
- `overs` - Over-by-over tracking
- `innings` - Cricket innings tracking
- `live_scoreboard` - Real-time match scoring
- `matches` - Individual cricket matches
- `series` - Cricket tournaments/competitions
- `schema_version` - Schema version tracking

## Safety Features

- **Confirmation prompt** - Must type 'YES' to proceed
- **Clear warnings** - Shows exactly what will be deleted
- **Ordered deletion** - Respects foreign key constraints
- **Error handling** - Graceful error handling for missing tables
- **Progress logging** - Shows progress as each table is cleared
- **Supabase safety** - Uses proper WHERE clauses to comply with Supabase DELETE requirements
- **Missing table handling** - Gracefully handles tables that don't exist in the schema
- **Schema visibility** - Clearly displays which schema is being targeted for clearing

## Environment Requirements

The script uses the same environment variables as the main application:

- `SUPABASE_URL` - Your Supabase project URL
- `SUPABASE_API_KEY` - Your Supabase API key
- `DATABASE_SCHEMA` - Database schema (defaults to 'prod_v1')

## Example Output

```
=== SPARK PARK CRICKET - DATABASE CLEAR SCRIPT ===
This script will clear ALL data from the database tables
==================================================

✅ Connected to Supabase database
Database Schema: testing_db

============================================================
🗄️  CLEARING DATA FROM SCHEMA: TESTING_DB
============================================================

=== TABLE COUNTS IN SCHEMA 'testing_db' ===
📊 testing_db.balls: 150 records
📊 testing_db.overs: 25 records
📊 testing_db.innings: 2 records
📊 testing_db.live_scoreboard: 1 records
📊 testing_db.matches: 1 records
📊 testing_db.series: 1 records
📊 testing_db.schema_version: 1 records
=============================

⚠️  WARNING: This will permanently delete ALL data from the following tables in schema 'testing_db':
   - balls
   - overs
   - innings
   - live_scoreboard
   - matches
   - series
   - schema_version

Are you sure you want to continue? Type 'YES' to confirm: YES

=== STARTING TABLE CLEARING ===
🗑️  Clearing table: balls
✅ Cleared table: balls
🗑️  Clearing table: overs
✅ Cleared table: overs
🗑️  Clearing table: innings
✅ Cleared table: innings
🗑️  Clearing table: live_scoreboard
✅ Cleared table: live_scoreboard
🗑️  Clearing table: matches
✅ Cleared table: matches
🗑️  Clearing table: series
✅ Cleared table: series
🗑️  Clearing table: schema_version
✅ Cleared table: schema_version
✅ All tables cleared successfully

=== CLEARING COMPLETED ===
🗄️  Schema 'testing_db' has been cleared
=== TABLE COUNTS IN SCHEMA 'testing_db' ===
📊 testing_db.balls: 0 records
📊 testing_db.overs: 0 records
📊 testing_db.innings: 0 records
📊 testing_db.live_scoreboard: 0 records
📊 testing_db.matches: 0 records
📊 testing_db.series: 0 records
📊 testing_db.schema_version: 0 records
=============================
✅ All tables have been cleared successfully!
```
