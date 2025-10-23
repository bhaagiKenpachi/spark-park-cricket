import { createSlice, PayloadAction } from '@reduxjs/toolkit';

export interface BallEvent {
    event_type: string;
    match_id: string;
    innings_number: number;
    ball_number: number;
    ball_type: string;
    run_type: string;
    runs: number;
    byes: number;
    total_runs: number;
    is_wicket: boolean;
    wicket_type: string;
    innings_runs: number;
    innings_wickets: number;
    innings_overs: string;
    timestamp: string;
    stream_id: string;
}

interface MatchEventsState {
    events: BallEvent[];
    loading: boolean;
    error: string | null;
}

export interface EventState {
    // Events organized by match ID
    eventsByMatchId: Record<string, MatchEventsState>;
    // Global state for the event slice
    loading: boolean;
    error: string | null;
}

const initialState: EventState = {
    eventsByMatchId: {},
    loading: false,
    error: null,
};

export const eventSlice = createSlice({
    name: 'events',
    initialState,
    reducers: {
        // Initialize events state for a match
        initializeMatchEvents: (state, action: PayloadAction<string>) => {
            const matchId = action.payload;
            if (!state.eventsByMatchId[matchId]) {
                state.eventsByMatchId[matchId] = {
                    events: [],
                    loading: false,
                    error: null,
                };
            }
        },

        // Clear events for a specific match
        clearMatchEvents: (state, action: PayloadAction<string>) => {
            const matchId = action.payload;
            if (state.eventsByMatchId[matchId]) {
                state.eventsByMatchId[matchId].events = [];
                state.eventsByMatchId[matchId].error = null;
            }
        },

        // Add a new event to a specific match
        addEvent: (state, action: PayloadAction<{ matchId: string; event: BallEvent }>) => {
            const { matchId, event } = action.payload;

            // Initialize match events if not exists
            if (!state.eventsByMatchId[matchId]) {
                state.eventsByMatchId[matchId] = {
                    events: [],
                    loading: false,
                    error: null,
                };
            }

            // Add event to the beginning and keep only last 50 events
            state.eventsByMatchId[matchId].events = [event, ...state.eventsByMatchId[matchId].events.slice(0, 49)];
        },

        // Add multiple events (for loading previous events)
        addPreviousEvents: (state, action: PayloadAction<{ matchId: string; events: BallEvent[] }>) => {
            const { matchId, events } = action.payload;

            // Initialize match events if not exists
            if (!state.eventsByMatchId[matchId]) {
                state.eventsByMatchId[matchId] = {
                    events: [],
                    loading: false,
                    error: null,
                };
            }

            // Add previous events to the end (they are already sorted chronologically)
            state.eventsByMatchId[matchId].events = [...events, ...state.eventsByMatchId[matchId].events];
        },

        // Set loading state for a match
        setMatchLoading: (state, action: PayloadAction<{ matchId: string; loading: boolean }>) => {
            const { matchId, loading } = action.payload;

            if (!state.eventsByMatchId[matchId]) {
                state.eventsByMatchId[matchId] = {
                    events: [],
                    loading: false,
                    error: null,
                };
            }

            state.eventsByMatchId[matchId].loading = loading;
        },

        // Set error state for a match
        setMatchError: (state, action: PayloadAction<{ matchId: string; error: string | null }>) => {
            const { matchId, error } = action.payload;

            if (!state.eventsByMatchId[matchId]) {
                state.eventsByMatchId[matchId] = {
                    events: [],
                    loading: false,
                    error: null,
                };
            }

            state.eventsByMatchId[matchId].error = error;
        },

        // Clear all events (useful for cleanup)
        clearAllEvents: (state) => {
            state.eventsByMatchId = {};
            state.loading = false;
            state.error = null;
        },

        // Remove events for a specific match (cleanup)
        removeMatchEvents: (state, action: PayloadAction<string>) => {
            const matchId = action.payload;
            delete state.eventsByMatchId[matchId];
        },
    },
});

export const {
    initializeMatchEvents,
    clearMatchEvents,
    addEvent,
    addPreviousEvents,
    setMatchLoading,
    setMatchError,
    clearAllEvents,
    removeMatchEvents,
} = eventSlice.actions;

// Selectors
export const selectEventsByMatchId = (state: { events: EventState }, matchId: string) => {
    return state.events.eventsByMatchId[matchId] || { events: [], loading: false, error: null };
};

export const selectEventsForMatch = (state: { events: EventState }, matchId: string): BallEvent[] => {
    return state.events.eventsByMatchId[matchId]?.events || [];
};

export const selectIsLoadingEventsForMatch = (state: { events: EventState }, matchId: string): boolean => {
    return state.events.eventsByMatchId[matchId]?.loading || false;
};

export const selectErrorForMatch = (state: { events: EventState }, matchId: string): string | null => {
    return state.events.eventsByMatchId[matchId]?.error || null;
};

export default eventSlice.reducer;
