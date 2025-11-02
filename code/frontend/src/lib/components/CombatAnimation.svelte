<script lang="ts">
	import type { Piece } from '../types';
	import PieceDisplay from './Piece.svelte';

	interface Props {
		attacker: Piece | null;
		defender: Piece | null;
		attackerWon: boolean;
		defenderWon: boolean;
		onComplete?: () => void;
	}

	let { 
		attacker, 
		defender, 
		attackerWon, 
		defenderWon,
		onComplete 
	}: Props = $props();

	let animationStage = $state<'reveal' | 'clash' | 'result' | 'hide'>('reveal');

	// Animation sequence
	$effect(() => {
		if (attacker && defender) {
			// Stage 1: Reveal (0.5s)
			animationStage = 'reveal';
			
			setTimeout(() => {
				// Stage 2: Clash (0.5s)
				animationStage = 'clash';
				
				setTimeout(() => {
					// Stage 3: Result (1s)
					animationStage = 'result';
					
					setTimeout(() => {
						// Stage 4: Hide and complete (0.5s)
						animationStage = 'hide';
						setTimeout(() => {
							onComplete?.();
						}, 500);
					}, 1000);
				}, 500);
			}, 500);
		}
	});
</script>

{#if attacker && defender}
	<div class="combat-overlay">
		<div class="combat-arena" class:reveal={animationStage === 'reveal'} 
			class:clash={animationStage === 'clash'} 
			class:result={animationStage === 'result'}
			class:hide={animationStage === 'hide'}>
			
			<div class="combat-title">COMBAT!</div>
			
			<div class="fighters">
				<div class="fighter attacker" class:winner={attackerWon} class:loser={!attackerWon && defenderWon}>
					<div class="fighter-label">{attacker.ownerName || 'Attacker'}</div>
					<div class="piece-wrapper">
						<PieceDisplay piece={attacker} viewerId={attacker.ownerId} />
					</div>
					<div class="piece-info">
						{#if attacker.type}
							<div class="type">{attacker.type}</div>
						{/if}
						{#if attacker.rank}
							<div class="rank">Rank: {attacker.rank}</div>
						{/if}
					</div>
				</div>

				<div class="vs">VS</div>

				<div class="fighter defender" class:winner={defenderWon} class:loser={!defenderWon && attackerWon}>
					<div class="fighter-label">{defender.ownerName || 'Defender'}</div>
					<div class="piece-wrapper">
						<PieceDisplay piece={defender} viewerId={defender.ownerId} />
					</div>
					<div class="piece-info">
						{#if defender.type}
							<div class="type">{defender.type}</div>
						{/if}
						{#if defender.rank}
							<div class="rank">Rank: {defender.rank}</div>
						{/if}
					</div>
				</div>
			</div>

			{#if animationStage === 'result'}
				<div class="result-text">
					{#if attackerWon && defenderWon}
						Both Eliminated!
					{:else if attackerWon}
						Attacker Wins!
					{:else if defenderWon}
						Defender Wins!
					{/if}
				</div>
			{/if}
		</div>
	</div>
{/if}

<style>
	.combat-overlay {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background: rgba(0, 0, 0, 0.85);
		display: flex;
		justify-content: center;
		align-items: center;
		z-index: 1000;
		animation: fadeIn 0.3s ease-in;
	}

	@keyframes fadeIn {
		from {
			opacity: 0;
		}
		to {
			opacity: 1;
		}
	}

	.combat-arena {
		background: linear-gradient(135deg, #1a202c 0%, #2d3748 100%);
		padding: 40px;
		border-radius: 16px;
		box-shadow: 0 8px 32px rgba(0, 0, 0, 0.5);
		border: 3px solid #4a5568;
		min-width: 600px;
	}

	.combat-arena.reveal {
		animation: zoomIn 0.5s ease-out;
	}

	.combat-arena.clash .fighters {
		animation: shake 0.5s ease-in-out;
	}

	.combat-arena.hide {
		animation: fadeOut 0.5s ease-in;
	}

	@keyframes zoomIn {
		from {
			transform: scale(0.5);
			opacity: 0;
		}
		to {
			transform: scale(1);
			opacity: 1;
		}
	}

	@keyframes shake {
		0%, 100% { transform: translateX(0); }
		25% { transform: translateX(-10px); }
		75% { transform: translateX(10px); }
	}

	@keyframes fadeOut {
		from {
			opacity: 1;
		}
		to {
			opacity: 0;
		}
	}

	.combat-title {
		text-align: center;
		font-size: 2.5rem;
		font-weight: bold;
		color: #ffd700;
		margin-bottom: 30px;
		text-shadow: 0 0 10px rgba(255, 215, 0, 0.5);
		animation: pulse 1s ease-in-out infinite;
	}

	.fighters {
		display: flex;
		justify-content: space-around;
		align-items: center;
		gap: 40px;
	}

	.fighter {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 15px;
		padding: 20px;
		background: rgba(255, 255, 255, 0.05);
		border-radius: 12px;
		min-width: 200px;
		transition: all 0.3s ease;
	}

	.fighter.winner {
		background: rgba(144, 238, 144, 0.2);
		border: 2px solid #90ee90;
		box-shadow: 0 0 20px rgba(144, 238, 144, 0.5);
		transform: scale(1.1);
	}

	.fighter.loser {
		opacity: 0.5;
		filter: grayscale(0.7);
		transform: scale(0.9);
	}

	.fighter-label {
		font-size: 1.2rem;
		font-weight: bold;
		color: #e2e8f0;
	}

	.piece-wrapper {
		width: 100px;
		height: 100px;
		display: flex;
		align-items: center;
		justify-content: center;
		background: #2d3748;
		border-radius: 8px;
		border: 2px solid #4a5568;
	}

	.piece-info {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 5px;
	}

	.type {
		font-size: 1.1rem;
		font-weight: bold;
		color: #ffd700;
	}

	.rank {
		font-size: 0.9rem;
		color: #cbd5e0;
	}

	.vs {
		font-size: 2rem;
		font-weight: bold;
		color: #ff6b6b;
		text-shadow: 0 0 10px rgba(255, 107, 107, 0.5);
		animation: pulse 0.5s ease-in-out infinite;
	}

	.result-text {
		text-align: center;
		font-size: 1.8rem;
		font-weight: bold;
		margin-top: 30px;
		color: #90ee90;
		animation: slideUp 0.5s ease-out;
	}

	@keyframes slideUp {
		from {
			transform: translateY(20px);
			opacity: 0;
		}
		to {
			transform: translateY(0);
			opacity: 1;
		}
	}

	@keyframes pulse {
		0%, 100% {
			transform: scale(1);
		}
		50% {
			transform: scale(1.1);
		}
	}

	@media (max-width: 768px) {
		.combat-arena {
			min-width: 90vw;
			padding: 20px;
		}

		.fighters {
			flex-direction: column;
			gap: 20px;
		}

		.fighter {
			min-width: 150px;
		}

		.vs {
			transform: rotate(90deg);
		}
	}
</style>
