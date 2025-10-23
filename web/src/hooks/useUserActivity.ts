import { useEffect, useRef, useState, useCallback } from 'react';

// User Activity Tracking Hook
// Tracks mouse movement, clicks, keyboard events, and window focus to determine user activity

export interface UseUserActivityOptions {
    activityTimeout?: number; // Time in milliseconds after which user is considered inactive (default: 5 minutes)
    debounceMs?: number; // Debounce time for activity events (default: 1000ms)
}

export interface UserActivityState {
    isUserActive: boolean;
    lastActivityTime: Date | null;
    timeSinceLastActivity: number; // milliseconds since last activity
}

export function useUserActivity(options: UseUserActivityOptions = {}): UserActivityState {
    const {
        activityTimeout = 5 * 60 * 1000, // 5 minutes default
        debounceMs = 1000 // 1 second default debounce
    } = options;

    const [isUserActive, setIsUserActive] = useState(true);
    const [lastActivityTime, setLastActivityTime] = useState<Date | null>(new Date());
    const [timeSinceLastActivity, setTimeSinceLastActivity] = useState(0);

    const activityTimeoutRef = useRef<NodeJS.Timeout | null>(null);
    const debounceTimeoutRef = useRef<NodeJS.Timeout | null>(null);
    const isMountedRef = useRef(true);

    // Debounced activity handler
    const handleActivity = useCallback(() => {
        if (!isMountedRef.current) return;

        // Clear existing debounce timeout
        if (debounceTimeoutRef.current) {
            clearTimeout(debounceTimeoutRef.current);
        }

        // Set debounced activity update
        debounceTimeoutRef.current = setTimeout(() => {
            if (!isMountedRef.current) return;

            const now = new Date();
            setLastActivityTime(now);
            setIsUserActive(true);

            console.log('👤 USER ACTIVITY: User became active at', now.toISOString());

            // Clear existing activity timeout
            if (activityTimeoutRef.current) {
                clearTimeout(activityTimeoutRef.current);
            }

            // Set new activity timeout
            activityTimeoutRef.current = setTimeout(() => {
                if (!isMountedRef.current) return;

                setIsUserActive(false);
                console.log('👤 USER ACTIVITY: User became inactive after', activityTimeout / 1000, 'seconds');
            }, activityTimeout);

        }, debounceMs);
    }, [activityTimeout, debounceMs]);

    // Activity event handlers
    const handleMouseMove = useCallback(() => handleActivity(), [handleActivity]);
    const handleMouseClick = useCallback(() => handleActivity(), [handleActivity]);
    const handleKeyPress = useCallback(() => handleActivity(), [handleActivity]);
    const handleWindowFocus = useCallback(() => handleActivity(), [handleActivity]);
    const handleWindowBlur = useCallback(() => {
        // When window loses focus, mark as inactive immediately
        if (!isMountedRef.current) return;
        setIsUserActive(false);
        console.log('👤 USER ACTIVITY: User became inactive (window blur)');
    }, []);

    // Update time since last activity
    useEffect(() => {
        const updateTimeSinceLastActivity = () => {
            if (lastActivityTime) {
                const now = new Date();
                const timeDiff = now.getTime() - lastActivityTime.getTime();
                setTimeSinceLastActivity(timeDiff);
            }
        };

        // Update immediately
        updateTimeSinceLastActivity();

        // Update every 10 seconds
        const interval = setInterval(updateTimeSinceLastActivity, 10000);

        return () => clearInterval(interval);
    }, [lastActivityTime]);

    // Set up event listeners
    useEffect(() => {
        isMountedRef.current = true;

        // Add event listeners
        const events = [
            { target: document, event: 'mousemove', handler: handleMouseMove },
            { target: document, event: 'click', handler: handleMouseClick },
            { target: document, event: 'keydown', handler: handleKeyPress },
            { target: window, event: 'focus', handler: handleWindowFocus },
            { target: window, event: 'blur', handler: handleWindowBlur },
        ];

        events.forEach(({ target, event, handler }) => {
            target.addEventListener(event, handler);
        });

        // Set initial activity timeout
        activityTimeoutRef.current = setTimeout(() => {
            if (!isMountedRef.current) return;
            setIsUserActive(false);
            console.log('👤 USER ACTIVITY: Initial activity timeout reached');
        }, activityTimeout);

        // Cleanup function
        return () => {
            isMountedRef.current = false;

            // Remove event listeners
            events.forEach(({ target, event, handler }) => {
                target.removeEventListener(event, handler);
            });

            // Clear timeouts
            if (activityTimeoutRef.current) {
                clearTimeout(activityTimeoutRef.current);
            }
            if (debounceTimeoutRef.current) {
                clearTimeout(debounceTimeoutRef.current);
            }
        };
    }, [handleMouseMove, handleMouseClick, handleKeyPress, handleWindowFocus, handleWindowBlur, activityTimeout]);

    return {
        isUserActive,
        lastActivityTime,
        timeSinceLastActivity,
    };
}
