<script lang="ts">
	import type { Piece } from '../types';

	interface Props {
		piece: Piece | null;
		isSelected?: boolean;
		isHighlighted?: boolean;
		isLake?: boolean;
		viewerId?: number;
	}

	let { 
		piece, 
		isSelected = false, 
		isHighlighted = false,
		isLake = false,
		viewerId = 0  // Default to player 0 (human player)
	}: Props = $props();

	const getPlayerColor = (ownerId: number): string => {
		return ownerId === 0 ? '#ff6b6b' : '#4dabf7';
	};

	const canSeePiece = $derived(() => {
		if (!piece || !piece.ownerName) return false; // Empty cell (no ownerName)
		// Can see own pieces or revealed enemy pieces
		return piece.ownerId === viewerId || piece.revealed;
	});

	const cellStyle = $derived(() => {
		if (isLake) return 'background: #2b6cb0;'; // Blue for lakes
		if (!piece || !piece.ownerName) return 'background: #2d3748;'; // Dark gray for empty
		const color = getPlayerColor(piece.ownerId);
		return `background: ${color};`;
	});
</script>

<div 
	class="piece"
	class:selected={isSelected}
	class:highlighted={isHighlighted}
	class:empty={!piece || !piece.ownerName}
	class:lake={isLake}
	style={cellStyle()}
>
	{#if isLake}
		<span class="lake-icon">🌊</span>
	{:else if piece && piece.ownerName}
		<!-- Show piece details if we can see it, otherwise show ? -->
		{#if canSeePiece()}
			{#if piece.icon}
				<span class="icon">{piece.icon}</span>
			{:else if piece.rank}
				<span class="rank">{piece.rank}</span>
			{/if}
		{:else}
			<!-- Enemy piece we can't see yet -->
			<span class="hidden">?</span>
		{/if}
	{/if}
</div>

<style>
	.piece {
		width: 100%;
		height: 100%;
		display: flex;
		align-items: center;
		justify-content: center;
		border: 1px solid #1a202c;
		border-radius: 4px;
		cursor: pointer;
		transition: all 0.2s;
		font-size: 1.5rem;
	}

	.piece.empty {
		cursor: default;
		background: #2d3748;
	}

	.piece.lake {
		cursor: default;
		background: #2b6cb0;
	}

	.piece.selected {
		border: 3px solid #ffd700;
		box-shadow: 0 0 10px #ffd700;
	}

	.piece.highlighted {
		border: 2px solid #90ee90;
	}

	.piece:not(.empty):not(.lake):hover {
		transform: scale(1.05);
		box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
	}

	.icon {
		font-size: 2rem;
	}

	.lake-icon {
		font-size: 2rem;
	}

	.hidden {
		font-size: 2rem;
		color: white;
		font-weight: bold;
	}

	.rank {
		font-size: 1.2rem;
		font-weight: bold;
		color: white;
	}
</style>
