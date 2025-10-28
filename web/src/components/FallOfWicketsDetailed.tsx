'use client';

import React, { useState, useEffect } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { RefreshCw } from 'lucide-react';
import { FallOfWickets, FallOfWicketsSummary } from '@/types/fallOfWickets';
import { apiService } from '@/services/api';

interface FallOfWicketsDetailedProps {
    matchId: string;
    inningsId?: string;
    className?: string;
}

export const FallOfWicketsDetailed: React.FC<FallOfWicketsDetailedProps> = ({
    matchId,
    inningsId,
    className = '',
}) => {
    const [fallOfWickets, setFallOfWickets] = useState<FallOfWickets[]>([]);
    const [summary, setSummary] = useState<FallOfWicketsSummary | null>(null);
    const [loading, setLoading] = useState(true);
    const [refreshing, setRefreshing] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [showAll, setShowAll] = useState(false);

    const fetchFallOfWickets = async (isRefresh = false) => {
        try {
            if (isRefresh) {
                setRefreshing(true);
            } else {
                setLoading(true);
            }
            setError(null);

            // Fetch both detailed list and summary
            const [listResponse, summaryResponse] = await Promise.all([
                inningsId
                    ? apiService.getFallOfWicketsByInnings(inningsId)
                    : apiService.getFallOfWicketsByMatch(matchId),
                apiService.getFallOfWicketsSummary(matchId, inningsId)
            ]);

            // Handle null responses gracefully
            setFallOfWickets(listResponse.data || []);
            setSummary(summaryResponse.data || null);
        } catch (err) {
            // Improved error handling to prevent null reference errors
            let errorMessage = 'Failed to fetch fall of wickets';

            if (err) {
                if (err instanceof Error) {
                    errorMessage = err.message || errorMessage;
                } else if (typeof err === 'object' && 'message' in err && err.message) {
                    errorMessage = String(err.message);
                } else if (typeof err === 'string') {
                    errorMessage = err;
                }
            }

            setError(errorMessage);
            console.error('Error fetching fall of wickets:', errorMessage);

            // Set empty state on error
            setFallOfWickets([]);
            setSummary(null);
        } finally {
            setLoading(false);
            setRefreshing(false);
        }
    };

    useEffect(() => {
        if (matchId) {
            fetchFallOfWickets();
        }
    }, [matchId, inningsId]);

    const handleRefresh = () => {
        fetchFallOfWickets(true);
    };

    if (loading) {
        return (
            <Card className={className}>
                <CardHeader>
                    <CardTitle className="text-lg font-semibold">Fall of Wickets</CardTitle>
                </CardHeader>
                <CardContent>
                    <div className="space-y-3">
                        {[...Array(5)].map((_, i) => (
                            <div key={i} className="flex items-center space-x-4 p-3 bg-gray-50 rounded-lg">
                                <div className="h-8 w-8 bg-gray-200 rounded-full animate-pulse" />
                                <div className="h-4 w-20 bg-gray-200 rounded animate-pulse" />
                                <div className="h-4 w-16 bg-gray-200 rounded animate-pulse" />
                                <div className="h-4 w-12 bg-gray-200 rounded animate-pulse" />
                                <div className="h-4 w-16 bg-gray-200 rounded animate-pulse" />
                            </div>
                        ))}
                    </div>
                </CardContent>
            </Card>
        );
    }

    if (error) {
        return (
            <Card className={className}>
                <CardHeader>
                    <CardTitle className="text-lg font-semibold">Fall of Wickets</CardTitle>
                </CardHeader>
                <CardContent>
                    <div className="text-center text-red-500 py-4">
                        <p>Error loading fall of wickets</p>
                        <p className="text-sm text-gray-500">{error}</p>
                    </div>
                </CardContent>
            </Card>
        );
    }

    if (!fallOfWickets || fallOfWickets.length === 0) {
        return (
            <Card className={className}>
                <CardHeader>
                    <CardTitle className="text-lg font-semibold">Fall of Wickets</CardTitle>
                </CardHeader>
                <CardContent>
                    <div className="text-center text-gray-500 py-4">
                        <p>No wickets have fallen yet</p>
                        <p className="text-sm text-gray-400 mt-1">Wickets will appear here as they fall during the match</p>
                    </div>
                </CardContent>
            </Card>
        );
    }

    const displayedWickets = showAll ? fallOfWickets : fallOfWickets.slice(0, 5);

    return (
        <Card className={className}>
            <CardHeader>
                <CardTitle className="text-lg font-semibold flex items-center justify-between">
                    <span>Fall of Wickets</span>
                    <div className="flex items-center space-x-2">
                        <Badge variant="secondary" className="text-xs">
                            {summary?.total_wickets || fallOfWickets.length} wickets
                        </Badge>
                        <Button
                            variant="ghost"
                            size="sm"
                            onClick={handleRefresh}
                            disabled={refreshing}
                            className="h-8 w-8 p-0"
                        >
                            <RefreshCw className={`h-4 w-4 ${refreshing ? 'animate-spin' : ''}`} />
                        </Button>
                    </div>
                </CardTitle>
            </CardHeader>
            <CardContent>
                <div className="flex flex-wrap gap-2">
                    {displayedWickets.map((wicket: FallOfWickets) => (
                        <div key={wicket.id} className="flex items-center space-x-2 px-3 py-2 bg-gray-100 rounded-lg border text-sm">
                            <div className="flex items-center justify-center w-6 h-6 bg-red-100 text-red-600 rounded-full text-xs font-semibold">
                                {wicket.wicket_number}
                            </div>
                            <span className="text-gray-900 font-medium">
                                {wicket.score}/{wicket.wicket_number} - {wicket.over_number}.{wicket.ball_number}
                            </span>
                        </div>
                    ))}
                </div>

                {fallOfWickets.length > 5 && (
                    <div className="text-center pt-2">
                        <Button
                            variant="outline"
                            size="sm"
                            onClick={() => setShowAll(!showAll)}
                            className="text-xs"
                        >
                            {showAll ? 'Show Less' : `Show All ${fallOfWickets.length} Wickets`}
                        </Button>
                    </div>
                )}
            </CardContent>
        </Card>
    );
};

export default FallOfWicketsDetailed;
