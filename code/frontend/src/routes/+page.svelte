<script lang="ts">
    import Button from "$lib/components/ui/Button.svelte";
    import Card from "$lib/components/ui/Card.svelte";
    import Input from "$lib/components/ui/Input.svelte";

    let name = $state("");
</script>

<div class="space-y-12">
    <header>
        <h1
            class="text-4xl font-extrabold text-white mb-2 uppercase tracking-wider text-brand-accent"
        >
            Stratego Command Center
        </h1>
        <p class="text-white/60 text-lg">
            A premium, adventurous foundation for the battlefield.
        </p>
    </header>

    <section class="grid grid-cols-1 md:grid-cols-2 gap-8">
        <Card class="space-y-6">
            <div>
                <h2
                    class="text-xl font-semibold text-brand-accent mb-1 uppercase tracking-wide"
                >
                    Tactical Controls
                </h2>
                <p class="text-white/60 text-sm">
                    Command buttons with micro-animations and variant support.
                </p>
            </div>

            <div class="flex flex-wrap gap-4">
                <Button variant="primary">Deploy Units</Button>
                <Button variant="secondary">Scout Ahead</Button>
                <Button variant="outline">Hold Position</Button>
                <Button variant="ghost">Retreat</Button>
            </div>

            <div class="flex flex-wrap gap-4">
                <Button size="sm">Small</Button>
                <Button size="md">Medium</Button>
                <Button size="lg">Large Scale</Button>
            </div>
        </Card>

        <Card class="space-y-6">
            <div>
                <h2
                    class="text-xl font-semibold text-brand-accent mb-1 uppercase tracking-wide"
                >
                    Intelligence Briefing
                </h2>
                <p class="text-white/60 text-sm">
                    Standardized input fields to gather enemy intel.
                </p>
            </div>

            <div class="space-y-4">
                <Input
                    label="Commander Name"
                    placeholder="Enter your rank & name..."
                    bind:value={name}
                />
                {#if name}
                    <p
                        class="text-brand-accent text-sm animate-in fade-in slide-in-from-left-2 transition-all border-l-2 border-brand-accent pl-2"
                    >
                        Intel received, <span class="font-bold">{name}</span>.
                        Awaiting deployment orders.
                    </p>
                {/if}
                <Input
                    label="Email Address"
                    type="email"
                    placeholder="name@example.com"
                />
            </div>
        </Card>

        <Card class="col-span-full space-y-4">
            <h2
                class="text-xl font-semibold text-brand-accent uppercase tracking-wide"
            >
                Thematic Atmosphere
            </h2>
            <p class="text-white/60">
                The interface utilizes deep crimsons, navy blues, and brass
                accents to evoke the classic Stratego board game, modernized
                with high-quality blur and subtle lighting effects.
            </p>
            <div class="grid grid-cols-1 sm:grid-cols-3 gap-4 pt-4">
                <div
                    class="h-24 rounded-xl bg-brand-primary/20 border border-brand-primary/30 flex items-center justify-center"
                >
                    <span class="text-brand-primary font-bold"
                        >Primary Accent</span
                    >
                </div>
                <div
                    class="h-24 rounded-xl bg-brand-secondary/20 border border-brand-secondary/30 flex items-center justify-center"
                >
                    <span class="text-brand-secondary font-bold"
                        >Secondary Accent</span
                    >
                </div>
                <div
                    class="h-24 rounded-xl bg-brand-accent/20 border border-brand-accent/30 flex items-center justify-center"
                >
                    <span class="text-brand-accent font-bold"
                        >Tertiary Accent</span
                    >
                </div>
            </div>
        </Card>
    </section>
</div>
