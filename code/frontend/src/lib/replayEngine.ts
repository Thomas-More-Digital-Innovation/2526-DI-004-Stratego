import type { Position } from './types/game';

export interface PieceData {
    type: string;
    rank: string;
    ownerId: number;
}

export interface HistoricalMove {
    moveIndex: number;
    playerId: number;
    fromX: number;
    fromY: number;
    toX: number;
    toY: number;
    attacker?: PieceData;
    defender?: PieceData;
    result: 'move' | 'win' | 'loss' | 'tie' | 'capture';
}

export interface GameHistory {
    gameId: string;
    initialState: (PieceData | null)[][];
    moves: HistoricalMove[];
    winnerId?: number | null;
}

/**
 * Reconstructs the board state at a specific move index.
 * @param history The full game history containing initial setup and moves.
 * @param targetMoveIndex The index (1-based) representing which state to return. 
 *                        0 returns the initial layout.
 */
export function getBoardAtMove(history: GameHistory, targetMoveIndex: number): (PieceData | null)[][] {
    // Deep clone initial setup (10x10)
    const board: (PieceData | null)[][] = history.initialState.map(row => row.map(cell => cell ? { ...cell } : null));

    // Apply moves sequentially
    for (let i = 0; i < targetMoveIndex && i < history.moves.length; i++) {
        applyMove(board, history.moves[i]);
    }

    return board;
}

function applyMove(board: (PieceData | null)[][], move: HistoricalMove) {
    const { fromX, fromY, toX, toY, result } = move;
    const attacker = board[fromY][fromX];

    if (!attacker) return;

    // Movement: Clear the original square
    board[fromY][fromX] = null;

    switch (result) {
        case 'move':
        case 'win':
        case 'capture':
            board[toY][toX] = attacker;
            break;
        case 'loss':
            break;
        case 'tie':
            board[toY][toX] = null;
            break;
    }
}
