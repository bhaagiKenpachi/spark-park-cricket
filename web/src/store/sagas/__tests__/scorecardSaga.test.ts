/* eslint-disable */
import { put } from 'redux-saga/effects';
import {
  addBallSaga,
  fetchInningsScoreSummarySaga,
  fetchLatestOverSaga,
} from '../scorecardSaga';
import {
  addBallRequest,
  addBallSuccess,
  addBallFailure,
  fetchInningsScoreSummarySuccess,
  fetchLatestOverSuccess,
  fetchUpdatedScorecardRequest,
} from '../../reducers/scorecardSlice';
import { ApiService, ApiError } from '../../../services/api';
import { graphqlService } from '../../../services/graphqlService';

// Mock the services
jest.mock('../../../services/api', () => ({
  ApiService: jest.fn(),
  ApiError: class ApiError extends Error {
    status: number;
    details?: unknown;
    constructor(message: string, status: number, details?: unknown) {
      super(message);
      this.name = 'ApiError';
      this.status = status;
      this.details = details;
    }
  },
}));
jest.mock('../../../services/graphqlService', () => ({
  graphqlService: {
    getInningsScoreSummary: jest.fn(),
    getLatestOverOnly: jest.fn(),
    getLiveScorecard: jest.fn(),
  },
}));

const mockApiService = ApiService as jest.MockedClass<typeof ApiService>;
const mockGraphqlService = graphqlService as jest.Mocked<typeof graphqlService>;

describe('scorecardSaga', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('addBallSaga', () => {
    const mockBallEvent = {
      match_id: 'test-match-id',
      innings_number: 1,
      ball_type: 'good',
      run_type: 'runs',
      runs: 4,
      byes: 0,
      is_wicket: false,
    };

    it('should successfully add ball and fetch updated data sequentially', () => {
      const generator = addBallSaga(addBallRequest(mockBallEvent));

      // Mock successful API call
      const mockApiInstance = {
        addBall: jest.fn().mockResolvedValue({ success: true }),
      };
      mockApiService.mockImplementation(() => mockApiInstance as any);

      // Also mock the constructor to return our mock instance
      (mockApiService as any).mockImplementation(() => mockApiInstance);

      // Mock successful GraphQL responses
      const mockInningsResponse = {
        success: true,
        data: {
          innings_number: 1,
          batting_team: 'A',
          total_runs: 4,
          total_wickets: 0,
          total_overs: 1,
          total_balls: 1,
          status: 'in_progress',
          extras: { total: 0 },
        },
      };

      const mockOverResponse = {
        success: true,
        data: {
          over_number: 1,
          total_runs: 4,
          total_balls: 1,
          total_wickets: 0,
          status: 'in_progress',
          balls: [
            {
              ball_number: 1,
              ball_type: 'good',
              run_type: 'runs',
              runs: 4,
              byes: 0,
              is_wicket: false,
              wicket_type: null,
            },
          ],
        },
      };

      // Mock GraphQL service methods
      (
        mockGraphqlService.getInningsScoreSummary as jest.Mock
      ).mockResolvedValue(mockInningsResponse);
      (mockGraphqlService.getLatestOverOnly as jest.Mock).mockResolvedValue(
        mockOverResponse
      );

      // Test the saga execution - API call first
      const apiCallResult = generator.next().value;

      // The saga calls: apiService.addBall.bind(apiService), ballEvent
      // So we need to check if it's calling the bound method
      expect(apiCallResult.type).toBe('CALL');
      expect(apiCallResult.payload.args[0]).toEqual(mockBallEvent);

      // Ball success
      expect(generator.next().value).toEqual(put(addBallSuccess()));

      // Test sequential GraphQL calls
      const inningsCallResult = generator.next().value;

      // Check the GraphQL call structure
      expect(inningsCallResult.type).toBe('CALL');
      expect(inningsCallResult.payload.args).toEqual(['test-match-id', 1]);

      expect(generator.next(mockInningsResponse).value).toEqual(
        put(fetchInningsScoreSummarySuccess(mockInningsResponse.data))
      );

      const overCallResult = generator.next().value;

      // Check the over call structure
      expect(overCallResult.type).toBe('CALL');
      expect(overCallResult.payload.args).toEqual(['test-match-id', 1]);

      expect(generator.next(mockOverResponse).value).toEqual(
        put(
          fetchLatestOverSuccess({
            inningsNumber: 1,
            over: mockOverResponse.data,
          })
        )
      );

      // Should not fetch updated scorecard since innings is not completed
      expect(generator.next().done).toBe(true);
    });

    it('should fetch updated scorecard when innings is completed', () => {
      const generator = addBallSaga(addBallRequest(mockBallEvent));

      // Mock successful API call
      const mockApiInstance = {
        addBall: jest.fn().mockResolvedValue({ success: true }),
      };
      mockApiService.mockImplementation(() => mockApiInstance as any);

      // Mock completed innings response
      const mockInningsResponse = {
        success: true,
        data: {
          innings_number: 1,
          batting_team: 'A',
          total_runs: 120,
          total_wickets: 10,
          total_overs: 20,
          total_balls: 120,
          status: 'completed',
          extras: { total: 5 },
        },
      };

      const mockOverResponse = {
        success: true,
        data: {
          over_number: 20,
          total_runs: 6,
          total_balls: 6,
          total_wickets: 1,
          status: 'completed',
          balls: [],
        },
      };

      mockGraphqlService.getInningsScoreSummary.mockResolvedValue(
        mockInningsResponse
      );
      mockGraphqlService.getLatestOverOnly.mockResolvedValue(mockOverResponse);

      // Test the saga execution
      const apiCallResult = generator.next().value;
      expect(apiCallResult.type).toBe('CALL');
      expect(apiCallResult.payload.args[0]).toEqual(mockBallEvent);

      expect(generator.next().value).toEqual(put(addBallSuccess()));

      // Test sequential GraphQL calls
      const inningsCallResult = generator.next().value;
      expect(inningsCallResult.type).toBe('CALL');
      expect(inningsCallResult.payload.args).toEqual(['test-match-id', 1]);

      expect(generator.next(mockInningsResponse).value).toEqual(
        put(fetchInningsScoreSummarySuccess(mockInningsResponse.data))
      );

      const overCallResult = generator.next().value;
      expect(overCallResult.type).toBe('CALL');
      expect(overCallResult.payload.args).toEqual(['test-match-id', 1]);

      expect(generator.next(mockOverResponse).value).toEqual(
        put(
          fetchLatestOverSuccess({
            inningsNumber: 1,
            over: mockOverResponse.data,
          })
        )
      );

      // Should fetch updated scorecard since innings is completed
      expect(generator.next().value).toEqual(
        put(fetchUpdatedScorecardRequest('test-match-id'))
      );

      expect(generator.next().done).toBe(true);
    });

    it('should handle API failure and not proceed with GraphQL calls', () => {
      const generator = addBallSaga(addBallRequest(mockBallEvent));

      // Mock API failure
      const mockApiInstance = {
        addBall: jest.fn().mockRejectedValue(new ApiError('API Error')),
      };
      mockApiService.mockImplementation(() => mockApiInstance as any);

      // Test the saga execution
      const apiCallResult = generator.next().value;
      expect(apiCallResult.type).toBe('CALL');
      expect(apiCallResult.payload.args[0]).toEqual(mockBallEvent);

      // The saga should catch the error and dispatch failure
      const errorResult = generator.throw(new ApiError('API Error')).value;
      expect(errorResult).toEqual(put(addBallFailure('API Error')));

      expect(generator.next().done).toBe(true);
    });

    it('should handle GraphQL errors gracefully without failing ball addition', () => {
      const generator = addBallSaga(addBallRequest(mockBallEvent));

      // Mock successful API call
      const mockApiInstance = {
        addBall: jest.fn().mockResolvedValue({ success: true }),
      };
      mockApiService.mockImplementation(() => mockApiInstance as any);

      // Mock GraphQL failure
      mockGraphqlService.getInningsScoreSummary.mockRejectedValue(
        new Error('GraphQL Error')
      );

      // Test the saga execution
      const apiCallResult = generator.next().value;
      expect(apiCallResult.type).toBe('CALL');
      expect(apiCallResult.payload.args[0]).toEqual(mockBallEvent);

      expect(generator.next().value).toEqual(put(addBallSuccess()));

      // Should handle GraphQL error gracefully
      const inningsCallResult = generator.next().value;
      expect(inningsCallResult.type).toBe('CALL');
      expect(inningsCallResult.payload.args).toEqual(['test-match-id', 1]);

      // Should catch the error and continue
      expect(generator.next().done).toBe(true);
    });

    it('should handle partial GraphQL success (innings success, over failure)', () => {
      const generator = addBallSaga(addBallRequest(mockBallEvent));

      // Mock successful API call
      const mockApiInstance = {
        addBall: jest.fn().mockResolvedValue({ success: true }),
      };
      mockApiService.mockImplementation(() => mockApiInstance as any);

      // Mock partial GraphQL responses
      const mockInningsResponse = {
        success: true,
        data: {
          innings_number: 1,
          batting_team: 'A',
          total_runs: 4,
          total_wickets: 0,
          total_overs: 1,
          total_balls: 1,
          status: 'in_progress',
          extras: { total: 0 },
        },
      };

      mockGraphqlService.getInningsScoreSummary.mockResolvedValue(
        mockInningsResponse
      );
      mockGraphqlService.getLatestOverOnly.mockRejectedValue(
        new Error('Over fetch failed')
      );

      // Test the saga execution
      const apiCallResult = generator.next().value;
      expect(apiCallResult.type).toBe('CALL');
      expect(apiCallResult.payload.args[0]).toEqual(mockBallEvent);

      expect(generator.next().value).toEqual(put(addBallSuccess()));

      // Test sequential GraphQL calls
      const inningsCallResult = generator.next().value;
      expect(inningsCallResult.type).toBe('CALL');
      expect(inningsCallResult.payload.args).toEqual(['test-match-id', 1]);

      expect(generator.next(mockInningsResponse).value).toEqual(
        put(fetchInningsScoreSummarySuccess(mockInningsResponse.data))
      );

      const overCallResult = generator.next().value;
      expect(overCallResult.type).toBe('CALL');
      expect(overCallResult.payload.args).toEqual(['test-match-id', 1]);

      // Should handle over fetch error gracefully
      expect(generator.next().done).toBe(true);
    });
  });

  describe('fetchInningsScoreSummarySaga', () => {
    it('should successfully fetch innings score summary', () => {
      const generator = fetchInningsScoreSummarySaga({
        type: 'scorecard/fetchInningsScoreSummaryRequest',
        payload: { matchId: 'test-match', inningsNumber: 1 },
      });

      const mockResponse = {
        success: true,
        data: {
          innings_number: 1,
          batting_team: 'A',
          total_runs: 50,
          total_wickets: 2,
          total_overs: 10,
          total_balls: 60,
          status: 'in_progress',
          extras: { total: 3 },
        },
      };

      mockGraphqlService.getInningsScoreSummary.mockResolvedValue(mockResponse);

      const callResult = generator.next().value;
      expect(callResult.type).toBe('CALL');
      expect(callResult.payload.args).toEqual(['test-match', 1]);

      expect(generator.next(mockResponse).value).toEqual(
        put(fetchInningsScoreSummarySuccess(mockResponse.data))
      );

      expect(generator.next().done).toBe(true);
    });

    it('should handle GraphQL failure', () => {
      const generator = fetchInningsScoreSummarySaga({
        type: 'scorecard/fetchInningsScoreSummaryRequest',
        payload: { matchId: 'test-match', inningsNumber: 1 },
      });

      const mockResponse = {
        success: false,
        error: 'Innings not found',
      };

      mockGraphqlService.getInningsScoreSummary.mockResolvedValue(mockResponse);

      const callResult = generator.next().value;
      expect(callResult.type).toBe('CALL');
      expect(callResult.payload.args).toEqual(['test-match', 1]);

      expect(generator.next(mockResponse).value).toEqual(
        put({
          type: 'scorecard/fetchInningsScoreSummaryFailure',
          payload: 'Innings not found',
        })
      );

      expect(generator.next().done).toBe(true);
    });
  });

  describe('fetchLatestOverSaga', () => {
    it('should successfully fetch latest over', () => {
      const generator = fetchLatestOverSaga({
        type: 'scorecard/fetchLatestOverRequest',
        payload: { matchId: 'test-match', inningsNumber: 1 },
      });

      const mockResponse = {
        success: true,
        data: {
          over_number: 5,
          total_runs: 8,
          total_balls: 6,
          total_wickets: 0,
          status: 'completed',
          balls: [
            {
              ball_number: 1,
              ball_type: 'good',
              run_type: 'runs',
              runs: 2,
              byes: 0,
              is_wicket: false,
              wicket_type: null,
            },
          ],
        },
      };

      mockGraphqlService.getLatestOverOnly.mockResolvedValue(mockResponse);

      const callResult = generator.next().value;
      expect(callResult.type).toBe('CALL');
      expect(callResult.payload.args).toEqual(['test-match', 1]);

      expect(generator.next(mockResponse).value).toEqual(
        put(
          fetchLatestOverSuccess({
            inningsNumber: 1,
            over: mockResponse.data,
          })
        )
      );

      expect(generator.next().done).toBe(true);
    });

    it('should handle GraphQL failure', () => {
      const generator = fetchLatestOverSaga({
        type: 'scorecard/fetchLatestOverRequest',
        payload: { matchId: 'test-match', inningsNumber: 1 },
      });

      const mockResponse = {
        success: false,
        error: 'No over found',
      };

      mockGraphqlService.getLatestOverOnly.mockResolvedValue(mockResponse);

      const callResult = generator.next().value;
      expect(callResult.type).toBe('CALL');
      expect(callResult.payload.args).toEqual(['test-match', 1]);

      expect(generator.next(mockResponse).value).toEqual(
        put({
          type: 'scorecard/fetchLatestOverFailure',
          payload: 'No over found',
        })
      );

      expect(generator.next().done).toBe(true);
    });
  });
});
