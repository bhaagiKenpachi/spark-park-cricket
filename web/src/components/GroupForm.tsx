'use client';

import { useState, useEffect } from 'react';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import {
    createGroupRequest,
    updateGroupRequest,
    clearError,
} from '@/store/reducers/groupSlice';
import { Group, GroupType, CreateGroupRequest, UpdateGroupRequest } from '@/types/group';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Save, X } from 'lucide-react';

interface GroupFormProps {
    group?: Group | undefined;
    onSuccess?: () => void;
    onCancel?: () => void;
}

interface FormData {
    name: string;
    description: string;
    type: GroupType;
}

interface FormErrors {
    name?: string;
    description?: string;
    type?: string;
}

const GROUP_TYPES: { value: GroupType; label: string }[] = [
    { value: 'custom', label: 'Custom' },
    { value: 'team', label: 'Team' },
    { value: 'series', label: 'Series' },
    { value: 'match', label: 'Match' },
    { value: 'location', label: 'Location' },
    { value: 'skill', label: 'Skill Level' },
];

export function GroupForm({
    group,
    onSuccess,
    onCancel,
}: GroupFormProps): React.JSX.Element {
    const dispatch = useAppDispatch();
    const { loading, error } = useAppSelector(state => state.group);

    const [formData, setFormData] = useState<FormData>({
        name: group?.name || '',
        description: group?.description || '',
        type: group?.type || 'custom',
    });

    const [errors, setErrors] = useState<FormErrors>({});

    useEffect(() => {
        if (error) {
            setErrors({ name: error });
        }
    }, [error]);

    const validateForm = (): boolean => {
        const newErrors: FormErrors = {};

        if (!formData.name.trim()) {
            newErrors.name = 'Group name is required';
        } else if (formData.name.trim().length < 3) {
            newErrors.name = 'Group name must be at least 3 characters';
        } else if (formData.name.trim().length > 100) {
            newErrors.name = 'Group name must be less than 100 characters';
        }

        if (formData.description && formData.description.length > 500) {
            newErrors.description = 'Description must be less than 500 characters';
        }

        if (!formData.type) {
            newErrors.type = 'Group type is required';
        }

        setErrors(newErrors);
        return Object.keys(newErrors).length === 0;
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        if (!validateForm()) {
            return;
        }

        dispatch(clearError());

        try {
            if (group) {
                // Update existing group
                const updateData: UpdateGroupRequest = {
                    name: formData.name.trim(),
                    ...(formData.description.trim() && { description: formData.description.trim() }),
                };
                await dispatch(updateGroupRequest({ id: group.id, groupData: updateData }));
            } else {
                // Create new group
                const createData: CreateGroupRequest = {
                    name: formData.name.trim(),
                    ...(formData.description.trim() && { description: formData.description.trim() }),
                    type: formData.type,
                };
                await dispatch(createGroupRequest(createData));
            }

            onSuccess?.();
        } catch (err) {
            console.error('Failed to save group:', err);
        }
    };

    const handleCancel = () => {
        dispatch(clearError());
        onCancel?.();
    };

    const handleInputChange = (field: keyof FormData, value: string) => {
        setFormData(prev => ({ ...prev, [field]: value }));

        // Clear error for this field when user starts typing
        if (errors[field as keyof FormErrors]) {
            setErrors(prev => ({ ...prev, [field]: undefined }));
        }
    };

    return (
        <Card className="w-full max-w-2xl mx-auto">
            <CardHeader>
                <CardTitle className="flex items-center gap-2">
                    {group ? 'Edit Group' : 'Create New Group'}
                </CardTitle>
            </CardHeader>
            <CardContent>
                <form onSubmit={handleSubmit} className="space-y-6">
                    {/* Group Name */}
                    <div className="space-y-2">
                        <Label htmlFor="name">Group Name *</Label>
                        <Input
                            id="name"
                            type="text"
                            value={formData.name}
                            onChange={(e) => handleInputChange('name', e.target.value)}
                            placeholder="Enter group name"
                            className={errors.name ? 'border-red-500' : ''}
                            disabled={loading}
                        />
                        {errors.name && (
                            <p className="text-sm text-red-500">{errors.name}</p>
                        )}
                    </div>

                    {/* Group Type */}
                    <div className="space-y-2">
                        <Label htmlFor="type">Group Type *</Label>
                        <Select
                            value={formData.type}
                            onValueChange={(value) => handleInputChange('type', value)}
                            disabled={loading || !!group}
                        >
                            <SelectTrigger className={errors.type ? 'border-red-500' : ''}>
                                <SelectValue placeholder="Select group type" />
                            </SelectTrigger>
                            <SelectContent>
                                {GROUP_TYPES.map((type) => (
                                    <SelectItem key={type.value} value={type.value}>
                                        {type.label}
                                    </SelectItem>
                                ))}
                            </SelectContent>
                        </Select>
                        {errors.type && (
                            <p className="text-sm text-red-500">{errors.type}</p>
                        )}
                        {group && (
                            <p className="text-sm text-gray-500">
                                Group type cannot be changed after creation
                            </p>
                        )}
                    </div>

                    {/* Description */}
                    <div className="space-y-2">
                        <Label htmlFor="description">Description</Label>
                        <Textarea
                            id="description"
                            value={formData.description}
                            onChange={(e) => handleInputChange('description', e.target.value)}
                            placeholder="Enter group description (optional)"
                            className={errors.description ? 'border-red-500' : ''}
                            disabled={loading}
                            rows={3}
                        />
                        {errors.description && (
                            <p className="text-sm text-red-500">{errors.description}</p>
                        )}
                        <p className="text-sm text-gray-500">
                            {formData.description.length}/500 characters
                        </p>
                    </div>

                    {/* Error Message */}
                    {error && !Object.values(errors).some(e => e) && (
                        <div className="p-3 bg-red-50 border border-red-200 rounded-md">
                            <p className="text-sm text-red-600">{error}</p>
                        </div>
                    )}

                    {/* Action Buttons */}
                    <div className="flex justify-end gap-3">
                        <Button
                            type="button"
                            variant="outline"
                            onClick={handleCancel}
                            disabled={loading}
                        >
                            <X className="w-4 h-4 mr-2" />
                            Cancel
                        </Button>
                        <Button
                            type="submit"
                            disabled={loading}
                        >
                            <Save className="w-4 h-4 mr-2" />
                            {loading ? 'Saving...' : (group ? 'Update Group' : 'Create Group')}
                        </Button>
                    </div>
                </form>
            </CardContent>
        </Card>
    );
}
