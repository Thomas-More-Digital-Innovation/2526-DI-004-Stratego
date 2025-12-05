<script lang="ts">
	interface PieceInfo {
		name: string;
		icon: string;
		rank: string;
		count: number;
	}

	interface Props {
		pieceInventory: Record<string, PieceInfo>;
		remaining: Record<string, number>;
		selectedPiece: string | null;
		onSelect: (piece: string) => void;
	}

	let { pieceInventory, remaining, selectedPiece, onSelect }: Props = $props();

	function getTotalPlaced(): number {
		let total = 0;
		Object.entries(pieceInventory).forEach(([char, info]) => {
			const placed = info.count - (remaining[char] || 0);
			total += placed;
		});
		return total;
	}

	function isComplete(): boolean {
		return Object.values(remaining).every((count) => count === 0);
	}
</script>

<div class="piece-bank">
	<div class="bank-header">
		<h4>Piece Bank</h4>
		<div class="status">
			<span class="placed">{getTotalPlaced()}/40</span>
			{#if isComplete()}
				<span class="complete">✓ Complete</span>
			{/if}
		</div>
	</div>

	<div class="pieces">
		{#each Object.entries(pieceInventory) as [char, info]}
			{@const remainingCount = remaining[char] || 0}
			{@const isSelected = selectedPiece === char}
			{@const isAvailable = remainingCount > 0}

			<button
				class="piece-item"
				class:selected={isSelected}
				class:unavailable={!isAvailable}
				onclick={() => isAvailable && onSelect(char)}
				disabled={!isAvailable}
			>
				<div class="piece-display">
					<span class="icon">{info.icon}</span>
					<span class="rank">{info.rank}</span>
				</div>
				<div class="details">
					<span class="name">{info.name}</span>
					<span class="count">{remainingCount}/{info.count}</span>
				</div>
			</button>
		{/each}
	</div>
</div>

<style>
	.piece-bank {
		display: flex;
		flex-direction: column;
		background: var(--bg);
		padding: 15px;
		border-radius: 8px;
		height: 100%;
		box-sizing: border-box;
		min-height: 0; /* allow children to shrink inside flex parent */
	}

	.bank-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 12px;

		h4 {
			margin: 0;
			color: var(--text);
			font-size: 1rem;
		}

		.status {
			display: flex;
			gap: 8px;
			align-items: center;

			.placed {
				color: var(--muted, #999);
				font-size: 0.85rem;
			}

			.complete {
				background: #4a4;
				color: white;
				padding: 2px 8px;
				border-radius: 10px;
				font-size: 0.75rem;
				font-weight: 600;
			}
		}
	}

	.pieces {
		display: flex;
		flex-direction: column;
		gap: 6px;
		/* make pieces take remaining space and have a definite height for scrolling */
		height: calc(100% - 56px); /* 56px ~= header height + gaps; adjust if needed */
		overflow-y: auto;
		-webkit-overflow-scrolling: touch;
		min-height: 0;
	}

	.piece-item {
		display: flex;
		align-items: center;
		gap: 10px;
		padding: 10px;
		background: #2a2a2a;
		border: 2px solid transparent;
		border-radius: 6px;
		cursor: pointer;
		transition: all 0.2s;

		&:hover:not(.unavailable) {
			background: #3a3a3a;
			border-color: #66f;
		}

		&.selected {
			background: #3a3a5a;
			border-color: #66f;
		}

		&.unavailable {
			opacity: 0.4;
			cursor: not-allowed;
		}

		.piece-display {
			display: flex;
			flex-direction: column;
			align-items: center;
			justify-content: center;
			gap: 1px;
			min-width: 45px;
			background: #ff6b6b;
			padding: 6px;
			border-radius: 6px;
			box-sizing: border-box;

			.icon {
				font-size: 1.3rem;
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

		.details {
			flex: 1;
			display: flex;
			justify-content: space-between;
			align-items: center;

			.name {
				color: var(--text);
				font-weight: 500;
				font-size: 0.9rem;
			}

			.count {
				color: var(--muted, #999);
				font-size: 0.85rem;
				font-weight: 600;
			}
		}
	}

	/* Scrollbar styling */
	.pieces::-webkit-scrollbar {
		width: 6px;
	}

	.pieces::-webkit-scrollbar-track {
		background: #1a1a1a;
		border-radius: 3px;
	}

	.pieces::-webkit-scrollbar-thumb {
		background: #444;
		border-radius: 3px;

		&:hover {
			background: #555;
		}
	}
</style>
