-- Add match timing fields to matches table
-- This migration adds start_time and end_time columns to track match duration
-- Version: 8.0.0
-- Date: 2025-10-05

-- ============================================
-- MATCH TIMING FIELDS
-- ============================================

-- Add start_time column to track when match started
ALTER TABLE testing_db.matches 
ADD COLUMN IF NOT EXISTS start_time TIMESTAMP WITH TIME ZONE;

-- Add end_time column to track when match ended
ALTER TABLE testing_db.matches 
ADD COLUMN IF NOT EXISTS end_time TIMESTAMP WITH TIME ZONE;

-- Add comments for documentation
COMMENT ON COLUMN testing_db.matches.start_time IS 'Timestamp when the match started (set when status changes to live)';
COMMENT ON COLUMN testing_db.matches.end_time IS 'Timestamp when the match ended (set when status changes to completed)';

-- Create index on start_time for performance
CREATE INDEX IF NOT EXISTS idx_matches_start_time ON testing_db.matches(start_time);

-- Create index on end_time for performance  
CREATE INDEX IF NOT EXISTS idx_matches_end_time ON testing_db.matches(end_time);

-- Create composite index for timing queries
CREATE INDEX IF NOT EXISTS idx_matches_timing ON testing_db.matches(start_time, end_time) 
WHERE start_time IS NOT NULL OR end_time IS NOT NULL;
