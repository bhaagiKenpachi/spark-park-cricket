import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { Provider } from 'react-redux';
import { configureStore } from '@reduxjs/toolkit';
import { SeriesList } from '@/components/SeriesList';
import { SeriesForm } from '@/components/SeriesForm';
import { seriesSlice } from '@/store/reducers/seriesSlice';
import { Series } from '@/store/reducers/seriesSlice';

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

import { apiService } from '@/services/api';

// Mock store for testing
const createMockStore = (initialState: unknown) => {
    return configureStore({
        reducer: {
            series: seriesSlice.reducer,
        },
        preloadedState: initialState,
    });
};

describe('Series Integration Tests', () => {
    beforeEach(() => {
        jest.clearAllMocks();
    });

    describe('Series List and Form Integration', () => {
        it('should complete create series workflow', async () => {
            const mockStore = createMockStore({
                series: {
                    series: [],
                    currentSeries: null,
                    loading: false,
                    error: null,
                },
            });

            const mockCreatedSeries: Series = {
                id: '1',
                name: 'New Series',
                description: 'New Description',
                start_date: '2024-01-01',
                end_date: '2024-01-31',
                status: 'upcoming',
                created_at: '2024-01-01T00:00:00Z',
                updated_at: '2024-01-01T00:00:00Z',
            };

            (apiService.createSeries as jest.Mock).mockResolvedValue({
                data: mockCreatedSeries,
                success: true,
            });

            (apiService.getSeries as jest.Mock).mockResolvedValue({
                data: [mockCreatedSeries],
                success: true,
            });

            render(
                <Provider store={mockStore}>
                    <SeriesList />
                </Provider>
            );

            // Click create series button
            const createButton = screen.getByText('Create Your First Series');
            fireEvent.click(createButton);

            // Fill out the form
            const nameInput = screen.getByLabelText('Series Name *');
            const descriptionInput = screen.getByLabelText('Description');
            const startDateInput = screen.getByLabelText('Start Date *');
            const endDateInput = screen.getByLabelText('End Date *');

            fireEvent.change(nameInput, { target: { value: 'New Series' } });
            fireEvent.change(descriptionInput, { target: { value: 'New Description' } });
            fireEvent.change(startDateInput, { target: { value: '2024-01-01' } });
            fireEvent.change(endDateInput, { target: { value: '2024-01-31' } });

            // Submit the form
            const submitButton = screen.getByText('Create Series');
            fireEvent.click(submitButton);

            // Wait for the form to close and series to appear
            await waitFor(() => {
                expect(screen.getByText('New Series')).toBeInTheDocument();
            });
        });

        it('should complete edit series workflow', async () => {
            const mockSeries: Series = {
                id: '1',
                name: 'Original Series',
                description: 'Original Description',
                start_date: '2024-01-01',
                end_date: '2024-01-31',
                status: 'upcoming',
                created_at: '2024-01-01T00:00:00Z',
                updated_at: '2024-01-01T00:00:00Z',
            };

            const mockUpdatedSeries: Series = {
                ...mockSeries,
                name: 'Updated Series',
                description: 'Updated Description',
            };

            const mockStore = createMockStore({
                series: {
                    series: [mockSeries],
                    currentSeries: null,
                    loading: false,
                    error: null,
                },
            });

            (apiService.updateSeries as jest.Mock).mockResolvedValue({
                data: mockUpdatedSeries,
                success: true,
            });

            (apiService.getSeries as jest.Mock).mockResolvedValue({
                data: [mockUpdatedSeries],
                success: true,
            });

            render(
                <Provider store={mockStore}>
                    <SeriesList />
                </Provider>
            );

            // Click edit button
            const editButton = screen.getByText('Edit');
            fireEvent.click(editButton);

            // Update the form
            const nameInput = screen.getByLabelText('Series Name *');
            const descriptionInput = screen.getByLabelText('Description');

            fireEvent.change(nameInput, { target: { value: 'Updated Series' } });
            fireEvent.change(descriptionInput, { target: { value: 'Updated Description' } });

            // Submit the form
            const submitButton = screen.getByText('Update Series');
            fireEvent.click(submitButton);

            // Wait for the updated series to appear
            await waitFor(() => {
                expect(screen.getByText('Updated Series')).toBeInTheDocument();
                expect(screen.getByText('Updated Description')).toBeInTheDocument();
            });
        });

        it('should complete delete series workflow', async () => {
            const mockSeries: Series = {
                id: '1',
                name: 'Series to Delete',
                description: 'Description',
                start_date: '2024-01-01',
                end_date: '2024-01-31',
                status: 'upcoming',
                created_at: '2024-01-01T00:00:00Z',
                updated_at: '2024-01-01T00:00:00Z',
            };

            const mockStore = createMockStore({
                series: {
                    series: [mockSeries],
                    currentSeries: null,
                    loading: false,
                    error: null,
                },
            });

            (apiService.deleteSeries as jest.Mock).mockResolvedValue({
                data: undefined,
                success: true,
            });

            (apiService.getSeries as jest.Mock).mockResolvedValue({
                data: [],
                success: true,
            });

            // Mock window.confirm
            const mockConfirm = jest.fn().mockReturnValue(true);
            Object.defineProperty(window, 'confirm', {
                value: mockConfirm,
                writable: true,
            });

            render(
                <Provider store={mockStore}>
                    <SeriesList />
                </Provider>
            );

            // Click delete button
            const deleteButton = screen.getByText('Delete');
            fireEvent.click(deleteButton);

            // Confirm deletion
            expect(mockConfirm).toHaveBeenCalledWith('Are you sure you want to delete this series?');

            // Wait for the series to be removed
            await waitFor(() => {
                expect(screen.getByText('No series found.')).toBeInTheDocument();
            });
        });
    });

    describe('Error Handling Integration', () => {
        it('should handle API errors gracefully', async () => {
            const mockStore = createMockStore({
                series: {
                    series: [],
                    currentSeries: null,
                    loading: false,
                    error: null,
                },
            });

            (apiService.createSeries as jest.Mock).mockRejectedValue(
                new Error('Network error')
            );

            render(
                <Provider store={mockStore}>
                    <SeriesForm onSuccess={jest.fn()} onCancel={jest.fn()} />
                </Provider>
            );

            // Fill out the form
            const nameInput = screen.getByLabelText('Series Name *');
            const startDateInput = screen.getByLabelText('Start Date *');
            const endDateInput = screen.getByLabelText('End Date *');

            fireEvent.change(nameInput, { target: { value: 'Test Series' } });
            fireEvent.change(startDateInput, { target: { value: '2024-01-01' } });
            fireEvent.change(endDateInput, { target: { value: '2024-01-31' } });

            // Submit the form
            const submitButton = screen.getByText('Create Series');
            fireEvent.click(submitButton);

            // Wait for error to appear
            await waitFor(() => {
                expect(screen.getByText('Failed to create series')).toBeInTheDocument();
            });
        });

        it('should handle validation errors', async () => {
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
                    <SeriesForm onSuccess={jest.fn()} onCancel={jest.fn()} />
                </Provider>
            );

            // Try to submit without filling required fields
            const submitButton = screen.getByText('Create Series');
            fireEvent.click(submitButton);

            // Wait for validation errors
            await waitFor(() => {
                expect(screen.getByText('Name is required')).toBeInTheDocument();
                expect(screen.getByText('Start date is required')).toBeInTheDocument();
                expect(screen.getByText('End date is required')).toBeInTheDocument();
            });
        });
    });

    describe('Loading States Integration', () => {
        it('should show loading state during API calls', async () => {
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

        it('should disable form during submission', async () => {
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
                    <SeriesForm onSuccess={jest.fn()} onCancel={jest.fn()} />
                </Provider>
            );

            const submitButton = screen.getByText('Saving...');
            expect(submitButton).toBeDisabled();
        });
    });
});
