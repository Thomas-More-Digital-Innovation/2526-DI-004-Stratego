import { describe, it, expect } from 'vitest';
import {
	encodeBoard,
	decodeBoard,
	encodeSetup,
	decodeSetup,
	validateSetup,
	createEmptyBoard,
	createEmptySetup,
	PieceID,
	BitOccupied,
	BitColor,
	BitMoved,
	MaskPieceType,
	ShiftPieceType,
	type BoardSquare
} from './boardBinary';

describe('Board Encoding', () => {
	it('should encode empty board', () => {
		const board = createEmptyBoard();
		const data = encodeBoard(board);

		expect(data.length).toBe(100);
		expect(data.every((byte) => byte === 0)).toBe(true);
	});

	it('should encode pieces', () => {
		const board = createEmptyBoard();
		board[0][0] = { occupied: true, pieceType: PieceID.Flag, color: 1 };
		board[0][1] = { occupied: true, pieceType: PieceID.Marshal, color: 2, moved: true };

		const data = encodeBoard(board);

		const cell0 = data[0];
		expect(cell0 & BitOccupied).toBeTruthy();
		expect((cell0 & MaskPieceType) >> ShiftPieceType).toBe(PieceID.Flag);
		expect(cell0 & BitColor).toBe(0);
		expect(cell0 & BitMoved).toBe(0);

		const cell1 = data[1];
		expect(cell1 & BitOccupied).toBeTruthy();
		expect((cell1 & MaskPieceType) >> ShiftPieceType).toBe(PieceID.Marshal);
		expect(cell1 & BitColor).toBeTruthy();
		expect(cell1 & BitMoved).toBeTruthy();
	});

	it('should decode board', () => {
		const data = new Uint8Array(100);
		data[0] = BitOccupied | (PieceID.Flag << ShiftPieceType);
		data[1] = BitOccupied | (PieceID.Marshal << ShiftPieceType) | BitColor | BitMoved;

		const board = decodeBoard(data);

		expect(board[0][0].occupied).toBe(true);
		expect(board[0][0].pieceType).toBe(PieceID.Flag);
		expect(board[0][0].color).toBe(1);
		expect(board[0][0].moved).toBeFalsy();

		expect(board[0][1].occupied).toBe(true);
		expect(board[0][1].pieceType).toBe(PieceID.Marshal);
		expect(board[0][1].color).toBe(2);
		expect(board[0][1].moved).toBe(true);
	});
});

describe('Setup Encoding', () => {
	it('should encode setup', () => {
		const rows = ['BBBBBB2222', '4444555566', '6677788399', 'M311110333'];
		const data = encodeSetup(rows, 1);

		expect(data.length).toBe(40);

		const cell0 = data[0];
		expect(cell0 & BitOccupied).toBeTruthy();
		expect((cell0 & MaskPieceType) >> ShiftPieceType).toBe(PieceID.Bomb);
		expect(cell0 & BitColor).toBe(0);
	});

	it('should decode setup', () => {
		const originalRows = ['BBBBBB2222', '4444555566', '6677788399', 'M311110333'];
		const data = encodeSetup(originalRows, 1);
		const { rows, playerID } = decodeSetup(data);

		expect(playerID).toBe(1);
		expect(rows).toEqual(originalRows);
	});
});

describe('Validation', () => {
	it('should validate correct setup', () => {
		const setup = ['BBBBBB2222', '2222224444', '5555666677', '78899M3333'];
		const result = validateSetup(setup);

		expect(result.valid).toBe(true);
	});

	it('should reject invalid setup', () => {
		const setup = ['BBBBBB2222', '2222224444', '5555666677', '00899M3333'];
		const result = validateSetup(setup);

		expect(result.valid).toBe(false);
		expect(result.error).toBeDefined();
	});
});

describe('Helpers', () => {
	it('should create empty board', () => {
		const board = createEmptyBoard();

		expect(board.length).toBe(10);
		expect(board[0].length).toBe(10);
		expect(board[0][0].occupied).toBe(false);
	});

	it('should create empty setup', () => {
		const setup = createEmptySetup();

		expect(setup.length).toBe(4);
		expect(setup[0]).toBe('..........');
	});
});
