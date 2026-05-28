<script lang="ts">
    import {
        validatePasswordInput,
        getPasswordStrengthScore,
        isPasswordStrong,
    } from "$lib/utils/authValidation";

    interface Props {
        password: string;
        isValid?: boolean;
        strength?: number;
    }

    let {
        password,
        isValid = $bindable(false),
        strength = $bindable(0),
    }: Props = $props();

    const checks = $derived(validatePasswordInput(password));
    const strengthScore = $derived(getPasswordStrengthScore(password));
    const isStrong = $derived(isPasswordStrong(password));

    // Reactively update parents bound variables
    $effect(() => {
        isValid = isStrong;
        strength = strengthScore;
    });
</script>

<div
    class="space-y-2 animate-fade-in py-2.5 px-3 bg-black/45 border border-white/10 rounded-2xl"
>
    <div
        class="text-[9px] uppercase tracking-widest text-white/40 font-bold border-b border-white/5 pb-1.5 mb-2"
    >
        Tactical Security Protocol
    </div>
    <div
        class="grid grid-cols-1 sm:grid-cols-2 gap-x-3 gap-y-2 text-xs font-semibold"
    >
        <!-- 1. Length -->
        <div
            class="flex items-center gap-2 transition-colors duration-200 {checks.isMinLength
                ? 'text-emerald-400'
                : 'text-white/40'}"
        >
            {#if checks.isMinLength}
                <svg
                    class="w-3.5 h-3.5 text-emerald-400 shrink-0"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                    stroke-width="3"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        d="M5 13l4 4L19 7"
                    />
                </svg>
            {:else}
                <svg
                    class="w-3.5 h-3.5 text-red-500/70 shrink-0"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                    stroke-width="3"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        d="M6 18L18 6M6 6l12 12"
                    />
                </svg>
            {/if}
            <span class="text-[10px] tracking-wide uppercase"
                >8+ Characters</span
            >
        </div>

        <!-- 2. Uppercase -->
        <div
            class="flex items-center gap-2 transition-colors duration-200 {checks.hasUppercase
                ? 'text-emerald-400'
                : 'text-white/40'}"
        >
            {#if checks.hasUppercase}
                <svg
                    class="w-3.5 h-3.5 text-emerald-400 shrink-0"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                    stroke-width="3"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        d="M5 13l4 4L19 7"
                    />
                </svg>
            {:else}
                <svg
                    class="w-3.5 h-3.5 text-red-500/70 shrink-0"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                    stroke-width="3"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        d="M6 18L18 6M6 6l12 12"
                    />
                </svg>
            {/if}
            <span class="text-[10px] tracking-wide uppercase"
                >Uppercase Letter (A-Z)</span
            >
        </div>

        <!-- 3. Lowercase -->
        <div
            class="flex items-center gap-2 transition-colors duration-200 {checks.hasLowercase
                ? 'text-emerald-400'
                : 'text-white/40'}"
        >
            {#if checks.hasLowercase}
                <svg
                    class="w-3.5 h-3.5 text-emerald-400 shrink-0"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                    stroke-width="3"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        d="M5 13l4 4L19 7"
                    />
                </svg>
            {:else}
                <svg
                    class="w-3.5 h-3.5 text-red-500/70 shrink-0"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                    stroke-width="3"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        d="M6 18L18 6M6 6l12 12"
                    />
                </svg>
            {/if}
            <span class="text-[10px] tracking-wide uppercase"
                >Lowercase Letter (a-z)</span
            >
        </div>

        <!-- 4. Number -->
        <div
            class="flex items-center gap-2 transition-colors duration-200 {checks.hasNumber
                ? 'text-emerald-400'
                : 'text-white/40'}"
        >
            {#if checks.hasNumber}
                <svg
                    class="w-3.5 h-3.5 text-emerald-400 shrink-0"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                    stroke-width="3"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        d="M5 13l4 4L19 7"
                    />
                </svg>
            {:else}
                <svg
                    class="w-3.5 h-3.5 text-red-500/70 shrink-0"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                    stroke-width="3"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        d="M6 18L18 6M6 6l12 12"
                    />
                </svg>
            {/if}
            <span class="text-[10px] tracking-wide uppercase"
                >Numeric Digit (0-9)</span
            >
        </div>

        <!-- 5. Format -->
        <div
            class="flex items-center gap-2 transition-colors duration-200 {checks.isValidFormat
                ? 'text-emerald-400'
                : 'text-white/40'} col-span-1 sm:col-span-2 border-t border-white/5 pt-1.5 mt-0.5"
        >
            {#if checks.isValidFormat}
                <svg
                    class="w-3.5 h-3.5 text-emerald-400 shrink-0"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                    stroke-width="3"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        d="M5 13l4 4L19 7"
                    />
                </svg>
            {:else}
                <svg
                    class="w-3.5 h-3.5 text-red-500/70 shrink-0"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                    stroke-width="3"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        d="M6 18L18 6M6 6l12 12"
                    />
                </svg>
            {/if}
            <span class="text-[10px] tracking-wide uppercase"
                >Valid Characters Only</span
            >
        </div>
    </div>
</div>
