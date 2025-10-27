import { Series } from '@/store/reducers/seriesSlice';
import { Match } from '@/store/reducers/matchSlice';
import { trackApiError, trackPerformance } from '@/lib/analytics';
import {
  FallOfWickets,
  FallOfWicketsSummary,
  CreateFallOfWicketsRequest,
  UpdateFallOfWicketsRequest,
  FallOfWicketsFilters,
} from '@/types/fallOfWickets';

export interface PaginatedSeriesResult {
  series: Series[];
  total_items: number;
  page: number;
  page_size: number;
  total_pages: number;
}
import {
  ScorecardResponse,
  BallEventRequest,
  OverSummary,
  InningsSummary,
  BallType,
  RunType,
} from '@/store/reducers/scorecardSlice';

// Ball event response from the API
export interface BallEventResponse {
  event_type: string;
  match_id: string;
  innings_number: number;
  ball_number: number;
  ball_type: string;
  run_type: string;
  runs: number;
  byes: number;
  total_runs: number;
  is_wicket: boolean;
  wicket_type: string;
  innings_runs: number;
  innings_wickets: number;
  innings_overs: string;
  timestamp: string;
  stream_id: string;
}

// Backend API Configuration
// Set NEXT_PUBLIC_API_URL in .env.local file or environment variables
// Default: https://spark-park.dojima.foundation/api/v1
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL!;

export interface ApiResponse<T> {
  data: T;
  message?: string;
  success: boolean;
}

export interface ScorecardApiResponse<T> {
  data: {
    data: T;
  };
  message?: string;
  success: boolean;
}

export interface ApiErrorInterface {
  message: string;
  status: number;
  details?: unknown;
}

class ApiService {
  private baseURL: string;

  constructor(baseURL: string = API_BASE_URL) {
    this.baseURL = baseURL;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {},
    retryCount: number = 0
  ): Promise<ApiResponse<T>> {
    // Add cache-busting parameter
    const separator = endpoint.includes('?') ? '&' : '?';
    const url = `${this.baseURL}${endpoint}${separator}_t=${Date.now()}`;

    const defaultHeaders = {
      'Content-Type': 'application/json',
      Accept: 'application/json',
      'Cache-Control': 'no-cache, no-store, must-revalidate',
      Pragma: 'no-cache',
      Expires: '0',
    };

    const config: RequestInit = {
      ...options,
      headers: {
        ...defaultHeaders,
        ...options.headers,
      },
      mode: 'cors',
      credentials: 'include',
    };

    const startTime = performance.now();

    try {
      const response = await fetch(url, config);
      const endTime = performance.now();
      const requestDuration = endTime - startTime;

      // Track API response time
      trackPerformance({
        metric_name: 'api_response_time',
        metric_value: requestDuration,
        component: 'api_service'
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));

        // Track API error
        trackApiError({
          endpoint,
          status_code: response.status,
          error_message: errorData.message || `HTTP error! status: ${response.status}`,
          request_duration: requestDuration
        });

        // Retry logic for 503 errors
        if (response.status === 503 && retryCount < 3) {
          await new Promise(resolve =>
            setTimeout(resolve, 1000 * (retryCount + 1))
          ); // Exponential backoff
          return this.request<T>(endpoint, options, retryCount + 1);
        }

        throw new ApiError(
          errorData.message || `HTTP error! status: ${response.status}`,
          response.status,
          errorData
        );
      }

      const data = await response.json();
      return {
        data,
        success: true,
        message: data.message,
      };
    } catch (error) {
      const endTime = performance.now();
      const requestDuration = endTime - startTime;

      // Retry logic for network errors
      if (
        retryCount < 3 &&
        (error instanceof TypeError ||
          (error instanceof Error && error.message.includes('Failed to fetch')))
      ) {
        await new Promise(resolve =>
          setTimeout(resolve, 1000 * (retryCount + 1))
        ); // Exponential backoff
        return this.request<T>(endpoint, options, retryCount + 1);
      }

      // Track API error for network errors
      if (error instanceof ApiError) {
        trackApiError({
          endpoint,
          status_code: error.status,
          error_message: error.message,
          request_duration: requestDuration
        });
        throw error;
      }

      // Track network errors
      trackApiError({
        endpoint,
        status_code: 0,
        error_message: error instanceof Error ? error.message : 'Network error',
        request_duration: requestDuration
      });

      throw new ApiError(
        error instanceof Error ? error.message : 'Network error',
        0,
        error
      );
    }
  }

  private async scorecardRequest<T>(
    endpoint: string,
    options: RequestInit = {},
    retryCount: number = 0
  ): Promise<ScorecardApiResponse<T>> {
    // Add cache-busting parameter
    const separator = endpoint.includes('?') ? '&' : '?';
    const url = `${this.baseURL}${endpoint}${separator}_t=${Date.now()}`;

    const defaultHeaders = {
      'Content-Type': 'application/json',
    };

    const config: RequestInit = {
      ...options,
      headers: {
        ...defaultHeaders,
        ...options.headers,
      },
    };

    const startTime = performance.now();

    try {
      const response = await fetch(url, config);
      const endTime = performance.now();
      const requestDuration = endTime - startTime;

      // Track API response time
      trackPerformance({
        metric_name: 'api_response_time',
        metric_value: requestDuration,
        component: 'scorecard_api_service'
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));

        // Track API error
        trackApiError({
          endpoint,
          status_code: response.status,
          error_message: errorData.message || `HTTP error! status: ${response.status}`,
          request_duration: requestDuration
        });

        throw new ApiError(
          errorData.message || `HTTP error! status: ${response.status}`,
          response.status,
          errorData
        );
      }

      const data = await response.json();
      return {
        data,
        success: true,
        message: data.message,
      };
    } catch (error) {
      const endTime = performance.now();
      const requestDuration = endTime - startTime;

      // Retry logic for network errors
      if (retryCount < 3 && error instanceof TypeError) {
        await new Promise(resolve =>
          setTimeout(resolve, 1000 * (retryCount + 1))
        );
        return this.scorecardRequest<T>(endpoint, options, retryCount + 1);
      }

      if (error instanceof ApiError) {
        // Track API error for existing ApiError
        trackApiError({
          endpoint,
          status_code: error.status,
          error_message: error.message,
          request_duration: requestDuration
        });
        throw error;
      }

      // Track network errors
      trackApiError({
        endpoint,
        status_code: 0,
        error_message: error instanceof Error ? error.message : 'Unknown error occurred',
        request_duration: requestDuration
      });

      throw new ApiError(
        error instanceof Error ? error.message : 'Unknown error occurred',
        500,
        error
      );
    }
  }

  // Series API methods
  async getSeries(
    limit?: number,
    offset?: number
  ): Promise<ApiResponse<PaginatedSeriesResult>> {
    const params = new URLSearchParams();
    if (limit !== undefined) params.append('limit', limit.toString());
    if (offset !== undefined) params.append('offset', offset.toString());

    const queryString = params.toString();
    const endpoint = queryString ? `/series?${queryString}` : '/series';

    return this.request<PaginatedSeriesResult>(endpoint);
  }

  async getSeriesById(id: string): Promise<ApiResponse<Series>> {
    return this.request<Series>(`/series/${id}`);
  }

  async createSeries(
    seriesData: Omit<Series, 'id' | 'created_at' | 'updated_at'>
  ): Promise<ApiResponse<Series>> {
    return this.request<Series>('/series', {
      method: 'POST',
      body: JSON.stringify(seriesData),
    });
  }

  async updateSeries(
    id: string,
    seriesData: Partial<Series>
  ): Promise<ApiResponse<Series>> {
    return this.request<Series>(`/series/${id}`, {
      method: 'PUT',
      body: JSON.stringify(seriesData),
    });
  }

  async deleteSeries(id: string): Promise<ApiResponse<void>> {
    return this.request<void>(`/series/${id}`, {
      method: 'DELETE',
    });
  }

  // Match API methods
  async getMatches(): Promise<ApiResponse<Match[]>> {
    return this.request<Match[]>('/matches');
  }

  async getMatchById(id: string): Promise<ApiResponse<Match>> {
    return this.request<Match>(`/matches/${id}`);
  }

  async createMatch(
    matchData: Omit<Match, 'id' | 'created_at' | 'updated_at' | 'match_number'>
  ): Promise<ApiResponse<Match>> {
    // Backend will auto-generate match_number, so we don't need to send it
    return this.request<Match>('/matches', {
      method: 'POST',
      body: JSON.stringify(matchData),
    });
  }

  async updateMatch(
    id: string,
    matchData: Partial<Omit<Match, 'match_number'>>
  ): Promise<ApiResponse<Match>> {
    return this.request<Match>(`/matches/${id}`, {
      method: 'PUT',
      body: JSON.stringify(matchData),
    });
  }

  async deleteMatch(id: string): Promise<ApiResponse<void>> {
    return this.request<void>(`/matches/${id}`, {
      method: 'DELETE',
    });
  }

  async startMatch(id: string): Promise<ApiResponse<{ message: string; match_id: string }>> {
    return this.request<{ message: string; match_id: string }>('/scorecard/start', {
      method: 'POST',
      body: JSON.stringify({ match_id: id }),
    });
  }

  async getMatchesBySeries(seriesId: string): Promise<ApiResponse<Match[]>> {
    return this.request<Match[]>(`/matches/series/${seriesId}`);
  }

  // Scorecard API methods
  async getScorecard(
    matchId: string
  ): Promise<ScorecardApiResponse<ScorecardResponse>> {
    return this.scorecardRequest<ScorecardResponse>(`/scorecard/${matchId}`);
  }

  async startScoring(
    matchId: string
  ): Promise<ApiResponse<{ message: string; match_id: string }>> {
    try {
      return await this.request<{ message: string; match_id: string }>(
        '/scorecard/start',
        {
          method: 'POST',
          body: JSON.stringify({ match_id: matchId }),
        }
      );
    } catch (error) {
      // If scoring is already started, return success response instead of error
      if (
        error instanceof ApiError &&
        (error.message.includes('scoring already started') ||
          error.message.includes('already started for this match'))
      ) {
        return {
          data: { message: 'Scoring already active', match_id: matchId },
          success: true,
          message: 'Scoring already active',
        };
      }
      throw error;
    }
  }

  async addBall(ballEvent: BallEventRequest): Promise<
    ApiResponse<{
      message: string;
      match_id: string;
      innings_number: number;
      ball_type: string;
      run_type: string;
      runs: number;
      byes: number;
      is_wicket: boolean;
    }>
  > {
    return this.request('/scorecard/ball', {
      method: 'POST',
      body: JSON.stringify(ballEvent),
    });
  }

  // Common ball scoring function
  async scoreBall(
    matchId: string,
    inningsNumber: number,
    ballType: BallType,
    runType: RunType,
    runs: number,
    byes: number = 0,
    isWicket: boolean = false
  ): Promise<
    ApiResponse<{
      message: string;
      match_id: string;
      innings_number: number;
      ball_type: string;
      run_type: string;
      runs: number;
      byes: number;
      is_wicket: boolean;
    }>
  > {
    const ballEvent: BallEventRequest = {
      match_id: matchId,
      innings_number: inningsNumber,
      ball_type: ballType,
      run_type: runType,
      runs,
      byes,
      is_wicket: isWicket,
    };

    return this.addBall(ballEvent);
  }

  async undoBall(
    matchId: string,
    inningsNumber: number = 1
  ): Promise<
    ApiResponse<{
      message: string;
      match_id: string;
      innings_number: number;
    }>
  > {
    return this.request(`/scorecard/${matchId}/ball?innings=${inningsNumber}`, {
      method: 'DELETE',
    });
  }

  async getCurrentOver(
    matchId: string,
    inningsNumber: number = 1
  ): Promise<ScorecardApiResponse<OverSummary>> {
    return this.scorecardRequest<OverSummary>(
      `/scorecard/${matchId}/current-over?innings=${inningsNumber}`
    );
  }

  async getInnings(
    matchId: string,
    inningsNumber: number
  ): Promise<ScorecardApiResponse<InningsSummary>> {
    return this.scorecardRequest<InningsSummary>(
      `/scorecard/${matchId}/innings/${inningsNumber}`
    );
  }

  async getOver(
    matchId: string,
    inningsNumber: number,
    overNumber: number
  ): Promise<ScorecardApiResponse<OverSummary>> {
    return this.scorecardRequest<OverSummary>(
      `/scorecard/${matchId}/innings/${inningsNumber}/over/${overNumber}`
    );
  }

  // Time Tracking API methods
  async getTimeTracking(matchId: string): Promise<ApiResponse<any>> {
    const response = await this.request<any>(`/scorecard/${matchId}/time-tracking`);
    // Extract the nested data from the API response
    const result: ApiResponse<any> = {
      data: response.data.data,
      success: response.success
    };
    if (response.message) {
      result.message = response.message;
    }
    return result;
  }

  // Vote API methods
  async getVotes(filters?: { status?: string; type?: string; created_by?: string; limit?: number; offset?: number; page?: number; page_size?: number }): Promise<ApiResponse<any>> {
    const params = new URLSearchParams();

    if (filters?.status) params.append('status', filters.status);
    if (filters?.type) params.append('type', filters.type);
    if (filters?.created_by) params.append('created_by', filters.created_by);

    // Support both page-based and offset-based pagination
    if (filters?.page) params.append('page', filters.page.toString());
    if (filters?.page_size) params.append('page_size', filters.page_size.toString());
    if (filters?.limit) params.append('limit', filters.limit.toString());
    if (filters?.offset) params.append('offset', filters.offset.toString());

    const queryString = params.toString();
    const endpoint = queryString ? `/votes?${queryString}` : '/votes';

    return this.request<any>(endpoint);
  }

  async getVoteById(id: string): Promise<ApiResponse<any>> {
    return this.request<any>(`/votes/${id}`);
  }

  async getVoteWithResults(id: string): Promise<ApiResponse<any>> {
    return this.request<any>(`/votes/${id}/results`);
  }

  async createVote(voteData: any): Promise<ApiResponse<any>> {
    return this.request<any>('/votes', {
      method: 'POST',
      body: JSON.stringify(voteData),
    });
  }

  async updateVote(id: string, voteData: any): Promise<ApiResponse<any>> {
    return this.request<any>(`/votes/${id}`, {
      method: 'PUT',
      body: JSON.stringify(voteData),
    });
  }

  async deleteVote(id: string): Promise<ApiResponse<{ message: string }>> {
    return this.request<{ message: string }>(`/votes/${id}`, {
      method: 'DELETE',
    });
  }

  async castVote(voteId: string, voteData: any): Promise<ApiResponse<any>> {
    return this.request<any>(`/votes/${voteId}/vote`, {
      method: 'POST',
      body: JSON.stringify(voteData),
    });
  }

  async getUserVote(voteId: string): Promise<ApiResponse<any>> {
    return this.request<any>(`/votes/${voteId}/my-vote`);
  }

  async hasUserVoted(voteId: string): Promise<ApiResponse<{ has_voted: boolean }>> {
    return this.request<{ has_voted: boolean }>(`/votes/${voteId}/has-voted`);
  }

  async closeVote(voteId: string): Promise<ApiResponse<any>> {
    return this.request<any>(`/votes/${voteId}/close`, {
      method: 'POST',
    });
  }

  async cancelVote(voteId: string): Promise<ApiResponse<any>> {
    return this.request<any>(`/votes/${voteId}/cancel`, {
      method: 'POST',
    });
  }

  // User API methods
  async updateUserName(name: string): Promise<ApiResponse<any>> {
    return this.request<any>('/users/me', {
      method: 'PUT',
      body: JSON.stringify({ name }),
    });
  }

  // Team API methods
  async createTeam(voteId: string, teamData: any): Promise<ApiResponse<any>> {
    return this.request<any>(`/votes/${voteId}/teams`, {
      method: 'POST',
      body: JSON.stringify(teamData),
    });
  }

  async getTeamsByVoteId(voteId: string): Promise<ApiResponse<any>> {
    return this.request<any>(`/votes/${voteId}/teams`);
  }

  async getTeamById(teamId: string): Promise<ApiResponse<any>> {
    return this.request<any>(`/teams/${teamId}`);
  }

  async updateTeam(teamId: string, teamData: any): Promise<ApiResponse<any>> {
    return this.request<any>(`/teams/${teamId}`, {
      method: 'PUT',
      body: JSON.stringify(teamData),
    });
  }

  async deleteTeam(teamId: string): Promise<ApiResponse<{ message: string }>> {
    return this.request<{ message: string }>(`/teams/${teamId}`, {
      method: 'DELETE',
    });
  }

  async getTeamPlayers(teamId: string): Promise<ApiResponse<any>> {
    return this.request<any>(`/teams/${teamId}/players`);
  }

  async addPlayerToTeam(teamId: string, userId: string): Promise<ApiResponse<{ message: string }>> {
    return this.request<{ message: string }>(`/teams/${teamId}/players`, {
      method: 'POST',
      body: JSON.stringify({ user_id: userId }),
    });
  }

  async addPlayersToTeam(teamId: string, userIds: string[]): Promise<ApiResponse<{ message: string }>> {
    return this.request<{ message: string }>(`/teams/${teamId}/players/bulk`, {
      method: 'POST',
      body: JSON.stringify({ user_ids: userIds }),
    });
  }

  async removePlayerFromTeam(teamId: string, playerId: string): Promise<ApiResponse<{ message: string }>> {
    return this.request<{ message: string }>(`/teams/${teamId}/players/${playerId}`, {
      method: 'DELETE',
    });
  }

  // Get ball events for a match
  async getBallEvents(matchId: string): Promise<ApiResponse<BallEventResponse[]>> {
    return this.request<BallEventResponse[]>(`/scorecard/${matchId}/events`, {
      method: 'GET',
    });
  }

  // Fall of Wickets API methods
  async createFallOfWickets(fallOfWicketsData: CreateFallOfWicketsRequest): Promise<ApiResponse<FallOfWickets>> {
    return this.request<FallOfWickets>('/fall-of-wickets', {
      method: 'POST',
      body: JSON.stringify(fallOfWicketsData),
    });
  }

  async getFallOfWicketsById(id: string): Promise<ApiResponse<FallOfWickets>> {
    return this.request<FallOfWickets>(`/fall-of-wickets/${id}`);
  }

  async listFallOfWickets(filters?: FallOfWicketsFilters): Promise<ApiResponse<FallOfWickets[]>> {
    const params = new URLSearchParams();
    if (filters?.match_id) params.append('match_id', filters.match_id);
    if (filters?.innings_id) params.append('innings_id', filters.innings_id);
    if (filters?.limit) params.append('limit', filters.limit.toString());
    if (filters?.offset) params.append('offset', filters.offset.toString());

    const queryString = params.toString();
    const endpoint = queryString ? `/fall-of-wickets?${queryString}` : '/fall-of-wickets';

    return this.request<FallOfWickets[]>(endpoint);
  }

  async updateFallOfWickets(id: string, updateData: UpdateFallOfWicketsRequest): Promise<ApiResponse<FallOfWickets>> {
    return this.request<FallOfWickets>(`/fall-of-wickets/${id}`, {
      method: 'PUT',
      body: JSON.stringify(updateData),
    });
  }

  async deleteFallOfWickets(id: string): Promise<ApiResponse<void>> {
    return this.request<void>(`/fall-of-wickets/${id}`, {
      method: 'DELETE',
    });
  }

  async getFallOfWicketsByMatch(matchId: string): Promise<ApiResponse<FallOfWickets[]>> {
    return this.request<FallOfWickets[]>(`/matches/${matchId}/fall-of-wickets`);
  }

  async getFallOfWicketsByInnings(inningsId: string): Promise<ApiResponse<FallOfWickets[]>> {
    return this.request<FallOfWickets[]>(`/innings/${inningsId}/fall-of-wickets`);
  }

  async getFallOfWicketsByBall(ballId: string): Promise<ApiResponse<FallOfWickets>> {
    return this.request<FallOfWickets>(`/balls/${ballId}/fall-of-wickets`);
  }

  async getFallOfWicketsSummary(matchId: string, inningsId?: string): Promise<ApiResponse<FallOfWicketsSummary>> {
    const params = new URLSearchParams();
    params.append('match_id', matchId);
    if (inningsId) params.append('innings_id', inningsId);

    return this.request<FallOfWicketsSummary>(`/fall-of-wickets/summary?${params.toString()}`);
  }

  async createFallOfWicketsFromBall(ballId: string, score: number): Promise<ApiResponse<FallOfWickets>> {
    return this.request<FallOfWickets>('/fall-of-wickets/from-ball', {
      method: 'POST',
      body: JSON.stringify({ ball_id: ballId, score }),
    });
  }
}

export { ApiService };
export class ApiError extends Error {
  public status: number;
  public details?: unknown;

  constructor(message: string, status: number, details?: unknown) {
    super(message);
    this.name = 'ApiError';
    this.status = status;
    this.details = details;
  }
}

export const apiService = new ApiService();
export default apiService;
