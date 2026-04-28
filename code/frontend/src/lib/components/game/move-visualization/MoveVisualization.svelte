<script lang="ts">
    import type { HistoricalMove, Position } from "$lib/types/game";

    import { BOARD_CONFIG } from "$lib/data/board.data";
    import MoveArrow from "./MoveArrow.svelte";
    import CombatIndicator from "./CombatIndicator.svelte";
    import MoveHighlight from "./MoveHighlight.svelte";

    interface Props {
        move: HistoricalMove;
        cellSize: number;
        scale?: number;
    }

    let { move, cellSize, scale = 1 }: Props = $props();

    const padding = BOARD_CONFIG.padding;
    const gap = $derived(BOARD_CONFIG.gap * scale);

    const isCombat = $derived(move.result !== "move");
</script>

<div
    class="absolute inset-0 pointer-events-none z-20"
    style="--cell-size: {cellSize}px; --scale: {scale}"
>
    <!-- Highlights -->
    {#if isCombat}
        {#if move.result === "win" || move.result === "capture"}
            <MoveHighlight
                x={move.fromX}
                y={move.fromY}
                {cellSize}
                {gap}
                {padding}
                state="win"
            />
            <MoveHighlight
                x={move.toX}
                y={move.toY}
                {cellSize}
                {gap}
                {padding}
                state="loss"
            />
        {:else if move.result === "loss"}
            <MoveHighlight
                x={move.fromX}
                y={move.fromY}
                {cellSize}
                {gap}
                {padding}
                state="loss"
            />
            <MoveHighlight
                x={move.toX}
                y={move.toY}
                {cellSize}
                {gap}
                {padding}
                state="win"
            />
        {:else if move.result === "tie"}
            <MoveHighlight
                x={move.fromX}
                y={move.fromY}
                {cellSize}
                {gap}
                {padding}
                state="loss"
            />
            <MoveHighlight
                x={move.toX}
                y={move.toY}
                {cellSize}
                {gap}
                {padding}
                state="loss"
            />
        {/if}
    {:else}
        <MoveHighlight
            x={move.fromX}
            y={move.fromY}
            {cellSize}
            {gap}
            {padding}
            state="move"
        />
        <MoveHighlight
            x={move.toX}
            y={move.toY}
            {cellSize}
            {gap}
            {padding}
            state="move"
        />
    {/if}

    <MoveArrow {move} {cellSize} {scale} {gap} {padding} />

    <!-- Indicators at destination -->
    <div
        class="absolute flex items-center justify-center"
        style="left: {move.toX * (cellSize + gap) + padding}px; top: {move.toY *
            (cellSize + gap) +
            padding}px; width: {cellSize}px; height: {cellSize}px;"
    >
        {#if isCombat}
            <CombatIndicator {move} />
        {/if}
    </div>
</div>
