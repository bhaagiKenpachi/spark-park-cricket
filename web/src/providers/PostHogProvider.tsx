'use client';

import { Suspense, useEffect } from 'react';
import { usePathname, useSearchParams } from 'next/navigation';
import { initPostHog, posthog } from '@/lib/posthog';

interface PostHogProviderProps {
  children: React.ReactNode;
}

function PostHogPageView() {
  const pathname = usePathname();
  const searchParams = useSearchParams();

  // Track page views on route change
  useEffect(() => {
    if (pathname && typeof window !== 'undefined') {
      let url = window.origin + pathname;
      if (searchParams && searchParams.toString()) {
        url = url + `?${searchParams.toString()}`;
      }

      // Capture pageview
      if (posthog) {
        posthog.capture('$pageview', {
          $current_url: url,
        });
      }
    }
  }, [pathname, searchParams]);

  return null;
}

export function PostHogProvider({ children }: PostHogProviderProps) {
  // Initialize PostHog on mount
  useEffect(() => {
    initPostHog();
  }, []);

  return (
    <>
      <Suspense fallback={null}>
        <PostHogPageView />
      </Suspense>
      {children}
    </>
  );
}

