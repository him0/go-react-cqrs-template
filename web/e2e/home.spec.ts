import { test, expect } from '@playwright/test'

test.describe('Home Page', () => {
  test('should display the welcome heading', async ({ page }) => {
    await page.goto('/')

    await expect(
      page.getByRole('heading', { name: 'Welcome to User Management App' })
    ).toBeVisible()
  })

  test('should display the tech stack section', async ({ page }) => {
    await page.goto('/')

    await expect(page.getByRole('heading', { name: 'Tech Stack' })).toBeVisible()
    await expect(page.getByText('Go with DDD architecture')).toBeVisible()
    await expect(page.getByText('React 18 + TypeScript')).toBeVisible()
  })

  test('should display the getting started section', async ({ page }) => {
    await page.goto('/')

    await expect(page.getByRole('heading', { name: 'Getting Started' })).toBeVisible()
  })
})
