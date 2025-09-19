import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { Provider } from 'react-redux';
import { configureStore } from '@reduxjs/toolkit';
import { MatchForm } from '@/components/MatchForm';
import {
  matchSlice,
  createMatchRequest,
  updateMatchRequest,
} from '../../store/reducers/matchSlice';
import { Match } from '../../store/reducers/matchSlice';

// Mock the API service
jest.mock('../../services/api', () => ({
  apiService: {
    getMatches: jest.fn(),
    createMatch: jest.fn(),
    updateMatch: jest.fn(),
    deleteMatch: jest.fn(),
    getMatchesBySeries: jest.fn(),
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

import { apiService } from '../../services/api';

// Mock Redux actions
jest.mock('../../store/reducers/matchSlice', () => ({
  ...jest.requireActual('../../store/reducers/matchSlice'),
  createMatchRequest: jest
    .fn()
    .mockReturnValue({ type: 'match/createMatchRequest' }),
  updateMatchRequest: jest
    .fn()
    .mockReturnValue({ type: 'match/updateMatchRequest' }),
}));

// Mock store for testing
const createMockStore = (initialState: unknown) => {
  return configureStore({
    reducer: {
      match: matchSlice.reducer,
    },
    preloadedState: initialState,
  });
};

describe('Match Integration Tests', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('Match CRUD Operations Integration', () => {
    it('should complete create match workflow', async () => {
      const mockStore = createMockStore({
        match: {
          matches: [],
          currentMatch: null,
          loading: false,
          error: null,
        },
      });

      const mockCreatedMatch: Match = {
        id: '1',
        series_id: 'series-1',
        match_number: 1,
        date: '2024-01-01T00:00:00Z',
        status: 'live',
        team_a_player_count: 11,
        team_b_player_count: 11,
        total_overs: 20,
        toss_winner: 'A',
        toss_type: 'H',
        batting_team: 'A',
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z',
      };

      (apiService.createMatch as jest.Mock).mockResolvedValue({
        data: mockCreatedMatch,
        success: true,
      });

      (apiService.getMatches as jest.Mock).mockResolvedValue({
        data: [mockCreatedMatch],
        success: true,
      });

      render(
        <Provider store={mockStore}>
          <MatchForm
            seriesId="series-1"
            onSuccess={jest.fn()}
            onCancel={jest.fn()}
          />
        </Provider>
      );

      // Fill out the form with all required fields
      const dateInput = screen.getByLabelText(/Date \*/);
      const playerCountInput = screen.getByLabelText(/Team Player Count \*/);
      const oversInput = screen.getByLabelText(/Total Overs \*/);
      const matchNumberInput = screen.getByLabelText(
        /Match Number \(Optional\)/
      );

      fireEvent.change(dateInput, { target: { value: '2024-01-01' } });
      fireEvent.change(playerCountInput, { target: { value: '11' } });
      fireEvent.change(oversInput, { target: { value: '20' } });
      fireEvent.change(matchNumberInput, { target: { value: '1' } });

      // Submit the form
      const submitButton = screen.getByText('Match');
      fireEvent.click(submitButton);

      // Wait for the Redux action to be dispatched
      await waitFor(() => {
        expect(createMatchRequest).toHaveBeenCalledWith({
          series_id: 'series-1',
          match_number: 1,
          date: '2024-01-01T00:00:00Z',
          status: 'live',
          team_a_player_count: 11,
          team_b_player_count: 11,
          total_overs: 20,
          toss_winner: 'A',
          toss_type: 'H',
          batting_team: 'A',
        });
      });
    });

    it('should complete edit match workflow', async () => {
      const mockMatch: Match = {
        id: '1',
        series_id: 'series-1',
        match_number: 1,
        date: '2024-01-01T00:00:00Z',
        status: 'live',
        team_a_player_count: 11,
        team_b_player_count: 11,
        total_overs: 20,
        toss_winner: 'A',
        toss_type: 'H',
        batting_team: 'A',
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z',
      };

      const mockUpdatedMatch: Match = {
        ...mockMatch,
        team_a_player_count: 10,
        team_b_player_count: 10,
        total_overs: 15,
      };

      const mockStore = createMockStore({
        match: {
          matches: [mockMatch],
          currentMatch: null,
          loading: false,
          error: null,
        },
      });

      (apiService.updateMatch as jest.Mock).mockResolvedValue({
        data: mockUpdatedMatch,
        success: true,
      });

      (apiService.getMatches as jest.Mock).mockResolvedValue({
        data: [mockUpdatedMatch],
        success: true,
      });

      render(
        <Provider store={mockStore}>
          <MatchForm
            match={mockMatch}
            onSuccess={jest.fn()}
            onCancel={jest.fn()}
          />
        </Provider>
      );

      // Update the form
      const playerCountInput = screen.getByLabelText(/Team Player Count \*/);
      const oversInput = screen.getByLabelText(/Total Overs \*/);

      fireEvent.change(playerCountInput, { target: { value: '10' } });
      fireEvent.change(oversInput, { target: { value: '15' } });

      // Submit the form
      const submitButton = screen.getByText('Match');
      fireEvent.click(submitButton);

      // Wait for the Redux action to be dispatched
      await waitFor(() => {
        expect(updateMatchRequest).toHaveBeenCalledWith({
          id: '1',
          matchData: {
            series_id: 'series-1',
            match_number: 1,
            date: '2024-01-01T00:00:00Z',
            status: 'live',
            team_a_player_count: 10,
            team_b_player_count: 10,
            total_overs: 15,
            toss_winner: 'A',
            toss_type: 'H',
            batting_team: 'A',
          },
        });
      });
    });

    it('should handle match number auto-generation', async () => {
      const mockStore = createMockStore({
        match: {
          matches: [],
          currentMatch: null,
          loading: false,
          error: null,
        },
      });

      const mockCreatedMatch: Match = {
        id: '1',
        series_id: 'series-1',
        match_number: 1, // Backend auto-generated
        date: '2024-01-01T00:00:00Z',
        status: 'live',
        team_a_player_count: 11,
        team_b_player_count: 11,
        total_overs: 20,
        toss_winner: 'A',
        toss_type: 'H',
        batting_team: 'A',
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z',
      };

      (apiService.createMatch as jest.Mock).mockResolvedValue({
        data: mockCreatedMatch,
        success: true,
      });

      render(
        <Provider store={mockStore}>
          <MatchForm
            seriesId="series-1"
            onSuccess={jest.fn()}
            onCancel={jest.fn()}
          />
        </Provider>
      );

      // Fill out the form without match number (should auto-generate)
      const dateInput = screen.getByLabelText(/Date \*/);
      const playerCountInput = screen.getByLabelText(/Team Player Count \*/);
      const oversInput = screen.getByLabelText(/Total Overs \*/);

      fireEvent.change(dateInput, { target: { value: '2024-01-01' } });
      fireEvent.change(playerCountInput, { target: { value: '11' } });
      fireEvent.change(oversInput, { target: { value: '20' } });

      // Submit the form
      const submitButton = screen.getByText('Match');
      fireEvent.click(submitButton);

      // Wait for the Redux action to be dispatched with match_number: 1 (default for auto-generation)
      await waitFor(() => {
        expect(createMatchRequest).toHaveBeenCalledWith({
          series_id: 'series-1',
          match_number: 1,
          date: '2024-01-01T00:00:00Z',
          status: 'live',
          team_a_player_count: 11,
          team_b_player_count: 11,
          total_overs: 20,
          toss_winner: 'A',
          toss_type: 'H',
          batting_team: 'A',
        });
      });
    });
  });

  describe('Error Handling Integration', () => {
    it('should handle API errors gracefully', async () => {
      const mockStore = createMockStore({
        match: {
          matches: [],
          currentMatch: null,
          loading: false,
          error: null,
        },
      });

      (apiService.createMatch as jest.Mock).mockRejectedValue(
        new Error('Network error')
      );

      render(
        <Provider store={mockStore}>
          <MatchForm
            seriesId="series-1"
            onSuccess={jest.fn()}
            onCancel={jest.fn()}
          />
        </Provider>
      );

      // Fill out the form
      const dateInput = screen.getByLabelText(/Date \*/);
      const playerCountInput = screen.getByLabelText(/Team Player Count \*/);
      const oversInput = screen.getByLabelText(/Total Overs \*/);

      fireEvent.change(dateInput, { target: { value: '2024-01-01' } });
      fireEvent.change(playerCountInput, { target: { value: '11' } });
      fireEvent.change(oversInput, { target: { value: '20' } });

      // Submit the form
      const submitButton = screen.getByText('Match');
      fireEvent.click(submitButton);

      // Verify Redux action was dispatched (the saga will handle the API error)
      await waitFor(() => {
        expect(createMatchRequest).toHaveBeenCalled();
      });
    });

    it('should handle validation errors', async () => {
      const mockStore = createMockStore({
        match: {
          matches: [],
          currentMatch: null,
          loading: false,
          error: null,
        },
      });

      render(
        <Provider store={mockStore}>
          <MatchForm
            seriesId="series-1"
            onSuccess={jest.fn()}
            onCancel={jest.fn()}
          />
        </Provider>
      );

      // Try to submit without filling required fields
      const submitButton = screen.getByText('Match');
      fireEvent.click(submitButton);

      // Form should not submit due to validation
      expect(apiService.createMatch).not.toHaveBeenCalled();
    });

    it('should handle invalid data ranges', async () => {
      const mockStore = createMockStore({
        match: {
          matches: [],
          currentMatch: null,
          loading: false,
          error: null,
        },
      });

      render(
        <Provider store={mockStore}>
          <MatchForm
            seriesId="series-1"
            onSuccess={jest.fn()}
            onCancel={jest.fn()}
          />
        </Provider>
      );

      // Fill out form with invalid data
      const dateInput = screen.getByLabelText(/Date \*/);
      const playerCountInput = screen.getByLabelText(/Team Player Count \*/);
      const oversInput = screen.getByLabelText(/Total Overs \*/);
      const submitButton = screen.getByText('Match');

      fireEvent.change(dateInput, { target: { value: '2024-01-01' } });
      fireEvent.change(playerCountInput, { target: { value: '15' } }); // Invalid: > 11
      fireEvent.change(oversInput, { target: { value: '25' } }); // Invalid: > 20
      fireEvent.click(submitButton);

      // Form should not submit due to validation
      expect(apiService.createMatch).not.toHaveBeenCalled();
    });
  });

  describe('Loading States Integration', () => {
    it('should show loading state during API calls', async () => {
      const mockStore = createMockStore({
        match: {
          matches: [],
          currentMatch: null,
          loading: true,
          error: null,
        },
      });

      render(
        <Provider store={mockStore}>
          <MatchForm
            seriesId="series-1"
            onSuccess={jest.fn()}
            onCancel={jest.fn()}
          />
        </Provider>
      );

      expect(screen.getByText('Saving...')).toBeInTheDocument();
      expect(screen.getByText('Saving...')).toBeDisabled();
    });

    it('should disable form during submission', async () => {
      const mockStore = createMockStore({
        match: {
          matches: [],
          currentMatch: null,
          loading: true,
          error: null,
        },
      });

      render(
        <Provider store={mockStore}>
          <MatchForm
            seriesId="series-1"
            onSuccess={jest.fn()}
            onCancel={jest.fn()}
          />
        </Provider>
      );

      const submitButton = screen.getByText('Saving...');
      expect(submitButton).toBeDisabled();
    });
  });

  describe('Form State Management Integration', () => {
    it('should render select components', () => {
      const mockStore = createMockStore({
        match: {
          matches: [],
          currentMatch: null,
          loading: false,
          error: null,
        },
      });

      render(
        <Provider store={mockStore}>
          <MatchForm
            seriesId="series-1"
            onSuccess={jest.fn()}
            onCancel={jest.fn()}
          />
        </Provider>
      );

      // Verify select components are rendered
      const comboboxes = screen.getAllByRole('combobox');
      expect(comboboxes).toHaveLength(2); // Toss winner and toss type selects
    });

    it('should populate form with existing match data', () => {
      const mockMatch: Match = {
        id: '1',
        series_id: 'series-1',
        match_number: 2,
        date: '2024-02-01T00:00:00Z',
        status: 'live',
        team_a_player_count: 10,
        team_b_player_count: 10,
        total_overs: 15,
        toss_winner: 'B',
        toss_type: 'T',
        batting_team: 'B',
        created_at: '2024-02-01T00:00:00Z',
        updated_at: '2024-02-01T00:00:00Z',
      };

      const mockStore = createMockStore({
        match: {
          matches: [],
          currentMatch: null,
          loading: false,
          error: null,
        },
      });

      render(
        <Provider store={mockStore}>
          <MatchForm
            match={mockMatch}
            onSuccess={jest.fn()}
            onCancel={jest.fn()}
          />
        </Provider>
      );

      // Verify basic fields are populated correctly
      expect(screen.getByDisplayValue('2024-02-01')).toBeInTheDocument();
      expect(screen.getByDisplayValue('10')).toBeInTheDocument();
      expect(screen.getByDisplayValue('15')).toBeInTheDocument();
      expect(screen.getByDisplayValue('2')).toBeInTheDocument();
    });
  });
});
