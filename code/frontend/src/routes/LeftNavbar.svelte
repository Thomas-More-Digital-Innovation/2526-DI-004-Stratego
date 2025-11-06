<script lang="ts">
</script>

<nav class="left-navbar">
    <div>
        <a href="/profile">Profile</a>
        <a href="/board-setup">Board Setup</a>
    </div>
    <div>
        <a href="/replays">Replays</a>
        <a href="/settings">Settings</a>
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

                &:hover {
                    transition: 400ms;
                    margin-right: 0px;
                    background-color: var(--primary);
                }
            }
        }
    }
</style>
