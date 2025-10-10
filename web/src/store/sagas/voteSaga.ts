import {
    call,
    put,
    takeEvery,
    takeLatest,
    CallEffect,
    PutEffect,
    delay,
} from 'redux-saga/effects';
import {
    fetchVotesRequest,
    fetchVotesSuccess,
    fetchVotesFailure,
    fetchVoteRequest,
    fetchVoteSuccess,
    fetchVoteFailure,
    fetchVoteWithResultsRequest,
    fetchVoteWithResultsSuccess,
    fetchVoteWithResultsFailure,
    createVoteRequest,
    createVoteSuccess,
    createVoteFailure,
    updateVoteRequest,
    updateVoteSuccess,
    updateVoteFailure,
    deleteVoteRequest,
    deleteVoteSuccess,
    deleteVoteFailure,
    castVoteRequest,
    castVoteSuccess,
    castVoteFailure,
    fetchUserVoteRequest,
    fetchUserVoteSuccess,
    fetchUserVoteFailure,
    checkUserVotedRequest,
    checkUserVotedSuccess,
    checkUserVotedFailure,
    closeVoteRequest,
    closeVoteSuccess,
    closeVoteFailure,
    cancelVoteRequest,
    cancelVoteSuccess,
    cancelVoteFailure,
} from '../reducers/voteSlice';
import { VoteFilters } from '@/types/vote';
import { ApiService } from '@/services/api';

// Fetch votes saga
export function* fetchVotesSaga(
    action: ReturnType<typeof fetchVotesRequest>
): Generator<CallEffect | PutEffect, void, any> {
    try {
        const apiService = new ApiService();
        const filters = action.payload;
        const response = yield call(apiService.getVotes.bind(apiService), filters);

        // Extract paginated data from response - handle nested data structure
        let paginatedData = response.data;

        // Check if data is nested (response.data.data)
        if (paginatedData && typeof paginatedData === 'object' && 'data' in paginatedData) {
            paginatedData = paginatedData.data;
        }

        // Check if response is paginated format
        if (paginatedData && typeof paginatedData === 'object' && 'votes' in paginatedData) {
            // New paginated format
            yield put(
                fetchVotesSuccess({
                    votes: paginatedData.votes || [],
                    totalItems: paginatedData.total_items || 0,
                    page: paginatedData.page || 1,
                    pageSize: paginatedData.page_size || 20,
                    totalPages: paginatedData.total_pages || 1,
                })
            );
        } else {
            // Legacy array format (fallback)
            const votesData = Array.isArray(paginatedData) ? paginatedData : [];
            yield put(
                fetchVotesSuccess({
                    votes: votesData,
                    totalItems: votesData.length,
                })
            );
        }
    } catch (error) {
        const errorMessage = error instanceof Error ? error.message : 'Failed to fetch votes';
        yield delay(100);
        yield put(fetchVotesFailure(errorMessage));
    }
}

// Fetch single vote saga
export function* fetchVoteSaga(
    action: ReturnType<typeof fetchVoteRequest>
): Generator<CallEffect | PutEffect, void, any> {
    try {
        const apiService = new ApiService();
        const voteId = action.payload;
        const response = yield call(apiService.getVoteById.bind(apiService), voteId);

        const voteData = response.data.data || response.data;
        yield put(fetchVoteSuccess(voteData));
    } catch (error) {
        const errorMessage = error instanceof Error ? error.message : 'Failed to fetch vote';
        yield put(fetchVoteFailure(errorMessage));
    }
}

// Fetch vote with results saga
export function* fetchVoteWithResultsSaga(
    action: ReturnType<typeof fetchVoteWithResultsRequest>
): Generator<CallEffect | PutEffect, void, any> {
    try {
        const apiService = new ApiService();
        const voteId = action.payload;
        const response = yield call(apiService.getVoteWithResults.bind(apiService), voteId);

        const voteWithResultsData = response.data.data || response.data;
        yield put(fetchVoteWithResultsSuccess(voteWithResultsData));
    } catch (error) {
        const errorMessage = error instanceof Error ? error.message : 'Failed to fetch vote with results';
        yield put(fetchVoteWithResultsFailure(errorMessage));
    }
}

// Create vote saga
export function* createVoteSaga(
    action: ReturnType<typeof createVoteRequest>
): Generator<CallEffect | PutEffect, void, any> {
    try {
        const apiService = new ApiService();
        const { title, description, type, options } = action.payload;
        const response = yield call(
            apiService.createVote.bind(apiService),
            { title, description, type, options }
        );

        const voteData = response.data.data || response.data;
        yield put(createVoteSuccess(voteData));

        // Refetch votes list after successful creation
        yield put(fetchVotesRequest(undefined));
    } catch (error) {
        const errorMessage = error instanceof Error ? error.message : 'Failed to create vote';
        yield put(createVoteFailure(errorMessage));
    }
}

// Update vote saga
export function* updateVoteSaga(
    action: ReturnType<typeof updateVoteRequest>
): Generator<CallEffect | PutEffect, void, any> {
    try {
        const apiService = new ApiService();
        const { id, voteData } = action.payload;
        const response = yield call(
            apiService.updateVote.bind(apiService),
            id,
            voteData
        );

        const updatedVote = response.data.data || response.data;
        yield put(updateVoteSuccess(updatedVote));

        // Refetch votes list after successful update
        yield put(fetchVotesRequest(undefined));
    } catch (error) {
        const errorMessage = error instanceof Error ? error.message : 'Failed to update vote';
        yield put(updateVoteFailure(errorMessage));
    }
}

// Delete vote saga
export function* deleteVoteSaga(
    action: ReturnType<typeof deleteVoteRequest>
): Generator<CallEffect | PutEffect, void, any> {
    try {
        const apiService = new ApiService();
        const voteId = action.payload;
        yield call(apiService.deleteVote.bind(apiService), voteId);

        yield put(deleteVoteSuccess(voteId));

        // Refetch votes list after successful deletion
        yield put(fetchVotesRequest(undefined));
    } catch (error) {
        const errorMessage = error instanceof Error ? error.message : 'Failed to delete vote';
        yield put(deleteVoteFailure(errorMessage));
    }
}

// Cast vote saga
export function* castVoteSaga(
    action: ReturnType<typeof castVoteRequest>
): Generator<CallEffect | PutEffect, void, any> {
    try {
        const apiService = new ApiService();
        const { voteId, optionIds } = action.payload;
        const response = yield call(
            apiService.castVote.bind(apiService),
            voteId,
            { selected_options: optionIds }
        );

        const userVoteData = response.data.data || response.data;
        yield put(castVoteSuccess(userVoteData));
    } catch (error) {
        const errorMessage = error instanceof Error ? error.message : 'Failed to cast vote';
        yield put(castVoteFailure(errorMessage));
    }
}

// Get user vote saga
export function* fetchUserVoteSaga(
    action: ReturnType<typeof fetchUserVoteRequest>
): Generator<CallEffect | PutEffect, void, any> {
    try {
        const apiService = new ApiService();
        const voteId = action.payload;
        const response = yield call(apiService.getUserVote.bind(apiService), voteId);

        const userVoteData = response.data.data || response.data;
        yield put(fetchUserVoteSuccess(userVoteData));
    } catch (error) {
        const errorMessage = error instanceof Error ? error.message : 'Failed to fetch user vote';
        yield put(fetchUserVoteFailure(errorMessage));
    }
}

// Check if user voted saga
export function* checkUserVotedSaga(
    action: ReturnType<typeof checkUserVotedRequest>
): Generator<CallEffect | PutEffect, void, any> {
    try {
        const apiService = new ApiService();
        const voteId = action.payload;
        const response = yield call(apiService.hasUserVoted.bind(apiService), voteId);

        const hasVotedData = response.data.data || response.data;
        yield put(
            checkUserVotedSuccess({
                voteId,
                hasVoted: hasVotedData.has_voted,
            })
        );
    } catch (error) {
        const errorMessage = error instanceof Error ? error.message : 'Failed to check user vote status';
        yield put(checkUserVotedFailure(errorMessage));
    }
}

// Close vote saga
export function* closeVoteSaga(
    action: ReturnType<typeof closeVoteRequest>
): Generator<CallEffect | PutEffect, void, any> {
    try {
        const apiService = new ApiService();
        const voteId = action.payload;
        const response = yield call(apiService.closeVote.bind(apiService), voteId);

        const voteData = response.data.data || response.data;
        yield put(closeVoteSuccess(voteData));

        // Refetch votes list after successful close
        yield put(fetchVotesRequest(undefined));
    } catch (error) {
        const errorMessage = error instanceof Error ? error.message : 'Failed to close vote';
        yield put(closeVoteFailure(errorMessage));
    }
}

// Cancel vote saga
export function* cancelVoteSaga(
    action: ReturnType<typeof cancelVoteRequest>
): Generator<CallEffect | PutEffect, void, any> {
    try {
        const apiService = new ApiService();
        const voteId = action.payload;
        const response = yield call(apiService.cancelVote.bind(apiService), voteId);

        const voteData = response.data.data || response.data;
        yield put(cancelVoteSuccess(voteData));

        // Refetch votes list after successful cancellation
        yield put(fetchVotesRequest(undefined));
    } catch (error) {
        const errorMessage = error instanceof Error ? error.message : 'Failed to cancel vote';
        yield put(cancelVoteFailure(errorMessage));
    }
}

// Main vote saga
export function* voteSaga() {
    yield takeLatest(fetchVotesRequest.type, fetchVotesSaga);
    yield takeEvery(fetchVoteRequest.type, fetchVoteSaga);
    yield takeEvery(fetchVoteWithResultsRequest.type, fetchVoteWithResultsSaga);
    yield takeEvery(createVoteRequest.type, createVoteSaga);
    yield takeEvery(updateVoteRequest.type, updateVoteSaga);
    yield takeEvery(deleteVoteRequest.type, deleteVoteSaga);
    yield takeEvery(castVoteRequest.type, castVoteSaga);
    yield takeEvery(fetchUserVoteRequest.type, fetchUserVoteSaga);
    yield takeEvery(checkUserVotedRequest.type, checkUserVotedSaga);
    yield takeEvery(closeVoteRequest.type, closeVoteSaga);
    yield takeEvery(cancelVoteRequest.type, cancelVoteSaga);
}
