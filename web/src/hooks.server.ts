import '$lib/server/coverage.handler';
import { env } from '$env/dynamic/private';
import { validateEnv } from '$lib/server/env';

validateEnv(env); // ensure mandatory env vars are present
