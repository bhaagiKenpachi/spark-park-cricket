# Spark Park Cricket - Database Migration Guide

## Overview

This guide explains how to reset and migrate the Supabase database to the new simplified cricket scoring system.

## What Changed

### 🗑️ **Removed Tables**
- `teams` - No longer needed (using Team A and Team B)
- `players` - No longer needed (simplified player management)

### 🔄 **Updated Tables**
- `matches` - Added toss functionality and team simplification
- `live_scoreboard` - Simplified to use TeamType instead of team IDs
- `overs` - Updated to use TeamType
- `balls` - Updated to use new run types and removed player references

### ✨ **New Features**
- **Toss System**: Heads/Tails toss with winner batting first
- **Team Simplification**: Team A vs Team B with configurable player counts
- **Run Types**: 1-9, NB (No Ball), WD (Wide), LB (Leg Byes)
- **Live by Default**: All matches start with "live" status

## Migration Steps

### 1. **Backup Current Data** (Optional)
```bash
# Export current data if needed
pg_dump your_database > backup_before_migration.sql
```

### 2. **Reset Database**
```bash
# Option 1: Using the comprehensive Go script
go run scripts/reset_and_migrate.go

# Option 2: Using the SQL script directly in Supabase dashboard
# Copy and paste the content of scripts/reset_supabase.sql
```

### 3. **Execute Schema Update** (CRITICAL)
**You MUST execute this manually in Supabase Dashboard:**

1. Go to your Supabase Dashboard
2. Navigate to SQL Editor
3. Copy and paste the content of `internal/database/migrations/003_update_matches_schema.sql`
4. Click 'Run' to execute

**This fixes the `batting_team` column issue and all other schema problems.**

### 4. **Verify Migration**
```bash
# Check if tables were created correctly
curl http://localhost:8080/health/database

# Test API endpoints
curl http://localhost:8080/api/v1/series
curl http://localhost:8080/api/v1/matches
```

## Database Schema

### **Series Table**
```sql
CREATE TABLE series (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    start_date TIMESTAMP WITH TIME ZONE NOT NULL,
    end_date TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### **Matches Table**
```sql
CREATE TABLE matches (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    series_id UUID REFERENCES series(id) ON DELETE CASCADE,
    match_number INTEGER NOT NULL,
    date TIMESTAMP WITH TIME ZONE NOT NULL,
    status VARCHAR(20) DEFAULT 'live' CHECK (status IN ('live', 'completed', 'cancelled')),
    team_a_player_count INTEGER NOT NULL DEFAULT 11 CHECK (team_a_player_count >= 1 AND team_a_player_count <= 11),
    team_b_player_count INTEGER NOT NULL DEFAULT 11 CHECK (team_b_player_count >= 1 AND team_b_player_count <= 11),
    total_overs INTEGER NOT NULL DEFAULT 20 CHECK (total_overs >= 1 AND total_overs <= 20),
    toss_winner VARCHAR(1) NOT NULL CHECK (toss_winner IN ('A', 'B')),
    toss_type VARCHAR(1) NOT NULL CHECK (toss_type IN ('H', 'T')),
    batting_team VARCHAR(1) NOT NULL DEFAULT 'A' CHECK (batting_team IN ('A', 'B')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### **Live Scoreboard Table**
```sql
CREATE TABLE live_scoreboard (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    match_id UUID REFERENCES matches(id) ON DELETE CASCADE,
    batting_team VARCHAR(1) NOT NULL CHECK (batting_team IN ('A', 'B')),
    score INTEGER DEFAULT 0,
    wickets INTEGER DEFAULT 0,
    overs DECIMAL(4,1) DEFAULT 0.0,
    balls INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### **Overs Table**
```sql
CREATE TABLE overs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    match_id UUID REFERENCES matches(id) ON DELETE CASCADE,
    over_number INTEGER NOT NULL,
    batting_team VARCHAR(1) NOT NULL CHECK (batting_team IN ('A', 'B')),
    total_runs INTEGER DEFAULT 0,
    total_balls INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### **Balls Table**
```sql
CREATE TABLE balls (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    over_id UUID REFERENCES overs(id) ON DELETE CASCADE,
    ball_number INTEGER NOT NULL,
    ball_type VARCHAR(20) NOT NULL CHECK (ball_type IN ('good', 'wide', 'no_ball', 'dead_ball')),
    run_type VARCHAR(2) NOT NULL CHECK (run_type IN ('1', '2', '3', '4', '5', '6', '7', '8', '9', 'NB', 'WD', 'LB')),
    is_wicket BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

## Sample Data

The migration creates **empty tables** with no sample data:

- **Series**: Empty table ready for new tournaments
- **Matches**: Empty table ready for new matches
- **Live Scoreboards**: Empty table ready for live scoring
- **Overs**: Empty table ready for over tracking
- **Balls**: Empty table ready for ball-by-ball events

This provides a clean slate for testing and development.

## API Changes

### **Removed Endpoints**
- `GET /api/v1/teams` - Teams API removed
- `POST /api/v1/teams` - Create team removed
- `GET /api/v1/players` - Players API removed
- `POST /api/v1/players` - Create player removed
- `GET /api/v1/scoreboard/{match_id}/ball` - Ball API removed
- `PUT /api/v1/scoreboard/{match_id}/score` - Score API removed
- `PUT /api/v1/scoreboard/{match_id}/wicket` - Wicket API removed

### **Updated Endpoints**
- `POST /api/v1/matches` - Now includes toss functionality
- `PUT /api/v1/matches/{id}` - Can update player counts and overs

### **New Fields**
- `team_a_player_count`: Number of players in Team A (1-11)
- `team_b_player_count`: Number of players in Team B (1-11)
- `total_overs`: Total overs for the match (1-20)
- `toss_winner`: Team that won the toss ("A" or "B")
- `toss_type`: Toss result ("H" for Heads, "T" for Tails)
- `batting_team`: Team currently batting ("A" or "B")

## Testing the Migration

### 1. **Health Check**
```bash
curl http://localhost:8080/health/database
```

### 2. **Create Series**
```bash
curl -X POST http://localhost:8080/api/v1/series \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Series",
    "start_date": "2024-03-22T00:00:00Z",
    "end_date": "2024-03-23T23:59:59Z"
  }'
```

### 3. **Create Match with Toss**
```bash
curl -X POST http://localhost:8080/api/v1/matches \
  -H "Content-Type: application/json" \
  -d '{
    "series_id": "your-series-id",
    "match_number": 1,
    "date": "2024-03-22T20:00:00Z",
    "team_a_player_count": 11,
    "team_b_player_count": 11,
    "total_overs": 20,
    "toss_winner": "A",
    "toss_type": "H"
  }'
```

### 4. **List Matches**
```bash
curl http://localhost:8080/api/v1/matches
```

## Troubleshooting

### **Common Issues**

1. **"Could not find the 'batting_team' column" Error**
   - **Solution**: Execute the schema update migration (`003_update_matches_schema.sql`)
   - This adds the missing `batting_team` column to the `matches` table
   - **Critical**: Must be done manually in Supabase Dashboard

2. **UUID Generation Error**
   - Ensure `uuid-ossp` extension is enabled
   - Check if the extension is properly installed

3. **Foreign Key Constraints**
   - Verify that referenced records exist
   - Check if cascade deletes are working

4. **Check Constraints**
   - Ensure enum values match the constraints
   - Verify player counts are within 1-11 range
   - Check overs are within 1-20 range

### **Verification Queries**

```sql
-- Check table counts (should all be 0)
SELECT 'series' as table_name, COUNT(*) as count FROM series
UNION ALL
SELECT 'matches', COUNT(*) FROM matches
UNION ALL
SELECT 'live_scoreboard', COUNT(*) FROM live_scoreboard
UNION ALL
SELECT 'overs', COUNT(*) FROM overs
UNION ALL
SELECT 'balls', COUNT(*) FROM balls;

-- Verify tables are empty
SELECT * FROM series LIMIT 5; -- Should return no rows
SELECT * FROM matches LIMIT 5; -- Should return no rows
SELECT * FROM live_scoreboard LIMIT 5; -- Should return no rows
```

## Rollback Plan

If you need to rollback:

1. **Restore from backup** (if you created one)
2. **Recreate old schema** using the previous migration files
3. **Restore data** from backup

## Support

For issues or questions:
1. Check the server logs for detailed error messages
2. Verify environment variables are set correctly
3. Ensure Supabase connection is working
4. Check database permissions and constraints

## Next Steps

After successful migration:
1. Update your Postman collection
2. Test all API endpoints
3. Update your frontend application
4. Deploy to production