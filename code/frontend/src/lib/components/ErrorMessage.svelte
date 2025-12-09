<script lang="ts">
    let { errorMessage } = $props();

    let timeOutId: ReturnType<typeof setTimeout> | null = null;

    $effect(() => {
        if (errorMessage) {
            timeOutId = setTimeout(() => {
                errorMessage = "";
            }, 5000);
        }
    });

    function clearError() {
        if (timeOutId != null) {
            clearTimeout(timeOutId);
            timeOutId = null;
        }
        errorMessage = "";
    }

</script>

{#if errorMessage}
    <button class="error-banner" onclick={clearError}>
        ⚠️ {errorMessage}
        <span>X</span>
    </button>
{/if}

<style>
    .error-banner {
        cursor: pointer;
        background: #fc8181;
        color: white;
        padding: 15px;
        border-radius: 8px;
        margin-bottom: 30px;
        text-align: center;
        font-weight: 600;
        max-width: 800px;
        width: 100%;
        margin: auto;
        position: fixed;
        z-index: 100;
        left: 0;
        right: 0;
        bottom: 8px;

        span {
            color: white;
            position: absolute;
            top: 50%;
            transform: translateY(-50%);
            right: 8px;
            width: 32px;
            height: 32px;
            font-size: x-large;
            line-height: 24px;
            border-radius: 8px;
            text-align: center;
            background-color: transparent;
            border: none;
        }
    }
</style>
