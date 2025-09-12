import { call, put, takeEvery, takeLatest, CallEffect, PutEffect, delay } from 'redux-saga/effects';
import {
    fetchSeriesRequest,
    fetchSeriesSuccess,
    fetchSeriesFailure,
    createSeriesRequest,
    createSeriesSuccess,
    createSeriesFailure,
    updateSeriesRequest,
    updateSeriesSuccess,
    updateSeriesFailure,
    deleteSeriesRequest,
    deleteSeriesSuccess,
    deleteSeriesFailure,
} from '../reducers/seriesSlice';
import { Series } from '../reducers/seriesSlice';
import { ApiService, ApiError, ApiResponse } from '@/services/api';

export function* fetchSeriesSaga(): Generator<CallEffect | PutEffect, void, ApiResponse<Series[]>> {
    try {
        const apiService = new ApiService();
        const response = yield call(apiService.getSeries.bind(apiService));

        // Extract the actual array from the response
        const seriesArray = response.data;

        yield put(fetchSeriesSuccess(seriesArray));
    } catch (error) {
        const errorMessage = error instanceof ApiError
            ? error.message
            : 'Failed to fetch series';

        // Add a small delay before showing error to allow for retries
        yield delay(100);
        yield put(fetchSeriesFailure(errorMessage));
    }
}

export function* createSeriesSaga(action: ReturnType<typeof createSeriesRequest>): Generator<CallEffect | PutEffect, void, ApiResponse<Series>> {
    try {
        const apiService = new ApiService();
        const response = yield call(apiService.createSeries.bind(apiService), action.payload);
        yield put(createSeriesSuccess(response.data));
    } catch (error) {
        const errorMessage = error instanceof ApiError
            ? error.message
            : 'Failed to create series';
        yield put(createSeriesFailure(errorMessage));
    }
}

export function* updateSeriesSaga(action: ReturnType<typeof updateSeriesRequest>): Generator<CallEffect | PutEffect, void, ApiResponse<Series>> {
    try {
        const apiService = new ApiService();
        const { id, seriesData } = action.payload;
        const response = yield call(apiService.updateSeries.bind(apiService), id, seriesData);
        yield put(updateSeriesSuccess(response.data));
    } catch (error) {
        const errorMessage = error instanceof ApiError
            ? error.message
            : 'Failed to update series';
        yield put(updateSeriesFailure(errorMessage));
    }
}

export function* deleteSeriesSaga(action: ReturnType<typeof deleteSeriesRequest>): Generator<CallEffect | PutEffect, void, ApiResponse<void>> {
    try {
        const apiService = new ApiService();
        yield call(apiService.deleteSeries.bind(apiService), action.payload);
        yield put(deleteSeriesSuccess(action.payload));
    } catch (error) {
        const errorMessage = error instanceof ApiError
            ? error.message
            : 'Failed to delete series';
        yield put(deleteSeriesFailure(errorMessage));
    }
}

export function* seriesSaga() {
    yield takeLatest(fetchSeriesRequest.type, fetchSeriesSaga);
    yield takeEvery(createSeriesRequest.type, createSeriesSaga);
    yield takeEvery(updateSeriesRequest.type, updateSeriesSaga);
    yield takeEvery(deleteSeriesRequest.type, deleteSeriesSaga);
}
