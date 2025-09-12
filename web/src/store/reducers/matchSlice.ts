import { createSlice, PayloadAction } from '@reduxjs/toolkit';

export interface Match {
    id: string;
    series_id: string;
    match_number: number;
    date: string;
    status: 'live' | 'completed' | 'cancelled';
    team_a_player_count: number;
    team_b_player_count: number;
    total_overs: number;
    toss_winner: 'A' | 'B';
    toss_type: 'H' | 'T';
    batting_team: 'A' | 'B';
    created_at: string;
    updated_at: string;
}

interface MatchState {
    matches: Match[];
    currentMatch: Match | null;
    loading: boolean;
    error: string | null;
}

const initialState: MatchState = {
    matches: [],
    currentMatch: null,
    loading: false,
    error: null,
};

export const matchSlice = createSlice({
    name: 'match',
    initialState,
    reducers: {
        fetchMatchesRequest: (state) => {
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
        createMatchRequest: (state, _action: PayloadAction<Omit<Match, 'id' | 'created_at' | 'updated_at'>>) => {
            state.loading = true;
            state.error = null;
        },
        createMatchSuccess: (state, action: PayloadAction<Match>) => {
            state.loading = false;
            state.matches.push(action.payload);
        },
        createMatchFailure: (state, action: PayloadAction<string>) => {
            state.loading = false;
            state.error = action.payload;
        },
        updateMatchRequest: (state, _action: PayloadAction<{ id: string; matchData: Partial<Match> }>) => {
            state.loading = true;
            state.error = null;
        },
        updateMatchSuccess: (state, action: PayloadAction<Match>) => {
            state.loading = false;
            const index = state.matches.findIndex(match => match.id === action.payload.id);
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
            state.matches = state.matches.filter(match => match.id !== action.payload);
        },
        deleteMatchFailure: (state, action: PayloadAction<string>) => {
            state.loading = false;
            state.error = action.payload;
        },
    },
});

export const {
    fetchMatchesRequest,
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
} = matchSlice.actions;

export default matchSlice.reducer;