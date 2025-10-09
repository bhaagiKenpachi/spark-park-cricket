import { call, put, takeLatest } from 'redux-saga/effects';
import { PayloadAction } from '@reduxjs/toolkit';
import { apiService, ApiResponse } from '@/services/api';
import { VoteTeamWithPlayers, VoteTeam, CreateVoteTeamRequest, UpdateVoteTeamRequest } from '@/types/team';
import {
  fetchTeamsRequest,
  fetchTeamsSuccess,
  fetchTeamsFailure,
  fetchTeamRequest,
  fetchTeamSuccess,
  fetchTeamFailure,
  createTeamRequest,
  createTeamSuccess,
  createTeamFailure,
  updateTeamRequest,
  updateTeamSuccess,
  updateTeamFailure,
  deleteTeamRequest,
  deleteTeamSuccess,
  deleteTeamFailure,
  addPlayerRequest,
  addPlayerSuccess,
  addPlayerFailure,
  addPlayersRequest,
  addPlayersSuccess,
  addPlayersFailure,
  removePlayerRequest,
  removePlayerSuccess,
  removePlayerFailure,
} from '../reducers/voteTeamSlice';

// Fetch teams for a vote
function* fetchTeamsSaga(action: PayloadAction<{ voteId: string }>) {
  try {
    const response: ApiResponse<VoteTeamWithPlayers[]> = yield call([apiService, apiService.getTeamsByVoteId], action.payload.voteId);

    // The API returns { data: { data: [...] } }, so we need to extract the nested data
    const teams = (response.data as any)?.data || response.data || [];

    yield put(fetchTeamsSuccess(Array.isArray(teams) ? teams : []));
  } catch (error: unknown) {
    const message = error instanceof Error ? error.message : 'Failed to fetch teams';
    yield put(fetchTeamsFailure(message));
  }
}

// Fetch single team
function* fetchTeamSaga(action: PayloadAction<{ teamId: string }>) {
  try {
    const response: ApiResponse<VoteTeamWithPlayers> = yield call([apiService, apiService.getTeamById], action.payload.teamId);
    // Extract nested data
    const team = (response.data as any)?.data || response.data;
    yield put(fetchTeamSuccess(team));
  } catch (error: unknown) {
    const message = error instanceof Error ? error.message : 'Failed to fetch team';
    yield put(fetchTeamFailure(message));
  }
}

// Create team
function* createTeamSaga(action: PayloadAction<{ voteId: string; teamData: CreateVoteTeamRequest }>) {
  try {
    yield call([apiService, apiService.createTeam], action.payload.voteId, action.payload.teamData);
    yield put(createTeamSuccess());
    // Refetch all teams for the vote
    yield put(fetchTeamsRequest({ voteId: action.payload.voteId }));
  } catch (error: unknown) {
    const message = error instanceof Error ? error.message : 'Failed to create team';
    yield put(createTeamFailure(message));
  }
}

// Update team
function* updateTeamSaga(action: PayloadAction<{ teamId: string; teamData: UpdateVoteTeamRequest }>) {
  try {
    yield call([apiService, apiService.updateTeam], action.payload.teamId, action.payload.teamData);
    yield put(updateTeamSuccess());
    // Refetch team to get updated data
    yield put(fetchTeamRequest({ teamId: action.payload.teamId }));
  } catch (error: unknown) {
    const message = error instanceof Error ? error.message : 'Failed to update team';
    yield put(updateTeamFailure(message));
  }
}

// Delete team
function* deleteTeamSaga(action: PayloadAction<{ teamId: string }>) {
  try {
    yield call([apiService, apiService.deleteTeam], action.payload.teamId);
    yield put(deleteTeamSuccess({ teamId: action.payload.teamId }));
  } catch (error: unknown) {
    const message = error instanceof Error ? error.message : 'Failed to delete team';
    yield put(deleteTeamFailure(message));
  }
}

// Add player to team
function* addPlayerSaga(action: PayloadAction<{ teamId: string; voteId: string; playerData: { user_id: string } }>) {
  try {
    yield call([apiService, apiService.addPlayerToTeam], action.payload.teamId, action.payload.playerData.user_id);
    yield put(addPlayerSuccess());
    // Refetch all teams for the vote to update both team cards
    yield put(fetchTeamsRequest({ voteId: action.payload.voteId }));
  } catch (error: unknown) {
    const message = error instanceof Error ? error.message : 'Failed to add player';
    yield put(addPlayerFailure(message));
  }
}

// Add multiple players to team
function* addPlayersSaga(action: PayloadAction<{ teamId: string; voteId: string; playerData: { user_ids: string[] } }>) {
  try {
    yield call([apiService, apiService.addPlayersToTeam], action.payload.teamId, action.payload.playerData.user_ids);
    yield put(addPlayersSuccess());
    // Refetch all teams for the vote to update both team cards
    yield put(fetchTeamsRequest({ voteId: action.payload.voteId }));
  } catch (error: unknown) {
    const message = error instanceof Error ? error.message : 'Failed to add players';
    yield put(addPlayersFailure(message));
  }
}

// Remove player from team
function* removePlayerSaga(action: PayloadAction<{ teamId: string; voteId: string; playerId: string }>) {
  try {
    yield call([apiService, apiService.removePlayerFromTeam], action.payload.teamId, action.payload.playerId);
    yield put(removePlayerSuccess());
    // Refetch all teams for the vote to update both team cards
    yield put(fetchTeamsRequest({ voteId: action.payload.voteId }));
  } catch (error: unknown) {
    const message = error instanceof Error ? error.message : 'Failed to remove player';
    yield put(removePlayerFailure(message));
  }
}

// Root saga
export default function* voteTeamSaga() {
  yield takeLatest(fetchTeamsRequest.type, fetchTeamsSaga);
  yield takeLatest(fetchTeamRequest.type, fetchTeamSaga);
  yield takeLatest(createTeamRequest.type, createTeamSaga);
  yield takeLatest(updateTeamRequest.type, updateTeamSaga);
  yield takeLatest(deleteTeamRequest.type, deleteTeamSaga);
  yield takeLatest(addPlayerRequest.type, addPlayerSaga);
  yield takeLatest(addPlayersRequest.type, addPlayersSaga);
  yield takeLatest(removePlayerRequest.type, removePlayerSaga);
}
