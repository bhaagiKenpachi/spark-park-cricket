-- ============================================
-- SPARK PARK CRICKET - PRODUCTION BACKUP
-- ============================================
-- Complete backup from prod_v1 to backup_prod_v1
-- Version: 2.0.0 (New - Based on all production migrations)
-- Date: 2025-10-18
-- Description: Simple backup using SELECT * - matches production schema exactly
-- Migrations applied: 001-005
-- ============================================

-- ============================================
-- SCHEMA CREATION
-- ============================================

CREATE SCHEMA IF NOT EXISTS backup_prod_v1;

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================
-- DROP EXISTING BACKUP TABLES (for clean backup)
-- ============================================

DROP TABLE IF EXISTS backup_prod_v1.team_players CASCADE;
DROP TABLE IF EXISTS backup_prod_v1.vote_teams CASCADE;
DROP TABLE IF EXISTS backup_prod_v1.user_votes CASCADE;
DROP TABLE IF EXISTS backup_prod_v1.vote_options CASCADE;
DROP TABLE IF EXISTS backup_prod_v1.votes CASCADE;
DROP TABLE IF EXISTS backup_prod_v1.balls CASCADE;
DROP TABLE IF EXISTS backup_prod_v1.overs CASCADE;
DROP TABLE IF EXISTS backup_prod_v1.innings CASCADE;
DROP TABLE IF EXISTS backup_prod_v1.live_scoreboard CASCADE;
DROP TABLE IF EXISTS backup_prod_v1.matches CASCADE;
DROP TABLE IF EXISTS backup_prod_v1.series CASCADE;
DROP TABLE IF EXISTS backup_prod_v1.user_sessions CASCADE;
DROP TABLE IF EXISTS backup_prod_v1.oauth_states CASCADE;
DROP TABLE IF EXISTS backup_prod_v1.schema_version CASCADE;
DROP TABLE IF EXISTS backup_prod_v1.users CASCADE;

-- ============================================
-- COPY TABLE STRUCTURES FROM PRODUCTION
-- ============================================
-- This creates exact copies of production table structures

-- Copy users table structure
CREATE TABLE backup_prod_v1.users (LIKE prod_v1.users INCLUDING ALL);

-- Copy user_sessions table structure
CREATE TABLE backup_prod_v1.user_sessions (LIKE prod_v1.user_sessions INCLUDING ALL);

-- Copy oauth_states table structure
CREATE TABLE backup_prod_v1.oauth_states (LIKE prod_v1.oauth_states INCLUDING ALL);

-- Copy schema_version table structure
CREATE TABLE backup_prod_v1.schema_version (LIKE prod_v1.schema_version INCLUDING ALL);

-- Copy series table structure
CREATE TABLE backup_prod_v1.series (LIKE prod_v1.series INCLUDING ALL);

-- Copy matches table structure
CREATE TABLE backup_prod_v1.matches (LIKE prod_v1.matches INCLUDING ALL);

-- Copy live_scoreboard table structure
CREATE TABLE backup_prod_v1.live_scoreboard (LIKE prod_v1.live_scoreboard INCLUDING ALL);

-- Copy innings table structure
CREATE TABLE backup_prod_v1.innings (LIKE prod_v1.innings INCLUDING ALL);

-- Copy overs table structure
CREATE TABLE backup_prod_v1.overs (LIKE prod_v1.overs INCLUDING ALL);

-- Copy balls table structure
CREATE TABLE backup_prod_v1.balls (LIKE prod_v1.balls INCLUDING ALL);

-- Copy voting tables if they exist
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'prod_v1' AND table_name = 'votes') THEN
        EXECUTE 'CREATE TABLE backup_prod_v1.votes (LIKE prod_v1.votes INCLUDING ALL)';
    END IF;
    
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'prod_v1' AND table_name = 'vote_options') THEN
        EXECUTE 'CREATE TABLE backup_prod_v1.vote_options (LIKE prod_v1.vote_options INCLUDING ALL)';
    END IF;
    
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'prod_v1' AND table_name = 'user_votes') THEN
        EXECUTE 'CREATE TABLE backup_prod_v1.user_votes (LIKE prod_v1.user_votes INCLUDING ALL)';
    END IF;
    
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'prod_v1' AND table_name = 'vote_teams') THEN
        EXECUTE 'CREATE TABLE backup_prod_v1.vote_teams (LIKE prod_v1.vote_teams INCLUDING ALL)';
    END IF;
    
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'prod_v1' AND table_name = 'team_players') THEN
        EXECUTE 'CREATE TABLE backup_prod_v1.team_players (LIKE prod_v1.team_players INCLUDING ALL)';
    END IF;
END $$;

-- ============================================
-- RE-CREATE FOREIGN KEYS
-- ============================================
-- Foreign keys are not included in LIKE clause, so we recreate them

-- User sessions foreign key
ALTER TABLE backup_prod_v1.user_sessions 
    ADD CONSTRAINT user_sessions_user_id_fkey 
    FOREIGN KEY (user_id) REFERENCES backup_prod_v1.users(id) ON DELETE CASCADE;

-- Series foreign keys (if created_by column exists and is UUID)
DO $$
DECLARE
    created_by_type TEXT;
BEGIN
    SELECT data_type INTO created_by_type
    FROM information_schema.columns
    WHERE table_schema = 'prod_v1'
    AND table_name = 'series'
    AND column_name = 'created_by';
    
    IF created_by_type = 'uuid' THEN
        ALTER TABLE backup_prod_v1.series 
            ADD CONSTRAINT series_created_by_fkey 
            FOREIGN KEY (created_by) REFERENCES backup_prod_v1.users(id) ON DELETE SET NULL;
    END IF;
END $$;

-- Matches foreign keys
ALTER TABLE backup_prod_v1.matches 
    ADD CONSTRAINT matches_series_id_fkey 
    FOREIGN KEY (series_id) REFERENCES backup_prod_v1.series(id) ON DELETE CASCADE;

DO $$
DECLARE
    created_by_type TEXT;
BEGIN
    SELECT data_type INTO created_by_type
    FROM information_schema.columns
    WHERE table_schema = 'prod_v1'
    AND table_name = 'matches'
    AND column_name = 'created_by';
    
    IF created_by_type = 'uuid' THEN
        ALTER TABLE backup_prod_v1.matches 
            ADD CONSTRAINT matches_created_by_fkey 
            FOREIGN KEY (created_by) REFERENCES backup_prod_v1.users(id) ON DELETE SET NULL;
    END IF;
END $$;

-- Live scoreboard foreign keys
ALTER TABLE backup_prod_v1.live_scoreboard 
    ADD CONSTRAINT live_scoreboard_match_id_fkey 
    FOREIGN KEY (match_id) REFERENCES backup_prod_v1.matches(id) ON DELETE CASCADE;

-- Innings foreign keys
ALTER TABLE backup_prod_v1.innings 
    ADD CONSTRAINT innings_match_id_fkey 
    FOREIGN KEY (match_id) REFERENCES backup_prod_v1.matches(id) ON DELETE CASCADE;

-- Overs foreign keys
ALTER TABLE backup_prod_v1.overs 
    ADD CONSTRAINT overs_innings_id_fkey 
    FOREIGN KEY (innings_id) REFERENCES backup_prod_v1.innings(id) ON DELETE CASCADE;

-- Balls foreign keys
ALTER TABLE backup_prod_v1.balls 
    ADD CONSTRAINT balls_over_id_fkey 
    FOREIGN KEY (over_id) REFERENCES backup_prod_v1.overs(id) ON DELETE CASCADE;

-- Voting tables foreign keys (if they exist)
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'backup_prod_v1' AND table_name = 'votes') THEN
        EXECUTE 'ALTER TABLE backup_prod_v1.votes 
                 ADD CONSTRAINT votes_created_by_fkey 
                 FOREIGN KEY (created_by) REFERENCES backup_prod_v1.users(id) ON DELETE CASCADE';
    END IF;
    
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'backup_prod_v1' AND table_name = 'vote_options') THEN
        EXECUTE 'ALTER TABLE backup_prod_v1.vote_options 
                 ADD CONSTRAINT vote_options_vote_id_fkey 
                 FOREIGN KEY (vote_id) REFERENCES backup_prod_v1.votes(id) ON DELETE CASCADE';
    END IF;
    
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'backup_prod_v1' AND table_name = 'user_votes') THEN
        EXECUTE 'ALTER TABLE backup_prod_v1.user_votes 
                 ADD CONSTRAINT user_votes_vote_id_fkey 
                 FOREIGN KEY (vote_id) REFERENCES backup_prod_v1.votes(id) ON DELETE CASCADE';
        EXECUTE 'ALTER TABLE backup_prod_v1.user_votes 
                 ADD CONSTRAINT user_votes_user_id_fkey 
                 FOREIGN KEY (user_id) REFERENCES backup_prod_v1.users(id) ON DELETE CASCADE';
    END IF;
    
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'backup_prod_v1' AND table_name = 'vote_teams') THEN
        EXECUTE 'ALTER TABLE backup_prod_v1.vote_teams 
                 ADD CONSTRAINT vote_teams_vote_id_fkey 
                 FOREIGN KEY (vote_id) REFERENCES backup_prod_v1.votes(id) ON DELETE CASCADE';
        EXECUTE 'ALTER TABLE backup_prod_v1.vote_teams 
                 ADD CONSTRAINT vote_teams_captain_id_fkey 
                 FOREIGN KEY (captain_id) REFERENCES backup_prod_v1.users(id) ON DELETE CASCADE';
        EXECUTE 'ALTER TABLE backup_prod_v1.vote_teams 
                 ADD CONSTRAINT vote_teams_created_by_fkey 
                 FOREIGN KEY (created_by) REFERENCES backup_prod_v1.users(id) ON DELETE CASCADE';
    END IF;
    
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'backup_prod_v1' AND table_name = 'team_players') THEN
        EXECUTE 'ALTER TABLE backup_prod_v1.team_players 
                 ADD CONSTRAINT team_players_team_id_fkey 
                 FOREIGN KEY (team_id) REFERENCES backup_prod_v1.vote_teams(id) ON DELETE CASCADE';
        EXECUTE 'ALTER TABLE backup_prod_v1.team_players 
                 ADD CONSTRAINT team_players_user_id_fkey 
                 FOREIGN KEY (user_id) REFERENCES backup_prod_v1.users(id) ON DELETE CASCADE';
    END IF;
END $$;

-- ============================================
-- COPY DATA USING SELECT *
-- ============================================

-- Copy data in dependency order
INSERT INTO backup_prod_v1.users SELECT * FROM prod_v1.users;
INSERT INTO backup_prod_v1.oauth_states SELECT * FROM prod_v1.oauth_states;
INSERT INTO backup_prod_v1.schema_version SELECT * FROM prod_v1.schema_version;
INSERT INTO backup_prod_v1.user_sessions SELECT * FROM prod_v1.user_sessions;
INSERT INTO backup_prod_v1.series SELECT * FROM prod_v1.series;
INSERT INTO backup_prod_v1.matches SELECT * FROM prod_v1.matches;
INSERT INTO backup_prod_v1.live_scoreboard SELECT * FROM prod_v1.live_scoreboard;
INSERT INTO backup_prod_v1.innings SELECT * FROM prod_v1.innings;
INSERT INTO backup_prod_v1.overs SELECT * FROM prod_v1.overs;
INSERT INTO backup_prod_v1.balls SELECT * FROM prod_v1.balls;

-- Copy voting data if tables exist
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'backup_prod_v1' AND table_name = 'votes') THEN
        EXECUTE 'INSERT INTO backup_prod_v1.votes SELECT * FROM prod_v1.votes';
    END IF;
    
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'backup_prod_v1' AND table_name = 'vote_options') THEN
        EXECUTE 'INSERT INTO backup_prod_v1.vote_options SELECT * FROM prod_v1.vote_options';
    END IF;
    
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'backup_prod_v1' AND table_name = 'user_votes') THEN
        EXECUTE 'INSERT INTO backup_prod_v1.user_votes SELECT * FROM prod_v1.user_votes';
    END IF;
    
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'backup_prod_v1' AND table_name = 'vote_teams') THEN
        EXECUTE 'INSERT INTO backup_prod_v1.vote_teams SELECT * FROM prod_v1.vote_teams';
    END IF;
    
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'backup_prod_v1' AND table_name = 'team_players') THEN
        EXECUTE 'INSERT INTO backup_prod_v1.team_players SELECT * FROM prod_v1.team_players';
    END IF;
END $$;

-- ============================================
-- COPY FUNCTIONS FROM PRODUCTION
-- ============================================

-- Update trigger function
CREATE OR REPLACE FUNCTION backup_prod_v1.update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Copy cleanup functions
CREATE OR REPLACE FUNCTION backup_prod_v1.cleanup_expired_oauth_states()
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    DELETE FROM backup_prod_v1.oauth_states 
    WHERE expires_at < NOW() 
    AND (used_at IS NOT NULL OR expires_at < NOW() - INTERVAL '1 hour');
    
    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION backup_prod_v1.cleanup_expired_user_sessions()
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    DELETE FROM backup_prod_v1.user_sessions 
    WHERE expires_at < NOW();
    
    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION backup_prod_v1.get_user_statistics()
RETURNS TABLE (
    total_users BIGINT,
    active_sessions BIGINT,
    expired_sessions BIGINT,
    oauth_states_pending BIGINT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        (SELECT COUNT(*) FROM backup_prod_v1.users) as total_users,
        (SELECT COUNT(*) FROM backup_prod_v1.user_sessions WHERE expires_at > NOW()) as active_sessions,
        (SELECT COUNT(*) FROM backup_prod_v1.user_sessions WHERE expires_at <= NOW()) as expired_sessions,
        (SELECT COUNT(*) FROM backup_prod_v1.oauth_states WHERE expires_at > NOW() AND used_at IS NULL) as oauth_states_pending;
END;
$$ LANGUAGE plpgsql;

-- Copy vote team validation functions if tables exist
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'backup_prod_v1' AND table_name = 'vote_teams') THEN
        EXECUTE 'CREATE OR REPLACE FUNCTION backup_prod_v1.check_user_voted_in_vote(p_user_id UUID, p_vote_id UUID)
                 RETURNS BOOLEAN AS $func$
                 BEGIN
                     RETURN EXISTS (
                         SELECT 1 FROM backup_prod_v1.user_votes 
                         WHERE user_id = p_user_id AND vote_id = p_vote_id
                     );
                 END;
                 $func$ LANGUAGE plpgsql';
        
        EXECUTE 'CREATE OR REPLACE FUNCTION backup_prod_v1.check_captain_voted()
                 RETURNS TRIGGER AS $func$
                 BEGIN
                     IF NOT backup_prod_v1.check_user_voted_in_vote(NEW.captain_id, NEW.vote_id) THEN
                         RAISE EXCEPTION ''Captain must be one of the voters'';
                     END IF;
                     RETURN NEW;
                 END;
                 $func$ LANGUAGE plpgsql';
        
        EXECUTE 'CREATE OR REPLACE FUNCTION backup_prod_v1.check_team_player_voted()
                 RETURNS TRIGGER AS $func$
                 DECLARE
                     v_vote_id UUID;
                 BEGIN
                     SELECT vote_id INTO v_vote_id 
                     FROM backup_prod_v1.vote_teams 
                     WHERE id = NEW.team_id;
                     
                     IF NOT backup_prod_v1.check_user_voted_in_vote(NEW.user_id, v_vote_id) THEN
                         RAISE EXCEPTION ''Only voters can be assigned to teams'';
                     END IF;
                     RETURN NEW;
                 END;
                 $func$ LANGUAGE plpgsql';
        
        EXECUTE 'CREATE OR REPLACE FUNCTION backup_prod_v1.check_team_max_players()
                 RETURNS TRIGGER AS $func$
                 DECLARE
                     player_count INTEGER;
                 BEGIN
                     SELECT COUNT(*) INTO player_count 
                     FROM backup_prod_v1.team_players 
                     WHERE team_id = NEW.team_id;
                     
                     IF player_count >= 20 THEN
                         RAISE EXCEPTION ''Team cannot have more than 20 players'';
                     END IF;
                     RETURN NEW;
                 END;
                 $func$ LANGUAGE plpgsql';
        
        EXECUTE 'CREATE OR REPLACE FUNCTION backup_prod_v1.check_user_not_in_other_team()
                 RETURNS TRIGGER AS $func$
                 DECLARE
                     v_vote_id UUID;
                     other_team_count INTEGER;
                 BEGIN
                     SELECT vote_id INTO v_vote_id 
                     FROM backup_prod_v1.vote_teams 
                     WHERE id = NEW.team_id;
                     
                     SELECT COUNT(*) INTO other_team_count
                     FROM backup_prod_v1.team_players tp
                     JOIN backup_prod_v1.vote_teams vt ON tp.team_id = vt.id
                     WHERE tp.user_id = NEW.user_id 
                     AND vt.vote_id = v_vote_id
                     AND tp.team_id != NEW.team_id;
                     
                     IF other_team_count > 0 THEN
                         RAISE EXCEPTION ''User cannot be in both teams'';
                     END IF;
                     RETURN NEW;
                 END;
                 $func$ LANGUAGE plpgsql';
    END IF;
END $$;

-- ============================================
-- CREATE TRIGGERS
-- ============================================

-- Core table triggers
DROP TRIGGER IF EXISTS update_users_updated_at ON backup_prod_v1.users;
CREATE TRIGGER update_users_updated_at 
    BEFORE UPDATE ON backup_prod_v1.users
    FOR EACH ROW EXECUTE FUNCTION backup_prod_v1.update_updated_at_column();

DROP TRIGGER IF EXISTS update_user_sessions_updated_at ON backup_prod_v1.user_sessions;
CREATE TRIGGER update_user_sessions_updated_at 
    BEFORE UPDATE ON backup_prod_v1.user_sessions
    FOR EACH ROW EXECUTE FUNCTION backup_prod_v1.update_updated_at_column();

DROP TRIGGER IF EXISTS update_schema_version_updated_at ON backup_prod_v1.schema_version;
CREATE TRIGGER update_schema_version_updated_at 
    BEFORE UPDATE ON backup_prod_v1.schema_version
    FOR EACH ROW EXECUTE FUNCTION backup_prod_v1.update_updated_at_column();

DROP TRIGGER IF EXISTS update_series_updated_at ON backup_prod_v1.series;
CREATE TRIGGER update_series_updated_at 
    BEFORE UPDATE ON backup_prod_v1.series
    FOR EACH ROW EXECUTE FUNCTION backup_prod_v1.update_updated_at_column();

DROP TRIGGER IF EXISTS update_matches_updated_at ON backup_prod_v1.matches;
CREATE TRIGGER update_matches_updated_at 
    BEFORE UPDATE ON backup_prod_v1.matches
    FOR EACH ROW EXECUTE FUNCTION backup_prod_v1.update_updated_at_column();

DROP TRIGGER IF EXISTS update_live_scoreboard_updated_at ON backup_prod_v1.live_scoreboard;
CREATE TRIGGER update_live_scoreboard_updated_at 
    BEFORE UPDATE ON backup_prod_v1.live_scoreboard
    FOR EACH ROW EXECUTE FUNCTION backup_prod_v1.update_updated_at_column();

DROP TRIGGER IF EXISTS update_innings_updated_at ON backup_prod_v1.innings;
CREATE TRIGGER update_innings_updated_at 
    BEFORE UPDATE ON backup_prod_v1.innings
    FOR EACH ROW EXECUTE FUNCTION backup_prod_v1.update_updated_at_column();

DROP TRIGGER IF EXISTS update_overs_updated_at ON backup_prod_v1.overs;
CREATE TRIGGER update_overs_updated_at 
    BEFORE UPDATE ON backup_prod_v1.overs
    FOR EACH ROW EXECUTE FUNCTION backup_prod_v1.update_updated_at_column();

-- Voting table triggers (if tables exist)
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'backup_prod_v1' AND table_name = 'votes') THEN
        EXECUTE 'DROP TRIGGER IF EXISTS update_votes_updated_at ON backup_prod_v1.votes';
        EXECUTE 'CREATE TRIGGER update_votes_updated_at 
                 BEFORE UPDATE ON backup_prod_v1.votes
                 FOR EACH ROW EXECUTE FUNCTION backup_prod_v1.update_updated_at_column()';
    END IF;
    
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'backup_prod_v1' AND table_name = 'vote_options') THEN
        EXECUTE 'DROP TRIGGER IF EXISTS update_vote_options_updated_at ON backup_prod_v1.vote_options';
        EXECUTE 'CREATE TRIGGER update_vote_options_updated_at 
                 BEFORE UPDATE ON backup_prod_v1.vote_options
                 FOR EACH ROW EXECUTE FUNCTION backup_prod_v1.update_updated_at_column()';
    END IF;
    
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'backup_prod_v1' AND table_name = 'vote_teams') THEN
        EXECUTE 'DROP TRIGGER IF EXISTS update_vote_teams_updated_at ON backup_prod_v1.vote_teams';
        EXECUTE 'CREATE TRIGGER update_vote_teams_updated_at 
                 BEFORE UPDATE ON backup_prod_v1.vote_teams
                 FOR EACH ROW EXECUTE FUNCTION backup_prod_v1.update_updated_at_column()';
        
        -- Validation triggers
        EXECUTE 'DROP TRIGGER IF EXISTS validate_captain_voted ON backup_prod_v1.vote_teams';
        EXECUTE 'CREATE TRIGGER validate_captain_voted
                 BEFORE INSERT OR UPDATE ON backup_prod_v1.vote_teams
                 FOR EACH ROW EXECUTE FUNCTION backup_prod_v1.check_captain_voted()';
    END IF;
    
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'backup_prod_v1' AND table_name = 'team_players') THEN
        EXECUTE 'DROP TRIGGER IF EXISTS validate_team_player_voted ON backup_prod_v1.team_players';
        EXECUTE 'CREATE TRIGGER validate_team_player_voted
                 BEFORE INSERT ON backup_prod_v1.team_players
                 FOR EACH ROW EXECUTE FUNCTION backup_prod_v1.check_team_player_voted()';
        
        EXECUTE 'DROP TRIGGER IF EXISTS validate_team_max_players ON backup_prod_v1.team_players';
        EXECUTE 'CREATE TRIGGER validate_team_max_players
                 BEFORE INSERT ON backup_prod_v1.team_players
                 FOR EACH ROW EXECUTE FUNCTION backup_prod_v1.check_team_max_players()';
        
        EXECUTE 'DROP TRIGGER IF EXISTS validate_user_not_in_other_team ON backup_prod_v1.team_players';
        EXECUTE 'CREATE TRIGGER validate_user_not_in_other_team
                 BEFORE INSERT ON backup_prod_v1.team_players
                 FOR EACH ROW EXECUTE FUNCTION backup_prod_v1.check_user_not_in_other_team()';
    END IF;
END $$;

-- ============================================
-- GRANT PERMISSIONS
-- ============================================

DO $$
DECLARE
    role_name TEXT;
    roles TEXT[] := ARRAY['anon', 'authenticated', 'service_role', 'postgres', 'supabase_auth_admin', 'supabase_storage_admin'];
BEGIN
    FOREACH role_name IN ARRAY roles
    LOOP
        BEGIN
            EXECUTE format('GRANT USAGE ON SCHEMA backup_prod_v1 TO %I', role_name);
            EXECUTE format('GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA backup_prod_v1 TO %I', role_name);
            EXECUTE format('GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA backup_prod_v1 TO %I', role_name);
            EXECUTE format('GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA backup_prod_v1 TO %I', role_name);
            RAISE NOTICE 'Granted permissions to role: %', role_name;
        EXCEPTION
            WHEN OTHERS THEN
                RAISE NOTICE 'Could not grant permissions to role %: %', role_name, SQLERRM;
        END;
    END LOOP;
END $$;

ALTER DEFAULT PRIVILEGES IN SCHEMA backup_prod_v1 GRANT ALL ON TABLES TO anon, authenticated, service_role, postgres;
ALTER DEFAULT PRIVILEGES IN SCHEMA backup_prod_v1 GRANT ALL ON SEQUENCES TO anon, authenticated, service_role, postgres;
ALTER DEFAULT PRIVILEGES IN SCHEMA backup_prod_v1 GRANT EXECUTE ON FUNCTIONS TO anon, authenticated, service_role, postgres;

-- ============================================
-- VERIFICATION
-- ============================================

SELECT 'Backup completed successfully!' as status;
SELECT 'All tables copied from prod_v1 to backup_prod_v1 using CREATE TABLE LIKE and SELECT *' as method;

-- Count records
SELECT 'Record counts:' as info;
SELECT 'users' as table_name, COUNT(*) as records FROM backup_prod_v1.users
UNION ALL SELECT 'series', COUNT(*) FROM backup_prod_v1.series
UNION ALL SELECT 'matches', COUNT(*) FROM backup_prod_v1.matches
UNION ALL SELECT 'innings', COUNT(*) FROM backup_prod_v1.innings
UNION ALL SELECT 'overs', COUNT(*) FROM backup_prod_v1.overs
UNION ALL SELECT 'balls', COUNT(*) FROM backup_prod_v1.balls
ORDER BY table_name;

-- Show backup tables
SELECT 'Backup tables created:' as info;
SELECT table_name 
FROM information_schema.tables 
WHERE table_schema = 'backup_prod_v1' 
ORDER BY table_name;

