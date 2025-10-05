-- Update matches status constraint to include 'not_started'
-- This migration updates the existing constraint to allow the new status

-- Drop the existing constraint
ALTER TABLE testing_db.matches DROP CONSTRAINT IF EXISTS matches_status_check;

-- Add the new constraint with 'not_started' included
ALTER TABLE testing_db.matches ADD CONSTRAINT matches_status_check 
CHECK (status IN ('not_started', 'live', 'completed', 'cancelled'));

-- Update the default value to 'not_started' for new matches
ALTER TABLE testing_db.matches ALTER COLUMN status SET DEFAULT 'not_started';
