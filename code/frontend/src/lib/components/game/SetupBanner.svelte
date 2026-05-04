<script lang="ts">
    import { onMount } from "svelte";
    import Button from "$lib/components/ui/Button.svelte";
    import { boardSetups } from "$lib/api/client";
    import { flipSetup } from "$lib/utils/board-binary";
    import type { BoardSetup } from "$lib/types/board-setup";
    import BoardSetupCard from "$lib/components/setup/BoardSetupCard.svelte";
    import type { GameMode } from "$lib/types/game";
    import { gamemodes } from "$lib/data/gamemodes.data";
    import { gameStore } from "$lib/state/game.svelte";
    import LoadSavedSetup from "./LoadSavedSetup.svelte";

    interface Props {
        onRandomize: (player?: number) => void;
        onStart: (headless?: boolean) => void;
        onLoadSetup: (setupData: string, player?: number) => void;
        viewerId?: number;
        gameMode?: GameMode;
        selectedPlayer?: number;
        onSelectPlayer?: (player: number) => void;
    }

    let {
        onRandomize,
        onStart,
        onLoadSetup,
        viewerId = 0,
        gameMode = gamemodes.human_vs_ai,
        selectedPlayer = 0,
        onSelectPlayer,
    }: Props = $props();

    const ownerId = $derived(
        gameMode.mode === gamemodes.ai_vs_ai.mode
            ? (selectedPlayer === 0 ? 2 : 1)
            : viewerId === -1
              ? 2
              : (viewerId === 0 ? 2 : 1),
    );

    let savedSetups = $state<BoardSetup[]>([]);
    let loadingSetups = $state(true);
    let showSelector = $state(false);
    let headless = $state(false);

    onMount(async () => {
        try {
            const result = await boardSetups.list();
            savedSetups = result ?? [];
        } catch {
            // Silently fail
        } finally {
            loadingSetups = false;
        }
    });

    function selectSetup(setupData: string) {
        let finalSetup = setupData;
        if (gameMode.mode === gamemodes.ai_vs_ai.mode && selectedPlayer === 1) {
            finalSetup = flipSetup(setupData);
        }
        onLoadSetup(
            finalSetup,
            gameMode.mode === gamemodes.ai_vs_ai.mode
                ? selectedPlayer
                : undefined,
        );
        showSelector = false;
    }

    function formatTime(seconds: number | null): string {
        if (seconds === null) return "00:00";
        const mins = Math.floor(seconds / 60);
        const secs = seconds % 60;
        return `${mins}:${secs.toString().padStart(2, "0")}`;
    }
</script>

{#snippet TimeLimit()}
    <div class="group cursor-help relative">
        <p
            class="text-white/20 text-[9px] font-bold uppercase tracking-[0.2em] flex items-center gap-1 w-fit"
        >
            ⏳ TIME LIMIT: {formatTime(gameStore.setupRemainingSecs)}
        </p>
        <div
            class="hidden group-hover:block absolute top-full left-0 mt-2 bg-surface-elevated rounded-lg p-4 text-xs w-64 text-white shadow-[0_10px_30px_rgba(0,0,0,0.5)] border border-white/10 z-50 animate-in fade-in slide-in-from-top-1 duration-200"
        >
            A time limit is set for this game to prevent infinite setup times.
            The game will start automatically when the time runs out.
        </div>
    </div>
{/snippet}

<div class="fixed top-0 left-0 right-0 z-50 pointer-events-none">
    <div
        class="glass pointer-events-auto flex items-center justify-between gap-6 px-8 py-4 border-b border-white/10"
    >
        <div class="flex items-center gap-3">
            <Button
                variant="outline"
                onclick={() => (window.location.href = "/")}
            >
                Back To Menu
            </Button>

            <div>
                <div class="flex gap-4 items-center">
                    <h2
                        class="text-lg font-black text-white uppercase tracking-tighter"
                    >
                        Setup Phase
                    </h2>
                    {@render TimeLimit()}
                </div>

                {#if gameMode.mode === gamemodes.ai_vs_ai.mode}
                    <div class="flex items-center gap-3 mt-2">
                        <button
                            class="flex items-center gap-2 px-3 py-1.5 rounded-lg border-2 transition-all {selectedPlayer ===
                            0
                                ? 'border-red-500 bg-red-500/10 shadow-[0_0_15px_rgba(239,68,68,0.2)]'
                                : 'border-white/5 bg-white/5 hover:border-white/20 opacity-40 hover:opacity-100'}"
                            onclick={() => onSelectPlayer?.(0)}
                        >
                            <div
                                class="w-2 h-2 rounded-full bg-red-500 shadow-[0_0_5px_rgba(239,68,68,1)]"
                            ></div>
                            <span
                                class="text-[10px] font-black text-white uppercase tracking-wider"
                            >
                                {gameStore.gameState?.player1Username ||
                                    "AI Red"}
                            </span>
                        </button>
                        <div
                            class="text-[10px] font-black text-white/10 uppercase italic"
                        >
                            vs
                        </div>
                        <button
                            class="flex items-center gap-2 px-3 py-1.5 rounded-lg border-2 transition-all {selectedPlayer ===
                            1
                                ? 'border-blue-500 bg-blue-500/10 shadow-[0_0_15px_rgba(59,130,246,0.2)]'
                                : 'border-white/5 bg-white/5 hover:border-white/20 opacity-40 hover:opacity-100'}"
                            onclick={() => onSelectPlayer?.(1)}
                        >
                            <span
                                class="text-[10px] font-black text-white uppercase tracking-wider"
                            >
                                {gameStore.gameState?.player2Username ||
                                    "AI Blue"}
                            </span>
                            <div
                                class="w-2 h-2 rounded-full bg-blue-500 shadow-[0_0_5px_rgba(59,130,246,1)]"
                            ></div>
                        </button>
                    </div>
                {:else}
                    <div class="flex flex-col gap-0.5">
                        <p class="text-white/40 text-xs font-medium">
                            Arrange your pieces or load a formation
                        </p>
                    </div>
                {/if}
            </div>
        </div>

        <div class="flex gap-4 items-center">
            {#if gameMode.mode === gamemodes.ai_vs_ai.mode}
                <div
                    class="flex items-center gap-2 px-3 py-2 bg-white/5 rounded-xl border border-white/5"
                >
                    <input
                        type="checkbox"
                        id="headless-mode"
                        bind:checked={headless}
                        class="accent-brand-primary cursor-pointer"
                    />
                    <label
                        for="headless-mode"
                        class="text-[10px] font-bold text-white/60 uppercase tracking-tighter cursor-pointer"
                    >
                        Headless Mode
                    </label>
                </div>
            {/if}
            <Button
                variant="ghost"
                onclick={() => (showSelector = true)}
                class="gap-3 bg-white/5! hover:bg-white/10!"
                disabled={loadingSetups || savedSetups.length === 0}
            >
                📁 Load Saved Setup
            </Button>

            <Button
                variant="outline"
                onclick={() =>
                    onRandomize(
                        gameMode.mode === gamemodes.ai_vs_ai.mode
                            ? selectedPlayer
                            : undefined,
                    )}
            >
                🎲 Randomize
            </Button>
            <Button variant="primary" onclick={() => onStart(headless)}
                >Start Game</Button
            >
        </div>
    </div>
</div>

{#if showSelector}
    <LoadSavedSetup
        {savedSetups}
        {ownerId}
        {selectedPlayer}
        onSelectSetup={selectSetup}
        bind:showSelector
    />
{/if}
