// Helper to check authentication status by fetching current user from backend
import { user } from './store.svelte';

export async function checkAuthStatus(): Promise<boolean> {
    try {
        const response = await fetch('http://localhost:8080/api/users/me', {
            credentials: 'include'
        });

        if (response.ok) {
            const userData = await response.json();
            user.set(userData);
            localStorage.setItem('user', JSON.stringify(userData));
            return true;
        } else {
            // Session expired or not logged in
            user.clear();
            localStorage.removeItem('user');
            return false;
        }
    } catch (error) {
        console.error('Failed to check auth status:', error);
        return false;
    }
}

export async function logout(): Promise<void> {
    try {
        await fetch('http://localhost:8080/api/users/logout', {
            method: 'POST',
            credentials: 'include'
        });
    } catch (error) {
        console.error('Logout failed:', error);
    } finally {
        user.clear();
        localStorage.removeItem('user');
    }
}
