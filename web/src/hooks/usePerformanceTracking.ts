import { useEffect, useRef, useCallback } from 'react';
import { trackPerformance } from '@/lib/analytics';

interface PerformanceTrackingOptions {
    componentName?: string;
    trackMountTime?: boolean;
    trackRenderTime?: boolean;
    trackUnmountTime?: boolean;
    trackUserInteractions?: boolean;
}

export const usePerformanceTracking = (options: PerformanceTrackingOptions = {}) => {
    const {
        componentName = 'UnknownComponent',
        trackMountTime = true,
        trackRenderTime = true,
        trackUnmountTime = true,
        trackUserInteractions = false
    } = options;

    const mountTimeRef = useRef<number>(0);
    const renderStartTimeRef = useRef<number>(0);
    const renderCountRef = useRef<number>(0);

    // Track component mount time
    useEffect(() => {
        if (trackMountTime) {
            mountTimeRef.current = performance.now();
            trackPerformance({
                metric_name: 'component_mount_time',
                metric_value: mountTimeRef.current,
                component: componentName
            });
        }

        // Track component unmount time
        return () => {
            if (trackUnmountTime && mountTimeRef.current > 0) {
                const unmountTime = performance.now();
                const totalLifetime = unmountTime - mountTimeRef.current;
                trackPerformance({
                    metric_name: 'component_lifetime',
                    metric_value: totalLifetime,
                    component: componentName
                });
            }
        };
    }, [componentName, trackMountTime, trackUnmountTime]);

    // Track render time
    useEffect(() => {
        if (trackRenderTime) {
            renderStartTimeRef.current = performance.now();
            renderCountRef.current += 1;

            // Use requestAnimationFrame to measure after render
            requestAnimationFrame(() => {
                const renderEndTime = performance.now();
                const renderTime = renderEndTime - renderStartTimeRef.current;

                trackPerformance({
                    metric_name: 'component_render_time',
                    metric_value: renderTime,
                    component: componentName
                });

                // Track render count
                trackPerformance({
                    metric_name: 'component_render_count',
                    metric_value: renderCountRef.current,
                    component: componentName
                });
            });
        }
    });

    // Track user interactions
    const trackUserInteraction = useCallback((interactionType: string, metadata?: Record<string, unknown>) => {
        if (trackUserInteractions) {
            trackPerformance({
                metric_name: 'user_interaction',
                metric_value: performance.now(),
                component: componentName,
                ...metadata
            });
        }
    }, [componentName, trackUserInteractions]);

    // Track custom performance metrics
    const trackCustomMetric = useCallback((metricName: string, value: number, metadata?: Record<string, unknown>) => {
        trackPerformance({
            metric_name: metricName,
            metric_value: value,
            component: componentName,
            ...metadata
        });
    }, [componentName]);

    return {
        trackUserInteraction,
        trackCustomMetric
    };
};

// Hook for tracking specific user interactions
export const useInteractionTracking = (componentName: string) => {
    const trackClick = useCallback((element?: string) => {
        trackPerformance({
            metric_name: 'user_click',
            metric_value: performance.now(),
            component: componentName
        });
    }, [componentName]);

    const trackHover = useCallback((element?: string) => {
        trackPerformance({
            metric_name: 'user_hover',
            metric_value: performance.now(),
            component: componentName
        });
    }, [componentName]);

    const trackFocus = useCallback((element?: string) => {
        trackPerformance({
            metric_name: 'user_focus',
            metric_value: performance.now(),
            component: componentName
        });
    }, [componentName]);

    const trackScroll = useCallback((element?: string) => {
        trackPerformance({
            metric_name: 'user_scroll',
            metric_value: performance.now(),
            component: componentName
        });
    }, [componentName]);

    return {
        trackClick,
        trackHover,
        trackFocus,
        trackScroll
    };
};

// Hook for tracking time on page
export const useTimeOnPageTracking = (pageName: string) => {
    const startTimeRef = useRef<number>(0);

    useEffect(() => {
        startTimeRef.current = performance.now();

        return () => {
            if (startTimeRef.current > 0) {
                const timeOnPage = performance.now() - startTimeRef.current;
                trackPerformance({
                    metric_name: 'time_on_page',
                    metric_value: timeOnPage,
                    page: pageName
                });
            }
        };
    }, [pageName]);
};
