-- Migration: Update wicket constraint from 10 to 20
-- This migration updates the innings table constraint to allow up to 20 wickets
-- Date: 2025-10-13
-- Environment: Development (dev_v1 schema)

-- Drop the old constraint
ALTER TABLE dev_v1.innings DROP CONSTRAINT IF EXISTS innings_total_wickets_check;

-- Add the new constraint allowing up to 20 wickets
ALTER TABLE dev_v1.innings ADD CONSTRAINT innings_total_wickets_check CHECK (total_wickets >= 0 AND total_wickets <= 20);

-- Update the comment
COMMENT ON COLUMN dev_v1.innings.total_wickets IS 'Total wickets fallen in this innings (0-20)';

-- Verification query (uncomment to check after migration):
-- SELECT constraint_name, check_clause 
-- FROM information_schema.check_constraints 
-- WHERE constraint_name = 'innings_total_wickets_check';

