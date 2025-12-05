<script lang="ts">
    import PieceBank from "../../../../lib/components/PieceBank.svelte";
    import BoardGrid from "../../../../lib/components/BoardGrid.svelte";
    import { PIECE_INVENTORY, type BoardSetup } from "$lib/types/boardSetup";
    import { createBoardSetupState } from "./BoardSetupEditorState.svelte";
    import BoardSetupNamer from "./BoardSetupNamer.svelte";

    interface Props {
        initialSetup?: string;
        onSave: (setup: BoardSetup) => void;
        onCancel: () => void;
    }

    let { initialSetup, onSave, onCancel }: Props = $props();

    let state = $state(createBoardSetupState(initialSetup, onSave));
</script>

<div class="editor">
    <div class="editor-header">
        <div class="header-left">
            <h3>Board Setup Editor</h3>
            <div class="debug-inline">
                <code>{state.board.join('')}</code>
                <button
                    class="btn-copy"
                    onclick={() => {
                        navigator.clipboard.writeText(state.board.join(''));
                        alert("Copied to clipboard!");
                    }}
                >
                    📋
                </button>
            </div>
        </div>
        <div class="actions">
            <button class="btn-secondary" onclick={onCancel}>Cancel</button>
            <button class="btn-primary" onclick={state.handleSave}
                >Save Setup</button
            >
        </div>
    </div>

    <div class="editor-content">
        <div class="main">
            <BoardSetupNamer bind:state />
            <BoardGrid
                board={state.board}
                selectedPiece={state.selectedPiece}
                selectedPieceIndex={state.selectedPieceIndex}
                onCellClick={state.handleCellClick}
                onSwap={state.swapCells}
            />
        </div>

        <div class="sidebar">
            <PieceBank
                pieceInventory={PIECE_INVENTORY}
                remaining={state.getRemainingPieces()}
                selectedPiece={state.selectedPiece}
                onSelect={state.handlePieceSelect}
            />
        </div>
    </div>
</div>

<style>
    .editor {
        display: flex;
        flex-direction: column;
        background: var(--bg-accent);
        border-radius: 12px;
        padding: 20px;
        height: calc(100vh - 40px);
        max-width: 1400px;
        margin: 24px auto;
        overflow: hidden;
        min-height: 0;
    }

    .editor-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 20px;
        gap: 20px;

        .header-left {
            display: flex;
            align-items: center;
            gap: 20px;
            flex: 1;

            h3 {
                margin: 0;
                color: var(--text);
                white-space: nowrap;
            }

            .debug-inline {
                display: flex;
                align-items: center;
                gap: 8px;
                background: #1a1a1a;
                padding: 8px 12px;
                border-radius: 6px;
                flex: 1;
                overflow: hidden;

                code {
                    color: #0f0;
                    font-size: 0.75rem;
                    word-break: break-all;
                    overflow: hidden;
                    text-overflow: ellipsis;
                    white-space: nowrap;
                }

                .btn-copy {
                    background: #4a4;
                    color: white;
                    border: none;
                    padding: 4px 8px;
                    border-radius: 4px;
                    cursor: pointer;
                    font-size: 1rem;
                    transition: background 0.2s;
                    flex-shrink: 0;

                    &:hover {
                        background: #393;
                    }
                }
            }
        }

        .actions {
            display: flex;
            gap: 10px;
        }
    }

    .editor-content {
        display: grid;
        grid-template-columns: 1fr 350px;
        gap: 20px;
        align-items: start;
        flex: 1 1 auto;
        min-height: 0;
        overflow: hidden;
        align-items: center;

        @media (max-width: 1024px) {
            grid-template-columns: 1fr;
        }

        .main {
            display: flex;
            flex-direction: column;
                    margin: 0 auto;

            gap: 20px;
        }
        .sidebar {
            display: flex;
            flex-direction: column;
            gap: 20px;
            min-height: 0; /* allow .piece-bank to shrink/scroll */
            height: 100%;

            
        }
    }

    .sidebar :global(.piece-bank) :global(.pieces) {
        flex: 1 1 auto;
        overflow-y: auto;
        min-height: 0;
    }

    .btn-secondary {
        background: transparent;
        color: var(--text);
        border: 1px solid #444;
        padding: 10px 20px;
        border-radius: 6px;
        cursor: pointer;
        transition: all 0.2s;

        &:hover {
            border-color: #66f;
            color: #66f;
        }
    }
</style>
