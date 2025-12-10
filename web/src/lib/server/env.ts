import { env } from '$env/dynamic/private';

const requiredVars = ['API_BASE_URL', 'ACCESS_TOKEN', 'TELEMETRY_ENABLED'] as const;

export function validateEnv(env: Record<string, string | undefined>) {
	const missing = requiredVars.filter((key) => !env[key]);
	if (missing.length > 0) {
		throw new Error(`Missing required environement variables: ${missing.join(', ')}`);
	}
}

export function isTelemetryEnabled(): boolean {
	return env.TELEMETRY_ENABLED === 'true';
}
