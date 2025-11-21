import { redirect } from '@sveltejs/kit';
import type { LayoutServerLoad } from './$types';

interface User {
	id: number;
	username: string;
	profile_picture?: string;
	created_at: string;
	updated_at: string;
}

export const load: LayoutServerLoad = async ({ cookies, fetch }) => {
	const sessionId = cookies.get('session_id');

	if (!sessionId) {
		throw redirect(303, '/');
	}

	// Verify session with backend
	try {
		const response = await fetch('http://backend-dev:8080/api/users/me', {
			headers: {
				Cookie: `session_id=${sessionId}`
			}
		});

		if (!response.ok) {
			// Session invalid, clear cookie and redirect
			cookies.delete('session_id', { path: '/' });
			throw redirect(303, '/');
		}

		const user: User = await response.json();

		return {
			user
		};
	} catch (error) {
		// Network error or invalid session
		cookies.delete('session_id', { path: '/' });
		throw redirect(303, '/');
	}
};
