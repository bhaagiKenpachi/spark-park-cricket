'use client';

import { useState, useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { AppDispatch, RootState } from '@/store';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Separator } from '@/components/ui/separator';
import {
    Wifi,
    WifiOff,
    Activity,
    Clock,
    Target,
    Zap,
    RefreshCw,
    X,
    Download
} from 'lucide-react';
import { useSSE, BallEvent } from '@/hooks/useSSE';
import { ApiService, BallEventResponse } from '@/services/api';
import { sseConfig } from '@/config/sseConfig';
import { fetchInningsScoreSummaryThunk, fetchLatestOverThunk } from '@/store/reducers/scorecardSlice';
import {
    initializeMatchEvents,
    clearMatchEvents,
    addEvent,
    addPreviousEvents,
    setMatchLoading,
    setMatchError,
    selectEventsForMatch,
    selectIsLoadingEventsForMatch,
    selectErrorForMatch,
    BallEvent as ReduxBallEvent
} from '@/store/reducers/eventSlice';

interface LiveUpdatesProps {
    matchId: string;
    onEvent?: (event: BallEvent) => void;
    className?: string;
}

export function LiveUpdates({ matchId, onEvent, className = '' }: LiveUpdatesProps) {
    const dispatch = useDispatch<AppDispatch>();

    // Use Redux selectors for events
    const events = useSelector((state: RootState) => selectEventsForMatch(state, matchId));
    const isLoadingPrevious = useSelector((state: RootState) => selectIsLoadingEventsForMatch(state, matchId));
    const error = useSelector((state: RootState) => selectErrorForMatch(state, matchId));


    // Function to fetch previous ball events from the dedicated ball events API
    const fetchPreviousEvents = async () => {
        if (isLoadingPrevious) return;

        dispatch(setMatchLoading({ matchId, loading: true }));
        dispatch(setMatchError({ matchId, error: null }));

        // Clear existing events for this match before loading new ones
        dispatch(clearMatchEvents(matchId));

        try {
            const apiService = new ApiService();
            const response = await apiService.getBallEvents(matchId);
            // The response structure is { data: { data: [...], success: true }, success: true }
            // We need to extract the actual ball events array from response.data.data
            const ballEventsResponse = response.data as any;
            const ballEventsArray = Array.isArray(ballEventsResponse?.data)
                ? ballEventsResponse.data
                : (Array.isArray(ballEventsResponse) ? ballEventsResponse : []);

            console.log(`✅ Loaded ${ballEventsArray.length} previous events for match ${matchId}`);

            // Convert API response to Redux ball events
            const previousEvents: ReduxBallEvent[] = ballEventsArray.map((event: BallEventResponse) => ({
                event_type: event.event_type,
                match_id: event.match_id,
                innings_number: event.innings_number,
                ball_number: event.ball_number,
                ball_type: event.ball_type,
                run_type: event.run_type,
                runs: event.runs,
                byes: event.byes,
                total_runs: event.total_runs,
                is_wicket: event.is_wicket,
                wicket_type: event.wicket_type,
                innings_runs: event.innings_runs,
                innings_wickets: event.innings_wickets,
                innings_overs: event.innings_overs,
                timestamp: event.timestamp,
                stream_id: event.stream_id,
            }));

            // Dispatch the previous events to Redux store (replaces all existing events)
            dispatch(addPreviousEvents({ matchId, events: previousEvents }));
            console.log(`📡 Dispatched ${previousEvents.length} events to Redux for match ${matchId}`);

        } catch (error) {
            console.error('❌ Error fetching previous ball events:', error);
            dispatch(setMatchError({ matchId, error: 'Failed to fetch previous events' }));
        } finally {
            dispatch(setMatchLoading({ matchId, loading: false }));
        }
    };

    const {
        isConnected,
        isConnecting,
        error: sseError,
        lastEvent,
        eventCount,
        connect,
        disconnect,
        manualReconnect,
        isIdle,
        needsManualRefresh,
        disconnectReason,
        lastEventTime,
        timeUntilDisconnect,
    } = useSSE(matchId, {
        onEvent: (event) => {
            // Add event to Redux store for this specific match
            dispatch(addEvent({ matchId, event }));

            // Trigger optimized API calls to refresh scorecard data
            // Same as what happens after manual ball addition
            try {
                // Fetch innings score summary for the current innings
                dispatch(fetchInningsScoreSummaryThunk({
                    matchId: event.match_id,
                    inningsNumber: event.innings_number,
                }));

                // Fetch latest over for the current innings
                dispatch(fetchLatestOverThunk({
                    matchId: event.match_id,
                    inningsNumber: event.innings_number,
                }));
            } catch (error) {
                console.error('❌ Error dispatching API calls:', error);
            }

            onEvent?.(event);
        },
        onConnect: () => {
            // Connected to SSE
        },
        onDisconnect: () => {
            // Disconnected from SSE
        },
        onError: (error) => {
            console.error('❌ LiveUpdates: SSE error for match:', matchId, error);
        },
    });

    // Initialize match events when component mounts or matchId changes
    useEffect(() => {
        if (matchId) {
            console.log(`🔄 Initializing events for match: ${matchId}`);
            dispatch(initializeMatchEvents(matchId));
        }
    }, [matchId]);

    const formatTimestamp = (timestamp: string) => {
        try {
            return new Date(timestamp).toLocaleTimeString();
        } catch {
            return 'Invalid time';
        }
    };

    const getBallTypeColor = (ballType: string) => {
        switch (ballType.toLowerCase()) {
            case 'good':
                return 'bg-green-100 text-green-800';
            case 'wide':
                return 'bg-yellow-100 text-yellow-800';
            case 'no_ball':
                return 'bg-red-100 text-red-800';
            case 'dead_ball':
                return 'bg-gray-100 text-gray-800';
            default:
                return 'bg-blue-100 text-blue-800';
        }
    };

    const getRunTypeColor = (runType: string) => {
        switch (runType.toLowerCase()) {
            case '1':
                return 'bg-blue-100 text-blue-800';
            case '2':
                return 'bg-indigo-100 text-indigo-800';
            case '3':
                return 'bg-purple-100 text-purple-800';
            case '4':
                return 'bg-pink-100 text-pink-800';
            case '6':
                return 'bg-orange-100 text-orange-800';
            case 'wd':
                return 'bg-yellow-100 text-yellow-800';
            case 'nb':
                return 'bg-red-100 text-red-800';
            default:
                return 'bg-gray-100 text-gray-800';
        }
    };

    const getConnectionStatus = () => {
        if (isConnecting) {
            return {
                icon: <RefreshCw className="h-4 w-4 animate-spin" />,
                text: 'Connecting...',
                color: 'text-blue-600',
                bgColor: 'bg-blue-50 border-blue-200',
            };
        }
        if (needsManualRefresh) {
            return {
                icon: <WifiOff className="h-4 w-4" />,
                text: 'Paused',
                color: 'text-orange-600',
                bgColor: 'bg-orange-50 border-orange-200',
            };
        }
        if (isConnected && isIdle) {
            return {
                icon: <Wifi className="h-4 w-4" />,
                text: 'Connected - Idle',
                color: 'text-yellow-600',
                bgColor: 'bg-yellow-50 border-yellow-200',
            };
        }
        if (isConnected) {
            return {
                icon: <Wifi className="h-4 w-4" />,
                text: 'Live',
                color: 'text-green-600',
                bgColor: 'bg-green-50 border-green-200',
            };
        }
        return {
            icon: <WifiOff className="h-4 w-4" />,
            text: 'Disconnected',
            color: 'text-red-600',
            bgColor: 'bg-red-50 border-red-200',
        };
    };

    const status = getConnectionStatus();


    return (
        <div className={`space-y-4 ${className}`}>
            {/* Mobile-First Control Panel - Compact */}
            <div className="bg-gradient-to-r from-slate-50 to-gray-50 border border-slate-200 rounded-lg p-2 sm:p-3">
                {/* Mobile Layout */}
                <div className="flex flex-col space-y-2 sm:hidden">
                    {/* Status Row */}
                    <div className="flex items-center justify-between">
                        <div className="flex items-center gap-2">
                            <div className={`p-1.5 rounded-full ${status.bgColor}`}>
                                {status.icon}
                            </div>
                            <div>
                                <span className={`text-sm font-semibold ${status.color}`}>{status.text}</span>
                                {eventCount > 0 && (
                                    <Badge variant="secondary" className="ml-2 text-xs">
                                        {eventCount}
                                    </Badge>
                                )}
                            </div>
                        </div>

                        {/* Action Button */}
                        <button
                            onClick={needsManualRefresh ? manualReconnect : (isConnected ? disconnect : connect)}
                            disabled={isConnecting}
                            className={`p-2 rounded-full transition-colors disabled:opacity-50 disabled:cursor-not-allowed ${isConnected
                                ? 'bg-red-50 hover:bg-red-100 border border-red-200 text-red-700'
                                : 'bg-green-50 hover:bg-green-100 border border-green-200 text-green-700'
                                }`}
                            title={needsManualRefresh ? 'Reconnect to Live Events' : isConnected ? 'Disconnect from Live Events' : 'Connect to Live Events'}
                        >
                            {needsManualRefresh ? (
                                <RefreshCw className={`h-4 w-4 ${isConnecting ? 'animate-spin' : ''}`} />
                            ) : isConnected ? (
                                <WifiOff className="h-4 w-4" />
                            ) : (
                                <Wifi className="h-4 w-4" />
                            )}
                        </button>
                    </div>

                    {/* Last Event Info */}
                    {lastEvent && (
                        <div className="flex items-center gap-1 text-xs text-gray-500">
                            <Clock className="h-3 w-3" />
                            <span>Last: {formatTimestamp(lastEvent.timestamp)}</span>
                        </div>
                    )}
                </div>

                {/* Desktop Layout */}
                <div className="hidden sm:flex items-center justify-between">
                    {/* Status and Info */}
                    <div className="flex items-center gap-4">
                        <div className="flex items-center gap-3">
                            <div className={`p-2 rounded-full ${status.bgColor}`}>
                                {status.icon}
                            </div>
                            <div>
                                <div className="flex items-center gap-2">
                                    <span className={`font-semibold ${status.color}`}>{status.text}</span>
                                    {eventCount > 0 && (
                                        <Badge variant="secondary" className="ml-2">
                                            {eventCount} events
                                        </Badge>
                                    )}
                                </div>
                                {lastEvent && (
                                    <div className="flex items-center gap-1 mt-1">
                                        <Clock className="h-3 w-3 text-gray-400" />
                                        <span className="text-xs text-gray-500">
                                            Last: {formatTimestamp(lastEvent.timestamp)}
                                        </span>
                                    </div>
                                )}
                            </div>
                        </div>
                    </div>

                    {/* Connect/Disconnect Icon Only */}
                    <div className="flex items-center">
                        <button
                            onClick={needsManualRefresh ? manualReconnect : (isConnected ? disconnect : connect)}
                            disabled={isConnecting}
                            className={`p-2 rounded-full transition-colors disabled:opacity-50 disabled:cursor-not-allowed ${isConnected
                                ? 'bg-red-50 hover:bg-red-100 border border-red-200 text-red-700'
                                : 'bg-green-50 hover:bg-green-100 border border-green-200 text-green-700'
                                }`}
                            title={needsManualRefresh ? 'Reconnect' : isConnected ? 'Disconnect' : 'Connect'}
                        >
                            {needsManualRefresh ? (
                                <RefreshCw className="h-4 w-4" />
                            ) : isConnected ? (
                                <WifiOff className="h-4 w-4" />
                            ) : (
                                <Wifi className="h-4 w-4" />
                            )}
                        </button>
                    </div>
                </div>

                {(error || sseError) && (
                    <div className="mt-4 p-3 bg-red-50 border border-red-200 rounded-lg text-red-700 text-sm">
                        <div className="flex items-center gap-2">
                            <X className="h-4 w-4" />
                            <span>{error || sseError}</span>
                        </div>
                    </div>
                )}
            </div>

            {/* Timeout Countdown Timer */}
            {isConnected && timeUntilDisconnect !== null && (
                <div className="text-xs text-gray-500 text-center flex items-center justify-center gap-1">
                    <Clock className="h-3 w-3" />
                    <span>
                        {timeUntilDisconnect > 0
                            ? `Auto-disconnect in ${timeUntilDisconnect}s`
                            : 'Disconnecting...'}
                    </span>
                </div>
            )}

            {/* Inactivity Warning Banner */}
            {isIdle && !needsManualRefresh && (
                <Card className="bg-yellow-50 border-yellow-200 border">
                    <CardContent className="pt-6">
                        <div className="flex items-center justify-between">
                            <div className="flex items-center gap-3">
                                <Clock className="h-5 w-5 text-yellow-600" />
                                <div>
                                    <h3 className="font-medium text-yellow-800">Connection Idle</h3>
                                    <p className="text-sm text-yellow-700">
                                        No recent ball events. Connection will pause automatically to save resources.
                                    </p>
                                </div>
                            </div>
                            <Button
                                variant="outline"
                                size="sm"
                                onClick={manualReconnect}
                                className="border-yellow-300 text-yellow-700 hover:bg-yellow-100"
                            >
                                <RefreshCw className="h-4 w-4 mr-1" />
                                Reconnect Now
                            </Button>
                        </div>
                    </CardContent>
                </Card>
            )}

            {/* Manual Refresh Prompt */}
            {needsManualRefresh && (
                <Card className="bg-orange-50 border-orange-200 border">
                    <CardContent className="pt-6">
                        <div className="flex items-center justify-between">
                            <div className="flex items-center gap-3">
                                <Zap className="h-5 w-5 text-orange-600" />
                                <div>
                                    <h3 className="font-medium text-orange-800">Connection Paused</h3>
                                    <p className="text-sm text-orange-700">
                                        Connection paused to save resources. {disconnectReason === 'inactivity' ? 'No recent activity detected.' : 'Click to reconnect.'}
                                    </p>
                                    {lastEventTime && (
                                        <p className="text-xs text-orange-600 mt-1">
                                            Last update: {lastEventTime.toLocaleTimeString()}
                                        </p>
                                    )}
                                </div>
                            </div>
                            <Button
                                onClick={manualReconnect}
                                className="bg-orange-600 hover:bg-orange-700 text-white"
                            >
                                <RefreshCw className="h-4 w-4 mr-1" />
                                Reconnect
                            </Button>
                        </div>
                    </CardContent>
                </Card>
            )}

            {/* Live Events List - Compact */}
            {events.length > 0 && (
                <Card className="bg-white shadow-sm border-gray-200">
                    <CardHeader className="pb-1 pt-2 px-2">
                        <div className="flex items-center justify-between">
                            <CardTitle className="text-sm font-semibold flex items-center gap-2">
                                <Activity className="h-4 w-4 text-blue-600" />
                                Live Events
                            </CardTitle>
                            <div className="flex items-center gap-1">
                                <button
                                    onClick={fetchPreviousEvents}
                                    disabled={isLoadingPrevious}
                                    className="p-1.5 rounded-full bg-blue-50 hover:bg-blue-100 border border-blue-200 text-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                                    title="Load Previous Events"
                                >
                                    <Download className={`h-3 w-3 ${isLoadingPrevious ? 'animate-pulse' : ''}`} />
                                </button>

                                <button
                                    onClick={needsManualRefresh ? manualReconnect : (isConnected ? disconnect : connect)}
                                    disabled={isConnecting}
                                    className={`p-1.5 rounded-full transition-colors disabled:opacity-50 disabled:cursor-not-allowed ${isConnected
                                        ? 'bg-red-50 hover:bg-red-100 border border-red-200 text-red-700'
                                        : 'bg-green-50 hover:bg-green-100 border border-green-200 text-green-700'
                                        }`}
                                    title={needsManualRefresh ? 'Reconnect to Live Events' : isConnected ? 'Disconnect from Live Events' : 'Connect to Live Events'}
                                >
                                    {needsManualRefresh ? (
                                        <RefreshCw className={`h-3 w-3 ${isConnecting ? 'animate-spin' : ''}`} />
                                    ) : isConnected ? (
                                        <WifiOff className="h-3 w-3" />
                                    ) : (
                                        <Wifi className="h-3 w-3" />
                                    )}
                                </button>

                                <Button
                                    variant="outline"
                                    size="sm"
                                    onClick={() => dispatch(clearMatchEvents(matchId))}
                                    className="bg-red-50 hover:bg-red-100 border-red-200 text-red-700 h-7 px-2"
                                    title="Clear All Events"
                                >
                                    <X className="h-3 w-3 mr-1" />
                                    <span className="text-xs">Clear</span>
                                </Button>
                            </div>
                        </div>
                    </CardHeader>
                    <CardContent className="p-1 px-2 pb-2">
                        <ScrollArea className="h-48 sm:h-64 w-full">
                            {/* Header */}
                            <div className="bg-gray-100 p-1 border-b border-gray-300 sticky top-0">
                                <div className="flex items-center justify-between text-xs font-semibold text-gray-700">
                                    <div className="w-20 text-center">Runs</div>
                                    <div className="w-16 text-center">Over</div>
                                    <div className="w-20 text-center">Score</div>
                                </div>
                            </div>

                            <div className="space-y-0.5">
                                {[...events].reverse().map((event, index) => {
                                    const isPreviousEvent = event.stream_id.startsWith('scorecard-');
                                    const isLatestEvent = index === 0; // First item in reversed array is latest
                                    return (
                                        <div key={`${event.stream_id}-${index}`} className={`p-1.5 border-l-4 transition-all duration-200 relative ${isLatestEvent
                                            ? 'bg-green-50 border-green-400 shadow-sm ring-1 ring-green-200'
                                            : isPreviousEvent
                                                ? 'bg-blue-50 border-blue-400'
                                                : 'bg-gray-50 border-gray-300'
                                            }`}>
                                            {/* Latest Event Badge */}
                                            {isLatestEvent && (
                                                <div className="absolute -top-1 -right-1">
                                                    <div className="bg-green-500 text-white text-xs px-2 py-0.5 rounded-full font-bold animate-pulse">
                                                        LATEST
                                                    </div>
                                                </div>
                                            )}
                                            {/* Compact horizontal layout */}
                                            <div className="flex items-center justify-between">
                                                {/* Latest Event Indicator */}
                                                {isLatestEvent && (
                                                    <div className="absolute -left-1 top-1/2 transform -translate-y-1/2">
                                                        <div className="w-2 h-2 bg-green-500 rounded-full animate-pulse"></div>
                                                    </div>
                                                )}

                                                {/* Runs - Circular ball */}
                                                <div className="w-20 flex justify-center relative">
                                                    <div className={`w-8 h-8 rounded-full border-2 flex flex-col items-center justify-center text-xs font-medium ${isLatestEvent
                                                        ? 'bg-green-100 border-green-400 shadow-md'
                                                        : 'bg-gray-100 border-gray-300'
                                                        }`}>
                                                        {(() => {
                                                            // Follow the same logic as scorecard ball display
                                                            if (event.is_wicket) {
                                                                return <span className="text-[11px] leading-none font-bold">W</span>;
                                                            }

                                                            // Handle byes with different ball types
                                                            if (event.byes > 0) {
                                                                const byeText = `B${event.byes}`;

                                                                if (event.ball_type === 'NO_BALL' || event.ball_type === 'no_ball') {
                                                                    return (
                                                                        <div className="text-[9px] leading-none text-gray-700 font-bold">
                                                                            <div>{byeText}</div>
                                                                            <div className="text-[6px] leading-none text-gray-600">+</div>
                                                                            <div className="text-[9px] leading-none text-gray-700 font-bold">nb</div>
                                                                        </div>
                                                                    );
                                                                }
                                                                if (event.ball_type === 'WIDE' || event.ball_type === 'wide') {
                                                                    return (
                                                                        <div className="text-[9px] leading-none text-gray-700 font-bold">
                                                                            <div>{byeText}</div>
                                                                            <div className="text-[6px] leading-none text-gray-600">+</div>
                                                                            <div className="text-[9px] leading-none text-gray-700 font-bold">wd</div>
                                                                        </div>
                                                                    );
                                                                }
                                                                if (event.run_type === 'LB') {
                                                                    return (
                                                                        <div className="text-[9px] leading-none text-gray-700 font-bold">
                                                                            <div>{byeText}</div>
                                                                            <div className="text-[6px] leading-none text-gray-600">+</div>
                                                                            <div className="text-[9px] leading-none text-gray-700 font-bold">lb</div>
                                                                        </div>
                                                                    );
                                                                }
                                                                // Regular byes on good ball
                                                                return <span className="text-[11px] leading-none text-gray-700 font-bold">{byeText}</span>;
                                                            }

                                                            // Handle other ball types without byes
                                                            switch (event.ball_type) {
                                                                case 'WIDE':
                                                                case 'wide':
                                                                    return <span className="text-[10px] leading-none font-bold">Wd</span>;
                                                                case 'NO_BALL':
                                                                case 'no_ball':
                                                                    return <span className="text-[10px] leading-none font-bold">nb</span>;
                                                                case 'DEAD_BALL':
                                                                case 'dead_ball':
                                                                    return <span className="text-[10px] leading-none font-bold">Db</span>;
                                                                default:
                                                                    // Handle run types
                                                                    switch (event.run_type) {
                                                                        case 'LB':
                                                                            return (
                                                                                <div className="text-[10px] leading-none font-bold">
                                                                                    <div>Lb</div>
                                                                                    <div className="text-[8px] leading-none text-gray-600">+</div>
                                                                                    <div className="text-[10px] leading-none font-bold">{event.runs || 0}</div>
                                                                                </div>
                                                                            );
                                                                        case 'WC':
                                                                            return <span className="text-[10px] leading-none font-bold">W</span>;
                                                                        default:
                                                                            return <span className="text-[11px] leading-none font-bold">{event.runs?.toString() || '0'}</span>;
                                                                    }
                                                            }
                                                        })()}
                                                    </div>
                                                </div>

                                                {/* Over */}
                                                <div className="w-16 text-center">
                                                    <span className={`text-sm font-bold ${isLatestEvent ? 'text-green-700' : 'text-blue-700'
                                                        }`}>
                                                        {event.innings_overs}
                                                    </span>
                                                </div>

                                                {/* Score */}
                                                <div className="w-20 text-center">
                                                    <span className={`text-sm font-bold ${isLatestEvent ? 'text-green-700' : 'text-purple-700'
                                                        }`}>
                                                        {event.innings_runs}/{event.innings_wickets}
                                                    </span>
                                                </div>
                                            </div>
                                        </div>
                                    );
                                })}
                            </div>
                        </ScrollArea>
                    </CardContent>
                </Card>
            )}

            {/* No Events Message - Mobile First */}
            {events.length === 0 && (
                <Card className="bg-gradient-to-br from-gray-50 to-slate-50 border-gray-200">
                    <CardContent className="py-6 sm:py-8 text-center">
                        <div className="flex flex-col items-center">
                            <div className="relative">
                                <Activity className="h-10 w-10 sm:h-12 sm:w-12 text-gray-300 mb-3 animate-pulse" />
                                {isConnected && (
                                    <div className="absolute -top-1 -right-1">
                                        <div className="w-3 h-3 bg-green-400 rounded-full animate-ping"></div>
                                    </div>
                                )}
                            </div>
                            <h3 className="text-base sm:text-lg font-semibold text-gray-700 mb-2">
                                {isConnected ? 'Waiting for Ball Events' : 'Connect to See Live Events'}
                            </h3>
                            <p className="text-gray-500 max-w-md text-sm px-4">
                                {isConnected
                                    ? 'Live ball events will appear here automatically when balls are added to the match.'
                                    : 'Click the Connect button above to start receiving live ball events.'
                                }
                            </p>
                            {isConnected && (
                                <div className="mt-3 flex items-center gap-2 text-sm text-green-600">
                                    <Wifi className="h-4 w-4" />
                                    <span className="hidden sm:inline">Live connection established</span>
                                    <span className="sm:hidden">Connected</span>
                                </div>
                            )}
                        </div>
                    </CardContent>
                </Card>
            )}
        </div>
    );
}
