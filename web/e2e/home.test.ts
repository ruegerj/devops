import { expect, test } from '@playwright/test';

test('home page has expected title', async ({ page }) => {
	await page.goto('/');
	const titleLocator = page.locator('[data-slot="card-header"]');
	await expect(titleLocator).toBeVisible();
	await expect(titleLocator).toHaveText('Secret Vault');
});
