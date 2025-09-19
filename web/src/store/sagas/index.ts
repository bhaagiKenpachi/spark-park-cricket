import { all, fork } from 'redux-saga/effects';
import { seriesSaga } from './seriesSaga';
import { matchSaga } from './matchSaga';
import { teamSaga } from './teamSaga';
import { playerSaga } from './playerSaga';
import { scoreboardSaga } from './scoreboardSaga';
import { scorecardSaga } from './scorecardSaga';

export function* rootSaga() {
  yield all([
    fork(seriesSaga),
    fork(matchSaga),
    fork(teamSaga),
    fork(playerSaga),
    fork(scoreboardSaga),
    fork(scorecardSaga),
  ]);
}
