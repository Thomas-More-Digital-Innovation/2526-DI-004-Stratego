<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import Board from '$lib/components/Board.svelte';
	import GameInfo from '$lib/components/GameInfo.svelte';
	import GameHistory from '$lib/components/GameHistory.svelte';
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
		
		// Request initial state
		setTimeout(() => api.requestState(), 500);
	}

	function setupMessageHandlers() {
		api.onMessage('gameState', (data) => {
			console.log('Received gameState:', data);
			gameStore.updateGameState(data);
		});

		api.onMessage('boardState', (data) => {
			console.log('Received boardState:', data);
			console.log('Sample piece from board[0][0]:', data.board?.[0]?.[0]);
			console.log('Sample piece from board[9][0] (enemy):', data.board?.[9]?.[0]);
			gameStore.updateBoardState(data);
		});

		api.onMessage('moveResult', (data) => {
			console.log('Move result received:', data);
			if (!data.success) {
				errorMessage = data.error || 'Move failed';
				setTimeout(() => errorMessage = '', 3000);
			} else {
				console.log('✓ Move successful, requesting updated state');
				// Request fresh board state after successful move
				setTimeout(() => api.requestState(), 100);
			}
			gameStore.setSelectedPosition(null);
			validMoves = [];
		});

		api.onMessage('aiMove', (data) => {
			console.log('AI Move:', data);
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

	function calculateValidMoves(x: number, y: number) {
		const board = gameStore.boardState?.board;
		if (!board) {
			console.log('No board in calculateValidMoves');
			return [];
		}

		const piece = board[y][x];
		console.log('calculateValidMoves for position', { x, y }, 'piece:', piece);
		
		// Check if cell is actually empty (no ownerName means it's an empty cell from backend)
		// Backend always sets ownerName for real pieces, even if type is hidden
		if (!piece || !piece.ownerName) {
			console.log('Empty cell (no ownerName)');
			return [];
		}
		
		// Check if piece belongs to player
		if (piece.ownerId !== 0) {
			console.log('Piece belongs to opponent, ownerId:', piece.ownerId);
			return [];
		}

		// Bombs and Flags cannot move
		if (piece.type === 'Bomb' || piece.type === 'Flag') {
			console.log('Piece cannot move:', piece.type);
			return [];
		}

		const moves: Position[] = [];
		const directions = [
			{ dx: 0, dy: -1 },  // up
			{ dx: 0, dy: 1 },   // down
			{ dx: -1, dy: 0 },  // left
			{ dx: 1, dy: 0 }    // right
		];

		// Check if piece is Scout (can move multiple spaces)
		const isScout = piece.type === 'Scout';

		for (const dir of directions) {
			let distance = 1;
			const maxDistance = isScout ? 10 : 1;

			while (distance <= maxDistance) {
				const newX = x + (dir.dx * distance);
				const newY = y + (dir.dy * distance);

				// Out of bounds
				if (newX < 0 || newX >= 10 || newY < 0 || newY >= 10) break;

				// Check for lakes
				const lakes = [
					{ x: 2, y: 4 }, { x: 3, y: 4 }, { x: 2, y: 5 }, { x: 3, y: 5 },
					{ x: 6, y: 4 }, { x: 7, y: 4 }, { x: 6, y: 5 }, { x: 7, y: 5 }
				];
				const isLake = lakes.some(l => l.x === newX && l.y === newY);
				if (isLake) break;

				const targetPiece = board[newY][newX];

				// Check if target cell has a piece (has ownerName field)
				if (targetPiece && targetPiece.ownerName) {
					// Can attack enemy pieces
					if (targetPiece.ownerId !== 0) {
						moves.push({ x: newX, y: newY });
					}
					// Stop at any piece
					break;
				}

				// Empty space - valid move
				moves.push({ x: newX, y: newY });
				distance++;
			}
		}

		console.log('Calculated moves:', moves);
		return moves;
	}

	function handleCellClick(x: number, y: number) {
		console.log('=== Cell clicked ===');
		console.log('Position:', { x, y });
		console.log('isReplaying:', gameStore.isReplaying);
		console.log('isHumanTurn:', isHumanTurn);
		console.log('gameMode:', gameMode);
		console.log('currentPlayerId:', gameStore.gameState?.currentPlayerId);
		
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
		// Backend always sets ownerName for real pieces, even if type is hidden
		const hasValidPiece = clickedPiece && clickedPiece.ownerName;

		// If no piece selected, select this one if it belongs to player
		if (!selected) {
			console.log('No piece currently selected');
			
			if (hasValidPiece && clickedPiece.ownerId === 0) {
				console.log('✓ Attempting to select own piece');
				const moves = calculateValidMoves(x, y);
				
				// Only select if piece can actually move
				if (moves.length > 0) {
					console.log('✓ Piece can move, selecting it');
					gameStore.setSelectedPosition({ x, y });
					validMoves = moves;
				} else {
					console.log('❌ Piece has no valid moves');
				}
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
			// Select different piece if it's yours and it can move
			if (hasValidPiece && clickedPiece.ownerId === 0) {
				console.log('Selecting different piece');
				const moves = calculateValidMoves(x, y);
				if (moves.length > 0) {
					gameStore.setSelectedPosition({ x, y });
					validMoves = moves;
				} else {
					// Piece can't move, deselect current
					gameStore.setSelectedPosition(null);
					validMoves = [];
				}
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
