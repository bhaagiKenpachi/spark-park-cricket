'use client';

import { useState, useEffect, useCallback, useMemo } from 'react';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import {
  fetchMatchesBySeriesRequest,
  deleteMatchRequest,
  startMatchRequest,
  Match,
} from '@/store/reducers/matchSlice';
import { useCompletedMatchesCache } from '@/hooks/useCompletedMatchesCache';
import { fetchScorecardRequest } from '@/store/reducers/scorecardSlice';
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
import {
  RefreshCw,
  Plus,
  Edit,
  Trash2,
  Calendar,
  Play,
  MoreVertical,
  RotateCcw,
  ChevronDown,
  ChevronUp,
  Clock,
  Trophy,
  Share2,
  CheckCircle,
  Circle,
} from 'lucide-react';
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
  expanded?: boolean;
  onToggleExpanded?: (expanded: boolean) => void;
}

export function SeriesWithMatches({
  series,
  onEditSeries,
  onDeleteSeries,
  onViewScorecard,
  currentUser,
  isAuthenticated,
  expanded = false,
  onToggleExpanded,
}: SeriesWithMatchesProps): React.JSX.Element {
  const dispatch = useAppDispatch();
  const {
    matches,
    loading: matchesLoading,
    error: matchesError,
  } = useAppSelector(state => state.match);
  const { scorecard } = useAppSelector(state => state.scorecard);
  const { getCachedData, setCachedData, clearExpired, isCached, clearMatch } =
    useCompletedMatchesCache();
  const [showMatchForm, setShowMatchForm] = useState(false);
  const [editingMatch, setEditingMatch] = useState<Match | undefined>();
  const [scorecardData, setScorecardData] = useState<{
    [matchId: string]: ScorecardData;
  }>({});
  const [expandedMatchDetails, setExpandedMatchDetails] = useState<{
    [matchId: string]: boolean;
  }>({});

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

  // Fetch scorecard data for completed matches with caching
  const fetchMatchScorecard = useCallback(
    async (matchId: string) => {
      // Check if data exists in cache and is not expired
      const cachedData = getCachedData(matchId);

      if (cachedData) {
        // Use cached data - cast to ScorecardData since we know it's the right type
        setScorecardData(prev => ({
          ...prev,
          [matchId]: cachedData as ScorecardData,
        }));
        return;
      }

      // Check if data already exists in local state
      if (!scorecardData[matchId]) {
        try {
          dispatch(fetchScorecardRequest(matchId));
        } catch (error) {
          console.error('Error fetching scorecard:', error);
        }
      }
    },
    [dispatch, scorecardData, getCachedData]
  );

  // Filter and sort matches for this series
  const seriesMatches = useMemo(() => {
    const filteredMatches = matches?.filter(match => match.series_id === series.id) || [];
    
    // Sort matches by status priority: live first, then completed, then not_started
    return filteredMatches.sort((a, b) => {
      const statusPriority = {
        'live': 0,
        'completed': 1,
        'not_started': 2
      };
      
      const aPriority = statusPriority[a.status as keyof typeof statusPriority] ?? 3;
      const bPriority = statusPriority[b.status as keyof typeof statusPriority] ?? 3;
      
      if (aPriority !== bPriority) {
        return aPriority - bPriority;
      }
      
      // If same status, sort by match number
      return a.match_number - b.match_number;
    });
  }, [matches, series.id]);

  // Check if current user owns this series
  const isOwner =
    isAuthenticated && currentUser && series.created_by === currentUser.id;

  useEffect(() => {
    if (expanded) {
      dispatch(fetchMatchesBySeriesRequest(series.id));
    }
  }, [dispatch, expanded, series.id]);

  // Clear expired cache entries on component mount
  useEffect(() => {
    clearExpired();
  }, [clearExpired]);

  // Fetch scorecard data for completed matches
  useEffect(() => {
    if (seriesMatches.length > 0) {
      seriesMatches.forEach(match => {
        if (
          match.status === 'completed' &&
          !scorecardData[match.id] &&
          !isCached(match.id)
        ) {
          fetchMatchScorecard(match.id);
        }
      });
    }
  }, [seriesMatches, scorecardData, fetchMatchScorecard, isCached]);

  // Update scorecard data when scorecard changes and cache it
  useEffect(() => {
    if (scorecard && scorecard.match_id) {
      setScorecardData(prev => ({
        ...prev,
        [scorecard.match_id]: scorecard,
      }));

      // Cache the scorecard data for completed matches
      const match = seriesMatches.find(m => m.id === scorecard.match_id);
      if (match && match.status === 'completed') {
        setCachedData(scorecard.match_id, scorecard, 10 * 60 * 1000); // 10 minutes for completed matches
      }
    }
  }, [scorecard, dispatch, seriesMatches, setCachedData]);

  const handleDeleteMatch = (id: string) => {
    if (window.confirm('Are you sure you want to delete this match?')) {
      dispatch(deleteMatchRequest(id));
    }
  };

  const handleStartMatch = (id: string) => {
    if (window.confirm('Are you sure you want to start this match?')) {
      dispatch(startMatchRequest(id));
    }
  };

  const handleEditMatch = (match: Match) => {
    setEditingMatch(match);
    setShowMatchForm(true);
  };

  const handleMatchFormSuccess = () => {
    setShowMatchForm(false);
    setEditingMatch(undefined);
    dispatch(fetchMatchesBySeriesRequest(series.id));
  };

  const handleMatchFormCancel = () => {
    setShowMatchForm(false);
    setEditingMatch(undefined);
  };

  const handleShareMatch = async (matchId: string, matchNumber: number, event: React.MouseEvent) => {
    // Prevent navigation when clicking share
    event.stopPropagation();
    
    const url = `${window.location.origin}/?match=${matchId}`;
    
    try {
      await navigator.clipboard.writeText(url);
      alert('MatchURL copied!');
    } catch (error) {
      // If clipboard API fails, use the old method
      const textArea = document.createElement('textarea');
      textArea.value = url;
      document.body.appendChild(textArea);
      textArea.select();
      document.execCommand('copy');
      document.body.removeChild(textArea);
      alert('Match URL copied!');
    }
  };

  // Smart refresh that respects cache
  const handleSmartRefresh = () => {
    console.log('🔄 Smart Refresh - Checking cache for completed matches');

    // Always refresh matches list (this is lightweight)
    dispatch(fetchMatchesBySeriesRequest(series.id));

    // Only fetch scorecard data for completed matches that aren't cached
    if (seriesMatches.length > 0) {
      const completedMatches = seriesMatches.filter(
        match => match.status === 'completed'
      );
      const cachedMatches = completedMatches.filter(match =>
        isCached(match.id)
      );
      const uncachedMatches = completedMatches.filter(
        match => !isCached(match.id)
      );

      console.log(
        `📊 Cache Status: ${cachedMatches.length} cached, ${uncachedMatches.length} need API call`
      );

      uncachedMatches.forEach(match => {
        console.log(
          `🌐 API Call needed for match ${match.id} (${match.match_number})`
        );
        fetchMatchScorecard(match.id);
      });

      if (cachedMatches.length > 0) {
        console.log(
          `✅ Using cache for matches: ${cachedMatches.map(m => m.match_number).join(', ')}`
        );
      }
    }
  };

  // Force refresh that bypasses cache
  const handleForceRefresh = () => {
    console.log('🔄 Force Refresh - Bypassing cache for all completed matches');

    // Always refresh matches list
    dispatch(fetchMatchesBySeriesRequest(series.id));

    // Force fetch scorecard data for all completed matches (bypass cache)
    if (seriesMatches.length > 0) {
      const completedMatches = seriesMatches.filter(
        match => match.status === 'completed'
      );
      console.log(
        `🌐 Force API calls for ${completedMatches.length} completed matches`
      );

      completedMatches.forEach(match => {
        console.log(
          `🔄 Force refreshing match ${match.id} (${match.match_number})`
        );
        // Clear cache for this match first
        clearMatch(match.id);
        // Then fetch fresh data
        dispatch(fetchScorecardRequest(match.id));
      });
    }
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
    <Card
      key={series.id}
      data-cy="series-item"
      className="shadow-sm hover:shadow-md transition-shadow duration-200 border-0 bg-gradient-to-br from-white to-gray-50/30"
    >
      <CardHeader className="pb-4">
        <div className="flex items-center justify-between">
          <div className="flex-1">
            <CardTitle
              data-cy="series-name"
              className="text-xl font-semibold text-gray-900 mb-1"
            >
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
              onClick={() => onToggleExpanded?.(!expanded)}
              className="bg-blue-50 hover:bg-blue-100 border-blue-200 text-blue-700 font-medium shadow-sm"
            >
              {expanded ? (
                <ChevronUp className="h-4 w-4" />
              ) : (
                <ChevronDown className="h-4 w-4" />
              )}
            </Button>
            {isOwner && (
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button
                    variant="ghost"
                    size="sm"
                    className="hover:bg-gray-100"
                  >
                    <MoreVertical className="h-4 w-4" />
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end">
                  {/* Series management options - at the top */}
                  {isOwner && (
                    <>
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
                    </>
                  )}
                  {/* Cache refresh options - at the bottom, only show when matches are expanded */}
                  {expanded && (
                    <>
                      <div className="border-t border-gray-200 my-1"></div>
                      <DropdownMenuItem
                        onClick={handleSmartRefresh}
                        disabled={matchesLoading}
                        data-cy="refresh-matches-button"
                      >
                        <RefreshCw
                          className={`h-4 w-4 mr-2 ${matchesLoading ? 'animate-spin' : ''}`}
                        />
                        Smart Refresh (Cache)
                      </DropdownMenuItem>
                      <DropdownMenuItem
                        onClick={handleForceRefresh}
                        disabled={matchesLoading}
                        data-cy="force-refresh-button"
                        className="text-orange-600 focus:text-orange-600"
                      >
                        <RotateCcw
                          className={`h-4 w-4 mr-2 ${matchesLoading ? 'animate-spin' : ''}`}
                        />
                        Force Refresh (No Cache)
                      </DropdownMenuItem>
                    </>
                  )}
                </DropdownMenuContent>
              </DropdownMenu>
            )}
          </div>
        </div>
      </CardHeader>
      <CardContent className="px-3 pt-0">
        {/* Team Names Section - Show if either team name is provided */}
        {(series.team_a_name || series.team_b_name) && (
          <div className="grid grid-cols-2 md:grid-cols-2 gap-4 p-4 bg-blue-50/50 rounded-lg border border-blue-100 mb-4">
            {series.team_a_name && (
              <div className="flex items-center space-x-3">
                <div className="p-2 bg-blue-100 rounded-lg">
                  <Trophy className="h-4 w-4 text-blue-600" />
                </div>
                <div>
                  <span className="text-xs font-medium text-gray-500 uppercase tracking-wide">
                    Team A
                  </span>
                  <p className="text-sm font-semibold text-gray-900">
                    {series.team_a_name}
                  </p>
                </div>
              </div>
            )}
            {series.team_b_name && (
              <div className="flex items-center space-x-3">
                <div className="p-2 bg-purple-100 rounded-lg">
                  <Trophy className="h-4 w-4 text-purple-600" />
                </div>
                <div>
                  <span className="text-xs font-medium text-gray-500 uppercase tracking-wide">
                    Team B
                  </span>
                  <p className="text-sm font-semibold text-gray-900">
                    {series.team_b_name}
                  </p>
                </div>
              </div>
            )}
          </div>
        )}

        {/* Dates Section */}
        <div className="grid grid-cols-2 md:grid-cols-2 gap-4 p-4 bg-white/50 rounded-lg border border-gray-100">
          <div className="flex items-center space-x-3">
            <div className="p-2 bg-green-100 rounded-lg">
              <Calendar className="h-4 w-4 text-green-600" />
            </div>
            <div>
              <span className="text-xs font-medium text-gray-500 uppercase tracking-wide">
                Start Date
              </span>
              <p className="text-sm font-semibold text-gray-900">
                {formatDate(series.start_date)}
              </p>
            </div>
          </div>
          <div className="flex items-center space-x-3">
            <div className="p-2 bg-red-100 rounded-lg">
              <Calendar className="h-4 w-4 text-red-600" />
            </div>
            <div>
              <span className="text-xs font-medium text-gray-500 uppercase tracking-wide">
                End Date
              </span>
              <p className="text-sm font-semibold text-gray-900">
                {formatDate(series.end_date)}
              </p>
            </div>
          </div>
        </div>

        {expanded && (
          <div className="mt-6">
            <div className="flex items-center justify-between mb-6">
              <div>
                <h3 className="text-lg font-semibold text-gray-900">Matches</h3>
                <p className="text-sm text-gray-500">
                  {seriesMatches.length}{' '}
                  {seriesMatches.length === 1 ? 'match' : 'matches'} in this
                  series
                </p>
              </div>
              <div className="flex space-x-2">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={handleSmartRefresh}
                  disabled={matchesLoading}
                  title="Refresh Matches"
                  className="border-blue-200 text-blue-700 hover:bg-blue-50"
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
                    New
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
                  <p className="text-gray-500 mb-4">
                    No matches found for this series.
                  </p>
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
                    <CardContent className="flex flex-col">
                      <div className="flex justify-between">
                        <div
                          className="space-y-5 cursor-pointer flex items-start w-full justify-between group"
                          onClick={() =>
                            onViewScorecard?.(match.id, series.created_by || '')
                          }
                        >
                          {/* Match Header */}
                          <div className="flex items-start justify-between">
                            <div className="flex-1">
                          <div className="flex items-center ">
                            {' '}
                            <h4 className="font-bold text-gray-900 group-hover:text-blue-600 transition-colors text-xl">
                              Match #{match.match_number}
                            </h4>{' '}
                            <div className="ml-4">
                              {match.status === 'live' && (
                                <div title="Live">
                                  <Circle className="h-3 w-3 text-green-500 fill-green-500" />
                                </div>
                              )}
                              {match.status === 'completed' && (
                                <div title="Completed">
                                  <CheckCircle className="h-4 w-4 text-green-600" />
                                </div>
                              )}
                              {match.status === 'not_started' && (
                                <div title="Not Started">
                                  <Circle className="h-3 w-3 text-blue-500" />
                                </div>
                              )}
                            </div>
                          </div>

                              <p className="text-sm text-gray-600 mt-1">
                                {new Date(match.date).toLocaleDateString(
                                  'en-US',
                                  {
                                    weekday: 'short',
                                    year: 'numeric',
                                    month: 'short',
                                    day: 'numeric',
                                  }
                                )}
                              </p>
                            </div>
                          </div>

                          {/* Match Actions */}
                          <div className="flex items-center space-x-2">
                            <Button
                              variant="ghost"
                              size="sm"
                              onClick={(e) => handleShareMatch(match.id, match.match_number, e)}
                              className="hover:bg-green-50 text-green-600 hover:text-green-700"
                              title="Share Match"
                            >
                              <Share2 className="h-4 w-4" />
                            </Button>
                            <DropdownMenu>
                              <DropdownMenuTrigger asChild>
                                <Button
                                  variant="ghost"
                                  size="sm"
                                  className="hover:bg-gray-100"
                                >
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
                                  {match.status === 'not_started' && (
                                    <DropdownMenuItem
                                      onClick={() => handleStartMatch(match.id)}
                                      className="text-green-600 focus:text-green-600"
                                    >
                                      <Play className="h-4 w-4 mr-2" />
                                      Start Match
                                    </DropdownMenuItem>
                                  )}
                                  <DropdownMenuItem
                                    onClick={() => handleEditMatch(match)}
                                  >
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
                        </div>
                      </div>
                      <div>
                        {/* Team Names - Show if available */}
                        {(match.team_a_name || match.team_b_name) && (
                          <div className="bg-gradient-to-r from-blue-50 to-purple-50 rounded-lg p-4 border border-blue-100 mb-3">
                            <div className="flex items-center justify-center space-x-4">
                              <div className="flex-1 text-center">
                                <span className="text-xs font-medium text-gray-600 uppercase tracking-wide block mb-1">
                                  Team A
                                </span>
                                <p className="text-md font-bold text-blue-900">
                                  {match.team_a_name || 'Team A'}
                                </p>
                              </div>
                              <span className="text-2xl font-bold text-gray-400">
                                vs
                              </span>
                              <div className="flex-1 text-center">
                                <span className="text-xs font-medium text-gray-600 uppercase tracking-wide block mb-1">
                                  Team B
                                </span>
                                <p className="text-md font-bold text-purple-900">
                                  {match.team_b_name || 'Team B'}
                                </p>
                              </div>
                            </div>
                          </div>
                        )}

                        {/* Match Quick Info */}
                        <div className="bg-gradient-to-r from-blue-50 to-purple-50 rounded-lg p-5 border border-blue-100">
                          <div className="flex items-center justify-between">
                            <div>
                              <span className="text-xs font-medium text-gray-600 uppercase tracking-wide">
                                Teams
                              </span>
                              <p className="text-lg font-bold text-blue-900">
                                {match.team_a_player_count}v
                                {match.team_b_player_count}
                              </p>
                            </div>
                            <div>
                              <span className="text-xs font-medium text-gray-600 uppercase tracking-wide">
                                Format
                              </span>
                              <p className="text-lg font-bold text-purple-900">
                                {match.total_overs} Over
                                {match.total_overs !== 1 ? 's' : ''}
                              </p>
                            </div>
                          </div>
                        </div>

                        {/* Show More Info Button */}
                        <div className="mt-4">
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={e => {
                              e.stopPropagation();
                              setExpandedMatchDetails(prev => ({
                                ...prev,
                                [match.id]: !prev[match.id],
                              }));
                            }}
                            className="w-full h-10 border-blue-200 text-blue-700 hover:bg-blue-50 hover:border-blue-300 transition-colors"
                          >
                            {expandedMatchDetails[match.id] ? (
                              <>
                                <ChevronUp className="h-4 w-4 mr-2" />
                                Hide Details
                              </>
                            ) : (
                              <>
                                <ChevronDown className="h-4 w-4 mr-2" />
                                Show More Info
                              </>
                            )}
                          </Button>
                        </div>

                        {/* Expandable Details */}
                        {expandedMatchDetails[match.id] && (
                          <div className="space-y-4 mt-4">
                            {/* Toss Information */}
                            <div className="bg-orange-50 rounded-lg p-3 border border-orange-200">
                              <div className="flex items-center space-x-2 mb-2">
                                <Trophy className="h-4 w-4 text-orange-600" />
                                <span className="text-xs font-semibold text-orange-700 uppercase tracking-wide">
                                  Toss Result
                                </span>
                              </div>
                              <p className="text-sm font-medium text-orange-900">
                                Team {match.toss_winner} won and chose to{' '}
                                {match.toss_type === 'H'
                                  ? 'bat first'
                                  : 'bowl first'}
                              </p>
                            </div>

                            {/* Match Timing Information */}
                            {(match.start_time || match.end_time) && (
                              <div className="bg-green-50 rounded-lg p-3 border border-green-200">
                                <div className="flex items-center space-x-2 mb-2">
                                  <Clock className="h-4 w-4 text-green-600" />
                                  <span className="text-xs font-semibold text-green-700 uppercase tracking-wide">
                                    Match Timing
                                  </span>
                                </div>
                                <div className="space-y-1 text-sm">
                                  {match.start_time && (
                                    <p className="text-green-800">
                                      <span className="font-semibold">
                                        Started:
                                      </span>{' '}
                                      {new Date(
                                        match.start_time
                                      ).toLocaleString('en-US', {
                                        month: 'short',
                                        day: 'numeric',
                                        hour: '2-digit',
                                        minute: '2-digit',
                                      })}
                                    </p>
                                  )}
                                  {match.end_time && (
                                    <p className="text-green-800">
                                      <span className="font-semibold">
                                        Ended:
                                      </span>{' '}
                                      {new Date(match.end_time).toLocaleString(
                                        'en-US',
                                        {
                                          month: 'short',
                                          day: 'numeric',
                                          hour: '2-digit',
                                          minute: '2-digit',
                                        }
                                      )}
                                    </p>
                                  )}
                                  {match.start_time && match.end_time && (
                                    <p className="text-green-900 font-semibold">
                                      Duration:{' '}
                                      {(() => {
                                        const start = new Date(
                                          match.start_time
                                        );
                                        const end = new Date(match.end_time);
                                        const duration =
                                          end.getTime() - start.getTime();
                                        const hours = Math.floor(
                                          duration / (1000 * 60 * 60)
                                        );
                                        const minutes = Math.floor(
                                          (duration % (1000 * 60 * 60)) /
                                            (1000 * 60)
                                        );
                                        return `${hours}h ${minutes}m`;
                                      })()}
                                    </p>
                                  )}
                                </div>
                              </div>
                            )}

                            {/* Match Result - Only for Completed Matches */}
                            {match.status === 'completed' &&
                              (() => {
                                const matchScorecard = scorecardData[match.id];
                                if (
                                  matchScorecard &&
                                  matchScorecard.innings &&
                                  Array.isArray(matchScorecard.innings)
                                ) {
                                  const teamAInnings =
                                    matchScorecard.innings.find(
                                      innings => innings.batting_team === 'A'
                                    );
                                  const teamBInnings =
                                    matchScorecard.innings.find(
                                      innings => innings.batting_team === 'B'
                                    );

                                  if (teamAInnings && teamBInnings) {
                                    const teamARuns = teamAInnings.total_runs;
                                    const teamBRuns = teamBInnings.total_runs;
                                    const winner =
                                      teamARuns > teamBRuns
                                        ? matchScorecard.team_a
                                        : matchScorecard.team_b;
                                    const margin = Math.abs(
                                      teamARuns - teamBRuns
                                    );

                                    return (
                                      <div className="bg-gradient-to-r from-green-50 to-emerald-50 rounded-lg p-4 border border-green-200 shadow-sm">
                                        <div className="flex items-center space-x-2 mb-3">
                                          <Trophy className="h-5 w-5 text-green-600" />
                                          <span className="text-sm font-bold text-green-700 uppercase tracking-wide">
                                            Match Result
                                          </span>
                                        </div>
                                        <div className="space-y-3">
                                          <div className="bg-white rounded-md p-3 flex justify-between items-center">
                                            <span className="text-base font-semibold text-gray-800">
                                              {matchScorecard.team_a}
                                            </span>
                                            <span className="text-lg font-bold text-blue-900">
                                              {teamARuns}/
                                              {teamAInnings.total_wickets}
                                            </span>
                                          </div>
                                          <div className="bg-white rounded-md p-3 flex justify-between items-center">
                                            <span className="text-base font-semibold text-gray-800">
                                              {matchScorecard.team_b}
                                            </span>
                                            <span className="text-lg font-bold text-purple-900">
                                              {teamBRuns}/
                                              {teamBInnings.total_wickets}
                                            </span>
                                          </div>
                                          <div className="bg-gradient-to-r from-green-100 to-emerald-100 rounded-md p-3">
                                            <p className="text-base font-bold text-green-800 text-center">
                                              🏆 {winner} won by {margin} run
                                              {margin !== 1 ? 's' : ''}
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
                                      <Trophy className="h-4 w-4 text-green-600" />
                                      <span className="text-xs font-semibold text-green-700 uppercase tracking-wide">
                                        Match Result
                                      </span>
                                    </div>
                                    <p className="text-sm font-medium text-green-800">
                                      Match completed - Loading scores...
                                    </p>
                                  </div>
                                );
                              })()}
                          </div>
                        )}
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
