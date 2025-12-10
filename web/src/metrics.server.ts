import { env } from '$env/dynamic/private';
import { isTelemetryEnabled } from '$lib/server/env';
import { getNodeAutoInstrumentations } from '@opentelemetry/auto-instrumentations-node';
import { PrometheusExporter } from '@opentelemetry/exporter-prometheus';
import { OTLPTraceExporter } from '@opentelemetry/exporter-trace-otlp-proto';
import { NodeSDK } from '@opentelemetry/sdk-node';
import { createAddHookMessageChannel } from 'import-in-the-middle';
import { register } from 'node:module';

const { registerOptions } = createAddHookMessageChannel();
register('import-in-the-middle/hook.mjs', import.meta.url, registerOptions);

let prometheusPort = 4318;
if (env.TELEMETRY_PORT) {
	prometheusPort = parseInt(env.TELEMETRY_PORT);
}

export const prometheusExporter = new PrometheusExporter({
	preventServerStart: !isTelemetryEnabled(),
	endpoint: '/metrics',
	port: prometheusPort
});

export const otelSdk = new NodeSDK({
	serviceName: 'sveltekit',
	traceExporter: new OTLPTraceExporter(),
	instrumentations: [getNodeAutoInstrumentations()],
	metricReader: prometheusExporter
});
