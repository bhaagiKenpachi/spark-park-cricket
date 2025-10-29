-- Check vote_groups table data and relationships
-- Run this against testing_db schema

-- 1. Check all entries in vote_groups table
SELECT 
    vg.id as vote_group_id,
    vg.vote_id,
    v.title as vote_title,
    vg.group_id,
    g.name as group_name,
    g.type as group_type,
    vg.created_at as assignment_date
FROM testing_db.vote_groups vg
LEFT JOIN testing_db.votes v ON vg.vote_id = v.id
LEFT JOIN testing_db.groups g ON vg.group_id = g.id
ORDER BY vg.created_at DESC;

-- 2. Check if "Groups testing" vote is assigned to any group
SELECT 
    vg.id as vote_group_id,
    vg.vote_id,
    v.title as vote_title,
    vg.group_id,
    g.name as group_name,
    g.type as group_type,
    vg.created_at as assignment_date
FROM testing_db.vote_groups vg
JOIN testing_db.votes v ON vg.vote_id = v.id
LEFT JOIN testing_db.groups g ON vg.group_id = g.id
WHERE v.id = '16714960-5d83-4758-bc39-93f489e88d5f'
   OR v.title = 'Groups testing';

-- 3. Count how many votes are assigned to each group
SELECT 
    g.id as group_id,
    g.name as group_name,
    g.type as group_type,
    COUNT(vg.vote_id) as assigned_votes_count,
    STRING_AGG(v.title, ', ') as assigned_vote_titles
FROM testing_db.groups g
LEFT JOIN testing_db.vote_groups vg ON g.id = vg.group_id
LEFT JOIN testing_db.votes v ON vg.vote_id = v.id
GROUP BY g.id, g.name, g.type
ORDER BY assigned_votes_count DESC, g.name;

-- 4. Check for "spark-park" group specifically
SELECT 
    g.id as group_id,
    g.name as group_name,
    g.type as group_type,
    COUNT(vg.vote_id) as assigned_votes_count,
    STRING_AGG(v.title, ', ') as assigned_vote_titles
FROM testing_db.groups g
LEFT JOIN testing_db.vote_groups vg ON g.id = vg.group_id
LEFT JOIN testing_db.votes v ON vg.vote_id = v.id
WHERE g.name LIKE '%spark%park%' OR g.name LIKE '%spark-park%'
GROUP BY g.id, g.name, g.type;

-- 5. List all groups that exist
SELECT 
    id,
    name,
    type,
    status,
    created_at
FROM testing_db.groups
ORDER BY created_at DESC;

