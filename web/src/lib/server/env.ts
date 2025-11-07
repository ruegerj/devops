const requiredVars = ['API_BASE_URL', 'ACCESS_TOKEN'] as const;

export function validateEnv(env: Record<string, string | undefined>) {
	const missing = requiredVars.filter((key) => !env[key]);
	console.log(missing);
	if (missing.length > 0) {
		throw new Error(`Missing required environement variables: ${missing.join(', ')}`);
	}
}
