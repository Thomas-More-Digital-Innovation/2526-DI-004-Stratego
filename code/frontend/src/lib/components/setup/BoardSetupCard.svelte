<script lang="ts">
    import type { Snippet } from "svelte";
    import { decodeSetup } from "$lib/utils/board-binary";
    import { PIECE_INVENTORY } from "$lib/types/board-setup";
    import type { BoardSetup } from "$lib/types/board-setup";
    import type { Piece as PieceType } from "$lib/types/game";
    import Board from "$lib/components/game/Board.svelte";
    import Card from "$lib/components/ui/Card.svelte";

    interface Props {
        setup: BoardSetup;
        ownerId?: number;
        onclick?: () => void;
        class?: string;
        actions?: Snippet;
        isInteractive?: boolean;
    }

    let {
        setup,
        ownerId = 1,
        onclick,
        class: className = "",
        actions,
        isInteractive = false,
    }: Props = $props();

    const isRedSide = $derived(ownerId === 1);

    function getDecodedBoard(setupData: string): (PieceType | null)[][] {
        const decoded = decodeSetup(setupData);
        return decoded.map((row, y) =>
            row.split("").map((char, x) => {
                if (char === "." || char === " ") return null;
                const info = PIECE_INVENTORY[char];
                if (!info) return null;
                return {
                    rank: info.rank,
                    ownerId: ownerId,
                    revealed: true,
                    position: { x, y },
                } as PieceType;
            }),
        );
    }
</script>

<Card
    class="group flex flex-col w-full  md:w-md h-full bg-white/5 border-white/5 hover:border-white/10 transition-all {isInteractive
        ? 'hover:scale-[1.02] cursor-pointer ring-0 hover:ring-2'
        : ''} {isInteractive
        ? isRedSide
            ? 'hover:ring-brand-secondary'
            : 'hover:ring-brand-primary'
        : ''} {className}"
    {onclick}
>
    <div class="flex justify-between items-start mb-4">
        <div>
            <h3
                class="font-bold text-white text-lg group-hover:text-brand-accent transition-colors"
            >
                {setup.name}
            </h3>
            {#if setup.description}
                <p class="text-white/40 text-xs mt-1 line-clamp-1">
                    {setup.description}
                </p>
            {/if}
        </div>
        {#if setup.is_default}
            <span
                class="text-[10px] font-bold bg-brand-accent/20 text-brand-accent px-2 py-0.5 rounded-full uppercase"
            >
                Default
            </span>
        {/if}
    </div>

    <Board
        board={getDecodedBoard(setup.setup_data)}
        rows={4}
        cols={10}
        responsive={true}
        isInteractive={false}
    />

    <div class="mt-auto flex flex-col gap-3 pt-4 border-t border-white/5">
        <div class="flex items-center justify-between">
            <div
                class="text-white/20 text-[10px] uppercase font-bold tracking-wider"
            >
                Updated {new Date(setup.updated_at).toLocaleDateString()}
            </div>
            {#if !actions && isInteractive}
                <div
                    class="text-brand-accent opacity-0 group-hover:opacity-100 transition-opacity font-bold text-xs uppercase tracking-widest"
                >
                    Select →
                </div>
            {/if}
        </div>

        {#if actions}
            {@render actions()}
        {/if}
    </div>
</Card>
