<script lang="ts">
    import Button from "$lib/components/ui/Button.svelte";
    import { gameStore } from "$lib/state/game.svelte";
    interface Props {
        currentMoveIndex: number;
        totalMoves: number;
        isReplaying: boolean;
        onPrevious: () => void;
        onNext: () => void;
        onGoToMove: (index: number) => void;
        onExitReplay: () => void;
    }

    let {
        currentMoveIndex,
        totalMoves,
        isReplaying,
        onPrevious,
        onNext,
        onGoToMove,
        onExitReplay,
    }: Props = $props();

    const canGoPrevious = $derived(currentMoveIndex > 0);
    const canGoNext = $derived(currentMoveIndex < totalMoves - 1);
</script>

<div class="flex items-center justify-between">
    <h3 class="text-sm font-bold text-brand-accent uppercase tracking-wider">
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
{:else}
    <p class="text-white/30 text-center py-4 text-sm">No moves yet</p>
{/if}
