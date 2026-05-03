<script lang="ts">
    import { fly, fade } from "svelte/transition";
    import type { ToastType } from "$lib/state/toast.svelte";
    import { toastStore } from "$lib/state/toast.svelte";

    interface Props {
        id: string;
        message: string;
        type: ToastType;
    }

    let { id, message, type }: Props = $props();

    const config = {
        success: {
            bg: "bg-emerald-500/10",
            border: "border-emerald-500/20",
            text: "text-emerald-400",
            icon: `<path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/>`,
        },
        error: {
            bg: "bg-rose-500/10",
            border: "border-rose-500/20",
            text: "text-rose-400",
            icon: `<circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/>`,
        },
        warning: {
            bg: "bg-amber-500/10",
            border: "border-amber-500/20",
            text: "text-amber-400",
            icon: `<path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3Z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/>`,
        },
        info: {
            bg: "bg-blue-500/10",
            border: "border-blue-500/20",
            text: "text-blue-400",
            icon: `<circle cx="12" cy="12" r="10"/><line x1="12" y1="16" x2="12" y2="12"/><line x1="12" y1="8" x2="12.01" y2="8"/>`,
        },
    };

    const current = $derived(config[type]);
</script>

<div
    in:fly={{ y: 20, duration: 300 }}
    out:fade={{ duration: 200 }}
    class="pointer-events-auto flex items-center gap-3 rounded-2xl border {current.border} {current.bg} p-4 {current.text} shadow-2xl backdrop-blur-xl"
>
    <svg
        xmlns="http://www.w3.org/2000/svg"
        width="18"
        height="18"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
    >
        {@html current.icon}
    </svg>
    <p class="text-sm font-medium leading-tight">{message}</p>
    <button
        aria-label="Close toast"
        onclick={() => toastStore.remove(id)}
        class="ml-auto rounded-lg p-1 opacity-50 transition-all hover:bg-white/5 hover:opacity-100"
    >
        <svg
            xmlns="http://www.w3.org/2000/svg"
            width="14"
            height="14"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
            ><path d="M18 6 6 18" /><path d="m6 6 12 12" /></svg
        >
    </button>
</div>
