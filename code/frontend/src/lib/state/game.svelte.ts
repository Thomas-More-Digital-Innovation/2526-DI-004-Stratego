import type { GameState, BoardState, HistoryMove, Piece, CombatAnimation } from '$lib/types/game';

class GameStore {
    gameState = $state<GameState | null>(null);
    boardState = $state<BoardState | null>(null);
    history = $state<HistoryMove[]>([]);
    currentHistoryIndex = $state(-1);
    isReplaying = $state(false);
    selectedPosition = $state<{ x: number; y: number } | null>(null);
    combatAnimation = $state<CombatAnimation | null>(null);

    updateGameState(state: GameState) {
        this.gameState = state;
    }

    updateBoardState(board: BoardState) {
        this.boardState = board;

        if (!this.isReplaying && this.gameState && !this.gameState.isSetupPhase) {
            this.addToHistory(board.board);
        }
    }

    private addToHistory(board: Piece[][]) {
        const boardCopy = board.map(row => row.map(cell => ({ ...cell })));

        this.history.push({
            moveNumber: this.history.length,
            boardState: boardCopy,
            move: { from: { x: 0, y: 0 }, to: { x: 0, y: 0 } },
        });

        this.currentHistoryIndex = this.history.length - 1;
    }

    setSelectedPosition(pos: { x: number; y: number } | null) {
        this.selectedPosition = pos;
    }

    showCombatAnimation(combat: CombatAnimation) {
        this.combatAnimation = combat;
    }

    hideCombatAnimation() {
        this.combatAnimation = null;
    }

    loadMoveHistory(moves: Array<{ from: { x: number; y: number }; to: { x: number; y: number } }>) {
        this.history = moves.map((move, index) => ({
            moveNumber: index,
            move,
            boardState: [],
        }));
        this.currentHistoryIndex = this.history.length > 0 ? this.history.length - 1 : -1;
    }

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
            this.currentHistoryIndex = this.history.length - 1;
            const historyItem = this.history[this.currentHistoryIndex];
            if (this.boardState) {
                this.boardState.board = historyItem.boardState;
            }
        }
    }

    reset() {
        this.gameState = null;
        this.boardState = null;
        this.history = [];
        this.currentHistoryIndex = -1;
        this.isReplaying = false;
        this.selectedPosition = null;
        this.combatAnimation = null;
    }

    exportGame() {
        return JSON.stringify({
            history: this.history,
            gameState: this.gameState,
        }, null, 2);
    }
}

export const gameStore = new GameStore();
