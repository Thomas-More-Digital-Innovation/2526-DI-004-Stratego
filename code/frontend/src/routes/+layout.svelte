<script lang="ts">
    import favicon from "$lib/assets/favicon.png";
	import "$lib/colors.css";

    let { children } = $props();
</script>

<svelte:head>
    <link rel="icon" href={favicon} />
</svelte:head>

{@render children?.()}

<style>
    :global(.tile) {
        background: var(--bg);
        padding: 32px;
        border-radius: 12px;
		margin: 16px;
        box-shadow: 0 0 48px var(--bg-accent);
    }

    :global(*) {
        box-sizing: border-box;
    }
</style>
