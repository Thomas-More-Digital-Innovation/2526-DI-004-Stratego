import type { GameInfo, GameMode, User, UserStats } from '$lib/types/game';
import type { BoardSetup } from '$lib/types/board-setup';
import { parseApiMessage } from '$lib/utils/api';

const API_BASE = import.meta.env.VITE_API_BASE || 'http://localhost:8080';

function getCookie(name: string): string {
    if (typeof document === 'undefined') return '';
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop()?.split(';').shift() || '';
    return '';
}

async function request<T>(path: string, options?: RequestInit): Promise<T> {
    const response = await fetch(`${API_BASE}${path}`, {
        credentials: 'include',
        ...options,
        headers: {
            'Content-Type': 'application/json',
            'X-XSRF-TOKEN': getCookie('XSRF-TOKEN'),
            ...options?.headers,
        },
    });

    if (response.status === 401 && !path.includes('/refresh') && !path.includes('/login') && !path.includes('/password')) {
        try {
            await auth.refresh();
            return request<T>(path, options);
        } catch (err) {
            throw new Error('Session expired');
        }
    }

    if (!response.ok) {
        const text = await response.text();
        const { message, type } = parseApiMessage(text);
        const error = new Error(message) as any;
        error.type = type;
        throw error;
    }

    return response.json();
}

async function requestVoid(path: string, options?: RequestInit): Promise<void> {
    const response = await fetch(`${API_BASE}${path}`, {
        credentials: 'include',
        ...options,
        headers: {
            'Content-Type': 'application/json',
            'X-XSRF-TOKEN': getCookie('XSRF-TOKEN'),
            ...options?.headers,
        },
    });

    if (response.status === 401 && !path.includes('/refresh') && !path.includes('/login') && !path.includes('/password')) {
        try {
            await auth.refresh();
            return requestVoid(path, options);
        } catch (err) {
            throw new Error('Session expired');
        }
    }

    if (!response.ok) {
        const text = await response.text();
        const { message, type } = parseApiMessage(text);
        const error = new Error(message) as any;
        error.type = type;
        throw error;
    }
}

// Auth
export const auth = {
    login: (username: string, password: string) =>
        requestVoid('/users/login', {
            method: 'POST',
            body: JSON.stringify({ username, password }),
        }),

    register: (username: string, password: string) =>
        requestVoid('/users/register', {
            method: 'POST',
            body: JSON.stringify({ username, password }),
        }),

    refresh: () => requestVoid('/users/refresh', { method: 'POST' }),

    logout: () => requestVoid('/users/logout', { method: 'POST' }),

    getMe: () => request<User>('/users/me'),

    changePassword: (oldPassword: string, newPassword: string, confirmPassword: string) =>
        requestVoid('/users/me/password', {
            method: 'POST',
            body: JSON.stringify({
                old_password: oldPassword,
                new_password: newPassword,
                confirm_password: confirmPassword
            }),
        }),
};

// Games
export const games = {
    create: (gameType: string, ai1: string, ai2: string) =>
        request<GameInfo>('/games', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ gameType, ai1, ai2 }),
        }),

    list: () => request<GameInfo[]>('/games'),

    getReconnectable: () =>
        request<{ hasGame: boolean; gameId?: string; gameType?: string }>('/users/me/reconnectable'),
};

// Stats
export const stats = {
    getMine: () => request<UserStats>('/users/me/stats'),
};

// Board Setups
export const boardSetups = {
    list: () => request<BoardSetup[]>('/board-setups'),

    getOne: (id: number) => request<BoardSetup>(`/board-setups/${id}`),

    create: (data: { name: string; description: string; setup_data: string; is_default: boolean }) =>
        requestVoid('/board-setups', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data),
        }),

    update: (id: number, data: { name: string; description: string; setup_data: string; is_default: boolean }) =>
        requestVoid(`/board-setups/${id}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data),
        }),

    delete: (id: number) =>
        requestVoid(`/board-setups/${id}`, { method: 'DELETE' }),
};

// Monitoring
export const monitoring = {
    health: () => request<{ status: string }>('/health'),
};
