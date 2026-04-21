import type { Position } from '$lib/types/game';

const WS_BASE = import.meta.env.VITE_WS_BASE || 'ws://localhost:8080';

type MessageHandler = (data: any) => void;

export class GameSocket {
    private ws: WebSocket | null = null;
    private handlers = new Map<string, MessageHandler>();

    connect(gameId: string, playerId: number = -1): Promise<void> {
        return new Promise((resolve, reject) => {
            const url = `${WS_BASE}/game/${gameId}?player=${playerId}`;
            this.ws = new WebSocket(url);

            this.ws.onopen = () => resolve();
            this.ws.onerror = (e) => reject(e);

            this.ws.onmessage = (event) => {
                try {
                    const message = JSON.parse(event.data);
                    this.handlers.get(message.type)?.(message.data);
                } catch (error) {
                    console.error('Failed to parse WebSocket message:', error);
                }
            };

            this.ws.onclose = () => {
                console.log('WebSocket closed');
            };
        });
    }

    on(type: string, handler: MessageHandler) {
        this.handlers.set(type, handler);
    }

    private send(type: string, data?: Record<string, unknown>) {
        if (!this.ws) return;
        this.ws.send(JSON.stringify({ type, data }));
    }

    sendMove(from: Position, to: Position) {
        this.send('move', { from, to });
    }

    requestValidMoves(position: Position) {
        this.send('getValidMoves', { position });
    }

    sendSwapPieces(pos1: Position, pos2: Position) {
        this.send('swapPieces', { pos1, pos2 });
    }

    sendRandomizeSetup(playerId?: number) {
        this.send('randomizeSetup', { playerId });
    }

    sendStartGame(headless: boolean = false) {
        this.send('startGame', { headless });
    }

    sendAnimationComplete() {
        this.send('animationComplete');
    }

    sendPause() {
        this.send('pause');
    }

    sendUnpause() {
        this.send('unpause');
    }

    sendLoadSetup(setupData: string, playerId?: number) {
        this.send('loadSetup', { setupData, playerId });
    }

    sendSetSpeed(speedMs: number) {
        this.send('setSpeed', { speedMs });
    }

    sendStep() {
        this.send('step');
    }

    disconnect() {
        this.ws?.close();
        this.ws = null;
        this.handlers.clear();
    }
}
