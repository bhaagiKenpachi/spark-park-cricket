'use client';

import { useState, useEffect } from 'react';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import {
    fetchVoteResultsByGroupsRequest,
    fetchVoteWithGroupResultsRequest,
    clearError,
} from '@/store/reducers/groupSlice';
import { VoteGroupResults, VoteWithGroupResults } from '@/types/group';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import {
    BarChart3,
    Users,
    RefreshCw
} from 'lucide-react';

interface VoteResultsByGroupsProps {
    voteId: string;
    showDetailed?: boolean;
}

const GROUP_TYPES: Record<string, string> = {
    custom: 'Custom',
    team: 'Team',
    series: 'Series',
    match: 'Match',
    location: 'Location',
    skill: 'Skill Level',
};

export function VoteResultsByGroups({ voteId, showDetailed = false }: VoteResultsByGroupsProps): React.JSX.Element {
    const dispatch = useAppDispatch();
    const {
        voteGroupResults,
        voteWithGroupResults,
        loading,
        error
    } = useAppSelector(state => state.group);

    const [activeTab, setActiveTab] = useState<'overview' | 'groups'>('overview');

    useEffect(() => {
        if (showDetailed) {
            dispatch(fetchVoteWithGroupResultsRequest(voteId));
        } else {
            dispatch(fetchVoteResultsByGroupsRequest(voteId));
        }
    }, [dispatch, voteId, showDetailed]);

    const handleRefresh = () => {
        if (showDetailed) {
            dispatch(fetchVoteWithGroupResultsRequest(voteId));
        } else {
            dispatch(fetchVoteResultsByGroupsRequest(voteId));
        }
    };

    const getTypeColor = (type: string): string => {
        const colors = {
            custom: 'bg-blue-100 text-blue-800',
            team: 'bg-green-100 text-green-800',
            series: 'bg-purple-100 text-purple-800',
            match: 'bg-orange-100 text-orange-800',
            location: 'bg-pink-100 text-pink-800',
            skill: 'bg-yellow-100 text-yellow-800',
        };
        return colors[type as keyof typeof colors] || 'bg-gray-100 text-gray-800';
    };

    const renderOverallResults = (results: VoteGroupResults | VoteWithGroupResults) => {
        return (
            <Card>
                <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                        <BarChart3 className="w-5 h-5" />
                        Overall Results
                    </CardTitle>
                </CardHeader>
                <CardContent>
                    <div className="space-y-4">
                        <div className="flex items-center justify-between">
                            <span className="text-sm font-medium">Total Votes</span>
                            <Badge variant="outline">{results.overall_results.total_votes}</Badge>
                        </div>

                        <div className="space-y-3">
                            {results.overall_results.options.map((option) => (
                                <div key={option.id} className="space-y-2">
                                    <div className="flex justify-between items-center">
                                        <span className="font-medium">{option.text}</span>
                                        <div className="flex items-center gap-2">
                                            <span className="text-sm text-gray-600">
                                                {option.vote_count} votes
                                            </span>
                                            {'percentage' in option && (
                                                <Badge variant="secondary">
                                                    {option.percentage.toFixed(1)}%
                                                </Badge>
                                            )}
                                        </div>
                                    </div>
                                    <div className="w-full bg-gray-200 rounded-full h-2">
                                        <div
                                            className="bg-blue-600 h-2 rounded-full transition-all duration-300"
                                            style={{
                                                width: `${('percentage' in option ? option.percentage :
                                                    (option.vote_count / results.overall_results.total_votes) * 100) || 0}%`
                                            }}
                                        ></div>
                                    </div>
                                </div>
                            ))}
                        </div>
                    </div>
                </CardContent>
            </Card>
        );
    };

    const renderGroupResults = (results: VoteGroupResults | VoteWithGroupResults) => {
        const groupResults = 'group_results' in results ? results.group_results : [];

        return (
            <Card>
                <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                        <Users className="w-5 h-5" />
                        Results by Groups
                    </CardTitle>
                </CardHeader>
                <CardContent>
                    <div className="space-y-6">
                        {groupResults.map((groupResult) => (
                            <div key={groupResult.group.id} className="border rounded-lg p-4">
                                <div className="flex items-center justify-between mb-4">
                                    <div>
                                        <h4 className="font-semibold">{groupResult.group.name}</h4>
                                        <div className="flex gap-1 mt-1">
                                            <Badge className={getTypeColor(groupResult.group.type)}>
                                                {GROUP_TYPES[groupResult.group.type]}
                                            </Badge>
                                            <Badge variant="outline">
                                                {groupResult.total_votes} votes
                                            </Badge>
                                        </div>
                                    </div>
                                </div>

                                <div className="space-y-3">
                                    {groupResult.options.map((option) => (
                                        <div key={option.id} className="space-y-2">
                                            <div className="flex justify-between items-center">
                                                <span className="text-sm font-medium">{option.text}</span>
                                                <div className="flex items-center gap-2">
                                                    <span className="text-xs text-gray-600">
                                                        {option.vote_count} votes
                                                    </span>
                                                    {'percentage' in option && (
                                                        <Badge variant="secondary" className="text-xs">
                                                            {option.percentage.toFixed(1)}%
                                                        </Badge>
                                                    )}
                                                </div>
                                            </div>
                                            <div className="w-full bg-gray-200 rounded-full h-1.5">
                                                <div
                                                    className="bg-green-600 h-1.5 rounded-full transition-all duration-300"
                                                    style={{
                                                        width: `${('percentage' in option ? option.percentage :
                                                            (option.vote_count / groupResult.total_votes) * 100) || 0}%`
                                                    }}
                                                ></div>
                                            </div>
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

    if (loading) {
        return (
            <div className="space-y-4">
                <div className="flex justify-between items-center">
                    <h2 className="text-xl font-semibold">Vote Results by Groups</h2>
                    <Button variant="outline" disabled>
                        <RefreshCw className="w-4 h-4 mr-2 animate-spin" />
                        Loading...
                    </Button>
                </div>

                <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                    <Card>
                        <CardContent className="p-6">
                            <div className="space-y-4 animate-pulse">
                                <div className="h-4 bg-gray-200 rounded"></div>
                                <div className="h-3 bg-gray-200 rounded"></div>
                                <div className="h-3 bg-gray-200 rounded w-3/4"></div>
                            </div>
                        </CardContent>
                    </Card>
                    <Card>
                        <CardContent className="p-6">
                            <div className="space-y-4 animate-pulse">
                                <div className="h-4 bg-gray-200 rounded"></div>
                                <div className="h-3 bg-gray-200 rounded"></div>
                                <div className="h-3 bg-gray-200 rounded w-3/4"></div>
                            </div>
                        </CardContent>
                    </Card>
                </div>
            </div>
        );
    }

    if (error) {
        return (
            <div className="space-y-4">
                <div className="flex justify-between items-center">
                    <h2 className="text-xl font-semibold">Vote Results by Groups</h2>
                    <Button variant="outline" onClick={handleRefresh}>
                        <RefreshCw className="w-4 h-4 mr-2" />
                        Retry
                    </Button>
                </div>

                <div className="p-4 bg-red-50 border border-red-200 rounded-md">
                    <p className="text-sm text-red-600">{error}</p>
                    <Button
                        size="sm"
                        variant="outline"
                        onClick={() => dispatch(clearError())}
                        className="mt-2"
                    >
                        Dismiss
                    </Button>
                </div>
            </div>
        );
    }

    const results = showDetailed ? voteWithGroupResults : voteGroupResults;

    if (!results) {
        return (
            <div className="space-y-4">
                <div className="flex justify-between items-center">
                    <h2 className="text-xl font-semibold">Vote Results by Groups</h2>
                    <Button variant="outline" onClick={handleRefresh}>
                        <RefreshCw className="w-4 h-4 mr-2" />
                        Refresh
                    </Button>
                </div>

                <Card>
                    <CardContent className="p-8 text-center">
                        <BarChart3 className="w-12 h-12 text-gray-400 mx-auto mb-4" />
                        <h3 className="text-lg font-medium text-gray-900 mb-2">No results available</h3>
                        <p className="text-gray-500">
                            This vote may not have any groups assigned or no votes have been cast yet.
                        </p>
                    </CardContent>
                </Card>
            </div>
        );
    }

    return (
        <div className="space-y-6">
            {/* Header */}
            <div className="flex justify-between items-center">
                <h2 className="text-xl font-semibold">Vote Results by Groups</h2>
                <div className="flex gap-2">
                    {showDetailed && (
                        <div className="flex gap-1 border rounded-md">
                            <Button
                                variant={activeTab === 'overview' ? 'default' : 'ghost'}
                                size="sm"
                                onClick={() => setActiveTab('overview')}
                            >
                                <BarChart3 className="w-4 h-4 mr-2" />
                                Overview
                            </Button>
                            <Button
                                variant={activeTab === 'groups' ? 'default' : 'ghost'}
                                size="sm"
                                onClick={() => setActiveTab('groups')}
                            >
                                <Users className="w-4 h-4 mr-2" />
                                Groups
                            </Button>
                        </div>
                    )}
                    <Button variant="outline" onClick={handleRefresh}>
                        <RefreshCw className="w-4 h-4 mr-2" />
                        Refresh
                    </Button>
                </div>
            </div>

            {/* Vote Info */}
            <Card>
                <CardHeader>
                    <CardTitle>{results.vote.title}</CardTitle>
                </CardHeader>
                <CardContent>
                    {results.vote.description && (
                        <p className="text-gray-600 mb-2">{results.vote.description}</p>
                    )}
                    <div className="flex gap-2">
                        <Badge variant="outline">{results.vote.type}</Badge>
                        <Badge className={
                            results.vote.status === 'active' ? 'bg-green-100 text-green-800' :
                                results.vote.status === 'closed' ? 'bg-gray-100 text-gray-800' :
                                    'bg-red-100 text-red-800'
                        }>
                            {results.vote.status}
                        </Badge>
                    </div>
                </CardContent>
            </Card>

            {/* Results */}
            {showDetailed ? (
                <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                    {activeTab === 'overview' && renderOverallResults(results)}
                    {activeTab === 'groups' && renderGroupResults(results)}
                </div>
            ) : (
                <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                    {renderOverallResults(results)}
                    {renderGroupResults(results)}
                </div>
            )}
        </div>
    );
}
