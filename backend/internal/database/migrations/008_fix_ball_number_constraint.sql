-- Fix ball number constraint to allow more than 6 balls per over
-- This is needed because illegal balls (wides, no balls) don't count towards over completion
-- but still need to be stored in the database

-- Drop the existing constraint
ALTER TABLE balls DROP CONSTRAINT IF EXISTS balls_ball_number_check;

-- Add new constraint that allows up to 20 balls per over (reasonable limit)
-- This accounts for scenarios with many wides/no balls
ALTER TABLE balls ADD CONSTRAINT balls_ball_number_check CHECK (ball_number >= 1 AND ball_number <= 20);

-- Add comment explaining the constraint
COMMENT ON CONSTRAINT balls_ball_number_check ON balls IS 'Ball number constraint: 1-20 to allow for illegal balls (wides, no balls) that dont count towards over completion';

-- Verify the constraint was updated
SELECT 'Ball number constraint updated successfully!' as status;
