<script lang="ts">
    import { page } from "$app/stores";
    import { goto } from "$app/navigation";
    import { games } from "$lib/api/client";
    import { AIs } from "$lib/data/AI.data";
    import Button from "$lib/components/ui/Button.svelte";
    import { gamemodes } from "$lib/data/gamemodes.data";
    import SelectionSummary from "$lib/components/setup/SelectionSummary.svelte";

    const gameMode = gamemodes.fromString(
        $page.url.searchParams.get("mode") || "",
    );

    let step = $state<"ai1" | "ai2">("ai1");
    let ai1 = $state("");
    let ai2 = $state("");
    let creating = $state(false);
    let error = $state("");

    const isAiVsAi = $derived(gameMode.mode === gamemodes.ai_vs_ai.mode);

    async function selectAi(aiId: string) {
        if (step === "ai1") {
            ai1 = aiId;
            if (gameMode.mode === gamemodes.human_vs_ai.mode) {
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
            const info = await games.create(gameMode.mode, ai1, ai2);
            goto(`/game/${info.gameId}?mode=${gameMode.mode}`);
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
    <title>GoStrategy — Select AI</title>
</svelte:head>

<div class="space-y-10 max-w-5xl mx-auto py-10 px-4">
    <header
        class="flex flex-col md:flex-row md:items-center justify-between gap-6"
    >
        <div class="space-y-1">
            <h1
                class="text-4xl font-black text-white uppercase tracking-tighter"
            >
                {#if isAiVsAi}
                    {step === "ai1" ? "First Combatant" : "Second Combatant"}
                {:else}
                    Choose your Foe
                {/if}
            </h1>
            <p
                class="text-white/40 font-medium uppercase tracking-widest text-[10px]"
            >
                {gameMode.title} Configuration
            </p>
        </div>
        <div class="flex items-center gap-3">
            <Button variant="ghost" onclick={back}>
                {step === "ai2" ? "Change First AI" : "Cancel"}
            </Button>
        </div>
    </header>

    <SelectionSummary {gameMode} {ai1} {ai2} {step} />

    {#if error}
        <div
            class="bg-brand-secondary/20 border border-brand-secondary/30 text-brand-secondary rounded-2xl px-6 py-4 text-sm font-medium text-center"
        >
            {error}
        </div>
    {/if}

    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        {#each AIs as ai}
            <button
                class="text-left p-6 rounded-2xl border-2 transition-all duration-300 flex items-center gap-6 group relative overflow-hidden {ai1 ===
                    ai.id || ai2 === ai.id
                    ? 'border-brand-primary bg-brand-primary/10 shadow-[0_0_30px_rgba(var(--brand-primary-rgb),0.1)]'
                    : 'border-white/5 bg-surface-elevated/20 hover:border-white/20 hover:bg-white/5'}"
                onclick={() => selectAi(ai.id)}
                disabled={creating}
            >
                {#if ai.image}
                    <div class="relative">
                        <img
                            src={ai.image}
                            alt={ai.name}
                            class="w-24 h-24 rounded-xl object-cover bg-black/40 border border-white/10 shadow-xl group-hover:scale-105 transition-transform duration-500"
                        />
                        <div
                            class="absolute inset-0 rounded-xl bg-linear-to-t from-black/60 to-transparent opacity-0 group-hover:opacity-100 transition-opacity"
                        ></div>
                    </div>
                {/if}
                <div class="flex-1 space-y-2">
                    <div class="flex items-center justify-between">
                        <h3
                            class="text-2xl font-black text-white uppercase tracking-tighter"
                        >
                            {ai.name}
                        </h3>
                        <div class="flex gap-2">
                            {#if ai1 === ai.id}
                                <span
                                    class="text-[9px] font-black bg-red-500 text-white px-2 py-1 rounded shadow-sm border border-white/10 uppercase tracking-wider"
                                    >P1</span
                                >
                            {/if}
                            {#if ai2 === ai.id}
                                <span
                                    class="text-[9px] font-black bg-blue-500 text-white px-2 py-1 rounded shadow-sm border border-white/10 uppercase tracking-wider"
                                    >P2</span
                                >
                            {/if}
                        </div>
                    </div>
                    <p
                        class="text-white/40 text-sm leading-relaxed font-medium"
                    >
                        {ai.description}
                    </p>
                </div>

                {#if (step === "ai1" && ai1 === ai.id) || (step === "ai2" && ai2 === ai.id)}
                    <div
                        class="absolute top-0 right-0 w-2 h-full bg-brand-primary"
                    ></div>
                {/if}
            </button>
        {/each}
    </div>

    {#if creating}
        <div
            class="fixed inset-0 bg-black/60 backdrop-blur-md z-50 flex flex-col items-center justify-center gap-6"
        >
            <div
                class="w-16 h-16 border-4 border-white/10 border-t-brand-primary rounded-full animate-spin"
            ></div>
            <div class="text-center space-y-2">
                <h2
                    class="text-2xl font-black text-white uppercase tracking-tighter"
                >
                    Initializing Engines
                </h2>
                <p
                    class="text-white/40 text-sm font-medium animate-pulse uppercase tracking-[0.2em]"
                >
                    Compiling strategic protocols...
                </p>
            </div>
        </div>
    {/if}
</div>
