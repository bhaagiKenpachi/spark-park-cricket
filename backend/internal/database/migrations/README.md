# Spark Park Cricket Database Schema

This directory contains comprehensive SQL migration scripts for the Spark Park Cricket tournament management system.

## Overview

The database schema supports multiple environments and includes:
- **User Authentication**: Google OAuth integration with session management
- **Cricket Tournament Management**: Series, matches, and live scoring
- **Ball-by-Ball Tracking**: Complete cricket scoring with all ball types and wicket types
- **Real-time Updates**: WebSocket-ready schema for live scoreboards

## Schema Structure

### Global Schema (Public)
Shared across all environments:
- `users` - User authentication and profile data
- `user_sessions` - Session management
- `oauth_states` - OAuth security state parameters

### Environment-Specific Schemas
- `dev_v1` - Development environment
- `testing_db` - Testing environment  
- `prod_v1` - Production environment

Each environment schema contains:
- `series` - Cricket tournaments and competitions
- `matches` - Individual cricket matches
- `live_scoreboard` - Real-time match scoring
- `innings` - Cricket innings tracking
- `overs` - Over-by-over tracking
- `balls` - Ball-by-ball events

## Migration Files

### Complete Schema Scripts

1. **`complete_schema_all_environments.sql`** - Template script for any environment
   - Replace `{SCHEMA_NAME}` with your target schema
   - Use for custom environments or as a template

2. **`complete_schema_dev_v1.sql`** - Development environment
   - Ready-to-run script for development
   - Creates `dev_v1` schema tables
   - Automatically grants permissions to all Supabase roles

3. **`complete_schema_testing_db.sql`** - Testing environment
   - Ready-to-run script for testing
   - Creates `testing_db` schema tables
   - Automatically grants permissions to all Supabase roles

4. **`complete_schema_prod_v1.sql`** - Production environment
   - Ready-to-run script for production
   - Creates `prod_v1` schema tables
   - Automatically grants permissions to all Supabase roles

### Legacy Migration Files

- `001_complete_schema.sql` - Original complete schema (dev_v1 only)
- `002_create_users_tables.sql` - User authentication tables
- `003_create_oauth_states_table.sql` - OAuth state management
- `004_add_created_by_fields.sql` - User ownership tracking
- `005_add_created_by_fields_testing_db.sql` - Testing schema ownership
- `v1.1.0_complete_schema_with_auth.sql` - Combined migration (all environments)

## Usage Instructions

### For Development Environment

```sql
-- 1. Create the schema (if it doesn't exist)
CREATE SCHEMA IF NOT EXISTS dev_v1;

-- 2. Run the complete schema script
\i complete_schema_dev_v1.sql
```

### For Testing Environment

```sql
-- 1. Create the schema (if it doesn't exist)
CREATE SCHEMA IF NOT EXISTS testing_db;

-- 2. Run the complete schema script
\i complete_schema_testing_db.sql
```

### For Production Environment

```sql
-- 1. Create the schema (if it doesn't exist)
CREATE SCHEMA IF NOT EXISTS prod_v1;

-- 2. Run the complete schema script
\i complete_schema_prod_v1.sql
```

### For Custom Environment

```sql
-- 1. Create your custom schema
CREATE SCHEMA IF NOT EXISTS your_custom_schema;

-- 2. Edit complete_schema_all_environments.sql
-- Replace all instances of {SCHEMA_NAME} with your_custom_schema

-- 3. Run the modified script
\i complete_schema_all_environments.sql
```

## Schema Features

### Constraints and Validations

- **Match Constraints**: Team player counts (1-20), overs (1-20), status validation
- **Toss Constraints**: Winner (A/B), type (H/T), batting team validation
- **Innings Constraints**: Innings number (1-2), wickets (0-10), status validation
- **Over Constraints**: Over number (≥1), balls (0-6), status validation
- **Ball Constraints**: Ball number (1-20), ball types, run types, wicket validation
- **Wicket Logic**: Ensures wicket_type is NULL when is_wicket is false

### Indexes for Performance

- **User Authentication**: Google ID, email, session lookups
- **Series**: Date ranges, creator tracking
- **Matches**: Series relationships, status, dates, toss information
- **Scoreboard**: Match relationships, batting team
- **Innings**: Match relationships, batting team, status
- **Overs**: Innings relationships, status
- **Balls**: Over relationships, run types, wicket tracking

### Automatic Features

- **UUID Generation**: All primary keys use UUID with automatic generation
- **Timestamps**: Automatic created_at and updated_at timestamps
- **Triggers**: Automatic updated_at timestamp updates on record changes
- **Cascading Deletes**: Proper foreign key relationships with cascade options
- **Permission Grants**: Automatic permission grants to all Supabase roles (all schemas)

### Cricket-Specific Features

- **Ball Types**: good, wide, no_ball, dead_ball
- **Run Types**: 0-9 (runs), NB (No Ball), WD (Wide), LB (Leg Byes), WC (Wicket)
- **Wicket Types**: bowled, caught, lbw, run_out, stumped, hit_wicket
- **Match States**: live, completed, cancelled
- **Innings States**: in_progress, completed
- **Over States**: in_progress, completed

## Verification

Each script includes verification queries that will:
1. Confirm successful schema creation
2. List all created tables in global schema
3. List all created tables in environment-specific schema
4. Display any errors or warnings

## Environment Configuration

The application uses these schemas based on configuration:

- **Development**: `DATABASE_SCHEMA=dev_v1`
- **Testing**: `TEST_SCHEMA=testing_db`
- **Production**: `DATABASE_SCHEMA=prod_v1`

## Best Practices

1. **Always create schemas first** before running migration scripts
2. **Use environment-specific scripts** for easier maintenance
3. **Verify schema creation** using the built-in verification queries
4. **Test in development** before applying to production
5. **Backup production data** before running migrations
6. **Use transactions** for production deployments

## Troubleshooting

### Common Issues

1. **Schema doesn't exist**: Create the schema first using `CREATE SCHEMA IF NOT EXISTS schema_name;`
2. **Permission errors**: Ensure your database user has CREATE privileges
3. **Extension errors**: Ensure `uuid-ossp` extension is available
4. **Constraint violations**: Check that your data meets the defined constraints
5. **Supabase permission errors**: All schemas (dev_v1, testing_db, prod_v1) automatically grant permissions to all Supabase roles

### Verification Commands

```sql
-- Check if schemas exist
SELECT schema_name FROM information_schema.schemata 
WHERE schema_name IN ('dev_v1', 'testing_db', 'prod_v1');

-- Check if tables exist in a schema
SELECT table_name FROM information_schema.tables 
WHERE table_schema = 'your_schema_name';

-- Check if indexes exist
SELECT indexname FROM pg_indexes 
WHERE schemaname = 'your_schema_name';
```

## Support

For issues or questions about the database schema:
1. Check the verification output for errors
2. Review the constraint definitions
3. Ensure proper schema creation
4. Verify database permissions
