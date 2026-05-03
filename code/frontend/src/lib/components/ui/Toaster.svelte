<script lang="ts">
    import { toastStore } from "$lib/state/toast.svelte";
    import Toast from "./Toast.svelte";
</script>

<div
    class="pointer-events-none fixed bottom-0 right-0 z-1000 flex w-full max-w-sm flex-col gap-3 p-6"
>
    {#each toastStore.toasts as toast (toast.id)}
        <Toast {...toast} />
    {/each}
</div>
