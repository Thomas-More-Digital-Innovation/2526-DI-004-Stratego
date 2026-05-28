<script lang="ts">
    // TODO: refactor this page
    import { goto } from "$app/navigation";
    import { authStore } from "$lib/state/auth.svelte";
    import { toastStore } from "$lib/state/toast.svelte";
    import Button from "$lib/components/ui/Button.svelte";
    import TacticalBriefing from "./_components/TacticalBriefing.svelte";
    import PasswordChecker from "./_components/PasswordChecker.svelte";

    let username = $state("");
    let password = $state("");
    let isLogin = $state(true);
    let loading = $state(false);
    let showPassword = $state(false);
    let isPasswordValid = $state(false);
    let passwordStrength = $state(0);
    let isBlueCommander = $state(true);

    // Reactive Username validation rules
    const isUsernameFormatValid = $derived(/^[a-zA-Z0-9_]*$/.test(username));
    const isUsernameLengthValid = $derived(
        username.length >= 3 && username.length <= 50,
    );
    const isUsernameValid = $derived(
        isUsernameFormatValid && isUsernameLengthValid,
    );

    const isFormValid = $derived(
        username.length > 0 &&
            isUsernameFormatValid &&
            (isLogin || isUsernameValid) &&
            password.length > 0 &&
            (isLogin || isPasswordValid),
    );

    async function handleSubmit() {
        if (!isFormValid) return;
        loading = true;

        try {
            if (isLogin) {
                await authStore.login(username, password);
            } else {
                await authStore.register(username, password);
            }
            goto("/");
        } catch (e: any) {
            toastStore.handleApiMessage(e, "Authentication failed");
        } finally {
            loading = false;
        }
    }
</script>

<svelte:head>
    <title>GoStrategy — {isLogin ? "Secure Login" : "Field Deployment"}</title>
</svelte:head>

{#snippet PillButton(
    onclick: () => void,
    text: string,
    disabled: boolean,
    isBlueCommander: boolean,
)}
    <button
        {onclick}
        {disabled}
        class="cursor-pointer flex-1 py-2.5 text-xs font-bold uppercase tracking-wider text-center z-10 transition-colors duration-200 {isBlueCommander
            ? 'text-brand-primary hover:text-brand-primary/60'
            : 'text-brand-secondary hover:text-brand-secondary/60'}"
    >
        {text}
    </button>
{/snippet}

<div class="flex items-center justify-center min-h-[80vh] p-2 sm:p-4">
    <!-- main dual-pane container -->
    <div
        class="w-full max-w-4xl grid grid-cols-1 md:grid-cols-12 glass rounded-3xl overflow-hidden shadow-2xl relative border border-white/10 animate-fade-in-slide"
    >
        <!-- decorative background ambient light blobs -->
        <div
            class="absolute -top-40 -left-40 w-80 h-80 {isBlueCommander
                ? 'bg-brand-primary/10'
                : 'bg-brand-secondary/10'} rounded-full blur-[100px] pointer-events-none"
        ></div>
        <div
            class="absolute -bottom-40 -right-40 w-80 h-80 {isBlueCommander
                ? 'bg-brand-secondary/10'
                : 'bg-brand-primary/10'} rounded-full blur-[100px] pointer-events-none"
        ></div>

        <!-- Left Pane: Tactical Briefing -->
        <TacticalBriefing bind:isBlueCommander />

        <!-- Right Pane: Authentication Form -->
        <div
            class="col-span-1 md:col-span-7 p-6 sm:p-10 flex flex-col justify-center relative z-10"
        >
            <div
                class="text-center md:text-left space-y-2 mb-6 animate-fade-in-slide delay-100 opacity-0"
            >
                <h1
                    class="text-2xl font-black text-white uppercase tracking-widest"
                >
                    {isLogin ? "Commander Who?" : "Enlisting Commander"}
                </h1>
                <p class="text-white/50 text-xs">
                    {isLogin
                        ? "Enter your secure credentials to command your military assets"
                        : "Construct your battle identity to start playing GoStrategy"}
                </p>
            </div>

            <!-- dynamic sliding pill switcher -->
            <div
                class="relative flex p-1 bg-black/40 border border-white/10 rounded-2xl mb-6 animate-fade-in-slide delay-200 opacity-0"
            >
                <div
                    class="absolute top-1 bottom-1 left-1 rounded-xl {isBlueCommander
                        ? 'bg-brand-primary/10 border-brand-primary/30'
                        : 'bg-brand-secondary/10 border-brand-secondary/30'} transition-all duration-300 ease-out shadow-inner"
                    style:width="calc(50% - 4px)"
                    style:transform="translateX({isLogin ? '0' : '100%'})"
                ></div>

                {@render PillButton(
                    () => {
                        isLogin = true;
                        showPassword = false;
                    },
                    "Sign In",
                    false,
                    isBlueCommander,
                )}
                {@render PillButton(
                    () => {
                        isLogin = false;
                        showPassword = false;
                    },
                    "Sign Up",
                    false,
                    isBlueCommander,
                )}
            </div>

            <!-- credentials form -->
            <form
                onsubmit={(e) => {
                    e.preventDefault();
                    handleSubmit();
                }}
                class="space-y-4 animate-fade-in-slide delay-300 opacity-0"
            >
                <!-- username group -->
                <div class="flex flex-col gap-1.5">
                    <label
                        for="username_input"
                        class="text-[10px] font-bold text-brand-accent uppercase tracking-widest ml-1"
                    >
                        Username
                    </label>
                    <div class="relative group">
                        <span
                            class="absolute left-4 top-1/2 -translate-y-1/2 flex items-center pointer-events-none"
                        >
                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                class="w-4.5 h-4.5 text-white/30 {isBlueCommander
                                    ? 'group-focus-within:text-brand-primary'
                                    : 'group-focus-within:text-brand-secondary'} transition-colors"
                                viewBox="0 0 24 24"
                                fill="none"
                                stroke="currentColor"
                                stroke-width="2"
                                stroke-linecap="round"
                                stroke-linejoin="round"
                            >
                                <path
                                    d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2"
                                />
                                <circle cx="12" cy="7" r="4" />
                            </svg>
                        </span>
                        <input
                            id="username_input"
                            type="text"
                            placeholder="Enter command username"
                            disabled={loading}
                            bind:value={username}
                            class="w-full bg-white/5 border transition-all duration-200 rounded-xl pl-11 pr-4 py-2.5 text-xs text-white placeholder:text-white/20 focus:outline-none focus:ring-2 {username.length >
                                0 &&
                            (!isUsernameFormatValid ||
                                (!isLogin && !isUsernameLengthValid))
                                ? 'border-red-500/50 focus:ring-red-500/35 focus:border-red-500/35'
                                : isBlueCommander
                                  ? 'border-white/10 focus:ring-brand-primary/45 focus:border-brand-primary/45'
                                  : 'border-white/10 focus:ring-brand-secondary/45 focus:border-brand-secondary/45'}"
                        />
                    </div>
                    <!-- username validation message -->
                    {#if username.length > 0 && (!isUsernameFormatValid || (!isLogin && !isUsernameLengthValid))}
                        <p
                            class="text-[9px] text-red-400 font-bold uppercase tracking-wider mt-1.5 ml-1 animate-fade-in"
                        >
                            {#if !isUsernameFormatValid}
                                Only letters, numbers, and underscores allowed
                            {:else if username.length < 3}
                                Username must be at least 3 characters
                            {:else if username.length > 50}
                                Username must be 50 characters or less
                            {/if}
                        </p>
                    {/if}
                </div>

                <!-- password group -->
                <div class="flex flex-col gap-1.5">
                    <label
                        for="password_input"
                        class="text-[10px] font-bold text-brand-accent uppercase tracking-widest ml-1"
                    >
                        Password
                    </label>
                    <div class="relative group">
                        <span
                            class="absolute left-4 top-1/2 -translate-y-1/2 flex items-center pointer-events-none"
                        >
                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                class="w-4.5 h-4.5 text-white/30 {isBlueCommander
                                    ? 'group-focus-within:text-brand-primary'
                                    : 'group-focus-within:text-brand-secondary'} transition-colors"
                                viewBox="0 0 24 24"
                                fill="none"
                                stroke="currentColor"
                                stroke-width="2"
                                stroke-linecap="round"
                                stroke-linejoin="round"
                            >
                                <rect
                                    x="3"
                                    y="11"
                                    width="18"
                                    height="11"
                                    rx="2"
                                    ry="2"
                                />
                                <path d="M7 11V7a5 5 0 0 1 10 0v4" />
                            </svg>
                        </span>
                        <input
                            id="password_input"
                            type={showPassword ? "text" : "password"}
                            placeholder="Enter password"
                            disabled={loading}
                            bind:value={password}
                            class="w-full bg-white/5 border transition-all duration-200 rounded-xl pl-11 pr-11 py-2.5 text-xs text-white placeholder:text-white/20 focus:outline-none focus:ring-2 {!isLogin &&
                            password.length > 0 &&
                            !isPasswordValid
                                ? 'border-red-500/50 focus:ring-red-500/35 focus:border-red-500/35'
                                : isBlueCommander
                                  ? 'border-white/10 focus:ring-brand-primary/45 focus:border-brand-primary/45'
                                  : 'border-white/10 focus:ring-brand-secondary/45 focus:border-brand-secondary/45'}"
                        />
                        <!-- password visibility toggle button -->
                        <button
                            type="button"
                            onclick={() => (showPassword = !showPassword)}
                            disabled={!password}
                            class="absolute right-3 top-1/2 -translate-y-1/2 text-white/30 hover:text-white active:scale-90 transition-all cursor-pointer disabled:opacity-0 disabled:pointer-events-none p-1 rounded-md hover:bg-white/5"
                        >
                            {#if showPassword}
                                <svg
                                    xmlns="http://www.w3.org/2000/svg"
                                    class="w-4 h-4"
                                    viewBox="0 0 24 24"
                                    fill="none"
                                    stroke="currentColor"
                                    stroke-width="2"
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                >
                                    <path
                                        d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"
                                    />
                                    <line x1="1" y1="1" x2="23" y2="23" />
                                </svg>
                            {:else}
                                <svg
                                    xmlns="http://www.w3.org/2000/svg"
                                    class="w-4 h-4"
                                    viewBox="0 0 24 24"
                                    fill="none"
                                    stroke="currentColor"
                                    stroke-width="2"
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                >
                                    <path
                                        d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"
                                    />
                                    <circle cx="12" cy="12" r="3" />
                                </svg>
                            {/if}
                        </button>
                    </div>
                </div>

                <!-- password strength checker (only in register mode) -->
                {#if !isLogin && password.length > 0}
                    <PasswordChecker
                        {password}
                        bind:isValid={isPasswordValid}
                        bind:strength={passwordStrength}
                    />
                {/if}

                <!-- submit authorization -->
                <div class="pt-2">
                    <Button
                        type="submit"
                        variant={isBlueCommander ? "primary" : "secondary"}
                        class="w-full relative overflow-hidden"
                        disabled={loading || !isFormValid}
                        disabledMessage={!isLogin &&
                        password.length > 0 &&
                        passwordStrength < 4
                            ? "Complete security standards first"
                            : "Enter credentials"}
                        {loading}
                    >
                        {isLogin
                            ? "Return to The Battlefield"
                            : "Deploy To The Battlefield"}
                    </Button>
                </div>
            </form>
        </div>
    </div>
</div>
