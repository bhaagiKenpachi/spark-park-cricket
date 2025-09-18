import {
  scorecardSlice,
  initialState,
  InningsSummary,
  OverSummary,
} from '../scorecardSlice';

describe('scorecardSlice', () => {
  describe('fetchInningsScoreSummarySuccess', () => {
    it('should create new innings with empty overs array when innings does not exist', () => {
      const mockInningsData: InningsSummary = {
        innings_number: 1,
        batting_team: 'A',
        total_runs: 10,
        total_wickets: 1,
        total_overs: 2,
        total_balls: 12,
        status: 'in_progress',
        extras: { byes: 0, leg_byes: 0, wides: 1, no_balls: 0, total: 2 },
      };

      const state = {
        ...initialState,
        scorecard: {
          match_id: 'test-match',
          series_id: 'test-series',
          series_name: 'Test Series',
          match_number: 1,
          date: '2025-01-01',
          status: 'live',
          team_a: 'Team A',
          team_b: 'Team B',
          total_overs: 20,
          toss_winner: 'Team A',
          toss_type: 'bat',
          venue: 'Test Venue',
          innings: [],
        },
      };

      const action = {
        type: 'scorecard/fetchInningsScoreSummarySuccess',
        payload: mockInningsData,
      };

      const newState = scorecardSlice.reducer(state, action);

      expect(newState.scorecard?.innings).toHaveLength(1);
      expect(newState.scorecard?.innings[0]).toEqual({
        ...mockInningsData,
        overs: [],
      });
    });

    it('should update existing innings and preserve overs array', () => {
      const existingOvers: OverSummary[] = [
        {
          over_number: 1,
          total_runs: 5,
          total_balls: 6,
          total_wickets: 0,
          status: 'completed',
          balls: [],
        },
      ];

      const mockInningsData: InningsSummary = {
        innings_number: 1,
        batting_team: 'A',
        total_runs: 15,
        total_wickets: 2,
        total_overs: 3,
        total_balls: 18,
        status: 'in_progress',
        extras: { byes: 0, leg_byes: 0, wides: 1, no_balls: 0, total: 3 },
      };

      const state = {
        ...initialState,
        scorecard: {
          match_id: 'test-match',
          series_id: 'test-series',
          match_number: 1,
          date: '2025-01-01',
          status: 'live',
          team_a: 'Team A',
          team_b: 'Team B',
          innings: [
            {
              innings_number: 1,
              batting_team: 'A',
              total_runs: 10,
              total_wickets: 1,
              total_overs: 2,
              total_balls: 12,
              status: 'in_progress',
              extras: { byes: 0, leg_byes: 0, wides: 1, no_balls: 0, total: 2 },
              overs: existingOvers,
            },
          ],
        },
      };

      const action = {
        type: 'scorecard/fetchInningsScoreSummarySuccess',
        payload: mockInningsData,
      };

      const newState = scorecardSlice.reducer(state, action);

      expect(newState.scorecard?.innings).toHaveLength(1);
      expect(newState.scorecard?.innings[0].total_runs).toBe(15);
      expect(newState.scorecard?.innings[0].total_wickets).toBe(2);
      expect(newState.scorecard?.innings[0].overs).toEqual(existingOvers);
    });

    it('should initialize overs array if it does not exist in existing innings', () => {
      const mockInningsData: InningsSummary = {
        innings_number: 1,
        batting_team: 'A',
        total_runs: 15,
        total_wickets: 2,
        total_overs: 3,
        total_balls: 18,
        status: 'in_progress',
        extras: { byes: 0, leg_byes: 0, wides: 1, no_balls: 0, total: 3 },
      };

      const state = {
        ...initialState,
        scorecard: {
          match_id: 'test-match',
          series_id: 'test-series',
          match_number: 1,
          date: '2025-01-01',
          status: 'live',
          team_a: 'Team A',
          team_b: 'Team B',
          innings: [
            {
              innings_number: 1,
              batting_team: 'A',
              total_runs: 10,
              total_wickets: 1,
              total_overs: 2,
              total_balls: 12,
              status: 'in_progress',
              extras: { byes: 0, leg_byes: 0, wides: 1, no_balls: 0, total: 2 },
              // Note: no overs property
            },
          ],
        },
      };

      const action = {
        type: 'scorecard/fetchInningsScoreSummarySuccess',
        payload: mockInningsData,
      };

      const newState = scorecardSlice.reducer(state, action);

      expect(newState.scorecard?.innings[0].overs).toEqual([]);
    });
  });

  describe('fetchLatestOverSuccess', () => {
    it('should add new over to existing innings', () => {
      const mockOverData: OverSummary = {
        over_number: 1,
        total_runs: 6,
        total_balls: 6,
        total_wickets: 0,
        status: 'completed',
        balls: [
          {
            ball_number: 1,
            ball_type: 'good',
            run_type: 'boundary',
            runs: 1,
            byes: 0,
            is_wicket: false,
            wicket_type: undefined,
          },
        ],
      };

      const state = {
        ...initialState,
        scorecard: {
          match_id: 'test-match',
          series_id: 'test-series',
          match_number: 1,
          date: '2025-01-01',
          status: 'live',
          team_a: 'Team A',
          team_b: 'Team B',
          innings: [
            {
              innings_number: 1,
              batting_team: 'A',
              total_runs: 0,
              total_wickets: 0,
              total_overs: 0,
              total_balls: 0,
              status: 'in_progress',
              extras: { byes: 0, leg_byes: 0, wides: 0, no_balls: 0, total: 0 },
              overs: [],
            },
          ],
        },
      };

      const action = {
        type: 'scorecard/fetchLatestOverSuccess',
        payload: {
          inningsNumber: 1,
          over: mockOverData,
        },
      };

      const newState = scorecardSlice.reducer(state, action);

      expect(newState.scorecard?.innings[0].overs).toHaveLength(1);
      expect(newState.scorecard?.innings[0].overs[0]).toEqual(mockOverData);
    });

    it('should update existing over in innings', () => {
      const existingOver: OverSummary = {
        over_number: 1,
        total_runs: 3,
        total_balls: 3,
        total_wickets: 0,
        status: 'in_progress',
        balls: [],
      };

      const updatedOver: OverSummary = {
        over_number: 1,
        total_runs: 6,
        total_balls: 6,
        total_wickets: 0,
        status: 'completed',
        balls: [
          {
            ball_number: 1,
            ball_type: 'good',
            run_type: 'boundary',
            runs: 1,
            byes: 0,
            is_wicket: false,
            wicket_type: undefined,
          },
        ],
      };

      const state = {
        ...initialState,
        scorecard: {
          match_id: 'test-match',
          series_id: 'test-series',
          match_number: 1,
          date: '2025-01-01',
          status: 'live',
          team_a: 'Team A',
          team_b: 'Team B',
          innings: [
            {
              innings_number: 1,
              batting_team: 'A',
              total_runs: 0,
              total_wickets: 0,
              total_overs: 0,
              total_balls: 0,
              status: 'in_progress',
              extras: { byes: 0, leg_byes: 0, wides: 0, no_balls: 0, total: 0 },
              overs: [existingOver],
            },
          ],
        },
      };

      const action = {
        type: 'scorecard/fetchLatestOverSuccess',
        payload: {
          inningsNumber: 1,
          over: updatedOver,
        },
      };

      const newState = scorecardSlice.reducer(state, action);

      expect(newState.scorecard?.innings[0].overs).toHaveLength(1);
      expect(newState.scorecard?.innings[0].overs[0]).toEqual(updatedOver);
    });

    it('should create overs array if it does not exist', () => {
      const mockOverData: OverSummary = {
        over_number: 1,
        total_runs: 6,
        total_balls: 6,
        total_wickets: 0,
        status: 'completed',
        balls: [],
      };

      const state = {
        ...initialState,
        scorecard: {
          match_id: 'test-match',
          series_id: 'test-series',
          match_number: 1,
          date: '2025-01-01',
          status: 'live',
          team_a: 'Team A',
          team_b: 'Team B',
          innings: [
            {
              innings_number: 1,
              batting_team: 'A',
              total_runs: 0,
              total_wickets: 0,
              total_overs: 0,
              total_balls: 0,
              status: 'in_progress',
              extras: { byes: 0, leg_byes: 0, wides: 0, no_balls: 0, total: 0 },
              // Note: no overs property
            },
          ],
        },
      };

      const action = {
        type: 'scorecard/fetchLatestOverSuccess',
        payload: {
          inningsNumber: 1,
          over: mockOverData,
        },
      };

      const newState = scorecardSlice.reducer(state, action);

      expect(newState.scorecard?.innings[0].overs).toBeDefined();
      expect(newState.scorecard?.innings[0].overs).toHaveLength(1);
      expect(newState.scorecard?.innings[0].overs[0]).toEqual(mockOverData);
    });

    it('should not update if innings does not exist', () => {
      const mockOverData: OverSummary = {
        over_number: 1,
        total_runs: 6,
        total_balls: 6,
        total_wickets: 0,
        status: 'completed',
        balls: [],
      };

      const state = {
        ...initialState,
        scorecard: {
          match_id: 'test-match',
          series_id: 'test-series',
          series_name: 'Test Series',
          match_number: 1,
          date: '2025-01-01',
          status: 'live',
          team_a: 'Team A',
          team_b: 'Team B',
          total_overs: 20,
          toss_winner: 'Team A',
          toss_type: 'bat',
          venue: 'Test Venue',
          innings: [],
        },
      };

      const action = {
        type: 'scorecard/fetchLatestOverSuccess',
        payload: {
          inningsNumber: 1,
          over: mockOverData,
        },
      };

      const newState = scorecardSlice.reducer(state, action);

      expect(newState.scorecard?.innings).toHaveLength(0);
    });
  });

  describe('addBallSuccess', () => {
    it('should set scoring to false', () => {
      const state = {
        ...initialState,
        scoring: true,
        error: 'Some error',
      };

      const action = {
        type: 'scorecard/addBallSuccess',
        payload: undefined,
      };

      const newState = scorecardSlice.reducer(state, action);

      expect(newState.scoring).toBe(false);
      expect(newState.error).toBe('Some error'); // Error should remain unchanged
    });
  });

  describe('addBallFailure', () => {
    it('should set scoring to false and set error message', () => {
      const state = {
        ...initialState,
        scoring: true,
        error: null,
      };

      const action = {
        type: 'scorecard/addBallFailure',
        payload: 'Failed to add ball',
      };

      const newState = scorecardSlice.reducer(state, action);

      expect(newState.scoring).toBe(false);
      expect(newState.error).toBe('Failed to add ball');
    });
  });
});
