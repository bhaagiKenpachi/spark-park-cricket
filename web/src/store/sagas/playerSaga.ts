import { call, put, takeEvery } from 'redux-saga/effects';
import {
  fetchPlayersRequest,
  fetchPlayersSuccess,
  fetchPlayersFailure,
} from '../reducers/playerSlice';
import { Player } from '../reducers/playerSlice';

// API functions (to be implemented)
const api = {
  fetchPlayers: async (): Promise<Player[]> => {
    const response = await fetch('/api/v1/players');
    if (!response.ok) {
      throw new Error('Failed to fetch players');
    }
    return response.json();
  },
};

function* fetchPlayersSaga() {
  try {
    const players: Player[] = yield call(api.fetchPlayers);
    yield put(fetchPlayersSuccess(players));
  } catch (error) {
    yield put(
      fetchPlayersFailure(
        error instanceof Error ? error.message : 'Unknown error'
      )
    );
  }
}

export function* playerSaga() {
  yield takeEvery(fetchPlayersRequest.type, fetchPlayersSaga);
}
