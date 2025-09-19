'use client';

import { useState, useEffect } from 'react';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import {
  fetchSeriesRequest,
  deleteSeriesRequest,
  Series,
} from '@/store/reducers/seriesSlice';
import { SeriesForm } from './SeriesForm';
import { SeriesWithMatches } from './SeriesWithMatches';
import { ScorecardView } from './ScorecardView';
import { Button } from '@/components/ui/button';
import { RefreshCw, Plus } from 'lucide-react';
import { User } from '@/services/authService';

export function SeriesList(): React.JSX.Element {
  const dispatch = useAppDispatch();
  const { series, loading, error } = useAppSelector(state => state.series);
  const { user: currentUser, isAuthenticated } = useAppSelector(
    state => state.auth
  );
  const [showForm, setShowForm] = useState(false);
  const [editingSeries, setEditingSeries] = useState<Series | undefined>();
  const [viewingScorecard, setViewingScorecard] = useState<string | null>(null);
  const [currentSeriesCreatedBy, setCurrentSeriesCreatedBy] = useState<
    string | null
  >(null);

  useEffect(() => {
    dispatch(fetchSeriesRequest());
  }, [dispatch]);

  const handleCreateSeries = () => {
    if (!isAuthenticated) {
      // Find the sign-in button and add a blinking red border effect
      const signInButton = document.querySelector(
        '[data-cy="login-button"]'
      ) as HTMLElement;
      if (signInButton) {
        console.log('Sign-in button found:', signInButton);
        console.log('Starting blink effect...');

        // Focus and scroll to the button
        signInButton.focus();
        signInButton.scrollIntoView({ behavior: 'smooth', block: 'center' });

        // Add blinking red border effect
        let blinkCount = 0;
        const blinkInterval = setInterval(() => {
          blinkCount++;
          console.log(`Blink ${blinkCount}/6`);

          if (blinkCount % 2 === 1) {
            // Red border ON
            console.log('Red border ON');
            signInButton.style.setProperty(
              'border',
              '2px solid red',
              'important'
            );
            signInButton.style.setProperty(
              'box-shadow',
              '0 0 10px rgba(255, 0, 0, 0.5)',
              'important'
            );
            signInButton.style.setProperty(
              'background-color',
              '#fee2e2',
              'important'
            );
          } else {
            // Red border OFF
            console.log('Red border OFF');
            signInButton.style.removeProperty('border');
            signInButton.style.removeProperty('box-shadow');
            signInButton.style.removeProperty('background-color');
          }

          if (blinkCount >= 6) {
            clearInterval(blinkInterval);
            console.log('Blink effect completed');
            // Clean up any remaining styles
            signInButton.style.removeProperty('border');
            signInButton.style.removeProperty('box-shadow');
            signInButton.style.removeProperty('background-color');
          }
        }, 500);
      } else {
        console.log('Sign-in button not found');
      }
      return;
    }
    setShowForm(true);
  };

  const handleDelete = (id: string) => {
    if (window.confirm('Are you sure you want to delete this series?')) {
      dispatch(deleteSeriesRequest(id));
    }
  };

  const handleEdit = (series: Series) => {
    setEditingSeries(series);
    setShowForm(true);
  };

  const handleFormSuccess = () => {
    setShowForm(false);
    setEditingSeries(undefined);
    dispatch(fetchSeriesRequest());
  };

  const handleFormCancel = () => {
    setShowForm(false);
    setEditingSeries(undefined);
  };

  const handleViewScorecard = (matchId: string, seriesCreatedBy: string) => {
    console.log('=== HANDLE VIEW SCORECARD ===');
    console.log('Match ID:', matchId);
    console.log('Series Created By:', seriesCreatedBy);
    setViewingScorecard(matchId);
    setCurrentSeriesCreatedBy(seriesCreatedBy);
  };

  const handleBackFromScorecard = () => {
    setViewingScorecard(null);
    setCurrentSeriesCreatedBy(null);
  };

  if (loading && (!series || !Array.isArray(series) || series.length === 0)) {
    return (
      <div
        className="
        w-full max-w-sm mx-auto px-4 py-8
        sm:max-w-md sm:px-6
        md:max-w-lg md:px-8
      "
      >
        <div className="flex flex-col items-center justify-center space-y-4">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
          <span className="text-sm text-gray-600">Loading series...</span>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div
        className="
        w-full max-w-sm mx-auto px-4 py-6
        sm:max-w-md sm:px-6
        md:max-w-lg md:px-8
      "
      >
        <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded-lg">
          <strong className="font-bold">Error:</strong>
          <span className="block sm:inline"> {error}</span>
          <div className="mt-3">
            <button
              onClick={() => dispatch(fetchSeriesRequest())}
              className="
                                w-full py-2 px-4 bg-red-600 text-white rounded-lg font-medium
                                active:bg-red-700 focus:outline-none focus:ring-2 focus:ring-red-500
                                sm:w-auto sm:px-6
                            "
            >
              Retry
            </button>
          </div>
        </div>
      </div>
    );
  }

  if (viewingScorecard) {
    return (
      <ScorecardView
        matchId={viewingScorecard}
        onBack={handleBackFromScorecard}
        {...(currentSeriesCreatedBy && {
          seriesCreatedBy: currentSeriesCreatedBy,
        })}
        currentUser={currentUser}
        isAuthenticated={isAuthenticated}
      />
    );
  }

  if (showForm) {
    return (
      <SeriesForm
        series={editingSeries || undefined}
        onSuccess={handleFormSuccess}
        onCancel={handleFormCancel}
      />
    );
  }

  return (
    <div className="w-full max-w-4xl mx-auto p-6" data-cy="series-list">
      <div className="flex flex-col items-center space-y-4 mb-6 sm:flex-row sm:justify-between sm:space-y-0">
        <h2 className="text-2xl font-bold text-center">Cricket Series</h2>
        <div className="flex space-x-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => dispatch(fetchSeriesRequest())}
            disabled={loading}
            data-cy="refresh-series-button"
            title="Refresh"
          >
            <RefreshCw className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
          </Button>
          <Button
            onClick={handleCreateSeries}
            data-cy="create-series-button"
            title="Create Series"
          >
            <Plus className="h-4 w-4 mr-2" />
            Series
          </Button>
        </div>
      </div>

      {!series || !Array.isArray(series) || series.length === 0 ? (
        <div className="text-center py-8">
          <p className="text-muted-foreground mb-4">No series found.</p>
          <Button
            onClick={handleCreateSeries}
            data-cy="create-first-series-button"
            title="Create Your First Series"
          >
            <Plus className="h-4 w-4 mr-2" />
            Your First Series
          </Button>
        </div>
      ) : (
        <div className="space-y-4">
          {Array.isArray(series) &&
            series.map((seriesItem, index) => (
              <SeriesWithMatches
                key={seriesItem.id || `series-${index}`}
                series={seriesItem}
                onEditSeries={handleEdit}
                onDeleteSeries={handleDelete}
                onViewScorecard={handleViewScorecard}
                currentUser={currentUser}
                isAuthenticated={isAuthenticated}
              />
            ))}
        </div>
      )}
    </div>
  );
}
