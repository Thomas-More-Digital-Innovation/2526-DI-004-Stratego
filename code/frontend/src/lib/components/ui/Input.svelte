<script lang="ts">
    interface Props {
        type?: string;
        placeholder?: string;
        value?: string;
        label?: string;
        id?: string;
        class?: string;
        disabled?: boolean;
        sanitize?: "username" | "password" | "generic";
    }

    let {
        type = "text",
        placeholder = "",
        value = $bindable(""),
        label = "",
        id = Math.random().toString(36).substring(7),
        class: className = "",
        disabled = false,
        sanitize = undefined,
    }: Props = $props();

    function handleInput(e: Event) {
        const input = e.target as HTMLInputElement;
        let val = input.value;

        if (sanitize === "username") {
            val = val.replace(/[^a-zA-Z0-9_]/g, "");
        } else if (sanitize === "password") {
            val = val.replace(/[^a-zA-Z0-9!@#$%^&*()_+=\-.]/g, "");
        } else if (sanitize === "generic") {
            val = val.replace(/[<>"'%;]/g, "");
        }

        if (val !== input.value) {
            input.value = val;
            value = val;
        }
    }
</script>

<div class="flex flex-col gap-1.5 {className}">
    {#if label}
        <label
            for={id}
            class="text-xs font-bold text-brand-accent uppercase tracking-widest ml-1"
        >
            {label}
        </label>
    {/if}
    <input
        {id}
        {type}
        {placeholder}
        {disabled}
        bind:value
        oninput={handleInput}
        class="w-full bg-white/5 border border-white/10 rounded-xl px-4 py-2.5 text-white placeholder:text-white/20 focus:outline-none focus:ring-2 focus:ring-brand-primary/50 focus:border-brand-primary/50 transition-all duration-200"
    />
</div>
