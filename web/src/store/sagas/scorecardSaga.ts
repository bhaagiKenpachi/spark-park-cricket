import {
  call,
  put,
  takeEvery,
  takeLatest,
  take,
  select,
  fork,
  CallEffect,
  PutEffect,
  SelectEffect,
  TakeEffect,
} from 'redux-saga/effects';
import { RootState } from '../store';
import {
  fetchScorecardRequest,
  fetchScorecardSuccess,
  fetchScorecardFailure,
  startScoringRequest,
  startScoringSuccess,
  startScoringFailure,
  addBallRequest,
  addBallSuccess,
  addBallFailure,
  undoBallThunk,
  fetchInningsRequest,
  fetchInningsSuccess,
  fetchInningsFailure,
  fetchInningsScoreSummaryThunk,
  fetchLatestOverThunk,
  clearInningsTransition,
  triggerInningsTransition,
  ScorecardResponse,
  InningsSummary,
  OverSummary,
} from '../reducers/scorecardSlice';
import { ApiService, ApiError, ApiResponse } from '@/services/api';
import { graphqlService } from '@/services/graphqlService';
import { LiveScorecardResponse, InningsDetailsResponse } from '@/lib/graphql';

export function* fetchScorecardSaga(
  action: ReturnType<typeof fetchScorecardRequest>
): Generator<
  CallEffect | PutEffect,
  void,
  { success: boolean; data?: LiveScorecardResponse['liveScorecard']; error?: string }
> {
  try {
    const response = yield call(
      graphqlService.getLiveScorecard.bind(graphqlService),
      action.payload
    );

    if (response.success && response.data) {
      // Transform GraphQL response to match ScorecardResponse format
      const scorecardData: ScorecardResponse = {
        match_id: response.data.match_id,
        match_number: response.data.match_number,
        series_name: response.data.series_name,
        team_a: response.data.team_a,
        team_b: response.data.team_b,
        total_overs: response.data.total_overs,
        toss_winner: response.data.toss_winner as 'A' | 'B',
        toss_type: response.data.toss_type as 'H' | 'T',
        current_innings: response.data.current_innings,
        match_status: response.data.match_status,
        innings: response.data.innings as InningsSummary[] || [],
      };
      yield put(fetchScorecardSuccess(scorecardData));
    } else {
      yield put(fetchScorecardFailure(response.error || 'Failed to fetch scorecard'));
    }
  } catch (error) {
    const errorMessage =
      error instanceof Error ? error.message : 'Failed to fetch scorecard';
    yield put(fetchScorecardFailure(errorMessage));
  }
}

export function* startScoringSaga(
  action: ReturnType<typeof startScoringRequest>
): Generator<
  CallEffect | PutEffect,
  void,
  ApiResponse<{ message: string; match_id: string }>
> {
  try {
    const apiService = new ApiService();
    yield call(apiService.startScoring.bind(apiService), action.payload);
    yield put(startScoringSuccess());
    // Refresh scorecard after starting scoring
    yield put(fetchScorecardRequest(action.payload));
  } catch (error) {
    const errorMessage =
      error instanceof ApiError ? error.message : 'Failed to start scoring';
    yield put(startScoringFailure(errorMessage));
  }
}

export function* addBallSaga(
  action: ReturnType<typeof addBallRequest>
): Generator<
  CallEffect | PutEffect,
  void,
  ApiResponse<{
    message: string;
    match_id: string;
    innings_number: number;
    ball_type: string;
    run_type: string;
    runs: number;
    byes: number;
    is_wicket: boolean;
  }>
> {
  try {
    const apiService = new ApiService();
    const ballEvent = action.payload;

    // Use the addBall method directly with the complete ball event
    yield call(apiService.addBall.bind(apiService), ballEvent);
    yield put(addBallSuccess());

    // Check if innings data exists in state before using GraphQL methods
    const currentInningsData = yield select((state: RootState) =>
      state.scorecard.scorecard?.innings?.find(
        innings => innings.innings_number === ballEvent.innings_number
      )
    );

    if (currentInningsData) {
      // Use optimized GraphQL methods to refresh specific data after adding ball
      try {
        // Fetch innings score summary for the current innings
        yield put(fetchInningsScoreSummaryThunk({
          matchId: ballEvent.match_id,
          inningsNumber: ballEvent.innings_number
        }));

        // Fetch latest over for the current innings
        yield put(fetchLatestOverThunk({
          matchId: ballEvent.match_id,
          inningsNumber: ballEvent.innings_number
        }));

        // Note: Full scorecard refresh is not needed on every ball
        // It will be called automatically when innings transitions occur
        // (e.g., when first innings completes and second innings starts)
      } catch (error) {
        console.warn('GraphQL refresh failed, falling back to full scorecard fetch:', error);
        // Fallback to full scorecard fetch if GraphQL fails
        yield put(fetchScorecardRequest(ballEvent.match_id));
      }
    } else {
      // If no innings data exists in state, do a full scorecard refresh
      // This happens for the first ball of a new innings
      console.log('No innings data in state, doing full scorecard refresh for first ball');
      yield put(fetchScorecardRequest(ballEvent.match_id));
    }

  } catch (error) {
    const errorMessage =
      error instanceof ApiError ? error.message : 'Failed to add ball';
    yield put(addBallFailure(errorMessage));
  }
}

export function* undoBallSaga(
  action: ReturnType<typeof undoBallThunk.fulfilled>
): Generator<CallEffect | PutEffect, void, any> {
  try {
    console.log('=== UNDO BALL SAGA ===');
    const { matchId, inningsNumber } = action.payload;

    // Check if innings data exists in state before using GraphQL methods
    const currentInningsData = yield select((state: RootState) =>
      state.scorecard.scorecard?.innings?.find(
        innings => innings.innings_number === inningsNumber
      )
    );

    if (currentInningsData) {
      // Use optimized GraphQL methods to refresh specific data after undoing ball
      try {
        // Fetch innings score summary for the current innings
        yield put(fetchInningsScoreSummaryThunk({
          matchId,
          inningsNumber
        }));

        // Fetch latest over for the current innings
        yield put(fetchLatestOverThunk({
          matchId,
          inningsNumber
        }));

        console.log('GraphQL refresh completed for undo ball');
      } catch (error) {
        console.warn('GraphQL refresh failed for undo ball, falling back to full scorecard fetch:', error);
        // Fallback to full scorecard fetch if GraphQL fails
        yield put(fetchScorecardRequest(matchId));
      }
    } else {
      // If no innings data exists in state, do a full scorecard refresh
      console.log('No innings data in state, doing full scorecard refresh for undo ball');
      yield put(fetchScorecardRequest(matchId));
    }
  } catch (error) {
    console.error('Error in undo ball saga:', error);
  }
}

export function* fetchInningsSaga(
  action: ReturnType<typeof fetchInningsRequest>
): Generator<
  CallEffect | PutEffect,
  void,
  { success: boolean; data?: InningsDetailsResponse['inningsDetails']; error?: string }
> {
  try {
    const { matchId, inningsNumber } = action.payload;
    const response = yield call(
      graphqlService.getInningsDetails.bind(graphqlService),
      matchId,
      inningsNumber
    );

    if (response.success && response.data) {
      // Transform GraphQL response to match InningsSummary format
      const inningsData: InningsSummary = {
        innings_number: response.data.innings_number,
        batting_team: response.data.batting_team as 'A' | 'B',
        total_runs: response.data.total_runs,
        total_wickets: response.data.total_wickets,
        total_overs: response.data.total_overs,
        total_balls: response.data.total_balls,
        status: response.data.status,
        extras: response.data.extras,
        overs: response.data.overs as OverSummary[] || [],
      };
      yield put(fetchInningsSuccess(inningsData));
    } else {
      yield put(fetchInningsFailure(response.error || 'Failed to fetch innings'));
    }
  } catch (error) {
    const errorMessage =
      error instanceof Error ? error.message : 'Failed to fetch innings';
    yield put(fetchInningsFailure(errorMessage));
  }
}

// Saga to handle innings transitions
export function* handleInningsTransitionSaga(): Generator<
  TakeEffect | SelectEffect | PutEffect,
  void,
  { scorecard: { scorecard: ScorecardResponse | null } }
> {
  while (true) {
    // Watch for innings transition detection
    yield take(triggerInningsTransition.type);

    // Get the current state to find the match ID
    const state = yield select();
    const scorecard = state.scorecard?.scorecard;

    if (scorecard?.match_id) {
      // Fetch full scorecard to get the new innings information
      yield put(fetchScorecardRequest(scorecard.match_id));

      // Clear the innings transition flag
      yield put(clearInningsTransition());
    }
  }
}

// Saga to check for innings transitions after thunk completion
export function* checkInningsTransitionSaga(): Generator<
  TakeEffect | SelectEffect | PutEffect,
  void,
  { scorecard: { inningsTransitionDetected: boolean; scorecard: ScorecardResponse | null } }
> {
  while (true) {
    // Wait for innings score summary thunk to complete
    yield take(fetchInningsScoreSummaryThunk.fulfilled.type);

    // Check if innings transition was detected
    const state = yield select();
    if (state.scorecard?.inningsTransitionDetected) {
      const scorecard = state.scorecard?.scorecard;
      if (scorecard?.match_id) {
        // Fetch full scorecard to get the new innings information
        yield put(fetchScorecardRequest(scorecard.match_id));

        // Clear the innings transition flag
        yield put(clearInningsTransition());
      }
    }
  }
}

export function* scorecardSaga() {
  yield takeLatest(fetchScorecardRequest.type, fetchScorecardSaga);
  yield takeEvery(startScoringRequest.type, startScoringSaga);
  yield takeEvery(addBallRequest.type, addBallSaga);
  yield takeEvery(undoBallThunk.fulfilled.type, undoBallSaga);
  yield takeLatest(fetchInningsRequest.type, fetchInningsSaga);
  yield fork(handleInningsTransitionSaga);
  yield fork(checkInningsTransitionSaga);
}
