<script lang="ts">
    import favicon from "$lib/assets/favicon.png";
	import "$lib/colors.css";
    import "$lib/main.css";

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

    :global(html, body) {
        height: 100%;
        margin: 0;
        padding: 0;
    }

    :global(*) {
        box-sizing: border-box;
    }
</style>
