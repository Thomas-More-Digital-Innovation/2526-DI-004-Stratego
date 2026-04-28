<script lang="ts">
    import { fly } from "svelte/transition";
    import type { HistoricalMove } from "$lib/types/game";
    import { getPieceIcon } from "$lib/utils/game";

    interface Props {
        move: HistoricalMove;
    }

    let { move }: Props = $props();

    const attackerIcon = $derived(getPieceIcon(move.attacker));
    const defenderIcon = $derived(getPieceIcon(move.defender));
</script>

<div
    class="unselectable flex items-center gap-1 bg-black/70 backdrop-blur-sm p-1 rounded-lg border border-white/20 shadow-xl -translate-y-10 opacity-80 hover:opacity-100 hover:scale-200 transition-transform duration-200 pointer-events-auto"
    in:fly={{ y: 10, duration: 400 }}
>
    <div
        class="relative flex flex-col items-center w-8"
        class:opacity-40={move.result === "loss" || move.result === "tie"}
    >
        {#if attackerIcon}
            <img
                src={attackerIcon}
                alt="Attacker"
                class="w-8 h-8 rounded bg-white/10"
                class:grayscale={move.result === "loss" ||
                    move.result === "tie"}
            />
        {:else}
            <div
                class="w-8 h-8 rounded bg-white/10 flex items-center justify-center font-bold text-white"
            >
                ?
            </div>
        {/if}
        {#if move.result === "loss" || move.result === "tie"}
            <div
                class="absolute inset-0 flex items-center justify-center text-xl drop-shadow-md"
            >
                ❌
            </div>
        {/if}
        <span class="text-[8px] font-bold text-white/70 mt-0.5">ATK</span>
    </div>

    <div class="text-[10px] font-bold text-white/50">VS</div>

    <div
        class="relative flex flex-col items-center w-8"
        class:opacity-40={move.result === "win" ||
            move.result === "tie" ||
            move.result === "capture"}
    >
        {#if defenderIcon}
            <img
                src={defenderIcon}
                alt="Defender"
                class="w-8 h-8 rounded bg-white/10"
                class:grayscale={move.result === "win" ||
                    move.result === "tie" ||
                    move.result === "capture"}
            />
        {:else}
            <div
                class="w-8 h-8 rounded bg-white/10 flex items-center justify-center font-bold text-white"
            >
                ?
            </div>
        {/if}
        {#if move.result === "win" || move.result === "tie" || move.result === "capture"}
            <div
                class="absolute inset-0 flex items-center justify-center text-xl drop-shadow-md"
            >
                ❌
            </div>
        {/if}
        <span class="text-[8px] font-bold text-white/70 mt-0.5">DEF</span>
    </div>
</div>
