'use client';

import { useState, useEffect } from 'react';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import {
    fetchGroupWithMembersRequest,
    fetchGroupStatsRequest,
    fetchGroupVotesRequest,
    joinGroupRequest,
    leaveGroupRequest,
    clearError,
} from '@/store/reducers/groupSlice';
import { GroupWithCreator, GroupMemberRole } from '@/types/group';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import {
    ArrowLeft,
    Edit,
    Users,
    Calendar,
    BarChart3,
    Vote,
    UserPlus,
    UserMinus,
    Crown,
    Shield,
    User
} from 'lucide-react';
import { GroupForm } from './GroupForm';

interface GroupViewProps {
    group: GroupWithCreator;
    onBack: () => void;
    onEdit: (group: GroupWithCreator) => void;
}

const GROUP_TYPES: Record<string, string> = {
    custom: 'Custom',
    team: 'Team',
    series: 'Series',
    match: 'Match',
    location: 'Location',
    skill: 'Skill Level',
};

const ROLE_ICONS: Record<GroupMemberRole, React.ComponentType<{ className?: string }>> = {
    admin: Crown,
    moderator: Shield,
    member: User,
};

const ROLE_COLORS: Record<GroupMemberRole, string> = {
    admin: 'bg-red-100 text-red-800',
    moderator: 'bg-blue-100 text-blue-800',
    member: 'bg-gray-100 text-gray-800',
};

export function GroupView({ group, onBack, onEdit }: GroupViewProps): React.JSX.Element {
    const dispatch = useAppDispatch();
    const {
        groupMembers,
        groupStats,
        groupVotes,
        loading,
        error
    } = useAppSelector(state => state.group);

    const [activeTab, setActiveTab] = useState<'overview' | 'members' | 'votes'>('overview');
    const [isEditing, setIsEditing] = useState(false);

    useEffect(() => {
        dispatch(fetchGroupWithMembersRequest(group.id));
        dispatch(fetchGroupStatsRequest(group.id));
        dispatch(fetchGroupVotesRequest(group.id));
    }, [dispatch, group.id]);

    const handleJoinGroup = async () => {
        await dispatch(joinGroupRequest(group.id));
        // Refresh members after joining
        dispatch(fetchGroupWithMembersRequest(group.id));
        dispatch(fetchGroupStatsRequest(group.id));
    };

    const handleLeaveGroup = async () => {
        await dispatch(leaveGroupRequest(group.id));
        // Refresh members after leaving
        dispatch(fetchGroupWithMembersRequest(group.id));
        dispatch(fetchGroupStatsRequest(group.id));
    };

    const getStatusColor = (status: string): string => {
        const colors = {
            active: 'bg-green-100 text-green-800',
            inactive: 'bg-yellow-100 text-yellow-800',
            archived: 'bg-gray-100 text-gray-800',
        };
        return colors[status as keyof typeof colors] || 'bg-gray-100 text-gray-800';
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

    const isUserMember = () => {
        // This would need to be implemented based on your auth system
        // For now, we'll assume the user is not a member
        return false;
    };

    const getUserRole = (): GroupMemberRole | null => {
        // This would need to be implemented based on your auth system
        // For now, we'll return null
        return null;
    };

    if (isEditing) {
        return (
            <GroupForm
                group={group}
                onSuccess={() => {
                    setIsEditing(false);
                    // Refresh group data
                    dispatch(fetchGroupWithMembersRequest(group.id));
                }}
                onCancel={() => setIsEditing(false)}
            />
        );
    }

    return (
        <div className="space-y-6">
            {/* Header */}
            <div className="flex items-center gap-4">
                <Button variant="outline" onClick={onBack}>
                    <ArrowLeft className="w-4 h-4 mr-2" />
                    Back
                </Button>
                <div className="flex-1">
                    <h1 className="text-2xl font-bold">{group.name}</h1>
                    <p className="text-gray-600">Group Details</p>
                </div>
                <Button variant="outline" onClick={() => setIsEditing(true)}>
                    <Edit className="w-4 h-4 mr-2" />
                    Edit
                </Button>
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

            {/* Group Info */}
            <Card>
                <CardHeader>
                    <div className="flex justify-between items-start">
                        <div>
                            <CardTitle className="text-xl">{group.name}</CardTitle>
                            <p className="text-gray-600 mt-1">Created by {group.creator.display_name}</p>
                        </div>
                        <div className="flex gap-2">
                            <Badge className={getTypeColor(group.type)}>
                                {GROUP_TYPES[group.type]}
                            </Badge>
                            <Badge className={getStatusColor(group.status)}>
                                {group.status}
                            </Badge>
                        </div>
                    </div>
                </CardHeader>
                <CardContent>
                    {group.description && (
                        <p className="text-gray-700 mb-4">{group.description}</p>
                    )}

                    <div className="flex items-center gap-4 text-sm text-gray-500">
                        <div className="flex items-center gap-1">
                            <Calendar className="w-4 h-4" />
                            <span>Created {new Date(group.created_at).toLocaleDateString()}</span>
                        </div>
                        <div className="flex items-center gap-1">
                            <Calendar className="w-4 h-4" />
                            <span>Updated {new Date(group.updated_at).toLocaleDateString()}</span>
                        </div>
                    </div>

                    {/* Action Buttons */}
                    <div className="flex gap-2 mt-4">
                        {isUserMember() ? (
                            <Button variant="outline" onClick={handleLeaveGroup}>
                                <UserMinus className="w-4 h-4 mr-2" />
                                Leave Group
                            </Button>
                        ) : (
                            <Button onClick={handleJoinGroup}>
                                <UserPlus className="w-4 h-4 mr-2" />
                                Join Group
                            </Button>
                        )}
                    </div>
                </CardContent>
            </Card>

            {/* Tabs */}
            <div className="flex gap-1 border-b">
                <Button
                    variant={activeTab === 'overview' ? 'default' : 'ghost'}
                    onClick={() => setActiveTab('overview')}
                >
                    <BarChart3 className="w-4 h-4 mr-2" />
                    Overview
                </Button>
                <Button
                    variant={activeTab === 'members' ? 'default' : 'ghost'}
                    onClick={() => setActiveTab('members')}
                >
                    <Users className="w-4 h-4 mr-2" />
                    Members
                </Button>
                <Button
                    variant={activeTab === 'votes' ? 'default' : 'ghost'}
                    onClick={() => setActiveTab('votes')}
                >
                    <Vote className="w-4 h-4 mr-2" />
                    Votes
                </Button>
            </div>

            {/* Tab Content */}
            {activeTab === 'overview' && (
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                    {groupStats && (
                        <>
                            <Card>
                                <CardContent className="p-6">
                                    <div className="flex items-center gap-3">
                                        <Users className="w-8 h-8 text-blue-600" />
                                        <div>
                                            <p className="text-2xl font-bold">{groupStats.member_count}</p>
                                            <p className="text-sm text-gray-600">Members</p>
                                        </div>
                                    </div>
                                </CardContent>
                            </Card>

                            <Card>
                                <CardContent className="p-6">
                                    <div className="flex items-center gap-3">
                                        <Vote className="w-8 h-8 text-green-600" />
                                        <div>
                                            <p className="text-2xl font-bold">{groupStats.vote_count}</p>
                                            <p className="text-sm text-gray-600">Votes</p>
                                        </div>
                                    </div>
                                </CardContent>
                            </Card>

                            <Card>
                                <CardContent className="p-6">
                                    <div className="flex items-center gap-3">
                                        <Crown className="w-8 h-8 text-yellow-600" />
                                        <div>
                                            <p className="text-2xl font-bold">{groupStats.admin_count}</p>
                                            <p className="text-sm text-gray-600">Admins</p>
                                        </div>
                                    </div>
                                </CardContent>
                            </Card>
                        </>
                    )}
                </div>
            )}

            {activeTab === 'members' && (
                <Card>
                    <CardHeader>
                        <CardTitle>Group Members</CardTitle>
                    </CardHeader>
                    <CardContent>
                        {loading ? (
                            <div className="space-y-3">
                                {[...Array(3)].map((_, i) => (
                                    <div key={i} className="flex items-center gap-3 animate-pulse">
                                        <div className="w-10 h-10 bg-gray-200 rounded-full"></div>
                                        <div className="flex-1">
                                            <div className="h-4 bg-gray-200 rounded mb-1"></div>
                                            <div className="h-3 bg-gray-200 rounded w-1/2"></div>
                                        </div>
                                        <div className="h-6 bg-gray-200 rounded w-16"></div>
                                    </div>
                                ))}
                            </div>
                        ) : groupMembers?.members.length === 0 ? (
                            <div className="text-center py-8">
                                <Users className="w-12 h-12 text-gray-400 mx-auto mb-4" />
                                <p className="text-gray-500">No members found</p>
                            </div>
                        ) : (
                            <div className="space-y-3">
                                {groupMembers?.members.map((member) => {
                                    const RoleIcon = ROLE_ICONS[member.role];
                                    return (
                                        <div key={member.id} className="flex items-center justify-between p-3 border rounded-lg">
                                            <div className="flex items-center gap-3">
                                                <div className="w-10 h-10 bg-gray-200 rounded-full flex items-center justify-center">
                                                    <User className="w-5 h-5 text-gray-600" />
                                                </div>
                                                <div>
                                                    <p className="font-medium">{member.user?.display_name || 'Unknown User'}</p>
                                                    <p className="text-sm text-gray-500">
                                                        Joined {new Date(member.joined_at).toLocaleDateString()}
                                                    </p>
                                                </div>
                                            </div>
                                            <Badge className={ROLE_COLORS[member.role]}>
                                                <RoleIcon className="w-3 h-3 mr-1" />
                                                {member.role}
                                            </Badge>
                                        </div>
                                    );
                                })}
                            </div>
                        )}
                    </CardContent>
                </Card>
            )}

            {activeTab === 'votes' && (
                <Card>
                    <CardHeader>
                        <CardTitle>Group Votes</CardTitle>
                    </CardHeader>
                    <CardContent>
                        {loading ? (
                            <div className="space-y-3">
                                {[...Array(3)].map((_, i) => (
                                    <div key={i} className="animate-pulse">
                                        <div className="h-4 bg-gray-200 rounded mb-2"></div>
                                        <div className="h-3 bg-gray-200 rounded w-3/4"></div>
                                    </div>
                                ))}
                            </div>
                        ) : groupVotes?.votes.length === 0 ? (
                            <div className="text-center py-8">
                                <Vote className="w-12 h-12 text-gray-400 mx-auto mb-4" />
                                <p className="text-gray-500">No votes assigned to this group</p>
                            </div>
                        ) : (
                            <div className="space-y-3">
                                {groupVotes?.votes.map((vote) => (
                                    <div key={vote.id} className="p-3 border rounded-lg">
                                        <h4 className="font-medium">{vote.title}</h4>
                                        {vote.description && (
                                            <p className="text-sm text-gray-600 mt-1">{vote.description}</p>
                                        )}
                                        <div className="flex items-center gap-2 mt-2">
                                            <Badge variant="outline">{vote.type}</Badge>
                                            <Badge className={getStatusColor(vote.status)}>
                                                {vote.status}
                                            </Badge>
                                            <span className="text-sm text-gray-500">
                                                {new Date(vote.created_at).toLocaleDateString()}
                                            </span>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        )}
                    </CardContent>
                </Card>
            )}
        </div>
    );
}
