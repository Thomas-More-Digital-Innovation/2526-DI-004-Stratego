<script lang="ts">
    import { goto } from "$app/navigation";
    import { authStore } from "$lib/state/auth.svelte";
    import Card from "$lib/components/ui/Card.svelte";
    import Button from "$lib/components/ui/Button.svelte";
    import type { GameMode } from "$lib/types/game";
    import { gamemodes } from "$lib/data/gamemodes.data";

    let selectedMode = $state<GameMode>("human_vs_ai");
    let error = $state("");

    function startFlow() {
        error = "";
        if (selectedMode === "human_vs_human") {
            error = "Coming Soon";
            return;
        }
        goto(`/select-ai?mode=${selectedMode}`);
    }
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

    {#if selectedMode}
        <!-- Game Mode Selection -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
            {#each gamemodes as m}
                <button
                    class="text-left disabled:opacity-50 disabled:cursor-not-allowed rounded-2xl p-6 border-2 transition-all duration-200 {selectedMode ===
                    m.mode
                        ? 'border-brand-primary bg-brand-primary/10 scale-102 shadow-glow'
                        : 'border-white/5 bg-surface-elevated/20 hover:border-white/20 hover:bg-white/5'}"
                    onclick={() => (selectedMode = m.mode)}
                    disabled={m.disabled}
                >
                    <div class="text-3xl mb-3">{m.icon}</div>
                    <h3
                        class="font-bold text-white text-lg lowercase tracking-widest"
                    >
                        {m.title}
                    </h3>
                    <p
                        class="text-white/50 text-xs mt-1 leading-relaxed italic"
                    >
                        {m.desc}
                    </p>
                </button>
            {/each}
        </div>

        <Button
            variant="primary"
            size="lg"
            class="w-full"
            onclick={startFlow}
            disabled={!authStore.user}
        >
            {#if !authStore.user}
                Sign In to Play
            {:else}
                Proceed to Intelligence Selection
            {/if}
        </Button>
    {/if}
</div>
