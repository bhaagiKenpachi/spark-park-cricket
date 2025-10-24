import { useEffect, useRef, useCallback } from 'react';
import { trackEvent, trackPerformance } from '@/lib/analytics';

interface SessionData {
    sessionId: string;
    startTime: number;
    endTime?: number;
    pagesVisited: string[];
    actionsPerformed: string[];
    sessionQuality: number;
    deviceInfo: {
        userAgent: string;
        screenResolution: string;
        timezone: string;
        language: string;
    };
}

interface SessionTrackingOptions {
    trackPageViews?: boolean;
    trackUserActions?: boolean;
    trackSessionQuality?: boolean;
    sessionTimeout?: number; // in milliseconds
}

export const useSessionTracking = (options: SessionTrackingOptions = {}) => {
    const {
        trackPageViews = true,
        trackUserActions = true,
        trackSessionQuality = true,
        sessionTimeout = 30 * 60 * 1000 // 30 minutes default
    } = options;

    const sessionDataRef = useRef<SessionData>({
        sessionId: generateSessionId(),
        startTime: Date.now(),
        pagesVisited: [],
        actionsPerformed: [],
        sessionQuality: 0,
        deviceInfo: {
            userAgent: typeof window !== 'undefined' ? window.navigator.userAgent : 'unknown',
            screenResolution: typeof window !== 'undefined' ? `${window.screen.width}x${window.screen.height}` : 'unknown',
            timezone: Intl.DateTimeFormat().resolvedOptions().timeZone,
            language: typeof window !== 'undefined' ? window.navigator.language : 'unknown'
        }
    });

    const lastActivityRef = useRef<number>(Date.now());
    const sessionTimeoutRef = useRef<NodeJS.Timeout | null>(null);

    // Generate unique session ID
    function generateSessionId(): string {
        return `session_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
    }

    // Track page view
    const trackPageView = useCallback((page: string) => {
        if (!trackPageViews) return;

        sessionDataRef.current.pagesVisited.push(page);
        lastActivityRef.current = Date.now();

        trackEvent('session_page_view', {
            session_id: sessionDataRef.current.sessionId,
            page,
            timestamp: Date.now(),
            page_sequence: sessionDataRef.current.pagesVisited.length
        });
    }, [trackPageViews]);

    // Track user action
    const trackUserAction = useCallback((action: string, metadata?: Record<string, unknown>) => {
        if (!trackUserActions) return;

        sessionDataRef.current.actionsPerformed.push(action);
        lastActivityRef.current = Date.now();

        trackEvent('session_user_action', {
            session_id: sessionDataRef.current.sessionId,
            action,
            timestamp: Date.now(),
            action_sequence: sessionDataRef.current.actionsPerformed.length,
            ...metadata
        });
    }, [trackUserActions]);

    // Calculate session quality score
    const calculateSessionQuality = useCallback(() => {
        if (!trackSessionQuality) return 0;

        const pagesVisited = sessionDataRef.current.pagesVisited.length;
        const actionsPerformed = sessionDataRef.current.actionsPerformed.length;
        const sessionDuration = Date.now() - sessionDataRef.current.startTime;

        // Simple quality score based on engagement
        let qualityScore = 0;

        // Base score for visiting pages
        qualityScore += Math.min(pagesVisited * 10, 50);

        // Bonus for performing actions
        qualityScore += Math.min(actionsPerformed * 5, 30);

        // Bonus for longer sessions (up to 20 points)
        const durationMinutes = sessionDuration / (1000 * 60);
        qualityScore += Math.min(durationMinutes * 2, 20);

        // Normalize to 0-100
        sessionDataRef.current.sessionQuality = Math.min(qualityScore, 100);

        return sessionDataRef.current.sessionQuality;
    }, [trackSessionQuality]);

    // Track session end
    const trackSessionEnd = useCallback(() => {
        const sessionDuration = Date.now() - sessionDataRef.current.startTime;
        const qualityScore = calculateSessionQuality();

        trackEvent('session_ended', {
            session_id: sessionDataRef.current.sessionId,
            duration: sessionDuration,
            pages_visited: sessionDataRef.current.pagesVisited.length,
            actions_performed: sessionDataRef.current.actionsPerformed.length,
            quality_score: qualityScore,
            device_info: sessionDataRef.current.deviceInfo
        });

        // Track session performance metrics
        trackPerformance({
            metric_name: 'session_duration',
            metric_value: sessionDuration
        });

        trackPerformance({
            metric_name: 'session_quality_score',
            metric_value: qualityScore
        });
    }, [calculateSessionQuality]);

    // Handle session timeout
    const handleSessionTimeout = useCallback(() => {
        trackSessionEnd();

        // Reset session data for new session
        sessionDataRef.current = {
            sessionId: generateSessionId(),
            startTime: Date.now(),
            pagesVisited: [],
            actionsPerformed: [],
            sessionQuality: 0,
            deviceInfo: sessionDataRef.current.deviceInfo
        };
    }, [trackSessionEnd]);

    // Set up session timeout
    useEffect(() => {
        const resetTimeout = () => {
            if (sessionTimeoutRef.current) {
                clearTimeout(sessionTimeoutRef.current);
            }

            sessionTimeoutRef.current = setTimeout(handleSessionTimeout, sessionTimeout);
        };

        // Reset timeout on activity
        const handleActivity = () => {
            lastActivityRef.current = Date.now();
            resetTimeout();
        };

        // Listen for user activity
        if (typeof window !== 'undefined') {
            window.addEventListener('mousedown', handleActivity);
            window.addEventListener('keydown', handleActivity);
            window.addEventListener('scroll', handleActivity);
            window.addEventListener('touchstart', handleActivity);
        }

        // Initial timeout setup
        resetTimeout();

        return () => {
            if (sessionTimeoutRef.current) {
                clearTimeout(sessionTimeoutRef.current);
            }

            if (typeof window !== 'undefined') {
                window.removeEventListener('mousedown', handleActivity);
                window.removeEventListener('keydown', handleActivity);
                window.removeEventListener('scroll', handleActivity);
                window.removeEventListener('touchstart', handleActivity);
            }
        };
    }, [sessionTimeout, handleSessionTimeout]);

    // Track session start
    useEffect(() => {
        trackEvent('session_started', {
            session_id: sessionDataRef.current.sessionId,
            start_time: sessionDataRef.current.startTime,
            device_info: sessionDataRef.current.deviceInfo
        });
    }, []);

    // Track session end on unmount
    useEffect(() => {
        return () => {
            trackSessionEnd();
        };
    }, [trackSessionEnd]);

    // Track page visibility changes
    useEffect(() => {
        if (typeof window === 'undefined') return;

        const handleVisibilityChange = () => {
            if (document.hidden) {
                trackEvent('session_paused', {
                    session_id: sessionDataRef.current.sessionId,
                    timestamp: Date.now()
                });
            } else {
                trackEvent('session_resumed', {
                    session_id: sessionDataRef.current.sessionId,
                    timestamp: Date.now()
                });
            }
        };

        document.addEventListener('visibilitychange', handleVisibilityChange);

        return () => {
            document.removeEventListener('visibilitychange', handleVisibilityChange);
        };
    }, []);

    return {
        trackPageView,
        trackUserAction,
        getSessionData: () => sessionDataRef.current,
        getSessionQuality: calculateSessionQuality
    };
};

// Hook for tracking specific session events
export const useSessionEventTracking = () => {
    const trackSessionEvent = useCallback((eventType: string, metadata?: Record<string, unknown>) => {
        trackEvent('session_event', {
            event_type: eventType,
            timestamp: Date.now(),
            ...metadata
        });
    }, []);

    const trackSessionError = useCallback((error: Error, context?: string) => {
        trackEvent('session_error', {
            error_message: error.message,
            error_stack: error.stack,
            context,
            timestamp: Date.now()
        });
    }, []);

    const trackSessionConversion = useCallback((conversionType: string, value?: number) => {
        trackEvent('session_conversion', {
            conversion_type: conversionType,
            value,
            timestamp: Date.now()
        });
    }, []);

    return {
        trackSessionEvent,
        trackSessionError,
        trackSessionConversion
    };
};
