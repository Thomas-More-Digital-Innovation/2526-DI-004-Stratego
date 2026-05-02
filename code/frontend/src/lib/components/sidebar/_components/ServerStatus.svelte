<script lang="ts">
    import { onMount } from "svelte";
    import { serverStore } from "$lib/state/server.svelte";

    onMount(() => {
        serverStore.startPolling(5000);
        return () => serverStore.stopPolling();
    });
</script>

<div class="flex items-center gap-1.5">
    <span
        class="text-[10px] {serverStore.isOnline
            ? 'text-white/50'
            : 'text-red-400 font-medium'}"
    >
        {serverStore.isOnline ? "Server Online" : "Server Offline"}
    </span>
    <div
        class="w-1.5 h-1.5 rounded-full {serverStore.isOnline
            ? 'bg-green-500 shadow-[0_0_8px_rgba(34,197,94,0.6)]'
            : 'bg-red-500 shadow-[0_0_8px_rgba(239,68,68,0.6)]'}"
    ></div>
</div>
