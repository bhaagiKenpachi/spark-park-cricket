import { Series } from '@/store/reducers/seriesSlice';
import { Match } from '@/store/reducers/matchSlice';
import {
    ScorecardResponse,
    BallEventRequest,
    OverSummary,
    InningsSummary,
    BallType,
    RunType,
} from '@/store/reducers/scorecardSlice';

const API_BASE_URL =
    process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

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

        try {
            const response = await fetch(url, config);

            if (!response.ok) {
                const errorData = await response.json().catch(() => ({}));

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

            if (error instanceof ApiError) {
                throw error;
            }
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

        try {
            const response = await fetch(url, config);

            if (!response.ok) {
                const errorData = await response.json().catch(() => ({}));
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
            // Retry logic for network errors
            if (retryCount < 3 && error instanceof TypeError) {
                await new Promise(resolve =>
                    setTimeout(resolve, 1000 * (retryCount + 1))
                );
                return this.scorecardRequest<T>(endpoint, options, retryCount + 1);
            }

            if (error instanceof ApiError) {
                throw error;
            }

            throw new ApiError(
                error instanceof Error ? error.message : 'Unknown error occurred',
                500,
                error
            );
        }
    }

    // Series API methods
    async getSeries(): Promise<ApiResponse<Series[]>> {
        return this.request<Series[]>('/series');
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
