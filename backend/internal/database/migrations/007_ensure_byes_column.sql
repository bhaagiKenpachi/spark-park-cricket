-- ============================================
-- ENSURE BYES COLUMN EXISTS
-- ============================================

-- Add byes column if it doesn't exist
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'balls' AND column_name = 'byes'
    ) THEN
        ALTER TABLE balls ADD COLUMN byes INTEGER DEFAULT 0 CHECK (byes >= 0 AND byes <= 6);
    END IF;
END $$;

-- Add comment to byes column
COMMENT ON COLUMN balls.byes IS 'Additional runs from byes (0-6)';
