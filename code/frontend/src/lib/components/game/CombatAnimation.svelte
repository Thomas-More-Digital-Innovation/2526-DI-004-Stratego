<script lang="ts">
    import type { Piece } from "$lib/types/game";
    import PieceDisplay from "./Piece.svelte";

    interface Props {
        attacker: Piece | null;
        defender: Piece | null;
        attackerWon: boolean;
        defenderWon: boolean;
        onComplete?: () => void;
    }

    let { attacker, defender, attackerWon, defenderWon, onComplete }: Props =
        $props();

    let stage = $state<"reveal" | "clash" | "result" | "hide">("reveal");

    $effect(() => {
        if (attacker && defender) {
            stage = "reveal";
            setTimeout(() => {
                stage = "clash";
                setTimeout(() => {
                    stage = "result";
                    setTimeout(() => {
                        stage = "hide";
                        setTimeout(() => onComplete?.(), 500);
                    }, 1000);
                }, 500);
            }, 500);
        }
    });
</script>

{#if attacker && defender}
    <div class="overlay" class:hide={stage === "hide"}>
        <div
            class="arena glass rounded-2xl"
            class:reveal={stage === "reveal"}
            class:clash={stage === "clash"}
        >
            <h2
                class="text-3xl font-black text-brand-accent text-center mb-6 uppercase tracking-widest animate-pulse"
            >
                Combat!
            </h2>

            <div class="flex items-center justify-center gap-8">
                <!-- Attacker -->
                <div
                    class="fighter"
                    class:winner={attackerWon}
                    class:loser={!attackerWon && defenderWon}
                >
                    <span
                        class="text-xs font-bold text-white/60 uppercase tracking-wider"
                    >
                        {attacker.ownerName || "Attacker"}
                    </span>
                    <div
                        class="w-20 h-20 flex items-center justify-center rounded-xl bg-surface-elevated border border-white/10"
                    >
                        <PieceDisplay
                            piece={attacker}
                            viewerId={attacker.ownerId}
                        />
                    </div>
                    {#if attacker.type}
                        <span class="text-sm font-bold text-brand-accent"
                            >{attacker.type}</span
                        >
                    {/if}
                    {#if attacker.rank}
                        <span class="text-xs text-white/50"
                            >Rank: {attacker.rank}</span
                        >
                    {/if}
                </div>

                <span
                    class="text-2xl font-black text-brand-secondary animate-pulse"
                    >VS</span
                >

                <!-- Defender -->
                <div
                    class="fighter"
                    class:winner={defenderWon}
                    class:loser={!defenderWon && attackerWon}
                >
                    <span
                        class="text-xs font-bold text-white/60 uppercase tracking-wider"
                    >
                        {defender.ownerName || "Defender"}
                    </span>
                    <div
                        class="w-20 h-20 flex items-center justify-center rounded-xl bg-surface-elevated border border-white/10"
                    >
                        <PieceDisplay
                            piece={defender}
                            viewerId={defender.ownerId}
                        />
                    </div>
                    {#if defender.type}
                        <span class="text-sm font-bold text-brand-accent"
                            >{defender.type}</span
                        >
                    {/if}
                    {#if defender.rank}
                        <span class="text-xs text-white/50"
                            >Rank: {defender.rank}</span
                        >
                    {/if}
                </div>
            </div>

            {#if stage === "result"}
                <p
                    class="text-center text-xl font-bold mt-6 text-brand-accent animate-slide-up"
                >
                    {#if attackerWon && defenderWon}
                        Both Eliminated!
                    {:else if attackerWon}
                        Attacker Wins!
                    {:else if defenderWon}
                        Defender Wins!
                    {/if}
                </p>
            {/if}
        </div>
    </div>
{/if}

<style>
    .overlay {
        position: fixed;
        inset: 0;
        background: rgba(0, 0, 0, 0.85);
        display: flex;
        justify-content: center;
        align-items: center;
        z-index: 1000;
        animation: fadeIn 0.3s ease-in;
    }

    .overlay.hide {
        animation: fadeOut 0.5s ease-in forwards;
    }

    .arena {
        padding: 2.5rem;
        min-width: 500px;
    }

    .arena.reveal {
        animation: zoomIn 0.5s ease-out;
    }

    .arena.clash {
        animation: shake 0.5s ease-in-out;
    }

    .fighter {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 8px;
        padding: 1rem;
        border-radius: 1rem;
        background: rgba(255, 255, 255, 0.03);
        transition: all 0.3s;
        min-width: 140px;
    }

    .fighter.winner {
        background: rgba(144, 238, 144, 0.15);
        border: 1px solid rgba(144, 238, 144, 0.4);
        transform: scale(1.05);
    }

    .fighter.loser {
        opacity: 0.4;
        filter: grayscale(0.7);
        transform: scale(0.9);
    }

    @keyframes fadeIn {
        from {
            opacity: 0;
        }
        to {
            opacity: 1;
        }
    }
    @keyframes fadeOut {
        from {
            opacity: 1;
        }
        to {
            opacity: 0;
        }
    }
    @keyframes zoomIn {
        from {
            transform: scale(0.5);
            opacity: 0;
        }
        to {
            transform: scale(1);
            opacity: 1;
        }
    }
    @keyframes shake {
        0%,
        100% {
            transform: translateX(0);
        }
        25% {
            transform: translateX(-10px);
        }
        75% {
            transform: translateX(10px);
        }
    }

    @media (max-width: 768px) {
        .arena {
            min-width: 90vw;
            padding: 1.5rem;
        }
    }
</style>
