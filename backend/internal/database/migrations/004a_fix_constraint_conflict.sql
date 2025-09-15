-- Fix Constraint Conflict
-- Drops conflicting constraints before migration 004 can run
-- Version: 2.0.2a
-- Date: 2025-09-15

-- ============================================
-- DROP CONFLICTING CONSTRAINTS
-- ============================================

-- Drop the balls_wicket_type_check constraint if it exists
-- This fixes the issue where migration 004 fails due to existing constraint
ALTER TABLE balls DROP CONSTRAINT IF EXISTS balls_wicket_type_check;

-- ============================================
-- VERIFICATION
-- ============================================

SELECT 'Constraint conflict resolved successfully!' as status;
