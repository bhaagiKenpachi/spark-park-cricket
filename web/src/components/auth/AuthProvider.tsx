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
    console.log('=== AUTH PROVIDER: useEffect ===');
    console.log('isInitialized:', isInitialized);
    console.log('Document cookies:', document.cookie);
    console.log(
      'LocalStorage auth state:',
      localStorage.getItem('auth_authenticated')
    );

    if (!isInitialized) {
      console.log('Dispatching initializeAuth...');
      dispatch(initializeAuth());
    } else {
      console.log('Auth already initialized, skipping...');
    }
  }, [dispatch, isInitialized]);

  // Show loading state while initializing auth
  if (!isInitialized) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto mb-4"></div>
          <p className="text-gray-600">Loading...</p>
        </div>
      </div>
    );
  }

  return <>{children}</>;
}
