/* eslint-disable @typescript-eslint/no-require-imports */
const { readFile, writeFile, readdir, access, unlink } = require('node:fs/promises');
const { join } = require('node:path');
const { cwd } = require('node:process');
const { createCoverageMap } = require('istanbul-lib-coverage');
const { createSourceMapStore } = require('istanbul-lib-source-maps');
const { readFileSync } = require('node:fs');

const UNIT_TEST_COVERAGE_FILE = 'coverage-final.json';

// merges both unit and e2e test coverage / ensures proper source mappings
(async () => {
	const coverageDir = join(cwd(), 'coverage');
	const map = createCoverageMap({});

	await loadUnitTestCoverageDataInto(map, coverageDir);
	await loadE2ETestCoverageDataInto(map, coverageDir);

	const sourceMapStore = createSourceMapStore();
	const transformed = await sourceMapStore.transformCoverage(map);

	await removeUnitTestCoverageFile(coverageDir);

	const targetFilePath = join(coverageDir, 'coverage-combined.json');
	await writeFile(targetFilePath, JSON.stringify(transformed.toJSON()), 'utf-8');

	console.log(`Wrote combined test coverage data to ${targetFilePath}`);
})();

async function loadUnitTestCoverageDataInto(map, baseDir) {
	const content = await readFile(join(baseDir, UNIT_TEST_COVERAGE_FILE), 'utf-8');
	if (!content) {
		return;
	}

	const coverageData = JSON.parse(content);
	map.merge(coverageData);
}

async function loadE2ETestCoverageDataInto(map, baseDir) {
	const e2eCoverageDir = join(baseDir, 'e2e');
	const e2eCoverageFiles = await readdir(e2eCoverageDir);
	const coverageFileContents = e2eCoverageFiles
		.filter((f) => f.endsWith('.json'))
		.map((f) => {
			const content = readFileSync(join(e2eCoverageDir, f), 'utf-8');
			return JSON.parse(content);
		});

	coverageFileContents.forEach((fc) => map.merge(fc));
}

async function removeUnitTestCoverageFile(baseDir) {
	const path = join(baseDir, UNIT_TEST_COVERAGE_FILE);
	const error = await access(path);
	if (error) {
		return;
	}

	await unlink(path);
}
