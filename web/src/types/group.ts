// Group types for the frontend group voting system

export type GroupType = 'custom' | 'team' | 'series' | 'match' | 'location' | 'skill';
export type GroupStatus = 'active' | 'inactive' | 'archived';
export type GroupMemberRole = 'admin' | 'moderator' | 'member';

export interface Group {
    id: string;
    name: string;
    description?: string;
    type: GroupType;
    status: GroupStatus;
    created_by: string;
    created_at: string;
    updated_at: string;
}

export interface GroupWithCreator extends Group {
    creator: {
        id: string;
        display_name: string;
        email: string;
    };
}

export interface GroupMember {
    id: string;
    group_id: string;
    user_id: string;
    role: GroupMemberRole;
    joined_at: string;
    user?: {
        id: string;
        display_name: string;
        email: string;
    };
}

export interface GroupWithMembers {
    group: Group;
    members: GroupMember[];
}

export interface GroupStats {
    group_id: string;
    member_count: number;
    vote_count: number;
    admin_count: number;
    moderator_count: number;
}

export interface CreateGroupRequest {
    name: string;
    description?: string;
    type: GroupType;
}

export interface UpdateGroupRequest {
    name?: string;
    description?: string;
    status?: GroupStatus;
}

export interface AddGroupMemberRequest {
    user_id: string;
    role: GroupMemberRole;
}

export interface UpdateGroupMemberRoleRequest {
    role: GroupMemberRole;
}

export interface GroupFilters {
    type?: GroupType;
    status?: GroupStatus;
    search?: string;
    limit?: number;
    offset?: number;
    page?: number;
    page_size?: number;
}

export interface PaginatedGroupsResult {
    groups: GroupWithCreator[];
    total_items: number;
    page: number;
    page_size: number;
    total_pages: number;
}

export interface GroupApiResponse<T> {
    data: T;
    message?: string;
    success: boolean;
}

export interface GroupError {
    message: string;
    status: number;
    details?: unknown;
}

// Group Voting specific types
export interface VoteGroupAssignment {
    vote_id: string;
    group_ids: string[];
}

export interface VoteGroupResults {
    vote: {
        id: string;
        title: string;
        description?: string;
        type: string;
        status: string;
    };
    overall_results: {
        total_votes: number;
        options: Array<{
            id: string;
            text: string;
            vote_count: number;
            percentage?: number;
        }>;
    };
    group_results: Array<{
        group: {
            id: string;
            name: string;
            type: GroupType;
        };
        total_votes: number;
        options: Array<{
            id: string;
            text: string;
            vote_count: number;
            percentage?: number;
            voters?: Array<{
                user_id: string;
                display_name: string;
                voted_at: string;
            }>;
        }>;
    }>;
}

export interface VoteWithGroupResults {
    vote: {
        id: string;
        title: string;
        description?: string;
        type: string;
        status: string;
        created_by: string;
        created_at: string;
        updated_at: string;
        closed_at?: string;
    };
    options: Array<{
        id: string;
        text: string;
        vote_count: number;
    }>;
    overall_results: {
        total_votes: number;
        options: Array<{
            id: string;
            text: string;
            vote_count: number;
            percentage: number;
        }>;
    };
    group_results: Array<{
        group: {
            id: string;
            name: string;
            type: GroupType;
        };
        total_votes: number;
        options: Array<{
            id: string;
            text: string;
            vote_count: number;
            percentage: number;
        }>;
    }>;
}

export interface GroupVotesResult {
    group: {
        id: string;
        name: string;
        type: GroupType;
    };
    votes: Array<{
        id: string;
        title: string;
        description?: string;
        type: string;
        status: string;
        created_by: string;
        created_at: string;
        updated_at: string;
        closed_at?: string;
    }>;
}
