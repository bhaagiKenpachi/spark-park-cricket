'use client';

import { useState, useEffect } from 'react';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import {
    fetchVotesRequest,
    deleteVoteRequest,
    closeVoteRequest,
    cancelVoteRequest,
    setFilters,
} from '@/store/reducers/voteSlice';
import { VoteFilters } from '@/types/vote';
import { Vote, VoteStatus, VoteType, VoteWithCreator } from '@/types/vote';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuTrigger,
    DropdownMenuSeparator,
} from '@/components/ui/dropdown-menu';
import {
    Trash2,
    Eye,
    Edit,
    Lock,
    X,
    Filter,
    Plus,
    Users,
    Calendar,
    CheckCircle,
    XCircle,
    ChevronLeft,
    ChevronRight,
    MoreVertical,
    Share2
} from 'lucide-react';

interface VoteListProps {
    onCreateVote?: () => void;
    onViewVote?: (voteId: string) => void;
    onEditVote?: (vote: Vote) => void;
}

export function VoteList({
    onCreateVote,
    onViewVote,
    onEditVote,
}: VoteListProps): React.JSX.Element {
    const dispatch = useAppDispatch();
    const { votes, loading, error, filters, pagination } = useAppSelector(state => state.vote);
    const { user: currentUser, isAuthenticated } = useAppSelector(state => state.auth);

    const [localFilters, setLocalFilters] = useState<VoteFilters>({
        ...filters,
        page: 1,
        page_size: 20,
    });

    // Ensure votes is always an array
    const votesList = Array.isArray(votes) ? votes : [];

    useEffect(() => {
        dispatch(fetchVotesRequest(localFilters));
    }, [dispatch, localFilters]);

    const handleFilterChange = (key: keyof VoteFilters, value: string | undefined) => {
        const newFilters = { ...localFilters, page: 1 }; // Reset to page 1 when filters change
        if (value) {
            (newFilters as any)[key] = value;
        } else {
            delete (newFilters as any)[key];
        }
        setLocalFilters(newFilters);
        dispatch(setFilters(newFilters));
    };

    const handlePageChange = (newPage: number) => {
        const newFilters = { ...localFilters, page: newPage };
        setLocalFilters(newFilters);
    };

    const handlePageSizeChange = (newPageSize: number) => {
        const newFilters = { ...localFilters, page: 1, page_size: newPageSize }; // Reset to page 1 when page size changes
        setLocalFilters(newFilters);
    };

    const handleDeleteVote = (voteId: string) => {
        if (window.confirm('Are you sure you want to delete this vote? This action cannot be undone.')) {
            dispatch(deleteVoteRequest(voteId));
        }
    };

    const handleCloseVote = (voteId: string) => {
        if (window.confirm('Are you sure you want to close this vote? Users will no longer be able to vote.')) {
            dispatch(closeVoteRequest(voteId));
        }
    };

    const handleCancelVote = (voteId: string) => {
        if (window.confirm('Are you sure you want to cancel this vote? This action cannot be undone.')) {
            dispatch(cancelVoteRequest(voteId));
        }
    };

    const handleShareVote = async (voteId: string, event: React.MouseEvent) => {
        // Prevent navigation when clicking share
        event.stopPropagation();

        const url = `${window.location.origin}/vote/${voteId}`;

        try {
            await navigator.clipboard.writeText(url);
            alert('URL copied to clipboard!');
        } catch (error) {
            // If clipboard API fails, use the old method
            const textArea = document.createElement('textarea');
            textArea.value = url;
            document.body.appendChild(textArea);
            textArea.select();
            document.execCommand('copy');
            document.body.removeChild(textArea);
            alert('URL copied to clipboard!');
        }
    };

    const getStatusBadge = (status: VoteStatus) => {
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

    const getTypeBadge = (type: VoteType) => {
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

    if (loading) {
        return (
            <div className="flex items-center justify-center p-8">
                <div className="text-center">
                    <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto mb-4"></div>
                    <p className="text-gray-600">Loading votes...</p>
                </div>
            </div>
        );
    }

    return (
        <div className="w-full max-w-md mx-auto px-4 py-4">
            <div className="mb-6">
                <h1 className="text-xl font-bold text-gray-900">Votes</h1>
                <p className="text-sm text-gray-600">Manage and participate in voting polls</p>
            </div>

            {/* Filters */}
            <div className="bg-white rounded-lg shadow-sm border mb-3">
                <div className="p-2">
                    <div className="flex items-center gap-1.5 mb-2">
                        <Filter className="h-3.5 w-3.5 text-gray-600" />
                        <span className="text-xs font-medium text-gray-700">Filters</span>
                    </div>
                    <div className="grid grid-cols-2 gap-2">
                        <div>
                            <label className="block text-[10px] font-medium text-gray-600 mb-0.5">
                                Status
                            </label>
                            <Select
                                value={localFilters.status || 'all'}
                                onValueChange={(value) => handleFilterChange('status', value === 'all' ? undefined : value)}
                            >
                                <SelectTrigger className="w-full h-8 text-xs">
                                    <SelectValue placeholder="All" />
                                </SelectTrigger>
                                <SelectContent>
                                    <SelectItem value="all">All</SelectItem>
                                    <SelectItem value="active">Active</SelectItem>
                                    <SelectItem value="closed">Closed</SelectItem>
                                    <SelectItem value="cancelled">Cancelled</SelectItem>
                                </SelectContent>
                            </Select>
                        </div>

                        <div>
                            <label className="block text-[10px] font-medium text-gray-600 mb-0.5">
                                Type
                            </label>
                            <Select
                                value={localFilters.type || 'all'}
                                onValueChange={(value) => handleFilterChange('type', value === 'all' ? undefined : value)}
                            >
                                <SelectTrigger className="w-full h-8 text-xs">
                                    <SelectValue placeholder="All" />
                                </SelectTrigger>
                                <SelectContent>
                                    <SelectItem value="all">All</SelectItem>
                                    <SelectItem value="single">Single</SelectItem>
                                    <SelectItem value="multiple">Multiple</SelectItem>
                                </SelectContent>
                            </Select>
                        </div>
                    </div>
                </div>
            </div>

            {error && (
                <div className="mb-4 p-3 bg-red-50 border border-red-200 text-red-700 rounded-lg text-sm">
                    {error}
                </div>
            )}

            {votesList.length === 0 ? (
                <div className="bg-white rounded-lg shadow-sm border">
                    <div className="text-center py-12 px-4">
                        <Users className="h-12 w-12 text-gray-400 mx-auto mb-4" />
                        <h3 className="text-base font-medium text-gray-900 mb-2">No votes found</h3>
                        <p className="text-sm text-gray-600 mb-4">
                            {Object.keys(localFilters).length > 0
                                ? 'No votes match your current filters.'
                                : 'Get started by creating your first vote.'}
                        </p>
                    </div>
                </div>
            ) : (
                <div className="space-y-3">
                    {votesList.map((vote) => (
                        <div key={vote.id} className="bg-white rounded-lg shadow-sm border">
                            <div className="p-4">
                                {/* Header */}
                                <div className="flex items-start justify-between mb-2">
                                    <h3 className="text-base font-semibold text-gray-900 line-clamp-2 flex-1 pr-2">
                                        {vote.title}
                                    </h3>
                                    <div className="flex flex-col gap-1">
                                        {getStatusBadge(vote.status)}
                                        {getTypeBadge(vote.type)}
                                    </div>
                                </div>

                                {/* Description */}
                                {vote.description && (
                                    <p className="text-sm text-gray-600 line-clamp-2 mb-3">
                                        {vote.description}
                                    </p>
                                )}

                                {/* Date and Creator */}
                                <div className="flex items-center justify-between text-xs text-gray-500 mb-3">
                                    <div className="flex items-center">
                                        <Calendar className="h-3 w-3 mr-1" />
                                        Created {formatDate(vote.created_at)}
                                    </div>
                                    <div className="flex items-center">
                                        <Users className="h-3 w-3 mr-1" />
                                        by {(vote as VoteWithCreator).creator_name || 'Unknown User'}
                                    </div>
                                </div>

                                {/* Actions */}
                                <div className="flex items-center justify-between">
                                    <div className="flex items-center space-x-2">
                                        {onViewVote && (
                                            <Button
                                                variant="outline"
                                                size="sm"
                                                onClick={() => onViewVote(vote.id)}
                                                className="flex items-center gap-1 text-xs px-2 py-1"
                                            >
                                                <Eye className="h-3 w-3" />
                                                View
                                            </Button>
                                        )}

                                        <Button
                                            variant="ghost"
                                            size="sm"
                                            onClick={(e) => handleShareVote(vote.id, e)}
                                            className="hover:bg-green-50 text-green-600 hover:text-green-700"
                                            title="Share Vote"
                                        >
                                            <Share2 className="h-4 w-4" />
                                        </Button>
                                    </div>

                                    {/* Only show actions menu if user is authenticated and is the creator */}
                                    {isAuthenticated && currentUser && currentUser.id === vote.created_by && (
                                        <DropdownMenu>
                                            <DropdownMenuTrigger asChild>
                                                <Button
                                                    variant="ghost"
                                                    size="sm"
                                                    className="h-8 w-8 p-0"
                                                >
                                                    <MoreVertical className="h-4 w-4" />
                                                </Button>
                                            </DropdownMenuTrigger>
                                            <DropdownMenuContent align="end">
                                                {onEditVote && vote.status === 'active' && (
                                                    <DropdownMenuItem
                                                        onClick={() => onEditVote(vote)}
                                                        className="cursor-pointer"
                                                    >
                                                        <Edit className="h-4 w-4 mr-2" />
                                                        Edit
                                                    </DropdownMenuItem>
                                                )}

                                                {vote.status === 'active' && (
                                                    <>
                                                        <DropdownMenuItem
                                                            onClick={() => handleCloseVote(vote.id)}
                                                            className="cursor-pointer text-green-600"
                                                        >
                                                            <CheckCircle className="h-4 w-4 mr-2" />
                                                            Close Vote
                                                        </DropdownMenuItem>
                                                        <DropdownMenuItem
                                                            onClick={() => handleCancelVote(vote.id)}
                                                            className="cursor-pointer text-orange-600"
                                                        >
                                                            <XCircle className="h-4 w-4 mr-2" />
                                                            Cancel Vote
                                                        </DropdownMenuItem>
                                                    </>
                                                )}

                                                {(onEditVote || vote.status === 'active') && (
                                                    <DropdownMenuSeparator />
                                                )}

                                                <DropdownMenuItem
                                                    onClick={() => handleDeleteVote(vote.id)}
                                                    className="cursor-pointer text-red-600"
                                                >
                                                    <Trash2 className="h-4 w-4 mr-2" />
                                                    Delete
                                                </DropdownMenuItem>
                                            </DropdownMenuContent>
                                        </DropdownMenu>
                                    )}
                                </div>
                            </div>
                        </div>
                    ))}
                </div>
            )}

            {/* Pagination Controls */}
            {!loading && votesList.length > 0 && (
                <div className="mt-4 border-t pt-4">
                    <div className="flex items-center justify-between text-sm">
                        <div className="text-gray-600">
                            Page {pagination.currentPage} of {pagination.totalPages} ({pagination.totalItems} total)
                        </div>
                        <div className="flex items-center gap-2">
                            <Button
                                variant="outline"
                                size="sm"
                                onClick={() => handlePageChange(pagination.currentPage - 1)}
                                disabled={pagination.currentPage === 1 || loading}
                                className="h-8 px-2"
                            >
                                <ChevronLeft className="h-4 w-4" />
                            </Button>
                            <span className="text-xs text-gray-600 min-w-[60px] text-center">
                                {pagination.currentPage} / {pagination.totalPages}
                            </span>
                            <Button
                                variant="outline"
                                size="sm"
                                onClick={() => handlePageChange(pagination.currentPage + 1)}
                                disabled={pagination.currentPage >= pagination.totalPages || loading}
                                className="h-8 px-2"
                            >
                                <ChevronRight className="h-4 w-4" />
                            </Button>
                        </div>
                    </div>

                    {/* Page Size Selector */}
                    <div className="mt-3 flex items-center justify-between">
                        <span className="text-xs text-gray-600">Items per page:</span>
                        <Select
                            value={String(localFilters.page_size || 20)}
                            onValueChange={(value) => handlePageSizeChange(Number(value))}
                        >
                            <SelectTrigger className="w-20 h-8 text-xs">
                                <SelectValue />
                            </SelectTrigger>
                            <SelectContent>
                                <SelectItem value="10">10</SelectItem>
                                <SelectItem value="20">20</SelectItem>
                                <SelectItem value="50">50</SelectItem>
                                <SelectItem value="100">100</SelectItem>
                            </SelectContent>
                        </Select>
                    </div>
                </div>
            )}
        </div>
    );
}
