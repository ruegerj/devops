import '$lib/server/coverage.handler';
import { env } from '$env/dynamic/private';
import { isTelemetryEnabled, validateEnv } from '$lib/server/env';
import { otelSdk } from './metrics.server';

validateEnv(env); // ensure mandatory env vars are present

if (isTelemetryEnabled()) {
	otelSdk.start();
}
