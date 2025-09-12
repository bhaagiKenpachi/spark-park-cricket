import { call, put, takeEvery } from 'redux-saga/effects';
import {
    fetchScoreboardRequest,
    fetchScoreboardSuccess,
    fetchScoreboardFailure,
} from '../reducers/scoreboardSlice';
import { Scoreboard } from '../reducers/scoreboardSlice';

// API functions (to be implemented)
const api = {
    fetchScoreboard: async (matchId: string): Promise<Scoreboard> => {
        const response = await fetch(`/api/v1/scoreboard/${matchId}`);
        if (!response.ok) {
            throw new Error('Failed to fetch scoreboard');
        }
        return response.json();
    },
};

function* fetchScoreboardSaga(action: { type: string; payload: string }) {
    try {
        const scoreboard: Scoreboard = yield call(api.fetchScoreboard, action.payload);
        yield put(fetchScoreboardSuccess(scoreboard));
    } catch (error) {
        yield put(fetchScoreboardFailure(error instanceof Error ? error.message : 'Unknown error'));
    }
}

export function* scoreboardSaga() {
    yield takeEvery(fetchScoreboardRequest.type, fetchScoreboardSaga);
}
