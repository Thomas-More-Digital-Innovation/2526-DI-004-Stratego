<script lang="ts">
    import { PIECE_INVENTORY, type PieceInfo } from "$lib/types/board-setup";
    import type { Piece as PieceType, Position } from "$lib/types/game";
    import { decodeSetup, encodeSetup } from "$lib/utils/board-binary";
    import Board from "../game/Board.svelte";
    import Card from "../ui/Card.svelte";
    import Button from "../ui/Button.svelte";

    interface Props {
        initialSetup?: string; // Base64 encoded 40 bytes
        onSave: (setupData: string) => void;
        onCancel: () => void;
    }

    let { initialSetup, onSave, onCancel }: Props = $props();

    // Internal state: 4 rows of 10 cells
    // We store the rank char ('0', 'B', '1'-'9', 'M')
    let grid = $state<string[][]>([]);
    let selectedPieceRank = $state<string | null>(null);
    let selectedGridCell = $state<Position | null>(null);

    // Initialize grid
    $effect(() => {
        if (initialSetup) {
            const decoded = decodeSetup(initialSetup);
            grid = decoded.map((row) =>
                row.split("").map((c) => (c === "." ? " " : c)),
            );
        } else if (grid.length === 0) {
            grid = Array(4)
                .fill(null)
                .map(() => Array(10).fill(" "));
        }
    });

    // Reset board selection when inventory selection changes
    $effect(() => {
        if (selectedPieceRank) {
            selectedGridCell = null;
        }
    });

    // Calculate remaining counts
    const remainingCounts = $derived.by(() => {
        const counts: Record<string, number> = {};
        Object.keys(PIECE_INVENTORY).forEach((k) => {
            counts[k] = PIECE_INVENTORY[k].count;
        });

        grid.forEach((row) => {
            row.forEach((cell) => {
                if (cell !== " " && counts[cell] !== undefined) {
                    counts[cell]--;
                }
            });
        });
        return counts;
    });

    const boardItems = $derived.by(() => {
        const lakeRow = Array(10).fill(null);
        const setupRows = grid.map((row, y) =>
            row.map((char, x) => {
                if (char === " ") return null;
                const info = PIECE_INVENTORY[char];
                if (!info) return null;
                return {
                    rank: info.rank,
                    ownerId: 1,
                    revealed: true,
                    iconBlue: info.icon_blue,
                    iconRed: info.icon_red,
                    position: { x, y: y + 1 },
                } as PieceType;
            }),
        );
        return [lakeRow, ...setupRows];
    });

    function handleCellClick(x: number, y: number) {
        if (y === 0) return; // Top row is context only
        const gridY = y - 1;
        const current = grid[gridY][x];

        if (selectedPieceRank) {
            // Place piece from inventory
            if (
                remainingCounts[selectedPieceRank] > 0 ||
                current === selectedPieceRank
            ) {
                grid[gridY][x] = selectedPieceRank;
            }
        } else {
            // Board-to-board interaction
            if (!selectedGridCell) {
                // First click: select piece if not empty
                if (current !== " ") {
                    selectedGridCell = { x, y };
                }
            } else {
                // Second click: swap or move
                if (selectedGridCell.x === x && selectedGridCell.y === y) {
                    // Clicked same cell: remove/clear cell
                    grid[gridY][x] = " ";
                    selectedGridCell = null;
                } else {
                    // Swap contents
                    const sourceY = selectedGridCell.y - 1;
                    const temp = grid[gridY][x];
                    grid[gridY][x] = grid[sourceY][selectedGridCell.x];
                    grid[sourceY][selectedGridCell.x] = temp;
                    selectedGridCell = null;
                }
            }
        }
    }

    // Drag and Drop Handlers
    function handleInventoryDragStart(e: DragEvent, rank: string) {
        if (e.dataTransfer) {
            e.dataTransfer.setData("application/stratego-rank", rank);
            e.dataTransfer.effectAllowed = "move";
            // Highlight the piece being dragged
            selectedPieceRank = rank;
        }
    }

    function handleBoardDragStart(e: DragEvent, x: number, y: number) {
        if (y === 0) return; // Cannot drag from lake row
        if (e.dataTransfer) {
            e.dataTransfer.setData(
                "application/stratego-pos",
                JSON.stringify({ x, y }),
            );
            e.dataTransfer.effectAllowed = "move";
            selectedGridCell = { x, y };
            selectedPieceRank = null;
        }
    }

    function handleBoardDrop(e: DragEvent, x: number, y: number) {
        e.preventDefault();
        if (y === 0) return; // Cannot drop on lake row
        if (!e.dataTransfer) return;

        const gridY = y - 1;
        const rank = e.dataTransfer.getData("application/stratego-rank");
        const posStr = e.dataTransfer.getData("application/stratego-pos");

        if (rank) {
            // Drop from inventory
            if (remainingCounts[rank] > 0 || grid[gridY][x] === rank) {
                grid[gridY][x] = rank;
            }
        } else if (posStr) {
            // Drop from board (move or swap)
            const sourcePos = JSON.parse(posStr);
            if (sourcePos.x === x && sourcePos.y === y) return;

            const sourceY = sourcePos.y - 1;
            const temp = grid[gridY][x];
            grid[gridY][x] = grid[sourceY][sourcePos.x];
            grid[sourceY][sourcePos.x] = temp;
        }

        selectedGridCell = null;
    }

    const isSetupLake = (x: number, y: number) =>
        y === 0 && [2, 3, 6, 7].includes(x);

    function handleRandomize() {
        const pieces: string[] = [];
        Object.entries(PIECE_INVENTORY).forEach(([rank, info]) => {
            for (let i = 0; i < info.count; i++) {
                pieces.push(rank);
            }
        });

        for (let i = pieces.length - 1; i > 0; i--) {
            const j = Math.floor(Math.random() * (i + 1));
            [pieces[i], pieces[j]] = [pieces[j], pieces[i]];
        }

        const newGrid: string[][] = [];
        for (let i = 0; i < 4; i++) {
            newGrid.push(pieces.slice(i * 10, (i + 1) * 10));
        }
        grid = newGrid;

        selectedPieceRank = null;
        selectedGridCell = null;
    }

    function handleSave() {
        const rows = grid.map((row) => row.join(""));
        onSave(encodeSetup(rows));
    }

    const isComplete = $derived(
        Object.values(remainingCounts).every((count) => count === 0),
    );
</script>

<div class="flex flex-col lg:flex-row gap-6">
    <div class="flex-1 space-y-4">
        <div class="flex items-center justify-between">
            <h2 class="text-xl font-bold text-white">Place your pieces</h2>
            <div class="text-sm text-white/40">
                {grid.flat().filter((c) => c !== " ").length} / 40 pieces placed
            </div>
        </div>

        <div class="bg-black/20 rounded-3xl p-6 border border-white/5">
            <Board
                board={boardItems}
                rows={5}
                cols={10}
                disabledRows={[0]}
                visualDisabledRows={[0]}
                isLakeCell={isSetupLake}
                selectedPosition={selectedGridCell}
                onCellClick={handleCellClick}
                onCellDragStart={handleBoardDragStart}
                onCellDrop={handleBoardDrop}
                isInteractive={true}
                scale={1.3}
            />
        </div>

        <div class="flex justify-between gap-3">
            <Button variant="outline" onclick={handleRandomize}
                >Randomize</Button
            >

            <div class="flex gap-3">
                <Button variant="ghost" onclick={onCancel}>Cancel</Button>
                <Button
                    variant="primary"
                    onclick={handleSave}
                    disabled={!isComplete}
                >
                    Save Setup
                </Button>
            </div>
        </div>
    </div>

    <Card class="w-full lg:w-80 h-fit sticky top-6">
        <h3 class="font-bold text-white mb-4">Inventory</h3>
        <div class="grid grid-cols-2 gap-2">
            {#each Object.entries(PIECE_INVENTORY) as [rank, info]}
                {@const count = remainingCounts[rank]}
                <button
                    class="inventory-item flex items-center gap-2 p-2 rounded-xl border transition-all"
                    class:active={selectedPieceRank === rank}
                    class:out-of-stock={count === 0}
                    onclick={() =>
                        (selectedPieceRank =
                            selectedPieceRank === rank ? null : rank)}
                    draggable={count > 0}
                    ondragstart={(e) => handleInventoryDragStart(e, rank)}
                >
                    <div
                        class="w-8 h-8 flex items-center justify-center text-xl bg-white/5 rounded-lg overflow-hidden p-1"
                    >
                        {#if info.icon_blue && (info.icon_blue.includes("/") || info.icon_blue.includes("."))}
                            <img
                                src={info.icon_blue}
                                alt={info.name}
                                class="w-full h-full object-contain"
                            />
                        {:else}
                            {info.rank}
                        {/if}
                    </div>
                    <div class="flex-1 text-left">
                        <div
                            class="text-[10px] font-bold text-white/40 uppercase leading-none mb-1"
                        >
                            {info.name}
                        </div>
                        <div
                            class="text-sm font-bold text-white flex justify-between"
                        >
                            <span>{info.rank}</span>
                            <span
                                class={count > 0
                                    ? "text-brand-accent"
                                    : "text-white/20"}>x{count}</span
                            >
                        </div>
                    </div>
                </button>
            {/each}
        </div>

        {#if !isComplete}
            <p class="mt-4 text-[10px] text-white/30 italic">
                * You must place all 40 pieces to save the setup.
            </p>
        {/if}
    </Card>
</div>

<style>
    .inventory-item {
        background: rgba(255, 255, 255, 0.03);
        border-color: rgba(255, 255, 255, 0.05);
    }

    .inventory-item:hover:not(.out-of-stock) {
        background: rgba(255, 255, 255, 0.07);
        border-color: rgba(255, 255, 255, 0.1);
    }

    .inventory-item.active {
        background: oklch(0.7 0.2 150 / 0.1);
        border-color: oklch(0.7 0.2 150);
    }

    .inventory-item.out-of-stock {
        opacity: 0.5;
        cursor: default;
    }
</style>
