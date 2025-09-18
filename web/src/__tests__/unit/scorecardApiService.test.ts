import { ApiService, ApiError } from '../../services/api';
import {
  ScorecardResponse,
  BallEventRequest,
  OverSummary,
  InningsSummary,
} from '../../store/reducers/scorecardSlice';

// Mock fetch globally
const mockFetch = jest.fn();
global.fetch = mockFetch;

// Mock console methods to avoid noise in tests
const originalConsoleError = console.error;
const originalConsoleLog = console.log;

beforeAll(() => {
  console.error = jest.fn();
  console.log = jest.fn();
});

afterAll(() => {
  console.error = originalConsoleError;
  console.log = originalConsoleLog;
});

describe('ApiService - Scorecard Endpoints', () => {
  let apiService: ApiService;

  beforeEach(() => {
    apiService = new ApiService();
    mockFetch.mockClear();
  });

  describe('getScorecard', () => {
    const mockScorecardResponse: ScorecardResponse = {
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
      const mockResponse = {
        ok: true,
        status: 200,
        json: async () => ({ data: mockScorecardResponse }),
      };

      mockFetch.mockResolvedValueOnce(mockResponse);

      const result = await apiService.getScorecard('match-1');

      expect(mockFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/scorecard/match-1?_t=',
        expect.objectContaining({
          method: 'GET',
          headers: expect.objectContaining({
            'Content-Type': 'application/json',
          }),
        })
      );

      expect(result.data).toEqual(mockScorecardResponse);
    });

    it('should handle API errors', async () => {
      const mockResponse = {
        ok: false,
        status: 404,
        json: async () => ({ error: 'Scorecard not found' }),
      };

      mockFetch.mockResolvedValueOnce(mockResponse);

      await expect(
        apiService.getScorecard('nonexistent-match')
      ).rejects.toThrow(ApiError);
    });

    it('should handle network errors', async () => {
      mockFetch.mockRejectedValueOnce(new Error('Network error'));

      await expect(apiService.getScorecard('match-1')).rejects.toThrow(
        'Network error'
      );
    });
  });

  describe('startScoring', () => {
    it('should start scoring successfully', async () => {
      const mockResponse = {
        ok: true,
        status: 200,
        json: async () => ({
          message: 'Scoring started successfully',
          match_id: 'match-1',
        }),
      };

      mockFetch.mockResolvedValueOnce(mockResponse);

      const result = await apiService.startScoring('match-1');

      expect(mockFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/scorecard/start?_t=',
        expect.objectContaining({
          method: 'POST',
          headers: expect.objectContaining({
            'Content-Type': 'application/json',
          }),
          body: JSON.stringify({ match_id: 'match-1' }),
        })
      );

      expect(result.data.message).toBe('Scoring started successfully');
      expect(result.data.match_id).toBe('match-1');
    });

    it('should handle start scoring errors', async () => {
      const mockResponse = {
        ok: false,
        status: 400,
        json: async () => ({ error: 'Match is not ready for scoring' }),
      };

      mockFetch.mockResolvedValueOnce(mockResponse);

      await expect(apiService.startScoring('match-1')).rejects.toThrow(
        ApiError
      );
    });
  });

  describe('addBall', () => {
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
      const mockResponse = {
        ok: true,
        status: 200,
        json: async () => ({
          message: 'Ball added successfully',
          match_id: 'match-1',
          innings_number: 1,
          ball_type: 'good',
          run_type: '4',
          runs: 4,
          byes: 0,
          is_wicket: false,
        }),
      };

      mockFetch.mockResolvedValueOnce(mockResponse);

      const result = await apiService.addBall(mockBallEvent);

      expect(mockFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/scorecard/ball?_t=',
        expect.objectContaining({
          method: 'POST',
          headers: expect.objectContaining({
            'Content-Type': 'application/json',
          }),
          body: JSON.stringify(mockBallEvent),
        })
      );

      expect(result.data.message).toBe('Ball added successfully');
      expect(result.data.match_id).toBe('match-1');
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

      const mockResponse = {
        ok: true,
        status: 200,
        json: async () => ({
          message: 'Wicket ball added successfully',
          match_id: 'match-1',
          innings_number: 1,
          ball_type: 'good',
          run_type: 'WC',
          runs: 0,
          byes: 0,
          is_wicket: true,
          wicket_type: 'bowled',
        }),
      };

      mockFetch.mockResolvedValueOnce(mockResponse);

      const result = await apiService.addBall(wicketBallEvent);

      expect(result.data.is_wicket).toBe(true);
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

      const mockResponse = {
        ok: true,
        status: 200,
        json: async () => ({
          message: 'Wide ball added successfully',
          match_id: 'match-1',
          innings_number: 1,
          ball_type: 'wide',
          run_type: 'WD',
          runs: 1,
          byes: 0,
          is_wicket: false,
        }),
      };

      mockFetch.mockResolvedValueOnce(mockResponse);

      const result = await apiService.addBall(wideBallEvent);

      expect(result.data.ball_type).toBe('wide');
      expect(result.data.run_type).toBe('WD');
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

      const mockResponse = {
        ok: true,
        status: 200,
        json: async () => ({
          message: 'No ball added successfully',
          match_id: 'match-1',
          innings_number: 1,
          ball_type: 'no_ball',
          run_type: 'NB',
          runs: 1,
          byes: 0,
          is_wicket: false,
        }),
      };

      mockFetch.mockResolvedValueOnce(mockResponse);

      const result = await apiService.addBall(noBallEvent);

      expect(result.data.ball_type).toBe('no_ball');
      expect(result.data.run_type).toBe('NB');
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

      const mockResponse = {
        ok: true,
        status: 200,
        json: async () => ({
          message: 'Ball with byes added successfully',
          match_id: 'match-1',
          innings_number: 1,
          ball_type: 'good',
          run_type: '1',
          runs: 1,
          byes: 2,
          is_wicket: false,
        }),
      };

      mockFetch.mockResolvedValueOnce(mockResponse);

      const result = await apiService.addBall(ballWithByes);

      expect(result.data.byes).toBe(2);
    });

    it('should handle add ball errors', async () => {
      const mockResponse = {
        ok: false,
        status: 400,
        json: async () => ({ error: 'Invalid ball data' }),
      };

      mockFetch.mockResolvedValueOnce(mockResponse);

      await expect(apiService.addBall(mockBallEvent)).rejects.toThrow(ApiError);
    });

    it('should handle innings completion error', async () => {
      const mockResponse = {
        ok: false,
        status: 400,
        json: async () => ({ error: 'Innings already completed' }),
      };

      mockFetch.mockResolvedValueOnce(mockResponse);

      await expect(apiService.addBall(mockBallEvent)).rejects.toThrow(ApiError);
    });
  });

  describe('getCurrentOver', () => {
    const mockOverResponse: OverSummary = {
      over_number: 5,
      total_runs: 8,
      total_balls: 6,
      total_wickets: 0,
      status: 'in_progress',
      balls: [
        {
          ball_number: 1,
          ball_type: 'good',
          run_type: '1',
          runs: 1,
          byes: 0,
          is_wicket: false,
        },
        {
          ball_number: 2,
          ball_type: 'good',
          run_type: '4',
          runs: 4,
          byes: 0,
          is_wicket: false,
        },
        {
          ball_number: 3,
          ball_type: 'wide',
          run_type: 'WD',
          runs: 1,
          byes: 0,
          is_wicket: false,
        },
      ],
    };

    it('should fetch current over successfully', async () => {
      const mockResponse = {
        ok: true,
        status: 200,
        json: async () => ({ data: mockOverResponse }),
      };

      mockFetch.mockResolvedValueOnce(mockResponse);

      const result = await apiService.getCurrentOver('match-1', 1);

      expect(mockFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/scorecard/match-1/current-over?innings=1&_t=',
        expect.objectContaining({
          method: 'GET',
          headers: expect.objectContaining({
            'Content-Type': 'application/json',
          }),
        })
      );

      expect(result.data).toEqual(mockOverResponse);
    });

    it('should fetch current over with default innings', async () => {
      const mockResponse = {
        ok: true,
        status: 200,
        json: async () => ({ data: mockOverResponse }),
      };

      mockFetch.mockResolvedValueOnce(mockResponse);

      await apiService.getCurrentOver('match-1');

      expect(mockFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/scorecard/match-1/current-over?innings=1&_t=',
        expect.any(Object)
      );
    });

    it('should handle get current over errors', async () => {
      const mockResponse = {
        ok: false,
        status: 404,
        json: async () => ({ error: 'Current over not found' }),
      };

      mockFetch.mockResolvedValueOnce(mockResponse);

      await expect(apiService.getCurrentOver('match-1', 1)).rejects.toThrow(
        ApiError
      );
    });
  });

  describe('getInnings', () => {
    const mockInningsResponse: InningsSummary = {
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
      const mockResponse = {
        ok: true,
        status: 200,
        json: async () => ({ data: mockInningsResponse }),
      };

      mockFetch.mockResolvedValueOnce(mockResponse);

      const result = await apiService.getInnings('match-1', 1);

      expect(mockFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/scorecard/match-1/innings/1?_t=',
        expect.objectContaining({
          method: 'GET',
          headers: expect.objectContaining({
            'Content-Type': 'application/json',
          }),
        })
      );

      expect(result.data).toEqual(mockInningsResponse);
    });

    it('should handle get innings errors', async () => {
      const mockResponse = {
        ok: false,
        status: 404,
        json: async () => ({ error: 'Innings not found' }),
      };

      mockFetch.mockResolvedValueOnce(mockResponse);

      await expect(apiService.getInnings('match-1', 1)).rejects.toThrow(
        ApiError
      );
    });
  });

  describe('getOver', () => {
    const mockOverResponse: OverSummary = {
      over_number: 3,
      total_runs: 12,
      total_balls: 6,
      total_wickets: 1,
      status: 'completed',
      balls: [
        {
          ball_number: 1,
          ball_type: 'good',
          run_type: '4',
          runs: 4,
          byes: 0,
          is_wicket: false,
        },
        {
          ball_number: 2,
          ball_type: 'good',
          run_type: '2',
          runs: 2,
          byes: 0,
          is_wicket: false,
        },
        {
          ball_number: 3,
          ball_type: 'good',
          run_type: 'WC',
          runs: 0,
          byes: 0,
          is_wicket: true,
          wicket_type: 'bowled',
        },
        {
          ball_number: 4,
          ball_type: 'good',
          run_type: '6',
          runs: 6,
          byes: 0,
          is_wicket: false,
        },
      ],
    };

    it('should fetch specific over successfully', async () => {
      const mockResponse = {
        ok: true,
        status: 200,
        json: async () => ({ data: mockOverResponse }),
      };

      mockFetch.mockResolvedValueOnce(mockResponse);

      const result = await apiService.getOver('match-1', 1, 3);

      expect(mockFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/scorecard/match-1/innings/1/over/3?_t=',
        expect.objectContaining({
          method: 'GET',
          headers: expect.objectContaining({
            'Content-Type': 'application/json',
          }),
        })
      );

      expect(result.data).toEqual(mockOverResponse);
    });

    it('should handle get over errors', async () => {
      const mockResponse = {
        ok: false,
        status: 404,
        json: async () => ({ error: 'Over not found' }),
      };

      mockFetch.mockResolvedValueOnce(mockResponse);

      await expect(apiService.getOver('match-1', 1, 3)).rejects.toThrow(
        ApiError
      );
    });
  });

  describe('Network Error Handling', () => {
    it('should retry on network errors', async () => {
      // First call fails, second succeeds
      mockFetch
        .mockRejectedValueOnce(new Error('Network error'))
        .mockResolvedValueOnce({
          ok: true,
          status: 200,
          json: async () => ({ data: { match_id: 'match-1' } }),
        });

      const result = await apiService.getScorecard('match-1');

      expect(mockFetch).toHaveBeenCalledTimes(2);
      expect(result.data.match_id).toBe('match-1');
    }, 10000);

    it('should throw ApiError after max retries', async () => {
      mockFetch.mockRejectedValue(new Error('Network error'));

      await expect(apiService.getScorecard('match-1')).rejects.toThrow(
        'Network error'
      );
    }, 10000);
  });

  describe('Request Configuration', () => {
    it('should include proper headers', async () => {
      const mockResponse = {
        ok: true,
        status: 200,
        json: async () => ({ data: { match_id: 'match-1' } }),
      };

      mockFetch.mockResolvedValueOnce(mockResponse);

      await apiService.getScorecard('match-1');

      expect(mockFetch).toHaveBeenCalledWith(
        expect.any(String),
        expect.objectContaining({
          headers: expect.objectContaining({
            'Content-Type': 'application/json',
            'Cache-Control': 'no-cache',
            Pragma: 'no-cache',
            Expires: '0',
          }),
          mode: 'cors',
          credentials: 'include',
        })
      );
    });

    it('should add cache-busting parameter', async () => {
      const mockResponse = {
        ok: true,
        status: 200,
        json: async () => ({ data: { match_id: 'match-1' } }),
      };

      mockFetch.mockResolvedValueOnce(mockResponse);

      await apiService.getScorecard('match-1');

      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringMatching(/\?_t=/),
        expect.any(Object)
      );
    });
  });
});
