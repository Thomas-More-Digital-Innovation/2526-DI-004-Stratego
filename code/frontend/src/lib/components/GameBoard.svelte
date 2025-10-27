<script lang="ts">
	import type { GameState, Position } from '$lib/types';
	import { getValidMoves, makeMove } from '$lib/gameState';
	import BoardCell from './BoardCell.svelte';

	interface Props {
		state: GameState;
		onStateChange: (state: GameState) => void;
	}

	let { state, onStateChange }: Props = $props();

	function handleCellClick(x: number, y: number) {
		if (state.phase !== 'playing') return;

		const clickedPiece = state.board[y][x];
		const currentPlayer = state.players[state.currentPlayerIndex];

		// If a piece is selected
		if (state.selectedPiece) {
			// Check if this is a valid move
			const isValidMove = state.validMoves.some(move => move.x === x && move.y === y);

			if (isValidMove) {
				// Make the move
				const newState = makeMove(
					state,
					state.selectedPiece.position,
					{ x, y }
				);
				onStateChange(newState);
			} else if (clickedPiece && clickedPiece.owner.id === currentPlayer.id) {
				// Select a different piece
				state.selectedPiece = clickedPiece;
				state.validMoves = getValidMoves(clickedPiece, state.board);
			} else {
				// Deselect
				state.selectedPiece = null;
				state.validMoves = [];
			}
		} else {
			// No piece selected, try to select one
			if (clickedPiece && clickedPiece.owner.id === currentPlayer.id && clickedPiece.alive) {
				state.selectedPiece = clickedPiece;
				state.validMoves = getValidMoves(clickedPiece, state.board);
			}
		}
	}

	function isValidMove(x: number, y: number): boolean {
		return state.validMoves.some(move => move.x === x && move.y === y);
	}

	function isSelected(x: number, y: number): boolean {
		return state.selectedPiece?.position.x === x && state.selectedPiece?.position.y === y;
	}
</script>

<div class="board-container">
	<div class="board">
		{#each Array(10) as _, y}
			{#each Array(10) as _, x}
				<BoardCell
					piece={state.board[y][x]}
					position={{ x, y }}
					isSelected={isSelected(x, y)}
					isValidMove={isValidMove(x, y)}
					currentPlayer={state.currentPlayerIndex}
					onClick={() => handleCellClick(x, y)}
				/>
			{/each}
		{/each}
	</div>
</div>

<style>
	.board-container {
		display: flex;
		justify-content: center;
		align-items: center;
		padding: 1rem;
	}

	.board {
		display: grid;
		grid-template-columns: repeat(10, 1fr);
		grid-template-rows: repeat(10, 1fr);
		gap: 2px;
		max-width: 600px;
		width: 100%;
		aspect-ratio: 1;
		background: #5a4a3a;
		padding: 8px;
		border-radius: 12px;
		box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
	}

	@media (max-width: 768px) {
		.board {
			max-width: 100%;
			padding: 4px;
			gap: 1px;
		}
	}
</style>
