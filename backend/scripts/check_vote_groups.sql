-- Script to check vote_groups table and relationships
-- Run this against your testing_db schema

-- 1. Check if vote_groups table exists
SELECT 
    table_name,
    table_schema
FROM information_schema.tables 
WHERE table_schema = 'testing_db' 
  AND table_name = 'vote_groups';

-- 2. Check all vote_groups relationships
SELECT 
    vg.id as vote_group_id,
    vg.vote_id,
    v.title as vote_title,
    vg.group_id,
    g.name as group_name,
    vg.created_at
FROM testing_db.vote_groups vg
LEFT JOIN testing_db.votes v ON vg.vote_id = v.id
LEFT JOIN testing_db.groups g ON vg.group_id = g.id
ORDER BY vg.created_at DESC;

-- 3. Count votes per group
SELECT 
    g.id as group_id,
    g.name as group_name,
    COUNT(vg.vote_id) as vote_count
FROM testing_db.groups g
LEFT JOIN testing_db.vote_groups vg ON g.id = vg.group_id
GROUP BY g.id, g.name
ORDER BY vote_count DESC, g.name;

-- 4. Find votes without any group assignment
SELECT 
    v.id,
    v.title,
    v.created_at
FROM testing_db.votes v
LEFT JOIN testing_db.vote_groups vg ON v.id = vg.vote_id
WHERE vg.id IS NULL
ORDER BY v.created_at DESC;

-- 5. Find specific vote-group assignment (replace with actual IDs)
-- SELECT * FROM testing_db.vote_groups 
-- WHERE vote_id = 'YOUR_VOTE_ID' OR group_id = 'YOUR_GROUP_ID';

