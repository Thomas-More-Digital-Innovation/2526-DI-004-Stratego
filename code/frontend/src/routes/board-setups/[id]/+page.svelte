<script lang="ts">
    import { onMount } from "svelte";
    import { page } from "$app/state";
    import { goto } from "$app/navigation";
    import SetupEditor from "$lib/components/setup/SetupEditor.svelte";
    import Loading from "$lib/components/ui/Loading.svelte";
    import { boardSetups } from "$lib/api/client";
    import Card from "$lib/components/ui/Card.svelte";
    import Button from "$lib/components/ui/Button.svelte";
    import Input from "$lib/components/ui/Input.svelte";
    import { toastStore } from "$lib/state/toast.svelte";
    import type { BoardSetup } from "$lib/types/board-setup";

    let id = $derived(Number(page.params.id));
    let setup = $state<BoardSetup | null>(null);
    let name = $state("");
    let description = $state("");
    let isDefault = $state(false);
    let loading = $state(true);
    let saving = $state(false);
    let loadError = $state("");

    onMount(async () => {
        try {
            const found = await boardSetups.getOne(id);
            if (found) {
                setup = found;
                name = found.name;
                description = found.description;
                isDefault = found.is_default;
            } else {
                loadError = "Setup not found";
                toastStore.error("Setup not found");
            }
        } catch (e: any) {
            loadError = e.message || "Failed to load setup";
            toastStore.handleApiMessage(e, "Failed to load setup");
        } finally {
            loading = false;
        }
    });

    async function handleSave(setupData: string) {
        saving = true;
        try {
            await boardSetups.update(id, {
                name,
                description,
                setup_data: setupData,
                is_default: isDefault,
            });
            goto("/board-setups");
        } catch (e: any) {
            toastStore.handleApiMessage(e, "Failed to update setup");
            saving = false;
        }
    }
</script>

<svelte:head>
    <title>GoStrategy — Edit Setup</title>
</svelte:head>

<div class="max-w-6xl mx-auto space-y-8">
    <div class="flex flex-col md:flex-row md:items-end justify-between gap-4">
        <div class="flex-1 space-y-1">
            <h1
                class="text-3xl font-black text-white uppercase tracking-tighter"
            >
                Edit Setup
            </h1>
            <p class="text-white/40">
                Adjust your strategy for the upcoming battles.
            </p>
        </div>
    </div>

    {#if loading}
        <Loading
            title="Loading Setup"
            description="Fetching your strategic configuration..."
            subtitle="Synchronizing"
        />
    {:else if loadError && !setup}
        <Card class="text-center py-12">
            <p class="text-brand-secondary">{loadError}</p>
            <Button
                variant="ghost"
                class="mt-4"
                onclick={() => goto("/board-setups")}>Back to List</Button
            >
        </Card>
    {:else if setup}
        <Card class="space-y-6">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <Input
                    label="Setup Name"
                    placeholder="e.g., Aggressive Scout Rush"
                    bind:value={name}
                    sanitize="generic"
                />
                <Input
                    label="Description (Optional)"
                    placeholder="Briefly describe your strategy..."
                    bind:value={description}
                    sanitize="generic"
                />
            </div>

            <div class="flex items-center gap-3">
                <input
                    id="isDefault"
                    type="checkbox"
                    bind:checked={isDefault}
                    class="w-5 h-5 rounded border-white/10 bg-white/5 text-brand-accent focus:ring-brand-accent/50"
                />
                <label
                    for="isDefault"
                    class="text-sm font-medium text-white/70 cursor-pointer select-none"
                    >Set as default setup</label
                >
            </div>
        </Card>

        {#if saving}
            <Loading
                title="Saving Setup"
                description="Updating your encrypted battle plans..."
                subtitle="Encrypting"
            />
        {:else}
            <SetupEditor
                initialSetup={setup.setup_data}
                onSave={handleSave}
                onCancel={() => goto("/board-setups")}
            />
        {/if}
    {/if}
</div>
