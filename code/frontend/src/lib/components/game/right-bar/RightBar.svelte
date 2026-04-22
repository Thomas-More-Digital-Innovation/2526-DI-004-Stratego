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

    import GameControls from "./GameControls.svelte";
    import GameHistory from "./GameHistory.svelte";
</script>

<Card class="space-y-3 max-h-[500px] flex flex-col">
    <GameControls
        {isReplaying}
        {onSetSpeed}
        {onTogglePause}
        {onStep}
        {onExitReplay}
    />
    <hr class="border-white/10" />
    <GameHistory
        {currentMoveIndex}
        {totalMoves}
        {isReplaying}
        {onPrevious}
        {onNext}
        {onGoToMove}
        {onExitReplay}
    />
</Card>
