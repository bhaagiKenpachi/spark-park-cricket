'use client';

import { useState, useEffect } from 'react';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import {
    fetchVoteWithResultsRequest,
    castVoteRequest,
    checkUserVotedRequest,
} from '@/store/reducers/voteSlice';
import { VoteWithResults, VoteOption } from '@/types/vote';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import {
    ArrowLeft,
    Users,
    Calendar,
    CheckCircle,
    XCircle,
    BarChart3,
    Vote as VoteIcon,
    ChevronDown,
    ChevronUp,
    RefreshCw
} from 'lucide-react';
import TeamManagement from './TeamManagement';
import { User } from '@/types/team';

interface VoteViewProps {
    voteId: string;
    onBack?: () => void;
}

export function VoteView({ voteId, onBack }: VoteViewProps): React.JSX.Element {
    const dispatch = useAppDispatch();
    const { currentVote, loading, error, hasVotedStatus, userVotes } = useAppSelector(state => state.vote);
    const { user, isAuthenticated } = useAppSelector(state => state.auth);

    const [selectedOptions, setSelectedOptions] = useState<string[]>([]);
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [expandedOptions, setExpandedOptions] = useState<Set<string>>(new Set());
    const [isRefreshing, setIsRefreshing] = useState(false);
    const [showTeams, setShowTeams] = useState(false);
    const [showAuthPrompt, setShowAuthPrompt] = useState(false);

    useEffect(() => {
        // Add a small delay to allow backend to initialize vote results
        // This prevents showing 0 results immediately after vote creation
        const timer = setTimeout(() => {
            dispatch(fetchVoteWithResultsRequest(voteId));
            dispatch(checkUserVotedRequest(voteId));
        }, 300);

        return () => clearTimeout(timer);
    }, [dispatch, voteId]);

    useEffect(() => {
        if (currentVote?.user_vote) {
            setSelectedOptions(currentVote.user_vote.selected_options);
        }
    }, [currentVote]);

    const handleOptionToggle = (optionId: string) => {
        if (!currentVote) return;

        if (currentVote.vote.type === 'single') {
            setSelectedOptions([optionId]);
        } else {
            setSelectedOptions(prev =>
                prev.includes(optionId)
                    ? prev.filter(id => id !== optionId)
                    : [...prev, optionId]
            );
        }
    };

    const handleSubmitVote = async () => {
        // Check if user is authenticated
        if (!isAuthenticated) {
            setShowAuthPrompt(true);
            // Auto-hide after 5 seconds
            setTimeout(() => setShowAuthPrompt(false), 5000);
            return;
        }

        if (selectedOptions.length === 0) {
            alert('Please select at least one option');
            return;
        }

        if (currentVote?.vote.type === 'single' && selectedOptions.length > 1) {
            alert('Please select only one option');
            return;
        }

        setIsSubmitting(true);
        try {
            dispatch(castVoteRequest({
                voteId,
                optionIds: selectedOptions,
            }));

            // Reset and refresh after successful submission
            setTimeout(() => {
                setSelectedOptions([]);
                handleRefresh();
                setIsSubmitting(false);
            }, 500);
        } catch (error) {
            setIsSubmitting(false);
        }
    };

    const toggleOptionExpanded = (optionId: string) => {
        const newExpanded = new Set(expandedOptions);
        if (newExpanded.has(optionId)) {
            newExpanded.delete(optionId);
        } else {
            newExpanded.add(optionId);
        }
        setExpandedOptions(newExpanded);
    };

    const handleRefresh = () => {
        setIsRefreshing(true);
        dispatch(fetchVoteWithResultsRequest(voteId));
        dispatch(checkUserVotedRequest(voteId));
        setTimeout(() => setIsRefreshing(false), 500);
    };

    const getStatusBadge = (status: string) => {
        switch (status) {
            case 'active':
                return <Badge variant="default" className="bg-green-100 text-green-800">Active</Badge>;
            case 'closed':
                return <Badge variant="secondary" className="bg-gray-100 text-gray-800">Closed</Badge>;
            case 'cancelled':
                return <Badge variant="destructive">Cancelled</Badge>;
            default:
                return <Badge variant="outline">{status}</Badge>;
        }
    };

    const getTypeBadge = (type: string) => {
        switch (type) {
            case 'single':
                return <Badge variant="outline" className="bg-blue-50 text-blue-700">Single Choice</Badge>;
            case 'multiple':
                return <Badge variant="outline" className="bg-purple-50 text-purple-700">Multiple Choice</Badge>;
            default:
                return <Badge variant="outline">{type}</Badge>;
        }
    };

    const formatDate = (dateString: string) => {
        return new Date(dateString).toLocaleDateString('en-US', {
            year: 'numeric',
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit',
        });
    };

    const getPercentage = (count: number, total: number) => {
        if (total === 0) return 0;
        return Math.round((count / total) * 100);
    };

    if (loading) {
        return (
            <div className="w-full max-w-md mx-auto flex items-center justify-center p-8">
                <div className="text-center">
                    <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto mb-4"></div>
                    <p className="text-gray-600">Loading vote...</p>
                </div>
            </div>
        );
    }

    if (error || !currentVote) {
        return (
            <div className="w-full max-w-md mx-auto px-4 py-4">
                <div className="text-center">
                    <XCircle className="h-12 w-12 text-red-500 mx-auto mb-4" />
                    <h2 className="text-lg font-semibold text-gray-900 mb-2">Vote Not Found</h2>
                    <p className="text-sm text-gray-600 mb-4">{error || 'The requested vote could not be found.'}</p>
                    {onBack && (
                        <Button onClick={onBack} variant="outline" size="sm">
                            <ArrowLeft className="h-4 w-4 mr-2" />
                            Go Back
                        </Button>
                    )}
                </div>
            </div>
        );
    }

    const hasVoted = hasVotedStatus[voteId] || !!currentVote.user_vote;
    const canVote = currentVote.vote.status === 'active'; // Allow voting even if already voted (to update)

    return (
        <div className="w-full max-w-md mx-auto px-4 py-4">
            <div className="mb-4">
                <div className="flex items-center justify-between mb-3">
                    <h1 className="text-xl font-bold text-gray-900 line-clamp-2 flex-1 pr-2">{currentVote.vote.title}</h1>
                    <div className="flex flex-col gap-1">
                        {getStatusBadge(currentVote.vote.status)}
                        {getTypeBadge(currentVote.vote.type)}
                    </div>
                </div>
            </div>

            <div className="space-y-4">
                {/* Vote Details */}
                <div className="bg-white rounded-lg shadow-sm border">
                    <div className="p-3 border-b">
                        <div className="flex items-center gap-2">
                            <VoteIcon className="h-4 w-4 text-gray-600" />
                            <span className="text-sm font-medium text-gray-700">Details</span>
                        </div>
                    </div>
                    <div className="p-4 space-y-3">
                        {currentVote.vote.description && (
                            <div>
                                <p className="text-sm text-gray-600">{currentVote.vote.description}</p>
                            </div>
                        )}

                        <div className="flex items-center justify-between text-xs text-gray-500">
                            <div className="flex items-center">
                                <Calendar className="h-3 w-3 mr-1" />
                                {formatDate(currentVote.vote.created_at)}
                            </div>
                            <div className="flex items-center">
                                <Users className="h-3 w-3 mr-1" />
                                {currentVote.total_votes} votes
                            </div>
                        </div>

                        {/* Creator Info */}
                        <div className="flex items-center text-xs text-gray-500">
                            <Users className="h-3 w-3 mr-1" />
                            Created by {currentVote.creator_name || 'Unknown User'}
                        </div>
                    </div>
                </div>

                {/* Voting Options */}
                <div className="bg-white rounded-lg shadow-sm border">
                    <div className="p-3 border-b">
                        <div className="flex items-center gap-2">
                            <BarChart3 className="h-4 w-4 text-gray-600" />
                            <span className="text-sm font-medium text-gray-700">
                                {hasVoted ? 'Your Vote' : 'Vote Options'}
                            </span>
                        </div>
                    </div>
                    <div className="p-4">
                        {canVote ? (
                            <div className="space-y-3">
                                <p className="text-xs text-gray-600">
                                    {currentVote.vote.type === 'single'
                                        ? 'Select one option:'
                                        : 'Select one or more options:'}
                                </p>

                                <div className="space-y-2">
                                    {currentVote.options.map((option) => (
                                        <div
                                            key={option.id}
                                            className={`p-3 border rounded-lg cursor-pointer transition-colors ${selectedOptions.includes(option.id)
                                                ? 'border-blue-500 bg-blue-50'
                                                : 'border-gray-200 hover:border-gray-300'
                                                }`}
                                            onClick={() => handleOptionToggle(option.id)}
                                        >
                                            <div className="flex items-center justify-between">
                                                <span className="text-sm font-medium">{option.text}</span>
                                                {currentVote.vote.type === 'single' ? (
                                                    <div className={`w-4 h-4 rounded-full border-2 ${selectedOptions.includes(option.id)
                                                        ? 'border-blue-500 bg-blue-500'
                                                        : 'border-gray-300'
                                                        }`}>
                                                        {selectedOptions.includes(option.id) && (
                                                            <div className="w-2 h-2 bg-white rounded-full m-0.5"></div>
                                                        )}
                                                    </div>
                                                ) : (
                                                    <div className={`w-4 h-4 border-2 rounded ${selectedOptions.includes(option.id)
                                                        ? 'border-blue-500 bg-blue-500'
                                                        : 'border-gray-300'
                                                        }`}>
                                                        {selectedOptions.includes(option.id) && (
                                                            <CheckCircle className="w-3 h-3 text-white m-0.5" />
                                                        )}
                                                    </div>
                                                )}
                                            </div>
                                        </div>
                                    ))}
                                </div>

                                {!isAuthenticated && showAuthPrompt && (
                                    <div className="p-3 bg-yellow-50 border-2 border-yellow-400 rounded-lg animate-pulse">
                                        <p className="text-sm font-medium text-yellow-800 text-center">
                                            🔐 Please sign in to vote
                                        </p>
                                    </div>
                                )}

                                <Button
                                    onClick={handleSubmitVote}
                                    disabled={selectedOptions.length === 0 || isSubmitting}
                                    className="w-full"
                                >
                                    {isSubmitting ? 'Submitting...' : !isAuthenticated ? 'Sign in to Vote' : (hasVoted ? 'Update Vote' : 'Submit Vote')}
                                </Button>
                            </div>
                        ) : hasVoted ? (
                            <div className="space-y-3">
                                <div className="p-3 bg-green-50 border border-green-200 rounded-lg">
                                    <div className="flex items-center gap-2 text-green-800">
                                        <CheckCircle className="h-4 w-4" />
                                        <span className="text-sm font-medium">You have already voted</span>
                                    </div>
                                    <p className="text-xs text-green-700 mt-1">
                                        You selected: {currentVote.user_vote?.selected_options
                                            .map(optionId =>
                                                currentVote.options.find(opt => opt.id === optionId)?.text
                                            )
                                            .filter(Boolean)
                                            .join(', ')}
                                    </p>
                                </div>
                                <p className="text-xs text-gray-600 text-center">
                                    You can change your vote anytime while the poll is active
                                </p>
                            </div>
                        ) : (
                            <div className="p-3 bg-gray-50 border border-gray-200 rounded-lg">
                                <div className="flex items-center gap-2 text-gray-600">
                                    <XCircle className="h-4 w-4" />
                                    <span className="text-sm font-medium">Voting is closed</span>
                                </div>
                                <p className="text-xs text-gray-500 mt-1">
                                    This vote is no longer accepting new votes.
                                </p>
                            </div>
                        )}
                    </div>
                </div>

                {/* Results */}
                <div className="bg-white rounded-lg shadow-sm border">
                    <div className="p-3 border-b">
                        <div className="flex items-center justify-between">
                            <div className="flex items-center gap-2">
                                <BarChart3 className="h-4 w-4 text-gray-600" />
                                <span className="text-sm font-medium text-gray-700">Results</span>
                            </div>
                            <Button
                                variant="ghost"
                                size="sm"
                                onClick={handleRefresh}
                                disabled={isRefreshing}
                                className="h-7 px-2 flex items-center gap-1"
                                title="Refresh results"
                            >
                                <RefreshCw className={`h-3.5 w-3.5 ${isRefreshing ? 'animate-spin' : ''}`} />
                                <span className="text-xs">Refresh</span>
                            </Button>
                        </div>
                    </div>
                    <div className="p-4">
                        <div className="space-y-4">
                            {currentVote.options.map((option) => {
                                const count = currentVote.results[option.id] || 0;
                                const percentage = getPercentage(count, currentVote.total_votes);
                                const isSelected = currentVote.user_vote?.selected_options.includes(option.id);
                                const voters = currentVote.results_with_names?.[option.id] || [];

                                const isExpanded = expandedOptions.has(option.id);

                                return (
                                    <div key={option.id} className="space-y-2">
                                        <div className="flex items-center justify-between">
                                            <span className={`text-xs font-medium ${isSelected ? 'text-blue-600' : 'text-gray-900'}`}>
                                                {option.text}
                                                {isSelected && <span className="ml-1">(Your choice)</span>}
                                            </span>
                                            <span className="text-xs text-gray-600">
                                                {count} ({percentage}%)
                                            </span>
                                        </div>
                                        <div className="w-full bg-gray-200 rounded-full h-2">
                                            <div
                                                className={`h-2 rounded-full transition-all duration-300 ${isSelected ? 'bg-blue-500' : 'bg-gray-400'}`}
                                                style={{ width: `${percentage}%` }}
                                            ></div>
                                        </div>

                                        {/* View Voters Button */}
                                        {voters.length > 0 && (
                                            <div>
                                                <Button
                                                    variant="ghost"
                                                    size="sm"
                                                    onClick={() => toggleOptionExpanded(option.id)}
                                                    className="w-full h-7 text-xs justify-between px-2 hover:bg-gray-50"
                                                >
                                                    <span className="flex items-center gap-1">
                                                        <Users className="h-3 w-3" />
                                                        View {count} {count === 1 ? 'voter' : 'voters'}
                                                    </span>
                                                    {isExpanded ? (
                                                        <ChevronUp className="h-3 w-3" />
                                                    ) : (
                                                        <ChevronDown className="h-3 w-3" />
                                                    )}
                                                </Button>

                                                {/* Voter List (Collapsible) */}
                                                {isExpanded && (
                                                    <div className="mt-2 pl-2 space-y-1 border-l-2 border-gray-200">
                                                        {voters.map((voter, idx) => (
                                                            <div key={`${voter.user_id}-${idx}`} className="flex items-center justify-between py-1">
                                                                <span className="flex items-center gap-1.5 text-xs text-gray-700">
                                                                    <div className="w-1.5 h-1.5 rounded-full bg-blue-500"></div>
                                                                    {voter.user_name}
                                                                </span>
                                                                <span className="text-xs text-gray-400">
                                                                    {new Date(voter.voted_at).toLocaleDateString('en-US', {
                                                                        month: 'short',
                                                                        day: 'numeric',
                                                                        hour: '2-digit',
                                                                        minute: '2-digit'
                                                                    })}
                                                                </span>
                                                            </div>
                                                        ))}
                                                    </div>
                                                )}
                                            </div>
                                        )}
                                    </div>
                                );
                            })}
                        </div>

                        <div className="mt-4 pt-3 border-t border-gray-200">
                            <div className="text-center">
                                <p className="text-xs text-gray-600">
                                    Total votes: <span className="font-medium">{currentVote.total_votes}</span>
                                </p>
                            </div>
                        </div>
                    </div>
                </div>

                {/* Team Management Section - Only show if team formation is enabled */}
                {currentVote.vote.team_formation_enabled && (
                    <div className="mt-6">
                        <Button
                            variant="outline"
                            onClick={() => setShowTeams(!showTeams)}
                            className="w-full mb-4"
                        >
                            <Users className="h-4 w-4 mr-2" />
                            {showTeams ? 'Hide Teams' : 'Manage Teams'}
                        </Button>

                        {showTeams && (
                            <TeamManagement
                                vote={currentVote}
                                voters={getVotersFromResults(currentVote)}
                                isAuthenticated={isAuthenticated}
                                currentUserId={user?.id || undefined} // CRITICAL: Required for team creator authorization
                            />
                        )}
                    </div>
                )}
            </div>
        </div>
    );
}

// Helper function to extract unique voters from results
function getVotersFromResults(vote: VoteWithResults): User[] {
    const votersMap = new Map<string, User>();

    if (vote.results_with_names) {
        // results_with_names is an object with option_id as keys and voter arrays as values
        Object.values(vote.results_with_names).forEach((voters: any) => {
            if (Array.isArray(voters)) {
                voters.forEach((voter: any) => {
                    if (voter && voter.user_id && !votersMap.has(voter.user_id)) {
                        votersMap.set(voter.user_id, {
                            id: voter.user_id,
                            name: voter.user_name || 'Unknown',
                            email: '', // Not available in results
                        });
                    }
                });
            }
        });
    }

    return Array.from(votersMap.values());
}
