-- Add time tracking fields to innings and overs tables
-- This migration adds start_time and end_time fields to track duration

-- Add time tracking fields to innings table
ALTER TABLE dev_v1.innings 
ADD COLUMN IF NOT EXISTS start_time TIMESTAMP WITH TIME ZONE,
ADD COLUMN IF NOT EXISTS end_time TIMESTAMP WITH TIME ZONE,
ADD COLUMN IF NOT EXISTS duration_seconds INTEGER DEFAULT 0 CHECK (duration_seconds >= 0);

-- Add time tracking fields to overs table  
ALTER TABLE dev_v1.overs 
ADD COLUMN IF NOT EXISTS start_time TIMESTAMP WITH TIME ZONE,
ADD COLUMN IF NOT EXISTS end_time TIMESTAMP WITH TIME ZONE,
ADD COLUMN IF NOT EXISTS duration_seconds INTEGER DEFAULT 0 CHECK (duration_seconds >= 0);

-- Add indexes for time-based queries
CREATE INDEX IF NOT EXISTS idx_innings_start_time ON dev_v1.innings(start_time);
CREATE INDEX IF NOT EXISTS idx_innings_end_time ON dev_v1.innings(end_time);
CREATE INDEX IF NOT EXISTS idx_overs_start_time ON dev_v1.overs(start_time);
CREATE INDEX IF NOT EXISTS idx_overs_end_time ON dev_v1.overs(end_time);

-- Add comments for documentation
COMMENT ON COLUMN dev_v1.innings.start_time IS 'Timestamp when innings started';
COMMENT ON COLUMN dev_v1.innings.end_time IS 'Timestamp when innings ended';
COMMENT ON COLUMN dev_v1.innings.duration_seconds IS 'Duration of innings in seconds';
COMMENT ON COLUMN dev_v1.overs.start_time IS 'Timestamp when over started';
COMMENT ON COLUMN dev_v1.overs.end_time IS 'Timestamp when over ended';
COMMENT ON COLUMN dev_v1.overs.duration_seconds IS 'Duration of over in seconds';