<script lang="ts">
	import type { GameState } from '$lib/types';

	interface Props {
		state: GameState;
	}

	let { state }: Props = $props();

	let currentPlayer = $derived(state.players[state.currentPlayerIndex]);
	let playerPieceCount = $derived(state.players.map(p => p.pieces.filter(piece => piece.alive).length));
</script>

<div class="game-info">
	<div class="info-section">
		<h2>Round {state.round}</h2>
		{#if state.phase === 'playing'}
			<div class="current-turn" class:red={currentPlayer.color === 'red'} class:blue={currentPlayer.color === 'blue'}>
				<span class="turn-indicator">▶</span>
				{currentPlayer.name}'s Turn
			</div>
		{/if}
	</div>

	<div class="players-info">
		{#each state.players as player, i}
			<div class="player-card" class:active={state.currentPlayerIndex === i} class:red={player.color === 'red'} class:blue={player.color === 'blue'}>
				<h3>{player.name}</h3>
				<div class="stats">
					<div class="stat">
						<span class="stat-label">Pieces</span>
						<span class="stat-value">{playerPieceCount[i]}/40</span>
					</div>
				</div>
			</div>
		{/each}
	</div>

	{#if state.selectedPiece}
		<div class="selected-piece-info">
			<h3>Selected Piece</h3>
			<div class="piece-details">
				<span class="piece-icon-large">{state.selectedPiece.type.icon}</span>
				<div class="piece-info">
					<strong>{state.selectedPiece.type.name}</strong>
					<p>{state.selectedPiece.type.description}</p>
					<small>Rank: {state.selectedPiece.type.rank}</small>
				</div>
			</div>
		</div>
	{/if}

	{#if state.phase === 'gameOver' && state.winner}
		<div class="game-over">
			<h2>🎉 Game Over! 🎉</h2>
			<p class="winner" class:red={state.winner.color === 'red'} class:blue={state.winner.color === 'blue'}>
				{state.winner.name} wins!
			</p>
			{#if state.winCause}
				<p class="win-cause">
					{#if state.winCause === 'flag_captured'}
						Flag Captured!
					{:else if state.winCause === 'no_movable_pieces'}
						Opponent has no movable pieces
					{:else if state.winCause === 'max_turns'}
						Maximum turns reached
					{/if}
				</p>
			{/if}
		</div>
	{/if}

	<div class="move-history">
		<h3>Move History ({state.moveHistory.length} moves)</h3>
		<div class="history-list">
			{#each state.moveHistory.slice(-10).reverse() as move, i}
				<div class="history-item">
					<span class="move-number">{state.moveHistory.length - i}</span>
					<span>({move.from.x},{move.from.y}) → ({move.to.x},{move.to.y})</span>
				</div>
			{/each}
		</div>
	</div>
</div>

<style>
	.game-info {
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
		padding: 1rem;
		max-width: 400px;
	}

	.info-section {
		background: white;
		padding: 1.5rem;
		border-radius: 12px;
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
	}

	.info-section h2 {
		margin: 0 0 1rem 0;
		color: #333;
		font-size: 1.5rem;
	}

	.current-turn {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.75rem 1rem;
		border-radius: 8px;
		font-weight: bold;
		font-size: 1.1rem;
		color: white;
		animation: glow 2s ease-in-out infinite;
	}

	.current-turn.red {
		background: linear-gradient(135deg, #ff6b6b 0%, #ee5a52 100%);
	}

	.current-turn.blue {
		background: linear-gradient(135deg, #4a90e2 0%, #357abd 100%);
	}

	@keyframes glow {
		0%, 100% {
			box-shadow: 0 0 20px rgba(255, 215, 0, 0.5);
		}
		50% {
			box-shadow: 0 0 30px rgba(255, 215, 0, 0.8);
		}
	}

	.turn-indicator {
		font-size: 1.2rem;
		animation: bounce 1s ease-in-out infinite;
	}

	@keyframes bounce {
		0%, 100% {
			transform: translateX(0);
		}
		50% {
			transform: translateX(5px);
		}
	}

	.players-info {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.player-card {
		background: white;
		padding: 1rem;
		border-radius: 12px;
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
		transition: all 0.3s ease;
		border: 3px solid transparent;
	}

	.player-card.active {
		border-color: #ffd700;
		transform: scale(1.02);
	}

	.player-card.red h3 {
		color: #ee5a52;
	}

	.player-card.blue h3 {
		color: #357abd;
	}

	.player-card h3 {
		margin: 0 0 0.5rem 0;
		font-size: 1.2rem;
	}

	.stats {
		display: flex;
		gap: 1rem;
	}

	.stat {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}

	.stat-label {
		font-size: 0.8rem;
		color: #666;
		text-transform: uppercase;
		letter-spacing: 0.5px;
	}

	.stat-value {
		font-size: 1.3rem;
		font-weight: bold;
		color: #333;
	}

	.selected-piece-info {
		background: linear-gradient(135deg, #fff7e6 0%, #ffe6cc 100%);
		padding: 1rem;
		border-radius: 12px;
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
		border: 2px solid #ffd700;
	}

	.selected-piece-info h3 {
		margin: 0 0 1rem 0;
		color: #8b6914;
	}

	.piece-details {
		display: flex;
		gap: 1rem;
		align-items: center;
	}

	.piece-icon-large {
		font-size: 3rem;
		filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.2));
	}

	.piece-info {
		flex: 1;
	}

	.piece-info strong {
		display: block;
		font-size: 1.1rem;
		color: #333;
		margin-bottom: 0.25rem;
	}

	.piece-info p {
		margin: 0.25rem 0;
		font-size: 0.9rem;
		color: #666;
	}

	.piece-info small {
		color: #999;
	}

	.game-over {
		background: linear-gradient(135deg, #ffd700 0%, #ffed4e 100%);
		padding: 2rem;
		border-radius: 12px;
		box-shadow: 0 8px 24px rgba(0, 0, 0, 0.2);
		text-align: center;
		animation: celebrate 0.5s ease-in-out;
	}

	@keyframes celebrate {
		0% {
			transform: scale(0.8) rotate(-5deg);
			opacity: 0;
		}
		50% {
			transform: scale(1.05) rotate(2deg);
		}
		100% {
			transform: scale(1) rotate(0deg);
			opacity: 1;
		}
	}

	.game-over h2 {
		margin: 0 0 1rem 0;
		font-size: 2rem;
		color: #333;
	}

	.winner {
		font-size: 1.5rem;
		font-weight: bold;
		margin: 0.5rem 0;
		padding: 0.5rem;
		border-radius: 8px;
		color: white;
	}

	.winner.red {
		background: #ee5a52;
	}

	.winner.blue {
		background: #357abd;
	}

	.win-cause {
		margin: 0.5rem 0 0 0;
		color: #666;
		font-size: 1rem;
	}

	.move-history {
		background: white;
		padding: 1rem;
		border-radius: 12px;
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
	}

	.move-history h3 {
		margin: 0 0 0.75rem 0;
		color: #333;
		font-size: 1rem;
	}

	.history-list {
		max-height: 200px;
		overflow-y: auto;
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}

	.history-item {
		display: flex;
		gap: 0.5rem;
		padding: 0.5rem;
		background: #f5f5f5;
		border-radius: 6px;
		font-size: 0.85rem;
		color: #666;
	}

	.move-number {
		font-weight: bold;
		color: #333;
		min-width: 2rem;
	}

	@media (max-width: 768px) {
		.game-info {
			max-width: 100%;
		}
	}
</style>
