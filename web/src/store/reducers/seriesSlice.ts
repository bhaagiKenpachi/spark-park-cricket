/* eslint-disable @typescript-eslint/no-unused-vars */
import { createSlice, PayloadAction } from '@reduxjs/toolkit';

export interface Series {
  id: string;
  name: string;
  description?: string;
  team_a_name?: string;
  team_b_name?: string;
  start_date: string;
  end_date: string;
  status: 'upcoming' | 'ongoing' | 'completed';
  created_by?: string;
  created_at: string;
  updated_at: string;
}

interface SeriesState {
  series: Series[];
  currentSeries: Series | null;
  loading: boolean;
  error: string | null;
  pagination: {
    currentPage: number;
    pageSize: number;
    totalItems: number;
    totalPages: number;
  };
}

const initialState: SeriesState = {
  series: [],
  currentSeries: null,
  loading: false,
  error: null,
  pagination: {
    currentPage: 1,
    pageSize: 20,
    totalItems: 0,
    totalPages: 0,
  },
};

export const seriesSlice = createSlice({
  name: 'series',
  initialState,
  reducers: {
    fetchSeriesRequest: (
      state,
      action: PayloadAction<{ page?: number; pageSize?: number } | undefined>
    ) => {
      state.loading = true;
      state.error = null;

      // Update pagination parameters if provided
      if (action.payload) {
        if (action.payload.page !== undefined) {
          state.pagination.currentPage = action.payload.page;
        }
        if (action.payload.pageSize !== undefined) {
          state.pagination.pageSize = action.payload.pageSize;
        }
      }
    },
    fetchSeriesSuccess: (
      state,
      action: PayloadAction<{ series: Series[]; totalItems: number }>
    ) => {
      state.loading = false;
      state.series = action.payload.series;
      state.pagination.totalItems = action.payload.totalItems;
      state.pagination.totalPages = Math.ceil(
        action.payload.totalItems / state.pagination.pageSize
      );
    },
    fetchSeriesFailure: (state, action: PayloadAction<string>) => {
      state.loading = false;
      state.error = action.payload;
    },
    setCurrentSeries: (state, action: PayloadAction<Series>) => {
      state.currentSeries = action.payload;
    },
    createSeriesRequest: (
      state,
      _action: PayloadAction<Omit<Series, 'id' | 'created_at' | 'updated_at'>>
    ) => {
      state.loading = true;
      state.error = null;
    },
    createSeriesSuccess: (state, action: PayloadAction<Series>) => {
      state.loading = false;
      state.series.unshift(action.payload); // Add to beginning so new series appear at top
      state.pagination.totalItems += 1; // Increment total count
      state.pagination.totalPages = Math.ceil(
        state.pagination.totalItems / state.pagination.pageSize
      );
    },
    createSeriesFailure: (state, action: PayloadAction<string>) => {
      state.loading = false;
      state.error = action.payload;
    },
    updateSeriesRequest: (
      state,
      _action: PayloadAction<{ id: string; seriesData: Partial<Series> }>
    ) => {
      state.loading = true;
      state.error = null;
    },
    updateSeriesSuccess: (state, action: PayloadAction<Series>) => {
      state.loading = false;
      const index = state.series.findIndex(
        series => series.id === action.payload.id
      );
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
      state.series = state.series.filter(
        series => series.id !== action.payload
      );
      state.pagination.totalItems = Math.max(0, state.pagination.totalItems - 1); // Decrement total count
      state.pagination.totalPages = Math.ceil(
        state.pagination.totalItems / state.pagination.pageSize
      );
    },
    deleteSeriesFailure: (state, action: PayloadAction<string>) => {
      state.loading = false;
      state.error = action.payload;
    },
    setPage: (state, action: PayloadAction<number>) => {
      state.pagination.currentPage = action.payload;
    },
    setPageSize: (state, action: PayloadAction<number>) => {
      state.pagination.pageSize = action.payload;
      state.pagination.currentPage = 1; // Reset to first page when page size changes
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
  setPage,
  setPageSize,
} = seriesSlice.actions;

export default seriesSlice.reducer;
