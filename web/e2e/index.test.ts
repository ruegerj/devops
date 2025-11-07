import { expect, test } from './base-fixtures';
import { resolveLocatorFor } from './util';

test.beforeEach(async ({ page }) => {
	await page.goto('/');
});

test('home page has expected title', async ({ page }) => {
	const title = resolveLocatorFor(page, 'card-title');
	await expect(title).toBeVisible();
	await expect(title).toHaveText('Secret Vault');
});

test('vault can be unlocked', async ({ page }) => {
	const expectedContent = JSON.stringify(
		{ message: 'Life, Universe and everything', number: 42 },
		null,
		3
	);

	await page.goto('/');

	const unlockBtn = resolveLocatorFor(page, 'unlock-btn');
	await expect(unlockBtn).toBeVisible();
	await unlockBtn.click();

	const codeContainer = resolveLocatorFor(page, 'code-container');
	await expect(codeContainer).toBeVisible();
	await expect(codeContainer).toHaveText(expectedContent);
});

test('vault can be unlocked and locked again', async ({ page }) => {
	const unlockBtn = resolveLocatorFor(page, 'unlock-btn');
	await expect(unlockBtn).toBeVisible();
	await unlockBtn.click();

	const lockBtn = resolveLocatorFor(page, 'lock-btn');
	await expect(lockBtn).toBeVisible();
	await lockBtn.click();

	const lockedInfoText = resolveLocatorFor(page, 'locked-info');
	await expect(lockedInfoText).toBeVisible();
	await expect(lockedInfoText).toHaveText('Open it up to see its contents...');
});
