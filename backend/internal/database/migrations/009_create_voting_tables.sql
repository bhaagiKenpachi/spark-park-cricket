-- Spark Park Cricket - Voting System Migration
-- Version: 009
-- Date: 2025-01-27
-- Description: Create voting system tables for polls and user voting

-- ============================================
-- VOTING SYSTEM TABLES
-- ============================================

-- Create votes table
CREATE TABLE IF NOT EXISTS testing_db.votes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('single', 'multiple')),
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'closed', 'cancelled')),
    created_by UUID NOT NULL REFERENCES testing_db.users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    closed_at TIMESTAMP WITH TIME ZONE NULL
);

-- Create vote_options table
CREATE TABLE IF NOT EXISTS testing_db.vote_options (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vote_id UUID NOT NULL REFERENCES testing_db.votes(id) ON DELETE CASCADE,
    text VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create user_votes table
CREATE TABLE IF NOT EXISTS testing_db.user_votes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vote_id UUID NOT NULL REFERENCES testing_db.votes(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES testing_db.users(id) ON DELETE CASCADE,
    selected_options UUID[] NOT NULL DEFAULT '{}',
    voted_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(vote_id, user_id) -- Ensure one vote per user per poll
);

-- ============================================
-- INDEXES FOR PERFORMANCE
-- ============================================

-- Indexes for votes table
CREATE INDEX IF NOT EXISTS idx_votes_status ON testing_db.votes(status);
CREATE INDEX IF NOT EXISTS idx_votes_created_by ON testing_db.votes(created_by);
CREATE INDEX IF NOT EXISTS idx_votes_created_at ON testing_db.votes(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_votes_type ON testing_db.votes(type);

-- Indexes for vote_options table
CREATE INDEX IF NOT EXISTS idx_vote_options_vote_id ON testing_db.vote_options(vote_id);

-- Indexes for user_votes table
CREATE INDEX IF NOT EXISTS idx_user_votes_vote_id ON testing_db.user_votes(vote_id);
CREATE INDEX IF NOT EXISTS idx_user_votes_user_id ON testing_db.user_votes(user_id);
CREATE INDEX IF NOT EXISTS idx_user_votes_voted_at ON testing_db.user_votes(voted_at DESC);

-- ============================================
-- TRIGGERS FOR UPDATED_AT
-- ============================================

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Triggers for votes table
CREATE TRIGGER update_votes_updated_at 
    BEFORE UPDATE ON testing_db.votes 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Triggers for vote_options table
CREATE TRIGGER update_vote_options_updated_at 
    BEFORE UPDATE ON testing_db.vote_options 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ============================================
-- COMMENTS
-- ============================================

COMMENT ON TABLE testing_db.votes IS 'Stores voting polls with title, description, and type';
COMMENT ON TABLE testing_db.vote_options IS 'Stores options for each vote poll';
COMMENT ON TABLE testing_db.user_votes IS 'Stores user voting selections with selected options';

COMMENT ON COLUMN testing_db.votes.type IS 'Vote type: single (one option) or multiple (multiple options)';
COMMENT ON COLUMN testing_db.votes.status IS 'Vote status: active, closed, or cancelled';
COMMENT ON COLUMN testing_db.user_votes.selected_options IS 'Array of selected option IDs';
