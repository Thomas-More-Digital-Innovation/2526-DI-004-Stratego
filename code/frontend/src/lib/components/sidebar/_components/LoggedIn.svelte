<script lang="ts">
    import { authStore } from "$lib/state/auth.svelte";
    import ServerStatus from "./ServerStatus.svelte";
</script>

<a
    href="/profile"
    class="flex items-center gap-3 px-4 py-2 text-white hover:text-brand-secondary transition-colors"
>
    <div
        class="w-8 h-8 rounded-full bg-brand-secondary/30 border border-brand-accent/50 flex items-center justify-center text-brand-accent text-xs font-bold"
    >
        {authStore.user?.username[0]?.toUpperCase() || "?"}
    </div>
    <div class="flex flex-col">
        <span class="text-xs font-bold uppercase tracking-wider"
            >{authStore.user?.username}</span
        >
        <ServerStatus />
    </div>
</a>
