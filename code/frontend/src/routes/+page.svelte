<script lang="ts">
    import Button from "$lib/components/ui/Button.svelte";
    import Card from "$lib/components/ui/Card.svelte";
    import Input from "$lib/components/ui/Input.svelte";

    let name = $state("");
</script>

<div class="space-y-12">
    <header>
        <h1 class="text-4xl font-bold text-white mb-2">Frontend Foundations</h1>
        <p class="text-white/40 text-lg">
            A premium, functional foundation for the new Stratego experience.
        </p>
    </header>

    <section class="grid grid-cols-1 md:grid-cols-2 gap-8">
        <Card class="space-y-6">
            <div>
                <h2 class="text-xl font-semibold text-white mb-1">
                    Interactive Elements
                </h2>
                <p class="text-white/40 text-sm">
                    Buttons with micro-animations and variant support.
                </p>
            </div>

            <div class="flex flex-wrap gap-4">
                <Button variant="primary">Primary Button</Button>
                <Button variant="secondary">Secondary</Button>
                <Button variant="outline">Outline</Button>
                <Button variant="ghost">Ghost</Button>
            </div>

            <div class="flex flex-wrap gap-4">
                <Button size="sm">Small</Button>
                <Button size="md">Medium</Button>
                <Button size="lg">Large Scale</Button>
            </div>
        </Card>

        <Card class="space-y-6">
            <div>
                <h2 class="text-xl font-semibold text-white mb-1">
                    Form Controls
                </h2>
                <p class="text-white/40 text-sm">
                    Standardized input fields with focus states.
                </p>
            </div>

            <div class="space-y-4">
                <Input
                    label="Your Name"
                    placeholder="Enter your name..."
                    bind:value={name}
                />
                {#if name}
                    <p
                        class="text-brand-primary text-sm animate-in fade-in slide-in-from-left-2 transition-all"
                    >
                        Hello, <span class="font-bold">{name}</span>! Welcome to
                        the new foundation.
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
            <h2 class="text-xl font-semibold text-white">
                Glassmorphism & Layout
            </h2>
            <p class="text-white/40">
                The entire interface uses a refined glassmorphism aesthetic with
                high-quality blur and subtle border highlights. This creates a
                depth-heavy, premium feel that modern users expect.
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
        </Card>Section
    </section>
</div>
