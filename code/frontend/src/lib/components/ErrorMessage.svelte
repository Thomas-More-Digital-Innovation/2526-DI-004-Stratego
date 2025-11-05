<script lang="ts">
    let { errorMessage } = $props();

</script>

{#if errorMessage}
    <div class="error-banner">
        ⚠️ {errorMessage}
        <button onclick={() => errorMessage = ''}>X</button>
    </div>
{/if}

<style>
    .error-banner {
        background: #fc8181;
        color: white;
        padding: 15px;
        border-radius: 8px;
        margin-bottom: 30px;
        text-align: center;
        font-weight: 600;
        width: 800px;
        margin: auto;
        position: fixed;
        left: 0;
        right: 0;
        bottom: 8px;

        button {
            position: absolute;
            top: 50%;
            transform: translateY(-50%);
            right: 8px;
            width: 32px;
            height: 32px;
            font-size: large;
            line-height: 24px;
            border: 1px solid black;
            border-radius: 8px;
            text-align: center;
        }
    }
</style>
