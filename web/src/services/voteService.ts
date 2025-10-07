import {
    Vote,
    VoteWithResults,
    CreateVoteRequest,
    UpdateVoteRequest,
    VoteRequest,
    VoteFilters,
    PaginatedVotesResult,
    UserVote,
    VoteApiResponse,
} from '@/types/vote';
import { ApiService, ApiError } from './api';

class VoteService {
    private apiService: ApiService;

    constructor() {
        this.apiService = new ApiService();
    }

    // Vote CRUD operations
    async getVotes(filters?: VoteFilters): Promise<VoteApiResponse<Vote[]>> {
        return this.apiService.getVotes(filters);
    }

    async getVoteById(id: string): Promise<VoteApiResponse<Vote>> {
        return this.apiService.getVoteById(id);
    }

    async getVoteWithResults(id: string): Promise<VoteApiResponse<VoteWithResults>> {
        return this.apiService.getVoteWithResults(id);
    }

    async createVote(voteData: CreateVoteRequest): Promise<VoteApiResponse<Vote>> {
        return this.apiService.createVote(voteData);
    }

    async updateVote(
        id: string,
        voteData: UpdateVoteRequest
    ): Promise<VoteApiResponse<Vote>> {
        return this.apiService.updateVote(id, voteData);
    }

    async deleteVote(id: string): Promise<VoteApiResponse<{ message: string }>> {
        return this.apiService.deleteVote(id);
    }

    // Voting operations
    async castVote(voteId: string, voteData: VoteRequest): Promise<VoteApiResponse<UserVote>> {
        return this.apiService.castVote(voteId, voteData);
    }

    async getUserVote(voteId: string): Promise<VoteApiResponse<UserVote>> {
        return this.apiService.getUserVote(voteId);
    }

    async hasUserVoted(voteId: string): Promise<VoteApiResponse<{ has_voted: boolean }>> {
        return this.apiService.hasUserVoted(voteId);
    }

    // Vote management operations
    async closeVote(voteId: string): Promise<VoteApiResponse<Vote>> {
        return this.apiService.closeVote(voteId);
    }

    async cancelVote(voteId: string): Promise<VoteApiResponse<Vote>> {
        return this.apiService.cancelVote(voteId);
    }

    // Helper methods
    async getActiveVotes(limit?: number, offset?: number): Promise<VoteApiResponse<Vote[]>> {
        const filters: VoteFilters = { status: 'active' };
        if (limit !== undefined) filters.limit = limit;
        if (offset !== undefined) filters.offset = offset;
        return this.getVotes(filters);
    }

    async getSingleChoiceVotes(limit?: number, offset?: number): Promise<VoteApiResponse<Vote[]>> {
        const filters: VoteFilters = { type: 'single' };
        if (limit !== undefined) filters.limit = limit;
        if (offset !== undefined) filters.offset = offset;
        return this.getVotes(filters);
    }

    async getMultipleChoiceVotes(limit?: number, offset?: number): Promise<VoteApiResponse<Vote[]>> {
        const filters: VoteFilters = { type: 'multiple' };
        if (limit !== undefined) filters.limit = limit;
        if (offset !== undefined) filters.offset = offset;
        return this.getVotes(filters);
    }

    async getUserVotes(userId: string, limit?: number, offset?: number): Promise<VoteApiResponse<Vote[]>> {
        const filters: VoteFilters = { created_by: userId };
        if (limit !== undefined) filters.limit = limit;
        if (offset !== undefined) filters.offset = offset;
        return this.getVotes(filters);
    }
}

export const voteService = new VoteService();
export { VoteService };
export default voteService;
