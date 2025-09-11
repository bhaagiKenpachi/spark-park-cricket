-- Test Database Setup Script
-- This script sets up the testing_db schema for running integration and e2e tests

-- Create testing_db schema if it doesn't exist
CREATE SCHEMA IF NOT EXISTS testing_db;

-- Set search path to testing_db
SET search_path TO testing_db;

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create series table
CREATE TABLE IF NOT EXISTS testing_db.series (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    start_date TIMESTAMP WITH TIME ZONE NOT NULL,
    end_date TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create matches table
CREATE TABLE IF NOT EXISTS testing_db.matches (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    series_id UUID NOT NULL REFERENCES testing_db.series(id) ON DELETE CASCADE,
    match_number INTEGER NOT NULL,
    date TIMESTAMP WITH TIME ZONE NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'scheduled',
    team_a_player_count INTEGER NOT NULL DEFAULT 11,
    team_b_player_count INTEGER NOT NULL DEFAULT 11,
    total_overs INTEGER NOT NULL DEFAULT 20,
    toss_winner VARCHAR(10) NOT NULL,
    toss_type VARCHAR(10) NOT NULL,
    batting_team VARCHAR(10) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(series_id, match_number)
);

-- Create innings table
CREATE TABLE IF NOT EXISTS testing_db.innings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    match_id UUID NOT NULL REFERENCES testing_db.matches(id) ON DELETE CASCADE,
    innings_number INTEGER NOT NULL,
    batting_team VARCHAR(10) NOT NULL,
    total_runs INTEGER NOT NULL DEFAULT 0,
    total_wickets INTEGER NOT NULL DEFAULT 0,
    total_overs DECIMAL(4,1) NOT NULL DEFAULT 0.0,
    total_balls INTEGER NOT NULL DEFAULT 0,
    status VARCHAR(50) NOT NULL DEFAULT 'in_progress',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(match_id, innings_number)
);

-- Create scorecard_overs table
CREATE TABLE IF NOT EXISTS testing_db.scorecard_overs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    innings_id UUID NOT NULL REFERENCES testing_db.innings(id) ON DELETE CASCADE,
    over_number INTEGER NOT NULL,
    total_runs INTEGER NOT NULL DEFAULT 0,
    total_balls INTEGER NOT NULL DEFAULT 0,
    status VARCHAR(50) NOT NULL DEFAULT 'in_progress',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(innings_id, over_number)
);

-- Create scorecard_balls table
CREATE TABLE IF NOT EXISTS testing_db.scorecard_balls (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    over_id UUID NOT NULL REFERENCES testing_db.scorecard_overs(id) ON DELETE CASCADE,
    ball_number INTEGER NOT NULL,
    ball_type VARCHAR(20) NOT NULL,
    run_type VARCHAR(10) NOT NULL,
    runs INTEGER NOT NULL DEFAULT 0,
    byes INTEGER NOT NULL DEFAULT 0,
    is_wicket BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(over_id, ball_number),
    CONSTRAINT balls_ball_type_check CHECK (ball_type IN ('good', 'wide', 'no_ball', 'dead_ball')),
    CONSTRAINT balls_run_type_check CHECK (run_type IN ('0', '1', '2', '3', '4', '5', '6', 'NB', 'WD', 'WC', 'LB')),
    CONSTRAINT balls_wicket_type_check CHECK (
        (is_wicket = true AND run_type = 'WC') OR 
        (is_wicket = false AND run_type != 'WC')
    ),
    CONSTRAINT balls_ball_number_check CHECK (ball_number >= 1 AND ball_number <= 20)
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_series_name ON testing_db.series(name);
CREATE INDEX IF NOT EXISTS idx_matches_series_id ON testing_db.matches(series_id);
CREATE INDEX IF NOT EXISTS idx_matches_status ON testing_db.matches(status);
CREATE INDEX IF NOT EXISTS idx_innings_match_id ON testing_db.innings(match_id);
CREATE INDEX IF NOT EXISTS idx_innings_number ON testing_db.innings(innings_number);
CREATE INDEX IF NOT EXISTS idx_overs_innings_id ON testing_db.scorecard_overs(innings_id);
CREATE INDEX IF NOT EXISTS idx_overs_number ON testing_db.scorecard_overs(over_number);
CREATE INDEX IF NOT EXISTS idx_balls_over_id ON testing_db.scorecard_balls(over_id);
CREATE INDEX IF NOT EXISTS idx_balls_number ON testing_db.scorecard_balls(ball_number);

-- Create functions for updating timestamps
CREATE OR REPLACE FUNCTION testing_db.update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for updating timestamps
CREATE TRIGGER update_series_updated_at BEFORE UPDATE ON testing_db.series FOR EACH ROW EXECUTE FUNCTION testing_db.update_updated_at_column();
CREATE TRIGGER update_matches_updated_at BEFORE UPDATE ON testing_db.matches FOR EACH ROW EXECUTE FUNCTION testing_db.update_updated_at_column();
CREATE TRIGGER update_innings_updated_at BEFORE UPDATE ON testing_db.innings FOR EACH ROW EXECUTE FUNCTION testing_db.update_updated_at_column();
CREATE TRIGGER update_overs_updated_at BEFORE UPDATE ON testing_db.scorecard_overs FOR EACH ROW EXECUTE FUNCTION testing_db.update_updated_at_column();

-- Grant permissions (adjust as needed for your setup)
-- GRANT ALL PRIVILEGES ON SCHEMA testing_db TO your_test_user;
-- GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA testing_db TO your_test_user;
-- GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA testing_db TO your_test_user;
