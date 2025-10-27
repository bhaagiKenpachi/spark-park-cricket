-- Migration: Create fall_of_wickets table
-- Description: Creates a simplified table to track wicket falls with only essential information

CREATE TABLE IF NOT EXISTS testing_db.fall_of_wickets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    match_id UUID NOT NULL REFERENCES testing_db.matches(id) ON DELETE CASCADE,
    innings_id UUID NOT NULL REFERENCES testing_db.innings(id) ON DELETE CASCADE,
    over_id UUID NOT NULL REFERENCES testing_db.overs(id) ON DELETE CASCADE,
    ball_id UUID NOT NULL REFERENCES testing_db.balls(id) ON DELETE CASCADE,
    wicket_number INTEGER NOT NULL CHECK (wicket_number >= 1 AND wicket_number <= 20),
    score INTEGER NOT NULL CHECK (score >= 0),
    over_number INTEGER NOT NULL CHECK (over_number >= 1),
    ball_number INTEGER NOT NULL CHECK (ball_number >= 1 AND ball_number <= 20),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Ensure unique wicket numbers per innings
    UNIQUE(innings_id, wicket_number),
    
    -- Ensure unique ball per wicket (one wicket per ball)
    UNIQUE(ball_id)
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_fall_of_wickets_match_id ON testing_db.fall_of_wickets(match_id);
CREATE INDEX IF NOT EXISTS idx_fall_of_wickets_innings_id ON testing_db.fall_of_wickets(innings_id);
CREATE INDEX IF NOT EXISTS idx_fall_of_wickets_over_id ON testing_db.fall_of_wickets(over_id);
CREATE INDEX IF NOT EXISTS idx_fall_of_wickets_ball_id ON testing_db.fall_of_wickets(ball_id);
CREATE INDEX IF NOT EXISTS idx_fall_of_wickets_wicket_number ON testing_db.fall_of_wickets(wicket_number);
CREATE INDEX IF NOT EXISTS idx_fall_of_wickets_created_at ON testing_db.fall_of_wickets(created_at);

-- Add comments for documentation
COMMENT ON TABLE testing_db.fall_of_wickets IS 'Tracks wicket falls with score and over position information';
COMMENT ON COLUMN testing_db.fall_of_wickets.wicket_number IS 'Sequential wicket number within the innings (1-20)';
COMMENT ON COLUMN testing_db.fall_of_wickets.score IS 'Total score when the wicket fell';
COMMENT ON COLUMN testing_db.fall_of_wickets.over_number IS 'Over number when wicket fell';
COMMENT ON COLUMN testing_db.fall_of_wickets.ball_number IS 'Ball number within the over when wicket fell';