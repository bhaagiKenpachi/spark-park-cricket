-- Migration: Update wicket constraint from 10 to 20
-- This migration updates the innings table constraint to allow up to 20 wickets

-- Drop the old constraint
ALTER TABLE testing_db.innings DROP CONSTRAINT IF EXISTS innings_total_wickets_check;

-- Add the new constraint allowing up to 20 wickets
ALTER TABLE testing_db.innings ADD CONSTRAINT innings_total_wickets_check CHECK (total_wickets >= 0 AND total_wickets <= 20);

-- Update the comment
COMMENT ON COLUMN testing_db.innings.total_wickets IS 'Total wickets fallen in this innings (0-20)';

