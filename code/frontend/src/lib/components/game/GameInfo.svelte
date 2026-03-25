<script lang="ts">
    import Card from "$lib/components/ui/Card.svelte";
    import type { GameState } from "$lib/types/game";

    interface Props {
        gameState: GameState | null;
        gameMode: string;
    }

    let { gameState, gameMode }: Props = $props();
</script>

<Card class="space-y-4">
    <div class="border-b border-white/10 pb-3">
        <h2
            class="text-lg font-bold text-brand-accent uppercase tracking-wider"
        >
            {gameMode.replace(/_/g, " ")}
        </h2>
    </div>

    {#if gameState}
        <div class="space-y-2 text-sm">
            <div
                class="flex justify-between items-center py-1 border-b border-white/5"
            >
                <span class="text-white/50">Round</span>
                <span class="font-semibold text-white">{gameState.round}</span>
            </div>

            <div
                class="flex justify-between items-center py-1 border-b border-white/5"
            >
                <span class="text-white/50">Current</span>
                <span
                    class="font-semibold"
                    class:text-brand-secondary={gameState.currentPlayerId === 0}
                    class:text-brand-primary={gameState.currentPlayerId === 1}
                >
                    {gameState.currentPlayerName}
                </span>
            </div>

            <div
                class="flex justify-between items-center py-1 border-b border-white/5"
            >
                <span class="text-white/50">P1 Pieces</span>
                <span class="font-semibold text-brand-secondary"
                    >{gameState.player1AlivePieces}</span
                >
            </div>

            <div
                class="flex justify-between items-center py-1 border-b border-white/5"
            >
                <span class="text-white/50">P2 Pieces</span>
                <span class="font-semibold text-brand-primary"
                    >{gameState.player2AlivePieces}</span
                >
            </div>

            {#if !gameState.isSetupPhase}
                <div class="flex justify-between items-center py-1">
                    <span class="text-white/50">Moves</span>
                    <span class="font-semibold text-white"
                        >{gameState.moveCount}</span
                    >
                </div>
            {/if}
        </div>

        {#if gameState.isGameOver}
            <div
                class="rounded-xl bg-brand-accent/20 border border-brand-accent/30 p-4 text-center space-y-1"
            >
                <h3
                    class="text-brand-accent font-bold uppercase tracking-wider"
                >
                    Game Over
                </h3>
                <p class="text-white text-sm">
                    Winner: <strong>{gameState.winnerName || "Draw"}</strong>
                </p>
                <p class="text-white/60 text-xs">{gameState.winCause}</p>
            </div>
        {/if}
    {:else}
        <p class="text-white/30 text-center py-4">No game data</p>
    {/if}
</Card>
