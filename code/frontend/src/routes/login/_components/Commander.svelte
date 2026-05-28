<script lang="ts">
    import CommanderBlue from "$lib/assets/pieces/blue/10.webp";
    import CommanderRed from "$lib/assets/pieces/red/10.webp";

    let { isBlueCommander = $bindable(true) } = $props();

    const activeCommander = $derived(
        isBlueCommander ? CommanderBlue : CommanderRed,
    );

    function switchCommander() {
        isBlueCommander = !isBlueCommander;
    }
</script>

<div
    class="flex justify-center items-center py-6 relative z-10 opacity-70 hover:opacity-95 transition-opacity duration-500"
>
    <button
        onclick={switchCommander}
        class="w-32 h-32 rounded-xl cursor-pointer hover:scale-105 transition-transform duration-500"
    >
        <img
            src={activeCommander}
            alt="Commander"
            class="w-32 h-32 rounded-xl"
        />
    </button>
</div>
