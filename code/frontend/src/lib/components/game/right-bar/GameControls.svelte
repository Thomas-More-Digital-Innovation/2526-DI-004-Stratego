<script lang="ts">
    import Button from "$lib/components/ui/Button.svelte";
    import { gameStore } from "$lib/state/game.svelte";

    interface Props {
        isReplaying: boolean;
        onSetSpeed: (speedMs: number) => void;
        onTogglePause: () => void;
        onStep: () => void;
        onExitReplay: () => void;
    }

    let {
        isReplaying,
        onSetSpeed,
        onTogglePause,
        onStep,
        onExitReplay,
    }: Props = $props();

    let speedMs = $state(1000);
</script>

<h3 class="text-sm font-bold text-brand-accent uppercase tracking-wider">
    Controls
</h3>
<div class="flex flex-col gap-2">
    <div class="flex gap-2">
        {#if isReplaying && !gameStore.gameState?.isGameOver}
            <Button
                variant="secondary"
                size="sm"
                class="flex-1"
                onclick={onExitReplay}
            >
                Exit Replay
            </Button>
        {:else if !gameStore.gameState?.isGameOver}
            <Button
                variant="outline"
                size="sm"
                onclick={onTogglePause}
                class="flex-1"
            >
                {gameStore.isPaused ? "▶ Resume" : "⏸ Pause"}
            </Button>
            {#if gameStore.isPaused && gameStore.gameMode == "ai_vs_ai"}
                <Button
                    variant="ghost"
                    size="sm"
                    onclick={onStep}
                    loading={gameStore.isStepping}
                    disabled={gameStore.isStepping}
                >
                    Step ⏭️
                </Button>
            {/if}
        {/if}
    </div>

    <div class="space-y-1.5 px-1">
        <div class="flex justify-between text-[10px] text-white/40">
            <span>Speed</span>
            <span>{speedMs}ms</span>
        </div>
        <input
            type="range"
            min="500"
            max="5000"
            step="100"
            bind:value={speedMs}
            onchange={() => onSetSpeed(speedMs)}
            class="w-full accent-brand-primary h-1 bg-white/10 rounded-lg appearance-none cursor-pointer"
        />
    </div>
</div>
