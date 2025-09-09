# Database Migration Guide

This guide explains how to use the Go migration tool to manage your Supabase database schema.

## Quick Start

### 1. Configure Environment
Create `.env` file with your Supabase credentials:
```env
SUPABASE_URL=https://qehkpqubnnpbaejhcwvx.supabase.co
SUPABASE_API_KEY=your_anon_public_key_here
DATABASE_PASSWORD=your_database_password
PORT=8080
```

### 2. Run Migrations
```bash
# Install dependencies
go mod tidy

# Run all pending migrations
go run cmd/migrate/main.go
```

## Your Supabase Project

Based on your project reference `qehkpqubnnpbaejhcwvx`, your configuration should be:

```env
SUPABASE_URL=https://qehkpqubnnpbaejhcwvx.supabase.co
SUPABASE_API_KEY=your_anon_public_key_here
DATABASE_PASSWORD=your_database_password
```

## Migration Files

The tool will automatically run all SQL files in `internal/database/migrations/` in order:

- `001_initial_schema.sql` - Creates all tables, indexes, and triggers
- `002_sample_data.sql` - Inserts sample data (if exists)

## Expected Output

```
🏏 Spark Park Cricket - Database Migration Runner
================================================
Connecting to database...
✓ Database connection established
Starting database migrations...
Applying migration: 001 (001_initial_schema)
✓ Migration 001 applied successfully
✓ Applied 1 new migrations
🎉 Database migrations completed successfully!
```

## What Gets Created

The initial migration creates:

### Tables
- `series` - Cricket tournaments/competitions
- `teams` - Cricket teams
- `players` - Individual players
- `matches` - Cricket matches
- `live_scoreboard` - Live match scoring
- `overs` - Over-by-over tracking
- `balls` - Ball-by-ball events
- `schema_migrations` - Migration tracking

### Features
- UUID primary keys
- Automatic timestamps
- Foreign key constraints
- Performance indexes
- Update triggers

## Troubleshooting

### Connection Issues
- Verify `SUPABASE_URL` format
- Check `DATABASE_PASSWORD` is correct
- Ensure Supabase project is active

### Migration Errors
- Check SQL syntax in migration files
- Verify no conflicting constraints
- Check for existing tables

### Environment Issues
- Ensure `.env` file exists in backend root
- Restart terminal after changing `.env`
- Check environment variable names

## Next Steps

After running migrations:

1. **Start Server**: `go run cmd/server/main.go`
2. **Test Health**: `curl http://localhost:8080/health`
3. **Use Postman**: Import collection and test APIs
4. **Create Data**: Use API endpoints to create series, teams, matches

## Advanced Usage

### Adding New Migrations

1. Create new SQL file: `003_add_new_feature.sql`
2. Use version number higher than existing
3. Run migration tool: `go run cmd/migrate/main.go`

### Migration Best Practices

- Use descriptive filenames
- Make migrations idempotent
- Test on development first
- Backup before major changes
- Use transactions for data changes

## Integration

The migration tool integrates with:
- **Supabase**: Direct PostgreSQL connection
- **Go Modules**: Uses `github.com/lib/pq` driver
- **Environment**: Reads from `.env` file
- **Version Control**: Tracks applied migrations

## Support

For issues:
1. Check migration tool README: `cmd/migrate/README.md`
2. Verify Supabase connection
3. Check migration file syntax
4. Review error messages
