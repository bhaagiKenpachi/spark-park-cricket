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
    created_at: string;
    updated_at: string;
    closed_at?: string;
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

export interface VoteWithResults {
    vote: Vote;
    options: VoteOption[];
    results: Record<string, number>;
    user_vote?: UserVote;
    total_votes: number;
    voted_users: string[];
}

export interface CreateVoteRequest {
    title: string;
    description?: string;
    type: VoteType;
    options: string[];
}

export interface UpdateVoteRequest {
    title?: string;
    description?: string;
    status?: VoteStatus;
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
}

export interface PaginatedVotesResult {
    votes: Vote[];
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
