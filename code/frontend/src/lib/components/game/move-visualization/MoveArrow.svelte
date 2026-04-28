<script lang="ts">
    import type { HistoricalMove, Move } from "$lib/types/game";
    import { fade } from "svelte/transition";

    interface Props {
        move: HistoricalMove;
        cellSize: number;
        scale: number;
        gap: number;
        padding: number;
    }

    let { move, cellSize, scale, gap, padding }: Props = $props();

    let fromX = $derived(
        move.fromX * (cellSize + gap) + padding + cellSize / 2,
    );
    let fromY = $derived(
        move.fromY * (cellSize + gap) + padding + cellSize / 2,
    );
    let toX = $derived(move.toX * (cellSize + gap) + padding + cellSize / 2);
    let toY = $derived(move.toY * (cellSize + gap) + padding + cellSize / 2);
    let angle = $derived(Math.atan2(toY - fromY, toX - fromX));
    let arrowLength = $derived(
        Math.sqrt((toX - fromX) ** 2 + (toY - fromY) ** 2),
    );
</script>

<svg
    class="absolute inset-0 w-full h-full"
    viewBox="0 0 {10 * (cellSize + gap) + padding * 2} {10 * (cellSize + gap) +
        padding * 2}"
>
    <defs>
        <marker
            id="arrowhead"
            markerWidth="5"
            markerHeight="4"
            refX="3.5"
            refY="2"
            orient="auto"
        >
            <polygon
                points="0 0, 5 2, 0 4"
                fill="oklch(0.7 0.2 150)"
                fill-opacity="0.8"
            />
        </marker>
    </defs>

    <line
        x1={fromX}
        y1={fromY}
        x2={fromX + Math.cos(angle) * arrowLength}
        y2={fromY + Math.sin(angle) * arrowLength}
        stroke="oklch(0.7 0.2 150)"
        stroke-width={4 * scale}
        stroke-opacity="0.6"
        marker-end="url(#arrowhead)"
        in:fade={{ duration: 300 }}
    />
</svg>
