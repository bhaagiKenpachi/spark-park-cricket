import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { ScorecardView } from '../../components/ScorecardView';
import { ApiService } from '../../services/api';

// Mock the services
jest.mock('../../services/api');
jest.mock('../../services/graphqlService', () => ({
  graphqlService: {
    getInningsScoreSummary: jest.fn(),
    getLatestOverOnly: jest.fn(),
    getLiveScorecard: jest.fn(),
  },
}));

const mockApiService = ApiService as jest.MockedClass<typeof ApiService>;
const mockGraphqlService = {
  getInningsScoreSummary: jest.fn(),
  getLatestOverOnly: jest.fn(),
  getLiveScorecard: jest.fn(),
};

// Mock the useAppDispatch and useAppSelector hooks
const mockDispatch = jest.fn();
const mockUseAppSelector = jest.fn();

jest.mock('../../store/hooks', () => ({
  useAppDispatch: () => mockDispatch,
  useAppSelector: (selector: (state: unknown) => unknown) =>
    mockUseAppSelector(selector),
}));

// Mock Apollo Client
jest.mock('@apollo/client', () => ({
  ApolloClient: jest.fn(),
  InMemoryCache: jest.fn(),
  createHttpLink: jest.fn(),
  gql: jest.fn(query => query),
}));

describe('Ball Addition Workflow Integration', () => {
  const mockMatchId = 'test-match-id';

  const initialScorecardData = {
    match_id: mockMatchId,
    series_id: 'test-series',
    match_number: 1,
    date: '2025-01-01',
    status: 'live',
    match_status: 'live', // Add the missing field
    team_a: 'Team A',
    team_b: 'Team B',
    innings: [
      {
        innings_number: 1,
        batting_team: 'A',
        total_runs: 0,
        total_wickets: 0,
        total_overs: 0,
        total_balls: 0,
        status: 'in_progress',
        extras: { total: 0 },
        overs: [],
      },
    ],
  };

  beforeEach(() => {
    jest.clearAllMocks();

    // Mock successful API responses
    const mockApiInstance = {
      addBall: jest.fn().mockResolvedValue({ success: true }),
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
    } as any;
    mockApiService.mockImplementation(() => mockApiInstance);

    // Mock successful GraphQL responses
    mockGraphqlService.getInningsScoreSummary.mockResolvedValue({
      success: true,
      data: {
        innings_number: 1,
        batting_team: 'A',
        total_runs: 4,
        total_wickets: 0,
        total_overs: 1,
        total_balls: 1,
        status: 'in_progress',
        extras: { total: 0 },
      },
    });

    mockGraphqlService.getLatestOverOnly.mockResolvedValue({
      success: true,
      data: {
        over_number: 1,
        total_runs: 4,
        total_balls: 1,
        total_wickets: 0,
        status: 'in_progress',
        balls: [
          {
            ball_number: 1,
            ball_type: 'good',
            run_type: 'runs',
            runs: 4,
            byes: 0,
            is_wicket: false,
            wicket_type: null,
          },
        ],
      },
    });

    // Default mock selector implementation
    mockUseAppSelector.mockImplementation(selector => {
      const mockState = {
        scorecard: {
          scorecard: initialScorecardData,
          latestOver: null,
          loading: false,
          error: null,
        },
      };
      return selector(mockState);
    });
  });

  it('should complete full ball addition workflow successfully', async () => {
    render(<ScorecardView matchId={mockMatchId} onBack={() => {}} />);

    // Initial state - no latest over data
    expect(screen.getByText('Latest Over: No over data')).toBeInTheDocument();

    // Click on 4 runs button
    const runButton = screen.getByText('4');
    fireEvent.click(runButton);

    // Verify that addBallRequest was dispatched
    expect(mockDispatch).toHaveBeenCalledWith({
      type: 'scorecard/addBallRequest',
      payload: {
        match_id: mockMatchId,
        innings_number: 1,
        ball_type: 'good',
        run_type: 'runs',
        runs: 4,
        byes: 0,
        is_wicket: false,
      },
    });

    // Simulate the saga execution by updating the mock state
    mockUseAppSelector.mockImplementation(selector => {
      const mockState = {
        scorecard: {
          scorecard: {
            ...initialScorecardData,
            innings: [
              {
                ...initialScorecardData.innings[0],
                total_runs: 4,
                total_overs: 1,
                total_balls: 1,
                overs: [
                  {
                    over_number: 1,
                    total_runs: 4,
                    total_balls: 1,
                    total_wickets: 0,
                    status: 'in_progress',
                    balls: [
                      {
                        ball_number: 1,
                        ball_type: 'good',
                        run_type: 'runs',
                        runs: 4,
                        byes: 0,
                        is_wicket: false,
                        wicket_type: null,
                      },
                    ],
                  },
                ],
              },
            ],
          },
          latestOver: {
            over_number: 1,
            total_runs: 4,
            total_balls: 1,
            total_wickets: 0,
            status: 'in_progress',
            balls: [
              {
                ball_number: 1,
                ball_type: 'good',
                run_type: 'runs',
                runs: 4,
                byes: 0,
                is_wicket: false,
                wicket_type: null,
              },
            ],
          },
          loading: false,
          error: null,
        },
      };
      return selector(mockState);
    });

    // Re-render to reflect the updated state
    render(<ScorecardView matchId={mockMatchId} onBack={() => {}} />);

    // Verify that the latest over data is now displayed
    expect(screen.getByText('Latest Over: 1')).toBeInTheDocument();
    expect(screen.getByText('4 runs, 1 balls')).toBeInTheDocument();
    expect(screen.getByText('4/0')).toBeInTheDocument();
  });

  it('should handle sequential ball additions correctly', async () => {
    render(<ScorecardView matchId={mockMatchId} onBack={() => {}} />);

    // Add first ball (4 runs)
    const runButton4 = screen.getByText('4');
    fireEvent.click(runButton4);

    // Update mock state for first ball
    mockUseAppSelector.mockImplementation(selector => {
      const mockState = {
        scorecard: {
          scorecard: {
            ...initialScorecardData,
            innings: [
              {
                ...initialScorecardData.innings[0],
                total_runs: 4,
                total_overs: 1,
                total_balls: 1,
                overs: [
                  {
                    over_number: 1,
                    total_runs: 4,
                    total_balls: 1,
                    total_wickets: 0,
                    status: 'in_progress',
                    balls: [
                      {
                        ball_number: 1,
                        ball_type: 'good',
                        run_type: 'runs',
                        runs: 4,
                        byes: 0,
                        is_wicket: false,
                        wicket_type: null,
                      },
                    ],
                  },
                ],
              },
            ],
          },
          latestOver: {
            over_number: 1,
            total_runs: 4,
            total_balls: 1,
            total_wickets: 0,
            status: 'in_progress',
            balls: [
              {
                ball_number: 1,
                ball_type: 'good',
                run_type: 'runs',
                runs: 4,
                byes: 0,
                is_wicket: false,
                wicket_type: null,
              },
            ],
          },
          loading: false,
          error: null,
        },
      };
      return selector(mockState);
    });

    // Re-render and verify first ball
    render(<ScorecardView matchId={mockMatchId} onBack={() => {}} />);
    expect(screen.getByText('4/0')).toBeInTheDocument();
    expect(screen.getByText('4 runs, 1 balls')).toBeInTheDocument();

    // Add second ball (2 runs)
    const runButton2 = screen.getByText('2');
    fireEvent.click(runButton2);

    // Update mock state for second ball
    mockUseAppSelector.mockImplementation(selector => {
      const mockState = {
        scorecard: {
          scorecard: {
            ...initialScorecardData,
            innings: [
              {
                ...initialScorecardData.innings[0],
                total_runs: 6,
                total_overs: 1,
                total_balls: 2,
                overs: [
                  {
                    over_number: 1,
                    total_runs: 6,
                    total_balls: 2,
                    total_wickets: 0,
                    status: 'in_progress',
                    balls: [
                      {
                        ball_number: 1,
                        ball_type: 'good',
                        run_type: 'runs',
                        runs: 4,
                        byes: 0,
                        is_wicket: false,
                        wicket_type: null,
                      },
                      {
                        ball_number: 2,
                        ball_type: 'good',
                        run_type: 'runs',
                        runs: 2,
                        byes: 0,
                        is_wicket: false,
                        wicket_type: null,
                      },
                    ],
                  },
                ],
              },
            ],
          },
          latestOver: {
            over_number: 1,
            total_runs: 6,
            total_balls: 2,
            total_wickets: 0,
            status: 'in_progress',
            balls: [
              {
                ball_number: 1,
                ball_type: 'good',
                run_type: 'runs',
                runs: 4,
                byes: 0,
                is_wicket: false,
                wicket_type: null,
              },
              {
                ball_number: 2,
                ball_type: 'good',
                run_type: 'runs',
                runs: 2,
                byes: 0,
                is_wicket: false,
                wicket_type: null,
              },
            ],
          },
          loading: false,
          error: null,
        },
      };
      return selector(mockState);
    });

    // Re-render and verify second ball
    render(<ScorecardView matchId={mockMatchId} onBack={() => {}} />);
    expect(screen.getByText('6/0')).toBeInTheDocument();
    expect(screen.getByText('6 runs, 2 balls')).toBeInTheDocument();
  });

  it('should handle innings completion and transition correctly', async () => {
    // Mock completed innings response
    mockGraphqlService.getInningsScoreSummary.mockResolvedValue({
      success: true,
      data: {
        innings_number: 1,
        batting_team: 'A',
        total_runs: 120,
        total_wickets: 10,
        total_overs: 20,
        total_balls: 120,
        status: 'completed',
        extras: { total: 5 },
      },
    });

    render(<ScorecardView matchId={mockMatchId} onBack={() => {}} />);

    // Add a ball that completes the innings
    const runButton = screen.getByText('4');
    fireEvent.click(runButton);

    // Verify that the saga would dispatch fetchUpdatedScorecardRequest for innings completion
    // This would be handled by the saga when it detects status: 'completed'
    expect(mockDispatch).toHaveBeenCalledWith({
      type: 'scorecard/addBallRequest',
      payload: expect.objectContaining({
        match_id: mockMatchId,
        innings_number: 1,
      }),
    });
  });

  it('should handle API failure gracefully', async () => {
    // Mock API failure
    const mockApiInstance = {
      addBall: jest.fn().mockRejectedValue(new Error('API Error')),
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
    } as any;
    mockApiService.mockImplementation(() => mockApiInstance);

    render(<ScorecardView matchId={mockMatchId} onBack={() => {}} />);

    const runButton = screen.getByText('4');
    fireEvent.click(runButton);

    // Verify that addBallRequest was still dispatched
    expect(mockDispatch).toHaveBeenCalledWith({
      type: 'scorecard/addBallRequest',
      payload: expect.objectContaining({
        match_id: mockMatchId,
        innings_number: 1,
      }),
    });

    // The saga would handle the API error and dispatch addBallFailure
    // This would be tested in the saga unit tests
  });

  it('should handle GraphQL errors gracefully without affecting ball addition', async () => {
    // Mock GraphQL failure
    mockGraphqlService.getInningsScoreSummary.mockRejectedValue(
      new Error('GraphQL Error')
    );

    render(<ScorecardView matchId={mockMatchId} onBack={() => {}} />);

    const runButton = screen.getByText('4');
    fireEvent.click(runButton);

    // Verify that addBallRequest was still dispatched
    expect(mockDispatch).toHaveBeenCalledWith({
      type: 'scorecard/addBallRequest',
      payload: expect.objectContaining({
        match_id: mockMatchId,
        innings_number: 1,
      }),
    });

    // The saga would handle GraphQL errors gracefully and still mark ball addition as successful
    // This ensures that GraphQL data fetch failures don't affect the core ball addition functionality
  });

  it('should handle refresh functionality correctly', async () => {
    render(<ScorecardView matchId={mockMatchId} onBack={() => {}} />);

    const refreshButton = screen.getByTitle('Refresh Scorecard');
    fireEvent.click(refreshButton);

    // Verify that all refresh actions are dispatched
    expect(mockDispatch).toHaveBeenCalledWith({
      type: 'scorecard/fetchUpdatedScorecardRequest',
      payload: mockMatchId,
    });
    expect(mockDispatch).toHaveBeenCalledWith({
      type: 'scorecard/fetchInningsScoreSummaryRequest',
      payload: { matchId: mockMatchId, inningsNumber: 1 },
    });
    expect(mockDispatch).toHaveBeenCalledWith({
      type: 'scorecard/fetchLatestOverRequest',
      payload: { matchId: mockMatchId, inningsNumber: 1 },
    });
  });

  it('should prevent scoring on completed match', async () => {
    const completedMatchData = {
      ...initialScorecardData,
      match_status: 'completed', // Ensure this field exists
      innings: [
        {
          ...initialScorecardData.innings[0],
          status: 'completed',
        },
        {
          innings_number: 2,
          batting_team: 'B',
          total_runs: 45,
          total_wickets: 10,
          total_overs: 15,
          total_balls: 90,
          status: 'completed',
          extras: { total: 2 },
          overs: [],
        },
      ],
    };

    mockUseAppSelector.mockImplementation(selector => {
      const mockState = {
        scorecard: {
          scorecard: completedMatchData,
          latestOver: null,
          loading: false,
          error: null,
        },
      };
      return selector(mockState);
    });

    render(<ScorecardView matchId={mockMatchId} onBack={() => {}} />);

    // Should show match completed
    expect(screen.getByText('Match Completed')).toBeInTheDocument();

    // Should disable scoring buttons
    const runButton = screen.getByText('4');
    expect(runButton).toBeDisabled();
  });
});
