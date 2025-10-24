'use client';

import { useState, useEffect } from 'react';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import { fetchTimeTrackingRequest, clearTimeTracking } from '@/store/reducers/timeTrackingSlice';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import {
    Clock,
    Timer,
    Play,
    Pause,
    CheckCircle,
    ChevronDown,
    ChevronUp,
    RefreshCw
} from 'lucide-react';
import {
    TimeTrackingResponse,
    InningsTimeTracking,
    OverTimeTracking,
    formatDuration,
    formatTime,
    formatDateTime
} from '@/types/timeTracking';

interface TimeTrackingViewProps {
    matchId: string;
    onBack: () => void;
}

export function TimeTrackingView({ matchId, onBack }: TimeTrackingViewProps): React.JSX.Element {
    const dispatch = useAppDispatch();
    const { timeTracking, loading, error } = useAppSelector(state => state.timeTracking);

    const [expandedInnings, setExpandedInnings] = useState<{ [key: number]: boolean }>({});
    const [expandedOvers, setExpandedOvers] = useState<{ [key: string]: boolean }>({});

    useEffect(() => {
        dispatch(fetchTimeTrackingRequest(matchId));
        return () => {
            dispatch(clearTimeTracking());
        };
    }, [dispatch, matchId]);

    const toggleInningsExpansion = (inningsNumber: number) => {
        setExpandedInnings(prev => ({
            ...prev,
            [inningsNumber]: !prev[inningsNumber]
        }));
    };

    const toggleOverExpansion = (inningsNumber: number, overNumber: number) => {
        const key = `${inningsNumber}-${overNumber}`;
        setExpandedOvers(prev => ({
            ...prev,
            [key]: !prev[key]
        }));
    };

    const getStatusColor = (status: string) => {
        switch (status) {
            case 'completed':
                return 'bg-green-100 text-green-800 border-green-200';
            case 'in_progress':
                return 'bg-blue-100 text-blue-800 border-blue-200';
            default:
                return 'bg-gray-100 text-gray-800 border-gray-200';
        }
    };

    const getStatusIcon = (status: string) => {
        switch (status) {
            case 'completed':
                return <CheckCircle className="h-4 w-4" />;
            case 'in_progress':
                return <Play className="h-4 w-4" />;
            default:
                return <Pause className="h-4 w-4" />;
        }
    };

    if (loading) {
        return (
            <div className="min-h-screen bg-gray-50 p-4">
                <div className="max-w-4xl mx-auto">
                    <div className="flex items-center justify-center h-64">
                        <div className="flex items-center space-x-2">
                            <RefreshCw className="h-6 w-6 animate-spin text-blue-600" />
                            <span className="text-lg text-gray-600">Loading time tracking data...</span>
                        </div>
                    </div>
                </div>
            </div>
        );
    }

    if (error) {
        return (
            <div className="min-h-screen bg-gray-50 p-4">
                <div className="max-w-4xl mx-auto">
                    <Card className="border-red-200 bg-red-50">
                        <CardContent className="p-6">
                            <div className="text-center">
                                <div className="text-red-600 mb-2">
                                    <Clock className="h-12 w-12 mx-auto" />
                                </div>
                                <h3 className="text-lg font-semibold text-red-800 mb-2">Error Loading Time Tracking</h3>
                                <p className="text-red-600 mb-4">{error}</p>
                                <Button
                                    onClick={() => dispatch(fetchTimeTrackingRequest(matchId))}
                                    variant="outline"
                                    className="border-red-300 text-red-700 hover:bg-red-100"
                                >
                                    <RefreshCw className="h-4 w-4 mr-2" />
                                    Retry
                                </Button>
                            </div>
                        </CardContent>
                    </Card>
                </div>
            </div>
        );
    }

    if (!timeTracking || !timeTracking.innings) {
        return (
            <div className="min-h-screen bg-gray-50 p-4">
                <div className="max-w-4xl mx-auto">
                    <Card>
                        <CardContent className="p-6">
                            <div className="text-center">
                                <Clock className="h-12 w-12 mx-auto text-gray-400 mb-4" />
                                <h3 className="text-lg font-semibold text-gray-800 mb-2">No Time Tracking Data</h3>
                                <p className="text-gray-600">Time tracking data is not available for this match.</p>
                            </div>
                        </CardContent>
                    </Card>
                </div>
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-gray-50 p-4">
            <div className="max-w-4xl mx-auto space-y-6">
                {/* Header */}
                <div className="flex items-center justify-between">
                    <div className="flex items-center space-x-3">
                        <Button
                            variant="ghost"
                            size="sm"
                            onClick={onBack}
                            className="h-8 px-2"
                        >
                            ← Back
                        </Button>
                        <div className="flex items-center space-x-2">
                            <Clock className="h-6 w-6 text-blue-600" />
                            <h1 className="text-2xl font-bold text-gray-900">Time Tracking</h1>
                        </div>
                    </div>
                    <Button
                        variant="outline"
                        size="sm"
                        onClick={() => dispatch(fetchTimeTrackingRequest(matchId))}
                        className="h-8 px-3"
                    >
                        <RefreshCw className="h-4 w-4 mr-2" />
                        Refresh
                    </Button>
                </div>

                {/* Match Summary */}
                <Card>
                    <CardHeader>
                        <CardTitle className="flex items-center space-x-2">
                            <Timer className="h-5 w-5 text-blue-600" />
                            <span>Match Summary</span>
                        </CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                            <div className="text-center p-4 bg-blue-50 rounded-lg">
                                <div className="text-2xl font-bold text-blue-600">
                                    {formatDuration(timeTracking.total_match_time_seconds)}
                                </div>
                                <div className="text-sm text-blue-800">Total Match Time</div>
                            </div>
                            <div className="text-center p-4 bg-green-50 rounded-lg">
                                <div className="text-2xl font-bold text-green-600">
                                    {timeTracking.innings?.length || 0}
                                </div>
                                <div className="text-sm text-green-800">Innings</div>
                            </div>
                            <div className="text-center p-4 bg-purple-50 rounded-lg">
                                <div className="text-2xl font-bold text-purple-600">
                                    {timeTracking.innings?.reduce((total, innings) => total + (innings.overs?.length || 0), 0) || 0}
                                </div>
                                <div className="text-sm text-purple-800">Total Overs</div>
                            </div>
                        </div>
                    </CardContent>
                </Card>

                {/* Innings Details */}
                {timeTracking.innings?.map((innings: InningsTimeTracking) => (
                    <Card key={innings.innings_number}>
                        <CardHeader>
                            <div className="flex items-center justify-between">
                                <CardTitle className="flex items-center space-x-2">
                                    <div className="flex items-center space-x-2">
                                        {getStatusIcon(innings.status)}
                                        <span>Innings {innings.innings_number}</span>
                                        <Badge className={getStatusColor(innings.status)}>
                                            Team {innings.batting_team}
                                        </Badge>
                                    </div>
                                </CardTitle>
                                <Button
                                    variant="ghost"
                                    size="sm"
                                    onClick={() => toggleInningsExpansion(innings.innings_number)}
                                >
                                    {expandedInnings[innings.innings_number] ? (
                                        <ChevronUp className="h-4 w-4" />
                                    ) : (
                                        <ChevronDown className="h-4 w-4" />
                                    )}
                                </Button>
                            </div>
                        </CardHeader>
                        <CardContent>
                            <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-4">
                                <div className="text-center p-3 bg-gray-50 rounded-lg">
                                    <div className="text-lg font-semibold text-gray-800">
                                        {formatDuration(innings.duration_seconds)}
                                    </div>
                                    <div className="text-sm text-gray-600">Duration</div>
                                </div>
                                <div className="text-center p-3 bg-gray-50 rounded-lg">
                                    <div className="text-lg font-semibold text-gray-800">
                                        {formatTime(innings.start_time)}
                                    </div>
                                    <div className="text-sm text-gray-600">Start Time</div>
                                </div>
                                <div className="text-center p-3 bg-gray-50 rounded-lg">
                                    <div className="text-lg font-semibold text-gray-800">
                                        {formatTime(innings.end_time)}
                                    </div>
                                    <div className="text-sm text-gray-600">End Time</div>
                                </div>
                                <div className="text-center p-3 bg-gray-50 rounded-lg">
                                    <div className="text-lg font-semibold text-gray-800">
                                        {innings.overs?.length || 0}
                                    </div>
                                    <div className="text-sm text-gray-600">Overs</div>
                                </div>
                            </div>

                            {expandedInnings[innings.innings_number] && (
                                <div className="space-y-3">
                                    <h4 className="font-semibold text-gray-800 mb-3">Overs Breakdown</h4>
                                    {innings.overs?.map((over: OverTimeTracking) => (
                                        <Card key={`${innings.innings_number}-${over.over_number}`} className="bg-gray-50">
                                            <CardHeader className="pb-2">
                                                <div className="flex items-center justify-between">
                                                    <div className="flex items-center space-x-2">
                                                        <span className="font-medium">Over {over.over_number}</span>
                                                        <Badge className={getStatusColor(over.status)}>
                                                            {over.status}
                                                        </Badge>
                                                    </div>
                                                    <Button
                                                        variant="ghost"
                                                        size="sm"
                                                        onClick={() => toggleOverExpansion(innings.innings_number, over.over_number)}
                                                    >
                                                        {expandedOvers[`${innings.innings_number}-${over.over_number}`] ? (
                                                            <ChevronUp className="h-4 w-4" />
                                                        ) : (
                                                            <ChevronDown className="h-4 w-4" />
                                                        )}
                                                    </Button>
                                                </div>
                                            </CardHeader>
                                            <CardContent>
                                                <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
                                                    <div className="text-center">
                                                        <div className="text-sm font-semibold text-gray-800">
                                                            {formatDuration(over.duration_seconds)}
                                                        </div>
                                                        <div className="text-xs text-gray-600">Duration</div>
                                                    </div>
                                                    <div className="text-center">
                                                        <div className="text-sm font-semibold text-gray-800">
                                                            {over.total_runs}
                                                        </div>
                                                        <div className="text-xs text-gray-600">Runs</div>
                                                    </div>
                                                    <div className="text-center">
                                                        <div className="text-sm font-semibold text-gray-800">
                                                            {over.total_balls}
                                                        </div>
                                                        <div className="text-xs text-gray-600">Balls</div>
                                                    </div>
                                                    <div className="text-center">
                                                        <div className="text-sm font-semibold text-gray-800">
                                                            {over.total_wickets}
                                                        </div>
                                                        <div className="text-xs text-gray-600">Wickets</div>
                                                    </div>
                                                </div>

                                                {expandedOvers[`${innings.innings_number}-${over.over_number}`] && (
                                                    <div className="mt-3 pt-3 border-t border-gray-200">
                                                        <div className="grid grid-cols-1 md:grid-cols-2 gap-3 text-sm">
                                                            <div>
                                                                <span className="font-medium text-gray-700">Start Time:</span>
                                                                <span className="ml-2 text-gray-600">
                                                                    {formatDateTime(over.start_time)}
                                                                </span>
                                                            </div>
                                                            <div>
                                                                <span className="font-medium text-gray-700">End Time:</span>
                                                                <span className="ml-2 text-gray-600">
                                                                    {formatDateTime(over.end_time)}
                                                                </span>
                                                            </div>
                                                        </div>
                                                    </div>
                                                )}
                                            </CardContent>
                                        </Card>
                                    ))}
                                </div>
                            )}
                        </CardContent>
                    </Card>
                ))}
            </div>
        </div>
    );
}


