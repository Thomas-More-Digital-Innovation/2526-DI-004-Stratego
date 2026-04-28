import type { GameMode, GameState, BoardState, HistoryMove, Piece, CombatAnimation } from '$lib/types/game';
import { getBoardAtMove, type PieceData, type GameHistory, type HistoricalMove } from '$lib/replayEngine';
import { gamemodes } from '$lib/data/gamemodes.data';
class GameStore {
    gameState = $state<GameState | null>(null);
    boardState = $state<BoardState | null>(null);
    history = $state<HistoryMove[]>([]);
    currentHistoryIndex = $state(-1);
    isReplaying = $state(false);
    selectedPosition = $state<{ x: number; y: number } | null>(null);
    combatAnimation = $state<CombatAnimation | null>(null);
    isStepping = $state(false);
    gameMode = $state<GameMode>(gamemodes.unknown);
    lastLiveBoard = $state<BoardState | null>(null);
    rawHistory = $state<GameHistory | null>(null);
    lastMove = $state<HistoricalMove | null>(null);


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
            this.lastMove = board.lastMove ?? null;
        }

        if (this.gameState && !this.gameState.isSetupPhase) {
            this.addToHistory(board.board, board.lastMove ?? null, viewerId);
        }
    }

    private addToHistory(board: Piece[][], lastMove: HistoricalMove | null, viewerId: number) {
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

        // Redact opponent piece info in history move to prevent leakage when scrubbing back
        let historyMove = lastMove;
        if (historyMove && viewerId !== -1) {
            if (historyMove.attacker && historyMove.attacker.ownerId !== viewerId) {
                historyMove = {
                    ...historyMove,
                    attacker: { ...historyMove.attacker, type: '', rank: '' }
                };
            }
            if (historyMove.defender && historyMove.defender.ownerId !== viewerId) {
                historyMove = {
                    ...historyMove,
                    defender: { ...historyMove.defender, type: '', rank: '' }
                };
            }
        }

        this.history.push({
            moveNumber: this.history.length,
            boardState: boardCopy,
            move: historyMove || { fromX: 0, fromY: 0, toX: 0, toY: 0, result: 'move', playerId: 0, moveIndex: this.history.length },
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
                move: {
                    moveIndex: index,
                    fromX: move.from.x,
                    fromY: move.from.y,
                    toX: move.to.x,
                    toY: move.to.y,
                    result: 'move',
                    playerId: 0,
                },
                boardState: [],
            }));
        } else {
            const gameHistory: GameHistory = {
                gameId: gameId,
                initialState: data.initialState,
                moves: data.fullHistory,
            };
            this.rawHistory = gameHistory;


            this.history = data.moves.map((move, index) => {
                const moveBoard = getBoardAtMove(gameHistory, index + 1);

                // Map PieceData back to Piece interface
                const mappedBoard = moveBoard.map((row: (PieceData | null)[], y: number) => row.map((cell: PieceData | null, x: number) => {
                    if (!cell) return null;
                    const revealed = this.gameMode.mode === gamemodes.ai_vs_ai.mode || (this.gameState?.isGameOver ?? false) || cell.ownerId === viewerId;
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
                    move: data.fullHistory[index],
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
            this.lastMove = historyItem.move;
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
            this.lastMove = this.lastLiveBoard.lastMove ?? null;
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
        this.gameMode = gamemodes.unknown;
        this.rawHistory = null;
        this.lastMove = null;
    }


    exportGame() {
        if (!this.rawHistory) return null;
        return JSON.stringify(this.rawHistory, null, 2);
    }

}

export const gameStore = new GameStore();
