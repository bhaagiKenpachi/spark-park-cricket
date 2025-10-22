import { useEffect, useRef, useState, useCallback } from 'react';

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

    const [state, setState] = useState<Omit<SSEConnectionState, 'connect' | 'disconnect'>>({
        isConnected: false,
        isConnecting: false,
        error: null,
        lastEvent: null,
        eventCount: 0,
    });

    const eventSourceRef = useRef<EventSource | null>(null);
    const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null);
    const reconnectAttemptsRef = useRef(0);
    const isMountedRef = useRef(true);
    const hasConnectedRef = useRef(false);

    const cleanup = useCallback(() => {
        console.log('🔌 SSE CONNECTION: Cleaning up connection');
        if (reconnectTimeoutRef.current) {
            clearTimeout(reconnectTimeoutRef.current);
            reconnectTimeoutRef.current = null;
        }
        if (eventSourceRef.current) {
            console.log('🔌 SSE CONNECTION: Closing EventSource');
            eventSourceRef.current.close();
            eventSourceRef.current = null;
        }
        hasConnectedRef.current = false;
    }, []);

    const connect = useCallback(() => {
        console.log('🔌 SSE CONNECTION: Starting connection for match:', matchId);

        if (!matchId || !isMountedRef.current) {
            console.log('❌ SSE CONNECTION: Failed - missing matchId or not mounted');
            return;
        }

        // Don't connect if already connected or connecting
        if (state.isConnected || state.isConnecting) {
            console.log('⚠️ SSE CONNECTION: Skipped - already connected or connecting');
            return;
        }

        setState(prev => ({ ...prev, isConnecting: true, error: null }));

        try {
            // Get API base URL from environment
            const apiBaseUrl = process.env.NEXT_PUBLIC_API_URL!;
            const sseUrl = `${apiBaseUrl}/sse/matches/${matchId}/balls`;

            console.log('🔌 SSE CONNECTION: Connecting to URL:', sseUrl);
            console.log('🔌 SSE CONNECTION: Environment API URL:', apiBaseUrl);

            const eventSource = new EventSource(sseUrl);
            eventSourceRef.current = eventSource;
            console.log('🔌 SSE CONNECTION: EventSource created, readyState:', eventSource.readyState);

            eventSource.onopen = () => {
                if (!isMountedRef.current) return;

                console.log('✅ SSE CONNECTION: Successfully connected for match:', matchId);
                setState(prev => ({
                    ...prev,
                    isConnected: true,
                    isConnecting: false,
                    error: null,
                }));

                reconnectAttemptsRef.current = 0;
                onConnect?.();
            };

            eventSource.onmessage = (event) => {
                if (!isMountedRef.current) return;

                try {
                    const data = JSON.parse(event.data);
                    console.log('🏏 SSE EVENT: Ball event received:', data);

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
                    console.error('❌ Error parsing SSE event:', error);
                    const errorMessage = `Failed to parse event: ${error instanceof Error ? error.message : 'Unknown error'}`;
                    setState(prev => ({ ...prev, error: errorMessage }));
                    onError?.(errorMessage);
                }
            };

            eventSource.addEventListener('ball_added', (event) => {
                if (!isMountedRef.current) return;

                try {
                    const data = JSON.parse(event.data);
                    console.log('🏏 Ball added event:', data);

                    const ballEvent: BallEvent = {
                        event_type: 'ball_added',
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
                    console.error('❌ SSE EVENT: Error parsing ball event:', error);
                }
            });

            eventSource.onerror = (error) => {
                if (!isMountedRef.current) return;

                console.error('❌ SSE CONNECTION ERROR:', error);
                console.error('❌ SSE CONNECTION: EventSource readyState:', eventSource.readyState);
                console.error('❌ SSE CONNECTION: EventSource URL:', eventSource.url);

                if (eventSource.readyState === EventSource.CLOSED) {
                    console.error('❌ SSE CONNECTION: Connection was closed by server');
                } else if (eventSource.readyState === EventSource.CONNECTING) {
                    console.error('❌ SSE CONNECTION: Connection failed during initial connection');
                }

                setState(prev => ({
                    ...prev,
                    isConnected: false,
                    isConnecting: false,
                    error: 'Connection error occurred',
                }));

                cleanup();
                onError?.('Connection error occurred');

                // Attempt to reconnect if we haven't exceeded max attempts
                if (reconnectAttemptsRef.current < maxReconnectAttempts) {
                    reconnectAttemptsRef.current++;
                    console.log(`🔄 SSE CONNECTION: Attempting to reconnect (${reconnectAttemptsRef.current}/${maxReconnectAttempts})...`);

                    reconnectTimeoutRef.current = setTimeout(() => {
                        if (isMountedRef.current) {
                            connect();
                        }
                    }, reconnectInterval);
                } else {
                    console.log('❌ SSE CONNECTION: Max reconnection attempts reached');
                    setState(prev => ({
                        ...prev,
                        error: 'Max reconnection attempts reached',
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
        console.log('🔌 Disconnecting SSE for match:', matchId);
        cleanup();
        setState(prev => ({
            ...prev,
            isConnected: false,
            isConnecting: false,
        }));
        onDisconnect?.();
    }, [matchId, cleanup, onDisconnect]);

    // Auto-connect on mount if enabled
    useEffect(() => {
        console.log('🔌 SSE CONNECTION: useEffect triggered - autoConnect:', autoConnect, 'matchId:', matchId);

        // Set mounted flag at the start of effect
        isMountedRef.current = true;

        if (autoConnect && matchId && !hasConnectedRef.current) {
            console.log('🔌 SSE CONNECTION: Auto-connecting to SSE for match:', matchId);
            hasConnectedRef.current = true;
            connect();
        } else {
            console.log('🔌 SSE CONNECTION: Not auto-connecting - autoConnect:', autoConnect, 'matchId:', matchId, 'hasConnected:', hasConnectedRef.current);
        }

        return () => {
            console.log('🔌 SSE CONNECTION: useEffect cleanup');
            isMountedRef.current = false;
            cleanup();
        };
    }, [autoConnect, matchId]); // Remove connect and cleanup from dependencies to prevent re-runs

    // Cleanup on unmount
    useEffect(() => {
        return () => {
            console.log('🔍 SSE DEBUG: Component unmounting - cleanup');
            isMountedRef.current = false;
            cleanup();
        };
    }, [cleanup]);

    return {
        ...state,
        connect,
        disconnect,
    };
}