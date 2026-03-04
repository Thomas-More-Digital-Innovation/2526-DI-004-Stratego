<script lang="ts">
	import Piece from './Piece.svelte';
	import type { BoardState, Position } from '../types';

	interface Props {
		boardState: BoardState | null;
		selectedPosition: Position | null;
		onCellClick: (x: number, y: number) => void;
		isInteractive?: boolean;
		viewerId?: number;
		validMoves?: Position[];
	}

	let { 
		boardState, 
		selectedPosition, 
		onCellClick,
		isInteractive = true,
		viewerId = 0,
		validMoves = []
	}: Props = $props();

	// Lake positions in Stratego (2x2 squares at specific locations)
	const lakePositions = [
		{ x: 2, y: 4 }, { x: 3, y: 4 }, { x: 2, y: 5 }, { x: 3, y: 5 },
		{ x: 6, y: 4 }, { x: 7, y: 4 }, { x: 6, y: 5 }, { x: 7, y: 5 }
	];

	const isLake = (x: number, y: number): boolean => {
		return lakePositions.some(pos => pos.x === x && pos.y === y);
	};

	const isSelected = (x: number, y: number): boolean => {
		return selectedPosition?.x === x && selectedPosition?.y === y;
	};

	const isValidMove = (x: number, y: number): boolean => {
		return validMoves.some(move => move.x === x && move.y === y);
	};
</script>

<div class="board-container">
	{#if boardState}
		<div class="board">
			{#each boardState.board as row, y}
				{#each row as piece, x}
					<div 
						class="cell"
						class:interactive={isInteractive && !isLake(x, y)}
						class:valid-move={isValidMove(x, y)}
						onclick={() => isInteractive && !isLake(x, y) && onCellClick(x, y)}
					>
						<Piece 
							piece={piece} 
							isSelected={isSelected(x, y)}
							isHighlighted={isValidMove(x, y)}
							isLake={isLake(x, y)}
							{viewerId}
						/>
					</div>
				{/each}
			{/each}
		</div>
	{:else}
		<div class="loading">Waiting for board state...</div>
	{/if}
</div>

<style>
	.board-container {
		display: flex;
		justify-content: center;
		align-items: center;
		padding: 20px;
	}

	.board {
		display: grid;
		grid-template-columns: repeat(10, 50px);
		grid-template-rows: repeat(10, 50px);
		gap: 2px;
		background: #1a202c;
		padding: 10px;
		border-radius: 8px;
		box-shadow: 0 4px 6px rgba(0, 0, 0, 0.3);
	}

	.cell {
		width: 50px;
		height: 50px;
		position: relative;
	}

	.cell.interactive {
		cursor: pointer;
	}

	.cell.valid-move::after {
		content: '';
		position: absolute;
		inset: 4px;
		border: 3px solid #90ee90;
		border-radius: 4px;
		pointer-events: none;
		animation: pulse 1.5s ease-in-out infinite;
	}

	@keyframes pulse {
		0%, 100% {
			opacity: 1;
			transform: scale(1);
		}
		50% {
			opacity: 0.6;
			transform: scale(0.95);
		}
	}

	.loading {
		padding: 40px;
		text-align: center;
		color: #718096;
		font-size: 1.2rem;
	}

	@media (max-width: 768px) {
		.board {
			grid-template-columns: repeat(10, 35px);
			grid-template-rows: repeat(10, 35px);
		}

		.cell {
			width: 35px;
			height: 35px;
		}
	}
</style>
