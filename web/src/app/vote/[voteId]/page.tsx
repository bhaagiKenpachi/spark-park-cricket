'use client';

import { useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import { fetchVoteWithResultsRequest } from '@/store/reducers/voteSlice';
import { VoteView } from '@/components/VoteView';
import { Button } from '@/components/ui/button';
import { LoginButton } from '@/components/auth/LoginButton';
import { ArrowLeft, Home } from 'lucide-react';

export default function VoteDetailPage(): React.JSX.Element {
  const dispatch = useAppDispatch();
  const router = useRouter();
  const params = useParams();
  const { currentVote, loading, error } = useAppSelector(state => state.vote);
  const { user: currentUser, isAuthenticated } = useAppSelector(state => state.auth);
  const [localVote, setLocalVote] = useState<any>(null);

  const voteId = params.voteId as string;

  useEffect(() => {
    if (voteId) {
      console.log('Fetching vote with ID:', voteId);
      dispatch(fetchVoteWithResultsRequest(voteId));
    }
  }, [dispatch, voteId]);

  useEffect(() => {
    console.log('Current vote state:', currentVote);
    console.log('Loading state:', loading);
    console.log('Error state:', error);
    if (currentVote && currentVote.vote && currentVote.vote.id === voteId) {
      setLocalVote(currentVote);
    }
  }, [currentVote, voteId, loading, error]);

  const handleBack = () => {
    router.push('/votes');
  };

  const handleGoHome = () => {
    router.push('/');
  };

  if (loading && !currentVote) {
    return (
      <div className="min-h-screen bg-gray-50">
        {/* Top Navigation Bar */}
        <nav className="bg-white border-b sticky top-0 z-50 shadow-sm">
          <div className="w-full max-w-md mx-auto px-3 py-2">
            <div className="flex items-center justify-between">
              {/* Left: Back Button + Logo */}
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
                  <Home className="h-4 w-4 text-white" />
                </div>
                <span className="text-sm font-semibold text-gray-900">Spark Park</span>
              </div>

              {/* Right: Auth */}
              <div className="flex items-center gap-2">
                {isAuthenticated ? (
                  <span className="text-sm text-gray-600">{currentUser?.name}</span>
                ) : (
                  <LoginButton />
                )}
              </div>
            </div>
          </div>
        </nav>

        <main className="w-full max-w-md mx-auto">
          <div className="flex items-center justify-center py-8">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
            <span className="ml-2 text-sm text-gray-600">
              Loading vote...
            </span>
          </div>
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

  if (error) {
    return (
      <div className="min-h-screen bg-gray-50">
        {/* Top Navigation Bar */}
        <nav className="bg-white border-b sticky top-0 z-50 shadow-sm">
          <div className="w-full max-w-md mx-auto px-3 py-2">
            <div className="flex items-center justify-between">
              {/* Left: Back Button + Logo */}
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
                  <Home className="h-4 w-4 text-white" />
                </div>
                <span className="text-sm font-semibold text-gray-900">Spark Park</span>
              </div>

              {/* Right: Auth */}
              <div className="flex items-center gap-2">
                {isAuthenticated ? (
                  <span className="text-sm text-gray-600">{currentUser?.name}</span>
                ) : (
                  <LoginButton />
                )}
              </div>
            </div>
          </div>
        </nav>

        <main className="w-full max-w-md mx-auto">
          <div className="text-center py-8">
            <div className="bg-red-50 border border-red-200 rounded-lg p-6 mx-4">
              <h2 className="text-lg font-semibold text-red-800 mb-2">Error: Vote not found</h2>
              <p className="text-red-600 mb-4">
                The vote you're looking for doesn't exist or you don't have permission to view it.
              </p>
              <div className="flex gap-2 justify-center">
                <Button
                  variant="outline"
                  onClick={handleBack}
                  className="text-red-600 border-red-300 hover:bg-red-50"
                >
                  Go Back
                </Button>
                <Button
                  onClick={handleGoHome}
                  className="bg-blue-600 hover:bg-blue-700"
                >
                  Go Home
                </Button>
              </div>
            </div>
          </div>
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

  if (!localVote) {
    return (
      <div className="min-h-screen bg-gray-50">
        {/* Top Navigation Bar */}
        <nav className="bg-white border-b sticky top-0 z-50 shadow-sm">
          <div className="w-full max-w-md mx-auto px-3 py-2">
            <div className="flex items-center justify-between">
              {/* Left: Back Button + Logo */}
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
                  <Home className="h-4 w-4 text-white" />
                </div>
                <span className="text-sm font-semibold text-gray-900">Spark Park</span>
              </div>

              {/* Right: Auth */}
              <div className="flex items-center gap-2">
                {isAuthenticated ? (
                  <span className="text-sm text-gray-600">{currentUser?.name}</span>
                ) : (
                  <LoginButton />
                )}
              </div>
            </div>
          </div>
        </nav>

        <main className="w-full max-w-md mx-auto">
          <div className="text-center py-8">
            <p className="text-muted-foreground mb-4">
              No vote found for this ID.
            </p>
            <Button onClick={handleBack}>Back to Home</Button>
          </div>
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

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Top Navigation Bar */}
      <nav className="bg-white border-b sticky top-0 z-50 shadow-sm">
        <div className="w-full max-w-md mx-auto px-3 py-2">
          <div className="flex items-center justify-between">
            {/* Left: Back Button + Logo */}
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
                <Home className="h-4 w-4 text-white" />
              </div>
              <span className="text-sm font-semibold text-gray-900">Spark Park</span>
            </div>

            {/* Right: Auth */}
            <div className="flex items-center gap-2">
              {isAuthenticated ? (
                <span className="text-sm text-gray-600">{currentUser?.name}</span>
              ) : (
                <LoginButton />
              )}
            </div>
          </div>
        </div>
      </nav>

      <main className="w-full max-w-md mx-auto">
        <VoteView
          voteId={voteId}
          onBack={handleBack}
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
