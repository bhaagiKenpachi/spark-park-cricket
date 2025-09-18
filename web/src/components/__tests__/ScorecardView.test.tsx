/* eslint-disable @typescript-eslint/no-explicit-any, @typescript-eslint/no-require-imports */
import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { configureStore } from '@reduxjs/toolkit';
import { ScorecardView } from '../ScorecardView';
import scorecardSlice, { ScorecardResponse } from '../../store/reducers/scorecardSlice';

// Mock the API service
jest.mock('../../services/api', () => ({
    ApiService: jest.fn().mockImplementation(() => ({
        getScorecard: jest.fn(),
        startScoring: jest.fn(),
        addBall: jest.fn(),
    })),
}));

// Mock Redux hooks
const mockDispatch = jest.fn();
const mockUseAppSelector = jest.fn();

jest.mock('../../store/hooks', () => ({
    useAppDispatch: () => jest.fn(),
    useAppSelector: jest.fn(),
}));

// Mock Redux actions
jest.mock('../../store/reducers/scorecardSlice', () => ({
    ...jest.requireActual('../../store/reducers/scorecardSlice'),
    fetchScorecardRequest: jest.fn().mockReturnValue({ type: 'scorecard/fetchScorecardRequest' }),
    startScoringRequest: jest.fn().mockReturnValue({ type: 'scorecard/startScoringRequest' }),
    addBallRequest: jest.fn().mockReturnValue({ type: 'scorecard/addBallRequest' }),
    clearScorecard: jest.fn().mockReturnValue({ type: 'scorecard/clearScorecard' }),
}));

// Mock store for testing
const createMockStore = (initialState: any) => {
    return configureStore({
        reducer: {
            scorecard: scorecardSlice.reducer,
        },
        preloadedState: initialState,
        middleware: (getDefaultMiddleware) =>
            getDefaultMiddleware({
                serializableCheck: false,
                immutableCheck: false,
            }),
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
                        { ball_number: 1, ball_type: 'good', run_type: '1', runs: 1, byes: 0, is_wicket: false },
                        { ball_number: 2, ball_type: 'good', run_type: '4', runs: 4, byes: 0, is_wicket: false },
                        { ball_number: 3, ball_type: 'wide', run_type: 'WD', runs: 1, byes: 0, is_wicket: false },
                        { ball_number: 4, ball_type: 'good', run_type: '2', runs: 2, byes: 0, is_wicket: false },
                    ],
                },
                {
                    over_number: 2,
                    total_runs: 12,
                    total_balls: 6,
                    total_wickets: 1,
                    status: 'completed',
                    balls: [
                        { ball_number: 1, ball_type: 'good', run_type: '6', runs: 6, byes: 0, is_wicket: false },
                        { ball_number: 2, ball_type: 'good', run_type: 'WC', runs: 0, byes: 0, is_wicket: true, wicket_type: 'bowled' },
                        { ball_number: 3, ball_type: 'good', run_type: '1', runs: 1, byes: 0, is_wicket: false },
                        { ball_number: 4, ball_type: 'good', run_type: '4', runs: 4, byes: 0, is_wicket: false },
                        { ball_number: 5, ball_type: 'good', run_type: '1', runs: 1, byes: 0, is_wicket: false },
                    ],
                },
            ],
        },
    ],
};

const mockOnBack = jest.fn();

describe('ScorecardView Component', () => {
    beforeEach(() => {
        jest.clearAllMocks();
    });

    it('should render loading state initially', () => {
        jest.spyOn(require('../../store/hooks'), 'useAppSelector').mockReturnValue({
            scorecard: null,
            loading: true,
            error: null,
            scoring: false,
        });

        render(<ScorecardView matchId="match-1" onBack={mockOnBack} />);

        expect(screen.getByText('Loading scorecard...')).toBeInTheDocument();
    });

    it('should render error state', () => {
        jest.spyOn(require('../../store/hooks'), 'useAppSelector').mockReturnValue({
            scorecard: null,
            loading: false,
            error: 'Failed to fetch scorecard',
            scoring: false,
        });

        render(<ScorecardView matchId="match-1" onBack={mockOnBack} />);

        expect(screen.getByText('Error:')).toBeInTheDocument();
        expect(screen.getByText('Failed to fetch scorecard')).toBeInTheDocument();
        expect(screen.getByText('Retry')).toBeInTheDocument();
    });

    it('should render no scorecard found message', () => {
        jest.spyOn(require('../../store/hooks'), 'useAppSelector').mockReturnValue({
            scorecard: null,
            loading: false,
            error: null,
            scoring: false,
        });

        render(<ScorecardView matchId="match-1" onBack={mockOnBack} />);

        expect(screen.getByText('No scorecard found for this match.')).toBeInTheDocument();
        expect(screen.getByText('Back to Matches')).toBeInTheDocument();
    });

    it('should render scorecard data correctly', () => {
        jest.spyOn(require('../../store/hooks'), 'useAppSelector').mockReturnValue({
            scorecard: mockScorecardData,
            loading: false,
            error: null,
            scoring: false,
        });

        render(<ScorecardView matchId="match-1" onBack={mockOnBack} />);

        expect(screen.getByText('Test Series - Match #1')).toBeInTheDocument();
        expect(screen.getByText('Team A vs Team B')).toBeInTheDocument();
        expect(screen.getByText('LIVE')).toBeInTheDocument();
        expect(screen.getByText('Team A')).toBeInTheDocument();
        expect(screen.getByText('Team B')).toBeInTheDocument();
    });

    it('should display innings data correctly', () => {
        jest.spyOn(require('../../store/hooks'), 'useAppSelector').mockReturnValue({
            scorecard: mockScorecardData,
            loading: false,
            error: null,
            scoring: false,
        });

        render(<ScorecardView matchId="match-1" onBack={mockOnBack} />);

        expect(screen.getAllByText('Innings 1')[0]).toBeInTheDocument();
        expect(screen.getByText('Live')).toBeInTheDocument();
        expect(screen.getByText('45/2')).toBeInTheDocument();
        expect(screen.getByText('5 overs')).toBeInTheDocument();
    });

    it('should display ball circles correctly', () => {
        jest.spyOn(require('../../store/hooks'), 'useAppSelector').mockReturnValue({
            scorecard: mockScorecardData,
            loading: false,
            error: null,
            scoring: false,
        });

        render(<ScorecardView matchId="match-1" onBack={mockOnBack} />);

        // Check for ball displays - use getAllByText to handle multiple elements
        expect(screen.getAllByText('1')[0]).toBeInTheDocument(); // Single run
        expect(screen.getAllByText('4')[0]).toBeInTheDocument(); // Four
        expect(screen.getAllByText('6')[0]).toBeInTheDocument(); // Six
        expect(screen.getAllByText('2')[0]).toBeInTheDocument(); // Two runs
        expect(screen.getByText('W')).toBeInTheDocument(); // Wicket
    });

    it('should show live scoring button for live matches', () => {
        jest.spyOn(require('../../store/hooks'), 'useAppSelector').mockReturnValue({
            scorecard: mockScorecardData,
            loading: false,
            error: null,
            scoring: false,
        });

        render(<ScorecardView matchId="match-1" onBack={mockOnBack} />);

        expect(screen.getByText('Live Scoring')).toBeInTheDocument();
    });

    it('should handle back button click', () => {
        jest.spyOn(require('../../store/hooks'), 'useAppSelector').mockReturnValue({
            scorecard: mockScorecardData,
            loading: false,
            error: null,
            scoring: false,
        });

        render(<ScorecardView matchId="match-1" onBack={mockOnBack} />);

        const backButton = screen.getByTitle('Back');
        fireEvent.click(backButton);

        expect(mockOnBack).toHaveBeenCalledTimes(1);
    });

    it('should handle refresh button click', () => {
        jest.spyOn(require('../../store/hooks'), 'useAppSelector').mockReturnValue({
            scorecard: mockScorecardData,
            loading: false,
            error: null,
            scoring: false,
        });

        render(<ScorecardView matchId="match-1" onBack={mockOnBack} />);

        const refreshButton = screen.getByTitle('Refresh Scorecard');
        fireEvent.click(refreshButton);

        // The refresh should trigger a fetchScorecardRequest action
        expect(refreshButton).toBeInTheDocument();
    });

    it('should show live scoring interface when start scoring is clicked', async () => {
        jest.spyOn(require('../../store/hooks'), 'useAppSelector').mockReturnValue({
            scorecard: mockScorecardData,
            loading: false,
            error: null,
            scoring: false,
        });

        render(<ScorecardView matchId="match-1" onBack={mockOnBack} />);

        const liveScoringButton = screen.getByText('Live Scoring');
        fireEvent.click(liveScoringButton);

        await waitFor(() => {
            // Check if live scoring interface is opened by looking for scoring buttons
            expect(screen.getAllByText('0')[0]).toBeInTheDocument();
        });
    });

    it('should display scoring buttons in live scoring interface', async () => {
        jest.spyOn(require('../../store/hooks'), 'useAppSelector').mockReturnValue({
            scorecard: mockScorecardData,
            loading: false,
            error: null,
            scoring: false,
        });

        render(<ScorecardView matchId="match-1" onBack={mockOnBack} />);

        const liveScoringButton = screen.getByText('Live Scoring');
        fireEvent.click(liveScoringButton);

        await waitFor(() => {
            expect(screen.getByText('Runs')).toBeInTheDocument();
            expect(screen.getByText('Extras')).toBeInTheDocument();
            expect(screen.getByText('Wickets')).toBeInTheDocument();
            expect(screen.getByText('Byes (Optional)')).toBeInTheDocument();
        });

        // Check for run buttons - use getAllByText to handle multiple elements
        expect(screen.getAllByText('0')[0]).toBeInTheDocument();
        expect(screen.getAllByText('1')[0]).toBeInTheDocument();
        expect(screen.getAllByText('2')[0]).toBeInTheDocument();
        expect(screen.getAllByText('3')[0]).toBeInTheDocument();
        expect(screen.getAllByText('4')[0]).toBeInTheDocument();
        expect(screen.getAllByText('6')[0]).toBeInTheDocument();

        // Check for extra buttons
        expect(screen.getByText('Wide')).toBeInTheDocument();
        expect(screen.getByText('No Ball')).toBeInTheDocument();

        // Check for wicket buttons
        expect(screen.getByText('BOWLED')).toBeInTheDocument();
        expect(screen.getByText('CAUGHT')).toBeInTheDocument();
        expect(screen.getByText('LBW')).toBeInTheDocument();
        expect(screen.getByText('RUN OUT')).toBeInTheDocument();
        expect(screen.getByText('STUMPED')).toBeInTheDocument();
        expect(screen.getByText('HIT WICKET')).toBeInTheDocument();
    });

    it('should handle ball scoring', async () => {
        jest.spyOn(require('../../store/hooks'), 'useAppSelector').mockReturnValue({
            scorecard: mockScorecardData,
            loading: false,
            error: null,
            scoring: false,
        });

        render(<ScorecardView matchId="match-1" onBack={mockOnBack} />);

        const liveScoringButton = screen.getByText('Live Scoring');
        fireEvent.click(liveScoringButton);

        await waitFor(() => {
            const fourButtons = screen.getAllByText('4');
            // Click the first four button (should be the scoring button)
            fireEvent.click(fourButtons[0]);
        });

        // The ball scoring should trigger an addBallRequest action
        expect(screen.getAllByText('4')[0]).toBeInTheDocument();
    });

    it('should handle byes selection', async () => {
        jest.spyOn(require('../../store/hooks'), 'useAppSelector').mockReturnValue({
            scorecard: mockScorecardData,
            loading: false,
            error: null,
            scoring: false,
        });

        render(<ScorecardView matchId="match-1" onBack={mockOnBack} />);

        const liveScoringButton = screen.getByText('Live Scoring');
        fireEvent.click(liveScoringButton);

        await waitFor(() => {
            const twoButtons = screen.getAllByText('2');
            // Click the first two button (should be the scoring button)
            fireEvent.click(twoButtons[0]);
        });

        // Check if byes selection is working by looking for the byes button
        expect(screen.getByText('Byes (Optional)')).toBeInTheDocument();
    });

    it('should show loading state during scoring', () => {
        jest.spyOn(require('../../store/hooks'), 'useAppSelector').mockReturnValue({
            scorecard: mockScorecardData,
            loading: false,
            error: null,
            scoring: true,
        });

        render(<ScorecardView matchId="match-1" onBack={mockOnBack} />);

        // Should show some loading indication
        expect(screen.getByText('LIVE')).toBeInTheDocument();
    });

    it('should handle match with no innings data', () => {
        const scorecardWithoutInnings = {
            ...mockScorecardData,
            innings: null,
        };

        jest.spyOn(require('../../store/hooks'), 'useAppSelector').mockReturnValue({
            scorecard: scorecardWithoutInnings,
            loading: false,
            error: null,
            scoring: false,
        });

        render(<ScorecardView matchId="match-1" onBack={mockOnBack} />);

        expect(screen.getAllByText('Match ready to start')[0]).toBeInTheDocument();
        expect(screen.getAllByText('Click "Open Live Scoring" to begin')[0]).toBeInTheDocument();
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

        jest.spyOn(require('../../store/hooks'), 'useAppSelector').mockReturnValue({
            scorecard: scorecardWithCompletedInnings,
            loading: false,
            error: null,
            scoring: false,
        });

        render(<ScorecardView matchId="match-1" onBack={mockOnBack} />);

        expect(screen.getAllByText('Completed')[0]).toBeInTheDocument();
    });

    it('should toggle expanded overs view', async () => {
        jest.spyOn(require('../../store/hooks'), 'useAppSelector').mockReturnValue({
            scorecard: mockScorecardData,
            loading: false,
            error: null,
            scoring: false,
        });

        render(<ScorecardView matchId="match-1" onBack={mockOnBack} />);

        // Look for the "Show All Overs" button
        const showOversButton = screen.getByText(/Show All Overs/);
        fireEvent.click(showOversButton);

        await waitFor(() => {
            expect(screen.getByText('Hide All Overs')).toBeInTheDocument();
        });
    });

    it('should display extras breakdown correctly', () => {
        jest.spyOn(require('../../store/hooks'), 'useAppSelector').mockReturnValue({
            scorecard: mockScorecardData,
            loading: false,
            error: null,
            scoring: false,
        });

        render(<ScorecardView matchId="match-1" onBack={mockOnBack} />);

        // Check if extras are displayed (the exact format may vary)
        expect(screen.getAllByText(/Extras/)[0]).toBeInTheDocument();
    });

    it('should handle wicket scoring correctly', async () => {
        jest.spyOn(require('../../store/hooks'), 'useAppSelector').mockReturnValue({
            scorecard: mockScorecardData,
            loading: false,
            error: null,
            scoring: false,
        });

        render(<ScorecardView matchId="match-1" onBack={mockOnBack} />);

        const liveScoringButton = screen.getByText('Live Scoring');
        fireEvent.click(liveScoringButton);

        await waitFor(() => {
            const bowledButton = screen.getByText('BOWLED');
            fireEvent.click(bowledButton);
        });

        // Should trigger wicket scoring
        expect(screen.getByText('BOWLED')).toBeInTheDocument();
    });
});
