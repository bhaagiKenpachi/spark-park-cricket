-- Scorecard Schema Migration
-- Creates tables for comprehensive cricket scoring
-- Version: 2.0.2
-- Date: 2025-09-10

-- ============================================
-- INNINGS TABLE
-- ============================================

CREATE TABLE IF NOT EXISTS innings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    match_id UUID REFERENCES matches(id) ON DELETE CASCADE,
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
-- OVERS TABLE (Updated for scorecard)
-- ============================================

-- Drop existing overs table if it exists and recreate with new structure
DROP TABLE IF EXISTS overs CASCADE;

CREATE TABLE overs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    innings_id UUID REFERENCES innings(id) ON DELETE CASCADE,
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
-- BALLS TABLE (Updated for scorecard)
-- ============================================

-- Drop existing balls table if it exists and recreate with new structure
DROP TABLE IF EXISTS balls CASCADE;

CREATE TABLE balls (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    over_id UUID REFERENCES overs(id) ON DELETE CASCADE,
    ball_number INTEGER NOT NULL CHECK (ball_number >= 1 AND ball_number <= 6),
    ball_type VARCHAR(20) NOT NULL CHECK (ball_type IN ('good', 'wide', 'no_ball', 'dead_ball')),
    run_type VARCHAR(2) NOT NULL CHECK (run_type IN ('0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'NB', 'WD', 'LB', 'WC')),
    runs INTEGER DEFAULT 0 CHECK (runs >= 0),
    byes INTEGER DEFAULT 0 CHECK (byes >= 0 AND byes <= 6), -- Additional runs from byes
    is_wicket BOOLEAN DEFAULT FALSE,
    wicket_type VARCHAR(20) CHECK (wicket_type IN ('bowled', 'caught', 'lbw', 'run_out', 'stumped', 'hit_wicket')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Ensure only one ball per over per ball number
    UNIQUE(over_id, ball_number),
    
    -- Ensure wicket_type is NULL when is_wicket is false
    CONSTRAINT balls_wicket_type_check CHECK (
        (is_wicket = true AND wicket_type IS NOT NULL) OR 
        (is_wicket = false AND wicket_type IS NULL)
    )
);

-- ============================================
-- INDEXES FOR PERFORMANCE
-- ============================================

-- Innings indexes
CREATE INDEX IF NOT EXISTS idx_innings_match_id ON innings(match_id);
CREATE INDEX IF NOT EXISTS idx_innings_batting_team ON innings(batting_team);
CREATE INDEX IF NOT EXISTS idx_innings_status ON innings(status);

-- Overs indexes
CREATE INDEX IF NOT EXISTS idx_overs_innings_id ON overs(innings_id);
CREATE INDEX IF NOT EXISTS idx_overs_status ON overs(status);

-- Balls indexes
CREATE INDEX IF NOT EXISTS idx_balls_over_id ON balls(over_id);
CREATE INDEX IF NOT EXISTS idx_balls_run_type ON balls(run_type);
CREATE INDEX IF NOT EXISTS idx_balls_is_wicket ON balls(is_wicket);

-- ============================================
-- COMMENTS FOR DOCUMENTATION
-- ============================================

COMMENT ON TABLE innings IS 'Cricket innings tracking with runs, wickets, and overs';
COMMENT ON TABLE overs IS 'Over-by-over tracking within innings';
COMMENT ON TABLE balls IS 'Ball-by-ball events with run types and wickets';

COMMENT ON COLUMN innings.innings_number IS 'Innings number: 1 (first innings) or 2 (second innings)';
COMMENT ON COLUMN innings.batting_team IS 'Team currently batting: A or B';
COMMENT ON COLUMN innings.total_runs IS 'Total runs scored in this innings';
COMMENT ON COLUMN innings.total_wickets IS 'Total wickets fallen in this innings (0-10)';
COMMENT ON COLUMN innings.total_overs IS 'Total overs completed in this innings (decimal)';
COMMENT ON COLUMN innings.total_balls IS 'Total balls bowled in this innings';
COMMENT ON COLUMN innings.status IS 'Innings status: in_progress or completed';

COMMENT ON COLUMN overs.over_number IS 'Over number within the innings (1, 2, 3, etc.)';
COMMENT ON COLUMN overs.total_runs IS 'Total runs scored in this over';
COMMENT ON COLUMN overs.total_balls IS 'Total balls bowled in this over (0-6)';
COMMENT ON COLUMN overs.total_wickets IS 'Total wickets fallen in this over';
COMMENT ON COLUMN overs.status IS 'Over status: in_progress or completed';

COMMENT ON COLUMN balls.ball_number IS 'Ball number within the over (1-6)';
COMMENT ON COLUMN balls.ball_type IS 'Type of ball: good, wide, no_ball, dead_ball';
COMMENT ON COLUMN balls.run_type IS 'Run type: 0-9 (runs), NB (No Ball), WD (Wide), LB (Leg Byes)';
COMMENT ON COLUMN balls.runs IS 'Actual runs scored from this ball';
COMMENT ON COLUMN balls.is_wicket IS 'Whether this ball resulted in a wicket';
COMMENT ON COLUMN balls.wicket_type IS 'Type of wicket: bowled, caught, lbw, run_out, stumped, hit_wicket';

-- ============================================
-- VERIFICATION
-- ============================================

SELECT 'Scorecard schema created successfully!' as status;
SELECT 'Tables created:' as info;
SELECT table_name FROM information_schema.tables 
WHERE table_schema = 'public' 
AND table_name IN ('innings', 'overs', 'balls')
ORDER BY table_name;
