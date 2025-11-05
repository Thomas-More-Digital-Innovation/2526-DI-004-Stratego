<script lang="ts">
    import logoUrl from "$lib/assets/favicon.png";
    import Footer from "$lib/components/Footer.svelte";
    import ErrorMessage from "$lib/components/ErrorMessage.svelte";
    import GameModes from "./GameModes.svelte";

    let errorMessage = $state<string>("");
    let loadedGameData = $state<string | null>(null);
    let fileInput: HTMLInputElement;

    function handleFileSelect(event: Event) {
        const target = event.target as HTMLInputElement;
        const file = target.files?.[0];

        if (file) {
            const reader = new FileReader();
            reader.onload = (e) => {
                try {
                    const content = e.target?.result as string;
                    loadedGameData = content;
                    alert("Game loaded! (Replay feature coming soon)");
                } catch (error) {
                    errorMessage = "Failed to load game file";
                }
            };
            reader.readAsText(file);
        }
    }
</script>

<svelte:head>
    <title>Stratego - Menu</title>
</svelte:head>

<ErrorMessage {errorMessage} />

<main>
    <div class="saved-games">
        <h2>Load Saved Game</h2>
        <p class="description">
            Load a previously saved game to review the moves and strategy.
        </p>

        <button class="load-btn" onclick={() => fileInput.click()}>
            📂 Load Game File
        </button>

        <input
            bind:this={fileInput}
            type="file"
            accept=".json"
            onchange={handleFileSelect}
            style="display: none;"
        />
    </div>
    <header>
        <span>
            <img src={logoUrl} alt="logo" width="128" />
            <h1>StrateGO</h1>
        </span>
        <p class="subtitle">Choose your game mode and start playing</p>
    </header>
    <GameModes bind:errorMessage/>
    <!-- <Footer /> -->
</main>

<style>
    :global(body) {
        margin: 0;
        padding: 0;
        color: #e2e8f0;
        font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto,
		sans-serif;
		background-image: url("$lib/assets/background.png");
    }
	
    main {
		height: 100vh;
		width: 100%;
		justify-content: space-between;
        display: flex;
		background: #1a202c4f;
    }

    header {
        text-align: center;
        margin-bottom: 50px;
    }

    header h1 {
        margin: 0;
        font-size: 3.5rem;
        color: #e2e8f0;
        text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.5);
        margin-bottom: 10px;
    }

    .subtitle {
        color: #a0aec0;
        font-size: 1.2rem;
        margin: 0;
    }

    .menu-content {
        display: flex;
        gap: 40px;
    }

    h2 {
        margin: 0 0 20px 0;
        color: #e2e8f0;
        font-size: 1.8rem;
        border-bottom: 2px solid #4a5568;
        padding-bottom: 10px;
    }

    .description {
        color: #a0aec0;
        margin: 0 0 20px 0;
    }

    .load-btn {
        width: 100%;
        padding: 15px;
        background: #4299e1;
        color: white;
        border: none;
        border-radius: 8px;
        font-size: 1.1rem;
        font-weight: 600;
        cursor: pointer;
        transition: all 0.2s;
    }

    .load-btn:hover {
        background: #3182ce;
        transform: translateY(-2px);
        box-shadow: 0 4px 8px rgba(66, 153, 225, 0.4);
    }

    @media (max-width: 768px) {
        header h1 {
            font-size: 2.5rem;
        }

        .mode-cards {
            grid-template-columns: 1fr;
        }
    }
</style>
