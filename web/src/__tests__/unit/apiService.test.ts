/* eslint-disable @typescript-eslint/no-explicit-any */
import { ApiService, ApiError } from '@/services/api';
import { Match } from '@/store/reducers/matchSlice';

// Mock fetch globally
global.fetch = jest.fn();

describe('ApiService - Match Endpoints', () => {
  let apiService: ApiService;
  const mockFetch = fetch as jest.MockedFunction<typeof fetch>;

  beforeEach(() => {
    apiService = new ApiService('http://localhost:8080/api/v1');
    mockFetch.mockClear();
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  describe('getMatches', () => {
    it('should fetch all matches successfully', async () => {
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
        ok: true,
        json: jest.fn().mockResolvedValue(mockMatches),
      };

      mockFetch.mockResolvedValue(mockResponse as any);

      const result = await apiService.getMatches();

      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/matches?_t='),
        expect.objectContaining({
          headers: expect.objectContaining({
            'Content-Type': 'application/json',
            'Accept': 'application/json',
          }),
          mode: 'cors',
          credentials: 'omit',
        })
      );

      expect(result).toEqual({
        data: mockMatches,
        success: true,
        message: undefined,
      });
    });

    it('should handle API errors', async () => {
      const mockResponse = {
        ok: false,
        status: 500,
        json: jest.fn().mockResolvedValue({ message: 'Internal server error' }),
      };

      mockFetch.mockResolvedValue(mockResponse as any);

      await expect(apiService.getMatches()).rejects.toThrow(ApiError);
    });

    it('should retry on 503 errors', async () => {
      const mockMatches: Match[] = [];

      const mockResponse = {
        ok: true,
        json: jest.fn().mockResolvedValue(mockMatches),
      };

      // First call returns 503, second call succeeds
      mockFetch
        .mockResolvedValueOnce({
          ok: false,
          status: 503,
          json: jest.fn().mockResolvedValue({ message: 'Service unavailable' }),
        } as any)
        .mockResolvedValueOnce(mockResponse as any);

      const result = await apiService.getMatches();

      expect(mockFetch).toHaveBeenCalledTimes(2);
      expect(result).toEqual({
        data: mockMatches,
        success: true,
        message: undefined,
      });
    });
  });

  describe('getMatchById', () => {
    it('should fetch a specific match by ID', async () => {
      const mockMatch: Match = {
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
      };

      const mockResponse = {
        ok: true,
        json: jest.fn().mockResolvedValue(mockMatch),
      };

      mockFetch.mockResolvedValue(mockResponse as any);

      const result = await apiService.getMatchById('1');

      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/matches/1?_t='),
        expect.objectContaining({
          headers: expect.objectContaining({
            'Content-Type': 'application/json',
            'Accept': 'application/json',
          }),
          mode: 'cors',
          credentials: 'omit',
        })
      );

      expect(result).toEqual({
        data: mockMatch,
        success: true,
        message: undefined,
      });
    });

    it('should handle 404 errors', async () => {
      const mockResponse = {
        ok: false,
        status: 404,
        json: jest.fn().mockResolvedValue({ message: 'Match not found' }),
      };

      mockFetch.mockResolvedValue(mockResponse as any);

      await expect(apiService.getMatchById('999')).rejects.toThrow(ApiError);
    });
  });

  describe('createMatch', () => {
    it('should create a new match successfully', async () => {
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
        ok: true,
        status: 201,
        json: jest.fn().mockResolvedValue(mockCreatedMatch),
      };

      mockFetch.mockResolvedValue(mockResponse as any);

      const result = await apiService.createMatch(matchData);

      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/matches?_t='),
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify({
            series_id: 'series-1',
            date: '2024-01-01T00:00:00Z',
            status: 'live',
            team_a_player_count: 11,
            team_b_player_count: 11,
            total_overs: 20,
            toss_winner: 'A',
            toss_type: 'H',
            batting_team: 'A',
          }),
          headers: expect.objectContaining({
            'Content-Type': 'application/json',
          }),
          mode: 'cors',
          credentials: 'omit',
        })
      );

      expect(result).toEqual({
        data: mockCreatedMatch,
        success: true,
        message: undefined,
      });
    });

    it('should remove match_number when it is 1 for auto-generation', async () => {
      const matchData: Omit<Match, 'id' | 'created_at' | 'updated_at'> = {
        series_id: 'series-1',
        match_number: 1, // Should be removed for auto-generation
        date: '2024-01-01T00:00:00Z',
        status: 'live',
        team_a_player_count: 11,
        team_b_player_count: 11,
        total_overs: 20,
        toss_winner: 'A',
        toss_type: 'H',
        batting_team: 'A',
      };

      const mockResponse = {
        ok: true,
        status: 201,
        json: jest.fn().mockResolvedValue({}),
      };

      mockFetch.mockResolvedValue(mockResponse as any);

      await apiService.createMatch(matchData);

      const expectedData = { ...matchData };
      delete (expectedData as any).match_number;

      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/matches?_t='),
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify(expectedData),
        })
      );
    });

    it('should handle validation errors', async () => {
      const invalidMatchData = {
        series_id: '',
        match_number: 1,
        date: 'invalid-date',
        status: 'live',
        team_a_player_count: 0,
        team_b_player_count: 0,
        total_overs: 0,
        toss_winner: 'A',
        toss_type: 'H',
        batting_team: 'A',
      };

      const mockResponse = {
        ok: false,
        status: 400,
        json: jest.fn().mockResolvedValue({ message: 'Validation failed' }),
      };

      mockFetch.mockResolvedValue(mockResponse as any);

      await expect(apiService.createMatch(invalidMatchData as any)).rejects.toThrow(ApiError);
    });
  });

  describe('updateMatch', () => {
    it('should update an existing match successfully', async () => {
      const matchId = '1';
      const updateData: Partial<Match> = {
        team_a_player_count: 10,
        team_b_player_count: 10,
        total_overs: 15,
      };

      const mockUpdatedMatch: Match = {
        id: '1',
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
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z',
      };

      const mockResponse = {
        ok: true,
        status: 200,
        json: jest.fn().mockResolvedValue(mockUpdatedMatch),
      };

      mockFetch.mockResolvedValue(mockResponse as any);

      const result = await apiService.updateMatch(matchId, updateData);

      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining(`/matches/${matchId}?_t=`),
        expect.objectContaining({
          method: 'PUT',
          body: JSON.stringify(updateData),
          headers: expect.objectContaining({
            'Content-Type': 'application/json',
          }),
        })
      );

      expect(result).toEqual({
        data: mockUpdatedMatch,
        success: true,
        message: undefined,
      });
    });

    it('should handle update errors', async () => {
      const matchId = '999';
      const updateData: Partial<Match> = {
        team_a_player_count: 10,
      };

      const mockResponse = {
        ok: false,
        status: 404,
        json: jest.fn().mockResolvedValue({ message: 'Match not found' }),
      };

      mockFetch.mockResolvedValue(mockResponse as any);

      await expect(apiService.updateMatch(matchId, updateData)).rejects.toThrow(ApiError);
    });
  });

  describe('deleteMatch', () => {
    it('should delete a match successfully', async () => {
      const matchId = '1';

      const mockResponse = {
        ok: true,
        status: 200,
        json: jest.fn().mockResolvedValue({ success: true }),
      };

      mockFetch.mockResolvedValue(mockResponse as any);

      const result = await apiService.deleteMatch(matchId);

      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining(`/matches/${matchId}?_t=`),
        expect.objectContaining({
          method: 'DELETE',
        })
      );

      expect(result).toEqual({
        data: { success: true },
        success: true,
        message: undefined,
      });
    });

    it('should handle delete errors', async () => {
      const matchId = '999';

      const mockResponse = {
        ok: false,
        status: 404,
        json: jest.fn().mockResolvedValue({ message: 'Match not found' }),
      };

      mockFetch.mockResolvedValue(mockResponse as any);

      await expect(apiService.deleteMatch(matchId)).rejects.toThrow(ApiError);
    });
  });

  describe('getMatchesBySeries', () => {
    it('should fetch matches by series ID', async () => {
      const seriesId = 'series-1';
      const mockMatches: Match[] = [
        {
          id: '1',
          series_id: seriesId,
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
        ok: true,
        json: jest.fn().mockResolvedValue(mockMatches),
      };

      mockFetch.mockResolvedValue(mockResponse as any);

      const result = await apiService.getMatchesBySeries(seriesId);

      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining(`/matches/series/${seriesId}?_t=`),
        expect.objectContaining({
          headers: expect.objectContaining({
            'Content-Type': 'application/json',
            'Accept': 'application/json',
          }),
          mode: 'cors',
          credentials: 'omit',
        })
      );

      expect(result).toEqual({
        data: mockMatches,
        success: true,
        message: undefined,
      });
    });

    it('should handle empty series', async () => {
      const seriesId = 'empty-series';

      const mockResponse = {
        ok: true,
        json: jest.fn().mockResolvedValue([]),
      };

      mockFetch.mockResolvedValue(mockResponse as any);

      const result = await apiService.getMatchesBySeries(seriesId);

      expect(result).toEqual({
        data: [],
        success: true,
        message: undefined,
      });
    });
  });

  describe('Network Error Handling', () => {
    it('should retry on network errors', async () => {
      const mockMatches: Match[] = [];

      const mockResponse = {
        ok: true,
        json: jest.fn().mockResolvedValue(mockMatches),
      };

      // First call fails with network error, second call succeeds
      mockFetch
        .mockRejectedValueOnce(new TypeError('Failed to fetch'))
        .mockResolvedValueOnce(mockResponse as any);

      const result = await apiService.getMatches();

      expect(mockFetch).toHaveBeenCalledTimes(2);
      expect(result).toEqual({
        data: mockMatches,
        success: true,
        message: undefined,
      });
    });

    it('should throw ApiError after max retries', async () => {
      // All calls fail with network error
      mockFetch.mockRejectedValue(new TypeError('Failed to fetch'));

      await expect(apiService.getMatches()).rejects.toThrow(ApiError);
      expect(mockFetch).toHaveBeenCalledTimes(4); // Initial + 3 retries
    }, 10000); // Increase timeout to 10 seconds
  });

  describe('Request Configuration', () => {
    it('should include proper headers', async () => {
      const mockResponse = {
        ok: true,
        json: jest.fn().mockResolvedValue([]),
      };

      mockFetch.mockResolvedValue(mockResponse as any);

      await apiService.getMatches();

      expect(mockFetch).toHaveBeenCalledWith(
        expect.any(String),
        expect.objectContaining({
          headers: expect.objectContaining({
            'Content-Type': 'application/json',
            'Accept': 'application/json',
            'Cache-Control': 'no-cache, no-store, must-revalidate',
            'Pragma': 'no-cache',
            'Expires': '0',
          }),
          mode: 'cors',
          credentials: 'omit',
        })
      );
    });

    it('should add cache-busting parameter', async () => {
      const mockResponse = {
        ok: true,
        json: jest.fn().mockResolvedValue([]),
      };

      mockFetch.mockResolvedValue(mockResponse as any);

      await apiService.getMatches();

      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringMatching(/\?_t=\d+/),
        expect.any(Object)
      );
    });
  });
});
