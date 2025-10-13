-- Migration: Add team names to matches table
-- Description: Adds optional team_a_name and team_b_name columns for custom team naming
-- Version: 011
-- Date: 2025-10-13

-- Add team_a_name column (nullable for backward compatibility)
ALTER TABLE matches 
ADD COLUMN IF NOT EXISTS team_a_name VARCHAR(255);

-- Add team_b_name column (nullable for backward compatibility)
ALTER TABLE matches 
ADD COLUMN IF NOT EXISTS team_b_name VARCHAR(255);

-- Add comments for documentation
COMMENT ON COLUMN matches.team_a_name IS 'Custom name for Team A (defaults to "Team A" if null)';
COMMENT ON COLUMN matches.team_b_name IS 'Custom name for Team B (defaults to "Team B" if null)';

