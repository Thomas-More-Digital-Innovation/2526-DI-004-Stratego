// Game state management using Svelte 5 runes
import type { GameState, Player, Piece, Position, Move, PieceType } from './types';
import { PIECE_TYPES } from './types';

const BOARD_SIZE = 10;

// Lake positions (water squares that can't be moved through)
const LAKE_POSITIONS = [
	{ x: 2, y: 4 }, { x: 3, y: 4 },
	{ x: 2, y: 5 }, { x: 3, y: 5 },
	{ x: 6, y: 4 }, { x: 7, y: 4 },
	{ x: 6, y: 5 }, { x: 7, y: 5 }
];

export function createGameState(): GameState {
	const board: (Piece | null)[][] = Array(BOARD_SIZE).fill(null).map(() => Array(BOARD_SIZE).fill(null));
	
	// Add lakes to the board
	for (const lake of LAKE_POSITIONS) {
		// We'll handle lakes separately in the UI
	}

	const player1: Player = {
		id: 0,
		name: 'Red Player',
		color: 'red',
		pieces: []
	};

	const player2: Player = {
		id: 1,
		name: 'Blue Player',
		color: 'blue',
		pieces: []
	};

	// Initialize pieces for both players in setup positions
	initializePlayerPieces(player1, board, 0);
	initializePlayerPieces(player2, board, 6);

	return {
		phase: 'setup',
		players: [player1, player2],
		currentPlayerIndex: 0,
		board,
		selectedPiece: null,
		validMoves: [],
		winner: null,
		winCause: null,
		round: 1,
		moveHistory: []
	};
}

function initializePlayerPieces(player: Player, board: (Piece | null)[][], startRow: number) {
	// Create a default setup with all pieces
	const piecesToPlace: PieceType[] = [];
	
	Object.values(PIECE_TYPES).forEach(type => {
		for (let i = 0; i < type.count; i++) {
			piecesToPlace.push(type);
		}
	});

	// Shuffle pieces for random placement
	piecesToPlace.sort(() => Math.random() - 0.5);

	let pieceIndex = 0;
	for (let y = startRow; y < startRow + 4; y++) {
		for (let x = 0; x < BOARD_SIZE; x++) {
			if (pieceIndex < piecesToPlace.length) {
				const type = piecesToPlace[pieceIndex];
				const piece: Piece = {
					type,
					owner: player,
					alive: true,
					revealed: false,
					position: { x, y }
				};
				board[y][x] = piece;
				player.pieces.push(piece);
				pieceIndex++;
			}
		}
	}
}

export function isLake(pos: Position): boolean {
	return LAKE_POSITIONS.some(lake => lake.x === pos.x && lake.y === pos.y);
}

export function getValidMoves(piece: Piece, board: (Piece | null)[][]): Position[] {
	if (!piece.type.movable) {
		return [];
	}

	const moves: Position[] = [];
	const { x, y } = piece.position;

	// Check all four directions
	const directions = [
		{ dx: 0, dy: -1 }, // up
		{ dx: 0, dy: 1 },  // down
		{ dx: -1, dy: 0 }, // left
		{ dx: 1, dy: 0 }   // right
	];

	for (const { dx, dy } of directions) {
		// Scouts can move multiple spaces
		const maxDistance = piece.type.rank === '2' ? BOARD_SIZE : 1;

		for (let distance = 1; distance <= maxDistance; distance++) {
			const newX = x + dx * distance;
			const newY = y + dy * distance;
			const newPos = { x: newX, y: newY };

			// Check bounds
			if (newX < 0 || newX >= BOARD_SIZE || newY < 0 || newY >= BOARD_SIZE) {
				break;
			}

			// Check if it's a lake
			if (isLake(newPos)) {
				break;
			}

			const targetPiece = board[newY][newX];

			// Can't move through own pieces
			if (targetPiece && targetPiece.owner.id === piece.owner.id) {
				break;
			}

			// Can move to empty square or attack enemy
			moves.push(newPos);

			// Can't move past an enemy piece
			if (targetPiece && targetPiece.owner.id !== piece.owner.id) {
				break;
			}
		}
	}

	return moves;
}

export function resolveCombat(attacker: Piece, defender: Piece): { winner: Piece | null; loser: Piece | null } {
	// Special rules
	// Spy defeats Marshal
	if (attacker.type.rank === '1' && defender.type.rank === 'M') {
		defender.alive = false;
		attacker.revealed = true;
		defender.revealed = true;
		return { winner: attacker, loser: defender };
	}

	// Miner defeats Bomb
	if (attacker.type.rank === '3' && defender.type.rank === 'B') {
		defender.alive = false;
		attacker.revealed = true;
		defender.revealed = true;
		return { winner: attacker, loser: defender };
	}

	// Bomb defeats all except Miner
	if (defender.type.rank === 'B') {
		attacker.alive = false;
		attacker.revealed = true;
		defender.revealed = true;
		return { winner: null, loser: attacker };
	}

	// Flag is captured
	if (defender.type.rank === '0') {
		defender.alive = false;
		attacker.revealed = true;
		defender.revealed = true;
		return { winner: attacker, loser: defender };
	}

	// Normal combat - compare strategic values
	const attackerValue = attacker.type.strategicValue;
	const defenderValue = defender.type.strategicValue;

	attacker.revealed = true;
	defender.revealed = true;

	if (attackerValue > defenderValue) {
		defender.alive = false;
		return { winner: attacker, loser: defender };
	} else if (defenderValue > attackerValue) {
		attacker.alive = false;
		return { winner: defender, loser: attacker };
	} else {
		// Tie - both die
		attacker.alive = false;
		defender.alive = false;
		return { winner: null, loser: null };
	}
}

export function makeMove(state: GameState, from: Position, to: Position): GameState {
	const piece = state.board[from.y][from.x];
	if (!piece) return state;

	const targetPiece = state.board[to.y][to.x];

	if (targetPiece) {
		// Combat!
		const result = resolveCombat(piece, targetPiece);
		
		// Update board based on combat result
		if (result.winner === piece) {
			// Attacker wins, move to new position
			state.board[from.y][from.x] = null;
			state.board[to.y][to.x] = piece;
			piece.position = to;
			
			// Remove defender from their pieces
			const defenderOwner = targetPiece.owner;
			defenderOwner.pieces = defenderOwner.pieces.filter(p => p !== targetPiece);
		} else if (result.winner === targetPiece) {
			// Defender wins, attacker removed
			state.board[from.y][from.x] = null;
			
			// Remove attacker from their pieces
			piece.owner.pieces = piece.owner.pieces.filter(p => p !== piece);
		} else {
			// Both die
			state.board[from.y][from.x] = null;
			state.board[to.y][to.x] = null;
			
			piece.owner.pieces = piece.owner.pieces.filter(p => p !== piece);
			targetPiece.owner.pieces = targetPiece.owner.pieces.filter(p => p !== targetPiece);
		}

		// Check if flag was captured
		if (targetPiece.type.rank === '0') {
			state.winner = piece.owner;
			state.winCause = 'flag_captured';
			state.phase = 'gameOver';
		}
	} else {
		// Simple move
		state.board[from.y][from.x] = null;
		state.board[to.y][to.x] = piece;
		piece.position = to;
	}

	// Record move
	state.moveHistory.push({ from, to });

	// Switch players
	state.currentPlayerIndex = state.currentPlayerIndex === 0 ? 1 : 0;
	if (state.currentPlayerIndex === 0) {
		state.round++;
	}

	// Clear selection
	state.selectedPiece = null;
	state.validMoves = [];

	// Check for game over conditions
	checkGameOver(state);

	return state;
}

function checkGameOver(state: GameState) {
	// Check if either player has no movable pieces
	for (const player of state.players) {
		const hasMovablePieces = player.pieces.some(piece => {
			if (!piece.alive || !piece.type.movable) return false;
			const moves = getValidMoves(piece, state.board);
			return moves.length > 0;
		});

		if (!hasMovablePieces) {
			const otherPlayer = state.players.find(p => p.id !== player.id);
			state.winner = otherPlayer || null;
			state.winCause = 'no_movable_pieces';
			state.phase = 'gameOver';
			return;
		}
	}

	// Check max turns (optional)
	if (state.round > 500) {
		state.winCause = 'max_turns';
		state.phase = 'gameOver';
		// Draw - no winner
	}
}

export function startGame(state: GameState): GameState {
	state.phase = 'playing';
	return state;
}
