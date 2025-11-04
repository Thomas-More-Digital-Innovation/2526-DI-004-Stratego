<script lang="ts">
	import { GameAPI } from '$lib/api';
	import type { GameMode } from '$lib/types';

	let api = new GameAPI();
	let selectedMode = $state<GameMode>('human_vs_ai');
	let isCreating = $state<boolean>(false);
	let errorMessage = $state<string>('');
	let loadedGameData = $state<string | null>(null);
	let fileInput: HTMLInputElement;

	async function startNewGame() {
		isCreating = true;
		errorMessage = '';

		try {
			const gameInfo = await api.createGame(selectedMode);
			// Use window.location for navigation
			window.location.href = `/game/${gameInfo.gameId}?mode=${selectedMode}`;
		} catch (error) {
			errorMessage = `Failed to create game: ${error}`;
			isCreating = false;
		}
	}

	function handleFileSelect(event: Event) {
		const target = event.target as HTMLInputElement;
		const file = target.files?.[0];
		
		if (file) {
			const reader = new FileReader();
			reader.onload = (e) => {
				try {
					const content = e.target?.result as string;
					loadedGameData = content;
					alert('Game loaded! (Replay feature coming soon)');
				} catch (error) {
					errorMessage = 'Failed to load game file';
				}
			};
			reader.readAsText(file);
		}
	}
</script>

<svelte:head>
	<title>Stratego - Menu</title>
</svelte:head>

<main>
	<div class="container">
		<header>
			<h1>🎮 Stratego</h1>
			<p class="subtitle">Choose your game mode and start playing</p>
		</header>

		{#if errorMessage}
			<div class="error-banner">
				⚠️ {errorMessage}
			</div>
		{/if}

		<div class="menu-content">
			<div class="game-modes">
				<h2>New Game</h2>
				
				<div class="mode-cards">
					<button 
						class="mode-card"
						class:selected={selectedMode === 'human_vs_ai'}
						onclick={() => selectedMode = 'human_vs_ai'}
						disabled={isCreating}
					>
						<div class="mode-icon">🧑 vs 🤖</div>
						<h3>Human vs AI</h3>
						<p>Play against the computer. You control the red pieces, AI controls blue.</p>
						<div class="mode-features">
							<span>• Strategic gameplay</span>
							<span>• AI opponent</span>
							<span>• Turn-based</span>
						</div>
					</button>

					<button 
						class="mode-card"
						class:selected={selectedMode === 'ai_vs_ai'}
						onclick={() => selectedMode = 'ai_vs_ai'}
						disabled={isCreating}
					>
						<div class="mode-icon">🤖 vs 🤖</div>
						<h3>AI vs AI</h3>
						<p>Watch two AI players battle it out. All pieces visible for spectating.</p>
						<div class="mode-features">
							<span>• Auto-play mode</span>
							<span>• Watch & learn</span>
							<span>• Fast-paced</span>
						</div>
					</button>
				</div>

				<button 
					class="start-btn"
					onclick={startNewGame}
					disabled={isCreating}
				>
					{#if isCreating}
						Creating game...
					{:else}
						Start Game
					{/if}
				</button>
			</div>

			<div class="divider"></div>

			<div class="saved-games">
				<h2>Load Saved Game</h2>
				<p class="description">Load a previously saved game to review the moves and strategy.</p>
				
				<button 
					class="load-btn"
					onclick={() => fileInput.click()}
				>
					📂 Load Game File
				</button>
				
				<input 
					bind:this={fileInput}
					type="file" 
					accept=".json"
					onchange={handleFileSelect}
					style="display: none;"
				/>
			</div>

			<div class="info-section">
				<h2>How to Play</h2>
				<div class="instructions">
					<div class="instruction">
						<span class="number">1</span>
						<div>
							<strong>Choose Mode</strong>
							<p>Select Human vs AI to play, or AI vs AI to watch</p>
						</div>
					</div>
					<div class="instruction">
						<span class="number">2</span>
						<div>
							<strong>Make Moves</strong>
							<p>Click your piece, then click where to move (highlighted cells)</p>
						</div>
					</div>
					<div class="instruction">
						<span class="number">3</span>
						<div>
							<strong>Capture Flag</strong>
							<p>Find and capture the enemy flag to win the game</p>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</main>

<style>
	:global(body) {
		margin: 0;
		padding: 0;
		background: #1a202c;
		color: #e2e8f0;
		font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
	}

	main {
		min-height: 100vh;
		padding: 40px 20px;
	}

	.container {
		max-width: 1200px;
		margin: 0 auto;
	}

	header {
		text-align: center;
		margin-bottom: 50px;
	}

	header h1 {
		margin: 0;
		font-size: 3.5rem;
		color: #e2e8f0;
		text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.5);
		margin-bottom: 10px;
	}

	.subtitle {
		color: #a0aec0;
		font-size: 1.2rem;
		margin: 0;
	}

	.error-banner {
		background: #fc8181;
		color: white;
		padding: 15px;
		border-radius: 8px;
		margin-bottom: 30px;
		text-align: center;
		font-weight: 600;
	}

	.menu-content {
		display: flex;
		flex-direction: column;
		gap: 40px;
	}

	.game-modes, .saved-games, .info-section {
		background: #2d3748;
		padding: 30px;
		border-radius: 12px;
		box-shadow: 0 4px 6px rgba(0, 0, 0, 0.3);
	}

	h2 {
		margin: 0 0 20px 0;
		color: #e2e8f0;
		font-size: 1.8rem;
		border-bottom: 2px solid #4a5568;
		padding-bottom: 10px;
	}

	.description {
		color: #a0aec0;
		margin: 0 0 20px 0;
	}

	.mode-cards {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
		gap: 20px;
		margin-bottom: 30px;
	}

	.mode-card {
		background: #4a5568;
		border: 3px solid transparent;
		border-radius: 12px;
		padding: 25px;
		cursor: pointer;
		transition: all 0.3s;
		text-align: left;
	}

	.mode-card:hover:not(:disabled) {
		transform: translateY(-5px);
		box-shadow: 0 8px 16px rgba(0, 0, 0, 0.4);
		border-color: #4299e1;
	}

	.mode-card.selected {
		border-color: #48bb78;
		background: #2f855a;
	}

	.mode-card:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}

	.mode-icon {
		font-size: 3rem;
		margin-bottom: 15px;
	}

	.mode-card h3 {
		margin: 0 0 10px 0;
		color: #e2e8f0;
		font-size: 1.4rem;
	}

	.mode-card p {
		margin: 0 0 15px 0;
		color: #cbd5e0;
		line-height: 1.5;
	}

	.mode-features {
		display: flex;
		flex-direction: column;
		gap: 5px;
		font-size: 0.9rem;
		color: #a0aec0;
	}

	.start-btn {
		width: 100%;
		padding: 18px;
		background: #48bb78;
		color: white;
		border: none;
		border-radius: 8px;
		font-size: 1.3rem;
		font-weight: 700;
		cursor: pointer;
		transition: all 0.2s;
		text-transform: uppercase;
		letter-spacing: 1px;
	}

	.start-btn:hover:not(:disabled) {
		background: #38a169;
		transform: translateY(-2px);
		box-shadow: 0 6px 12px rgba(72, 187, 120, 0.4);
	}

	.start-btn:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}

	.divider {
		height: 2px;
		background: #4a5568;
	}

	.load-btn {
		width: 100%;
		padding: 15px;
		background: #4299e1;
		color: white;
		border: none;
		border-radius: 8px;
		font-size: 1.1rem;
		font-weight: 600;
		cursor: pointer;
		transition: all 0.2s;
	}

	.load-btn:hover {
		background: #3182ce;
		transform: translateY(-2px);
		box-shadow: 0 4px 8px rgba(66, 153, 225, 0.4);
	}

	.instructions {
		display: flex;
		flex-direction: column;
		gap: 20px;
	}

	.instruction {
		display: flex;
		gap: 20px;
		align-items: flex-start;
	}

	.number {
		width: 40px;
		height: 40px;
		background: #4299e1;
		color: white;
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
		font-weight: bold;
		font-size: 1.2rem;
		flex-shrink: 0;
	}

	.instruction strong {
		display: block;
		color: #e2e8f0;
		margin-bottom: 5px;
		font-size: 1.1rem;
	}

	.instruction p {
		margin: 0;
		color: #cbd5e0;
		line-height: 1.5;
	}

	@media (max-width: 768px) {
		header h1 {
			font-size: 2.5rem;
		}

		.mode-cards {
			grid-template-columns: 1fr;
		}
	}
</style>
