/**
 * React hook for SSE (Server-Sent Events) integration
 * Provides easy SSE connection management for cricket match updates
 */

import { useState, useEffect, useRef, useCallback } from 'react';
import {
    SSEConnection,
    createSSEConnection,
    SSEConnectionConfig,
    SSEEventCallbacks,
    SSEBallEvent,
    SSEConnectedEvent,
    SSEMatchEvent
} from '@/services/sseService';

export interface UseSSEConfig extends SSEConnectionConfig {
    onConnected?: (data: SSEConnectedEvent) => void;
    onBallAdded?: (data: SSEBallEvent) => void;
    onMatchEvent?: (data: SSEMatchEvent) => void;
    onError?: (error: Error) => void;
    onDisconnect?: () => void;
}

export interface UseSSEReturn {
    isConnected: boolean;
    connectionState: 'CONNECTING' | 'CONNECTED' | 'DISCONNECTED' | 'ERROR';
    error: Error | null;
    reconnectAttempts: number;
    disconnect: () => void;
    reconnect: () => void;
}

export function useSSE(config: UseSSEConfig): UseSSEReturn {
    const [isConnected, setIsConnected] = useState(false);
    const [connectionState, setConnectionState] = useState<'CONNECTING' | 'CONNECTED' | 'DISCONNECTED' | 'ERROR'>('DISCONNECTED');
    const [error, setError] = useState<Error | null>(null);
    const [reconnectAttempts, setReconnectAttempts] = useState(0);

    const sseConnectionRef = useRef<SSEConnection | null>(null);

    // Use refs for callbacks to avoid dependency issues
    const callbacksRef = useRef({
        onConnected: config.onConnected,
        onBallAdded: config.onBallAdded,
        onMatchEvent: config.onMatchEvent,
        onError: config.onError,
        onDisconnect: config.onDisconnect,
    });

    // Update callbacks ref when they change
    useEffect(() => {
        callbacksRef.current = {
            onConnected: config.onConnected,
            onBallAdded: config.onBallAdded,
            onMatchEvent: config.onMatchEvent,
            onError: config.onError,
            onDisconnect: config.onDisconnect,
        };
    }, [config.onConnected, config.onBallAdded, config.onMatchEvent, config.onError, config.onDisconnect]);

    // Stable callback handlers
    const handleConnected = useCallback((data: SSEConnectedEvent) => {
        setIsConnected(true);
        setConnectionState('CONNECTED');
        setError(null);
        setReconnectAttempts(0);
        callbacksRef.current.onConnected?.(data);
    }, []);

    const handleBallAdded = useCallback((data: SSEBallEvent) => {
        callbacksRef.current.onBallAdded?.(data);
    }, []);

    const handleMatchEvent = useCallback((data: SSEMatchEvent) => {
        callbacksRef.current.onMatchEvent?.(data);
    }, []);

    const handleError = useCallback((err: Error) => {
        setError(err);
        setConnectionState('ERROR');
        callbacksRef.current.onError?.(err);
    }, []);

    const handleDisconnect = useCallback(() => {
        setIsConnected(false);
        setConnectionState('DISCONNECTED');
        callbacksRef.current.onDisconnect?.();
    }, []);

    const disconnect = useCallback(() => {
        if (sseConnectionRef.current) {
            sseConnectionRef.current.disconnect();
            sseConnectionRef.current = null;
        }
        setIsConnected(false);
        setConnectionState('DISCONNECTED');
    }, []);

    const reconnect = useCallback(() => {
        disconnect();
        if (config.enabled !== false && config.matchId) {
            const connection = createSSEConnection(
                {
                    matchId: config.matchId,
                    endpoint: config.endpoint,
                    enabled: config.enabled ?? true,
                    autoReconnect: config.autoReconnect ?? true,
                    maxReconnectAttempts: config.maxReconnectAttempts ?? 5,
                    reconnectDelay: config.reconnectDelay ?? 3000,
                },
                {
                    onConnected: handleConnected,
                    onBallAdded: handleBallAdded,
                    onMatchEvent: handleMatchEvent,
                    onError: handleError,
                    onDisconnect: handleDisconnect,
                }
            );
            sseConnectionRef.current = connection;
            connection.connect();
        }
    }, [
        config.matchId,
        config.endpoint,
        config.enabled,
        config.autoReconnect,
        config.maxReconnectAttempts,
        config.reconnectDelay,
        handleConnected,
        handleBallAdded,
        handleMatchEvent,
        handleError,
        handleDisconnect,
        disconnect
    ]);

    // Main effect for connection management
    useEffect(() => {
        if (config.enabled === false || !config.matchId) {
            disconnect();
            return;
        }

        // Prevent duplicate connections
        if (sseConnectionRef.current) {
            return;
        }

        const connection = createSSEConnection(
            {
                matchId: config.matchId,
                endpoint: config.endpoint,
                enabled: config.enabled ?? true,
                autoReconnect: config.autoReconnect ?? true,
                maxReconnectAttempts: config.maxReconnectAttempts ?? 5,
                reconnectDelay: config.reconnectDelay ?? 3000,
            },
            {
                onConnected: handleConnected,
                onBallAdded: handleBallAdded,
                onMatchEvent: handleMatchEvent,
                onError: handleError,
                onDisconnect: handleDisconnect,
            }
        );

        sseConnectionRef.current = connection;
        connection.connect();

        return () => {
            connection.disconnect();
            sseConnectionRef.current = null;
        };
    }, [
        config.matchId,
        config.endpoint,
        config.enabled,
        config.autoReconnect,
        config.maxReconnectAttempts,
        config.reconnectDelay,
        handleConnected,
        handleBallAdded,
        handleMatchEvent,
        handleError,
        handleDisconnect,
        disconnect
    ]);

    // Update reconnect attempts from connection
    useEffect(() => {
        if (sseConnectionRef.current) {
            setReconnectAttempts(sseConnectionRef.current.getReconnectAttempts());
        }
    }, [isConnected, connectionState]);

    return {
        isConnected,
        connectionState,
        error,
        reconnectAttempts,
        disconnect,
        reconnect,
    };
}
