-- ============================================
-- SPARK PARK CRICKET - COMPLETE SCHEMA (PROD_V1)
-- ============================================
-- Comprehensive Cricket Tournament Management System
-- Version: 2.1.0
-- Date: 2025-01-27
-- Environment: Production (prod_v1)
-- Description: Complete schema for production environment with proper permissions
-- Features: Tables, indexes, triggers, functions, and automatic permission grants
-- ============================================

-- ============================================
-- EXTENSIONS
-- ============================================

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================
-- USER AUTHENTICATION TABLES (PROD_V1 Schema)
-- ============================================
-- These tables are created in the prod_v1 schema

-- Create users table (prod_v1 schema)
CREATE TABLE IF NOT EXISTS prod_v1.users (
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

-- Create user_sessions table (prod_v1 schema)
CREATE TABLE IF NOT EXISTS prod_v1.user_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES prod_v1.users(id) ON DELETE CASCADE,
    session_id VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create oauth_states table for storing OAuth state parameters (prod_v1 schema)
CREATE TABLE IF NOT EXISTS prod_v1.oauth_states (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    state VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    used_at TIMESTAMP WITH TIME ZONE NULL
);

-- Create schema_version table for tracking migrations (prod_v1 schema)
CREATE TABLE IF NOT EXISTS prod_v1.schema_version (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    version VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    applied_by VARCHAR(255),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Add missing columns if they don't exist (for existing tables)
DO $$
BEGIN
    -- Add applied_by column if it doesn't exist
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'schema_version' 
        AND column_name = 'applied_by'
        AND table_schema = 'prod_v1'
    ) THEN
        ALTER TABLE prod_v1.schema_version ADD COLUMN applied_by VARCHAR(255);
    END IF;
    
    -- Add updated_at column if it doesn't exist
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'schema_version' 
        AND column_name = 'updated_at'
        AND table_schema = 'prod_v1'
    ) THEN
        ALTER TABLE prod_v1.schema_version ADD COLUMN updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW();
    END IF;
END $$;

-- ============================================
-- CRICKET SCHEMA TABLES (PROD_V1)
-- ============================================

-- Create series table
CREATE TABLE IF NOT EXISTS prod_v1.series (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    start_date TIMESTAMP WITH TIME ZONE NOT NULL,
    end_date TIMESTAMP WITH TIME ZONE NOT NULL,
    created_by UUID REFERENCES prod_v1.users(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create matches table
CREATE TABLE IF NOT EXISTS prod_v1.matches (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    series_id UUID REFERENCES prod_v1.series(id) ON DELETE CASCADE,
    match_number INTEGER NOT NULL,
    date TIMESTAMP WITH TIME ZONE NOT NULL,
    status VARCHAR(20) DEFAULT 'live' CHECK (status IN ('live', 'completed', 'cancelled')),
    team_a_player_count INTEGER NOT NULL DEFAULT 11 CHECK (team_a_player_count >= 1 AND team_a_player_count <= 20),
    team_b_player_count INTEGER NOT NULL DEFAULT 11 CHECK (team_b_player_count >= 1 AND team_b_player_count <= 20),
    total_overs INTEGER NOT NULL DEFAULT 20 CHECK (total_overs >= 1 AND total_overs <= 20),
    toss_winner VARCHAR(1) NOT NULL CHECK (toss_winner IN ('A', 'B')),
    toss_type VARCHAR(1) NOT NULL CHECK (toss_type IN ('H', 'T')),
    batting_team VARCHAR(1) NOT NULL DEFAULT 'A' CHECK (batting_team IN ('A', 'B')),
    created_by UUID REFERENCES prod_v1.users(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create live_scoreboard table
CREATE TABLE IF NOT EXISTS prod_v1.live_scoreboard (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    match_id UUID REFERENCES prod_v1.matches(id) ON DELETE CASCADE,
    batting_team VARCHAR(1) NOT NULL CHECK (batting_team IN ('A', 'B')),
    score INTEGER DEFAULT 0,
    wickets INTEGER DEFAULT 0,
    overs DECIMAL(4,1) DEFAULT 0.0,
    balls INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create innings table
CREATE TABLE IF NOT EXISTS prod_v1.innings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    match_id UUID REFERENCES prod_v1.matches(id) ON DELETE CASCADE,
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

-- Create overs table
CREATE TABLE IF NOT EXISTS prod_v1.overs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    innings_id UUID REFERENCES prod_v1.innings(id) ON DELETE CASCADE,
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

-- Create balls table
CREATE TABLE IF NOT EXISTS prod_v1.balls (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    over_id UUID REFERENCES prod_v1.overs(id) ON DELETE CASCADE,
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
    CONSTRAINT prod_v1_balls_wicket_type_check CHECK (
        (is_wicket = true AND wicket_type IS NOT NULL) OR 
        (is_wicket = false AND wicket_type IS NULL)
    )
);

-- ============================================
-- INDEXES FOR PERFORMANCE
-- ============================================

-- User authentication indexes (prod_v1 schema)
CREATE INDEX IF NOT EXISTS idx_prod_v1_users_google_id ON prod_v1.users(google_id);
CREATE INDEX IF NOT EXISTS idx_prod_v1_users_email ON prod_v1.users(email);
CREATE INDEX IF NOT EXISTS idx_prod_v1_user_sessions_session_id ON prod_v1.user_sessions(session_id);
CREATE INDEX IF NOT EXISTS idx_prod_v1_user_sessions_user_id ON prod_v1.user_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_prod_v1_user_sessions_expires_at ON prod_v1.user_sessions(expires_at);
CREATE INDEX IF NOT EXISTS idx_prod_v1_oauth_states_state ON prod_v1.oauth_states(state);
CREATE INDEX IF NOT EXISTS idx_prod_v1_oauth_states_expires_at ON prod_v1.oauth_states(expires_at);
CREATE INDEX IF NOT EXISTS idx_prod_v1_oauth_states_used_at ON prod_v1.oauth_states(used_at);
CREATE INDEX IF NOT EXISTS idx_prod_v1_schema_version_version ON prod_v1.schema_version(version);

-- Cricket schema indexes (prod_v1)
-- Series indexes
CREATE INDEX IF NOT EXISTS idx_prod_v1_series_start_date ON prod_v1.series(start_date);
CREATE INDEX IF NOT EXISTS idx_prod_v1_series_end_date ON prod_v1.series(end_date);
CREATE INDEX IF NOT EXISTS idx_prod_v1_series_created_by ON prod_v1.series(created_by);

-- Matches indexes
CREATE INDEX IF NOT EXISTS idx_prod_v1_matches_series_id ON prod_v1.matches(series_id);
CREATE INDEX IF NOT EXISTS idx_prod_v1_matches_status ON prod_v1.matches(status);
CREATE INDEX IF NOT EXISTS idx_prod_v1_matches_date ON prod_v1.matches(date);
CREATE INDEX IF NOT EXISTS idx_prod_v1_matches_toss_winner ON prod_v1.matches(toss_winner);
CREATE INDEX IF NOT EXISTS idx_prod_v1_matches_batting_team ON prod_v1.matches(batting_team);
CREATE INDEX IF NOT EXISTS idx_prod_v1_matches_created_by ON prod_v1.matches(created_by);

-- Live scoreboard indexes
CREATE INDEX IF NOT EXISTS idx_prod_v1_live_scoreboard_match_id ON prod_v1.live_scoreboard(match_id);
CREATE INDEX IF NOT EXISTS idx_prod_v1_live_scoreboard_batting_team ON prod_v1.live_scoreboard(batting_team);

-- Innings indexes
CREATE INDEX IF NOT EXISTS idx_prod_v1_innings_match_id ON prod_v1.innings(match_id);
CREATE INDEX IF NOT EXISTS idx_prod_v1_innings_batting_team ON prod_v1.innings(batting_team);
CREATE INDEX IF NOT EXISTS idx_prod_v1_innings_status ON prod_v1.innings(status);

-- Overs indexes
CREATE INDEX IF NOT EXISTS idx_prod_v1_overs_innings_id ON prod_v1.overs(innings_id);
CREATE INDEX IF NOT EXISTS idx_prod_v1_overs_status ON prod_v1.overs(status);

-- Balls indexes
CREATE INDEX IF NOT EXISTS idx_prod_v1_balls_over_id ON prod_v1.balls(over_id);
CREATE INDEX IF NOT EXISTS idx_prod_v1_balls_run_type ON prod_v1.balls(run_type);
CREATE INDEX IF NOT EXISTS idx_prod_v1_balls_is_wicket ON prod_v1.balls(is_wicket);

-- ============================================
-- FUNCTIONS AND TRIGGERS
-- ============================================

-- Create function to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create function to clean up expired OAuth states
CREATE OR REPLACE FUNCTION cleanup_expired_oauth_states()
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    DELETE FROM prod_v1.oauth_states 
    WHERE expires_at < NOW() 
    AND (used_at IS NOT NULL OR expires_at < NOW() - INTERVAL '1 hour');
    
    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    RETURN deleted_count;
END;
$$ language 'plpgsql';

-- Create function to clean up expired user sessions
CREATE OR REPLACE FUNCTION cleanup_expired_user_sessions()
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    DELETE FROM prod_v1.user_sessions 
    WHERE expires_at < NOW();
    
    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    RETURN deleted_count;
END;
$$ language 'plpgsql';

-- Create function to get user statistics
CREATE OR REPLACE FUNCTION get_user_statistics()
RETURNS TABLE (
    total_users BIGINT,
    active_sessions BIGINT,
    expired_sessions BIGINT,
    oauth_states_pending BIGINT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        (SELECT COUNT(*) FROM prod_v1.users) as total_users,
        (SELECT COUNT(*) FROM prod_v1.user_sessions WHERE expires_at > NOW()) as active_sessions,
        (SELECT COUNT(*) FROM prod_v1.user_sessions WHERE expires_at <= NOW()) as expired_sessions,
        (SELECT COUNT(*) FROM prod_v1.oauth_states WHERE expires_at > NOW() AND used_at IS NULL) as oauth_states_pending;
END;
$$ language 'plpgsql';

-- Create triggers to automatically update updated_at for prod_v1 schema
DROP TRIGGER IF EXISTS update_prod_v1_users_updated_at ON prod_v1.users;
CREATE TRIGGER update_prod_v1_users_updated_at BEFORE UPDATE ON prod_v1.users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_prod_v1_user_sessions_updated_at ON prod_v1.user_sessions;
CREATE TRIGGER update_prod_v1_user_sessions_updated_at BEFORE UPDATE ON prod_v1.user_sessions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_prod_v1_schema_version_updated_at ON prod_v1.schema_version;
CREATE TRIGGER update_prod_v1_schema_version_updated_at BEFORE UPDATE ON prod_v1.schema_version
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Create triggers for prod_v1 schema
DROP TRIGGER IF EXISTS update_prod_v1_series_updated_at ON prod_v1.series;
CREATE TRIGGER update_prod_v1_series_updated_at BEFORE UPDATE ON prod_v1.series
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_prod_v1_matches_updated_at ON prod_v1.matches;
CREATE TRIGGER update_prod_v1_matches_updated_at BEFORE UPDATE ON prod_v1.matches
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_prod_v1_live_scoreboard_updated_at ON prod_v1.live_scoreboard;
CREATE TRIGGER update_prod_v1_live_scoreboard_updated_at BEFORE UPDATE ON prod_v1.live_scoreboard
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_prod_v1_innings_updated_at ON prod_v1.innings;
CREATE TRIGGER update_prod_v1_innings_updated_at BEFORE UPDATE ON prod_v1.innings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_prod_v1_overs_updated_at ON prod_v1.overs;
CREATE TRIGGER update_prod_v1_overs_updated_at BEFORE UPDATE ON prod_v1.overs
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ============================================
-- DOCUMENTATION AND COMMENTS
-- ============================================

-- User authentication table comments (prod_v1 schema)
COMMENT ON TABLE prod_v1.users IS 'Stores user information from Google OAuth';
COMMENT ON TABLE prod_v1.user_sessions IS 'Stores user session information for authentication';
COMMENT ON TABLE prod_v1.oauth_states IS 'Stores OAuth state parameters for security';
COMMENT ON TABLE prod_v1.schema_version IS 'Tracks database schema migration versions';

-- Cricket schema table comments (prod_v1)
COMMENT ON TABLE prod_v1.series IS 'Cricket tournaments and competitions';
COMMENT ON TABLE prod_v1.matches IS 'Individual cricket matches with Team A vs Team B and toss functionality';
COMMENT ON TABLE prod_v1.live_scoreboard IS 'Real-time match scoring and statistics';
COMMENT ON TABLE prod_v1.innings IS 'Cricket innings tracking with runs, wickets, and overs';
COMMENT ON TABLE prod_v1.overs IS 'Over-by-over tracking within innings';
COMMENT ON TABLE prod_v1.balls IS 'Ball-by-ball events with run types and wickets';

-- User authentication column comments (prod_v1 schema)
COMMENT ON COLUMN prod_v1.users.google_id IS 'Google OAuth user ID';
COMMENT ON COLUMN prod_v1.users.email IS 'User email address';
COMMENT ON COLUMN prod_v1.users.name IS 'User display name';
COMMENT ON COLUMN prod_v1.users.picture IS 'User profile picture URL';
COMMENT ON COLUMN prod_v1.users.email_verified IS 'Whether the email is verified by Google';
COMMENT ON COLUMN prod_v1.users.last_login_at IS 'Timestamp of last successful login';

COMMENT ON COLUMN prod_v1.user_sessions.user_id IS 'Reference to users table';
COMMENT ON COLUMN prod_v1.user_sessions.session_id IS 'Unique session identifier';
COMMENT ON COLUMN prod_v1.user_sessions.expires_at IS 'Session expiration timestamp';

COMMENT ON COLUMN prod_v1.oauth_states.state IS 'OAuth state parameter for security';
COMMENT ON COLUMN prod_v1.oauth_states.expires_at IS 'State expiration timestamp';
COMMENT ON COLUMN prod_v1.oauth_states.used_at IS 'Timestamp when state was used';

COMMENT ON COLUMN prod_v1.schema_version.version IS 'Migration version identifier';
COMMENT ON COLUMN prod_v1.schema_version.description IS 'Description of the migration';
COMMENT ON COLUMN prod_v1.schema_version.applied_at IS 'Timestamp when migration was applied';
COMMENT ON COLUMN prod_v1.schema_version.applied_by IS 'User or system that applied the migration';
COMMENT ON COLUMN prod_v1.schema_version.updated_at IS 'Timestamp when record was last updated';

-- Function comments
COMMENT ON FUNCTION cleanup_expired_oauth_states() IS 'Cleans up expired OAuth states and returns count of deleted records';
COMMENT ON FUNCTION cleanup_expired_user_sessions() IS 'Cleans up expired user sessions and returns count of deleted records';
COMMENT ON FUNCTION get_user_statistics() IS 'Returns statistics about users, sessions, and OAuth states';

-- Cricket schema column comments (prod_v1)
COMMENT ON COLUMN prod_v1.series.created_by IS 'User who created this series';
COMMENT ON COLUMN prod_v1.matches.created_by IS 'User who created this match';
COMMENT ON COLUMN prod_v1.matches.toss_winner IS 'Team that won the toss: A or B';
COMMENT ON COLUMN prod_v1.matches.toss_type IS 'Toss result: H (Heads) or T (Tails)';
COMMENT ON COLUMN prod_v1.matches.batting_team IS 'Team currently batting: A or B';
COMMENT ON COLUMN prod_v1.matches.team_a_player_count IS 'Number of players in Team A (1-20)';
COMMENT ON COLUMN prod_v1.matches.team_b_player_count IS 'Number of players in Team B (1-20)';
COMMENT ON COLUMN prod_v1.matches.total_overs IS 'Total overs for the match (1-20)';

COMMENT ON COLUMN prod_v1.innings.innings_number IS 'Innings number: 1 (first innings) or 2 (second innings)';
COMMENT ON COLUMN prod_v1.innings.batting_team IS 'Team currently batting: A or B';
COMMENT ON COLUMN prod_v1.innings.total_runs IS 'Total runs scored in this innings';
COMMENT ON COLUMN prod_v1.innings.total_wickets IS 'Total wickets fallen in this innings (0-10)';
COMMENT ON COLUMN prod_v1.innings.total_overs IS 'Total overs completed in this innings (decimal)';
COMMENT ON COLUMN prod_v1.innings.total_balls IS 'Total balls bowled in this innings';
COMMENT ON COLUMN prod_v1.innings.status IS 'Innings status: in_progress or completed';

COMMENT ON COLUMN prod_v1.overs.over_number IS 'Over number within the innings (1, 2, 3, etc.)';
COMMENT ON COLUMN prod_v1.overs.total_runs IS 'Total runs scored in this over';
COMMENT ON COLUMN prod_v1.overs.total_balls IS 'Total balls bowled in this over (0-6)';
COMMENT ON COLUMN prod_v1.overs.total_wickets IS 'Total wickets fallen in this over';
COMMENT ON COLUMN prod_v1.overs.status IS 'Over status: in_progress or completed';

COMMENT ON COLUMN prod_v1.balls.ball_number IS 'Ball number within the over (1-20 to allow for illegal balls)';
COMMENT ON COLUMN prod_v1.balls.ball_type IS 'Type of ball: good, wide, no_ball, dead_ball';
COMMENT ON COLUMN prod_v1.balls.run_type IS 'Run type: 0-9 (runs), NB (No Ball), WD (Wide), LB (Leg Byes), WC (Wicket)';
COMMENT ON COLUMN prod_v1.balls.runs IS 'Actual runs scored from this ball';
COMMENT ON COLUMN prod_v1.balls.byes IS 'Additional runs from byes (0-6)';
COMMENT ON COLUMN prod_v1.balls.is_wicket IS 'Whether this ball resulted in a wicket';
COMMENT ON COLUMN prod_v1.balls.wicket_type IS 'Type of wicket: bowled, caught, lbw, run_out, stumped, hit_wicket';

-- ============================================
-- PERMISSIONS SETUP
-- ============================================

-- Grant permissions to all common Supabase roles
DO $$
DECLARE
    role_name TEXT;
    roles TEXT[] := ARRAY['anon', 'authenticated', 'service_role', 'postgres', 'supabase_auth_admin', 'supabase_storage_admin'];
BEGIN
    FOREACH role_name IN ARRAY roles
    LOOP
        BEGIN
            EXECUTE format('GRANT USAGE ON SCHEMA prod_v1 TO %I', role_name);
            EXECUTE format('GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA prod_v1 TO %I', role_name);
            EXECUTE format('GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA prod_v1 TO %I', role_name);
            RAISE NOTICE 'Granted permissions to role: %', role_name;
        EXCEPTION
            WHEN OTHERS THEN
                RAISE NOTICE 'Could not grant permissions to role %: %', role_name, SQLERRM;
        END;
    END LOOP;
END $$;

-- Set default privileges for future tables and sequences
ALTER DEFAULT PRIVILEGES IN SCHEMA prod_v1 GRANT ALL ON TABLES TO anon, authenticated, service_role, postgres;
ALTER DEFAULT PRIVILEGES IN SCHEMA prod_v1 GRANT ALL ON SEQUENCES TO anon, authenticated, service_role, postgres;

-- ============================================
-- VERIFICATION
-- ============================================

SELECT 'Complete schema created successfully for prod_v1!' as status;
SELECT 'This migration includes all tables, constraints, indexes, triggers, functions, and permissions' as info;

-- Verify prod_v1 schema tables (including auth tables)
SELECT 'prod_v1 schema tables created (including auth):' as info;
SELECT table_name FROM information_schema.tables 
WHERE table_schema = 'prod_v1' 
AND table_name IN ('users', 'user_sessions', 'oauth_states', 'schema_version', 'series', 'matches', 'live_scoreboard', 'innings', 'overs', 'balls')
ORDER BY table_name;

-- Verify prod_v1 schema tables
SELECT 'prod_v1 schema tables created:' as info;
SELECT table_name FROM information_schema.tables 
WHERE table_schema = 'prod_v1' 
AND table_name IN ('series', 'matches', 'live_scoreboard', 'innings', 'overs', 'balls')
ORDER BY table_name;

-- Verify permissions on prod_v1 schema
SELECT 'prod_v1 schema permissions granted to:' as info;
SELECT 
    n.nspname as schema_name,
    r.rolname as role_name,
    has_schema_privilege(r.oid, n.oid, 'USAGE') as has_usage,
    has_schema_privilege(r.oid, n.oid, 'CREATE') as has_create
FROM pg_namespace n
CROSS JOIN pg_roles r
WHERE n.nspname = 'prod_v1'
AND r.rolname IN ('anon', 'authenticated', 'service_role', 'postgres')
ORDER BY r.rolname;
