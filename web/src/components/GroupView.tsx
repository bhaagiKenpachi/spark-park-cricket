'use client';

import { GroupWithCreator } from '@/types/group';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import {
    ArrowLeft,
    Edit,
    Calendar,
} from 'lucide-react';

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

export function GroupView({ group, onBack, onEdit }: GroupViewProps): React.JSX.Element {
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
                <Button variant="outline" onClick={() => onEdit(group)}>
                    <Edit className="w-4 h-4 mr-2" />
                    Edit
                </Button>
            </div>

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
                </CardContent>
            </Card>
        </div>
    );
}
