<script lang="ts">
	import type { GameState } from '../types';

	interface Props {
		gameState: GameState | null;
		gameMode: string;
	}

	let { gameState, gameMode }: Props = $props();

	const getPlayerColor = (playerId: number): string => {
		return playerId === 0 ? '#ff6b6b' : '#4dabf7';
	};
</script>

<div class="game-info">
	<div class="info-header">
		<h2>Stratego - {gameMode}</h2>
	</div>

	{#if gameState}
		<div class="info-grid">
			<div class="info-item">
				<span class="label">Round:</span>
				<span class="value">{gameState.round}</span>
			</div>

			<div class="info-item">
				<span class="label">Current Player:</span>
				<span 
					class="value player-name"
					style="color: {getPlayerColor(gameState.currentPlayerId)}"
				>
					{gameState.currentPlayerName}
				</span>
			</div>

			<div class="info-item">
				<span class="label">Player 1 Pieces:</span>
				<span class="value">{gameState.player1AlivePieces}</span>
			</div>

			<div class="info-item">
				<span class="label">Player 2 Pieces:</span>
				<span class="value">{gameState.player2AlivePieces}</span>
			</div>

			<div class="info-item">
				<span class="label">Moves:</span>
				<span class="value">{gameState.moveCount}</span>
			</div>

			{#if gameState.isGameOver}
				<div class="game-over">
					<h3>🎉 Game Over!</h3>
					<p>Winner: <strong>{gameState.winnerName || 'Draw'}</strong></p>
					<p>Cause: {gameState.winCause}</p>
				</div>
			{/if}
		</div>
	{:else}
		<p class="no-data">No game data available</p>
	{/if}
</div>

<style>
	.game-info {
		background: #2d3748;
		padding: 20px;
		border-radius: 8px;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
	}

	.info-header {
		margin-bottom: 15px;
		border-bottom: 2px solid #4a5568;
		padding-bottom: 10px;
	}

	.info-header h2 {
		margin: 0;
		color: #e2e8f0;
		font-size: 1.5rem;
	}

	.info-grid {
		display: grid;
		gap: 10px;
	}

	.info-item {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 8px 0;
		border-bottom: 1px solid #4a5568;
	}

	.label {
		color: #a0aec0;
		font-weight: 500;
	}

	.value {
		color: #e2e8f0;
		font-weight: 600;
	}

	.player-name {
		font-size: 1.1rem;
	}

	.game-over {
		margin-top: 20px;
		padding: 15px;
		background: #48bb78;
		border-radius: 6px;
		text-align: center;
	}

	.game-over h3 {
		margin: 0 0 10px 0;
		color: white;
	}

	.game-over p {
		margin: 5px 0;
		color: white;
	}

	.no-data {
		color: #718096;
		text-align: center;
		padding: 20px 0;
	}
</style>
