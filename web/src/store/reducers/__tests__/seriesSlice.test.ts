import { seriesSlice, Series } from '../seriesSlice';

describe('seriesSlice', () => {
    const initialState = {
        series: [],
        currentSeries: null,
        loading: false,
        error: null,
    };

    it('should handle initial state', () => {
        expect(seriesSlice.reducer(undefined, { type: 'unknown' })).toEqual(initialState);
    });

    it('should handle fetchSeriesRequest', () => {
        const actual = seriesSlice.reducer(initialState, seriesSlice.actions.fetchSeriesRequest());
        expect(actual.loading).toBe(true);
        expect(actual.error).toBe(null);
    });

    it('should handle fetchSeriesSuccess', () => {
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

        const actual = seriesSlice.reducer(
            { ...initialState, loading: true },
            seriesSlice.actions.fetchSeriesSuccess(mockSeries)
        );

        expect(actual.loading).toBe(false);
        expect(actual.series).toEqual(mockSeries);
        expect(actual.error).toBe(null);
    });

    it('should handle fetchSeriesFailure', () => {
        const errorMessage = 'Failed to fetch series';
        const actual = seriesSlice.reducer(
            { ...initialState, loading: true },
            seriesSlice.actions.fetchSeriesFailure(errorMessage)
        );

        expect(actual.loading).toBe(false);
        expect(actual.error).toBe(errorMessage);
    });

    it('should handle setCurrentSeries', () => {
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

        const actual = seriesSlice.reducer(
            initialState,
            seriesSlice.actions.setCurrentSeries(mockSeries)
        );

        expect(actual.currentSeries).toEqual(mockSeries);
    });

    it('should handle createSeriesRequest', () => {
        const seriesData = {
            name: 'New Series',
            description: 'New Description',
            start_date: '2024-02-01',
            end_date: '2024-02-28',
            status: 'upcoming' as const,
        };

        const actual = seriesSlice.reducer(
            initialState,
            seriesSlice.actions.createSeriesRequest(seriesData)
        );

        expect(actual.loading).toBe(true);
        expect(actual.error).toBe(null);
    });

    it('should handle createSeriesSuccess', () => {
        const newSeries: Series = {
            id: '2',
            name: 'New Series',
            description: 'New Description',
            start_date: '2024-02-01',
            end_date: '2024-02-28',
            status: 'upcoming',
            created_at: '2024-02-01T00:00:00Z',
            updated_at: '2024-02-01T00:00:00Z',
        };

        const actual = seriesSlice.reducer(
            { ...initialState, loading: true },
            seriesSlice.actions.createSeriesSuccess(newSeries)
        );

        expect(actual.loading).toBe(false);
        expect(actual.series).toContain(newSeries);
    });

    it('should handle createSeriesFailure', () => {
        const errorMessage = 'Failed to create series';
        const actual = seriesSlice.reducer(
            { ...initialState, loading: true },
            seriesSlice.actions.createSeriesFailure(errorMessage)
        );

        expect(actual.loading).toBe(false);
        expect(actual.error).toBe(errorMessage);
    });

    it('should handle updateSeriesSuccess', () => {
        const existingSeries: Series = {
            id: '1',
            name: 'Original Series',
            description: 'Original Description',
            start_date: '2024-01-01',
            end_date: '2024-01-31',
            status: 'upcoming',
            created_at: '2024-01-01T00:00:00Z',
            updated_at: '2024-01-01T00:00:00Z',
        };

        const updatedSeries: Series = {
            ...existingSeries,
            name: 'Updated Series',
            description: 'Updated Description',
        };

        const stateWithSeries = {
            ...initialState,
            series: [existingSeries],
        };

        const actual = seriesSlice.reducer(
            { ...stateWithSeries, loading: true },
            seriesSlice.actions.updateSeriesSuccess(updatedSeries)
        );

        expect(actual.loading).toBe(false);
        expect(actual.series[0]).toEqual(updatedSeries);
    });

    it('should handle deleteSeriesSuccess', () => {
        const seriesToDelete: Series = {
            id: '1',
            name: 'Series to Delete',
            description: 'Description',
            start_date: '2024-01-01',
            end_date: '2024-01-31',
            status: 'upcoming',
            created_at: '2024-01-01T00:00:00Z',
            updated_at: '2024-01-01T00:00:00Z',
        };

        const remainingSeries: Series = {
            id: '2',
            name: 'Remaining Series',
            description: 'Description',
            start_date: '2024-02-01',
            end_date: '2024-02-28',
            status: 'upcoming',
            created_at: '2024-02-01T00:00:00Z',
            updated_at: '2024-02-01T00:00:00Z',
        };

        const stateWithSeries = {
            ...initialState,
            series: [seriesToDelete, remainingSeries],
        };

        const actual = seriesSlice.reducer(
            { ...stateWithSeries, loading: true },
            seriesSlice.actions.deleteSeriesSuccess('1')
        );

        expect(actual.loading).toBe(false);
        expect(actual.series).toHaveLength(1);
        expect(actual.series[0]).toEqual(remainingSeries);
    });
});
