<script lang="ts">
    import { GameAPI } from "$lib/api";
    import type { GameMode } from "$lib/types";

    let api = new GameAPI();
    let selectedMode = $state<GameMode>("human_vs_ai");
    let isCreating = $state<boolean>(false);

    let { errorMessage = $bindable("") } = $props();

    async function startNewGame() {
        isCreating = true;
        errorMessage = "";

        if (selectedMode === "human_vs_human") {
            // TODO coming soon
            errorMessage = "Coming Soon";
            isCreating = false;
            return;
        }

        try {
            const gameInfo = await api.createGame(selectedMode);
            // Use window.location for navigation
            window.location.href = `/game/${gameInfo.gameId}?mode=${selectedMode}`;
        } catch (error) {
            errorMessage = `Failed to create game: ${error}`;
            isCreating = false;
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
            disabled={isCreating}
        >
            <div class="mode-icon">🧑 vs 🧑</div>
            <h3>Human vs Human</h3>
            <p>Coming Soon</p>
        </button>

        <button
            class="game-mode-card"
            class:selected={selectedMode === "human_vs_ai"}
            onclick={() => (selectedMode = "human_vs_ai")}
            disabled={isCreating}
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
            disabled={isCreating}
        >
            <div class="mode-icon">🤖 vs 🤖</div>
            <h3>AI vs AI</h3>
            <p>
                Watch two AI players battle it out. All pieces visible for
                spectating.
            </p>
        </button>
    </div>

    <button class="start-btn" onclick={startNewGame} disabled={isCreating}>
        {#if isCreating}
            Creating game...
        {:else}
            Start Game
        {/if}
    </button>
</div>

<style>
    .game-modes {
        background: var(--bg);
        padding: 32px;
    }
    .mode-cards {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax( 300px, 700px));
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
            box-shadow: 0 6px 12px color-mix(in srgb, var(--primary-dark), transparent 50%);
        }

        &:disabled {
            opacity: 0.6;
            cursor: not-allowed;
        }
    }
</style>
