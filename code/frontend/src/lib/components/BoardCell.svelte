<script lang="ts">
	import type { Piece, Position } from '$lib/types';
	import { isLake } from '$lib/gameState';

	interface Props {
		piece: Piece | null;
		position: Position;
		isSelected?: boolean;
		isValidMove?: boolean;
		currentPlayer: number;
		onClick: () => void;
	}

	let { piece, position, isSelected = false, isValidMove = false, currentPlayer, onClick }: Props = $props();

	let lake = $derived(isLake(position));
	let isOwnPiece = $derived(piece && piece.owner.id === currentPlayer);
	let isEnemyPiece = $derived(piece && piece.owner.id !== currentPlayer);
	let canInteract = $derived((isOwnPiece || isValidMove) && !lake);
</script>

<button
	class="cell"
	class:lake
	class:selected={isSelected}
	class:valid-move={isValidMove}
	class:own-piece={isOwnPiece}
	class:enemy-piece={isEnemyPiece}
	class:interactive={canInteract}
	disabled={lake}
	onclick={onClick}
	type="button"
>
	{#if lake}
		<span class="lake-icon">🌊</span>
	{:else if piece}
		<div class="piece" class:red={piece.owner.color === 'red'} class:blue={piece.owner.color === 'blue'}>
			{#if piece.revealed || piece.owner.id === currentPlayer}
				<span class="piece-icon">{piece.type.icon}</span>
				<span class="piece-rank">{piece.type.rank}</span>
			{:else}
				<span class="piece-hidden">?</span>
			{/if}
		</div>
	{:else if isValidMove}
		<span class="move-indicator">•</span>
	{/if}
</button>

<style>
	.cell {
		aspect-ratio: 1;
		background: linear-gradient(135deg, #f0e6d2 0%, #e8dcc8 100%);
		border: 2px solid #8b7355;
		display: flex;
		align-items: center;
		justify-content: center;
		cursor: pointer;
		transition: all 0.2s ease;
		position: relative;
		font-size: 0.75rem;
	}

	.cell:hover:not(:disabled) {
		transform: scale(1.05);
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
		z-index: 10;
	}

	.cell.lake {
		background: linear-gradient(135deg, #4a90e2 0%, #357abd 100%);
		cursor: not-allowed;
		border-color: #2c5a8a;
	}

	.lake-icon {
		font-size: 1.5rem;
		filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.2));
	}

	.cell.selected {
		border-color: #ffd700;
		border-width: 3px;
		box-shadow: 0 0 20px rgba(255, 215, 0, 0.6);
	}

	.cell.valid-move {
		background: linear-gradient(135deg, #90ee90 0%, #7ed87e 100%);
		border-color: #4caf50;
		animation: pulse 1s ease-in-out infinite;
	}

	@keyframes pulse {
		0%, 100% {
			opacity: 1;
		}
		50% {
			opacity: 0.7;
		}
	}

	.move-indicator {
		font-size: 2rem;
		color: #4caf50;
		font-weight: bold;
	}

	.piece {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		width: 100%;
		height: 100%;
		padding: 0.25rem;
		border-radius: 8px;
		background: rgba(255, 255, 255, 0.9);
		box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
	}

	.piece.red {
		background: linear-gradient(135deg, #ff6b6b 0%, #ee5a52 100%);
		color: white;
	}

	.piece.blue {
		background: linear-gradient(135deg, #4a90e2 0%, #357abd 100%);
		color: white;
	}

	.piece-icon {
		font-size: 1.5rem;
		line-height: 1;
		filter: drop-shadow(0 1px 2px rgba(0, 0, 0, 0.3));
	}

	.piece-rank {
		font-size: 0.7rem;
		font-weight: bold;
		margin-top: 0.1rem;
		text-shadow: 0 1px 2px rgba(0, 0, 0, 0.3);
	}

	.piece-hidden {
		font-size: 2rem;
		font-weight: bold;
		text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
	}

	.own-piece .piece {
		cursor: pointer;
	}

	.enemy-piece .piece {
		cursor: not-allowed;
	}

	.cell.interactive {
		cursor: pointer;
	}

	@media (max-width: 768px) {
		.piece-icon {
			font-size: 1.2rem;
		}
		.piece-rank {
			font-size: 0.6rem;
		}
	}
</style>
