<script lang="ts">
    import { goto } from "$app/navigation";
    import { authStore } from "$lib/state/auth.svelte";
    import Card from "$lib/components/ui/Card.svelte";
    import Input from "$lib/components/ui/Input.svelte";
    import Button from "$lib/components/ui/Button.svelte";

    let username = $state("");
    let password = $state("");
    let isLogin = $state(true);
    let error = $state("");
    let loading = $state(false);

    async function handleSubmit() {
        error = "";
        loading = true;

        try {
            if (isLogin) {
                await authStore.login(username, password);
            } else {
                await authStore.register(username, password);
            }
            goto("/");
        } catch (e: any) {
            error = e.message || "Authentication failed";
        } finally {
            loading = false;
        }
    }
</script>

<svelte:head>
    <title>Stratego — {isLogin ? "Login" : "Register"}</title>
</svelte:head>

<div class="flex items-center justify-center min-h-[80vh]">
    <Card class="w-full max-w-md space-y-6">
        <div class="text-center space-y-2">
            <h1
                class="text-2xl font-extrabold text-brand-accent uppercase tracking-widest"
            >
                {isLogin ? "Login" : "Register"}
            </h1>
            <p class="text-white/50 text-sm">
                {isLogin
                    ? "Sign in to track your stats and save board setups"
                    : "Create an account to start playing Stratego"}
            </p>
        </div>

        {#if error}
            <div
                class="bg-brand-secondary/20 border border-brand-secondary/30 text-brand-secondary rounded-xl px-4 py-3 text-sm text-center"
            >
                {error}
            </div>
        {/if}

        <form
            onsubmit={(e) => {
                e.preventDefault();
                handleSubmit();
            }}
            class="space-y-4"
        >
            <Input
                label="Username"
                placeholder="Enter your username"
                bind:value={username}
            />
            <Input
                label="Password"
                type="password"
                placeholder="Enter your password"
                bind:value={password}
            />
            <Button
                type="submit"
                variant="primary"
                class="w-full"
                disabled={loading || !username || !password}
            >
                {#if loading}
                    Processing...
                {:else}
                    {isLogin ? "Login" : "Register"}
                {/if}
            </Button>
        </form>

        <Button
            variant="ghost"
            class="w-full"
            onclick={() => {
                isLogin = !isLogin;
                error = "";
            }}
        >
            {isLogin
                ? "Don't have an account? Register"
                : "Already have an account? Login"}
        </Button>
    </Card>
</div>
