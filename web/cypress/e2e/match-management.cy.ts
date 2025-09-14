describe('Match Management E2E Tests', () => {
    beforeEach(() => {
        // Mock API responses for matches
        cy.intercept('GET', '/api/v1/matches', {
            statusCode: 200,
            body: [
                {
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
                },
            ],
        }).as('fetchMatches');

        cy.intercept('GET', '/api/v1/matches/series/series-1', {
            statusCode: 200,
            body: [
                {
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
                },
            ],
        }).as('fetchMatchesBySeries');

        cy.intercept('POST', '/api/v1/matches', {
            statusCode: 201,
            body: {
                id: '2',
                series_id: 'series-1',
                match_number: 2,
                date: '2024-02-01T00:00:00Z',
                status: 'live',
                team_a_player_count: 11,
                team_b_player_count: 11,
                total_overs: 20,
                toss_winner: 'B',
                toss_type: 'T',
                batting_team: 'B',
                created_at: '2024-02-01T00:00:00Z',
                updated_at: '2024-02-01T00:00:00Z',
            },
        }).as('createMatch');

        cy.intercept('PUT', '/api/v1/matches/1', {
            statusCode: 200,
            body: {
                id: '1',
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
                created_at: '2024-01-01T00:00:00Z',
                updated_at: '2024-01-01T00:00:00Z',
            },
        }).as('updateMatch');

        cy.intercept('DELETE', '/api/v1/matches/1', {
            statusCode: 200,
            body: { success: true },
        }).as('deleteMatch');
    });

    it('should display match form on page load', () => {
        cy.visit('/series/series-1/matches');
        cy.wait('@fetchMatchesBySeries');

        cy.get('[data-cy=match-form]').should('be.visible');
        cy.get('[data-cy=match-date]').should('be.visible');
        cy.get('[data-cy=team-player-count]').should('be.visible');
        cy.get('[data-cy=total-overs]').should('be.visible');
        cy.get('[data-cy=match-number]').should('be.visible');
    });

    it('should create a new match', () => {
        cy.visit('/series/series-1/matches');
        cy.wait('@fetchMatchesBySeries');

        // Fill out the form
        cy.get('[data-cy=match-date]').type('2024-02-01');
        cy.get('[data-cy=team-player-count]').type('11');
        cy.get('[data-cy=total-overs]').type('20');
        cy.get('[data-cy=match-number]').type('2');

        // Select toss winner
        cy.get('[data-cy=toss-winner]').click();
        cy.get('[data-cy=toss-winner-option-B]').click();

        // Select toss type
        cy.get('[data-cy=toss-type]').click();
        cy.get('[data-cy=toss-type-option-T]').click();

        // Submit the form
        cy.get('[data-cy=create-match-button]').click();
        cy.wait('@createMatch');

        // Verify success
        cy.get('[data-cy=success-message]').should('contain', 'Match created successfully');
    });

    it('should edit an existing match', () => {
        cy.visit('/series/series-1/matches');
        cy.wait('@fetchMatchesBySeries');

        // Click edit button on first match
        cy.get('[data-cy=match-item]').first().find('[data-cy=edit-match-button]').click();

        // Update the form
        cy.get('[data-cy=team-player-count]').clear().type('10');
        cy.get('[data-cy=total-overs]').clear().type('15');

        // Submit the form
        cy.get('[data-cy=update-match-button]').click();
        cy.wait('@updateMatch');

        // Verify success
        cy.get('[data-cy=success-message]').should('contain', 'Match updated successfully');
    });

    it('should delete a match', () => {
        cy.visit('/series/series-1/matches');
        cy.wait('@fetchMatchesBySeries');

        // Click delete button on first match
        cy.get('[data-cy=match-item]').first().find('[data-cy=delete-match-button]').click();

        // Confirm deletion
        cy.on('window:confirm', () => true);

        cy.wait('@deleteMatch');

        // Verify success
        cy.get('[data-cy=success-message]').should('contain', 'Match deleted successfully');
    });

    it('should validate required fields', () => {
        cy.visit('/series/series-1/matches');
        cy.wait('@fetchMatchesBySeries');

        // Try to submit without filling required fields
        cy.get('[data-cy=create-match-button]').click();

        // Check for validation errors
        cy.get('[data-cy=error-message]').should('contain', 'Date is required');
        cy.get('[data-cy=error-message]').should('contain', 'Team player count is required');
        cy.get('[data-cy=error-message]').should('contain', 'Total overs is required');
    });

    it('should validate data ranges', () => {
        cy.visit('/series/series-1/matches');
        cy.wait('@fetchMatchesBySeries');

        // Fill out form with invalid data
        cy.get('[data-cy=match-date]').type('2024-01-01');
        cy.get('[data-cy=team-player-count]').type('15'); // Invalid: > 11
        cy.get('[data-cy=total-overs]').type('25'); // Invalid: > 20

        // Submit the form
        cy.get('[data-cy=create-match-button]').click();

        // Check for validation errors
        cy.get('[data-cy=error-message]').should('contain', 'Team player count must be between 1 and 11');
        cy.get('[data-cy=error-message]').should('contain', 'Total overs must be between 1 and 20');
    });

    it('should handle API errors gracefully', () => {
        // Mock API error
        cy.intercept('POST', '/api/v1/matches', {
            statusCode: 400,
            body: { error: 'Invalid match data' },
        }).as('createMatchError');

        cy.visit('/series/series-1/matches');
        cy.wait('@fetchMatchesBySeries');

        // Fill out the form
        cy.get('[data-cy=match-date]').type('2024-01-01');
        cy.get('[data-cy=team-player-count]').type('11');
        cy.get('[data-cy=total-overs]').type('20');

        // Submit the form
        cy.get('[data-cy=create-match-button]').click();
        cy.wait('@createMatchError');

        // Check for error message
        cy.get('[data-cy=error-message]').should('contain', 'Invalid match data');
    });

    it('should show loading states', () => {
        // Mock slow API response
        cy.intercept('GET', '/api/v1/matches/series/series-1', {
            delay: 1000,
            statusCode: 200,
            body: [],
        }).as('slowFetchMatches');

        cy.visit('/series/series-1/matches');

        // Check for loading state
        cy.get('[data-cy=loading-spinner]').should('be.visible');
        cy.get('[data-cy=loading-text]').should('contain', 'Loading matches...');

        cy.wait('@slowFetchMatches');

        // Loading state should be gone
        cy.get('[data-cy=loading-spinner]').should('not.exist');
    });

    it('should be responsive on mobile devices', () => {
        // Set mobile viewport
        cy.viewport(375, 667);

        cy.visit('/series/series-1/matches');
        cy.wait('@fetchMatchesBySeries');

        // Check mobile layout
        cy.get('[data-cy=match-form]').should('be.visible');
        cy.get('[data-cy=create-match-button]').should('be.visible');

        // Check that buttons are touch-friendly
        cy.get('[data-cy=create-match-button]').should('have.css', 'min-height', '44px');
    });

    it('should navigate between match pages', () => {
        cy.visit('/series/series-1/matches');
        cy.wait('@fetchMatchesBySeries');

        // Click on a match item
        cy.get('[data-cy=match-item]').first().click();

        // Should navigate to match detail page
        cy.url().should('include', '/matches/');

        // Check for back button
        cy.get('[data-cy=back-button]').should('be.visible');

        // Click back button
        cy.get('[data-cy=back-button]').click();

        // Should return to matches list
        cy.url().should('include', '/matches');
    });

    it('should filter matches by status', () => {
        // Mock matches with different statuses
        cy.intercept('GET', '/api/v1/matches/series/series-1', {
            statusCode: 200,
            body: [
                {
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
                },
                {
                    id: '2',
                    series_id: 'series-1',
                    match_number: 2,
                    date: '2024-02-01T00:00:00Z',
                    status: 'completed',
                    team_a_player_count: 11,
                    team_b_player_count: 11,
                    total_overs: 20,
                    toss_winner: 'B',
                    toss_type: 'T',
                    batting_team: 'B',
                    created_at: '2024-02-01T00:00:00Z',
                    updated_at: '2024-02-01T00:00:00Z',
                },
            ],
        }).as('fetchMatchesWithStatus');

        cy.visit('/series/series-1/matches');
        cy.wait('@fetchMatchesWithStatus');

        // Check that both matches are displayed
        cy.get('[data-cy=match-item]').should('have.length', 2);

        // Filter by live status
        cy.get('[data-cy=status-filter]').select('live');

        // Should show only live matches
        cy.get('[data-cy=match-item]').should('have.length', 1);
        cy.get('[data-cy=match-status]').should('contain', 'live');
    });

    it('should handle match number auto-generation', () => {
        cy.visit('/series/series-1/matches');
        cy.wait('@fetchMatchesBySeries');

        // Fill out the form without match number
        cy.get('[data-cy=match-date]').type('2024-02-01');
        cy.get('[data-cy=team-player-count]').type('11');
        cy.get('[data-cy=total-overs]').type('20');
        // Don't fill match number - should auto-generate

        // Submit the form
        cy.get('[data-cy=create-match-button]').click();
        cy.wait('@createMatch');

        // Verify success
        cy.get('[data-cy=success-message]').should('contain', 'Match created successfully');
    });

    it('should display match details correctly', () => {
        cy.visit('/series/series-1/matches');
        cy.wait('@fetchMatchesBySeries');

        // Check match details are displayed
        cy.get('[data-cy=match-item]').first().within(() => {
            cy.get('[data-cy=match-number]').should('contain', '1');
            cy.get('[data-cy=match-date]').should('contain', '2024-01-01');
            cy.get('[data-cy=match-status]').should('contain', 'live');
            cy.get('[data-cy=team-player-count]').should('contain', '11');
            cy.get('[data-cy=total-overs]').should('contain', '20');
            cy.get('[data-cy=toss-winner]').should('contain', 'Team A');
            cy.get('[data-cy=toss-type]').should('contain', 'Heads');
        });
    });

    it('should handle form cancellation', () => {
        cy.visit('/series/series-1/matches');
        cy.wait('@fetchMatchesBySeries');

        // Fill out some form data
        cy.get('[data-cy=match-date]').type('2024-02-01');
        cy.get('[data-cy=team-player-count]').type('11');

        // Click cancel button
        cy.get('[data-cy=cancel-button]').click();

        // Form should be reset or closed
        cy.get('[data-cy=match-date]').should('have.value', '');
        cy.get('[data-cy=team-player-count]').should('have.value', '');
    });
});
