<script lang="ts">
    import Button from "$lib/components/ui/Button.svelte";

    interface Props {
        onRandomize: () => void;
        onStart: () => void;
    }

    let { onRandomize, onStart }: Props = $props();
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
                Click two pieces in the bottom 4 rows to swap them
            </p>
        </div>

        <div class="flex gap-3">
            <Button variant="outline" onclick={onRandomize}>
                🎲 Randomize
            </Button>
            <Button variant="primary" onclick={onStart}>▶️ Start Game</Button>
        </div>
    </div>
</div>
