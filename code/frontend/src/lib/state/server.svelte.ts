import { monitoring } from "$lib/api/client";

class ServerState {
    #isOnline = $state(true);
    #lastChecked = $state<Date | null>(null);
    #checkInterval: ReturnType<typeof setInterval> | null = null;

    get isOnline() {
        return this.#isOnline;
    }

    get lastChecked() {
        return this.#lastChecked;
    }

    async check() {
        // Quick check for browser's offline state
        if (typeof navigator !== "undefined" && !navigator.onLine) {
            this.#isOnline = false;
            this.#lastChecked = new Date();
            return false;
        }

        try {
            const response = await monitoring.health();
            this.#isOnline = response.status === "ok";
        } catch (error) {
            this.#isOnline = false;
        } finally {
            this.#lastChecked = new Date();
        }

        return this.#isOnline;
    }

    startPolling(ms = 30000) {
        if (this.#checkInterval) return;
        
        // Initial check
        this.check();

        this.#checkInterval = setInterval(() => {
            this.check();
        }, ms);

        if (typeof window !== "undefined") {
            window.addEventListener("online", () => this.check());
            window.addEventListener("offline", () => {
                this.#isOnline = false;
                this.#lastChecked = new Date();
            });
        }
    }

    stopPolling() {
        if (this.#checkInterval) {
            clearInterval(this.#checkInterval);
            this.#checkInterval = null;
        }
    }
}

export const serverStore = new ServerState();
