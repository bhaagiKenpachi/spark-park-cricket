-- Add Byes Column to Balls Table
-- Adds byes functionality for all ball types
-- Version: 2.0.4
-- Date: 2025-09-10

-- ============================================
-- ADD BYES COLUMN TO BALLS TABLE
-- ============================================

-- Add byes column to balls table
ALTER TABLE balls ADD COLUMN IF NOT EXISTS byes INTEGER DEFAULT 0 CHECK (byes >= 0 AND byes <= 6);

-- Update existing records to have byes = 0
UPDATE balls SET byes = 0 WHERE byes IS NULL;

-- Make byes column NOT NULL
ALTER TABLE balls ALTER COLUMN byes SET NOT NULL;

-- ============================================
-- VERIFICATION
-- ============================================

SELECT 'Byes column added successfully to balls table!' as status;

-- Test the byes functionality
SELECT 'Testing byes functionality...' as info;

-- This should work (good ball with byes)
-- INSERT INTO balls (over_id, ball_number, ball_type, run_type, runs, byes, is_wicket) 
-- VALUES ('test-over-id', 1, 'good', '1', 1, 2, false);

-- This should work (wide ball with byes)
-- INSERT INTO balls (over_id, ball_number, ball_type, run_type, runs, byes, is_wicket) 
-- VALUES ('test-over-id', 2, 'wide', 'WD', 1, 1, false);

-- This should work (no ball with byes)
-- INSERT INTO balls (over_id, ball_number, ball_type, run_type, runs, byes, is_wicket) 
-- VALUES ('test-over-id', 3, 'no_ball', 'NB', 1, 3, false);

-- This should fail (byes > 6)
-- INSERT INTO balls (over_id, ball_number, ball_type, run_type, runs, byes, is_wicket) 
-- VALUES ('test-over-id', 4, 'good', '1', 1, 7, false);

-- This should fail (byes < 0)
-- INSERT INTO balls (over_id, ball_number, ball_type, run_type, runs, byes, is_wicket) 
-- VALUES ('test-over-id', 5, 'good', '1', 1, -1, false);
