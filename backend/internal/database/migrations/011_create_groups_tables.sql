-- Spark Park Cricket - Groups Migration
-- Version: 011
-- Date: 2025-01-27
-- Description: Create groups and group members tables for flexible group voting

-- ============================================
-- GROUPS SYSTEM TABLES
-- ============================================

-- Create groups table
CREATE TABLE IF NOT EXISTS testing_db.groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('custom', 'team', 'series', 'match', 'location', 'skill')),
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'archived')),
    created_by UUID NOT NULL REFERENCES testing_db.users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create group_members table
CREATE TABLE IF NOT EXISTS testing_db.group_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    group_id UUID NOT NULL REFERENCES testing_db.groups(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES testing_db.users(id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL DEFAULT 'member' CHECK (role IN ('member', 'admin', 'moderator')),
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Ensure one user can only be in a group once
    UNIQUE(group_id, user_id)
);

-- Create vote_groups table (many-to-many relationship between votes and groups)
CREATE TABLE IF NOT EXISTS testing_db.vote_groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vote_id UUID NOT NULL REFERENCES testing_db.votes(id) ON DELETE CASCADE,
    group_id UUID NOT NULL REFERENCES testing_db.groups(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Ensure unique vote-group combination
    UNIQUE(vote_id, group_id)
);

-- ============================================
-- INDEXES FOR PERFORMANCE
-- ============================================

-- Indexes for groups table
CREATE INDEX IF NOT EXISTS idx_groups_status ON testing_db.groups(status);
CREATE INDEX IF NOT EXISTS idx_groups_type ON testing_db.groups(type);
CREATE INDEX IF NOT EXISTS idx_groups_created_by ON testing_db.groups(created_by);
CREATE INDEX IF NOT EXISTS idx_groups_name ON testing_db.groups(name);

-- Indexes for group_members table
CREATE INDEX IF NOT EXISTS idx_group_members_group_id ON testing_db.group_members(group_id);
CREATE INDEX IF NOT EXISTS idx_group_members_user_id ON testing_db.group_members(user_id);
CREATE INDEX IF NOT EXISTS idx_group_members_role ON testing_db.group_members(role);

-- Indexes for vote_groups table
CREATE INDEX IF NOT EXISTS idx_vote_groups_vote_id ON testing_db.vote_groups(vote_id);
CREATE INDEX IF NOT EXISTS idx_vote_groups_group_id ON testing_db.vote_groups(group_id);

-- ============================================
-- TRIGGERS FOR AUTOMATIC TIMESTAMP UPDATES
-- ============================================

-- Trigger for groups table updated_at
CREATE OR REPLACE FUNCTION testing_db.update_groups_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_groups_updated_at
    BEFORE UPDATE ON testing_db.groups
    FOR EACH ROW
    EXECUTE FUNCTION testing_db.update_groups_updated_at();

-- ============================================
-- SAMPLE DATA FOR TESTING
-- ============================================

-- Insert sample groups (only if they don't exist)
INSERT INTO testing_db.groups (id, name, description, type, created_by)
SELECT 
    gen_random_uuid(),
    'Team A Players',
    'Players assigned to Team A for voting',
    'team',
    (SELECT id FROM testing_db.users LIMIT 1)
WHERE NOT EXISTS (SELECT 1 FROM testing_db.groups WHERE name = 'Team A Players');

INSERT INTO testing_db.groups (id, name, description, type, created_by)
SELECT 
    gen_random_uuid(),
    'Team B Players',
    'Players assigned to Team B for voting',
    'team',
    (SELECT id FROM testing_db.users LIMIT 1)
WHERE NOT EXISTS (SELECT 1 FROM testing_db.groups WHERE name = 'Team B Players');

INSERT INTO testing_db.groups (id, name, description, type, created_by)
SELECT 
    gen_random_uuid(),
    'Cricket Enthusiasts',
    'General group for cricket fans and enthusiasts',
    'custom',
    (SELECT id FROM testing_db.users LIMIT 1)
WHERE NOT EXISTS (SELECT 1 FROM testing_db.groups WHERE name = 'Cricket Enthusiasts');

INSERT INTO testing_db.groups (id, name, description, type, created_by)
SELECT 
    gen_random_uuid(),
    'Beginner Players',
    'Group for players with beginner level skills',
    'skill',
    (SELECT id FROM testing_db.users LIMIT 1)
WHERE NOT EXISTS (SELECT 1 FROM testing_db.groups WHERE name = 'Beginner Players');

INSERT INTO testing_db.groups (id, name, description, type, created_by)
SELECT 
    gen_random_uuid(),
    'Advanced Players',
    'Group for players with advanced level skills',
    'skill',
    (SELECT id FROM testing_db.users LIMIT 1)
WHERE NOT EXISTS (SELECT 1 FROM testing_db.groups WHERE name = 'Advanced Players');
