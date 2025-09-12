import { createSlice, PayloadAction } from '@reduxjs/toolkit';

export interface Player {
    id: string;
    team_id: string;
    name: string;
    role: 'batsman' | 'bowler' | 'all_rounder' | 'wicket_keeper';
    created_at: string;
    updated_at: string;
}

interface PlayerState {
    players: Player[];
    currentPlayer: Player | null;
    loading: boolean;
    error: string | null;
}

const initialState: PlayerState = {
    players: [],
    currentPlayer: null,
    loading: false,
    error: null,
};

export const playerSlice = createSlice({
    name: 'player',
    initialState,
    reducers: {
        fetchPlayersRequest: (state) => {
            state.loading = true;
            state.error = null;
        },
        fetchPlayersSuccess: (state, action: PayloadAction<Player[]>) => {
            state.loading = false;
            state.players = action.payload;
        },
        fetchPlayersFailure: (state, action: PayloadAction<string>) => {
            state.loading = false;
            state.error = action.payload;
        },
        setCurrentPlayer: (state, action: PayloadAction<Player>) => {
            state.currentPlayer = action.payload;
        },
    },
});

export const {
    fetchPlayersRequest,
    fetchPlayersSuccess,
    fetchPlayersFailure,
    setCurrentPlayer,
} = playerSlice.actions;

export default playerSlice.reducer;
