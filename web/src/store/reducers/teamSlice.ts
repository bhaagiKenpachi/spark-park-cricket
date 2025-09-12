import { createSlice, PayloadAction } from '@reduxjs/toolkit';

export interface Team {
    id: string;
    name: string;
    description?: string;
    created_at: string;
    updated_at: string;
}

interface TeamState {
    teams: Team[];
    currentTeam: Team | null;
    loading: boolean;
    error: string | null;
}

const initialState: TeamState = {
    teams: [],
    currentTeam: null,
    loading: false,
    error: null,
};

export const teamSlice = createSlice({
    name: 'team',
    initialState,
    reducers: {
        fetchTeamsRequest: (state) => {
            state.loading = true;
            state.error = null;
        },
        fetchTeamsSuccess: (state, action: PayloadAction<Team[]>) => {
            state.loading = false;
            state.teams = action.payload;
        },
        fetchTeamsFailure: (state, action: PayloadAction<string>) => {
            state.loading = false;
            state.error = action.payload;
        },
        setCurrentTeam: (state, action: PayloadAction<Team>) => {
            state.currentTeam = action.payload;
        },
    },
});

export const {
    fetchTeamsRequest,
    fetchTeamsSuccess,
    fetchTeamsFailure,
    setCurrentTeam,
} = teamSlice.actions;

export default teamSlice.reducer;
