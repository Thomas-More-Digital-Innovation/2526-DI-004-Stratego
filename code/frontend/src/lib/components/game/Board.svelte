<script lang="ts">
    import Piece from "./Piece.svelte";
    import type {
        BoardState,
        Position,
        Piece as PieceType,
        HistoricalMove,
    } from "$lib/types/game";
    import MoveVisualization from "./move-visualization/MoveVisualization.svelte";
    import { BOARD_CONFIG } from "$lib/data/board.data";

    interface Props {
        boardState?: BoardState | null;
        board?: (PieceType | null)[][]; // Allow passing raw board array
        selectedPosition?: Position | null;
        onCellClick?: (x: number, y: number) => void;
        isInteractive?: boolean;
        viewerId?: number;
        disabledRows?: number[];
        validMoves?: Position[];
        rows?: number;
        cols?: number;
        scale?: number; // Scale factor (e.g., 0.5 for small preview)
        onCellDragStart?: (e: DragEvent, x: number, y: number) => void;
        onCellDragOver?: (e: DragEvent, x: number, y: number) => void;
        onCellDragLeave?: (e: DragEvent, x: number, y: number) => void;
        onCellDrop?: (e: DragEvent, x: number, y: number) => void;
        isLakeCell?: (x: number, y: number) => boolean;
        responsive?: boolean;
        visualDisabledRows?: number[];
        lastMove?: HistoricalMove | null;
    }

    let {
        boardState,
        board,
        selectedPosition,
        onCellClick,
        isInteractive = true,
        viewerId = 0,
        disabledRows = [],
        validMoves = [],
        rows = BOARD_CONFIG.rows,
        cols = BOARD_CONFIG.cols,
        scale = 1,
        onCellDragStart,
        onCellDragOver,
        onCellDragLeave,
        onCellDrop,
        isLakeCell,
        responsive = false,
        visualDisabledRows = [],
        lastMove = null,
    }: Props = $props();

    const displayBoard = $derived(board || boardState?.board || []);

    const lakePositions = [
        { x: 2, y: 4 },
        { x: 3, y: 4 },
        { x: 2, y: 5 },
        { x: 3, y: 5 },
        { x: 6, y: 4 },
        { x: 7, y: 4 },
        { x: 6, y: 5 },
        { x: 7, y: 5 },
    ];

    const isLake = (x: number, y: number) =>
        isLakeCell?.(x, y) ??
        (rows === 10 &&
            lakePositions.some((pos) => pos.x === x && pos.y === y));

    const isSelected = (x: number, y: number) =>
        selectedPosition?.x === x && selectedPosition?.y === y;

    const isValidMove = (x: number, y: number) =>
        validMoves.some((m) => m.x === x && m.y === y);

    const cellSize = $derived(BOARD_CONFIG.baseCellSize * scale);
</script>

<div
    class="board-wrapper"
    style="--cell-size: {cellSize}px; --cols: {cols}; --rows: {rows}; --scale: {scale}; --gap: {BOARD_CONFIG.gap}px;"
    class:non-interactive={!isInteractive}
>
    {#if displayBoard.length > 0}
        <div class="board glass rounded-2xl p-3" class:responsive>
            {#each Array(rows) as _, y}
                {#each Array(cols) as _, x}
                    {@const piece = displayBoard[y]?.[x]}
                    <button
                        class="cell"
                        class:lake={isLake(x, y)}
                        class:interactive={isInteractive && !isLake(x, y)}
                        class:valid-move={isValidMove(x, y)}
                        class:visual-disabled={visualDisabledRows.includes(y)}
                        onclick={() => onCellClick?.(x, y)}
                        disabled={!isInteractive ||
                            isLake(x, y) ||
                            disabledRows.includes(y)}
                        draggable={isInteractive && !!piece && !isLake(x, y)}
                        ondragstart={(e) => onCellDragStart?.(e, x, y)}
                        ondragover={(e) => {
                            if (onCellDrop) {
                                e.preventDefault();
                                onCellDragOver?.(e, x, y);
                            }
                        }}
                        ondragleave={(e) => onCellDragLeave?.(e, x, y)}
                        ondrop={(e) => onCellDrop?.(e, x, y)}
                    >
                        <Piece
                            {piece}
                            isSelected={isSelected(x, y)}
                            isHighlighted={isValidMove(x, y)}
                            isLake={isLake(x, y)}
                            {viewerId}
                            {scale}
                        />
                    </button>
                {/each}
            {/each}

            {#if lastMove}
                <MoveVisualization move={lastMove} {cellSize} {scale} />
            {/if}
        </div>
    {:else}
        <div class="flex items-center justify-center p-10 text-white/40">
            No board data available
        </div>
    {/if}
</div>

<style>
    .board-wrapper {
        display: flex;
        justify-content: center;
        align-items: center;
    }

    .board-wrapper.non-interactive {
        pointer-events: none;
    }

    .board {
        position: relative;
        display: grid;
        grid-template-columns: repeat(var(--cols), var(--cell-size));
        grid-template-rows: repeat(var(--rows), var(--cell-size));
        gap: calc(var(--gap) * var(--scale, 1));
    }

    .board.responsive {
        grid-template-columns: repeat(var(--cols), 1fr);
        grid-template-rows: auto;
        width: 100%;
        gap: 2px;
    }

    .board.responsive .cell {
        width: 100%;
        height: auto;
        aspect-ratio: 1;
    }

    .cell {
        width: var(--cell-size);
        height: var(--cell-size);
        position: relative;
        padding: 0;
        border: none;
        border-radius: 6px;
        background: none;
        cursor: default;
    }

    .cell.interactive {
        cursor: pointer;
    }

    .cell.visual-disabled {
        background-color: rgba(255, 0, 0, 0.1);
    }

    .cell.visual-disabled:not(.lake)::before {
        content: "";
        position: absolute;
        inset: 0;
        background: linear-gradient(
            to top right,
            transparent calc(50% - 1px),
            rgba(255, 255, 255, 0.2) 50%,
            transparent calc(50% + 1px)
        );
        border-radius: 6px;
        pointer-events: none;
        z-index: 1;
    }

    .cell.valid-move::after {
        content: "";
        position: absolute;
        inset: 3px;
        border: 2px solid oklch(0.75 0.18 145);
        border-radius: 4px;
        pointer-events: none;
        animation: pulse 1.5s ease-in-out infinite;
    }

    @keyframes pulse {
        0%,
        100% {
            opacity: 1;
        }
        50% {
            opacity: 0.4;
        }
    }
</style>
