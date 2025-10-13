-- Migration: Add team names to series table
-- Description: Adds team_a_name and team_b_name columns to series table for custom team naming

-- Add team_a_name column to series table (nullable for backward compatibility)
ALTER TABLE series ADD COLUMN IF NOT EXISTS team_a_name VARCHAR(100);

-- Add team_b_name column to series table (nullable for backward compatibility)
ALTER TABLE series ADD COLUMN IF NOT EXISTS team_b_name VARCHAR(100);

-- Add comments for documentation
COMMENT ON COLUMN series.team_a_name IS 'Custom name for Team A in this series (optional)';
COMMENT ON COLUMN series.team_b_name IS 'Custom name for Team B in this series (optional)';

