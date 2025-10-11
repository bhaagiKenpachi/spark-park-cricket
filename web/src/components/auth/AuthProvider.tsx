'use client';

import { useEffect } from 'react';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import { initializeAuth, checkAuthStatus } from '@/store/reducers/authSlice';

interface AuthProviderProps {
  children: React.ReactNode;
}

export function AuthProvider({ children }: AuthProviderProps) {
  const dispatch = useAppDispatch();
  const { isInitialized } = useAppSelector(state => state.auth);

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

  // Always render children, don't block the app for auth initialization
  // Auth can be initialized in the background
  return <>{children}</>;
}
