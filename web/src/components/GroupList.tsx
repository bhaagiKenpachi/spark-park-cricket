'use client';

import { useState, useEffect } from 'react';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import {
    fetchGroupsRequest,
    deleteGroupRequest,
    clearError,
} from '@/store/reducers/groupSlice';
import { GroupWithCreator, GroupType, GroupFilters } from '@/types/group';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Badge } from '@/components/ui/badge';
import {
    Plus,
    Search,
    Users,
    Calendar,
    MoreHorizontal,
    Edit,
    Trash2,
    Eye,
    Filter
} from 'lucide-react';
import { GroupForm } from './GroupForm';
import { GroupView } from './GroupView';

const GROUP_TYPES: { value: GroupType | 'all'; label: string }[] = [
    { value: 'all', label: 'All Types' },
    { value: 'custom', label: 'Custom' },
    { value: 'team', label: 'Team' },
    { value: 'series', label: 'Series' },
    { value: 'match', label: 'Match' },
    { value: 'location', label: 'Location' },
    { value: 'skill', label: 'Skill Level' },
];

const GROUP_STATUSES = [
    { value: 'all', label: 'All Status' },
    { value: 'active', label: 'Active' },
    { value: 'inactive', label: 'Inactive' },
    { value: 'archived', label: 'Archived' },
];

interface GroupListProps {
    onGroupSelect?: (group: GroupWithCreator) => void;
}

export function GroupList({ onGroupSelect }: GroupListProps): React.JSX.Element {
    const dispatch = useAppDispatch();
    const { groups, loading, error, pagination } = useAppSelector(state => state.group);

    const [showForm, setShowForm] = useState(false);
    const [editingGroup, setEditingGroup] = useState<GroupWithCreator | null>(null);
    const [viewingGroup, setViewingGroup] = useState<GroupWithCreator | null>(null);
    const [filters, setFilters] = useState<GroupFilters>({
        limit: 20,
        offset: 0,
    });
    const [searchQuery, setSearchQuery] = useState('');

    useEffect(() => {
        dispatch(fetchGroupsRequest(filters));
    }, [dispatch, filters]);

    const handleSearch = () => {
        setFilters(prev => ({
            ...prev,
            ...(searchQuery.trim() ? { search: searchQuery.trim() } : {}),
            offset: 0,
        }));
    };

    const handleFilterChange = (key: keyof GroupFilters, value: string) => {
        if (value === 'all') {
            setFilters(prev => {
                const newFilters = { ...prev };
                delete newFilters[key as keyof GroupFilters];
                return { ...newFilters, offset: 0 };
            });
        } else {
            setFilters(prev => ({
                ...prev,
                [key]: value,
                offset: 0,
            }));
        }
    };

    const handlePageChange = (page: number) => {
        if (!pagination) return;

        const newOffset = (page - 1) * pagination.page_size;
        setFilters(prev => ({
            ...prev,
            offset: newOffset,
        }));
    };

    const handleDeleteGroup = async (group: GroupWithCreator) => {
        if (window.confirm(`Are you sure you want to delete "${group.name}"? This action cannot be undone.`)) {
            await dispatch(deleteGroupRequest(group.id));
        }
    };

    const handleFormSuccess = () => {
        setShowForm(false);
        setEditingGroup(null);
        dispatch(fetchGroupsRequest(filters));
    };

    const handleFormCancel = () => {
        setShowForm(false);
        setEditingGroup(null);
    };

    const getTypeColor = (type: GroupType): string => {
        const colors = {
            custom: 'bg-blue-100 text-blue-800',
            team: 'bg-green-100 text-green-800',
            series: 'bg-purple-100 text-purple-800',
            match: 'bg-orange-100 text-orange-800',
            location: 'bg-pink-100 text-pink-800',
            skill: 'bg-yellow-100 text-yellow-800',
        };
        return colors[type] || 'bg-gray-100 text-gray-800';
    };

    const getStatusColor = (status: string): string => {
        const colors = {
            active: 'bg-green-100 text-green-800',
            inactive: 'bg-yellow-100 text-yellow-800',
            archived: 'bg-gray-100 text-gray-800',
        };
        return colors[status as keyof typeof colors] || 'bg-gray-100 text-gray-800';
    };

    if (showForm || editingGroup) {
        return (
            <GroupForm
                group={editingGroup || undefined}
                onSuccess={handleFormSuccess}
                onCancel={handleFormCancel}
            />
        );
    }

    if (viewingGroup) {
        return (
            <GroupView
                group={viewingGroup}
                onBack={() => setViewingGroup(null)}
                onEdit={(group) => {
                    setViewingGroup(null);
                    setEditingGroup(group);
                }}
            />
        );
    }

    return (
        <div className="space-y-6">
            {/* Header */}
            <div className="flex justify-between items-center">
                <h1 className="text-2xl font-bold">Groups</h1>
                <Button onClick={() => setShowForm(true)}>
                    <Plus className="w-4 h-4 mr-2" />
                    Create Group
                </Button>
            </div>

            {/* Filters */}
            <Card>
                <CardContent className="pt-6">
                    <div className="flex flex-col sm:flex-row gap-4">
                        {/* Search */}
                        <div className="flex-1">
                            <div className="relative">
                                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-4 h-4" />
                                <Input
                                    placeholder="Search groups..."
                                    value={searchQuery}
                                    onChange={(e) => setSearchQuery(e.target.value)}
                                    onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
                                    className="pl-10"
                                />
                            </div>
                        </div>

                        {/* Type Filter */}
                        <Select
                            value={filters.type || 'all'}
                            onValueChange={(value) => handleFilterChange('type', value)}
                        >
                            <SelectTrigger className="w-full sm:w-40">
                                <SelectValue placeholder="Type" />
                            </SelectTrigger>
                            <SelectContent>
                                {GROUP_TYPES.map((type) => (
                                    <SelectItem key={type.value} value={type.value}>
                                        {type.label}
                                    </SelectItem>
                                ))}
                            </SelectContent>
                        </Select>

                        {/* Status Filter */}
                        <Select
                            value={filters.status || 'all'}
                            onValueChange={(value) => handleFilterChange('status', value)}
                        >
                            <SelectTrigger className="w-full sm:w-40">
                                <SelectValue placeholder="Status" />
                            </SelectTrigger>
                            <SelectContent>
                                {GROUP_STATUSES.map((status) => (
                                    <SelectItem key={status.value} value={status.value}>
                                        {status.label}
                                    </SelectItem>
                                ))}
                            </SelectContent>
                        </Select>

                        <Button onClick={handleSearch} variant="outline">
                            <Filter className="w-4 h-4 mr-2" />
                            Filter
                        </Button>
                    </div>
                </CardContent>
            </Card>

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

            {/* Groups List */}
            {loading ? (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                    {[...Array(6)].map((_, i) => (
                        <Card key={i} className="animate-pulse">
                            <CardContent className="p-6">
                                <div className="h-4 bg-gray-200 rounded mb-2"></div>
                                <div className="h-3 bg-gray-200 rounded mb-4"></div>
                                <div className="flex gap-2">
                                    <div className="h-6 bg-gray-200 rounded w-16"></div>
                                    <div className="h-6 bg-gray-200 rounded w-20"></div>
                                </div>
                            </CardContent>
                        </Card>
                    ))}
                </div>
            ) : groups.length === 0 ? (
                <Card>
                    <CardContent className="p-8 text-center">
                        <Users className="w-12 h-12 text-gray-400 mx-auto mb-4" />
                        <h3 className="text-lg font-medium text-gray-900 mb-2">No groups found</h3>
                        <p className="text-gray-500 mb-4">
                            {searchQuery || filters.type || filters.status
                                ? 'Try adjusting your search criteria'
                                : 'Get started by creating your first group'
                            }
                        </p>
                        {!searchQuery && !filters.type && !filters.status && (
                            <Button onClick={() => setShowForm(true)}>
                                <Plus className="w-4 h-4 mr-2" />
                                Create Group
                            </Button>
                        )}
                    </CardContent>
                </Card>
            ) : (
                <>
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                        {groups.map((group) => (
                            <Card key={group.id} className="hover:shadow-md transition-shadow">
                                <CardHeader className="pb-3">
                                    <div className="flex justify-between items-start">
                                        <CardTitle className="text-lg">{group.name}</CardTitle>
                                        <div className="flex gap-1">
                                            <Button
                                                size="sm"
                                                variant="ghost"
                                                onClick={() => setViewingGroup(group)}
                                            >
                                                <Eye className="w-4 h-4" />
                                            </Button>
                                            <Button
                                                size="sm"
                                                variant="ghost"
                                                onClick={() => setEditingGroup(group)}
                                            >
                                                <Edit className="w-4 h-4" />
                                            </Button>
                                            <Button
                                                size="sm"
                                                variant="ghost"
                                                onClick={() => handleDeleteGroup(group)}
                                            >
                                                <Trash2 className="w-4 h-4" />
                                            </Button>
                                        </div>
                                    </div>
                                </CardHeader>
                                <CardContent>
                                    {group.description && (
                                        <p className="text-sm text-gray-600 mb-3 line-clamp-2">
                                            {group.description}
                                        </p>
                                    )}

                                    <div className="flex flex-wrap gap-2 mb-3">
                                        <Badge className={getTypeColor(group.type)}>
                                            {GROUP_TYPES.find(t => t.value === group.type)?.label}
                                        </Badge>
                                        <Badge className={getStatusColor(group.status)}>
                                            {group.status}
                                        </Badge>
                                    </div>

                                    <div className="flex items-center justify-between text-sm text-gray-500">
                                        <span>Created by {group.creator.display_name}</span>
                                        <span>{new Date(group.created_at).toLocaleDateString()}</span>
                                    </div>

                                    <Button
                                        className="w-full mt-3"
                                        variant="outline"
                                        onClick={() => onGroupSelect?.(group)}
                                    >
                                        View Details
                                    </Button>
                                </CardContent>
                            </Card>
                        ))}
                    </div>

                    {/* Pagination */}
                    {pagination && pagination.total_pages > 1 && (
                        <div className="flex justify-center items-center gap-2">
                            <Button
                                variant="outline"
                                size="sm"
                                onClick={() => handlePageChange(pagination.page - 1)}
                                disabled={pagination.page <= 1}
                            >
                                Previous
                            </Button>

                            <span className="text-sm text-gray-600">
                                Page {pagination.page} of {pagination.total_pages}
                            </span>

                            <Button
                                variant="outline"
                                size="sm"
                                onClick={() => handlePageChange(pagination.page + 1)}
                                disabled={pagination.page >= pagination.total_pages}
                            >
                                Next
                            </Button>
                        </div>
                    )}
                </>
            )}
        </div>
    );
}
