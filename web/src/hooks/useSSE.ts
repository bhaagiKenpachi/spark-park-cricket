import { useEffect, useRef, useState, useCallback } from 'react';
import { useUserActivity } from './useUserActivity';
import { sseConfig } from '@/config/sseConfig';

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

export interface SSEConnectionState {
    isConnected: boolean;
    isConnecting: boolean;
    error: string | null;
    lastEvent: BallEvent | null;
    eventCount: number;
    connect: () => void;
    disconnect: () => void;
    manualReconnect: () => void;
    // Inactivity management states
    isIdle: boolean;
    needsManualRefresh: boolean;
    disconnectReason: string | null;
    lastEventTime: Date | null;
    connectionStartTime: Date | null;
}

export interface UseSSEOptions {
    onEvent?: (event: BallEvent) => void;
    onConnect?: () => void;
    onDisconnect?: () => void;
    onError?: (error: string) => void;
    autoConnect?: boolean;
    reconnectInterval?: number;
    maxReconnectAttempts?: number;
}

export function useSSE(
    matchId: string,
    options: UseSSEOptions = {}
): SSEConnectionState {
    const {
        onEvent,
        onConnect,
        onDisconnect,
        onError,
        autoConnect = true,
        reconnectInterval = 3000,
        maxReconnectAttempts = 5,
    } = options;

    const [state, setState] = useState<Omit<SSEConnectionState, 'connect' | 'disconnect' | 'manualReconnect'>>({
        isConnected: false,
        isConnecting: false,
        error: null,
        lastEvent: null,
        eventCount: 0,
        // Inactivity management states
        isIdle: false,
        needsManualRefresh: false,
        disconnectReason: null,
        lastEventTime: null,
        connectionStartTime: null,
    });

    const eventSourceRef = useRef<EventSource | null>(null);
    const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null);
    const reconnectAttemptsRef = useRef(0);
    const isMountedRef = useRef(true);
    const hasConnectedRef = useRef(false);
    const inactivityCheckIntervalRef = useRef<NodeJS.Timeout | null>(null);

    // Track user activity
    const { isUserActive } = useUserActivity({
        activityTimeout: sseConfig.userActivityTimeout,
        debounceMs: sseConfig.debounceMs,
    });

    const cleanup = useCallback(() => {
        if (reconnectTimeoutRef.current) {
            clearTimeout(reconnectTimeoutRef.current);
            reconnectTimeoutRef.current = null;
        }
        if (inactivityCheckIntervalRef.current) {
            clearInterval(inactivityCheckIntervalRef.current);
            inactivityCheckIntervalRef.current = null;
        }
        if (eventSourceRef.current) {
            eventSourceRef.current.close();
            eventSourceRef.current = null;
        }
        hasConnectedRef.current = false;
    }, []);

    // Inactivity detection function
    const checkInactivity = useCallback(() => {
        if (!isMountedRef.current) return;

        setState(currentState => {
            if (!currentState.isConnected) return currentState;

            const now = Date.now();

            // Use connection time if no ball events have been received yet
            const referenceTime = currentState.lastEventTime || currentState.connectionStartTime;
            if (!referenceTime) return currentState;

            const timeSinceLastActivity = now - referenceTime.getTime();

            // Check for inactivity warning
            if (timeSinceLastActivity >= sseConfig.inactivityWarningThreshold && !currentState.isIdle) {
                return { ...currentState, isIdle: true };
            }

            // Check for auto-disconnect
            if (timeSinceLastActivity >= sseConfig.inactivityDisconnectThreshold && currentState.isConnected) {
                cleanup();
                return {
                    ...currentState,
                    needsManualRefresh: true,
                    disconnectReason: 'inactivity',
                    isConnected: false
                };
            }

            return currentState;
        });
    }, [cleanup]);

    const connect = useCallback(() => {

        if (!matchId || !isMountedRef.current) {
            return;
        }

        // Don't connect if already connected or connecting
        if (state.isConnected || state.isConnecting) {
            return;
        }

        // Reset inactivity states when connecting
        setState(prev => ({
            ...prev,
            isConnecting: true,
            error: null,
            isIdle: false,
            needsManualRefresh: false,
            disconnectReason: null,
            connectionStartTime: null,
            lastEventTime: null,
        }));

        try {
            // Get API base URL from environment
            const apiBaseUrl = process.env.NEXT_PUBLIC_API_URL!;
            const sseUrl = `${apiBaseUrl}/sse/matches/${matchId}/balls`;

            const eventSource = new EventSource(sseUrl);
            eventSourceRef.current = eventSource;

            eventSource.onopen = () => {
                if (!isMountedRef.current) return;

                const connectionTime = new Date();
                setState(prev => ({
                    ...prev,
                    isConnected: true,
                    isConnecting: false,
                    error: null,
                    connectionStartTime: connectionTime,
                    lastEventTime: connectionTime, // Initialize with connection time
                }));

                reconnectAttemptsRef.current = 0;

                // Start inactivity monitoring
                inactivityCheckIntervalRef.current = setInterval(checkInactivity, sseConfig.activityCheckInterval);

                onConnect?.();
            };

            // Listen for ball_added events (custom event type)
            eventSource.addEventListener('ball_added', (event) => {
                if (!isMountedRef.current) return;

                try {
                    const data = JSON.parse(event.data);

                    // Update last event time
                    const now = new Date();
                    setState(prev => ({
                        ...prev,
                        lastEventTime: now,
                        isIdle: false // Reset idle state on any event
                    }));

                    const ballEvent: BallEvent = {
                        event_type: data.event_type || 'ball_added',
                        match_id: data.match_id,
                        innings_number: parseInt(data.innings_number) || 1,
                        ball_number: parseInt(data.ball_number) || 1,
                        ball_type: data.ball_type || 'good',
                        run_type: data.run_type || '1',
                        runs: parseInt(data.runs) || 0,
                        byes: parseInt(data.byes) || 0,
                        total_runs: parseInt(data.total_runs) || 0,
                        is_wicket: data.is_wicket === 'true' || data.is_wicket === true,
                        wicket_type: data.wicket_type || '',
                        innings_runs: parseInt(data.innings_runs) || 0,
                        innings_wickets: parseInt(data.innings_wickets) || 0,
                        innings_overs: data.innings_overs || '0.0',
                        timestamp: data.timestamp || new Date().toISOString(),
                        stream_id: data.stream_id || '',
                    };

                    setState(prev => ({
                        ...prev,
                        lastEvent: ballEvent,
                        eventCount: prev.eventCount + 1,
                    }));

                    onEvent?.(ballEvent);
                } catch (error) {
                    console.error('❌ Error parsing ball_added event:', error);
                }
            });

            // Listen for timeout events
            eventSource.addEventListener('timeout', (event) => {
                if (!isMountedRef.current) return;

                setState(prev => ({
                    ...prev,
                    isConnected: false,
                    needsManualRefresh: true,
                    disconnectReason: 'server_timeout'
                }));
                cleanup();
            });

            // Listen for connected events
            eventSource.addEventListener('connected', (event) => {
                if (!isMountedRef.current) return;
                // Connection confirmed
            });

            // Keep onmessage for any generic messages
            eventSource.onmessage = (event) => {
                if (!isMountedRef.current) return;
                // Generic message received
            };


            eventSource.onerror = (error) => {
                if (!isMountedRef.current) return;

                setState(prev => ({
                    ...prev,
                    isConnected: false,
                    isConnecting: false,
                    error: 'Connection error occurred',
                }));

                cleanup();
                onError?.('Connection error occurred');

                // Check if user is active before attempting reconnection
                if (!isUserActive) {
                    setState(prev => ({
                        ...prev,
                        needsManualRefresh: true,
                        disconnectReason: 'user_inactive',
                        error: 'Connection paused due to user inactivity',
                    }));
                    return;
                }

                // Attempt to reconnect if we haven't exceeded max attempts and user is active
                if (reconnectAttemptsRef.current < maxReconnectAttempts) {
                    reconnectAttemptsRef.current++;

                    reconnectTimeoutRef.current = setTimeout(() => {
                        if (isMountedRef.current) {
                            connect();
                        }
                    }, reconnectInterval);
                } else {
                    setState(prev => ({
                        ...prev,
                        needsManualRefresh: true,
                        disconnectReason: 'max_attempts_reached',
                        error: 'Max reconnect attempts reached',
                    }));
                }
            };

        } catch (error) {
            console.error('❌ SSE CONNECTION: Failed to create EventSource:', error);
            const errorMessage = `Failed to connect: ${error instanceof Error ? error.message : 'Unknown error'}`;
            setState(prev => ({
                ...prev,
                isConnecting: false,
                error: errorMessage,
            }));
            onError?.(errorMessage);
        }
    }, [matchId, state.isConnected, state.isConnecting, onEvent, onConnect, onError, reconnectInterval, maxReconnectAttempts, cleanup]);

    const disconnect = useCallback(() => {
        cleanup();
        setState(prev => ({
            ...prev,
            isConnected: false,
            isConnecting: false,
        }));
        onDisconnect?.();
    }, [matchId, cleanup, onDisconnect]);

    // Manual reconnect function that resets all states
    const manualReconnect = useCallback(() => {
        cleanup();
        reconnectAttemptsRef.current = 0;
        setState(prev => ({
            ...prev,
            isConnected: false,
            isConnecting: false,
            error: null,
            isIdle: false,
            needsManualRefresh: false,
            disconnectReason: null,
            connectionStartTime: null,
            lastEventTime: null,
        }));
        connect();
    }, [cleanup, connect]);

    // Auto-connect on mount if enabled
    useEffect(() => {

        // Set mounted flag at the start of effect
        isMountedRef.current = true;

        if (autoConnect && matchId && !hasConnectedRef.current) {
            hasConnectedRef.current = true;
            connect();
        } else {
        }

        return () => {
            isMountedRef.current = false;
            cleanup();
        };
    }, [autoConnect, matchId]); // Remove connect and cleanup from dependencies to prevent re-runs

    // Cleanup on unmount
    useEffect(() => {
        return () => {
            isMountedRef.current = false;
            cleanup();
        };
    }, [cleanup]);

    return {
        ...state,
        connect,
        disconnect,
        manualReconnect,
    };
}