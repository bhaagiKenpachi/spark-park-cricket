import { createSlice, PayloadAction } from '@reduxjs/toolkit';

export interface Scoreboard {
    id: string;
    match_id: string;
    team1_runs: number;
    team1_wickets: number;
    team1_overs: number;
    team2_runs: number;
    team2_wickets: number;
    team2_overs: number;
    current_innings: number;
    current_batting_team: string;
    current_bowling_team: string;
    created_at: string;
    updated_at: string;
}

interface ScoreboardState {
    scoreboard: Scoreboard | null;
    loading: boolean;
    error: string | null;
}

const initialState: ScoreboardState = {
    scoreboard: null,
    loading: false,
    error: null,
};

export const scoreboardSlice = createSlice({
    name: 'scoreboard',
    initialState,
    reducers: {
        fetchScoreboardRequest: (state) => {
            state.loading = true;
            state.error = null;
        },
        fetchScoreboardSuccess: (state, action: PayloadAction<Scoreboard>) => {
            state.loading = false;
            state.scoreboard = action.payload;
        },
        fetchScoreboardFailure: (state, action: PayloadAction<string>) => {
            state.loading = false;
            state.error = action.payload;
        },
        updateScoreboard: (state, action: PayloadAction<Partial<Scoreboard>>) => {
            if (state.scoreboard) {
                state.scoreboard = { ...state.scoreboard, ...action.payload };
            }
        },
    },
});

export const {
    fetchScoreboardRequest,
    fetchScoreboardSuccess,
    fetchScoreboardFailure,
    updateScoreboard,
} = scoreboardSlice.actions;

export default scoreboardSlice.reducer;
