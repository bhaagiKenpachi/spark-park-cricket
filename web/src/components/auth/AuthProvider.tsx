'use client';

import { useEffect } from 'react';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import { initializeAuth } from '@/store/reducers/authSlice';

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

  // Always render children, don't block the app for auth initialization
  // Auth can be initialized in the background
  return <>{children}</>;
}
