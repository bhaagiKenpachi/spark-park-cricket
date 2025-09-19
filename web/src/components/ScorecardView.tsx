'use client';

import { useState, useEffect } from 'react';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import {
  fetchScorecardRequest,
  startScoringRequest,
  addBallRequest,
  undoBallThunk,
  clearScorecard,
  fetchAllOversDetailsThunk,
  BallEventRequest,
  BallType,
  RunType,
  BallSummary,
  OverSummary,
  InningsSummary,
} from '@/store/reducers/scorecardSlice';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import {
  ArrowLeft,
  Play,
  X,
  ChevronDown,
  ChevronUp,
  RefreshCw,
  Undo2,
} from 'lucide-react';
import { User } from '@/services/authService';

interface ScorecardViewProps {
  matchId: string;
  onBack: () => void;
  seriesCreatedBy?: string;
  currentUser?: User | null;
  isAuthenticated: boolean;
}

export function ScorecardView({
  matchId,
  onBack,
  seriesCreatedBy,
  currentUser,
  isAuthenticated,
}: ScorecardViewProps): React.JSX.Element {
  const dispatch = useAppDispatch();
  const { scorecard, loading, error, scoring } = useAppSelector(
    state => state.scorecard
  );

  const [showLiveScoring, setShowLiveScoring] = useState(false);
  const [currentInn, setCurrentInn] = useState(1);
  const [currentByes, setCurrentByes] = useState(0);
  const [scoringMessage, setScoringMessage] = useState<string | null>(null);
  const [expandedOvers, setExpandedOvers] = useState<{
    [key: string]: boolean;
  }>({});

  // Check if current user owns the series
  const isOwner =
    isAuthenticated && currentUser && seriesCreatedBy === currentUser.id;

  // Check if match is completed or both innings are completed
  const isMatchCompleted = scorecard?.match_status === 'completed';
  const isBothInnCompleted = scorecard?.innings &&
    scorecard.innings.length >= 2 &&
    scorecard.innings.every(innings => innings.status === 'completed');

  // Determine if scoring should be available
  const isScoringAvailable = isOwner && !isMatchCompleted && !isBothInnCompleted;


  useEffect(() => {
    dispatch(fetchScorecardRequest(matchId));
    return () => {
      dispatch(clearScorecard());
    };
  }, [dispatch, matchId]);

  // Auto-show live scoring if match is already live
  useEffect(() => {
    if (scorecard?.match_status === 'live') {
      setShowLiveScoring(true);
    }
  }, [scorecard]);

  // Auto-detect current innings from scorecard data
  useEffect(() => {

    if (scorecard?.innings && Array.isArray(scorecard.innings) && scorecard.innings.length > 0) {
      const currentInnData = scorecard.innings.find(
        innings => innings.status === 'in_progress'
      );
      if (currentInnData) {
        setCurrentInn(currentInnData.innings_number);
      }
    } else if (scorecard?.innings === null || (Array.isArray(scorecard?.innings) && scorecard.innings.length === 0)) {
      // If no innings exist yet, start with innings 1
      setCurrentInn(1);
    }
  }, [scorecard]);

  // Handle scoring success
  useEffect(() => {
    if (scoring === false && showLiveScoring) {
      setScoringMessage('Live scoring is now active!');
      setTimeout(() => setScoringMessage(null), 3000);
    }
  }, [scoring, showLiveScoring]);

  // Handle ball scoring success
  useEffect(() => {
    if (scoring === false && showLiveScoring && scorecard) {
      setScoringMessage('Ball scored successfully!');
      setTimeout(() => setScoringMessage(null), 2000);
    }
  }, [scoring, showLiveScoring, scorecard]);

  const handleStartScoring = () => {

    // Check if scoring is available (ownership + match not completed)
    if (!isScoringAvailable) {
      if (!isOwner) {
        setScoringMessage('Only the series creator can start scoring.');
      } else if (isMatchCompleted) {
        setScoringMessage('Cannot score on completed match.');
      } else if (isBothInnCompleted) {
        setScoringMessage('Cannot score when both innings are completed.');
      }
      setTimeout(() => setScoringMessage(null), 3000);
      return;
    }


    // If match is already live, just show the interface without calling the API
    if (scorecardData?.match_status === 'live') {
      setShowLiveScoring(true);
      setScoringMessage('Live scoring interface opened!');
      setTimeout(() => setScoringMessage(null), 3000);
    } else {
      // Only call the API if match is not live yet
      dispatch(startScoringRequest(matchId));
      setShowLiveScoring(true);
      setScoringMessage('Live scoring started!');
      setTimeout(() => setScoringMessage(null), 3000);
    }
  };

  const handleBallScore = (runs: number, ballType: string) => {


    // Check if scoring is available (ownership + match not completed)
    if (!isScoringAvailable) {
      if (!isOwner) {
        setScoringMessage('Only the series creator can score balls.');
      } else if (isMatchCompleted) {
        setScoringMessage('Cannot score on completed match.');
      } else if (isBothInnCompleted) {
        setScoringMessage('Cannot score when both innings are completed.');
      }
      setTimeout(() => setScoringMessage(null), 3000);
      return;
    }


    // Check if current innings is still in progress
    const currentInnDataForScoring = scorecardData?.innings?.find(
      innings => innings.innings_number === currentInn
    );

    // If no innings exist yet (null or empty array), allow scoring to create the first innings
    if (scorecardData?.innings === null || (Array.isArray(scorecardData?.innings) && scorecardData.innings.length === 0)) {
      // Allow scoring - this will create the first innings
    } else if (currentInnDataForScoring?.status !== 'in_progress') {
      setScoringMessage(
        'Cannot score on completed innings. Please check innings status.'
      );
      setTimeout(() => setScoringMessage(null), 5000);
      return;
    }

    const isWicket = [
      'bowled',
      'caught',
      'lbw',
      'run_out',
      'stumped',
      'hit_wicket',
    ].includes(ballType);
    const runType = isWicket ? 'WC' : (runs.toString() as RunType);

    // For wickets, the ball type should be 'good' and wicket type should be the actual wicket type
    const actualBallType = isWicket ? 'good' : (ballType as BallType);
    const wicketType = isWicket ? ballType : undefined;

    const ballEvent: BallEventRequest = {
      match_id: matchId,
      innings_number: currentInn,
      ball_type: actualBallType,
      run_type: runType,
      runs,
      byes: currentByes,
      is_wicket: isWicket,
      ...(wicketType && { wicket_type: wicketType }),
    };

    dispatch(addBallRequest(ballEvent));
    setCurrentByes(0); // Reset byes after scoring

    // Ball counting is handled by the backend
  };

  // Helper function to check if it's the first ball of the current innings
  const isFirstBallOfInn = () => {

    const currentInnData = scorecardData?.innings?.find(
      innings => innings.innings_number === currentInn
    );
    if (!currentInnData) {
      return true; // If no innings data, consider it first ball
    }

    // Count total balls across all overs in this innings
    const totalBalls = currentInnData.overs?.reduce((total, over) => {
      return total + (over.balls ? over.balls.length : 0);
    }, 0) || 0;

    // If there's exactly 1 ball, it's the first ball (and we can't undo it)
    return totalBalls === 1;
  };

  const handleUndoBall = () => {

    // Check if scoring is available (ownership + match not completed)
    if (!isScoringAvailable) {
      if (!isOwner) {
        setScoringMessage('Only the series creator can undo balls.');
      } else if (isMatchCompleted) {
        setScoringMessage('Cannot undo ball on completed match.');
      } else if (isBothInnCompleted) {
        setScoringMessage('Cannot undo ball when both innings are completed.');
      }
      setTimeout(() => setScoringMessage(null), 3000);
      return;
    }

    // Check if current innings is still in progress
    const currentInnDataForUndo = scorecardData?.innings?.find(
      innings => innings.innings_number === currentInn
    );

    if (!currentInnDataForUndo || currentInnDataForUndo.status !== 'in_progress') {
      setScoringMessage(
        'Cannot undo ball on completed innings. Please check innings status.'
      );
      setTimeout(() => setScoringMessage(null), 5000);
      return;
    }

    // Check if it's the first ball of the innings
    if (isFirstBallOfInn()) {
      setScoringMessage('Cannot undo ball - this is the first ball of the innings.');
      setTimeout(() => setScoringMessage(null), 3000);
      return;
    }

    dispatch(undoBallThunk({ matchId, inningsNumber: currentInn }));
    setScoringMessage('Undoing last ball...');

    // Automatically refresh data after undo
    setTimeout(() => {
      dispatch(fetchScorecardRequest(matchId));
      setScoringMessage(null);
    }, 1000);
  };

  const handleByesChange = (byes: number) => {
    setCurrentByes(byes);
  };

  const handleRefresh = () => {
    dispatch(fetchScorecardRequest(matchId));
  };

  const toggleExpandedOvers = (inningsKey: string) => {
    const isCurrentlyExpanded = expandedOvers[inningsKey];

    setExpandedOvers(prev => ({
      ...prev,
      [inningsKey]: !prev[inningsKey],
    }));

    // If expanding (not collapsing), fetch all overs details
    if (!isCurrentlyExpanded) {
      // Extract innings number from the key (format: "A-1" or "B-1")
      const inningsNumberStr = inningsKey.split('-')[1];
      if (inningsNumberStr) {
        const inningsNumber = parseInt(inningsNumberStr);

        // Fetch all overs details for this innings
        dispatch(fetchAllOversDetailsThunk({
          matchId,
          inningsNumber,
        }));
      }
    }
  };

  const renderBallCircle = (ball: BallSummary, index: number) => {
    const isWicket = ball.is_wicket;

    // Determine display based on ball type and run type
    let display: string;
    if (isWicket) {
      display = 'W';
    } else {
      // Check ball_type first for special deliveries
      switch (ball.ball_type) {
        case 'WIDE':
        case 'wide':
          display = 'Wd';
          break;
        case 'NO_BALL':
        case 'no_ball':
          display = 'Nb';
          break;
        case 'DEAD_BALL':
        case 'dead_ball':
          display = 'Db';
          break;
        default:
          // For good balls, check run_type for special cases
          switch (ball.run_type) {
            case 'LB':
              display = 'Lb';
              break;
            case 'WC':
              display = 'W';
              break;
            default:
              display = ball.runs.toString();
              break;
          }
          break;
      }
    }

    // Special handling for no ball display
    if (ball.ball_type === 'NO_BALL' || ball.ball_type === 'no_ball') {
      const noBallRuns = ball.runs || 0;
      const noBallByes = ball.byes || 0;

      return (
        <div
          key={index}
          className="w-8 h-8 rounded-full border-2 border-orange-500 bg-orange-100 flex flex-col items-center justify-center text-xs font-medium"
        >
          <div className="text-[10px] leading-none text-orange-700 font-bold">nb</div>
          <div className="text-[6px] leading-none text-orange-600">+</div>
          <div className="text-[8px] leading-none text-orange-700 font-bold">
            {noBallRuns + noBallByes}
          </div>
        </div>
      );
    }

    // Special handling for balls with byes
    if (ball.byes > 0) {
      const totalRuns = ball.runs + ball.byes;

      return (
        <div
          key={index}
          className="w-8 h-8 rounded-full border-2 border-slate-400 bg-slate-100 flex flex-col items-center justify-center text-xs font-medium"
        >
          <div className="text-[10px] leading-none text-slate-700 font-bold">B</div>
          <div className="text-[8px] leading-none text-slate-600">+</div>
          <div className="text-[10px] leading-none text-slate-700 font-bold">
            {totalRuns}
          </div>
        </div>
      );
    }

    // Handle display for other ball types without byes
    return (
      <div
        key={index}
        className={`w-8 h-8 rounded-full border-2 flex items-center justify-center text-xs font-medium ${isWicket
          ? 'border-red-500 bg-red-100 text-red-700'
          : ball.ball_type === 'WIDE' || ball.ball_type === 'wide' || ball.run_type === 'LB'
            ? 'border-slate-400 bg-slate-100 text-slate-700'
            : ball.ball_type === 'DEAD_BALL' || ball.ball_type === 'dead_ball'
              ? 'border-gray-500 bg-gray-100 text-gray-700'
              : ball.runs === 4
                ? 'border-blue-500 bg-blue-100 text-blue-700'
                : ball.runs === 6
                  ? 'border-purple-500 bg-purple-100 text-purple-700'
                  : ball.runs === 0
                    ? 'border-gray-300 bg-gray-100 text-gray-600'
                    : 'border-green-500 bg-green-100 text-green-700'
          }`}
      >
        {display}
      </div>
    );
  };

  const renderOverDetails = (over: OverSummary) => (
    <div key={over.over_number} className="mb-4">
      <div className="flex items-center justify-between mb-2">
        <h5 className="font-medium text-sm">Over {over.over_number}</h5>
        <div className="text-xs text-gray-600">
          {over.total_runs} runs, {over.total_wickets} wickets
        </div>
      </div>
      <div className="flex flex-wrap gap-1">
        {over.balls && Array.isArray(over.balls) && over.balls.length > 0 ? (
          over.balls.map((ball: BallSummary, index: number) =>
            renderBallCircle(ball, index)
          )
        ) : (
          <div className="text-xs text-gray-400">Over not started</div>
        )}
      </div>
    </div>
  );

  if (loading && !scorecard) {
    return (
      <div className="w-full max-w-6xl mx-auto p-6">
        <div className="flex items-center justify-center py-8">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
          <span className="ml-2 text-sm text-gray-600">
            Loading scorecard...
          </span>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="w-full max-w-6xl mx-auto p-6">
        <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded-lg">
          <strong className="font-bold">Error:</strong>
          <span className="block sm:inline"> {error}</span>
          <div className="mt-3">
            <Button
              onClick={() => dispatch(fetchScorecardRequest(matchId))}
              variant="destructive"
              size="sm"
            >
              Retry
            </Button>
          </div>
        </div>
      </div>
    );
  }

  if (!scorecard) {
    return (
      <div className="w-full max-w-6xl mx-auto p-6">
        <div className="text-center py-8">
          <p className="text-muted-foreground mb-4">
            No scorecard found for this match.
          </p>
          <Button onClick={onBack}>Back to Matches</Button>
        </div>
      </div>
    );
  }

  const scorecardData = scorecard;

  return (
    <div className="w-full max-w-6xl mx-auto p-6">
      {/* Header */}
      <div className="mb-6">
        <div className="text-center mb-4">
          <h1 className="text-2xl lg:text-3xl font-bold">
            {scorecardData.series_name} - Match #{scorecardData.match_number}
          </h1>
          <p className="text-sm lg:text-base text-gray-600">
            {scorecardData.team_a} vs {scorecardData.team_b}
          </p>
        </div>
        <div className="flex justify-center space-x-2">
          <Button variant="outline" onClick={onBack} title="Back">
            <ArrowLeft className="h-4 w-4" />
          </Button>
          <Button
            variant="outline"
            onClick={handleRefresh}
            title="Refresh Scorecard"
            disabled={loading}
          >
            {loading ? (
              <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-gray-600"></div>
            ) : (
              <RefreshCw className="h-4 w-4" />
            )}
          </Button>
        </div>
      </div>

      {/* Match Status and Live Scoring Button */}
      <div className="mb-6 text-center">
        <div className="flex items-center justify-center space-x-4">
          <Badge
            variant={
              scorecardData.match_status === 'live' ? 'default' : 'secondary'
            }
            className={
              scorecardData.match_status === 'live' ? 'bg-green-600' : ''
            }
          >
            {scorecardData.match_status.toUpperCase()}
          </Badge>
          {scorecardData.match_status === 'live' &&
            !showLiveScoring &&
            isScoringAvailable && (
              <Button
                onClick={handleStartScoring}
                className="bg-green-600 hover:bg-green-700"
                title={
                  scorecardData.match_status === 'live'
                    ? 'Open Live Scoring'
                    : 'Start Live Scoring'
                }
                disabled={scoring}
              >
                {scoring ? (
                  <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
                ) : (
                  <>
                    <Play className="h-4 w-4 mr-2" />
                    Live Scoring
                  </>
                )}
              </Button>
            )}
        </div>
        {scoringMessage && (
          <div className="mt-3">
            <Badge variant="default" className="bg-green-600">
              {scoringMessage}
            </Badge>
          </div>
        )}
        {isMatchCompleted && (
          <div className="mt-3">
            <Badge variant="secondary" className="bg-gray-500 text-white">
              Match Completed - Scoring Not Available
            </Badge>
          </div>
        )}
        {isBothInnCompleted && !isMatchCompleted && (
          <div className="mt-3">
            <Badge variant="secondary" className="bg-gray-500 text-white">
              Both Inn Completed - Scoring Not Available
            </Badge>
          </div>
        )}
      </div>

      {/* Match Completion Summary */}
      {isMatchCompleted && scorecardData.innings && Array.isArray(scorecardData.innings) && (
        <Card className="mb-6 border-2 border-green-200 bg-gradient-to-r from-green-50 to-blue-50">
          <CardHeader className="text-center">
            <CardTitle className="text-2xl font-bold text-gray-800 flex items-center justify-center">
              <div className="w-4 h-4 bg-green-500 rounded-full mr-3"></div>
              Match Completed
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              {/* Team A Summary */}
              {(() => {
                const teamAInnings = scorecardData.innings.find(innings => innings.batting_team === 'A');
                return teamAInnings ? (
                  <div className="text-center p-4 bg-white rounded-lg shadow-sm border">
                    <h3 className="text-lg font-semibold text-gray-800 mb-2">{scorecardData.team_a}</h3>
                    <div className="text-3xl font-bold text-blue-600 mb-1">
                      {teamAInnings.total_runs}/{teamAInnings.total_wickets}
                    </div>
                    <div className="text-sm text-gray-600">
                      {teamAInnings.total_overs} overs
                    </div>
                    {teamAInnings.extras && (
                      <div className="text-xs text-gray-500 mt-1">
                        Extras: {teamAInnings.extras.total}
                      </div>
                    )}
                  </div>
                ) : null;
              })()}

              {/* Team B Summary */}
              {(() => {
                const teamBInnings = scorecardData.innings.find(innings => innings.batting_team === 'B');
                return teamBInnings ? (
                  <div className="text-center p-4 bg-white rounded-lg shadow-sm border">
                    <h3 className="text-lg font-semibold text-gray-800 mb-2">{scorecardData.team_b}</h3>
                    <div className="text-3xl font-bold text-blue-600 mb-1">
                      {teamBInnings.total_runs}/{teamBInnings.total_wickets}
                    </div>
                    <div className="text-sm text-gray-600">
                      {teamBInnings.total_overs} overs
                    </div>
                    {teamBInnings.extras && (
                      <div className="text-xs text-gray-500 mt-1">
                        Extras: {teamBInnings.extras.total}
                      </div>
                    )}
                  </div>
                ) : null;
              })()}
            </div>

            {/* Match Result */}
            {(() => {
              const teamAInnings = scorecardData.innings.find(innings => innings.batting_team === 'A');
              const teamBInnings = scorecardData.innings.find(innings => innings.batting_team === 'B');

              if (teamAInnings && teamBInnings) {
                const teamARuns = teamAInnings.total_runs;
                const teamBRuns = teamBInnings.total_runs;
                const winner = teamARuns > teamBRuns ? scorecardData.team_a : scorecardData.team_b;
                const margin = Math.abs(teamARuns - teamBRuns);

                return (
                  <div className="mt-6 text-center">
                    <div className="inline-block bg-white rounded-lg px-6 py-4 shadow-sm border">
                      <div className="text-lg font-semibold text-gray-800 mb-1">Result</div>
                      <div className="text-xl font-bold text-green-600">
                        {winner} won by {margin} run{margin !== 1 ? 's' : ''}
                      </div>
                    </div>
                  </div>
                );
              }
              return null;
            })()}
          </CardContent>
        </Card>
      )}

      {/* Teams Scorecard - Horizontal Layout */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4 mb-6">
        {/* Team A */}
        <Card>
          <CardHeader className="pb-3">
            <CardTitle className="text-lg">{scorecardData.team_a}</CardTitle>
          </CardHeader>
          <CardContent>
            {scorecardData.innings &&
              Array.isArray(scorecardData.innings) &&
              scorecardData.innings.length > 0 ? (
              scorecardData.innings
                .filter(
                  (innings: InningsSummary) => innings.batting_team === 'A'
                )
                .map((innings: InningsSummary) => {
                  const inningsKey = `A-${innings.innings_number}`;
                  const latestOver =
                    innings.overs &&
                      Array.isArray(innings.overs) &&
                      innings.overs.length > 0
                      ? innings.overs.reduce((latest: OverSummary, current: OverSummary) =>
                        current.over_number > latest.over_number
                          ? current
                          : latest
                      )
                      : null;
                  const isExpanded = expandedOvers[inningsKey];

                  return (
                    <div key={innings.innings_number} className="mb-3">
                      <div className="flex items-center justify-between mb-2">
                        <div className="flex items-center space-x-2">
                          <h4 className="font-medium text-sm">
                            Inn {innings.innings_number}
                          </h4>
                          <Badge
                            variant={
                              innings.status === 'in_progress'
                                ? 'default'
                                : 'secondary'
                            }
                            className={
                              innings.status === 'in_progress'
                                ? 'bg-green-600 text-white'
                                : 'bg-gray-500 text-white'
                            }
                          >
                            {innings.status === 'in_progress'
                              ? 'Live'
                              : 'Completed'}
                          </Badge>
                        </div>
                        <div className="text-xl font-bold flex items-center">
                          {scoring ? (
                            <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-gray-600 mr-2"></div>
                          ) : null}
                          {innings.total_runs}/{innings.total_wickets}
                        </div>
                      </div>
                      <div className="text-xs text-gray-600 mb-2 flex items-center">
                        {scoring ? (
                          <div className="animate-spin rounded-full h-3 w-3 border-b-2 border-gray-500 mr-2"></div>
                        ) : null}
                        {innings.total_overs} overs
                      </div>

                      {/* Extras Display */}
                      {innings.extras && (
                        <div className="text-xs text-gray-500 mb-2">
                          <span className="font-medium">Extras:</span>{' '}
                          {innings.extras.total}
                          {innings.extras.byes > 0 &&
                            ` (b ${innings.extras.byes})`}
                          {innings.extras.leg_byes > 0 &&
                            ` (lb ${innings.extras.leg_byes})`}
                          {innings.extras.wides > 0 &&
                            ` (w ${innings.extras.wides})`}
                          {innings.extras.no_balls > 0 &&
                            ` (nb ${innings.extras.no_balls})`}
                        </div>
                      )}

                      {/* Latest Over Only */}
                      {latestOver && (
                        <div className="mb-2">
                          <div className="flex items-center justify-between mb-1">
                            <span className="text-sm font-medium">
                              Over {latestOver.over_number}
                            </span>
                            <span className="text-xs text-gray-600 flex items-center">
                              {scoring ? (
                                <div className="animate-spin rounded-full h-3 w-3 border-b-2 border-gray-500 mr-1"></div>
                              ) : null}
                              {latestOver.total_runs} runs,{' '}
                              {latestOver.total_wickets} wickets
                            </span>
                          </div>
                          <div className="flex flex-wrap gap-1">
                            {latestOver.balls &&
                              Array.isArray(latestOver.balls) &&
                              latestOver.balls.length > 0 ? (
                              latestOver.balls.map(
                                (ball: BallSummary, index: number) =>
                                  renderBallCircle(ball, index)
                              )
                            ) : (
                              <div className="text-xs text-gray-400">
                                Over not started
                              </div>
                            )}
                          </div>
                        </div>
                      )}

                      {/* Show All Overs Button */}
                      {innings.overs &&
                        Array.isArray(innings.overs) &&
                        innings.overs.length > 1 && (
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => toggleExpandedOvers(inningsKey)}
                            className="text-xs h-6 px-2"
                          >
                            {isExpanded ? (
                              <>
                                <ChevronUp className="h-3 w-3 mr-1" />
                                Hide All Overs
                              </>
                            ) : (
                              <>
                                <ChevronDown className="h-3 w-3 mr-1" />
                                Show All Overs ({innings.overs.length})
                              </>
                            )}
                          </Button>
                        )}

                      {/* All Overs (Expanded) */}
                      {isExpanded &&
                        innings.overs &&
                        Array.isArray(innings.overs) &&
                        innings.overs.length > 0 && (
                          <div className="mt-2 space-y-2 border-t pt-2">
                            {[...innings.overs]
                              .sort((a: OverSummary, b: OverSummary) => b.over_number - a.over_number)
                              .map((over: OverSummary) =>
                                renderOverDetails(over)
                              )}
                          </div>
                        )}
                    </div>
                  );
                })
            ) : scorecardData.innings === null ? (
              <div className="text-sm text-gray-500 text-center py-4">
                <div className="mb-2">Match ready to start</div>
                <div className="text-xs">
                  Click &quot;Open Live Scoring&quot; to begin
                </div>
              </div>
            ) : (
              <div className="text-sm text-gray-400">No innings data</div>
            )}
          </CardContent>
        </Card>

        {/* Team B */}
        <Card>
          <CardHeader className="pb-3">
            <CardTitle className="text-lg">{scorecardData.team_b}</CardTitle>
          </CardHeader>
          <CardContent>
            {scorecardData.innings &&
              Array.isArray(scorecardData.innings) &&
              scorecardData.innings.length > 0 ? (
              scorecardData.innings
                .filter(
                  (innings: InningsSummary) => innings.batting_team === 'B'
                )
                .map((innings: InningsSummary) => {
                  const inningsKey = `B-${innings.innings_number}`;
                  const latestOver =
                    innings.overs &&
                      Array.isArray(innings.overs) &&
                      innings.overs.length > 0
                      ? innings.overs.reduce((latest: OverSummary, current: OverSummary) =>
                        current.over_number > latest.over_number
                          ? current
                          : latest
                      )
                      : null;
                  const isExpanded = expandedOvers[inningsKey];

                  return (
                    <div key={innings.innings_number} className="mb-3">
                      <div className="flex items-center justify-between mb-2">
                        <div className="flex items-center space-x-2">
                          <h4 className="font-medium text-sm">
                            Inn {innings.innings_number}
                          </h4>
                          <Badge
                            variant={
                              innings.status === 'in_progress'
                                ? 'default'
                                : 'secondary'
                            }
                            className={
                              innings.status === 'in_progress'
                                ? 'bg-green-600 text-white'
                                : 'bg-gray-500 text-white'
                            }
                          >
                            {innings.status === 'in_progress'
                              ? 'Live'
                              : 'Completed'}
                          </Badge>
                        </div>
                        <div className="text-xl font-bold flex items-center">
                          {scoring ? (
                            <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-gray-600 mr-2"></div>
                          ) : null}
                          {innings.total_runs}/{innings.total_wickets}
                        </div>
                      </div>
                      <div className="text-xs text-gray-600 mb-2 flex items-center">
                        {scoring ? (
                          <div className="animate-spin rounded-full h-3 w-3 border-b-2 border-gray-500 mr-2"></div>
                        ) : null}
                        {innings.total_overs} overs
                      </div>

                      {/* Runs Required for Team B */}
                      {innings.batting_team === 'B' && innings.status === 'in_progress' && (() => {
                        // Find Team A's total runs to calculate target
                        const teamAInn = scorecardData?.innings?.find(
                          (teamInn: InningsSummary) => teamInn.batting_team === 'A'
                        );
                        if (teamAInn) {
                          const target = teamAInn.total_runs + 1;
                          const runsRequired = target - innings.total_runs;
                          const totalBalls = (scorecardData?.total_overs || 0) * 6;
                          // Calculate balls bowled more accurately
                          const completedOvers = Math.floor(innings.total_overs);
                          const ballsInCurrentOver = latestOver?.balls?.length || 0;
                          const ballsBowled = (completedOvers * 6) + ballsInCurrentOver;
                          const ballsRemaining = totalBalls - ballsBowled;

                          return (
                            <div className="text-xs text-red-600 font-medium mb-2">
                              {runsRequired > 0 ? (
                                <>
                                  Need {runsRequired} runs in {ballsRemaining} balls
                                </>
                              ) : (
                                <>
                                  Won by {Math.abs(runsRequired)} runs
                                </>
                              )}
                            </div>
                          );
                        }
                        return null;
                      })()}

                      {/* Extras Display */}
                      {innings.extras && (
                        <div className="text-xs text-gray-500 mb-2">
                          <span className="font-medium">Extras:</span>{' '}
                          {innings.extras.total}
                          {innings.extras.byes > 0 &&
                            ` (b ${innings.extras.byes})`}
                          {innings.extras.leg_byes > 0 &&
                            ` (lb ${innings.extras.leg_byes})`}
                          {innings.extras.wides > 0 &&
                            ` (w ${innings.extras.wides})`}
                          {innings.extras.no_balls > 0 &&
                            ` (nb ${innings.extras.no_balls})`}
                        </div>
                      )}

                      {/* Latest Over Only */}
                      {latestOver && (
                        <div className="mb-2">
                          <div className="flex items-center justify-between mb-1">
                            <span className="text-sm font-medium">
                              Over {latestOver.over_number}
                            </span>
                            <span className="text-xs text-gray-600 flex items-center">
                              {scoring ? (
                                <div className="animate-spin rounded-full h-3 w-3 border-b-2 border-gray-500 mr-1"></div>
                              ) : null}
                              {latestOver.total_runs} runs,{' '}
                              {latestOver.total_wickets} wickets
                            </span>
                          </div>
                          <div className="flex flex-wrap gap-1">
                            {latestOver.balls &&
                              Array.isArray(latestOver.balls) &&
                              latestOver.balls.length > 0 ? (
                              latestOver.balls.map(
                                (ball: BallSummary, index: number) =>
                                  renderBallCircle(ball, index)
                              )
                            ) : (
                              <div className="text-xs text-gray-400">
                                Over not started
                              </div>
                            )}
                          </div>
                        </div>
                      )}

                      {/* Show All Overs Button */}
                      {innings.overs &&
                        Array.isArray(innings.overs) &&
                        innings.overs.length > 1 && (
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => toggleExpandedOvers(inningsKey)}
                            className="text-xs h-6 px-2"
                          >
                            {isExpanded ? (
                              <>
                                <ChevronUp className="h-3 w-3 mr-1" />
                                Hide All Overs
                              </>
                            ) : (
                              <>
                                <ChevronDown className="h-3 w-3 mr-1" />
                                Show All Overs ({innings.overs.length})
                              </>
                            )}
                          </Button>
                        )}

                      {/* All Overs (Expanded) */}
                      {isExpanded &&
                        innings.overs &&
                        Array.isArray(innings.overs) &&
                        innings.overs.length > 0 && (
                          <div className="mt-2 space-y-2 border-t pt-2">
                            {[...innings.overs]
                              .sort((a: OverSummary, b: OverSummary) => b.over_number - a.over_number)
                              .map((over: OverSummary) =>
                                renderOverDetails(over)
                              )}
                          </div>
                        )}
                    </div>
                  );
                })
            ) : scorecardData.innings === null ? (
              <div className="text-sm text-gray-500 text-center py-4">
                <div className="mb-2">Match ready to start</div>
                <div className="text-xs">
                  Click &quot;Open Live Scoring&quot; to begin
                </div>
              </div>
            ) : (
              <div className="text-sm text-gray-400">No innings data</div>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Live Scoring Interface */}
      {showLiveScoring && isScoringAvailable && (
        <Card className="border border-gray-200 shadow-lg">
          <CardHeader className="border-b border-gray-200">
            <CardTitle className="flex items-center justify-between">
              <div className="flex items-center space-x-3">
                <div className="w-3 h-3 bg-green-500 rounded-full animate-pulse"></div>
                <span className="text-xl font-bold text-gray-800">Live Score</span>
                {scoring && (
                  <div className="ml-3">
                    <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-green-600"></div>
                  </div>
                )}
              </div>
              <div className="flex items-center space-x-3">
                <Badge variant="default" className="bg-green-600 text-white font-semibold px-1.5 py-1">
                  Inn {currentInn}
                </Badge>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setShowLiveScoring(false)}
                  className="h-9 w-9 p-0 border-red-300 hover:bg-red-50 hover:border-red-400"
                  disabled={scoring}
                >
                  <X className="h-4 w-4 text-red-600" />
                </Button>
              </div>
            </CardTitle>
            {/* Current Inn Info */}
            {scorecardData.innings && Array.isArray(scorecardData.innings) && (
              <div className="mt-2 text-sm text-gray-600">
                {scorecardData.innings.map((innings: InningsSummary) => (
                  <div
                    key={innings.innings_number}
                    className="flex items-center space-x-2"
                  >
                    <span>Inn {innings.innings_number}:</span>
                    <Badge
                      variant={
                        innings.status === 'in_progress'
                          ? 'default'
                          : 'secondary'
                      }
                      className={
                        innings.status === 'in_progress' ? 'bg-green-600' : ''
                      }
                    >
                      {innings.status === 'in_progress'
                        ? 'In Progress'
                        : 'Completed'}
                    </Badge>
                    <span className="flex items-center">
                      {scoring ? (
                        <div className="animate-spin rounded-full h-3 w-3 border-b-2 border-gray-500 mr-2"></div>
                      ) : null}
                      {innings.batting_team === 'A'
                        ? scorecardData.team_a
                        : scorecardData.team_b}
                      - {innings.total_runs}/{innings.total_wickets} (
                      {innings.total_overs} overs)
                    </span>
                  </div>
                ))}
              </div>
            )}
          </CardHeader>
          <CardContent>
            {/* Runs Actions */}
            <div className="mb-8">
              <div className="flex items-center mb-4">
                <div className="w-2 h-2 bg-green-500 rounded-full mr-2"></div>
                <h4 className="font-semibold text-lg text-gray-800">Runs</h4>
              </div>
              <div className="grid grid-cols-3 gap-3">
                {[0, 1, 2, 3, 4, 6].map(runs => (
                  <Button
                    key={runs}
                    onClick={() => handleBallScore(runs, 'good')}
                    size="lg"
                    variant={
                      runs === 4
                        ? 'default'
                        : runs === 6
                          ? 'secondary'
                          : 'outline'
                    }
                    className={`h-14 text-lg font-bold transition-all duration-200 ${runs === 4
                      ? 'bg-blue-600 hover:bg-blue-700 text-white shadow-lg hover:shadow-xl'
                      : runs === 6
                        ? 'bg-purple-600 hover:bg-purple-700 text-white shadow-lg hover:shadow-xl'
                        : 'border-2 hover:border-green-400 hover:bg-green-50 hover:text-green-700'
                      }`}
                    disabled={scoring}
                  >
                    {scoring ? (
                      <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-white"></div>
                    ) : (
                      <span className="text-xl font-bold">{runs}</span>
                    )}
                  </Button>
                ))}
              </div>
            </div>

            {/* Extras Actions */}
            <div className="mb-8">
              <div className="flex items-center mb-4">
                <div className="w-2 h-2 bg-yellow-500 rounded-full mr-2"></div>
                <h4 className="font-semibold text-lg text-gray-800">Extras</h4>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <Button
                  onClick={() => handleBallScore(1, 'wide')}
                  size="lg"
                  variant="outline"
                  className="h-14 border-2 border-yellow-500 text-yellow-700 hover:bg-yellow-50 hover:border-yellow-600 transition-all duration-200 font-semibold text-lg"
                  disabled={scoring}
                >
                  {scoring ? (
                    <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-yellow-700"></div>
                  ) : (
                    'Wide'
                  )}
                </Button>
                <Button
                  onClick={() => handleBallScore(1, 'no_ball')}
                  size="lg"
                  variant="outline"
                  className="h-14 border-2 border-orange-500 text-orange-700 hover:bg-orange-50 hover:border-orange-600 transition-all duration-200 font-semibold text-lg"
                  disabled={scoring}
                >
                  {scoring ? (
                    <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-orange-700"></div>
                  ) : (
                    'No Ball'
                  )}
                </Button>
              </div>
            </div>

            {/* Wicket Actions */}
            <div className="mb-8">
              <div className="flex items-center mb-4">
                <div className="w-2 h-2 bg-red-500 rounded-full mr-2"></div>
                <h4 className="font-semibold text-lg text-gray-800">Wickets</h4>
              </div>
              <div className="grid grid-cols-2 gap-3">
                {[
                  'bowled',
                  'caught',
                  'lbw',
                  'run_out',
                  'stumped',
                  'hit_wicket',
                ].map(wicketType => (
                  <Button
                    key={wicketType}
                    onClick={() => handleBallScore(0, wicketType)}
                    size="lg"
                    variant="destructive"
                    className="h-12 bg-red-600 hover:bg-red-700 text-white font-semibold transition-all duration-200 shadow-lg hover:shadow-xl"
                    disabled={scoring}
                  >
                    {scoring ? (
                      <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-white"></div>
                    ) : (
                      <span className="text-sm font-bold">
                        {wicketType.replace('_', ' ').toUpperCase()}
                      </span>
                    )}
                  </Button>
                ))}
              </div>
            </div>

            {/* Byes Selection - Moved to Bottom */}
            <div className="border-t border-gray-200 pt-6">
              <div className="flex items-center mb-4">
                <div className="w-2 h-2 bg-blue-500 rounded-full mr-2"></div>
                <h4 className="font-semibold text-lg text-gray-800">
                  Byes (Optional)
                </h4>
              </div>
              <div className="flex flex-col items-center space-y-4">
                <div className="flex flex-wrap justify-center gap-3 max-w-full">
                  {[0, 1, 2, 3, 4, 5, 6].map(byes => (
                    <button
                      key={byes}
                      onClick={() => handleByesChange(byes)}
                      disabled={scoring}
                      className={`w-12 h-12 rounded-full border-2 flex items-center justify-center text-lg font-bold transition-all duration-200 ${byes === currentByes
                        ? 'border-blue-500 bg-blue-100 text-blue-700 shadow-lg scale-110'
                        : 'border-gray-300 bg-white text-gray-500 hover:bg-gray-50 hover:border-gray-400 hover:scale-105'
                        } ${scoring ? 'opacity-50 cursor-not-allowed' : ''}`}
                    >
                      {byes}
                    </button>
                  ))}
                </div>
                <div className="text-sm text-gray-600 text-center font-medium">
                  {currentByes > 0
                    ? `+${currentByes} byes selected`
                    : 'No byes selected'}
                </div>
              </div>
            </div>

            {/* Undo Ball Action - Only show if there are balls to undo */}
            {(() => {
              const currentInnData = scorecardData?.innings?.find(
                innings => innings.innings_number === currentInn
              );
              const totalBalls = currentInnData?.overs?.reduce((total, over) => {
                return total + (over.balls ? over.balls.length : 0);
              }, 0) || 0;

              // Only show undo button if there are balls in the innings
              if (totalBalls === 0) {
                return null;
              }

              return (
                <div className="border-t border-gray-200 pt-6 mt-6">
                  <div className="flex justify-center">
                    <Button
                      onClick={handleUndoBall}
                      variant="outline"
                      size="lg"
                      className={`h-12 border-2 border-red-500 text-red-700 hover:bg-red-50 hover:border-red-600 transition-all duration-200 font-semibold shadow-lg hover:shadow-xl ${isFirstBallOfInn() ? 'opacity-50 cursor-not-allowed' : ''
                        }`}
                      disabled={scoring || isFirstBallOfInn()}
                    >
                      {scoring ? (
                        <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-red-700"></div>
                      ) : (
                        <>
                          <Undo2 className="h-5 w-5 mr-2" />
                          Undo Last Ball
                        </>
                      )}
                    </Button>
                  </div>
                  {isFirstBallOfInn() && (
                    <div className="text-center mt-3">
                      <span className="text-sm text-gray-500 font-medium">
                        Cannot undo - first ball of innings
                      </span>
                    </div>
                  )}
                </div>
              );
            })()}
          </CardContent>
        </Card>
      )}
    </div>
  );
}
