<script lang="ts">
	import type { GameMode } from '../types';

	interface Props {
		onStartGame: (mode: GameMode) => void;
		onSaveGame: () => void;
		onLoadGame: (file: File) => void;
		isGameActive: boolean;
	}

	let { onStartGame, onSaveGame, onLoadGame, isGameActive }: Props = $props();

	let selectedMode = $state<GameMode>('human-vs-ai');
	let fileInput: HTMLInputElement;

	const handleFileSelect = (event: Event) => {
		const target = event.target as HTMLInputElement;
		const file = target.files?.[0];
		if (file) {
			onLoadGame(file);
		}
	};
</script>

<div class="controls">
	<div class="control-section">
		<h3>Start New Game</h3>
		<div class="mode-selector">
			<label>
				<input 
					type="radio" 
					bind:group={selectedMode} 
					value="human-vs-ai"
					disabled={isGameActive}
				/>
				Human vs AI
			</label>
			<label>
				<input 
					type="radio" 
					bind:group={selectedMode} 
					value="ai-vs-ai"
					disabled={isGameActive}
				/>
				AI vs AI
			</label>
		</div>
		<button 
			class="btn btn-primary"
			onclick={() => onStartGame(selectedMode)}
			disabled={isGameActive}
		>
			Start Game
		</button>
	</div>

	<div class="control-section">
		<h3>Game History</h3>
		<div class="history-buttons">
			<button 
				class="btn btn-secondary"
				onclick={onSaveGame}
				disabled={!isGameActive}
			>
				💾 Save Game
			</button>
			<button 
				class="btn btn-secondary"
				onclick={() => fileInput.click()}
			>
				📂 Load Game
			</button>
			<input 
				bind:this={fileInput}
				type="file" 
				accept=".json"
				onchange={handleFileSelect}
				style="display: none;"
			/>
		</div>
	</div>
</div>

<style>
	.controls {
		background: #2d3748;
		padding: 20px;
		border-radius: 8px;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
	}

	.control-section {
		margin-bottom: 25px;
	}

	.control-section:last-child {
		margin-bottom: 0;
	}

	h3 {
		margin: 0 0 15px 0;
		color: #e2e8f0;
		font-size: 1.2rem;
		border-bottom: 2px solid #4a5568;
		padding-bottom: 8px;
	}

	.mode-selector {
		display: flex;
		flex-direction: column;
		gap: 10px;
		margin-bottom: 15px;
	}

	.mode-selector label {
		display: flex;
		align-items: center;
		gap: 8px;
		color: #e2e8f0;
		cursor: pointer;
		padding: 8px;
		border-radius: 4px;
		transition: background 0.2s;
	}

	.mode-selector label:hover {
		background: #4a5568;
	}

	.mode-selector input[type="radio"] {
		cursor: pointer;
	}

	.history-buttons {
		display: flex;
		flex-direction: column;
		gap: 10px;
	}

	.btn {
		padding: 10px 20px;
		border: none;
		border-radius: 6px;
		font-size: 1rem;
		font-weight: 600;
		cursor: pointer;
		transition: all 0.2s;
		width: 100%;
	}

	.btn:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.btn-primary {
		background: #48bb78;
		color: white;
	}

	.btn-primary:hover:not(:disabled) {
		background: #38a169;
		transform: translateY(-2px);
		box-shadow: 0 4px 8px rgba(72, 187, 120, 0.3);
	}

	.btn-secondary {
		background: #4299e1;
		color: white;
	}

	.btn-secondary:hover:not(:disabled) {
		background: #3182ce;
		transform: translateY(-2px);
		box-shadow: 0 4px 8px rgba(66, 153, 225, 0.3);
	}
</style>
