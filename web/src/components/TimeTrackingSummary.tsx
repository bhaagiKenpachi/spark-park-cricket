'use client';

import { useEffect } from 'react';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import { fetchTimeTrackingRequest } from '@/store/reducers/timeTrackingSlice';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Clock, Timer, Play, CheckCircle } from 'lucide-react';
import { formatDuration, formatTime } from '@/types/timeTracking';

interface TimeTrackingSummaryProps {
    matchId: string;
}

export function TimeTrackingSummary({ matchId }: TimeTrackingSummaryProps): React.JSX.Element {
    const dispatch = useAppDispatch();
    const { timeTracking, loading, error } = useAppSelector(state => state.timeTracking);

    useEffect(() => {
        dispatch(fetchTimeTrackingRequest(matchId));
    }, [dispatch, matchId]);

    if (loading) {
        return (
            <Card>
                <CardContent className="p-4">
                    <div className="flex items-center justify-center">
                        <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-blue-600"></div>
                        <span className="ml-2 text-sm text-gray-600">Loading time data...</span>
                    </div>
                </CardContent>
            </Card>
        );
    }

    if (error || !timeTracking) {
        return (
            <Card>
                <CardContent className="p-4">
                    <div className="text-center text-gray-500">
                        <Clock className="h-8 w-8 mx-auto mb-2 text-gray-400" />
                        <p className="text-sm">Time tracking data not available</p>
                    </div>
                </CardContent>
            </Card>
        );
    }

    const completedInnings = timeTracking.innings?.filter(innings => innings.status === 'completed') || [];
    const inProgressInnings = timeTracking.innings?.filter(innings => innings.status === 'in_progress') || [];

    return (
        <Card>
            <CardHeader className="pb-3">
                <CardTitle className="flex items-center space-x-2 text-lg">
                    <Timer className="h-5 w-5 text-blue-600" />
                    <span>Match Timing</span>
                </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
                {/* Total Match Time */}
                <div className="text-center p-4 bg-blue-50 rounded-lg">
                    <div className="text-2xl font-bold text-blue-600">
                        {formatDuration(timeTracking.total_match_time_seconds)}
                    </div>
                    <div className="text-sm text-blue-800">Total Match Time</div>
                </div>

                {/* Innings Summary */}
                <div className="space-y-3">
                    {timeTracking.innings?.map((innings, index) => (
                        <div key={innings.innings_number} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                            <div className="flex items-center space-x-3">
                                <div className="flex items-center space-x-2">
                                    {innings.status === 'completed' ? (
                                        <CheckCircle className="h-4 w-4 text-green-600" />
                                    ) : (
                                        <Play className="h-4 w-4 text-blue-600" />
                                    )}
                                    <span className="font-medium">Innings {innings.innings_number}</span>
                                    <Badge variant="outline" className="text-xs">
                                        Team {innings.batting_team}
                                    </Badge>
                                </div>
                            </div>
                            <div className="flex items-center space-x-4 text-sm">
                                <div className="text-center">
                                    <div className="font-semibold text-gray-800">
                                        {formatDuration(innings.duration_seconds)}
                                    </div>
                                    <div className="text-xs text-gray-600">Duration</div>
                                </div>
                                <div className="text-center">
                                    <div className="font-semibold text-gray-800">
                                        {innings.overs?.length || 0}
                                    </div>
                                    <div className="text-xs text-gray-600">Overs</div>
                                </div>
                                <div className="text-center">
                                    <div className="font-semibold text-gray-800">
                                        {formatTime(innings.start_time)}
                                    </div>
                                    <div className="text-xs text-gray-600">Started</div>
                                </div>
                            </div>
                        </div>
                    ))}
                </div>

                {/* Quick Stats */}
                <div className="grid grid-cols-3 gap-3 pt-3 border-t">
                    <div className="text-center">
                        <div className="text-lg font-bold text-gray-800">
                            {completedInnings.length}
                        </div>
                        <div className="text-xs text-gray-600">Completed</div>
                    </div>
                    <div className="text-center">
                        <div className="text-lg font-bold text-gray-800">
                            {inProgressInnings.length}
                        </div>
                        <div className="text-xs text-gray-600">In Progress</div>
                    </div>
                    <div className="text-center">
                        <div className="text-lg font-bold text-gray-800">
                            {timeTracking.innings?.reduce((total, innings) => total + (innings.overs?.length || 0), 0) || 0}
                        </div>
                        <div className="text-xs text-gray-600">Total Overs</div>
                    </div>
                </div>
            </CardContent>
        </Card>
    );
}


