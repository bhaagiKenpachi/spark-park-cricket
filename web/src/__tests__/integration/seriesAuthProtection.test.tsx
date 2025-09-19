import { render, screen } from '@testing-library/react';
import { Provider } from 'react-redux';
import { configureStore } from '@reduxjs/toolkit';
import { SeriesList } from '@/components/SeriesList';
import { ProtectedRoute } from '@/components/auth/ProtectedRoute';
import { seriesSlice } from '@/store/reducers/seriesSlice';
import { matchSlice } from '@/store/reducers/matchSlice';
import authReducer from '@/store/reducers/authSlice';

// Mock the API service
jest.mock('@/services/api', () => ({
  apiService: {
    getSeries: jest.fn(),
    createSeries: jest.fn(),
    updateSeries: jest.fn(),
    deleteSeries: jest.fn(),
  },
  ApiError: class ApiError extends Error {
    status: number;
    details?: unknown;
    constructor(message: string, status: number, details?: unknown) {
      super(message);
      this.name = 'ApiError';
      this.status = status;
      this.details = details;
    }
  },
}));

// Mock the auth service to prevent initialization
jest.mock('@/services/authService', () => ({
  authService: {
    getAuthStatus: jest.fn(),
    getCurrentUser: jest.fn(),
    logout: jest.fn(),
    getStoredUser: jest.fn(),
    isAuthenticated: jest.fn(),
    clearAuthState: jest.fn(),
    setAuthState: jest.fn(),
  },
}));

// Mock store for testing
const createMockStore = (authState: unknown) => {
  return configureStore({
    reducer: {
      series: seriesSlice.reducer,
      match: matchSlice.reducer,
      auth: authReducer,
    },
    preloadedState: {
      auth: authState,
    },
  });
};

describe('Series Authentication Protection', () => {
  describe('Unauthenticated User', () => {
    it('should show authentication required message when not logged in', () => {
      const mockStore = createMockStore({
        auth: {
          user: null,
          isAuthenticated: false,
          isLoading: false,
          error: null,
          isInitialized: true,
        },
        series: {
          series: [],
          loading: false,
          error: null,
        },
        match: {
          matches: [],
          loading: false,
          error: null,
        },
      });

      render(
        <Provider store={mockStore}>
          <ProtectedRoute>
            <SeriesList />
          </ProtectedRoute>
        </Provider>
      );

      // Should show authentication required message
      expect(screen.getByText('Authentication Required')).toBeInTheDocument();
      expect(screen.getByText('Please sign in to access this feature.')).toBeInTheDocument();
      
      // Should not show series list
      expect(screen.queryByText('Cricket Series')).not.toBeInTheDocument();
      expect(screen.queryByText('Create Series')).not.toBeInTheDocument();
    });

    it('should show login button when not authenticated', () => {
      const mockStore = createMockStore({
        auth: {
          user: null,
          isAuthenticated: false,
          isLoading: false,
          error: null,
          isInitialized: true,
        },
        series: {
          series: [],
          loading: false,
          error: null,
        },
        match: {
          matches: [],
          loading: false,
          error: null,
        },
      });

      render(
        <Provider store={mockStore}>
          <ProtectedRoute>
            <SeriesList />
          </ProtectedRoute>
        </Provider>
      );

      // Should show login button
      expect(screen.getByRole('button', { name: /sign in with google/i })).toBeInTheDocument();
    });
  });

  describe('Authenticated User', () => {
    it('should show series list when authenticated', () => {
      const mockStore = createMockStore({
        auth: {
          user: {
            id: '1',
            name: 'Test User',
            email: 'test@example.com',
            picture: 'https://example.com/avatar.jpg',
          },
          isAuthenticated: true,
          isLoading: false,
          error: null,
          isInitialized: true,
        },
        series: {
          series: [],
          loading: false,
          error: null,
        },
        match: {
          matches: [],
          loading: false,
          error: null,
        },
      });

      render(
        <Provider store={mockStore}>
          <ProtectedRoute>
            <SeriesList />
          </ProtectedRoute>
        </Provider>
      );

      // Should show series list
      expect(screen.getByText('Cricket Series')).toBeInTheDocument();
      expect(screen.getByText('Create Series')).toBeInTheDocument();
      
      // Should not show authentication required message
      expect(screen.queryByText('Authentication Required')).not.toBeInTheDocument();
    });
  });

  describe('Loading State', () => {
    it('should show loading spinner while auth is initializing', () => {
      const mockStore = createMockStore({
        auth: {
          user: null,
          isAuthenticated: false,
          isLoading: false,
          error: null,
          isInitialized: false, // Auth not initialized yet
        },
        series: {
          series: [],
          loading: false,
          error: null,
        },
        match: {
          matches: [],
          loading: false,
          error: null,
        },
      });

      render(
        <Provider store={mockStore}>
          <ProtectedRoute>
            <SeriesList />
          </ProtectedRoute>
        </Provider>
      );

      // Should show loading state
      expect(screen.getByText('Loading...')).toBeInTheDocument();
      
      // Should not show series list or auth message yet
      expect(screen.queryByText('Cricket Series')).not.toBeInTheDocument();
      expect(screen.queryByText('Authentication Required')).not.toBeInTheDocument();
    });
  });
});
