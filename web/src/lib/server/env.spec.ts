import { afterEach, beforeEach, describe, expect, it } from 'vitest';
import { validateEnv } from './env';

describe('validateEnv', () => {
	let originEnv: Record<string, string | undefined>;

	beforeEach(() => {
		originEnv = { ...process.env };
	});

	afterEach(() => {
		process.env = originEnv;
	});

	it("should throw if no required var isn't present", () => {
		// GIVEN
		delete process.env.API_BASE_URL;
		delete process.env.ACCESS_TOKEN;

		// WHEN & THEN
		expect(() => validateEnv(process.env)).toThrow();
	});

	it('should throw if some required vars are missing', () => {
		// GIVEN
		delete process.env.API_BASE_URL;
		process.env.ACCESS_TOKEN = 'foo.bar';

		// WHEN & THEN
		expect(() => validateEnv(process.env)).toThrow();

		// GIVEN
		process.env.API_BASE_URL = 'http://foo.bar';
		delete process.env.ACCESS_TOKEN;

		// WHEN & THEN
		expect(() => validateEnv(process.env)).toThrow();
	});

	it('should pass if all required vars are present', () => {
		// GIVEN
		process.env.API_BASE_URL = 'http://foo.bar';
		process.env.ACCESS_TOKEN = 'foo.bar';

		// WHEN & THEN
		expect(() => validateEnv(process.env)).not.toThrow();
	});
});
