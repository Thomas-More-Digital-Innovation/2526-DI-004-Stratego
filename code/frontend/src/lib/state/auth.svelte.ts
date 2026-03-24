import type { User } from '$lib/types/game';
import { auth } from '$lib/api/client';

class AuthStore {
    user = $state<User | null>(null);
    loading = $state(true);

    get isLoggedIn() {
        return this.user !== null;
    }

    async check() {
        this.loading = true;
        try {
            this.user = await auth.getMe();
        } catch {
            this.user = null;
        } finally {
            this.loading = false;
        }
    }

    async login(username: string, password: string) {
        await auth.login(username, password);
        await this.check();
    }

    async register(username: string, password: string) {
        await auth.register(username, password);
        await this.check();
    }

    async logout() {
        await auth.logout();
        this.user = null;
    }
}

export const authStore = new AuthStore();
