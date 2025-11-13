import { writeFileSync, mkdirSync } from 'fs';
import { join } from 'path';

// dump coverage data for instrumented builds
if (process.env.USE_PLUGIN_ISTANBUL) {
	const coverageDir = join(process.cwd(), 'coverage', 'e2e');
	mkdirSync(coverageDir, { recursive: true });

	const dumpCoverage = () => {
		const coverage = globalThis.__coverage__;
		if (!coverage) {
			console.warn('[coverage] No coverage data found on exit');
			return;
		}

		const filePath = join(coverageDir, `server_coverage_${Date.now()}.json`);
		writeFileSync(filePath, JSON.stringify(coverage));
		console.log(`[coverage] Dumped server coverage to: ${filePath}`);
	};

	// Handle all common exit signals
	process.on('exit', dumpCoverage);
	process.on('SIGINT', () => {
		dumpCoverage();
		process.exit(0);
	});
	process.on('SIGTERM', () => {
		dumpCoverage();
		process.exit(0);
	});
	process.on('uncaughtException', (err) => {
		console.error('[coverage] Uncaught exception:', err);
		dumpCoverage();
		process.exit(1);
	});
}
