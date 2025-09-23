import { test, expect } from '@playwright/test';

test.describe('Authentication Flow', () => {
  test.beforeEach(async ({ page }) => {
    // Navigate to the home page before each test
    await page.goto('/');
    await page.waitForLoadState('networkidle');
  });

  test('should show login button when not authenticated', async ({ page }) => {
    // Should show login button
    await expect(page.locator('[data-cy="login-button"]')).toBeVisible();
    await expect(page.getByText('Sign in with Google')).toBeVisible();

    // Should not show user menu
    await expect(page.locator('[data-testid="user-menu"]')).not.toBeVisible();
  });

  test('should handle Google OAuth redirect', async ({ page, browserName }) => {
    // Skip this test on WebKit as it doesn't support redirect status codes
    if (browserName === 'webkit') {
      test.skip(
        true,
        'WebKit does not support redirect status codes in route.fulfill()'
      );
      return;
    }

    // Mock the Google OAuth redirect by intercepting the request
    let redirected = false;
    await page.route('**/auth/google', async route => {
      redirected = true;
      await route.fulfill({
        status: 302,
        headers: { Location: '/auth/callback?code=test' },
      });
    });

    // Click login button
    await page.locator('[data-cy="login-button"]').click();

    // Should attempt to redirect to Google OAuth
    await page.waitForTimeout(1000); // Give time for redirect attempt
    expect(redirected).toBe(true);
  });

  test('should display page title and basic content', async ({ page }) => {
    // Check page title
    await expect(page).toHaveTitle('Spark Park Cricket');

    // Check main heading (use first() to avoid multiple element issue)
    await expect(
      page.getByRole('heading', { name: 'Spark Park Cricket', exact: true })
    ).toBeVisible();

    // Check welcome message
    await expect(page.getByText('Welcome to Spark Park Cricket')).toBeVisible();

    // Check description
    await expect(
      page.getByText('Manage your cricket tournaments')
    ).toBeVisible();
  });

  test('should show loading state initially', async ({ page }) => {
    // Should show loading series text initially
    await expect(page.getByText('Loading series...')).toBeVisible();

    // Should show login button
    await expect(page.locator('[data-cy="login-button"]')).toBeVisible();
  });
});

test.describe('Responsive Design', () => {
  test('should work on mobile devices', async ({ page }) => {
    // Set mobile viewport
    await page.setViewportSize({ width: 375, height: 667 });

    await page.goto('/');
    await page.waitForLoadState('networkidle');

    // Check that the page is responsive - use exact match for heading
    await expect(
      page.getByRole('heading', { name: 'Spark Park Cricket', exact: true })
    ).toBeVisible();

    // Check that buttons are still visible on mobile
    await expect(page.locator('[data-cy="login-button"]')).toBeVisible();
  });

  test('should work on tablet devices', async ({ page }) => {
    // Set tablet viewport
    await page.setViewportSize({ width: 768, height: 1024 });

    await page.goto('/');
    await page.waitForLoadState('networkidle');

    // Check that the page works on tablet - use exact match for heading
    await expect(
      page.getByRole('heading', { name: 'Spark Park Cricket', exact: true })
    ).toBeVisible();
    await expect(page.locator('[data-cy="login-button"]')).toBeVisible();
  });
});
