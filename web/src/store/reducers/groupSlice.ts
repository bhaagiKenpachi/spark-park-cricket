import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit';
import { groupService } from '@/services/groupService';
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
    VoteGroupResults,
    VoteWithGroupResults,
    GroupVotesResult,
} from '@/types/group';

interface GroupState {
    groups: GroupWithCreator[];
    currentGroup: GroupWithCreator | null;
    groupMembers: GroupWithMembers | null;
    groupStats: GroupStats | null;
    groupVotes: GroupVotesResult | null;
    voteGroupResults: VoteGroupResults | null;
    voteWithGroupResults: VoteWithGroupResults | null;
    loading: boolean;
    error: string | null;
    pagination: {
        total_items: number;
        page: number;
        page_size: number;
        total_pages: number;
    } | null;
}

const initialState: GroupState = {
    groups: [],
    currentGroup: null,
    groupMembers: null,
    groupStats: null,
    groupVotes: null,
    voteGroupResults: null,
    voteWithGroupResults: null,
    loading: false,
    error: null,
    pagination: null,
};

// Async thunks for group operations
export const fetchGroupsRequest = createAsyncThunk(
    'group/fetchGroups',
    async (filters?: GroupFilters) => {
        const response = await groupService.getGroups(filters);
        return response.data;
    }
);

export const fetchGroupByIdRequest = createAsyncThunk(
    'group/fetchGroupById',
    async (id: string) => {
        const response = await groupService.getGroupById(id);
        return response.data;
    }
);

export const fetchGroupWithMembersRequest = createAsyncThunk(
    'group/fetchGroupWithMembers',
    async (id: string) => {
        const response = await groupService.getGroupWithMembers(id);
        return response.data;
    }
);

export const createGroupRequest = createAsyncThunk(
    'group/createGroup',
    async (groupData: CreateGroupRequest) => {
        const response = await groupService.createGroup(groupData);
        return response.data;
    }
);

export const updateGroupRequest = createAsyncThunk(
    'group/updateGroup',
    async ({ id, groupData }: { id: string; groupData: UpdateGroupRequest }) => {
        const response = await groupService.updateGroup(id, groupData);
        return response.data;
    }
);

export const deleteGroupRequest = createAsyncThunk(
    'group/deleteGroup',
    async (id: string) => {
        const response = await groupService.deleteGroup(id);
        return { id, message: response.data.message };
    }
);

export const joinGroupRequest = createAsyncThunk(
    'group/joinGroup',
    async (groupId: string) => {
        const response = await groupService.joinGroup(groupId);
        return { groupId, message: response.data.message };
    }
);

export const leaveGroupRequest = createAsyncThunk(
    'group/leaveGroup',
    async (groupId: string) => {
        const response = await groupService.leaveGroup(groupId);
        return { groupId, message: response.data.message };
    }
);

export const addGroupMemberRequest = createAsyncThunk(
    'group/addGroupMember',
    async ({ groupId, memberData }: { groupId: string; memberData: AddGroupMemberRequest }) => {
        const response = await groupService.addGroupMember(groupId, memberData);
        return { groupId, memberData, message: response.data.message };
    }
);

export const removeGroupMemberRequest = createAsyncThunk(
    'group/removeGroupMember',
    async ({ groupId, userId }: { groupId: string; userId: string }) => {
        const response = await groupService.removeGroupMember(groupId, userId);
        return { groupId, userId, message: response.data.message };
    }
);

export const updateGroupMemberRoleRequest = createAsyncThunk(
    'group/updateGroupMemberRole',
    async ({ groupId, userId, roleData }: { groupId: string; userId: string; roleData: UpdateGroupMemberRoleRequest }) => {
        const response = await groupService.updateGroupMemberRole(groupId, userId, roleData);
        return { groupId, userId, roleData, message: response.data.message };
    }
);

export const fetchGroupStatsRequest = createAsyncThunk(
    'group/fetchGroupStats',
    async (groupId: string) => {
        const response = await groupService.getGroupStats(groupId);
        return response.data;
    }
);

// Group voting async thunks
export const assignGroupsToVoteRequest = createAsyncThunk(
    'group/assignGroupsToVote',
    async ({ voteId, groupIds }: { voteId: string; groupIds: string[] }) => {
        const response = await groupService.assignGroupsToVote(voteId, groupIds);
        return { voteId, groupIds, message: response.data.message };
    }
);

export const removeGroupsFromVoteRequest = createAsyncThunk(
    'group/removeGroupsFromVote',
    async ({ voteId, groupIds }: { voteId: string; groupIds: string[] }) => {
        const response = await groupService.removeGroupsFromVote(voteId, groupIds);
        return { voteId, groupIds, message: response.data.message };
    }
);

export const fetchVoteGroupsRequest = createAsyncThunk(
    'group/fetchVoteGroups',
    async (voteId: string) => {
        const response = await groupService.getVoteGroups(voteId);
        return response.data;
    }
);

export const fetchVoteResultsByGroupsRequest = createAsyncThunk(
    'group/fetchVoteResultsByGroups',
    async (voteId: string) => {
        const response = await groupService.getVoteResultsByGroups(voteId);
        return response.data;
    }
);

export const fetchGroupVoteResultsRequest = createAsyncThunk(
    'group/fetchGroupVoteResults',
    async ({ voteId, groupId }: { voteId: string; groupId: string }) => {
        const response = await groupService.getGroupVoteResults(voteId, groupId);
        return response.data;
    }
);

export const fetchVoteWithGroupResultsRequest = createAsyncThunk(
    'group/fetchVoteWithGroupResults',
    async (voteId: string) => {
        const response = await groupService.getVoteWithGroupResults(voteId);
        return response.data;
    }
);

export const fetchGroupVotesRequest = createAsyncThunk(
    'group/fetchGroupVotes',
    async (groupId: string) => {
        const response = await groupService.getGroupVotes(groupId);
        return response.data;
    }
);

const groupSlice = createSlice({
    name: 'group',
    initialState,
    reducers: {
        clearError: (state) => {
            state.error = null;
        },
        clearCurrentGroup: (state) => {
            state.currentGroup = null;
        },
        clearGroupMembers: (state) => {
            state.groupMembers = null;
        },
        clearGroupStats: (state) => {
            state.groupStats = null;
        },
        clearGroupVotes: (state) => {
            state.groupVotes = null;
        },
        clearVoteGroupResults: (state) => {
            state.voteGroupResults = null;
        },
        clearVoteWithGroupResults: (state) => {
            state.voteWithGroupResults = null;
        },
    },
    extraReducers: (builder) => {
        // Fetch groups
        builder
            .addCase(fetchGroupsRequest.pending, (state) => {
                state.loading = true;
                state.error = null;
            })
            .addCase(fetchGroupsRequest.fulfilled, (state, action: PayloadAction<PaginatedGroupsResult>) => {
                state.loading = false;
                state.groups = action.payload.groups;
                state.pagination = {
                    total_items: action.payload.total_items,
                    page: action.payload.page,
                    page_size: action.payload.page_size,
                    total_pages: action.payload.total_pages,
                };
            })
            .addCase(fetchGroupsRequest.rejected, (state, action) => {
                state.loading = false;
                state.error = action.error.message || 'Failed to fetch groups';
            });

        // Fetch group by ID
        builder
            .addCase(fetchGroupByIdRequest.pending, (state) => {
                state.loading = true;
                state.error = null;
            })
            .addCase(fetchGroupByIdRequest.fulfilled, (state, action: PayloadAction<GroupWithCreator>) => {
                state.loading = false;
                state.currentGroup = action.payload;
            })
            .addCase(fetchGroupByIdRequest.rejected, (state, action) => {
                state.loading = false;
                state.error = action.error.message || 'Failed to fetch group';
            });

        // Fetch group with members
        builder
            .addCase(fetchGroupWithMembersRequest.pending, (state) => {
                state.loading = true;
                state.error = null;
            })
            .addCase(fetchGroupWithMembersRequest.fulfilled, (state, action: PayloadAction<GroupWithMembers>) => {
                state.loading = false;
                state.groupMembers = action.payload;
            })
            .addCase(fetchGroupWithMembersRequest.rejected, (state, action) => {
                state.loading = false;
                state.error = action.error.message || 'Failed to fetch group members';
            });

        // Create group
        builder
            .addCase(createGroupRequest.pending, (state) => {
                state.loading = true;
                state.error = null;
            })
            .addCase(createGroupRequest.fulfilled, (state, action: PayloadAction<GroupWithCreator>) => {
                state.loading = false;
                state.groups.unshift(action.payload);
                state.currentGroup = action.payload;
            })
            .addCase(createGroupRequest.rejected, (state, action) => {
                state.loading = false;
                state.error = action.error.message || 'Failed to create group';
            });

        // Update group
        builder
            .addCase(updateGroupRequest.pending, (state) => {
                state.loading = true;
                state.error = null;
            })
            .addCase(updateGroupRequest.fulfilled, (state, action: PayloadAction<GroupWithCreator>) => {
                state.loading = false;
                const index = state.groups.findIndex(group => group.id === action.payload.id);
                if (index !== -1) {
                    state.groups[index] = action.payload;
                }
                if (state.currentGroup?.id === action.payload.id) {
                    state.currentGroup = action.payload;
                }
            })
            .addCase(updateGroupRequest.rejected, (state, action) => {
                state.loading = false;
                state.error = action.error.message || 'Failed to update group';
            });

        // Delete group
        builder
            .addCase(deleteGroupRequest.pending, (state) => {
                state.loading = true;
                state.error = null;
            })
            .addCase(deleteGroupRequest.fulfilled, (state, action) => {
                state.loading = false;
                state.groups = state.groups.filter(group => group.id !== action.payload.id);
                if (state.currentGroup?.id === action.payload.id) {
                    state.currentGroup = null;
                }
            })
            .addCase(deleteGroupRequest.rejected, (state, action) => {
                state.loading = false;
                state.error = action.error.message || 'Failed to delete group';
            });

        // Join group
        builder
            .addCase(joinGroupRequest.pending, (state) => {
                state.loading = true;
                state.error = null;
            })
            .addCase(joinGroupRequest.fulfilled, (state) => {
                state.loading = false;
                // Refresh group members if we have them loaded
                if (state.groupMembers) {
                    // The actual member addition will be handled by refetching
                }
            })
            .addCase(joinGroupRequest.rejected, (state, action) => {
                state.loading = false;
                state.error = action.error.message || 'Failed to join group';
            });

        // Leave group
        builder
            .addCase(leaveGroupRequest.pending, (state) => {
                state.loading = true;
                state.error = null;
            })
            .addCase(leaveGroupRequest.fulfilled, (state) => {
                state.loading = false;
                // Refresh group members if we have them loaded
                if (state.groupMembers) {
                    // The actual member removal will be handled by refetching
                }
            })
            .addCase(leaveGroupRequest.rejected, (state, action) => {
                state.loading = false;
                state.error = action.error.message || 'Failed to leave group';
            });

        // Fetch group stats
        builder
            .addCase(fetchGroupStatsRequest.pending, (state) => {
                state.loading = true;
                state.error = null;
            })
            .addCase(fetchGroupStatsRequest.fulfilled, (state, action: PayloadAction<GroupStats>) => {
                state.loading = false;
                state.groupStats = action.payload;
            })
            .addCase(fetchGroupStatsRequest.rejected, (state, action) => {
                state.loading = false;
                state.error = action.error.message || 'Failed to fetch group stats';
            });

        // Fetch vote groups
        builder
            .addCase(fetchVoteGroupsRequest.pending, (state) => {
                state.loading = true;
                state.error = null;
            })
            .addCase(fetchVoteGroupsRequest.fulfilled, (state, action) => {
                state.loading = false;
                // This will be handled by the vote slice or component
            })
            .addCase(fetchVoteGroupsRequest.rejected, (state, action) => {
                state.loading = false;
                state.error = action.error.message || 'Failed to fetch vote groups';
            });

        // Fetch vote results by groups
        builder
            .addCase(fetchVoteResultsByGroupsRequest.pending, (state) => {
                state.loading = true;
                state.error = null;
            })
            .addCase(fetchVoteResultsByGroupsRequest.fulfilled, (state, action: PayloadAction<VoteGroupResults>) => {
                state.loading = false;
                state.voteGroupResults = action.payload;
            })
            .addCase(fetchVoteResultsByGroupsRequest.rejected, (state, action) => {
                state.loading = false;
                state.error = action.error.message || 'Failed to fetch vote results by groups';
            });

        // Fetch vote with group results
        builder
            .addCase(fetchVoteWithGroupResultsRequest.pending, (state) => {
                state.loading = true;
                state.error = null;
            })
            .addCase(fetchVoteWithGroupResultsRequest.fulfilled, (state, action: PayloadAction<VoteWithGroupResults>) => {
                state.loading = false;
                state.voteWithGroupResults = action.payload;
            })
            .addCase(fetchVoteWithGroupResultsRequest.rejected, (state, action) => {
                state.loading = false;
                state.error = action.error.message || 'Failed to fetch vote with group results';
            });

        // Fetch group votes
        builder
            .addCase(fetchGroupVotesRequest.pending, (state) => {
                state.loading = true;
                state.error = null;
            })
            .addCase(fetchGroupVotesRequest.fulfilled, (state, action: PayloadAction<GroupVotesResult>) => {
                state.loading = false;
                state.groupVotes = action.payload;
            })
            .addCase(fetchGroupVotesRequest.rejected, (state, action) => {
                state.loading = false;
                state.error = action.error.message || 'Failed to fetch group votes';
            });
    },
});

export const {
    clearError,
    clearCurrentGroup,
    clearGroupMembers,
    clearGroupStats,
    clearGroupVotes,
    clearVoteGroupResults,
    clearVoteWithGroupResults,
} = groupSlice.actions;

export default groupSlice.reducer;
