'use client';

import { useEffect } from 'react';
import { usePathname, useSearchParams } from 'next/navigation';
import { posthog } from '@/lib/posthog';
import { trackPerformance } from '@/lib/analytics';

export function PostHogPageView() {
  const pathname = usePathname();
  const searchParams = useSearchParams();

  // Track page views on route change
  useEffect(() => {
    if (pathname) {
      let url = window.origin + pathname;
      if (searchParams && searchParams.toString()) {
        url = url + `?${searchParams.toString()}`;
      }

      // Capture pageview (allowed in dev when override flag is enabled)
      const allowInDev = process.env.NODE_ENV !== 'development' || process.env.NEXT_PUBLIC_POSTHOG_ENABLE_IN_DEV === 'true';
      if (allowInDev && posthog) {
        posthog.capture('$pageview', {
          $current_url: url,
        });
      }
    }
  }, [pathname, searchParams]);

  // Track page performance metrics
  useEffect(() => {
    if ('performance' in window) {
      const navigation = performance.getEntriesByType('navigation')[0] as PerformanceNavigationTiming;
      if (navigation) {
        // Track page load time
        const pageLoadTime = navigation.loadEventEnd - navigation.fetchStart;
        trackPerformance({
          metric_name: 'page_load_time',
          metric_value: pageLoadTime,
          page: pathname
        });

        // Track Time to First Byte (TTFB)
        const ttfb = navigation.responseStart - navigation.fetchStart;
        trackPerformance({
          metric_name: 'time_to_first_byte',
          metric_value: ttfb,
          page: pathname
        });

        // Track First Contentful Paint (FCP)
        const fcpEntries = performance.getEntriesByName('first-contentful-paint');
        if (fcpEntries.length > 0) {
          trackPerformance({
            metric_name: 'first_contentful_paint',
            metric_value: fcpEntries[0]?.startTime || 0,
            page: pathname
          });
        }

        // Track Largest Contentful Paint (LCP)
        const lcpEntries = performance.getEntriesByType('largest-contentful-paint');
        if (lcpEntries.length > 0) {
          trackPerformance({
            metric_name: 'largest_contentful_paint',
            metric_value: lcpEntries[lcpEntries.length - 1]?.startTime || 0,
            page: pathname
          });
        }
      }
    }
  }, [pathname]);

  return null;
}
