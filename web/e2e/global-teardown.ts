import { execSync } from 'child_process';

export default async function globalTeardown() {
	const sveltePid = process.env.SVELTE_PID;
	const apiContainerId = process.env.API_CONTAINER_ID;

	if (sveltePid) {
		console.log(`Stopping Svelte app (pid: ${sveltePid})...`);
		process.kill(-parseInt(sveltePid), 'SIGINT'); // kill whole process group (including children)
		await new Promise((res) => setTimeout(res, 2000));
		console.log('Stopped app, grace period for shutdown elapsed...');
	}

	if (apiContainerId) {
		console.log(`Stopping api container (id: ${apiContainerId})...`);
		execSync(`docker stop ${apiContainerId}`, { stdio: 'inherit' });
		execSync(`docker rm -f ${apiContainerId}`, { stdio: 'inherit' });
		console.log('Stopped & removed container');
	}
}
