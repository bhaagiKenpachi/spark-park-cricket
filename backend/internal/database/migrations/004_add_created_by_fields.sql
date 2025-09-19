-- Migration: Add created_by fields to series and matches tables
-- Description: Adds user ownership tracking to series and matches

-- Add created_by field to series table
ALTER TABLE series ADD COLUMN IF NOT EXISTS created_by UUID REFERENCES users(id) ON DELETE SET NULL;

-- Add created_by field to matches table
ALTER TABLE matches ADD COLUMN IF NOT EXISTS created_by UUID REFERENCES users(id) ON DELETE SET NULL;

-- Create indexes for better performance on created_by fields
CREATE INDEX IF NOT EXISTS idx_series_created_by ON series(created_by);
CREATE INDEX IF NOT EXISTS idx_matches_created_by ON matches(created_by);

-- Add comments for documentation
COMMENT ON COLUMN series.created_by IS 'User who created this series';
COMMENT ON COLUMN matches.created_by IS 'User who created this match';
