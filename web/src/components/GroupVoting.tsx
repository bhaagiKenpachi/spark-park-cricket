'use client';

import { useState, useEffect } from 'react';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import {
    fetchGroupsRequest,
    fetchVoteGroupsRequest,
    assignGroupsToVoteRequest,
    removeGroupsFromVoteRequest,
    clearError,
} from '@/store/reducers/groupSlice';
import { GroupWithCreator, GroupFilters } from '@/types/group';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Checkbox } from '@/components/ui/checkbox';
import {
    Users,
    Plus,
    Minus,
    Search,
    Filter,
    CheckCircle,
    XCircle
} from 'lucide-react';

interface GroupVotingProps {
    voteId: string;
    onSuccess?: () => void;
}

const GROUP_TYPES: Record<string, string> = {
    custom: 'Custom',
    team: 'Team',
    series: 'Series',
    match: 'Match',
    location: 'Location',
    skill: 'Skill Level',
};

export function GroupVoting({ voteId, onSuccess }: GroupVotingProps): React.JSX.Element {
    const dispatch = useAppDispatch();
    const {
        groups,
        loading,
        error
    } = useAppSelector(state => state.group);

    const [availableGroups, setAvailableGroups] = useState<GroupWithCreator[]>([]);
    const [assignedGroups, setAssignedGroups] = useState<GroupWithCreator[]>([]);
    const [selectedGroups, setSelectedGroups] = useState<string[]>([]);
    const [filters, setFilters] = useState<GroupFilters>({
        limit: 50,
        offset: 0,
    });
    const [searchQuery, setSearchQuery] = useState('');
    const [isAssigning, setIsAssigning] = useState(false);
    const [isRemoving, setIsRemoving] = useState(false);

    useEffect(() => {
        dispatch(fetchGroupsRequest(filters));
        dispatch(fetchVoteGroupsRequest(voteId));
    }, [dispatch, filters, voteId]);

    useEffect(() => {
        // This would be populated from the vote groups response
        // For now, we'll use empty arrays
        setAssignedGroups([]);
        setAvailableGroups(groups.filter(group =>
            !assignedGroups.some(assigned => assigned.id === group.id)
        ));
    }, [groups, assignedGroups]);

    const handleSearch = () => {
        setFilters(prev => ({
            ...prev,
            ...(searchQuery.trim() ? { search: searchQuery.trim() } : {}),
            offset: 0,
        }));
    };

    const handleGroupSelect = (groupId: string, checked: boolean) => {
        if (checked) {
            setSelectedGroups(prev => [...prev, groupId]);
        } else {
            setSelectedGroups(prev => prev.filter(id => id !== groupId));
        }
    };

    const handleSelectAll = (checked: boolean) => {
        if (checked) {
            setSelectedGroups(availableGroups.map(group => group.id));
        } else {
            setSelectedGroups([]);
        }
    };

    const handleAssignGroups = async () => {
        if (selectedGroups.length === 0) return;

        setIsAssigning(true);
        try {
            await dispatch(assignGroupsToVoteRequest({
                voteId,
                groupIds: selectedGroups
            }));
            setSelectedGroups([]);
            // Refresh vote groups
            dispatch(fetchVoteGroupsRequest(voteId));
            onSuccess?.();
        } catch (err) {
            console.error('Failed to assign groups:', err);
        } finally {
            setIsAssigning(false);
        }
    };

    const handleRemoveGroups = async (groupIds: string[]) => {
        setIsRemoving(true);
        try {
            await dispatch(removeGroupsFromVoteRequest({
                voteId,
                groupIds
            }));
            // Refresh vote groups
            dispatch(fetchVoteGroupsRequest(voteId));
            onSuccess?.();
        } catch (err) {
            console.error('Failed to remove groups:', err);
        } finally {
            setIsRemoving(false);
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

    const getStatusColor = (status: string): string => {
        const colors = {
            active: 'bg-green-100 text-green-800',
            inactive: 'bg-yellow-100 text-yellow-800',
            archived: 'bg-gray-100 text-gray-800',
        };
        return colors[status as keyof typeof colors] || 'bg-gray-100 text-gray-800';
    };

    return (
        <div className="space-y-6">
            {/* Header */}
            <div className="flex justify-between items-center">
                <h2 className="text-xl font-semibold">Group Voting Management</h2>
                <div className="text-sm text-gray-600">
                    {assignedGroups.length} groups assigned
                </div>
            </div>

            {/* Error Message */}
            {error && (
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
            )}

            {/* Assigned Groups */}
            {assignedGroups.length > 0 && (
                <Card>
                    <CardHeader>
                        <CardTitle className="flex items-center gap-2">
                            <CheckCircle className="w-5 h-5 text-green-600" />
                            Assigned Groups
                        </CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
                            {assignedGroups.map((group) => (
                                <div key={group.id} className="flex items-center justify-between p-3 border rounded-lg">
                                    <div className="flex-1">
                                        <h4 className="font-medium">{group.name}</h4>
                                        <div className="flex gap-1 mt-1">
                                            <Badge className={getTypeColor(group.type)}>
                                                {GROUP_TYPES[group.type]}
                                            </Badge>
                                            <Badge className={getStatusColor(group.status)}>
                                                {group.status}
                                            </Badge>
                                        </div>
                                    </div>
                                    <Button
                                        size="sm"
                                        variant="outline"
                                        onClick={() => handleRemoveGroups([group.id])}
                                        disabled={isRemoving}
                                    >
                                        <Minus className="w-4 h-4" />
                                    </Button>
                                </div>
                            ))}
                        </div>
                    </CardContent>
                </Card>
            )}

            {/* Available Groups */}
            <Card>
                <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                        <Users className="w-5 h-5" />
                        Available Groups
                    </CardTitle>
                </CardHeader>
                <CardContent>
                    {/* Search and Filter */}
                    <div className="flex gap-2 mb-4">
                        <div className="flex-1 relative">
                            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-4 h-4" />
                            <input
                                type="text"
                                placeholder="Search groups..."
                                value={searchQuery}
                                onChange={(e) => setSearchQuery(e.target.value)}
                                onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
                                className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                            />
                        </div>
                        <Button onClick={handleSearch} variant="outline">
                            <Filter className="w-4 h-4 mr-2" />
                            Filter
                        </Button>
                    </div>

                    {/* Select All */}
                    {availableGroups.length > 0 && (
                        <div className="flex items-center gap-2 mb-4 p-2 bg-gray-50 rounded-md">
                            <Checkbox
                                checked={selectedGroups.length === availableGroups.length}
                                onCheckedChange={handleSelectAll}
                            />
                            <span className="text-sm font-medium">
                                Select All ({availableGroups.length} groups)
                            </span>
                        </div>
                    )}

                    {/* Groups List */}
                    {loading ? (
                        <div className="space-y-3">
                            {[...Array(3)].map((_, i) => (
                                <div key={i} className="flex items-center gap-3 animate-pulse">
                                    <div className="w-4 h-4 bg-gray-200 rounded"></div>
                                    <div className="flex-1">
                                        <div className="h-4 bg-gray-200 rounded mb-1"></div>
                                        <div className="h-3 bg-gray-200 rounded w-1/2"></div>
                                    </div>
                                </div>
                            ))}
                        </div>
                    ) : availableGroups.length === 0 ? (
                        <div className="text-center py-8">
                            <Users className="w-12 h-12 text-gray-400 mx-auto mb-4" />
                            <p className="text-gray-500">No available groups found</p>
                        </div>
                    ) : (
                        <div className="space-y-2">
                            {availableGroups.map((group) => (
                                <div key={group.id} className="flex items-center gap-3 p-3 border rounded-lg hover:bg-gray-50">
                                    <Checkbox
                                        checked={selectedGroups.includes(group.id)}
                                        onCheckedChange={(checked) =>
                                            handleGroupSelect(group.id, checked as boolean)
                                        }
                                    />
                                    <div className="flex-1">
                                        <h4 className="font-medium">{group.name}</h4>
                                        {group.description && (
                                            <p className="text-sm text-gray-600 mt-1 line-clamp-1">
                                                {group.description}
                                            </p>
                                        )}
                                        <div className="flex gap-1 mt-1">
                                            <Badge className={getTypeColor(group.type)}>
                                                {GROUP_TYPES[group.type]}
                                            </Badge>
                                            <Badge className={getStatusColor(group.status)}>
                                                {group.status}
                                            </Badge>
                                        </div>
                                    </div>
                                </div>
                            ))}
                        </div>
                    )}

                    {/* Assign Button */}
                    {selectedGroups.length > 0 && (
                        <div className="mt-4 pt-4 border-t">
                            <Button
                                onClick={handleAssignGroups}
                                disabled={isAssigning}
                                className="w-full"
                            >
                                <Plus className="w-4 h-4 mr-2" />
                                {isAssigning
                                    ? 'Assigning...'
                                    : `Assign ${selectedGroups.length} Group${selectedGroups.length > 1 ? 's' : ''}`
                                }
                            </Button>
                        </div>
                    )}
                </CardContent>
            </Card>
        </div>
    );
}
