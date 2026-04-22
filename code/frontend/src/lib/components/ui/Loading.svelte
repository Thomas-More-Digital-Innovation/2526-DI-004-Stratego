<script lang="ts">
    import { onMount } from "svelte";
    import fafo from "$lib/assets/ai/fafo.webp";
    import fato from "$lib/assets/ai/fato.webp";

    interface Props {
        title?: string;
        description?: string;
        subtitle?: string;
    }

    let {
        title = "Something is happening",
        description = "",
        subtitle = "Shhh... The computer is thinking",
    }: Props = $props();

    const images = [fafo, fato];
    let selectedImage = $state(fafo);

    onMount(() => {
        selectedImage = images[Math.floor(Math.random() * images.length)];
    });
</script>

<div
    class="flex flex-col items-center justify-center min-h-[50vh] gap-8 animate-in fade-in zoom-in-95 duration-700"
>
    <!-- Image Container with Glow -->
    <div class="relative group">
        <div
            class="absolute -inset-4 bg-brand-primary/20 rounded-full blur-2xl opacity-50 group-hover:opacity-80 transition-opacity duration-1000 animate-pulse"
        ></div>
        <div
            class="relative w-48 h-48 md:w-64 md:h-64 rounded-full overflow-hidden border-4 border-white/10 shadow-2xl z-10"
        >
            <img
                src={selectedImage}
                alt={subtitle}
                class="w-full h-full object-contain animate-spin"
            />
        </div>

        <!-- Decorative Rings -->
        <div
            class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[110%] h-[110%] border border-white/5 rounded-full animate-ping-slow"
        ></div>
        <div
            class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[125%] h-[125%] border border-white/5 rounded-full animate-ping-slower"
        ></div>
    </div>

    <!-- Text Content -->
    <div class="text-center space-y-3 max-w-sm">
        <div class="space-y-1">
            <h2
                class="text-3xl font-black text-white uppercase tracking-tighter"
            >
                {title}
            </h2>
            <p class="text-white/40 text-sm font-medium tracking-wide">
                {description}
            </p>
        </div>

        <div class="flex items-center justify-center gap-3 pt-4">
            <div
                class="h-px w-8 bg-linear-to-r from-transparent to-white/20"
            ></div>
            <span
                class="text-[10px] font-black text-brand-accent uppercase tracking-[0.3em] animate-pulse"
            >
                {subtitle}
            </span>
            <div
                class="h-px w-8 bg-linear-to-l from-transparent to-white/20"
            ></div>
        </div>
    </div>
</div>
