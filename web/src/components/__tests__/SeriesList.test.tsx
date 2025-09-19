import { render, screen, fireEvent } from '@testing-library/react';
import { Provider } from 'react-redux';
import { configureStore } from '@reduxjs/toolkit';
import { SeriesList } from '../SeriesList';
import { seriesSlice } from '@/store/reducers/seriesSlice';
import { Series } from '@/store/reducers/seriesSlice';

// Mock the hooks
const mockDispatch = jest.fn();
jest.mock('@/store/hooks', () => ({
  useAppDispatch: () => mockDispatch,
  useAppSelector: jest.fn(),
}));

import { useAppSelector } from '@/store/hooks';

// Mock store for testing
const createMockStore = (initialState: unknown) => {
  return configureStore({
    reducer: {
      series: seriesSlice.reducer,
    },
    preloadedState: initialState,
  });
};

// Mock window.confirm
const mockConfirm = jest.fn();
Object.defineProperty(window, 'confirm', {
  value: mockConfirm,
  writable: true,
});

describe('SeriesList', () => {
  beforeEach(() => {
    mockConfirm.mockClear();
    mockDispatch.mockClear();
  });

  it('should render loading state when loading is true and no series', () => {
    (useAppSelector as jest.Mock).mockReturnValue({
      series: [],
      currentSeries: null,
      loading: true,
      error: null,
    });

    const mockStore = createMockStore({
      series: {
        series: [],
        currentSeries: null,
        loading: true,
        error: null,
      },
    });

    render(
      <Provider store={mockStore}>
        <SeriesList />
      </Provider>
    );

    expect(screen.getByText('Loading series...')).toBeInTheDocument();
  });

  it('should render error state when error exists', () => {
    (useAppSelector as jest.Mock).mockReturnValue({
      series: [],
      currentSeries: null,
      loading: false,
      error: 'Failed to fetch series',
    });

    const mockStore = createMockStore({
      series: {
        series: [],
        currentSeries: null,
        loading: false,
        error: 'Failed to fetch series',
      },
    });

    render(
      <Provider store={mockStore}>
        <SeriesList />
      </Provider>
    );

    expect(screen.getByText('Error:')).toBeInTheDocument();
    expect(screen.getByText('Failed to fetch series')).toBeInTheDocument();
  });

  it('should render empty state when no series and no loading', () => {
    (useAppSelector as jest.Mock).mockReturnValue({
      series: [],
      currentSeries: null,
      loading: false,
      error: null,
    });

    const mockStore = createMockStore({
      series: {
        series: [],
        currentSeries: null,
        loading: false,
        error: null,
      },
    });

    render(
      <Provider store={mockStore}>
        <SeriesList />
      </Provider>
    );

    expect(screen.getByText('No series found.')).toBeInTheDocument();
    expect(screen.getByText('Create Your First Series')).toBeInTheDocument();
  });

  it('should render series list when series exist', () => {
    const mockSeries: Series[] = [
      {
        id: '1',
        name: 'Test Series',
        description: 'Test Description',
        start_date: '2024-01-01',
        end_date: '2024-01-31',
        status: 'upcoming',
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z',
      },
    ];

    (useAppSelector as jest.Mock).mockReturnValue({
      series: mockSeries,
      currentSeries: null,
      loading: false,
      error: null,
    });

    const mockStore = createMockStore({
      series: {
        series: mockSeries,
        currentSeries: null,
        loading: false,
        error: null,
      },
    });

    render(
      <Provider store={mockStore}>
        <SeriesList />
      </Provider>
    );

    expect(screen.getByText('Cricket Series')).toBeInTheDocument();
    expect(screen.getByText('Test Series')).toBeInTheDocument();
  });

  it('should show create series form when create button is clicked', () => {
    (useAppSelector as jest.Mock).mockReturnValue({
      series: [],
      currentSeries: null,
      loading: false,
      error: null,
    });

    const mockStore = createMockStore({
      series: {
        series: [],
        currentSeries: null,
        loading: false,
        error: null,
      },
    });

    render(
      <Provider store={mockStore}>
        <SeriesList />
      </Provider>
    );

    const createButton = screen.getByText('Create Your First Series');
    fireEvent.click(createButton);

    expect(screen.getByText('Create New Series')).toBeInTheDocument();
  });

  it('should show edit series form when edit button is clicked', () => {
    const mockSeries: Series[] = [
      {
        id: '1',
        name: 'Test Series',
        description: 'Test Description',
        start_date: '2024-01-01',
        end_date: '2024-01-31',
        status: 'upcoming',
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z',
      },
    ];

    (useAppSelector as jest.Mock).mockReturnValue({
      series: mockSeries,
      currentSeries: null,
      loading: false,
      error: null,
    });

    const mockStore = createMockStore({
      series: {
        series: mockSeries,
        currentSeries: null,
        loading: false,
        error: null,
      },
    });

    render(
      <Provider store={mockStore}>
        <SeriesList />
      </Provider>
    );

    const editButton = screen.getByText('Edit');
    fireEvent.click(editButton);

    expect(screen.getByText('Edit Series')).toBeInTheDocument();
  });

  it('should call delete action when delete button is clicked and confirmed', () => {
    const mockSeries: Series[] = [
      {
        id: '1',
        name: 'Test Series',
        description: 'Test Description',
        start_date: '2024-01-01',
        end_date: '2024-01-31',
        status: 'upcoming',
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z',
      },
    ];

    (useAppSelector as jest.Mock).mockReturnValue({
      series: mockSeries,
      currentSeries: null,
      loading: false,
      error: null,
    });

    const mockStore = createMockStore({
      series: {
        series: mockSeries,
        currentSeries: null,
        loading: false,
        error: null,
      },
    });

    mockConfirm.mockReturnValue(true);

    render(
      <Provider store={mockStore}>
        <SeriesList />
      </Provider>
    );

    const deleteButton = screen.getByText('Delete');
    fireEvent.click(deleteButton);

    expect(mockConfirm).toHaveBeenCalledWith(
      'Are you sure you want to delete this series?'
    );
  });

  it('should display correct status badges', () => {
    const mockSeries: Series[] = [
      {
        id: '1',
        name: 'Upcoming Series',
        description: 'Test Description',
        start_date: '2024-01-01',
        end_date: '2024-01-31',
        status: 'upcoming',
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z',
      },
      {
        id: '2',
        name: 'Ongoing Series',
        description: 'Test Description',
        start_date: '2024-02-01',
        end_date: '2024-02-28',
        status: 'ongoing',
        created_at: '2024-02-01T00:00:00Z',
        updated_at: '2024-02-01T00:00:00Z',
      },
      {
        id: '3',
        name: 'Completed Series',
        description: 'Test Description',
        start_date: '2024-03-01',
        end_date: '2024-03-31',
        status: 'completed',
        created_at: '2024-03-01T00:00:00Z',
        updated_at: '2024-03-01T00:00:00Z',
      },
    ];

    (useAppSelector as jest.Mock).mockReturnValue({
      series: mockSeries,
      currentSeries: null,
      loading: false,
      error: null,
    });

    const mockStore = createMockStore({
      series: {
        series: mockSeries,
        currentSeries: null,
        loading: false,
        error: null,
      },
    });

    render(
      <Provider store={mockStore}>
        <SeriesList />
      </Provider>
    );

    expect(screen.getByText('upcoming')).toBeInTheDocument();
    expect(screen.getByText('ongoing')).toBeInTheDocument();
    expect(screen.getByText('completed')).toBeInTheDocument();
  });
});
