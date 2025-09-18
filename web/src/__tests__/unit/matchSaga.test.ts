/* eslint-disable @typescript-eslint/no-explicit-any */
import { runSaga } from 'redux-saga';
import {
  fetchMatchesSaga,
  createMatchSaga,
  updateMatchSaga,
  deleteMatchSaga,
} from '@/store/sagas/matchSaga';
import {
  // fetchMatchesRequest,
  fetchMatchesSuccess,
  fetchMatchesFailure,
  createMatchRequest,
  createMatchSuccess,
  createMatchFailure,
  updateMatchRequest,
  updateMatchSuccess,
  updateMatchFailure,
  deleteMatchRequest,
  deleteMatchSuccess,
  deleteMatchFailure,
} from '@/store/reducers/matchSlice';
import { Match } from '@/store/reducers/matchSlice';
import { ApiService, ApiError } from '../../services/api';

// Mock the API service
jest.mock('../../services/api', () => ({
  ApiService: jest.fn().mockImplementation(() => ({
    getMatches: jest.fn(),
    createMatch: jest.fn(),
    updateMatch: jest.fn(),
    deleteMatch: jest.fn(),
  })),
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

describe('Match Saga Tests', () => {
  let dispatched: any[];
  let mockApiService: any;

  beforeEach(() => {
    dispatched = [];
    mockApiService = {
      getMatches: jest.fn(),
      createMatch: jest.fn(),
      updateMatch: jest.fn(),
      deleteMatch: jest.fn(),
    };

    // Mock ApiService constructor to return our mock
    (ApiService as jest.Mock).mockImplementation(() => mockApiService);
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  describe('fetchMatchesSaga', () => {
    it('should fetch matches successfully', async () => {
      const mockMatches: Match[] = [
        {
          id: '1',
          series_id: 'series-1',
          match_number: 1,
          date: '2024-01-01T00:00:00Z',
          status: 'live',
          team_a_player_count: 11,
          team_b_player_count: 11,
          total_overs: 20,
          toss_winner: 'A',
          toss_type: 'H',
          batting_team: 'A',
          created_at: '2024-01-01T00:00:00Z',
          updated_at: '2024-01-01T00:00:00Z',
        },
      ];

      const mockResponse = {
        data: { data: mockMatches },
        success: true,
      };

      mockApiService.getMatches.mockResolvedValue(mockResponse);

      await runSaga(
        {
          dispatch: action => dispatched.push(action),
        },
        fetchMatchesSaga
      ).toPromise();

      expect(mockApiService.getMatches).toHaveBeenCalledTimes(1);
      expect(dispatched).toContainEqual(fetchMatchesSuccess(mockMatches));
    });

    it('should handle nested response structure', async () => {
      const mockMatches: Match[] = [
        {
          id: '1',
          series_id: 'series-1',
          match_number: 1,
          date: '2024-01-01T00:00:00Z',
          status: 'live',
          team_a_player_count: 11,
          team_b_player_count: 11,
          total_overs: 20,
          toss_winner: 'A',
          toss_type: 'H',
          batting_team: 'A',
          created_at: '2024-01-01T00:00:00Z',
          updated_at: '2024-01-01T00:00:00Z',
        },
      ];

      const mockResponse = {
        data: mockMatches, // Direct array without nested data
        success: true,
      };

      mockApiService.getMatches.mockResolvedValue(mockResponse);

      await runSaga(
        {
          dispatch: action => dispatched.push(action),
        },
        fetchMatchesSaga
      ).toPromise();

      expect(dispatched).toContainEqual(fetchMatchesSuccess(mockMatches));
    });

    it('should handle API errors', async () => {
      const error = new ApiError('Network error', 500);
      mockApiService.getMatches.mockRejectedValue(error);

      await runSaga(
        {
          dispatch: action => dispatched.push(action),
        },
        fetchMatchesSaga
      ).toPromise();

      expect(dispatched).toContainEqual(fetchMatchesFailure('Network error'));
    });

    it('should handle generic errors', async () => {
      const error = new Error('Generic error');
      mockApiService.getMatches.mockRejectedValue(error);

      await runSaga(
        {
          dispatch: action => dispatched.push(action),
        },
        fetchMatchesSaga
      ).toPromise();

      expect(dispatched).toContainEqual(
        fetchMatchesFailure('Failed to fetch matches')
      );
    });
  });

  describe('createMatchSaga', () => {
    it('should create match successfully', async () => {
      const matchData: Omit<Match, 'id' | 'created_at' | 'updated_at'> = {
        series_id: 'series-1',
        match_number: 1,
        date: '2024-01-01T00:00:00Z',
        status: 'live',
        team_a_player_count: 11,
        team_b_player_count: 11,
        total_overs: 20,
        toss_winner: 'A',
        toss_type: 'H',
        batting_team: 'A',
      };

      const mockCreatedMatch: Match = {
        ...matchData,
        id: '1',
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z',
      };

      const mockResponse = {
        data: { data: mockCreatedMatch },
        success: true,
      };

      mockApiService.createMatch.mockResolvedValue(mockResponse);

      const action = createMatchRequest(matchData);

      await runSaga(
        {
          dispatch: action => dispatched.push(action),
        },
        createMatchSaga,
        action
      ).toPromise();

      expect(mockApiService.createMatch).toHaveBeenCalledWith(matchData);
      expect(dispatched).toContainEqual(createMatchSuccess(mockCreatedMatch));
    });

    it('should handle nested response structure', async () => {
      const matchData: Omit<Match, 'id' | 'created_at' | 'updated_at'> = {
        series_id: 'series-1',
        match_number: 1,
        date: '2024-01-01T00:00:00Z',
        status: 'live',
        team_a_player_count: 11,
        team_b_player_count: 11,
        total_overs: 20,
        toss_winner: 'A',
        toss_type: 'H',
        batting_team: 'A',
      };

      const mockCreatedMatch: Match = {
        ...matchData,
        id: '1',
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z',
      };

      const mockResponse = {
        data: mockCreatedMatch, // Direct object without nested data
        success: true,
      };

      mockApiService.createMatch.mockResolvedValue(mockResponse);

      const action = createMatchRequest(matchData);

      await runSaga(
        {
          dispatch: action => dispatched.push(action),
        },
        createMatchSaga,
        action
      ).toPromise();

      expect(dispatched).toContainEqual(createMatchSuccess(mockCreatedMatch));
    });

    it('should handle API errors', async () => {
      const matchData: Omit<Match, 'id' | 'created_at' | 'updated_at'> = {
        series_id: 'series-1',
        match_number: 1,
        date: '2024-01-01T00:00:00Z',
        status: 'live',
        team_a_player_count: 11,
        team_b_player_count: 11,
        total_overs: 20,
        toss_winner: 'A',
        toss_type: 'H',
        batting_team: 'A',
      };

      const error = new ApiError('Validation failed', 400);
      mockApiService.createMatch.mockRejectedValue(error);

      const action = createMatchRequest(matchData);

      await runSaga(
        {
          dispatch: action => dispatched.push(action),
        },
        createMatchSaga,
        action
      ).toPromise();

      expect(dispatched).toContainEqual(
        createMatchFailure('Validation failed')
      );
    });

    it('should handle generic errors', async () => {
      const matchData: Omit<Match, 'id' | 'created_at' | 'updated_at'> = {
        series_id: 'series-1',
        match_number: 1,
        date: '2024-01-01T00:00:00Z',
        status: 'live',
        team_a_player_count: 11,
        team_b_player_count: 11,
        total_overs: 20,
        toss_winner: 'A',
        toss_type: 'H',
        batting_team: 'A',
      };

      const error = new Error('Network error');
      mockApiService.createMatch.mockRejectedValue(error);

      const action = createMatchRequest(matchData);

      await runSaga(
        {
          dispatch: action => dispatched.push(action),
        },
        createMatchSaga,
        action
      ).toPromise();

      expect(dispatched).toContainEqual(
        createMatchFailure('Failed to create match')
      );
    });
  });

  describe('updateMatchSaga', () => {
    it('should update match successfully', async () => {
      const matchId = '1';
      const matchData: Omit<Match, 'id' | 'created_at' | 'updated_at'> = {
        series_id: 'series-1',
        match_number: 1,
        date: '2024-01-01T00:00:00Z',
        status: 'live',
        team_a_player_count: 10,
        team_b_player_count: 10,
        total_overs: 15,
        toss_winner: 'A',
        toss_type: 'H',
        batting_team: 'A',
      };

      const mockUpdatedMatch: Match = {
        ...matchData,
        id: matchId,
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z',
      };

      const mockResponse = {
        data: { data: mockUpdatedMatch },
        success: true,
      };

      mockApiService.updateMatch.mockResolvedValue(mockResponse);

      const action = updateMatchRequest({ id: matchId, matchData });

      await runSaga(
        {
          dispatch: action => dispatched.push(action),
        },
        updateMatchSaga,
        action
      ).toPromise();

      expect(mockApiService.updateMatch).toHaveBeenCalledWith(
        matchId,
        matchData
      );
      expect(dispatched).toContainEqual(updateMatchSuccess(mockUpdatedMatch));
    });

    it('should handle nested response structure', async () => {
      const matchId = '1';
      const matchData: Omit<Match, 'id' | 'created_at' | 'updated_at'> = {
        series_id: 'series-1',
        match_number: 1,
        date: '2024-01-01T00:00:00Z',
        status: 'live',
        team_a_player_count: 10,
        team_b_player_count: 10,
        total_overs: 15,
        toss_winner: 'A',
        toss_type: 'H',
        batting_team: 'A',
      };

      const mockUpdatedMatch: Match = {
        ...matchData,
        id: matchId,
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z',
      };

      const mockResponse = {
        data: mockUpdatedMatch, // Direct object without nested data
        success: true,
      };

      mockApiService.updateMatch.mockResolvedValue(mockResponse);

      const action = updateMatchRequest({ id: matchId, matchData });

      await runSaga(
        {
          dispatch: action => dispatched.push(action),
        },
        updateMatchSaga,
        action
      ).toPromise();

      expect(dispatched).toContainEqual(updateMatchSuccess(mockUpdatedMatch));
    });

    it('should handle API errors', async () => {
      const matchId = '999';
      const matchData: Omit<Match, 'id' | 'created_at' | 'updated_at'> = {
        series_id: 'series-1',
        match_number: 1,
        date: '2024-01-01T00:00:00Z',
        status: 'live',
        team_a_player_count: 10,
        team_b_player_count: 10,
        total_overs: 15,
        toss_winner: 'A',
        toss_type: 'H',
        batting_team: 'A',
      };

      const error = new ApiError('Match not found', 404);
      mockApiService.updateMatch.mockRejectedValue(error);

      const action = updateMatchRequest({ id: matchId, matchData });

      await runSaga(
        {
          dispatch: action => dispatched.push(action),
        },
        updateMatchSaga,
        action
      ).toPromise();

      expect(dispatched).toContainEqual(updateMatchFailure('Match not found'));
    });

    it('should handle generic errors', async () => {
      const matchId = '1';
      const matchData: Omit<Match, 'id' | 'created_at' | 'updated_at'> = {
        series_id: 'series-1',
        match_number: 1,
        date: '2024-01-01T00:00:00Z',
        status: 'live',
        team_a_player_count: 10,
        team_b_player_count: 10,
        total_overs: 15,
        toss_winner: 'A',
        toss_type: 'H',
        batting_team: 'A',
      };

      const error = new Error('Network error');
      mockApiService.updateMatch.mockRejectedValue(error);

      const action = updateMatchRequest({ id: matchId, matchData });

      await runSaga(
        {
          dispatch: action => dispatched.push(action),
        },
        updateMatchSaga,
        action
      ).toPromise();

      expect(dispatched).toContainEqual(
        updateMatchFailure('Failed to update match')
      );
    });
  });

  describe('deleteMatchSaga', () => {
    it('should delete match successfully', async () => {
      const matchId = '1';

      const mockResponse = {
        data: { data: undefined },
        success: true,
      };

      mockApiService.deleteMatch.mockResolvedValue(mockResponse);

      const action = deleteMatchRequest(matchId);

      await runSaga(
        {
          dispatch: action => dispatched.push(action),
        },
        deleteMatchSaga,
        action
      ).toPromise();

      expect(mockApiService.deleteMatch).toHaveBeenCalledWith(matchId);
      expect(dispatched).toContainEqual(deleteMatchSuccess(matchId));
    });

    it('should handle nested response structure', async () => {
      const matchId = '1';

      const mockResponse = {
        data: undefined, // Direct undefined without nested data
        success: true,
      };

      mockApiService.deleteMatch.mockResolvedValue(mockResponse);

      const action = deleteMatchRequest(matchId);

      await runSaga(
        {
          dispatch: action => dispatched.push(action),
        },
        deleteMatchSaga,
        action
      ).toPromise();

      expect(dispatched).toContainEqual(deleteMatchSuccess(matchId));
    });

    it('should handle API errors', async () => {
      const matchId = '999';

      const error = new ApiError('Match not found', 404);
      mockApiService.deleteMatch.mockRejectedValue(error);

      const action = deleteMatchRequest(matchId);

      await runSaga(
        {
          dispatch: action => dispatched.push(action),
        },
        deleteMatchSaga,
        action
      ).toPromise();

      expect(dispatched).toContainEqual(deleteMatchFailure('Match not found'));
    });

    it('should handle generic errors', async () => {
      const matchId = '1';

      const error = new Error('Network error');
      mockApiService.deleteMatch.mockRejectedValue(error);

      const action = deleteMatchRequest(matchId);

      await runSaga(
        {
          dispatch: action => dispatched.push(action),
        },
        deleteMatchSaga,
        action
      ).toPromise();

      expect(dispatched).toContainEqual(
        deleteMatchFailure('Failed to delete match')
      );
    });
  });

  describe('Saga Error Handling', () => {
    it('should handle ApiError with status and details', async () => {
      const error = new ApiError('Custom error', 422, { field: 'validation' });
      mockApiService.getMatches.mockRejectedValue(error);

      await runSaga(
        {
          dispatch: action => dispatched.push(action),
        },
        fetchMatchesSaga
      ).toPromise();

      expect(dispatched).toContainEqual(fetchMatchesFailure('Custom error'));
    });

    it('should handle non-Error objects', async () => {
      const error = 'String error';
      mockApiService.getMatches.mockRejectedValue(error);

      await runSaga(
        {
          dispatch: action => dispatched.push(action),
        },
        fetchMatchesSaga
      ).toPromise();

      expect(dispatched).toContainEqual(
        fetchMatchesFailure('Failed to fetch matches')
      );
    });

    it('should handle null/undefined errors', async () => {
      mockApiService.getMatches.mockRejectedValue(null);

      await runSaga(
        {
          dispatch: action => dispatched.push(action),
        },
        fetchMatchesSaga
      ).toPromise();

      expect(dispatched).toContainEqual(
        fetchMatchesFailure('Failed to fetch matches')
      );
    });
  });
});
