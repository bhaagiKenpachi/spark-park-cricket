-- Migration: Add created_by fields to series and matches tables in testing_db schema
-- Description: Adds user ownership tracking to series and matches in the testing_db schema

-- Add created_by field to series table in testing_db schema
ALTER TABLE testing_db.series ADD COLUMN IF NOT EXISTS created_by UUID REFERENCES testing_db.users(id) ON DELETE SET NULL;

-- Add created_by field to matches table in testing_db schema
ALTER TABLE testing_db.matches ADD COLUMN IF NOT EXISTS created_by UUID REFERENCES testing_db.users(id) ON DELETE SET NULL;

-- Create indexes for better performance on created_by fields in testing_db schema
CREATE INDEX IF NOT EXISTS idx_testing_db_series_created_by ON testing_db.series(created_by);
CREATE INDEX IF NOT EXISTS idx_testing_db_matches_created_by ON testing_db.matches(created_by);

-- Add comments for documentation
COMMENT ON COLUMN testing_db.series.created_by IS 'User who created this series';
COMMENT ON COLUMN testing_db.matches.created_by IS 'User who created this match';
