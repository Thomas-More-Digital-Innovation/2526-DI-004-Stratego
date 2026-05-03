import { parseApiMessage } from '$lib/utils/api';

export type ToastType = 'info' | 'success' | 'warning' | 'error';

export interface Toast {
    id: string;
    message: string;
    type: ToastType;
    duration?: number;
}

class ToastStore {
    toasts = $state<Toast[]>([]);

    add(message: string, type: ToastType = 'info', duration: number = 3000) {
        const id = Math.random().toString(36).substring(2, 9);
        const toast: Toast = { id, message, type, duration };
        this.toasts = [...this.toasts, toast];

        if (duration > 0) {
            setTimeout(() => {
                this.remove(id);
            }, duration);
        }

        return id;
    }

    /**
     * Handles an API message, using fallback if needed.
     */
    handleApiMessage(error: any, fallbackMessage: string = 'An error occurred') {
        // If error is an object but doesn't have a message and stringifies to '{}', use fallback
        let messageToParse = typeof error === 'string' ? error : error?.message;

        if (!messageToParse && error) {
            const stringified = JSON.stringify(error);
            messageToParse = stringified === '{}' ? '' : stringified;
        }

        const { message, type: parsedType } = parseApiMessage(messageToParse);
        const type = typeof error?.type === 'string' ? error.type : parsedType;
        // Map types to supported toast types
        let toastType: ToastType = 'error';
        if (type === 'warning') toastType = 'warning';
        if (type === 'info') toastType = 'info';
        if (type === 'success') toastType = 'success';

        const id = this.add(message || fallbackMessage, toastType);
        return { id, message: message || fallbackMessage, type };
    }

    success(message: string, duration?: number) {
        return this.add(message, 'success', duration);
    }

    error(message: string, duration?: number) {
        return this.add(message, 'error', duration);
    }

    warning(message: string, duration?: number) {
        return this.add(message, 'warning', duration);
    }

    info(message: string, duration?: number) {
        return this.add(message, 'info', duration);
    }

    remove(id: string) {
        this.toasts = this.toasts.filter(t => t.id !== id);
    }
}

export const toastStore = new ToastStore();
