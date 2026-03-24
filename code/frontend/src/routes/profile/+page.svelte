<script lang="ts">
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import { authStore } from "$lib/state/auth.svelte";
    import { stats } from "$lib/api/client";
    import Card from "$lib/components/ui/Card.svelte";
    import Button from "$lib/components/ui/Button.svelte";
    import type { UserStats } from "$lib/types/game";

    let userStats = $state<UserStats | null>(null);
    let error = $state("");

    onMount(async () => {
        if (!authStore.user) {
            goto("/login");
            return;
        }
        try {
            userStats = await stats.getMine();
        } catch (e: any) {
            error = e.message || "Failed to load stats";
        }
    });

    function formatDuration(seconds: number): string {
        const mins = Math.floor(seconds / 60);
        const secs = Math.floor(seconds % 60);
        return `${mins}m ${secs}s`;
    }

    function getWinRate(): string {
        if (!userStats || userStats.total_games === 0) return "0%";
        return (
            ((userStats.wins / userStats.total_games) * 100).toFixed(1) + "%"
        );
    }

    async function handleLogout() {
        await authStore.logout();
        goto("/login");
    }
</script>

<svelte:head>
    <title>Stratego — Profile</title>
</svelte:head>

{#if authStore.user}
    <div class="space-y-8">
        <!-- Header -->
        <Card class="flex items-center gap-6">
            <div
                class="w-16 h-16 rounded-full bg-gradient-to-br from-brand-primary to-brand-secondary flex items-center justify-center text-white text-2xl font-bold shrink-0"
            >
                {authStore.user.username[0]?.toUpperCase() || "?"}
            </div>
            <div class="flex-1">
                <h1 class="text-xl font-bold text-white">
                    {authStore.user.username}
                </h1>
                <p class="text-white/40 text-sm">
                    Joined {new Date(
                        authStore.user.created_at,
                    ).toLocaleDateString()}
                </p>
            </div>
            <Button variant="secondary" onclick={handleLogout}>Logout</Button>
        </Card>

        {#if error}
            <div
                class="bg-brand-secondary/20 border border-brand-secondary/30 text-brand-secondary rounded-xl px-4 py-3 text-sm text-center"
            >
                {error}
            </div>
        {:else if userStats}
            <!-- Stats Grid -->
            <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
                <Card class="text-center">
                    <div class="text-3xl font-bold text-white">
                        {userStats.total_games}
                    </div>
                    <div
                        class="text-xs text-white/40 uppercase tracking-wider mt-1"
                    >
                        Total Games
                    </div>
                </Card>
                <Card class="text-center">
                    <div class="text-3xl font-bold text-brand-accent">
                        {userStats.wins}
                    </div>
                    <div
                        class="text-xs text-white/40 uppercase tracking-wider mt-1"
                    >
                        Wins
                    </div>
                </Card>
                <Card class="text-center">
                    <div class="text-3xl font-bold text-brand-secondary">
                        {userStats.losses}
                    </div>
                    <div
                        class="text-xs text-white/40 uppercase tracking-wider mt-1"
                    >
                        Losses
                    </div>
                </Card>
                <Card class="text-center">
                    <div class="text-3xl font-bold text-white">
                        {userStats.draws}
                    </div>
                    <div
                        class="text-xs text-white/40 uppercase tracking-wider mt-1"
                    >
                        Draws
                    </div>
                </Card>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                <Card
                    class="text-center bg-gradient-to-br from-brand-primary/10 to-brand-secondary/10"
                >
                    <div class="text-3xl font-bold text-brand-accent">
                        {getWinRate()}
                    </div>
                    <div
                        class="text-xs text-white/40 uppercase tracking-wider mt-1"
                    >
                        Win Rate
                    </div>
                </Card>
                <Card class="text-center">
                    <div class="text-3xl font-bold text-white">
                        {userStats.total_moves}
                    </div>
                    <div
                        class="text-xs text-white/40 uppercase tracking-wider mt-1"
                    >
                        Total Moves
                    </div>
                </Card>
                <Card class="text-center">
                    <div class="text-3xl font-bold text-white">
                        {formatDuration(userStats.avg_game_duration_seconds)}
                    </div>
                    <div
                        class="text-xs text-white/40 uppercase tracking-wider mt-1"
                    >
                        Avg Duration
                    </div>
                </Card>
            </div>
        {/if}

        <!-- Actions -->
        <div class="flex gap-4 justify-center">
            <Button variant="primary" onclick={() => goto("/board-setups")}>
                Manage Board Setups
            </Button>
            <Button variant="outline" onclick={() => goto("/")}>Home</Button>
        </div>
    </div>
{/if}
