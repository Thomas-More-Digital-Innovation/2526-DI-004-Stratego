<script lang="ts">
    import { onMount, onDestroy } from "svelte";
    import { page } from "$app/stores";
    import { GameSocket } from "$lib/api/websocket";
    import { gameStore } from "$lib/state/game.svelte";
    import Board from "$lib/components/game/Board.svelte";
    import GameInfo from "$lib/components/game/GameInfo.svelte";
    import GameHistory from "$lib/components/game/GameHistory.svelte";
    import CombatAnimation from "$lib/components/game/CombatAnimation.svelte";
    import SetupBanner from "$lib/components/game/SetupBanner.svelte";
    import Button from "$lib/components/ui/Button.svelte";
    import type { Position } from "$lib/types/game";

    let socket = new GameSocket();
    let gameId = $state("");
    let gameMode = $state("");
    let error = $state("");
    let connected = $state(false);
    let validMoves = $state<Position[]>([]);
    let setupSwapPos1 = $state<Position | null>(null);

    const isSetupPhase = $derived(gameStore.gameState?.isSetupPhase ?? false);

    const isHumanTurn = $derived.by(() => {
        return (
            gameStore.gameState?.currentPlayerId === 0 &&
            gameMode === "human_vs_ai" &&
            !gameStore.gameState?.isGameOver &&
            !isSetupPhase
        );
    });

    const viewerId = $derived(gameMode === "human_vs_ai" ? 0 : -1);

    onMount(async () => {
        gameId = $page.params.id || "";
        const params = new URLSearchParams(window.location.search);
        gameMode = params.get("mode") || "human_vs_ai";

        setupHandlers();

        try {
            const playerId = gameMode === "human_vs_ai" ? 0 : -1;
            await socket.connect(gameId, playerId);
            connected = true;
        } catch (e) {
            error = `Failed to connect: ${e}`;
            setTimeout(() => (window.location.href = "/"), 3000);
        }
    });

    onDestroy(() => {
        socket.disconnect();
        gameStore.reset();
    });

    function setupHandlers() {
        socket.on("gameState", (data) => gameStore.updateGameState(data));
        socket.on("boardState", (data) => gameStore.updateBoardState(data));
        socket.on("moveHistory", (data) =>
            gameStore.loadMoveHistory(data.moves),
        );

        socket.on("moveResult", (data) => {
            if (!data.success) {
                error = data.error || "Move failed";
                setTimeout(() => (error = ""), 3000);
                gameStore.setSelectedPosition(null);
                validMoves = [];
            }
        });

        socket.on("validMoves", (data) => {
            validMoves = data.validMoves || [];
        });

        socket.on("combat", (data) => {
            gameStore.showCombatAnimation({
                attacker: data.attacker,
                defender: data.defender,
                attackerWon: data.attackerWon,
                defenderWon: data.defenderWon,
            });
        });

        socket.on("gameOver", () => {
            setTimeout(() => {
                if (confirm("Game Over! Return to menu?")) {
                    window.location.href = "/";
                }
            }, 2000);
        });

        socket.on("error", (data) => {
            error = data.error;
            setTimeout(() => (error = ""), 3000);
        });
    }

    function handleCellClick(x: number, y: number) {
        if (isSetupPhase && gameMode === "human_vs_ai") {
            handleSetupClick(x, y);
            return;
        }

        if (gameStore.isReplaying || !isHumanTurn) return;

        const board = gameStore.boardState?.board;
        if (!board) return;

        const clickedPiece = board[y][x];
        const selected = gameStore.selectedPosition;
        const hasValidPiece = clickedPiece && clickedPiece.ownerName;

        if (!selected) {
            if (hasValidPiece && clickedPiece.ownerId === 0) {
                gameStore.setSelectedPosition({ x, y });
                socket.requestValidMoves({ x, y });
            }
            return;
        }

        if (selected.x === x && selected.y === y) {
            gameStore.setSelectedPosition(null);
            validMoves = [];
            return;
        }

        const isValid = validMoves.some((m) => m.x === x && m.y === y);

        if (isValid) {
            socket.sendMove(selected, { x, y });
            gameStore.setSelectedPosition(null);
            validMoves = [];
        } else if (hasValidPiece && clickedPiece.ownerId === 0) {
            gameStore.setSelectedPosition({ x, y });
            socket.requestValidMoves({ x, y });
        } else {
            gameStore.setSelectedPosition(null);
            validMoves = [];
        }
    }

    function handleSetupClick(x: number, y: number) {
        if (y < 6 || y > 9) return;

        if (!setupSwapPos1) {
            setupSwapPos1 = { x, y };
        } else {
            socket.sendSwapPieces(setupSwapPos1, { x, y });
            setupSwapPos1 = null;
        }
    }

    function handleCellDragStart(e: DragEvent, x: number, y: number) {
        if (!isSetupPhase) return;
        e.dataTransfer?.setData("text/plain", JSON.stringify({ x, y }));
    }

    function handleCellDrop(e: DragEvent, x: number, y: number) {
        if (!isSetupPhase) return;
        const data = e.dataTransfer?.getData("text/plain");
        if (!data) return;

        try {
            const from = JSON.parse(data) as Position;
            if (from.x === x && from.y === y) return;
            socket.sendSwapPieces(from, { x, y });
        } catch (e) {
            console.error("Failed to parse drop data", e);
        }
    }

    function handleRandomize() {
        socket.sendRandomizeSetup();
        setupSwapPos1 = null;
    }

    function handleStartGame() {
        socket.sendStartGame();
        setupSwapPos1 = null;
    }

    function handleLoadSetup(setupData: string) {
        socket.sendLoadSetup(setupData);
        setupSwapPos1 = null;
    }

    function saveGame() {
        try {
            const data = gameStore.exportGame();
            const blob = new Blob([data], { type: "application/json" });
            const url = URL.createObjectURL(blob);
            const a = document.createElement("a");
            a.href = url;
            a.download = `stratego-${gameId}-${Date.now()}.json`;
            a.click();
            URL.revokeObjectURL(url);
        } catch {
            error = "Failed to save game";
        }
    }
</script>

<svelte:head>
    <title>Stratego — Game {gameId}</title>
</svelte:head>

{#if error}
    <div
        class="bg-brand-secondary/20 border border-brand-secondary/30 text-brand-secondary rounded-xl px-4 py-3 text-sm text-center mb-4"
    >
        ⚠️ {error}
    </div>
{/if}

{#if !connected}
    <div class="flex flex-col items-center justify-center min-h-[60vh] gap-4">
        <div
            class="w-10 h-10 border-3 border-white/20 border-t-brand-primary rounded-full animate-spin"
        ></div>
        <p class="text-white/40">Connecting to game...</p>
    </div>
{:else}
    <div class="flex items-center justify-between mb-6">
        <Button variant="ghost" onclick={() => (window.location.href = "/")}>
            ← Back
        </Button>
        <h1 class="text-lg font-bold text-white uppercase tracking-wider">
            Game
        </h1>
        <Button
            variant="outline"
            size="sm"
            onclick={saveGame}
            disabled={!connected}
        >
            💾 Save
        </Button>
    </div>

    <div class="grid grid-cols-[280px_1fr_280px] gap-6 items-start">
        <div>
            <GameInfo gameState={gameStore.gameState} {gameMode} />
        </div>

        <div class="flex justify-center">
            <Board
                boardState={gameStore.boardState}
                selectedPosition={isSetupPhase
                    ? setupSwapPos1
                    : gameStore.selectedPosition}
                onCellClick={handleCellClick}
                onCellDragStart={handleCellDragStart}
                onCellDrop={handleCellDrop}
                isInteractive={!gameStore.isReplaying &&
                    (isHumanTurn || isSetupPhase)}
                {viewerId}
                {validMoves}
                disabledRows={isSetupPhase ? [0, 1, 2, 3, 4, 5] : []}
                visualDisabledRows={isSetupPhase ? [4, 5] : []}
            />
        </div>

        <div>
            {#if !isSetupPhase}
                <GameHistory
                    currentMoveIndex={gameStore.currentHistoryIndex}
                    totalMoves={gameStore.history.length}
                    isReplaying={gameStore.isReplaying}
                    onPrevious={() => gameStore.previousMove()}
                    onNext={() => gameStore.nextMove()}
                    onGoToMove={(index) => gameStore.goToMove(index)}
                    onExitReplay={() => gameStore.exitReplay()}
                />
            {:else}
                <div
                    class="glass rounded-2xl p-6 space-y-3 border border-white/10"
                >
                    <h3
                        class="text-sm font-bold text-brand-accent uppercase tracking-wider"
                    >
                        Setup Instructions
                    </h3>
                    <ul class="text-white/50 text-sm space-y-2">
                        <li>Click two pieces to swap them</li>
                        <li>Use "Randomize" for a random setup</li>
                        <li>Click "Start Game" when ready</li>
                    </ul>
                </div>
            {/if}
        </div>
    </div>
{/if}

{#if gameStore.combatAnimation}
    <CombatAnimation
        attacker={gameStore.combatAnimation.attacker}
        defender={gameStore.combatAnimation.defender}
        attackerWon={gameStore.combatAnimation.attackerWon}
        defenderWon={gameStore.combatAnimation.defenderWon}
        onComplete={() => {
            socket.sendAnimationComplete();
            gameStore.hideCombatAnimation();
        }}
    />
{/if}

{#if isSetupPhase && gameMode === "human_vs_ai"}
    <SetupBanner
        onRandomize={handleRandomize}
        onStart={handleStartGame}
        onLoadSetup={handleLoadSetup}
        {viewerId}
    />
{/if}
