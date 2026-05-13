<script lang="ts">
    import { goto } from "$app/navigation";
    import SetupEditor from "$lib/components/setup/SetupEditor.svelte";
    import { boardSetups } from "$lib/api/client";
    import { toastStore } from "$lib/state/toast.svelte";
    import Card from "$lib/components/ui/Card.svelte";
    import Input from "$lib/components/ui/Input.svelte";

    let name = $state("My New Setup");
    let description = $state("");
    let isDefault = $state(false);
    let saving = $state(false);

    async function handleSave(setupData: string) {
        saving = true;
        try {
            await boardSetups.create({
                name,
                description,
                setup_data: setupData,
                is_default: isDefault,
            });
            goto("/board-setups");
        } catch (e: any) {
            toastStore.handleApiMessage(e, "Failed to save setup");
            saving = false;
        }
    }
</script>

<svelte:head>
    <title>GoStrategy — New Setup</title>
</svelte:head>

<div class="max-w-6xl mx-auto space-y-8">
    <div class="flex flex-col md:flex-row md:items-end justify-between gap-4">
        <div class="flex-1 space-y-1">
            <h1
                class="text-3xl font-black text-white uppercase tracking-tighter"
            >
                Create New Setup
            </h1>
            <p class="text-white/40">
                Design your starting positions for the battlefield.
            </p>
        </div>
    </div>

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
        <div class="text-center py-12 text-white/30 italic">
            Saving your setup...
        </div>
    {:else}
        <SetupEditor
            onSave={handleSave}
            onCancel={() => goto("/board-setups")}
        />
    {/if}
</div>
