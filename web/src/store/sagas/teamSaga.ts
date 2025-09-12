import { call, put, takeEvery } from 'redux-saga/effects';
import {
    fetchTeamsRequest,
    fetchTeamsSuccess,
    fetchTeamsFailure,
} from '../reducers/teamSlice';
import { Team } from '../reducers/teamSlice';

// API functions (to be implemented)
const api = {
    fetchTeams: async (): Promise<Team[]> => {
        const response = await fetch('/api/v1/teams');
        if (!response.ok) {
            throw new Error('Failed to fetch teams');
        }
        return response.json();
    },
};

function* fetchTeamsSaga() {
    try {
        const teams: Team[] = yield call(api.fetchTeams);
        yield put(fetchTeamsSuccess(teams));
    } catch (error) {
        yield put(fetchTeamsFailure(error instanceof Error ? error.message : 'Unknown error'));
    }
}

export function* teamSaga() {
    yield takeEvery(fetchTeamsRequest.type, fetchTeamsSaga);
}
