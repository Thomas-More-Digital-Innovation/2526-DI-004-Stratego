<script lang="ts">
    import "./layout.css";
    import { page } from "$app/state";
    import logo from "$lib/assets/favicon.png";
    import baseBg from "$lib/assets/background.webp";
    import profileBg from "$lib/assets/background-profile.webp";
    import boardBg from "$lib/assets/background-board-setup.webp";
    import EarlyAccessDisclaimer from "$lib/components/EarlyAccessDisclaimer.svelte";
    import SideBar from "$lib/components/sidebar/SideBar.svelte";
    import Toaster from "$lib/components/ui/Toaster.svelte";

    let { children } = $props();

    const isFullPage = $derived(page.url.pathname.startsWith("/game/"));

    const backgroundImage = $derived(() => {
        if (page.url.pathname.startsWith("/profile")) return profileBg;
        if (page.url.pathname.startsWith("/board-setups")) return boardBg;
        return baseBg;
    });
</script>

<svelte:head>
    <link rel="icon" href={logo} />
</svelte:head>

<div
    class="min-h-screen flex bg-surface-base/20 relative overflow-hidden"
    style:--bg-image="url({backgroundImage()})"
>
    <!-- Dynamic Background Layer -->
    <div class="app-background">
        <div class="app-background-overlay"></div>
    </div>

    {#if !isFullPage}
        <SideBar />
    {/if}

    <main class="flex-1 {isFullPage ? '' : 'ml-64'} p-10">
        <div class="max-w-6xl mx-auto">
            {@render children()}
        </div>
    </main>

    <EarlyAccessDisclaimer />
    <Toaster />
</div>
