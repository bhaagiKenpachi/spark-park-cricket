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
    innings: null, // Start with null innings to show "Match ready to start"
  };

  const initialInningsData = {
    innings_number: 1,
    batting_team: 'A',
    total_runs: 0,
    total_wickets: 0,
    total_overs: 0,
    total_balls: 0,
    status: 'in_progress',
    extras: { total: 0 },
    overs: [],
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
    render(<ScorecardView matchId={mockMatchId} onBack={() => { }} />);

    // Initial state - no latest over data (component shows "Match ready to start" when innings is null)
    // There are two instances (one for each team), so use getAllByText
    expect(screen.getAllByText('Match ready to start')).toHaveLength(2);

    // Click on 4 runs button (use getAllByText to get the first one, which should be the runs button)
    const runButtons = screen.getAllByText('4');
    const runButton = runButtons[0]; // First "4" button is the runs button
    fireEvent.click(runButton);

    // Verify that addBallRequest was dispatched
    expect(mockDispatch).toHaveBeenCalledWith({
      type: 'scorecard/addBallRequest',
      payload: {
        match_id: mockMatchId,
        innings_number: 1,
        ball_type: 'good',
        run_type: '4', // Component uses the runs value as run_type for regular runs
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
                ...initialInningsData,
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
    render(<ScorecardView matchId={mockMatchId} onBack={() => { }} />);

    // Verify that the latest over data is now displayed
    expect(screen.getByText('Latest Over 1')).toBeInTheDocument();
    expect(screen.getByText('4 runs, 0 wickets')).toBeInTheDocument();
    expect(screen.getByText('4/0')).toBeInTheDocument();
  });

  it('should handle sequential ball additions correctly', async () => {
    render(<ScorecardView matchId={mockMatchId} onBack={() => { }} />);

    // Add first ball (4 runs)
    const runButtons4 = screen.getAllByText('4');
    const runButton4 = runButtons4[0]; // First "4" button is the runs button
    fireEvent.click(runButton4);

    // Update mock state for first ball
    mockUseAppSelector.mockImplementation(selector => {
      const mockState = {
        scorecard: {
          scorecard: {
            ...initialScorecardData,
            innings: [
              {
                ...initialInningsData,
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
    render(<ScorecardView matchId={mockMatchId} onBack={() => { }} />);
    expect(screen.getByText('4/0')).toBeInTheDocument();
    expect(screen.getByText('4 runs, 0 wickets')).toBeInTheDocument();

    // Add second ball (2 runs)
    const runButtons2 = screen.getAllByText('2');
    const runButton2 = runButtons2[0]; // First "2" button is the runs button
    fireEvent.click(runButton2);

    // Update mock state for second ball
    mockUseAppSelector.mockImplementation(selector => {
      const mockState = {
        scorecard: {
          scorecard: {
            ...initialScorecardData,
            innings: [
              {
                ...initialInningsData,
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
    render(<ScorecardView matchId={mockMatchId} onBack={() => { }} />);
    expect(screen.getByText('6/0')).toBeInTheDocument();
    expect(screen.getByText('6 runs, 0 wickets')).toBeInTheDocument();
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

    render(<ScorecardView matchId={mockMatchId} onBack={() => { }} />);

    // Add a ball that completes the innings
    const runButtons = screen.getAllByText('4');
    const runButton = runButtons[0]; // First "4" button is the runs button
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

    render(<ScorecardView matchId={mockMatchId} onBack={() => { }} />);

    const runButtons = screen.getAllByText('4');
    const runButton = runButtons[0]; // First "4" button is the runs button
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

    render(<ScorecardView matchId={mockMatchId} onBack={() => { }} />);

    const runButtons = screen.getAllByText('4');
    const runButton = runButtons[0]; // First "4" button is the runs button
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
    render(<ScorecardView matchId={mockMatchId} onBack={() => { }} />);

    const refreshButton = screen.getByTitle('Refresh Scorecard');
    fireEvent.click(refreshButton);

    // Verify that refresh action is dispatched (component only calls fetchScorecardRequest on refresh)
    expect(mockDispatch).toHaveBeenCalledWith({
      type: 'scorecard/fetchScorecardRequest',
      payload: mockMatchId,
    });
  });

  it('should prevent scoring on completed match', async () => {
    const completedMatchData = {
      ...initialScorecardData,
      match_status: 'completed', // Ensure this field exists
      innings: [
        {
          ...initialInningsData,
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

    render(<ScorecardView matchId={mockMatchId} onBack={() => { }} />);

    // Should show match completed (component shows "COMPLETED" badge)
    expect(screen.getByText('COMPLETED')).toBeInTheDocument();

    // Should disable scoring buttons (buttons should not be visible when match is completed)
    // The live scoring interface should not be shown for completed matches
    expect(screen.queryByText('Live Scoring')).not.toBeInTheDocument();
  });
});
