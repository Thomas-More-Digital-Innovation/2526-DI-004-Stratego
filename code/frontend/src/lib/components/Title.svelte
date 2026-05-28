<script lang="ts">
    import goLogo from "$lib/assets/go.svg";

    type size = "sm" | "md" | "lg" | "xl" | "2xl";

    const sizes = {
        sm: 30,
        md: 32,
        lg: 36,
        xl: 40,
        "2xl": 64,
    };

    let { size = "xl" }: { size?: size } = $props();
</script>

<div
    class="text-{size} font-extrabold tracking-widest text-white flex gap-1 items-center drop-shadow-md"
>
    <span class="flex items-center"
        ><a href="https://go.dev/">
            <img
                src={goLogo}
                alt="go"
                class="mr-1"
                width={sizes[size]}
                height={sizes[size]}
            /></a
        >STRATEGY</span
    >
</div>
