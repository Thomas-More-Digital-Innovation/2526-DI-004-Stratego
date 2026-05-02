<script lang="ts">
    import { page } from "$app/state";

    let showFullDisclaimer = $derived(page.url.pathname === "/");
</script>

<div
    class="group fixed bottom-2 right-2 p-4 text-red-300 text-center bg-red-500/10 rounded-md"
>
    <p class="text-sm">This game is early access, expect bugs.</p>
    <div class={showFullDisclaimer ? "block" : "hidden group-hover:block"}>
        <p class="text-xs">Loss of data is possible. Thank you for testing!</p>
        <p class="text-sm">
            Report any bugs <a
                class="underline text-brand-primary hover:text-brand-accent"
                href="https://github.com/Thomas-More-Digital-Innovation/2526-DI-004-Stratego/issues"
                >here</a
            >.
        </p>
    </div>
</div>
