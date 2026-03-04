import type { Handle } from '@sveltejs/kit';

export const handle: Handle = async ({ event, resolve }) => {
	// Pass through all requests
	const response = await resolve(event);
	return response;
};
