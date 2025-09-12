import { combineReducers } from '@reduxjs/toolkit';
import { seriesSlice } from './seriesSlice';
import { matchSlice } from './matchSlice';
import { teamSlice } from './teamSlice';
import { playerSlice } from './playerSlice';
import { scoreboardSlice } from './scoreboardSlice';
import scorecardReducer from './scorecardSlice';

export const rootReducer = combineReducers({
    series: seriesSlice.reducer,
    match: matchSlice.reducer,
    team: teamSlice.reducer,
    player: playerSlice.reducer,
    scoreboard: scoreboardSlice.reducer,
    scorecard: scorecardReducer,
});

export type RootState = ReturnType<typeof rootReducer>;
