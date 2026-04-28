import type { GameState, BoardState, HistoryMove, Piece, CombatAnimation } from '$lib/types/game';
import { getBoardAtMove, type PieceData, type GameHistory } from '$lib/replayEngine';
class GameStore {
    gameState = $state<GameState | null>(null);
    boardState = $state<BoardState | null>(null);
    history = $state<HistoryMove[]>([]);
    currentHistoryIndex = $state(-1);
    isReplaying = $state(false);
    selectedPosition = $state<{ x: number; y: number } | null>(null);
    combatAnimation = $state<CombatAnimation | null>(null);
    isStepping = $state(false);
    gameMode = $state<string>("");
    lastLiveBoard = $state<BoardState | null>(null);

    get isPaused() {
        return this.gameState?.paused ?? false;
    }

    updateGameState(state: GameState) {
        this.gameState = state;
        this.isStepping = false;
    }

    updateBoardState(board: BoardState, viewerId: number = -1) {
        this.lastLiveBoard = board;

        if (!this.isReplaying) {
            this.boardState = board;
        }

        if (this.gameState && !this.gameState.isSetupPhase) {
            this.addToHistory(board.board, viewerId);
        }
    }

    private addToHistory(board: Piece[][], viewerId: number) {
        const boardCopy = board.map(row => row.map(cell => {
            const newCell = { ...cell };
            // Strip opponent piece identity in history if it's Human vs AI (viewerId !== -1)
            if (viewerId !== -1 && newCell.ownerId !== viewerId) {
                newCell.type = undefined;
                newCell.rank = undefined;
                newCell.iconBlue = undefined;
                newCell.iconRed = undefined;
            }
            return newCell;
        }));

        this.history.push({
            moveNumber: this.history.length,
            boardState: boardCopy,
            move: { from: { x: 0, y: 0 }, to: { x: 0, y: 0 } },
        });

        if (!this.isReplaying) {
            this.currentHistoryIndex = this.history.length - 1;
        }
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

    loadMoveHistory(data: { moves: Array<{ from: { x: number; y: number }; to: { x: number; y: number } }>; fullHistory: any[]; initialState: any[][] }, gameId: string = '', viewerId: number = -1) {
        if (!data.fullHistory || !data.initialState) {
            this.history = data.moves.map((move, index) => ({
                moveNumber: index,
                move,
                boardState: [],
            }));
        } else {
            const gameHistory: GameHistory = {
                gameId: gameId,
                initialState: data.initialState,
                moves: data.fullHistory,
            };

            this.history = data.moves.map((move, index) => {
                const moveBoard = getBoardAtMove(gameHistory, index + 1);

                // Map PieceData back to Piece interface
                const mappedBoard = moveBoard.map((row: (PieceData | null)[], y: number) => row.map((cell: PieceData | null, x: number) => {
                    if (!cell) return null;
                    const revealed = this.gameMode === "ai_vs_ai" || (this.gameState?.isGameOver ?? false) || cell.ownerId === viewerId;
                    return {
                        type: cell.type,
                        rank: cell.rank,
                        ownerId: cell.ownerId,
                        revealed: revealed,
                        position: { x, y }
                    } as Piece;
                }));

                return {
                    moveNumber: index,
                    move,
                    boardState: mappedBoard as Piece[][],
                };
            });
        }
        this.currentHistoryIndex = this.history.length > 0 ? this.history.length - 1 : -1;

        // auto go to last move if game is over
        if (this.gameState?.isGameOver && !this.isReplaying && this.history.length > 0) {
            this.goToMove(this.history.length - 1);
        }
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
        if (this.gameState?.isGameOver) return;
        this.isReplaying = false;
        if (this.lastLiveBoard && this.boardState) {
            this.boardState.board = this.lastLiveBoard.board;
            this.currentHistoryIndex = this.history.length - 1;
        }
    }

    reset() {
        this.gameState = null;
        this.boardState = null;
        this.lastLiveBoard = null;
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
