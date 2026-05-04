<script lang="ts">
    import Button from "$lib/components/ui/Button.svelte";
    import type { BoardSetup } from "$lib/types/board-setup";
    import { gameStore } from "$lib/state/game.svelte";
    import BoardSetupCard from "../setup/BoardSetupCard.svelte";

    interface Props {
        savedSetups: BoardSetup[];
        ownerId: number;
        selectedPlayer: number;
        onSelectSetup: (setupData: string) => void;
        showSelector: boolean;
    }

    let {
        savedSetups,
        ownerId,
        selectedPlayer,
        showSelector = $bindable(),
        onSelectSetup,
    }: Props = $props();
</script>

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
        class:border-brand-primary={ownerId === 1}
        class:border-brand-secondary={ownerId === 2}
    >
        <div
            class="px-8 py-6 border-b border-white/10 flex justify-between items-center bg-white/5"
        >
            <div>
                <h2
                    class="text-2xl font-black text-white uppercase tracking-tighter"
                >
                    Formation for {selectedPlayer === 0
                        ? gameStore.gameState?.player1Username
                        : gameStore.gameState?.player2Username}
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
                        onclick={() => onSelectSetup(setup.setup_data)}
                    />
                {/each}
            </div>
        </div>
    </div>
</div>
