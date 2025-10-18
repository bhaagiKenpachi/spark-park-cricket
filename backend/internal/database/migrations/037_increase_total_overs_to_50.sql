-- Migration: Increase total_overs limit from 20 to 50
-- Purpose: Support longer format cricket matches (50 overs per side)
-- Date: 2025-10-18

-- Drop existing constraint on dev_v1.matches
ALTER TABLE dev_v1.matches DROP CONSTRAINT IF EXISTS matches_total_overs_check;

-- Add new constraint allowing up to 50 overs
ALTER TABLE dev_v1.matches ADD CONSTRAINT matches_total_overs_check 
    CHECK (total_overs >= 1 AND total_overs <= 50);

-- Update comment
COMMENT ON COLUMN dev_v1.matches.total_overs IS 'Total overs for the match (1-50)';

-- Update default value remains 20 for backward compatibility
ALTER TABLE dev_v1.matches ALTER COLUMN total_overs SET DEFAULT 20;

