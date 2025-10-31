import '$lib/server/coverage.handler';
import { validateEnv } from '$lib/server/env';

validateEnv(); // ensure mandatory env vars are present
