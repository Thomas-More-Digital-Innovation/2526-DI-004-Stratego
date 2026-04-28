import type { GameInfo, GameMode, User, UserStats } from '$lib/types/game';
import type { BoardSetup } from '$lib/types/board-setup';

const API_BASE = import.meta.env.VITE_API_BASE || 'http://localhost:8080';

async function request<T>(path: string, options?: RequestInit): Promise<T> {
    const headers = {
        'Content-Type': 'application/json',
        'X-XSRF-TOKEN': '1',
        ...options?.headers,
    };

    const response = await fetch(`${API_BASE}${path}`, {
        credentials: 'include',
        ...options,
        headers,
    });

    if (!response.ok) {
        const text = await response.text();
        throw new Error(text || `Request failed: ${response.status}`);
    }

    return response.json();
}

async function requestVoid(path: string, options?: RequestInit): Promise<void> {
    const headers = {
        'Content-Type': 'application/json',
        'X-XSRF-TOKEN': '1',
        ...options?.headers,
    };

    const response = await fetch(`${API_BASE}${path}`, {
        credentials: 'include',
        ...options,
        headers,
    });

    if (!response.ok) {
        const text = await response.text();
        throw new Error(text || `Request failed: ${response.status}`);
    }
}

// Auth
export const auth = {
    login: (username: string, password: string) =>
        requestVoid('/users/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password }),
        }),

    register: (username: string, password: string) =>
        requestVoid('/users/register', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password }),
        }),

    logout: () => requestVoid('/users/logout', { method: 'POST' }),

    getMe: () => request<User>('/users/me'),
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
};

// Stats
export const stats = {
    getMine: () => request<UserStats>('/users/me/stats'),
};

// Board Setups
export const boardSetups = {
    list: () => request<BoardSetup[]>('/board-setups'),

    create: (data: { name: string; description: string; setup_data: string; is_default: boolean }) =>
        requestVoid('/board-setups', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data),
        }),

    update: (id: number, data: { name: string; description: string; setup_data: string; is_default: boolean }) =>
        requestVoid(`/board-setups?id=${id}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data),
        }),

    delete: (id: number) =>
        requestVoid(`/board-setups?id=${id}`, { method: 'DELETE' }),
};
