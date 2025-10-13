-- Migration: Add team names to series and matches tables in testing_db schema
-- Description: Adds team_a_name and team_b_name columns to series and matches tables in testing_db

-- Add team_a_name column to series table (nullable for backward compatibility)
ALTER TABLE testing_db.series ADD COLUMN IF NOT EXISTS team_a_name VARCHAR(100);

-- Add team_b_name column to series table (nullable for backward compatibility)
ALTER TABLE testing_db.series ADD COLUMN IF NOT EXISTS team_b_name VARCHAR(100);

-- Add team_a_name column to matches table (nullable for backward compatibility)
ALTER TABLE testing_db.matches ADD COLUMN IF NOT EXISTS team_a_name VARCHAR(100);

-- Add team_b_name column to matches table (nullable for backward compatibility)
ALTER TABLE testing_db.matches ADD COLUMN IF NOT EXISTS team_b_name VARCHAR(100);

-- Add comments for documentation
COMMENT ON COLUMN testing_db.series.team_a_name IS 'Custom name for Team A in this series (optional)';
COMMENT ON COLUMN testing_db.series.team_b_name IS 'Custom name for Team B in this series (optional)';
COMMENT ON COLUMN testing_db.matches.team_a_name IS 'Custom name for Team A in this match (copied from series or overridden)';
COMMENT ON COLUMN testing_db.matches.team_b_name IS 'Custom name for Team B in this match (copied from series or overridden)';

