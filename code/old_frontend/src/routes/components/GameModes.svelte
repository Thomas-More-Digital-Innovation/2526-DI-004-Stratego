<script lang="ts">
    import { GameAPI } from "$lib/api";
    import type { GameMode } from "$lib/types";
    import SelectAI from "./SelectAI.svelte";

    interface User {
        id: number;
        username: string;
        profile_picture?: string;
        created_at: string;
        updated_at: string;
    }

    type StartGameState = "create" | "selectAi1" | "selectAi2" | "creating";

    let api = new GameAPI();
    let selectedMode = $state<GameMode>("human_vs_ai");
    let startGameState = $state<StartGameState>("create");
    let ai1 = $state<string>("");
    let ai2 = $state<string>("");

    let { errorMessage = $bindable(""), user }: { errorMessage: string; user: User } = $props();

    function cancel() {
        startGameState = "create";
        ai1 = "";
        ai2 = "";
    }

    async function startNewGame() {
        errorMessage = "";

        if (startGameState === "create") {
            return selectMode();
        } 

        try {
            const gameInfo = await api.createGame(selectedMode, ai1, ai2);
            // Use window.location for navigation
            window.location.href = `/game/${gameInfo.gameId}?mode=${selectedMode}`;
        } catch (error) {
            errorMessage = `Failed to create game: ${error}`;
            cancel();
        }
    }

    function selectMode() {
        if (selectedMode === "human_vs_human") {
            // TODO coming soon
            errorMessage = "Coming Soon";
            startGameState = "create";
            return;
        }
        
        startGameState = "selectAi1";
        ai1 = "";
        ai2 = "";
    }

    function selectAi() {
        if (selectedMode === "human_vs_ai" || startGameState === "selectAi2") {
            startGameState = "creating";
            startNewGame();
        } else {
            startGameState = "selectAi2";
        }
    }
</script>

<div class="game-modes">
    <h2>Game Modes</h2>

    <div class="mode-cards">
        <button
            class="game-mode-card"
            class:selected={selectedMode === "human_vs_human"}
            onclick={() => (selectedMode = "human_vs_human")}
            disabled={startGameState === "creating"}
        >
            <div class="mode-icon">🧑 vs 🧑</div>
            <h3>Human vs Human</h3>
            <p>Coming Soon</p>
        </button>

        <button
            class="game-mode-card"
            class:selected={selectedMode === "human_vs_ai"}
            onclick={() => (selectedMode = "human_vs_ai")}
            disabled={startGameState === "creating"}
        >
            <div class="mode-icon">🧑 vs 🤖</div>
            <h3>Human vs AI</h3>
            <p>
                Play against the computer. You control the red pieces, AI
                controls blue.
            </p>
        </button>

        <button
            class="game-mode-card"
            class:selected={selectedMode === "ai_vs_ai"}
            onclick={() => (selectedMode = "ai_vs_ai")}
            disabled={startGameState === "creating"}
        >
            <div class="mode-icon">🤖 vs 🤖</div>
            <h3>AI vs AI</h3>
            <p>
                Watch two AI players battle it out. All pieces visible for
                spectating.
            </p>
        </button>
    </div>

    <button
        class="start-btn"
        onclick={startNewGame}
        disabled={startGameState === "creating"}
    >
        {#if startGameState === "creating"}
            Creating game...
        {:else if startGameState === "selectAi1" || startGameState === "selectAi2"}
            Select AI
        {:else}
            Start Game
        {/if}
    </button>
</div>

{#if startGameState === "selectAi1"}
    <SelectAI
        title="Select AI 1"
        onSelectAI={(ai) => {
            ai1 = ai;
            selectAi();
        }}
        onClose={cancel}
    />
{:else if startGameState === "selectAi2"}
    <SelectAI
        title="Select AI 2"
        onSelectAI={(ai) => {
            ai2 = ai;
            selectAi();
        }}
        onClose={cancel}
    />
{/if}

<style>
    .game-modes {
        background: var(--bg);
        padding: 32px;
    }
    .mode-cards {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(300px, 700px));
        gap: 20px;
        margin-bottom: 30px;

        .game-mode-card {
            background: var(--bg-accent);
            border: 3px solid transparent;
            border-radius: 12px;
            padding: 25px;
            cursor: pointer;
            transition: all 0.3s;
            text-align: left;

            &:hover:not(:disabled) {
                transform: translateY(-5px);
                box-shadow: 0 8px 16px rgba(0, 0, 0, 0.4);
                border-color: var(--secondary);
            }

            &.selected {
                border-color: var(--primary);
                background: var(--primary-dark);
            }

            &:disabled {
                opacity: 0.6;
                cursor: not-allowed;
            }

            .mode-icon {
                font-size: 3rem;
                margin-bottom: 15px;
            }

            h3 {
                margin: 0 0 10px 0;
                color: #e2e8f0;
                font-size: 1.4rem;
            }

            p {
                margin: 0 0 15px 0;
                color: #cbd5e0;
                line-height: 1.5;
            }
        }
    }

    .start-btn {
        width: 100%;
        padding: 18px;
        background: var(--primary-dark);
        color: var(--text);
        border: 4px solid var(--text);

        border-radius: 8px;
        font-size: 1.3rem;
        font-weight: 700;
        cursor: pointer;
        transition: all 0.2s;
        text-transform: uppercase;
        letter-spacing: 1px;

        &:hover:not(:disabled) {
            background: var(--primary);
            transform: translateY(-2px);
            box-shadow: 0 6px 12px
                color-mix(in srgb, var(--primary-dark), transparent 50%);
        }

        &:disabled {
            opacity: 0.6;
            cursor: not-allowed;
        }
    }
</style>
