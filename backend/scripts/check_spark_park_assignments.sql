-- Check if "Spark-park" group (id: 838f9d9d-fd26-4dd1-8493-1a0ef098ddb7) has any vote assignments

-- 1. Check all vote assignments for Spark-park group
SELECT 
    vg.id as vote_group_id,
    vg.vote_id,
    v.title as vote_title,
    v.status as vote_status,
    vg.group_id,
    g.name as group_name,
    vg.created_at as assignment_date
FROM testing_db.vote_groups vg
JOIN testing_db.votes v ON vg.vote_id = v.id
JOIN testing_db.groups g ON vg.group_id = g.id
WHERE g.id = '838f9d9d-fd26-4dd1-8493-1a0ef098ddb7'
   OR g.name ILIKE '%spark%park%'
ORDER BY vg.created_at DESC;

-- 2. Check if "Groups testing" vote is assigned to any group
SELECT 
    vg.id as vote_group_id,
    vg.vote_id,
    v.title as vote_title,
    vg.group_id,
    g.name as group_name,
    vg.created_at as assignment_date
FROM testing_db.vote_groups vg
JOIN testing_db.votes v ON vg.vote_id = v.id
LEFT JOIN testing_db.groups g ON vg.group_id = g.id
WHERE v.id = '16714960-5d83-4758-bc39-93f489e88d5f'
   OR v.title ILIKE '%groups testing%';

-- 3. Total count of vote assignments for Spark-park
SELECT 
    COUNT(*) as total_assignments
FROM testing_db.vote_groups
WHERE group_id = '838f9d9d-fd26-4dd1-8493-1a0ef098ddb7';

-- 4. Check if vote_groups table has ANY data at all
SELECT COUNT(*) as total_vote_groups_entries FROM testing_db.vote_groups;

-- 5. Show all vote-group assignments (if any exist)
SELECT 
    vg.*,
    v.title as vote_title,
    g.name as group_name
FROM testing_db.vote_groups vg
LEFT JOIN testing_db.votes v ON vg.vote_id = v.id
LEFT JOIN testing_db.groups g ON vg.group_id = g.id
ORDER BY vg.created_at DESC
LIMIT 20;

