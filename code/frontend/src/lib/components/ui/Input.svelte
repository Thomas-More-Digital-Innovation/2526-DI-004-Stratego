<script lang="ts">
    interface Props {
        type?: string;
        placeholder?: string;
        value?: string;
        label?: string;
        id?: string;
        class?: string;
    }

    let {
        type = "text",
        placeholder = "",
        value = $bindable(""),
        label = "",
        id = Math.random().toString(36).substring(7),
        class: className = "",
    }: Props = $props();
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
        bind:value
        class="w-full bg-white/5 border border-white/10 rounded-xl px-4 py-2.5 text-white placeholder:text-white/20 focus:outline-none focus:ring-2 focus:ring-brand-primary/50 focus:border-brand-primary/50 transition-all duration-200"
    />
</div>
