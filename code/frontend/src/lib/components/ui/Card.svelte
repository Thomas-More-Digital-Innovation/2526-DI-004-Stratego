<script lang="ts">
    import type { Snippet } from "svelte";

    interface Props {
        children: Snippet;
        class?: string;
        padding?: "none" | "sm" | "md" | "lg";
    }

    let { children, class: className = "", padding = "md" }: Props = $props();

    const paddings = {
        none: "",
        sm: "p-3",
        md: "p-6",
        lg: "p-10",
    };
</script>

<div class="glass rounded-2xl overflow-hidden {paddings[padding]} {className}">
    {@render children()}
</div>
