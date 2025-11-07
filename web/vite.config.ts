import tailwindcss from '@tailwindcss/vite';
import { defineConfig } from 'vitest/config';
import { sveltekit } from '@sveltejs/kit/vite';
import Icons from 'unplugin-icons/vite';
import IstanbulPlugin from 'vite-plugin-istanbul';

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export default defineConfig(({ mode }) => {
	const isVitest = process.env.VITEST === 'true';
	const useIstanbul = process.env.USE_PLUGIN_ISTANBUL === 'true';

	return {
		plugins: [
			tailwindcss(),
			sveltekit(),
			Icons({
				compiler: 'svelte'
			}),
			// instrument files for clientside coverage collection in e2e
			...(useIstanbul
				? [
						IstanbulPlugin({
							forceBuildInstrument: true,
							requireEnv: isVitest
						})
					]
				: [])
		],
		build: {
			sourcemap: process.env.NODE_ENV != 'production'
		},
		resolve: process.env.VITEST
			? {
					conditions: ['browser']
				}
			: undefined,
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
	};
});
