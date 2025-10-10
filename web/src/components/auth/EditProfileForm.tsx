'use client';

import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import { updateUserName } from '@/store/reducers/authSlice';
import { X, Check, Loader2 } from 'lucide-react';

interface EditProfileFormProps {
    onCancel: () => void;
    onSuccess: () => void;
}

export default function EditProfileForm({ onCancel, onSuccess }: EditProfileFormProps): React.JSX.Element {
    const dispatch = useAppDispatch();
    const { user, isLoading, error } = useAppSelector(state => state.auth);
    const [name, setName] = useState(user?.name || '');
    const [validationError, setValidationError] = useState('');

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        // Validation
        const trimmedName = name.trim();
        if (!trimmedName) {
            setValidationError('Name cannot be empty');
            return;
        }
        if (trimmedName.length < 2) {
            setValidationError('Name must be at least 2 characters long');
            return;
        }
        if (trimmedName.length > 255) {
            setValidationError('Name must be at most 255 characters long');
            return;
        }

        setValidationError('');

        try {
            await dispatch(updateUserName(trimmedName)).unwrap();
            onSuccess();
        } catch (err) {
            // Error is handled by Redux
            console.error('Failed to update name:', err);
        }
    };

    const handleCancel = () => {
        setName(user?.name || '');
        setValidationError('');
        onCancel();
    };

    return (
        <Card className="w-full">
            <CardHeader className="pb-3">
                <CardTitle className="text-base font-bold">Edit Profile</CardTitle>
            </CardHeader>
            <CardContent>
                <form onSubmit={handleSubmit} className="space-y-4">
                    <div className="space-y-2">
                        <Label htmlFor="name" className="text-sm font-medium">
                            Display Name
                        </Label>
                        <Input
                            id="name"
                            type="text"
                            value={name}
                            onChange={(e) => setName(e.target.value)}
                            placeholder="Enter your name"
                            className="w-full"
                            disabled={isLoading}
                            maxLength={255}
                        />
                        {validationError && (
                            <p className="text-xs text-red-600">{validationError}</p>
                        )}
                        {error && (
                            <p className="text-xs text-red-600">{error}</p>
                        )}
                        <p className="text-xs text-gray-500">
                            Your name will be displayed publicly across the platform
                        </p>
                    </div>

                    <div className="flex gap-2 pt-2">
                        <Button
                            type="submit"
                            disabled={isLoading || !name.trim() || name.trim() === user?.name}
                            className="flex-1 flex items-center gap-2 bg-blue-600 hover:bg-blue-700"
                        >
                            {isLoading ? (
                                <>
                                    <Loader2 className="h-4 w-4 animate-spin" />
                                    Saving...
                                </>
                            ) : (
                                <>
                                    <Check className="h-4 w-4" />
                                    Save Changes
                                </>
                            )}
                        </Button>
                        <Button
                            type="button"
                            variant="outline"
                            onClick={handleCancel}
                            disabled={isLoading}
                            className="flex-1 flex items-center gap-2"
                        >
                            <X className="h-4 w-4" />
                            Cancel
                        </Button>
                    </div>
                </form>
            </CardContent>
        </Card>
    );
}

