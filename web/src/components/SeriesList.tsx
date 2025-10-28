'use client';

import { useState, useEffect } from 'react';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import {
  fetchSeriesRequest,
  deleteSeriesRequest,
  setPage,
  setPageSize,
  Series,
} from '@/store/reducers/seriesSlice';
import { SeriesForm } from './SeriesForm';
import { SeriesWithMatches } from './SeriesWithMatches';
import { ScorecardView } from './ScorecardView';
import { Button } from '@/components/ui/button';
import { Pagination } from '@/components/ui/pagination';
import { RefreshCw, Plus } from 'lucide-react';
import { InFeedAd } from '@/components/ads/InFeedAd';
import {
  trackSeriesViewed,
  trackSeriesDeleted,
  trackSeriesPaginationChanged,
} from '@/lib/analytics';

export function SeriesList(): React.JSX.Element {
  const dispatch = useAppDispatch();
  const { series, loading, error, pagination } = useAppSelector(
    state => state.series
  );
  const { user: currentUser, isAuthenticated } = useAppSelector(
    state => state.auth
  );
  const [showForm, setShowForm] = useState(false);
  const [editingSeries, setEditingSeries] = useState<Series | undefined>();
  const [viewingScorecard, setViewingScorecard] = useState<string | null>(null);
  const [currentSeriesCreatedBy, setCurrentSeriesCreatedBy] = useState<
    string | null
  >(null);
  const [expandedSeriesId, setExpandedSeriesId] = useState<string | null>(null);

  // Fetch series data on component mount
  useEffect(() => {
    dispatch(
      fetchSeriesRequest({
        page: pagination.currentPage,
        pageSize: pagination.pageSize,
      })
    );
  }, [dispatch, pagination.currentPage, pagination.pageSize]);

  // Track series view when component mounts or series data changes
  useEffect(() => {
    if (series && series.length > 0) {
      series.forEach(seriesItem => {
        trackSeriesViewed({
          series_id: seriesItem.id,
          series_name: seriesItem.name,
          start_date: seriesItem.start_date,
          end_date: seriesItem.end_date,
        });
      });
    }
  }, [series]);

  useEffect(() => {
  // Only run when not viewing scorecard
  if (!viewingScorecard) {
    const scrollId = sessionStorage.getItem('scrollToSeriesId');
    if (scrollId) {
      // Use setTimeout to ensure DOM is ready
      setTimeout(() => {
        const el = document.getElementById(`series-card-${scrollId}`);
        if (el) {
          el.scrollIntoView({ behavior: 'smooth', block: 'center' });
        }
        sessionStorage.removeItem('scrollToSeriesId');
      }, 100);
    }
  }
}, [series, viewingScorecard]);

  const handleDelete = (id: string) => {
    if (window.confirm('Are you sure you want to delete this series?')) {
      // Track deletion
      const seriesToDelete = series.find(s => s.id === id);
      if (seriesToDelete) {
        trackSeriesDeleted({
          series_id: id,
          series_name: seriesToDelete.name,
        });
      }
      dispatch(deleteSeriesRequest(id));
    }
  };

  const handleEdit = (series: Series) => {
    if (!isAuthenticated) {
      // Find the sign-in button and add a blinking red border effect
      const signInButton = document.querySelector(
        '[data-cy="login-button"]'
      ) as HTMLElement;
      if (signInButton) {
        // Focus and scroll to the button
        signInButton.focus();
        signInButton.scrollIntoView({ behavior: 'smooth', block: 'center' });

        // Add blinking red border effect
        let blinkCount = 0;
        const blinkInterval = setInterval(() => {
          blinkCount++;

          if (blinkCount % 2 === 1) {
            // Red border ON
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
            signInButton.style.removeProperty('border');
            signInButton.style.removeProperty('box-shadow');
            signInButton.style.removeProperty('background-color');
          }

          if (blinkCount >= 6) {
            clearInterval(blinkInterval);
            // Clean up any remaining styles
            signInButton.style.removeProperty('border');
            signInButton.style.removeProperty('box-shadow');
            signInButton.style.removeProperty('background-color');
          }
        }, 500);
      } else {
      }
      return;
    }
    setEditingSeries(series);
    setShowForm(true);
  };

  const handleCreateSeries = () => {
    if (!isAuthenticated) {
      // Find the sign-in button and add a blinking red border effect
      const signInButton = document.querySelector(
        '[data-cy="login-button"]'
      ) as HTMLElement;
      if (signInButton) {
        // Focus and scroll to the button
        signInButton.focus();
        signInButton.scrollIntoView({ behavior: 'smooth', block: 'center' });

        // Add blinking red border effect
        let blinkCount = 0;
        const blinkInterval = setInterval(() => {
          blinkCount++;

          if (blinkCount % 2 === 1) {
            // Red border ON
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
            signInButton.style.removeProperty('border');
            signInButton.style.removeProperty('box-shadow');
            signInButton.style.removeProperty('background-color');
          }

          if (blinkCount >= 6) {
            clearInterval(blinkInterval);
            // Clean up any remaining styles
            signInButton.style.removeProperty('border');
            signInButton.style.removeProperty('box-shadow');
            signInButton.style.removeProperty('background-color');
          }
        }, 500);
      } else {
      }
      return;
    }
    setShowForm(true);
  };

  const handleFormSuccess = () => {
    setShowForm(false);
    setEditingSeries(undefined);
    dispatch(
      fetchSeriesRequest({
        page: pagination.currentPage,
        pageSize: pagination.pageSize,
      })
    );
  };

  const handlePageChange = (page: number) => {
    trackSeriesPaginationChanged({
      page,
      page_size: pagination.pageSize,
      total_items: pagination.totalItems,
    });
    dispatch(setPage(page));
  };

  const handlePageSizeChange = (pageSize: number) => {
    trackSeriesPaginationChanged({
      page: pagination.currentPage,
      page_size: pageSize,
      total_items: pagination.totalItems,
    });
    dispatch(setPageSize(pageSize));
  };

  const handleRefresh = () => {
    dispatch(
      fetchSeriesRequest({
        page: pagination.currentPage,
        pageSize: pagination.pageSize,
      })
    );
  };

  const handleFormCancel = () => {
    setShowForm(false);
    setEditingSeries(undefined);
  };

  const handleViewScorecard = (matchId: string, seriesCreatedBy: string) => {
    setViewingScorecard(matchId);
    setCurrentSeriesCreatedBy(seriesCreatedBy);
  };

  const handleBackFromScorecard = () => {
    setViewingScorecard(null);
    setCurrentSeriesCreatedBy(null);
    if (expandedSeriesId) {
      sessionStorage.setItem('scrollToSeriesId', expandedSeriesId);
      window.history.replaceState(null, '', `/series/${expandedSeriesId}`)
    }
  };

  // Show loading state only when actually loading
  if (loading || !series || !Array.isArray(series)) {
    return (
      <div className="w-full max-w-sm mx-auto px-4 py-8 sm:max-w-md sm:px-6 md:max-w-lg md:px-8">
        <div className="flex flex-col items-center justify-center space-y-4">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
          <span className="text-sm text-gray-600">Loading series...</span>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="w-full max-w-sm mx-auto px-4 py-8 sm:max-w-md sm:px-6 md:max-w-lg md:px-8">
        <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded-lg">
          <strong className="font-bold">Error:</strong>
          <span className="block sm:inline"> {error}</span>
          <div className="mt-3">
            <button
              onClick={handleRefresh}
              className="w-full py-2 px-4 bg-red-600 text-white rounded-lg font-medium active:bg-red-700 focus:outline-none focus:ring-2 focus:ring-red-500 sm:w-auto sm:px-6"
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
    <div className="w-full max-w-4xl mx-auto p-2" data-cy="series-list">
      <div className="flex flex-col items-center space-y-4 mb-6 sm:flex-row sm:justify-between sm:space-y-0">
        <h2 className="text-2xl font-bold text-center">Cricket Series</h2>
        <div className="flex space-x-2">
          <Button
            variant="outline"
            size="sm"
            onClick={handleRefresh}
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
            New
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
            New
          </Button>
        </div>
      ) : (
        <div className="space-y-6">
          <div className="space-y-4">
            {Array.isArray(series) &&
              series.map((seriesItem, index) => (
                <div key={seriesItem.id || `series-${index}`} id={`series-card-${seriesItem.id}`}>
                  <SeriesWithMatches
                    series={seriesItem}
                    onEditSeries={handleEdit}
                    onDeleteSeries={handleDelete}
                    onViewScorecard={handleViewScorecard}
                    currentUser={currentUser}
                    isAuthenticated={isAuthenticated}
                    expanded={expandedSeriesId === seriesItem.id}
                    onToggleExpanded={(isExpanded) => {
                      setExpandedSeriesId(isExpanded ? seriesItem.id : null);
                    }}
                  />
                  {/* Insert ad after every 3 series items (only if ads enabled) */}
                  {process.env.NEXT_PUBLIC_ENABLE_ADS !== 'false' &&
                    (index + 1) % 3 === 0 &&
                    index < series.length - 1 && (
                      <InFeedAd
                        adSlot="9963510764"
                        adLayout="in-article"
                        className="my-6"
                      />
                    )}
                </div>
              ))}
          </div>

          {/* Pagination */}
          <Pagination
            currentPage={pagination.currentPage}
            totalPages={pagination.totalPages}
            totalItems={pagination.totalItems}
            pageSize={pagination.pageSize}
            onPageChange={handlePageChange}
            onPageSizeChange={handlePageSizeChange}
            pageSizeOptions={[10, 20, 50]}
            showPageSizeSelector={true}
            showTotalInfo={true}
          />
        </div>
      )}
    </div>
  );
}
