-- Spark Park Cricket - Groups Migration
-- Version: 011
-- Date: 2025-01-27
-- Description: Create groups and group members tables for flexible group voting

-- ============================================
-- GROUPS SYSTEM TABLES
-- ============================================

-- Create groups table
CREATE TABLE IF NOT EXISTS prod_v1.groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('custom', 'team', 'series', 'match', 'location', 'skill')),
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'archived')),
    created_by UUID NOT NULL REFERENCES prod_v1.users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create group_members table
CREATE TABLE IF NOT EXISTS prod_v1.group_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    group_id UUID NOT NULL REFERENCES prod_v1.groups(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES prod_v1.users(id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL DEFAULT 'member' CHECK (role IN ('member', 'admin', 'moderator')),
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Ensure one user can only be in a group once
    UNIQUE(group_id, user_id)
);

-- Create vote_groups table (many-to-many relationship between votes and groups)
CREATE TABLE IF NOT EXISTS prod_v1.vote_groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vote_id UUID NOT NULL REFERENCES prod_v1.votes(id) ON DELETE CASCADE,
    group_id UUID NOT NULL REFERENCES prod_v1.groups(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Ensure unique vote-group combination
    UNIQUE(vote_id, group_id)
);

-- ============================================
-- INDEXES FOR PERFORMANCE
-- ============================================

-- Indexes for groups table
CREATE INDEX IF NOT EXISTS idx_groups_status ON prod_v1.groups(status);
CREATE INDEX IF NOT EXISTS idx_groups_type ON prod_v1.groups(type);
CREATE INDEX IF NOT EXISTS idx_groups_created_by ON prod_v1.groups(created_by);
CREATE INDEX IF NOT EXISTS idx_groups_name ON prod_v1.groups(name);

-- Indexes for group_members table
CREATE INDEX IF NOT EXISTS idx_group_members_group_id ON prod_v1.group_members(group_id);
CREATE INDEX IF NOT EXISTS idx_group_members_user_id ON prod_v1.group_members(user_id);
CREATE INDEX IF NOT EXISTS idx_group_members_role ON prod_v1.group_members(role);

-- Indexes for vote_groups table
CREATE INDEX IF NOT EXISTS idx_vote_groups_vote_id ON prod_v1.vote_groups(vote_id);
CREATE INDEX IF NOT EXISTS idx_vote_groups_group_id ON prod_v1.vote_groups(group_id);

-- ============================================
-- TRIGGERS FOR AUTOMATIC TIMESTAMP UPDATES
-- ============================================

-- Trigger for groups table updated_at
CREATE OR REPLACE FUNCTION prod_v1.update_groups_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_groups_updated_at
    BEFORE UPDATE ON prod_v1.groups
    FOR EACH ROW
    EXECUTE FUNCTION prod_v1.update_groups_updated_at();