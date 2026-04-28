<script lang="ts">
    import type { MoveVisualizationHighlightState } from "$lib/types/game";
    import { fade } from "svelte/transition";

    interface Props {
        x: number;
        y: number;
        cellSize: number;
        gap: number;
        padding: number;
        state: MoveVisualizationHighlightState;
    }

    let { x, y, cellSize, gap, padding, state }: Props = $props();

    const highlightColors: Record<MoveVisualizationHighlightState, string> = {
        move: "bg-brand-accent/20 border-brand-accent/40",
        win: "bg-green-900/40 border-green-500/50",
        loss: "bg-red-500/30 border-red-400/50",
    };
</script>

<div
    class="absolute rounded-md border-2 {highlightColors[state]}"
    style="left: {x * (cellSize + gap) + padding}px; top: {y *
        (cellSize + gap) +
        padding}px; width: {cellSize}px; height: {cellSize}px;"
    in:fade
></div>
