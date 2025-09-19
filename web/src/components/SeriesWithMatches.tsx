'use client';

import { useState, useEffect, useCallback, useMemo } from 'react';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import {
  fetchMatchesRequest,
  deleteMatchRequest,
  Match,
} from '@/store/reducers/matchSlice';
import {
  fetchScorecardRequest,
} from '@/store/reducers/scorecardSlice';
import { MatchForm } from './MatchForm';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { RefreshCw, Plus, Edit, Trash2, Calendar, Play, MoreVertical } from 'lucide-react';
import { Series } from '@/store/reducers/seriesSlice';
import { User } from '@/services/authService';

interface ScorecardData {
  match_id: string;
  team_a: string;
  team_b: string;
  innings: Array<{
    batting_team: string;
    total_runs: number;
    total_wickets: number;
  }>;
}

interface SeriesWithMatchesProps {
  series: Series;
  onEditSeries: (series: Series) => void;
  onDeleteSeries: (id: string) => void;
  onViewScorecard?: (matchId: string, seriesCreatedBy: string) => void;
  currentUser?: User | null;
  isAuthenticated: boolean;
}

export function SeriesWithMatches({
  series,
  onEditSeries,
  onDeleteSeries,
  onViewScorecard,
  currentUser,
  isAuthenticated,
}: SeriesWithMatchesProps): React.JSX.Element {
  const dispatch = useAppDispatch();
  const {
    matches,
    loading: matchesLoading,
    error: matchesError,
  } = useAppSelector(state => state.match);
  const { scorecard } = useAppSelector(state => state.scorecard);
  const [showMatchForm, setShowMatchForm] = useState(false);
  const [editingMatch, setEditingMatch] = useState<Match | undefined>();
  const [expanded, setExpanded] = useState(false);
  const [scorecardData, setScorecardData] = useState<{ [matchId: string]: ScorecardData }>({});

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

  // Fetch scorecard data for completed matches
  const fetchMatchScorecard = useCallback(async (matchId: string) => {
    if (!scorecardData[matchId]) {
      try {
        dispatch(fetchScorecardRequest(matchId));
      } catch (error) {
        console.error('Error fetching scorecard:', error);
      }
    }
  }, [dispatch, scorecardData]);

  // Filter matches for this series
  const seriesMatches = useMemo(() =>
    matches?.filter(match => match.series_id === series.id) || [],
    [matches, series.id]
  );

  // Check if current user owns this series
  const isOwner =
    isAuthenticated && currentUser && series.created_by === currentUser.id;

  useEffect(() => {
    if (expanded) {
      dispatch(fetchMatchesRequest());
    }
  }, [dispatch, expanded]);

  // Fetch scorecard data for completed matches
  useEffect(() => {
    if (seriesMatches.length > 0) {
      seriesMatches.forEach(match => {
        if (match.status === 'completed' && !scorecardData[match.id]) {
          fetchMatchScorecard(match.id);
        }
      });
    }
  }, [seriesMatches, scorecardData, fetchMatchScorecard]);

  // Update scorecard data when scorecard changes
  useEffect(() => {
    if (scorecard && scorecard.match_id) {
      setScorecardData(prev => ({
        ...prev,
        [scorecard.match_id]: scorecard
      }));
    }
  }, [scorecard]);

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
    <Card key={series.id} data-cy="series-item" className="shadow-sm hover:shadow-md transition-shadow duration-200 border-0 bg-gradient-to-br from-white to-gray-50/30">
      <CardHeader className="pb-4">
        <div className="flex items-center justify-between">
          <div className="flex-1">
            <CardTitle data-cy="series-name" className="text-xl font-semibold text-gray-900 mb-1">
              {series.name}
            </CardTitle>
            {series.description && (
              <p className="text-sm text-gray-600 mt-1">{series.description}</p>
            )}
          </div>
          <div className="flex items-center space-x-2">
            <Button
              variant="outline"
              size="sm"
              onClick={() => setExpanded(!expanded)}
              className="bg-blue-50 hover:bg-blue-100 border-blue-200 text-blue-700 font-medium shadow-sm"
            >
              {expanded ? 'Hide Matches' : 'Show Matches'}
            </Button>
            {isOwner && (
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button variant="ghost" size="sm" className="hover:bg-gray-100">
                    <MoreVertical className="h-4 w-4" />
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end">
                  <DropdownMenuItem onClick={() => onEditSeries(series)}>
                    <Edit className="h-4 w-4 mr-2" />
                    Edit Series
                  </DropdownMenuItem>
                  <DropdownMenuItem
                    onClick={() => onDeleteSeries(series.id)}
                    className="text-red-600 focus:text-red-600"
                  >
                    <Trash2 className="h-4 w-4 mr-2" />
                    Delete Series
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            )}
          </div>
        </div>
      </CardHeader>
      <CardContent className="pt-0">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 p-4 bg-white/50 rounded-lg border border-gray-100">
          <div className="flex items-center space-x-3">
            <div className="p-2 bg-green-100 rounded-lg">
              <Calendar className="h-4 w-4 text-green-600" />
            </div>
            <div>
              <span className="text-xs font-medium text-gray-500 uppercase tracking-wide">Start Date</span>
              <p className="text-sm font-semibold text-gray-900">{formatDate(series.start_date)}</p>
            </div>
          </div>
          <div className="flex items-center space-x-3">
            <div className="p-2 bg-red-100 rounded-lg">
              <Calendar className="h-4 w-4 text-red-600" />
            </div>
            <div>
              <span className="text-xs font-medium text-gray-500 uppercase tracking-wide">End Date</span>
              <p className="text-sm font-semibold text-gray-900">{formatDate(series.end_date)}</p>
            </div>
          </div>
        </div>

        {expanded && (
          <div className="mt-6">
            <div className="flex items-center justify-between mb-6">
              <div>
                <h3 className="text-lg font-semibold text-gray-900">
                  Matches
                </h3>
                <p className="text-sm text-gray-500">
                  {seriesMatches.length} {seriesMatches.length === 1 ? 'match' : 'matches'} in this series
                </p>
              </div>
              <div className="flex space-x-2">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => dispatch(fetchMatchesRequest())}
                  disabled={matchesLoading}
                  data-cy="refresh-matches-button"
                  title="Refresh"
                  className="hover:bg-gray-50"
                >
                  <RefreshCw
                    className={`h-4 w-4 ${matchesLoading ? 'animate-spin' : ''}`}
                  />
                </Button>
                {isOwner && (
                  <Button
                    size="sm"
                    onClick={() => setShowMatchForm(true)}
                    data-cy="create-match-button"
                    title="Add Match"
                    className="bg-blue-600 hover:bg-blue-700 text-white shadow-sm"
                  >
                    <Plus className="h-4 w-4 mr-2" />
                    Add Match
                  </Button>
                )}
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
              <div className="text-center py-12">
                <div className="p-4 bg-gray-50 rounded-lg border-2 border-dashed border-gray-200">
                  <Play className="h-12 w-12 text-gray-400 mx-auto mb-4" />
                  <p className="text-gray-500 mb-4">No matches found for this series.</p>
                  {isOwner && (
                    <Button
                      size="sm"
                      onClick={() => setShowMatchForm(true)}
                      className="bg-blue-600 hover:bg-blue-700 text-white"
                    >
                      <Plus className="h-4 w-4 mr-2" />
                      Create First Match
                    </Button>
                  )}
                </div>
              </div>
            ) : (
              <div className="grid gap-4">
                {seriesMatches.map((match, index) => (
                  <Card
                    key={match.id || `match-${index}`}
                    className="bg-white border border-gray-200 hover:border-gray-300 transition-colors duration-200 shadow-sm hover:shadow-md"
                  >
                    <CardContent className="p-4">
                      <div className="flex items-start justify-between">
                        <div
                          className="space-y-4 cursor-pointer flex-1 group"
                          onClick={() =>
                            onViewScorecard?.(
                              match.id,
                              series.created_by || ''
                            )
                          }
                        >
                          {/* Match Header */}
                          <div className="flex items-center justify-between">
                            <div>
                              <h4 className="font-semibold text-gray-900 group-hover:text-blue-600 transition-colors text-lg">
                                Match #{match.match_number}
                              </h4>
                              <p className="text-sm text-gray-500 mt-1">
                                {new Date(match.date).toLocaleDateString('en-US', {
                                  weekday: 'short',
                                  year: 'numeric',
                                  month: 'short',
                                  day: 'numeric'
                                })}
                              </p>
                            </div>
                            <Badge
                              variant={
                                match.status === 'live'
                                  ? 'default'
                                  : match.status === 'completed'
                                    ? 'secondary'
                                    : 'outline'
                              }
                              className={
                                match.status === 'live'
                                  ? 'bg-green-500 text-white border-green-500 font-semibold'
                                  : match.status === 'completed'
                                    ? 'bg-gray-100 text-gray-800 border-gray-200 font-semibold'
                                    : 'bg-yellow-100 text-yellow-800 border-yellow-200 font-semibold'
                              }
                            >
                              {match.status.toUpperCase()}
                            </Badge>
                          </div>

                          {/* Match Details Grid */}
                          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                            {/* Teams Info */}
                            <div className="bg-gray-50 rounded-lg p-3">
                              <div className="flex items-center space-x-2 mb-2">
                                <div className="w-2 h-2 bg-blue-500 rounded-full"></div>
                                <span className="text-xs font-medium text-gray-500 uppercase tracking-wide">Teams</span>
                              </div>
                              <p className="text-sm font-semibold text-gray-900">
                                {match.team_a_player_count} vs {match.team_b_player_count} players
                              </p>
                            </div>

                            {/* Match Format */}
                            <div className="bg-gray-50 rounded-lg p-3">
                              <div className="flex items-center space-x-2 mb-2">
                                <div className="w-2 h-2 bg-purple-500 rounded-full"></div>
                                <span className="text-xs font-medium text-gray-500 uppercase tracking-wide">Format</span>
                              </div>
                              <p className="text-sm font-semibold text-gray-900">
                                {match.total_overs} overs
                              </p>
                            </div>
                          </div>

                          {/* Toss Information */}
                          <div className="bg-blue-50 rounded-lg p-3 border border-blue-100">
                            <div className="flex items-center space-x-2 mb-2">
                              <div className="w-2 h-2 bg-orange-500 rounded-full"></div>
                              <span className="text-xs font-medium text-blue-600 uppercase tracking-wide">Toss Result</span>
                            </div>
                            <p className="text-sm font-semibold text-blue-900">
                              Team {match.toss_winner} won the toss and chose to {match.toss_type === 'H' ? 'bat first' : 'bowl first'}
                            </p>
                          </div>

                          {/* Match Completion Summary */}
                          {match.status === 'completed' && (() => {
                            const matchScorecard = scorecardData[match.id];
                            if (matchScorecard && matchScorecard.innings && Array.isArray(matchScorecard.innings)) {
                              const teamAInnings = matchScorecard.innings.find(innings => innings.batting_team === 'A');
                              const teamBInnings = matchScorecard.innings.find(innings => innings.batting_team === 'B');

                              if (teamAInnings && teamBInnings) {
                                const teamARuns = teamAInnings.total_runs;
                                const teamBRuns = teamBInnings.total_runs;
                                const winner = teamARuns > teamBRuns ? matchScorecard.team_a : matchScorecard.team_b;
                                const margin = Math.abs(teamARuns - teamBRuns);

                                return (
                                  <div className="bg-green-50 rounded-lg p-3 border border-green-200">
                                    <div className="flex items-center space-x-2 mb-3">
                                      <div className="w-2 h-2 bg-green-500 rounded-full"></div>
                                      <span className="text-xs font-medium text-green-600 uppercase tracking-wide">Match Result</span>
                                    </div>
                                    <div className="space-y-2">
                                      <div className="flex justify-between items-center">
                                        <span className="text-sm font-medium text-gray-700">{matchScorecard.team_a}</span>
                                        <span className="text-sm font-bold text-gray-900">{teamARuns}/{teamAInnings.total_wickets}</span>
                                      </div>
                                      <div className="flex justify-between items-center">
                                        <span className="text-sm font-medium text-gray-700">{matchScorecard.team_b}</span>
                                        <span className="text-sm font-bold text-gray-900">{teamBRuns}/{teamBInnings.total_wickets}</span>
                                      </div>
                                      <div className="pt-2 border-t border-green-200">
                                        <p className="text-sm font-semibold text-green-800 text-center">
                                          {winner} won by {margin} run{margin !== 1 ? 's' : ''}
                                        </p>
                                      </div>
                                    </div>
                                  </div>
                                );
                              }
                            }

                            // Fallback for when scorecard data is not available yet
                            return (
                              <div className="bg-green-50 rounded-lg p-3 border border-green-200">
                                <div className="flex items-center space-x-2 mb-2">
                                  <div className="w-2 h-2 bg-green-500 rounded-full"></div>
                                  <span className="text-xs font-medium text-green-600 uppercase tracking-wide">Match Result</span>
                                </div>
                                <p className="text-sm font-semibold text-green-800">
                                  Match completed - Loading scores...
                                </p>
                              </div>
                            );
                          })()}
                        </div>

                        {/* Match Actions Dropdown */}
                        <DropdownMenu>
                          <DropdownMenuTrigger asChild>
                            <Button variant="ghost" size="sm" className="hover:bg-gray-100">
                              <MoreVertical className="h-4 w-4" />
                            </Button>
                          </DropdownMenuTrigger>
                          <DropdownMenuContent align="end">
                            <DropdownMenuItem
                              onClick={() =>
                                onViewScorecard?.(
                                  match.id,
                                  series.created_by || ''
                                )
                              }
                            >
                              <Play className="h-4 w-4 mr-2" />
                              View Scorecard
                            </DropdownMenuItem>
                            {isOwner && (
                              <>
                                <DropdownMenuItem onClick={() => handleEditMatch(match)}>
                                  <Edit className="h-4 w-4 mr-2" />
                                  Edit Match
                                </DropdownMenuItem>
                                <DropdownMenuItem
                                  onClick={() => handleDeleteMatch(match.id)}
                                  className="text-red-600 focus:text-red-600"
                                >
                                  <Trash2 className="h-4 w-4 mr-2" />
                                  Delete Match
                                </DropdownMenuItem>
                              </>
                            )}
                          </DropdownMenuContent>
                        </DropdownMenu>
                      </div>
                    </CardContent>
                  </Card>
                ))}
              </div>
            )}
          </div>
        )}

      </CardContent>
    </Card>
  );
}
