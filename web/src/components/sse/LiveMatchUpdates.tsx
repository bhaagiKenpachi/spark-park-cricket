'use client';

import React, { useState, useCallback, useMemo } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Separator } from '@/components/ui/separator';
import {
    Activity,
    Wifi,
    WifiOff,
    RefreshCw,
    Clock,
    Target,
    Zap,
    AlertCircle,
    CheckCircle,
    XCircle
} from 'lucide-react';
import { useSSE, UseSSEConfig } from '@/hooks/useSSE';
import { SSEBallEvent, SSEConnectedEvent, SSEMatchEvent } from '@/services/sseService';

interface LiveMatchUpdatesProps {
    matchId: string;
    enabled?: boolean;
}

interface EventItem {
    id: string;
    type: 'connected' | 'ball_added' | 'match_event' | 'error' | 'disconnect';
    timestamp: Date;
    data: any;
}

export function LiveMatchUpdates({ matchId, enabled = true }: LiveMatchUpdatesProps) {
    const [events, setEvents] = useState<EventItem[]>([]);
    const [maxEvents] = useState(100); // Keep last 100 events

    const addEvent = useCallback((type: EventItem['type'], data: any) => {
        const event: EventItem = {
            id: `${Date.now()}-${Math.random()}`,
            type,
            timestamp: new Date(),
            data,
        };

        setEvents(prev => {
            const newEvents = [event, ...prev].slice(0, maxEvents);
            return newEvents;
        });
    }, [maxEvents]);

    const handleConnected = useCallback((data: SSEConnectedEvent) => {
        addEvent('connected', data);
    }, [addEvent]);

    const handleBallAdded = useCallback((data: SSEBallEvent) => {
        addEvent('ball_added', data);
    }, [addEvent]);

    const handleMatchEvent = useCallback((data: SSEMatchEvent) => {
        addEvent('match_event', data);
    }, [addEvent]);

    const handleError = useCallback((error: Error) => {
        addEvent('error', { message: error.message });
    }, [addEvent]);

    const handleDisconnect = useCallback(() => {
        addEvent('disconnect', { message: 'Connection lost' });
    }, [addEvent]);

    const sseConfig: UseSSEConfig = {
        matchId,
        endpoint: 'balls',
        enabled,
        onConnected: handleConnected,
        onBallAdded: handleBallAdded,
        onMatchEvent: handleMatchEvent,
        onError: handleError,
        onDisconnect: handleDisconnect,
        autoReconnect: true,
        maxReconnectAttempts: 5,
        reconnectDelay: 3000,
    };

    const { isConnected, connectionState, error, reconnectAttempts, disconnect, reconnect } = useSSE(sseConfig);

    const clearEvents = useCallback(() => {
        setEvents([]);
    }, []);

    const getConnectionStatus = useMemo(() => {
        switch (connectionState) {
            case 'CONNECTED':
                return {
                    icon: <Wifi className="h-4 w-4 text-green-500" />,
                    text: 'Connected',
                    color: 'bg-green-100 text-green-800 border-green-200'
                };
            case 'CONNECTING':
                return {
                    icon: <RefreshCw className="h-4 w-4 text-blue-500 animate-spin" />,
                    text: 'Connecting...',
                    color: 'bg-blue-100 text-blue-800 border-blue-200'
                };
            case 'ERROR':
                return {
                    icon: <WifiOff className="h-4 w-4 text-red-500" />,
                    text: `Error (${reconnectAttempts}/5)`,
                    color: 'bg-red-100 text-red-800 border-red-200'
                };
            default:
                return {
                    icon: <WifiOff className="h-4 w-4 text-gray-500" />,
                    text: 'Disconnected',
                    color: 'bg-gray-100 text-gray-800 border-gray-200'
                };
        }
    }, [connectionState, reconnectAttempts]);

    const formatTimestamp = (timestamp: Date) => {
        return timestamp.toLocaleTimeString('en-US', {
            hour12: false,
            hour: '2-digit',
            minute: '2-digit',
            second: '2-digit'
        }) + '.' + timestamp.getMilliseconds().toString().padStart(3, '0');
    };

    const getEventIcon = (type: EventItem['type']) => {
        switch (type) {
            case 'connected':
                return <CheckCircle className="h-4 w-4 text-green-500" />;
            case 'ball_added':
                return <Target className="h-4 w-4 text-blue-500" />;
            case 'match_event':
                return <Zap className="h-4 w-4 text-yellow-500" />;
            case 'error':
                return <XCircle className="h-4 w-4 text-red-500" />;
            case 'disconnect':
                return <AlertCircle className="h-4 w-4 text-orange-500" />;
            default:
                return <Activity className="h-4 w-4 text-gray-500" />;
        }
    };

    const getEventColor = (type: EventItem['type']) => {
        switch (type) {
            case 'connected':
                return 'bg-green-50 border-green-200';
            case 'ball_added':
                return 'bg-blue-50 border-blue-200';
            case 'match_event':
                return 'bg-yellow-50 border-yellow-200';
            case 'error':
                return 'bg-red-50 border-red-200';
            case 'disconnect':
                return 'bg-orange-50 border-orange-200';
            default:
                return 'bg-gray-50 border-gray-200';
        }
    };

    const renderBallEvent = (data: SSEBallEvent) => (
        <div className="space-y-2">
            <div className="flex items-center gap-2">
                <Badge variant="outline" className="text-xs">
                    Innings {data.innings_number}
                </Badge>
                <Badge variant="outline" className="text-xs">
                    Ball {data.ball_number}
                </Badge>
                <Badge variant="outline" className="text-xs">
                    {data.ball_type}
                </Badge>
            </div>

            <div className="grid grid-cols-2 gap-2 text-sm">
                <div>
                    <span className="font-medium">Runs:</span> {data.runs}
                    {data.byes > 0 && <span className="text-muted-foreground"> + {data.byes} byes</span>}
                </div>
                <div>
                    <span className="font-medium">Total:</span> {data.total_runs}
                </div>
                <div>
                    <span className="font-medium">Type:</span> {data.run_type}
                </div>
                <div>
                    <span className="font-medium">Wicket:</span> {data.is_wicket ? 'Yes' : 'No'}
                    {data.is_wicket && data.wicket_type && (
                        <span className="text-muted-foreground"> ({data.wicket_type})</span>
                    )}
                </div>
            </div>

            <div className="pt-2 border-t">
                <div className="text-sm text-muted-foreground">
                    <span className="font-medium">Innings Total:</span> {data.innings_runs}/{data.innings_wickets}
                    <span className="ml-2">({data.innings_overs} overs)</span>
                </div>
            </div>
        </div>
    );

    const renderEventContent = (event: EventItem) => {
        switch (event.type) {
            case 'connected':
                return (
                    <div className="text-sm">
                        <div className="font-medium text-green-700">Connected to match stream</div>
                        <div className="text-muted-foreground">{event.data.message}</div>
                    </div>
                );

            case 'ball_added':
                return renderBallEvent(event.data);

            case 'match_event':
                return (
                    <div className="text-sm">
                        <div className="font-medium text-yellow-700">Match Event</div>
                        <div className="text-muted-foreground">{event.data.message}</div>
                    </div>
                );

            case 'error':
                return (
                    <div className="text-sm">
                        <div className="font-medium text-red-700">Connection Error</div>
                        <div className="text-muted-foreground">{event.data.message}</div>
                    </div>
                );

            case 'disconnect':
                return (
                    <div className="text-sm">
                        <div className="font-medium text-orange-700">Disconnected</div>
                        <div className="text-muted-foreground">{event.data.message}</div>
                    </div>
                );

            default:
                return (
                    <div className="text-sm text-muted-foreground">
                        Unknown event type: {event.type}
                    </div>
                );
        }
    };

    return (
        <Card className="h-full">
            <CardHeader className="pb-3">
                <div className="flex items-center justify-between">
                    <CardTitle className="flex items-center gap-2">
                        <Activity className="h-5 w-5" />
                        Live Updates
                    </CardTitle>

                    <div className="flex items-center gap-2">
                        <Badge className={getConnectionStatus.color}>
                            {getConnectionStatus.icon}
                            <span className="ml-1">{getConnectionStatus.text}</span>
                        </Badge>

                        <div className="flex gap-1">
                            <Button
                                variant="outline"
                                size="sm"
                                onClick={reconnect}
                                disabled={isConnected}
                                className="h-8"
                            >
                                <RefreshCw className="h-3 w-3" />
                            </Button>

                            <Button
                                variant="outline"
                                size="sm"
                                onClick={clearEvents}
                                className="h-8"
                            >
                                Clear
                            </Button>
                        </div>
                    </div>
                </div>

                {error && (
                    <div className="text-sm text-red-600 bg-red-50 p-2 rounded border">
                        <strong>Error:</strong> {error.message}
                    </div>
                )}
            </CardHeader>

            <CardContent className="p-0">
                <ScrollArea className="h-[500px]">
                    <div className="p-4 space-y-3">
                        {events.length === 0 ? (
                            <div className="text-center text-muted-foreground py-8">
                                <Activity className="h-12 w-12 mx-auto mb-4 opacity-50" />
                                <p>No events yet</p>
                                <p className="text-sm">Live updates will appear here</p>
                            </div>
                        ) : (
                            events.map((event, index) => (
                                <div key={event.id}>
                                    <div className={`p-3 rounded-lg border ${getEventColor(event.type)}`}>
                                        <div className="flex items-start gap-3">
                                            <div className="flex-shrink-0 mt-0.5">
                                                {getEventIcon(event.type)}
                                            </div>

                                            <div className="flex-1 min-w-0">
                                                <div className="flex items-center gap-2 mb-2">
                                                    <span className="text-xs font-medium uppercase tracking-wide">
                                                        {event.type.replace('_', ' ')}
                                                    </span>
                                                    <Separator orientation="vertical" className="h-3" />
                                                    <div className="flex items-center gap-1 text-xs text-muted-foreground">
                                                        <Clock className="h-3 w-3" />
                                                        {formatTimestamp(event.timestamp)}
                                                    </div>
                                                </div>

                                                {renderEventContent(event)}
                                            </div>
                                        </div>
                                    </div>

                                    {index < events.length - 1 && <Separator className="my-2" />}
                                </div>
                            ))
                        )}
                    </div>
                </ScrollArea>
            </CardContent>
        </Card>
    );
}
