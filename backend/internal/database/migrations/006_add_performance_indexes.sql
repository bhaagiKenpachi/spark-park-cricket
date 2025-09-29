-- Migration: Add Performance Indexes for Add Ball API Optimization
-- Description: Add database indexes to optimize frequent queries in the add ball API
-- Date: 2024-01-XX
-- Author: Development Team

-- Index for balls table - optimize GetBallsByOver and GetBallCountByOver queries
-- This index will significantly speed up queries filtering by over_id
CREATE INDEX IF NOT EXISTS idx_balls_over_id ON balls(over_id);

-- Index for balls table - optimize ball number queries
-- This index will help with ordering and finding max ball numbers
CREATE INDEX IF NOT EXISTS idx_balls_over_id_ball_number ON balls(over_id, ball_number);

-- Index for overs table - optimize GetCurrentOver queries
-- This index will speed up queries filtering by innings_id and status
CREATE INDEX IF NOT EXISTS idx_overs_innings_status ON overs(innings_id, status);

-- Index for overs table - optimize GetOversByInnings queries
-- This index will help with ordering overs by over_number
CREATE INDEX IF NOT EXISTS idx_overs_innings_over_number ON overs(innings_id, over_number);

-- Index for innings table - optimize GetInningsByMatchAndNumber queries
-- This index will speed up queries filtering by match_id and innings_number
CREATE INDEX IF NOT EXISTS idx_innings_match_number ON innings(match_id, innings_number);

-- Index for matches table - optimize GetByID queries
-- This index will speed up match lookups (though primary key should already be indexed)
CREATE INDEX IF NOT EXISTS idx_matches_status ON matches(status);

-- Composite index for balls table - optimize complex queries
-- This index will help with queries that filter by over_id and ball_type
CREATE INDEX IF NOT EXISTS idx_balls_over_type ON balls(over_id, ball_type);

-- Index for balls table - optimize queries filtering by ball_type
-- This index will help with legal ball counting
CREATE INDEX IF NOT EXISTS idx_balls_ball_type ON balls(ball_type);

-- Add comments for documentation
COMMENT ON INDEX idx_balls_over_id IS 'Optimizes GetBallsByOver and GetBallCountByOver queries';
COMMENT ON INDEX idx_balls_over_id_ball_number IS 'Optimizes ball number ordering and max ball number queries';
COMMENT ON INDEX idx_overs_innings_status IS 'Optimizes GetCurrentOver queries';
COMMENT ON INDEX idx_overs_innings_over_number IS 'Optimizes GetOversByInnings ordering';
COMMENT ON INDEX idx_innings_match_number IS 'Optimizes GetInningsByMatchAndNumber queries';
COMMENT ON INDEX idx_matches_status IS 'Optimizes match status filtering';
COMMENT ON INDEX idx_balls_over_type IS 'Optimizes complex ball queries with type filtering';
COMMENT ON INDEX idx_balls_ball_type IS 'Optimizes legal ball counting queries';
