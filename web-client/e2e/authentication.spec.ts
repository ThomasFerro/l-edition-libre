import { test as base, expect, type Page } from '@playwright/test';

const authenticatedTest = base.extend<{ page: Page }>({
  page: async ({page}, use) => {
    await page.goto("/");
    await page.locator('[data-test="Authenticate button"]').click()
    await page.waitForLoadState("networkidle")
    use(page)
  },
});

authenticatedTest('Authenticated user is welcomed', async ({ page }) => {
  await expect(page.locator('h1')).toHaveText('Bienvenue John Doe');
})
