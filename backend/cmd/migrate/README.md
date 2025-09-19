# Database Migration Tool

This Go script automates the execution of database migrations for the Spark Park Cricket backend.

## Overview

The migration tool reads SQL files from `internal/database/migrations/` and attempts to execute them against your Supabase database. It provides detailed logging and fallback instructions for manual execution when needed.

## Usage

### Environment Variables

The script requires the following environment variables:

- `SUPABASE_URL`: Your Supabase project URL
- `SUPABASE_API_KEY`: Your Supabase API key (service role key)

### Running the Migration

```bash
# Set environment variables
export SUPABASE_URL="https://your-project.supabase.co"
export SUPABASE_API_KEY="your-service-role-key"

# Run the migration
go run cmd/migrate/main.go
```

### In GitHub Actions

The migration script is automatically executed in the CI/CD pipeline for integration and E2E tests.

## Features

- **Automatic SQL Processing**: Reads all `.sql` files from the migrations directory
- **Statement Parsing**: Splits SQL into individual statements for better error handling
- **DDL Detection**: Automatically handles CREATE TABLE, CREATE INDEX, and COMMENT statements
- **Detailed Logging**: Provides comprehensive output about what's being executed
- **Fallback Instructions**: When automatic execution fails, provides manual execution instructions
- **Error Handling**: Gracefully handles missing files, invalid SQL, and connection issues

## Migration Files

Place your SQL migration files in `internal/database/migrations/` with the following naming convention:

- `001_initial_schema.sql`
- `002_add_indexes.sql`
- `003_update_constraints.sql`

## Supported SQL Statements

The script automatically handles:

- `CREATE EXTENSION`
- `CREATE TABLE`
- `CREATE INDEX`
- `COMMENT ON TABLE`
- `COMMENT ON COLUMN`

Other statements may require manual execution in the Supabase Dashboard.

## Output Example

```
🗄️ Starting database migration...
📍 Supabase URL: https://your-project.supabase.co
🔑 API Key: ********
📁 Found 1 migration files: [001_complete_schema.sql]

🔄 Processing migration: 001_complete_schema.sql
📝 Executing SQL from 001_complete_schema.sql...
   Statement 1: CREATE EXTENSION IF NOT EXISTS "uuid-ossp"
   ✅ DDL statement executed successfully
   Statement 2: CREATE TABLE IF NOT EXISTS dev_v1.series (...)
   ✅ DDL statement executed successfully
   ...

🎉 Database migration completed!
📋 Summary:
   - Processed 1 migration files
   - Successfully executed: 1
   - Manual execution required: 0
```

## Troubleshooting

### Missing Environment Variables

```
❌ SUPABASE_URL environment variable is required
❌ SUPABASE_API_KEY environment variable is required
```

**Solution**: Set the required environment variables before running the script.

### Migration Directory Not Found

```
❌ Migration directory not found: internal/database/migrations
```

**Solution**: Ensure you're running the script from the backend directory root.

### Manual Execution Required

When the script cannot automatically execute certain SQL statements, it will:

1. Display the full SQL content
2. Provide a direct link to the Supabase Dashboard
3. Suggest manual execution steps

## Integration with CI/CD

The migration script is integrated into the GitHub Actions workflow:

- **Integration Tests**: Runs migrations before integration tests
- **E2E Tests**: Runs migrations before end-to-end tests
- **Error Handling**: Fails gracefully with clear error messages
- **Logging**: Provides detailed output for debugging

## Security Notes

- The script masks API keys in output for security
- Uses service role keys for database operations
- Validates environment variables before execution
- Provides secure fallback instructions

## Development

To modify the migration script:

1. Edit `cmd/migrate/main.go`
2. Test locally with dummy environment variables
3. Ensure it compiles without errors
4. Test in the CI/CD pipeline

## Related Files

- `internal/database/migrations/`: SQL migration files
- `MIGRATION_GUIDE.md`: Detailed migration instructions
- `.github/workflows/backend-ci.yml`: CI/CD integration
