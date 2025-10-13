-- Migration: Add team names to matches table
-- Description: Adds team_a_name and team_b_name columns to matches table for custom team naming

-- Add team_a_name column to matches table (nullable for backward compatibility)
ALTER TABLE matches ADD COLUMN IF NOT EXISTS team_a_name VARCHAR(100);

-- Add team_b_name column to matches table (nullable for backward compatibility)
ALTER TABLE matches ADD COLUMN IF NOT EXISTS team_b_name VARCHAR(100);

-- Add comments for documentation
COMMENT ON COLUMN matches.team_a_name IS 'Custom name for Team A in this match (copied from series or overridden)';
COMMENT ON COLUMN matches.team_b_name IS 'Custom name for Team B in this match (copied from series or overridden)';

