'use client';

import { useParams, useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import { fetchMatchByIdRequest } from '@/store/reducers/matchSlice';
import { ScorecardView } from '@/components/ScorecardView';
import { Button } from '@/components/ui/button';
import { ArrowLeft, Home as HomeIcon } from 'lucide-react';

export default function MatchPage(): React.JSX.Element {
  const params = useParams();
  const router = useRouter();
  const dispatch = useAppDispatch();
  const { currentMatch, loading, error } = useAppSelector(state => state.match);
  const { user: currentUser, isAuthenticated } = useAppSelector(state => state.auth);
  const [localMatch, setLocalMatch] = useState<any>(null);

  const matchId = params.matchId as string;

  useEffect(() => {
    if (matchId) {
      dispatch(fetchMatchByIdRequest(matchId));
    }
  }, [dispatch, matchId]);

  useEffect(() => {
    if (currentMatch && currentMatch.id === matchId) {
      setLocalMatch(currentMatch);
    }
  }, [currentMatch, matchId, loading, error]);

  const handleBack = () => {
    router.push('/');
  };

  const handleGoHome = () => {
    router.push('/');
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50">
        {/* Navigation */}
        <nav className="bg-white border-b sticky top-0 z-50 shadow-sm">
          <div className="w-full max-w-md mx-auto px-3 py-2">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={handleBack}
                  className="h-8 px-2"
                >
                  <ArrowLeft className="h-4 w-4" />
                </Button>
                <div className="bg-blue-600 p-1.5 rounded-lg">
                  <HomeIcon className="h-4 w-4 text-white" />
                </div>
                <span className="text-sm font-semibold text-gray-900">Spark Park</span>
              </div>
              <Button
                variant="ghost"
                size="sm"
                onClick={handleGoHome}
                className="h-8 px-2"
              >
                Home
              </Button>
            </div>
          </div>
        </nav>

        <main className="w-full max-w-md mx-auto p-4">
          <div className="flex items-center justify-center py-8">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
            <span className="ml-2 text-sm text-gray-600">Loading match...</span>
          </div>
        </main>
      </div>
    );
  }

  if (error || !localMatch) {
    return (
      <div className="min-h-screen bg-gray-50">
        {/* Navigation */}
        <nav className="bg-white border-b sticky top-0 z-50 shadow-sm">
          <div className="w-full max-w-md mx-auto px-3 py-2">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={handleBack}
                  className="h-8 px-2"
                >
                  <ArrowLeft className="h-4 w-4" />
                </Button>
                <div className="bg-blue-600 p-1.5 rounded-lg">
                  <HomeIcon className="h-4 w-4 text-white" />
                </div>
                <span className="text-sm font-semibold text-gray-900">Spark Park</span>
              </div>
              <Button
                variant="ghost"
                size="sm"
                onClick={handleGoHome}
                className="h-8 px-2"
              >
                Home
              </Button>
            </div>
          </div>
        </nav>

        <main className="w-full max-w-md mx-auto p-4">
          <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded-lg">
            <strong className="font-bold">Error:</strong>
            <span className="block sm:inline"> {error || 'Match not found'}</span>
            <div className="mt-3 space-x-2">
              <Button
                onClick={handleBack}
                variant="outline"
                size="sm"
                className="bg-red-600 text-white border-red-600 hover:bg-red-700"
              >
                Go Back
              </Button>
              <Button
                onClick={handleGoHome}
                size="sm"
                className="bg-blue-600 text-white hover:bg-blue-700"
              >
                Go Home
              </Button>
            </div>
          </div>
        </main>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Navigation */}
      <nav className="bg-white border-b sticky top-0 z-50 shadow-sm">
        <div className="w-full max-w-md mx-auto px-3 py-2">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <Button
                variant="ghost"
                size="sm"
                onClick={handleBack}
                className="h-8 px-2"
              >
                <ArrowLeft className="h-4 w-4" />
              </Button>
              <div className="bg-blue-600 p-1.5 rounded-lg">
                <HomeIcon className="h-4 w-4 text-white" />
              </div>
              <span className="text-sm font-semibold text-gray-900">Spark Park</span>
            </div>
            <Button
              variant="ghost"
              size="sm"
              onClick={handleGoHome}
              className="h-8 px-2"
            >
            </Button>
          </div>
        </div>
      </nav>

      <main className="w-full max-w-md mx-auto">
        <ScorecardView
          matchId={localMatch.id}
          onBack={handleBack}
          currentUser={currentUser}
          isAuthenticated={isAuthenticated}
        />
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
