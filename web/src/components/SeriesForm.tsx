'use client';

import { useState, useEffect } from 'react';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import {
    createSeriesRequest,
    updateSeriesRequest,
    Series
} from '@/store/reducers/seriesSlice';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Calendar, Clock, Save, X } from 'lucide-react';

interface SeriesFormProps {
    series?: Series | undefined;
    onSuccess?: () => void;
    onCancel?: () => void;
}

interface FormData {
    name: string;
    start_date: string;
    end_date: string;
}

export function SeriesForm({ series, onSuccess, onCancel }: SeriesFormProps): React.JSX.Element {
    const dispatch = useAppDispatch();
    const { loading, error } = useAppSelector((state) => state.series);

    const [formData, setFormData] = useState<FormData>({
        name: series?.name || '',
        start_date: series?.start_date ? series.start_date.split('T')[0] : '',
        end_date: series?.end_date ? series.end_date.split('T')[0] : '',
    });

    const [formErrors, setFormErrors] = useState<Partial<FormData>>({});

    useEffect(() => {
        if (series) {
            // Convert RFC3339 dates to YYYY-MM-DD format for HTML date inputs
            const formatDateForInput = (dateString: string) => {
                return dateString.split('T')[0];
            };

            setFormData({
                name: series.name,
                start_date: formatDateForInput(series.start_date),
                end_date: formatDateForInput(series.end_date),
            });
        }
    }, [series]);

    const validateForm = (): boolean => {
        const errors: Partial<FormData> = {};

        if (!formData.name.trim()) {
            errors.name = 'Name is required';
        }

        if (!formData.start_date) {
            errors.start_date = 'Start date is required';
        }

        if (!formData.end_date) {
            errors.end_date = 'End date is required';
        }

        if (formData.start_date && formData.end_date && formData.start_date >= formData.end_date) {
            errors.end_date = 'End date must be after start date';
        }

        setFormErrors(errors);
        return Object.keys(errors).length === 0;
    };

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();

        if (!validateForm()) {
            return;
        }

        // Convert date strings to RFC3339 format for the API
        const apiData = {
            ...formData,
            start_date: `${formData.start_date}T00:00:00Z`,
            end_date: `${formData.end_date}T00:00:00Z`,
        };

        if (series) {
            dispatch(updateSeriesRequest({
                id: series.id,
                seriesData: apiData,
            }));
        } else {
            dispatch(createSeriesRequest(apiData));
        }

        if (onSuccess) {
            onSuccess();
        }
    };

    const handleInputChange = (field: keyof FormData, value: string) => {
        setFormData(prev => ({ ...prev, [field]: value }));
        if (formErrors[field]) {
            setFormErrors(prev => ({ ...prev, [field]: undefined }));
        }
    };

    return (
        <div className="w-full max-w-md mx-auto p-6">
            <Card>
                <CardHeader>
                    <CardTitle className="text-center">
                        {series ? 'Edit Series' : 'Create New Series'}
                    </CardTitle>
                </CardHeader>
                <CardContent>
                    {error && (
                        <div className="mb-4 p-3 bg-red-100 border border-red-400 text-red-700 rounded" data-cy="error-message">
                            {error}
                        </div>
                    )}

                    <form onSubmit={handleSubmit} className="space-y-4">
                        <div className="space-y-2">
                            <Label htmlFor="name">Series Name *</Label>
                            <Input
                                type="text"
                                id="name"
                                value={formData.name}
                                onChange={(e) => handleInputChange('name', e.target.value)}
                                placeholder="Enter series name"
                                data-cy="series-name"
                                className={formErrors.name ? 'border-red-500' : ''}
                            />
                            {formErrors.name && (
                                <p className="text-sm text-red-600">{formErrors.name}</p>
                            )}
                        </div>


                        <div className="space-y-2">
                            <Label htmlFor="start_date" className="flex items-center">
                                <Calendar className="h-4 w-4 mr-2" />
                                Start Date *
                            </Label>
                            <Input
                                type="date"
                                id="start_date"
                                value={formData.start_date}
                                onChange={(e) => handleInputChange('start_date', e.target.value)}
                                data-cy="start-date"
                                className={formErrors.start_date ? 'border-red-500' : ''}
                            />
                            {formErrors.start_date && (
                                <p className="text-sm text-red-600">{formErrors.start_date}</p>
                            )}
                        </div>

                        <div className="space-y-2">
                            <Label htmlFor="end_date" className="flex items-center">
                                <Calendar className="h-4 w-4 mr-2" />
                                End Date *
                            </Label>
                            <Input
                                type="date"
                                id="end_date"
                                value={formData.end_date}
                                onChange={(e) => handleInputChange('end_date', e.target.value)}
                                data-cy="end-date"
                                className={formErrors.end_date ? 'border-red-500' : ''}
                            />
                            {formErrors.end_date && (
                                <p className="text-sm text-red-600">{formErrors.end_date}</p>
                            )}
                        </div>


                        <div className="flex flex-col space-y-3 pt-4 sm:flex-row sm:space-y-0 sm:space-x-3 sm:justify-center">
                            <Button
                                type="submit"
                                disabled={loading}
                                className="w-full sm:w-auto"
                                data-cy={series ? 'update-series-button' : 'create-series-button'}
                                title={loading ? 'Saving...' : (series ? 'Update Series' : 'Create Series')}
                            >
                                <Save className="h-4 w-4 mr-2" />
                                {loading ? 'Saving...' : 'Series'}
                            </Button>

                            {onCancel && (
                                <Button
                                    type="button"
                                    variant="outline"
                                    onClick={onCancel}
                                    className="w-full sm:w-auto"
                                    title="Cancel"
                                >
                                    <X className="h-4 w-4" />
                                </Button>
                            )}
                        </div>
                    </form>
                </CardContent>
            </Card>
        </div>
    );
}
