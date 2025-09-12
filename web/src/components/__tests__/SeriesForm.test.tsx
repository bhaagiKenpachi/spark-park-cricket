import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { Provider } from 'react-redux';
import { configureStore } from '@reduxjs/toolkit';
import { SeriesForm } from '../SeriesForm';
import { seriesSlice } from '@/store/reducers/seriesSlice';
import { Series } from '@/store/reducers/seriesSlice';

// Mock store for testing
const createMockStore = (initialState: unknown) => {
    return configureStore({
        reducer: {
            series: seriesSlice.reducer,
        },
        preloadedState: initialState,
    });
};

describe('SeriesForm', () => {
    const mockOnSuccess = jest.fn();
    const mockOnCancel = jest.fn();

    beforeEach(() => {
        mockOnSuccess.mockClear();
        mockOnCancel.mockClear();
    });

    it('should render create form', () => {
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
                <SeriesForm onSuccess={mockOnSuccess} onCancel={mockOnCancel} />
            </Provider>
        );

        expect(screen.getByText('Create New Series')).toBeInTheDocument();
        expect(screen.getByLabelText('Series Name *')).toBeInTheDocument();
        expect(screen.getByLabelText('Description')).toBeInTheDocument();
        expect(screen.getByLabelText('Start Date *')).toBeInTheDocument();
        expect(screen.getByLabelText('End Date *')).toBeInTheDocument();
        expect(screen.getByLabelText('Status')).toBeInTheDocument();
        expect(screen.getByText('Create Series')).toBeInTheDocument();
        expect(screen.getByText('Cancel')).toBeInTheDocument();
    });

    it('should render edit form with existing data', () => {
        const mockSeries: Series = {
            id: '1',
            name: 'Test Series',
            description: 'Test Description',
            start_date: '2024-01-01',
            end_date: '2024-01-31',
            status: 'upcoming',
            created_at: '2024-01-01T00:00:00Z',
            updated_at: '2024-01-01T00:00:00Z',
        };

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
                <SeriesForm series={mockSeries} onSuccess={mockOnSuccess} onCancel={mockOnCancel} />
            </Provider>
        );

        expect(screen.getByText('Edit Series')).toBeInTheDocument();
        expect(screen.getByDisplayValue('Test Series')).toBeInTheDocument();
        expect(screen.getByDisplayValue('Test Description')).toBeInTheDocument();
        expect(screen.getByDisplayValue('2024-01-01')).toBeInTheDocument();
        expect(screen.getByDisplayValue('2024-01-31')).toBeInTheDocument();
        expect(screen.getByDisplayValue('upcoming')).toBeInTheDocument();
        expect(screen.getByText('Update Series')).toBeInTheDocument();
    });

    it('should show loading state when submitting', () => {
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
                <SeriesForm onSuccess={mockOnSuccess} onCancel={mockOnCancel} />
            </Provider>
        );

        expect(screen.getByText('Saving...')).toBeInTheDocument();
        expect(screen.getByText('Saving...')).toBeDisabled();
    });

    it('should show error message', () => {
        const mockStore = createMockStore({
            series: {
                series: [],
                currentSeries: null,
                loading: false,
                error: 'Failed to create series',
            },
        });

        render(
            <Provider store={mockStore}>
                <SeriesForm onSuccess={mockOnSuccess} onCancel={mockOnCancel} />
            </Provider>
        );

        expect(screen.getByText('Failed to create series')).toBeInTheDocument();
    });

    it('should validate required fields', async () => {
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
                <SeriesForm onSuccess={mockOnSuccess} onCancel={mockOnCancel} />
            </Provider>
        );

        const submitButton = screen.getByText('Create Series');
        fireEvent.click(submitButton);

        await waitFor(() => {
            expect(screen.getByText('Name is required')).toBeInTheDocument();
            expect(screen.getByText('Start date is required')).toBeInTheDocument();
            expect(screen.getByText('End date is required')).toBeInTheDocument();
        });
    });

    it('should validate end date is after start date', async () => {
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
                <SeriesForm onSuccess={mockOnSuccess} onCancel={mockOnCancel} />
            </Provider>
        );

        const nameInput = screen.getByLabelText('Series Name *');
        const startDateInput = screen.getByLabelText('Start Date *');
        const endDateInput = screen.getByLabelText('End Date *');
        const submitButton = screen.getByText('Create Series');

        fireEvent.change(nameInput, { target: { value: 'Test Series' } });
        fireEvent.change(startDateInput, { target: { value: '2024-01-31' } });
        fireEvent.change(endDateInput, { target: { value: '2024-01-01' } });
        fireEvent.click(submitButton);

        await waitFor(() => {
            expect(screen.getByText('End date must be after start date')).toBeInTheDocument();
        });
    });

    it('should clear validation errors when user types', async () => {
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
                <SeriesForm onSuccess={mockOnSuccess} onCancel={mockOnCancel} />
            </Provider>
        );

        const submitButton = screen.getByText('Create Series');
        fireEvent.click(submitButton);

        await waitFor(() => {
            expect(screen.getByText('Name is required')).toBeInTheDocument();
        });

        const nameInput = screen.getByLabelText('Series Name *');
        fireEvent.change(nameInput, { target: { value: 'Test Series' } });

        await waitFor(() => {
            expect(screen.queryByText('Name is required')).not.toBeInTheDocument();
        });
    });

    it('should call onCancel when cancel button is clicked', () => {
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
                <SeriesForm onSuccess={mockOnSuccess} onCancel={mockOnCancel} />
            </Provider>
        );

        const cancelButton = screen.getByText('Cancel');
        fireEvent.click(cancelButton);

        expect(mockOnCancel).toHaveBeenCalledTimes(1);
    });

    it('should submit form with valid data', async () => {
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
                <SeriesForm onSuccess={mockOnSuccess} onCancel={mockOnCancel} />
            </Provider>
        );

        const nameInput = screen.getByLabelText('Series Name *');
        const descriptionInput = screen.getByLabelText('Description');
        const startDateInput = screen.getByLabelText('Start Date *');
        const endDateInput = screen.getByLabelText('End Date *');
        const statusSelect = screen.getByLabelText('Status');
        const submitButton = screen.getByText('Create Series');

        fireEvent.change(nameInput, { target: { value: 'Test Series' } });
        fireEvent.change(descriptionInput, { target: { value: 'Test Description' } });
        fireEvent.change(startDateInput, { target: { value: '2024-01-01' } });
        fireEvent.change(endDateInput, { target: { value: '2024-01-31' } });
        fireEvent.change(statusSelect, { target: { value: 'upcoming' } });
        fireEvent.click(submitButton);

        expect(mockOnSuccess).toHaveBeenCalledTimes(1);
    });
});
