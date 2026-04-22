import type { Piece } from '$lib/types/game';

// Binary format: 8 bits per square
// Bit 7: Occupied | Bits 6-2: Piece type | Bit 1: Color | Bit 0: Moved
const BitOccupied = 1 << 7;
const BitColor = 1 << 1;
const MaskPieceType = 0x7C;
const ShiftPieceType = 2;

// Backend char -> piece ID (matches boardBinary.go rankToPieceID)
const CharToPieceID: Record<string, number> = {
    '0': 1, 'B': 2, '1': 3, '2': 4,
    '3': 5, '4': 6, '5': 7, '6': 8,
    '7': 9, '8': 10, '9': 11, 'M': 12,
};

const PieceIDToChar: Record<number, string> = Object.fromEntries(
    Object.entries(CharToPieceID).map(([k, v]) => [v, k])
);

/**
 * Encodes a 4x10 setup to a 40-character string of piece ranks.
 * Input rows use backend chars: 0=Flag, B=Bomb, 1=Spy, 2-9=Ranks, M=Marshal.
 */
export function encodeSetup(rows: string[]): string {
    // Ensure we have 4 rows of 10 chars, join into 40 chars
    let combined = '';
    for (let i = 0; i < 4; i++) {
        const row = rows[i] || '.'.repeat(10);
        combined += row.padEnd(10, '.').slice(0, 10);
    }
    return combined;
}

/**
 * Decodes a 40-character rank string to a 4x10 setup.
 * Returns rows using backend chars.
 */
export function decodeSetup(encoded: string): string[] {
    if (encoded.length !== 40) {
        // Fallback or handle error if needed, but we expect 40
        return Array(4).fill('.'.repeat(10));
    }
    const rows: string[] = [];
    for (let i = 0; i < 4; i++) {
        rows.push(encoded.slice(i * 10, (i + 1) * 10));
    }
    return rows;
}

/**
 * Flips a setup.
 */
export function flipSetup(encoded: string): string {
    const rows = decodeSetup(encoded);
    const flippedRows = rows.map(row => row.split('').reverse().join('')).reverse();
    return encodeSetup(flippedRows);
}

/**
 * Decodes a base64 10x10 board to a 2D array of Piece objects.
 */
export function decodeBoard(encoded: string): (Piece | null)[][] {
    const binary = atob(encoded);
    const data = new Uint8Array(binary.length);
    for (let i = 0; i < binary.length; i++) {
        data[i] = binary.charCodeAt(i);
    }

    const board: (Piece | null)[][] = [];
    for (let y = 0; y < 10; y++) {
        const row: (Piece | null)[] = [];
        for (let x = 0; x < 10; x++) {
            const cell = data[y * 10 + x];
            if ((cell & BitOccupied) === 0) {
                row.push(null);
                continue;
            }

            const pieceID = (cell & MaskPieceType) >> ShiftPieceType;
            const rank = PieceIDToChar[pieceID];
            const ownerId = (cell & BitColor) ? 2 : 1;

            row.push({
                rank,
                ownerId,
                revealed: true,
                position: { x, y }
            });
        }
        board.push(row);
    }
    return board;
}
