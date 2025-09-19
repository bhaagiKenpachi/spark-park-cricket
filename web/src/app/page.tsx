'use client';

import { useEffect } from 'react';
import { SeriesList } from '@/components/SeriesList';
import { LoginButton } from '@/components/auth/LoginButton';
import { UserMenu } from '@/components/auth/UserMenu';
import { useAppSelector, useAppDispatch } from '@/store/hooks';
import { checkAuthStatus } from '@/store/reducers/authSlice';

export default function Home(): React.JSX.Element {
  const { isAuthenticated, user, isLoading, error, isInitialized } =
    useAppSelector(state => state.auth);
  const dispatch = useAppDispatch();

  // Debug auth state
  console.log('=== MAIN PAGE AUTH STATE ===');
  console.log('isAuthenticated:', isAuthenticated);
  console.log('user:', user);
  console.log('isLoading:', isLoading);
  console.log('error:', error);
  console.log('isInitialized:', isInitialized);

  // Handle authentication success callback
  useEffect(() => {
    console.log('=== FRONTEND PAGE LOAD ===');
    console.log('Current URL:', window.location.href);
    console.log('Document cookies:', document.cookie);
    console.log(
      'LocalStorage auth state:',
      localStorage.getItem('auth_authenticated')
    );
    console.log('LocalStorage user:', localStorage.getItem('auth_user'));

    const urlParams = new URLSearchParams(window.location.search);
    console.log('URL parameters:', Object.fromEntries(urlParams.entries()));

    if (urlParams.get('auth') === 'success') {
      console.log('=== OAUTH SUCCESS CALLBACK ===');
      console.log('Clearing URL parameter and checking auth status');

      // Clear the URL parameter
      window.history.replaceState({}, document.title, window.location.pathname);
      // Check authentication status
      dispatch(checkAuthStatus());
    }
  }, [dispatch]);

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white shadow-sm border-b">
        <div
          className="
          w-full max-w-sm mx-auto px-4 py-4
          sm:max-w-md sm:px-6 sm:py-5
          md:max-w-lg md:px-8 md:py-6
        "
        >
          <div
            className="flex flex-col items-center space-y-4
            sm:flex-row sm:justify-between sm:space-y-0
          "
          >
            <h1
              className="
              text-xl font-bold text-gray-900 text-center
              sm:text-2xl md:text-3xl
            "
            >
              Spark Park Cricket
            </h1>

            {/* Authentication Section */}
            <div className="flex items-center">
              {isAuthenticated ? <UserMenu /> : <LoginButton />}
            </div>
          </div>
        </div>
      </header>

      <main
        className="
        w-full max-w-sm mx-auto px-4 py-6
        sm:max-w-md sm:px-6 sm:py-8
        md:max-w-lg md:px-8 md:py-10
        lg:max-w-xl lg:py-12
      "
      >
        <div className="mb-8 text-center">
          <h2
            className="
            text-xl font-bold text-gray-900 mb-4
            sm:text-2xl md:text-3xl
          "
          >
            Welcome to Spark Park Cricket
          </h2>
          <p
            className="
            text-sm text-gray-600
            sm:text-base md:text-lg
          "
          >
            Manage your cricket tournaments, matches, and teams with our
            comprehensive tournament management system.
          </p>
        </div>

        <SeriesList />
      </main>

      <footer className="bg-white border-t">
        <div
          className="
          w-full max-w-sm mx-auto py-4 px-4
          sm:max-w-md sm:px-6 sm:py-6
          md:max-w-lg md:px-8
        "
        >
          <p
            className="text-center text-gray-500 text-xs
            sm:text-sm
          "
          >
            © 2024 Spark Park Cricket. All rights reserved.
          </p>
        </div>
      </footer>
    </div>
  );
}
