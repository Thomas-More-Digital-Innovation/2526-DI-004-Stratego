<script lang="ts">
    import { gameStore } from "$lib/state/game.svelte";

    interface Props {
        selectedPlayer: number;
        onSelectPlayer?: (player: number) => void;
    }

    let { selectedPlayer, onSelectPlayer }: Props = $props();
</script>

<button
    class="flex items-center gap-3 mt-2 cursor-pointer group"
    aria-label="Select AI player"
    onclick={() => onSelectPlayer?.(selectedPlayer === 0 ? 1 : 0)}
>
    <div
        class="flex items-center gap-2 px-3 py-1.5 rounded-lg border-2 transition-all {selectedPlayer ===
        0
            ? 'border-red-500 bg-red-500/10 group-hover:border-red-500/40 shadow-[0_0_15px_rgba(239,68,68,0.2)]'
            : 'border-white/5 bg-white/5 group-hover:border-white/20 opacity-40 group-hover:opacity-100'}"
    >
        <div
            class="w-2 h-2 rounded-full bg-red-500 shadow-[0_0_5px_rgba(239,68,68,1)]"
        ></div>
        <span
            class="text-[10px] font-black text-white uppercase tracking-wider"
        >
            {gameStore.gameState?.player1Username || "AI Red"}
        </span>
    </div>
    <div class="text-[10px] font-black text-white/10 uppercase italic">vs</div>
    <div
        class="flex items-center gap-2 px-3 py-1.5 rounded-lg border-2 transition-all {selectedPlayer ===
        1
            ? 'border-blue-500 bg-blue-500/10 group-hover:border-blue-500/40 shadow-[0_0_15px_rgba(59,130,246,0.2)]'
            : 'border-white/5 bg-white/5 group-hover:border-white/20 opacity-40 group-hover:opacity-100'}"
    >
        <span
            class="text-[10px] font-black text-white uppercase tracking-wider"
        >
            {gameStore.gameState?.player2Username || "AI Blue"}
        </span>
        <div
            class="w-2 h-2 rounded-full bg-blue-500 shadow-[0_0_5px_rgba(59,130,246,1)]"
        ></div>
    </div>
</button>
