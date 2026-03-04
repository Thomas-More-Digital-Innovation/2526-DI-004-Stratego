import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ cookies, fetch }: any) => {
	const sessionId = cookies.get('session_id');

	if (!sessionId) {
		return {
			user: null
		};
	}

	// Verify session with backend
	try {
		const response = await fetch('http://backend-dev:8080/api/users/me', {
			headers: {
				Cookie: `session_id=${sessionId}`
			}
		});

		if (!response.ok) {
			// Session invalid, clear cookie
			cookies.delete('session_id', { path: '/' });
			return {
				user: null
			};
		}

		const user = await response.json();

		return {
			user
		};
	} catch (error) {
		// Network error or invalid session
		cookies.delete('session_id', { path: '/' });
		return {
			user: null
		};
	}
};
