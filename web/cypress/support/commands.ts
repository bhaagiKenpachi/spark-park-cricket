/// <reference types="cypress" />

// Custom command for login
Cypress.Commands.add('login', (email: string, password: string) => {
  cy.session([email, password], () => {
    cy.visit('/login');
    cy.get('[data-cy=email-input]').type(email);
    cy.get('[data-cy=password-input]').type(password);
    cy.get('[data-cy=login-button]').click();
    cy.url().should('not.include', '/login');
  });
});

// Custom command for creating series
Cypress.Commands.add('createSeries', (seriesData: any) => {
  cy.visit('/series/new');
  cy.get('[data-cy=series-name]').type(seriesData.name);
  cy.get('[data-cy=series-description]').type(seriesData.description);
  cy.get('[data-cy=start-date]').type(seriesData.start_date);
  cy.get('[data-cy=end-date]').type(seriesData.end_date);
  cy.get('[data-cy=create-series-button]').click();
});

// Custom command for creating team
Cypress.Commands.add('createTeam', (teamData: any) => {
  cy.visit('/teams/new');
  cy.get('[data-cy=team-name]').type(teamData.name);
  cy.get('[data-cy=team-description]').type(teamData.description);
  cy.get('[data-cy=create-team-button]').click();
});

// Custom command for creating match
Cypress.Commands.add('createMatch', (matchData: any) => {
  cy.visit('/matches/new');
  cy.get('[data-cy=series-select]').select(matchData.series_id);
  cy.get('[data-cy=team1-select]').select(matchData.team1_id);
  cy.get('[data-cy=team2-select]').select(matchData.team2_id);
  cy.get('[data-cy=venue-input]').type(matchData.venue);
  cy.get('[data-cy=match-date]').type(matchData.match_date);
  cy.get('[data-cy=create-match-button]').click();
});

// Custom command for mocking API responses
Cypress.Commands.add('mockApiResponse', (endpoint: string, response: any) => {
  cy.intercept('GET', endpoint, response).as('apiCall');
});
