import { call, put, takeEvery, takeLatest, CallEffect, PutEffect } from 'redux-saga/effects';
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
    fetchInningsRequest,
    fetchInningsSuccess,
    fetchInningsFailure,
    ScorecardResponse,
    InningsSummary,
} from '../reducers/scorecardSlice';
import { ApiService, ApiError, ApiResponse } from '@/services/api';

export function* fetchScorecardSaga(action: ReturnType<typeof fetchScorecardRequest>): Generator<CallEffect | PutEffect, void, ApiResponse<ScorecardResponse>> {
    try {
        const apiService = new ApiService();
        const response = yield call(apiService.getScorecard.bind(apiService), action.payload);
        yield put(fetchScorecardSuccess(response.data));
    } catch (error) {
        const errorMessage = error instanceof ApiError
            ? error.message
            : 'Failed to fetch scorecard';
        yield put(fetchScorecardFailure(errorMessage));
    }
}

export function* startScoringSaga(action: ReturnType<typeof startScoringRequest>): Generator<CallEffect | PutEffect, void, ApiResponse<{ message: string; match_id: string }>> {
    try {
        const apiService = new ApiService();
        yield call(apiService.startScoring.bind(apiService), action.payload);
        yield put(startScoringSuccess());
        // Refresh scorecard after starting scoring
        yield put(fetchScorecardRequest(action.payload));
    } catch (error) {
        const errorMessage = error instanceof ApiError
            ? error.message
            : 'Failed to start scoring';
        yield put(startScoringFailure(errorMessage));
    }
}

export function* addBallSaga(action: ReturnType<typeof addBallRequest>): Generator<CallEffect | PutEffect, void, ApiResponse<{ message: string; match_id: string; innings_number: number; ball_type: string; run_type: string; runs: number; byes: number; is_wicket: boolean }>> {
    try {
        const apiService = new ApiService();
        const ballEvent = action.payload;

        // Use the addBall method directly with the complete ball event
        yield call(apiService.addBall.bind(apiService), ballEvent);
        yield put(addBallSuccess());
        // Refresh scorecard after adding ball
        yield put(fetchScorecardRequest(ballEvent.match_id));
    } catch (error) {
        const errorMessage = error instanceof ApiError
            ? error.message
            : 'Failed to add ball';
        yield put(addBallFailure(errorMessage));
    }
}

export function* fetchInningsSaga(action: ReturnType<typeof fetchInningsRequest>): Generator<CallEffect | PutEffect, void, ApiResponse<InningsSummary>> {
    try {
        const apiService = new ApiService();
        const { matchId, inningsNumber } = action.payload;
        const response = yield call(apiService.getInnings.bind(apiService), matchId, inningsNumber);
        yield put(fetchInningsSuccess(response.data));
    } catch (error) {
        const errorMessage = error instanceof ApiError
            ? error.message
            : 'Failed to fetch innings';
        yield put(fetchInningsFailure(errorMessage));
    }
}

export function* scorecardSaga() {
    yield takeLatest(fetchScorecardRequest.type, fetchScorecardSaga);
    yield takeEvery(startScoringRequest.type, startScoringSaga);
    yield takeEvery(addBallRequest.type, addBallSaga);
    yield takeLatest(fetchInningsRequest.type, fetchInningsSaga);
}
