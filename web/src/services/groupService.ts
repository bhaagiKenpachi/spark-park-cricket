import {
    GroupWithCreator,
    GroupWithMembers,
    GroupStats,
    CreateGroupRequest,
    UpdateGroupRequest,
    AddGroupMemberRequest,
    UpdateGroupMemberRoleRequest,
    GroupFilters,
    PaginatedGroupsResult,
    GroupApiResponse,
    VoteGroupResults,
    VoteWithGroupResults,
    GroupVotesResult,
} from '@/types/group';
import { ApiService } from './api';

class GroupService {
    private apiService: ApiService;

    constructor() {
        this.apiService = new ApiService();
    }

    // Group CRUD operations
    async getGroups(filters?: GroupFilters): Promise<GroupApiResponse<PaginatedGroupsResult>> {
        return this.apiService.getGroups(filters);
    }

    async getGroupById(id: string): Promise<GroupApiResponse<GroupWithCreator>> {
        return this.apiService.getGroupById(id);
    }

    async getGroupWithMembers(id: string): Promise<GroupApiResponse<GroupWithMembers>> {
        return this.apiService.getGroupWithMembers(id);
    }

    async createGroup(groupData: CreateGroupRequest): Promise<GroupApiResponse<GroupWithCreator>> {
        return this.apiService.createGroup(groupData);
    }

    async updateGroup(
        id: string,
        groupData: UpdateGroupRequest
    ): Promise<GroupApiResponse<GroupWithCreator>> {
        return this.apiService.updateGroup(id, groupData);
    }

    async deleteGroup(id: string): Promise<GroupApiResponse<{ message: string }>> {
        return this.apiService.deleteGroup(id);
    }

    // Group member operations
    async joinGroup(groupId: string): Promise<GroupApiResponse<{ message: string }>> {
        return this.apiService.joinGroup(groupId);
    }

    async leaveGroup(groupId: string): Promise<GroupApiResponse<{ message: string }>> {
        return this.apiService.leaveGroup(groupId);
    }

    async addGroupMember(
        groupId: string,
        memberData: AddGroupMemberRequest
    ): Promise<GroupApiResponse<{ message: string }>> {
        return this.apiService.addGroupMember(groupId, memberData);
    }

    async removeGroupMember(
        groupId: string,
        userId: string
    ): Promise<GroupApiResponse<{ message: string }>> {
        return this.apiService.removeGroupMember(groupId, userId);
    }

    async updateGroupMemberRole(
        groupId: string,
        userId: string,
        roleData: UpdateGroupMemberRoleRequest
    ): Promise<GroupApiResponse<{ message: string }>> {
        return this.apiService.updateGroupMemberRole(groupId, userId, roleData);
    }

    // Group search and filtering
    async searchGroups(query: string, limit?: number, offset?: number): Promise<GroupApiResponse<PaginatedGroupsResult>> {
        const filters: GroupFilters = { search: query };
        if (limit !== undefined) filters.limit = limit;
        if (offset !== undefined) filters.offset = offset;
        return this.getGroups(filters);
    }

    async getGroupsByType(type: string, limit?: number, offset?: number): Promise<GroupApiResponse<PaginatedGroupsResult>> {
        const filters: GroupFilters = { type: type as any };
        if (limit !== undefined) filters.limit = limit;
        if (offset !== undefined) filters.offset = offset;
        return this.getGroups(filters);
    }

    // Group statistics
    async getGroupStats(groupId: string): Promise<GroupApiResponse<GroupStats>> {
        return this.apiService.getGroupStats(groupId);
    }

    // Group voting operations
    async assignGroupsToVote(voteId: string, groupIds: string[]): Promise<GroupApiResponse<{ message: string }>> {
        return this.apiService.assignGroupsToVote(voteId, groupIds);
    }

    async removeGroupsFromVote(voteId: string, groupIds: string[]): Promise<GroupApiResponse<{ message: string }>> {
        return this.apiService.removeGroupsFromVote(voteId, groupIds);
    }

    async getVoteGroups(voteId: string): Promise<GroupApiResponse<{ groups: GroupWithCreator[] }>> {
        return this.apiService.getVoteGroups(voteId);
    }

    async getVoteResultsByGroups(voteId: string): Promise<GroupApiResponse<VoteGroupResults>> {
        return this.apiService.getVoteResultsByGroups(voteId);
    }

    async getGroupVoteResults(voteId: string, groupId: string): Promise<GroupApiResponse<VoteGroupResults>> {
        return this.apiService.getGroupVoteResults(voteId, groupId);
    }

    async getVoteWithGroupResults(voteId: string): Promise<GroupApiResponse<VoteWithGroupResults>> {
        return this.apiService.getVoteWithGroupResults(voteId);
    }

    async getGroupVotes(groupId: string): Promise<GroupApiResponse<GroupVotesResult>> {
        return this.apiService.getGroupVotes(groupId);
    }

    // Helper methods
    async getActiveGroups(limit?: number, offset?: number): Promise<GroupApiResponse<PaginatedGroupsResult>> {
        const filters: GroupFilters = { status: 'active' };
        if (limit !== undefined) filters.limit = limit;
        if (offset !== undefined) filters.offset = offset;
        return this.getGroups(filters);
    }

    async getCustomGroups(limit?: number, offset?: number): Promise<GroupApiResponse<PaginatedGroupsResult>> {
        return this.getGroupsByType('custom', limit, offset);
    }

    async getTeamGroups(limit?: number, offset?: number): Promise<GroupApiResponse<PaginatedGroupsResult>> {
        return this.getGroupsByType('team', limit, offset);
    }

    async getSeriesGroups(limit?: number, offset?: number): Promise<GroupApiResponse<PaginatedGroupsResult>> {
        return this.getGroupsByType('series', limit, offset);
    }

    async getMatchGroups(limit?: number, offset?: number): Promise<GroupApiResponse<PaginatedGroupsResult>> {
        return this.getGroupsByType('match', limit, offset);
    }

    async getLocationGroups(limit?: number, offset?: number): Promise<GroupApiResponse<PaginatedGroupsResult>> {
        return this.getGroupsByType('location', limit, offset);
    }

    async getSkillGroups(limit?: number, offset?: number): Promise<GroupApiResponse<PaginatedGroupsResult>> {
        return this.getGroupsByType('skill', limit, offset);
    }
}

export const groupService = new GroupService();
export { GroupService };
export default groupService;
