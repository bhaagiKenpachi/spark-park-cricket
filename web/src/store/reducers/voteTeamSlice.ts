import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { VoteTeamWithPlayers, CreateVoteTeamRequest, UpdateVoteTeamRequest, AddPlayerRequest, TeamAssignmentRequest } from '@/types/team';

interface TeamState {
  teams: VoteTeamWithPlayers[];
  currentTeam: VoteTeamWithPlayers | null;
  loading: boolean;
  error: string | null;
}

const initialState: TeamState = {
  teams: [],
  currentTeam: null,
  loading: false,
  error: null,
};

const teamSlice = createSlice({
  name: 'team',
  initialState,
  reducers: {
    // Fetch teams for a vote
    fetchTeamsRequest: (state, action: PayloadAction<{ voteId: string }>) => {
      state.loading = true;
      state.error = null;
    },
    fetchTeamsSuccess: (state, action: PayloadAction<VoteTeamWithPlayers[]>) => {
      state.loading = false;
      state.teams = action.payload;
      state.error = null;
    },
    fetchTeamsFailure: (state, action: PayloadAction<string>) => {
      state.loading = false;
      state.error = action.payload;
    },

    // Fetch single team
    fetchTeamRequest: (state, action: PayloadAction<{ teamId: string }>) => {
      state.loading = true;
      state.error = null;
    },
    fetchTeamSuccess: (state, action: PayloadAction<VoteTeamWithPlayers>) => {
      state.loading = false;
      state.currentTeam = action.payload;
      state.error = null;
    },
    fetchTeamFailure: (state, action: PayloadAction<string>) => {
      state.loading = false;
      state.error = action.payload;
    },

    // Create team
    createTeamRequest: (state, action: PayloadAction<{ voteId: string; teamData: CreateVoteTeamRequest }>) => {
      state.loading = true;
      state.error = null;
    },
    createTeamSuccess: (state) => {
      state.loading = false;
      state.error = null;
      // Will refetch teams to get updated list
    },
    createTeamFailure: (state, action: PayloadAction<string>) => {
      state.loading = false;
      state.error = action.payload;
    },

    // Update team
    updateTeamRequest: (state, action: PayloadAction<{ teamId: string; teamData: UpdateVoteTeamRequest }>) => {
      state.loading = true;
      state.error = null;
    },
    updateTeamSuccess: (state) => {
      state.loading = false;
      state.error = null;
      // Will refetch team to get updated data
    },
    updateTeamFailure: (state, action: PayloadAction<string>) => {
      state.loading = false;
      state.error = action.payload;
    },

    // Delete team
    deleteTeamRequest: (state, action: PayloadAction<{ teamId: string }>) => {
      state.loading = true;
      state.error = null;
    },
    deleteTeamSuccess: (state, action: PayloadAction<{ teamId: string }>) => {
      state.loading = false;
      state.teams = state.teams.filter(t => t.id !== action.payload.teamId);
      if (state.currentTeam?.id === action.payload.teamId) {
        state.currentTeam = null;
      }
      state.error = null;
    },
    deleteTeamFailure: (state, action: PayloadAction<string>) => {
      state.loading = false;
      state.error = action.payload;
    },

    // Add player to team
    addPlayerRequest: (state, action: PayloadAction<{ teamId: string; voteId: string; playerData: AddPlayerRequest }>) => {
      state.loading = true;
      state.error = null;
    },
    addPlayerSuccess: (state) => {
      state.loading = false;
      state.error = null;
      // Will refetch teams to get updated player list
    },
    addPlayerFailure: (state, action: PayloadAction<string>) => {
      state.loading = false;
      state.error = action.payload;
    },

    // Add multiple players to team
    addPlayersRequest: (state, action: PayloadAction<{ teamId: string; voteId: string; playerData: TeamAssignmentRequest }>) => {
      state.loading = true;
      state.error = null;
    },
    addPlayersSuccess: (state) => {
      state.loading = false;
      state.error = null;
      // Will refetch teams to get updated player list
    },
    addPlayersFailure: (state, action: PayloadAction<string>) => {
      state.loading = false;
      state.error = action.payload;
    },

    // Remove player from team
    removePlayerRequest: (state, action: PayloadAction<{ teamId: string; voteId: string; playerId: string }>) => {
      state.loading = true;
      state.error = null;
    },
    removePlayerSuccess: (state) => {
      state.loading = false;
      state.error = null;
      // Will refetch teams to get updated player list
    },
    removePlayerFailure: (state, action: PayloadAction<string>) => {
      state.loading = false;
      state.error = action.payload;
    },

    // Clear error
    clearTeamError: (state) => {
      state.error = null;
    },

    // Clear teams
    clearTeams: (state) => {
      state.teams = [];
      state.currentTeam = null;
    },
  },
});

export const {
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
  clearTeamError,
  clearTeams,
} = teamSlice.actions;

export default teamSlice.reducer;
