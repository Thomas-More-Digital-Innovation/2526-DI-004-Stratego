<script lang="ts">
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import { authStore } from "$lib/state/auth.svelte";
    import { boardSetups } from "$lib/api/client";
    import Card from "$lib/components/ui/Card.svelte";
    import Button from "$lib/components/ui/Button.svelte";
    import Board from "$lib/components/game/Board.svelte";
    import type { BoardSetup } from "$lib/types/board-setup";
    import { PIECE_INVENTORY, MAX_BOARD_SETUPS } from "$lib/types/board-setup";
    import { decodeSetup } from "$lib/utils/board-binary";
    import type { Piece as PieceType } from "$lib/types/game";

    let setups = $state<BoardSetup[]>([]);
    let error = $state("");
    let loading = $state(true);

    onMount(async () => {
        if (!authStore.user) {
            goto("/login");
            return;
        }
        await loadSetups();
    });

    async function loadSetups() {
        loading = true;
        try {
            const result = await boardSetups.list();
            setups = result ?? [];
        } catch (e: any) {
            error = e.message || "Failed to load setups";
        } finally {
            loading = false;
        }
    }

    async function deleteSetup(id: number) {
        if (!confirm("Delete this setup?")) return;
        try {
            await boardSetups.delete(id);
            await loadSetups();
        } catch (e: any) {
            error = "Failed to delete: " + e.message;
        }
    }

    function getDecodedBoard(setupData: string): (PieceType | null)[][] {
        const decoded = decodeSetup(setupData);
        return decoded.map((row, y) =>
            row.split("").map((char, x) => {
                if (char === "." || char === " ") return null;
                const info = PIECE_INVENTORY[char];
                if (!info) return null;
                return {
                    rank: info.rank,
                    ownerId: 1,
                    revealed: true,
                    position: { x, y },
                } as PieceType;
            }),
        );
    }
</script>

<svelte:head>
    <title>Stratego — Board Setups</title>
</svelte:head>

<div class="space-y-6">
    <div class="flex items-center justify-between">
        <div>
            <h1
                class="text-2xl font-extrabold text-white uppercase tracking-widest"
            >
                Board Setups
            </h1>
            <p class="text-white/40 text-sm mt-1">
                {setups.length}/{MAX_BOARD_SETUPS} setups saved
            </p>
        </div>
        <Button
            variant="primary"
            disabled={setups.length >= MAX_BOARD_SETUPS}
            onclick={() => goto("/board-setups/new")}
        >
            + Create New
        </Button>
    </div>

    {#if error}
        <div
            class="bg-brand-secondary/20 border border-brand-secondary/30 text-brand-secondary rounded-xl px-4 py-3 text-sm text-center"
        >
            {error}
        </div>
    {/if}

    {#if loading}
        <div class="text-center py-12 text-white/30">Loading...</div>
    {:else if setups.length === 0}
        <Card class="text-center py-12">
            <p class="text-white/30">
                No board setups yet. Create your first setup!
            </p>
        </Card>
    {:else}
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {#each setups as setup}
                <Card
                    class="flex flex-col h-full bg-white/5 border-white/5 hover:border-white/10 transition-colors"
                >
                    <div class="flex justify-between items-start mb-4">
                        <div>
                            <h3 class="font-bold text-white text-lg">
                                {setup.name}
                            </h3>
                            {#if setup.description}
                                <p
                                    class="text-white/40 text-xs mt-1 line-clamp-1"
                                >
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

                    <div
                        class="flex-1 bg-black/20 rounded-2xl p-3 scale-[0.6] origin-top-left -mr-[40%] -mb-[80px]"
                    >
                        <Board
                            board={getDecodedBoard(setup.setup_data)}
                            rows={4}
                            cols={10}
                            scale={1}
                            isInteractive={false}
                        />
                    </div>

                    <div class="mt-auto space-y-3">
                        <div
                            class="text-white/20 text-[10px] uppercase font-bold tracking-wider"
                        >
                            Updated {new Date(
                                setup.updated_at,
                            ).toLocaleDateString()}
                        </div>
                        <div class="flex gap-2">
                            <Button
                                variant="outline"
                                size="sm"
                                class="flex-1"
                                onclick={() =>
                                    goto(`/board-setups/${setup.id}`)}
                            >
                                Edit
                            </Button>
                            <Button
                                variant="ghost"
                                size="sm"
                                onclick={() => deleteSetup(setup.id)}
                            >
                                Delete
                            </Button>
                        </div>
                    </div></Card
                >
            {/each}
        </div>
    {/if}
</div>
