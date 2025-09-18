import {
  call,
  put,
  takeEvery,
  takeLatest,
  CallEffect,
  PutEffect,
} from 'redux-saga/effects';
import {
  fetchMatchesRequest,
  fetchMatchesSuccess,
  fetchMatchesFailure,
  createMatchRequest,
  createMatchSuccess,
  createMatchFailure,
  updateMatchRequest,
  updateMatchSuccess,
  updateMatchFailure,
  deleteMatchRequest,
  deleteMatchSuccess,
  deleteMatchFailure,
} from '../reducers/matchSlice';
import { ApiService, ApiError, ApiResponse } from '@/services/api';
import { Match } from '../reducers/matchSlice';

export function* fetchMatchesSaga(): Generator<
  CallEffect | PutEffect,
  void,
  ApiResponse<{ data: Match[] }>
> {
  try {
    const apiService = new ApiService();
    const response = yield call(apiService.getMatches.bind(apiService));

    // Extract the actual array from the nested response structure
    const matchesArray = response.data.data || response.data;
    yield put(fetchMatchesSuccess(matchesArray));
  } catch (error) {
    const errorMessage =
      error instanceof ApiError ? error.message : 'Failed to fetch matches';
    yield put(fetchMatchesFailure(errorMessage));
  }
}

export function* createMatchSaga(
  action: ReturnType<typeof createMatchRequest>
): Generator<CallEffect | PutEffect, void, ApiResponse<{ data: Match }>> {
  try {
    const apiService = new ApiService();
    const response = yield call(
      apiService.createMatch.bind(apiService),
      action.payload
    );
    yield put(createMatchSuccess(response.data.data || response.data));
  } catch (error) {
    const errorMessage =
      error instanceof ApiError ? error.message : 'Failed to create match';
    yield put(createMatchFailure(errorMessage));
  }
}

export function* updateMatchSaga(
  action: ReturnType<typeof updateMatchRequest>
): Generator<CallEffect | PutEffect, void, ApiResponse<{ data: Match }>> {
  try {
    const apiService = new ApiService();
    const { id, matchData } = action.payload;
    const response = yield call(
      apiService.updateMatch.bind(apiService),
      id,
      matchData
    );
    yield put(updateMatchSuccess(response.data.data || response.data));
  } catch (error) {
    const errorMessage =
      error instanceof ApiError ? error.message : 'Failed to update match';
    yield put(updateMatchFailure(errorMessage));
  }
}

export function* deleteMatchSaga(
  action: ReturnType<typeof deleteMatchRequest>
): Generator<CallEffect | PutEffect, void, ApiResponse<{ data: void }>> {
  try {
    const apiService = new ApiService();
    yield call(apiService.deleteMatch.bind(apiService), action.payload);
    yield put(deleteMatchSuccess(action.payload));
  } catch (error) {
    const errorMessage =
      error instanceof ApiError ? error.message : 'Failed to delete match';
    yield put(deleteMatchFailure(errorMessage));
  }
}

export function* matchSaga() {
  yield takeLatest(fetchMatchesRequest.type, fetchMatchesSaga);
  yield takeEvery(createMatchRequest.type, createMatchSaga);
  yield takeEvery(updateMatchRequest.type, updateMatchSaga);
  yield takeEvery(deleteMatchRequest.type, deleteMatchSaga);
}
