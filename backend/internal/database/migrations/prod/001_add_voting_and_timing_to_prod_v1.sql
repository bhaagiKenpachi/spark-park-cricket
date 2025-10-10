-- ============================================
-- MIGRATION 001: Add Voting System and Match Timing to prod_v1
-- ============================================
-- Version: 001
-- Date: 2025-10-08
-- Description: Adds voting tables (votes, vote_options, user_votes) 
--              and match timing fields (start_time, end_time) to prod_v1 schema
-- Safe to run multiple times - uses IF NOT EXISTS checks

-- First, ensure prod_v1 schema exists
CREATE SCHEMA IF NOT EXISTS prod_v1;

-- Grant permissions
GRANT USAGE ON SCHEMA prod_v1 TO anon, authenticated, service_role, postgres;

-- ============================================
-- ADD MISSING COLUMNS TO EXISTING TABLES
-- ============================================

-- Add start_time and end_time to matches if they don't exist
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_schema = 'prod_v1' 
        AND table_name = 'matches' 
        AND column_name = 'start_time'
    ) THEN
        ALTER TABLE prod_v1.matches ADD COLUMN start_time TIMESTAMP WITH TIME ZONE;
    END IF;
    
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_schema = 'prod_v1' 
        AND table_name = 'matches' 
        AND column_name = 'end_time'
    ) THEN
        ALTER TABLE prod_v1.matches ADD COLUMN end_time TIMESTAMP WITH TIME ZONE;
    END IF;
END $$;

-- Update matches status constraint to include 'not_started'
DO $$
BEGIN
    -- Drop existing constraint if it exists
    ALTER TABLE prod_v1.matches DROP CONSTRAINT IF EXISTS matches_status_check;
    
    -- Add new constraint
    ALTER TABLE prod_v1.matches ADD CONSTRAINT matches_status_check 
    CHECK (status IN ('not_started', 'live', 'completed', 'cancelled'));
    
    -- Update default value
    ALTER TABLE prod_v1.matches ALTER COLUMN status SET DEFAULT 'not_started';
END $$;

-- ============================================
-- CREATE VOTING TABLES IF THEY DON'T EXIST
-- ============================================

-- Create votes table
CREATE TABLE IF NOT EXISTS prod_v1.votes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('single', 'multiple')),
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'closed', 'cancelled')),
    created_by UUID NOT NULL REFERENCES prod_v1.users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    closed_at TIMESTAMP WITH TIME ZONE NULL
);

-- Create vote_options table
CREATE TABLE IF NOT EXISTS prod_v1.vote_options (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vote_id UUID NOT NULL REFERENCES prod_v1.votes(id) ON DELETE CASCADE,
    text VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create user_votes table
CREATE TABLE IF NOT EXISTS prod_v1.user_votes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vote_id UUID NOT NULL REFERENCES prod_v1.votes(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES prod_v1.users(id) ON DELETE CASCADE,
    selected_options UUID[] NOT NULL DEFAULT '{}',
    voted_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(vote_id, user_id)
);

-- ============================================
-- CREATE INDEXES
-- ============================================

-- Match timing indexes
CREATE INDEX IF NOT EXISTS idx_prod_v1_matches_start_time ON prod_v1.matches(start_time);
CREATE INDEX IF NOT EXISTS idx_prod_v1_matches_end_time ON prod_v1.matches(end_time);
CREATE INDEX IF NOT EXISTS idx_prod_v1_matches_timing ON prod_v1.matches(start_time, end_time) 
WHERE start_time IS NOT NULL OR end_time IS NOT NULL;

-- Voting indexes
CREATE INDEX IF NOT EXISTS idx_prod_v1_votes_status ON prod_v1.votes(status);
CREATE INDEX IF NOT EXISTS idx_prod_v1_votes_created_by ON prod_v1.votes(created_by);
CREATE INDEX IF NOT EXISTS idx_prod_v1_votes_created_at ON prod_v1.votes(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_prod_v1_votes_type ON prod_v1.votes(type);
CREATE INDEX IF NOT EXISTS idx_prod_v1_vote_options_vote_id ON prod_v1.vote_options(vote_id);
CREATE INDEX IF NOT EXISTS idx_prod_v1_user_votes_vote_id ON prod_v1.user_votes(vote_id);
CREATE INDEX IF NOT EXISTS idx_prod_v1_user_votes_user_id ON prod_v1.user_votes(user_id);
CREATE INDEX IF NOT EXISTS idx_prod_v1_user_votes_voted_at ON prod_v1.user_votes(voted_at DESC);

-- ============================================
-- CREATE TRIGGERS
-- ============================================

-- Ensure function exists
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Voting table triggers
DROP TRIGGER IF EXISTS update_prod_v1_votes_updated_at ON prod_v1.votes;
CREATE TRIGGER update_prod_v1_votes_updated_at 
    BEFORE UPDATE ON prod_v1.votes 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_prod_v1_vote_options_updated_at ON prod_v1.vote_options;
CREATE TRIGGER update_prod_v1_vote_options_updated_at 
    BEFORE UPDATE ON prod_v1.vote_options 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ============================================
-- ADD COMMENTS
-- ============================================

COMMENT ON TABLE prod_v1.votes IS 'Stores voting polls with title, description, and type';
COMMENT ON TABLE prod_v1.vote_options IS 'Stores options for each vote poll';
COMMENT ON TABLE prod_v1.user_votes IS 'Stores user voting selections with selected options';

COMMENT ON COLUMN prod_v1.votes.type IS 'Vote type: single (one option) or multiple (multiple options)';
COMMENT ON COLUMN prod_v1.votes.status IS 'Vote status: active, closed, or cancelled';
COMMENT ON COLUMN prod_v1.user_votes.selected_options IS 'Array of selected option IDs';

COMMENT ON COLUMN prod_v1.matches.start_time IS 'Timestamp when the match started (set when status changes to live)';
COMMENT ON COLUMN prod_v1.matches.end_time IS 'Timestamp when the match ended (set when status changes to completed)';

-- ============================================
-- GRANT PERMISSIONS
-- ============================================

GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA prod_v1 TO anon, authenticated, service_role, postgres;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA prod_v1 TO anon, authenticated, service_role, postgres;

-- ============================================
-- VERIFICATION
-- ============================================

SELECT 'Migration completed successfully!' as status;

SELECT 'Voting tables:' as info;
SELECT table_name FROM information_schema.tables 
WHERE table_schema = 'prod_v1' 
AND table_name IN ('votes', 'vote_options', 'user_votes')
ORDER BY table_name;
