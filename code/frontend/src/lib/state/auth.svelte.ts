import type { User } from '$lib/types/game';
import { auth } from '$lib/api/client';

class AuthStore {
    user = $state<User | null>(null);
    loading = $state(true);

    get isLoggedIn() {
        return this.user !== null;
    }

    private checkPromise: Promise<void> | null = null;

    async check() {
        if (this.checkPromise) return this.checkPromise;

        this.checkPromise = (async () => {
            this.loading = true;
            try {
                this.user = await auth.getMe();
            } catch {
                this.user = null;
            } finally {
                this.loading = false;
            }
        })();

        return this.checkPromise;
    }

    async login(username: string, password: string) {
        await auth.login(username, password);
        this.checkPromise = null;
        await this.check();
    }

    async register(username: string, password: string) {
        await auth.register(username, password);
        this.checkPromise = null;
        await this.check();
    }

    async logout() {
        await auth.logout();
        this.user = null;
        this.checkPromise = null;
    }
}

export const authStore = new AuthStore();
