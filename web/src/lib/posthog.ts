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

            // Session Recording Configuration
            session_recording: {
                recordCrossOriginIframes: true,
                maskAllInputs: false,
                maskInputOptions: {
                    password: true,
                },
            },

            // Auto-capture configuration
            capture_pageview: true,
            capture_pageleave: true,
            autocapture: true,

            // Performance and privacy settings
            persistence: 'localStorage+cookie',
            disable_session_recording: false,

            // Advanced options
            sanitize_properties: (properties) => {
                // Remove sensitive data
                const sanitized = { ...properties };
                if (sanitized.password) delete sanitized.password;
                if (sanitized.token) delete sanitized.token;
                if (sanitized.api_key) delete sanitized.api_key;
                return sanitized;
            },

            // Debug mode based on environment
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


