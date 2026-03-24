<script lang="ts">
    import Piece from "./Piece.svelte";
    import type { BoardState, Position } from "$lib/types/game";

    interface Props {
        boardState: BoardState | null;
        selectedPosition: Position | null;
        onCellClick: (x: number, y: number) => void;
        isInteractive?: boolean;
        viewerId?: number;
        validMoves?: Position[];
    }

    let {
        boardState,
        selectedPosition,
        onCellClick,
        isInteractive = true,
        viewerId = 0,
        validMoves = [],
    }: Props = $props();

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
        lakePositions.some((pos) => pos.x === x && pos.y === y);

    const isSelected = (x: number, y: number) =>
        selectedPosition?.x === x && selectedPosition?.y === y;

    const isValidMove = (x: number, y: number) =>
        validMoves.some((m) => m.x === x && m.y === y);
</script>

<div class="board-wrapper">
    {#if boardState}
        <div class="board glass rounded-2xl p-3">
            {#each boardState.board as row, y}
                {#each row as piece, x}
                    <button
                        class="cell"
                        class:interactive={isInteractive && !isLake(x, y)}
                        class:valid-move={isValidMove(x, y)}
                        onclick={() =>
                            isInteractive && !isLake(x, y) && onCellClick(x, y)}
                        disabled={!isInteractive || isLake(x, y)}
                    >
                        <Piece
                            {piece}
                            isSelected={isSelected(x, y)}
                            isHighlighted={isValidMove(x, y)}
                            isLake={isLake(x, y)}
                            {viewerId}
                        />
                    </button>
                {/each}
            {/each}
        </div>
    {:else}
        <div class="flex items-center justify-center p-10 text-white/40">
            Waiting for board state...
        </div>
    {/if}
</div>

<style>
    .board-wrapper {
        display: flex;
        justify-content: center;
        align-items: center;
    }

    .board {
        display: grid;
        grid-template-columns: repeat(10, 48px);
        grid-template-rows: repeat(10, 48px);
        gap: 2px;
    }

    .cell {
        width: 48px;
        height: 48px;
        position: relative;
        padding: 0;
        border: none;
        background: none;
    }

    .cell.interactive {
        cursor: pointer;
    }

    .cell:disabled {
        cursor: default;
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

    @media (max-width: 768px) {
        .board {
            grid-template-columns: repeat(10, 34px);
            grid-template-rows: repeat(10, 34px);
        }

        .cell {
            width: 34px;
            height: 34px;
        }
    }
</style>
