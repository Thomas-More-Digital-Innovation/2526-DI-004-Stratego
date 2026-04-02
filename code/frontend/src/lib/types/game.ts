export interface Position {
    x: number;
    y: number;
}

export interface Piece {
    type?: string;
    rank?: string;
    ownerId: number;
    ownerName?: string;
    revealed: boolean;
    iconBlue?: string;
    iconRed?: string;
    position: Position;
}

export interface Move {
    from: Position;
    to: Position;
}

export interface GameState {
    round: number;
    currentPlayerId: number;
    currentPlayerName: string;
    isGameOver: boolean;
    winnerId?: number;
    winnerName?: string;
    winCause?: string;
    player1Score: number;
    player2Score: number;
    waitingForInput: boolean;
    moveCount: number;
    player1AlivePieces: number;
    player2AlivePieces: number;
    isSetupPhase: boolean;
}

export interface BoardState {
    board: Piece[][];
    width: number;
    height: number;
}

export interface GameInfo {
    gameId: string;
    gameType: GameMode;
    wsUrl: string;
}

export interface HistoryMove {
    moveNumber: number;
    move: Move;
    piece?: Piece;
    boardState: Piece[][];
}

export interface CombatAnimation {
    attacker: Piece;
    defender: Piece;
    attackerWon: boolean;
    defenderWon: boolean;
}

export type GameMode = 'human_vs_ai' | 'ai_vs_ai' | 'human_vs_human';

export interface AI {
    name: string;
    id: string;
    description: string;
    image?: string;
}

export interface User {
    id: number;
    username: string;
    profile_picture?: string;
    created_at: string;
    updated_at: string;
}

export interface UserStats {
    total_games: number;
    wins: number;
    losses: number;
    draws: number;
    total_moves: number;
    avg_game_duration_seconds: number;
}
