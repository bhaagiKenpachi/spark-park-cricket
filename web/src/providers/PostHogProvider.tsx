'use client';

import { useEffect } from 'react';
import { initPostHog } from '@/lib/posthog';

interface PostHogProviderProps {
  children: React.ReactNode;
}

export function PostHogProvider({ children }: PostHogProviderProps) {
  // Initialize PostHog on mount (only on client side)
  useEffect(() => {
    if (typeof window !== 'undefined') {
      initPostHog().catch((error) => {
        console.error('Failed to initialize PostHog:', error);
      });
    }
  }, []);

  return <>{children}</>;
}

