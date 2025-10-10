/* eslint-disable @typescript-eslint/no-unused-vars */
import { createSlice, PayloadAction } from '@reduxjs/toolkit';

export interface Match {
  id: string;
  series_id: string;
  match_number: number;
  date: string;
  status: 'not_started' | 'live' | 'completed' | 'cancelled';
  team_a_player_count: number;
  team_b_player_count: number;
  total_overs: number;
  toss_winner: 'A' | 'B';
  toss_type: 'H' | 'T';
  batting_team: 'A' | 'B';
  start_time?: string;
  end_time?: string;
  created_at: string;
  updated_at: string;
}

interface MatchState {
  matches: Match[];
  currentMatch: Match | null;
  loading: boolean;
  error: string | null;
  // Cache for completed matches scorecard data
  completedMatchesCache: {
    [matchId: string]: {
      data: unknown;
      timestamp: number;
      expiresAt: number;
    };
  };
}

const initialState: MatchState = {
  matches: [],
  currentMatch: null,
  loading: false,
  error: null,
  completedMatchesCache: {},
};

export const matchSlice = createSlice({
  name: 'match',
  initialState,
  reducers: {
    fetchMatchesRequest: state => {
      state.loading = true;
      state.error = null;
    },
    fetchMatchesBySeriesRequest: (state, _action: PayloadAction<string>) => {
      state.loading = true;
      state.error = null;
    },
    fetchMatchesSuccess: (state, action: PayloadAction<Match[]>) => {
      state.loading = false;
      state.matches = action.payload;
    },
    fetchMatchesFailure: (state, action: PayloadAction<string>) => {
      state.loading = false;
      state.error = action.payload;
    },
    setCurrentMatch: (state, action: PayloadAction<Match>) => {
      state.currentMatch = action.payload;
    },
    createMatchRequest: (
      state,
      _action: PayloadAction<
        Omit<Match, 'id' | 'created_at' | 'updated_at' | 'match_number'>
      >
    ) => {
      state.loading = true;
      state.error = null;
    },
    createMatchSuccess: (state, action: PayloadAction<Match>) => {
      state.loading = false;
      state.matches.unshift(action.payload); // Add to beginning so new matches appear at top
    },
    createMatchFailure: (state, action: PayloadAction<string>) => {
      state.loading = false;
      state.error = action.payload;
    },
    updateMatchRequest: (
      state,
      _action: PayloadAction<{
        id: string;
        matchData: Omit<
          Match,
          'id' | 'created_at' | 'updated_at' | 'match_number'
        >;
      }>
    ) => {
      state.loading = true;
      state.error = null;
    },
    updateMatchSuccess: (state, action: PayloadAction<Match>) => {
      state.loading = false;
      const index = state.matches.findIndex(
        match => match.id === action.payload.id
      );
      if (index !== -1) {
        state.matches[index] = action.payload;
      }
    },
    updateMatchFailure: (state, action: PayloadAction<string>) => {
      state.loading = false;
      state.error = action.payload;
    },
    deleteMatchRequest: (state, _action: PayloadAction<string>) => {
      state.loading = true;
      state.error = null;
    },
    deleteMatchSuccess: (state, action: PayloadAction<string>) => {
      state.loading = false;
      state.matches = state.matches.filter(
        match => match.id !== action.payload
      );
    },
    deleteMatchFailure: (state, action: PayloadAction<string>) => {
      state.loading = false;
      state.error = action.payload;
    },
    startMatchRequest: (state, _action: PayloadAction<string>) => {
      state.loading = true;
      state.error = null;
    },
    startMatchSuccess: (state, action: PayloadAction<Match>) => {
      state.loading = false;
      const index = state.matches.findIndex(
        match => match.id === action.payload.id
      );
      if (index !== -1 && state.matches[index]) {
        // Update only the status and updated_at fields
        state.matches[index].status = action.payload.status;
        state.matches[index].updated_at = action.payload.updated_at;
      }
    },
    startMatchFailure: (state, action: PayloadAction<string>) => {
      state.loading = false;
      state.error = action.payload;
    },
    // Cache management for completed matches
    cacheCompletedMatchData: (
      state,
      action: PayloadAction<{
        matchId: string;
        data: unknown;
        cacheDuration?: number; // in milliseconds, default 5 minutes
      }>
    ) => {
      const { matchId, data, cacheDuration = 5 * 60 * 1000 } = action.payload;
      const now = Date.now();
      state.completedMatchesCache[matchId] = {
        data,
        timestamp: now,
        expiresAt: now + cacheDuration,
      };
    },
    clearExpiredCache: (state) => {
      const now = Date.now();
      Object.keys(state.completedMatchesCache).forEach(matchId => {
        const cacheEntry = state.completedMatchesCache[matchId];
        if (cacheEntry && cacheEntry.expiresAt < now) {
          delete state.completedMatchesCache[matchId];
        }
      });
    },
    clearMatchCache: (state, action: PayloadAction<string>) => {
      delete state.completedMatchesCache[action.payload];
    },
    clearAllCache: (state) => {
      state.completedMatchesCache = {};
    },
  },
});

export const {
  fetchMatchesRequest,
  fetchMatchesBySeriesRequest,
  fetchMatchesSuccess,
  fetchMatchesFailure,
  setCurrentMatch,
  createMatchRequest,
  createMatchSuccess,
  createMatchFailure,
  updateMatchRequest,
  updateMatchSuccess,
  updateMatchFailure,
  deleteMatchRequest,
  deleteMatchSuccess,
  deleteMatchFailure,
  startMatchRequest,
  startMatchSuccess,
  startMatchFailure,
  cacheCompletedMatchData,
  clearExpiredCache,
  clearMatchCache,
  clearAllCache,
} = matchSlice.actions;

export default matchSlice.reducer;
