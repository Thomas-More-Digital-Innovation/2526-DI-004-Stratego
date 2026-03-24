<script lang="ts">
    import type { Snippet } from "svelte";

    interface Props {
        type?: "button" | "submit" | "reset";
        variant?: "primary" | "secondary" | "ghost" | "outline";
        size?: "sm" | "md" | "lg";
        disabled?: boolean;
        onclick?: () => void;
        children: Snippet;
        class?: string;
    }

    let {
        type = "button",
        variant = "primary",
        size = "md",
        disabled = false,
        onclick,
        children,
        class: className = "",
    }: Props = $props();

    const variants = {
        primary: `bg-brand-primary border-b-4 border-black/30 hover:border-brand-primary hover:translate-y-1 
            text-white shadow-lg glow-primary uppercase tracking-wider font-bold`,
        secondary: `bg-brand-secondary border-b-4 border-black/30 hover:border-brand-secondary hover:translate-y-1 
            text-white shadow-md uppercase tracking-wider font-bold`,
        outline: `border-2 border-brand-accent/50 hover:translate-y-1 hover:border-brand-accent/50 shadow-xs hover:shadow-none shadow-brand-accent/50 text-brand-accent hover:bg-brand-accent/10 
            uppercase tracking-wider font-bold focus:ring-2 focus:ring-brand-accent`,
        ghost: `hover:bg-brand-accent/10 text-brand-accent/80 hover:text-brand-accent 
            uppercase tracking-wider font-bold`,
    };

    const sizes = {
        sm: "px-3 py-1.5 text-sm",
        md: "px-5 py-2.5 text-base",
        lg: "px-8 py-3.5 text-lg font-semibold",
    };
</script>

<button
    {type}
    {disabled}
    {onclick}
    class="inline-flex items-center justify-center rounded-xl transition-all duration-200 active:scale-95 cursor-pointer disabled:opacity-50 disabled:pointer-events-none {variants[
        variant
    ]} {sizes[size]} {className}"
>
    {@render children()}
</button>
