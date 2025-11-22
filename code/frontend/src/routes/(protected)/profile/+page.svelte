<script lang="ts">
    import { goto } from "$app/navigation";
    import type { PageData } from "./$types";

    let { data }: { data: PageData } = $props();

    async function handleLogout() {
        try {
            await fetch("http://localhost:8080/api/users/logout", {
                method: "POST",
                credentials: "include",
            });
            goto("/");
        } catch (e) {
            alert("Logout failed: " + e);
        }
    }

    function formatDuration(seconds: number): string {
        const mins = Math.floor(seconds / 60);
        const secs = Math.floor(seconds % 60);
        return `${mins}m ${secs}s`;
    }

    function getWinRate(): string {
        if (!data.stats || data.stats.total_games === 0) return "0%";
        return (
            ((data.stats.wins / data.stats.total_games) * 100).toFixed(1) + "%"
        );
    }
</script>

<main>
    <div class="profile-container">
        <div class="profile-card">
            <div class="profile-header">
                {#if data.user.profile_picture}
                    <img
                        src={data.user.profile_picture}
                        alt="Profile"
                        class="profile-pic"
                    />
                {:else}
                    <div class="profile-pic-placeholder">
                        {data.user.username[0]?.toUpperCase() || "?"}
                    </div>
                {/if}

                <div class="profile-info">
                    <h1>{data.user.username}</h1>
                    <p class="join-date">
                        Joined {new Date(
                            data.user.created_at
                        ).toLocaleDateString()}
                    </p>
                </div>

                <button class="logout-btn" onclick={handleLogout}>Logout</button
                >
            </div>

            {#if data.error}
                <div class="error">{data.error}</div>
            {:else if data.stats}
                <div class="stats-grid">
                    <div class="stat-card">
                        <div class="stat-value">{data.stats.total_games}</div>
                        <div class="stat-label">Total Games</div>
                    </div>

                    <div class="stat-card">
                        <div class="stat-value">{data.stats.wins}</div>
                        <div class="stat-label">Wins</div>
                    </div>

                    <div class="stat-card">
                        <div class="stat-value">{data.stats.losses}</div>
                        <div class="stat-label">Losses</div>
                    </div>

                    <div class="stat-card">
                        <div class="stat-value">{data.stats.draws}</div>
                        <div class="stat-label">Draws</div>
                    </div>

                    <div class="stat-card highlight">
                        <div class="stat-value">{getWinRate()}</div>
                        <div class="stat-label">Win Rate</div>
                    </div>

                    <div class="stat-card">
                        <div class="stat-value">{data.stats.total_moves}</div>
                        <div class="stat-label">Total Moves</div>
                    </div>

                    <div class="stat-card">
                        <div class="stat-value">
                            {formatDuration(
                                data.stats.avg_game_duration_seconds
                            )}
                        </div>
                        <div class="stat-label">Avg Game Duration</div>
                    </div>
                </div>
            {/if}

            <div class="actions">
                <a href="/board-setups" class="action-btn"
                    >Manage Board Setups</a
                >
                <a href="/" class="action-btn secondary">Home</a>
            </div>
        </div>
    </div>
</main>

<style>

	main {
		height: 100vh;
		overflow: hidden;
		width: 100%;
		background-image: url("$lib/assets/background-profile.png");
		background-size: cover;
		background-position: bottom;
	}
    .profile-container {
        padding: 20px;
        max-width: 900px;
        margin: 0 auto;
        min-height: 80vh;
        display: flex;
		height: 100%;
        align-items: center;
        justify-content: center;
    }

    .profile-card {
        background: var(--bg-accent-t);
        border-radius: 16px;
        box-shadow: 0 8px 24px rgba(0, 0, 0, 0.2);
        padding: 40px;
        width: 100%;
    }

    .profile-header {
        display: flex;
        align-items: center;
        gap: 24px;
        margin-bottom: 40px;
        padding-bottom: 24px;
        border-bottom: 1px solid #444;
    }

    .profile-pic {
        width: 100px;
        height: 100px;
        border-radius: 50%;
        object-fit: cover;
        border: 3px solid var(--primary);
    }

    .profile-pic-placeholder {
        width: 100px;
        height: 100px;
        border-radius: 50%;
        background: linear-gradient(135deg, var(--primary), var(--secondary));
        color: white;
        display: grid;
        place-items: center;
        font-size: 2.5rem;
        font-weight: 700;
    }

    .profile-info {
        flex: 1;
    }

    .profile-info h1 {
        margin: 0 0 8px 0;
        color: var(--text);
        font-size: 2rem;
    }

    .join-date {
        margin: 0;
        color: var(--muted, #999);
        font-size: 0.95rem;
    }

    .logout-btn {
        padding: 12px 24px;
        background: #d44;
        color: white;
        border: none;
        border-radius: 8px;
        font-weight: 600;
        cursor: pointer;
        transition: background 0.2s;
    }

    .logout-btn:hover {
        background: #c33;
    }

    .error {
        background: #f44;
        color: white;
        padding: 16px;
        border-radius: 8px;
        margin-bottom: 20px;
        text-align: center;
    }

    .stats-grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
        gap: 20px;
        margin-bottom: 32px;
    }

    .stat-card {
        background: var(--bg);
        padding: 24px;
        border-radius: 12px;
        text-align: center;
        transition: transform 0.2s;
    }

    .stat-card:hover {
        transform: translateY(-4px);
    }

    .stat-card.highlight {
        background: linear-gradient(135deg, var(--primary), var(--secondary));
        color: white;
    }

    .stat-value {
        font-size: 2.5rem;
        font-weight: 700;
        margin-bottom: 8px;
        color: var(--text);
    }

    .stat-card.highlight .stat-value,
    .stat-card.highlight .stat-label {
        color: white;
    }

    .stat-label {
        font-size: 0.95rem;
        color: var(--muted, #999);
        text-transform: uppercase;
        letter-spacing: 0.5px;
    }

    .actions {
        display: flex;
        gap: 16px;
        justify-content: center;
        padding-top: 24px;
        border-top: 1px solid #444;
    }

    .action-btn {
        padding: 14px 28px;
        background: var(--primary);
        color: white;
        border-radius: 8px;
        text-decoration: none;
        font-weight: 600;
        transition: all 0.2s;
        display: inline-block;
    }

    .action-btn:hover {
        background: var(--primary-dark);
        transform: translateY(-2px);
    }

    .action-btn.secondary {
        background: transparent;
        border: 2px solid var(--primary);
        color: var(--primary);
    }

    .action-btn.secondary:hover {
        background: var(--primary-dark);
        color: white;
    }

    @media (max-width: 640px) {
        .profile-header {
            flex-direction: column;
            text-align: center;
        }

        .stats-grid {
            grid-template-columns: repeat(2, 1fr);
        }

        .actions {
            flex-direction: column;
        }
    }
</style>
