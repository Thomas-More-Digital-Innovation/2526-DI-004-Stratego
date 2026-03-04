<script lang="ts">
    import { goto } from '$app/navigation';
    import logoUrl from "$lib/assets/favicon.png";
    import tmUrl from "$lib/assets/tm_logo.png";
    import ErrorMessage from "$lib/components/ErrorMessage.svelte";
    import GameModes from "./components/GameModes.svelte";
    import LeftNavbar from "./components/LeftNavbar.svelte";
    import type { PageData } from './$types';
    import AuthPanel from './components/AuthPanel.svelte';

    let { data }: { data: PageData } = $props();
    
    let errorMessage = $state<string>("");
</script>

<svelte:head>
    <title>Stratego - Menu</title>
</svelte:head>

<ErrorMessage {errorMessage} />

<main>
    <div class="left-side">
        <header>
            <span>
                <img src={logoUrl} alt="logo" width="64" />
                <h1>StrateGO</h1>
            </span>
        </header>
        <LeftNavbar isLoggedIn={!!data.user} />
        <footer>
            <img src={tmUrl} alt="Logo Thomas More" width="256" />
        </footer>
    </div>
    
    {#if data.user}
        <GameModes bind:errorMessage user={data.user} />
    {:else}
        <AuthPanel />
    {/if}
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
        max-height: 100vh;
        overflow: hidden;
        width: 100%;
        justify-content: space-between;
        display: flex;
        background-image: url("$lib/assets/background.png");

        .left-side {
            display: flex;
            flex-direction: column;
            justify-content: space-between;
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
