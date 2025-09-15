-- Complete Schema Update Migration
-- Fixes all schema issues for simplified cricket system
-- Version: 2.0.1
-- Date: 2025-09-10

-- ============================================
-- MATCHES TABLE UPDATES
-- ============================================

-- Drop existing foreign key constraints and columns that are no longer needed
ALTER TABLE matches DROP CONSTRAINT IF EXISTS matches_team1_id_fkey;
ALTER TABLE matches DROP CONSTRAINT IF EXISTS matches_team2_id_fkey;

-- Drop old columns
ALTER TABLE matches DROP COLUMN IF EXISTS team1_id;
ALTER TABLE matches DROP COLUMN IF EXISTS team2_id;

-- Add new columns for simplified structure
ALTER TABLE matches ADD COLUMN IF NOT EXISTS team_a_player_count INTEGER DEFAULT 11 CHECK (team_a_player_count >= 1 AND team_a_player_count <= 11);
ALTER TABLE matches ADD COLUMN IF NOT EXISTS team_b_player_count INTEGER DEFAULT 11 CHECK (team_b_player_count >= 1 AND team_b_player_count <= 11);
ALTER TABLE matches ADD COLUMN IF NOT EXISTS total_overs INTEGER DEFAULT 20 CHECK (total_overs >= 1 AND total_overs <= 20);
ALTER TABLE matches ADD COLUMN IF NOT EXISTS toss_winner VARCHAR(1) CHECK (toss_winner IN ('A', 'B'));
ALTER TABLE matches ADD COLUMN IF NOT EXISTS toss_type VARCHAR(1) CHECK (toss_type IN ('H', 'T'));
ALTER TABLE matches ADD COLUMN IF NOT EXISTS batting_team VARCHAR(1) DEFAULT 'A' CHECK (batting_team IN ('A', 'B'));

-- Update status constraint to remove 'scheduled' and add 'live' as default
ALTER TABLE matches DROP CONSTRAINT IF EXISTS matches_status_check;
ALTER TABLE matches ADD CONSTRAINT matches_status_check CHECK (status IN ('live', 'completed', 'cancelled'));

-- Update existing records to have default values
UPDATE matches SET 
    team_a_player_count = 11,
    team_b_player_count = 11,
    total_overs = 20,
    toss_winner = 'A',
    toss_type = 'H',
    batting_team = 'A'
WHERE team_a_player_count IS NULL OR team_b_player_count IS NULL OR total_overs IS NULL OR toss_winner IS NULL OR toss_type IS NULL OR batting_team IS NULL;

-- Make new columns NOT NULL after setting defaults
ALTER TABLE matches ALTER COLUMN team_a_player_count SET NOT NULL;
ALTER TABLE matches ALTER COLUMN team_b_player_count SET NOT NULL;
ALTER TABLE matches ALTER COLUMN total_overs SET NOT NULL;
ALTER TABLE matches ALTER COLUMN toss_winner SET NOT NULL;
ALTER TABLE matches ALTER COLUMN toss_type SET NOT NULL;
ALTER TABLE matches ALTER COLUMN batting_team SET NOT NULL;

-- ============================================
-- LIVE_SCOREBOARD TABLE UPDATES
-- ============================================

-- Update live_scoreboard table to use batting_team instead of batting_team_id
ALTER TABLE live_scoreboard DROP CONSTRAINT IF EXISTS live_scoreboard_batting_team_id_fkey;
ALTER TABLE live_scoreboard DROP COLUMN IF EXISTS batting_team_id;
ALTER TABLE live_scoreboard ADD COLUMN IF NOT EXISTS batting_team VARCHAR(1) DEFAULT 'A' CHECK (batting_team IN ('A', 'B'));
UPDATE live_scoreboard SET batting_team = 'A' WHERE batting_team IS NULL;
ALTER TABLE live_scoreboard ALTER COLUMN batting_team SET NOT NULL;

-- ============================================
-- OVERS TABLE UPDATES
-- ============================================

-- Update overs table to use batting_team instead of batting_team_id
ALTER TABLE overs DROP CONSTRAINT IF EXISTS overs_batting_team_id_fkey;
ALTER TABLE overs DROP COLUMN IF EXISTS batting_team_id;
ALTER TABLE overs ADD COLUMN IF NOT EXISTS batting_team VARCHAR(1) DEFAULT 'A' CHECK (batting_team IN ('A', 'B'));
UPDATE overs SET batting_team = 'A' WHERE batting_team IS NULL;
ALTER TABLE overs ALTER COLUMN batting_team SET NOT NULL;

-- ============================================
-- BALLS TABLE UPDATES
-- ============================================

-- Update balls table to use new run_type instead of runs and remove player references
ALTER TABLE balls DROP CONSTRAINT IF EXISTS balls_batsman_id_fkey;
ALTER TABLE balls DROP CONSTRAINT IF EXISTS balls_bowler_id_fkey;
ALTER TABLE balls DROP COLUMN IF EXISTS batsman_id;
ALTER TABLE balls DROP COLUMN IF EXISTS bowler_id;
ALTER TABLE balls DROP COLUMN IF EXISTS runs;
ALTER TABLE balls ADD COLUMN IF NOT EXISTS run_type VARCHAR(2) DEFAULT '0' CHECK (run_type IN ('1', '2', '3', '4', '5', '6', '7', '8', '9', 'NB', 'WD', 'LB'));
UPDATE balls SET run_type = '0' WHERE run_type IS NULL;
ALTER TABLE balls ALTER COLUMN run_type SET NOT NULL;

-- ============================================
-- INDEXES AND PERFORMANCE
-- ============================================

-- Create indexes for new columns
CREATE INDEX IF NOT EXISTS idx_matches_toss_winner ON matches(toss_winner);
CREATE INDEX IF NOT EXISTS idx_matches_batting_team ON matches(batting_team);
CREATE INDEX IF NOT EXISTS idx_live_scoreboard_batting_team ON live_scoreboard(batting_team);
CREATE INDEX IF NOT EXISTS idx_overs_batting_team ON overs(batting_team);
CREATE INDEX IF NOT EXISTS idx_balls_run_type ON balls(run_type);

-- ============================================
-- DOCUMENTATION AND COMMENTS
-- ============================================

-- Add comments for new columns
COMMENT ON COLUMN matches.team_a_player_count IS 'Number of players in Team A (1-11)';
COMMENT ON COLUMN matches.team_b_player_count IS 'Number of players in Team B (1-11)';
COMMENT ON COLUMN matches.total_overs IS 'Total overs for the match (1-20)';
COMMENT ON COLUMN matches.toss_winner IS 'Team that won the toss: A or B';
COMMENT ON COLUMN matches.toss_type IS 'Toss result: H (Heads) or T (Tails)';
COMMENT ON COLUMN matches.batting_team IS 'Team currently batting: A or B';
COMMENT ON COLUMN balls.run_type IS 'Run type: 1-9 (runs), NB (No Ball), WD (Wide), LB (Leg Byes)';

-- ============================================
-- VERIFICATION
-- ============================================

-- Verify the schema update
SELECT 'Schema update completed successfully!' as status;
SELECT 'Matches table columns:' as info;
SELECT column_name, data_type, is_nullable, column_default 
FROM information_schema.columns 
WHERE table_name = 'matches' 
ORDER BY ordinal_position;

SELECT 'Live scoreboard table columns:' as info;
SELECT column_name, data_type, is_nullable, column_default 
FROM information_schema.columns 
WHERE table_name = 'live_scoreboard' 
ORDER BY ordinal_position;
