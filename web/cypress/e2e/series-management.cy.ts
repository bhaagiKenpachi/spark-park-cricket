describe('Series Management E2E Tests', () => {
    beforeEach(() => {
        // Mock API responses
        cy.intercept('GET', '/api/v1/series', {
            statusCode: 200,
            body: [
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
            ],
        }).as('fetchSeries');

        cy.intercept('POST', '/api/v1/series', {
            statusCode: 201,
            body: {
                id: '2',
                name: 'New Series',
                description: 'New Description',
                start_date: '2024-02-01',
                end_date: '2024-02-28',
                status: 'upcoming',
                created_at: '2024-02-01T00:00:00Z',
                updated_at: '2024-02-01T00:00:00Z',
            },
        }).as('createSeries');

        cy.intercept('PUT', '/api/v1/series/1', {
            statusCode: 200,
            body: {
                id: '1',
                name: 'Updated Series',
                description: 'Updated Description',
                start_date: '2024-01-01',
                end_date: '2024-01-31',
                status: 'upcoming',
                created_at: '2024-01-01T00:00:00Z',
                updated_at: '2024-01-01T00:00:00Z',
            },
        }).as('updateSeries');

        cy.intercept('DELETE', '/api/v1/series/1', {
            statusCode: 200,
            body: { success: true },
        }).as('deleteSeries');
    });

    it('should display series list on page load', () => {
        cy.visit('/');
        cy.wait('@fetchSeries');

        cy.get('[data-cy=series-list]').should('be.visible');
        cy.get('[data-cy=series-item]').should('have.length.at.least', 1);
        cy.get('[data-cy=series-name]').should('contain', 'Test Series');
    });

    it('should create a new series', () => {
        cy.visit('/');
        cy.wait('@fetchSeries');

        // Click create series button
        cy.get('[data-cy=create-series-button]').click();

        // Fill out the form
        cy.get('[data-cy=series-name]').type('New Series');
        cy.get('[data-cy=series-description]').type('New Description');
        cy.get('[data-cy=start-date]').type('2024-02-01');
        cy.get('[data-cy=end-date]').type('2024-02-28');

        // Submit the form
        cy.get('[data-cy=create-series-button]').click();
        cy.wait('@createSeries');

        // Verify success
        cy.get('[data-cy=success-message]').should('contain', 'Series created successfully');
    });

    it('should edit an existing series', () => {
        cy.visit('/');
        cy.wait('@fetchSeries');

        // Click edit button on first series
        cy.get('[data-cy=series-item]').first().find('[data-cy=edit-button]').click();

        // Update the form
        cy.get('[data-cy=series-name]').clear().type('Updated Series');
        cy.get('[data-cy=series-description]').clear().type('Updated Description');

        // Submit the form
        cy.get('[data-cy=update-series-button]').click();
        cy.wait('@updateSeries');

        // Verify success
        cy.get('[data-cy=success-message]').should('contain', 'Series updated successfully');
    });

    it('should delete a series', () => {
        cy.visit('/');
        cy.wait('@fetchSeries');

        // Click delete button on first series
        cy.get('[data-cy=series-item]').first().find('[data-cy=delete-button]').click();

        // Confirm deletion
        cy.on('window:confirm', () => true);

        cy.wait('@deleteSeries');

        // Verify success
        cy.get('[data-cy=success-message]').should('contain', 'Series deleted successfully');
    });

    it('should validate required fields', () => {
        cy.visit('/');
        cy.wait('@fetchSeries');

        // Click create series button
        cy.get('[data-cy=create-series-button]').click();

        // Try to submit without filling required fields
        cy.get('[data-cy=create-series-button]').click();

        // Check for validation errors
        cy.get('[data-cy=error-message]').should('contain', 'Name is required');
        cy.get('[data-cy=error-message]').should('contain', 'Start date is required');
        cy.get('[data-cy=error-message]').should('contain', 'End date is required');
    });

    it('should validate date constraints', () => {
        cy.visit('/');
        cy.wait('@fetchSeries');

        // Click create series button
        cy.get('[data-cy=create-series-button]').click();

        // Fill out form with invalid dates
        cy.get('[data-cy=series-name]').type('Test Series');
        cy.get('[data-cy=start-date]').type('2024-01-31');
        cy.get('[data-cy=end-date]').type('2024-01-01');

        // Submit the form
        cy.get('[data-cy=create-series-button]').click();

        // Check for date validation error
        cy.get('[data-cy=error-message]').should('contain', 'End date must be after start date');
    });

    it('should handle API errors gracefully', () => {
        // Mock API error
        cy.intercept('POST', '/api/v1/series', {
            statusCode: 400,
            body: { error: 'Invalid data' },
        }).as('createSeriesError');

        cy.visit('/');
        cy.wait('@fetchSeries');

        // Click create series button
        cy.get('[data-cy=create-series-button]').click();

        // Fill out the form
        cy.get('[data-cy=series-name]').type('Test Series');
        cy.get('[data-cy=series-description]').type('Test Description');
        cy.get('[data-cy=start-date]').type('2024-01-01');
        cy.get('[data-cy=end-date]').type('2024-01-31');

        // Submit the form
        cy.get('[data-cy=create-series-button]').click();
        cy.wait('@createSeriesError');

        // Check for error message
        cy.get('[data-cy=error-message]').should('contain', 'Invalid data');
    });

    it('should show loading states', () => {
        // Mock slow API response
        cy.intercept('GET', '/api/v1/series', {
            delay: 1000,
            statusCode: 200,
            body: [],
        }).as('slowFetchSeries');

        cy.visit('/');

        // Check for loading state
        cy.get('[data-cy=loading-spinner]').should('be.visible');
        cy.get('[data-cy=loading-text]').should('contain', 'Loading series...');

        cy.wait('@slowFetchSeries');

        // Loading state should be gone
        cy.get('[data-cy=loading-spinner]').should('not.exist');
    });

    it('should be responsive on mobile devices', () => {
        // Set mobile viewport
        cy.viewport(375, 667);

        cy.visit('/');
        cy.wait('@fetchSeries');

        // Check mobile layout
        cy.get('[data-cy=series-list]').should('be.visible');
        cy.get('[data-cy=create-series-button]').should('be.visible');

        // Check that buttons are touch-friendly
        cy.get('[data-cy=create-series-button]').should('have.css', 'min-height', '44px');
    });

    it('should navigate between series pages', () => {
        cy.visit('/');
        cy.wait('@fetchSeries');

        // Click on a series item
        cy.get('[data-cy=series-item]').first().click();

        // Should navigate to series detail page
        cy.url().should('include', '/series/');

        // Check for back button
        cy.get('[data-cy=back-button]').should('be.visible');

        // Click back button
        cy.get('[data-cy=back-button]').click();

        // Should return to series list
        cy.url().should('include', '/');
    });

    it('should filter series by status', () => {
        // Mock series with different statuses
        cy.intercept('GET', '/api/v1/series', {
            statusCode: 200,
            body: [
                {
                    id: '1',
                    name: 'Upcoming Series',
                    description: 'Test Description',
                    start_date: '2024-01-01',
                    end_date: '2024-01-31',
                    status: 'upcoming',
                    created_at: '2024-01-01T00:00:00Z',
                    updated_at: '2024-01-01T00:00:00Z',
                },
                {
                    id: '2',
                    name: 'Ongoing Series',
                    description: 'Test Description',
                    start_date: '2024-02-01',
                    end_date: '2024-02-28',
                    status: 'ongoing',
                    created_at: '2024-02-01T00:00:00Z',
                    updated_at: '2024-02-01T00:00:00Z',
                },
            ],
        }).as('fetchSeriesWithStatus');

        cy.visit('/');
        cy.wait('@fetchSeriesWithStatus');

        // Check that both series are displayed
        cy.get('[data-cy=series-item]').should('have.length', 2);
        cy.get('[data-cy=series-name]').should('contain', 'Upcoming Series');
        cy.get('[data-cy=series-name]').should('contain', 'Ongoing Series');

        // Filter by upcoming status
        cy.get('[data-cy=status-filter]').select('upcoming');

        // Should show only upcoming series
        cy.get('[data-cy=series-item]').should('have.length', 1);
        cy.get('[data-cy=series-name]').should('contain', 'Upcoming Series');
        cy.get('[data-cy=series-name]').should('not.contain', 'Ongoing Series');
    });
});