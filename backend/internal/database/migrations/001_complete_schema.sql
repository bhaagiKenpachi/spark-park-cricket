-- Spark Park Cricket - Complete Schema Migration
-- Comprehensive Cricket Tournament Management System with Scorecard
-- Version: 3.0.0
-- Date: 2025-09-15

-- ============================================
-- EXTENSIONS
-- ============================================

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================
-- SERIES TABLE
-- ============================================

CREATE TABLE IF NOT EXISTS dev_v1.series (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    start_date TIMESTAMP WITH TIME ZONE NOT NULL,
    end_date TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- ============================================
-- MATCHES TABLE
-- ============================================

CREATE TABLE IF NOT EXISTS dev_v1.matches (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    series_id UUID REFERENCES dev_v1.series(id) ON DELETE CASCADE,
    match_number INTEGER NOT NULL,
    date TIMESTAMP WITH TIME ZONE NOT NULL,
    status VARCHAR(20) DEFAULT 'live' CHECK (status IN ('live', 'completed', 'cancelled')),
    team_a_player_count INTEGER NOT NULL DEFAULT 11 CHECK (team_a_player_count >= 1 AND team_a_player_count <= 20),
    team_b_player_count INTEGER NOT NULL DEFAULT 11 CHECK (team_b_player_count >= 1 AND team_b_player_count <= 20),
    total_overs INTEGER NOT NULL DEFAULT 20 CHECK (total_overs >= 1 AND total_overs <= 20),
    toss_winner VARCHAR(1) NOT NULL CHECK (toss_winner IN ('A', 'B')),
    toss_type VARCHAR(1) NOT NULL CHECK (toss_type IN ('H', 'T')),
    batting_team VARCHAR(1) NOT NULL DEFAULT 'A' CHECK (batting_team IN ('A', 'B')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- ============================================
-- LIVE SCOREBOARD TABLE
-- ============================================

CREATE TABLE IF NOT EXISTS dev_v1.live_scoreboard (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    match_id UUID REFERENCES dev_v1.matches(id) ON DELETE CASCADE,
    batting_team VARCHAR(1) NOT NULL CHECK (batting_team IN ('A', 'B')),
    score INTEGER DEFAULT 0,
    wickets INTEGER DEFAULT 0,
    overs DECIMAL(4,1) DEFAULT 0.0,
    balls INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- ============================================
-- INNINGS TABLE
-- ============================================

CREATE TABLE IF NOT EXISTS dev_v1.innings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    match_id UUID REFERENCES dev_v1.matches(id) ON DELETE CASCADE,
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

-- ============================================
-- OVERS TABLE
-- ============================================

CREATE TABLE IF NOT EXISTS dev_v1.overs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    innings_id UUID REFERENCES dev_v1.innings(id) ON DELETE CASCADE,
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

-- ============================================
-- BALLS TABLE
-- ============================================

CREATE TABLE IF NOT EXISTS dev_v1.balls (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    over_id UUID REFERENCES dev_v1.overs(id) ON DELETE CASCADE,
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
    CONSTRAINT dev_v1_balls_wicket_type_check CHECK (
        (is_wicket = true AND wicket_type IS NOT NULL) OR 
        (is_wicket = false AND wicket_type IS NULL)
    )
);

-- ============================================
-- PERFORMANCE INDEXES
-- ============================================

-- Series indexes
CREATE INDEX IF NOT EXISTS idx_series_start_date ON dev_v1.series(start_date);
CREATE INDEX IF NOT EXISTS idx_series_end_date ON dev_v1.series(end_date);

-- Matches indexes
CREATE INDEX IF NOT EXISTS idx_matches_series_id ON dev_v1.matches(series_id);
CREATE INDEX IF NOT EXISTS idx_matches_status ON dev_v1.matches(status);
CREATE INDEX IF NOT EXISTS idx_matches_date ON dev_v1.matches(date);
CREATE INDEX IF NOT EXISTS idx_matches_toss_winner ON dev_v1.matches(toss_winner);
CREATE INDEX IF NOT EXISTS idx_matches_batting_team ON dev_v1.matches(batting_team);

-- Live scoreboard indexes
CREATE INDEX IF NOT EXISTS idx_live_scoreboard_match_id ON dev_v1.live_scoreboard(match_id);
CREATE INDEX IF NOT EXISTS idx_live_scoreboard_batting_team ON dev_v1.live_scoreboard(batting_team);

-- Innings indexes
CREATE INDEX IF NOT EXISTS idx_innings_match_id ON dev_v1.innings(match_id);
CREATE INDEX IF NOT EXISTS idx_innings_batting_team ON dev_v1.innings(batting_team);
CREATE INDEX IF NOT EXISTS idx_innings_status ON dev_v1.innings(status);

-- Overs indexes
CREATE INDEX IF NOT EXISTS idx_overs_innings_id ON dev_v1.overs(innings_id);
CREATE INDEX IF NOT EXISTS idx_overs_status ON dev_v1.overs(status);

-- Balls indexes
CREATE INDEX IF NOT EXISTS idx_balls_over_id ON dev_v1.balls(over_id);
CREATE INDEX IF NOT EXISTS idx_balls_run_type ON dev_v1.balls(run_type);
CREATE INDEX IF NOT EXISTS idx_balls_is_wicket ON dev_v1.balls(is_wicket);

-- ============================================
-- DOCUMENTATION AND COMMENTS
-- ============================================

-- Table comments
COMMENT ON TABLE dev_v1.series IS 'Cricket tournaments and competitions';
COMMENT ON TABLE dev_v1.matches IS 'Individual cricket matches with Team A vs Team B and toss functionality';
COMMENT ON TABLE dev_v1.live_scoreboard IS 'Real-time match scoring and statistics';
COMMENT ON TABLE dev_v1.innings IS 'Cricket innings tracking with runs, wickets, and overs';
COMMENT ON TABLE dev_v1.overs IS 'Over-by-over tracking within innings';
COMMENT ON TABLE dev_v1.balls IS 'Ball-by-ball events with run types and wickets';

-- Column comments
COMMENT ON COLUMN dev_v1.matches.toss_winner IS 'Team that won the toss: A or B';
COMMENT ON COLUMN dev_v1.matches.toss_type IS 'Toss result: H (Heads) or T (Tails)';
COMMENT ON COLUMN dev_v1.matches.batting_team IS 'Team currently batting: A or B';
COMMENT ON COLUMN dev_v1.matches.team_a_player_count IS 'Number of players in Team A (1-20)';
COMMENT ON COLUMN dev_v1.matches.team_b_player_count IS 'Number of players in Team B (1-20)';
COMMENT ON COLUMN dev_v1.matches.total_overs IS 'Total overs for the match (1-20)';

COMMENT ON COLUMN dev_v1.innings.innings_number IS 'Innings number: 1 (first innings) or 2 (second innings)';
COMMENT ON COLUMN dev_v1.innings.batting_team IS 'Team currently batting: A or B';
COMMENT ON COLUMN dev_v1.innings.total_runs IS 'Total runs scored in this innings';
COMMENT ON COLUMN dev_v1.innings.total_wickets IS 'Total wickets fallen in this innings (0-10)';
COMMENT ON COLUMN dev_v1.innings.total_overs IS 'Total overs completed in this innings (decimal)';
COMMENT ON COLUMN dev_v1.innings.total_balls IS 'Total balls bowled in this innings';
COMMENT ON COLUMN dev_v1.innings.status IS 'Innings status: in_progress or completed';

COMMENT ON COLUMN dev_v1.overs.over_number IS 'Over number within the innings (1, 2, 3, etc.)';
COMMENT ON COLUMN dev_v1.overs.total_runs IS 'Total runs scored in this over';
COMMENT ON COLUMN dev_v1.overs.total_balls IS 'Total balls bowled in this over (0-6)';
COMMENT ON COLUMN dev_v1.overs.total_wickets IS 'Total wickets fallen in this over';
COMMENT ON COLUMN dev_v1.overs.status IS 'Over status: in_progress or completed';

COMMENT ON COLUMN dev_v1.balls.ball_number IS 'Ball number within the over (1-20 to allow for illegal balls)';
COMMENT ON COLUMN dev_v1.balls.ball_type IS 'Type of ball: good, wide, no_ball, dead_ball';
COMMENT ON COLUMN dev_v1.balls.run_type IS 'Run type: 0-9 (runs), NB (No Ball), WD (Wide), LB (Leg Byes), WC (Wicket)';
COMMENT ON COLUMN dev_v1.balls.runs IS 'Actual runs scored from this ball';
COMMENT ON COLUMN dev_v1.balls.byes IS 'Additional runs from byes (0-6)';
COMMENT ON COLUMN dev_v1.balls.is_wicket IS 'Whether this ball resulted in a wicket';
COMMENT ON COLUMN dev_v1.balls.wicket_type IS 'Type of wicket: bowled, caught, lbw, run_out, stumped, hit_wicket';

-- ============================================
-- VERIFICATION
-- ============================================

SELECT 'Complete schema created successfully!' as status;
SELECT 'Tables created:' as info;
SELECT table_name FROM information_schema.tables 
WHERE table_schema = 'dev_v1' 
AND table_name IN ('series', 'matches', 'live_scoreboard', 'innings', 'overs', 'balls')
ORDER BY table_name;
