<script lang="ts">
    let { state = $bindable() }: { state: any } = $props();
</script>

<div class="field">
    <label for="name">Name</label>
    <input
        id="name"
        type="text"
        bind:value={state.setupName}
        placeholder="My Setup"
    />
</div>

<div class="field">
    <label for="description">Description</label>
    <textarea
        id="description"
        bind:value={state.setupDescription}
        placeholder="Optional description"
        rows="3"
    ></textarea>
</div>

<div class="field-checkbox">
    <input id="default" type="checkbox" bind:checked={state.isDefault} />
    <label for="default">Set as default</label>
</div>

<style>
    .field {
        margin-bottom: 16px;

        label {
            display: block;
            margin-bottom: 6px;
            color: var(--text);
            font-weight: 500;
        }

        input,
        textarea {
            width: 100%;
            padding: 10px;
            border: 1px solid #444;
            border-radius: 6px;
            background: var(--bg-accent);
            color: var(--text);
            font-size: 1rem;
            font-family: inherit;
        }

        textarea {
            resize: vertical;
        }
    }

    .field-checkbox {
        display: flex;
        align-items: center;
        gap: 8px;
        margin-bottom: 20px;

        input {
            width: auto;
        }
    }
</style>