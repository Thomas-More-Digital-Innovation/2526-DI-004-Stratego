<script lang="ts">
    import type { PageData } from "./$types";
    import BoardSetupEditor from "./components/BoardSetupEditor.svelte";
    import SetupCard from "$lib/components/SetupCard.svelte";
    import { MAX_BOARD_SETUPS, type BoardSetup } from "$lib/types/boardSetup.type";
    import { invalidateAll } from "$app/navigation";

    let { data }: { data: PageData } = $props();

    let error = $state(data.error || "");
    let showEditor = $state(false);
    let editingSetup = $state<BoardSetup | null>(null);

    function openCreateEditor() {
        if (data.setups.length >= MAX_BOARD_SETUPS) {
            return;
        }
        editingSetup = null;
        showEditor = true;
    }

    function openEditEditor(setup: BoardSetup) {
        editingSetup = setup;
        showEditor = true;
    }

    function handleEditorSave(setupData: BoardSetup) {
        saveSetup(setupData);
        showEditor = false;
    }

    function handleEditorCancel() {
        showEditor = false;
        editingSetup = null;
    }

    async function saveSetup(setupData: BoardSetup) {
        if (!setupData) return;
        try {
            const url = editingSetup
                ? `http://localhost:8080/api/board-setups?id=${editingSetup.name}`
                : `http://localhost:8080/api/board-setups`;

            const method = editingSetup ? "PUT" : "POST";
            console.log(setupData);
            const response = await fetch(url, {
                method,
                headers: { "Content-Type": "application/json" },
                credentials: "include",
                body: JSON.stringify({
                    name: setupData.name,
                    description: setupData.description,
                    setup_data: setupData.setupData,
                    is_default: setupData.isDefault,
                }),
            });

            if (!response.ok) throw new Error("Failed to save setup");

            editingSetup = null;
            await invalidateAll();
        } catch (e) {
            alert("Failed to save: " + e);
        }
    }

    async function deleteSetup(setupId: number) {
        if (!confirm("Delete this setup?")) return;

        try {
            const response = await fetch(
                `http://localhost:8080/api/board-setups?id=${setupId}`,
                {
                    method: "DELETE",
                    credentials: "include",
                }
            );

            if (!response.ok) throw new Error("Failed to delete");

            await invalidateAll();
        } catch (e) {
            alert("Failed to delete: " + e);
        }
    }
</script>

<main>
    {#if showEditor}
        <BoardSetupEditor
            initialSetup={editingSetup
                ? JSON.stringify(editingSetup)
                : undefined}
            onSave={handleEditorSave}
            onCancel={handleEditorCancel}
        />
    {:else}
        <div class="container">
            <div class="header">
                <button class="btn-primary" onclick={() => window.location.href = '/'}>Home</button>
                <h1>Board Setups</h1>
                <div class="header-info">
                    <span class="count"
                        >{data.setups.length}/{MAX_BOARD_SETUPS} setups</span
                    >
                    <button class="btn-primary" onclick={openCreateEditor} disabled={data.setups.length >= MAX_BOARD_SETUPS}>
                        + Create New Setup
                    </button>
                </div>
            </div>

            {#if error}
                <div class="error">{error}</div>
            {:else if data.setups.length === 0}
                <div class="empty">
                    <p>No board setups yet. Create your first setup!</p>
                </div>
            {:else}
                <div class="grid">
                    {#each data.setups as setup}
                        <SetupCard
                            {setup}
                            onEdit={openEditEditor}
                            onDelete={deleteSetup}
                        />
                    {/each}
                </div>
            {/if}
        </div>
    {/if}
</main>

<style>
    main {
        height: 100vh;
        overflow: hidden;
        width: 100%;
        background-image: url("$lib/assets/background-board-setup.png");
        background-size: cover;
        background-position: bottom;
    }
    .container {
        padding: 20px;
        max-width: 1200px;
        margin: 0 auto;
        height: 100%;
        display: flex;
        flex-direction: column;
    }

    .header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 24px;

        h1 {
            margin: 0;
            color: var(--text);
        }

        .header-info {
            display: flex;
            align-items: center;
            gap: 16px;

            .count {
                color: var(--muted, #999);
                font-weight: 500;
            }
        }
    }

    .error {
        background: #f44;
        color: white;
        padding: 16px;
        border-radius: 8px;
        margin-bottom: 20px;
    }

    .empty {
        text-align: center;
        padding: 60px 20px;
        color: var(--muted, #999);
    }

    .grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
        gap: 20px;
        overflow-y: auto;
        overflow-x: hidden;
        border-radius: 8px;
    }
</style>
