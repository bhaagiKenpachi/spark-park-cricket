'use client';

import { useState, useEffect } from 'react';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import {
  createSeriesRequest,
  updateSeriesRequest,
  Series,
} from '@/store/reducers/seriesSlice';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Calendar, Save, X } from 'lucide-react';
import { trackSeriesCreated, trackSeriesEdited } from '@/lib/analytics';

interface SeriesFormProps {
  series?: Series | undefined;
  onSuccess?: () => void;
  onCancel?: () => void;
}

interface FormData {
  name: string;
  team_a_name: string;
  team_b_name: string;
  start_date: string;
  end_date: string;
}

export function SeriesForm({
  series,
  onSuccess,
  onCancel,
}: SeriesFormProps): React.JSX.Element {
  const dispatch = useAppDispatch();
  const { loading, error } = useAppSelector(state => state.series);

  // Helper function to format datetime for datetime-local input
  const formatDateTimeForInput = (dateString: string): string => {
    try {
      const date = new Date(dateString);
      const year = date.getFullYear();
      const month = String(date.getMonth() + 1).padStart(2, '0');
      const day = String(date.getDate()).padStart(2, '0');
      const hours = String(date.getHours()).padStart(2, '0');
      const minutes = String(date.getMinutes()).padStart(2, '0');
      return `${year}-${month}-${day}T${hours}:${minutes}`;
    } catch {
      return '';
    }
  };

  const getDefaultDateTime = (): string => {
    const now = new Date();
    return formatDateTimeForInput(now.toISOString());
  };

  const [formData, setFormData] = useState<FormData>({
    name: series?.name || '',
    team_a_name: series?.team_a_name || '',
    team_b_name: series?.team_b_name || '',
    start_date: series?.start_date
      ? formatDateTimeForInput(series.start_date)
      : getDefaultDateTime(),
    end_date: series?.end_date
      ? formatDateTimeForInput(series.end_date)
      : getDefaultDateTime(),
  });

  const [formErrors, setFormErrors] = useState<Partial<FormData>>({});

  useEffect(() => {
    if (series) {
      setFormData({
        name: series.name,
        team_a_name: series.team_a_name || '',
        team_b_name: series.team_b_name || '',
        start_date: formatDateTimeForInput(series.start_date),
        end_date: formatDateTimeForInput(series.end_date),
      });
    }
  }, [series]);

  const validateForm = (): boolean => {
    const errors: Partial<FormData> = {};

    if (!formData.name.trim()) {
      errors.name = 'Name is required';
    }

    if (!formData.start_date) {
      errors.start_date = 'Start date & time is required';
    }

    if (!formData.end_date) {
      errors.end_date = 'End date & time is required';
    }

    // Compare datetime values properly
    if (formData.start_date && formData.end_date) {
      const startDateTime = new Date(formData.start_date);
      const endDateTime = new Date(formData.end_date);

      // End datetime must be greater than or equal to start datetime
      if (endDateTime < startDateTime) {
        errors.end_date = 'End date & time cannot be earlier than start date & time';
      }
    }

    setFormErrors(errors);
    return Object.keys(errors).length === 0;
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    if (!validateForm()) {
      return;
    }

    // Convert datetime-local strings to RFC3339 format for the API
    // datetime-local format is YYYY-MM-DDTHH:mm
    const convertToRFC3339 = (dateTimeString: string): string => {
      // If already includes time, add seconds and Z
      if (dateTimeString.includes('T')) {
        return `${dateTimeString}:00Z`;
      }
      // If only date, add time
      return `${dateTimeString}T00:00:00Z`;
    };

    // Build the API data object
    const baseData = {
      name: formData.name,
      start_date: convertToRFC3339(formData.start_date),
      end_date: convertToRFC3339(formData.end_date),
      status: 'upcoming' as const,
    };

    // Add team names if they have values
    const apiData = {
      ...baseData,
      ...(formData.team_a_name && { team_a_name: formData.team_a_name }),
      ...(formData.team_b_name && { team_b_name: formData.team_b_name }),
    };

    if (series) {
      dispatch(
        updateSeriesRequest({
          id: series.id,
          seriesData: apiData,
        })
      );
      // Track series edited event
      trackSeriesEdited({
        series_id: series.id,
        series_name: apiData.name,
        start_date: apiData.start_date,
        end_date: apiData.end_date,
      });
    } else {
      dispatch(createSeriesRequest(apiData));
      // Track series created event
      trackSeriesCreated({
        series_id: '', // Will be set after creation
        series_name: apiData.name,
        start_date: apiData.start_date,
        end_date: apiData.end_date,
      });
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
            <div
              className="mb-4 p-3 bg-red-100 border border-red-400 text-red-700 rounded"
              data-cy="error-message"
            >
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
                onChange={e => handleInputChange('name', e.target.value)}
                placeholder="Enter series name"
                data-cy="series-name"
                className={formErrors.name ? 'border-red-500' : ''}
              />
              {formErrors.name && (
                <p className="text-sm text-red-600">{formErrors.name}</p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="team_a_name">Team A Name (Optional)</Label>
              <Input
                type="text"
                id="team_a_name"
                value={formData.team_a_name}
                onChange={e => handleInputChange('team_a_name', e.target.value)}
                placeholder="e.g., Mumbai Indians"
                data-cy="team-a-name"
                className={formErrors.team_a_name ? 'border-red-500' : ''}
              />
              {formErrors.team_a_name && (
                <p className="text-sm text-red-600">{formErrors.team_a_name}</p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="team_b_name">Team B Name (Optional)</Label>
              <Input
                type="text"
                id="team_b_name"
                value={formData.team_b_name}
                onChange={e => handleInputChange('team_b_name', e.target.value)}
                placeholder="e.g., Chennai Super Kings"
                data-cy="team-b-name"
                className={formErrors.team_b_name ? 'border-red-500' : ''}
              />
              {formErrors.team_b_name && (
                <p className="text-sm text-red-600">{formErrors.team_b_name}</p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="start_date" className="flex items-center">
                <Calendar className="h-4 w-4 mr-2" />
                Start Date & Time *
              </Label>
              <Input
                type="datetime-local"
                id="start_date"
                value={formData.start_date}
                onChange={e => handleInputChange('start_date', e.target.value)}
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
                End Date & Time *
              </Label>
              <Input
                type="datetime-local"
                id="end_date"
                value={formData.end_date}
                onChange={e => handleInputChange('end_date', e.target.value)}
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
                data-cy={
                  series ? 'update-series-button' : 'create-series-button'
                }
                title={
                  loading
                    ? 'Saving...'
                    : series
                      ? 'Update Series'
                      : 'Create Series'
                }
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
