export interface BoardSetup {
	name: string,
	description: string,
	isDefault: boolean,
	setupData: string
}
export interface PieceInfo {
	name: string;
	icon: string;
	rank: string;
	count: number;
}

export const PIECE_INVENTORY: Record<string, PieceInfo> = {
	'0': { name: 'Flag', icon: '🚩', rank: 'F', count: 1 },
	'B': { name: 'Bomb', icon: '💣', rank: 'B', count: 6 },
	'1': { name: 'Spy', icon: '🕵️', rank: 'S', count: 1 },
	'2': { name: 'Scout', icon: '🔭', rank: '2', count: 8 },
	'3': { name: 'Miner', icon: '⛏️', rank: '3', count: 5 },
	'4': { name: 'Sergeant', icon: '🎖️', rank: '4', count: 4 },
	'5': { name: 'Lieutenant', icon: '🎖️', rank: '5', count: 4 },
	'6': { name: 'Captain', icon: '👮', rank: '6', count: 4 },
	'7': { name: 'Major', icon: '⭐', rank: '7', count: 3 },
	'8': { name: 'Colonel', icon: '⭐⭐', rank: '8', count: 2 },
	'9': { name: 'General', icon: '👑', rank: '9', count: 1 },
	'M': { name: 'Marshal', icon: '👑👑', rank: '10', count: 1 }
};

export const MAX_BOARD_SETUPS = 10;
export const BOARD_ROWS = 4;
export const BOARD_COLS = 10;
export const BOARD_CELLS = BOARD_ROWS * BOARD_COLS;
