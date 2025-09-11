-- Spark Park Cricket - Initial Schema Migration
-- Simplified Cricket Tournament Management System
-- Version: 2.0.0
-- Date: 2025-09-10

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Series table - Cricket tournaments/competitions
CREATE TABLE IF NOT EXISTS series (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    start_date TIMESTAMP WITH TIME ZONE NOT NULL,
    end_date TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Matches table - Simplified with Team A vs Team B
CREATE TABLE IF NOT EXISTS matches (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    series_id UUID REFERENCES series(id) ON DELETE CASCADE,
    match_number INTEGER NOT NULL,
    date TIMESTAMP WITH TIME ZONE NOT NULL,
    status VARCHAR(20) DEFAULT 'live' CHECK (status IN ('live', 'completed', 'cancelled')),
    team_a_player_count INTEGER NOT NULL DEFAULT 11 CHECK (team_a_player_count >= 1 AND team_a_player_count <= 11),
    team_b_player_count INTEGER NOT NULL DEFAULT 11 CHECK (team_b_player_count >= 1 AND team_b_player_count <= 11),
    total_overs INTEGER NOT NULL DEFAULT 20 CHECK (total_overs >= 1 AND total_overs <= 20),
    toss_winner VARCHAR(1) NOT NULL CHECK (toss_winner IN ('A', 'B')),
    toss_type VARCHAR(1) NOT NULL CHECK (toss_type IN ('H', 'T')),
    batting_team VARCHAR(1) NOT NULL DEFAULT 'A' CHECK (batting_team IN ('A', 'B')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Live scoreboard table - Real-time match scoring
CREATE TABLE IF NOT EXISTS live_scoreboard (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    match_id UUID REFERENCES matches(id) ON DELETE CASCADE,
    batting_team VARCHAR(1) NOT NULL CHECK (batting_team IN ('A', 'B')),
    score INTEGER DEFAULT 0,
    wickets INTEGER DEFAULT 0,
    overs DECIMAL(4,1) DEFAULT 0.0,
    balls INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Overs table - Over-by-over tracking
CREATE TABLE IF NOT EXISTS overs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    match_id UUID REFERENCES matches(id) ON DELETE CASCADE,
    over_number INTEGER NOT NULL,
    batting_team VARCHAR(1) NOT NULL CHECK (batting_team IN ('A', 'B')),
    total_runs INTEGER DEFAULT 0,
    total_balls INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Balls table - Ball-by-ball tracking with run types
CREATE TABLE IF NOT EXISTS balls (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    over_id UUID REFERENCES overs(id) ON DELETE CASCADE,
    ball_number INTEGER NOT NULL,
    ball_type VARCHAR(20) NOT NULL CHECK (ball_type IN ('good', 'wide', 'no_ball', 'dead_ball')),
    run_type VARCHAR(2) NOT NULL CHECK (run_type IN ('1', '2', '3', '4', '5', '6', '7', '8', '9', 'NB', 'WD', 'LB')),
    is_wicket BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Performance indexes
CREATE INDEX IF NOT EXISTS idx_matches_series_id ON matches(series_id);
CREATE INDEX IF NOT EXISTS idx_matches_status ON matches(status);
CREATE INDEX IF NOT EXISTS idx_matches_date ON matches(date);
CREATE INDEX IF NOT EXISTS idx_live_scoreboard_match_id ON live_scoreboard(match_id);
CREATE INDEX IF NOT EXISTS idx_overs_match_id ON overs(match_id);
CREATE INDEX IF NOT EXISTS idx_balls_over_id ON balls(over_id);

-- Add comments for documentation
COMMENT ON TABLE series IS 'Cricket tournaments and competitions';
COMMENT ON TABLE matches IS 'Individual cricket matches with Team A vs Team B and toss functionality';
COMMENT ON TABLE live_scoreboard IS 'Real-time match scoring and statistics';
COMMENT ON TABLE overs IS 'Over-by-over match tracking';
COMMENT ON TABLE balls IS 'Ball-by-ball events with run types (1-9, NB, WD, LB)';

COMMENT ON COLUMN matches.toss_winner IS 'Team that won the toss: A or B';
COMMENT ON COLUMN matches.toss_type IS 'Toss result: H (Heads) or T (Tails)';
COMMENT ON COLUMN matches.batting_team IS 'Team currently batting: A or B';
COMMENT ON COLUMN matches.team_a_player_count IS 'Number of players in Team A (1-11)';
COMMENT ON COLUMN matches.team_b_player_count IS 'Number of players in Team B (1-11)';
COMMENT ON COLUMN matches.total_overs IS 'Total overs for the match (1-20)';

COMMENT ON COLUMN balls.run_type IS 'Run type: 1-9 (runs), NB (No Ball), WD (Wide), LB (Leg Byes)';
COMMENT ON COLUMN balls.ball_type IS 'Type of ball: good, wide, no_ball, dead_ball';
