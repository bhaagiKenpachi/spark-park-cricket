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
  const [currentInnings, setCurrentInnings] = useState(1);
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
  const isBothInningsCompleted = scorecard?.innings &&
    scorecard.innings.length >= 2 &&
    scorecard.innings.every(innings => innings.status === 'completed');

  // Determine if scoring should be available
  const isScoringAvailable = isOwner && !isMatchCompleted && !isBothInningsCompleted;

  console.log('=== SCORECARD VIEW OWNERSHIP CHECK ===');
  console.log('Is Authenticated:', isAuthenticated);
  console.log('Current User:', currentUser);
  console.log('Current User ID:', currentUser?.id);
  console.log('Series Created By:', seriesCreatedBy);
  console.log('Is Owner:', isOwner);
  console.log('Match ID:', matchId);
  console.log('Match Status:', scorecard?.match_status);
  console.log('Is Match Completed:', isMatchCompleted);
  console.log('Is Both Innings Completed:', isBothInningsCompleted);
  console.log('Is Scoring Available:', isScoringAvailable);

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
    console.log('=== SCORECARD DATA CHANGE ===');
    console.log('Scorecard:', scorecard);
    console.log('Scorecard innings:', scorecard?.innings);

    if (scorecard?.innings && Array.isArray(scorecard.innings) && scorecard.innings.length > 0) {
      const currentInningsData = scorecard.innings.find(
        innings => innings.status === 'in_progress'
      );
      console.log('Current innings data:', currentInningsData);
      if (currentInningsData) {
        console.log('Setting current innings to:', currentInningsData.innings_number);
        setCurrentInnings(currentInningsData.innings_number);

        // Debug: Log all balls in current innings
        const totalBalls = currentInningsData.overs?.reduce((total, over) => {
          console.log('Over:', over.over_number, 'Balls:', over.balls?.length || 0);
          return total + (over.balls ? over.balls.length : 0);
        }, 0) || 0;
        console.log('Total balls in current innings:', totalBalls);
      }
    } else if (scorecard?.innings === null || (Array.isArray(scorecard?.innings) && scorecard.innings.length === 0)) {
      // If no innings exist yet, start with innings 1
      console.log('No innings data, setting current innings to 1');
      setCurrentInnings(1);
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
    console.log('=== HANDLE START SCORING ===');
    console.log('Match ID:', matchId);
    console.log('Is Owner:', isOwner);
    console.log('Is Authenticated:', isAuthenticated);
    console.log('Current User:', currentUser);
    console.log('Series Created By:', seriesCreatedBy);
    console.log('Match Status:', scorecardData?.match_status);
    console.log('Is Scoring Available:', isScoringAvailable);
    console.log('Current cookies:', document.cookie);

    // Check if scoring is available (ownership + match not completed)
    if (!isScoringAvailable) {
      if (!isOwner) {
        console.log('❌ Not owner, blocking start scoring');
        setScoringMessage('Only the series creator can start scoring.');
      } else if (isMatchCompleted) {
        console.log('❌ Match completed, blocking start scoring');
        setScoringMessage('Cannot score on completed match.');
      } else if (isBothInningsCompleted) {
        console.log('❌ Both innings completed, blocking start scoring');
        setScoringMessage('Cannot score when both innings are completed.');
      }
      setTimeout(() => setScoringMessage(null), 3000);
      return;
    }

    console.log('✅ Scoring available, proceeding with start scoring');

    // If match is already live, just show the interface without calling the API
    if (scorecardData?.match_status === 'live') {
      console.log('Match already live, showing interface without API call');
      setShowLiveScoring(true);
      setScoringMessage('Live scoring interface opened!');
      setTimeout(() => setScoringMessage(null), 3000);
    } else {
      console.log('Match not live, calling startScoringRequest API');
      // Only call the API if match is not live yet
      dispatch(startScoringRequest(matchId));
      setShowLiveScoring(true);
      setScoringMessage('Live scoring started!');
      setTimeout(() => setScoringMessage(null), 3000);
    }
  };

  const handleBallScore = (runs: number, ballType: string) => {
    console.log('=== HANDLE BALL SCORE ===');
    console.log('Runs:', runs);
    console.log('Ball Type:', ballType);
    console.log('Current Innings:', currentInnings);
    console.log('Is Owner:', isOwner);
    console.log('Is Authenticated:', isAuthenticated);
    console.log('Current User:', currentUser);
    console.log('Series Created By:', seriesCreatedBy);
    console.log('Is Scoring Available:', isScoringAvailable);
    console.log('Current cookies:', document.cookie);

    // Debug: Check current scorecard state before scoring
    const currentInningsDataBefore = scorecardData?.innings?.find(
      innings => innings.innings_number === currentInnings
    );
    console.log('Current innings data before scoring:', currentInningsDataBefore);
    if (currentInningsDataBefore) {
      const totalBallsBefore = currentInningsDataBefore.overs?.reduce((total, over) => {
        return total + (over.balls ? over.balls.length : 0);
      }, 0) || 0;
      console.log('Total balls before scoring:', totalBallsBefore);
    }

    // Check if scoring is available (ownership + match not completed)
    if (!isScoringAvailable) {
      if (!isOwner) {
        console.log('❌ Not owner, blocking ball scoring');
        setScoringMessage('Only the series creator can score balls.');
      } else if (isMatchCompleted) {
        console.log('❌ Match completed, blocking ball scoring');
        setScoringMessage('Cannot score on completed match.');
      } else if (isBothInningsCompleted) {
        console.log('❌ Both innings completed, blocking ball scoring');
        setScoringMessage('Cannot score when both innings are completed.');
      }
      setTimeout(() => setScoringMessage(null), 3000);
      return;
    }

    console.log('✅ Scoring available, proceeding with ball scoring');

    // Check if current innings is still in progress
    const currentInningsDataForScoring = scorecardData?.innings?.find(
      innings => innings.innings_number === currentInnings
    );

    // If no innings exist yet (null or empty array), allow scoring to create the first innings
    if (scorecardData?.innings === null || (Array.isArray(scorecardData?.innings) && scorecardData.innings.length === 0)) {
      // Allow scoring - this will create the first innings
    } else if (currentInningsDataForScoring?.status !== 'in_progress') {
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
      innings_number: currentInnings,
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
  const isFirstBallOfInnings = () => {
    console.log('=== IS FIRST BALL CHECK ===');
    console.log('Current innings:', currentInnings);
    console.log('Scorecard data:', scorecardData);

    const currentInningsData = scorecardData?.innings?.find(
      innings => innings.innings_number === currentInnings
    );
    console.log('Current innings data:', currentInningsData);

    if (!currentInningsData) {
      console.log('No innings data found, returning true (first ball)');
      return true; // If no innings data, consider it first ball
    }

    // Count total balls across all overs in this innings
    const totalBalls = currentInningsData.overs?.reduce((total, over) => {
      console.log('Over', over.over_number, 'has', over.balls?.length || 0, 'balls');
      return total + (over.balls ? over.balls.length : 0);
    }, 0) || 0;

    console.log('Total balls in innings:', totalBalls);
    console.log('Is first ball (totalBalls === 1):', totalBalls === 1);

    // If there's exactly 1 ball, it's the first ball (and we can't undo it)
    return totalBalls === 1;
  };

  const handleUndoBall = () => {
    console.log('=== HANDLE UNDO BALL ===');
    console.log('Match ID:', matchId);
    console.log('Current Innings:', currentInnings);
    console.log('Is Owner:', isOwner);
    console.log('Is Authenticated:', isAuthenticated);
    console.log('Current User:', currentUser);
    console.log('Series Created By:', seriesCreatedBy);
    console.log('Is Scoring Available:', isScoringAvailable);

    // Check if scoring is available (ownership + match not completed)
    if (!isScoringAvailable) {
      if (!isOwner) {
        console.log('❌ Not owner, blocking undo ball');
        setScoringMessage('Only the series creator can undo balls.');
      } else if (isMatchCompleted) {
        console.log('❌ Match completed, blocking undo ball');
        setScoringMessage('Cannot undo ball on completed match.');
      } else if (isBothInningsCompleted) {
        console.log('❌ Both innings completed, blocking undo ball');
        setScoringMessage('Cannot undo ball when both innings are completed.');
      }
      setTimeout(() => setScoringMessage(null), 3000);
      return;
    }

    // Check if current innings is still in progress
    const currentInningsDataForUndo = scorecardData?.innings?.find(
      innings => innings.innings_number === currentInnings
    );

    if (!currentInningsDataForUndo || currentInningsDataForUndo.status !== 'in_progress') {
      setScoringMessage(
        'Cannot undo ball on completed innings. Please check innings status.'
      );
      setTimeout(() => setScoringMessage(null), 5000);
      return;
    }

    // Check if it's the first ball of the innings
    if (isFirstBallOfInnings()) {
      console.log('❌ First ball of innings, blocking undo ball');
      setScoringMessage('Cannot undo ball - this is the first ball of the innings.');
      setTimeout(() => setScoringMessage(null), 3000);
      return;
    }

    console.log('✅ Undo ball available, proceeding with undo');
    dispatch(undoBallThunk({ matchId, inningsNumber: currentInnings }));
    setScoringMessage('Undoing last ball...');
    setTimeout(() => setScoringMessage(null), 2000);
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

    // Handle display with byes
    let displayWithByes: string;
    if (ball.byes > 0) {
      // For balls with byes, show "B" + runs + "+" + byes
      if (ball.ball_type === 'WIDE' || ball.ball_type === 'wide' || ball.ball_type === 'NO_BALL' || ball.ball_type === 'no_ball') {
        // For wide/no ball with byes, show ball type + "+" + byes
        displayWithByes = `${display}+${ball.byes}`;
      } else {
        // For good balls with byes, show "B" + runs + "+" + byes
        displayWithByes = `B${ball.runs}+${ball.byes}`;
      }
    } else {
      displayWithByes = display;
    }

    return (
      <div
        key={index}
        className={`w-8 h-8 rounded-full border-2 flex items-center justify-center text-xs font-medium ${isWicket
          ? 'border-red-500 bg-red-100 text-red-700'
          : ball.ball_type === 'WIDE' || ball.ball_type === 'wide' || ball.ball_type === 'NO_BALL' || ball.ball_type === 'no_ball' || ball.run_type === 'LB' || ball.byes > 0
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
        {displayWithByes}
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
          <div className="text-xs text-gray-400">No balls</div>
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
        {isBothInningsCompleted && !isMatchCompleted && (
          <div className="mt-3">
            <Badge variant="secondary" className="bg-gray-500 text-white">
              Both Innings Completed - Scoring Not Available
            </Badge>
          </div>
        )}
      </div>

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
                      ? innings.overs.reduce((latest, current) =>
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
                            Innings {innings.innings_number}
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
                              Latest Over {latestOver.over_number}
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
                                No balls
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
                            {innings.overs.map((over: OverSummary) =>
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
                      ? innings.overs.reduce((latest, current) =>
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
                            Innings {innings.innings_number}
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
                              Latest Over {latestOver.over_number}
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
                                No balls
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
                            {innings.overs.map((over: OverSummary) =>
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
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center justify-between">
              <div className="flex items-center">
                <span>Live Scoring</span>
                {scoring && (
                  <div className="ml-3">
                    <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-green-600"></div>
                  </div>
                )}
              </div>
              <div className="flex items-center space-x-2">
                <Badge variant="default" className="bg-green-600">
                  Innings {currentInnings}
                </Badge>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setShowLiveScoring(false)}
                  className="h-8 w-8 p-0"
                  disabled={scoring}
                >
                  <X className="h-4 w-4" />
                </Button>
              </div>
            </CardTitle>
            {/* Current Innings Info */}
            {scorecardData.innings && Array.isArray(scorecardData.innings) && (
              <div className="mt-2 text-sm text-gray-600">
                {scorecardData.innings.map((innings: InningsSummary) => (
                  <div
                    key={innings.innings_number}
                    className="flex items-center space-x-2"
                  >
                    <span>Innings {innings.innings_number}:</span>
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
            <div className="mb-6">
              <h4 className="font-medium mb-3 text-gray-700">Runs</h4>
              <div className="grid grid-cols-4 gap-2">
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
                    className={
                      runs === 4
                        ? 'bg-blue-600 hover:bg-blue-700'
                        : runs === 6
                          ? 'bg-purple-600 hover:bg-purple-700 text-white'
                          : ''
                    }
                    disabled={scoring}
                  >
                    {scoring ? (
                      <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
                    ) : (
                      runs
                    )}
                  </Button>
                ))}
              </div>
            </div>

            {/* Extras Actions */}
            <div className="mb-6">
              <h4 className="font-medium mb-3 text-gray-700">Extras</h4>
              <div className="grid grid-cols-2 gap-2">
                <Button
                  onClick={() => handleBallScore(1, 'wide')}
                  size="lg"
                  variant="outline"
                  className="border-yellow-500 text-yellow-700 hover:bg-yellow-50"
                  disabled={scoring}
                >
                  {scoring ? (
                    <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-yellow-700"></div>
                  ) : (
                    'Wide'
                  )}
                </Button>
                <Button
                  onClick={() => handleBallScore(1, 'no_ball')}
                  size="lg"
                  variant="outline"
                  className="border-orange-500 text-orange-700 hover:bg-orange-50"
                  disabled={scoring}
                >
                  {scoring ? (
                    <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-orange-700"></div>
                  ) : (
                    'No Ball'
                  )}
                </Button>
              </div>
            </div>

            {/* Wicket Actions */}
            <div className="mb-6">
              <h4 className="font-medium mb-3 text-gray-700">Wickets</h4>
              <div className="grid grid-cols-2 gap-2">
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
                    disabled={scoring}
                  >
                    {scoring ? (
                      <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
                    ) : (
                      wicketType.replace('_', ' ').toUpperCase()
                    )}
                  </Button>
                ))}
              </div>
            </div>

            {/* Byes Selection - Moved to Bottom */}
            <div className="border-t pt-4">
              <h4 className="font-medium mb-3 text-gray-700">
                Byes (Optional)
              </h4>
              <div className="flex items-center justify-center space-x-2">
                <div className="flex space-x-1">
                  {[0, 1, 2, 3, 4].map(byes => (
                    <button
                      key={byes}
                      onClick={() => handleByesChange(byes)}
                      disabled={scoring}
                      className={`w-10 h-10 rounded-full border-2 flex items-center justify-center text-sm font-medium transition-colors ${byes === currentByes
                        ? 'border-blue-500 bg-blue-100 text-blue-700'
                        : 'border-gray-300 bg-white text-gray-500 hover:bg-gray-50'
                        } ${scoring ? 'opacity-50 cursor-not-allowed' : ''}`}
                    >
                      {byes}
                    </button>
                  ))}
                </div>
                <div className="ml-4 text-sm text-gray-600">
                  {currentByes > 0
                    ? `+${currentByes} byes selected`
                    : 'No byes'}
                </div>
              </div>
            </div>

            {/* Undo Ball Action - Only show if there are balls to undo */}
            {(() => {
              const currentInningsData = scorecardData?.innings?.find(
                innings => innings.innings_number === currentInnings
              );
              const totalBalls = currentInningsData?.overs?.reduce((total, over) => {
                return total + (over.balls ? over.balls.length : 0);
              }, 0) || 0;

              // Only show undo button if there are balls in the innings
              if (totalBalls === 0) {
                return null;
              }

              return (
                <div className="border-t pt-4 mt-4">
                  <div className="flex justify-center">
                    <Button
                      onClick={handleUndoBall}
                      variant="outline"
                      size="lg"
                      className={`border-red-500 text-red-700 hover:bg-red-50 ${isFirstBallOfInnings() ? 'opacity-50 cursor-not-allowed' : ''
                        }`}
                      disabled={scoring || isFirstBallOfInnings()}
                    >
                      {scoring ? (
                        <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-red-700"></div>
                      ) : (
                        <>
                          <Undo2 className="h-4 w-4 mr-2" />
                          Undo Last Ball
                        </>
                      )}
                    </Button>
                  </div>
                  {isFirstBallOfInnings() && (
                    <div className="text-center mt-2">
                      <span className="text-xs text-gray-500">
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
