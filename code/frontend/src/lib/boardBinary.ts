/**
 * Binary format for Stratego board encoding
 * 8 bits per square: Occupied | Piece type | Color | Moved
 */

export const BitOccupied = 0x80;
export const BitColor = 0x02;
export const BitMoved = 0x01;
export const MaskPieceType = 0x7C;
export const ShiftPieceType = 2;

export enum PieceID {
	Flag = 1,
	Bomb = 2,
	Spy = 3,
	Scout = 4,
	Miner = 5,
	Sergeant = 6,
	Lieutenant = 7,
	Captain = 8,
	Major = 9,
	Colonel = 10,
	General = 11,
	Marshal = 12
}

const rankToPieceID: Record<string, PieceID> = {
	'0': PieceID.Flag, 'B': PieceID.Bomb, '1': PieceID.Spy, '2': PieceID.Scout,
	'3': PieceID.Miner, '4': PieceID.Sergeant, '5': PieceID.Lieutenant, '6': PieceID.Captain,
	'7': PieceID.Major, '8': PieceID.Colonel, '9': PieceID.General, 'M': PieceID.Marshal
};

const pieceIDToRank: Record<number, string> = Object.fromEntries(
	Object.entries(rankToPieceID).map(([k, v]) => [v, k])
);

export interface BoardSquare {
	occupied: boolean;
	pieceType?: PieceID;
	color?: number;
	moved?: boolean;
}

// Encode 10x10 board to 100 bytes
export function encodeBoard(board: BoardSquare[][]): Uint8Array {
	const data = new Uint8Array(100);
	
	for (let y = 0; y < 10; y++) {
		for (let x = 0; x < 10; x++) {
			const square = board[y][x];
			if (!square.occupied || !square.pieceType) continue;
			
			let cell = BitOccupied;
			cell |= (square.pieceType << ShiftPieceType) & MaskPieceType;
			
			if (square.color === 2) {
				cell |= BitColor;
			}
			
			if (square.moved) {
				cell |= BitMoved;
			}
			
			data[y * 10 + x] = cell;
		}
	}
	
	return data;
}

// Decode 100 bytes to 10x10 board
export function decodeBoard(data: Uint8Array): BoardSquare[][] {
	const board: BoardSquare[][] = [];
	
	for (let y = 0; y < 10; y++) {
		const row: BoardSquare[] = [];
		for (let x = 0; x < 10; x++) {
			const cell = data[y * 10 + x];
			
			if (cell & BitOccupied) {
				const pieceType = (cell & MaskPieceType) >> ShiftPieceType;
				row.push({
					occupied: true,
					pieceType,
					color: (cell & BitColor) ? 2 : 1,
					moved: !!(cell & BitMoved)
				});
			} else {
				row.push({ occupied: false });
			}
		}
		board.push(row);
	}
	
	return board;
}

// Encode 4x10 setup to 40 bytes
export function encodeSetup(rows: string[], playerID: number): Uint8Array {
	if (rows.length !== 4) throw new Error(`Expected 4 rows, got ${rows.length}`);
	
	const data = new Uint8Array(40);
	
	for (let i = 0; i < 4; i++) {
		if (rows[i].length !== 10) throw new Error(`Row ${i} must be 10 chars`);
		
		for (let j = 0; j < 10; j++) {
			const char = rows[i][j];
			if (char === '.' || char === ' ') continue;
			
			let cell = BitOccupied;
			cell |= (rankToPieceID[char] << ShiftPieceType) & MaskPieceType;
			
			if (playerID === 2) {
				cell |= BitColor;
			}
			
			data[i * 10 + j] = cell;
		}
	}
	
	return data;
}

// Decode 40 bytes to 4x10 setup
export function decodeSetup(data: Uint8Array): { rows: string[]; playerID: number } {
	if (data.length !== 40) throw new Error(`Invalid data length: expected 40, got ${data.length}`);
	
	const rows: string[] = [];
	let playerID = 1;
	
	for (let i = 0; i < 4; i++) {
		let row = '';
		for (let j = 0; j < 10; j++) {
			const cell = data[i * 10 + j];
			
			if (cell & BitOccupied) {
				const pieceID = (cell & MaskPieceType) >> ShiftPieceType;
				row += pieceIDToRank[pieceID];
				
				if (cell & BitColor) {
					playerID = 2;
				}
			} else {
				row += '.';
			}
		}
		rows.push(row);
	}
	
	return { rows, playerID };
}

// Validate setup piece counts
export function validateSetup(rows: string[]): { valid: boolean; error?: string } {
	if (rows.length !== 4) {
		return { valid: false, error: 'Setup must have 4 rows' };
	}
	
	const counts: Record<string, number> = {};
	for (const row of rows) {
		if (row.length !== 10) {
			return { valid: false, error: 'Each row must be 10 chars' };
		}
		for (const c of row) {
			if (c !== '.' && c !== ' ') {
				counts[c] = (counts[c] || 0) + 1;
			}
		}
	}
	
	const expected: Record<string, number> = {
		'0': 1, 'B': 6, '1': 1, '2': 8, '3': 5, '4': 4,
		'5': 4, '6': 4, '7': 3, '8': 2, '9': 1, 'M': 1
	};
	
	for (const [piece, count] of Object.entries(expected)) {
		if (counts[piece] !== count) {
			return {
				valid: false,
				error: `Piece ${piece}: expected ${count}, got ${counts[piece] || 0}`
			};
		}
	}
	
	return { valid: true };
}

// Helper: Create empty board
export function createEmptyBoard(): BoardSquare[][] {
	return Array.from({ length: 10 }, () =>
		Array.from({ length: 10 }, () => ({ occupied: false }))
	);
}

// Helper: Create empty setup
export function createEmptySetup(): string[] {
	return Array(4).fill('.'.repeat(10));
}
