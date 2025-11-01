import tailwindcss from '@tailwindcss/vite';
import { defineConfig } from 'vitest/config';
import { sveltekit } from '@sveltejs/kit/vite';
import Icons from 'unplugin-icons/vite';
import IstanbulPlugin from 'vite-plugin-istanbul';

export default defineConfig({
	plugins: [
		tailwindcss(),
		sveltekit(),
		Icons({
			compiler: 'svelte'
		}),
		// instrument files for clientside coverage collection in e2e
		...(process.env.USE_PLUGIN_ISTANBUL
			? [
					IstanbulPlugin({
						forceBuildInstrument: true
					})
				]
			: [])
	],
	test: {
		expect: { requireAssertions: true },
		include: ['src/**/*.{test,spec}.{js,ts}'],
		coverage: {
			provider: 'istanbul',
			exclude: [
				'node_modules',
				'.svelte-kit',
				'build',
				'coverage',
				'e2e',
				'**/*.spec.ts',
				'src/lib/components/ui',
				'src/lib/server/coverage.handler.ts'
			],
			include: ['src/**/*.ts', 'src/**/*.svelte']
		}
	}
});
