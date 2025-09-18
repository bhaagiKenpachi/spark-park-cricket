'use client';

import { useState, useEffect } from 'react';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import {
  createMatchRequest,
  updateMatchRequest,
  Match,
} from '@/store/reducers/matchSlice';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Save, X, Calendar, Users, Target, Trophy, Coins } from 'lucide-react';

interface MatchFormProps {
  match?: Match | undefined;
  seriesId?: string;
  onSuccess?: () => void;
  onCancel?: () => void;
}

interface FormData {
  series_id: string;
  match_number: number;
  date: string;
  team_player_count: number;
  total_overs: number;
  toss_winner: 'A' | 'B';
  toss_type: 'H' | 'T';
}

interface FormErrors {
  series_id?: string;
  match_number?: string;
  date?: string;
  team_player_count?: string;
  total_overs?: string;
  toss_winner?: string;
  toss_type?: string;
}

export function MatchForm({
  match,
  seriesId,
  onSuccess,
  onCancel,
}: MatchFormProps): React.JSX.Element {
  const dispatch = useAppDispatch();
  const { loading, error } = useAppSelector(state => state.match);

  const [formData, setFormData] = useState<FormData>({
    series_id: seriesId || match?.series_id || '',
    match_number: match?.match_number || 0, // 0 means auto-generate
    date:
      (match?.date
        ? match.date.split('T')[0]
        : new Date().toISOString().split('T')[0]) || '',
    team_player_count: match?.team_a_player_count || 0,
    total_overs: match?.total_overs ?? 0,
    toss_winner: match?.toss_winner || 'A',
    toss_type: match?.toss_type || 'H',
  });

  const [formErrors, setFormErrors] = useState<FormErrors>({});

  useEffect(() => {
    if (match) {
      setFormData({
        series_id: match.series_id,
        match_number: match.match_number,
        date: match.date.split('T')[0] || '',
        team_player_count: match.team_a_player_count,
        total_overs: match.total_overs,
        toss_winner: match.toss_winner,
        toss_type: match.toss_type,
      });
    }
  }, [match]);

  const validateForm = (): boolean => {
    const errors: FormErrors = {};

    if (!formData.series_id.trim()) {
      errors.series_id = 'Series ID is required';
    }

    if (!formData.date) {
      errors.date = 'Date is required';
    }

    if (formData.team_player_count === 0) {
      errors.team_player_count = 'Team player count is required';
    } else if (
      formData.team_player_count < 1 ||
      formData.team_player_count > 11
    ) {
      errors.team_player_count = 'Team player count must be between 1 and 11';
    }

    if (formData.total_overs === 0) {
      errors.total_overs = 'Total overs is required';
    } else if (formData.total_overs < 1 || formData.total_overs > 20) {
      errors.total_overs = 'Total overs must be between 1 and 20';
    }

    if (!formData.toss_winner) {
      errors.toss_winner = 'Toss winner is required';
    }

    if (!formData.toss_type) {
      errors.toss_type = 'Toss type is required';
    }

    setFormErrors(errors);
    return Object.keys(errors).length === 0;
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    if (!validateForm()) {
      return;
    }

    // Convert date string to RFC3339 format for the API
    const apiData: Omit<Match, 'id' | 'created_at' | 'updated_at'> = {
      series_id: formData.series_id,
      match_number: formData.match_number > 0 ? formData.match_number : 1, // Use provided number or default to 1
      date: `${formData.date}T00:00:00Z`,
      status: 'live' as const,
      team_a_player_count: formData.team_player_count,
      team_b_player_count: formData.team_player_count,
      total_overs: formData.total_overs,
      toss_winner: formData.toss_winner,
      toss_type: formData.toss_type,
      batting_team: formData.toss_winner, // Default batting team is toss winner
    };

    console.log('Form data being submitted:', formData);
    console.log('API data being sent:', apiData);

    if (match) {
      dispatch(
        updateMatchRequest({
          id: match.id,
          matchData: apiData,
        })
      );
    } else {
      dispatch(createMatchRequest(apiData));
    }

    if (onSuccess) {
      onSuccess();
    }
  };

  const handleInputChange = (field: keyof FormData, value: string | number) => {
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
            {match ? 'Edit Match' : 'Create New Match'}
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
              <Label htmlFor="date" className="flex items-center">
                <Calendar className="h-4 w-4 mr-2" />
                Date *
              </Label>
              <Input
                type="date"
                id="date"
                value={formData.date}
                onChange={e => handleInputChange('date', e.target.value)}
                data-cy="match-date"
                className={formErrors.date ? 'border-red-500' : ''}
              />
              {formErrors.date && (
                <p className="text-sm text-red-600">{formErrors.date}</p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="team_player_count" className="flex items-center">
                <Users className="h-4 w-4 mr-2" />
                Team Player Count *
              </Label>
              <Input
                type="number"
                id="team_player_count"
                value={formData.team_player_count || ''}
                onChange={e =>
                  handleInputChange(
                    'team_player_count',
                    e.target.value === '' ? 0 : parseInt(e.target.value) || 0
                  )
                }
                min="1"
                max="11"
                data-cy="team-player-count"
                className={formErrors.team_player_count ? 'border-red-500' : ''}
              />
              {formErrors.team_player_count && (
                <p className="text-sm text-red-600">
                  {formErrors.team_player_count}
                </p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="total_overs" className="flex items-center">
                <Target className="h-4 w-4 mr-2" />
                Total Overs *
              </Label>
              <Input
                type="number"
                id="total_overs"
                value={formData.total_overs || ''}
                onChange={e =>
                  handleInputChange(
                    'total_overs',
                    e.target.value === '' ? 0 : parseInt(e.target.value) || 0
                  )
                }
                min="1"
                max="20"
                data-cy="total-overs"
                className={formErrors.total_overs ? 'border-red-500' : ''}
              />
              {formErrors.total_overs && (
                <p className="text-sm text-red-600">{formErrors.total_overs}</p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="toss_winner" className="flex items-center">
                <Trophy className="h-4 w-4 mr-2" />
                Toss Winner *
              </Label>
              <Select
                value={formData.toss_winner}
                onValueChange={value =>
                  handleInputChange('toss_winner', value as 'A' | 'B')
                }
              >
                <SelectTrigger>
                  <SelectValue placeholder="Select toss winner" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="A">Team A</SelectItem>
                  <SelectItem value="B">Team B</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-2">
              <Label htmlFor="toss_type" className="flex items-center">
                <Coins className="h-4 w-4 mr-2" />
                Toss Type *
              </Label>
              <Select
                value={formData.toss_type}
                onValueChange={value =>
                  handleInputChange('toss_type', value as 'H' | 'T')
                }
              >
                <SelectTrigger>
                  <SelectValue placeholder="Select toss type" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="H">Heads</SelectItem>
                  <SelectItem value="T">Tails</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-2">
              <Label htmlFor="match_number">Match Number (Optional)</Label>
              <Input
                type="number"
                id="match_number"
                value={formData.match_number || ''}
                onChange={e =>
                  handleInputChange(
                    'match_number',
                    e.target.value ? parseInt(e.target.value) : 0
                  )
                }
                min="1"
                placeholder="Leave empty for auto-generation"
                data-cy="match-number"
              />
            </div>

            <div className="flex flex-col space-y-3 pt-4 sm:flex-row sm:space-y-0 sm:space-x-3 sm:justify-center">
              <Button
                type="submit"
                disabled={loading}
                className="w-full sm:w-auto"
                data-cy={match ? 'update-match-button' : 'create-match-button'}
                title={
                  loading
                    ? 'Saving...'
                    : match
                      ? 'Update Match'
                      : 'Create Match'
                }
              >
                <Save className="h-4 w-4 mr-2" />
                {loading ? 'Saving...' : 'Match'}
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
