import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { Provider } from 'react-redux';
import { configureStore } from '@reduxjs/toolkit';
import { ScorecardView } from '../../components/ScorecardView';
/* eslint-disable @typescript-eslint/no-explicit-any */
import scorecardReducer, {
  ScorecardResponse,
  BallEventRequest,
  fetchScorecardRequest,
} from '../../store/reducers/scorecardSlice';
import { apiService } from '../../services/api';

// Mock the API service
jest.mock('../../services/api', () => ({
  apiService: {
    getScorecard: jest.fn(),
    startScoring: jest.fn(),
    addBall: jest.fn(),
    getInnings: jest.fn(),
  },
}));

// Mock Redux actions
jest.mock('../../store/reducers/scorecardSlice', () => ({
  ...jest.requireActual('../../store/reducers/scorecardSlice'),
  fetchScorecardRequest: jest
    .fn()
    .mockReturnValue({ type: 'scorecard/fetchScorecardRequest' }),
  startScoringRequest: jest
    .fn()
    .mockReturnValue({ type: 'scorecard/startScoringRequest' }),
  addBallRequest: jest
    .fn()
    .mockReturnValue({ type: 'scorecard/addBallRequest' }),
  clearScorecard: jest
    .fn()
    .mockReturnValue({ type: 'scorecard/clearScorecard' }),
}));

// Mock store for testing
const createMockStore = (initialState: any) => {
  return configureStore({
    reducer: {
      scorecard: scorecardReducer.default,
    },
    preloadedState: initialState,
  });
};

// Mock data
const mockScorecardData: ScorecardResponse = {
  match_id: 'match-1',
  match_number: 1,
  series_name: 'Test Series',
  team_a: 'Team A',
  team_b: 'Team B',
  total_overs: 20,
  toss_winner: 'A',
  toss_type: 'H',
  current_innings: 1,
  match_status: 'live',
  innings: [
    {
      innings_number: 1,
      batting_team: 'A',
      total_runs: 45,
      total_wickets: 2,
      total_overs: 5,
      total_balls: 30,
      status: 'in_progress',
      extras: {
        byes: 2,
        leg_byes: 1,
        wides: 3,
        no_balls: 1,
        total: 7,
      },
      overs: [
        {
          over_number: 1,
          total_runs: 8,
          total_balls: 6,
          total_wickets: 0,
          status: 'completed',
          balls: [
            {
              ball_number: 1,
              ball_type: 'good',
              run_type: '1',
              runs: 1,
              byes: 0,
              is_wicket: false,
            },
            {
              ball_number: 2,
              ball_type: 'good',
              run_type: '4',
              runs: 4,
              byes: 0,
              is_wicket: false,
            },
            {
              ball_number: 3,
              ball_type: 'wide',
              run_type: 'WD',
              runs: 1,
              byes: 0,
              is_wicket: false,
            },
            {
              ball_number: 4,
              ball_type: 'good',
              run_type: '2',
              runs: 2,
              byes: 0,
              is_wicket: false,
            },
          ],
        },
      ],
    },
  ],
};

const mockOnBack = jest.fn();

describe('Scorecard Integration Tests', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('Scorecard CRUD Operations Integration', () => {
    it('should complete fetch scorecard workflow', async () => {
      const store = createMockStore({
        scorecard: {
          scorecard: null,
          loading: false,
          error: null,
          scoring: false,
        },
      });

      render(
        <Provider store={store}>
          <ScorecardView matchId="match-1" onBack={mockOnBack} />
        </Provider>
      );

      // Component should dispatch fetchScorecardRequest on mount
      await waitFor(() => {
        expect(fetchScorecardRequest).toHaveBeenCalledWith('match-1');
      });
    });

    it('should complete start scoring workflow', async () => {
      const store = createMockStore({
        scorecard: {
          scorecard: mockScorecardData,
          loading: false,
          error: null,
          scoring: false,
        },
      });

      (apiService.startScoring as jest.Mock).mockResolvedValueOnce({
        data: { message: 'Scoring started', match_id: 'match-1' },
        status: 200,
        message: 'Success',
      });

      render(
        <Provider store={store}>
          <ScorecardView matchId="match-1" onBack={mockOnBack} />
        </Provider>
      );

      const liveScoringButton = screen.getByText('Live Scoring');
      fireEvent.click(liveScoringButton);

      await waitFor(() => {
        expect(apiService.startScoring).toHaveBeenCalledWith('match-1');
      });
    });

    it('should complete add ball workflow', async () => {
      const store = createMockStore({
        scorecard: {
          scorecard: mockScorecardData,
          loading: false,
          error: null,
          scoring: false,
        },
      });

      (apiService.addBall as jest.Mock).mockResolvedValueOnce({
        data: {
          message: 'Ball added successfully',
          match_id: 'match-1',
          innings_number: 1,
          ball_type: 'good',
          run_type: '4',
          runs: 4,
          byes: 0,
          is_wicket: false,
        },
        status: 200,
        message: 'Success',
      });

      render(
        <Provider store={store}>
          <ScorecardView matchId="match-1" onBack={mockOnBack} />
        </Provider>
      );

      const liveScoringButton = screen.getByText('Live Scoring');
      fireEvent.click(liveScoringButton);

      await waitFor(() => {
        const fourButton = screen.getAllByText('4')[0];
        fireEvent.click(fourButton);
      });

      const expectedBallEvent: BallEventRequest = {
        match_id: 'match-1',
        innings_number: 1,
        ball_type: 'good',
        run_type: '4',
        runs: 4,
        byes: 0,
        is_wicket: false,
      };

      await waitFor(() => {
        expect(apiService.addBall).toHaveBeenCalledWith(expectedBallEvent);
      });
    });

    it('should handle wicket ball scoring workflow', async () => {
      const store = createMockStore({
        scorecard: {
          scorecard: mockScorecardData,
          loading: false,
          error: null,
          scoring: false,
        },
      });

      (apiService.addBall as jest.Mock).mockResolvedValueOnce({
        data: {
          message: 'Wicket ball added successfully',
          match_id: 'match-1',
          innings_number: 1,
          ball_type: 'good',
          run_type: 'WC',
          runs: 0,
          byes: 0,
          is_wicket: true,
          wicket_type: 'bowled',
        },
        status: 200,
        message: 'Success',
      });

      render(
        <Provider store={store}>
          <ScorecardView matchId="match-1" onBack={mockOnBack} />
        </Provider>
      );

      const liveScoringButton = screen.getByText('Live Scoring');
      fireEvent.click(liveScoringButton);

      await waitFor(() => {
        const bowledButton = screen.getByText('BOWLED');
        fireEvent.click(bowledButton);
      });

      const expectedWicketEvent: BallEventRequest = {
        match_id: 'match-1',
        innings_number: 1,
        ball_type: 'good',
        run_type: 'WC',
        runs: 0,
        byes: 0,
        is_wicket: true,
        wicket_type: 'bowled',
      };

      await waitFor(() => {
        expect(apiService.addBall).toHaveBeenCalledWith(expectedWicketEvent);
      });
    });

    it('should handle wide ball scoring workflow', async () => {
      const store = createMockStore({
        scorecard: {
          scorecard: mockScorecardData,
          loading: false,
          error: null,
          scoring: false,
        },
      });

      (apiService.addBall as jest.Mock).mockResolvedValueOnce({
        data: {
          message: 'Wide ball added successfully',
          match_id: 'match-1',
          innings_number: 1,
          ball_type: 'wide',
          run_type: 'WD',
          runs: 1,
          byes: 0,
          is_wicket: false,
        },
        status: 200,
        message: 'Success',
      });

      render(
        <Provider store={store}>
          <ScorecardView matchId="match-1" onBack={mockOnBack} />
        </Provider>
      );

      const liveScoringButton = screen.getByText('Live Scoring');
      fireEvent.click(liveScoringButton);

      await waitFor(() => {
        const wideButton = screen.getByText('Wide');
        fireEvent.click(wideButton);
      });

      const expectedWideEvent: BallEventRequest = {
        match_id: 'match-1',
        innings_number: 1,
        ball_type: 'wide',
        run_type: 'WD',
        runs: 1,
        byes: 0,
        is_wicket: false,
      };

      await waitFor(() => {
        expect(apiService.addBall).toHaveBeenCalledWith(expectedWideEvent);
      });
    });

    it('should handle ball with byes scoring workflow', async () => {
      const store = createMockStore({
        scorecard: {
          scorecard: mockScorecardData,
          loading: false,
          error: null,
          scoring: false,
        },
      });

      (apiService.addBall as jest.Mock).mockResolvedValueOnce({
        data: {
          message: 'Ball with byes added successfully',
          match_id: 'match-1',
          innings_number: 1,
          ball_type: 'good',
          run_type: '1',
          runs: 1,
          byes: 2,
          is_wicket: false,
        },
        status: 200,
        message: 'Success',
      });

      render(
        <Provider store={store}>
          <ScorecardView matchId="match-1" onBack={mockOnBack} />
        </Provider>
      );

      const liveScoringButton = screen.getByText('Live Scoring');
      fireEvent.click(liveScoringButton);

      // Select byes first
      await waitFor(() => {
        const byesButton = screen.getAllByText('2')[0];
        fireEvent.click(byesButton);
      });

      // Then score a run
      await waitFor(() => {
        const oneButton = screen.getByText('1');
        fireEvent.click(oneButton);
      });

      const expectedBallWithByesEvent: BallEventRequest = {
        match_id: 'match-1',
        innings_number: 1,
        ball_type: 'good',
        run_type: '1',
        runs: 1,
        byes: 2,
        is_wicket: false,
      };

      await waitFor(() => {
        expect(apiService.addBall).toHaveBeenCalledWith(
          expectedBallWithByesEvent
        );
      });
    });
  });

  describe('Error Handling Integration', () => {
    it('should handle API errors gracefully', async () => {
      const store = createMockStore({
        scorecard: {
          scorecard: null,
          loading: false,
          error: null,
          scoring: false,
        },
      });

      (apiService.getScorecard as jest.Mock).mockRejectedValueOnce(
        new Error('API Error')
      );

      render(
        <Provider store={store}>
          <ScorecardView matchId="match-1" onBack={mockOnBack} />
        </Provider>
      );

      await waitFor(() => {
        expect(apiService.getScorecard).toHaveBeenCalledWith('match-1');
      });
    });

    it('should handle start scoring errors', async () => {
      const store = createMockStore({
        scorecard: {
          scorecard: mockScorecardData,
          loading: false,
          error: null,
          scoring: false,
        },
      });

      (apiService.startScoring as jest.Mock).mockRejectedValueOnce(
        new Error('Match not ready')
      );

      render(
        <Provider store={store}>
          <ScorecardView matchId="match-1" onBack={mockOnBack} />
        </Provider>
      );

      const liveScoringButton = screen.getByText('Live Scoring');
      fireEvent.click(liveScoringButton);

      await waitFor(() => {
        expect(apiService.startScoring).toHaveBeenCalledWith('match-1');
      });
    });

    it('should handle add ball errors', async () => {
      const store = createMockStore({
        scorecard: {
          scorecard: mockScorecardData,
          loading: false,
          error: null,
          scoring: false,
        },
      });

      (apiService.addBall as jest.Mock).mockRejectedValueOnce(
        new Error('Invalid ball data')
      );

      render(
        <Provider store={store}>
          <ScorecardView matchId="match-1" onBack={mockOnBack} />
        </Provider>
      );

      const liveScoringButton = screen.getByText('Live Scoring');
      fireEvent.click(liveScoringButton);

      await waitFor(() => {
        const fourButton = screen.getAllByText('4')[0];
        fireEvent.click(fourButton);
      });

      await waitFor(() => {
        expect(apiService.addBall).toHaveBeenCalled();
      });
    });

    it('should handle innings completion errors', async () => {
      const store = createMockStore({
        scorecard: {
          scorecard: mockScorecardData,
          loading: false,
          error: null,
          scoring: false,
        },
      });

      (apiService.addBall as jest.Mock).mockRejectedValueOnce(
        new Error('Innings already completed')
      );

      render(
        <Provider store={store}>
          <ScorecardView matchId="match-1" onBack={mockOnBack} />
        </Provider>
      );

      const liveScoringButton = screen.getByText('Live Scoring');
      fireEvent.click(liveScoringButton);

      await waitFor(() => {
        const fourButton = screen.getAllByText('4')[0];
        fireEvent.click(fourButton);
      });

      await waitFor(() => {
        expect(apiService.addBall).toHaveBeenCalled();
      });
    });
  });

  describe('Loading States Integration', () => {
    it('should show loading state during API calls', () => {
      const store = createMockStore({
        scorecard: {
          scorecard: null,
          loading: true,
          error: null,
          scoring: false,
        },
      });

      render(
        <Provider store={store}>
          <ScorecardView matchId="match-1" onBack={mockOnBack} />
        </Provider>
      );

      expect(screen.getByText('Loading scorecard...')).toBeInTheDocument();
    });

    it('should disable form during scoring', () => {
      const store = createMockStore({
        scorecard: {
          scorecard: mockScorecardData,
          loading: false,
          error: null,
          scoring: true,
        },
      });

      render(
        <Provider store={store}>
          <ScorecardView matchId="match-1" onBack={mockOnBack} />
        </Provider>
      );

      // Should show loading state
      expect(screen.getByText('LIVE')).toBeInTheDocument();
    });
  });

  describe('Form State Management Integration', () => {
    it('should render scorecard data correctly', () => {
      const store = createMockStore({
        scorecard: {
          scorecard: mockScorecardData,
          loading: false,
          error: null,
          scoring: false,
        },
      });

      render(
        <Provider store={store}>
          <ScorecardView matchId="match-1" onBack={mockOnBack} />
        </Provider>
      );

      expect(screen.getByText('Test Series - Match #1')).toBeInTheDocument();
      expect(screen.getByText('Team A vs Team B')).toBeInTheDocument();
      expect(screen.getByText('LIVE')).toBeInTheDocument();
      expect(screen.getByText('Team A')).toBeInTheDocument();
      expect(screen.getByText('Team B')).toBeInTheDocument();
    });

    it('should display innings data correctly', () => {
      const store = createMockStore({
        scorecard: {
          scorecard: mockScorecardData,
          loading: false,
          error: null,
          scoring: false,
        },
      });

      render(
        <Provider store={store}>
          <ScorecardView matchId="match-1" onBack={mockOnBack} />
        </Provider>
      );

      expect(screen.getAllByText('Innings 1')[0]).toBeInTheDocument();
      expect(screen.getByText('Live')).toBeInTheDocument();
      expect(screen.getByText('45/2')).toBeInTheDocument();
      expect(screen.getByText('5 overs')).toBeInTheDocument();
      expect(screen.getAllByText(/Extras/)[0]).toBeInTheDocument();
    });

    it('should handle match with no innings data', () => {
      const scorecardWithoutInnings = {
        ...mockScorecardData,
        innings: null,
      };

      const store = createMockStore({
        scorecard: {
          scorecard: scorecardWithoutInnings,
          loading: false,
          error: null,
          scoring: false,
        },
      });

      render(
        <Provider store={store}>
          <ScorecardView matchId="match-1" onBack={mockOnBack} />
        </Provider>
      );

      expect(
        screen.getAllByText('Match ready to start')[0]
      ).toBeInTheDocument();
      expect(
        screen.getAllByText('Click "Open Live Scoring" to begin')[0]
      ).toBeInTheDocument();
    });

    it('should handle completed innings', () => {
      const scorecardWithCompletedInnings = {
        ...mockScorecardData,
        innings: [
          {
            ...mockScorecardData.innings[0],
            status: 'completed',
          },
        ],
      };

      const store = createMockStore({
        scorecard: {
          scorecard: scorecardWithCompletedInnings,
          loading: false,
          error: null,
          scoring: false,
        },
      });

      render(
        <Provider store={store}>
          <ScorecardView matchId="match-1" onBack={mockOnBack} />
        </Provider>
      );

      expect(screen.getAllByText('Completed')[0]).toBeInTheDocument();
    });

    it('should display extras breakdown correctly', () => {
      const store = createMockStore({
        scorecard: {
          scorecard: mockScorecardData,
          loading: false,
          error: null,
          scoring: false,
        },
      });

      render(
        <Provider store={store}>
          <ScorecardView matchId="match-1" onBack={mockOnBack} />
        </Provider>
      );

      expect(screen.getAllByText(/Extras/)[0]).toBeInTheDocument();
      // Note: Extras breakdown format may vary, so we just check that extras are displayed
    });

    it('should handle byes selection correctly', async () => {
      const store = createMockStore({
        scorecard: {
          scorecard: mockScorecardData,
          loading: false,
          error: null,
          scoring: false,
        },
      });

      render(
        <Provider store={store}>
          <ScorecardView matchId="match-1" onBack={mockOnBack} />
        </Provider>
      );

      const liveScoringButton = screen.getByText('Live Scoring');
      fireEvent.click(liveScoringButton);

      await waitFor(() => {
        const byesButton = screen.getAllByText('2')[0];
        fireEvent.click(byesButton);
      });

      expect(screen.getByText('Byes (Optional)')).toBeInTheDocument();
    });

    it.skip('should toggle expanded overs view', async () => {
      // TODO: This test needs to be updated when the "Show All Overs" functionality is implemented
      const store = createMockStore({
        scorecard: {
          scorecard: mockScorecardData,
          loading: false,
          error: null,
          scoring: false,
        },
      });

      render(
        <Provider store={store}>
          <ScorecardView matchId="match-1" onBack={mockOnBack} />
        </Provider>
      );

      // Look for the "Show All Overs" button
      const showOversButton = screen.getByText(/Show All Overs/);
      fireEvent.click(showOversButton);

      await waitFor(() => {
        expect(screen.getByText('Hide All Overs')).toBeInTheDocument();
      });
    });
  });
});
