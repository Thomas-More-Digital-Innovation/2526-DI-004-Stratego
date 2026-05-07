<script lang="ts" module>
    declare const __VERSION__: string;
</script>

<script lang="ts">
    import Changelog from "./Changelog.svelte";
    import { fade, scale } from "svelte/transition";

    let isOpen = $state(false);

    function toggleModal() {
        isOpen = !isOpen;
    }

    function handleKeydown(e: KeyboardEvent) {
        if (e.key === "Escape" && isOpen) {
            isOpen = false;
        }
    }
</script>

<svelte:window onkeydown={handleKeydown} />

<button
    onclick={toggleModal}
    class="fixed right-4 top-4 z-40 cursor-pointer rounded-full border border-white/10 bg-black/40 px-4 py-2 text-sm font-medium text-white/70 backdrop-blur-md transition-all hover:border-white/20 hover:bg-black/60 hover:text-white"
>
    v{__VERSION__} — Changelog
</button>

{#if isOpen}
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <div
        transition:fade={{ duration: 200 }}
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/80 p-4 backdrop-blur-sm"
        onclick={toggleModal}
    >
        <div
            transition:scale={{ duration: 300, start: 0.95 }}
            class="relative flex h-full max-h-[80vh] w-full max-w-2xl flex-col overflow-hidden rounded-2xl border border-white/10 bg-[#0a0a0a] shadow-2xl"
            onclick={(e) => e.stopPropagation()}
        >
            <div
                class="flex items-center justify-between border-b border-white/10 p-6"
            >
                <h2 class="text-xl font-bold">Changelog</h2>
                <button
                    aria-label="Close changelog"
                    onclick={toggleModal}
                    class="rounded-lg p-2 text-white/50 transition-colors hover:bg-white/5 hover:text-white"
                >
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        width="20"
                        height="20"
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

            <div class="flex-1 overflow-y-auto p-8 pt-4 custom-scrollbar">
                <Changelog />
            </div>

            <div
                class="border-t border-white/10 p-4 text-center text-xs text-white/50"
            >
                GoStrategy v{__VERSION__} — THOMAS MORE DIGITAL INNOVATION
                <br />
                Made by Sem Van Broekhoven
            </div>
        </div>
    </div>
{/if}

<style>
    .custom-scrollbar::-webkit-scrollbar {
        width: 8px;
    }
    .custom-scrollbar::-webkit-scrollbar-track {
        background: transparent;
    }
    .custom-scrollbar::-webkit-scrollbar-thumb {
        background: rgba(255, 255, 255, 0.1);
        border-radius: 4px;
    }
    .custom-scrollbar::-webkit-scrollbar-thumb:hover {
        background: rgba(255, 255, 255, 0.2);
    }
</style>
