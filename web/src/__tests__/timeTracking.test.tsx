import React from 'react';
import { render, screen } from '@testing-library/react';
import { Provider } from 'react-redux';
import { configureStore } from '@reduxjs/toolkit';
import { TimeTrackingView } from '@/components/TimeTrackingView';
import { TimeTrackingSummary } from '@/components/TimeTrackingSummary';
import { formatDuration, formatTime, formatDateTime } from '@/types/timeTracking';

// Mock store
const mockStore = configureStore({
    reducer: {
        timeTracking: (state = { timeTracking: null, loading: false, error: null }, action) => state,
    },
});

describe('Time Tracking Components', () => {
    const mockMatchId = 'test-match-123';

    test('TimeTrackingView renders loading state', () => {
        const mockStoreWithLoading = configureStore({
            reducer: {
                timeTracking: () => ({ timeTracking: null, loading: true, error: null }),
            },
        });

        render(
            <Provider store={mockStoreWithLoading}>
                <TimeTrackingView matchId={mockMatchId} onBack={() => { }} />
            </Provider>
        );

        expect(screen.getByText('Loading time tracking data...')).toBeInTheDocument();
    });

    test('TimeTrackingView renders error state', () => {
        const mockStoreWithError = configureStore({
            reducer: {
                timeTracking: () => ({
                    timeTracking: null,
                    loading: false,
                    error: 'Failed to fetch data'
                }),
            },
        });

        render(
            <Provider store={mockStoreWithError}>
                <TimeTrackingView matchId={mockMatchId} onBack={() => { }} />
            </Provider>
        );

        expect(screen.getByText('Error Loading Time Tracking')).toBeInTheDocument();
        expect(screen.getByText('Failed to fetch data')).toBeInTheDocument();
    });

    test('TimeTrackingView renders no data state', () => {
        render(
            <Provider store={mockStore}>
                <TimeTrackingView matchId={mockMatchId} onBack={() => { }} />
            </Provider>
        );

        expect(screen.getByText('No Time Tracking Data')).toBeInTheDocument();
    });

    test('TimeTrackingSummary renders loading state', () => {
        const mockStoreWithLoading = configureStore({
            reducer: {
                timeTracking: () => ({ timeTracking: null, loading: true, error: null }),
            },
        });

        render(
            <Provider store={mockStoreWithLoading}>
                <TimeTrackingSummary matchId={mockMatchId} />
            </Provider>
        );

        expect(screen.getByText('Loading time data...')).toBeInTheDocument();
    });

    test('TimeTrackingSummary renders no data state', () => {
        render(
            <Provider store={mockStore}>
                <TimeTrackingSummary matchId={mockMatchId} />
            </Provider>
        );

        expect(screen.getByText('Time tracking data not available')).toBeInTheDocument();
    });
});

describe('Time Tracking Utility Functions', () => {
    test('formatDuration formats seconds correctly', () => {
        expect(formatDuration(30)).toBe('30s');
        expect(formatDuration(90)).toBe('1m 30s');
        expect(formatDuration(3600)).toBe('1h');
        expect(formatDuration(3661)).toBe('1h 1m');
    });

    test('formatTime formats time strings correctly', () => {
        const testTime = '2024-01-15T10:30:00Z';
        const formatted = formatTime(testTime);
        expect(formatted).toMatch(/\d{2}:\d{2}:\d{2}/);
    });

    test('formatTime handles null input', () => {
        expect(formatTime(null)).toBe('Not started');
    });

    test('formatDateTime formats date and time correctly', () => {
        const testTime = '2024-01-15T10:30:00Z';
        const formatted = formatDateTime(testTime);
        expect(formatted).toContain('2024');
        expect(formatted).toContain('Jan');
    });

    test('formatDateTime handles null input', () => {
        expect(formatDateTime(null)).toBe('Not started');
    });
});


