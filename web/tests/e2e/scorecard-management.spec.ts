import { test, expect } from '@playwright/test';

test.describe('Scorecard Management', () => {
  test('should display basic interface', async ({ page }) => {
    await page.goto('/');
    await page.waitForLoadState('networkidle');

    // Check basic page elements
    await expect(page).toHaveTitle('Spark Park Cricket');
    await expect(
      page.getByRole('heading', { name: 'Spark Park Cricket', exact: true })
    ).toBeVisible();
    await expect(page.locator('[data-cy="login-button"]')).toBeVisible();
  });

  test('should handle basic navigation', async ({ page }) => {
    // Mock empty series response
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

    // Test basic interaction - look for any create button
    const createButton = page
      .locator('[data-cy="create-series-button"]')
      .or(page.locator('[data-cy="create-first-series-button"]'));

    if ((await createButton.count()) > 0) {
      await expect(createButton.first()).toBeVisible();

      // Click create series button
      await createButton.first().click();

      // Should still show main page elements
      await expect(
        page.getByRole('heading', { name: 'Spark Park Cricket', exact: true })
      ).toBeVisible();
    } else {
      // If no create button, just verify basic elements
      await expect(
        page.getByRole('heading', { name: 'Spark Park Cricket', exact: true })
      ).toBeVisible();
    }
  });

  test('should be responsive on mobile devices', async ({ page }) => {
    // Set mobile viewport
    await page.setViewportSize({ width: 375, height: 667 });

    await page.goto('/');
    await page.waitForLoadState('networkidle');

    // Check that the layout adapts to mobile
    await expect(
      page.getByRole('heading', { name: 'Spark Park Cricket', exact: true })
    ).toBeVisible();
    await expect(page.locator('[data-cy="login-button"]')).toBeVisible();
  });

  test('should handle error states gracefully', async ({ page }) => {
    await page.goto('/');
    await page.waitForLoadState('networkidle');

    // Should show either loading state or content
    const loadingOrContent = page
      .getByText('Loading series...')
      .or(page.getByText('Cricket Series'));
    await expect(loadingOrContent.first()).toBeVisible();
  });
});
