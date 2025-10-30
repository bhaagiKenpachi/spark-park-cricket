-- Migration: Add team_formation_enabled column to votes table (Testing DB)
-- Description: Adds a boolean column to control whether team formation is enabled for a vote
-- Date: 2025-01-04
-- Database: prod_v1

-- Add team_formation_enabled column to votes table (only if it doesn't exist)
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_schema = 'prod_v1' 
        AND table_name = 'votes' 
        AND column_name = 'team_formation_enabled'
    ) THEN
        ALTER TABLE prod_v1.votes 
        ADD COLUMN team_formation_enabled BOOLEAN NOT NULL DEFAULT true;
        
        -- Add comment to explain the column
        COMMENT ON COLUMN prod_v1.votes.team_formation_enabled IS 'Controls whether team formation is enabled for this vote. When true, users can create and manage teams from voted users.';
    END IF;
END $$;

-- Update existing votes to have team formation enabled by default
-- This ensures backward compatibility with existing team formation data
UPDATE prod_v1.votes SET team_formation_enabled = true WHERE team_formation_enabled IS NULL;
