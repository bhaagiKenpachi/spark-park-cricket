// Time Tracking Types for Cricket Match Analysis

export interface TimeTrackingResponse {
    match_id: string;
    innings?: InningsTimeTracking[];
    total_match_time_seconds: number;
}

export interface InningsTimeTracking {
    innings_number: number;
    batting_team: 'A' | 'B';
    start_time: string | null;
    end_time: string | null;
    duration_seconds: number;
    status: 'in_progress' | 'completed';
    overs?: OverTimeTracking[];
}

export interface OverTimeTracking {
    over_number: number;
    start_time: string | null;
    end_time: string | null;
    duration_seconds: number;
    status: 'in_progress' | 'completed';
    total_runs: number;
    total_balls: number;
    total_wickets: number;
}

// Utility functions for time formatting
export const formatDuration = (seconds: number): string => {
    if (seconds < 60) {
        return `${seconds}s`;
    }

    const minutes = Math.floor(seconds / 60);
    const remainingSeconds = seconds % 60;

    if (minutes < 60) {
        return remainingSeconds > 0 ? `${minutes}m ${remainingSeconds}s` : `${minutes}m`;
    }

    const hours = Math.floor(minutes / 60);
    const remainingMinutes = minutes % 60;

    return remainingMinutes > 0 ? `${hours}h ${remainingMinutes}m` : `${hours}h`;
};

export const formatTime = (timeString: string | null): string => {
    if (!timeString) return 'Not started';

    try {
        const date = new Date(timeString);
        return date.toLocaleTimeString('en-US', {
            hour: '2-digit',
            minute: '2-digit',
            second: '2-digit',
            hour12: false
        });
    } catch {
        return 'Invalid time';
    }
};

export const formatDateTime = (timeString: string | null): string => {
    if (!timeString) return 'Not started';

    try {
        const date = new Date(timeString);
        return date.toLocaleString('en-US', {
            year: 'numeric',
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit',
            second: '2-digit',
            hour12: false
        });
    } catch {
        return 'Invalid time';
    }
};


