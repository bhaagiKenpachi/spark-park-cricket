-- ============================================
-- SPARK PARK CRICKET - BACKUP SCHEMA CREATION
-- ============================================
-- Comprehensive Data Backup from prod_v1 to backup_prod_v1
-- Version: 1.0.1 (Fixed - Correct Schema References)
-- Date: 2025-01-27
-- Description: Complete backup of prod_v1 schema including all data
-- Features: Schema creation, data copying, indexes, triggers, and functions
-- ============================================

-- ============================================
-- SCHEMA CREATION
-- ============================================

-- Create backup schema if it doesn't exist
CREATE SCHEMA IF NOT EXISTS backup_prod_v1;

-- ============================================
-- TABLE STRUCTURE CREATION (WITH CORRECT REFERENCES)
-- ============================================

-- Create users table (backup_prod_v1 schema)
CREATE TABLE IF NOT EXISTS backup_prod_v1.users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    google_id VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    picture TEXT,
    email_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_login_at TIMESTAMP WITH TIME ZONE
);

-- Create user_sessions table (backup_prod_v1 schema)
CREATE TABLE IF NOT EXISTS backup_prod_v1.user_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES backup_prod_v1.users(id) ON DELETE CASCADE,
    session_id VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create oauth_states table (backup_prod_v1 schema)
CREATE TABLE IF NOT EXISTS backup_prod_v1.oauth_states (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    state VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    used_at TIMESTAMP WITH TIME ZONE NULL
);

-- Create schema_version table (backup_prod_v1 schema)
CREATE TABLE IF NOT EXISTS backup_prod_v1.schema_version (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    version VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    applied_by VARCHAR(255),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create series table (backup_prod_v1 schema)
CREATE TABLE IF NOT EXISTS backup_prod_v1.series (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    start_date TIMESTAMP WITH TIME ZONE NOT NULL,
    end_date TIMESTAMP WITH TIME ZONE NOT NULL,
    created_by UUID REFERENCES backup_prod_v1.users(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create matches table (backup_prod_v1 schema)
CREATE TABLE IF NOT EXISTS backup_prod_v1.matches (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    series_id UUID REFERENCES backup_prod_v1.series(id) ON DELETE CASCADE,
    match_number INTEGER NOT NULL,
    date TIMESTAMP WITH TIME ZONE NOT NULL,
    status VARCHAR(20) DEFAULT 'live' CHECK (status IN ('live', 'completed', 'cancelled')),
    team_a_player_count INTEGER NOT NULL DEFAULT 11 CHECK (team_a_player_count >= 1 AND team_a_player_count <= 20),
    team_b_player_count INTEGER NOT NULL DEFAULT 11 CHECK (team_b_player_count >= 1 AND team_b_player_count <= 20),
    total_overs INTEGER NOT NULL DEFAULT 20 CHECK (total_overs >= 1 AND total_overs <= 20),
    toss_winner VARCHAR(1) NOT NULL CHECK (toss_winner IN ('A', 'B')),
    toss_type VARCHAR(1) NOT NULL CHECK (toss_type IN ('H', 'T')),
    batting_team VARCHAR(1) NOT NULL DEFAULT 'A' CHECK (batting_team IN ('A', 'B')),
    created_by UUID REFERENCES backup_prod_v1.users(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create live_scoreboard table (backup_prod_v1 schema)
CREATE TABLE IF NOT EXISTS backup_prod_v1.live_scoreboard (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    match_id UUID REFERENCES backup_prod_v1.matches(id) ON DELETE CASCADE,
    batting_team VARCHAR(1) NOT NULL CHECK (batting_team IN ('A', 'B')),
    score INTEGER DEFAULT 0,
    wickets INTEGER DEFAULT 0,
    overs DECIMAL(4,1) DEFAULT 0.0,
    balls INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create innings table (backup_prod_v1 schema)
CREATE TABLE IF NOT EXISTS backup_prod_v1.innings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    match_id UUID REFERENCES backup_prod_v1.matches(id) ON DELETE CASCADE,
    innings_number INTEGER NOT NULL CHECK (innings_number IN (1, 2)),
    batting_team VARCHAR(1) NOT NULL CHECK (batting_team IN ('A', 'B')),
    total_runs INTEGER DEFAULT 0 CHECK (total_runs >= 0),
    total_wickets INTEGER DEFAULT 0 CHECK (total_wickets >= 0 AND total_wickets <= 10),
    total_overs DECIMAL(4,1) DEFAULT 0.0 CHECK (total_overs >= 0),
    total_balls INTEGER DEFAULT 0 CHECK (total_balls >= 0),
    status VARCHAR(20) DEFAULT 'in_progress' CHECK (status IN ('in_progress', 'completed')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Ensure only one innings per match per innings number
    UNIQUE(match_id, innings_number)
);

-- Create overs table (backup_prod_v1 schema)
CREATE TABLE IF NOT EXISTS backup_prod_v1.overs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    innings_id UUID REFERENCES backup_prod_v1.innings(id) ON DELETE CASCADE,
    over_number INTEGER NOT NULL CHECK (over_number >= 1),
    total_runs INTEGER DEFAULT 0 CHECK (total_runs >= 0),
    total_balls INTEGER DEFAULT 0 CHECK (total_balls >= 0 AND total_balls <= 6),
    total_wickets INTEGER DEFAULT 0 CHECK (total_wickets >= 0),
    status VARCHAR(20) DEFAULT 'in_progress' CHECK (status IN ('in_progress', 'completed')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Ensure only one over per innings per over number
    UNIQUE(innings_id, over_number)
);

-- Create balls table (backup_prod_v1 schema)
CREATE TABLE IF NOT EXISTS backup_prod_v1.balls (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    over_id UUID REFERENCES backup_prod_v1.overs(id) ON DELETE CASCADE,
    ball_number INTEGER NOT NULL CHECK (ball_number >= 1 AND ball_number <= 20),
    ball_type VARCHAR(20) NOT NULL CHECK (ball_type IN ('good', 'wide', 'no_ball', 'dead_ball')),
    run_type VARCHAR(2) NOT NULL CHECK (run_type IN ('0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'NB', 'WD', 'LB', 'WC')),
    runs INTEGER DEFAULT 0 CHECK (runs >= 0),
    byes INTEGER DEFAULT 0 CHECK (byes >= 0 AND byes <= 6),
    is_wicket BOOLEAN DEFAULT FALSE,
    wicket_type VARCHAR(20) CHECK (wicket_type IN ('bowled', 'caught', 'lbw', 'run_out', 'stumped', 'hit_wicket')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Ensure only one ball per over per ball number
    UNIQUE(over_id, ball_number),
    
    -- Ensure wicket_type is NULL when is_wicket is false
    CONSTRAINT backup_prod_v1_balls_wicket_type_check CHECK (
        (is_wicket = true AND wicket_type IS NOT NULL) OR 
        (is_wicket = false AND wicket_type IS NULL)
    )
);

-- ============================================
-- DATA COPYING (WITH PROPER DEPENDENCY ORDER)
-- ============================================

-- Clear existing data in backup tables (in reverse dependency order)
TRUNCATE TABLE backup_prod_v1.balls CASCADE;
TRUNCATE TABLE backup_prod_v1.overs CASCADE;
TRUNCATE TABLE backup_prod_v1.innings CASCADE;
TRUNCATE TABLE backup_prod_v1.live_scoreboard CASCADE;
TRUNCATE TABLE backup_prod_v1.matches CASCADE;
TRUNCATE TABLE backup_prod_v1.series CASCADE;
TRUNCATE TABLE backup_prod_v1.user_sessions CASCADE;
TRUNCATE TABLE backup_prod_v1.oauth_states CASCADE;
TRUNCATE TABLE backup_prod_v1.schema_version CASCADE;
TRUNCATE TABLE backup_prod_v1.users CASCADE;

-- Copy data in dependency order
-- 1. Users (no dependencies)
INSERT INTO backup_prod_v1.users 
SELECT * FROM prod_v1.users;

-- 2. OAuth states (no dependencies)
INSERT INTO backup_prod_v1.oauth_states 
SELECT * FROM prod_v1.oauth_states;

-- 3. Schema version (no dependencies)
INSERT INTO backup_prod_v1.schema_version 
SELECT * FROM prod_v1.schema_version;

-- 4. User sessions (depends on users)
INSERT INTO backup_prod_v1.user_sessions 
SELECT * FROM prod_v1.user_sessions;

-- 5. Series (depends on users)
INSERT INTO backup_prod_v1.series 
SELECT * FROM prod_v1.series;

-- 6. Matches (depends on series and users)
INSERT INTO backup_prod_v1.matches 
SELECT * FROM prod_v1.matches;

-- 7. Live scoreboard (depends on matches)
INSERT INTO backup_prod_v1.live_scoreboard 
SELECT * FROM prod_v1.live_scoreboard;

-- 8. Innings (depends on matches)
INSERT INTO backup_prod_v1.innings 
SELECT * FROM prod_v1.innings;

-- 9. Overs (depends on innings)
INSERT INTO backup_prod_v1.overs 
SELECT * FROM prod_v1.overs;

-- 10. Balls (depends on overs)
INSERT INTO backup_prod_v1.balls 
SELECT * FROM prod_v1.balls;

-- ============================================
-- INDEXES CREATION
-- ============================================

-- User authentication indexes (backup_prod_v1 schema)
CREATE INDEX IF NOT EXISTS idx_backup_prod_v1_users_google_id ON backup_prod_v1.users(google_id);
CREATE INDEX IF NOT EXISTS idx_backup_prod_v1_users_email ON backup_prod_v1.users(email);
CREATE INDEX IF NOT EXISTS idx_backup_prod_v1_user_sessions_session_id ON backup_prod_v1.user_sessions(session_id);
CREATE INDEX IF NOT EXISTS idx_backup_prod_v1_user_sessions_user_id ON backup_prod_v1.user_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_backup_prod_v1_user_sessions_expires_at ON backup_prod_v1.user_sessions(expires_at);
CREATE INDEX IF NOT EXISTS idx_backup_prod_v1_oauth_states_state ON backup_prod_v1.oauth_states(state);
CREATE INDEX IF NOT EXISTS idx_backup_prod_v1_oauth_states_expires_at ON backup_prod_v1.oauth_states(expires_at);
CREATE INDEX IF NOT EXISTS idx_backup_prod_v1_oauth_states_used_at ON backup_prod_v1.oauth_states(used_at);
CREATE INDEX IF NOT EXISTS idx_backup_prod_v1_schema_version_version ON backup_prod_v1.schema_version(version);

-- Cricket schema indexes (backup_prod_v1)
-- Series indexes
CREATE INDEX IF NOT EXISTS idx_backup_prod_v1_series_start_date ON backup_prod_v1.series(start_date);
CREATE INDEX IF NOT EXISTS idx_backup_prod_v1_series_end_date ON backup_prod_v1.series(end_date);
CREATE INDEX IF NOT EXISTS idx_backup_prod_v1_series_created_by ON backup_prod_v1.series(created_by);

-- Matches indexes
CREATE INDEX IF NOT EXISTS idx_backup_prod_v1_matches_series_id ON backup_prod_v1.matches(series_id);
CREATE INDEX IF NOT EXISTS idx_backup_prod_v1_matches_status ON backup_prod_v1.matches(status);
CREATE INDEX IF NOT EXISTS idx_backup_prod_v1_matches_date ON backup_prod_v1.matches(date);
CREATE INDEX IF NOT EXISTS idx_backup_prod_v1_matches_toss_winner ON backup_prod_v1.matches(toss_winner);
CREATE INDEX IF NOT EXISTS idx_backup_prod_v1_matches_batting_team ON backup_prod_v1.matches(batting_team);
CREATE INDEX IF NOT EXISTS idx_backup_prod_v1_matches_created_by ON backup_prod_v1.matches(created_by);

-- Live scoreboard indexes
CREATE INDEX IF NOT EXISTS idx_backup_prod_v1_live_scoreboard_match_id ON backup_prod_v1.live_scoreboard(match_id);
CREATE INDEX IF NOT EXISTS idx_backup_prod_v1_live_scoreboard_batting_team ON backup_prod_v1.live_scoreboard(batting_team);

-- Innings indexes
CREATE INDEX IF NOT EXISTS idx_backup_prod_v1_innings_match_id ON backup_prod_v1.innings(match_id);
CREATE INDEX IF NOT EXISTS idx_backup_prod_v1_innings_batting_team ON backup_prod_v1.innings(batting_team);
CREATE INDEX IF NOT EXISTS idx_backup_prod_v1_innings_status ON backup_prod_v1.innings(status);

-- Overs indexes
CREATE INDEX IF NOT EXISTS idx_backup_prod_v1_overs_innings_id ON backup_prod_v1.overs(innings_id);
CREATE INDEX IF NOT EXISTS idx_backup_prod_v1_overs_status ON backup_prod_v1.overs(status);

-- Balls indexes
CREATE INDEX IF NOT EXISTS idx_backup_prod_v1_balls_over_id ON backup_prod_v1.balls(over_id);
CREATE INDEX IF NOT EXISTS idx_backup_prod_v1_balls_run_type ON backup_prod_v1.balls(run_type);
CREATE INDEX IF NOT EXISTS idx_backup_prod_v1_balls_is_wicket ON backup_prod_v1.balls(is_wicket);

-- ============================================
-- FUNCTIONS CREATION
-- ============================================

-- Create function to automatically update updated_at timestamp (if not exists)
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create backup-specific cleanup functions
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
$$ language 'plpgsql';

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
$$ language 'plpgsql';

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
$$ language 'plpgsql';

-- ============================================
-- TRIGGERS CREATION
-- ============================================

-- Create triggers for backup_prod_v1 schema
DROP TRIGGER IF EXISTS update_backup_prod_v1_users_updated_at ON backup_prod_v1.users;
CREATE TRIGGER update_backup_prod_v1_users_updated_at BEFORE UPDATE ON backup_prod_v1.users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_backup_prod_v1_user_sessions_updated_at ON backup_prod_v1.user_sessions;
CREATE TRIGGER update_backup_prod_v1_user_sessions_updated_at BEFORE UPDATE ON backup_prod_v1.user_sessions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_backup_prod_v1_schema_version_updated_at ON backup_prod_v1.schema_version;
CREATE TRIGGER update_backup_prod_v1_schema_version_updated_at BEFORE UPDATE ON backup_prod_v1.schema_version
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_backup_prod_v1_series_updated_at ON backup_prod_v1.series;
CREATE TRIGGER update_backup_prod_v1_series_updated_at BEFORE UPDATE ON backup_prod_v1.series
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_backup_prod_v1_matches_updated_at ON backup_prod_v1.matches;
CREATE TRIGGER update_backup_prod_v1_matches_updated_at BEFORE UPDATE ON backup_prod_v1.matches
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_backup_prod_v1_live_scoreboard_updated_at ON backup_prod_v1.live_scoreboard;
CREATE TRIGGER update_backup_prod_v1_live_scoreboard_updated_at BEFORE UPDATE ON backup_prod_v1.live_scoreboard
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_backup_prod_v1_innings_updated_at ON backup_prod_v1.innings;
CREATE TRIGGER update_backup_prod_v1_innings_updated_at BEFORE UPDATE ON backup_prod_v1.innings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_backup_prod_v1_overs_updated_at ON backup_prod_v1.overs;
CREATE TRIGGER update_backup_prod_v1_overs_updated_at BEFORE UPDATE ON backup_prod_v1.overs
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ============================================
-- COMMENTS AND DOCUMENTATION
-- ============================================

-- Copy table comments
COMMENT ON TABLE backup_prod_v1.users IS 'Backup: Stores user information from Google OAuth';
COMMENT ON TABLE backup_prod_v1.user_sessions IS 'Backup: Stores user session information for authentication';
COMMENT ON TABLE backup_prod_v1.oauth_states IS 'Backup: Stores OAuth state parameters for security';
COMMENT ON TABLE backup_prod_v1.schema_version IS 'Backup: Tracks database schema migration versions';
COMMENT ON TABLE backup_prod_v1.series IS 'Backup: Cricket tournaments and competitions';
COMMENT ON TABLE backup_prod_v1.matches IS 'Backup: Individual cricket matches with Team A vs Team B and toss functionality';
COMMENT ON TABLE backup_prod_v1.live_scoreboard IS 'Backup: Real-time match scoring and statistics';
COMMENT ON TABLE backup_prod_v1.innings IS 'Backup: Cricket innings tracking with runs, wickets, and overs';
COMMENT ON TABLE backup_prod_v1.overs IS 'Backup: Over-by-over tracking within innings';
COMMENT ON TABLE backup_prod_v1.balls IS 'Backup: Ball-by-ball events with run types and wickets';

-- Copy function comments
COMMENT ON FUNCTION backup_prod_v1.cleanup_expired_oauth_states() IS 'Backup: Cleans up expired OAuth states and returns count of deleted records';
COMMENT ON FUNCTION backup_prod_v1.cleanup_expired_user_sessions() IS 'Backup: Cleans up expired user sessions and returns count of deleted records';
COMMENT ON FUNCTION backup_prod_v1.get_user_statistics() IS 'Backup: Returns statistics about users, sessions, and OAuth states';

-- ============================================
-- PERMISSIONS SETUP
-- ============================================

-- Grant permissions to all common Supabase roles for backup schema
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
            RAISE NOTICE 'Granted backup permissions to role: %', role_name;
        EXCEPTION
            WHEN OTHERS THEN
                RAISE NOTICE 'Could not grant backup permissions to role %: %', role_name, SQLERRM;
        END;
    END LOOP;
END $$;

-- Set default privileges for future tables and sequences in backup schema
ALTER DEFAULT PRIVILEGES IN SCHEMA backup_prod_v1 GRANT ALL ON TABLES TO anon, authenticated, service_role, postgres;
ALTER DEFAULT PRIVILEGES IN SCHEMA backup_prod_v1 GRANT ALL ON SEQUENCES TO anon, authenticated, service_role, postgres;

-- ============================================
-- VERIFICATION AND STATISTICS
-- ============================================

SELECT 'Backup schema created successfully!' as status;
SELECT 'This backup includes all tables, data, indexes, triggers, functions, and permissions' as info;

-- Verify backup schema tables
SELECT 'backup_prod_v1 schema tables created:' as info;
SELECT table_name FROM information_schema.tables 
WHERE table_schema = 'backup_prod_v1' 
ORDER BY table_name;

-- Count records in each table for verification
SELECT 'Data copy verification - Record counts:' as info;
SELECT 'users' as table_name, COUNT(*) as record_count FROM backup_prod_v1.users
UNION ALL
SELECT 'user_sessions', COUNT(*) FROM backup_prod_v1.user_sessions
UNION ALL
SELECT 'oauth_states', COUNT(*) FROM backup_prod_v1.oauth_states
UNION ALL
SELECT 'schema_version', COUNT(*) FROM backup_prod_v1.schema_version
UNION ALL
SELECT 'series', COUNT(*) FROM backup_prod_v1.series
UNION ALL
SELECT 'matches', COUNT(*) FROM backup_prod_v1.matches
UNION ALL
SELECT 'live_scoreboard', COUNT(*) FROM backup_prod_v1.live_scoreboard
UNION ALL
SELECT 'innings', COUNT(*) FROM backup_prod_v1.innings
UNION ALL
SELECT 'overs', COUNT(*) FROM backup_prod_v1.overs
UNION ALL
SELECT 'balls', COUNT(*) FROM backup_prod_v1.balls
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

SELECT 'Backup completed successfully!' as final_status;