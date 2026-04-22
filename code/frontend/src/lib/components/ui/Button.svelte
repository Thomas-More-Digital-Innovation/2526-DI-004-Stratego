<script lang="ts">
    import type { Snippet } from "svelte";

    interface Props {
        type?: "button" | "submit" | "reset";
        variant?: "primary" | "secondary" | "ghost" | "outline";
        size?: "sm" | "md" | "lg";
        disabled?: boolean;
        loading?: boolean;
        onclick?: () => void;
        children: Snippet;
        class?: string;
    }

    let {
        type = "button",
        variant = "primary",
        size = "md",
        disabled = false,
        loading = false,
        onclick,
        children,
        class: className = "",
    }: Props = $props();

    const shared = `inline-flex items-center justify-center rounded-xl 
    transition-all duration-200 active:scale-95 cursor-pointer 
    disabled:opacity-50 disabled:pointer-events-none`;

    const variants = {
        primary: `bg-brand-primary border-b-4 border-black/30 hover:border-brand-primary 
            text-white shadow-lg glow-primary uppercase tracking-wider font-bold hover:translate-y-1`,
        secondary: `bg-brand-secondary border-b-4 border-black/30 hover:border-brand-secondary 
            text-white shadow-md uppercase tracking-wider font-bold hover:translate-y-1`,
        outline: `border-2 border-brand-accent/50 hover:border-brand-accent/50 shadow-xs hover:shadow-none shadow-brand-accent/50 text-brand-accent hover:bg-brand-accent/10 
            uppercase tracking-wider font-bold focus:ring-2 focus:ring-brand-accent hover:translate-y-0.5`,
        ghost: `hover:bg-brand-accent/10 text-brand-accent/80 hover:text-brand-accent 
            uppercase tracking-wider font-bold hover:translate-y-0.5`,
    };

    const sizes = {
        sm: "px-3 py-1.5 text-sm",
        md: "px-5 py-2.5 text-base",
        lg: "px-8 py-3.5 text-lg font-semibold",
    };
</script>

<button
    {type}
    disabled={disabled || loading}
    {onclick}
    class="{shared} {variants[variant]} {sizes[size]} {className} relative"
>
    {#if loading}
        <div class="absolute inset-0 flex items-center justify-center">
            <div
                class="w-4 h-4 border-2 border-white/20 border-t-current rounded-full animate-spin"
            ></div>
        </div>
        <span class="opacity-0">
            {@render children()}
        </span>
    {:else}
        {@render children()}
    {/if}
</button>
