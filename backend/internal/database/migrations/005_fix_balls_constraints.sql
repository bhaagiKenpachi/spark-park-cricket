-- Fix Balls Table Constraints
-- Fixes the wicket_type constraint issue
-- Version: 2.0.3
-- Date: 2025-09-10

-- ============================================
-- FIX BALLS TABLE CONSTRAINTS
-- ============================================

-- Drop existing constraint if it exists
ALTER TABLE balls DROP CONSTRAINT IF EXISTS balls_wicket_type_check;

-- Update run_type constraint to include 'WC' (wicket)
ALTER TABLE balls DROP CONSTRAINT IF EXISTS balls_run_type_check;
ALTER TABLE balls ADD CONSTRAINT balls_run_type_check CHECK (run_type IN ('0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'NB', 'WD', 'LB', 'WC'));

-- Add proper wicket_type constraint
ALTER TABLE balls ADD CONSTRAINT balls_wicket_type_check CHECK (
    (is_wicket = true AND wicket_type IS NOT NULL) OR 
    (is_wicket = false AND wicket_type IS NULL)
);

-- ============================================
-- VERIFICATION
-- ============================================

SELECT 'Balls table constraints fixed successfully!' as status;

-- Test the constraints
SELECT 'Testing constraints...' as info;

-- This should work (wicket with wicket_type)
-- INSERT INTO balls (over_id, ball_number, ball_type, run_type, runs, is_wicket, wicket_type) 
-- VALUES ('test-over-id', 1, 'good', 'WC', 0, true, 'bowled');

-- This should work (no wicket, no wicket_type)
-- INSERT INTO balls (over_id, ball_number, ball_type, run_type, runs, is_wicket) 
-- VALUES ('test-over-id', 2, 'good', '1', 1, false);

-- This should fail (wicket without wicket_type)
-- INSERT INTO balls (over_id, ball_number, ball_type, run_type, runs, is_wicket) 
-- VALUES ('test-over-id', 3, 'good', 'WC', 0, true);

-- This should fail (no wicket with wicket_type)
-- INSERT INTO balls (over_id, ball_number, ball_type, run_type, runs, is_wicket, wicket_type) 
-- VALUES ('test-over-id', 4, 'good', '1', 1, false, 'bowled');
