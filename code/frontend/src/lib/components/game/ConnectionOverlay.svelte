<script lang="ts">
    import Button from "$lib/components/ui/Button.svelte";

    interface Props {
        isReconnecting: boolean;
        gameId: string;
        reconnectAttempts: number;
        maxReconnectAttempts: number;
        onRetry?: () => void;
        onReturnToMenu?: () => void;
    }

    let {
        isReconnecting,
        gameId,
        reconnectAttempts,
        maxReconnectAttempts,
        onRetry = () => window.location.reload(),
        onReturnToMenu = () => (window.location.href = "/"),
    }: Props = $props();
</script>

<div class="flex flex-col items-center justify-center min-h-[60vh] gap-6 max-w-md mx-auto text-center px-4">
    {#if isReconnecting}
        <div class="relative flex items-center justify-center w-20 h-20">
            <div class="absolute inset-0 rounded-full border-4 border-brand-accent/20 border-t-brand-accent animate-spin"></div>
            <div class="text-brand-accent text-3xl animate-pulse">⚡</div>
        </div>
        <div class="space-y-2">
            <h2 class="text-2xl font-black text-white tracking-wider uppercase">Connection Lost</h2>
            <p class="text-white/60 text-sm">
                Reconnecting to game session <span class="font-mono text-brand-accent font-semibold">{gameId}</span>...
            </p>
            <p class="text-white/40 text-xs font-mono">
                Attempt {reconnectAttempts} of {maxReconnectAttempts}
            </p>
        </div>
    {:else}
        <div class="flex items-center justify-center w-16 h-16 rounded-2xl bg-brand-secondary/10 border border-brand-secondary/20 text-brand-secondary text-3xl">
            🔌
        </div>
        <div class="space-y-2">
            <h2 class="text-2xl font-black text-white tracking-wider uppercase">Offline</h2>
            <p class="text-white/60 text-sm">
                Failed to establish a secure link with the game server.
            </p>
        </div>
        <div class="flex flex-col sm:flex-row gap-3 w-full mt-4">
            <Button variant="outline" class="w-full" onclick={onRetry}>
                Retry Connection
            </Button>
            <Button variant="secondary" class="w-full" onclick={onReturnToMenu}>
                Return to Menu
            </Button>
        </div>
    {/if}
</div>
