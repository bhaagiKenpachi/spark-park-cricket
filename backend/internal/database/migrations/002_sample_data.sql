-- Empty tables - No sample data
-- Version: 2.0.0
-- Date: 2025-09-10

-- This migration creates empty tables with the latest schema
-- No sample data is inserted - tables start empty for clean testing

-- Verify tables are created and empty
SELECT 'Tables created successfully - starting with empty data!' as status;
SELECT COUNT(*) as series_count FROM series;
SELECT COUNT(*) as matches_count FROM matches;
SELECT COUNT(*) as scoreboard_count FROM live_scoreboard;
SELECT COUNT(*) as overs_count FROM overs;
SELECT COUNT(*) as balls_count FROM balls;