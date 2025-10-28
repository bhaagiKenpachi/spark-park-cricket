// Vote types for the frontend voting system

export type VoteStatus = 'active' | 'closed' | 'cancelled';
export type VoteType = 'single' | 'multiple';

export interface Vote {
    id: string;
    title: string;
    description?: string;
    type: VoteType;
    status: VoteStatus;
    created_by: string;
    team_formation_enabled: boolean;
    created_at: string;
    updated_at: string;
    closed_at?: string;
}

export interface VoteWithCreator extends Vote {
    creator_name: string;
}

export interface VoteOption {
    id: string;
    vote_id: string;
    text: string;
    created_at: string;
    updated_at: string;
}

export interface UserVote {
    id: string;
    vote_id: string;
    user_id: string;
    selected_options: string[];
    voted_at: string;
}

export interface VoterInfo {
    user_id: string;
    user_name: string;
    voted_at: string;
}

export interface VoteWithResults {
    vote: Vote;
    options: VoteOption[];
    results: Record<string, number>;
    results_with_names: Record<string, VoterInfo[]>;
    user_vote?: UserVote;
    total_votes: number;
    voted_users: string[];
    creator_name: string;
}

export interface CreateVoteRequest {
    title: string;
    description?: string;
    type: VoteType;
    options: string[];
    team_formation_enabled?: boolean;
}

export interface UpdateVoteRequest {
    title?: string;
    description?: string;
    status?: VoteStatus;
    team_formation_enabled?: boolean;
}

export interface VoteRequest {
    selected_options: string[];
}

export interface VoteFilters {
    status?: VoteStatus;
    type?: VoteType;
    created_by?: string;
    limit?: number;
    offset?: number;
    page?: number;
    page_size?: number;
}

export interface PaginatedVotesResult {
    votes: VoteWithCreator[];
    total_items: number;
    page: number;
    page_size: number;
    total_pages: number;
}

export interface VoteApiResponse<T> {
    data: T;
    message?: string;
    success: boolean;
}

export interface VoteError {
    message: string;
    status: number;
    details?: unknown;
}
