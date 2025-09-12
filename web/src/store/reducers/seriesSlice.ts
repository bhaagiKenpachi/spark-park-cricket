import { createSlice, PayloadAction } from '@reduxjs/toolkit';

export interface Series {
    id: string;
    name: string;
    description?: string;
    start_date: string;
    end_date: string;
    status: 'upcoming' | 'ongoing' | 'completed';
    created_at: string;
    updated_at: string;
}

interface SeriesState {
    series: Series[];
    currentSeries: Series | null;
    loading: boolean;
    error: string | null;
}

const initialState: SeriesState = {
    series: [],
    currentSeries: null,
    loading: false,
    error: null,
};

export const seriesSlice = createSlice({
    name: 'series',
    initialState,
    reducers: {
        fetchSeriesRequest: (state) => {
            state.loading = true;
            state.error = null;
        },
        fetchSeriesSuccess: (state, action: PayloadAction<Series[]>) => {
            state.loading = false;
            state.series = action.payload;
        },
        fetchSeriesFailure: (state, action: PayloadAction<string>) => {
            state.loading = false;
            state.error = action.payload;
        },
        setCurrentSeries: (state, action: PayloadAction<Series>) => {
            state.currentSeries = action.payload;
        },
        createSeriesRequest: (state, _action: PayloadAction<Omit<Series, 'id' | 'created_at' | 'updated_at'>>) => {
            state.loading = true;
            state.error = null;
        },
        createSeriesSuccess: (state, action: PayloadAction<Series>) => {
            state.loading = false;
            state.series.push(action.payload);
        },
        createSeriesFailure: (state, action: PayloadAction<string>) => {
            state.loading = false;
            state.error = action.payload;
        },
        updateSeriesRequest: (state, _action: PayloadAction<{ id: string; seriesData: Partial<Series> }>) => {
            state.loading = true;
            state.error = null;
        },
        updateSeriesSuccess: (state, action: PayloadAction<Series>) => {
            state.loading = false;
            const index = state.series.findIndex(series => series.id === action.payload.id);
            if (index !== -1) {
                state.series[index] = action.payload;
            }
        },
        updateSeriesFailure: (state, action: PayloadAction<string>) => {
            state.loading = false;
            state.error = action.payload;
        },
        deleteSeriesRequest: (state, _action: PayloadAction<string>) => {
            state.loading = true;
            state.error = null;
        },
        deleteSeriesSuccess: (state, action: PayloadAction<string>) => {
            state.loading = false;
            state.series = state.series.filter(series => series.id !== action.payload);
        },
        deleteSeriesFailure: (state, action: PayloadAction<string>) => {
            state.loading = false;
            state.error = action.payload;
        },
    },
});

export const {
    fetchSeriesRequest,
    fetchSeriesSuccess,
    fetchSeriesFailure,
    setCurrentSeries,
    createSeriesRequest,
    createSeriesSuccess,
    createSeriesFailure,
    updateSeriesRequest,
    updateSeriesSuccess,
    updateSeriesFailure,
    deleteSeriesRequest,
    deleteSeriesSuccess,
    deleteSeriesFailure,
} = seriesSlice.actions;

export default seriesSlice.reducer;
