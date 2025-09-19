describe('Scorecard Management E2E Tests', () => {
  const baseUrl = 'http://localhost:3000';
  const matchId = 'match-1';
  const seriesId = 'series-1';

  beforeEach(() => {
    // Mock API responses for scorecard endpoints
    cy.intercept('GET', `/api/v1/scorecard/${matchId}`, {
      statusCode: 200,
      body: {
        data: {
          match_id: matchId,
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
        },
      },
    }).as('getScorecard');

    cy.intercept('POST', '/api/v1/scorecard/start', {
      statusCode: 200,
      body: {
        message: 'Scoring started successfully',
        match_id: matchId,
      },
    }).as('startScoring');

    cy.intercept('POST', '/api/v1/scorecard/ball', {
      statusCode: 200,
      body: {
        message: 'Ball added successfully',
        match_id: matchId,
        innings_number: 1,
        ball_type: 'good',
        run_type: '4',
        runs: 4,
        byes: 0,
        is_wicket: false,
      },
    }).as('addBall');

    cy.intercept('GET', `/api/v1/scorecard/${matchId}/current-over*`, {
      statusCode: 200,
      body: {
        data: {
          over_number: 5,
          total_runs: 8,
          total_balls: 6,
          total_wickets: 0,
          status: 'in_progress',
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
          ],
        },
      },
    }).as('getCurrentOver');

    cy.intercept('GET', `/api/v1/scorecard/${matchId}/innings/1`, {
      statusCode: 200,
      body: {
        data: {
          innings_number: 1,
          batting_team: 'A',
          total_runs: 120,
          total_wickets: 3,
          total_overs: 10,
          total_balls: 60,
          status: 'completed',
          extras: {
            byes: 5,
            leg_byes: 2,
            wides: 8,
            no_balls: 1,
            total: 16,
          },
          overs: [],
        },
      },
    }).as('getInnings');
  });

  describe('Scorecard Display', () => {
    it('should display scorecard on page load', () => {
      cy.visit(`${baseUrl}/series/${seriesId}/matches/${matchId}/scorecard`);

      cy.wait('@getScorecard');

      cy.get('[data-cy="scorecard-header"]').should(
        'contain',
        'Test Series - Match #1'
      );
      cy.get('[data-cy="teams-display"]').should('contain', 'Team A vs Team B');
      cy.get('[data-cy="match-status"]').should('contain', 'LIVE');
    });

    it('should display team scorecards correctly', () => {
      cy.visit(`${baseUrl}/series/${seriesId}/matches/${matchId}/scorecard`);

      cy.wait('@getScorecard');

      cy.get('[data-cy="team-a-card"]').should('contain', 'Team A');
      cy.get('[data-cy="team-b-card"]').should('contain', 'Team B');
      cy.get('[data-cy="innings-1"]').should('contain', 'Innings 1');
      cy.get('[data-cy="score-display"]').should('contain', '45/2');
      cy.get('[data-cy="overs-display"]').should('contain', '5 overs');
    });

    it('should display ball-by-ball details', () => {
      cy.visit(`${baseUrl}/series/${seriesId}/matches/${matchId}/scorecard`);

      cy.wait('@getScorecard');

      // Check for ball circles
      cy.get('[data-cy="ball-circle"]').should('have.length.at.least', 4);
      cy.get('[data-cy="ball-circle"]').should('contain', '1'); // Single run
      cy.get('[data-cy="ball-circle"]').should('contain', '4'); // Four
      cy.get('[data-cy="ball-circle"]').should('contain', '2'); // Two runs
    });

    it('should display extras breakdown', () => {
      cy.visit(`${baseUrl}/series/${seriesId}/matches/${matchId}/scorecard`);

      cy.wait('@getScorecard');

      cy.get('[data-cy="extras-display"]').should('contain', 'Extras: 7');
      cy.get('[data-cy="extras-breakdown"]').should('contain', '(b 2)');
      cy.get('[data-cy="extras-breakdown"]').should('contain', '(lb 1)');
      cy.get('[data-cy="extras-breakdown"]').should('contain', '(w 3)');
      cy.get('[data-cy="extras-breakdown"]').should('contain', '(nb 1)');
    });
  });

  describe('Live Scoring Interface', () => {
    it('should open live scoring interface', () => {
      cy.visit(`${baseUrl}/series/${seriesId}/matches/${matchId}/scorecard`);

      cy.wait('@getScorecard');

      cy.get('[data-cy="live-scoring-button"]').click();

      cy.get('[data-cy="live-scoring-interface"]').should('be.visible');
      cy.get('[data-cy="runs-section"]').should('contain', 'Runs');
      cy.get('[data-cy="extras-section"]').should('contain', 'Extras');
      cy.get('[data-cy="wickets-section"]').should('contain', 'Wickets');
    });

    it('should score runs correctly', () => {
      cy.visit(`${baseUrl}/series/${seriesId}/matches/${matchId}/scorecard`);

      cy.wait('@getScorecard');

      cy.get('[data-cy="live-scoring-button"]').click();

      // Score a four
      cy.get('[data-cy="run-button-4"]').click();
      cy.wait('@addBall');

      // Verify the API call was made with correct data
      cy.get('@addBall').should('have.been.calledWith', {
        method: 'POST',
        url: '/api/v1/scorecard/ball',
        body: {
          match_id: matchId,
          innings_number: 1,
          ball_type: 'good',
          run_type: '4',
          runs: 4,
          byes: 0,
          is_wicket: false,
        },
      });
    });

    it('should score wickets correctly', () => {
      cy.visit(`${baseUrl}/series/${seriesId}/matches/${matchId}/scorecard`);

      cy.wait('@getScorecard');

      cy.get('[data-cy="live-scoring-button"]').click();

      // Score a wicket
      cy.get('[data-cy="wicket-button-bowled"]').click();

      // Mock the wicket response
      cy.intercept('POST', '/api/v1/scorecard/ball', {
        statusCode: 200,
        body: {
          message: 'Wicket ball added successfully',
          match_id: matchId,
          innings_number: 1,
          ball_type: 'good',
          run_type: 'WC',
          runs: 0,
          byes: 0,
          is_wicket: true,
          wicket_type: 'bowled',
        },
      }).as('addWicketBall');

      cy.wait('@addWicketBall');
    });

    it('should score extras correctly', () => {
      cy.visit(`${baseUrl}/series/${seriesId}/matches/${matchId}/scorecard`);

      cy.wait('@getScorecard');

      cy.get('[data-cy="live-scoring-button"]').click();

      // Score a wide
      cy.get('[data-cy="wide-button"]').click();

      // Mock the wide response
      cy.intercept('POST', '/api/v1/scorecard/ball', {
        statusCode: 200,
        body: {
          message: 'Wide ball added successfully',
          match_id: matchId,
          innings_number: 1,
          ball_type: 'wide',
          run_type: 'WD',
          runs: 1,
          byes: 0,
          is_wicket: false,
        },
      }).as('addWideBall');

      cy.wait('@addWideBall');
    });

    it('should score no ball correctly', () => {
      cy.visit(`${baseUrl}/series/${seriesId}/matches/${matchId}/scorecard`);

      cy.wait('@getScorecard');

      cy.get('[data-cy="live-scoring-button"]').click();

      // Score a no ball
      cy.get('[data-cy="no-ball-button"]').click();

      // Mock the no ball response
      cy.intercept('POST', '/api/v1/scorecard/ball', {
        statusCode: 200,
        body: {
          message: 'No ball added successfully',
          match_id: matchId,
          innings_number: 1,
          ball_type: 'no_ball',
          run_type: 'NB',
          runs: 1,
          byes: 0,
          is_wicket: false,
        },
      }).as('addNoBall');

      cy.wait('@addNoBall');
    });

    it('should handle byes selection', () => {
      cy.visit(`${baseUrl}/series/${seriesId}/matches/${matchId}/scorecard`);

      cy.wait('@getScorecard');

      cy.get('[data-cy="live-scoring-button"]').click();

      // Select byes
      cy.get('[data-cy="byes-button-2"]').click();
      cy.get('[data-cy="byes-selected"]').should('contain', '+2 byes selected');

      // Score a run with byes
      cy.get('[data-cy="run-button-1"]').click();

      // Mock the ball with byes response
      cy.intercept('POST', '/api/v1/scorecard/ball', {
        statusCode: 200,
        body: {
          message: 'Ball with byes added successfully',
          match_id: matchId,
          innings_number: 1,
          ball_type: 'good',
          run_type: '1',
          runs: 1,
          byes: 2,
          is_wicket: false,
        },
      }).as('addBallWithByes');

      cy.wait('@addBallWithByes');
    });

    it('should close live scoring interface', () => {
      cy.visit(`${baseUrl}/series/${seriesId}/matches/${matchId}/scorecard`);

      cy.wait('@getScorecard');

      cy.get('[data-cy="live-scoring-button"]').click();
      cy.get('[data-cy="live-scoring-interface"]').should('be.visible');

      cy.get('[data-cy="close-scoring-button"]').click();
      cy.get('[data-cy="live-scoring-interface"]').should('not.be.visible');
    });
  });

  describe('Over Management', () => {
    it('should expand and collapse overs view', () => {
      cy.visit(`${baseUrl}/series/${seriesId}/matches/${matchId}/scorecard`);

      cy.wait('@getScorecard');

      // Initially should show "Show All Overs" button
      cy.get('[data-cy="show-all-overs-button"]').should('be.visible');

      // Click to expand
      cy.get('[data-cy="show-all-overs-button"]').click();
      cy.get('[data-cy="hide-all-overs-button"]').should('be.visible');

      // Click to collapse
      cy.get('[data-cy="hide-all-overs-button"]').click();
      cy.get('[data-cy="show-all-overs-button"]').should('be.visible');
    });

    it('should display latest over details', () => {
      cy.visit(`${baseUrl}/series/${seriesId}/matches/${matchId}/scorecard`);

      cy.wait('@getScorecard');

      cy.get('[data-cy="latest-over"]').should('contain', 'Latest Over 1');
      cy.get('[data-cy="over-summary"]').should('contain', '8 runs, 0 wickets');
    });
  });

  describe('Navigation and Actions', () => {
    it('should handle back button', () => {
      cy.visit(`${baseUrl}/series/${seriesId}/matches/${matchId}/scorecard`);

      cy.wait('@getScorecard');

      cy.get('[data-cy="back-button"]').click();

      // Should navigate back to matches page
      cy.url().should('include', `/series/${seriesId}/matches`);
    });

    it('should handle refresh button', () => {
      cy.visit(`${baseUrl}/series/${seriesId}/matches/${matchId}/scorecard`);

      cy.wait('@getScorecard');

      cy.get('[data-cy="refresh-button"]').click();

      // Should trigger another API call
      cy.wait('@getScorecard');
    });
  });

  describe('Error Handling', () => {
    it('should handle scorecard not found error', () => {
      cy.intercept('GET', `/api/v1/scorecard/${matchId}`, {
        statusCode: 404,
        body: { error: 'Scorecard not found' },
      }).as('getScorecardError');

      cy.visit(`${baseUrl}/series/${seriesId}/matches/${matchId}/scorecard`);

      cy.wait('@getScorecardError');

      cy.get('[data-cy="error-message"]').should('contain', 'Error:');
      cy.get('[data-cy="retry-button"]').should('be.visible');
    });

    it('should handle start scoring error', () => {
      cy.intercept('POST', '/api/v1/scorecard/start', {
        statusCode: 400,
        body: { error: 'Match is not ready for scoring' },
      }).as('startScoringError');

      cy.visit(`${baseUrl}/series/${seriesId}/matches/${matchId}/scorecard`);

      cy.wait('@getScorecard');

      cy.get('[data-cy="live-scoring-button"]').click();
      cy.wait('@startScoringError');

      cy.get('[data-cy="error-message"]').should(
        'contain',
        'Match is not ready for scoring'
      );
    });

    it('should handle add ball error', () => {
      cy.intercept('POST', '/api/v1/scorecard/ball', {
        statusCode: 400,
        body: { error: 'Invalid ball data' },
      }).as('addBallError');

      cy.visit(`${baseUrl}/series/${seriesId}/matches/${matchId}/scorecard`);

      cy.wait('@getScorecard');

      cy.get('[data-cy="live-scoring-button"]').click();
      cy.get('[data-cy="run-button-4"]').click();
      cy.wait('@addBallError');

      cy.get('[data-cy="error-message"]').should(
        'contain',
        'Invalid ball data'
      );
    });

    it('should handle innings completion error', () => {
      cy.intercept('POST', '/api/v1/scorecard/ball', {
        statusCode: 400,
        body: { error: 'Innings already completed' },
      }).as('inningsCompletedError');

      cy.visit(`${baseUrl}/series/${seriesId}/matches/${matchId}/scorecard`);

      cy.wait('@getScorecard');

      cy.get('[data-cy="live-scoring-button"]').click();
      cy.get('[data-cy="run-button-4"]').click();
      cy.wait('@inningsCompletedError');

      cy.get('[data-cy="error-message"]').should(
        'contain',
        'Innings already completed'
      );
    });
  });

  describe('Loading States', () => {
    it('should show loading state during scorecard fetch', () => {
      cy.intercept('GET', `/api/v1/scorecard/${matchId}`, {
        delay: 1000,
        statusCode: 200,
        body: { data: {} },
      }).as('getScorecardDelayed');

      cy.visit(`${baseUrl}/series/${seriesId}/matches/${matchId}/scorecard`);

      cy.get('[data-cy="loading-spinner"]').should('be.visible');
      cy.get('[data-cy="loading-text"]').should(
        'contain',
        'Loading scorecard...'
      );

      cy.wait('@getScorecardDelayed');
    });

    it('should show loading state during scoring', () => {
      cy.intercept('POST', '/api/v1/scorecard/ball', {
        delay: 1000,
        statusCode: 200,
        body: { message: 'Ball added successfully' },
      }).as('addBallDelayed');

      cy.visit(`${baseUrl}/series/${seriesId}/matches/${matchId}/scorecard`);

      cy.wait('@getScorecard');

      cy.get('[data-cy="live-scoring-button"]').click();
      cy.get('[data-cy="run-button-4"]').click();

      // Should show loading state on buttons
      cy.get('[data-cy="run-button-4"]').should('be.disabled');
      cy.get('[data-cy="scoring-loading-spinner"]').should('be.visible');

      cy.wait('@addBallDelayed');
    });
  });

  describe('Responsive Design', () => {
    it('should be responsive on mobile devices', () => {
      cy.viewport('iphone-x');

      cy.visit(`${baseUrl}/series/${seriesId}/matches/${matchId}/scorecard`);

      cy.wait('@getScorecard');

      // Check that the layout adapts to mobile
      cy.get('[data-cy="scorecard-container"]').should('be.visible');
      cy.get('[data-cy="team-a-card"]').should('be.visible');
      cy.get('[data-cy="team-b-card"]').should('be.visible');

      // Live scoring interface should still work on mobile
      cy.get('[data-cy="live-scoring-button"]').click();
      cy.get('[data-cy="live-scoring-interface"]').should('be.visible');
    });

    it('should be responsive on tablet devices', () => {
      cy.viewport('ipad-2');

      cy.visit(`${baseUrl}/series/${seriesId}/matches/${matchId}/scorecard`);

      cy.wait('@getScorecard');

      // Check that the layout works on tablet
      cy.get('[data-cy="scorecard-container"]').should('be.visible');
      cy.get('[data-cy="team-a-card"]').should('be.visible');
      cy.get('[data-cy="team-b-card"]').should('be.visible');
    });
  });

  describe('Match Status Handling', () => {
    it('should handle scheduled match status', () => {
      cy.intercept('GET', `/api/v1/scorecard/${matchId}`, {
        statusCode: 200,
        body: {
          data: {
            match_id: matchId,
            match_number: 1,
            series_name: 'Test Series',
            team_a: 'Team A',
            team_b: 'Team B',
            total_overs: 20,
            toss_winner: 'A',
            toss_type: 'H',
            current_innings: 1,
            match_status: 'scheduled',
            innings: null,
          },
        },
      }).as('getScheduledMatch');

      cy.visit(`${baseUrl}/series/${seriesId}/matches/${matchId}/scorecard`);

      cy.wait('@getScheduledMatch');

      cy.get('[data-cy="match-status"]').should('contain', 'SCHEDULED');
      cy.get('[data-cy="match-ready-message"]').should(
        'contain',
        'Match ready to start'
      );
    });

    it('should handle completed match status', () => {
      cy.intercept('GET', `/api/v1/scorecard/${matchId}`, {
        statusCode: 200,
        body: {
          data: {
            match_id: matchId,
            match_number: 1,
            series_name: 'Test Series',
            team_a: 'Team A',
            team_b: 'Team B',
            total_overs: 20,
            toss_winner: 'A',
            toss_type: 'H',
            current_innings: 2,
            match_status: 'completed',
            innings: [
              {
                innings_number: 1,
                batting_team: 'A',
                total_runs: 150,
                total_wickets: 5,
                total_overs: 20,
                total_balls: 120,
                status: 'completed',
                extras: {
                  byes: 5,
                  leg_byes: 2,
                  wides: 8,
                  no_balls: 1,
                  total: 16,
                },
                overs: [],
              },
              {
                innings_number: 2,
                batting_team: 'B',
                total_runs: 145,
                total_wickets: 8,
                total_overs: 19,
                total_balls: 114,
                status: 'completed',
                extras: {
                  byes: 3,
                  leg_byes: 1,
                  wides: 6,
                  no_balls: 2,
                  total: 12,
                },
                overs: [],
              },
            ],
          },
        },
      }).as('getCompletedMatch');

      cy.visit(`${baseUrl}/series/${seriesId}/matches/${matchId}/scorecard`);

      cy.wait('@getCompletedMatch');

      cy.get('[data-cy="match-status"]').should('contain', 'COMPLETED');
      cy.get('[data-cy="innings-1"]').should('contain', 'Completed');
      cy.get('[data-cy="innings-2"]').should('contain', 'Completed');
    });
  });
});
