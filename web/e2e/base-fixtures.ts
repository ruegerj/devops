import { join } from 'node:path';
import { cwd } from 'node:process';
import { randomBytes } from 'node:crypto';
import { writeFileSync } from 'node:fs';
import { mkdir } from 'node:fs/promises';
import { test as baseTest } from '@playwright/test';

const istanbulCLIOutput = join(cwd(), 'coverage', 'e2e');

export function generateUUID(): string {
	return randomBytes(16).toString('hex');
}

// monkey patch test method in order to capture coverage
export const test = baseTest.extend({
	context: async ({ context }, use) => {
		await context.addInitScript(() => {
			window.addEventListener('beforeunload', () => {
				// eslint-disable-next-line @typescript-eslint/no-explicit-any
				(window as any).collectIstanbulCoverage(JSON.stringify((window as any).__coverage__));
			});
		});
		await mkdir(istanbulCLIOutput, { recursive: true });
		await context.exposeFunction('collectIstanbulCoverage', (coverageJSON: string) => {
			if (coverageJSON)
				writeFileSync(
					join(istanbulCLIOutput, `playwright_coverage_${generateUUID()}.json`),
					coverageJSON
				);
		});
		await use(context);
		for (const page of context.pages()) {
			await page.evaluate(() =>
				// eslint-disable-next-line @typescript-eslint/no-explicit-any
				(window as any).collectIstanbulCoverage(JSON.stringify((window as any).__coverage__))
			);
		}
	}
});

export const expect = test.expect;
