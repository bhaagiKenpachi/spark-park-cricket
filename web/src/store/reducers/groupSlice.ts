import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit';
import { groupService } from '@/services/groupService';
import {
    GroupWithCreator,
    CreateGroupRequest,
    UpdateGroupRequest,
    GroupFilters,
    PaginatedGroupsResult,
} from '@/types/group';

interface GroupState {
    groups: GroupWithCreator[];
    currentGroup: GroupWithCreator | null;
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
    loading: false,
    error: null,
    pagination: null,
};

// Async thunks for group operations
export const fetchGroupsRequest = createAsyncThunk(
    'group/fetchGroups',
    async (filters?: GroupFilters) => {
        const response = await groupService.getGroups(filters);
        const payload = response.data as any;

        // Normalize various backend response shapes to PaginatedGroupsResult
        // Possible shapes:
        // 1) { data: [ {..group..} ] }
        // 2) { groups: [ {..group..} ], total_items, page, page_size, total_pages }
        // 3) [ {..group..} ]
        let rawGroups: any[] = [];
        if (Array.isArray(payload)) {
            rawGroups = payload;
        } else if (Array.isArray(payload?.data)) {
            rawGroups = payload.data;
        } else if (Array.isArray(payload?.groups)) {
            rawGroups = payload.groups;
        }

        // Map to GroupWithCreator expected by UI
        const groupsMapped: GroupWithCreator[] = rawGroups.map((g: any) => {
            // creator_name -> creator.display_name fallback
            const creatorName: string = g.creator_name || g.creator?.display_name || 'Unknown';
            const creatorId: string = g.created_by || g.creator?.id || '';
            const creatorEmail: string = g.creator?.email || '';

            return {
                id: g.id,
                name: g.name,
                description: g.description,
                type: g.type,
                status: g.status,
                created_by: g.created_by,
                created_at: g.created_at,
                updated_at: g.updated_at,
                creator: {
                    id: creatorId,
                    display_name: creatorName,
                    email: creatorEmail,
                },
            } as GroupWithCreator;
        });

        const totalItems = typeof (payload && payload.total_items) === 'number' ? payload.total_items : groupsMapped.length;
        const pageSize = typeof (payload && payload.page_size) === 'number'
            ? payload.page_size
            : (filters && typeof filters.limit === 'number'
                ? filters.limit
                : (groupsMapped.length || 20));
        const page = typeof (payload && payload.page) === 'number' ? payload.page : 1;
        const totalPages = typeof (payload && payload.total_pages) === 'number'
            ? payload.total_pages
            : (pageSize > 0 ? Math.max(1, Math.ceil(totalItems / pageSize)) : 1);

        const normalized: PaginatedGroupsResult = {
            groups: groupsMapped,
            total_items: totalItems,
            page,
            page_size: pageSize,
            total_pages: totalPages,
        };

        return normalized;
    }
);

export const fetchGroupByIdRequest = createAsyncThunk(
    'group/fetchGroupById',
    async (id: string) => {
        const response = await groupService.getGroupById(id);
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
                state.groups = []; // Reset groups to empty array on error
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

        // Create group
        builder
            .addCase(createGroupRequest.pending, (state) => {
                state.loading = true;
                state.error = null;
            })
            .addCase(createGroupRequest.fulfilled, (state, action: PayloadAction<GroupWithCreator>) => {
                state.loading = false;
                if (state.groups) {
                    state.groups.unshift(action.payload);
                } else {
                    state.groups = [action.payload];
                }
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
                if (state.groups) {
                    const index = state.groups.findIndex(group => group.id === action.payload.id);
                    if (index !== -1) {
                        state.groups[index] = action.payload;
                    }
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
                if (state.groups) {
                    state.groups = state.groups.filter(group => group.id !== action.payload.id);
                }
                if (state.currentGroup?.id === action.payload.id) {
                    state.currentGroup = null;
                }
            })
            .addCase(deleteGroupRequest.rejected, (state, action) => {
                state.loading = false;
                state.error = action.error.message || 'Failed to delete group';
            });
    },
});

export const {
    clearError,
    clearCurrentGroup,
} = groupSlice.actions;

export default groupSlice.reducer;
