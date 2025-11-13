import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
	globalSetup: 'e2e/global-setup.ts',
	globalTeardown: 'e2e/global-teardown.ts',
	testDir: 'e2e',
	projects: [
		{
			name: 'chromium',
			use: { ...devices['Desktop Chrome'] }
		}
	],
	use: {
		baseURL: 'http://localhost:4173'
	}
});
