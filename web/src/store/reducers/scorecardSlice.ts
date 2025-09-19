/* eslint-disable @typescript-eslint/no-unused-vars */
import { createSlice, PayloadAction, createAsyncThunk } from '@reduxjs/toolkit';
import { graphqlService } from '@/services/graphqlService';
import { apiService } from '@/services/api';
import { BallSummary as GraphQLBallSummary, OverSummary as GraphQLOverSummary } from '@/lib/graphql';

// Enums
export type BallType = 'good' | 'wide' | 'no_ball' | 'dead_ball';
export type RunType =
  | '0'
  | '1'
  | '2'
  | '3'
  | '4'
  | '5'
  | '6'
  | '7'
  | '8'
  | '9'
  | 'NB'
  | 'WD'
  | 'LB'
  | 'WC';
export type TeamType = 'A' | 'B';
export type TossType = 'H' | 'T';
export type WicketType =
  | 'bowled'
  | 'caught'
  | 'lbw'
  | 'run_out'
  | 'stumped'
  | 'hit_wicket';

// Ball Summary
export interface BallSummary {
  ball_number: number;
  ball_type: BallType;
  run_type: RunType;
  runs: number;
  byes: number;
  is_wicket: boolean;
  wicket_type?: string;
}

// Over Summary
export interface OverSummary {
  over_number: number;
  total_runs: number;
  total_balls: number;
  total_wickets: number;
  status: string;
  balls: BallSummary[];
}

// Extras Summary
export interface ExtrasSummary {
  byes: number;
  leg_byes: number;
  wides: number;
  no_balls: number;
  total: number;
}

// Innings Summary
export interface InningsSummary {
  innings_number: number;
  batting_team: TeamType;
  total_runs: number;
  total_wickets: number;
  total_overs: number;
  total_balls: number;
  status: string;
  extras?: ExtrasSummary;
  overs: OverSummary[];
}

// Scorecard Response
export interface ScorecardResponse {
  match_id: string;
  match_number: number;
  series_name: string;
  team_a: string;
  team_b: string;
  total_overs: number;
  toss_winner: TeamType;
  toss_type: TossType;
  current_innings: number;
  innings: InningsSummary[];
  match_status: string;
}

// Ball Event Request
export interface BallEventRequest {
  match_id: string;
  innings_number: number;
  ball_type: BallType;
  run_type: RunType;
  runs: number;
  is_wicket: boolean;
  wicket_type?: string;
  byes?: number;
}

// Scorecard State
interface ScorecardState {
  scorecard: ScorecardResponse | null;
  loading: boolean;
  error: string | null;
  scoring: boolean;
  inningsTransitionDetected: boolean;
}

export const initialState: ScorecardState = {
  scorecard: null,
  loading: false,
  error: null,
  scoring: false,
  inningsTransitionDetected: false,
};

// Async thunks for GraphQL operations
export const fetchInningsScoreSummaryThunk = createAsyncThunk(
  'scorecard/fetchInningsScoreSummary',
  async ({ matchId, inningsNumber }: { matchId: string; inningsNumber: number }) => {
    const response = await graphqlService.getInningsScoreSummary(matchId, inningsNumber);
    if (response.success && response.data) {
      return response.data;
    }
    throw new Error(response.error || 'Failed to fetch innings score summary');
  }
);

export const fetchLatestOverThunk = createAsyncThunk(
  'scorecard/fetchLatestOver',
  async ({ matchId, inningsNumber }: { matchId: string; inningsNumber: number }) => {
    const response = await graphqlService.getLatestOverOnly(matchId, inningsNumber);
    if (response.success && response.data) {
      return { inningsNumber, over: response.data };
    }
    throw new Error(response.error || 'Failed to fetch latest over');
  }
);

export const fetchAllOversDetailsThunk = createAsyncThunk(
  'scorecard/fetchAllOversDetails',
  async ({ matchId, inningsNumber }: { matchId: string; inningsNumber: number }) => {
    const response = await graphqlService.getAllOversDetails(matchId, inningsNumber);
    if (response.success && response.data) {
      return { inningsNumber, overs: response.data };
    }
    throw new Error(response.error || 'Failed to fetch all overs details');
  }
);

export const undoBallThunk = createAsyncThunk(
  'scorecard/undoBall',
  async ({ matchId, inningsNumber }: { matchId: string; inningsNumber: number }) => {
    const response = await apiService.undoBall(matchId, inningsNumber);
    if (response.success) {
      return { matchId, inningsNumber, message: response.message };
    }
    throw new Error(response.message || 'Failed to undo ball');
  }
);

export const scorecardSlice = createSlice({
  name: 'scorecard',
  initialState,
  reducers: {
    fetchScorecardRequest: (state, _action: PayloadAction<string>) => {
      state.loading = true;
      state.error = null;
    },
    fetchScorecardSuccess: (
      state,
      action: PayloadAction<ScorecardResponse>
    ) => {
      state.loading = false;
      state.scorecard = action.payload;
    },
    fetchScorecardFailure: (state, action: PayloadAction<string>) => {
      state.loading = false;
      state.error = action.payload;
    },
    startScoringRequest: (state, _action: PayloadAction<string>) => {
      state.scoring = true;
      state.error = null;
    },
    startScoringSuccess: state => {
      state.scoring = false;
    },
    startScoringFailure: (state, action: PayloadAction<string>) => {
      state.scoring = false;
      state.error = action.payload;
    },
    addBallRequest: (state, _action: PayloadAction<BallEventRequest>) => {
      state.scoring = true;
      state.error = null;
    },
    addBallSuccess: state => {
      state.scoring = false;
    },
    addBallFailure: (state, action: PayloadAction<string>) => {
      state.scoring = false;
      state.error = action.payload;
    },
    undoBallRequest: (state, _action: PayloadAction<{ matchId: string; inningsNumber: number }>) => {
      state.scoring = true;
      state.error = null;
    },
    undoBallSuccess: state => {
      state.scoring = false;
    },
    undoBallFailure: (state, action: PayloadAction<string>) => {
      state.scoring = false;
      state.error = action.payload;
    },
    fetchInningsRequest: (
      state,
      _action: PayloadAction<{ matchId: string; inningsNumber: number }>
    ) => {
      state.loading = true;
      state.error = null;
    },
    fetchInningsSuccess: (state, action: PayloadAction<InningsSummary>) => {
      state.loading = false;
      // Update the specific innings in the scorecard
      if (state.scorecard) {
        const inningsIndex = state.scorecard.innings.findIndex(
          innings => innings.innings_number === action.payload.innings_number
        );
        if (inningsIndex !== -1) {
          state.scorecard.innings[inningsIndex] = action.payload;
        }
      }
    },
    fetchInningsScoreSummarySuccess: (
      state,
      action: PayloadAction<InningsSummary>
    ) => {
      state.loading = false;
      if (state.scorecard) {
        const inningsIndex = state.scorecard.innings.findIndex(
          innings => innings.innings_number === action.payload.innings_number
        );
        if (inningsIndex !== -1) {
          // Update existing innings and preserve overs array
          const existingOvers =
            state.scorecard.innings?.[inningsIndex]?.overs || [];
          state.scorecard.innings![inningsIndex] = {
            ...action.payload,
            overs: existingOvers,
          };
        } else {
          // Create new innings with empty overs array
          state.scorecard.innings.push({
            ...action.payload,
            overs: [],
          });
        }
      }
    },
    fetchLatestOverSuccess: (
      state,
      action: PayloadAction<{ inningsNumber: number; over: OverSummary }>
    ) => {
      state.loading = false;
      if (state.scorecard) {
        const inningsIndex = state.scorecard.innings.findIndex(
          innings => innings.innings_number === action.payload.inningsNumber
        );
        if (inningsIndex !== -1) {
          // Initialize overs array if it doesn't exist
          if (!state.scorecard.innings![inningsIndex]!.overs) {
            state.scorecard.innings![inningsIndex]!.overs = [];
          }

          const overIndex = state.scorecard.innings![
            inningsIndex
          ]!.overs!.findIndex(
            over => over.over_number === action.payload.over.over_number
          );

          if (overIndex !== -1) {
            // Update existing over
            state.scorecard.innings![inningsIndex]!.overs![overIndex] =
              action.payload.over;
          } else {
            // Add new over
            state.scorecard.innings![inningsIndex]!.overs!.push(
              action.payload.over
            );
          }
        }
      }
    },
    fetchInningsFailure: (state, action: PayloadAction<string>) => {
      state.loading = false;
      state.error = action.payload;
    },
    clearScorecard: state => {
      state.scorecard = null;
      state.error = null;
      state.inningsTransitionDetected = false;
    },
    clearInningsTransition: state => {
      state.inningsTransitionDetected = false;
    },
    triggerInningsTransition: state => {
      // This action will be handled by the saga
    },
  },
  extraReducers: builder => {
    builder
      .addCase(fetchInningsScoreSummaryThunk.fulfilled, (state, action) => {
        state.loading = false;
        if (state.scorecard) {
          const inningsIndex = state.scorecard.innings.findIndex(
            innings => innings.innings_number === action.payload.innings_number
          );

          // Check for innings transition (status changed from 'in_progress' to 'completed')
          let inningsTransitionDetected = false;
          if (inningsIndex !== -1) {
            const existingInnings = state.scorecard.innings[inningsIndex];
            const wasInProgress = existingInnings?.status === 'in_progress';
            const isNowCompleted = action.payload.status === 'completed';
            inningsTransitionDetected = wasInProgress && isNowCompleted;
          }

          if (inningsIndex !== -1) {
            // Update existing innings and preserve overs array
            const existingOvers =
              state.scorecard.innings?.[inningsIndex]?.overs || [];
            state.scorecard.innings![inningsIndex] = {
              ...action.payload,
              batting_team: action.payload.batting_team as 'A' | 'B',
              extras: {
                byes: 0,
                leg_byes: 0,
                wides: 0,
                no_balls: 0,
                total: action.payload.extras?.total || 0,
              },
              overs: existingOvers,
            };
          } else {
            // Create new innings with empty overs array
            state.scorecard.innings.push({
              ...action.payload,
              batting_team: action.payload.batting_team as 'A' | 'B',
              extras: {
                byes: 0,
                leg_byes: 0,
                wides: 0,
                no_balls: 0,
                total: action.payload.extras?.total || 0,
              },
              overs: [],
            });
          }

          // Set flag to indicate innings transition was detected
          if (inningsTransitionDetected) {
            state.inningsTransitionDetected = true;
            // Trigger the innings transition saga
            // Note: We can't dispatch actions from reducers, so we'll handle this differently
          }
        }
      })
      .addCase(fetchInningsScoreSummaryThunk.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to fetch innings score summary';
      })
      .addCase(fetchLatestOverThunk.fulfilled, (state, action) => {
        state.loading = false;
        if (state.scorecard) {
          const inningsIndex = state.scorecard.innings.findIndex(
            innings => innings.innings_number === action.payload.inningsNumber
          );
          if (inningsIndex !== -1) {
            // Initialize overs array if it doesn't exist
            if (!state.scorecard.innings![inningsIndex]!.overs) {
              state.scorecard.innings![inningsIndex]!.overs = [];
            }

            const overIndex = state.scorecard.innings![
              inningsIndex
            ]!.overs!.findIndex(
              over => over.over_number === action.payload.over.over_number
            );

            if (overIndex !== -1) {
              // Update existing over
              state.scorecard.innings![inningsIndex]!.overs![overIndex] = {
                ...action.payload.over,
                balls: action.payload.over.balls.map((ball: GraphQLBallSummary) => ({
                  ...ball,
                  ball_type: ball.ball_type as BallType,
                  run_type: ball.run_type as RunType,
                })),
              };
            } else {
              // Add new over
              state.scorecard.innings![inningsIndex]!.overs!.push({
                ...action.payload.over,
                balls: action.payload.over.balls.map((ball: GraphQLBallSummary) => ({
                  ...ball,
                  ball_type: ball.ball_type as BallType,
                  run_type: ball.run_type as RunType,
                })),
              });
            }
          }
        }
      })
      .addCase(fetchLatestOverThunk.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to fetch latest over';
      })
      .addCase(fetchAllOversDetailsThunk.fulfilled, (state, action) => {
        state.loading = false;
        if (state.scorecard) {
          const inningsIndex = state.scorecard.innings.findIndex(
            innings => innings.innings_number === action.payload.inningsNumber
          );
          if (inningsIndex !== -1) {
            // Update all overs for the innings
            state.scorecard.innings![inningsIndex]!.overs = action.payload.overs.map((over: GraphQLOverSummary) => ({
              ...over,
              balls: over.balls.map((ball: GraphQLBallSummary) => ({
                ...ball,
                ball_type: ball.ball_type as BallType,
                run_type: ball.run_type as RunType,
              })),
            }));
          }
        }
      })
      .addCase(fetchAllOversDetailsThunk.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to fetch all overs details';
      })
      .addCase(undoBallThunk.fulfilled, (state, action) => {
        state.scoring = false;
        // The scorecard will be refetched to get updated data
      })
      .addCase(undoBallThunk.rejected, (state, action) => {
        state.scoring = false;
        state.error = action.error.message || 'Failed to undo ball';
      });
  },
});

export const {
  fetchScorecardRequest,
  fetchScorecardSuccess,
  fetchScorecardFailure,
  startScoringRequest,
  startScoringSuccess,
  startScoringFailure,
  addBallRequest,
  addBallSuccess,
  addBallFailure,
  undoBallRequest,
  undoBallSuccess,
  undoBallFailure,
  fetchInningsRequest,
  fetchInningsSuccess,
  fetchInningsScoreSummarySuccess,
  fetchLatestOverSuccess,
  fetchInningsFailure,
  clearScorecard,
  clearInningsTransition,
  triggerInningsTransition,
} = scorecardSlice.actions;

export default scorecardSlice.reducer;
