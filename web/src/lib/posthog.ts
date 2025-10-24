import posthog from 'posthog-js';

let isInitialized = false;

export const initPostHog = (): void => {
    // Only run on client side
    if (typeof window === 'undefined') {
        return;
    }

    // Prevent multiple initializations
    if (isInitialized) {
        return;
    }

    const posthogKey = process.env.NEXT_PUBLIC_POSTHOG_KEY;
    const posthogHost = process.env.NEXT_PUBLIC_POSTHOG_HOST;

    // Don't initialize if no key is provided
    if (!posthogKey || !posthogHost) {
        console.warn('PostHog not initialized: Missing API key or host');
        return;
    }

    try {

        posthog.init(posthogKey, {
            api_host: posthogHost,
            capture_pageview: true,
            capture_pageleave: true,
            autocapture: true,
            persistence: 'localStorage+cookie',
            disable_session_recording: false,
            loaded: (posthogInstance) => {
                const debug = process.env.NEXT_PUBLIC_POSTHOG_DEBUG === 'true';
                if (debug && process.env.NODE_ENV === 'development') {
                    posthogInstance.debug();
                }
            },
        });

        isInitialized = true;
        console.log('PostHog initialized successfully');
    } catch (error) {
        console.error('Failed to initialize PostHog:', error);
    }
};

export const getPostHog = () => {
    if (typeof window === 'undefined') {
        return null;
    }
    return posthog;
};

export { posthog };


