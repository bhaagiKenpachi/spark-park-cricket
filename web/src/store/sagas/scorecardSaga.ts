import { call, put, takeEvery, takeLatest } from 'redux-saga/effects';
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
} from '../reducers/scorecardSlice';
import { ApiService, ApiError } from '@/services/api';

export function* fetchScorecardSaga(action: ReturnType<typeof fetchScorecardRequest>): Generator<any, void, any> {
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

export function* startScoringSaga(action: ReturnType<typeof startScoringRequest>): Generator<any, void, any> {
    try {
        const apiService = new ApiService();
        const response = yield call(apiService.startScoring.bind(apiService), action.payload);
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

export function* addBallSaga(action: ReturnType<typeof addBallRequest>): Generator<any, void, any> {
    try {
        const apiService = new ApiService();
        const ballEvent = action.payload;

        // Use the addBall method directly with the complete ball event
        const response = yield call(apiService.addBall.bind(apiService), ballEvent);
        yield put(addBallSuccess(response.data));
        // Refresh scorecard after adding ball
        yield put(fetchScorecardRequest(ballEvent.match_id));
    } catch (error) {
        const errorMessage = error instanceof ApiError
            ? error.message
            : 'Failed to add ball';
        yield put(addBallFailure(errorMessage));
    }
}

export function* fetchInningsSaga(action: ReturnType<typeof fetchInningsRequest>): Generator<any, void, any> {
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
