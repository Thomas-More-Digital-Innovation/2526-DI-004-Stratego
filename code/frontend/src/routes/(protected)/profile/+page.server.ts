import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ parent, fetch, cookies }) => {
	const { user } = await parent();
	const sessionId = cookies.get('session_id');

	try {
		const response = await fetch(`http://backend-dev:8080/api/users/stats?user_id=${user.id}`, {
			headers: {
				Cookie: `session_id=${sessionId}`
			}
		});

		if (!response.ok) {
			return {
				user,
				stats: null,
				error: 'Failed to load statistics'
			};
		}

		const stats = await response.json();

		return {
			user,
			stats,
			error: null
		};
	} catch (error) {
		return {
			user,
			stats: null,
			error: 'Network error while loading statistics'
		};
	}
};
