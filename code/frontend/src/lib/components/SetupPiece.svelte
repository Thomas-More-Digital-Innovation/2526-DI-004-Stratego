<script lang="ts">
	import type { PieceInfo } from '$lib/types/boardSetup.type';

	interface Props {
		pieceChar: string;
		pieceInfo?: PieceInfo;
		isSelected?: boolean;
		isEmpty?: boolean;
		canPlace?: boolean;
	}

	let { pieceChar, pieceInfo, isSelected = false, isEmpty = false, canPlace = false }: Props = $props();
</script>

<div
	class="setup-piece"
	class:selected={isSelected}
	class:empty={isEmpty}
	class:can-place={canPlace}
>
	{#if !isEmpty && pieceInfo}
		<div class="piece-content">
			<span class="icon">{pieceInfo.icon}</span>
			<span class="rank">{pieceInfo.rank}</span>
		</div>
	{/if}
</div>

<style>
	.setup-piece {
		width: 100%;
		height: 100%;
		display: flex;
		align-items: center;
		justify-content: center;
		border: 2px solid #1a202c;
		border-radius: 6px;
		cursor: pointer;
		transition: all 0.2s;
		background: #ff6b6b;
		position: relative;
		box-sizing: border-box;

		&.empty {
			background: #2d3748;
			cursor: default;
		}

		&.can-place {
			background: #2a4a2a;
			border-color: #4a4;
		}

		&.selected {
			border-color: #ffd700;
			box-shadow: 0 0 8px #ffd700;
		}

		&:not(.empty):hover {
			box-shadow: 0 2px 6px rgba(0, 0, 0, 0.4);
		}
	}

	.piece-content {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 1px;
		padding: 2px;

		.icon {
			font-size: 1.4rem;
			line-height: 1;
		}

		.rank {
			font-size: 0.7rem;
			font-weight: bold;
			color: white;
			background: rgba(0, 0, 0, 0.3);
			padding: 1px 4px;
			border-radius: 3px;
			line-height: 1;
		}
	}
</style>
