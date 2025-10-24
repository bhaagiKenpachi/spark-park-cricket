// PostHog Event Names
export enum PostHogEvent {
    // Series Events
    SERIES_VIEWED = 'series_viewed',
    SERIES_CREATED = 'series_created',
    SERIES_EDITED = 'series_edited',
    SERIES_DELETED = 'series_deleted',
    SERIES_PAGINATION_CHANGED = 'series_pagination_changed',

    // Match Events
    MATCH_VIEWED = 'match_viewed',
    MATCH_CREATED = 'match_created',
    MATCH_EDITED = 'match_edited',
    MATCH_DELETED = 'match_deleted',

    // Scorecard Events
    SCORECARD_VIEWED = 'scorecard_viewed',
    LIVE_SCORING_STARTED = 'live_scoring_started',
    BALL_ADDED = 'ball_added',
    INNINGS_COMPLETED = 'innings_completed',
    MATCH_COMPLETED = 'match_completed',
    OVER_COMPLETED = 'over_completed',

    // WebSocket Events
    WEBSOCKET_CONNECTED = 'websocket_connected',
    WEBSOCKET_DISCONNECTED = 'websocket_disconnected',
    WEBSOCKET_ERROR = 'websocket_error',

    // Authentication Events
    USER_LOGGED_IN = 'user_logged_in',
    USER_LOGGED_OUT = 'user_logged_out',
    USER_PROFILE_UPDATED = 'user_profile_updated',

    // Navigation Events
    PAGE_VIEWED = '$pageview',
    PAGE_LEFT = '$pageleave',

    // Error Events
    ERROR_OCCURRED = 'error_occurred',
    API_ERROR = 'api_error',

    // Performance Events
    PERFORMANCE_METRIC = 'performance_metric',
    PAGE_LOAD_TIME = 'page_load_time',
    API_RESPONSE_TIME = 'api_response_time',

    // User Engagement Events
    FEATURE_USED = 'feature_used',
    USER_ACTION = 'user_action',
    TIME_ON_PAGE = 'time_on_page',

    // Funnel and Journey Events
    FUNNEL_STEP = 'funnel_step',

    // Business Intelligence Events
    CONVERSION = 'conversion',
    ACHIEVEMENT_UNLOCKED = 'achievement_unlocked',
}

// Series Event Properties
export interface SeriesEventProperties {
    series_id: string;
    series_name?: string;
    start_date?: string;
    end_date?: string;
}

export interface SeriesPaginationProperties {
    page: number;
    page_size: number;
    total_items?: number;
}

// Match Event Properties
export interface MatchEventProperties {
    match_id: string;
    series_id: string;
    match_number?: number;
    match_status?: string;
    total_overs?: number;
    team_a_player_count?: number;
    team_b_player_count?: number;
}

// Scorecard Event Properties
export interface ScorecardEventProperties {
    match_id: string;
    innings_number?: number;
    total_runs?: number;
    total_wickets?: number;
    total_overs?: number;
    match_status?: string;
}

export interface BallEventProperties {
    match_id: string;
    innings_number: number;
    over_number: number;
    ball_number: number;
    ball_type: string;
    run_type: string;
    runs: number;
    is_wicket: boolean;
    wicket_type?: string;
}

export interface OverEventProperties {
    match_id: string;
    innings_number: number;
    over_number: number;
    total_runs: number;
    total_wickets: number;
    total_balls: number;
}

// WebSocket Event Properties
export interface WebSocketEventProperties {
    match_id: string;
    connection_status?: 'connected' | 'disconnected' | 'error';
    error_message?: string;
    reconnect_attempt?: number;
}

// User Event Properties
export interface UserEventProperties {
    user_id: string;
    email?: string;
    name?: string;
    provider?: string;
    role?: string;
}

// Feature Flags
export enum FeatureFlag {
    NEW_SCORECARD_UI = 'new-scorecard-ui',
    ENABLE_ADS = 'enable-ads',
    BETA_FEATURES = 'beta-features',
    ADVANCED_ANALYTICS = 'advanced-analytics',
    WEBSOCKET_AUTO_RECONNECT = 'websocket-auto-reconnect',
}

// User Properties Interface
export interface UserProperties {
    id: string;
    email?: string;
    name?: string;
    created_at?: string;
    role?: string;
    total_series_created?: number;
    total_matches_scored?: number;
}


