/* eslint-disable @typescript-eslint/no-explicit-any */
import { runSaga } from 'redux-saga';
import {
  fetchScorecardSaga,
  startScoringSaga,
  addBallSaga,
  fetchInningsSaga,
} from '../../store/sagas/scorecardSaga';
import {
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
  fetchInningsFailure,
  ScorecardResponse,
  InningsSummary,
  BallEventRequest,
} from '../../store/reducers/scorecardSlice';
import { ApiError } from '../../services/api';
import { graphqlService } from '../../services/graphqlService';

// Mock the API service
// const mockApiService = {
//   startScoring: jest.fn(),
//   addBall: jest.fn(),
// };

// Mock the GraphQL service
// const mockGraphqlService = {
//   getLiveScorecard: jest.fn(),
//   getInningsDetails: jest.fn(),
// };

jest.mock('../../services/api', () => ({
  ApiService: jest.fn().mockImplementation(() => ({
    startScoring: jest.fn(),
    addBall: jest.fn(),
  })),
  ApiError: class extends Error {
    public status?: number;
    constructor(message: string, status?: number) {
      super(message);
      this.name = 'ApiError';
      if (status !== undefined) {
        this.status = status;
      }
    }
  },
}));

jest.mock('../../services/graphqlService', () => ({
  graphqlService: {
    getLiveScorecard: jest.fn(),
    getInningsDetails: jest.fn(),
  },
}));

describe('Scorecard Sagas', () => {
  let dispatched: unknown[];

  beforeEach(() => {
    dispatched = [];
    jest.clearAllMocks();
  });

  const mockStore = {
    dispatch: (action: unknown) => dispatched.push(action),
    getState: (): unknown => ({
      scorecard: {
        scorecard: null, // No innings data exists, so saga will call fetchScorecardRequest
      },
    }),
  };

  describe('fetchScorecardSaga', () => {
    const mockScorecardData: ScorecardResponse = {
      match_id: 'match-1',
      match_number: 1,
      series_name: 'Test Series',
      team_a: 'Team A',
      team_b: 'Team B',
      total_overs: 20,
      toss_winner: 'A',
      toss_type: 'H',
      current_innings: 1,
      match_status: 'live',
      innings: [
        {
          innings_number: 1,
          batting_team: 'A',
          total_runs: 45,
          total_wickets: 2,
          total_overs: 5,
          total_balls: 30,
          status: 'in_progress',
          extras: {
            byes: 2,
            leg_byes: 1,
            wides: 3,
            no_balls: 1,
            total: 7,
          },
          overs: [],
        },
      ],
    };

    it('should fetch scorecard successfully', async () => {
      (graphqlService.getLiveScorecard as jest.Mock).mockResolvedValueOnce({
        success: true,
        data: mockScorecardData,
      });

      await (runSaga as any)(
        {
          dispatch: mockStore.dispatch,
          getState: mockStore.getState,
        },
        fetchScorecardSaga,
        fetchScorecardRequest('match-1') as any
      ).toPromise();

      expect(graphqlService.getLiveScorecard).toHaveBeenCalledWith('match-1');
      expect(dispatched).toContainEqual(
        fetchScorecardSuccess(mockScorecardData)
      );
    });

    it('should handle API errors', async () => {
      (graphqlService.getLiveScorecard as jest.Mock).mockResolvedValueOnce({
        success: false,
        error: 'Scorecard not found',
      });

      await (runSaga as any)(
        {
          dispatch: mockStore.dispatch,
          getState: mockStore.getState,
        },
        fetchScorecardSaga,
        fetchScorecardRequest('match-1') as any
      ).toPromise();

      expect(dispatched).toContainEqual(
        fetchScorecardFailure('Scorecard not found')
      );
    });

    it('should handle generic errors', async () => {
      (graphqlService.getLiveScorecard as jest.Mock).mockRejectedValueOnce(
        new Error('Network error')
      );

      await (runSaga as any)(
        {
          dispatch: mockStore.dispatch,
          getState: mockStore.getState,
        },
        fetchScorecardSaga,
        fetchScorecardRequest('match-1') as any
      ).toPromise();

      expect(dispatched).toContainEqual(fetchScorecardFailure('Network error'));
    });

    it('should handle nested data response', async () => {
      (graphqlService.getLiveScorecard as jest.Mock).mockResolvedValueOnce({
        success: true,
        data: mockScorecardData,
      });

      await (runSaga as any)(
        {
          dispatch: mockStore.dispatch,
          getState: mockStore.getState,
        },
        fetchScorecardSaga,
        fetchScorecardRequest('match-1') as any
      ).toPromise();

      expect(dispatched).toContainEqual(
        fetchScorecardSuccess(mockScorecardData)
      );
    });
  });

  describe('startScoringSaga', () => {
    it('should start scoring successfully', async () => {
      const mockStartScoring = jest.fn().mockResolvedValueOnce({
        data: { message: 'Scoring started', match_id: 'match-1' },
        status: 200,
        message: 'Success',
      });

      // eslint-disable-next-line @typescript-eslint/no-require-imports
      const ApiService = require('../../services/api').ApiService;
      ApiService.mockImplementation(() => ({
        startScoring: mockStartScoring,
      }));

      await (runSaga as any)(
        {
          dispatch: mockStore.dispatch,
          getState: mockStore.getState,
        },
        startScoringSaga,
        startScoringRequest('match-1') as any
      ).toPromise();

      expect(mockStartScoring).toHaveBeenCalledWith('match-1');
      expect(dispatched).toContainEqual(startScoringSuccess());
      expect(dispatched).toContainEqual(
        fetchScorecardRequest('match-1') as any
      );
    });

    it('should handle start scoring API errors', async () => {
      const error = new ApiError('Match is not ready for scoring', 400);
      const mockStartScoring = jest.fn().mockRejectedValueOnce(error);

      // eslint-disable-next-line @typescript-eslint/no-require-imports
      const ApiService = require('../../services/api').ApiService;
      ApiService.mockImplementation(() => ({
        startScoring: mockStartScoring,
      }));

      await (runSaga as any)(
        {
          dispatch: mockStore.dispatch,
          getState: mockStore.getState,
        },
        startScoringSaga,
        startScoringRequest('match-1') as any
      ).toPromise();

      expect(dispatched).toContainEqual(
        startScoringFailure('Match is not ready for scoring')
      );
    });

    it('should handle start scoring generic errors', async () => {
      const mockStartScoring = jest
        .fn()
        .mockRejectedValueOnce(new Error('Network error'));

      // eslint-disable-next-line @typescript-eslint/no-require-imports
      const ApiService = require('../../services/api').ApiService;
      ApiService.mockImplementation(() => ({
        startScoring: mockStartScoring,
      }));

      await (runSaga as any)(
        {
          dispatch: mockStore.dispatch,
          getState: mockStore.getState,
        },
        startScoringSaga,
        startScoringRequest('match-1') as any
      ).toPromise();

      expect(dispatched).toContainEqual(
        startScoringFailure('Failed to start scoring')
      );
    });
  });

  describe('addBallSaga', () => {
    const mockBallEvent: BallEventRequest = {
      match_id: 'match-1',
      innings_number: 1,
      ball_type: 'good',
      run_type: '4',
      runs: 4,
      is_wicket: false,
      byes: 0,
    };

    it('should add ball successfully', async () => {
      const mockAddBall = jest.fn().mockResolvedValueOnce({
        data: {
          message: 'Ball added successfully',
          match_id: 'match-1',
          innings_number: 1,
          ball_type: 'good',
          run_type: '4',
          runs: 4,
          byes: 0,
          is_wicket: false,
        },
        status: 200,
        message: 'Success',
      });

      // eslint-disable-next-line @typescript-eslint/no-require-imports
      const ApiService = require('../../services/api').ApiService;
      ApiService.mockImplementation(() => ({
        addBall: mockAddBall,
      }));

      await (runSaga as any)(
        {
          dispatch: mockStore.dispatch,
          getState: mockStore.getState,
        },
        addBallSaga,
        addBallRequest(mockBallEvent) as any
      ).toPromise();

      expect(mockAddBall).toHaveBeenCalledWith(mockBallEvent);
      expect(dispatched).toContainEqual(addBallSuccess());
      // The saga will call fetchScorecardRequest since no innings data exists in state
      expect(dispatched).toContainEqual(
        fetchScorecardRequest('match-1') as any
      );
    });

    it('should add wicket ball successfully', async () => {
      const wicketBallEvent: BallEventRequest = {
        match_id: 'match-1',
        innings_number: 1,
        ball_type: 'good',
        run_type: 'WC',
        runs: 0,
        is_wicket: true,
        wicket_type: 'bowled',
        byes: 0,
      };

      const mockAddBall = jest.fn().mockResolvedValueOnce({
        data: {
          message: 'Wicket ball added successfully',
          match_id: 'match-1',
          innings_number: 1,
          ball_type: 'good',
          run_type: 'WC',
          runs: 0,
          byes: 0,
          is_wicket: true,
          wicket_type: 'bowled',
        },
        status: 200,
        message: 'Success',
      });

      // eslint-disable-next-line @typescript-eslint/no-require-imports
      const ApiService = require('../../services/api').ApiService;
      ApiService.mockImplementation(() => ({
        addBall: mockAddBall,
      }));

      await (runSaga as any)(
        {
          dispatch: mockStore.dispatch,
          getState: mockStore.getState,
        },
        addBallSaga,
        addBallRequest(wicketBallEvent)
      ).toPromise();

      expect(mockAddBall).toHaveBeenCalledWith(wicketBallEvent);
      expect(dispatched).toContainEqual(addBallSuccess());
      expect(dispatched).toContainEqual(
        fetchScorecardRequest('match-1') as any
      );
    });

    it('should add wide ball successfully', async () => {
      const wideBallEvent: BallEventRequest = {
        match_id: 'match-1',
        innings_number: 1,
        ball_type: 'wide',
        run_type: 'WD',
        runs: 1,
        is_wicket: false,
        byes: 0,
      };

      const mockAddBall = jest.fn().mockResolvedValueOnce({
        data: {
          message: 'Wide ball added successfully',
          match_id: 'match-1',
          innings_number: 1,
          ball_type: 'wide',
          run_type: 'WD',
          runs: 1,
          byes: 0,
          is_wicket: false,
        },
        status: 200,
        message: 'Success',
      });

      // eslint-disable-next-line @typescript-eslint/no-require-imports
      const ApiService = require('../../services/api').ApiService;
      ApiService.mockImplementation(() => ({
        addBall: mockAddBall,
      }));

      await (runSaga as any)(
        {
          dispatch: mockStore.dispatch,
          getState: mockStore.getState,
        },
        addBallSaga,
        addBallRequest(wideBallEvent)
      ).toPromise();

      expect(mockAddBall).toHaveBeenCalledWith(wideBallEvent);
      expect(dispatched).toContainEqual(addBallSuccess());
    });

    it('should add no ball successfully', async () => {
      const noBallEvent: BallEventRequest = {
        match_id: 'match-1',
        innings_number: 1,
        ball_type: 'no_ball',
        run_type: 'NB',
        runs: 1,
        is_wicket: false,
        byes: 0,
      };

      const mockAddBall = jest.fn().mockResolvedValueOnce({
        data: {
          message: 'No ball added successfully',
          match_id: 'match-1',
          innings_number: 1,
          ball_type: 'no_ball',
          run_type: 'NB',
          runs: 1,
          byes: 0,
          is_wicket: false,
        },
        status: 200,
        message: 'Success',
      });

      // eslint-disable-next-line @typescript-eslint/no-require-imports
      const ApiService = require('../../services/api').ApiService;
      ApiService.mockImplementation(() => ({
        addBall: mockAddBall,
      }));

      await (runSaga as any)(
        {
          dispatch: mockStore.dispatch,
          getState: mockStore.getState,
        },
        addBallSaga,
        addBallRequest(noBallEvent)
      ).toPromise();

      expect(mockAddBall).toHaveBeenCalledWith(noBallEvent);
      expect(dispatched).toContainEqual(addBallSuccess());
    });

    it('should add ball with byes successfully', async () => {
      const ballWithByes: BallEventRequest = {
        match_id: 'match-1',
        innings_number: 1,
        ball_type: 'good',
        run_type: '1',
        runs: 1,
        is_wicket: false,
        byes: 2,
      };

      const mockAddBall = jest.fn().mockResolvedValueOnce({
        data: {
          message: 'Ball with byes added successfully',
          match_id: 'match-1',
          innings_number: 1,
          ball_type: 'good',
          run_type: '1',
          runs: 1,
          byes: 2,
          is_wicket: false,
        },
        status: 200,
        message: 'Success',
      });

      // eslint-disable-next-line @typescript-eslint/no-require-imports
      const ApiService = require('../../services/api').ApiService;
      ApiService.mockImplementation(() => ({
        addBall: mockAddBall,
      }));

      await (runSaga as any)(
        {
          dispatch: mockStore.dispatch,
          getState: mockStore.getState,
        },
        addBallSaga,
        addBallRequest(ballWithByes)
      ).toPromise();

      expect(mockAddBall).toHaveBeenCalledWith(ballWithByes);
      expect(dispatched).toContainEqual(addBallSuccess());
    });

    it('should handle add ball API errors', async () => {
      const error = new ApiError('Invalid ball data', 400);

      const mockAddBall = jest.fn().mockRejectedValueOnce(error);

      // eslint-disable-next-line @typescript-eslint/no-require-imports
      const ApiService = require('../../services/api').ApiService;
      ApiService.mockImplementation(() => ({
        addBall: mockAddBall,
      }));

      await (runSaga as any)(
        {
          dispatch: mockStore.dispatch,
          getState: mockStore.getState,
        },
        addBallSaga,
        addBallRequest(mockBallEvent) as any
      ).toPromise();

      expect(dispatched).toContainEqual(addBallFailure('Invalid ball data'));
    });

    it('should handle add ball generic errors', async () => {
      const mockAddBall = jest
        .fn()
        .mockRejectedValueOnce(new Error('Network error'));

      // eslint-disable-next-line @typescript-eslint/no-require-imports
      const ApiService = require('../../services/api').ApiService;
      ApiService.mockImplementation(() => ({
        addBall: mockAddBall,
      }));

      await (runSaga as any)(
        {
          dispatch: mockStore.dispatch,
          getState: mockStore.getState,
        },
        addBallSaga,
        addBallRequest(mockBallEvent) as any
      ).toPromise();

      expect(dispatched).toContainEqual(addBallFailure('Failed to add ball'));
    });

    it('should handle innings completion error', async () => {
      const error = new ApiError('Innings already completed', 409);

      const mockAddBall = jest.fn().mockRejectedValueOnce(error);

      // eslint-disable-next-line @typescript-eslint/no-require-imports
      const ApiService = require('../../services/api').ApiService;
      ApiService.mockImplementation(() => ({
        addBall: mockAddBall,
      }));

      await (runSaga as any)(
        {
          dispatch: mockStore.dispatch,
          getState: mockStore.getState,
        },
        addBallSaga,
        addBallRequest(mockBallEvent) as any
      ).toPromise();

      expect(dispatched).toContainEqual(
        addBallFailure('Innings already completed')
      );
    });
  });

  describe('fetchInningsSaga', () => {
    const mockInningsData: InningsSummary = {
      innings_number: 1,
      batting_team: 'A',
      total_runs: 120,
      total_wickets: 3,
      total_overs: 10,
      total_balls: 60,
      status: 'completed',
      extras: {
        byes: 5,
        leg_byes: 2,
        wides: 8,
        no_balls: 1,
        total: 16,
      },
      overs: [],
    };

    it('should fetch innings successfully', async () => {
      (graphqlService.getInningsDetails as jest.Mock).mockResolvedValueOnce({
        success: true,
        data: mockInningsData,
      });

      await (runSaga as any)(
        {
          dispatch: mockStore.dispatch,
          getState: mockStore.getState,
        },
        fetchInningsSaga,
        fetchInningsRequest({ matchId: 'match-1', inningsNumber: 1 })
      ).toPromise();

      expect(graphqlService.getInningsDetails).toHaveBeenCalledWith(
        'match-1',
        1
      );
      expect(dispatched).toContainEqual(fetchInningsSuccess(mockInningsData));
    });

    it('should handle fetch innings API errors', async () => {
      (graphqlService.getInningsDetails as jest.Mock).mockResolvedValueOnce({
        success: false,
        error: 'Innings not found',
      });

      await (runSaga as any)(
        {
          dispatch: mockStore.dispatch,
          getState: mockStore.getState,
        },
        fetchInningsSaga,
        fetchInningsRequest({ matchId: 'match-1', inningsNumber: 1 })
      ).toPromise();

      expect(dispatched).toContainEqual(
        fetchInningsFailure('Innings not found')
      );
    });

    it('should handle fetch innings generic errors', async () => {
      (graphqlService.getInningsDetails as jest.Mock).mockRejectedValueOnce(
        new Error('Network error')
      );

      await (runSaga as any)(
        {
          dispatch: mockStore.dispatch,
          getState: mockStore.getState,
        },
        fetchInningsSaga,
        fetchInningsRequest({ matchId: 'match-1', inningsNumber: 1 })
      ).toPromise();

      expect(dispatched).toContainEqual(fetchInningsFailure('Network error'));
    });

    it('should handle nested data response', async () => {
      (graphqlService.getInningsDetails as jest.Mock).mockResolvedValueOnce({
        success: true,
        data: mockInningsData,
      });

      await (runSaga as any)(
        {
          dispatch: mockStore.dispatch,
          getState: mockStore.getState,
        },
        fetchInningsSaga,
        fetchInningsRequest({ matchId: 'match-1', inningsNumber: 1 })
      ).toPromise();

      expect(dispatched).toContainEqual(fetchInningsSuccess(mockInningsData));
    });
  });
});
