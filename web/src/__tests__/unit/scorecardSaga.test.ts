import { runSaga } from 'redux-saga';
import { call, put } from 'redux-saga/effects';
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
import { ApiService, ApiError } from '../../services/api';

// Mock the API service
const mockApiService = {
  getScorecard: jest.fn(),
  startScoring: jest.fn(),
  addBall: jest.fn(),
  getInnings: jest.fn(),
};

jest.mock('../../services/api', () => ({
  ApiService: jest.fn().mockImplementation(() => mockApiService),
  ApiError: class extends Error {
    constructor(message: string) {
      super(message);
      this.name = 'ApiError';
    }
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
    getState: () => ({}),
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
      mockApiService.getScorecard.mockResolvedValueOnce({
        data: mockScorecardData,
        status: 200,
        message: 'Success',
      });

      await runSaga(
        mockStore,
        fetchScorecardSaga,
        fetchScorecardRequest('match-1')
      );

      expect(mockApiService.getScorecard).toHaveBeenCalledWith('match-1');
      expect(dispatched).toContainEqual(
        fetchScorecardSuccess(mockScorecardData)
      );
    });

    it('should handle API errors', async () => {
      const error = new ApiError('Scorecard not found');
      mockApiService.getScorecard.mockRejectedValueOnce(error);

      await runSaga(
        mockStore,
        fetchScorecardSaga,
        fetchScorecardRequest('match-1')
      );

      expect(dispatched).toContainEqual(
        fetchScorecardFailure('Scorecard not found')
      );
    });

    it('should handle generic errors', async () => {
      mockApiService.getScorecard.mockRejectedValueOnce(
        new Error('Network error')
      );

      await runSaga(
        mockStore,
        fetchScorecardSaga,
        fetchScorecardRequest('match-1')
      );

      expect(dispatched).toContainEqual(
        fetchScorecardFailure('Failed to fetch scorecard')
      );
    });

    it('should handle nested data response', async () => {
      mockApiService.getScorecard.mockResolvedValueOnce({
        data: { data: mockScorecardData },
        status: 200,
        message: 'Success',
      });

      await runSaga(
        mockStore,
        fetchScorecardSaga,
        fetchScorecardRequest('match-1')
      );

      expect(dispatched).toContainEqual(
        fetchScorecardSuccess(mockScorecardData)
      );
    });
  });

  describe('startScoringSaga', () => {
    it('should start scoring successfully', async () => {
      mockApiService.startScoring.mockResolvedValueOnce({
        data: { message: 'Scoring started', match_id: 'match-1' },
        status: 200,
        message: 'Success',
      });

      await runSaga(
        mockStore,
        startScoringSaga,
        startScoringRequest('match-1')
      ).toPromise();

      expect(mockApiService.startScoring).toHaveBeenCalledWith('match-1');
      expect(dispatched).toContainEqual(startScoringSuccess());
      expect(dispatched).toContainEqual(fetchScorecardRequest('match-1'));
    });

    it('should handle start scoring API errors', async () => {
      const error = new ApiError('Match is not ready for scoring');
      mockApiService.startScoring.mockRejectedValueOnce(error);

      await runSaga(
        mockStore,
        startScoringSaga,
        startScoringRequest('match-1')
      ).toPromise();

      expect(dispatched).toContainEqual(
        startScoringFailure('Match is not ready for scoring')
      );
    });

    it('should handle start scoring generic errors', async () => {
      mockApiService.startScoring.mockRejectedValueOnce(
        new Error('Network error')
      );

      await runSaga(
        mockStore,
        startScoringSaga,
        startScoringRequest('match-1')
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
      mockApiService.addBall.mockResolvedValueOnce({
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

      await runSaga(
        mockStore,
        addBallSaga,
        addBallRequest(mockBallEvent)
      ).toPromise();

      expect(mockApiService.addBall).toHaveBeenCalledWith(mockBallEvent);
      expect(dispatched).toContainEqual(addBallSuccess());
      expect(dispatched).toContainEqual(fetchScorecardRequest('match-1'));
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

      mockApiService.addBall.mockResolvedValueOnce({
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

      await runSaga(
        mockStore,
        addBallSaga,
        addBallRequest(wicketBallEvent)
      ).toPromise();

      expect(mockApiService.addBall).toHaveBeenCalledWith(wicketBallEvent);
      expect(dispatched).toContainEqual(addBallSuccess());
      expect(dispatched).toContainEqual(fetchScorecardRequest('match-1'));
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

      mockApiService.addBall.mockResolvedValueOnce({
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

      await runSaga(
        mockStore,
        addBallSaga,
        addBallRequest(wideBallEvent)
      ).toPromise();

      expect(mockApiService.addBall).toHaveBeenCalledWith(wideBallEvent);
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

      mockApiService.addBall.mockResolvedValueOnce({
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

      await runSaga(
        mockStore,
        addBallSaga,
        addBallRequest(noBallEvent)
      ).toPromise();

      expect(mockApiService.addBall).toHaveBeenCalledWith(noBallEvent);
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

      mockApiService.addBall.mockResolvedValueOnce({
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

      await runSaga(
        mockStore,
        addBallSaga,
        addBallRequest(ballWithByes)
      ).toPromise();

      expect(mockApiService.addBall).toHaveBeenCalledWith(ballWithByes);
      expect(dispatched).toContainEqual(addBallSuccess());
    });

    it('should handle add ball API errors', async () => {
      const error = new ApiError('Invalid ball data');
      mockApiService.addBall.mockRejectedValueOnce(error);

      await runSaga(
        mockStore,
        addBallSaga,
        addBallRequest(mockBallEvent)
      ).toPromise();

      expect(dispatched).toContainEqual(addBallFailure('Invalid ball data'));
    });

    it('should handle add ball generic errors', async () => {
      mockApiService.addBall.mockRejectedValueOnce(new Error('Network error'));

      await runSaga(
        mockStore,
        addBallSaga,
        addBallRequest(mockBallEvent)
      ).toPromise();

      expect(dispatched).toContainEqual(addBallFailure('Failed to add ball'));
    });

    it('should handle innings completion error', async () => {
      const error = new ApiError('Innings already completed');
      mockApiService.addBall.mockRejectedValueOnce(error);

      await runSaga(
        mockStore,
        addBallSaga,
        addBallRequest(mockBallEvent)
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
      mockApiService.getInnings.mockResolvedValueOnce({
        data: mockInningsData,
        status: 200,
        message: 'Success',
      });

      await runSaga(
        mockStore,
        fetchInningsSaga,
        fetchInningsRequest({ matchId: 'match-1', inningsNumber: 1 })
      ).toPromise();

      expect(mockApiService.getInnings).toHaveBeenCalledWith('match-1', 1);
      expect(dispatched).toContainEqual(fetchInningsSuccess(mockInningsData));
    });

    it('should handle fetch innings API errors', async () => {
      const error = new ApiError('Innings not found');
      mockApiService.getInnings.mockRejectedValueOnce(error);

      await runSaga(
        mockStore,
        fetchInningsSaga,
        fetchInningsRequest({ matchId: 'match-1', inningsNumber: 1 })
      ).toPromise();

      expect(dispatched).toContainEqual(
        fetchInningsFailure('Innings not found')
      );
    });

    it('should handle fetch innings generic errors', async () => {
      mockApiService.getInnings.mockRejectedValueOnce(
        new Error('Network error')
      );

      await runSaga(
        mockStore,
        fetchInningsSaga,
        fetchInningsRequest({ matchId: 'match-1', inningsNumber: 1 })
      ).toPromise();

      expect(dispatched).toContainEqual(
        fetchInningsFailure('Failed to fetch innings')
      );
    });

    it('should handle nested data response', async () => {
      mockApiService.getInnings.mockResolvedValueOnce({
        data: { data: mockInningsData },
        status: 200,
        message: 'Success',
      });

      await runSaga(
        mockStore,
        fetchInningsSaga,
        fetchInningsRequest({ matchId: 'match-1', inningsNumber: 1 })
      ).toPromise();

      expect(dispatched).toContainEqual(fetchInningsSuccess(mockInningsData));
    });
  });
});
