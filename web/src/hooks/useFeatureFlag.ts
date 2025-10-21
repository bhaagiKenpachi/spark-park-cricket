import { useState, useEffect } from 'react';
import { posthog } from '@/lib/posthog';
import { FeatureFlag } from '@/types/posthog';

/**
 * Custom hook to access PostHog feature flags
 * @param flagName - The name of the feature flag
 * @param defaultValue - Default value if flag is not available or loading
 * @returns Object containing flag value and loading state
 */
export function useFeatureFlag(
    flagName: FeatureFlag | string,
    defaultValue: boolean = false
): { enabled: boolean; loading: boolean } {
    const [enabled, setEnabled] = useState<boolean>(defaultValue);
    const [loading, setLoading] = useState<boolean>(true);

    useEffect(() => {
        // Only run on client side
        if (typeof window === 'undefined' || !posthog) {
            setEnabled(defaultValue);
            setLoading(false);
            return;
        }

        // Check if PostHog is loaded
        const checkFlag = () => {
            try {
                const flagValue = posthog.isFeatureEnabled(flagName);
                setEnabled(flagValue !== undefined ? flagValue : defaultValue);
                setLoading(false);
            } catch (error) {
                console.error('Error checking feature flag:', flagName, error);
                setEnabled(defaultValue);
                setLoading(false);
            }
        };

        // Initial check
        checkFlag();

        // Listen for feature flag changes
        const handleFlagChange = () => {
            checkFlag();
        };

        // PostHog emits 'feature_flags_received' event when flags are loaded
        posthog.onFeatureFlags?.(handleFlagChange);

        // Cleanup
        return () => {
            // PostHog doesn't have a removeListener, so nothing to cleanup
        };
    }, [flagName, defaultValue]);

    return { enabled, loading };
}

/**
 * Hook to get multiple feature flags at once
 * @param flags - Array of flag names
 * @returns Object with flag names as keys and their enabled state as values
 */
export function useFeatureFlags(
    flags: (FeatureFlag | string)[]
): Record<string, boolean> {
    const [flagValues, setFlagValues] = useState<Record<string, boolean>>({});

    useEffect(() => {
        if (typeof window === 'undefined' || !posthog) {
            return;
        }

        const checkFlags = () => {
            const values: Record<string, boolean> = {};
            flags.forEach(flag => {
                try {
                    const flagValue = posthog.isFeatureEnabled(flag);
                    values[flag] = flagValue !== undefined ? flagValue : false;
                } catch (error) {
                    console.error('Error checking feature flag:', flag, error);
                    values[flag] = false;
                }
            });
            setFlagValues(values);
        };

        checkFlags();
        posthog.onFeatureFlags?.(checkFlags);
    }, [flags]);

    return flagValues;
}

/**
 * Hook to get feature flag payload (for multivariate flags)
 * @param flagName - The name of the feature flag
 * @returns The flag payload or undefined
 */
export function useFeatureFlagPayload(
    flagName: FeatureFlag | string
): unknown | undefined {
    const [payload, setPayload] = useState<unknown | undefined>(undefined);

    useEffect(() => {
        if (typeof window === 'undefined' || !posthog) {
            return;
        }

        const checkPayload = () => {
            try {
                const flagPayload = posthog.getFeatureFlagPayload?.(flagName);
                setPayload(flagPayload);
            } catch (error) {
                console.error('Error getting feature flag payload:', flagName, error);
                setPayload(undefined);
            }
        };

        checkPayload();
        posthog.onFeatureFlags?.(checkPayload);
    }, [flagName]);

    return payload;
}


