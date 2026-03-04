<script lang="ts">
    let { isLoggedIn = false }: { isLoggedIn?: boolean } = $props();
</script>

<nav class="left-navbar">
    <div>
        <a href="/profile" class:disabled={!isLoggedIn} onclick={(e) => !isLoggedIn && e.preventDefault()}>Profile</a>
        <a href="/board-setups" class:disabled={!isLoggedIn} onclick={(e) => !isLoggedIn && e.preventDefault()}>Board Setups</a>
    </div>
    <div>
        <a href="/replays" class:disabled={!isLoggedIn} onclick={(e) => !isLoggedIn && e.preventDefault()}>Replays</a>
        <a href="/settings" class:disabled={!isLoggedIn} onclick={(e) => !isLoggedIn && e.preventDefault()}>Settings</a>
    </div>
</nav>

<style>
    nav.left-navbar {
        width: clamp(200px, 50%, 400px);
        height: 100%;
        padding: 16px 0;

        display: flex;
        flex-direction: column;
        justify-content: space-between;

        div {
            display: flex;
            flex-direction: column;
            a {
                transition: 400ms;
                color: var(--text);
                font-size: x-large;
                text-decoration: none;
                font-weight: bold;
                padding: 8px 0;
                padding-left: 16px;
                margin: 8px 0;
                margin-right: 16px;
                background-color: var(--primary-dark);
                cursor: pointer;
                border: 4px solid var(--text);
                border-left: 0px;
                position: relative;

                border-bottom-right-radius: 100px;
                border-top-right-radius: 100px;

                &:hover:not(.disabled) {
                    transition: 400ms;
                    margin-right: 0px;
                    background-color: var(--primary);
                }

                &.disabled {
                    opacity: 0.5;
                    cursor: not-allowed;
                    background-color: var(--bg-accent);
                }
            }
        }
    }
</style>
