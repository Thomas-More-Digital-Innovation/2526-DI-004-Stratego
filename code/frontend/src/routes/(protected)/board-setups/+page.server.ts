import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ parent, fetch, cookies }) => {
	const { user } = await parent();
	const sessionId = cookies.get('session_id');

	try {
		const response = await fetch('http://backend-dev:8080/api/board-setups', {
			headers: {
				Cookie: `session_id=${sessionId}`
			}
		});

		if (!response.ok) {
			return {
				user,
				setups: [],
				error: 'Failed to load board setups'
			};
		}

		const setups = await response.json();
		if (setups === null) {
			return {
				user,
				setups: [],
				error: null
			};
		}
		return {
			user,
			setups,
			error: null
		};
	} catch (error) {
		return {
			user,
			setups: [],
			error: 'Network error while loading setups'
		};
	}
};
