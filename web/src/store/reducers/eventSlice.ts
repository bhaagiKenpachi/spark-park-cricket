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

            const currentEvents = state.eventsByMatchId[matchId].events;

            // Check if event already exists by stream_id to avoid duplicates
            const eventExists = currentEvents.some(e => e.stream_id === event.stream_id);
            if (eventExists) {
                console.log(`⚠️  Event ${event.stream_id} already exists, skipping duplicate`);
                return;
            }

            // Append new events to the END of the array
            // Events are stored in reverse chronological order (oldest at end, newest at beginning after loading)
            // But for real-time updates, we add to the end so when reversed for display, it appears at top
            // Keep only the last 50 events to prevent memory issues
            state.eventsByMatchId[matchId].events = [...currentEvents, event].slice(-50);
        },

        // Add multiple events (for loading previous events - replaces all existing events)
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

            // Replace all existing events with the new previous events
            state.eventsByMatchId[matchId].events = [...events];
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
