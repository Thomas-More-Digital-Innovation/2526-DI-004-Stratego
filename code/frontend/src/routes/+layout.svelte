<script lang="ts">
    import "./layout.css";
    import { page } from "$app/stores";
    import { onMount } from "svelte";
    import { authStore } from "$lib/state/auth.svelte";
    import logo from "$lib/assets/favicon.png";
    import baseBg from "$lib/assets/background.png";
    import profileBg from "$lib/assets/background-profile.png";
    import boardBg from "$lib/assets/background-board-setup.png";

    let { children } = $props();

    onMount(() => {
        authStore.check();
    });

    const navItems = [
        { name: "Command Center", href: "/" },
        { name: "Profile", href: "/profile" },
        { name: "Board Setups", href: "/board-setups" },
    ];

    const isFullPage = $derived($page.url.pathname.startsWith("/game/"));

    const backgroundImage = $derived(() => {
        if ($page.url.pathname.startsWith("/profile")) return profileBg;
        if ($page.url.pathname.startsWith("/board-setups")) return boardBg;
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
        <aside
            class="w-64 border-r border-white/5 bg-surface-elevated/30 backdrop-blur-xl flex flex-col fixed inset-y-0"
        >
            <div class="px-4 py-8">
                <h1
                    class="text-2xl font-extrabold tracking-widest uppercase text-white flex items-center gap-3 drop-shadow-md"
                >
                    <img src={logo} alt="Logo" class="w-12 h-12" />Stratego
                </h1>
            </div>

            <nav class="flex-1 px-4 space-y-1">
                {#each navItems as item}
                    <a
                        href={item.href}
                        class="flex items-center px-4 py-3 text-sm font-medium rounded-xl transition-all duration-200 group {item.href ===
                        $page.url.pathname
                            ? 'bg-brand-primary/10 text-brand-primary'
                            : 'text-white/50 hover:bg-white/5 hover:text-white'}"
                    >
                        {item.name}
                    </a>
                {/each}
            </nav>

            <div class="p-4 border-t border-brand-accent/20 bg-black/20">
                {#if authStore.user}
                    <div class="flex items-center gap-3 px-4 py-2">
                        <div
                            class="w-8 h-8 rounded-full bg-brand-secondary/30 border border-brand-accent/50 flex items-center justify-center text-brand-accent text-xs font-bold"
                        >
                            {authStore.user.username[0]?.toUpperCase() || "?"}
                        </div>
                        <div class="flex flex-col">
                            <span
                                class="text-xs font-bold text-white uppercase tracking-wider"
                                >{authStore.user.username}</span
                            >
                            <span class="text-[10px] text-white/50">Online</span
                            >
                        </div>
                    </div>
                {:else}
                    <a
                        href="/login"
                        class="flex items-center gap-3 px-4 py-2 text-white/50 hover:text-white transition-colors"
                    >
                        <div
                            class="w-8 h-8 rounded-full bg-white/5 border border-white/10 flex items-center justify-center text-white/30 text-xs font-bold"
                        >
                            ?
                        </div>
                        <span class="text-xs font-bold uppercase tracking-wider"
                            >Sign In</span
                        >
                    </a>
                {/if}
            </div>
        </aside>
    {/if}

    <main class="flex-1 {isFullPage ? '' : 'ml-64'} p-10">
        <div class="max-w-6xl mx-auto">
            {@render children()}
        </div>
    </main>
</div>
