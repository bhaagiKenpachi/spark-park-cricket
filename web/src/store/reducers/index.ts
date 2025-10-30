import { combineReducers } from '@reduxjs/toolkit';
import { seriesSlice } from './seriesSlice';
import { matchSlice } from './matchSlice';
import { playerSlice } from './playerSlice';
import { scoreboardSlice } from './scoreboardSlice';
import scorecardReducer from './scorecardSlice';
import authReducer from './authSlice';
import voteReducer from './voteSlice';
import voteTeamReducer from './voteTeamSlice';
import eventReducer from './eventSlice';
import timeTrackingReducer from './timeTrackingSlice';
import groupReducer from './groupSlice';

export const rootReducer = combineReducers({
  series: seriesSlice.reducer,
  match: matchSlice.reducer,
  player: playerSlice.reducer,
  scoreboard: scoreboardSlice.reducer,
  scorecard: scorecardReducer,
  auth: authReducer,
  vote: voteReducer,
  voteTeam: voteTeamReducer,
  events: eventReducer,
  timeTracking: timeTrackingReducer,
  group: groupReducer,
});

export type RootState = ReturnType<typeof rootReducer>;
