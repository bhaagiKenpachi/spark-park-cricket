-- Spark Park Cricket - Vote Teams Migration
-- Version: 010
-- Date: 2025-10-08
-- Description: Create team tables for dividing voters into teams

-- ============================================
-- VOTE TEAMS TABLES
-- ============================================

-- Create vote_teams table
CREATE TABLE IF NOT EXISTS testing_db.vote_teams (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vote_id UUID NOT NULL REFERENCES testing_db.votes(id) ON DELETE CASCADE,
    team_name VARCHAR(100) NOT NULL,
    team_letter VARCHAR(1) NOT NULL CHECK (team_letter IN ('A', 'B')),
    captain_id UUID NOT NULL REFERENCES testing_db.users(id) ON DELETE CASCADE,
    created_by UUID NOT NULL REFERENCES testing_db.users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Ensure only two teams (A and B) per vote
    UNIQUE(vote_id, team_letter)
);

-- Create team_players table
CREATE TABLE IF NOT EXISTS testing_db.team_players (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    team_id UUID NOT NULL REFERENCES testing_db.vote_teams(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES testing_db.users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Ensure one user can only be in one team per vote
    UNIQUE(team_id, user_id)
);

-- ============================================
-- INDEXES FOR PERFORMANCE
-- ============================================

-- Indexes for vote_teams table
CREATE INDEX IF NOT EXISTS idx_vote_teams_vote_id ON testing_db.vote_teams(vote_id);
CREATE INDEX IF NOT EXISTS idx_vote_teams_captain_id ON testing_db.vote_teams(captain_id);
CREATE INDEX IF NOT EXISTS idx_vote_teams_created_by ON testing_db.vote_teams(created_by);
CREATE INDEX IF NOT EXISTS idx_vote_teams_team_letter ON testing_db.vote_teams(team_letter);

-- Indexes for team_players table
CREATE INDEX IF NOT EXISTS idx_team_players_team_id ON testing_db.team_players(team_id);
CREATE INDEX IF NOT EXISTS idx_team_players_user_id ON testing_db.team_players(user_id);

-- ============================================
-- TRIGGERS FOR UPDATED_AT
-- ============================================

-- Function already exists from previous migrations
-- Triggers for vote_teams table
DROP TRIGGER IF EXISTS update_vote_teams_updated_at ON testing_db.vote_teams;
CREATE TRIGGER update_vote_teams_updated_at 
    BEFORE UPDATE ON testing_db.vote_teams 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ============================================
-- FUNCTIONS FOR VALIDATION
-- ============================================

-- Function to check if user has voted
CREATE OR REPLACE FUNCTION testing_db.check_user_voted_in_vote(p_user_id UUID, p_vote_id UUID)
RETURNS BOOLEAN AS $$
BEGIN
    RETURN EXISTS (
        SELECT 1 FROM testing_db.user_votes 
        WHERE user_id = p_user_id AND vote_id = p_vote_id
    );
END;
$$ LANGUAGE plpgsql;

-- Function to check if captain has voted
CREATE OR REPLACE FUNCTION testing_db.check_captain_voted()
RETURNS TRIGGER AS $$
BEGIN
    IF NOT testing_db.check_user_voted_in_vote(NEW.captain_id, NEW.vote_id) THEN
        RAISE EXCEPTION 'Captain must be one of the voters';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Function to check if team player has voted
CREATE OR REPLACE FUNCTION testing_db.check_team_player_voted()
RETURNS TRIGGER AS $$
DECLARE
    v_vote_id UUID;
BEGIN
    -- Get vote_id from team
    SELECT vote_id INTO v_vote_id 
    FROM testing_db.vote_teams 
    WHERE id = NEW.team_id;
    
    -- Check if user voted
    IF NOT testing_db.check_user_voted_in_vote(NEW.user_id, v_vote_id) THEN
        RAISE EXCEPTION 'Only voters can be assigned to teams';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Function to check max players limit (20 per team)
CREATE OR REPLACE FUNCTION testing_db.check_team_max_players()
RETURNS TRIGGER AS $$
DECLARE
    player_count INTEGER;
BEGIN
    -- Count current players in team
    SELECT COUNT(*) INTO player_count 
    FROM testing_db.team_players 
    WHERE team_id = NEW.team_id;
    
    -- Check if limit exceeded
    IF player_count >= 20 THEN
        RAISE EXCEPTION 'Team cannot have more than 20 players';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Function to prevent user from being in both teams
CREATE OR REPLACE FUNCTION testing_db.check_user_not_in_other_team()
RETURNS TRIGGER AS $$
DECLARE
    v_vote_id UUID;
    other_team_count INTEGER;
BEGIN
    -- Get vote_id from current team
    SELECT vote_id INTO v_vote_id 
    FROM testing_db.vote_teams 
    WHERE id = NEW.team_id;
    
    -- Check if user is in another team for this vote
    SELECT COUNT(*) INTO other_team_count
    FROM testing_db.team_players tp
    JOIN testing_db.vote_teams vt ON tp.team_id = vt.id
    WHERE tp.user_id = NEW.user_id 
    AND vt.vote_id = v_vote_id
    AND tp.team_id != NEW.team_id;
    
    IF other_team_count > 0 THEN
        RAISE EXCEPTION 'User cannot be in both teams';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- ============================================
-- TRIGGERS FOR VALIDATION
-- ============================================

-- Trigger to validate captain is a voter
DROP TRIGGER IF EXISTS validate_captain_voted ON testing_db.vote_teams;
CREATE TRIGGER validate_captain_voted
    BEFORE INSERT OR UPDATE ON testing_db.vote_teams
    FOR EACH ROW EXECUTE FUNCTION testing_db.check_captain_voted();

-- Trigger to validate team player is a voter
DROP TRIGGER IF EXISTS validate_team_player_voted ON testing_db.team_players;
CREATE TRIGGER validate_team_player_voted
    BEFORE INSERT ON testing_db.team_players
    FOR EACH ROW EXECUTE FUNCTION testing_db.check_team_player_voted();

-- Trigger to validate max players limit
DROP TRIGGER IF EXISTS validate_team_max_players ON testing_db.team_players;
CREATE TRIGGER validate_team_max_players
    BEFORE INSERT ON testing_db.team_players
    FOR EACH ROW EXECUTE FUNCTION testing_db.check_team_max_players();

-- Trigger to prevent user from being in both teams
DROP TRIGGER IF EXISTS validate_user_not_in_other_team ON testing_db.team_players;
CREATE TRIGGER validate_user_not_in_other_team
    BEFORE INSERT ON testing_db.team_players
    FOR EACH ROW EXECUTE FUNCTION testing_db.check_user_not_in_other_team();

-- ============================================
-- COMMENTS
-- ============================================

COMMENT ON TABLE testing_db.vote_teams IS 'Stores teams for vote-based team division';
COMMENT ON TABLE testing_db.team_players IS 'Stores player assignments to teams';

COMMENT ON COLUMN testing_db.vote_teams.team_letter IS 'Team identifier: A or B (only 2 teams per vote)';
COMMENT ON COLUMN testing_db.vote_teams.captain_id IS 'User ID of team captain (must be a voter)';
COMMENT ON COLUMN testing_db.team_players.user_id IS 'User ID of team player (must be a voter)';

COMMENT ON FUNCTION testing_db.check_user_voted_in_vote(UUID, UUID) IS 'Checks if a user has voted in a specific vote';
COMMENT ON FUNCTION testing_db.check_captain_voted() IS 'Validates that captain is one of the voters';
COMMENT ON FUNCTION testing_db.check_team_player_voted() IS 'Validates that team player is one of the voters';
COMMENT ON FUNCTION testing_db.check_team_max_players() IS 'Ensures team does not exceed 20 players';
COMMENT ON FUNCTION testing_db.check_user_not_in_other_team() IS 'Prevents user from being in both teams';

