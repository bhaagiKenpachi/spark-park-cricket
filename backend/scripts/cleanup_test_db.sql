-- Test Database Cleanup Script
-- This script cleans up the testing_db schema after running tests

-- Set search path to testing_db
SET search_path TO testing_db;

-- Truncate all tables in reverse dependency order
TRUNCATE TABLE testing_db.scorecard_balls CASCADE;
TRUNCATE TABLE testing_db.scorecard_overs CASCADE;
TRUNCATE TABLE testing_db.innings CASCADE;
TRUNCATE TABLE testing_db.matches CASCADE;
TRUNCATE TABLE testing_db.series CASCADE;

-- Reset sequences (if any)
-- Note: UUID primary keys don't use sequences, but if you have any integer sequences, reset them here

-- Optional: Drop the entire schema (uncomment if you want to completely remove test data)
-- DROP SCHEMA IF EXISTS testing_db CASCADE;
