<script lang="ts">
    import { onMount } from "svelte";
    import Button from "$lib/components/ui/Button.svelte";
    import { boardSetups } from "$lib/api/client";
    import type { BoardSetup } from "$lib/types/board-setup";

    interface Props {
        onRandomize: () => void;
        onStart: () => void;
        onLoadSetup: (setupData: string) => void;
    }

    let { onRandomize, onStart, onLoadSetup }: Props = $props();

    let savedSetups = $state<BoardSetup[]>([]);
    let loadingSetups = $state(true);

    onMount(async () => {
        try {
            const result = await boardSetups.list();
            savedSetups = result ?? [];
        } catch {
            // Silently fail - user can still randomize
        } finally {
            loadingSetups = false;
        }
    });
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
                Click two pieces to swap, or load a saved setup
            </p>
        </div>

        <div class="flex gap-3 items-center">
            {#if !loadingSetups && savedSetups.length > 0}
                <select
                    class="bg-white/5 border border-white/10 rounded-xl px-3 py-2 text-sm text-white focus:outline-none focus:border-brand-accent/50 appearance-none cursor-pointer"
                    onchange={(e) => {
                        const val = (e.target as HTMLSelectElement).value;
                        if (val) {
                            onLoadSetup(val);
                            (e.target as HTMLSelectElement).value = "";
                        }
                    }}
                >
                    <option value="" class="bg-black"
                        >📁 Load Saved Setup</option
                    >
                    {#each savedSetups as setup}
                        <option value={setup.setup_data} class="bg-black">
                            {setup.name}
                            {setup.is_default ? "⭐" : ""}
                        </option>
                    {/each}
                </select>
            {/if}

            <Button variant="outline" onclick={onRandomize}>
                🎲 Randomize
            </Button>
            <Button variant="primary" onclick={onStart}>▶️ Start Game</Button>
        </div>
    </div>
</div>
