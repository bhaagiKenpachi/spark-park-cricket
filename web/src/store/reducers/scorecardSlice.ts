/* eslint-disable @typescript-eslint/no-unused-vars */
import { createSlice, PayloadAction } from '@reduxjs/toolkit';

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
}

export const initialState: ScorecardState = {
  scorecard: null,
  loading: false,
  error: null,
  scoring: false,
};

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
    },
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
  fetchInningsRequest,
  fetchInningsSuccess,
  fetchInningsScoreSummarySuccess,
  fetchLatestOverSuccess,
  fetchInningsFailure,
  clearScorecard,
} = scorecardSlice.actions;

export default scorecardSlice.reducer;
