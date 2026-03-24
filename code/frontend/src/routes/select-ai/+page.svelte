<script lang="ts">
    import { page } from "$app/stores";
    import { goto } from "$app/navigation";
    import { games } from "$lib/api/client";
    import { AIs } from "$lib/data/AI.data";
    import Card from "$lib/components/ui/Card.svelte";
    import Button from "$lib/components/ui/Button.svelte";
    import type { GameMode } from "$lib/types/game";

    const gameMode =
        ($page.url.searchParams.get("mode") as GameMode) || "human_vs_ai";

    let step = $state<"ai1" | "ai2">("ai1");
    let ai1 = $state("");
    let ai2 = $state("");
    let creating = $state(false);
    let error = $state("");

    async function selectAi(aiId: string) {
        if (step === "ai1") {
            ai1 = aiId;
            if (gameMode === "human_vs_ai") {
                await start();
            } else {
                step = "ai2";
            }
        } else {
            ai2 = aiId;
            await start();
        }
    }

    async function start() {
        creating = true;
        error = "";
        try {
            const info = await games.create(gameMode, ai1, ai2);
            goto(`/game/${info.gameId}?mode=${gameMode}`);
        } catch (e: any) {
            error = e.message || "Failed to create game";
            creating = false;
        }
    }

    function back() {
        if (step === "ai2") {
            step = "ai1";
            ai1 = "";
        } else {
            goto("/");
        }
    }
</script>

<svelte:head>
    <title>Stratego — Select AI</title>
</svelte:head>

<div class="space-y-8 max-w-4xl mx-auto py-10">
    <header class="flex items-center justify-between">
        <div class="space-y-1">
            <h1
                class="text-3xl font-extrabold text-white uppercase tracking-widest"
            >
                Select {step === "ai1" ? "AI Opponent" : "Second AI"}
            </h1>
            <p class="text-white/50">
                {gameMode.replace(/_/g, " ").toUpperCase()}
            </p>
        </div>
        <Button variant="ghost" onclick={back}>
            {step === "ai2" ? "Back to AI 1" : "Cancel"}
        </Button>
    </header>

    {#if error}
        <div
            class="bg-brand-secondary/20 border border-brand-secondary/30 text-brand-secondary rounded-xl px-4 py-3 text-sm text-center"
        >
            {error}
        </div>
    {/if}

    <div class="flex flex-wrap gap-4">
        {#each AIs as ai}
            <button
                class="text-left flex-1 p-6 rounded-2xl border-2 transition-all duration-200 flex items-center gap-6 {ai1 ===
                    ai.id || ai2 === ai.id
                    ? 'border-brand-primary bg-brand-primary/10 scale-102'
                    : 'border-white/5 bg-surface-elevated/20 hover:border-white/20 hover:bg-white/5'}"
                onclick={() => selectAi(ai.id)}
                disabled={creating}
            >
                {#if ai.image}
                    <img
                        src={ai.image}
                        alt={ai.name}
                        class="w-20 h-20 rounded-xl object-cover bg-black/40 border border-white/5 shadow-inner"
                    />
                {/if}
                <div class="flex-1">
                    <div class="flex items-center justify-between">
                        <h3
                            class="text-xl font-bold text-white uppercase tracking-wider"
                        >
                            {ai.name}
                        </h3>
                        {#if ai1 === ai.id}
                            <span
                                class="text-[10px] font-bold bg-brand-primary text-white px-2 py-1 rounded"
                                >SELECTED</span
                            >
                        {/if}
                    </div>
                    <p class="text-white/40 text-sm mt-1 leading-relaxed">
                        {ai.description}
                    </p>
                </div>
            </button>
        {/each}
    </div>

    {#if creating}
        <div class="flex flex-col items-center gap-4 py-6">
            <div
                class="w-10 h-10 border-4 border-white/10 border-t-brand-primary rounded-full animate-spin"
            ></div>
            <p class="text-white/40 text-sm animate-pulse">
                Initializing strategic algorithms...
            </p>
        </div>
    {/if}
</div>
