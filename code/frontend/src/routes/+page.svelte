<script lang="ts">
    import { goto } from "$app/navigation";
    import { authStore } from "$lib/state/auth.svelte";
    import { games } from "$lib/api/client";
    import Card from "$lib/components/ui/Card.svelte";
    import Button from "$lib/components/ui/Button.svelte";
    import type { GameMode, AI } from "$lib/types/game";

    const AIs: AI[] = [
        {
            name: "FAFO",
            id: "fafo",
            description:
                "The Fuck Around & Find Out AI is a simple random-move AI.",
        },
        {
            name: "FATO",
            id: "fato",
            description:
                "The Fuck Around & Try Out AI is a random-move AI that can remember the board and act on it.",
        },
    ];

    let selectedMode = $state<GameMode>("human_vs_ai");
    let error = $state("");
    let creating = $state(false);

    // AI selection flow
    let selectingAi = $state<null | "ai1" | "ai2">(null);
    let ai1 = $state("");
    let ai2 = $state("");

    function startFlow() {
        error = "";
        if (selectedMode === "human_vs_human") {
            error = "Coming Soon";
            return;
        }
        selectingAi = "ai1";
        ai1 = "";
        ai2 = "";
    }

    function selectAi(aiId: string) {
        if (selectingAi === "ai1") {
            ai1 = aiId;
            if (selectedMode === "human_vs_ai") {
                createGame();
            } else {
                selectingAi = "ai2";
            }
        } else {
            ai2 = aiId;
            createGame();
        }
    }

    async function createGame() {
        creating = true;
        selectingAi = null;
        try {
            const info = await games.create(selectedMode, ai1, ai2);
            goto(`/game/${info.gameId}?mode=${selectedMode}`);
        } catch (e: any) {
            error = e.message || "Failed to create game";
            creating = false;
        }
    }

    function cancel() {
        selectingAi = null;
        ai1 = "";
        ai2 = "";
    }

    const modes: {
        mode: GameMode;
        icon: string;
        title: string;
        desc: string;
    }[] = [
        {
            mode: "human_vs_ai",
            icon: "🧑 vs 🤖",
            title: "Human vs AI",
            desc: "Play against the computer.",
        },
        {
            mode: "ai_vs_ai",
            icon: "🤖 vs 🤖",
            title: "AI vs AI",
            desc: "Watch two AIs battle it out.",
        },
        {
            mode: "human_vs_human",
            icon: "🧑 vs 🧑",
            title: "Human vs Human",
            desc: "Coming soon.",
        },
    ];
</script>

<svelte:head>
    <title>Stratego — Command Center</title>
</svelte:head>

<div class="space-y-8">
    <header>
        <h1
            class="text-3xl font-extrabold text-white uppercase tracking-widest"
        >
            Command Center
        </h1>
        <p class="text-white/50 mt-1">
            {#if authStore.user}
                Welcome back, <span class="text-brand-accent font-semibold"
                    >{authStore.user.username}</span
                >.
            {:else}
                <a href="/login" class="text-brand-primary hover:underline"
                    >Sign in</a
                > to start playing.
            {/if}
        </p>
    </header>

    {#if error}
        <div
            class="bg-brand-secondary/20 border border-brand-secondary/30 text-brand-secondary rounded-xl px-4 py-3 text-sm text-center"
        >
            {error}
        </div>
    {/if}

    <!-- AI Selector Dialog -->
    {#if selectingAi}
        <Card class="space-y-4">
            <div class="flex justify-between items-center">
                <h2
                    class="text-lg font-bold text-brand-accent uppercase tracking-wider"
                >
                    Select {selectingAi === "ai1" ? "AI Opponent" : "Second AI"}
                </h2>
                <Button variant="ghost" size="sm" onclick={cancel}
                    >✕ Cancel</Button
                >
            </div>
            <div class="grid gap-3">
                {#each AIs as ai}
                    <button
                        class="text-left p-4 rounded-xl border border-white/10 bg-white/5 hover:bg-brand-primary/10 hover:border-brand-primary/30 transition-all"
                        onclick={() => selectAi(ai.id)}
                    >
                        <h3 class="font-bold text-white">{ai.name}</h3>
                        <p class="text-white/50 text-sm mt-1">
                            {ai.description}
                        </p>
                    </button>
                {/each}
            </div>
        </Card>
    {:else}
        <!-- Game Mode Selection -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
            {#each modes as m}
                <button
                    class="text-left rounded-2xl p-6 border-2 transition-all duration-200 {selectedMode ===
                    m.mode
                        ? 'border-brand-primary bg-brand-primary/10'
                        : 'border-white/10 bg-surface-glass hover:border-white/20 hover:bg-white/5'}"
                    onclick={() => (selectedMode = m.mode)}
                    disabled={creating}
                >
                    <div class="text-3xl mb-3">{m.icon}</div>
                    <h3 class="font-bold text-white text-lg">{m.title}</h3>
                    <p class="text-white/50 text-sm mt-1">{m.desc}</p>
                </button>
            {/each}
        </div>

        <Button
            variant="primary"
            size="lg"
            class="w-full"
            onclick={startFlow}
            disabled={creating || !authStore.user}
        >
            {#if creating}
                Creating Game...
            {:else if !authStore.user}
                Sign In to Play
            {:else}
                Start Game
            {/if}
        </Button>
    {/if}
</div>
