/* eslint-disable @typescript-eslint/no-unused-vars */
import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { Vote, VoteWithResults, UserVote, VoteFilters } from '@/types/vote';

interface VoteState {
  votes: Vote[];
  currentVote: VoteWithResults | null;
  userVotes: Record<string, UserVote>; // voteId -> UserVote
  hasVotedStatus: Record<string, boolean>; // voteId -> hasVoted
  loading: boolean;
  error: string | null;
  pagination: {
    currentPage: number;
    pageSize: number;
    totalItems: number;
    totalPages: number;
  };
  filters: VoteFilters;
}

const initialState: VoteState = {
  votes: [],
  currentVote: null,
  userVotes: {},
  hasVotedStatus: {},
  loading: false,
  error: null,
  pagination: {
    currentPage: 1,
    pageSize: 20,
    totalItems: 0,
    totalPages: 0,
  },
  filters: {},
};

export const voteSlice = createSlice({
  name: 'vote',
  initialState,
  reducers: {
    // Fetch votes
    fetchVotesRequest: (
      state,
      action: PayloadAction<VoteFilters | undefined>
    ) => {
      state.loading = true;
      state.error = null;
      if (action.payload) {
        state.filters = action.payload;
      }
    },
    fetchVotesSuccess: (
      state,
      action: PayloadAction<{ votes: Vote[]; totalItems: number; page?: number; pageSize?: number; totalPages?: number }>
    ) => {
      state.loading = false;
      state.votes = action.payload.votes;
      state.pagination.totalItems = action.payload.totalItems;
      state.pagination.currentPage = action.payload.page || state.pagination.currentPage;
      state.pagination.pageSize = action.payload.pageSize || state.pagination.pageSize;
      state.pagination.totalPages = action.payload.totalPages || Math.ceil(
        action.payload.totalItems / state.pagination.pageSize
      );
    },
    fetchVotesFailure: (state, action: PayloadAction<string>) => {
      state.loading = false;
      state.error = action.payload;
    },

    // Get single vote
    fetchVoteRequest: (state, _action: PayloadAction<string>) => {
      state.loading = true;
      state.error = null;
    },
    fetchVoteSuccess: (state, action: PayloadAction<Vote>) => {
      state.loading = false;
      const existingIndex = state.votes.findIndex(
        vote => vote.id === action.payload.id
      );
      if (existingIndex !== -1) {
        state.votes[existingIndex] = action.payload;
      } else {
        state.votes.push(action.payload);
      }
    },
    fetchVoteFailure: (state, action: PayloadAction<string>) => {
      state.loading = false;
      state.error = action.payload;
    },

    // Get vote with results
    fetchVoteWithResultsRequest: (state, _action: PayloadAction<string>) => {
      state.loading = true;
      state.error = null;
    },
    fetchVoteWithResultsSuccess: (state, action: PayloadAction<VoteWithResults>) => {
      state.loading = false;
      state.currentVote = action.payload;

      // Update votes array with the vote data
      const existingIndex = state.votes.findIndex(
        vote => vote.id === action.payload.vote.id
      );
      if (existingIndex !== -1) {
        state.votes[existingIndex] = action.payload.vote;
      } else {
        state.votes.push(action.payload.vote);
      }

      // Update user vote status
      if (action.payload.user_vote) {
        state.userVotes[action.payload.vote.id] = action.payload.user_vote;
        state.hasVotedStatus[action.payload.vote.id] = true;
      } else {
        state.hasVotedStatus[action.payload.vote.id] = false;
      }
    },
    fetchVoteWithResultsFailure: (state, action: PayloadAction<string>) => {
      state.loading = false;
      state.error = action.payload;
    },

    // Create vote
    createVoteRequest: (
      state,
      _action: PayloadAction<{ title: string; description?: string; type: 'single' | 'multiple'; options: string[]; team_formation_enabled?: boolean }>
    ) => {
      state.loading = true;
      state.error = null;
    },
    createVoteSuccess: (state, action: PayloadAction<Vote>) => {
      state.loading = false;
      state.votes.unshift(action.payload); // Add to beginning
    },
    createVoteFailure: (state, action: PayloadAction<string>) => {
      state.loading = false;
      state.error = action.payload;
    },

    // Update vote
    updateVoteRequest: (
      state,
      _action: PayloadAction<{ id: string; voteData: { title?: string; description?: string; status?: 'active' | 'closed' | 'cancelled'; team_formation_enabled?: boolean } }>
    ) => {
      state.loading = true;
      state.error = null;
    },
    updateVoteSuccess: (state, action: PayloadAction<Vote>) => {
      state.loading = false;
      const index = state.votes.findIndex(
        vote => vote.id === action.payload.id
      );
      if (index !== -1) {
        state.votes[index] = action.payload;
      }
      // Update current vote if it's the same vote
      if (state.currentVote && state.currentVote.vote.id === action.payload.id) {
        state.currentVote.vote = action.payload;
      }
    },
    updateVoteFailure: (state, action: PayloadAction<string>) => {
      state.loading = false;
      state.error = action.payload;
    },

    // Delete vote
    deleteVoteRequest: (state, _action: PayloadAction<string>) => {
      state.loading = true;
      state.error = null;
    },
    deleteVoteSuccess: (state, action: PayloadAction<string>) => {
      state.loading = false;
      state.votes = state.votes.filter(vote => vote.id !== action.payload);
      // Clear current vote if it's the deleted vote
      if (state.currentVote && state.currentVote.vote.id === action.payload) {
        state.currentVote = null;
      }
      // Clean up user vote data
      delete state.userVotes[action.payload];
      delete state.hasVotedStatus[action.payload];
    },
    deleteVoteFailure: (state, action: PayloadAction<string>) => {
      state.loading = false;
      state.error = action.payload;
    },

    // Cast vote
    castVoteRequest: (
      state,
      _action: PayloadAction<{ voteId: string; optionIds: string[] }>
    ) => {
      state.loading = true;
      state.error = null;
    },
    castVoteSuccess: (state, action: PayloadAction<UserVote>) => {
      state.loading = false;
      state.userVotes[action.payload.vote_id] = action.payload;
      state.hasVotedStatus[action.payload.vote_id] = true;

      // Update current vote results if it's the same vote
      if (state.currentVote && state.currentVote.vote.id === action.payload.vote_id) {
        state.currentVote.user_vote = action.payload;
        // Increment vote counts for selected options
        action.payload.selected_options.forEach(optionId => {
          state.currentVote!.results[optionId] = (state.currentVote!.results[optionId] || 0) + 1;
        });
        state.currentVote.total_votes += 1;
      }
    },
    castVoteFailure: (state, action: PayloadAction<string>) => {
      state.loading = false;
      state.error = action.payload;
    },

    // Get user vote
    fetchUserVoteRequest: (state, _action: PayloadAction<string>) => {
      state.loading = true;
      state.error = null;
    },
    fetchUserVoteSuccess: (state, action: PayloadAction<UserVote>) => {
      state.loading = false;
      state.userVotes[action.payload.vote_id] = action.payload;
      state.hasVotedStatus[action.payload.vote_id] = true;
    },
    fetchUserVoteFailure: (state, action: PayloadAction<string>) => {
      state.loading = false;
      state.error = action.payload;
    },

    // Check if user voted
    checkUserVotedRequest: (state, _action: PayloadAction<string>) => {
      state.loading = true;
      state.error = null;
    },
    checkUserVotedSuccess: (state, action: PayloadAction<{ voteId: string; hasVoted: boolean }>) => {
      state.loading = false;
      state.hasVotedStatus[action.payload.voteId] = action.payload.hasVoted;
    },
    checkUserVotedFailure: (state, action: PayloadAction<string>) => {
      state.loading = false;
      state.error = action.payload;
    },

    // Close vote
    closeVoteRequest: (state, _action: PayloadAction<string>) => {
      state.loading = true;
      state.error = null;
    },
    closeVoteSuccess: (state, action: PayloadAction<Vote>) => {
      state.loading = false;
      const index = state.votes.findIndex(
        vote => vote.id === action.payload.id
      );
      if (index !== -1) {
        state.votes[index] = action.payload;
      }
      // Update current vote if it's the same vote
      if (state.currentVote && state.currentVote.vote.id === action.payload.id) {
        state.currentVote.vote = action.payload;
      }
    },
    closeVoteFailure: (state, action: PayloadAction<string>) => {
      state.loading = false;
      state.error = action.payload;
    },

    // Cancel vote
    cancelVoteRequest: (state, _action: PayloadAction<string>) => {
      state.loading = true;
      state.error = null;
    },
    cancelVoteSuccess: (state, action: PayloadAction<Vote>) => {
      state.loading = false;
      const index = state.votes.findIndex(
        vote => vote.id === action.payload.id
      );
      if (index !== -1) {
        state.votes[index] = action.payload;
      }
      // Update current vote if it's the same vote
      if (state.currentVote && state.currentVote.vote.id === action.payload.id) {
        state.currentVote.vote = action.payload;
      }
    },
    cancelVoteFailure: (state, action: PayloadAction<string>) => {
      state.loading = false;
      state.error = action.payload;
    },

    // Pagination
    setPage: (state, action: PayloadAction<number>) => {
      state.pagination.currentPage = action.payload;
    },
    setPageSize: (state, action: PayloadAction<number>) => {
      state.pagination.pageSize = action.payload;
      state.pagination.currentPage = 1; // Reset to first page when page size changes
    },

    // Filters
    setFilters: (state, action: PayloadAction<VoteFilters>) => {
      state.filters = action.payload;
      state.pagination.currentPage = 1; // Reset to first page when filters change
    },

    // Clear current vote
    clearCurrentVote: (state) => {
      state.currentVote = null;
    },

    // Clear error
    clearError: (state) => {
      state.error = null;
    },

    // Reset state
    resetVoteState: () => initialState,
  },
});

export const {
  fetchVotesRequest,
  fetchVotesSuccess,
  fetchVotesFailure,
  fetchVoteRequest,
  fetchVoteSuccess,
  fetchVoteFailure,
  fetchVoteWithResultsRequest,
  fetchVoteWithResultsSuccess,
  fetchVoteWithResultsFailure,
  createVoteRequest,
  createVoteSuccess,
  createVoteFailure,
  updateVoteRequest,
  updateVoteSuccess,
  updateVoteFailure,
  deleteVoteRequest,
  deleteVoteSuccess,
  deleteVoteFailure,
  castVoteRequest,
  castVoteSuccess,
  castVoteFailure,
  fetchUserVoteRequest,
  fetchUserVoteSuccess,
  fetchUserVoteFailure,
  checkUserVotedRequest,
  checkUserVotedSuccess,
  checkUserVotedFailure,
  closeVoteRequest,
  closeVoteSuccess,
  closeVoteFailure,
  cancelVoteRequest,
  cancelVoteSuccess,
  cancelVoteFailure,
  setPage,
  setPageSize,
  setFilters,
  clearCurrentVote,
  clearError,
  resetVoteState,
} = voteSlice.actions;

export default voteSlice.reducer;
