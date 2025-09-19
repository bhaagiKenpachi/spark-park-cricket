describe('Authentication E2E Tests', () => {
  beforeEach(() => {
    // Clear localStorage and cookies before each test
    cy.clearLocalStorage();
    cy.clearCookies();

    // Visit the home page
    cy.visit('/');
  });

  describe('Authentication Flow', () => {
    it('should show login button when not authenticated', () => {
      // Should show login button
      cy.contains('Sign in with Google').should('be.visible');

      // Should not show user menu
      cy.get('[data-testid="user-menu"]').should('not.exist');
    });

    it('should handle Google OAuth redirect', () => {
      // Mock the Google OAuth redirect
      cy.window().then(win => {
        cy.stub(win, 'location', {
          href: '',
          assign: cy.stub(),
          replace: cy.stub(),
          reload: cy.stub(),
        });
      });

      // Click login button
      cy.contains('Sign in with Google').click();

      // Should redirect to Google OAuth
      cy.window().its('location.href').should('contain', '/auth/google');
    });

    it('should show loading state during authentication', () => {
      // Mock authentication in progress
      cy.window().then(win => {
        cy.stub(win, 'location', {
          href: '',
          assign: cy.stub(),
          replace: cy.stub(),
          reload: cy.stub(),
        });
      });

      // Click login button
      cy.contains('Sign in with Google').click();

      // Should show loading state
      cy.contains('Signing in...').should('be.visible');
    });
  });

  describe('Protected Routes', () => {
    it('should redirect to login when accessing protected route without authentication', () => {
      // Try to access a protected route
      cy.visit('/series/create');

      // Should show authentication required message
      cy.contains('Authentication Required').should('be.visible');
      cy.contains('Please sign in to access this feature.').should(
        'be.visible'
      );
      cy.contains('Sign in with Google').should('be.visible');
    });

    it('should show custom fallback for protected routes', () => {
      // This would require a specific page with custom fallback
      // For now, we'll test the default behavior
      cy.visit('/series/create');

      // Should show default authentication prompt
      cy.contains('Authentication Required').should('be.visible');
    });
  });

  describe('Authenticated User Experience', () => {
    beforeEach(() => {
      // Mock authenticated user state
      cy.window().then(win => {
        win.localStorage.setItem('auth_authenticated', 'true');
        win.localStorage.setItem(
          'auth_user',
          JSON.stringify({
            id: 'user-123',
            google_id: 'google-123',
            email: 'test@example.com',
            name: 'Test User',
            picture: 'https://example.com/picture.jpg',
            email_verified: true,
            created_at: '2024-01-01T00:00:00Z',
            updated_at: '2024-01-01T00:00:00Z',
          })
        );
      });
    });

    it('should show user menu when authenticated', () => {
      cy.visit('/');

      // Should show user menu
      cy.get('[data-testid="user-menu"]').should('be.visible');
      cy.contains('Test User').should('be.visible');
    });

    it('should show user dropdown when clicked', () => {
      cy.visit('/');

      // Click on user menu
      cy.contains('Test User').click();

      // Should show dropdown menu
      cy.contains('Profile').should('be.visible');
      cy.contains('Sign out').should('be.visible');
    });

    it('should allow access to protected routes when authenticated', () => {
      cy.visit('/series/create');

      // Should not show authentication required message
      cy.contains('Authentication Required').should('not.exist');

      // Should show the protected content (series creation form)
      cy.get('form').should('be.visible');
    });

    it('should handle logout', () => {
      // Mock logout API call
      cy.intercept('POST', '/api/v1/auth/logout', {
        statusCode: 200,
        body: {
          data: { message: 'Logged out successfully' },
          success: true,
          message: 'Success',
        },
      }).as('logout');

      cy.visit('/');

      // Click on user menu
      cy.contains('Test User').click();

      // Click sign out
      cy.contains('Sign out').click();

      // Wait for logout API call
      cy.wait('@logout');

      // Should show login button again
      cy.contains('Sign in with Google').should('be.visible');

      // Should not show user menu
      cy.get('[data-testid="user-menu"]').should('not.exist');
    });
  });

  describe('Session Management', () => {
    it('should handle expired session', () => {
      // Mock expired session
      cy.window().then(win => {
        win.localStorage.setItem('auth_authenticated', 'true');
        win.localStorage.setItem(
          'auth_user',
          JSON.stringify({
            id: 'user-123',
            google_id: 'google-123',
            email: 'test@example.com',
            name: 'Test User',
            picture: 'https://example.com/picture.jpg',
            email_verified: true,
            created_at: '2024-01-01T00:00:00Z',
            updated_at: '2024-01-01T00:00:00Z',
          })
        );
      });

      // Mock auth status check returning unauthenticated
      cy.intercept('GET', '/api/v1/auth/status', {
        statusCode: 200,
        body: {
          data: { authenticated: false },
          success: true,
          message: 'Success',
        },
      }).as('authStatus');

      cy.visit('/');

      // Wait for auth status check
      cy.wait('@authStatus');

      // Should show login button (local state cleared)
      cy.contains('Sign in with Google').should('be.visible');

      // Should not show user menu
      cy.get('[data-testid="user-menu"]').should('not.exist');
    });

    it('should handle network errors during auth status check', () => {
      // Mock network error
      cy.intercept('GET', '/api/v1/auth/status', {
        statusCode: 500,
        body: { error: 'Internal Server Error' },
      }).as('authStatusError');

      cy.visit('/');

      // Wait for auth status check to fail
      cy.wait('@authStatusError');

      // Should show login button (fallback to unauthenticated)
      cy.contains('Sign in with Google').should('be.visible');
    });
  });

  describe('Authentication State Persistence', () => {
    it('should persist authentication state across page reloads', () => {
      // Mock authenticated user
      cy.window().then(win => {
        win.localStorage.setItem('auth_authenticated', 'true');
        win.localStorage.setItem(
          'auth_user',
          JSON.stringify({
            id: 'user-123',
            google_id: 'google-123',
            email: 'test@example.com',
            name: 'Test User',
            picture: 'https://example.com/picture.jpg',
            email_verified: true,
            created_at: '2024-01-01T00:00:00Z',
            updated_at: '2024-01-01T00:00:00Z',
          })
        );
      });

      // Mock successful auth status check
      cy.intercept('GET', '/api/v1/auth/status', {
        statusCode: 200,
        body: {
          data: {
            authenticated: true,
            user: {
              id: 'user-123',
              google_id: 'google-123',
              email: 'test@example.com',
              name: 'Test User',
              picture: 'https://example.com/picture.jpg',
              email_verified: true,
              created_at: '2024-01-01T00:00:00Z',
              updated_at: '2024-01-01T00:00:00Z',
            },
          },
          success: true,
          message: 'Success',
        },
      }).as('authStatus');

      cy.visit('/');

      // Wait for auth status check
      cy.wait('@authStatus');

      // Should show user menu
      cy.get('[data-testid="user-menu"]').should('be.visible');
      cy.contains('Test User').should('be.visible');

      // Reload the page
      cy.reload();

      // Wait for auth status check again
      cy.wait('@authStatus');

      // Should still show user menu
      cy.get('[data-testid="user-menu"]').should('be.visible');
      cy.contains('Test User').should('be.visible');
    });
  });

  describe('Error Handling', () => {
    it('should handle authentication errors gracefully', () => {
      // Mock authentication error
      cy.intercept('GET', '/api/v1/auth/status', {
        statusCode: 401,
        body: {
          error: 'UNAUTHORIZED',
          message: 'Not authenticated',
        },
      }).as('authError');

      cy.visit('/');

      // Wait for auth status check to fail
      cy.wait('@authError');

      // Should show login button
      cy.contains('Sign in with Google').should('be.visible');

      // Should not show error message to user (handled gracefully)
      cy.contains('Error').should('not.exist');
    });

    it('should handle logout errors gracefully', () => {
      // Mock authenticated user
      cy.window().then(win => {
        win.localStorage.setItem('auth_authenticated', 'true');
        win.localStorage.setItem(
          'auth_user',
          JSON.stringify({
            id: 'user-123',
            google_id: 'google-123',
            email: 'test@example.com',
            name: 'Test User',
            picture: 'https://example.com/picture.jpg',
            email_verified: true,
            created_at: '2024-01-01T00:00:00Z',
            updated_at: '2024-01-01T00:00:00Z',
          })
        );
      });

      // Mock successful auth status check
      cy.intercept('GET', '/api/v1/auth/status', {
        statusCode: 200,
        body: {
          data: {
            authenticated: true,
            user: {
              id: 'user-123',
              google_id: 'google-123',
              email: 'test@example.com',
              name: 'Test User',
              picture: 'https://example.com/picture.jpg',
              email_verified: true,
              created_at: '2024-01-01T00:00:00Z',
              updated_at: '2024-01-01T00:00:00Z',
            },
          },
          success: true,
          message: 'Success',
        },
      }).as('authStatus');

      // Mock logout error
      cy.intercept('POST', '/api/v1/auth/logout', {
        statusCode: 500,
        body: { error: 'LOGOUT_ERROR', message: 'Failed to logout' },
      }).as('logoutError');

      cy.visit('/');

      // Wait for auth status check
      cy.wait('@authStatus');

      // Click on user menu
      cy.contains('Test User').click();

      // Click sign out
      cy.contains('Sign out').click();

      // Wait for logout API call to fail
      cy.wait('@logoutError');

      // Should still clear local state and show login button
      cy.contains('Sign in with Google').should('be.visible');
      cy.get('[data-testid="user-menu"]').should('not.exist');
    });
  });
});
