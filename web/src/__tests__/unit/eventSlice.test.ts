import { configureStore } from '@reduxjs/toolkit';
import eventReducer, {
    initializeMatchEvents,
    clearMatchEvents,
    addEvent,
    addPreviousEvents,
    setMatchLoading,
    setMatchError,
    clearAllEvents,
    removeMatchEvents,
    selectEventsForMatch,
    selectIsLoadingEventsForMatch,
    selectErrorForMatch,
    BallEvent,
    EventState,
} from '@/store/reducers/eventSlice';

describe('Event Slice', () => {
    let store: ReturnType<typeof configureStore>;

    beforeEach(() => {
        store = configureStore({
            reducer: {
                events: eventReducer,
            },
        });
    });

    const mockBallEvent: BallEvent = {
        event_type: 'ball_added',
        match_id: 'match-1',
        innings_number: 1,
        ball_number: 1,
        ball_type: 'good',
        run_type: '1',
        runs: 1,
        byes: 0,
        total_runs: 1,
        is_wicket: false,
        wicket_type: '',
        innings_runs: 10,
        innings_wickets: 1,
        innings_overs: '1.1',
        timestamp: '2024-01-01T10:00:00Z',
        stream_id: 'stream-1',
    };

    const mockBallEvent2: BallEvent = {
        ...mockBallEvent,
        match_id: 'match-2',
        ball_number: 2,
        stream_id: 'stream-2',
    };

    describe('Event Isolation by Match ID', () => {
        it('should initialize events for different matches separately', () => {
            store.dispatch(initializeMatchEvents('match-1'));
            store.dispatch(initializeMatchEvents('match-2'));

            const state = store.getState() as { events: EventState };
            const match1Events = selectEventsForMatch(state, 'match-1');
            const match2Events = selectEventsForMatch(state, 'match-2');

            expect(match1Events).toEqual([]);
            expect(match2Events).toEqual([]);
            expect(match1Events).not.toBe(match2Events); // Different arrays
        });

        it('should add events to the correct match only', () => {
            store.dispatch(initializeMatchEvents('match-1'));
            store.dispatch(initializeMatchEvents('match-2'));

            // Add event to match-1
            store.dispatch(addEvent({ matchId: 'match-1', event: mockBallEvent }));

            const state = store.getState() as { events: EventState };
            const match1Events = selectEventsForMatch(state, 'match-1');
            const match2Events = selectEventsForMatch(state, 'match-2');

            expect(match1Events).toHaveLength(1);
            expect(match1Events[0]).toEqual(mockBallEvent);
            expect(match2Events).toHaveLength(0);
        });

        it('should not mix events between different matches', () => {
            store.dispatch(initializeMatchEvents('match-1'));
            store.dispatch(initializeMatchEvents('match-2'));

            // Add events to both matches
            store.dispatch(addEvent({ matchId: 'match-1', event: mockBallEvent }));
            store.dispatch(addEvent({ matchId: 'match-2', event: mockBallEvent2 }));

            const state = store.getState() as { events: EventState };
            const match1Events = selectEventsForMatch(state, 'match-1');
            const match2Events = selectEventsForMatch(state, 'match-2');

            expect(match1Events).toHaveLength(1);
            expect(match1Events[0]?.match_id).toBe('match-1');
            expect(match2Events).toHaveLength(1);
            expect(match2Events[0]?.match_id).toBe('match-2');

            // Verify no cross-contamination
            expect(match1Events[0]?.match_id).not.toBe(match2Events[0]?.match_id);
        });

        it('should clear events for specific match only', () => {
            store.dispatch(initializeMatchEvents('match-1'));
            store.dispatch(initializeMatchEvents('match-2'));

            // Add events to both matches
            store.dispatch(addEvent({ matchId: 'match-1', event: mockBallEvent }));
            store.dispatch(addEvent({ matchId: 'match-2', event: mockBallEvent2 }));

            // Clear events for match-1 only
            store.dispatch(clearMatchEvents('match-1'));

            const state = store.getState() as { events: EventState };
            const match1Events = selectEventsForMatch(state, 'match-1');
            const match2Events = selectEventsForMatch(state, 'match-2');

            expect(match1Events).toHaveLength(0);
            expect(match2Events).toHaveLength(1);
            expect(match2Events[0]?.match_id).toBe('match-2');
        });

        it('should handle loading states per match', () => {
            store.dispatch(initializeMatchEvents('match-1'));
            store.dispatch(initializeMatchEvents('match-2'));

            // Set loading for match-1 only
            store.dispatch(setMatchLoading({ matchId: 'match-1', loading: true }));

            const state = store.getState() as { events: EventState };
            const match1Loading = selectIsLoadingEventsForMatch(state, 'match-1');
            const match2Loading = selectIsLoadingEventsForMatch(state, 'match-2');

            expect(match1Loading).toBe(true);
            expect(match2Loading).toBe(false);
        });

        it('should handle error states per match', () => {
            store.dispatch(initializeMatchEvents('match-1'));
            store.dispatch(initializeMatchEvents('match-2'));

            // Set error for match-1 only
            store.dispatch(setMatchError({ matchId: 'match-1', error: 'Test error' }));

            const state = store.getState() as { events: EventState };
            const match1Error = selectErrorForMatch(state, 'match-1');
            const match2Error = selectErrorForMatch(state, 'match-2');

            expect(match1Error).toBe('Test error');
            expect(match2Error).toBe(null);
        });

        it('should add previous events to the correct match', () => {
            store.dispatch(initializeMatchEvents('match-1'));
            store.dispatch(initializeMatchEvents('match-2'));

            const previousEvents = [mockBallEvent, mockBallEvent2];

            // Add previous events to match-1
            store.dispatch(addPreviousEvents({ matchId: 'match-1', events: previousEvents }));

            const state = store.getState() as { events: EventState };
            const match1Events = selectEventsForMatch(state, 'match-1');
            const match2Events = selectEventsForMatch(state, 'match-2');

            expect(match1Events).toHaveLength(2);
            expect(match2Events).toHaveLength(0);
        });

        it('should remove events for specific match', () => {
            store.dispatch(initializeMatchEvents('match-1'));
            store.dispatch(initializeMatchEvents('match-2'));

            // Add events to both matches
            store.dispatch(addEvent({ matchId: 'match-1', event: mockBallEvent }));
            store.dispatch(addEvent({ matchId: 'match-2', event: mockBallEvent2 }));

            // Remove match-1 events completely
            store.dispatch(removeMatchEvents('match-1'));

            const state = store.getState() as { events: EventState };
            const match1Events = selectEventsForMatch(state, 'match-1');
            const match2Events = selectEventsForMatch(state, 'match-2');

            expect(match1Events).toHaveLength(0);
            expect(match2Events).toHaveLength(1);
        });

        it('should clear all events across all matches', () => {
            store.dispatch(initializeMatchEvents('match-1'));
            store.dispatch(initializeMatchEvents('match-2'));

            // Add events to both matches
            store.dispatch(addEvent({ matchId: 'match-1', event: mockBallEvent }));
            store.dispatch(addEvent({ matchId: 'match-2', event: mockBallEvent2 }));

            // Clear all events
            store.dispatch(clearAllEvents());

            const state = store.getState() as { events: EventState };
            const match1Events = selectEventsForMatch(state, 'match-1');
            const match2Events = selectEventsForMatch(state, 'match-2');

            expect(match1Events).toHaveLength(0);
            expect(match2Events).toHaveLength(0);
        });

        it('should maintain event order within a match', () => {
            store.dispatch(initializeMatchEvents('match-1'));

            const event1 = { ...mockBallEvent, ball_number: 1, timestamp: '2024-01-01T10:00:00Z' };
            const event2 = { ...mockBallEvent, ball_number: 2, timestamp: '2024-01-01T10:01:00Z' };
            const event3 = { ...mockBallEvent, ball_number: 3, timestamp: '2024-01-01T10:02:00Z' };

            // Add events in order
            store.dispatch(addEvent({ matchId: 'match-1', event: event1 }));
            store.dispatch(addEvent({ matchId: 'match-1', event: event2 }));
            store.dispatch(addEvent({ matchId: 'match-1', event: event3 }));

            const state = store.getState() as { events: EventState };
            const match1Events = selectEventsForMatch(state, 'match-1');

            expect(match1Events).toHaveLength(3);
            expect(match1Events[0]).toEqual(event3); // Most recent first
            expect(match1Events[1]).toEqual(event2);
            expect(match1Events[2]).toEqual(event1);
        });

        it('should limit events to 50 per match', () => {
            store.dispatch(initializeMatchEvents('match-1'));

            // Add 55 events
            for (let i = 1; i <= 55; i++) {
                const event = { ...mockBallEvent, ball_number: i, stream_id: `stream-${i}` };
                store.dispatch(addEvent({ matchId: 'match-1', event }));
            }

            const state = store.getState() as { events: EventState };
            const match1Events = selectEventsForMatch(state, 'match-1');

            expect(match1Events).toHaveLength(50);
            // Should keep the most recent 50 events
            expect(match1Events[0]?.ball_number).toBe(55);
            expect(match1Events[49]?.ball_number).toBe(6);
        });
    });
});
