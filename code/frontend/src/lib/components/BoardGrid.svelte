<script lang="ts">
	// TODO right mouse click menu not working
	import { PIECE_INVENTORY } from '$lib/types/boardSetup';
	import SetupPiece from './SetupPiece.svelte';

	interface Props {
		board: string[];
		selectedPiece: string | null;
		selectedPieceIndex: number | null;
		onCellClick: (index: number) => void;
		onSwap: (from: number, to: number) => void;
		onRemove?: (index: number) => void;
	}

	let { board, selectedPiece, selectedPieceIndex, onCellClick, onSwap, onRemove }: Props = $props();

	let draggedIndex = $state<number | null>(null);
	let menuIndex = $state<number | null>(null);
	let menuX = $state<number>(0);
	let menuY = $state<number>(0);
	let showMenu = $state(false);

	function handleDragStart(index: number) {
		draggedIndex = index;
	}

	function handleDragOver(event: DragEvent) {
		event.preventDefault();
	}

	function handleDrop(index: number) {
		if (draggedIndex !== null && draggedIndex !== index) {
			onSwap(draggedIndex, index);
		}
		draggedIndex = null;
	}

	function handleContextMenu(event: MouseEvent, index: number) {
		event.preventDefault();
		menuIndex = index;
		// position relative to the board container so menu stays within editor
		const container = (event.currentTarget as HTMLElement)?.closest('.board-grid') as HTMLElement | null;
		const rect = container ? container.getBoundingClientRect() : { left: 0, top: 0 };
		menuX = event.clientX - rect.left;
		menuY = event.clientY - rect.top;
		showMenu = true;
	}

	function closeMenu() {
		showMenu = false;
		menuIndex = null;
	}

	function handleRemoveFromCell() {
		if (menuIndex !== null && typeof onRemove === 'function') {
			onRemove(menuIndex);
		}
		closeMenu();
	}

	function getCellClass(index: number): string {
		const classes = ['cell'];
		
		if (board[index] === '.') {
			classes.push('empty');
		}
		if (selectedPiece && board[index] === '.') {
			classes.push('can-place');
		}
		
		return classes.join(' ');
	}

	function getPieceInfo(piece: string) {
		// safe accessor to avoid template parsing edge-cases
		return PIECE_INVENTORY[piece];
	}
</script>

<div class="board-grid">
	<!-- Lakes Row (top) -->
	<div class="lake-row-display">
		{#each Array(10) as _, col}
			<!-- compute lake cells inline to avoid {@const} inside each -->
			{#if (col >= 2 && col <= 3) || (col >= 6 && col <= 7)}
				<div class="cell lake">
					<span class="lake-icon">🌊</span>
				</div>
			{:else}
				<div class="cell empty-row"></div>
			{/if}
		{/each}
	</div>

	<!-- Board Grid (4 rows of pieces) -->
	<div class="grid">
		{#each board as piece, index}
			<div
				class={getCellClass(index)}
				draggable={piece !== '.'}
				ondragstart={() => handleDragStart(index)}
				ondragover={handleDragOver}
				ondrop={() => handleDrop(index)}
				onclick={() => onCellClick(index)}
				onkeydown={(e) => e.key === 'Enter' || e.key === ' ' ? onCellClick(index) : null}
				oncontextmenu={(e) => handleContextMenu(e, index)}
				role="button"
				tabindex="0"
			>
				<SetupPiece
					pieceChar={piece}
					pieceInfo={getPieceInfo(piece)}
					isEmpty={piece === '.'}
					isSelected={selectedPieceIndex === index}
					canPlace={selectedPiece !== null && piece === '.'}
				/>
			</div>
		{/each}
	</div>

	{#if showMenu}
		<button class="context-backdrop" type="button" aria-label="Close context menu" onclick={closeMenu}></button>
		<div class="context-menu" style="left: {menuX}px; top: {menuY}px;">
			<button class="ctx-btn" onclick={handleRemoveFromCell}>Remove</button>
		</div>
	{/if}

	<div class="info">
		<p>💡 Drag & drop pieces to rearrange, or click to pick up and place.</p>
	</div>
</div>

<style>
	.board-grid {
		margin: 0px auto;
		display: flex;
		flex-direction: column;
		gap: 0;
		position: relative; /* allow absolutely-positioned context menu */
	}

	.lake-row-display {
		display: grid;
		grid-template-columns: repeat(10, 1fr);
		gap: 4px;
		background: var(--bg);
		padding: 10px;
		padding-bottom: 5px;
		border-radius: 8px 8px 0 0;
		max-width: 600px;
	}

	.grid {
		display: grid;
		grid-template-columns: repeat(10, 1fr);
		gap: 4px;
		background: var(--bg);
		padding: 10px;
		padding-top: 5px;
		border-radius: 0 0 8px 8px;
		max-width: 600px;
	}

	.cell {
		aspect-ratio: 1;
		display: flex;
		align-items: center;
		justify-content: center;
		border-radius: 6px;
		transition: all 0.2s;
		overflow: hidden;
		box-sizing: border-box;

		&.empty {
			background: #2d3748;
			border: 2px solid #1a202c;
		}

		&.empty-row {
			background: #3a3a2a;
			border: 2px solid #1a202c;
		}

		&.can-place {
			background: #2a4a2a;
			border-color: #4a4;
		}

		&.lake {
			background: linear-gradient(135deg, #1a3a5a 0%, #2a4a6a 100%);
			cursor: default;
			border: 2px solid #4a6a8a;

			.lake-icon {
				font-size: 1.5rem;
			}
		}
	}

	/* context menu */
	.context-backdrop {
		position: absolute;
		top: 0;
		right: 0;
		bottom: 0;
		left: 0;
		touch-action: none;
		background: none;
		border: none;
		cursor: pointer;
	}

	.context-menu {
		position: absolute;
		z-index: 30;
		background: var(--bg);
		border: 1px solid rgba(0,0,0,0.3);
		box-shadow: 0 6px 18px rgba(0,0,0,0.25);
		padding: 6px;
		border-radius: 6px;
	}

	.ctx-btn {
		background: transparent;
		border: none;
		padding: 6px 12px;
		cursor: pointer;
		color: var(--text, #fff);
	}

	.info {
		text-align: center;
		color: var(--muted, #999);
		font-size: 0.9rem;
		margin-top: 10px;

		p {
			margin: 0;
		}
	}

	@media (max-width: 768px) {
		.grid,
		.lake-row-display {
			max-width: 100%;
		}
	}
</style>
