import { put } from 'redux-saga/effects';
import {
  fetchSeriesSuccess,
  fetchSeriesFailure,
  createSeriesRequest,
  createSeriesSuccess,
  createSeriesFailure,
  updateSeriesRequest,
  updateSeriesSuccess,
  updateSeriesFailure,
  deleteSeriesRequest,
  deleteSeriesSuccess,
  deleteSeriesFailure,
} from '../../reducers/seriesSlice';
import { Series } from '../../reducers/seriesSlice';

// Mock the API service
jest.mock('@/services/api', () => ({
  ApiService: jest.fn().mockImplementation(() => ({
    getSeries: jest.fn(),
    createSeries: jest.fn(),
    updateSeries: jest.fn(),
    deleteSeries: jest.fn(),
  })),
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

// Import the mocked API service
import { ApiService, ApiError } from '@/services/api';

// Import saga functions (we need to export them from the saga file)
import {
  fetchSeriesSaga,
  createSeriesSaga,
  updateSeriesSaga,
  deleteSeriesSaga,
} from '../seriesSaga';

describe('seriesSaga', () => {
  let mockApiService: {
    getSeries: jest.Mock;
    createSeries: jest.Mock;
    updateSeries: jest.Mock;
    deleteSeries: jest.Mock;
  };

  beforeEach(() => {
    jest.clearAllMocks();
    // Get the mock instance
    mockApiService = {
      getSeries: jest.fn(),
      createSeries: jest.fn(),
      updateSeries: jest.fn(),
      deleteSeries: jest.fn(),
    };
    (ApiService as jest.Mock).mockImplementation(() => mockApiService);
  });

  describe('fetchSeriesSaga', () => {
    it('should fetch series successfully', () => {
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

      const mockResponse = { data: { data: mockSeries }, success: true };
      mockApiService.getSeries.mockResolvedValue(mockResponse);

      const generator = fetchSeriesSaga();
      const apiCall = generator.next().value;
      const putAction = generator.next(mockResponse).value;

      expect(apiCall).toMatchObject({
        '@@redux-saga/IO': true,
        combinator: false,
        type: 'CALL',
      });
      expect(putAction).toEqual(put(fetchSeriesSuccess(mockSeries)));
    });

    it('should handle fetch series failure with ApiError', () => {
      const error = new ApiError('Network error', 500);
      const generator = fetchSeriesSaga();

      generator.next(); // Skip the API call
      const delayAction = generator.throw(error).value;
      const putAction = generator.next().value;

      expect(delayAction).toMatchObject({
        '@@redux-saga/IO': true,
        combinator: false,
        type: 'CALL',
      });
      expect(putAction).toEqual(put(fetchSeriesFailure('Network error')));
    });

    it('should handle fetch series failure with generic error', () => {
      const error = new Error('Generic error');
      const generator = fetchSeriesSaga();

      generator.next(); // Skip the API call
      const delayAction = generator.throw(error).value;
      const putAction = generator.next().value;

      expect(delayAction).toMatchObject({
        '@@redux-saga/IO': true,
        combinator: false,
        type: 'CALL',
      });
      expect(putAction).toEqual(
        put(fetchSeriesFailure('Failed to fetch series'))
      );
    });
  });

  describe('createSeriesSaga', () => {
    it('should create series successfully', () => {
      const seriesData = {
        name: 'New Series',
        description: 'New Description',
        start_date: '2024-02-01',
        end_date: '2024-02-28',
        status: 'upcoming' as const,
      };

      const createdSeries: Series = {
        id: '2',
        ...seriesData,
        created_at: '2024-02-01T00:00:00Z',
        updated_at: '2024-02-01T00:00:00Z',
      };

      const mockResponse = { data: { data: createdSeries }, success: true };
      mockApiService.createSeries.mockResolvedValue(mockResponse);

      const action = createSeriesRequest(seriesData);
      const generator = createSeriesSaga(action);

      const apiCall = generator.next().value;
      const putAction = generator.next(mockResponse).value;

      expect(apiCall).toMatchObject({
        '@@redux-saga/IO': true,
        combinator: false,
        type: 'CALL',
      });
      expect(putAction).toEqual(put(createSeriesSuccess(createdSeries)));
    });

    it('should handle create series failure', () => {
      const seriesData = {
        name: 'New Series',
        description: 'New Description',
        start_date: '2024-02-01',
        end_date: '2024-02-28',
        status: 'upcoming' as const,
      };

      const error = new ApiError('Validation error', 400);
      const action = createSeriesRequest(seriesData);
      const generator = createSeriesSaga(action);

      generator.next(); // Skip the API call
      const putAction = generator.throw(error).value;

      expect(putAction).toEqual(put(createSeriesFailure('Validation error')));
    });
  });

  describe('updateSeriesSaga', () => {
    it('should update series successfully', () => {
      const seriesData = { name: 'Updated Series' };
      const updatedSeries: Series = {
        id: '1',
        name: 'Updated Series',
        description: 'Test Description',
        start_date: '2024-01-01',
        end_date: '2024-01-31',
        status: 'upcoming',
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z',
      };

      const mockResponse = { data: { data: updatedSeries }, success: true };
      mockApiService.updateSeries.mockResolvedValue(mockResponse);

      const action = updateSeriesRequest({ id: '1', seriesData });
      const generator = updateSeriesSaga(action);

      const apiCall = generator.next().value;
      const putAction = generator.next(mockResponse).value;

      expect(apiCall).toMatchObject({
        '@@redux-saga/IO': true,
        combinator: false,
        type: 'CALL',
      });
      expect(putAction).toEqual(put(updateSeriesSuccess(updatedSeries)));
    });

    it('should handle update series failure', () => {
      const seriesData = { name: 'Updated Series' };
      const error = new ApiError('Not found', 404);
      const action = updateSeriesRequest({ id: '1', seriesData });
      const generator = updateSeriesSaga(action);

      generator.next(); // Skip the API call
      const putAction = generator.throw(error).value;

      expect(putAction).toEqual(put(updateSeriesFailure('Not found')));
    });
  });

  describe('deleteSeriesSaga', () => {
    it('should delete series successfully', () => {
      const mockResponse = { data: { data: undefined }, success: true };
      mockApiService.deleteSeries.mockResolvedValue(mockResponse);

      const action = deleteSeriesRequest('1');
      const generator = deleteSeriesSaga(action);

      const apiCall = generator.next().value;
      const putAction = generator.next(mockResponse).value;

      expect(apiCall).toMatchObject({
        '@@redux-saga/IO': true,
        combinator: false,
        type: 'CALL',
      });
      expect(putAction).toEqual(put(deleteSeriesSuccess('1')));
    });

    it('should handle delete series failure', () => {
      const error = new ApiError('Not found', 404);
      const action = deleteSeriesRequest('1');
      const generator = deleteSeriesSaga(action);

      generator.next(); // Skip the API call
      const putAction = generator.throw(error).value;

      expect(putAction).toEqual(put(deleteSeriesFailure('Not found')));
    });
  });
});
