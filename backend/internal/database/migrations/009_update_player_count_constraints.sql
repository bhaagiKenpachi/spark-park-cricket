-- Update team player count constraints to allow up to 20 players
-- Version: 2.0.1
-- Date: 2025-01-15

-- Drop existing check constraints
ALTER TABLE matches DROP CONSTRAINT IF EXISTS matches_team_a_player_count_check;
ALTER TABLE matches DROP CONSTRAINT IF EXISTS matches_team_b_player_count_check;

-- Add new check constraints allowing up to 20 players
ALTER TABLE matches ADD CONSTRAINT matches_team_a_player_count_check 
    CHECK (team_a_player_count >= 1 AND team_a_player_count <= 20);

ALTER TABLE matches ADD CONSTRAINT matches_team_b_player_count_check 
    CHECK (team_b_player_count >= 1 AND team_b_player_count <= 20);

-- Update the default values to be more flexible (keeping 11 as default)
ALTER TABLE matches ALTER COLUMN team_a_player_count SET DEFAULT 11;
ALTER TABLE matches ALTER COLUMN team_b_player_count SET DEFAULT 11;

-- Update comments to reflect new limits
COMMENT ON COLUMN matches.team_a_player_count IS 'Number of players in Team A (1-20)';
COMMENT ON COLUMN matches.team_b_player_count IS 'Number of players in Team B (1-20)';

-- Verify the changes
SELECT 'Player count constraints updated successfully - now allowing 1-20 players per team!' as status;
