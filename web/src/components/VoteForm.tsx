'use client';

import { useState, useEffect } from 'react';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import {
    createVoteRequest,
    updateVoteRequest,
    fetchVoteWithResultsRequest,
} from '@/store/reducers/voteSlice';
import { Vote, VoteType } from '@/types/vote';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Save, X, Plus, Trash2, Users } from 'lucide-react';
import { GroupVoting } from './GroupVoting';
import { Checkbox } from '@/components/ui/checkbox';
import { groupService } from '@/services/groupService';
import voteService from '@/services/voteService';
import { createVoteSuccess, createVoteFailure } from '@/store/reducers/voteSlice';
import { assignGroupsToVoteRequest } from '@/store/reducers/groupSlice';

interface VoteFormProps {
    vote?: Vote | undefined;
    onSuccess?: () => void;
    onCancel?: () => void;
}

interface FormData {
    title: string;
    description: string;
    type: VoteType;
    options: string[];
    team_formation_enabled: boolean;
}

interface FormErrors {
    title?: string;
    description?: string;
    options?: string;
}

export function VoteForm({
    vote,
    onSuccess,
    onCancel,
}: VoteFormProps): React.JSX.Element {
    const dispatch = useAppDispatch();
    const { loading, error, currentVote } = useAppSelector(state => state.vote);

    const [formData, setFormData] = useState<FormData>({
        title: vote?.title || '',
        description: vote?.description || '',
        type: vote?.type || 'single',
        options: ['', ''], // Start with 2 empty options
        team_formation_enabled: vote?.team_formation_enabled ?? true,
    });

    const [showGroupAssignment, setShowGroupAssignment] = useState(!!vote);
    const [availableGroups, setAvailableGroups] = useState<any[]>([]);
    const [selectedGroupId, setSelectedGroupId] = useState<string>('');
    const [groupsLoading, setGroupsLoading] = useState(false);

    const [formErrors, setFormErrors] = useState<FormErrors>({});

    // Fetch vote with options when editing
    useEffect(() => {
        if (vote?.id) {
            dispatch(fetchVoteWithResultsRequest(vote.id));
        }
    }, [vote?.id, dispatch]);

    // Update form when vote data is loaded
    useEffect(() => {
        if (vote) {
            const voteOptions = currentVote?.options.map(opt => opt.text) || ['', ''];
            setFormData({
                title: vote.title,
                description: vote.description || '',
                type: vote.type,
                options: voteOptions.length >= 2 ? voteOptions : ['', ''],
                team_formation_enabled: vote.team_formation_enabled,
            });
        }
    }, [vote, currentVote]);

    // Load groups for assignment in create mode
    useEffect(() => {
        const loadGroups = async () => {
            setGroupsLoading(true);
            try {
                const res = await groupService.getGroups({ limit: 50, offset: 0 });
                const groupsData = (res.data as any).data || res.data;
                setAvailableGroups(Array.isArray(groupsData) ? groupsData : []);
            } catch (e) {
                // ignore load errors silently for now
                setAvailableGroups([]);
            } finally {
                setGroupsLoading(false);
            }
        };
        if (!vote) {
            loadGroups();
        }
    }, [vote]);

    const validateForm = (): boolean => {
        const errors: FormErrors = {};

        if (!formData.title.trim()) {
            errors.title = 'Title is required';
        } else if (formData.title.length < 3) {
            errors.title = 'Title must be at least 3 characters';
        } else if (formData.title.length > 255) {
            errors.title = 'Title must be less than 255 characters';
        }

        if (!formData.description.trim()) {
            errors.description = 'Description is required';
        } else if (formData.description.length < 10) {
            errors.description = 'Description must be at least 10 characters';
        } else if (formData.description.length > 1000) {
            errors.description = 'Description must be less than 1000 characters';
        }

        // Only validate options when creating (not editing)
        if (!vote) {
            const validOptions = formData.options.filter(option => option.trim() !== '');
            if (validOptions.length < 2) {
                errors.options = 'At least 2 options are required';
            }

            // Check for duplicate options
            const uniqueOptions = new Set(validOptions.map(option => option.toLowerCase()));
            if (uniqueOptions.size !== validOptions.length) {
                errors.options = 'Options must be unique';
            }
        }

        setFormErrors(errors);
        return Object.keys(errors).length === 0;
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        if (!validateForm()) {
            return;
        }

        const validOptions = formData.options.filter(option => option.trim() !== '');
        const apiData = {
            title: formData.title.trim(),
            description: formData.description.trim() || undefined,
            type: formData.type,
            options: validOptions,
            team_formation_enabled: formData.team_formation_enabled,
        };

        if (vote) {
            // When editing, only update title, description, and team formation setting (not options)
            const updateData: { title?: string; description?: string; team_formation_enabled?: boolean } = {
                title: apiData.title,
                team_formation_enabled: apiData.team_formation_enabled,
            };
            if (apiData.description) {
                updateData.description = apiData.description;
            }
            dispatch(
                updateVoteRequest({
                    id: vote.id,
                    voteData: updateData,
                })
            );
        } else {
            // Create and then assign groups if selected
            const createData = {
                title: apiData.title,
                description: apiData.description || 'No description provided',
                type: apiData.type,
                options: apiData.options,
                team_formation_enabled: apiData.team_formation_enabled,
            };
            try {
                const response = await voteService.createVote(createData as any);
                const created = (response.data as any).data || response.data;
                // Handle nested vote structure: {message: '...', vote: {...}}
                const voteObject = (created as any)?.vote || created;
                const createdVote = Array.isArray(voteObject) ? voteObject[0] : voteObject;
                if (createdVote?.id) {
                    dispatch(createVoteSuccess(createdVote));
                    // Only assign group if one is selected (not empty string or 'none')
                    const shouldAssign = selectedGroupId && selectedGroupId.trim() !== '' && selectedGroupId !== 'none';
                    if (shouldAssign) {
                        await dispatch(assignGroupsToVoteRequest({
                            voteId: createdVote.id,
                            groupIds: [selectedGroupId],
                        }));
                    }
                }
                if (onSuccess) {
                    onSuccess();
                }
            } catch (err: any) {
                dispatch(createVoteFailure(err?.message || 'Failed to create vote'));
            }
        }
    };

    const handleInputChange = (field: keyof FormData, value: string | VoteType | boolean) => {
        setFormData(prev => ({ ...prev, [field]: value }));
        if (field in formErrors) {
            setFormErrors(prev => {
                const newErrors = { ...prev };
                delete (newErrors as any)[field];
                return newErrors;
            });
        }
    };

    const handleOptionChange = (index: number, value: string) => {
        const newOptions = [...formData.options];
        newOptions[index] = value;
        setFormData(prev => ({ ...prev, options: newOptions }));
        if (formErrors.options) {
            setFormErrors(prev => {
                const newErrors = { ...prev };
                delete newErrors.options;
                return newErrors;
            });
        }
    };

    const addOption = () => {
        setFormData(prev => ({ ...prev, options: [...prev.options, ''] }));
    };

    const removeOption = (index: number) => {
        if (formData.options.length > 2) {
            const newOptions = formData.options.filter((_, i) => i !== index);
            setFormData(prev => ({ ...prev, options: newOptions }));
        }
    };

    return (
        <div className="w-full max-w-md mx-auto px-4 py-4">
            <div className="bg-white rounded-lg shadow-sm border">
                <div className="p-4 border-b">
                    <h2 className="text-lg font-semibold text-gray-900">
                        {vote ? 'Edit Vote' : 'Create New Vote'}
                    </h2>
                </div>
                <div className="p-4">
                    {error && (
                        <div
                            className="mb-4 p-3 bg-red-50 border border-red-200 text-red-700 rounded-lg text-sm"
                            data-cy="error-message"
                        >
                            {error}
                        </div>
                    )}

                    <form onSubmit={handleSubmit} className="space-y-4">
                        <div className="space-y-1">
                            <Label htmlFor="title" className="text-sm font-medium">Vote Title *</Label>
                            <Input
                                type="text"
                                id="title"
                                value={formData.title}
                                onChange={e => handleInputChange('title', e.target.value)}
                                placeholder="Enter vote title"
                                data-cy="vote-title"
                                className={`h-10 ${formErrors.title ? 'border-red-500' : ''}`}
                            />
                            {formErrors.title && (
                                <p className="text-xs text-red-600">{formErrors.title}</p>
                            )}
                        </div>

                        <div className="space-y-1">
                            <Label htmlFor="description" className="text-sm font-medium">Description *</Label>
                            <Textarea
                                id="description"
                                value={formData.description}
                                onChange={e => handleInputChange('description', e.target.value)}
                                placeholder="Enter vote description"
                                data-cy="vote-description"
                                className={formErrors.description ? 'border-red-500' : ''}
                                rows={2}
                            />
                            {formErrors.description && (
                                <p className="text-xs text-red-600">{formErrors.description}</p>
                            )}
                        </div>

                        <div className="space-y-1">
                            <Label htmlFor="type" className="text-sm font-medium">Vote Type *</Label>
                            <Select
                                value={formData.type}
                                onValueChange={(value: VoteType) => handleInputChange('type', value)}
                                disabled={!!vote}
                            >
                                <SelectTrigger data-cy="vote-type" className="h-10">
                                    <SelectValue placeholder="Select vote type" />
                                </SelectTrigger>
                                <SelectContent>
                                    <SelectItem value="single">Single Choice</SelectItem>
                                    <SelectItem value="multiple">Multiple Choice</SelectItem>
                                </SelectContent>
                            </Select>
                        </div>

                        <div className="space-y-1">
                            <div className="flex items-center justify-between">
                                <Label htmlFor="team-formation" className="text-sm font-medium">
                                    Enable Team Formation
                                </Label>
                                <input
                                    id="team-formation"
                                    type="checkbox"
                                    checked={formData.team_formation_enabled}
                                    onChange={(e) => handleInputChange('team_formation_enabled', e.target.checked)}
                                    data-cy="team-formation-toggle"
                                    className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                                />
                            </div>
                            <p className="text-xs text-gray-500">
                                Allow users to create and manage teams from voted users
                            </p>
                        </div>

                        <div className="space-y-3">
                            <div className="flex items-center justify-between">
                                <Label className="text-sm font-medium">
                                    Vote Options * {vote && <span className="text-xs text-gray-500">(read-only)</span>}
                                </Label>
                                {!vote && (
                                    <Button
                                        type="button"
                                        variant="outline"
                                        size="sm"
                                        onClick={addOption}
                                        className="flex items-center gap-1 text-xs px-2 py-1"
                                    >
                                        <Plus className="h-3 w-3" />
                                        Add
                                    </Button>
                                )}
                            </div>

                            {formData.options.map((option, index) => (
                                <div key={index} className="flex items-center gap-2">
                                    <Input
                                        type="text"
                                        value={option}
                                        onChange={e => handleOptionChange(index, e.target.value)}
                                        placeholder={`Option ${index + 1}`}
                                        data-cy={`vote-option-${index}`}
                                        className="flex-1 h-10"
                                        disabled={!!vote}
                                    />
                                    {!vote && formData.options.length > 2 && (
                                        <Button
                                            type="button"
                                            variant="outline"
                                            size="sm"
                                            onClick={() => removeOption(index)}
                                            className="text-red-600 hover:text-red-700 p-2"
                                        >
                                            <Trash2 className="h-4 w-4" />
                                        </Button>
                                    )}
                                </div>
                            ))}

                            {formErrors.options && (
                                <p className="text-xs text-red-600">{formErrors.options}</p>
                            )}
                        </div>

                        {/* Group Assignment Section */}
                        {vote ? (
                            <div className="space-y-3">
                                <div className="flex items-center justify-between">
                                    <Label className="text-sm font-medium">Assign Groups</Label>
                                    <Button
                                        type="button"
                                        variant="outline"
                                        size="sm"
                                        onClick={() => setShowGroupAssignment(!showGroupAssignment)}
                                    >
                                        <Users className="h-4 w-4 mr-2" />
                                        {showGroupAssignment ? 'Hide' : 'Manage Groups'}
                                    </Button>
                                </div>

                                {showGroupAssignment && (
                                    <div className="border rounded-lg p-4 bg-gray-50">
                                        <GroupVoting
                                            voteId={vote.id}
                                            onSuccess={() => {
                                                // Optionally refresh vote data
                                                dispatch(fetchVoteWithResultsRequest(vote.id));
                                            }}
                                        />
                                    </div>
                                )}
                            </div>
                        ) : (
                            <div className="space-y-2">
                                <Label className="text-sm font-medium flex items-center gap-2">
                                    <Users className="h-4 w-4" /> Assign to Group
                                </Label>
                                {groupsLoading ? (
                                    <div className="rounded-md bg-gray-50 border p-3 text-xs text-gray-600 flex items-center gap-2">
                                        <div className="animate-spin rounded-full h-3 w-3 border-b-2 border-gray-400"></div>
                                        Loading groups...
                                    </div>
                                ) : availableGroups.length === 0 ? (
                                    <div className="rounded-md bg-gray-50 border p-3 text-xs text-gray-600">
                                        No groups available.
                                    </div>
                                ) : (
                                    <Select
                                        value={selectedGroupId || 'none'}
                                        onValueChange={(value) => {
                                            const newValue = value === 'none' ? '' : value;
                                            setSelectedGroupId(newValue);
                                        }}
                                    >
                                        <SelectTrigger className="h-9 text-sm">
                                            <SelectValue placeholder="Select a group (optional)" />
                                        </SelectTrigger>
                                        <SelectContent>
                                            <SelectItem value="none">No group</SelectItem>
                                            {availableGroups.map((g: any) => (
                                                <SelectItem key={g.id} value={g.id}>
                                                    {g.name} ({g.type})
                                                </SelectItem>
                                            ))}
                                        </SelectContent>
                                    </Select>
                                )}
                                <p className="text-[11px] text-gray-500">The vote will be assigned to this group when created.</p>
                            </div>
                        )}

                        <div className="flex gap-2 pt-4">
                            <Button
                                type="submit"
                                disabled={loading}
                                className="flex-1"
                                data-cy={vote ? 'update-vote-button' : 'create-vote-button'}
                            >
                                <Save className="h-4 w-4 mr-2" />
                                {loading ? 'Saving...' : vote ? 'Update' : 'Create'}
                            </Button>

                            {onCancel && (
                                <Button
                                    type="button"
                                    variant="outline"
                                    onClick={onCancel}
                                    className="px-4"
                                >
                                    <X className="h-4 w-4" />
                                </Button>
                            )}
                        </div>
                    </form>
                </div>
            </div>
        </div>
    );
}
