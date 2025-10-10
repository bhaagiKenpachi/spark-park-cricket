import { all, fork } from 'redux-saga/effects';
import { seriesSaga } from './seriesSaga';
import { matchSaga } from './matchSaga';
import { playerSaga } from './playerSaga';
import { scoreboardSaga } from './scoreboardSaga';
import { scorecardSaga } from './scorecardSaga';
import { voteSaga } from './voteSaga';
import voteTeamSaga from './voteTeamSaga';

export function* rootSaga() {
  yield all([
    fork(seriesSaga),
    fork(matchSaga),
    fork(playerSaga),
    fork(scoreboardSaga),
    fork(scorecardSaga),
    fork(voteSaga),
    fork(voteTeamSaga),
  ]);
}
