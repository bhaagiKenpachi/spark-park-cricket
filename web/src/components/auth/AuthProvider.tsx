'use client';

import { useEffect } from 'react';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import { initializeAuth, checkAuthStatus } from '@/store/reducers/authSlice';
import { identifyUser, resetUser, trackUserLogin, trackUserLogout } from '@/lib/analytics';

interface AuthProviderProps {
  children: React.ReactNode;
}

export function AuthProvider({ children }: AuthProviderProps) {
  const dispatch = useAppDispatch();
  const { isInitialized, user, isAuthenticated } = useAppSelector(state => state.auth);

  useEffect(() => {
    if (!isInitialized) {
      dispatch(initializeAuth()).catch(error => {
        console.error('Auth initialization failed:', error);
        // Don't block the app if auth initialization fails
      });
    }
  }, [dispatch, isInitialized]);

  // Handle authentication success callback globally
  useEffect(() => {
    if (typeof window !== 'undefined') {
      const urlParams = new URLSearchParams(window.location.search);

      if (urlParams.get('auth') === 'success') {
        // Clear the auth parameter from URL
        const newUrl = new URL(window.location.href);
        newUrl.searchParams.delete('auth');
        window.history.replaceState({}, document.title, newUrl.pathname + newUrl.search);

        // Check authentication status
        dispatch(checkAuthStatus());
      }
    }
  }, [dispatch]);

  // Track user authentication state changes
  useEffect(() => {
    if (isAuthenticated && user) {
      // Calculate user properties
      const userProperties = {
        id: user.id,
        email: user.email,
        name: user.name,
        created_at: user.created_at,
        // Add calculated properties
        total_series_created: 0, // This would need to be calculated from user data
        total_matches_scored: 0, // This would need to be calculated from user data
        user_role: 'standard', // Default role, could be enhanced with actual roles
        signup_date: user.created_at,
        last_active: new Date().toISOString(),
        // Add session information
        session_start: new Date().toISOString(),
        timezone: Intl.DateTimeFormat().resolvedOptions().timeZone,
        user_agent: typeof window !== 'undefined' ? window.navigator.userAgent : 'unknown',
        screen_resolution: typeof window !== 'undefined' ? `${window.screen.width}x${window.screen.height}` : 'unknown',
        language: typeof window !== 'undefined' ? window.navigator.language : 'unknown'
      };

      // User logged in - identify them in PostHog
      identifyUser(user.id, userProperties);

      // Track login event
      trackUserLogin({
        user_id: user.id,
        email: user.email,
        name: user.name,
        provider: 'google', // Currently only Google OAuth is supported
      });
    } else if (!isAuthenticated && isInitialized) {
      // User logged out - reset PostHog user
      resetUser();
    }
  }, [isAuthenticated, user, isInitialized]);

  // Always render children, don't block the app for auth initialization
  // Auth can be initialized in the background
  return <>{children}</>;
}
