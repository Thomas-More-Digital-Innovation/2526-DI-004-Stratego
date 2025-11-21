<script lang="ts">
    let username = $state("");
    let password = $state("");
    let isLogin = $state(true);
    let authError = $state("");
    async function handleAuth() {
        authError = "";
        const endpoint = isLogin ? "/api/users/login" : "/api/users/register";

        try {
            const response = await fetch(`http://localhost:8080${endpoint}`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                credentials: "include",
                body: JSON.stringify({ username, password }),
            });

            if (!response.ok) {
                const text = await response.text();
                authError = text || "Authentication failed";
                return;
            }

            // Redirect to reload page with user data
            window.location.href = "/";
        } catch (e) {
            authError = "Network error: " + e;
        }
    }
</script>

<div class="auth-section">
    <div class="auth-card">
        <h2>{isLogin ? "Login" : "Register"} to Play</h2>
        <p class="auth-description">
            {#if isLogin}
                Sign in to track your stats and save your board setups
            {:else}
                Create an account to start playing Stratego
            {/if}
        </p>

        {#if authError}
            <div class="error">{authError}</div>
        {/if}

        <form
            onsubmit={(e) => {
                e.preventDefault();
                handleAuth();
            }}
        >
            <div class="field">
                <label for="username">Username</label>
                <input
                    id="username"
                    type="text"
                    bind:value={username}
                    placeholder="Enter username"
                    required
                />
            </div>

            <div class="field">
                <label for="password">Password</label>
                <input
                    id="password"
                    type="password"
                    bind:value={password}
                    placeholder="Enter password"
                    required
                />
            </div>

            <button type="submit" class="submit-btn">
                {isLogin ? "Login" : "Register"}
            </button>
        </form>

        <button
            type="button"
            class="toggle-btn"
            onclick={() => {
                isLogin = !isLogin;
                authError = "";
            }}
        >
            {isLogin
                ? "Don't have an account? Register"
                : "Already have an account? Login"}
        </button>
    </div>
</div>

<style>
    .auth-section {
        background: var(--bg);
        padding: 32px;
        display: flex;
        align-items: center;
        justify-content: center;
    }
    
    .auth-card {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(300px, 700px));
        background: var(--bg-accent);
        padding: 16px;
        border-radius: 12px;
        box-shadow: 0 8px 24px rgba(0, 0, 0, 0.2);
        max-width: 400px;
        width: 100%;

        h2 {
            margin: 0 0 8px 0;
            text-align: center;
            color: var(--text);
        }
    }

    .auth-description {
        text-align: center;
        color: var(--muted, #999);
        margin-bottom: 24px;
        font-size: 0.95rem;
    }

    .error {
        background: #f44;
        color: white;
        padding: 12px;
        border-radius: 6px;
        margin-bottom: 16px;
        text-align: center;
    }

    .field {
        margin-bottom: 20px;

        label {
            display: block;
            margin-bottom: 6px;
            color: var(--text);
            font-weight: 500;
        }

        input {
            width: 100%;
            padding: 11px;
            border: 1px solid #444;
            border-radius: 6px;
            background: var(--bg);
            color: var(--text);
            font-size: 1rem;

            &:focus {
                outline: none;
                border-color: var(--secondary);
            }
        }
    }

    .submit-btn {
        width: 100%;
        padding: 14px;
        background: var(--primary);
        color: white;
        border: none;
        border-radius: 6px;
        font-size: 1.05rem;
        font-weight: 600;
        cursor: pointer;
        transition: background 0.2s;

        &:hover {
            background: var(--primary-dark);
        }
    }

    .toggle-btn {
        width: 100%;
        margin-top: 16px;
        padding: 10px;
        background: transparent;
        color: var(--text);
        border: 1px solid #444;
        border-radius: 6px;
        cursor: pointer;
        transition: all 0.2s;

        &:hover {
            border-color: var(--secondary);
            color: var(--secondary);
        }
    }
</style>
