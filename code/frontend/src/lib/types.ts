// Type definitions for Stratego

export type PieceType = {
	name: string;
	rank: string;
	movable: boolean;
	description: string;
	icon: string;
	count: number;
	strategicValue: number;
};

export type Position = {
	x: number;
	y: number;
};

export type Piece = {
	type: PieceType;
	owner: Player;
	alive: boolean;
	revealed: boolean;
	position: Position;
};

export type Player = {
	id: number;
	name: string;
	color: string;
	pieces: Piece[];
};

export type Move = {
	from: Position;
	to: Position;
};

export type GamePhase = 'setup' | 'playing' | 'gameOver';

export type WinCause = 'flag_captured' | 'no_movable_pieces' | 'max_turns';

export type GameState = {
	phase: GamePhase;
	players: [Player, Player];
	currentPlayerIndex: number;
	board: (Piece | null)[][];
	selectedPiece: Piece | null;
	validMoves: Position[];
	winner: Player | null;
	winCause: WinCause | null;
	round: number;
	moveHistory: Move[];
};

// Piece type definitions matching the backend
export const PIECE_TYPES: Record<string, PieceType> = {
	FLAG: { name: 'Flag', rank: '0', movable: false, description: 'The piece you must capture to win.', icon: '🚩', count: 1, strategicValue: 0 },
	BOMB: { name: 'Bomb', rank: 'B', movable: false, description: 'Eliminates most attackers.', icon: '💣', count: 6, strategicValue: 7 },
	SPY: { name: 'Spy', rank: '1', movable: true, description: 'Can defeat the Marshal.', icon: '🕵️', count: 1, strategicValue: 7 },
	SCOUT: { name: 'Scout', rank: '2', movable: true, description: 'Can move multiple spaces.', icon: '🔭', count: 8, strategicValue: 3 },
	MINER: { name: 'Miner', rank: '3', movable: true, description: 'Can defuse bombs.', icon: '⛏️', count: 5, strategicValue: 6 },
	SERGEANT: { name: 'Sergeant', rank: '4', movable: true, description: 'Standard combat piece.', icon: '⚔️', count: 4, strategicValue: 4 },
	LIEUTENANT: { name: 'Lieutenant', rank: '5', movable: true, description: 'Standard combat piece.', icon: '🗡️', count: 4, strategicValue: 5 },
	CAPTAIN: { name: 'Captain', rank: '6', movable: true, description: 'Standard combat piece.', icon: '🛡️', count: 4, strategicValue: 6 },
	MAJOR: { name: 'Major', rank: '7', movable: true, description: 'Standard combat piece.', icon: '⚜️', count: 3, strategicValue: 7 },
	COLONEL: { name: 'Colonel', rank: '8', movable: true, description: 'High-ranking officer.', icon: '👑', count: 2, strategicValue: 8 },
	GENERAL: { name: 'General', rank: '9', movable: true, description: 'Second strongest piece.', icon: '🎖️', count: 1, strategicValue: 9 },
	MARSHAL: { name: 'Marshal', rank: 'M', movable: true, description: 'Strongest piece.', icon: '⭐', count: 1, strategicValue: 10 }
};
