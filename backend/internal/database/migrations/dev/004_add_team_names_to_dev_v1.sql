-- Migration: Add team names to series and matches tables (Dev Environment)
-- Description: Adds team_a_name and team_b_name columns to series and matches tables for custom team naming

-- Add team names to series table
ALTER TABLE dev_v1.series ADD COLUMN IF NOT EXISTS team_a_name VARCHAR(100);
ALTER TABLE dev_v1.series ADD COLUMN IF NOT EXISTS team_b_name VARCHAR(100);

-- Add team names to matches table
ALTER TABLE dev_v1.matches ADD COLUMN IF NOT EXISTS team_a_name VARCHAR(100);
ALTER TABLE dev_v1.matches ADD COLUMN IF NOT EXISTS team_b_name VARCHAR(100);

-- Add comments for documentation
COMMENT ON COLUMN dev_v1.series.team_a_name IS 'Custom name for Team A in this series (optional)';
COMMENT ON COLUMN dev_v1.series.team_b_name IS 'Custom name for Team B in this series (optional)';
COMMENT ON COLUMN dev_v1.matches.team_a_name IS 'Custom name for Team A in this match (copied from series or overridden)';
COMMENT ON COLUMN dev_v1.matches.team_b_name IS 'Custom name for Team B in this match (copied from series or overridden)';

