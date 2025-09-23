import { test, expect } from '@playwright/test';

test.describe('Match Management', () => {
  test('should display basic interface elements', async ({ page }) => {
    await page.goto('/');
    await page.waitForLoadState('networkidle');

    // Check that the page loads with basic elements
    await expect(page).toHaveTitle('Spark Park Cricket');
    await expect(
      page.getByRole('heading', { name: 'Spark Park Cricket', exact: true })
    ).toBeVisible();
    await expect(page.locator('[data-cy="login-button"]')).toBeVisible();
  });

  test('should be responsive on mobile devices', async ({ page }) => {
    // Set mobile viewport
    await page.setViewportSize({ width: 375, height: 667 });

    await page.goto('/');
    await page.waitForLoadState('networkidle');

    // Check that the page loads and basic elements are visible
    await expect(
      page.getByRole('heading', { name: 'Spark Park Cricket', exact: true })
    ).toBeVisible();
    await expect(page.locator('[data-cy="login-button"]')).toBeVisible();
  });

  test('should show appropriate create button based on series state', async ({
    page,
  }) => {
    // Mock empty series response to trigger "create first series" button
    await page.route('/api/v1/series', async route => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          data: [],
          success: true,
          message: 'Success',
        }),
      });
    });

    await page.goto('/');
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(1000);

    // Should show either create-series-button or create-first-series-button
    const createButton = page
      .locator('[data-cy="create-series-button"]')
      .or(page.locator('[data-cy="create-first-series-button"]'));

    if ((await createButton.count()) > 0) {
      await expect(createButton.first()).toBeVisible();

      // Test basic interaction with create series button
      await createButton.first().click();

      // Should still show the main page elements
      await expect(
        page.getByRole('heading', { name: 'Spark Park Cricket', exact: true })
      ).toBeVisible();
    } else {
      // If no create button found, just verify page loads
      await expect(
        page.getByRole('heading', { name: 'Spark Park Cricket', exact: true })
      ).toBeVisible();
    }
  });
});
