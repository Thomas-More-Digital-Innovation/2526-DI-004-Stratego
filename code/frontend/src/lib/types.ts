
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
	icon?: string;
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
}

export interface MoveResult {
	success: boolean;
	error?: string;
	move?: Move;
	attackerDied: boolean;
	defenderDied: boolean;
	combatResult?: CombatResult;
}

export interface CombatResult {
	attackerRank: string;
	defenderRank: string;
	attackerRevealed: boolean;
	defenderRevealed: boolean;
}

export interface BoardState {
	board: Piece[][];
	width: number;
	height: number;
}

export interface GameInfo {
	gameId: string;
	gameType: 'human-vs-ai' | 'ai-vs-ai' | 'human-vs-human';
	wsUrl: string;
}

export interface HistoryMove {
	moveNumber: number;
	move: Move;
	piece?: Piece;
	result?: MoveResult;
	boardState: Piece[][];
}

export type GameMode = 'human-vs-ai' | 'ai-vs-ai';
