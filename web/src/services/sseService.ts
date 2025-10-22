/**
 * SSE Service for real-time cricket match updates
 * Handles Server-Sent Events connection and event processing
 */

export interface SSEBallEvent {
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
    wicket_type?: string;
    innings_runs: number;
    innings_wickets: number;
    innings_overs: string;
    timestamp: string;
    stream_id: string;
}

export interface SSEConnectedEvent {
    message: string;
    timestamp: string;
}

export interface SSEMatchEvent {
    event_type: string;
    match_id: string;
    message: string;
    timestamp: string;
}

export interface SSEEventCallbacks {
    onConnected?: (data: SSEConnectedEvent) => void;
    onBallAdded?: (data: SSEBallEvent) => void;
    onMatchEvent?: (data: SSEMatchEvent) => void;
    onError?: (error: Error) => void;
    onDisconnect?: () => void;
}

export interface SSEConnectionConfig {
    matchId: string;
    endpoint: 'balls' | 'events';
    enabled?: boolean;
    autoReconnect?: boolean;
    maxReconnectAttempts?: number;
    reconnectDelay?: number;
}

export class SSEConnection {
    private eventSource: EventSource | null = null;
    private url: string;
    private config: SSEConnectionConfig;
    private callbacks: SSEEventCallbacks;
    private reconnectAttempts = 0;
    private isConnected = false;
    private reconnectTimer: NodeJS.Timeout | null = null;

    constructor(config: SSEConnectionConfig, callbacks: SSEEventCallbacks) {
        this.config = {
            enabled: true,
            autoReconnect: true,
            maxReconnectAttempts: 5,
            reconnectDelay: 3000,
            ...config,
        };
        this.callbacks = callbacks;

        // Build SSE URL
        const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';
        this.url = `${apiUrl}/sse/matches/${config.matchId}/${config.endpoint}`;
    }

    connect(): void {
        if (!this.config.enabled || this.isConnected) {
            return;
        }

        try {
            console.log(`📡 SSE: Connecting to ${this.url}`);
            this.eventSource = new EventSource(this.url);

            this.eventSource.onopen = () => {
                console.log('📡 SSE: Connection opened');
                this.isConnected = true;
                this.reconnectAttempts = 0;
            };

            this.eventSource.addEventListener('connected', (event) => {
                try {
                    const data = JSON.parse(event.data) as SSEConnectedEvent;
                    console.log('📡 SSE: Connected event received:', data);
                    this.callbacks.onConnected?.(data);
                } catch (error) {
                    console.error('📡 SSE: Error parsing connected event:', error);
                }
            });

            this.eventSource.addEventListener('ball_added', (event) => {
                try {
                    const data = JSON.parse(event.data) as SSEBallEvent;
                    console.log('📡 SSE: Ball event received:', data);
                    this.callbacks.onBallAdded?.(data);
                } catch (error) {
                    console.error('📡 SSE: Error parsing ball event:', error);
                }
            });

            this.eventSource.addEventListener('match_event', (event) => {
                try {
                    const data = JSON.parse(event.data) as SSEMatchEvent;
                    console.log('📡 SSE: Match event received:', data);
                    this.callbacks.onMatchEvent?.(data);
                } catch (error) {
                    console.error('📡 SSE: Error parsing match event:', error);
                }
            });

            this.eventSource.onerror = (error) => {
                console.error('📡 SSE: Connection error:', error);
                this.isConnected = false;
                this.callbacks.onError?.(new Error('SSE connection error'));

                if (this.config.autoReconnect && this.reconnectAttempts < this.config.maxReconnectAttempts!) {
                    this.scheduleReconnect();
                } else {
                    this.callbacks.onDisconnect?.();
                }
            };

            this.eventSource.addEventListener('close', () => {
                console.log('📡 SSE: Connection closed');
                this.isConnected = false;
                this.callbacks.onDisconnect?.();
            });

        } catch (error) {
            console.error('📡 SSE: Failed to create connection:', error);
            this.callbacks.onError?.(error as Error);
        }
    }

    private scheduleReconnect(): void {
        if (this.reconnectTimer) {
            clearTimeout(this.reconnectTimer);
        }

        this.reconnectAttempts++;
        const delay = this.config.reconnectDelay! * this.reconnectAttempts;

        console.log(`📡 SSE: Scheduling reconnect attempt ${this.reconnectAttempts}/${this.config.maxReconnectAttempts} in ${delay}ms`);

        this.reconnectTimer = setTimeout(() => {
            this.disconnect();
            this.connect();
        }, delay);
    }

    disconnect(): void {
        if (this.reconnectTimer) {
            clearTimeout(this.reconnectTimer);
            this.reconnectTimer = null;
        }

        if (this.eventSource) {
            console.log('📡 SSE: Disconnecting...');
            this.eventSource.close();
            this.eventSource = null;
        }

        this.isConnected = false;
    }

    getConnectionState(): 'CONNECTING' | 'CONNECTED' | 'DISCONNECTED' | 'ERROR' {
        if (!this.eventSource) return 'DISCONNECTED';
        if (this.isConnected) return 'CONNECTED';
        return 'CONNECTING';
    }

    getReconnectAttempts(): number {
        return this.reconnectAttempts;
    }
}

export function createSSEConnection(
    config: SSEConnectionConfig,
    callbacks: SSEEventCallbacks
): SSEConnection {
    return new SSEConnection(config, callbacks);
}
