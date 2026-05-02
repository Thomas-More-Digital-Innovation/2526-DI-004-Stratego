<script>
    import ServerStatus from "./ServerStatus.svelte";
</script>

<a
    href="/login"
    class="flex items-center gap-3 px-4 py-2 text-white/50 hover:text-white transition-colors"
>
    <div
        class="w-8 h-8 rounded-full bg-white/5 border border-white/10 flex items-center justify-center text-white/30 text-xs font-bold"
    >
        ?
    </div>
    <div class="flex flex-col">
        <span class="text-xs font-bold uppercase tracking-wider">Sign In</span>
        <ServerStatus />
    </div>
</a>
