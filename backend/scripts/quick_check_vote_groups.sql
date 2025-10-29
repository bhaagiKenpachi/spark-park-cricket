-- Quick check - Does vote_groups table exist and have data?

-- 1. Check if table exists
SELECT EXISTS (
    SELECT FROM information_schema.tables 
    WHERE table_schema = 'testing_db' 
    AND table_name = 'vote_groups'
) as table_exists;

-- 2. Count total entries
SELECT COUNT(*) as total_assignments FROM testing_db.vote_groups;

-- 3. Show first 10 entries (if any)
SELECT * FROM testing_db.vote_groups LIMIT 10;

-- 4. Check if "Groups testing" vote has any group assignment
SELECT vg.* 
FROM testing_db.vote_groups vg
WHERE vg.vote_id = '16714960-5d83-4758-bc39-93f489e88d5f';

