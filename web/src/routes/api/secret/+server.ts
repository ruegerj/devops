import { env } from '$env/dynamic/private';

export async function GET({ request }) {
	const baseUrl = env.API_BASE_URL;
	const url = new URL('/api/secret', baseUrl);

	const reqHeaders = request.headers;
	reqHeaders.delete('Host');

	const accessToken = env.ACCESS_TOKEN;
	reqHeaders.set('Authorization', `Bearer ${accessToken}`);

	const response = await fetch(url.toString(), {
		method: 'GET',
		headers: reqHeaders,
		duplex: 'half'
	} as RequestInit);

	const resHeaders = new Headers(response.headers);
	resHeaders.delete('content-encoding');
	resHeaders.delete('content-length');
	resHeaders.delete('transfer-encoding');

	return new Response(response.body, {
		status: response.status,
		statusText: response.statusText,
		headers: resHeaders
	});
}
