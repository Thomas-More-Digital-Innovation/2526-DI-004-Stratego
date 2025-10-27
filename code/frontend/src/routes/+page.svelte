<script lang="ts">
	import { createGameState, startGame } from '$lib/gameState';
	import GameBoard from '$lib/components/GameBoard.svelte';
	import GameInfo from '$lib/components/GameInfo.svelte';
	import type { GameState } from '$lib/types';

	// Use Svelte 5 rune for state
	let gameState = $state<GameState>(createGameState());

	function handleStartGame() {
		gameState = startGame(gameState);
	}

	function handleNewGame() {
		gameState = createGameState();
	}

	function handleStateChange(newState: GameState) {
		gameState = newState;
	}
</script>

<svelte:head>
	<title>Stratego - Strategic Board Game</title>
</svelte:head>

<div class="container">
	<header>
		<h1>⚔️ Stratego ⚔️</h1>
		<p class="subtitle">Strategic Board Game of Hidden Ranks and Tactical Warfare</p>
	</header>

	{#if gameState.phase === 'setup'}
		<div class="setup-screen">
			<div class="setup-content">
				<h2>Welcome to Stratego!</h2>
				<p>A classic game of strategy and deception.</p>
				
				<div class="rules">
					<h3>Quick Rules:</h3>
					<ul>
						<li>🎯 <strong>Objective:</strong> Capture your opponent's Flag</li>
						<li>♟️ <strong>Movement:</strong> Most pieces move one square (Scouts can move multiple)</li>
						<li>⚔️ <strong>Combat:</strong> Higher rank wins (with special exceptions)</li>
						<li>💣 <strong>Bombs:</strong> Destroy all attackers except Miners</li>
						<li>🕵️ <strong>Spy:</strong> Can defeat the Marshal if attacking first</li>
						<li>⛏️ <strong>Miners:</strong> Can defuse Bombs</li>
					</ul>
				</div>

				<div class="setup-info">
					<p>Pieces have been randomly placed for both players.</p>
					<p>Red player goes first!</p>
				</div>

				<button class="start-btn" onclick={handleStartGame}>
					Start Game
				</button>
			</div>
		</div>
	{:else}
		<div class="game-container">
			<div class="game-area">
				<GameBoard state={gameState} onStateChange={handleStateChange} />
			</div>
			
			<aside class="sidebar">
				<GameInfo state={gameState} />
				
				<div class="controls">
					<button class="btn btn-primary" onclick={handleNewGame}>
						New Game
					</button>
				</div>
			</aside>
		</div>
	{/if}

	<footer>
		<p>Created with Svelte 5 | Digital Innovation Stratego Project</p>
	</footer>
</div>

<style>
	:global(body) {
		margin: 0;
		padding: 0;
		font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
		min-height: 100vh;
	}

	.container {
		min-height: 100vh;
		display: flex;
		flex-direction: column;
	}

	header {
		text-align: center;
		padding: 2rem 1rem;
		background: rgba(255, 255, 255, 0.95);
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
	}

	header h1 {
		margin: 0;
		font-size: 3rem;
		color: #333;
		text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.1);
	}

	.subtitle {
		margin: 0.5rem 0 0 0;
		font-size: 1.1rem;
		color: #666;
		font-style: italic;
	}

	.setup-screen {
		flex: 1;
		display: flex;
		align-items: center;
		justify-content: center;
		padding: 2rem;
	}

	.setup-content {
		background: white;
		padding: 3rem;
		border-radius: 20px;
		max-width: 600px;
		box-shadow: 0 12px 40px rgba(0, 0, 0, 0.2);
		text-align: center;
	}

	.setup-content h2 {
		margin: 0 0 1rem 0;
		font-size: 2.5rem;
		color: #333;
	}

	.setup-content > p {
		font-size: 1.2rem;
		color: #666;
		margin-bottom: 2rem;
	}

	.rules {
		text-align: left;
		background: #f8f9fa;
		padding: 1.5rem;
		border-radius: 12px;
		margin-bottom: 2rem;
	}

	.rules h3 {
		margin: 0 0 1rem 0;
		color: #333;
		font-size: 1.3rem;
	}

	.rules ul {
		margin: 0;
		padding-left: 1.5rem;
		list-style: none;
	}

	.rules li {
		margin: 0.75rem 0;
		color: #555;
		line-height: 1.6;
	}

	.setup-info {
		background: #e3f2fd;
		padding: 1rem;
		border-radius: 8px;
		margin-bottom: 2rem;
	}

	.setup-info p {
		margin: 0.5rem 0;
		color: #1976d2;
		font-weight: 500;
	}

	.start-btn {
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
		color: white;
		border: none;
		padding: 1rem 3rem;
		font-size: 1.3rem;
		font-weight: bold;
		border-radius: 12px;
		cursor: pointer;
		transition: all 0.3s ease;
		box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
	}

	.start-btn:hover {
		transform: translateY(-2px);
		box-shadow: 0 6px 20px rgba(102, 126, 234, 0.6);
	}

	.start-btn:active {
		transform: translateY(0);
	}

	.game-container {
		flex: 1;
		display: flex;
		gap: 2rem;
		padding: 2rem;
		max-width: 1400px;
		margin: 0 auto;
		width: 100%;
	}

	.game-area {
		flex: 1;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.sidebar {
		width: 400px;
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.controls {
		background: white;
		padding: 1rem;
		border-radius: 12px;
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
	}

	.btn {
		width: 100%;
		padding: 0.75rem 1.5rem;
		font-size: 1rem;
		font-weight: bold;
		border: none;
		border-radius: 8px;
		cursor: pointer;
		transition: all 0.3s ease;
	}

	.btn-primary {
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
		color: white;
		box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
	}

	.btn-primary:hover {
		transform: translateY(-2px);
		box-shadow: 0 6px 16px rgba(102, 126, 234, 0.5);
	}

	footer {
		text-align: center;
		padding: 1.5rem;
		background: rgba(255, 255, 255, 0.95);
		color: #666;
		font-size: 0.9rem;
	}

	footer p {
		margin: 0;
	}

	@media (max-width: 1024px) {
		.game-container {
			flex-direction: column;
		}

		.sidebar {
			width: 100%;
		}
	}

	@media (max-width: 768px) {
		header h1 {
			font-size: 2rem;
		}

		.subtitle {
			font-size: 0.9rem;
		}

		.setup-content {
			padding: 2rem 1.5rem;
		}

		.setup-content h2 {
			font-size: 2rem;
		}

		.game-container {
			padding: 1rem;
			gap: 1rem;
		}
	}
</style>
