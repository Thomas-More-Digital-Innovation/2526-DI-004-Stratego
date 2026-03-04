<script lang="ts">
    import "./layout.css";
    import { page } from "$app/stores";

    let { children } = $props();

    const navItems = [
        { name: "Dashboard", href: "/" },
        { name: "Games", href: "/games" },
        { name: "Leaderboard", href: "/leaderboard" },
        { name: "Settings", href: "/settings" },
    ];
</script>

<div class="min-h-screen flex bg-surface-base">
    <!-- Sidebar -->
    <aside
        class="w-64 border-r border-white/5 bg-surface-elevated/30 backdrop-blur-xl flex flex-col fixed inset-y-0"
    >
        <div class="p-8">
            <h1
                class="text-2xl font-bold tracking-tight text-white flex items-center gap-2"
            >
                <span
                    class="w-8 h-8 rounded-lg bg-brand-primary glow-primary flex items-center justify-center text-sm"
                    >S</span
                >
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

        <div class="p-4 border-t border-white/5">
            <div class="flex items-center gap-3 px-4 py-2">
                <div
                    class="w-8 h-8 rounded-full bg-brand-secondary/20 border border-brand-secondary/30"
                ></div>
                <div class="flex flex-col">
                    <span class="text-xs font-semibold text-white"
                        >Guest User</span
                    >
                    <span class="text-[10px] text-white/40">Not logged in</span>
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
