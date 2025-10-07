import { env } from '$env/dynamic/private';

const requiredVars = ['API_BASE_URL', 'ACCESS_TOKEN'] as const;

export function validateEnv() {
	const missing = requiredVars.filter((key) => !env[key]);
	if (missing.length > 0) {
		throw new Error(`Missing required environement variables: ${missing.join(', ')}`);
	}
}
