# Database Migration Tool

This Go script runs database migrations against your Supabase PostgreSQL database.

## Features

- ✅ **Automatic Migration Tracking**: Tracks applied migrations in `schema_migrations` table
- ✅ **Transaction Safety**: Each migration runs in a transaction
- ✅ **Version Control**: Migrations are ordered by version number
- ✅ **Idempotent**: Safe to run multiple times
- ✅ **Error Handling**: Comprehensive error reporting
- ✅ **Progress Tracking**: Shows migration progress and status

## Prerequisites

1. **Environment Variables**: Set up your `.env` file with Supabase credentials
2. **Go Dependencies**: Install required Go packages
3. **Migration Files**: Ensure migration files exist in `internal/database/migrations/`

## Setup

### 1. Environment Configuration

Create a `.env` file in the backend root directory:

```env
# Supabase Configuration
SUPABASE_URL=https://qehkpqubnnpbaejhcwvx.supabase.co
SUPABASE_API_KEY=your_anon_public_key_here

# Database Configuration (for migrations)
DATABASE_PASSWORD=your_database_password

# Server Configuration
PORT=8080
```

### 2. Install Dependencies

```bash
# Install PostgreSQL driver
go mod tidy
```

### 3. Migration Files

Ensure your migration files are in `internal/database/migrations/` with the naming convention:
- `001_initial_schema.sql`
- `002_sample_data.sql`
- `003_add_indexes.sql`

## Usage

### Run All Migrations

```bash
# From the backend directory
go run cmd/migrate/main.go
```

### Expected Output

```
🏏 Spark Park Cricket - Database Migration Runner
================================================
Connecting to database...
✓ Database connection established
Starting database migrations...
Applying migration: 001 (001_initial_schema)
✓ Migration 001 applied successfully
Applying migration: 002 (002_sample_data)
✓ Migration 002 applied successfully
✓ Applied 2 new migrations
🎉 Database migrations completed successfully!
```

### Subsequent Runs

```
🏏 Spark Park Cricket - Database Migration Runner
================================================
Connecting to database...
✓ Database connection established
Starting database migrations...
⏭️  Migration 001 already applied
⏭️  Migration 002 already applied
✓ All migrations are up to date
🎉 Database migrations completed successfully!
```

## Migration File Format

### Naming Convention
- Format: `{version}_{description}.sql`
- Example: `001_initial_schema.sql`
- Version must be numeric and sortable

### File Structure
```sql
-- Migration: 001_initial_schema.sql
-- Description: Create initial database schema

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create tables
CREATE TABLE IF NOT EXISTS series (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    start_date TIMESTAMP WITH TIME ZONE NOT NULL,
    end_date TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Add indexes, constraints, etc.
```

## How It Works

### 1. Connection
- Connects to Supabase PostgreSQL database using connection string
- Extracts project reference from `SUPABASE_URL`
- Uses `DATABASE_PASSWORD` for authentication

### 2. Migration Tracking
- Creates `schema_migrations` table if it doesn't exist
- Tracks applied migrations by version number
- Prevents duplicate migration execution

### 3. Migration Execution
- Scans `internal/database/migrations/` directory
- Sorts migrations by version number
- Applies only pending migrations
- Each migration runs in a transaction

### 4. Error Handling
- Validates database connection
- Checks migration file existence
- Rolls back failed migrations
- Provides detailed error messages

## Troubleshooting

### Connection Issues

**Error**: `Failed to connect to database`
- Check `SUPABASE_URL` format
- Verify `DATABASE_PASSWORD` is correct
- Ensure Supabase project is active

**Error**: `Failed to ping database`
- Check network connectivity
- Verify database is accessible
- Check firewall settings

### Migration Issues

**Error**: `Migrations directory not found`
- Ensure `internal/database/migrations/` directory exists
- Check current working directory

**Error**: `Failed to execute migration`
- Check SQL syntax in migration file
- Verify table/column names
- Check for conflicting constraints

### Environment Issues

**Error**: `SUPABASE_URL environment variable is required`
- Create `.env` file in backend root
- Set `SUPABASE_URL` with your project URL
- Restart terminal/IDE to load environment

## Advanced Usage

### Custom Migration Directory

Modify the `migrationsDir` variable in `main.go`:
```go
migrationsDir := "custom/migrations/path"
```

### Connection String Customization

Modify the `GetConnectionString()` function for custom connection parameters:
```go
connStr := fmt.Sprintf("postgresql://postgres:%s@db.%s.supabase.co:5432/postgres?sslmode=require&connect_timeout=30",
    os.Getenv("DATABASE_PASSWORD"), projectRef)
```

### Migration Validation

Add validation before applying migrations:
```go
// Validate migration content
if strings.Contains(content, "DROP TABLE") {
    return fmt.Errorf("migration contains dangerous DROP TABLE statement")
}
```

## Integration with CI/CD

### GitHub Actions Example

```yaml
- name: Run Database Migrations
  run: |
    cd backend
    go run cmd/migrate/main.go
  env:
    SUPABASE_URL: ${{ secrets.SUPABASE_URL }}
    SUPABASE_API_KEY: ${{ secrets.SUPABASE_API_KEY }}
    DATABASE_PASSWORD: ${{ secrets.DATABASE_PASSWORD }}
```

### Docker Example

```dockerfile
# Add migration step to Dockerfile
COPY . .
RUN go mod tidy
CMD ["go", "run", "cmd/migrate/main.go"]
```

## Best Practices

1. **Version Numbers**: Use zero-padded numbers (001, 002, 003)
2. **Descriptive Names**: Use clear, descriptive migration names
3. **Idempotent**: Write migrations that can run multiple times safely
4. **Backup**: Always backup database before major migrations
5. **Test**: Test migrations on development database first
6. **Review**: Review migration SQL before applying to production

## Migration Examples

### Creating Tables
```sql
-- 001_create_series_table.sql
CREATE TABLE IF NOT EXISTS series (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    start_date TIMESTAMP WITH TIME ZONE NOT NULL,
    end_date TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### Adding Columns
```sql
-- 002_add_series_description.sql
ALTER TABLE series ADD COLUMN IF NOT EXISTS description TEXT;
```

### Creating Indexes
```sql
-- 003_add_series_indexes.sql
CREATE INDEX IF NOT EXISTS idx_series_name ON series(name);
CREATE INDEX IF NOT EXISTS idx_series_dates ON series(start_date, end_date);
```

### Data Migrations
```sql
-- 004_populate_initial_data.sql
INSERT INTO series (name, start_date, end_date) VALUES
('IPL 2024', '2024-03-22T00:00:00Z', '2024-05-26T23:59:59Z'),
('World Cup 2024', '2024-10-01T00:00:00Z', '2024-11-15T23:59:59Z')
ON CONFLICT DO NOTHING;
```
