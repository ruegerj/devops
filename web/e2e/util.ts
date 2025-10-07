import { Page } from '@playwright/test';

export function resolveLocatorFor(page: Page, testElemId: string) {
	return page.locator(`[data-test="${testElemId}"]`);
}
