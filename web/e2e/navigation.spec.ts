import { test, expect } from '@playwright/test';

test.describe('Navigation', () => {
  test('should have navigation links', async ({ page }) => {
    await page.goto('/');

    await expect(page.getByRole('link', { name: 'Home' })).toBeVisible();
    await expect(page.getByRole('link', { name: 'Users' })).toBeVisible();
    await expect(page.getByRole('link', { name: 'About' })).toBeVisible();
  });

  test('should navigate to Users page', async ({ page }) => {
    await page.goto('/');

    await page.getByRole('link', { name: 'Users' }).click();

    await expect(page).toHaveURL('/users');
  });

  test('should navigate to About page', async ({ page }) => {
    await page.goto('/');

    await page.getByRole('link', { name: 'About' }).click();

    await expect(page).toHaveURL('/about');
  });

  test('should navigate back to Home page', async ({ page }) => {
    await page.goto('/about');

    await page.getByRole('link', { name: 'Home' }).click();

    await expect(page).toHaveURL('/');
  });
});
