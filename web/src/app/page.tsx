'use client';

import { useEffect } from 'react';
import Link from 'next/link';
import { SeriesList } from '@/components/SeriesList';
import { LoginButton } from '@/components/auth/LoginButton';
import { UserMenu } from '@/components/auth/UserMenu';
import { useAppSelector, useAppDispatch } from '@/store/hooks';
import { checkAuthStatus } from '@/store/reducers/authSlice';
import { Button } from '@/components/ui/button';
import { Vote, Trophy, Home as HomeIcon } from 'lucide-react';

export default function Home(): React.JSX.Element {
  const { isAuthenticated } = useAppSelector(state => state.auth);
  const dispatch = useAppDispatch();

  // Handle authentication success callback
  useEffect(() => {
    const urlParams = new URLSearchParams(window.location.search);

    if (urlParams.get('auth') === 'success') {
      // Clear the URL parameter
      window.history.replaceState({}, document.title, window.location.pathname);
      // Check authentication status
      dispatch(checkAuthStatus());
    }
  }, [dispatch]);

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Top Navigation Bar */}
      <nav className="bg-white border-b sticky top-0 z-50 shadow-sm">
        <div className="w-full max-w-md mx-auto px-3 py-2">
          <div className="flex items-center justify-between">
            {/* Left: Logo/Home */}
            <div className="flex items-center gap-2">
              <div className="bg-blue-600 p-1.5 rounded-lg">
                <HomeIcon className="h-4 w-4 text-white" />
              </div>
              <span className="text-sm font-semibold text-gray-900">Spark Park</span>
            </div>

            {/* Right: Navigation + Auth */}
            <div className="flex items-center gap-2">
              <Link href="/votes">
                <Button variant="ghost" size="sm" className="h-8 px-2.5 flex items-center gap-1.5">
                  <Vote className="h-3.5 w-3.5" />
                  <span className="text-xs font-medium">Votes</span>
                </Button>
              </Link>
              {isAuthenticated ? <UserMenu /> : <LoginButton />}
            </div>
          </div>
        </div>
      </nav>

      {/* Sub Navigation Bar */}
      <div className="bg-white border-b">
        <div className="w-full max-w-md mx-auto px-3 py-3">
          <div className="flex items-center gap-2">
            <div className="bg-green-100 p-1.5 rounded-lg">
              <Trophy className="h-4 w-4 text-green-600" />
            </div>
            <div>
              <h1 className="text-base font-bold text-gray-900">Cricket Series</h1>
              <p className="text-xs text-gray-500">Tournaments & Matches</p>
            </div>
          </div>
        </div>
      </div>

      <main className="w-full max-w-md mx-auto px-4 py-4">
        <SeriesList />
      </main>

      <footer className="bg-white border-t">
        <div className="w-full max-w-md mx-auto py-3 px-4">
          <p className="text-center text-gray-500 text-xs">
            © 2024 Spark Park Cricket. All rights reserved.
          </p>
        </div>
      </footer>
    </div>
  );
}

