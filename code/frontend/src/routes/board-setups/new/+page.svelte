<script lang="ts">
    import { goto } from "$app/navigation";
    import SetupEditor from "$lib/components/setup/SetupEditor.svelte";
    import { boardSetups } from "$lib/api/client";
    import Card from "$lib/components/ui/Card.svelte";

    let name = $state("My New Setup");
    let description = $state("");
    let isDefault = $state(false);
    let error = $state("");
    let saving = $state(false);

    async function handleSave(setupData: string) {
        saving = true;
        error = "";
        try {
            await boardSetups.create({
                name,
                description,
                setup_data: setupData,
                is_default: isDefault,
            });
            goto("/board-setups");
        } catch (e: any) {
            error = e.message || "Failed to save setup";
            saving = false;
        }
    }
</script>

<svelte:head>
    <title>Stratego — New Setup</title>
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
            <div class="space-y-2">
                <label
                    for="name"
                    class="text-xs font-bold text-white/40 uppercase tracking-widest pl-1"
                    >Setup Name</label
                >
                <input
                    id="name"
                    type="text"
                    bind:value={name}
                    placeholder="e.g., Aggressive Scout Rush"
                    class="w-full bg-white/5 border border-white/10 rounded-xl px-4 py-3 text-white focus:outline-none focus:border-brand-accent/50 transition-colors"
                />
            </div>
            <div class="space-y-2">
                <label
                    for="description"
                    class="text-xs font-bold text-white/40 uppercase tracking-widest pl-1"
                    >Description (Optional)</label
                >
                <input
                    id="description"
                    type="text"
                    bind:value={description}
                    placeholder="Briefly describe your strategy..."
                    class="w-full bg-white/5 border border-white/10 rounded-xl px-4 py-3 text-white focus:outline-none focus:border-brand-accent/50 transition-colors"
                />
            </div>
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

        {#if error}
            <div
                class="bg-brand-secondary/20 border border-brand-secondary/30 text-brand-secondary rounded-xl px-4 py-3 text-sm text-center"
            >
                {error}
            </div>
        {/if}
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
