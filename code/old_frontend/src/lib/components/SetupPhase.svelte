<script lang="ts">
	import type { Position } from '$lib/types';
	
	interface Props {
		onSwapPieces: (pos1: Position, pos2: Position) => void;
		onRandomize: () => void;
		onStart: () => void;
	}
	
	let { onSwapPieces, onRandomize, onStart }: Props = $props();
	
	let selectedPos = $state<Position | null>(null);
</script>

<div class="setup-phase">
	<div class="setup-banner">
		<div class="setup-header">
			<h2>⚔️ Setup Phase</h2>
			<p>Click two pieces in your setup area (bottom 4 rows) to swap them</p>
		</div>
		
		<div class="setup-controls">
			<button class="btn-randomize" onclick={onRandomize}>
				🎲 Randomize
			</button>
			
			<button class="btn-start" onclick={onStart}>
				▶️ Start Game
			</button>
		</div>
	</div>
</div>

<style>
	.setup-phase {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		z-index: 100;
		pointer-events: none;
	}
	
	.setup-banner {
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
		padding: 1.5rem 2rem;
		box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
		display: flex;
		justify-content: space-between;
		align-items: center;
		gap: 2rem;
		pointer-events: auto;
	}
	
	.setup-header {
		flex: 1;
		color: white;
	}
	
	.setup-header h2 {
		font-size: 1.5rem;
		margin: 0 0 0.25rem 0;
		text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.3);
	}
	
	.setup-header p {
		margin: 0;
		opacity: 0.95;
		font-size: 0.95rem;
	}
	
	.setup-controls {
		display: flex;
		gap: 1rem;
	}
	
	.setup-controls button {
		padding: 0.75rem 1.5rem;
		font-size: 1rem;
		font-weight: 600;
		border: none;
		border-radius: 8px;
		cursor: pointer;
		transition: all 0.2s;
		box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
		white-space: nowrap;
	}
	
	.btn-randomize {
		background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
		color: white;
	}
	
	.btn-randomize:hover {
		transform: translateY(-2px);
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
	}
	
	.btn-start {
		background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
		color: white;
		font-size: 1.1rem;
		padding: 0.75rem 2rem;
	}
	
	.btn-start:hover {
		transform: translateY(-2px);
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
	}
</style>
