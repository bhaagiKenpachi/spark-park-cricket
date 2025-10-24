import { createSlice, PayloadAction, createAsyncThunk } from '@reduxjs/toolkit';
import { apiService } from '@/services/api';
import { TimeTrackingResponse } from '@/types/timeTracking';

// Time Tracking State
interface TimeTrackingState {
    timeTracking: TimeTrackingResponse | null;
    loading: boolean;
    error: string | null;
}

const initialState: TimeTrackingState = {
    timeTracking: null,
    loading: false,
    error: null,
};

// Async thunks
export const fetchTimeTrackingRequest = createAsyncThunk(
    'timeTracking/fetchTimeTrackingRequest',
    async (matchId: string) => {
        const response = await apiService.getTimeTracking(matchId);
        if (!response.success) {
            throw new Error(response.message || 'Failed to fetch time tracking data');
        }
        return response.data;
    }
);

export const timeTrackingSlice = createSlice({
    name: 'timeTracking',
    initialState,
    reducers: {
        clearTimeTracking: (state) => {
            state.timeTracking = null;
            state.error = null;
        },
        setTimeTrackingError: (state, action: PayloadAction<string>) => {
            state.error = action.payload;
            state.loading = false;
        },
    },
    extraReducers: (builder) => {
        builder
            .addCase(fetchTimeTrackingRequest.pending, (state) => {
                state.loading = true;
                state.error = null;
            })
            .addCase(fetchTimeTrackingRequest.fulfilled, (state, action) => {
                state.loading = false;
                state.timeTracking = action.payload;
                state.error = null;
            })
            .addCase(fetchTimeTrackingRequest.rejected, (state, action) => {
                state.loading = false;
                state.error = action.error.message || 'Failed to fetch time tracking data';
            });
    },
});

export const { clearTimeTracking, setTimeTrackingError } = timeTrackingSlice.actions;
export default timeTrackingSlice.reducer;


