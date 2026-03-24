<script lang="ts">
    import type { Piece as PieceType } from "$lib/types/game";

    interface Props {
        piece: PieceType | null;
        isSelected?: boolean;
        isHighlighted?: boolean;
        isLake?: boolean;
        viewerId?: number;
    }

    let {
        piece,
        isSelected = false,
        isHighlighted = false,
        isLake = false,
        viewerId = 0,
    }: Props = $props();

    const canSeePiece = $derived(() => {
        if (!piece || !piece.ownerName) return false;
        return piece.ownerId === viewerId || piece.revealed;
    });
</script>

<div
    class="piece"
    class:selected={isSelected}
    class:highlighted={isHighlighted}
    class:empty={!piece || !piece.ownerName}
    class:lake={isLake}
    class:player0={piece?.ownerName && piece.ownerId === 0}
    class:player1={piece?.ownerName && piece.ownerId === 1}
>
    {#if isLake}
        <span class="text-2xl">🌊</span>
    {:else if piece && piece.ownerName}
        {#if canSeePiece()}
            <div class="flex flex-col items-center justify-center gap-0.5">
                {#if piece.icon}
                    <span class="text-lg leading-none">{piece.icon}</span>
                {/if}
                {#if piece.rank}
                    <span
                        class="text-[10px] font-bold text-white bg-black/30 px-1.5 rounded leading-tight"
                        >{piece.rank}</span
                    >
                {/if}
            </div>
        {:else}
            <span class="text-xl font-bold text-white/80">?</span>
        {/if}
    {/if}
</div>

<style>
    .piece {
        width: 100%;
        height: 100%;
        display: flex;
        align-items: center;
        justify-content: center;
        border-radius: 6px;
        transition: all 0.15s;
        border: 1px solid rgba(255, 255, 255, 0.05);
    }

    .piece.empty {
        background: var(--color-surface-elevated);
    }

    .piece.lake {
        background: var(--color-brand-primary);
        opacity: 0.4;
        cursor: default;
    }

    .piece.player0 {
        background: var(--color-brand-secondary);
        cursor: pointer;
    }

    .piece.player1 {
        background: var(--color-brand-primary);
        cursor: pointer;
    }

    .piece.player0:hover,
    .piece.player1:hover {
        transform: scale(1.08);
        z-index: 1;
    }

    .piece.selected {
        border: 2px solid var(--color-brand-accent);
        box-shadow: 0 0 12px
            color-mix(in srgb, var(--color-brand-accent), transparent 40%);
    }

    .piece.highlighted {
        border: 2px solid oklch(0.75 0.18 145);
    }
</style>
