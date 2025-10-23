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
    X
} from 'lucide-react';
import { useSSE, BallEvent } from '@/hooks/useSSE';
import { ApiService } from '@/services/api';
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
    const [isExpanded, setIsExpanded] = useState(false);
    const [showRawData, setShowRawData] = useState(false);

    // Use Redux selectors for events
    const events = useSelector((state: RootState) => selectEventsForMatch(state, matchId));
    const isLoadingPrevious = useSelector((state: RootState) => selectIsLoadingEventsForMatch(state, matchId));
    const error = useSelector((state: RootState) => selectErrorForMatch(state, matchId));

    // Function to convert scorecard ball data to BallEvent format
    const convertScorecardBallToEvent = (ball: any, overNumber: number, inningsNumber: number): BallEvent => {
        return {
            event_type: 'ball_added',
            match_id: matchId,
            innings_number: inningsNumber,
            ball_number: ball.ball_number,
            ball_type: ball.ball_type,
            run_type: ball.run_type,
            runs: ball.runs,
            byes: ball.byes,
            total_runs: ball.runs + ball.byes,
            is_wicket: ball.is_wicket,
            wicket_type: ball.wicket_type || '',
            innings_runs: 0, // Will be calculated from scorecard
            innings_wickets: 0, // Will be calculated from scorecard
            innings_overs: `${overNumber}.${ball.ball_number}`, // Approximate
            timestamp: new Date().toISOString(), // Use current time as we don't have exact timestamp
            stream_id: `scorecard-${inningsNumber}-${overNumber}-${ball.ball_number}`
        };
    };

    // Function to fetch previous ball events from scorecard
    const fetchPreviousEvents = async () => {
        if (isLoadingPrevious) return;

        dispatch(setMatchLoading({ matchId, loading: true }));
        dispatch(setMatchError({ matchId, error: null }));

        try {
            const apiService = new ApiService();
            const response = await apiService.getScorecard(matchId);
            const scorecard = response.data.data; // Extract the actual scorecard data

            const previousEvents: ReduxBallEvent[] = [];

            // Extract ball events from all innings
            if (scorecard.innings && Array.isArray(scorecard.innings)) {
                scorecard.innings.forEach((innings: any) => {
                    if (innings.overs && Array.isArray(innings.overs)) {
                        innings.overs.forEach((over: any) => {
                            if (over.balls && Array.isArray(over.balls)) {
                                over.balls.forEach((ball: any) => {
                                    const event = convertScorecardBallToEvent(ball, over.over_number, innings.innings_number);
                                    // Update innings totals
                                    event.innings_runs = innings.total_runs;
                                    event.innings_wickets = innings.total_wickets;
                                    event.innings_overs = innings.total_overs.toString();
                                    previousEvents.push(event);
                                });
                            }
                        });
                    }
                });
            }

            // Sort by innings number, then by over number, then by ball number
            previousEvents.sort((a, b) => {
                if (a.innings_number !== b.innings_number) {
                    return a.innings_number - b.innings_number;
                }
                const aOver = parseFloat(a.innings_overs);
                const bOver = parseFloat(b.innings_overs);
                if (aOver !== bOver) {
                    return aOver - bOver;
                }
                return a.ball_number - b.ball_number;
            });

            // Dispatch the previous events to Redux store
            dispatch(addPreviousEvents({ matchId, events: previousEvents }));

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

    // Initialize match events when component mounts
    useEffect(() => {
        if (matchId) {
            dispatch(initializeMatchEvents(matchId));
            fetchPreviousEvents();
        }
    }, [matchId, dispatch]);

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
            {/* Connection Status Card */}
            <Card className={`${status.bgColor} border`}>
                <CardHeader className="pb-3">
                    <div className="flex items-center justify-between">
                        <CardTitle className="flex items-center gap-2 text-lg">
                            <Activity className="h-5 w-5" />
                            Live Updates
                        </CardTitle>
                        <div className="flex items-center gap-2">
                            <Badge variant="outline" className={`${status.color} ${status.bgColor}`}>
                                {status.icon}
                                <span className="ml-1">{status.text}</span>
                            </Badge>
                            {eventCount > 0 && (
                                <Badge variant="secondary">
                                    {eventCount} events
                                </Badge>
                            )}
                        </div>
                    </div>
                </CardHeader>
                <CardContent>
                    <div className="flex items-center justify-between">
                        <div className="flex items-center gap-4">
                            <div className="flex items-center gap-2">
                                <Clock className="h-4 w-4 text-gray-500" />
                                <span className="text-sm text-gray-600">
                                    Match ID: {matchId}
                                </span>
                            </div>
                            {lastEvent && (
                                <div className="flex items-center gap-2">
                                    <Target className="h-4 w-4 text-gray-500" />
                                    <span className="text-sm text-gray-600">
                                        Last: {formatTimestamp(lastEvent.timestamp)}
                                    </span>
                                </div>
                            )}
                        </div>
                        <div className="flex items-center gap-2">
                            <Button
                                variant="outline"
                                size="sm"
                                onClick={fetchPreviousEvents}
                                disabled={isLoadingPrevious}
                            >
                                <RefreshCw className={`h-4 w-4 mr-1 ${isLoadingPrevious ? 'animate-spin' : ''}`} />
                                {isLoadingPrevious ? 'Loading...' : 'Load Previous'}
                            </Button>
                            <Button
                                variant="outline"
                                size="sm"
                                onClick={() => setShowRawData(!showRawData)}
                            >
                                <Zap className="h-4 w-4 mr-1" />
                                {showRawData ? 'Hide' : 'Show'} Raw
                            </Button>
                            <Button
                                variant="outline"
                                size="sm"
                                onClick={needsManualRefresh ? manualReconnect : (isConnected ? disconnect : connect)}
                                disabled={isConnecting}
                            >
                                {needsManualRefresh ? (
                                    <>
                                        <RefreshCw className="h-4 w-4 mr-1" />
                                        Reconnect
                                    </>
                                ) : isConnected ? (
                                    <>
                                        <WifiOff className="h-4 w-4 mr-1" />
                                        Disconnect
                                    </>
                                ) : (
                                    <>
                                        <Wifi className="h-4 w-4 mr-1" />
                                        Connect
                                    </>
                                )}
                            </Button>
                        </div>
                    </div>
                    {(error || sseError) && (
                        <div className="mt-3 p-2 bg-red-50 border border-red-200 rounded text-red-700 text-sm">
                            {error || sseError}
                        </div>
                    )}
                </CardContent>
            </Card>

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

            {/* Events List */}
            {events.length > 0 && (
                <Card>
                    <CardHeader className="pb-3">
                        <div className="flex items-center justify-between">
                            <CardTitle className="text-lg">Recent Events</CardTitle>
                            <div className="flex items-center gap-2">
                                <Button
                                    variant="outline"
                                    size="sm"
                                    onClick={() => setIsExpanded(!isExpanded)}
                                >
                                    {isExpanded ? 'Collapse' : 'Expand'}
                                </Button>
                                <Button
                                    variant="outline"
                                    size="sm"
                                    onClick={() => dispatch(clearMatchEvents(matchId))}
                                >
                                    <X className="h-4 w-4 mr-1" />
                                    Clear
                                </Button>
                            </div>
                        </div>
                    </CardHeader>
                    <CardContent>
                        <ScrollArea className={`${isExpanded ? 'h-96' : 'h-48'} w-full`}>
                            <div className="space-y-2">
                                {events.map((event, index) => {
                                    const isPreviousEvent = event.stream_id.startsWith('scorecard-');
                                    return (
                                        <div key={`${event.stream_id}-${index}`} className={`p-3 rounded-lg ${isPreviousEvent ? 'bg-blue-50 border-l-4 border-blue-300' : 'bg-gray-50'}`}>
                                            <div className="flex items-center justify-between mb-2">
                                                <div className="flex items-center gap-2">
                                                    {isPreviousEvent && (
                                                        <Badge variant="outline" className="text-blue-600 border-blue-300">
                                                            Previous
                                                        </Badge>
                                                    )}
                                                    <Badge className={getBallTypeColor(event.ball_type)}>
                                                        {event.ball_type.toUpperCase()}
                                                    </Badge>
                                                    <Badge className={getRunTypeColor(event.run_type)}>
                                                        {event.run_type.toUpperCase()}
                                                    </Badge>
                                                    {event.is_wicket && (
                                                        <Badge variant="destructive">
                                                            WICKET
                                                        </Badge>
                                                    )}
                                                </div>
                                                <span className="text-xs text-gray-500">
                                                    {formatTimestamp(event.timestamp)}
                                                </span>
                                            </div>

                                            <div className="grid grid-cols-2 md:grid-cols-4 gap-2 text-sm">
                                                <div>
                                                    <span className="text-gray-600">Ball:</span>
                                                    <span className="ml-1 font-medium">{event.ball_number}</span>
                                                </div>
                                                <div>
                                                    <span className="text-gray-600">Runs:</span>
                                                    <span className="ml-1 font-medium">{event.runs}</span>
                                                </div>
                                                <div>
                                                    <span className="text-gray-600">Byes:</span>
                                                    <span className="ml-1 font-medium">{event.byes}</span>
                                                </div>
                                                <div>
                                                    <span className="text-gray-600">Total:</span>
                                                    <span className="ml-1 font-medium">{event.total_runs}</span>
                                                </div>
                                            </div>

                                            <div className="mt-2 pt-2 border-t border-gray-200">
                                                <div className="grid grid-cols-2 md:grid-cols-3 gap-2 text-sm">
                                                    <div>
                                                        <span className="text-gray-600">Innings:</span>
                                                        <span className="ml-1 font-medium">{event.innings_number}</span>
                                                    </div>
                                                    <div>
                                                        <span className="text-gray-600">Score:</span>
                                                        <span className="ml-1 font-medium">
                                                            {event.innings_runs}/{event.innings_wickets}
                                                        </span>
                                                    </div>
                                                    <div>
                                                        <span className="text-gray-600">Overs:</span>
                                                        <span className="ml-1 font-medium">{event.innings_overs}</span>
                                                    </div>
                                                </div>
                                            </div>

                                            {showRawData && (
                                                <>
                                                    <Separator className="my-2" />
                                                    <details className="text-xs">
                                                        <summary className="cursor-pointer text-gray-600 hover:text-gray-800">
                                                            Raw Event Data
                                                        </summary>
                                                        <pre className="mt-2 p-2 bg-gray-100 rounded text-xs overflow-x-auto">
                                                            {JSON.stringify(event, null, 2)}
                                                        </pre>
                                                    </details>
                                                </>
                                            )}
                                        </div>
                                    );
                                })}
                            </div>
                        </ScrollArea>
                    </CardContent>
                </Card>
            )}

            {/* No Events Message */}
            {events.length === 0 && isConnected && (
                <Card>
                    <CardContent className="py-8 text-center">
                        <Activity className="h-12 w-12 text-gray-400 mx-auto mb-4" />
                        <p className="text-gray-600">Waiting for ball events...</p>
                        <p className="text-sm text-gray-500 mt-1">
                            Events will appear here when balls are added to the match
                        </p>
                    </CardContent>
                </Card>
            )}
        </div>
    );
}
