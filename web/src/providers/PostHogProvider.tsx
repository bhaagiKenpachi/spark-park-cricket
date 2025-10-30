'use client';

import { useEffect } from 'react';
import dynamic from 'next/dynamic';
import { initPostHog } from '@/lib/posthog';

interface PostHogProviderProps {
  children: React.ReactNode;
}

// Dynamically import PostHogPageView to prevent SSR issues
const PostHogPageView = dynamic(
  () => import('./PostHogPageView').then((mod) => ({ default: mod.PostHogPageView })),
  {
    ssr: false,
    loading: () => null,
  }
);

export function PostHogProvider({ children }: PostHogProviderProps) {
  // Initialize PostHog on mount
  useEffect(() => {
    initPostHog();
  }, []);

  return (
    <>
      <PostHogPageView />
      {children}
    </>
  );
}

