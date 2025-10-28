// filepath: /home/sem/prog/go/2526-DI-004-Stratego/code/frontend/src/lib/store.svelte.ts

import type { GameState, BoardState, HistoryMove, Piece } from './types';

export class GameStore {
	gameState = $state<GameState | null>(null);
	boardState = $state<BoardState | null>(null);
	history = $state<HistoryMove[]>([]);
	currentHistoryIndex = $state<number>(-1);
	isReplaying = $state<boolean>(false);
	selectedPosition = $state<{ x: number; y: number } | null>(null);

	updateGameState(state: GameState) {
		this.gameState = state;
	}

	updateBoardState(board: BoardState) {
		this.boardState = board;
		
		// Add to history if not replaying
		if (!this.isReplaying && this.gameState) {
			this.addToHistory(board.board);
		}
	}

	private addToHistory(board: Piece[][]) {
		// Clone the board for history
		const boardCopy = board.map(row => row.map(cell => ({ ...cell })));
		
		this.history.push({
			moveNumber: this.history.length,
			boardState: boardCopy,
			move: { from: { x: 0, y: 0 }, to: { x: 0, y: 0 } } // Will be updated with actual move
		});
		
		this.currentHistoryIndex = this.history.length - 1;
	}

	setSelectedPosition(pos: { x: number; y: number } | null) {
		this.selectedPosition = pos;
	}

	// History navigation
	goToMove(index: number) {
		if (index < 0 || index >= this.history.length) return;
		
		this.isReplaying = true;
		this.currentHistoryIndex = index;
		
		const historyItem = this.history[index];
		if (this.boardState) {
			this.boardState.board = historyItem.boardState;
		}
	}

	nextMove() {
		if (this.currentHistoryIndex < this.history.length - 1) {
			this.goToMove(this.currentHistoryIndex + 1);
		}
	}

	previousMove() {
		if (this.currentHistoryIndex > 0) {
			this.goToMove(this.currentHistoryIndex - 1);
		}
	}

	exitReplay() {
		this.isReplaying = false;
		if (this.history.length > 0) {
			this.goToMove(this.history.length - 1);
		}
	}

	reset() {
		this.gameState = null;
		this.boardState = null;
		this.history = [];
		this.currentHistoryIndex = -1;
		this.isReplaying = false;
		this.selectedPosition = null;
	}

	// Save game to JSON
	exportGame() {
		return JSON.stringify({
			history: this.history,
			gameState: this.gameState
		}, null, 2);
	}

	// Load game from JSON
	importGame(jsonData: string) {
		try {
			const data = JSON.parse(jsonData);
			this.history = data.history || [];
			this.gameState = data.gameState || null;
			
			if (this.history.length > 0) {
				this.goToMove(0);
			}
		} catch (error) {
			console.error('Failed to import game:', error);
			throw error;
		}
	}
}

export const gameStore = new GameStore();
