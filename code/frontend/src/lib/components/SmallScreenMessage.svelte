<script lang="ts">
    import background from "$lib/assets/background.webp";
    let screenWidth = window.innerWidth;

    function handleResize() {
        screenWidth = window.innerWidth;
    }

    window.addEventListener("resize", handleResize);
</script>

{#if screenWidth < 1200}
    <div
        class="absolute top-0 left-0 p-8 text-center flex flex-col items-center justify-center z-50 h-screen w-screen bg-surface-base"
    >
        <img
            src={background}
            alt=""
            class="absolute top-0 left-0 w-screen h-screen object-cover blur-xs opacity-50"
        />
        <div class="z-10">
            <h1 class="text-4xl mb-4">Sorry! ☹️</h1>
            <p class="text-lg mb-2">
                This app is designed for larger screens to ensure the best
                experience.
            </p>
            <p class="text-lg">
                Please try again on a larger device or change the zoom level of
                your browser. (Ctrl + Scroll)
            </p>
            <p class="text-md mt-4">
                Thanks for understanding! I'm working hard on making it
                accessible for smaller screens too. - Sem
            </p>
        </div>
    </div>
{/if}
