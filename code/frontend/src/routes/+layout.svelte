<script lang="ts">
    import "./layout.css";
    import { page } from "$app/stores";

    let { children } = $props();

    const navItems = [
        { name: "Command Center", href: "/" },
        { name: "Battlefields", href: "/games" },
        { name: "Hall of Fame", href: "/leaderboard" },
        { name: "War Room", href: "/settings" },
    ];
</script>

<div class="min-h-screen flex bg-surface-base">
    <!-- Sidebar -->
    <aside
        class="w-64 border-r border-white/5 bg-surface-elevated/30 backdrop-blur-xl flex flex-col fixed inset-y-0"
    >
        <div class="p-8">
            <h1
                class="text-2xl font-extrabold tracking-widest uppercase text-white flex items-center gap-3 drop-shadow-md"
            >
                <div
                    class="relative w-8 h-8 rounded-md bg-brand-primary border-2 border-brand-accent flex items-center justify-center overflow-hidden shadow-[0_0_15px_rgba(var(--color-brand-primary),0.5)]"
                >
                    <span
                        class="relative z-10 text-brand-accent font-black text-lg"
                        >S</span
                    >
                    <div
                        class="absolute inset-0 bg-gradient-to-br from-white/20 to-transparent"
                    ></div>
                </div>
                Stratego
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
            <div class="flex items-center gap-3 px-4 py-2">
                <div
                    class="w-8 h-8 rounded-full bg-brand-secondary/30 border border-brand-accent/50 flex items-center justify-center text-brand-accent text-xs font-bold"
                >
                    ?
                </div>
                <div class="flex flex-col">
                    <span
                        class="text-xs font-bold text-white uppercase tracking-wider"
                        >Unknown Scout</span
                    >
                    <span class="text-[10px] text-white/50"
                        >Awaiting Orders</span
                    >
                </div>
            </div>
        </div>
    </aside>

    <!-- Main Content -->
    <main class="flex-1 ml-64 p-10">
        <div class="max-w-6xl mx-auto">
            {@render children()}
        </div>
    </main>
</div>
