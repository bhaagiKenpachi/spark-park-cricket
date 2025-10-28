-- ============================================
-- SPARK PARK CRICKET - COMPLETE PRODUCTION BACKUP
-- ============================================
-- Complete backup from prod_v1 to backup_prod_v1
-- Version: 3.0.0 (Complete - All Migrations)
-- Date: 2025-01-27
-- Description: Complete backup of all prod_v1 tables including all migrations
-- Migrations included: 001-008 + complete_schema_prod_v1.sql
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
-- Dropped in reverse dependency order to handle foreign keys

DROP TABLE IF EXISTS backup_prod_v1.fall_of_wickets CASCADE;
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
-- This creates exact copies of production table structures using LIKE INCLUDING ALL
-- This includes: columns, constraints, indexes, defaults, storage parameters

-- Core authentication tables
CREATE TABLE backup_prod_v1.users (LIKE prod_v1.users INCLUDING ALL);
CREATE TABLE backup_prod_v1.user_sessions (LIKE prod_v1.user_sessions INCLUDING ALL);
CREATE TABLE backup_prod_v1.oauth_states (LIKE prod_v1.oauth_states INCLUDING ALL);
CREATE TABLE backup_prod_v1.schema_version (LIKE prod_v1.schema_version INCLUDING ALL);

-- Cricket core tables
CREATE TABLE backup_prod_v1.series (LIKE prod_v1.series INCLUDING ALL);
CREATE TABLE backup_prod_v1.matches (LIKE prod_v1.matches INCLUDING ALL);
CREATE TABLE backup_prod_v1.live_scoreboard (LIKE prod_v1.live_scoreboard INCLUDING ALL);
CREATE TABLE backup_prod_v1.innings (LIKE prod_v1.innings INCLUDING ALL);
CREATE TABLE backup_prod_v1.overs (LIKE prod_v1.overs INCLUDING ALL);
CREATE TABLE backup_prod_v1.balls (LIKE prod_v1.balls INCLUDING ALL);

-- Voting tables (if they exist)
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'prod_v1' AND table_name = 'votes') THEN
        CREATE TABLE backup_prod_v1.votes (LIKE prod_v1.votes INCLUDING ALL);
    END IF;
    
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'prod_v1' AND table_name = 'vote_options') THEN
        CREATE TABLE backup_prod_v1.vote_options (LIKE prod_v1.vote_options INCLUDING ALL);
    END IF;
    
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'prod_v1' AND table_name = 'user_votes') THEN
        CREATE TABLE backup_prod_v1.user_votes (LIKE prod_v1.user_votes INCLUDING ALL);
    END IF;
    
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'prod_v1' AND table_name = 'vote_teams') THEN
        CREATE TABLE backup_prod_v1.vote_teams (LIKE prod_v1.vote_teams INCLUDING ALL);
    END IF;
    
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'prod_v1' AND table_name = 'team_players') THEN
        CREATE TABLE backup_prod_v1.team_players (LIKE prod_v1.team_players INCLUDING ALL);
    END IF;
    
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'prod_v1' AND table_name = 'fall_of_wickets') THEN
        CREATE TABLE backup_prod_v1.fall_of_wickets (LIKE prod_v1.fall_of_wickets INCLUDING ALL);
    END IF;
END $$;

-- ============================================
-- FIX FOREIGN KEY CONSTRAINTS
-- ============================================
-- LIKE INCLUDING ALL copies foreign keys pointing to prod_v1, we need to update them

-- Drop foreign keys that reference prod_v1
DO $$
DECLARE
    r RECORD;
BEGIN
    -- Loop through all foreign key constraints in backup_prod_v1
    FOR r IN 
        SELECT 
            tc.table_schema,
            tc.constraint_name,
            tc.table_name
        FROM information_schema.table_constraints AS tc
        JOIN information_schema.key_column_usage AS kcu
            ON tc.constraint_name = kcu.constraint_name
        JOIN information_schema.constraint_column_usage AS ccu
            ON ccu.constraint_name = tc.constraint_name
        WHERE tc.constraint_type = 'FOREIGN KEY'
        AND tc.table_schema = 'backup_prod_v1'
        AND ccu.table_schema = 'prod_v1'
    LOOP
        EXECUTE format('ALTER TABLE backup_prod_v1.%I DROP CONSTRAINT IF EXISTS %I CASCADE',
            r.table_name, r.constraint_name);
    END LOOP;
END $$;

-- Recreate foreign keys pointing to backup_prod_v1 tables
-- User sessions
ALTER TABLE backup_prod_v1.user_sessions 
    ADD CONSTRAINT user_sessions_user_id_fkey 
    FOREIGN KEY (user_id) REFERENCES backup_prod_v1.users(id) ON DELETE CASCADE;

-- Series
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'backup_prod_v1' 
        AND table_name = 'series' 
        AND column_name = 'created_by'
        AND data_type = 'uuid'
    ) THEN
        ALTER TABLE backup_prod_v1.series 
            ADD CONSTRAINT series_created_by_fkey 
            FOREIGN KEY (created_by) REFERENCES backup_prod_v1.users(id) ON DELETE SET NULL;
    END IF;
END $$;

-- Matches
ALTER TABLE backup_prod_v1.matches 
    ADD CONSTRAINT matches_series_id_fkey 
    FOREIGN KEY (series_id) REFERENCES backup_prod_v1.series(id) ON DELETE CASCADE;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'backup_prod_v1' 
        AND table_name = 'matches' 
        AND column_name = 'created_by'
        AND data_type = 'uuid'
    ) THEN
        ALTER TABLE backup_prod_v1.matches 
            ADD CONSTRAINT matches_created_by_fkey 
            FOREIGN KEY (created_by) REFERENCES backup_prod_v1.users(id) ON DELETE SET NULL;
    END IF;
END $$;

-- Live scoreboard
ALTER TABLE backup_prod_v1.live_scoreboard 
    ADD CONSTRAINT live_scoreboard_match_id_fkey 
    FOREIGN KEY (match_id) REFERENCES backup_prod_v1.matches(id) ON DELETE CASCADE;

-- Innings
ALTER TABLE backup_prod_v1.innings 
    ADD CONSTRAINT innings_match_id_fkey 
    FOREIGN KEY (match_id) REFERENCES backup_prod_v1.matches(id) ON DELETE CASCADE;

-- Overs
ALTER TABLE backup_prod_v1.overs 
    ADD CONSTRAINT overs_innings_id_fkey 
    FOREIGN KEY (innings_id) REFERENCES backup_prod_v1.innings(id) ON DELETE CASCADE;

-- Balls
ALTER TABLE backup_prod_v1.balls 
    ADD CONSTRAINT balls_over_id_fkey 
    FOREIGN KEY (over_id) REFERENCES backup_prod_v1.overs(id) ON DELETE CASCADE;

-- Fall of wickets
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'backup_prod_v1' AND table_name = 'fall_of_wickets') THEN
        ALTER TABLE backup_prod_v1.fall_of_wickets 
            ADD CONSTRAINT fall_of_wickets_match_id_fkey 
            FOREIGN KEY (match_id) REFERENCES backup_prod_v1.matches(id) ON DELETE CASCADE;
        
        ALTER TABLE backup_prod_v1.fall_of_wickets 
            ADD CONSTRAINT fall_of_wickets_innings_id_fkey 
            FOREIGN KEY (innings_id) REFERENCES backup_prod_v1.innings(id) ON DELETE CASCADE;
        
        ALTER TABLE backup_prod_v1.fall_of_wickets 
            ADD CONSTRAINT fall_of_wickets_over_id_fkey 
            FOREIGN KEY (over_id) REFERENCES backup_prod_v1.overs(id) ON DELETE CASCADE;
        
        ALTER TABLE backup_prod_v1.fall_of_wickets 
            ADD CONSTRAINT fall_of_wickets_ball_id_fkey 
            FOREIGN KEY (ball_id) REFERENCES backup_prod_v1.balls(id) ON DELETE CASCADE;
    END IF;
END $$;

-- Voting tables foreign keys (if they exist)
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'backup_prod_v1' AND table_name = 'votes') THEN
        ALTER TABLE backup_prod_v1.votes 
            ADD CONSTRAINT votes_created_by_fkey 
            FOREIGN KEY (created_by) REFERENCES backup_prod_v1.users(id) ON DELETE CASCADE;
    END IF;
    
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'backup_prod_v1' AND table_name = 'vote_options') THEN
        ALTER TABLE backup_prod_v1.vote_options 
            ADD CONSTRAINT vote_options_vote_id_fkey 
            FOREIGN KEY (vote_id) REFERENCES backup_prod_v1.votes(id) ON DELETE CASCADE;
    END IF;
    
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'backup_prod_v1' AND table_name = 'user_votes') THEN
        ALTER TABLE backup_prod_v1.user_votes 
            ADD CONSTRAINT user_votes_vote_id_fkey 
            FOREIGN KEY (vote_id) REFERENCES backup_prod_v1.votes(id) ON DELETE CASCADE;
        
        ALTER TABLE backup_prod_v1.user_votes 
            ADD CONSTRAINT user_votes_user_id_fkey 
            FOREIGN KEY (user_id) REFERENCES backup_prod_v1.users(id) ON DELETE CASCADE;
    END IF;
    
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'backup_prod_v1' AND table_name = 'vote_teams') THEN
        ALTER TABLE backup_prod_v1.vote_teams 
            ADD CONSTRAINT vote_teams_vote_id_fkey 
            FOREIGN KEY (vote_id) REFERENCES backup_prod_v1.votes(id) ON DELETE CASCADE;
        
        ALTER TABLE backup_prod_v1.vote_teams 
            ADD CONSTRAINT vote_teams_captain_id_fkey 
            FOREIGN KEY (captain_id) REFERENCES backup_prod_v1.users(id) ON DELETE CASCADE;
        
        ALTER TABLE backup_prod_v1.vote_teams 
            ADD CONSTRAINT vote_teams_created_by_fkey 
            FOREIGN KEY (created_by) REFERENCES backup_prod_v1.users(id) ON DELETE CASCADE;
    END IF;
    
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'backup_prod_v1' AND table_name = 'team_players') THEN
        ALTER TABLE backup_prod_v1.team_players 
            ADD CONSTRAINT team_players_team_id_fkey 
            FOREIGN KEY (team_id) REFERENCES backup_prod_v1.vote_teams(id) ON DELETE CASCADE;
        
        ALTER TABLE backup_prod_v1.team_players 
            ADD CONSTRAINT team_players_user_id_fkey 
            FOREIGN KEY (user_id) REFERENCES backup_prod_v1.users(id) ON DELETE CASCADE;
    END IF;
END $$;

-- ============================================
-- COPY DATA (IN DEPENDENCY ORDER)
-- ============================================
-- Copy data in correct dependency order to respect foreign key constraints

-- Level 1: No dependencies
INSERT INTO backup_prod_v1.users SELECT * FROM prod_v1.users;
INSERT INTO backup_prod_v1.oauth_states SELECT * FROM prod_v1.oauth_states;
INSERT INTO backup_prod_v1.schema_version SELECT * FROM prod_v1.schema_version;

-- Level 2: Depends on users
INSERT INTO backup_prod_v1.user_sessions SELECT * FROM prod_v1.user_sessions;
INSERT INTO backup_prod_v1.series SELECT * FROM prod_v1.series;

-- Level 3: Depends on series/users
INSERT INTO backup_prod_v1.matches SELECT * FROM prod_v1.matches;

-- Level 4: Depends on matches
INSERT INTO backup_prod_v1.live_scoreboard SELECT * FROM prod_v1.live_scoreboard;
INSERT INTO backup_prod_v1.innings SELECT * FROM prod_v1.innings;

-- Level 5: Depends on innings
INSERT INTO backup_prod_v1.overs SELECT * FROM prod_v1.overs;

-- Level 6: Depends on overs
INSERT INTO backup_prod_v1.balls SELECT * FROM prod_v1.balls;

-- Level 7: Depends on matches/innings/overs/balls
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'backup_prod_v1' AND table_name = 'fall_of_wickets') THEN
        INSERT INTO backup_prod_v1.fall_of_wickets SELECT * FROM prod_v1.fall_of_wickets;
    END IF;
END $$;

-- Voting tables (depends on users/votes)
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'backup_prod_v1' AND table_name = 'votes') THEN
        INSERT INTO backup_prod_v1.votes SELECT * FROM prod_v1.votes;
    END IF;
    
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'backup_prod_v1' AND table_name = 'vote_options') THEN
        INSERT INTO backup_prod_v1.vote_options SELECT * FROM prod_v1.vote_options;
    END IF;
    
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'backup_prod_v1' AND table_name = 'user_votes') THEN
        INSERT INTO backup_prod_v1.user_votes SELECT * FROM prod_v1.user_votes;
    END IF;
    
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'backup_prod_v1' AND table_name = 'vote_teams') THEN
        INSERT INTO backup_prod_v1.vote_teams SELECT * FROM prod_v1.vote_teams;
    END IF;
    
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'backup_prod_v1' AND table_name = 'team_players') THEN
        INSERT INTO backup_prod_v1.team_players SELECT * FROM prod_v1.team_players;
    END IF;
END $$;

-- ============================================
-- COPY FUNCTIONS FROM PRODUCTION
-- ============================================

-- Update trigger function (global function)
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Backup-specific cleanup functions
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
DROP TRIGGER IF EXISTS update_backup_prod_v1_users_updated_at ON backup_prod_v1.users;
CREATE TRIGGER update_backup_prod_v1_users_updated_at 
    BEFORE UPDATE ON backup_prod_v1.users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_backup_prod_v1_user_sessions_updated_at ON backup_prod_v1.user_sessions;
CREATE TRIGGER update_backup_prod_v1_user_sessions_updated_at 
    BEFORE UPDATE ON backup_prod_v1.user_sessions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_backup_prod_v1_schema_version_updated_at ON backup_prod_v1.schema_version;
CREATE TRIGGER update_backup_prod_v1_schema_version_updated_at 
    BEFORE UPDATE ON backup_prod_v1.schema_version
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_backup_prod_v1_series_updated_at ON backup_prod_v1.series;
CREATE TRIGGER update_backup_prod_v1_series_updated_at 
    BEFORE UPDATE ON backup_prod_v1.series
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_backup_prod_v1_matches_updated_at ON backup_prod_v1.matches;
CREATE TRIGGER update_backup_prod_v1_matches_updated_at 
    BEFORE UPDATE ON backup_prod_v1.matches
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_backup_prod_v1_live_scoreboard_updated_at ON backup_prod_v1.live_scoreboard;
CREATE TRIGGER update_backup_prod_v1_live_scoreboard_updated_at 
    BEFORE UPDATE ON backup_prod_v1.live_scoreboard
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_backup_prod_v1_innings_updated_at ON backup_prod_v1.innings;
CREATE TRIGGER update_backup_prod_v1_innings_updated_at 
    BEFORE UPDATE ON backup_prod_v1.innings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_backup_prod_v1_overs_updated_at ON backup_prod_v1.overs;
CREATE TRIGGER update_backup_prod_v1_overs_updated_at 
    BEFORE UPDATE ON backup_prod_v1.overs
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Voting table triggers (if tables exist)
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'backup_prod_v1' AND table_name = 'votes') THEN
        DROP TRIGGER IF EXISTS update_backup_prod_v1_votes_updated_at ON backup_prod_v1.votes;
        CREATE TRIGGER update_backup_prod_v1_votes_updated_at 
            BEFORE UPDATE ON backup_prod_v1.votes
            FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
    END IF;
    
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'backup_prod_v1' AND table_name = 'vote_options') THEN
        DROP TRIGGER IF EXISTS update_backup_prod_v1_vote_options_updated_at ON backup_prod_v1.vote_options;
        CREATE TRIGGER update_backup_prod_v1_vote_options_updated_at 
            BEFORE UPDATE ON backup_prod_v1.vote_options
            FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
    END IF;
    
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'backup_prod_v1' AND table_name = 'vote_teams') THEN
        DROP TRIGGER IF EXISTS update_backup_prod_v1_vote_teams_updated_at ON backup_prod_v1.vote_teams;
        CREATE TRIGGER update_backup_prod_v1_vote_teams_updated_at 
            BEFORE UPDATE ON backup_prod_v1.vote_teams
            FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
        
        -- Validation triggers
        DROP TRIGGER IF EXISTS validate_backup_prod_v1_captain_voted ON backup_prod_v1.vote_teams;
        CREATE TRIGGER validate_backup_prod_v1_captain_voted
            BEFORE INSERT OR UPDATE ON backup_prod_v1.vote_teams
            FOR EACH ROW EXECUTE FUNCTION backup_prod_v1.check_captain_voted();
    END IF;
    
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'backup_prod_v1' AND table_name = 'team_players') THEN
        DROP TRIGGER IF EXISTS validate_backup_prod_v1_team_player_voted ON backup_prod_v1.team_players;
        CREATE TRIGGER validate_backup_prod_v1_team_player_voted
            BEFORE INSERT ON backup_prod_v1.team_players
            FOR EACH ROW EXECUTE FUNCTION backup_prod_v1.check_team_player_voted();
        
        DROP TRIGGER IF EXISTS validate_backup_prod_v1_team_max_players ON backup_prod_v1.team_players;
        CREATE TRIGGER validate_backup_prod_v1_team_max_players
            BEFORE INSERT ON backup_prod_v1.team_players
            FOR EACH ROW EXECUTE FUNCTION backup_prod_v1.check_team_max_players();
        
        DROP TRIGGER IF EXISTS validate_backup_prod_v1_user_not_in_other_team ON backup_prod_v1.team_players;
        CREATE TRIGGER validate_backup_prod_v1_user_not_in_other_team
            BEFORE INSERT ON backup_prod_v1.team_players
            FOR EACH ROW EXECUTE FUNCTION backup_prod_v1.check_user_not_in_other_team();
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

-- Set default privileges for future tables and sequences
ALTER DEFAULT PRIVILEGES IN SCHEMA backup_prod_v1 GRANT ALL ON TABLES TO anon, authenticated, service_role, postgres;
ALTER DEFAULT PRIVILEGES IN SCHEMA backup_prod_v1 GRANT ALL ON SEQUENCES TO anon, authenticated, service_role, postgres;
ALTER DEFAULT PRIVILEGES IN SCHEMA backup_prod_v1 GRANT EXECUTE ON FUNCTIONS TO anon, authenticated, service_role, postgres;

-- ============================================
-- VERIFICATION AND STATISTICS
-- ============================================

SELECT 'Backup completed successfully!' as status;
SELECT 'All tables copied from prod_v1 to backup_prod_v1 using CREATE TABLE LIKE INCLUDING ALL' as method;
SELECT 'This backup includes all migrations: 001-008 + complete_schema_prod_v1.sql' as migrations_included;

-- Show all backup tables
SELECT 'Backup tables created:' as info;
SELECT table_name 
FROM information_schema.tables 
WHERE table_schema = 'backup_prod_v1' 
ORDER BY table_name;

-- Count records in each table for verification
SELECT 'Record counts in backup schema:' as info;
SELECT 'users' as table_name, COUNT(*) as records FROM backup_prod_v1.users
UNION ALL SELECT 'user_sessions', COUNT(*) FROM backup_prod_v1.user_sessions
UNION ALL SELECT 'oauth_states', COUNT(*) FROM backup_prod_v1.oauth_states
UNION ALL SELECT 'schema_version', COUNT(*) FROM backup_prod_v1.schema_version
UNION ALL SELECT 'series', COUNT(*) FROM backup_prod_v1.series
UNION ALL SELECT 'matches', COUNT(*) FROM backup_prod_v1.matches
UNION ALL SELECT 'live_scoreboard', COUNT(*) FROM backup_prod_v1.live_scoreboard
UNION ALL SELECT 'innings', COUNT(*) FROM backup_prod_v1.innings
UNION ALL SELECT 'overs', COUNT(*) FROM backup_prod_v1.overs
UNION ALL SELECT 'balls', COUNT(*) FROM backup_prod_v1.balls
UNION ALL SELECT 'votes', COUNT(*) FROM backup_prod_v1.votes WHERE EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'backup_prod_v1' AND table_name = 'votes')
UNION ALL SELECT 'vote_options', COUNT(*) FROM backup_prod_v1.vote_options WHERE EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'backup_prod_v1' AND table_name = 'vote_options')
UNION ALL SELECT 'user_votes', COUNT(*) FROM backup_prod_v1.user_votes WHERE EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'backup_prod_v1' AND table_name = 'user_votes')
UNION ALL SELECT 'vote_teams', COUNT(*) FROM backup_prod_v1.vote_teams WHERE EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'backup_prod_v1' AND table_name = 'vote_teams')
UNION ALL SELECT 'team_players', COUNT(*) FROM backup_prod_v1.team_players WHERE EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'backup_prod_v1' AND table_name = 'team_players')
UNION ALL SELECT 'fall_of_wickets', COUNT(*) FROM backup_prod_v1.fall_of_wickets WHERE EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'backup_prod_v1' AND table_name = 'fall_of_wickets')
ORDER BY table_name;

-- Compare record counts between original and backup
SELECT 'Record count comparison (prod_v1 vs backup_prod_v1):' as info;
SELECT 
    'users' as table_name,
    (SELECT COUNT(*) FROM prod_v1.users) as prod_v1_count,
    (SELECT COUNT(*) FROM backup_prod_v1.users) as backup_count,
    CASE 
        WHEN (SELECT COUNT(*) FROM prod_v1.users) = (SELECT COUNT(*) FROM backup_prod_v1.users) 
        THEN 'MATCH' 
        ELSE 'MISMATCH' 
    END as status
UNION ALL
SELECT 
    'matches',
    (SELECT COUNT(*) FROM prod_v1.matches),
    (SELECT COUNT(*) FROM backup_prod_v1.matches),
    CASE 
        WHEN (SELECT COUNT(*) FROM prod_v1.matches) = (SELECT COUNT(*) FROM backup_prod_v1.matches) 
        THEN 'MATCH' 
        ELSE 'MISMATCH' 
    END
UNION ALL
SELECT 
    'balls',
    (SELECT COUNT(*) FROM prod_v1.balls),
    (SELECT COUNT(*) FROM backup_prod_v1.balls),
    CASE 
        WHEN (SELECT COUNT(*) FROM prod_v1.balls) = (SELECT COUNT(*) FROM backup_prod_v1.balls) 
        THEN 'MATCH' 
        ELSE 'MISMATCH' 
    END
ORDER BY table_name;

SELECT 'Complete backup of prod_v1 schema finished successfully!' as final_status;

