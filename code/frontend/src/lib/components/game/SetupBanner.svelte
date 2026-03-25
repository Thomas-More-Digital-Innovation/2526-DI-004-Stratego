<script lang="ts">
    import { onMount } from "svelte";
    import Button from "$lib/components/ui/Button.svelte";
    import { boardSetups } from "$lib/api/client";
    import type { BoardSetup } from "$lib/types/board-setup";
    import BoardSetupCard from "$lib/components/setup/BoardSetupCard.svelte";

    interface Props {
        onRandomize: () => void;
        onStart: () => void;
        onLoadSetup: (setupData: string) => void;
        viewerId?: number;
    }

    let { onRandomize, onStart, onLoadSetup, viewerId = 0 }: Props = $props();

    const ownerId = $derived(viewerId === -1 ? 1 : viewerId + 1);

    let savedSetups = $state<BoardSetup[]>([]);
    let loadingSetups = $state(true);
    let showSelector = $state(false);

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
        onLoadSetup(setupData);
        showSelector = false;
    }
</script>

<div class="fixed top-0 left-0 right-0 z-50 pointer-events-none">
    <div
        class="glass pointer-events-auto flex items-center justify-between gap-6 px-8 py-4 border-b border-white/10"
    >
        <div>
            <h2
                class="text-lg font-bold text-brand-accent uppercase tracking-wider"
            >
                ⚔️ Setup Phase
            </h2>
            <p class="text-white/50 text-sm">
                Arrange your pieces or load a saved configuration
            </p>
        </div>

        <div class="flex gap-3 items-center">
            <Button
                variant="ghost"
                onclick={() => (showSelector = true)}
                class="gap-3 bg-white/5! hover:bg-white/10!"
                disabled={loadingSetups || savedSetups.length === 0}
            >
                📁 Load Saved Setup
            </Button>

            <Button variant="outline" onclick={onRandomize}>
                🎲 Randomize
            </Button>
            <Button variant="primary" onclick={onStart}>▶️ Start Game</Button>
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
