'use client';

import React, { useState, useEffect } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { RefreshCw } from 'lucide-react';
import { FallOfWicketsSummary, WicketFall, FallOfWickets } from '@/types/fallOfWickets';
import { apiService } from '@/services/api';

interface FallOfWicketsDisplayProps {
    matchId: string;
    className?: string;
}

interface InningsFallOfWickets {
    inningsNumber: number;
    wickets: WicketFall[];
    totalWickets: number;
}

export const FallOfWicketsDisplay: React.FC<FallOfWicketsDisplayProps> = ({
    matchId,
    className = '',
}) => {
    const [inningsData, setInningsData] = useState<InningsFallOfWickets[]>([]);
    const [loading, setLoading] = useState(true);
    const [refreshing, setRefreshing] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const fetchFallOfWickets = async (isRefresh = false) => {
        try {
            if (isRefresh) {
                setRefreshing(true);
            } else {
                setLoading(true);
            }
            setError(null);

            // Fetch fall of wickets for the entire match (both innings)
            const response = await apiService.getFallOfWicketsByMatch(matchId);

            // Handle null response gracefully
            const allWickets: FallOfWickets[] = response.data || [];

            // If no wickets data, show empty state
            if (!allWickets || allWickets.length === 0) {
                setInningsData([]);
                return;
            }

            // Group wickets by innings number
            const inningsMap = new Map<number, FallOfWickets[]>();
            allWickets.forEach(wicket => {
                const inningsNumber = wicket.innings_number;
                if (!inningsMap.has(inningsNumber)) {
                    inningsMap.set(inningsNumber, []);
                }
                inningsMap.get(inningsNumber)!.push(wicket);
            });

            // Convert to our format
            const inningsArray: InningsFallOfWickets[] = [];
            inningsMap.forEach((wickets, inningsNumber) => {
                const wicketsFormatted: WicketFall[] = wickets.map(w => ({
                    wicket_number: w.wicket_number,
                    score: w.score,
                    over_number: w.over_number,
                    ball_number: w.ball_number,
                    over_position: `${w.over_number}.${w.ball_number}`
                }));

                inningsArray.push({
                    inningsNumber: inningsNumber,
                    wickets: wicketsFormatted,
                    totalWickets: wicketsFormatted.length
                });
            });

            // Sort by innings number
            inningsArray.sort((a, b) => a.inningsNumber - b.inningsNumber);
            setInningsData(inningsArray);
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
            setInningsData([]);
        } finally {
            setLoading(false);
            setRefreshing(false);
        }
    };

    useEffect(() => {
        if (matchId) {
            fetchFallOfWickets();
        }
    }, [matchId]);

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
                        {[...Array(3)].map((_, i) => (
                            <div key={i} className="flex items-center space-x-4">
                                <div className="h-6 w-6 bg-gray-200 rounded-full animate-pulse" />
                                <div className="h-4 w-16 bg-gray-200 rounded animate-pulse" />
                                <div className="h-4 w-12 bg-gray-200 rounded animate-pulse" />
                                <div className="h-4 w-20 bg-gray-200 rounded animate-pulse" />
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

    if (!inningsData || inningsData.length === 0) {
        return (
            <Card className={className}>
                <CardHeader>
                    <CardTitle className="text-lg font-semibold flex items-center justify-between">
                        <span>Fall of Wickets</span>
                        <div className="flex items-center space-x-2">
                            <Badge variant="secondary" className="text-xs">
                                0 wickets
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
                    <div className="text-center text-gray-500 py-4">
                        <p>No wickets have fallen yet</p>
                        <p className="text-sm text-gray-400 mt-1">Wickets will appear here as they fall during the match</p>
                    </div>
                </CardContent>
            </Card>
        );
    }

    const totalWickets = inningsData.reduce((sum, innings) => sum + innings.totalWickets, 0);

    return (
        <Card className={className}>
            <CardHeader>
                <CardTitle className="text-lg font-semibold flex items-center justify-between">
                    <span>Fall of Wickets</span>
                    <div className="flex items-center space-x-2">
                        <Badge variant="secondary" className="text-xs">
                            {totalWickets} wickets
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
                <div className="space-y-4">
                    {inningsData.map((innings) => (
                        <div key={innings.inningsNumber}>
                            <div className="flex items-center justify-between mb-2">
                                <h4 className="text-sm font-medium text-gray-700">
                                    Innings {innings.inningsNumber}
                                </h4>
                                <Badge variant="outline" className="text-xs">
                                    {innings.totalWickets} wickets
                                </Badge>
                            </div>
                            <div className="flex flex-wrap gap-2">
                                {innings.wickets.map((wicket: WicketFall, index: number) => (
                                    <div key={index} className="flex items-center space-x-2 px-3 py-2 bg-gray-100 rounded-lg border text-sm">
                                        <div className="flex items-center justify-center w-6 h-6 bg-red-100 text-red-600 rounded-full text-xs font-semibold">
                                            {wicket.wicket_number}
                                        </div>
                                        <span className="text-gray-900 font-medium">
                                            {wicket.score}/{wicket.wicket_number} - {wicket.over_position}
                                        </span>
                                    </div>
                                ))}
                            </div>
                        </div>
                    ))}
                </div>
            </CardContent>
        </Card>
    );
};

export default FallOfWicketsDisplay;
