-- Migration: Add innings_number column to fall_of_wickets table
-- Description: Adds innings_number column to track which innings (1 or 2) the wicket fell in

-- Add innings_number column
ALTER TABLE prod_v1.fall_of_wickets 
ADD COLUMN innings_number INTEGER CHECK (innings_number >= 1 AND innings_number <= 2);

-- Update existing records with innings_number based on innings_id
-- This will populate the innings_number for existing records
UPDATE prod_v1.fall_of_wickets 
SET innings_number = (
    SELECT innings_number 
    FROM prod_v1.innings 
    WHERE prod_v1.innings.id = prod_v1.fall_of_wickets.innings_id
);

-- Make innings_number NOT NULL after populating existing records
ALTER TABLE prod_v1.fall_of_wickets 
ALTER COLUMN innings_number SET NOT NULL;

-- Add index for better performance
CREATE INDEX IF NOT EXISTS idx_fall_of_wickets_innings_number ON prod_v1.fall_of_wickets(innings_number);

-- Add comment for documentation
COMMENT ON COLUMN prod_v1.fall_of_wickets.innings_number IS 'Innings number (1 or 2) when the wicket fell';
