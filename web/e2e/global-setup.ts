import { GenericContainer, Wait } from 'testcontainers';
import jwt from 'jsonwebtoken';
import waitOn from 'wait-on';
import { join } from 'node:path';
import { cwd } from 'node:process';
import { spawn } from 'node:child_process';
import { execSync } from 'node:child_process';

async function globalSetup() {
	const apiPort = 3000;
	const jwtSigningKey = 'v3ryS3cure!';

	// build & start testcontainer for api
	const image = await GenericContainer.fromDockerfile(join(cwd(), '..', 'api')).build();
	const container = await image
		.withEnvironment({ JWT_KEY: jwtSigningKey })
		.withExposedPorts(apiPort)
		.withWaitStrategy(Wait.forHttp('/health', apiPort))
		.start();

	const apiUrl = `http://${container.getHost()}:${container.getMappedPort(apiPort)}`;
	console.log(`Api container started at: ${apiUrl}`);

	const accessToken = jwt.sign({ sub: 'john.doe' }, jwtSigningKey, { algorithm: 'HS256' });

	// ensure preview launches latest build
	execSync('pnpm run build');

	// launch sveltekit instance
	const svelteKitProcess = spawn('pnpm', ['run', 'preview'], {
		cwd: cwd(),
		env: {
			...process.env,
			API_BASE_URL: apiUrl,
			ACCESS_TOKEN: accessToken
		},
		stdio: 'inherit'
	});

	await waitOn({
		resources: ['http://localhost:4173'],
		timeout: 30_000, // 30sec
		validateStatus: (status) => status == 200
	});

	console.log('SvelteKit instance started');

	// setup teardown hook
	process.env.__API_CONTAINER_ID__ = container.getId();
	process.on('exit', async () => {
		svelteKitProcess.kill();
		await container.stop();
	});
}

export default globalSetup;
