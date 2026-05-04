<script lang="ts">
    import { AIs } from "$lib/data/AI.data";
    import { gamemodes } from "$lib/data/gamemodes.data";
    import type { GameMode } from "$lib/types/game";

    interface Props {
        gameMode: GameMode;
        ai1: string;
        ai2: string;
        step: "ai1" | "ai2";
    }

    let { gameMode, ai1, ai2, step }: Props = $props();

    const selectedAi1 = $derived(AIs.find((a) => a.id === ai1));
    const selectedAi2 = $derived(AIs.find((a) => a.id === ai2));
</script>

<div
    class="grid grid-cols-1 md:grid-cols-3 items-center gap-4 py-8 px-8 rounded-3xl bg-surface-elevated/10 border border-white/5 backdrop-blur-sm relative overflow-hidden"
>
    <div
        class="absolute inset-0 bg-linear-to-r from-red-500/5 via-transparent to-blue-500/5 pointer-events-none"
    ></div>

    <div class="flex flex-col items-center gap-3 text-center relative z-10">
        <span
            class="text-[10px] font-bold text-red-500 uppercase tracking-[0.2em]"
            >Player 1 (Red)</span
        >
        <div class="h-14 flex items-center justify-center">
            {#if gameMode.mode === gamemodes.human_vs_ai.mode}
                <div class="flex items-center gap-3">
                    <div
                        class="w-10 h-10 rounded-full bg-white/10 flex items-center justify-center border border-white/10"
                    >
                        <span class="text-lg">🧑</span>
                    </div>
                    <span class="text-xl font-black text-white uppercase"
                        >You</span
                    >
                </div>
            {:else if ai1}
                <div class="flex items-center gap-3">
                    <img
                        src={selectedAi1?.image}
                        alt=""
                        class="w-10 h-10 rounded-lg object-cover border border-white/10 shadow-lg"
                    />
                    <span class="text-xl font-black text-white uppercase"
                        >{selectedAi1?.name}</span
                    >
                </div>
            {:else}
                <span
                    class="text-xl font-black text-white/10 uppercase animate-pulse"
                    >Selecting...</span
                >
            {/if}
        </div>
    </div>

    <div class="hidden md:flex justify-center relative z-10">
        <div
            class="text-3xl font-black text-white/5 uppercase italic tracking-tighter"
        >
            VS
        </div>
    </div>

    <div class="flex flex-col items-center gap-3 text-center relative z-10">
        <span
            class="text-[10px] font-bold text-blue-500 uppercase tracking-[0.2em]"
            >Player 2 (Blue)</span
        >
        <div class="h-14 flex items-center justify-center">
            {#if ai2}
                <div class="flex items-center gap-3">
                    <img
                        src={selectedAi2?.image}
                        alt=""
                        class="w-10 h-10 rounded-lg object-cover border border-white/10 shadow-lg"
                    />
                    <span class="text-xl font-black text-white uppercase"
                        >{selectedAi2?.name}</span
                    >
                </div>
            {:else if step === "ai2" || (gameMode.mode === gamemodes.human_vs_ai.mode && step === "ai1")}
                <span
                    class="text-xl font-black text-white/10 uppercase animate-pulse"
                    >Selecting...</span
                >
            {:else}
                <span class="text-xl font-black text-white/5 uppercase"
                    >Pending</span
                >
            {/if}
        </div>
    </div>
</div>
