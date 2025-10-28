// Fall of Wickets Types
export interface FallOfWickets {
    id: string;
    match_id: string;
    innings_id: string;
    innings_number: number;
    over_id: string;
    ball_id: string;
    wicket_number: number;
    score: number;
    over_number: number;
    ball_number: number;
    created_at: string;
}

export interface WicketFall {
    wicket_number: number;
    score: number;
    over_number: number;
    ball_number: number;
    over_position: string; // e.g., "15.3"
}

export interface FallOfWicketsSummary {
    match_id: string;
    innings_id?: string;
    total_wickets: number;
    wickets: WicketFall[];
}

export interface CreateFallOfWicketsRequest {
    match_id: string;
    innings_id: string;
    innings_number: number;
    over_id: string;
    ball_id: string;
    wicket_number: number;
    score: number;
    over_number: number;
    ball_number: number;
}

export interface UpdateFallOfWicketsRequest {
    score?: number;
    over_number?: number;
    ball_number?: number;
}

export interface FallOfWicketsFilters {
    match_id?: string;
    innings_id?: string;
    limit?: number;
    offset?: number;
}