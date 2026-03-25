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
    import BoardSetupCard from "$lib/components/setup/BoardSetupCard.svelte";

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
        <div class="flex flex-wrap gap-6">
            {#each setups as setup}
                <BoardSetupCard {setup} ownerId={1}>
                    {#snippet actions()}
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
                    {/snippet}
                </BoardSetupCard>
            {/each}
        </div>
    {/if}
</div>
