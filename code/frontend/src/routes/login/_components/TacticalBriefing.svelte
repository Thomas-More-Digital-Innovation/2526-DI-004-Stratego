<script lang="ts">
    import Title from "$lib/components/Title.svelte";
    import GridOverlay from "$lib/components/GridOverlay.svelte";
    import Commander from "./Commander.svelte";

    let { isBlueCommander = $bindable(true) } = $props();

    function switchCommander() {
        isBlueCommander = !isBlueCommander;
    }
</script>

<div
    class="hidden md:flex md:col-span-5 bg-linear-to-br from-black/60 to-surface-base/30 p-8 border-r border-white/5 flex-col justify-between relative overflow-hidden select-none"
>
    <!-- Grid Overlay -->
    <GridOverlay opacity={5} size={20} color="white" />

    <!-- Briefing Header -->
    <div class="space-y-6 relative z-10">
        <div class="space-y-1">
            <Title size="2xl" />
            <p class="text-white/60 text-xs leading-relaxed font-sans">
                Step into a high-stakes arena of tactical warfare, hidden
                configurations, and battlefield calculations. Deploy your squad,
                outsmart the enemy commander, and capture their flag.
            </p>
        </div>
    </div>

    <!-- Command Badge SVG -->
    <Commander bind:isBlueCommander />

    <!-- Tactical Quote -->
    <div class="space-y-4 relative z-10">
        <div
            class="border-l-2 border-brand-secondary/40 pl-3 py-1 italic font-serif text-white/50 text-[11px]"
        >
            "Tactics is the art of using soldiers in battle; strategy is the art
            of using battles to win the war."
            <span
                class="block mt-1 font-mono text-[9px] text-brand-accent/60 uppercase tracking-widest not-italic"
                >— Carl von Clausewitz</span
            >
        </div>
    </div>
</div>
