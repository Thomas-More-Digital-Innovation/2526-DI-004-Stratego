<script lang="ts">
    import { goto } from "$app/navigation";
    import { authStore } from "$lib/state/auth.svelte";
    import Button from "$lib/components/ui/Button.svelte";
    import { games as gamesApi } from "$lib/api/client";
    import { toastStore } from "$lib/state/toast.svelte";

    let reconnectableGame = $state<{
        gameId: string;
        gameType: string;
        seatIndex: number;
    } | null>(null);
    let isReconnecting = $state(false);

    function checkGame() {
        if (!authStore.user) {
            reconnectableGame = null;
            return;
        }

        gamesApi
            .getReconnectable()
            .then((res: any) => {
                if (res.hasGame && res.gameId && res.gameType) {
                    reconnectableGame = {
                        gameId: res.gameId,
                        gameType: res.gameType,
                        seatIndex: res.seatIndex ?? 0,
                    };
                } else {
                    reconnectableGame = null;
                }
            })
            .catch((e) => {
                console.error("Failed to check for reconnectable game:", e);
            });
    }

    $effect(() => {
        if (authStore.user) {
            checkGame();

            const interval = setInterval(checkGame, 15000);

            const handleVisibilityChange = () => {
                if (document.visibilityState === "visible") {
                    checkGame();
                }
            };
            document.addEventListener(
                "visibilitychange",
                handleVisibilityChange,
            );

            return () => {
                clearInterval(interval);
                document.removeEventListener(
                    "visibilitychange",
                    handleVisibilityChange,
                );
            };
        } else {
            reconnectableGame = null;
        }
    });

    function reconnect() {
        if (!reconnectableGame || isReconnecting) return;
        isReconnecting = true;

        gamesApi
            .getReconnectable()
            .then((res: any) => {
                if (
                    reconnectableGame != null &&
                    res.hasGame &&
                    res.gameId === reconnectableGame?.gameId
                ) {
                    goto(
                        `/game/${reconnectableGame.gameId}?mode=${reconnectableGame.gameType}&seat=${reconnectableGame.seatIndex}`,
                    );
                } else {
                    toastStore.warning(
                        "Game session has ended or is no longer available.",
                    );
                    reconnectableGame = null;
                    isReconnecting = false;
                }
            })
            .catch((e) => {
                toastStore.error("Failed to verify active game session.");
                isReconnecting = false;
            });
    }
</script>

{#if reconnectableGame}
    <div
        class="relative overflow-hidden rounded-2xl glass border border-brand-accent/20 p-6 flex flex-col md:flex-row items-center justify-between gap-4 transition-all duration-300 shadow-lg glow-accent"
    >
        <div
            class="absolute -right-16 -top-16 w-32 h-32 bg-brand-accent/10 blur-3xl pointer-events-none rounded-full"
        ></div>
        <div
            class="absolute -left-16 -bottom-16 w-32 h-32 bg-brand-primary/10 blur-3xl pointer-events-none rounded-full"
        ></div>

        <div class="flex items-center gap-4 z-10">
            <div
                class="flex items-center justify-center w-12 h-12 rounded-xl bg-brand-accent/10 border border-brand-accent/30 text-brand-accent text-2xl animate-pulse"
            >
                ⚠️
            </div>
            <div>
                <h3 class="text-white font-bold text-lg tracking-widest">
                    Active Game In Progress
                </h3>
                <p class="text-white/60 text-sm mt-0.5">
                    You have an ongoing game in progress that will be cleaned up
                    shortly (<span class="font-mono text-brand-accent"
                        >{reconnectableGame.gameId}</span
                    >).
                </p>
            </div>
        </div>

        <div class="w-full md:w-auto z-10 shrink-0">
            <Button
                variant="outline"
                class="w-full md:w-auto border-brand-accent hover:bg-brand-accent/10 text-brand-accent font-bold"
                onclick={reconnect}
                loading={isReconnecting}
            >
                Reconnect to Game
            </Button>
        </div>
    </div>
{/if}
