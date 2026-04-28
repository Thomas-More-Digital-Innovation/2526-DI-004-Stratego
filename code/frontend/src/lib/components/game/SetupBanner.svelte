<script lang="ts">
    import { onMount } from "svelte";
    import Button from "$lib/components/ui/Button.svelte";
    import { boardSetups } from "$lib/api/client";
    import { flipSetup } from "$lib/utils/board-binary";
    import type { BoardSetup } from "$lib/types/board-setup";
    import BoardSetupCard from "$lib/components/setup/BoardSetupCard.svelte";
    import type { GameMode } from "$lib/types/game";
    import { gamemodes } from "$lib/data/gamemodes.data";

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
            ? selectedPlayer + 1
            : viewerId === -1
              ? 1
              : viewerId + 1,
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
</script>

<div class="fixed top-0 left-0 right-0 z-50 pointer-events-none">
    <div
        class="glass pointer-events-auto flex items-center justify-between gap-6 px-8 py-4 border-b border-white/10"
    >
        <div class="flex items-center gap-3">
            <Button
                variant="outline"
                onclick={() => (window.location.href = "/")}
            >
                ◀️ Back To Menu
            </Button>

            <div>
                <h2
                    class="text-lg font-bold text-brand-accent uppercase tracking-wider"
                >
                    ⚔️ Setup Phase
                </h2>
                {#if gameMode.mode === gamemodes.ai_vs_ai.mode}
                    <div class="flex gap-2 mt-1">
                        <Button
                            variant={selectedPlayer === 0 ? "primary" : "ghost"}
                            size="sm"
                            onclick={() => onSelectPlayer?.(0)}
                            class="py-1! h-auto! text-[10px]!"
                        >
                            AI Red
                        </Button>
                        <Button
                            variant={selectedPlayer === 1 ? "primary" : "ghost"}
                            size="sm"
                            onclick={() => onSelectPlayer?.(1)}
                            class="py-1! h-auto! text-[10px]!"
                        >
                            AI Blue
                        </Button>
                    </div>
                {:else}
                    <p class="text-white/50 text-sm">
                        Arrange your pieces or load a saved configuration
                    </p>
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
                >▶️ Start Game</Button
            >
        </div>
    </div>
</div>

{#if showSelector}
    <div
        class="fixed inset-0 z-60 flex items-center justify-center p-6 animate-in fade-in duration-300"
    >
        <!-- Backdrop -->
        <!-- svelte-ignore a11y_click_events_have_key_events -->
        <!-- svelte-ignore a11y_no_static_element_interactions -->
        <div
            class="absolute inset-0 bg-black/60 backdrop-blur-md"
            onclick={() => (showSelector = false)}
        ></div>

        <!-- Modal -->
        <div
            class="relative w-full max-w-5xl max-h-[85vh] flex flex-col glass rounded-3xl shadow-2xl overflow-hidden border-2"
            class:border-brand-primary={ownerId !== 1}
            class:border-brand-secondary={ownerId === 1}
        >
            <div
                class="px-8 py-6 border-b border-white/10 flex justify-between items-center bg-white/5"
            >
                <div>
                    <h2
                        class="text-2xl font-black text-white uppercase tracking-tighter"
                    >
                        Choose your formation
                    </h2>
                    <p class="text-white/40 text-sm mt-1">
                        Select one of your {savedSetups.length} saved board configurations
                    </p>
                </div>
                <Button
                    variant="ghost"
                    onclick={() => (showSelector = false)}
                    class="size-10! p-0! rounded-full hover:bg-white/10"
                >
                    ✕
                </Button>
            </div>

            <div class="flex-1 overflow-y-auto p-8 bg-black/20">
                <div class="flex justify-center flex-wrap gap-6">
                    {#each savedSetups as setup}
                        <BoardSetupCard
                            {setup}
                            {ownerId}
                            isInteractive={true}
                            onclick={() => selectSetup(setup.setup_data)}
                        />
                    {/each}
                </div>
            </div>
        </div>
    </div>
{/if}
