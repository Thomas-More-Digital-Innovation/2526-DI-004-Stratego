<script lang="ts">
    import { page } from "$app/state";
    import logoWebp from "$lib/assets/favicon.webp";
    import { onMount } from "svelte";
    import { authStore } from "$lib/state/auth.svelte";
    import LoggedOut from "./_components/LoggedOut.svelte";
    import LoggedIn from "./_components/LoggedIn.svelte";
    import goLogo from "$lib/assets/go.svg";

    const navItems = [
        { name: "Command Center", href: "/" },
        { name: "Profile", href: "/profile" },
        { name: "Board Setups", href: "/board-setups" },
    ];

    onMount(async () => {
        await authStore.check();
    });
</script>

<aside
    class="w-64 border-r border-white/5 bg-surface-elevated/30 backdrop-blur-xl flex flex-col fixed inset-y-0"
>
    <div class="px-4 py-8">
        <h1
            class="text-xl font-extrabold tracking-widest text-white flex gap-1 items-center drop-shadow-md"
        >
            <img src={logoWebp} alt="Logo" class="w-12 h-12" />
            <span class="flex items-center"
                ><a href="https://go.dev/">
                    <img src={goLogo} alt="go" class="w-10 h-10 mr-1" /></a
                >STRATEGY</span
            >
        </h1>
    </div>

    <nav class="flex-1 px-4 space-y-1">
        {#each navItems as item}
            <a
                href={item.href}
                class="flex items-center px-4 py-3 text-sm font-medium rounded-xl transition-all duration-200 group {item.href ===
                page.url.pathname
                    ? 'bg-brand-primary/10 text-brand-primary'
                    : 'text-white/50 hover:bg-white/5 hover:text-white'}"
            >
                {item.name}
            </a>
        {/each}
    </nav>

    <div class="group p-4 border-t border-brand-accent/20 bg-black/20">
        {#if authStore.loading}
            Who are you?
        {:else if authStore.user}
            <LoggedIn />
        {:else}
            <LoggedOut />
        {/if}
    </div>
</aside>
