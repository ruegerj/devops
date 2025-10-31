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
		// instrument files for coverage collection in e2e
		...(process.env.USE_PLUGIN_ISTANBUL
			? [
					IstanbulPlugin({
						include: 'src/**/*',
						exclude: ['node_modules', 'e2e', '**/*.spec.ts', '**/*.spec.ts'],
						extension: ['.svelte', '.ts'],
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
			exclude: ['node_modules', 'e2e', '**/*.spec.ts', 'src/lib/components/ui/*'],
			include: ['src/**/*.ts', 'src/**/*.svelte'],
			extension: ['.svelte', '.ts']
		}
	}
});
