<script lang="ts">
    import logoUrl from "$lib/assets/favicon.png";
    import tmUrl from "$lib/assets/tm_logo.png";
    import ErrorMessage from "$lib/components/ErrorMessage.svelte";
    import GameModes from "./GameModes.svelte";
    import LeftNavbar from "./LeftNavbar.svelte";

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
    <!-- <div class="saved-games">
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
    </div> -->
    <div class="left-side">
        <header>
            <span>
                <img src={logoUrl} alt="logo" width="64" />
                <h1>StrateGO</h1>
            </span>
        </header>
		<LeftNavbar />
        <footer>
            <img src={tmUrl} alt="Logo Thomas More" width="256" />
        </footer>
    </div>
    <GameModes bind:errorMessage />
</main>

<style>
    :global(body) {
        margin: 0;
        padding: 0;
        color: #e2e8f0;
        font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto,
            sans-serif;
    }

    main {
        height: 100vh;
        width: 100%;
        justify-content: space-between;
        display: flex;
        background: #1a202c4f;

        .left-side {
            display: flex;
            flex-direction: column;
            justify-content: space-between;
            background-image: url("$lib/assets/background.png");
            width: 100%;
        }
        header {
            text-align: center;
            margin-bottom: 50px;

            span {
                display: flex;
            }

            h1 {
                margin: 0;
                font-size: 3.5rem;
                color: #e2e8f0;
                text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.5);
            }
        }

        header,
        footer {
            padding: 8px;
        }
    }

    @media (max-width: 768px) {
        header h1 {
            font-size: 2.5rem;
        }
    }
</style>
