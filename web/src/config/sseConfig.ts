// SSE Configuration for Web App
// This file centralizes SSE timeout configuration using environment variables

export interface SSEConfig {
    inactivityWarningThreshold: number; // milliseconds
    inactivityDisconnectThreshold: number; // milliseconds
    activityCheckInterval: number; // milliseconds
    userActivityTimeout: number; // milliseconds
    debounceMs: number; // milliseconds
    ballEventTimeoutSeconds: number; // seconds - timeout for ball events
}

// Parse environment variables with fallbacks
const parseEnvNumber = (envVar: string | undefined, defaultValue: number): number => {
    if (!envVar) return defaultValue;
    const parsed = parseInt(envVar, 10);
    return isNaN(parsed) ? defaultValue : parsed;
};

// Parse duration strings (e.g., "2m", "30s") to milliseconds
const parseDurationToMs = (duration: string | undefined, defaultValue: number): number => {
    if (!duration) return defaultValue;

    // Handle simple number format (assume seconds)
    if (/^\d+$/.test(duration)) {
        return parseInt(duration, 10) * 1000;
    }

    // Handle duration format (e.g., "2m", "30s", "1m30s")
    const match = duration.match(/^(\d+)([sm])$/);
    if (match && match[1] && match[2]) {
        const value = parseInt(match[1], 10);
        const unit = match[2];
        if (unit === 'm') return value * 60 * 1000; // minutes to milliseconds
        if (unit === 's') return value * 1000; // seconds to milliseconds
    }

    return defaultValue;
};

export const sseConfig: SSEConfig = {
    // Inactivity warning threshold (default: 30 seconds for testing)
    inactivityWarningThreshold: parseDurationToMs(
        process.env.NEXT_PUBLIC_SSE_INACTIVITY_WARNING_THRESHOLD,
        30 * 1000 // 30 seconds for testing
    ),

    // Inactivity disconnect threshold (default: 1 minute for testing)
    inactivityDisconnectThreshold: parseDurationToMs(
        process.env.NEXT_PUBLIC_SSE_INACTIVITY_DISCONNECT_THRESHOLD,
        60 * 1000 // 1 minute for testing
    ),

    // Activity check interval (default: 10 seconds for testing)
    activityCheckInterval: parseDurationToMs(
        process.env.NEXT_PUBLIC_SSE_ACTIVITY_CHECK_INTERVAL,
        10 * 1000 // 10 seconds for testing
    ),

    // User activity timeout (default: 5 minutes)
    userActivityTimeout: parseDurationToMs(
        process.env.NEXT_PUBLIC_SSE_USER_ACTIVITY_TIMEOUT,
        5 * 60 * 1000 // 5 minutes
    ),

    // Debounce time for user activity (default: 1 second)
    debounceMs: parseEnvNumber(
        process.env.NEXT_PUBLIC_SSE_DEBOUNCE_MS,
        1000 // 1 second
    ),

    // Ball event timeout in seconds (default: 120 seconds / 2 minutes)
    ballEventTimeoutSeconds: parseEnvNumber(
        process.env.NEXT_PUBLIC_SSE_BALL_EVENT_TIMEOUT,
        120 // 120 seconds / 2 minutes
    ),
};

// Log configuration on client side (always show for debugging)
if (typeof window !== 'undefined') {
    console.log('🔧 SSE Configuration:', {
        inactivityWarningThreshold: `${sseConfig.inactivityWarningThreshold / 1000}s`,
        inactivityDisconnectThreshold: `${sseConfig.inactivityDisconnectThreshold / 1000}s`,
        activityCheckInterval: `${sseConfig.activityCheckInterval / 1000}s`,
        userActivityTimeout: `${sseConfig.userActivityTimeout / 1000}s`,
        debounceMs: `${sseConfig.debounceMs}ms`,
        ballEventTimeoutSeconds: `${sseConfig.ballEventTimeoutSeconds}s`,
    });
}
