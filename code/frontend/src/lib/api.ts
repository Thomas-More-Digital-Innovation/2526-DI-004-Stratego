import type { GameInfo, GameMode, Position, Move } from './types';

const API_BASE = 'http://localhost:8080';
const WS_BASE = 'ws://localhost:8080';

export class GameAPI {
	private ws: WebSocket | null = null;
	private messageHandlers: Map<string, (data: any) => void> = new Map();

	// REST API
	async createGame(gameType: GameMode): Promise<GameInfo> {
		const response = await fetch(`${API_BASE}/api/games`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ gameType })
		});

		if (!response.ok) {
			throw new Error('Failed to create game');
		}

		return response.json();
	}

	async listGames(): Promise<any[]> {
		const response = await fetch(`${API_BASE}/api/games`);
		if (!response.ok) {
			throw new Error('Failed to list games');
		}
		return response.json();
	}

	// WebSocket
	connectWebSocket(gameId: string, playerId: number = -1): Promise<void> {
		return new Promise((resolve, reject) => {
			const wsUrl = `${WS_BASE}/ws/game/${gameId}?player=${playerId}`;
			this.ws = new WebSocket(wsUrl);

			this.ws.onopen = () => {
				console.log('WebSocket connected');
				resolve();
			};

			this.ws.onerror = (error) => {
				console.error('WebSocket error:', error);
				reject(error);
			};

			this.ws.onmessage = (event) => {
				try {
					const message = JSON.parse(event.data);
					const handler = this.messageHandlers.get(message.type);
					if (handler) {
						handler(message.data);
					}
				} catch (error) {
					console.error('Failed to parse WebSocket message:', error);
				}
			};

			this.ws.onclose = () => {
				console.log('WebSocket closed');
			};
		});
	}

	onMessage(type: string, handler: (data: any) => void) {
		this.messageHandlers.set(type, handler);
	}

	sendMove(from: Position, to: Position) {
		if (!this.ws) return;
		this.ws.send(JSON.stringify({
			type: 'move',
			data: { from, to }
		}));
	}

	requestValidMoves(position: Position) {
		if (!this.ws) return;
		this.ws.send(JSON.stringify({
			type: 'getValidMoves',
			data: { position }
		}));
	}

	sendAnimationComplete() {
		if (!this.ws) return;
		console.log('Sending animation complete signal to backend');
		this.ws.send(JSON.stringify({ type: 'animationComplete' }));
	}

	disconnect() {
		if (this.ws) {
			this.ws.close();
			this.ws = null;
		}
		this.messageHandlers.clear();
	}
}
