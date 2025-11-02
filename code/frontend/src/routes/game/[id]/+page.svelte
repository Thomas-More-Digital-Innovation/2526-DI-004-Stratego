<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import Board from '$lib/components/Board.svelte';
	import GameInfo from '$lib/components/GameInfo.svelte';
	import GameHistory from '$lib/components/GameHistory.svelte';
	import CombatAnimation from '$lib/components/CombatAnimation.svelte';
	import { GameAPI } from '$lib/api';
	import { gameStore } from '$lib/store.svelte';
	import type { Position } from '$lib/types';

	let gameId = $state<string>('');
	let gameMode = $state<string>('');
	let api = new GameAPI();
	let errorMessage = $state<string>('');
	let isConnected = $state<boolean>(false);
	let validMoves = $state<Position[]>([]);

	const isHumanTurn = $derived.by(() => {
		const result = gameStore.gameState?.currentPlayerId === 0 && 
			gameMode === 'human-vs-ai' && 
			!gameStore.gameState?.isGameOver;
		console.log('🎮 isHumanTurn calculation:', {
			currentPlayerId: gameStore.gameState?.currentPlayerId,
			gameMode,
			isGameOver: gameStore.gameState?.isGameOver,
			waitingForInput: gameStore.gameState?.waitingForInput,
			result
		});
		return result;
	});

	const viewerId = $derived(gameMode === 'human-vs-ai' ? 0 : -1);

	onMount(async () => {
		// Extract game ID from URL path
		const pathParts = window.location.pathname.split('/');
		gameId = pathParts[pathParts.length - 1];
		
		// Get game mode from URL params
		const params = new URLSearchParams(window.location.search);
		gameMode = params.get('mode') || 'human-vs-ai';

		try {
			await connectToGame();
		} catch (error) {
			errorMessage = `Failed to connect: ${error}`;
			setTimeout(() => window.location.href = '/', 3000);
		}
	});

	onDestroy(() => {
		api.disconnect();
		gameStore.reset();
	});

	async function connectToGame() {
		setupMessageHandlers();
		
		const playerId = gameMode === 'human-vs-ai' ? 0 : -1;
		await api.connectWebSocket(gameId as string, playerId);
		isConnected = true;
		
		// Backend will automatically send initial state when we connect
	}

	function setupMessageHandlers() {
		api.onMessage('gameState', (data) => {
			console.log('Received gameState:', data);
			gameStore.updateGameState(data);
		});

		api.onMessage('boardState', (data) => {
			console.log('Received boardState:', data);
			gameStore.updateBoardState(data);
		});

		api.onMessage('moveResult', (data) => {
			console.log('Move result received:', data);
			if (!data.success) {
				errorMessage = data.error || 'Move failed';
				setTimeout(() => errorMessage = '', 3000);
				gameStore.setSelectedPosition(null);
				validMoves = [];
			}
			// Backend will automatically push updated state after move is processed
		});

		api.onMessage('validMoves', (data) => {
			console.log('Valid moves received:', data);
			validMoves = data.validMoves || [];
		});

		api.onMessage('combat', (data) => {
			console.log('💥 Combat message received:', data);
			// Show combat animation
			gameStore.showCombatAnimation({
				attacker: data.attacker,
				defender: data.defender,
				attackerWon: data.attackerWon,
				defenderWon: data.defenderWon
			});
		});

		api.onMessage('gameOver', (data) => {
			console.log('Game Over:', data);
			setTimeout(() => {
				if (confirm('Game Over! Return to menu?')) {
					window.location.href = '/';
				}
			}, 2000);
		});

		api.onMessage('error', (data) => {
			errorMessage = data.error;
			setTimeout(() => errorMessage = '', 3000);
		});
	}

	function handleCellClick(x: number, y: number) {
		console.log('=== Cell clicked ===');
		console.log('Position:', { x, y });
		console.log('isReplaying:', gameStore.isReplaying);
		console.log('isHumanTurn:', isHumanTurn);
		
		// Prevent moves during replay or if not human's turn
		if (gameStore.isReplaying || !isHumanTurn) {
			console.log('❌ Click ignored - not interactive');
			return;
		}

		const board = gameStore.boardState?.board;
		if (!board) {
			console.log('❌ No board state');
			return;
		}

		const clickedPiece = board[y][x];
		const selected = gameStore.selectedPosition;

		console.log('Clicked piece:', clickedPiece);
		console.log('Currently selected:', selected);
		
		// Check if clicked cell has a valid piece (must have ownerName field)
		const hasValidPiece = clickedPiece && clickedPiece.ownerName;

		// If no piece selected, select this one if it belongs to player
		if (!selected) {
			console.log('No piece currently selected');
			
			if (hasValidPiece && clickedPiece.ownerId === 0) {
				console.log('✓ Requesting valid moves from backend');
				gameStore.setSelectedPosition({ x, y });
				// Request valid moves from backend
				api.requestValidMoves({ x, y });
			} else if (hasValidPiece) {
				console.log('❌ Not your piece (ownerId:', clickedPiece.ownerId, ')');
			} else {
				console.log('❌ Empty cell');
			}
			return;
		}

		// If same piece clicked, deselect
		if (selected.x === x && selected.y === y) {
			console.log('Deselecting piece');
			gameStore.setSelectedPosition(null);
			validMoves = [];
			return;
		}

		// Check if clicked position is a valid move
		const isValidMove = validMoves.some(m => m.x === x && m.y === y);
		console.log('Is valid move?', isValidMove);
		
		if (isValidMove) {
			// Make the move
			console.log('✓ Making move from', selected, 'to', { x, y });
			api.sendMove(selected, { x, y });
			gameStore.setSelectedPosition(null);
			validMoves = [];
		} else {
			// Select different piece if it's yours
			if (hasValidPiece && clickedPiece.ownerId === 0) {
				console.log('Selecting different piece');
				gameStore.setSelectedPosition({ x, y });
				api.requestValidMoves({ x, y });
			} else {
				// Clicked on empty or enemy, deselect
				console.log('Deselecting (clicked empty or enemy)');
				gameStore.setSelectedPosition(null);
				validMoves = [];
			}
		}
		console.log('=== End click handler ===');
	}

	function saveGame() {
		try {
			const gameData = gameStore.exportGame();
			const blob = new Blob([gameData], { type: 'application/json' });
			const url = URL.createObjectURL(blob);
			const a = document.createElement('a');
			a.href = url;
			a.download = `stratego-${gameId}-${Date.now()}.json`;
			a.click();
			URL.revokeObjectURL(url);
		} catch (error) {
			errorMessage = 'Failed to save game';
			console.error('Failed to save game:', error);
		}
	}
</script>

<svelte:head>
	<title>Stratego - Game {gameId}</title>
</svelte:head>

<main>
	<div class="container">
		<header>
			<button class="back-btn" onclick={() => window.location.href = '/'}>
				← Back to Menu
			</button>
			<h1>🎮 Stratego Game</h1>
			<button class="save-btn" onclick={saveGame} disabled={!isConnected}>
				💾 Save
			</button>
		</header>

		{#if errorMessage}
			<div class="error-banner">
				⚠️ {errorMessage}
			</div>
		{/if}

		{#if !isConnected}
			<div class="loading-screen">
				<div class="spinner"></div>
				<p>Connecting to game...</p>
			</div>
		{:else}
			<div class="game-layout">
				<div class="left-panel">
					<GameInfo 
						gameState={gameStore.gameState}
						{gameMode}
					/>
				</div>

				<div class="center-panel">
					<Board 
						boardState={gameStore.boardState}
						selectedPosition={gameStore.selectedPosition}
						onCellClick={handleCellClick}
						isInteractive={!gameStore.isReplaying && isHumanTurn}
						{viewerId}
						{validMoves}
					/>
				</div>

				<div class="right-panel">
					<GameHistory
						currentMoveIndex={gameStore.currentHistoryIndex}
						totalMoves={gameStore.history.length}
						isReplaying={gameStore.isReplaying}
						onPrevious={() => gameStore.previousMove()}
						onNext={() => gameStore.nextMove()}
						onGoToMove={(index) => gameStore.goToMove(index)}
						onExitReplay={() => gameStore.exitReplay()}
					/>
				</div>
			</div>
		{/if}
		
		<!-- Combat Animation -->
		{#if gameStore.combatAnimation}
			<CombatAnimation
				attacker={gameStore.combatAnimation.attacker}
				defender={gameStore.combatAnimation.defender}
				attackerWon={gameStore.combatAnimation.attackerWon}
				defenderWon={gameStore.combatAnimation.defenderWon}
				onComplete={() => {
					console.log('Combat animation complete, notifying backend');
					api.sendAnimationComplete();
					gameStore.hideCombatAnimation();
				}}
			/>
		{/if}
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
		padding: 20px;
	}

	.container {
		max-width: 1600px;
		margin: 0 auto;
	}

	header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 30px;
		gap: 20px;
	}

	header h1 {
		margin: 0;
		font-size: 2rem;
		color: #e2e8f0;
		flex: 1;
		text-align: center;
	}

	.back-btn, .save-btn {
		padding: 10px 20px;
		border: none;
		border-radius: 6px;
		font-size: 1rem;
		font-weight: 600;
		cursor: pointer;
		transition: all 0.2s;
		background: #4a5568;
		color: #e2e8f0;
	}

	.back-btn:hover, .save-btn:hover:not(:disabled) {
		background: #718096;
		transform: translateY(-2px);
	}

	.save-btn:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.error-banner {
		background: #fc8181;
		color: white;
		padding: 15px;
		border-radius: 8px;
		margin-bottom: 20px;
		text-align: center;
		font-weight: 600;
		animation: slideDown 0.3s ease-out;
	}

	@keyframes slideDown {
		from {
			opacity: 0;
			transform: translateY(-20px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	.loading-screen {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		min-height: 400px;
		gap: 20px;
	}

	.spinner {
		width: 50px;
		height: 50px;
		border: 4px solid #4a5568;
		border-top-color: #4299e1;
		border-radius: 50%;
		animation: spin 1s linear infinite;
	}

	@keyframes spin {
		to { transform: rotate(360deg); }
	}

	.game-layout {
		display: grid;
		grid-template-columns: 300px 1fr 300px;
		gap: 20px;
		align-items: start;
	}

	.left-panel, .right-panel {
		display: flex;
		flex-direction: column;
		gap: 20px;
	}

	.center-panel {
		display: flex;
		justify-content: center;
	}

	@media (max-width: 1200px) {
		.game-layout {
			grid-template-columns: 1fr;
		}

		header h1 {
			font-size: 1.5rem;
		}
	}
</style>
