<script lang="ts">
    import Card from "$lib/components/ui/Card.svelte";
    import Button from "$lib/components/ui/Button.svelte";

    interface Props {
        currentMoveIndex: number;
        totalMoves: number;
        isReplaying: boolean;
        onPrevious: () => void;
        onNext: () => void;
        onGoToMove: (index: number) => void;
        onExitReplay: () => void;
        onTogglePause: () => void;
        onSetSpeed: (speedMs: number) => void;
        onStep: () => void;
    }

    let {
        currentMoveIndex,
        totalMoves,
        isReplaying,
        onPrevious,
        onNext,
        onGoToMove,
        onExitReplay,
        onTogglePause,
        onSetSpeed,
        onStep,
    }: Props = $props();

    let speedMs = $state(1000);

    import { gameStore } from "$lib/state/game.svelte";

    const canGoPrevious = $derived(currentMoveIndex > 0);
    const canGoNext = $derived(currentMoveIndex < totalMoves - 1);
</script>

<Card class="space-y-3 max-h-[500px] flex flex-col">
    <div
        class="flex justify-between items-center border-b border-white/10 pb-2"
    >
        <h3
            class="text-sm font-bold text-brand-accent uppercase tracking-wider"
        >
            Move History
        </h3>
        {#if isReplaying}
            <span
                class="text-[10px] font-bold bg-brand-secondary/20 text-brand-secondary px-2 py-0.5 rounded-full uppercase"
            >
                Replay
            </span>
        {/if}
        {#if gameStore.isPaused}
            <span
                class="text-[10px] font-bold bg-yellow-500/20 text-yellow-500 px-2 py-0.5 rounded-full uppercase"
            >
                Paused
            </span>
        {/if}
    </div>

    {#if totalMoves > 0}
        <p class="text-white/40 text-xs text-center">
            Move {currentMoveIndex + 1} of {totalMoves}
        </p>

        <div class="grid grid-cols-2 gap-2">
            <Button
                variant="outline"
                size="sm"
                onclick={onPrevious}
                disabled={!canGoPrevious}
            >
                ◀ Prev
            </Button>
            <Button
                variant="outline"
                size="sm"
                onclick={onNext}
                disabled={!canGoNext}
            >
                Next ▶
            </Button>
        </div>

        <div class="flex-1 overflow-y-auto space-y-1 min-h-0">
            {#each Array(totalMoves) as _, index}
                <button
                    class="w-full text-left px-3 py-1.5 rounded-lg text-xs transition-all {index ===
                    currentMoveIndex
                        ? 'bg-brand-primary/20 text-brand-primary font-semibold'
                        : 'text-white/40 hover:bg-white/5 hover:text-white/70'}"
                    onclick={() => onGoToMove(index)}
                >
                    Move {index + 1}
                </button>
            {/each}
        </div>

        <div class="flex flex-col gap-3 pt-2 border-t border-white/5">
            {#if isReplaying}
                <Button variant="secondary" size="sm" onclick={onExitReplay}>
                    Exit Replay
                </Button>
            {:else}
                <div class="flex gap-2">
                    <Button
                        variant="outline"
                        size="sm"
                        onclick={onTogglePause}
                        class="flex-1"
                    >
                        {gameStore.isPaused ? "▶ Resume" : "⏸ Pause"}
                    </Button>
                    {#if gameStore.isPaused}
                        <Button
                            variant="ghost"
                            size="sm"
                            onclick={onStep}
                            class="px-2"
                        >
                            Step ⏭️
                        </Button>
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
            {/if}
        </div>
    {:else}
        <p class="text-white/30 text-center py-4 text-sm">No moves yet</p>
    {/if}
</Card>
