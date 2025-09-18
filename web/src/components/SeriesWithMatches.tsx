'use client';

import { useState, useEffect } from 'react';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import {
  fetchMatchesRequest,
  deleteMatchRequest,
  Match,
} from '@/store/reducers/matchSlice';
import { MatchForm } from './MatchForm';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { RefreshCw, Plus, Edit, Trash2, Calendar, Play } from 'lucide-react';
import { Series } from '@/store/reducers/seriesSlice';

interface SeriesWithMatchesProps {
  series: Series;
  onEditSeries: (series: Series) => void;
  onDeleteSeries: (id: string) => void;
  onViewScorecard?: (matchId: string) => void;
}

export function SeriesWithMatches({
  series,
  onEditSeries,
  onDeleteSeries,
  onViewScorecard,
}: SeriesWithMatchesProps): React.JSX.Element {
  const dispatch = useAppDispatch();
  const {
    matches,
    loading: matchesLoading,
    error: matchesError,
  } = useAppSelector(state => state.match);
  const [showMatchForm, setShowMatchForm] = useState(false);
  const [editingMatch, setEditingMatch] = useState<Match | undefined>();
  const [expanded, setExpanded] = useState(false);

  // Format date to human readable format
  const formatDate = (dateString: string) => {
    try {
      const date = new Date(dateString);
      return date.toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
      });
    } catch {
      return dateString;
    }
  };

  // Filter matches for this series
  const seriesMatches = matches.filter(match => match.series_id === series.id);

  useEffect(() => {
    if (expanded) {
      dispatch(fetchMatchesRequest());
    }
  }, [dispatch, expanded]);

  const handleDeleteMatch = (id: string) => {
    if (window.confirm('Are you sure you want to delete this match?')) {
      dispatch(deleteMatchRequest(id));
    }
  };

  const handleEditMatch = (match: Match) => {
    setEditingMatch(match);
    setShowMatchForm(true);
  };

  const handleMatchFormSuccess = () => {
    setShowMatchForm(false);
    setEditingMatch(undefined);
    dispatch(fetchMatchesRequest());
  };

  const handleMatchFormCancel = () => {
    setShowMatchForm(false);
    setEditingMatch(undefined);
  };

  if (showMatchForm) {
    return (
      <MatchForm
        match={editingMatch || undefined}
        seriesId={series.id}
        onSuccess={handleMatchFormSuccess}
        onCancel={handleMatchFormCancel}
      />
    );
  }

  return (
    <Card key={series.id} data-cy="series-item">
      <CardHeader>
        <div className="flex items-center justify-between">
          <CardTitle data-cy="series-name">{series.name}</CardTitle>
          <Button
            variant="ghost"
            size="sm"
            onClick={() => setExpanded(!expanded)}
          >
            {expanded ? 'Hide Matches' : 'Show Matches'}
          </Button>
        </div>
      </CardHeader>
      <CardContent>
        <div className="space-y-2 text-sm">
          <div className="flex items-center">
            <Calendar className="h-4 w-4 mr-2 text-gray-500" />
            <span className="font-medium">Start:</span>
            <span className="ml-2">{formatDate(series.start_date)}</span>
          </div>
          <div className="flex items-center">
            <Calendar className="h-4 w-4 mr-2 text-gray-500" />
            <span className="font-medium">End:</span>
            <span className="ml-2">{formatDate(series.end_date)}</span>
          </div>
        </div>

        {expanded && (
          <div className="mt-6">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-lg font-semibold">
                Matches ({seriesMatches.length})
              </h3>
              <div className="flex space-x-2">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => dispatch(fetchMatchesRequest())}
                  disabled={matchesLoading}
                  data-cy="refresh-matches-button"
                  title="Refresh"
                >
                  <RefreshCw
                    className={`h-4 w-4 ${matchesLoading ? 'animate-spin' : ''}`}
                  />
                </Button>
                <Button
                  size="sm"
                  onClick={() => setShowMatchForm(true)}
                  data-cy="create-match-button"
                  title="Add Match"
                >
                  <Plus className="h-4 w-4 mr-2" />
                  Match
                </Button>
              </div>
            </div>

            {matchesLoading ? (
              <div className="flex items-center justify-center py-4">
                <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-blue-600"></div>
                <span className="ml-2 text-sm text-gray-600">
                  Loading matches...
                </span>
              </div>
            ) : matchesError ? (
              <div className="bg-red-100 border border-red-400 text-red-700 px-3 py-2 rounded text-sm">
                Error loading matches: {matchesError}
              </div>
            ) : seriesMatches.length === 0 ? (
              <div className="text-center py-4 text-gray-500">
                <p className="mb-2">No matches found for this series.</p>
                <Button
                  size="sm"
                  variant="outline"
                  onClick={() => setShowMatchForm(true)}
                >
                  Create First Match
                </Button>
              </div>
            ) : (
              <div className="space-y-3">
                {seriesMatches.map((match, index) => (
                  <Card
                    key={match.id || `match-${index}`}
                    className="bg-gray-50"
                  >
                    <CardContent className="pt-4">
                      <div className="space-y-3">
                        {/* Match Details */}
                        <div
                          className="space-y-1 cursor-pointer"
                          onClick={() => onViewScorecard?.(match.id)}
                        >
                          <div className="font-medium">
                            Match #{match.match_number}
                          </div>
                          <div className="text-sm text-gray-600">
                            {match.date.split('T')[0]} •{' '}
                            {match.team_a_player_count} vs{' '}
                            {match.team_b_player_count} players •{' '}
                            {match.total_overs} overs
                          </div>
                          <div className="flex items-center space-x-2">
                            <Badge
                              variant={
                                match.status === 'live'
                                  ? 'default'
                                  : match.status === 'completed'
                                    ? 'secondary'
                                    : 'outline'
                              }
                            >
                              {match.status}
                            </Badge>
                            <span className="text-xs text-gray-500">
                              Toss: Team {match.toss_winner} (
                              {match.toss_type === 'H' ? 'Heads' : 'Tails'})
                            </span>
                          </div>
                        </div>

                        {/* Action Buttons */}
                        <div className="flex flex-wrap gap-2 pt-2 border-t border-gray-200">
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => onViewScorecard?.(match.id)}
                            data-cy="view-scorecard-button"
                            title="View Scorecard"
                          >
                            <Play className="h-4 w-4 mr-2" />
                            Scorecard
                          </Button>
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => handleEditMatch(match)}
                            data-cy="edit-match-button"
                            title="Edit Match"
                          >
                            <Edit className="h-4 w-4" />
                          </Button>
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => handleDeleteMatch(match.id)}
                            data-cy="delete-match-button"
                            title="Delete Match"
                          >
                            <Trash2 className="h-4 w-4" />
                          </Button>
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                ))}
              </div>
            )}
          </div>
        )}

        <div className="flex space-x-2 mt-4">
          <Button
            variant="outline"
            size="sm"
            onClick={() => onEditSeries(series)}
            data-cy="edit-series-button"
            title="Edit Series"
          >
            <Edit className="h-4 w-4 mr-2" />
            Series
          </Button>
          <Button
            variant="destructive"
            size="sm"
            onClick={() => onDeleteSeries(series.id)}
            data-cy="delete-series-button"
            title="Delete Series"
          >
            <Trash2 className="h-4 w-4 mr-2" />
            Series
          </Button>
        </div>
      </CardContent>
    </Card>
  );
}
