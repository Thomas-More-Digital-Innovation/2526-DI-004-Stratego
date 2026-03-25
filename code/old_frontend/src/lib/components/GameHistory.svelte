<script lang="ts">
	interface Props {
		currentMoveIndex: number;
		totalMoves: number;
		isReplaying: boolean;
		onPrevious: () => void;
		onNext: () => void;
		onGoToMove: (index: number) => void;
		onExitReplay: () => void;
	}

	let { 
		currentMoveIndex,
		totalMoves,
		isReplaying,
		onPrevious,
		onNext,
		onGoToMove,
		onExitReplay
	}: Props = $props();

	const canGoPrevious = $derived(currentMoveIndex > 0);
	const canGoNext = $derived(currentMoveIndex < totalMoves - 1);
</script>

<div class="history">
	<div class="history-header">
		<h3>Move History</h3>
		{#if isReplaying}
			<span class="replay-badge">REPLAY MODE</span>
		{/if}
	</div>

	{#if totalMoves > 0}
		<div class="history-info">
			<p>Move {currentMoveIndex + 1} of {totalMoves}</p>
		</div>

		<div class="history-controls">
			<button 
				class="btn btn-nav"
				onclick={onPrevious}
				disabled={!canGoPrevious}
				title="Previous Move"
			>
				◀ Previous
			</button>

			<button 
				class="btn btn-nav"
				onclick={onNext}
				disabled={!canGoNext}
				title="Next Move"
			>
				Next ▶
			</button>
		</div>

		<div class="move-list">
			{#each Array(totalMoves) as _, index}
				<button
					class="move-item"
					class:active={index === currentMoveIndex}
					onclick={() => onGoToMove(index)}
				>
					Move {index + 1}
				</button>
			{/each}
		</div>

		{#if isReplaying}
			<button 
				class="btn btn-exit"
				onclick={onExitReplay}
			>
				Exit Replay Mode
			</button>
		{/if}
	{:else}
		<p class="no-history">No moves recorded yet</p>
	{/if}
</div>

<style>
	.history {
		background: #2d3748;
		padding: 20px;
		border-radius: 8px;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
		max-height: 600px;
		display: flex;
		flex-direction: column;
	}

	.history-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 15px;
		border-bottom: 2px solid #4a5568;
		padding-bottom: 10px;
	}

	h3 {
		margin: 0;
		color: #e2e8f0;
		font-size: 1.2rem;
	}

	.replay-badge {
		background: #f6ad55;
		color: #1a202c;
		padding: 4px 12px;
		border-radius: 12px;
		font-size: 0.8rem;
		font-weight: 600;
	}

	.history-info {
		text-align: center;
		color: #a0aec0;
		margin-bottom: 15px;
	}

	.history-info p {
		margin: 0;
	}

	.history-controls {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 10px;
		margin-bottom: 15px;
	}

	.move-list {
		flex: 1;
		overflow-y: auto;
		display: flex;
		flex-direction: column;
		gap: 5px;
		margin-bottom: 15px;
	}

	.move-item {
		padding: 10px;
		background: #4a5568;
		color: #e2e8f0;
		border: none;
		border-radius: 4px;
		cursor: pointer;
		transition: all 0.2s;
		text-align: left;
	}

	.move-item:hover {
		background: #718096;
	}

	.move-item.active {
		background: #4299e1;
		color: white;
		font-weight: 600;
	}

	.btn {
		padding: 10px 15px;
		border: none;
		border-radius: 6px;
		font-size: 0.95rem;
		font-weight: 600;
		cursor: pointer;
		transition: all 0.2s;
	}

	.btn:disabled {
		opacity: 0.3;
		cursor: not-allowed;
	}

	.btn-nav {
		background: #4a5568;
		color: #e2e8f0;
	}

	.btn-nav:hover:not(:disabled) {
		background: #718096;
	}

	.btn-exit {
		background: #fc8181;
		color: white;
		width: 100%;
	}

	.btn-exit:hover {
		background: #f56565;
	}

	.no-history {
		color: #718096;
		text-align: center;
		padding: 20px 0;
	}
</style>
